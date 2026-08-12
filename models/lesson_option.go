package models

type LessonOption struct {
	ID         int    `json:"id" pg:"id,pk"`
	QuestionID int    `json:"question_id" pg:"question_id,notnull"`
	OptionText string `json:"option_text" pg:"option_text,notnull"`
	IsCorrect  bool   `json:"is_correct" pg:"is_correct,notnull"`
	Order      int    `json:"order" pg:"order,notnull"`
}
