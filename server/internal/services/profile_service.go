package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/zhitu/server/internal/config"
	"github.com/zhitu/server/internal/models"
	"gorm.io/gorm"
)

// 档案相关错误
var (
	ErrProfileNotFound  = errors.New("profile not found, please create first")
	ErrSubResourceNotFound = errors.New("sub resource not found")
	ErrResumeFileEmpty  = errors.New("resume file content is empty or unparseable")
	ErrResumeParseFail  = errors.New("failed to parse resume via llm")
)

// ProfileService 用户档案业务逻辑
type ProfileService struct {
	db      *gorm.DB
	llm     *LLMService
	storage *config.StorageConfig
}

// NewProfileService 构造 ProfileService
func NewProfileService(db *gorm.DB, llm *LLMService, storage *config.StorageConfig) *ProfileService {
	return &ProfileService{db: db, llm: llm, storage: storage}
}

// FullProfile 完整档案（主表 + 所有子资源），用于一次性返回前端
type FullProfile struct {
	UserProfile     *models.UserProfile   `json:"profile"`
	Educations      []models.UserEducation `json:"educations"`
	Works           []models.UserWork      `json:"works"`
	Projects        []models.UserProject   `json:"projects"`
	Skills          []models.UserSkill     `json:"skills"`
	Honors          []models.UserHonor     `json:"honors"`
	Practices       []models.UserPractice  `json:"practices"`
}

// ---------- 主表 ----------

// GetFullProfile 获取完整档案（主表 + 所有子资源）
// 若主表不存在则返回空壳（不自动创建，由用户显式 PUT 触发）
func (s *ProfileService) GetFullProfile(userID uint) (*FullProfile, error) {
	fp := &FullProfile{
		Educations: []models.UserEducation{},
		Works:      []models.UserWork{},
		Projects:   []models.UserProject{},
		Skills:     []models.UserSkill{},
		Honors:     []models.UserHonor{},
		Practices:  []models.UserPractice{},
	}

	var profile models.UserProfile
	err := s.db.First(&profile, "user_id = ?", userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fp.UserProfile = &models.UserProfile{UserID: userID}
		return fp, nil
	}
	if err != nil {
		return nil, err
	}
	fp.UserProfile = &profile

	if err := s.db.Where("user_id = ?", userID).Order("id ASC").Find(&fp.Educations).Error; err != nil {
		return nil, err
	}
	if err := s.db.Where("user_id = ?", userID).Order("id ASC").Find(&fp.Works).Error; err != nil {
		return nil, err
	}
	if err := s.db.Where("user_id = ?", userID).Order("id ASC").Find(&fp.Projects).Error; err != nil {
		return nil, err
	}
	if err := s.db.Where("user_id = ?", userID).Order("id ASC").Find(&fp.Skills).Error; err != nil {
		return nil, err
	}
	if err := s.db.Where("user_id = ?", userID).Order("id ASC").Find(&fp.Honors).Error; err != nil {
		return nil, err
	}
	if err := s.db.Where("user_id = ?", userID).Order("id ASC").Find(&fp.Practices).Error; err != nil {
		return nil, err
	}
	return fp, nil
}

// UpsertProfile 创建或更新档案基础信息，并重算完成度
func (s *ProfileService) UpsertProfile(userID uint, updates map[string]interface{}) (*models.UserProfile, error) {
	var profile models.UserProfile
	err := s.db.First(&profile, "user_id = ?", userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建
		profile = models.UserProfile{UserID: userID}
		for k, v := range updates {
			applyProfileField(&profile, k, v)
		}
		if err := s.db.Create(&profile).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		// 更新
		for k, v := range updates {
			applyProfileField(&profile, k, v)
		}
		if err := s.db.Save(&profile).Error; err != nil {
			return nil, err
		}
	}

	// 重算完成度
	pct := s.calcCompletion(userID)
	profile.CompletionPct = pct
	s.db.Model(&profile).Update("completion_pct", pct)
	return &profile, nil
}

