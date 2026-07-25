package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/zhitu/server/internal/models"
	"gorm.io/gorm"
)

// 简历相关错误
var (
	ErrResumeNotFound    = errors.New("resume not found")
	ErrVersionNotFound   = errors.New("resume version not found")
	ErrEmptyResumeContent = errors.New("resume content is empty")
)

// ResumeService 简历业务逻辑（非 AI 部分）
type ResumeService struct {
	db       *gorm.DB
	profile  *ProfileService
}

// NewResumeService 构造 ResumeService
func NewResumeService(db *gorm.DB, profile *ProfileService) *ResumeService {
	return &ResumeService{db: db, profile: profile}
}

// ---------- 简历主表 CRUD ----------

// List 列出用户所有简历
func (s *ResumeService) List(userID uint) ([]models.Resume, error) {
	var list []models.Resume
	err := s.db.Where("user_id = ?", userID).Order("updated_at DESC").Find(&list).Error
	return list, err
}

// Get 获取单份简历（含当前版本内容）
func (s *ResumeService) Get(userID, id uint) (*models.Resume, error) {
	var resume models.Resume
	err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&resume).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrResumeNotFound
	}
	return &resume, err
}

// CreateResumeInput 创建简历入参
type CreateResumeInput struct {
	Name           string `json:"name" binding:"required"`
	TargetCompany  string `json:"target_company"`
	TargetPosition string `json:"target_position"`
	TargetJD       string `json:"target_jd"`
	Scene          string `json:"scene"`          // manual/jd/scenario
	InitialContent string `json:"initial_content"` // 可选初始内容 JSON
}

// Create 创建简历，并生成 v1.0 版本
func (s *ResumeService) Create(userID uint, in *CreateResumeInput) (*models.Resume, error) {
	if in.Scene == "" {
		in.Scene = "manual"
	}
	resume := &models.Resume{
		UserID:         userID,
		Name:           in.Name,
		TargetCompany:  in.TargetCompany,
		TargetPosition: in.TargetPosition,
		TargetJD:       in.TargetJD,
		Scene:          in.Scene,
	}
	if err := s.db.Create(resume).Error; err != nil {
		return nil, err
	}

	// 创建 v1.0 版本
	content := in.InitialContent
	if content == "" {
		content = DefaultResumeContent()
	}
	version, err := s.createVersion(resume.ID, "v1.0", content, "初始版本")
	if err != nil {
		return nil, err
	}

	// 更新当前版本指针
	resume.CurrentVersionID = version.ID
	if err := s.db.Model(resume).Update("current_version_id", version.ID).Error; err != nil {
		return nil, err
	}
	return resume, nil
}

// Update 更新简历元信息（不含 content）
func (s *ResumeService) Update(userID, id uint, updates map[string]interface{}) error {
	allowed := map[string]bool{
		"name": true, "target_company": true, "target_position": true, "target_jd": true, "scene": true,
	}
	filtered := map[string]interface{}{}
	for k, v := range updates {
		if allowed[k] {
			filtered[k] = v
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	result := s.db.Model(&models.Resume{}).Where("id = ? AND user_id = ?", id, userID).Updates(filtered)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrResumeNotFound
	}
	return nil
}

// Delete 删除简历（级联删除版本与 AI 操作记录）
func (s *ResumeService) Delete(userID, id uint) error {
	// 先校验归属
	var resume models.Resume
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&resume).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrResumeNotFound
	}
	// 级联删除
	if err := s.db.Where("resume_id = ?", id).Delete(&models.ResumeVersion{}).Error; err != nil {
		return err
	}
	if err := s.db.Where("resume_id = ?", id).Delete(&models.ResumeAIOperation{}).Error; err != nil {
		return err
	}
	return s.db.Delete(&resume).Error
}

// ---------- 版本管理 ----------

// ListVersions 列出简历的所有版本
func (s *ResumeService) ListVersions(userID, resumeID uint) ([]models.ResumeVersion, error) {
	if _, err := s.Get(userID, resumeID); err != nil {
		return nil, err
	}
	var list []models.ResumeVersion
	err := s.db.Where("resume_id = ?", resumeID).Order("id DESC").Find(&list).Error
	return list, err
}

