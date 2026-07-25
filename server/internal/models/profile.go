package models

import (
	"time"
)

// UserProfile 用户档案主表（1:1 与 users 关联）
// 注意：User 表仅保留认证字段，档案详细信息均在此表与子表
type UserProfile struct {
	UserID          uint       `gorm:"primaryKey;uniqueIndex" json:"user_id"`
	RealName        string     `gorm:"size:50" json:"real_name"`
	Gender          string     `gorm:"size:10" json:"gender"` // male/female/other
	BirthDate       *time.Time `json:"birth_date,omitempty"`
	Phone           string     `gorm:"size:20" json:"phone"`
	TargetPosition  string     `gorm:"size:100" json:"target_position"`
	TargetCity      string     `gorm:"size:200" json:"target_city"` // 逗号分隔多选
	ExpectedSalary  string     `gorm:"size:50" json:"expected_salary"`
	JobStatus       string     `gorm:"size:20" json:"job_status"` // fresh/graduated/employed/resigned
	SelfIntroduction string    `gorm:"type:text" json:"self_introduction"`
	CompletionPct   int        `gorm:"default:0" json:"completion_pct"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// TableName 显式指定表名
func (UserProfile) TableName() string { return "user_profiles" }

// UserEducation 教育背景
type UserEducation struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	School     string    `gorm:"size:100;not null" json:"school"`
	Major      string    `gorm:"size:100" json:"major"`
	Degree     string    `gorm:"size:20" json:"degree"` // 大专/本科/硕士/博士
	StartDate  string    `gorm:"size:10" json:"start_date"` // YYYY-MM
	EndDate    string    `gorm:"size:10" json:"end_date"`   // YYYY-MM 或 "至今"
	GPA        string    `gorm:"size:20" json:"gpa"`
	Courses    string    `gorm:"type:text" json:"courses"` // 主修课程，逗号分隔
	Exchange   string    `gorm:"size:200" json:"exchange"` // 交换经历
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (UserEducation) TableName() string { return "user_educations" }

// UserWork 工作/实习经历
type UserWork struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	Company     string    `gorm:"size:100;not null" json:"company"`
	Position    string    `gorm:"size:100" json:"position"`
	StartDate   string    `gorm:"size:10" json:"start_date"`
	EndDate     string    `gorm:"size:10" json:"end_date"`
	Description string    `gorm:"type:text" json:"description"` // STAR 法则描述
	LeaveReason string    `gorm:"size:200" json:"leave_reason"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (UserWork) TableName() string { return "user_works" }

// UserProject 项目经历
type UserProject struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Role        string    `gorm:"size:100" json:"role"`
	StartDate   string    `gorm:"size:10" json:"start_date"`
	EndDate     string    `gorm:"size:10" json:"end_date"`
	Description string    `gorm:"type:text" json:"description"`
	TechStack   string    `gorm:"type:text" json:"tech_stack"` // JSON 数组字符串 ["Go","React"]
	URL         string    `gorm:"size:500" json:"url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (UserProject) TableName() string { return "user_projects" }

// UserSkill 技能特长
type UserSkill struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	Category    string    `gorm:"size:30;not null" json:"category"` // 编程语言/框架工具/软技能/证书/语言能力
	Name        string    `gorm:"size:100;not null" json:"name"`
	Proficiency string    `gorm:"size:20" json:"proficiency"` // 了解/熟悉/熟练/精通
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (UserSkill) TableName() string { return "user_skills" }

// UserHonor 荣誉奖项
type UserHonor struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Issuer    string    `gorm:"size:100" json:"issuer"`
	AwardDate string    `gorm:"size:10" json:"award_date"`
	Level     string    `gorm:"size:20" json:"level"` // 校级/市级/省级/国家级/国际级
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (UserHonor) TableName() string { return "user_honors" }

// UserPractice 校内外实践
type UserPractice struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	Title       string    `gorm:"size:100;not null" json:"title"`
	Organization string   `gorm:"size:100" json:"organization"`
	StartDate   string    `gorm:"size:10" json:"start_date"`
	EndDate     string    `gorm:"size:10" json:"end_date"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (UserPractice) TableName() string { return "user_practices" }