// applyProfileField 将 map 中的字段应用到 profile 结构体
func applyProfileField(p *models.UserProfile, key string, val interface{}) {
	switch key {
	case "real_name":
		p.RealName = toString(val)
	case "gender":
		p.Gender = toString(val)
	case "phone":
		p.Phone = toString(val)
	case "target_position":
		p.TargetPosition = toString(val)
	case "target_city":
		p.TargetCity = toString(val)
	case "expected_salary":
		p.ExpectedSalary = toString(val)
	case "job_status":
		p.JobStatus = toString(val)
	case "self_introduction":
		p.SelfIntroduction = toString(val)
	}
	// birth_date 单独处理（time.Time）
	if key == "birth_date" {
		if s, ok := val.(string); ok && s != "" {
			if t, err := parseDate(s); err == nil {
				p.BirthDate = &t
			}
		} else if val == nil {
			p.BirthDate = nil
		}
	}
}

// toString 将 interface{} 安全转为 string
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case fmt.Stringer:
		return t.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

// GetCompletion 获取完成度百分比
func (s *ProfileService) GetCompletion(userID uint) (int, error) {
	return s.calcCompletion(userID), nil
}

// calcCompletion 计算完成度
// 权重：基础信息 30% + 教育 15% + 工作 15% + 项目 15% + 技能 10% + 荣誉 5% + 自我介绍 5% + 实践 5%
func (s *ProfileService) calcCompletion(userID uint) int {
	var profile models.UserProfile
	hasProfile := true
	if err := s.db.First(&profile, "user_id = ?", userID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		hasProfile = false
	}

	var eduCount, workCount, projCount, skillCount, honorCount, practiceCount int64
	s.db.Model(&models.UserEducation{}).Where("user_id = ?", userID).Count(&eduCount)
	s.db.Model(&models.UserWork{}).Where("user_id = ?", userID).Count(&workCount)
	s.db.Model(&models.UserProject{}).Where("user_id = ?", userID).Count(&projCount)
	s.db.Model(&models.UserSkill{}).Where("user_id = ?", userID).Count(&skillCount)
	s.db.Model(&models.UserHonor{}).Where("user_id = ?", userID).Count(&honorCount)
	s.db.Model(&models.UserPractice{}).Where("user_id = ?", userID).Count(&practiceCount)

	// 基础信息 30%：姓名+电话+目标岗位+期望薪资 等关键字段
	basicFields := 0
	if hasProfile {
		if profile.RealName != "" {
			basicFields++
		}
		if profile.Phone != "" {
			basicFields++
		}
		if profile.TargetPosition != "" {
			basicFields++
		}
		if profile.ExpectedSalary != "" {
			basicFields++
		}
		if profile.TargetCity != "" {
			basicFields++
		}
	}
	basicScore := float64(basicFields) / 5.0 * 30.0

	eduScore := 0.0
	if eduCount > 0 {
		eduScore = 15
	}
	workScore := 0.0
	if workCount > 0 {
		workScore = 15
	}
	projScore := 0.0
	if projCount > 0 {
		projScore = 15
	}
	skillScore := 0.0
	if skillCount > 0 {
		skillScore = 10
	}
	honorScore := 0.0
	if honorCount > 0 {
		honorScore = 5
	}
	practiceScore := 0.0
	if practiceCount > 0 {
		practiceScore = 5
	}
	introScore := 0.0
	if hasProfile && profile.SelfIntroduction != "" {
		introScore = 5
	}

	total := basicScore + eduScore + workScore + projScore + skillScore + honorScore + practiceScore + introScore
	if total > 100 {
		total = 100
	}
	return int(total)
}

// ---------- 子资源通用 CRUD ----------

// Educations
func (s *ProfileService) ListEducations(userID uint) ([]models.UserEducation, error) {
	var list []models.UserEducation
	err := s.db.Where("user_id = ?", userID).Order("id ASC").Find(&list).Error
	return list, err
}

func (s *ProfileService) CreateEducation(userID uint, m *models.UserEducation) (*models.UserEducation, error) {
	m.ID = 0
	m.UserID = userID
	if err := s.db.Create(m).Error; err != nil {
		return nil, err
	}
	s.touchCompletion(userID)
	return m, nil
}

