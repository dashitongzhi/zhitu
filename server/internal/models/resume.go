package models

import "time"

// Resume 简历主表（一个用户可有多份简历）
type Resume struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UserID           uint      `gorm:"index;not null" json:"user_id"`
	Name             string    `gorm:"size:100;not null" json:"name"` // 简历名称
	TargetCompany    string    `gorm:"size:100" json:"target_company"`
	TargetPosition   string    `gorm:"size:100" json:"target_position"`
	TargetJD         string    `gorm:"type:text" json:"target_jd"`         // 目标 JD 全文
	Scene            string    `gorm:"size:20;default:manual" json:"scene"` // manual / jd / scenario
	CurrentVersionID uint      `gorm:"default:0" json:"current_version_id"` // 当前生效版本 ID
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (Resume) TableName() string { return "resumes" }

// ResumeVersion 简历版本快照
// 每次保存或 AI 操作（生成/润色）后生成新版本
type ResumeVersion struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ResumeID    uint      `gorm:"index;not null" json:"resume_id"`
	VersionLabel string   `gorm:"size:20" json:"version_label"` // v1.0、v1.1
	Content     string    `gorm:"type:text" json:"content"`     // 完整简历数据 JSON 字符串
	ChangeNote  string    `gorm:"size:200" json:"change_note"`
	CreatedAt   time.Time `json:"created_at"`
}

func (ResumeVersion) TableName() string { return "resume_versions" }

// ResumeAIOperation 简历 AI 操作记录
type ResumeAIOperation struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ResumeID        uint      `gorm:"index;not null" json:"resume_id"`
	VersionIDBefore  uint     `gorm:"default:0" json:"version_id_before"`
	VersionIDAfter   uint     `gorm:"default:0" json:"version_id_after"`
	OperationType   string    `gorm:"size:20" json:"operation_type"` // generate/polish/score/jd_match
	InputJD         string    `gorm:"type:text" json:"input_jd"`
	Result          string    `gorm:"type:text" json:"result"` // AI 返回结果 JSON 字符串
	CreatedAt       time.Time `json:"created_at"`
}

func (ResumeAIOperation) TableName() string { return "resume_ai_operations" }
