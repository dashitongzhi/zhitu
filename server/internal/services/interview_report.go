package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhitu/server/internal/models"
	"gorm.io/gorm"
	"strings"
)

// generateScores 生成五维度评分
// 模型：cfg.ChatModel → "mimo-v2.5-pro"（文字分析，结构化评分）
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

候选人简历摘要：
%s

对话记录：
%s

评分维度（每项 0-100）：
1. professional：专业能力（技术准确性、深度、知识覆盖面）
2. expression：表达能力（语言组织、术语准确、条理清晰）
3. logic：逻辑思维（因果链、结构化思考、问题拆解）
4. adaptability：应变能力（对追问的应对、思路调整速度）
5. pace：语速仪态（通过文字推断语速与停顿）

如果候选人发送了简历，请在评分时参考简历内容，判断其回答是否与简历经历一致、是否充分展开。
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
只返回 JSON。`, interview.Scene, interview.TargetPosition,
		nonEmpty(interview.TargetJD, "无"),
		nonEmpty(summarizeResume(interview.ResumeSnapshot), "（候选人未发送简历）"),
		dialog.String())

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
// 模型：cfg.ChatModel → "mimo-v2.5-pro"（文字分析，复盘报告）
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

候选人简历摘要：
%s

对话记录：
%s

如果候选人发送了简历，请在复盘报告中结合简历内容评估其经历表达与岗位匹配度。

请逐题回顾候选人的回答表现，并把详细评价集中写入 question_feedback：
- 只评价已经收到候选人回答的题目，按真实对话顺序逐题匹配，不得虚构题目或回答。
- question_feedback 条目数量必须与候选人已回答题数一致，question_no 使用原题号。
- question 保留面试官原题；answer 准确概括候选人实际回答，不得替候选人补充不存在的经历。
- score 是该题独立百分制评分（0-100）。
- comment 具体指出内容完整性、逻辑、证据、表达上的亮点与不足，避免“回答不错”等空泛结论。
- suggestion 给出可执行的更优回答结构、应补数据或关键要点。
- 总体 summary、highlights、improvements 必须与逐题评价保持一致。

返回 JSON：
{
  "summary": "总体评价文本",
  "highlights": ["亮点1","亮点2"],
  "improvements": ["改进建议1","改进建议2"],
  "recommendations": ["推荐练习方向1","推荐练习方向2"],
  "word_cloud": [{"word":"高频词","count":5}],
  "question_feedback": [
    {
      "question_no": 1,
      "question": "面试官原题文本",
      "answer": "候选人回答的摘要（可截断）",
      "score": 82,
      "comment": "对这一题的点评，指出亮点与不足",
      "suggestion": "针对该题的更优回答方向或要点"
    }
  ]
}
只返回 JSON。`, interview.Scene, interview.TargetPosition, string(scoresJSON),
		nonEmpty(summarizeResume(interview.ResumeSnapshot), "（候选人未发送简历）"),
		dialog.String())

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
		QuestionFeedback []struct {
			QuestionNo int    `json:"question_no"`
			Question   string `json:"question"`
			Answer     string `json:"answer"`
			Score      int    `json:"score"`
			Comment    string `json:"comment"`
			Suggestion string `json:"suggestion"`
		} `json:"question_feedback"`
	}
	if err := s.llm.ChatJSON(ctx, messagesLLM, &reportData); err != nil {
		return err
	}

	highlights, _ := json.Marshal(reportData.Highlights)
	improvements, _ := json.Marshal(reportData.Improvements)
	recommendations, _ := json.Marshal(reportData.Recommendations)
	wordCloud, _ := json.Marshal(reportData.WordCloud)
	questionFeedback, _ := json.Marshal(reportData.QuestionFeedback)

	report := &models.InterviewReport{
		InterviewID:      interview.ID,
		OverallScore:     overallScore,
		Summary:          reportData.Summary,
		Highlights:       string(highlights),
		Improvements:     string(improvements),
		Recommendations:  string(recommendations),
		WordCloud:        string(wordCloud),
		QuestionFeedback: string(questionFeedback),
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
