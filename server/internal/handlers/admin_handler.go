package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhitu/server/internal/services"
	"github.com/zhitu/server/internal/utils"
)

// AdminHandler 管理端 HTTP 处理器
type AdminHandler struct {
	adminSvc *services.AdminService
}

// NewAdminHandler 构造函数
func NewAdminHandler(adminSvc *services.AdminService) *AdminHandler {
	return &AdminHandler{adminSvc: adminSvc}
}

// GetStats GET /api/admin/stats 仪表盘统计
func (h *AdminHandler) GetStats(c *gin.Context) {
	stats, err := h.adminSvc.GetDashboardStats()
	if err != nil {
		utils.Fail(c, 500, utils.CodeInternalError, "获取统计数据失败")
		return
	}
	utils.OK(c, stats)
}

// GetFunnel GET /api/admin/deliveries/funnel 全局投递漏斗
func (h *AdminHandler) GetFunnel(c *gin.Context) {
	funnel, err := h.adminSvc.GetDeliveryFunnel()
	if err != nil {
		utils.Fail(c, 500, utils.CodeInternalError, "获取漏斗数据失败")
		return
	}
	utils.OK(c, funnel)
}

// ListUsers GET /api/admin/users 用户列表
func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	result, err := h.adminSvc.ListUsers(page, pageSize, keyword, status)
	if err != nil {
		utils.Fail(c, 500, utils.CodeInternalError, "获取用户列表失败")
		return
	}
	utils.OK(c, result)
}

// GetUser GET /api/admin/users/:id 用户详情
func (h *AdminHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Fail(c, 400, utils.CodeInvalidParams, "无效的用户ID")
		return
	}

	detail, err := h.adminSvc.GetUserDetail(uint(id))
	if err != nil {
		utils.Fail(c, 404, utils.CodeNotFound, err.Error())
		return
	}
	utils.OK(c, detail)
}

// ToggleUserStatus PATCH /api/admin/users/:id/status 切换用户状态
func (h *AdminHandler) ToggleUserStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Fail(c, 400, utils.CodeInvalidParams, "无效的用户ID")
		return
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Fail(c, 400, utils.CodeInvalidParams, "参数错误")
		return
	}

	if err := h.adminSvc.ToggleUserStatus(uint(id), body.Status); err != nil {
		utils.Fail(c, 400, utils.CodeInvalidParams, err.Error())
		return
	}
	utils.OK(c, nil)
}

// ResetPassword POST /api/admin/users/:id/reset-password 重置用户密码
func (h *AdminHandler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Fail(c, 400, utils.CodeInvalidParams, "无效的用户ID")
		return
	}

	var body struct {
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Fail(c, 400, utils.CodeInvalidParams, "参数错误")
		return
	}

	if err := h.adminSvc.ResetUserPassword(uint(id), body.NewPassword); err != nil {
		utils.Fail(c, 400, utils.CodeInvalidParams, err.Error())
		return
	}
	utils.OK(c, nil)
}

// ListDeliveries GET /api/admin/deliveries 全局投递列表
func (h *AdminHandler) ListDeliveries(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	company := c.Query("company")
	userEmail := c.Query("user_email")

	result, err := h.adminSvc.ListDeliveries(page, pageSize, status, company, userEmail)
	if err != nil {
		utils.Fail(c, 500, utils.CodeInternalError, "获取投递列表失败")
		return
	}
	utils.OK(c, result)
}
