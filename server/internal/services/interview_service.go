package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zhitu/server/internal/config"
	"github.com/zhitu/server/internal/models"
	"github.com/zhitu/server/internal/utils"
	"gorm.io/gorm"
)

// 面试相关错误
var (
	ErrInterviewNotFound = errors.New("interview not found")
	ErrInterviewEnded    = errors.New("interview already ended")
	ErrMessageNotFound   = errors.New("interview message not found")
	ErrReportNotReady    = errors.New("report not generated yet, please end the interview first")
)

// 面试场景枚举
const (
	SceneTech     = "tech"
	SceneBehavior = "behavior"
	ScenePressure = "pressure"
	SceneHR       = "hr"
	SceneGroup    = "group"
	SceneTeaching = "teaching"
)

// 面试状态枚举
const (
	StatusOngoing   = "ongoing"
	StatusCompleted = "completed"
	StatusCancelled = "cancelled"
)

// 面试模式枚举
const (
	ModeText   = "text"
	ModeVoice  = "voice"
	ModeHybrid = "hybrid"
)

// InterviewService 面试业务逻辑
type InterviewService struct {
	db      *gorm.DB
	llm     *LLMService
	profile *ProfileService
	storage *config.StorageConfig
}

// NewInterviewService 构造 InterviewService
func NewInterviewService(db *gorm.DB, llm *LLMService, profile *ProfileService, storage *config.StorageConfig) *InterviewService {
	return &InterviewService{db: db, llm: llm, profile: profile, storage: storage}
}

// CreateInterviewInput 创建面试入参
type CreateInterviewInput struct {
	Scene          string `json:"scene" binding:"required"`
	TargetCompany  string `json:"target_company"`
	TargetPosition string `json:"target_position" binding:"required"`
	TargetJD       string `json:"target_jd"`
	Difficulty     string `json:"difficulty"`
	TotalQuestions int    `json:"total_questions"`
	Mode           string `json:"mode"`
}

// Create 创建面试会话，并自动生成第一道题
func (s *InterviewService) Create(ctx context.Context, userID uint, in *CreateInterviewInput) (*models.Interview, error) {
	if in.Scene != SceneTech && in.Scene != SceneBehavior && in.Scene != ScenePressure && in.Scene != SceneHR && in.Scene != SceneGroup && in.Scene != SceneTeaching {
		return nil, errors.New("invalid scene, must be one of: tech/behavior/pressure/hr/group/teaching")
	}
	if in.Difficulty == "" {
		in.Difficulty = "mid"
	}
	if in.TotalQuestions == 0 {
		in.TotalQuestions = 8
	}
	if in.TotalQuestions < 5 || in.TotalQuestions > 15 {
		return nil, errors.New("total_questions must be between 5 and 15")
	}
	if in.Mode == "" {
		in.Mode = ModeText
	}

	now := time.Now()
	interview := &models.Interview{
		UserID:         userID,
		Scene:          in.Scene,
		TargetCompany:  in.TargetCompany,
		TargetPosition: in.TargetPosition,
		TargetJD:       in.TargetJD,
		Difficulty:     in.Difficulty,
		TotalQuestions: in.TotalQuestions,
		Mode:           in.Mode,
		Status:         StatusOngoing,
		StartedAt:      &now,
	}
	if err := s.db.Create(interview).Error; err != nil {
		return nil, err
	}

	// 自动生成第一道题
	if err := s.askNextQuestion(ctx, interview, 1); err != nil {
		return nil, fmt.Errorf("generate first question: %w", err)
	}

	return interview, nil
}

// Get 获取面试详情（含所有消息）
func (s *InterviewService) Get(userID, id uint) (*models.Interview, error) {
	var interview models.Interview
	err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&interview).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInterviewNotFound
	}
	return &interview, err
}

// List 列出用户的所有面试
func (s *InterviewService) List(userID uint) ([]models.Interview, error) {
	var list []models.Interview
	err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&list).Error
	return list, err
}

