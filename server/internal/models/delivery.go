package models

import "time"

// 投递状态枚举
// pending(待响应) → written_test(笔试中) → interview(面试中)
// → waiting_offer(待Offer) → offer(已Offer) / rejected(已拒绝)
const (
	DeliveryStatusPending      = "pending"
	DeliveryStatusWrittenTest  = "written_test"
	DeliveryStatusInterview    = "interview"
	DeliveryStatusWaitingOffer = "waiting_offer"
	DeliveryStatusOffer        = "offer"
	DeliveryStatusRejected     = "rejected"
)

// 投递渠道枚举
const (
	ChannelBoss     = "boss"
	ChannelOfficial = "official"
	ChannelReferral = "referral"
	ChannelCampus   = "campus"
	ChannelHeadhunt = "headhunt"
	ChannelOther    = "other"
)

// 优先级枚举
const (
	PriorityHigh   = "high"
	PriorityMedium = "medium"
	PriorityLow    = "low"
)

// 面试轮次类型
const (
	RoundWrittenTest = "written_test"
	RoundFirstTech   = "first_tech"
	RoundSecondTech  = "second_tech"
	RoundThirdTech   = "third_tech"
	RoundCross       = "cross"
	RoundHR          = "hr"
	RoundAdditional  = "additional"
	RoundFinal       = "final"
)

// 轮次结果
const (
	RoundResultPass    = "pass"
	RoundResultPending = "pending"
	RoundResultReject  = "rejected"
)

// 面试形式
const (
	FormatOnsite = "onsite"
	FormatVideo  = "video"
	FormatPhone  = "phone"
)

// Delivery 投递记录主表
type Delivery struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"index;not null" json:"user_id"`
	Company        string    `gorm:"size:100;not null" json:"company"`
	Position       string    `gorm:"size:100;not null" json:"position"`
	Channel        string    `gorm:"size:20;not null" json:"channel"` // boss/official/referral/campus/headhunt/other
	ApplyDate      string    `gorm:"size:10;not null" json:"apply_date"` // YYYY-MM-DD
	Status         string    `gorm:"size:20;default:pending" json:"status"`
	Priority       string    `gorm:"size:10;default:medium" json:"priority"` // high/medium/low
	JDText         string    `gorm:"type:text" json:"jd_text"`
	ResumeVerID    uint      `gorm:"default:0" json:"resume_version_id"` // 关联简历版本快照
	HRContact      string    `gorm:"type:text" json:"hr_contact"`        // JSON: {name, wechat, phone, email}
	NextStep       string    `gorm:"type:text" json:"next_step"`         // JSON: {todo, deadline}
	OfferDetail    string    `gorm:"type:text" json:"offer_detail"`      // JSON: {salary_base, annual_bonus, stock, benefits, deadline}
	Remark         string    `gorm:"type:text" json:"remark"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (Delivery) TableName() string { return "deliveries" }

// DeliveryRound 面试轮次记录
type DeliveryRound struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	DeliveryID      uint       `gorm:"index;not null" json:"delivery_id"`
	RoundType       string     `gorm:"size:20;not null" json:"round_type"` // written_test/first_tech/second_tech/third_tech/cross/hr/additional/final
	InterviewTime   *time.Time `json:"interview_time,omitempty"`
	Format          string     `gorm:"size:10" json:"format"` // onsite/video/phone
	InterviewerName string     `gorm:"size:50" json:"interviewer_name"`
	InterviewerTitle string    `gorm:"size:50" json:"interviewer_title"`
	QuestionSummary string     `gorm:"type:text" json:"question_summary"`
	Feedback        string     `gorm:"type:text" json:"feedback"`
	Result          string     `gorm:"size:10;default:pending" json:"result"` // pass/pending/rejected
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (DeliveryRound) TableName() string { return "delivery_rounds" }

// DeliveryFeedback HR反馈记录
type DeliveryFeedback struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	DeliveryID   uint      `gorm:"index;not null" json:"delivery_id"`
	ContactTime  string    `gorm:"size:20;not null" json:"contact_time"` // YYYY-MM-DD HH:MM
	Method       string    `gorm:"size:10" json:"method"` // wechat/phone/email
	Summary      string    `gorm:"type:text" json:"summary"`
	NextAction   string    `gorm:"type:text" json:"next_action"`
	CreatedAt    time.Time `json:"created_at"`
}

func (DeliveryFeedback) TableName() string { return "delivery_feedbacks" }
