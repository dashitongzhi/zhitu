package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhitu/server/internal/middleware"
	"github.com/zhitu/server/internal/models"
	"github.com/zhitu/server/internal/services"
	"github.com/zhitu/server/internal/utils"
)

// ResumeHandler 简历路由处理器
type ResumeHandler struct {
	svc *services.ResumeService
	ai  *services.ResumeAIService
}

// NewResumeHandler 构造 ResumeHandler
func NewResumeHandler(svc *services.ResumeService, ai *services.ResumeAIService) *ResumeHandler {
	return &ResumeHandler{svc: svc, ai: ai}
}

// ---------- 简历主表 CRUD ----------

// List GET /api/v1/resumes
func (h *ResumeHandler) List(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	list, err := h.svc.List(userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	if list == nil {
		list = []models.Resume{}
	}
	utils.OK(c, list)
}

// Create POST /api/v1/resumes
func (h *ResumeHandler) Create(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	var in services.CreateResumeInput
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	resume, err := h.svc.Create(userID, &in)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "created", resume)
}

// Get GET /api/v1/resumes/:id
func (h *ResumeHandler) Get(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid id")
		return
	}
	resume, err := h.svc.Get(userID, uint(id))
	if err != nil {
		if err == services.ErrResumeNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, resume)
}

// Update PUT /api/v1/resumes/:id
func (h *ResumeHandler) Update(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.Update(userID, uint(id), req); err != nil {
		if err == services.ErrResumeNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "updated", nil)
}

// Delete DELETE /api/v1/resumes/:id
func (h *ResumeHandler) Delete(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(userID, uint(id)); err != nil {
		if err == services.ErrResumeNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "deleted", nil)
}

// ---------- 版本管理 ----------

// ListVersions GET /api/v1/resumes/:id/versions
func (h *ResumeHandler) ListVersions(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	resumeID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	list, err := h.svc.ListVersions(userID, uint(resumeID))
	if err != nil {
		if err == services.ErrResumeNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, list)
}

// CreateVersion POST /api/v1/resumes/:id/versions
func (h *ResumeHandler) CreateVersion(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	resumeID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in services.CreateVersionInput
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	version, err := h.svc.CreateVersion(userID, uint(resumeID), &in)
	if err != nil {
		if err == services.ErrResumeNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "version created", version)
}

// GetVersion GET /api/v1/resumes/:id/versions/:vid
func (h *ResumeHandler) GetVersion(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	resumeID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	versionID, _ := strconv.ParseUint(c.Param("vid"), 10, 64)
	version, err := h.svc.GetVersion(userID, uint(resumeID), uint(versionID))
	if err != nil {
		if err == services.ErrResumeNotFound || err == services.ErrVersionNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, version)
}

// RollbackVersion POST /api/v1/resumes/:id/rollback/:vid
func (h *ResumeHandler) RollbackVersion(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	resumeID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	versionID, _ := strconv.ParseUint(c.Param("vid"), 10, 64)
	version, err := h.svc.RollbackVersion(userID, uint(resumeID), uint(versionID))
	if err != nil {
		if err == services.ErrResumeNotFound || err == services.ErrVersionNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "rolled back", version)
}

// ---------- AI 操作 ----------

// AIGenerate POST /api/v1/resumes/:id/ai/generate （SSE 流式）
func (h *ResumeHandler) AIGenerate(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	resumeID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var in services.GenerateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		// 允许空 body
		in = services.GenerateInput{}
	}

	// SSE 头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	flusher, ok := c.Writer.(interface{ Flush() })
	if !ok {
		utils.InternalError(c, "streaming not supported")
		return
	}

	// 调用 AI 生成，流式推送进度
	version, err := h.ai.Generate(c.Request.Context(), userID, uint(resumeID), &in, func(delta string) {
		data, _ := json.Marshal(map[string]string{"type": "delta", "content": delta})
		fmt.Fprintf(c.Writer, "data: %s\n\n", data)
		flusher.Flush()
	})
	if err != nil {
		errData, _ := json.Marshal(map[string]string{"type": "error", "message": err.Error()})
		fmt.Fprintf(c.Writer, "data: %s\n\n", errData)
		flusher.Flush()
		return
	}

	// 推送最终版本
	doneData, _ := json.Marshal(map[string]interface{}{
		"type":    "done",
		"version": version,
	})
	fmt.Fprintf(c.Writer, "data: %s\n\n", doneData)
	flusher.Flush()
}

// AIPolish POST /api/v1/resumes/:id/ai/polish
func (h *ResumeHandler) AIPolish(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	resumeID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in services.PolishInput
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	version, err := h.ai.Polish(c.Request.Context(), userID, uint(resumeID), &in)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "polished", version)
}

// AIScore POST /api/v1/resumes/:id/ai/score
func (h *ResumeHandler) AIScore(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	resumeID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		JD string `json:"jd"`
	}
	_ = c.ShouldBindJSON(&req)

	result, err := h.ai.Score(c.Request.Context(), userID, uint(resumeID), req.JD)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, result)
}

// AIJDMatch POST /api/v1/resumes/:id/ai/jd-match
func (h *ResumeHandler) AIJDMatch(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	resumeID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		JD string `json:"jd" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "jd is required")
		return
	}
	result, err := h.ai.JDMatch(c.Request.Context(), userID, uint(resumeID), req.JD)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, result)
}

// ---------- 同步档案 ----------

// SyncProfile POST /api/v1/resumes/:id/sync-profile
func (h *ResumeHandler) SyncProfile(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	resumeID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.SyncToProfile(userID, uint(resumeID)); err != nil {
		if err == services.ErrResumeNotFound || err == services.ErrVersionNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "synced to profile", nil)
}
