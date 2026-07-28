package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/zhitu/server/internal/config"
	"github.com/zhitu/server/internal/models"
	"gorm.io/gorm"
)

func TestCreateSupportsAllSceneHallScenesWithoutLLM(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:scene_hall_test?mode=memory&cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.AutoMigrate(&models.Interview{}, &models.InterviewMessage{}); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	service := NewInterviewService(db, nil, nil, nil)
	scenes := []string{
		SceneTeaching,
		SceneCorporate,
		SceneGroup,
		SceneDefense,
		SceneClient,
		ScenePressure,
		ScenePublic,
		SceneMedical,
		SceneMedia,
		SceneRemote,
		SceneSystem,
		SceneAviation,
	}

	for _, scene := range scenes {
		t.Run(scene, func(t *testing.T) {
			interview, err := service.Create(context.Background(), 1, &CreateInterviewInput{
				Scene:          scene,
				TargetPosition: "测试岗位",
				TotalQuestions: 5,
				Mode:           ModeHybrid,
			})
			if err != nil {
				t.Fatalf("create scene %q: %v", scene, err)
			}
			if interview.CurrentQuestionNo != 1 {
				t.Fatalf("current question = %d, want 1", interview.CurrentQuestionNo)
			}

			var firstQuestion models.InterviewMessage
			if err := db.Where("interview_id = ? AND question_no = 1", interview.ID).First(&firstQuestion).Error; err != nil {
				t.Fatalf("load first question: %v", err)
			}
			if firstQuestion.Content == "" {
				t.Fatal("first question is empty")
			}
		})
	}
}

func TestAttachResumeUsesOwnedCurrentVersion(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:attach_resume_test?mode=memory&cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.AutoMigrate(
		&models.Interview{},
		&models.Resume{},
		&models.ResumeVersion{},
	); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	interview := models.Interview{UserID: 1, Status: StatusOngoing}
	if err := db.Create(&interview).Error; err != nil {
		t.Fatalf("create interview: %v", err)
	}
	resume := models.Resume{UserID: 1, Name: "后端工程师简历", Scene: "manual"}
	if err := db.Create(&resume).Error; err != nil {
		t.Fatalf("create resume: %v", err)
	}
	content := `{"personal":{"name":"候选人"},"project":[{"name":"搜索平台","role":"负责人","description":"负责架构设计"}]}`
	version := models.ResumeVersion{
		ResumeID:     resume.ID,
		VersionLabel: "v1.0",
		Content:      content,
	}
	if err := db.Create(&version).Error; err != nil {
		t.Fatalf("create resume version: %v", err)
	}
	if err := db.Model(&resume).Update("current_version_id", version.ID).Error; err != nil {
		t.Fatalf("set current version: %v", err)
	}

	service := NewInterviewService(db, nil, nil, nil)
	got, err := service.AttachResume(1, interview.ID, &AttachResumeInput{ResumeID: resume.ID})
	if err != nil {
		t.Fatalf("AttachResume() error = %v", err)
	}
	if got.ResumeSnapshot != content || got.ResumeName != resume.Name {
		t.Fatalf("attached resume = %#v", got)
	}

	var persisted models.Interview
	if err := db.First(&persisted, interview.ID).Error; err != nil {
		t.Fatalf("reload interview: %v", err)
	}
	if persisted.ResumeSnapshot != content || persisted.ResumeName != resume.Name {
		t.Fatalf("persisted resume snapshot/name = %q/%q", persisted.ResumeSnapshot, persisted.ResumeName)
	}

	if _, err := service.AttachResume(2, interview.ID, &AttachResumeInput{ResumeID: resume.ID}); !errors.Is(err, ErrInterviewNotFound) {
		t.Fatalf("other user AttachResume() error = %v, want ErrInterviewNotFound", err)
	}
}

func TestSummarizeResumeTruncatesPromptContent(t *testing.T) {
	content := `{"custom":[{"title":"补充经历","content":"` +
		strings.Repeat("项", maxResumePromptRunes+100) +
		`"}]}`
	summary := summarizeResume(content)
	if !strings.HasSuffix(summary, "（简历内容已截断）") {
		t.Fatalf("summary was not truncated")
	}
	if got := len([]rune(summary)); got > maxResumePromptRunes+20 {
		t.Fatalf("summary rune count = %d", got)
	}
}

func TestTranscribeVoiceDoesNotAdvanceInterview(t *testing.T) {
	requests := make(chan map[string]interface{}, 1)
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		requests <- payload
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"这是转写后的回答"}}]}`))
	}))
	defer upstream.Close()

	db, err := gorm.Open(sqlite.Open("file:transcribe_voice_test?mode=memory&cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.AutoMigrate(&models.Interview{}, &models.InterviewMessage{}); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	interview := models.Interview{
		UserID:            1,
		Status:            StatusOngoing,
		CurrentQuestionNo: 2,
		TotalQuestions:    5,
	}
	if err := db.Create(&interview).Error; err != nil {
		t.Fatalf("create interview: %v", err)
	}

	llm := NewLLMService(&config.LLMConfig{
		Provider:     "mimo",
		BaseURL:      upstream.URL,
		APIKey:       "test-key",
		WhisperModel: "mimo-v2.5-asr",
	})
	service := NewInterviewService(db, llm, nil, nil)
	text, err := service.TranscribeVoice(
		context.Background(),
		1,
		interview.ID,
		bytes.NewReader([]byte("RIFF-test-wav")),
		"answer.wav",
	)
	if err != nil {
		t.Fatalf("TranscribeVoice() error = %v", err)
	}
	if text != "这是转写后的回答" {
		t.Fatalf("TranscribeVoice() = %q", text)
	}

	payload := <-requests
	if payload["model"] != "mimo-v2.5-asr" {
		t.Fatalf("model = %#v", payload["model"])
	}
	messages, ok := payload["messages"].([]interface{})
	if !ok || len(messages) != 1 {
		t.Fatalf("messages = %#v", payload["messages"])
	}
	messagePayload := messages[0].(map[string]interface{})
	content := messagePayload["content"].([]interface{})
	audioItem := content[0].(map[string]interface{})
	inputAudio := audioItem["input_audio"].(map[string]interface{})
	if data, _ := inputAudio["data"].(string); !strings.HasPrefix(data, "data:audio/wav;base64,") {
		t.Fatalf("audio data url = %q", data)
	}

	var count int64
	if err := db.Model(&models.InterviewMessage{}).Where("interview_id = ?", interview.ID).Count(&count).Error; err != nil {
		t.Fatalf("count messages: %v", err)
	}
	if count != 0 {
		t.Fatalf("message count = %d, want 0", count)
	}
	var persisted models.Interview
	if err := db.First(&persisted, interview.ID).Error; err != nil {
		t.Fatalf("reload interview: %v", err)
	}
	if persisted.CurrentQuestionNo != 2 {
		t.Fatalf("current question = %d, want 2", persisted.CurrentQuestionNo)
	}
}
