package services

import (
	"errors"
	"fmt"
	"github.com/zhitu/server/internal/config"
	"github.com/zhitu/server/internal/models"
	"gorm.io/gorm"
)

// 档案相关错误
var (
	ErrProfileNotFound     = errors.New("profile not found, please create first")
	ErrSubResourceNotFound = errors.New("sub resource not found")
	ErrResumeFileEmpty     = errors.New("resume file content is empty or unparseable")
	ErrResumeParseFail     = errors.New("failed to parse resume via llm")
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
	UserProfile *models.UserProfile    `json:"profile"`
	Educations  []models.UserEducation `json:"educations"`
	Works       []models.UserWork      `json:"works"`
	Projects    []models.UserProject   `json:"projects"`
	Skills      []models.UserSkill     `json:"skills"`
	Honors      []models.UserHonor     `json:"honors"`
	Practices   []models.UserPractice  `json:"practices"`
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