// GetVersion 获取特定版本
func (s *ResumeService) GetVersion(userID, resumeID, versionID uint) (*models.ResumeVersion, error) {
	if _, err := s.Get(userID, resumeID); err != nil {
		return nil, err
	}
	var version models.ResumeVersion
	err := s.db.Where("id = ? AND resume_id = ?", versionID, resumeID).First(&version).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrVersionNotFound
	}
	return &version, err
}

// CreateVersionInput 创建版本入参
type CreateVersionInput struct {
	Content    string `json:"content" binding:"required"`
	ChangeNote string `json:"change_note"`
}

// CreateVersion 创建新版本（自动递增版本号），并把简历 current_version_id 指向新版本
func (s *ResumeService) CreateVersion(userID, resumeID uint, in *CreateVersionInput) (*models.ResumeVersion, error) {
	resume, err := s.Get(userID, resumeID)
	if err != nil {
		return nil, err
	}

	// 生成下一个版本号
	nextLabel := s.nextVersionLabel(resumeID)

	version, err := s.createVersion(resumeID, nextLabel, in.Content, in.ChangeNote)
	if err != nil {
		return nil, err
	}

	// 更新当前版本指针
	resume.CurrentVersionID = version.ID
	if err := s.db.Model(&models.Resume{}).Where("id = ?", resumeID).Update("current_version_id", version.ID).Error; err != nil {
		return nil, err
	}
	return version, nil
}

// RollbackVersion 回退到指定历史版本（实际是创建一个新版本，内容拷贝自目标版本）
func (s *ResumeService) RollbackVersion(userID, resumeID, versionID uint) (*models.ResumeVersion, error) {
	target, err := s.GetVersion(userID, resumeID, versionID)
	if err != nil {
		return nil, err
	}
	return s.CreateVersion(userID, resumeID, &CreateVersionInput{
		Content:    target.Content,
		ChangeNote: fmt.Sprintf("回退到 %s", target.VersionLabel),
	})
}

// createVersion 内部创建版本记录
func (s *ResumeService) createVersion(resumeID uint, label, content, note string) (*models.ResumeVersion, error) {
	v := &models.ResumeVersion{
		ResumeID:     resumeID,
		VersionLabel: label,
		Content:      content,
		ChangeNote:   note,
	}
	if err := s.db.Create(v).Error; err != nil {
		return nil, err
	}
	return v, nil
}

// nextVersionLabel 生成下一个版本号
// 规则：v1.0 → v1.1 → v1.2 ...；若已有 v1.x 则取最大 minor+1
func (s *ResumeService) nextVersionLabel(resumeID uint) string {
	var versions []models.ResumeVersion
	s.db.Where("resume_id = ?", resumeID).Order("id ASC").Find(&versions)

	if len(versions) == 0 {
		return "v1.0"
	}

	maxMajor := 1
	maxMinor := 0
	for _, v := range versions {
		major, minor := parseVersionLabel(v.VersionLabel)
		if major > maxMajor || (major == maxMajor && minor > maxMinor) {
			maxMajor = major
			maxMinor = minor
		}
	}
	return fmt.Sprintf("v%d.%d", maxMajor, maxMinor+1)
}

// parseVersionLabel 解析 v1.2 格式
func parseVersionLabel(label string) (major, minor int) {
	label = strings.TrimPrefix(label, "v")
	parts := strings.Split(label, ".")
	if len(parts) >= 1 {
		fmt.Sscanf(parts[0], "%d", &major)
	}
	if len(parts) >= 2 {
		fmt.Sscanf(parts[1], "%d", &minor)
	}
	return
}

// ---------- 简历区同步到档案 ----------