// ListMessages 列出面试的所有消息
func (s *InterviewService) ListMessages(userID, interviewID uint) ([]models.InterviewMessage, error) {
	if _, err := s.Get(userID, interviewID); err != nil {
		return nil, err
	}
	var list []models.InterviewMessage
	err := s.db.Where("interview_id = ?", interviewID).Order("id ASC").Find(&list).Error
	return list, err
}

// SendMessage 用户发送文字回答，AI 生成下一题
// 通过 onDelta 回调流式推送 AI 回复
func (s *InterviewService) SendMessage(ctx context.Context, userID, interviewID uint, userText string, onDelta func(string)) (*models.InterviewMessage, error) {
	interview, err := s.Get(userID, interviewID)
	if err != nil {
		return nil, err
	}
	if interview.Status != StatusOngoing {
		return nil, ErrInterviewEnded
	}

	// 1. 存用户回答
	nextQuestionNo := interview.CurrentQuestionNo + 1
	userMsg := &models.InterviewMessage{
		InterviewID: interviewID,
		Role:        "user",
		Content:     userText,
		QuestionNo:  interview.CurrentQuestionNo,
	}
	if err := s.db.Create(userMsg).Error; err != nil {
		return nil, err
	}

	// 2. 检查是否已答完所有题
	if interview.CurrentQuestionNo >= interview.TotalQuestions {
		// 自动结束并生成报告
		return nil, s.endAndGenerateReport(ctx, interview)
	}

	// 3. AI 生成下一题（追问或新题）
	return s.askNextQuestionWithStream(ctx, interview, nextQuestionNo, onDelta)
}

// SendVoice 用户发送语音回答，先 Whisper 转写，再走文字逻辑
// 返回 AI 下一题的消息（含 TTS audio_url）
func (s *InterviewService) SendVoice(ctx context.Context, userID, interviewID uint, audio io.Reader, filename string, onDelta func(string)) (*models.InterviewMessage, error) {
	// 1. 保存音频文件到磁盘
	relPath, absPath, err := utils.SaveFile(audio, s.storage.AudioDir, filename)
	if err != nil {
		return nil, fmt.Errorf("save audio: %w", err)
	}

	// 2. Whisper 转写（重新读取磁盘文件）
	transcribed, err := s.transcribeFromPath(ctx, absPath)
	if err != nil {
		return nil, fmt.Errorf("transcribe: %w", err)
	}

	// 3. 存用户回答（含音频 URL）
	interview, err := s.Get(userID, interviewID)
	if err != nil {
		return nil, err
	}
	if interview.Status != StatusOngoing {
		return nil, ErrInterviewEnded
	}

	userMsg := &models.InterviewMessage{
		InterviewID: interviewID,
		Role:        "user",
		Content:     transcribed,
		AudioURL:    "/static/" + relPath,
		QuestionNo:  interview.CurrentQuestionNo,
	}
	if err := s.db.Create(userMsg).Error; err != nil {
		return nil, err
	}

	if interview.CurrentQuestionNo >= interview.TotalQuestions {
		return nil, s.endAndGenerateReport(ctx, interview)
	}

	// 4. AI 生成下一题
	return s.askNextQuestionWithStream(ctx, interview, interview.CurrentQuestionNo+1, onDelta)
}

// transcribeFromPath 从磁盘文件路径读取并转写
func (s *InterviewService) transcribeFromPath(ctx context.Context, absPath string) (string, error) {
	f, err := os.Open(absPath)
	if err != nil {
		return "", fmt.Errorf("open audio file: %w", err)
	}
	defer f.Close()
	return s.llm.Transcribe(ctx, f, filepath.Base(absPath))
}

