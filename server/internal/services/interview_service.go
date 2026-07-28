package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/zhitu/server/internal/config"
	"github.com/zhitu/server/internal/models"
	"gorm.io/gorm"
	"time"
)

// 面试相关错误
var (
	ErrInterviewNotFound = errors.New("interview not found")
	ErrInterviewEnded    = errors.New("interview already ended")
	ErrMessageNotFound   = errors.New("interview message not found")
	ErrReportNotReady    = errors.New("report not generated yet, please end the interview first")
)

// 面试场景枚举
const (
	SceneTech      = "tech"
	SceneBehavior  = "behavior"
	ScenePressure  = "pressure"
	SceneHR        = "hr"
	SceneGroup     = "group"
	SceneTeaching  = "teaching"
	SceneCorporate = "corporate"
	SceneDefense   = "defense"
	SceneClient    = "client"
	ScenePublic    = "public"
	SceneMedical   = "medical"
	SceneMedia     = "media"
	SceneRemote    = "remote"
	SceneSystem    = "system"
	SceneAviation  = "aviation"
)

var validInterviewScenes = map[string]struct{}{
	SceneTech: {}, SceneBehavior: {}, ScenePressure: {}, SceneHR: {},
	SceneGroup: {}, SceneTeaching: {}, SceneCorporate: {}, SceneDefense: {},
	SceneClient: {}, ScenePublic: {}, SceneMedical: {}, SceneMedia: {},
	SceneRemote: {}, SceneSystem: {}, SceneAviation: {},
}

// 面试状态枚举
const (
	StatusOngoing   = "ongoing"
	StatusCompleted = "completed"
	StatusCancelled = "cancelled"
)

// 面试模式枚举
const (
	ModeText   = "text"
	ModeVoice  = "voice"
	ModeHybrid = "hybrid"
)

const (
	maxResumePromptRunes   = 12000
	maxFollowupPromptRunes = 2000
)

// InterviewService 面试业务逻辑
type InterviewService struct {
	db      *gorm.DB
	llm     *LLMService
	profile *ProfileService
	storage *config.StorageConfig
}

// NewInterviewService 构造 InterviewService
func NewInterviewService(db *gorm.DB, llm *LLMService, profile *ProfileService, storage *config.StorageConfig) *InterviewService {
	return &InterviewService{db: db, llm: llm, profile: profile, storage: storage}
}

// CreateInterviewInput 创建面试入参
type CreateInterviewInput struct {
	Scene          string `json:"scene" binding:"required"`
	TargetCompany  string `json:"target_company"`
	TargetPosition string `json:"target_position" binding:"required"`
	TargetJD       string `json:"target_jd"`
	Difficulty     string `json:"difficulty"`
	TotalQuestions int    `json:"total_questions"`
	Mode           string `json:"mode"`
}

// AttachResumeInput 面试发送简历入参
type AttachResumeInput struct {
	ResumeID  uint `json:"resume_id" binding:"required"`
	VersionID uint `json:"version_id"` // 可选，不传则用简历当前版本
}

// Create 创建面试会话，并写入一道无需大模型的开场题。
// 大模型只在用户开始作答后参与后续追问，不阻塞进入面试房间。
func (s *InterviewService) Create(ctx context.Context, userID uint, in *CreateInterviewInput) (*models.Interview, error) {
	if _, ok := validInterviewScenes[in.Scene]; !ok {
		return nil, errors.New("invalid interview scene")
	}
	if in.Difficulty == "" {
		in.Difficulty = "mid"
	}
	if in.TotalQuestions == 0 {
		in.TotalQuestions = 8
	}
	if in.TotalQuestions < 5 || in.TotalQuestions > 15 {
		return nil, errors.New("total_questions must be between 5 and 15")
	}
	if in.Mode == "" {
		in.Mode = ModeText
	}

	now := time.Now()
	interview := &models.Interview{
		UserID:            userID,
		Scene:             in.Scene,
		TargetCompany:     in.TargetCompany,
		TargetPosition:    in.TargetPosition,
		TargetJD:          in.TargetJD,
		Difficulty:        in.Difficulty,
		TotalQuestions:    in.TotalQuestions,
		Mode:              in.Mode,
		Status:            StatusOngoing,
		CurrentQuestionNo: 1,
		StartedAt:         &now,
	}
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(interview).Error; err != nil {
			return err
		}
		firstQuestion := &models.InterviewMessage{
			InterviewID:  interview.ID,
			Role:         "assistant",
			Content:      openingQuestion(interview),
			QuestionType: "opening",
			QuestionNo:   1,
		}
		return tx.Create(firstQuestion).Error
	}); err != nil {
		return nil, err
	}

	return interview, nil
}

func openingQuestion(interview *models.Interview) string {
	questions := map[string]string{
		SceneTeaching:  "欢迎进入模拟教室。请面向考官做一段简短的自我介绍，并说明你对%s课堂教学的理解。",
		SceneCorporate: "欢迎进入企业会议室。请简要介绍自己，并选择一段最能证明你胜任%s的经历。",
		SceneGroup:     "欢迎进入群面讨论室。请用一分钟陈述你对讨论目标的理解，以及你准备承担的团队角色。",
		SceneDefense:   "欢迎进入项目答辩室。请先概述与你申请的%s最相关的项目背景、职责和成果。",
		SceneClient:    "欢迎进入客户会议室。假设客户首次与你见面，请围绕%s做一段简洁、有说服力的价值介绍。",
		ScenePressure:  "欢迎进入压力面试室。请直面回答：与其他候选人相比，我们为什么应该选择你担任%s？",
		ScenePublic:    "欢迎进入结构化面试厅。请结合%s的职责，谈谈你如何看待服务意识与执行能力。",
		SceneMedical:   "欢迎进入医疗面试室。请结合%s岗位，说明你如何兼顾专业判断、患者感受与沟通效率。",
		SceneMedia:     "欢迎进入媒体演播室。请面向镜头完成一分钟自我介绍，并说明你应聘%s的核心优势。",
		SceneRemote:    "欢迎进入远程面试间。请用清晰简洁的方式介绍自己，并说明你胜任%s和远程协作的优势。",
		SceneSystem:    "欢迎进入系统设计室。请先说明你在%s相关工作中如何进行需求澄清与技术方案取舍。",
		SceneAviation:  "欢迎进入航空面试厅。请结合%s岗位，介绍一次你主动提供优质服务或处理突发情况的经历。",
	}
	if question, ok := questions[interview.Scene]; ok {
		return fmt.Sprintf(question, interview.TargetPosition)
	}
	return fmt.Sprintf("欢迎参加本次面试。请先做一段简短的自我介绍，并说明你与%s岗位的匹配之处。", interview.TargetPosition)
}

// Get 获取面试详情（含所有消息）
func (s *InterviewService) Get(userID, id uint) (*models.Interview, error) {
	var interview models.Interview
	err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&interview).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInterviewNotFound
	}
	return &interview, err
}

// List 列出用户的所有面试
func (s *InterviewService) List(userID uint) ([]models.Interview, error) {
	var list []models.Interview
	err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&list).Error
	return list, err
}

// ListMessages 列出面试的所有消息
func (s *InterviewService) ListMessages(userID, interviewID uint) ([]models.InterviewMessage, error) {
	if _, err := s.Get(userID, interviewID); err != nil {
		return nil, err
	}
	var list []models.InterviewMessage
	err := s.db.Where("interview_id = ?", interviewID).Order("id ASC").Find(&list).Error
	return list, err
}
