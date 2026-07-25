package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/zhitu/server/internal/models"
	"gorm.io/gorm"
)

// ResumeAIService 简历 AI 业务逻辑
type ResumeAIService struct {
	db       *gorm.DB
	llm      *LLMService
	resume   *ResumeService
	profile  *ProfileService
}

// NewResumeAIService 构造 ResumeAIService
func NewResumeAIService(db *gorm.DB, llm *LLMService, resume *ResumeService, profile *ProfileService) *ResumeAIService {
	return &ResumeAIService{db: db, llm: llm, resume: resume, profile: profile}
}

// GenerateInput AI 生成简历入参
type GenerateInput struct {
	TargetJD    string `json:"target_jd"`     // 目标 JD（可空，空则纯档案生成）
	Scene       string `json:"scene"`         // 生成场景描述（可空）
	ModuleHints string `json:"module_hints"`  // 模块生成提示（可空，如"重点突出项目经历"）
}

// Generate AI 参考用户档案 + JD/场景生成完整简历内容
// 通过 onDelta 回调流式推送生成过程文本（用于前端展示进度）
func (s *ResumeAIService) Generate(ctx context.Context, userID, resumeID uint, in *GenerateInput, onDelta func(string)) (*models.ResumeVersion, error) {
	// 1. 读取用户档案
	fp, err := s.profile.GetFullProfile(userID)
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err)
	}

	// 2. 读取简历元信息
	resume, err := s.resume.Get(userID, resumeID)
	if err != nil {
		return nil, err
	}

	// 3. 构造 prompt
	prompt := s.buildGeneratePrompt(fp, resume, in)
	messages := []ChatMessage{
		{Role: "system", Content: "你是资深简历顾问，擅长根据候选人档案与目标岗位 JD 生成高质量简历。严格输出 JSON。"},
		{Role: "user", Content: prompt},
	}

	// 4. 若有 onDelta 则先流式推送"生成中"提示（实际生成走 ChatJSON 保证结构）
	if onDelta != nil {
		onDelta("正在分析您的档案与目标岗位...\n")
	}

	// 5. 调 LLM 生成结构化简历
	var content ResumeContent
	if err := s.llm.ChatJSON(ctx, messages, &content); err != nil {
		return nil, fmt.Errorf("llm generate: %w", err)
	}

	// 确保 module_order 与 module_visibility 存在
	if len(content.ModuleOrder) == 0 {
		content.ModuleOrder = []string{"personal", "intention", "education", "work", "project", "skills", "honor"}
	}
	if content.ModuleVisibility == nil {
		content.ModuleVisibility = map[string]bool{}
		for _, m := range content.ModuleOrder {
			content.ModuleVisibility[m] = true
		}
	}

	contentBytes, err := json.Marshal(content)
	if err != nil {
		return nil, fmt.Errorf("marshal content: %w", err)
	}

	// 6. 创建新版本
	version, err := s.resume.CreateVersion(userID, resumeID, &CreateVersionInput{
		Content:    string(contentBytes),
		ChangeNote: "AI 生成",
	})
	if err != nil {
		return nil, err
	}

	// 7. 记录 AI 操作
	beforeID := uint(0)
	if resume.CurrentVersionID != 0 {
		beforeID = resume.CurrentVersionID
	}
	s.recordAIOperation(resumeID, beforeID, version.ID, "generate", in.TargetJD, contentBytes)

	if onDelta != nil {
		onDelta(fmt.Sprintf("\n生成完成，版本 %s 已保存。", version.VersionLabel))
	}

	return version, nil
}