// GetTTS 获取某条 AI 提问的 TTS 音频
// 若消息已有 audio_url 则直接返回路径，否则现场合成并保存
func (s *InterviewService) GetTTS(ctx context.Context, userID, interviewID, messageID uint) ([]byte, string, error) {
	if _, err := s.Get(userID, interviewID); err != nil {
		return nil, "", err
	}

	var msg models.InterviewMessage
	err := s.db.Where("id = ? AND interview_id = ? AND role = ?", messageID, interviewID, "assistant").First(&msg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", ErrMessageNotFound
	}
	if err != nil {
		return nil, "", err
	}

	// 若已有 audio_url，读取磁盘文件返回
	if msg.AudioURL != "" {
		absPath := strings.TrimPrefix(msg.AudioURL, "/static/")
		cwd, _ := os.Getwd()
		full := filepath.Join(cwd, absPath)
		data, err := os.ReadFile(full)
		if err == nil {
			return data, filepath.Base(full), nil
		}
		// 读取失败则重新合成
	}

	// 现场合成
	audio, err := s.llm.Synthesize(ctx, msg.Content)
	if err != nil {
		return nil, "", fmt.Errorf("synthesize: %w", err)
	}

	// 保存到磁盘
	filename := fmt.Sprintf("tts_%d.mp3", msg.ID)
	relPath, _, err := utils.SaveBytes(audio, s.storage.TTSDir, filename)
	if err != nil {
		return nil, "", fmt.Errorf("save tts: %w", err)
	}

	// 更新消息的 audio_url
	audioURL := "/static/" + relPath
	s.db.Model(&msg).Update("audio_url", audioURL)

	return audio, filename, nil
}

// End 结束面试并生成复盘报告
func (s *InterviewService) End(ctx context.Context, userID, id uint) (*models.InterviewReport, error) {
	interview, err := s.Get(userID, id)
	if err != nil {
		return nil, err
	}
	if interview.Status == StatusCompleted {
		// 已结束，直接返回报告
		return s.GetReport(userID, id)
	}

	if err := s.endAndGenerateReport(ctx, interview); err != nil {
		return nil, err
	}
	return s.GetReport(userID, id)
}

// endAndGenerateReport 内部结束面试并生成评分与报告
func (s *InterviewService) endAndGenerateReport(ctx context.Context, interview *models.Interview) error {
	now := time.Now()
	interview.Status = StatusCompleted
	interview.EndedAt = &now
	if err := s.db.Model(interview).Updates(map[string]interface{}{
		"status":   StatusCompleted,
		"ended_at": now,
	}).Error; err != nil {
		return err
	}

	// 生成评分
	if err := s.generateScores(ctx, interview); err != nil {
		// 评分失败不阻塞报告生成
		fmt.Printf("generate scores failed: %v\n", err)
	}

	// 生成报告
	if err := s.generateReport(ctx, interview); err != nil {
		return fmt.Errorf("generate report: %w", err)
	}
	return nil
}

// askNextQuestion 生成下一题（非流式，用于初始化第一题）
func (s *InterviewService) askNextQuestion(ctx context.Context, interview *models.Interview, questionNo int) error {
	_, err := s.askNextQuestionWithStream(ctx, interview, questionNo, nil)
	return err
}

