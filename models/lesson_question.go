package models

type LessonQuestion struct {
	ID           int    `json:"id" pg:"id,pk"`
	LessonID     int    `json:"lesson_id" pg:"lesson_id,notnull"`
	QuestionText string `json:"question_text" pg:"question_text,notnull"`
	Explanation  string `json:"explanation" pg:"explanation"`
	Order        int    `json:"order" pg:"order,notnull"`
}