// buildGeneratePrompt 构造 AI 生成简历的 prompt
func (s *ResumeAIService) buildGeneratePrompt(fp *FullProfile, resume *models.Resume, in *GenerateInput) string {
	profileJSON, _ := json.Marshal(struct {
		Profile     *models.UserProfile   `json:"profile"`
		Educations  []models.UserEducation `json:"educations"`
		Works       []models.UserWork      `json:"works"`
		Projects    []models.UserProject   `json:"projects"`
		Skills      []models.UserSkill     `json:"skills"`
		Honors      []models.UserHonor     `json:"honors"`
		Practices   []models.UserPractice  `json:"practices"`
	}{
		Profile:    fp.UserProfile,
		Educations: fp.Educations,
		Works:      fp.Works,
		Projects:   fp.Projects,
		Skills:     fp.Skills,
		Honors:     fp.Honors,
		Practices:  fp.Practices,
	})

	prompt := fmt.Sprintf(`请根据以下用户档案生成一份完整简历。

简历名称：%s
目标公司：%s
目标岗位：%s
目标 JD：
%s

场景描述：%s
模块生成提示：%s

用户档案 JSON：
%s

输出要求：
1. 严格返回 JSON，字段结构如下：
{
  "personal": {"name":"","gender":"","age":"","phone":"","email":"","github":"","avatar":"","city":""},
  "intention": {"position":"","city":"","salary":"","arrival":"","industry":""},
  "education": [{"school":"","major":"","degree":"","start":"","end":"","courses":"","gpa":""}],
  "work": [{"company":"","position":"","start":"","end":"","description":"","leave_reason":""}],
  "project": [{"name":"","role":"","start":"","end":"","description":"","tech_stack":[],"url":""}],
  "skills": [{"category":"","name":"","proficiency":""}],
  "honor": [{"name":"","issuer":"","date":"","level":""}],
  "custom": [{"title":"","content":""}],
  "module_order": ["personal","intention","education","work","project","skills","honor"],
  "module_visibility": {"personal":true,"work":true}
}
2. 个人信息从档案主表填充，缺失字段留空
3. 教育/工作/项目/技能/荣誉从档案子表填充，并用 STAR 法则润色描述
4. 若提供了 JD，针对性突出与 JD 匹配的经历与技能
5. description 字段用换行分隔多条要点
6. 日期统一 YYYY-MM 格式，"至今"用字符串"至今"
7. 只返回 JSON，不要任何解释`,
		resume.Name, resume.TargetCompany, resume.TargetPosition,
		nonEmpty(in.TargetJD, "（无）"),
		nonEmpty(in.Scene, "（无）"),
		nonEmpty(in.ModuleHints, "（无）"),
		string(profileJSON),
	)
	return prompt
}

// PolishInput AI 润色入参
type PolishInput struct {
	Module string `json:"module"` // 要润色的模块名：work/project/all
	JD     string `json:"jd"`     // 润色参考的 JD（可空）
}

// Polish AI 润色简历指定模块（或全部）
func (s *ResumeAIService) Polish(ctx context.Context, userID, resumeID uint, in *PolishInput) (*models.ResumeVersion, error) {
	resume, err := s.resume.Get(userID, resumeID)
	if err != nil {
		return nil, err
	}

	// 取当前版本内容
	version, err := s.resume.GetVersion(userID, resumeID, resume.CurrentVersionID)
	if err != nil {
		return nil, err
	}

	var content ResumeContent
	if err := json.Unmarshal([]byte(version.Content), &content); err != nil {
		return nil, fmt.Errorf("parse content: %w", err)
	}

	// 构造润色 prompt
	moduleDesc := in.Module
	if moduleDesc == "" || moduleDesc == "all" {
		moduleDesc = "全部模块"
	}

	contentJSON, _ := json.Marshal(content)
	prompt := fmt.Sprintf(`请润色以下简历的【%s】部分。

参考 JD：
%s

当前简历 JSON：
%s

润色要求：
1. 用 STAR 法则优化工作/项目描述，动词开头
2. 量化成果（如"提升 30%%"），若原文无数据可合理推断但不要编造
3. 保持原有 JSON 结构，只修改内容不改变字段
4. 返回完整的简历 JSON（包含所有模块，未润色的模块保持原样）
5. 只返回 JSON，不要解释`, moduleDesc, nonEmpty(in.JD, "（无）"), string(contentJSON))

	messages := []ChatMessage{
		{Role: "system", Content: "你是简历润色专家，严格输出 JSON。"},
		{Role: "user", Content: prompt},
	}

	var newContent ResumeContent
	if err := s.llm.ChatJSON(ctx, messages, &newContent); err != nil {
		return nil, fmt.Errorf("llm polish: %w", err)
	}

	newBytes, _ := json.Marshal(newContent)
	newVersion, err := s.resume.CreateVersion(userID, resumeID, &CreateVersionInput{
		Content:    string(newBytes),
		ChangeNote: fmt.Sprintf("AI 润色：%s", moduleDesc),
	})
	if err != nil {
		return nil, err
	}

	s.recordAIOperation(resumeID, resume.CurrentVersionID, newVersion.ID, "polish", in.JD, newBytes)
	return newVersion, nil
}

// ScoreResult AI 评分结果
type ScoreResult struct {
	OverallScore int             `json:"overall_score"` // 0-100
	Dimensions   []ScoreDimension `json:"dimensions"`
	Suggestions  []string         `json:"suggestions"`
}

