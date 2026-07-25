package services

import (
	"errors"
	"strings"

	"github.com/zhitu/server/internal/config"
	"github.com/zhitu/server/internal/models"
	"github.com/zhitu/server/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 邮箱长度限制
const (
	EmailMinLength = 3
	EmailMaxLength = 255
	NicknameMaxLen = 50
)

// 业务错误
var (
	ErrEmailEmpty       = errors.New("email is required")
	ErrEmailInvalidFmt  = errors.New("email format is invalid")
	ErrEmailTooLong     = errors.New("email is too long")
	ErrEmailRegistered  = errors.New("email already registered")
	ErrUserNotFound     = errors.New("user not found")
	ErrPasswordMismatch = errors.New("email or password is incorrect")
	ErrOldPwdMismatch   = errors.New("old password is incorrect")
	ErrAdminCredential  = errors.New("admin email or password is incorrect")
)

// AuthService 用户认证相关业务逻辑
type AuthService struct {
	db    *gorm.DB
	admin *config.AdminConfig
}

// NewAuthService 构造 AuthService
func NewAuthService(db *gorm.DB, adminCfg *config.AdminConfig) *AuthService {
	return &AuthService{
		db:    db,
		admin: adminCfg,
	}
}

// Register 注册新用户
func (s *AuthService) Register(email, password, nickname string) (*models.User, error) {
	email = normalizeEmail(email)
	if err := validateEmail(email); err != nil {
		return nil, err
	}
	if err := utils.ValidatePasswordStrength(password); err != nil {
		return nil, err
	}
	if len(nickname) > NicknameMaxLen {
		return nil, errors.New("nickname is too long")
	}

	// 检查邮箱是否已被注册
	var count int64
	if err := s.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrEmailRegistered
	}

	// 禁止使用管理员邮箱注册普通账号
	if s.admin != nil && email == normalizeEmail(s.admin.Email) {
		return nil, ErrEmailRegistered
	}

	hashed, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    email,
		Password: hashed,
		Nickname: strings.TrimSpace(nickname),
	}
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Login 普通用户邮箱+密码登录
func (s *AuthService) Login(email, password string) (*models.User, error) {
	email = normalizeEmail(email)
	if err := validateEmail(email); err != nil {
		return nil, err
	}
	if password == "" {
		return nil, ErrPasswordMismatch
	}

	var user models.User
	err := s.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPasswordMismatch
	}
	if err != nil {
		return nil, err
	}

	if !utils.ComparePassword(user.Password, password) {
		return nil, ErrPasswordMismatch
	}
	return &user, nil
}

// ChangePassword 修改密码（需提供旧密码）
func (s *AuthService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	if err := utils.ValidatePasswordStrength(newPassword); err != nil {
		return err
	}
	if oldPassword == newPassword {
		return errors.New("new password must be different from the old one")
	}

	var user models.User
	err := s.db.First(&user, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}

	if !utils.ComparePassword(user.Password, oldPassword) {
		return ErrOldPwdMismatch
	}

	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	return s.db.Model(&user).Update("password", hashed).Error
}

// AdminLogin 管理员登录：直接比对 config 中的凭据
// 兼容明文存储与 bcrypt 哈希存储（以 $2a$/$2b$/$2y$ 开头时按哈希处理）
func (s *AuthService) AdminLogin(email, password string) (string, error) {
	if s.admin == nil {
		return "", ErrAdminCredential
	}
	if normalizeEmail(email) != normalizeEmail(s.admin.Email) {
		return "", ErrAdminCredential
	}
	if !safeCompareAdminPassword(password, s.admin.Password) {
		return "", ErrAdminCredential
	}
	return s.admin.Email, nil
}

// GetProfile 根据 ID 获取用户信息
func (s *AuthService) GetProfile(userID uint) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// normalizeEmail 规范化邮箱（去空格 + 转小写）
func normalizeEmail(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// validateEmail 邮箱格式与长度校验
func validateEmail(email string) error {
	if email == "" {
		return ErrEmailEmpty
	}
	if len(email) > EmailMaxLength {
		return ErrEmailTooLong
	}
	// 简单邮箱格式校验：必须含 @ 且本地段和域名段非空，域名段含一个 .
	at := strings.IndexByte(email, '@')
	if at <= 0 || at == len(email)-1 {
		return ErrEmailInvalidFmt
	}
	domain := email[at+1:]
	if !strings.Contains(domain, ".") {
		return ErrEmailInvalidFmt
	}
	return nil
}

// safeCompareAdminPassword 兼容 config 中明文 / bcrypt 哈希两种存储方式
func safeCompareAdminPassword(plain, stored string) bool {
	if plain == stored {
		return true
	}
	// 若存储的是 bcrypt 哈希（$2a$ / $2b$ 开头），用 bcrypt 比较
	if strings.HasPrefix(stored, "$2a$") || strings.HasPrefix(stored, "$2b$") || strings.HasPrefix(stored, "$2y$") {
		return bcrypt.CompareHashAndPassword([]byte(stored), []byte(plain)) == nil
	}
	return false
}
