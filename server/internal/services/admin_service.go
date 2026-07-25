package services

import (
	"errors"
	"fmt"

	"github.com/zhitu/server/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminService 管理端服务
type AdminService struct {
	db *gorm.DB
}

// NewAdminService 构造函数
func NewAdminService(db *gorm.DB) *AdminService {
	return &AdminService{db: db}
}

// DashboardStats 仪表盘统计
type DashboardStats struct {
	TotalUsers     int64 `json:"total_users"`
	ActiveUsers    int64 `json:"active_users"`
	DisabledUsers  int64 `json:"disabled_users"`
	TotalResumes   int64 `json:"total_resumes"`
	TotalDeliveries int64 `json:"total_deliveries"`
	OfferCount     int64 `json:"offer_count"`
	RejectedCount  int64 `json:"rejected_count"`
	TotalInterviews int64 `json:"total_interviews"`
}

// GetDashboardStats 获取仪表盘统计数据
func (s *AdminService) GetDashboardStats() (*DashboardStats, error) {
	var stats DashboardStats

	if err := s.db.Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.User{}).Where("status = ?", "active").Count(&stats.ActiveUsers).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.User{}).Where("status = ?", "disabled").Count(&stats.DisabledUsers).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Resume{}).Count(&stats.TotalResumes).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Delivery{}).Count(&stats.TotalDeliveries).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Delivery{}).Where("status = ?", "offer").Count(&stats.OfferCount).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Delivery{}).Where("status = ?", "rejected").Count(&stats.RejectedCount).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Interview{}).Count(&stats.TotalInterviews).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

// AdminUserList 管理员用户列表结果
type AdminUserList struct {
	Total int64         `json:"total"`
	List  []models.User `json:"list"`
}

// ListUsers 分页获取用户列表
func (s *AdminService) ListUsers(page, pageSize int, keyword, status string) (*AdminUserList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := s.db.Model(&models.User{})

	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("email LIKE ? OR nickname LIKE ?", like, like)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var users []models.User
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, err
	}

	return &AdminUserList{Total: total, List: users}, nil
}

// GetUserDetail 获取用户详情（含档案统计、简历数、投递数、面试数）
type UserDetail struct {
	models.User
	ResumeCount    int64 `json:"resume_count"`
	DeliveryCount  int64 `json:"delivery_count"`
	InterviewCount int64 `json:"interview_count"`
}

func (s *AdminService) GetUserDetail(userID uint) (*UserDetail, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}

	detail := &UserDetail{User: user}

	s.db.Model(&models.Resume{}).Where("user_id = ?", userID).Count(&detail.ResumeCount)
	s.db.Model(&models.Delivery{}).Where("user_id = ?", userID).Count(&detail.DeliveryCount)
	s.db.Model(&models.Interview{}).Where("user_id = ?", userID).Count(&detail.InterviewCount)

	return detail, nil
}

// ToggleUserStatus 切换用户状态（启用/禁用）
func (s *AdminService) ToggleUserStatus(userID uint, status string) error {
	if status != "active" && status != "disabled" {
		return fmt.Errorf("无效的状态值")
	}

	result := s.db.Model(&models.User{}).Where("id = ?", userID).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在")
	}
	return nil
}

// ResetUserPassword 重置用户密码
func (s *AdminService) ResetUserPassword(userID uint, newPassword string) error {
	if len(newPassword) < 6 {
		return fmt.Errorf("密码长度不能少于6位")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	result := s.db.Model(&models.User{}).Where("id = ?", userID).Update("password", string(hashed))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在")
	}
	return nil
}

// AdminDeliveryList 管理员投递列表结果
type AdminDeliveryList struct {
	Total int64              `json:"total"`
	List  []models.Delivery  `json:"list"`
}

// ListDeliveries 分页获取全局投递列表
func (s *AdminService) ListDeliveries(page, pageSize int, status, company, userEmail string) (*AdminDeliveryList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := s.db.Model(&models.Delivery{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if company != "" {
		query = query.Where("company LIKE ?", "%"+company+"%")
	}
	if userEmail != "" {
		query = query.Where("user_id IN (SELECT id FROM users WHERE email LIKE ?)", "%"+userEmail+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var deliveries []models.Delivery
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&deliveries).Error; err != nil {
		return nil, err
	}

	return &AdminDeliveryList{Total: total, List: deliveries}, nil
}

// AdminDeliveryFunnel 全局投递漏斗统计
type AdminDeliveryFunnel struct {
	Applied        int64   `json:"applied"`
	WrittenTestPass int64  `json:"written_test_pass"`
	FirstPass      int64   `json:"first_pass"`
	SecondPass     int64   `json:"second_pass"`
	OfferCount     int64   `json:"offer_count"`
	WrittenTestRate float64 `json:"written_test_rate"`
	FirstRate      float64 `json:"first_rate"`
	SecondRate     float64 `json:"second_rate"`
	OfferRate      float64 `json:"offer_rate"`
}

// GetDeliveryFunnel 全局投递漏斗
func (s *AdminService) GetDeliveryFunnel() (*AdminDeliveryFunnel, error) {
	var funnel AdminDeliveryFunnel

	s.db.Model(&models.Delivery{}).Count(&funnel.Applied)

	// 通过面试轮次记录精确计算
	s.db.Model(&models.DeliveryRound{}).
		Where("round_type = ? AND result = ?", "written_test", "pass").
		Count(&funnel.WrittenTestPass)

	s.db.Model(&models.DeliveryRound{}).
		Where("round_type = ? AND result = ?", "first_tech", "pass").
		Count(&funnel.FirstPass)

	s.db.Model(&models.DeliveryRound{}).
		Where("round_type = ? AND result = ?", "second_tech", "pass").
		Count(&funnel.SecondPass)

	s.db.Model(&models.Delivery{}).Where("status = ?", "offer").Count(&funnel.OfferCount)

	if funnel.Applied > 0 {
		funnel.WrittenTestRate = float64(funnel.WrittenTestPass) / float64(funnel.Applied) * 100
	}
	if funnel.WrittenTestPass > 0 {
		funnel.FirstRate = float64(funnel.FirstPass) / float64(funnel.WrittenTestPass) * 100
	}
	if funnel.FirstPass > 0 {
		funnel.SecondRate = float64(funnel.SecondPass) / float64(funnel.FirstPass) * 100
	}
	if funnel.Applied > 0 {
		funnel.OfferRate = float64(funnel.OfferCount) / float64(funnel.Applied) * 100
	}

	return &funnel, nil
}