func (s *ProfileService) UpdateEducation(userID, id uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.UserEducation{}).Where("id = ? AND user_id = ?", id, userID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

func (s *ProfileService) DeleteEducation(userID, id uint) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.UserEducation{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

// Works
func (s *ProfileService) ListWorks(userID uint) ([]models.UserWork, error) {
	var list []models.UserWork
	err := s.db.Where("user_id = ?", userID).Order("id ASC").Find(&list).Error
	return list, err
}

func (s *ProfileService) CreateWork(userID uint, m *models.UserWork) (*models.UserWork, error) {
	m.ID = 0
	m.UserID = userID
	if err := s.db.Create(m).Error; err != nil {
		return nil, err
	}
	s.touchCompletion(userID)
	return m, nil
}

func (s *ProfileService) UpdateWork(userID, id uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.UserWork{}).Where("id = ? AND user_id = ?", id, userID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

func (s *ProfileService) DeleteWork(userID, id uint) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.UserWork{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

// Projects
func (s *ProfileService) ListProjects(userID uint) ([]models.UserProject, error) {
	var list []models.UserProject
	err := s.db.Where("user_id = ?", userID).Order("id ASC").Find(&list).Error
	return list, err
}

func (s *ProfileService) CreateProject(userID uint, m *models.UserProject) (*models.UserProject, error) {
	m.ID = 0
	m.UserID = userID
	if err := s.db.Create(m).Error; err != nil {
		return nil, err
	}
	s.touchCompletion(userID)
	return m, nil
}

func (s *ProfileService) UpdateProject(userID, id uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.UserProject{}).Where("id = ? AND user_id = ?", id, userID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

func (s *ProfileService) DeleteProject(userID, id uint) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.UserProject{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

// Skills
func (s *ProfileService) ListSkills(userID uint) ([]models.UserSkill, error) {
	var list []models.UserSkill
	err := s.db.Where("user_id = ?", userID).Order("id ASC").Find(&list).Error
	return list, err
}

func (s *ProfileService) CreateSkill(userID uint, m *models.UserSkill) (*models.UserSkill, error) {
	m.ID = 0
	m.UserID = userID
	if err := s.db.Create(m).Error; err != nil {
		return nil, err
	}
	s.touchCompletion(userID)
	return m, nil
}

func (s *ProfileService) UpdateSkill(userID, id uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.UserSkill{}).Where("id = ? AND user_id = ?", id, userID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

func (s *ProfileService) DeleteSkill(userID, id uint) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.UserSkill{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

// Honors
func (s *ProfileService) ListHonors(userID uint) ([]models.UserHonor, error) {
	var list []models.UserHonor
	err := s.db.Where("user_id = ?", userID).Order("id ASC").Find(&list).Error
	return list, err
}

func (s *ProfileService) CreateHonor(userID uint, m *models.UserHonor) (*models.UserHonor, error) {
	m.ID = 0
	m.UserID = userID
	if err := s.db.Create(m).Error; err != nil {
		return nil, err
	}
	s.touchCompletion(userID)
	return m, nil
}

func (s *ProfileService) UpdateHonor(userID, id uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.UserHonor{}).Where("id = ? AND user_id = ?", id, userID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

func (s *ProfileService) DeleteHonor(userID, id uint) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.UserHonor{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

// Practices
func (s *ProfileService) ListPractices(userID uint) ([]models.UserPractice, error) {
	var list []models.UserPractice
	err := s.db.Where("user_id = ?", userID).Order("id ASC").Find(&list).Error
	return list, err
}

func (s *ProfileService) CreatePractice(userID uint, m *models.UserPractice) (*models.UserPractice, error) {
	m.ID = 0
	m.UserID = userID
	if err := s.db.Create(m).Error; err != nil {
		return nil, err
	}
	s.touchCompletion(userID)
	return m, nil
}

func (s *ProfileService) UpdatePractice(userID, id uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.UserPractice{}).Where("id = ? AND user_id = ?", id, userID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

func (s *ProfileService) DeletePractice(userID, id uint) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.UserPractice{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

// touchCompletion 重算并更新完成度
func (s *ProfileService) touchCompletion(userID uint) {
	pct := s.calcCompletion(userID)
	s.db.Model(&models.UserProfile{}).Where("user_id = ?", userID).Update("completion_pct", pct)
}

// ---------- 简历上传解析 ----------

// ParsedProfile LLM 解析后的结构化档案
type ParsedProfile struct {
	RealName        string                  `json:"real_name"`
	Gender          string                  `json:"gender"`
	Phone           string                  `json:"phone"`
	TargetPosition  string                  `json:"target_position"`
	SelfIntroduction string                 `json:"self_introduction"`
	Educations      []models.UserEducation  `json:"educations"`
	Works           []models.UserWork       `json:"works"`
	Projects        []models.UserProject    `json:"projects"`
	Skills          []models.UserSkill      `json:"skills"`
	Honors          []models.UserHonor      `json:"honors"`
	Practices       []models.UserPractice   `json:"practices"`
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
