package services

import (
	"anti-scam-trainer/backend/models"
	"anti-scam-trainer/backend/repositories"

	"github.com/go-pg/pg"
)

func ListLessons(db *pg.DB) ([]models.Lesson, error) {
	return repositories.ListLessons(db)
}

func GetLessonByID(db *pg.DB, id int) (*models.Lesson, error) {
	return repositories.GetLessonByID(db, id)
}

func SubmitLessonQuiz(db *pg.DB, userID, lessonID, score int, passed bool) error {
	return repositories.SubmitLessonQuiz(db, userID, lessonID, score, passed)
}

func ListLessonChats(db *pg.DB, lessonID, userID int) ([]models.Chat, error) {
	return repositories.ListLessonChats(db, lessonID, userID)
}
