package handlers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhitu/server/internal/models"
	"github.com/zhitu/server/internal/services"
	"github.com/zhitu/server/internal/utils"
)

// DeliveryHandler 投递看板 HTTP 处理
type DeliveryHandler struct {
	svc *services.DeliveryService
}

// NewDeliveryHandler 构造 DeliveryHandler
func NewDeliveryHandler(svc *services.DeliveryService) *DeliveryHandler {
	return &DeliveryHandler{svc: svc}
}

// List GET /api/v1/deliveries?status=&channel=
func (h *DeliveryHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	status := c.Query("status")
	channel := c.Query("channel")
	list, err := h.svc.List(userID, status, channel)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, list)
}

// Get GET /api/v1/deliveries/:id  （含轮次与反馈）
func (h *DeliveryHandler) Get(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := parseUintParam(c, "id")
	d, rounds, feedbacks, err := h.svc.GetDetail(userID, id)
	if err != nil {
		respondDeliveryErr(c, err)
		return
	}
	utils.OK(c, gin.H{
		"delivery":  d,
		"rounds":    rounds,
		"feedbacks": feedbacks,
	})
}

// Create POST /api/v1/deliveries
func (h *DeliveryHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	var in services.CreateDeliveryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	d, err := h.svc.Create(userID, &in)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.OK(c, d)
}

// Update PUT /api/v1/deliveries/:id
func (h *DeliveryHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := parseUintParam(c, "id")
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.Update(userID, id, updates); err != nil {
		respondDeliveryErr(c, err)
		return
	}
	utils.OK(c, gin.H{"id": id})
}

// Delete DELETE /api/v1/deliveries/:id
func (h *DeliveryHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := parseUintParam(c, "id")
	if err := h.svc.Delete(userID, id); err != nil {
		respondDeliveryErr(c, err)
		return
	}
	utils.OK(c, gin.H{"id": id})
}

// ChangeStatus PATCH /api/v1/deliveries/:id/status
// body: {"status": "interview"}
func (h *DeliveryHandler) ChangeStatus(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := parseUintParam(c, "id")
	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	d, err := h.svc.ChangeStatus(userID, id, body.Status)
	if err != nil {
		if errors.Is(err, services.ErrInvalidStatus) || errors.Is(err, services.ErrInvalidTransition) {
			utils.BadRequest(c, err.Error())
			return
		}
		respondDeliveryErr(c, err)
		return
	}
	utils.OK(c, d)
}

// ---------- 面试轮次 ----------

// ListRounds GET /api/v1/deliveries/:id/rounds
func (h *DeliveryHandler) ListRounds(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := parseUintParam(c, "id")
	list, err := h.svc.ListRounds(userID, id)
	if err != nil {
		respondDeliveryErr(c, err)
		return
	}
	utils.OK(c, list)
}

// CreateRound POST /api/v1/deliveries/:id/rounds
func (h *DeliveryHandler) CreateRound(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := parseUintParam(c, "id")
	var r models.DeliveryRound
	if err := c.ShouldBindJSON(&r); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	round, err := h.svc.CreateRound(userID, id, &r)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.OK(c, round)
}

// UpdateRound PUT /api/v1/deliveries/:id/rounds/:rid
func (h *DeliveryHandler) UpdateRound(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := parseUintParam(c, "id")
	rid := parseUintParam(c, "rid")
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateRound(userID, id, rid, updates); err != nil {
		respondDeliveryErr(c, err)
		return
	}
	utils.OK(c, gin.H{"id": rid})
}

// DeleteRound DELETE /api/v1/deliveries/:id/rounds/:rid
func (h *DeliveryHandler) DeleteRound(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := parseUintParam(c, "id")
	rid := parseUintParam(c, "rid")
	if err := h.svc.DeleteRound(userID, id, rid); err != nil {
		respondDeliveryErr(c, err)
		return
	}
	utils.OK(c, gin.H{"id": rid})
}

// ---------- HR 反馈 ----------

// ListFeedbacks GET /api/v1/deliveries/:id/feedbacks
func (h *DeliveryHandler) ListFeedbacks(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := parseUintParam(c, "id")
	list, err := h.svc.ListFeedbacks(userID, id)
	if err != nil {
		respondDeliveryErr(c, err)
		return
	}
	utils.OK(c, list)
}

// CreateFeedback POST /api/v1/deliveries/:id/feedbacks
func (h *DeliveryHandler) CreateFeedback(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := parseUintParam(c, "id")
	var f models.DeliveryFeedback
	if err := c.ShouldBindJSON(&f); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	fb, err := h.svc.CreateFeedback(userID, id, &f)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.OK(c, fb)
}

// DeleteFeedback DELETE /api/v1/deliveries/:id/feedbacks/:fid
func (h *DeliveryHandler) DeleteFeedback(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := parseUintParam(c, "id")
	fid := parseUintParam(c, "fid")
	if err := h.svc.DeleteFeedback(userID, id, fid); err != nil {
		respondDeliveryErr(c, err)
		return
	}
	utils.OK(c, gin.H{"id": fid})
}

// ---------- 统计与漏斗 ----------

// GetStats GET /api/v1/deliveries/stats
func (h *DeliveryHandler) GetStats(c *gin.Context) {
	userID := c.GetUint("user_id")
	st, err := h.svc.GetStats(userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, st)
}

// GetFunnel GET /api/v1/deliveries/funnel
func (h *DeliveryHandler) GetFunnel(c *gin.Context) {
	userID := c.GetUint("user_id")
	f, err := h.svc.GetFunnel(userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, f)
}

// ---------- 辅助 ----------

func parseUintParam(c *gin.Context, key string) uint {
	v, _ := strconv.ParseUint(c.Param(key), 10, 64)
	return uint(v)
}

func respondDeliveryErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrDeliveryNotFound):
		utils.NotFound(c, "delivery not found")
	case errors.Is(err, services.ErrRoundNotFound):
		utils.NotFound(c, "delivery round not found")
	case errors.Is(err, services.ErrFeedbackNotFound):
		utils.NotFound(c, "delivery feedback not found")
	default:
		utils.InternalError(c, err.Error())
	}
}
