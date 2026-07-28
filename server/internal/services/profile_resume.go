package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zhitu/server/internal/models"
	"io"
	"strings"
)

// ---------- 简历上传解析 ----------

// ParsedProfile LLM 解析后的结构化档案
type ParsedProfile struct {
	RealName         string                 `json:"real_name"`
	Gender           string                 `json:"gender"`
	Phone            string                 `json:"phone"`
	TargetPosition   string                 `json:"target_position"`
	SelfIntroduction string                 `json:"self_introduction"`
	Educations       []models.UserEducation `json:"educations"`
	Works            []models.UserWork      `json:"works"`
	Projects         []models.UserProject   `json:"projects"`
	Skills           []models.UserSkill     `json:"skills"`
	Honors           []models.UserHonor     `json:"honors"`
	Practices        []models.UserPractice  `json:"practices"`
}

// ParseResumeFile 解析上传的简历文件，提取结构化档案并写入数据库
// reader 为文件内容，filename 用于判断格式
func (s *ProfileService) ParseResumeFile(ctx context.Context, userID uint, reader io.Reader, filename string) (*ParsedProfile, error) {
	// 1. 提取纯文本
	text, err := extractResumeText(reader, filename)
	if err != nil {
		return nil, fmt.Errorf("extract resume text: %w", err)
	}
	if strings.TrimSpace(text) == "" {
		return nil, ErrResumeFileEmpty
	}

	// 2. LLM 结构化
	parsed, err := s.parseResumeByLLM(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrResumeParseFail, err)
	}

	// 3. 写入档案（合并策略：主表 upsert，子表追加）
	if err := s.mergeParsedProfile(userID, parsed); err != nil {
		return nil, fmt.Errorf("merge parsed profile: %w", err)
	}

	return parsed, nil
}

// parseResumeByLLM 调 LLM 将简历文本转为结构化 JSON
func (s *ProfileService) parseResumeByLLM(ctx context.Context, text string) (*ParsedProfile, error) {
	prompt := `你是简历解析助手。请从下面的简历文本中提取结构化信息，严格返回 JSON 格式。
字段说明：
- real_name: 姓名
- gender: 性别（male/female/other，无法判断留空）
- phone: 电话
- target_position: 目标岗位（若简历中未明确，根据经历推断）
- self_introduction: 个人介绍/自我评价
- educations: 教育背景数组，每项含 school/major/degree(大专/本科/硕士/博士)/start_date(YYYY-MM)/end_date(YYYY-MM或"至今")/gpa/courses
- works: 工作经历数组，每项含 company/position/start_date/end_date/description(STAR法则)/leave_reason
- projects: 项目经历数组，每项含 name/role/start_date/end_date/description/tech_stack(填字符串，逗号分隔)
- skills: 技能数组，每项含 category(编程语言/框架工具/软技能/证书/语言能力)/name/proficiency(了解/熟悉/熟练/精通)
- honors: 荣誉奖项数组，每项含 name/issuer/award_date(YYYY-MM)/level(校级/市级/省级/国家级/国际级)
- practices: 校内外实践数组，每项含 title/organization/start_date/end_date/description

注意：
1. 只返回 JSON，不要任何解释或 markdown 代码块
2. 无法提取的字段留空字符串或空数组
3. 日期统一 YYYY-MM 格式，无法判断具体月份用 "-01"
4. tech_stack 字段填逗号分隔字符串（如 "Go,React,MySQL"），后端会转 JSON 存储

简历文本：
` + text

	messages := []ChatMessage{
		{Role: "system", Content: "你是一个严格遵循指令的简历解析助手，只输出 JSON。"},
		{Role: "user", Content: prompt},
	}

	var parsed ParsedProfile
	if err := s.llm.ChatJSON(ctx, messages, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

// mergeParsedProfile 将解析结果合并写入档案
// 主表 upsert，子表直接追加（不删原有数据）
func (s *ProfileService) mergeParsedProfile(userID uint, p *ParsedProfile) error {
	// 主表
	updates := map[string]interface{}{}
	if p.RealName != "" {
		updates["real_name"] = p.RealName
	}
	if p.Gender != "" {
		updates["gender"] = p.Gender
	}
	if p.Phone != "" {
		updates["phone"] = p.Phone
	}
	if p.TargetPosition != "" {
		updates["target_position"] = p.TargetPosition
	}
	if p.SelfIntroduction != "" {
		updates["self_introduction"] = p.SelfIntroduction
	}
	if len(updates) > 0 {
		if _, err := s.UpsertProfile(userID, updates); err != nil {
			return err
		}
	} else {
		// 确保主表存在以记录完成度
		s.db.Where("user_id = ?", userID).FirstOrCreate(&models.UserProfile{UserID: userID})
	}

	// 子表批量追加
	for i := range p.Educations {
		p.Educations[i].ID = 0
		p.Educations[i].UserID = userID
	}
	if len(p.Educations) > 0 {
		if err := s.db.Create(&p.Educations).Error; err != nil {
			return err
		}
	}
	for i := range p.Works {
		p.Works[i].ID = 0
		p.Works[i].UserID = userID
	}
	if len(p.Works) > 0 {
		if err := s.db.Create(&p.Works).Error; err != nil {
			return err
		}
	}
	for i := range p.Projects {
		p.Projects[i].ID = 0
		p.Projects[i].UserID = userID
	}
	if len(p.Projects) > 0 {
		if err := s.db.Create(&p.Projects).Error; err != nil {
			return err
		}
	}
	for i := range p.Skills {
		p.Skills[i].ID = 0
		p.Skills[i].UserID = userID
	}
	if len(p.Skills) > 0 {
		if err := s.db.Create(&p.Skills).Error; err != nil {
			return err
		}
	}
	for i := range p.Honors {
		p.Honors[i].ID = 0
		p.Honors[i].UserID = userID
	}
	if len(p.Honors) > 0 {
		if err := s.db.Create(&p.Honors).Error; err != nil {
			return err
		}
	}
	for i := range p.Practices {
		p.Practices[i].ID = 0
		p.Practices[i].UserID = userID
	}
	if len(p.Practices) > 0 {
		if err := s.db.Create(&p.Practices).Error; err != nil {
			return err
		}
	}

	s.touchCompletion(userID)
	return nil
}

// MarshalTechStack 将技能标签数组序列化为 JSON 字符串存储
func MarshalTechStack(tags []string) string {
	if len(tags) == 0 {
		return "[]"
	}
	b, _ := json.Marshal(tags)
	return string(b)
}

// UnmarshalTechStack 反序列化 tech_stack 字段
func UnmarshalTechStack(s string) []string {
	if s == "" {
		return []string{}
	}
	var tags []string
	if err := json.Unmarshal([]byte(s), &tags); err != nil {
		return []string{}
	}
	return tags
}
