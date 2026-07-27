package services

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
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