// askNextQuestionWithStream 生成下一题，支持流式推送
func (s *InterviewService) askNextQuestionWithStream(ctx context.Context, interview *models.Interview, questionNo int, onDelta func(string)) (*models.InterviewMessage, error) {
	// 1. 读取历史消息
	var history []models.InterviewMessage
	s.db.Where("interview_id = ?", interview.ID).Order("id ASC").Find(&history)

	// 2. 读取用户档案摘要
	profileSummary := ""
	if fp, err := s.profile.GetFullProfile(interview.UserID); err == nil && fp.UserProfile != nil {
		profileSummary = fmt.Sprintf("姓名:%s, 目标岗位:%s, 教育:%d条, 工作:%d条, 项目:%d条",
			fp.UserProfile.RealName, fp.UserProfile.TargetPosition,
			len(fp.Educations), len(fp.Works), len(fp.Projects))
	}

	// 3. 构造 system prompt
	sysPrompt := s.buildInterviewerPrompt(interview, questionNo, profileSummary)

	// 4. 构造 messages
	messages := []ChatMessage{{Role: "system", Content: sysPrompt}}
	for _, m := range history {
		if m.Role == "assistant" {
			messages = append(messages, ChatMessage{Role: "assistant", Content: m.Content})
		} else {
			messages = append(messages, ChatMessage{Role: "user", Content: m.Content})
		}
	}
	// 提示 AI 出下一题
	messages = append(messages, ChatMessage{Role: "user", Content: fmt.Sprintf("（请出第 %d 题，只问一个问题）", questionNo)})

	// 5. 调 LLM
	var aiContent string
	var err error
	if onDelta != nil {
		// 流式：先收集完整内容，同时推送 delta
		var buf strings.Builder
		err = s.llm.ChatStream(ctx, messages, func(delta string) {
			buf.WriteString(delta)
			onDelta(delta)
		})
		aiContent = buf.String()
	} else {
		aiContent, err = s.llm.Chat(ctx, messages)
	}
	if err != nil {
		return nil, fmt.Errorf("llm ask: %w", err)
	}
	if strings.TrimSpace(aiContent) == "" {
		aiContent = "请简单介绍一下你自己。"
	}

	// 6. 推断问题类型
	questionType := inferQuestionType(aiContent, interview.Scene)

	// 7. 存 AI 消息
	aiMsg := &models.InterviewMessage{
		InterviewID:  interview.ID,
		Role:         "assistant",
		Content:      aiContent,
		QuestionType: questionType,
		QuestionNo:   questionNo,
	}
	if err := s.db.Create(aiMsg).Error; err != nil {
		return nil, err
	}

	// 8. 更新面试当前题号
	s.db.Model(interview).Update("current_question_no", questionNo)

	return aiMsg, nil
}

// buildInterviewerPrompt 构造面试官 system prompt
func (s *InterviewService) buildInterviewerPrompt(interview *models.Interview, questionNo int, profileSummary string) string {
	sceneDesc := map[string]string{
		SceneTech:     "技术面（算法、项目深挖、系统设计）",
		SceneBehavior: "行为面（STAR 法则）",
		ScenePressure: "压力面（挑战性/陷阱题）",
		SceneHR:       "HR 面（薪资/规划/离职原因）",
		SceneGroup:    "群面模拟（多角色讨论）",
		SceneTeaching: "教资模拟教室（结构化问答、模拟试讲、考官答辩）",
	}[interview.Scene]

	diffDesc := map[string]string{
		"junior": "初级",
		"mid":    "中级",
		"senior": "高级",
		"mixed":  "混合自适应",
	}[interview.Difficulty]

	return fmt.Sprintf(`你是一位资深面试官，正在面试一位应聘【%s】【%s】的候选人。

面试场景：%s
当前第 %d 题（共 %d 题），难度等级：%s

面试规则：
1. 一次只问一个问题
2. 根据候选人回答质量决定追问或换题——回答模糊或浅显时连续追问，回答充分时切换话题
3. 结合以下 JD 关键词出题
4. 难度递增
5. 模拟真实面试官语气，不要透露你是 AI
6. 如果是教资模拟教室：前两题进行结构化问答，中间两题要求候选人围绕抽题主题完成试讲片段，最后一题以考官身份针对教学设计进行答辩追问

企业风格提示：根据你对 %s 面试风格的了解调整出题策略（如字节跳动注重底层原理与项目深挖，阿里注重系统设计，腾讯注重技术广度等；若不确定则按通用标准）。

JD：
%s

候选人档案摘要：%s

请开始提问。`, interview.TargetCompany, interview.TargetPosition,
		sceneDesc, questionNo, interview.TotalQuestions, diffDesc,
		nonEmpty(interview.TargetCompany, "通用公司"),
		nonEmpty(interview.TargetJD, "（无 JD，按通用标准出题）"),
		nonEmpty(profileSummary, "（档案未填写）"))
}

