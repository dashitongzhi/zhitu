// Package models 定义 GORM 数据模型
package models

import (
	"time"
)

// User 普通用户模型
// 注意：管理员账号不在此表中，凭据直接来自 configs/config.yaml
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"` // bcrypt 哈希，不出现在 JSON
	Nickname  string    `gorm:"size:50" json:"nickname"`
	Avatar    string    `gorm:"size:500" json:"avatar"`
	Status    string    `gorm:"size:20;default:active" json:"status"` // active / disabled
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 显式指定表名
func (User) TableName() string {
	return "users"
}
