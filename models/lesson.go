package models

type Lesson struct {
	ID            int    `json:"id" pg:"id,pk"`
	Title         string `json:"title" pg:"title,notnull"`
	Description   string `json:"description" pg:"description"`
	TheoryContent string `json:"theory_content" pg:"theory_content,notnull"`
	Icon          string `json:"icon" pg:"icon"`
	Order         int    `json:"order" pg:"order"`
	IsActive      bool   `json:"is_active" pg:"is_active"`
	PassingScore  int    `json:"passing_score" pg:"passing_score,notnull"`
}