// inferQuestionType 根据问题内容推断类型
func inferQuestionType(content string, scene string) string {
	c := strings.ToLower(content)
	if strings.Contains(c, "项目") || strings.Contains(c, "经历") {
		return "project"
	}
	if strings.Contains(c, "算法") || strings.Contains(c, "代码") || strings.Contains(c, "复杂度") {
		return "algorithm"
	}
	if strings.Contains(c, "设计") || strings.Contains(c, "架构") || strings.Contains(c, "系统") {
		return "system_design"
	}
	if strings.Contains(c, "为什么") || strings.Contains(c, "star") || strings.Contains(c, "冲突") {
		return "behavior"
	}
	if strings.Contains(c, "追问") || strings.Contains(c, "刚才") {
		return "followup"
	}
	return "open"
}

// generateScores 生成五维度评分
func (s *InterviewService) generateScores(ctx context.Context, interview *models.Interview) error {
	var messages []models.InterviewMessage
	s.db.Where("interview_id = ?", interview.ID).Order("id ASC").Find(&messages)

	// 拼接对话文本
	var dialog strings.Builder
	for _, m := range messages {
		role := "面试官"
		if m.Role == "user" {
			role = "候选人"
		}
		dialog.WriteString(fmt.Sprintf("%s：%s\n", role, m.Content))
	}

	prompt := fmt.Sprintf(`请对以下模拟面试对话进行五维度评分。

面试场景：%s
目标岗位：%s
JD：%s

对话记录：
%s

评分维度（每项 0-100）：
1. professional：专业能力（技术准确性、深度、知识覆盖面）
2. expression：表达能力（语言组织、术语准确、条理清晰）
3. logic：逻辑思维（因果链、结构化思考、问题拆解）
4. adaptability：应变能力（对追问的应对、思路调整速度）
5. pace：语速仪态（通过文字推断语速与停顿）

返回 JSON：
{
  "scores": [
    {"dimension":"professional","score":85,"comment":"..."},
    {"dimension":"expression","score":80,"comment":"..."},
    {"dimension":"logic","score":75,"comment":"..."},
    {"dimension":"adaptability","score":70,"comment":"..."},
    {"dimension":"pace","score":82,"comment":"..."}
  ]
}
只返回 JSON。`, interview.Scene, interview.TargetPosition, nonEmpty(interview.TargetJD, "无"), dialog.String())

	messagesLLM := []ChatMessage{
		{Role: "system", Content: "你是面试评分专家，严格输出 JSON。"},
		{Role: "user", Content: prompt},
	}

	var result struct {
		Scores []struct {
			Dimension string `json:"dimension"`
			Score     int    `json:"score"`
			Comment   string `json:"comment"`
		} `json:"scores"`
	}
	if err := s.llm.ChatJSON(ctx, messagesLLM, &result); err != nil {
		return err
	}

	// 写入评分表
	for _, sc := range result.Scores {
		score := &models.InterviewScore{
			InterviewID: interview.ID,
			Dimension:   sc.Dimension,
			Score:       sc.Score,
			Comment:     sc.Comment,
		}
		s.db.Create(score)
	}
	return nil
}

