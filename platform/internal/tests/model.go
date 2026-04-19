package tests

import "time"

type Test struct {
	ID          int64     `json:"id"`
	ModuleID    int64     `json:"module_id"`
	ChapterID   *int64    `json:"chapter_id"`
	CourseID    *int64    `json:"course_id"`
	Type        string    `json:"type"` // module | chapter | final
	Title       string    `json:"title"`
	TimeLimit   int       `json:"time_limit"`
	MaxAttempts int       `json:"max_attempts"`
	CreatedAt   time.Time `json:"created_at"`
}

type Question struct {
	ID           int64   `json:"id"`
	TestID       int64   `json:"test_id"`
	Type         string  `json:"type"` // mcq | open | matching
	QuestionText string  `json:"question_text"`
	OrderIndex   int     `json:"order_index"`
	CorrectText  *string `json:"correct_text,omitempty"`

	Pairs   []map[string]string `json:"pairs,omitempty"`
	Options []string            `json:"options,omitempty"`
}

type Answer struct {
	ID         int64  `json:"id"`
	QuestionID int64  `json:"question_id"`
	Text       string `json:"text"`
	IsCorrect  bool   `json:"-"`
}

type Attempt struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	TestID     int64     `json:"test_id"`
	Score      int       `json:"score"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
}