// ScoreDimension 评分维度
type ScoreDimension struct {
	Name    string `json:"name"`
	Score   int    `json:"score"`   // 0-100
	Comment string `json:"comment"`
}

// Score AI 评分简历（基于 JD 匹配 + 内容质量）
func (s *ResumeAIService) Score(ctx context.Context, userID, resumeID uint, jd string) (*ScoreResult, error) {
	resume, err := s.resume.Get(userID, resumeID)
	if err != nil {
		return nil, err
	}
	version, err := s.resume.GetVersion(userID, resumeID, resume.CurrentVersionID)
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(`请对以下简历进行评分。

目标 JD：
%s

简历 JSON：
%s

评分维度（每项 0-100）：
1. keyword_coverage：关键词覆盖（与 JD 匹配度）
2. project_depth：项目深度
3. experience_relevance：经验相关性
4. skill_breadth：技能广度
5. expression_quality：表达质量

返回 JSON：
{
  "overall_score": 85,
  "dimensions": [{"name":"keyword_coverage","score":80,"comment":"..."}],
  "suggestions": ["建议1","建议2"]
}
只返回 JSON。`, nonEmpty(jd, "（无 JD，按通用标准评分）"), version.Content)

	messages := []ChatMessage{
		{Role: "system", Content: "你是简历评分专家，严格输出 JSON。"},
		{Role: "user", Content: prompt},
	}

	var result ScoreResult
	if err := s.llm.ChatJSON(ctx, messages, &result); err != nil {
		return nil, fmt.Errorf("llm score: %w", err)
	}

	// 记录操作
	resultBytes, _ := json.Marshal(result)
	s.recordAIOperation(resumeID, resume.CurrentVersionID, 0, "score", jd, resultBytes)

	return &result, nil
}

// MatchResult AI JD 匹配结果
type MatchResult struct {
	MatchScore      int              `json:"match_score"` // 0-100
	MissingKeywords []string         `json:"missing_keywords"`
	Suggestions     []MatchSuggestion `json:"suggestions"`
}

// MatchSuggestion JD 匹配优化建议
type MatchSuggestion struct {
	Type     string `json:"type"`     // 补充/润色/优化/调整
	Module   string `json:"module"`   // 目标模块
	Content  string `json:"content"`  // 具体修改内容
	Reason   string `json:"reason"`
}

// JDMatch AI 分析简历与 JD 匹配度
func (s *ResumeAIService) JDMatch(ctx context.Context, userID, resumeID uint, jd string) (*MatchResult, error) {
	if jd == "" {
		return nil, errors.New("jd is required for match analysis")
	}

	resume, err := s.resume.Get(userID, resumeID)
	if err != nil {
		return nil, err
	}
	version, err := s.resume.GetVersion(userID, resumeID, resume.CurrentVersionID)
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(`请分析简历与目标 JD 的匹配度。

JD：
%s

简历 JSON：
%s

返回 JSON：
{
  "match_score": 75,
  "missing_keywords": ["React","TypeScript"],
  "suggestions": [
    {"type":"补充","module":"skills","content":"增加 React 技能","reason":"JD 要求 React 但简历未体现"},
    {"type":"润色","module":"work","content":"在工作经历中强调前端项目","reason":"提升相关性"}
  ]
}
只返回 JSON。`, jd, version.Content)

	messages := []ChatMessage{
		{Role: "system", Content: "你是 JD 匹配分析专家，严格输出 JSON。"},
		{Role: "user", Content: prompt},
	}

	var result MatchResult
	if err := s.llm.ChatJSON(ctx, messages, &result); err != nil {
		return nil, fmt.Errorf("llm jd_match: %w", err)
	}

	resultBytes, _ := json.Marshal(result)
	s.recordAIOperation(resumeID, resume.CurrentVersionID, 0, "jd_match", jd, resultBytes)

	return &result, nil
}

// recordAIOperation 记录 AI 操作
func (s *ResumeAIService) recordAIOperation(resumeID, beforeID, afterID uint, opType, jd string, result []byte) {
	op := &models.ResumeAIOperation{
		ResumeID:         resumeID,
		VersionIDBefore:  beforeID,
		VersionIDAfter:   afterID,
		OperationType:    opType,
		InputJD:          jd,
		Result:           string(result),
	}
	s.db.Create(op)
}

// nonEmpty 返回非空字符串，空则返回 fallback
func nonEmpty(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}
