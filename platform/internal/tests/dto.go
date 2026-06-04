package tests

type CreateTestRequest struct {
	ModuleID    *int64           `json:"module_id"`
	ChapterID   *int64           `json:"chapter_id"`
	CourseID    *int64           `json:"course_id"`
	Type        string           `json:"type"` // module | chapter | final
	Title       string           `json:"title"`
	TimeLimit   int              `json:"time_limit"`
	MaxAttempts int              `json:"max_attempts"`
	Questions   []CreateQuestion `json:"questions"`
}

type CreateQuestion struct {
	Type         string         `json:"type"` // mcq | open | matching
	QuestionText string         `json:"question_text"`
	OrderIndex   int            `json:"order_index"`
	Answers      []CreateAnswer `json:"answers,omitempty"`
	CorrectText  string         `json:"correct_text,omitempty"`
	Pairs        []MatchingPair `json:"pairs,omitempty"`
}

type MatchingPair struct {
	Left  string `json:"left"`
	Right string `json:"right"`
}

type CreateAnswer struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type TeacherAnswer struct {
	ID         int64  `json:"id"`
	QuestionID int64  `json:"question_id"`
	Text       string `json:"text"`
	IsCorrect  bool   `json:"is_correct"`
}

