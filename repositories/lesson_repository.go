package repositories

import (
	"anti-scam-trainer/backend/models"
	"time"

	"github.com/go-pg/pg"
)

func ListLessons(db *pg.DB) ([]models.Lesson, error) {
	var lessons []models.Lesson
	err := db.Model(&lessons).Order("\"order\"").Select()
	return lessons, err
}

func GetLessonByID(db *pg.DB, id int) (*models.Lesson, error) {
	var lesson models.Lesson
	err := db.Model(&lesson).Where("id = ?", id).Select()
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

func SubmitLessonQuiz(db *pg.DB, userID, lessonID, score int, passed bool) error {
	userLesson := models.UserLesson{
		UserID:   userID,
		LessonID: lessonID,
	}

	err := db.Model(&userLesson).
		Where("user_id = ? AND lesson_id = ?", userID, lessonID).
		Select()
	if err != nil && err != pg.ErrNoRows {
		return err
	}

	if err == pg.ErrNoRows {
		userLesson = models.UserLesson{
			UserID:      userID,
			LessonID:    lessonID,
			QuizScore:   score,
			QuizPassed:  passed,
			CompletedAt: time.Now(),
		}
		_, err = db.Model(&userLesson).Insert()
		return err
	}

	userLesson.QuizScore = score
	userLesson.QuizPassed = passed
	if passed {
		userLesson.CompletedAt = time.Now()
	}
	_, err = db.Model(&userLesson).
		Column("quiz_score", "quiz_passed", "completed_at").
		Where("id = ?", userLesson.ID).
		Update()
	return err
}

func ListLessonChats(db *pg.DB, lessonID int, userID int) ([]models.Chat, error) {
	var chats []models.Chat
	err := db.Model(&chats).
		Join("JOIN lesson_chats lc ON lc.chat_id = chats.id").
		Where("lc.lesson_id = ?", lessonID).
		Order("lc.\"order\"").
		Select()
	return chats, err
}

func ListLessonChatsWithSessions(db *pg.DB, lessonID, userID int) ([]models.Chat, error) {
	return ListLessonChats(db, lessonID, userID)
}
