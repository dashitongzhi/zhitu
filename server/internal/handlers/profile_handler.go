package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhitu/server/internal/middleware"
	"github.com/zhitu/server/internal/services"
	"github.com/zhitu/server/internal/utils"
)

// ProfileHandler 用户档案路由处理器
type ProfileHandler struct {
	svc *services.ProfileService
}

// NewProfileHandler 构造 ProfileHandler
func NewProfileHandler(svc *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

// ---------- 主表 ----------

// GetProfile GET /api/v1/profile 获取完整档案
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	fp, err := h.svc.GetFullProfile(userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, fp)
}

// UpdateProfile PUT /api/v1/profile 更新基础信息
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid request body: "+err.Error())
		return
	}
	// 过滤允许更新的字段
	allowed := map[string]bool{
		"real_name": true, "gender": true, "birth_date": true, "phone": true,
		"target_position": true, "target_city": true, "expected_salary": true,
		"job_status": true, "self_introduction": true,
	}
	updates := map[string]interface{}{}
	for k, v := range req {
		if allowed[k] {
			updates[k] = v
		}
	}
	if len(updates) == 0 {
		utils.BadRequest(c, "no updatable fields provided")
		return
	}

	profile, err := h.svc.UpsertProfile(userID, updates)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, profile)
}

// GetCompletion GET /api/v1/profile/completion 获取完成度
func (h *ProfileHandler) GetCompletion(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	pct, err := h.svc.GetCompletion(userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, gin.H{"completion_pct": pct})
}

// ParseResume POST /api/v1/profile/parse-resume 上传简历解析
func (h *ProfileHandler) ParseResume(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "file field is required: "+err.Error())
		return
	}
	defer file.Close()

	if !utils.IsResumeExt(header.Filename) {
		utils.BadRequest(c, "unsupported resume format, only .pdf/.docx/.txt allowed")
		return
	}
	if !utils.ValidateFileSize(header.Size, 20) {
		utils.BadRequest(c, "file size exceeds 20MB limit")
		return
	}

	parsed, err := h.svc.ParseResumeFile(c.Request.Context(), userID, file, header.Filename)
	if err != nil {
		utils.InternalError(c, "parse resume failed: "+err.Error())
		return
	}
	utils.OKWithMsg(c, "resume parsed and merged into profile", parsed)
}

// ---------- 子资源通用 CRUD ----------

// 通用：列表 + 创建 + 更新 + 删除，通过子资源类型分发
// 为减少重复代码，用类型参数化；但 Go 1.22 泛型对 gorm 处理较繁，这里按子资源分别实现

// Educations
func (h *ProfileHandler) ListEducations(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	list, err := h.svc.ListEducations(userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, list)
}

func (h *ProfileHandler) CreateEducation(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	var m services.EducationDTO
	if err := c.ShouldBindJSON(&m); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	created, err := h.svc.CreateEducation(userID, m.ToModel())
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "created", created)
}

func (h *ProfileHandler) UpdateEducation(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid id")
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateEducation(userID, uint(id), updates); err != nil {
		if err == services.ErrSubResourceNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "updated", nil)
}

func (h *ProfileHandler) DeleteEducation(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid id")
		return
	}
	if err := h.svc.DeleteEducation(userID, uint(id)); err != nil {
		if err == services.ErrSubResourceNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "deleted", nil)
}

// Works
func (h *ProfileHandler) ListWorks(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	list, err := h.svc.ListWorks(userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, list)
}

func (h *ProfileHandler) CreateWork(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	var m services.WorkDTO
	if err := c.ShouldBindJSON(&m); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	created, err := h.svc.CreateWork(userID, m.ToModel())
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "created", created)
}

func (h *ProfileHandler) UpdateWork(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateWork(userID, uint(id), updates); err != nil {
		if err == services.ErrSubResourceNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "updated", nil)
}

func (h *ProfileHandler) DeleteWork(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteWork(userID, uint(id)); err != nil {
		if err == services.ErrSubResourceNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "deleted", nil)
}

// Projects
func (h *ProfileHandler) ListProjects(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	list, err := h.svc.ListProjects(userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, list)
}

func (h *ProfileHandler) CreateProject(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	var m services.ProjectDTO
	if err := c.ShouldBindJSON(&m); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	created, err := h.svc.CreateProject(userID, m.ToModel())
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "created", created)
}

func (h *ProfileHandler) UpdateProject(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateProject(userID, uint(id), updates); err != nil {
		if err == services.ErrSubResourceNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "updated", nil)
}

func (h *ProfileHandler) DeleteProject(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteProject(userID, uint(id)); err != nil {
		if err == services.ErrSubResourceNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "deleted", nil)
}

// Skills
func (h *ProfileHandler) ListSkills(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	list, err := h.svc.ListSkills(userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, list)
}

func (h *ProfileHandler) CreateSkill(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	var m services.SkillDTO
	if err := c.ShouldBindJSON(&m); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	created, err := h.svc.CreateSkill(userID, m.ToModel())
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "created", created)
}

func (h *ProfileHandler) UpdateSkill(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateSkill(userID, uint(id), updates); err != nil {
		if err == services.ErrSubResourceNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "updated", nil)
}

func (h *ProfileHandler) DeleteSkill(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteSkill(userID, uint(id)); err != nil {
		if err == services.ErrSubResourceNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "deleted", nil)
}

// Honors
func (h *ProfileHandler) ListHonors(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	list, err := h.svc.ListHonors(userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, list)
}

func (h *ProfileHandler) CreateHonor(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	var m services.HonorDTO
	if err := c.ShouldBindJSON(&m); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	created, err := h.svc.CreateHonor(userID, m.ToModel())
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "created", created)
}

func (h *ProfileHandler) UpdateHonor(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateHonor(userID, uint(id), updates); err != nil {
		if err == services.ErrSubResourceNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "updated", nil)
}

func (h *ProfileHandler) DeleteHonor(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteHonor(userID, uint(id)); err != nil {
		if err == services.ErrSubResourceNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "deleted", nil)
}

// Practices
func (h *ProfileHandler) ListPractices(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	list, err := h.svc.ListPractices(userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, list)
}

func (h *ProfileHandler) CreatePractice(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	var m services.PracticeDTO
	if err := c.ShouldBindJSON(&m); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	created, err := h.svc.CreatePractice(userID, m.ToModel())
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "created", created)
}

func (h *ProfileHandler) UpdatePractice(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdatePractice(userID, uint(id), updates); err != nil {
		if err == services.ErrSubResourceNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "updated", nil)
}

func (h *ProfileHandler) DeletePractice(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeletePractice(userID, uint(id)); err != nil {
		if err == services.ErrSubResourceNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "deleted", nil)
}