// generateReport 生成复盘报告
func (s *InterviewService) generateReport(ctx context.Context, interview *models.Interview) error {
	var messages []models.InterviewMessage
	s.db.Where("interview_id = ?", interview.ID).Order("id ASC").Find(&messages)

	var scores []models.InterviewScore
	s.db.Where("interview_id = ?", interview.ID).Find(&scores)

	// 计算加权总分
	overallScore := calcOverallScore(scores, interview.Scene)

	// 拼接对话
	var dialog strings.Builder
	for _, m := range messages {
		role := "面试官"
		if m.Role == "user" {
			role = "候选人"
		}
		dialog.WriteString(fmt.Sprintf("%s：%s\n", role, m.Content))
	}

	scoresJSON, _ := json.Marshal(scores)
	prompt := fmt.Sprintf(`请为以下模拟面试生成复盘报告。

面试场景：%s
目标岗位：%s
五维度评分 JSON：%s

对话记录：
%s

返回 JSON：
{
  "summary": "总体评价文本",
  "highlights": ["亮点1","亮点2"],
  "improvements": ["改进建议1","改进建议2"],
  "recommendations": ["推荐练习方向1","推荐练习方向2"],
  "word_cloud": [{"word":"高频词","count":5}]
}
只返回 JSON。`, interview.Scene, interview.TargetPosition, string(scoresJSON), dialog.String())

	messagesLLM := []ChatMessage{
		{Role: "system", Content: "你是面试复盘专家，严格输出 JSON。"},
		{Role: "user", Content: prompt},
	}

	var reportData struct {
		Summary         string   `json:"summary"`
		Highlights      []string `json:"highlights"`
		Improvements    []string `json:"improvements"`
		Recommendations []string `json:"recommendations"`
		WordCloud       []struct {
			Word  string `json:"word"`
			Count int    `json:"count"`
		} `json:"word_cloud"`
	}
	if err := s.llm.ChatJSON(ctx, messagesLLM, &reportData); err != nil {
		return err
	}

	highlights, _ := json.Marshal(reportData.Highlights)
	improvements, _ := json.Marshal(reportData.Improvements)
	recommendations, _ := json.Marshal(reportData.Recommendations)
	wordCloud, _ := json.Marshal(reportData.WordCloud)

	report := &models.InterviewReport{
		InterviewID:     interview.ID,
		OverallScore:    overallScore,
		Summary:         reportData.Summary,
		Highlights:      string(highlights),
		Improvements:    string(improvements),
		Recommendations: string(recommendations),
		WordCloud:       string(wordCloud),
	}
	return s.db.Create(report).Error
}

// calcOverallScore 根据场景权重计算加权总分
func calcOverallScore(scores []models.InterviewScore, scene string) int {
	weights := map[string]map[string]float64{
		SceneTech:     {"professional": 0.4, "expression": 0.15, "logic": 0.25, "adaptability": 0.1, "pace": 0.1},
		SceneBehavior: {"professional": 0.2, "expression": 0.3, "logic": 0.15, "adaptability": 0.25, "pace": 0.1},
		ScenePressure: {"professional": 0.2, "expression": 0.15, "logic": 0.25, "adaptability": 0.35, "pace": 0.05},
		SceneHR:       {"professional": 0.1, "expression": 0.3, "logic": 0.15, "adaptability": 0.2, "pace": 0.25},
		SceneGroup:    {"professional": 0.2, "expression": 0.25, "logic": 0.15, "adaptability": 0.3, "pace": 0.1},
		SceneTeaching: {"professional": 0.3, "expression": 0.25, "logic": 0.15, "adaptability": 0.2, "pace": 0.1},
	}[scene]
	if weights == nil {
		weights = map[string]float64{"professional": 0.25, "expression": 0.2, "logic": 0.2, "adaptability": 0.2, "pace": 0.15}
	}

	scoreMap := map[string]int{}
	for _, s := range scores {
		scoreMap[s.Dimension] = s.Score
	}

	total := 0.0
	for dim, w := range weights {
		if sc, ok := scoreMap[dim]; ok {
			total += float64(sc) * w
		}
	}
	return int(total)
}

// GetReport 获取复盘报告
func (s *InterviewService) GetReport(userID, interviewID uint) (*models.InterviewReport, error) {
	if _, err := s.Get(userID, interviewID); err != nil {
		return nil, err
	}
	var report models.InterviewReport
	err := s.db.Where("interview_id = ?", interviewID).First(&report).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrReportNotReady
	}
	return &report, err
}

// GetScores 获取评分明细
func (s *InterviewService) GetScores(userID, interviewID uint) ([]models.InterviewScore, error) {
	if _, err := s.Get(userID, interviewID); err != nil {
		return nil, err
	}
	var list []models.InterviewScore
	err := s.db.Where("interview_id = ?", interviewID).Order("id ASC").Find(&list).Error
	return list, err
}
