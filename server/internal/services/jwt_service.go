// Package services 承载业务逻辑，对 handler 层提供无 HTTP 上下文的纯函数接口
package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zhitu/server/internal/config"
	"github.com/zhitu/server/internal/models"
)

// JWTService 负责 JWT 的签发与解析
// 用户与管理员使用各自独立的签名密钥，互不混用
type JWTService struct {
	secret      []byte // 普通用户签名密钥
	adminSecret []byte // 管理员签名密钥
	issuer      string
	ttl         time.Duration
}

// Claims 自定义 JWT 载荷
type Claims struct {
	UserID  uint   `json:"user_id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// NewJWTService 从配置构造 JWTService
func NewJWTService(cfg *config.JWTConfig) *JWTService {
	return &JWTService{
		secret:      []byte(cfg.Secret),
		adminSecret: []byte(cfg.AdminSecret),
		issuer:      cfg.Issuer,
		ttl:         time.Duration(cfg.ExpireHours) * time.Hour,
	}
}

// TTLSeconds 返回 Token 有效期（秒）
func (s *JWTService) TTLSeconds() int {
	return int(s.ttl.Seconds())
}

// GenerateForUser 为普通用户签发 Token（使用用户密钥签名）
func (s *JWTService) GenerateForUser(u *models.User) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:  u.ID,
		Email:   u.Email,
		IsAdmin: false,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   fmt.Sprintf("%d", u.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
		},
	}
	return s.signWith(claims, s.secret)
}

// GenerateForAdmin 为管理员签发 Token（使用管理员密钥签名，UserID 固定为 0）
func (s *JWTService) GenerateForAdmin(email string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:  0,
		Email:   email,
		IsAdmin: true,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   "admin",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
		},
	}
	return s.signWith(claims, s.adminSecret)
}

// Parse 解析并校验 Token，返回 Claims
// 依次尝试用户密钥与管理员密钥，任一成功即可
func (s *JWTService) Parse(tokenStr string) (*Claims, error) {
	// 先用用户密钥解析
	if claims, err := s.parseWith(tokenStr, s.secret); err == nil {
		return claims, nil
	}
	// 再用管理员密钥解析
	if claims, err := s.parseWith(tokenStr, s.adminSecret); err == nil {
		return claims, nil
	}
	return nil, errors.New("invalid or expired token")
}

func (s *JWTService) signWith(claims Claims, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (s *JWTService) parseWith(tokenStr string, secret []byte) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
