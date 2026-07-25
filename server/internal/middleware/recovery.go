package middleware

import (
	"log"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/zhitu/server/internal/utils"
)

// Recovery 捕获 panic 并返回 500，避免进程崩溃
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[RECOVERY] panic: %v\n%s", r, debug.Stack())
				utils.InternalError(c, "internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}
