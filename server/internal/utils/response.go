// Package utils 提供跨模块复用的工具函数
package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 标准业务错误码
const (
	CodeSuccess        = 0
	CodeInvalidParams  = 40001
	CodeUnauthorized   = 40101
	CodeForbidden      = 40301
	CodeNotFound       = 40401
	CodeConflict       = 40901
	CodeInternalError  = 50001
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// OK 返回成功响应
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// OKWithMsg 返回带自定义消息的成功响应
func OKWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: msg,
		Data:    data,
	})
}

// Fail 返回失败响应（HTTP 状态码由调用方决定）
func Fail(c *gin.Context, httpCode, bizCode int, msg string) {
	c.JSON(httpCode, Response{
		Code:    bizCode,
		Message: msg,
	})
}

// FailWithData 返回带 data 的失败响应（如校验错误字段）
func FailWithData(c *gin.Context, httpCode, bizCode int, msg string, data interface{}) {
	c.JSON(httpCode, Response{
		Code:    bizCode,
		Message: msg,
		Data:    data,
	})
}

// BadRequest 400 参数错误
func BadRequest(c *gin.Context, msg string) {
	Fail(c, http.StatusBadRequest, CodeInvalidParams, msg)
}

// Unauthorized 401 未认证
func Unauthorized(c *gin.Context, msg string) {
	Fail(c, http.StatusUnauthorized, CodeUnauthorized, msg)
}

// Forbidden 403 无权限
func Forbidden(c *gin.Context, msg string) {
	Fail(c, http.StatusForbidden, CodeForbidden, msg)
}

// NotFound 404 资源不存在
func NotFound(c *gin.Context, msg string) {
	Fail(c, http.StatusNotFound, CodeNotFound, msg)
}

// Conflict 409 冲突（如邮箱已存在）
func Conflict(c *gin.Context, msg string) {
	Fail(c, http.StatusConflict, CodeConflict, msg)
}

// InternalError 500 内部错误
func InternalError(c *gin.Context, msg string) {
	Fail(c, http.StatusInternalServerError, CodeInternalError, msg)
}
