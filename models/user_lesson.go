package models

import "time"

type UserLesson struct {
	ID             int       `json:"id" pg:"id,pk"`
	UserID         int       `json:"user_id" pg:"user_id,notnull"`
	LessonID       int       `json:"lesson_id" pg:"lesson_id,notnull"`
	TheoryRead     bool      `json:"theory_read" pg:"theory_read"`
	QuizScore      int       `json:"quiz_score" pg:"quiz_score"`
	QuizPassed     bool      `json:"quiz_passed" pg:"quiz_passed"`
	ChatsCompleted int       `json:"chats_completed" pg:"chats_completed"`
	StartedAt      time.Time `json:"started_at" pg:"started_at"`
	CompletedAt    time.Time `json:"completed_at" pg:"completed_at"`
}
