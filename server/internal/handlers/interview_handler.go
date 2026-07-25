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

// InterviewHandler 模拟面试路由处理器
type InterviewHandler struct {
	svc *services.InterviewService
}

// NewInterviewHandler 构造 InterviewHandler
func NewInterviewHandler(svc *services.InterviewService) *InterviewHandler {
	return &InterviewHandler{svc: svc}
}

// Create POST /api/v1/interviews 创建面试会话
func (h *InterviewHandler) Create(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	var in services.CreateInterviewInput
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	interview, err := h.svc.Create(c.Request.Context(), userID, &in)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "interview created", interview)
}

// List GET /api/v1/interviews
func (h *InterviewHandler) List(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	list, err := h.svc.List(userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	if list == nil {
		list = []models.Interview{}
	}
	utils.OK(c, list)
}

// Get GET /api/v1/interviews/:id
func (h *InterviewHandler) Get(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	interview, err := h.svc.Get(userID, uint(id))
	if err != nil {
		if err == services.ErrInterviewNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	// 同时返回消息列表
	messages, _ := h.svc.ListMessages(userID, uint(id))
	if messages == nil {
		messages = []models.InterviewMessage{}
	}
	utils.OK(c, gin.H{
		"interview": interview,
		"messages":  messages,
	})
}

// SendMessage POST /api/v1/interviews/:id/messages 文字回答（SSE 流式）
func (h *InterviewHandler) SendMessage(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	interviewID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
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

	// 调用 AI 流式生成
	aiMsg, err := h.svc.SendMessage(c.Request.Context(), userID, uint(interviewID), req.Content, func(delta string) {
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

	// 检查面试是否已结束（答完所有题会自动结束）
	interview, _ := h.svc.Get(userID, uint(interviewID))
	if interview != nil && interview.Status == services.StatusCompleted {
		doneData, _ := json.Marshal(map[string]interface{}{
			"type":      "interview_ended",
			"message":   aiMsg,
			"interview": interview,
		})
		fmt.Fprintf(c.Writer, "data: %s\n\n", doneData)
	} else if aiMsg != nil {
		doneData, _ := json.Marshal(map[string]interface{}{
			"type":    "done",
			"message": aiMsg,
		})
		fmt.Fprintf(c.Writer, "data: %s\n\n", doneData)
	}
	flusher.Flush()
}

// SendVoice POST /api/v1/interviews/:id/voice 语音回答（multipart 上传）
func (h *InterviewHandler) SendVoice(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	interviewID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	file, header, err := c.Request.FormFile("audio")
	if err != nil {
		utils.BadRequest(c, "audio field is required: "+err.Error())
		return
	}
	defer file.Close()

	if !utils.IsAudioExt(header.Filename) {
		utils.BadRequest(c, "unsupported audio format, only mp3/wav/m4a/webm/ogg allowed")
		return
	}
	if !utils.ValidateFileSize(header.Size, 25) {
		utils.BadRequest(c, "audio size exceeds 25MB limit")
		return
	}

	// SSE 头（语音模式也用 SSE 推送 AI 回复）
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	flusher, ok := c.Writer.(interface{ Flush() })
	if !ok {
		utils.InternalError(c, "streaming not supported")
		return
	}

	// 推送转写进度
	transcribeData, _ := json.Marshal(map[string]string{"type": "status", "message": "正在转写语音..."})
	fmt.Fprintf(c.Writer, "data: %s\n\n", transcribeData)
	flusher.Flush()

	aiMsg, err := h.svc.SendVoice(c.Request.Context(), userID, uint(interviewID), file, header.Filename, func(delta string) {
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

	if aiMsg != nil {
		doneData, _ := json.Marshal(map[string]interface{}{
			"type":    "done",
			"message": aiMsg,
		})
		fmt.Fprintf(c.Writer, "data: %s\n\n", doneData)
	} else {
		// 面试已结束
		doneData, _ := json.Marshal(map[string]string{"type": "interview_ended"})
		fmt.Fprintf(c.Writer, "data: %s\n\n", doneData)
	}
	flusher.Flush()
}

// GetTTS GET /api/v1/interviews/:id/tts/:msgId 获取 AI 提问的 TTS 音频
func (h *InterviewHandler) GetTTS(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	interviewID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	msgID, _ := strconv.ParseUint(c.Param("msgId"), 10, 64)

	audio, filename, err := h.svc.GetTTS(c.Request.Context(), userID, uint(interviewID), uint(msgID))
	if err != nil {
		if err == services.ErrInterviewNotFound || err == services.ErrMessageNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	c.Header("Content-Type", "audio/mpeg")
	c.Header("Content-Disposition", "inline; filename=\""+filename+"\"")
	c.Header("Content-Length", strconv.Itoa(len(audio)))
	c.Writer.WriteHeader(200)
	c.Writer.Write(audio)
}

// End POST /api/v1/interviews/:id/end 结束面试并生成复盘
func (h *InterviewHandler) End(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	report, err := h.svc.End(c.Request.Context(), userID, uint(id))
	if err != nil {
		if err == services.ErrInterviewNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	utils.OKWithMsg(c, "interview ended, report generated", report)
}

// GetReport GET /api/v1/interviews/:id/report
func (h *InterviewHandler) GetReport(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	report, err := h.svc.GetReport(userID, uint(id))
	if err != nil {
		if err == services.ErrInterviewNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		if err == services.ErrReportNotReady {
			utils.NotFound(c, "report not generated yet, please end the interview first")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, report)
}

// GetScores GET /api/v1/interviews/:id/scores
func (h *InterviewHandler) GetScores(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	scores, err := h.svc.GetScores(userID, uint(id))
	if err != nil {
		if err == services.ErrInterviewNotFound {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}
	if scores == nil {
		scores = []models.InterviewScore{}
	}
	utils.OK(c, scores)
}
