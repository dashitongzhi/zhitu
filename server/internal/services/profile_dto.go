package services

import (
	"github.com/zhitu/server/internal/models"
)

// EducationDTO 教育背景请求 DTO
type EducationDTO struct {
	School    string `json:"school" binding:"required"`
	Major     string `json:"major"`
	Degree    string `json:"degree"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	GPA       string `json:"gpa"`
	Courses   string `json:"courses"`
	Exchange  string `json:"exchange"`
}

// ToModel 转换为模型
func (d *EducationDTO) ToModel() *models.UserEducation {
	return &models.UserEducation{
		School:    d.School,
		Major:     d.Major,
		Degree:    d.Degree,
		StartDate: d.StartDate,
		EndDate:   d.EndDate,
		GPA:       d.GPA,
		Courses:   d.Courses,
		Exchange:  d.Exchange,
	}
}

// WorkDTO 工作经历请求 DTO
type WorkDTO struct {
	Company     string `json:"company" binding:"required"`
	Position    string `json:"position"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Description string `json:"description"`
	LeaveReason string `json:"leave_reason"`
}

func (d *WorkDTO) ToModel() *models.UserWork {
	return &models.UserWork{
		Company:     d.Company,
		Position:    d.Position,
		StartDate:   d.StartDate,
		EndDate:     d.EndDate,
		Description: d.Description,
		LeaveReason: d.LeaveReason,
	}
}

// ProjectDTO 项目经历请求 DTO
type ProjectDTO struct {
	Name        string   `json:"name" binding:"required"`
	Role        string   `json:"role"`
	StartDate   string   `json:"start_date"`
	EndDate     string   `json:"end_date"`
	Description string   `json:"description"`
	TechStack   []string `json:"tech_stack"`
	URL         string   `json:"url"`
}

func (d *ProjectDTO) ToModel() *models.UserProject {
	return &models.UserProject{
		Name:        d.Name,
		Role:        d.Role,
		StartDate:   d.StartDate,
		EndDate:     d.EndDate,
		Description: d.Description,
		TechStack:   MarshalTechStack(d.TechStack),
		URL:         d.URL,
	}
}

// SkillDTO 技能请求 DTO
type SkillDTO struct {
	Category    string `json:"category" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Proficiency string `json:"proficiency"`
}

func (d *SkillDTO) ToModel() *models.UserSkill {
	return &models.UserSkill{
		Category:    d.Category,
		Name:        d.Name,
		Proficiency: d.Proficiency,
	}
}

// HonorDTO 荣誉奖项请求 DTO
type HonorDTO struct {
	Name      string `json:"name" binding:"required"`
	Issuer    string `json:"issuer"`
	AwardDate string `json:"award_date"`
	Level     string `json:"level"`
}

func (d *HonorDTO) ToModel() *models.UserHonor {
	return &models.UserHonor{
		Name:      d.Name,
		Issuer:    d.Issuer,
		AwardDate: d.AwardDate,
		Level:     d.Level,
	}
}

// PracticeDTO 校内外实践请求 DTO
type PracticeDTO struct {
	Title        string `json:"title" binding:"required"`
	Organization string `json:"organization"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Description  string `json:"description"`
}

func (d *PracticeDTO) ToModel() *models.UserPractice {
	return &models.UserPractice{
		Title:        d.Title,
		Organization: d.Organization,
		StartDate:    d.StartDate,
		EndDate:      d.EndDate,
		Description:  d.Description,
	}
}
