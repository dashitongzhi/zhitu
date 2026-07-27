// Package handlers 实现 HTTP 请求处理，调用 services 完成业务并返回统一响应
package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zhitu/server/internal/services"
	"github.com/zhitu/server/internal/utils"
)

// AuthHandler 认证相关路由处理器
type AuthHandler struct {
	auth *services.AuthService
	jwt  *services.JWTService
}

// NewAuthHandler 构造 AuthHandler
func NewAuthHandler(auth *services.AuthService, jwt *services.JWTService) *AuthHandler {
	return &AuthHandler{auth: auth, jwt: jwt}
}

// ---------- DTO ----------

type registerReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=64"`
	Nickname string `json:"nickname" binding:"omitempty,max=50"`
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type adminLoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type changePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=64"`
}

type tokenResp struct {
	Token    string `json:"token"`
	TokenType string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"` // 秒
}

type userResp struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	IsAdmin  bool   `json:"is_admin"`
}

// ---------- Handlers ----------

// Register 普通用户注册
// POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	user, err := h.auth.Register(req.Email, req.Password, req.Nickname)
	if err != nil {
		switch err {
		case services.ErrEmailRegistered:
			utils.Conflict(c, err.Error())
		case services.ErrEmailEmpty, services.ErrEmailInvalidFmt, services.ErrEmailTooLong:
			utils.BadRequest(c, err.Error())
		default:
			// 密码强度错误等
			utils.BadRequest(c, err.Error())
		}
		return
	}

	token, err := h.jwt.GenerateForUser(user)
	if err != nil {
		utils.InternalError(c, "generate token failed")
		return
	}

	utils.OKWithMsg(c, "register success", tokenResp{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: h.jwt.TTLSeconds(),
	})
}

// Login 用户登录。
// 配置中的管理员凭据也可以从统一登录入口进入产品工作台，
// 同时保留 /api/auth/admin/login 供独立管理后台使用。
// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	user, err := h.auth.Login(req.Email, req.Password)
	if err != nil {
		adminEmail, adminErr := h.auth.AdminLogin(req.Email, req.Password)
		if adminErr != nil {
			utils.Unauthorized(c, err.Error())
			return
		}

		token, tokenErr := h.jwt.GenerateForAdmin(adminEmail)
		if tokenErr != nil {
			utils.InternalError(c, "generate token failed")
			return
		}

		utils.OKWithMsg(c, "admin login success", tokenResp{
			Token:     token,
			TokenType: "Bearer",
			ExpiresIn: h.jwt.TTLSeconds(),
		})
		return
	}

	token, err := h.jwt.GenerateForUser(user)
	if err != nil {
		utils.InternalError(c, "generate token failed")
		return
	}

	utils.OKWithMsg(c, "login success", tokenResp{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: h.jwt.TTLSeconds(),
	})
}

// AdminLogin 管理员登录（凭据来自 config，不入库）
// POST /api/auth/admin/login
func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var req adminLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	email, err := h.auth.AdminLogin(req.Email, req.Password)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	token, err := h.jwt.GenerateForAdmin(email)
	if err != nil {
		utils.InternalError(c, "generate token failed")
		return
	}

	utils.OKWithMsg(c, "admin login success", tokenResp{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: h.jwt.TTLSeconds(),
	})
}

// ChangePassword 修改当前用户密码（需登录）
// POST /api/auth/change-password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := getCurrentUserID(c)
	if !exists {
		utils.Unauthorized(c, "missing authentication")
		return
	}

	var req changePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	if err := h.auth.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		switch err {
		case services.ErrOldPwdMismatch:
			utils.BadRequest(c, err.Error())
		case services.ErrUserNotFound:
			utils.NotFound(c, err.Error())
		default:
			utils.BadRequest(c, err.Error())
		}
		return
	}

	utils.OKWithMsg(c, "password changed", nil)
}

// Me 获取当前登录用户信息（需登录）
// GET /api/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := getCurrentUserID(c)
	if !exists {
		utils.Unauthorized(c, "missing authentication")
		return
	}

	isAdmin, _ := c.Get("is_admin")
	isAdminBool, _ := isAdmin.(bool)

	if isAdminBool {
		emailAny, _ := c.Get("email")
		emailStr, _ := emailAny.(string)
		utils.OK(c, userResp{
			ID:       0,
			Email:    emailStr,
			Nickname: "admin",
			IsAdmin:  true,
		})
		return
	}

	user, err := h.auth.GetProfile(userID)
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			utils.NotFound(c, err.Error())
		default:
			utils.InternalError(c, err.Error())
		}
		return
	}

	utils.OK(c, userResp{
		ID:       user.ID,
		Email:    user.Email,
		Nickname: strings.TrimSpace(user.Nickname),
		Avatar:   user.Avatar,
		IsAdmin:  false,
	})
}

// getCurrentUserID 从 gin.Context 读取 JWT 中间件注入的 user_id
func getCurrentUserID(c *gin.Context) (uint, bool) {
	v, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	id, ok := v.(uint)
	return id, ok
}
