// Package middleware 集中放置所有 gin 中间件
package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zhitu/server/internal/services"
	"github.com/zhitu/server/internal/utils"
)

// Header 名称
const (
	HeaderAuthorization = "Authorization"
	ContextUserID       = "user_id"
	ContextEmail        = "email"
	ContextIsAdmin      = "is_admin"
)

// JWTAuth 校验 Authorization: Bearer <token>，将 user_id/email/is_admin 注入 context
// 未通过则直接返回 401
func JWTAuth(jwtSvc *services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := extractBearerToken(c)
		if !ok {
			utils.Unauthorized(c, "missing or invalid Authorization header")
			c.Abort()
			return
		}

		claims, err := jwtSvc.Parse(token)
		if err != nil {
			utils.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextEmail, claims.Email)
		c.Set(ContextIsAdmin, claims.IsAdmin)
		c.Next()
	}
}

// RequireAdmin 要求当前 token 为管理员身份（须在 JWTAuth 之后使用）
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		v, exists := c.Get(ContextIsAdmin)
		if !exists {
			utils.Unauthorized(c, "missing authentication")
			c.Abort()
			return
		}
		isAdmin, _ := v.(bool)
		if !isAdmin {
			utils.Forbidden(c, "admin only")
			c.Abort()
			return
		}
		c.Next()
	}
}

// extractBearerToken 从 Authorization 头中取出 Bearer token
func extractBearerToken(c *gin.Context) (string, bool) {
	auth := c.GetHeader(HeaderAuthorization)
	if auth == "" {
		return "", false
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(auth, prefix) {
		return "", false
	}
	token := strings.TrimSpace(auth[len(prefix):])
	if token == "" {
		return "", false
	}
	return token, true
}