// SyncToProfile 将简历 content 中的信息同步到用户档案
// 策略：仅同步非空字段到档案主表与子表（追加，不覆盖原有）
func (s *ResumeService) SyncToProfile(userID, resumeID uint) error {
	resume, err := s.Get(userID, resumeID)
	if err != nil {
		return err
	}

	// 取当前版本内容
	var version models.ResumeVersion
	if err := s.db.Where("id = ? AND resume_id = ?", resume.CurrentVersionID, resumeID).First(&version).Error; err != nil {
		return ErrVersionNotFound
	}

	var content ResumeContent
	if err := json.Unmarshal([]byte(version.Content), &content); err != nil {
		return fmt.Errorf("parse resume content: %w", err)
	}

	// 同步主表字段（仅非空）
	updates := map[string]interface{}{}
	if content.Personal.Name != "" {
		updates["real_name"] = content.Personal.Name
	}
	if content.Personal.Phone != "" {
		updates["phone"] = content.Personal.Phone
	}
	if content.Intention.Position != "" {
		updates["target_position"] = content.Intention.Position
	}
	if content.Intention.City != "" {
		updates["target_city"] = content.Intention.City
	}
	if content.Intention.Salary != "" {
		updates["expected_salary"] = content.Intention.Salary
	}
	if len(updates) > 0 {
		if _, err := s.profile.UpsertProfile(userID, updates); err != nil {
			return err
		}
	}
	return nil
}

// ---------- 简历 content 结构 ----------

// ResumeContent 简历版本 content JSON 结构
type ResumeContent struct {
	Personal         ResumePersonal     `json:"personal"`
	Intention        ResumeIntention    `json:"intention"`
	Education        []ResumeEducation  `json:"education"`
	Work             []ResumeWork       `json:"work"`
	Project          []ResumeProject    `json:"project"`
	Skills           []ResumeSkill      `json:"skills"`
	Honor            []ResumeHonor      `json:"honor"`
	Custom           []ResumeCustom     `json:"custom"`
	ModuleOrder      []string           `json:"module_order"`
	ModuleVisibility map[string]bool    `json:"module_visibility"`
}

// ResumePersonal 简历中的个人信息
type ResumePersonal struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Age    string `json:"age"`
	Phone  string `json:"phone"`
	Email  string `json:"email"`
	GitHub string `json:"github"`
	Avatar string `json:"avatar"`
	City   string `json:"city"`
}

// ResumeIntention 求职意向
type ResumeIntention struct {
	Position string `json:"position"`
	City     string `json:"city"`
	Salary   string `json:"salary"`
	Arrival  string `json:"arrival"`
	Industry string `json:"industry"`
}

// ResumeEducation 简历中的教育背景
type ResumeEducation struct {
	School   string `json:"school"`
	Major    string `json:"major"`
	Degree   string `json:"degree"`
	Start    string `json:"start"`
	End      string `json:"end"`
	Courses  string `json:"courses"`
	GPA      string `json:"gpa"`
}

// ResumeWork 简历中的工作经历
type ResumeWork struct {
	Company     string `json:"company"`
	Position    string `json:"position"`
	Start       string `json:"start"`
	End         string `json:"end"`
	Description string `json:"description"`
	LeaveReason string `json:"leave_reason"`
}

// ResumeProject 简历中的项目经历
type ResumeProject struct {
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	Start       string   `json:"start"`
	End         string   `json:"end"`
	Description string   `json:"description"`
	TechStack   []string `json:"tech_stack"`
	URL         string   `json:"url"`
}

// ResumeSkill 简历中的技能
type ResumeSkill struct {
	Category    string `json:"category"`
	Name        string `json:"name"`
	Proficiency string `json:"proficiency"`
}

// ResumeHonor 简历中的荣誉
type ResumeHonor struct {
	Name   string `json:"name"`
	Issuer string `json:"issuer"`
	Date   string `json:"date"`
	Level  string `json:"level"`
}

// ResumeCustom 自定义模块
type ResumeCustom struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// DefaultResumeContent 返回空白简历 content JSON
func DefaultResumeContent() string {
	defaultOrder := []string{"personal", "intention", "education", "work", "project", "skills", "honor"}
	visibility := map[string]bool{}
	for _, m := range defaultOrder {
		visibility[m] = true
	}
	content := ResumeContent{
		Personal:         ResumePersonal{},
		Intention:        ResumeIntention{},
		Education:        []ResumeEducation{},
		Work:             []ResumeWork{},
		Project:          []ResumeProject{},
		Skills:           []ResumeSkill{},
		Honor:            []ResumeHonor{},
		Custom:           []ResumeCustom{},
		ModuleOrder:      defaultOrder,
		ModuleVisibility: visibility,
	}
	b, _ := json.Marshal(content)
	return string(b)
}
