package tests

import (
	"context"
	"platform/internal/database"
)

func GetTestByModule(moduleID int64) (*Test, error) {
	row := database.DB.QueryRow(context.Background(),
		`SELECT id, module_id, title, time_limit, max_attempts, created_at
		 FROM tests WHERE module_id=$1`, moduleID)

	var t Test
	err := row.Scan(&t.ID, &t.ModuleID, &t.Title, &t.TimeLimit, &t.MaxAttempts, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func GetQuestions(testID int64) ([]Question, error) {
	rows, _ := database.DB.Query(context.Background(),
		`SELECT id, test_id, question_text, order_index
		 FROM questions WHERE test_id=$1 ORDER BY order_index`, testID)

	defer rows.Close()

	q := make([]Question, 0)

	for rows.Next() {
		var item Question
		rows.Scan(&item.ID, &item.TestID, &item.QuestionText, &item.OrderIndex)
		q = append(q, item)
	}
	return q, nil
}

func GetAnswers(questionID int64) ([]Answer, error) {

	rows, _ := database.DB.Query(context.Background(),
		`SELECT id, question_id, text, is_correct
		 FROM answers WHERE question_id=$1`, questionID)

	defer rows.Close()

	ans := make([]Answer, 0)

	for rows.Next() {
		var a Answer
		rows.Scan(&a.ID, &a.QuestionID, &a.Text, &a.IsCorrect)
		ans = append(ans, a)
	}
	return ans, nil
}

func CreateAttempt(userID, testID int64) (*Attempt, error) {

	row := database.DB.QueryRow(context.Background(),
		`INSERT INTO attempts (user_id, test_id, started_at)
		 VALUES ($1,$2,NOW())
		 RETURNING id, started_at`, userID, testID)

	var a Attempt
	row.Scan(&a.ID, &a.StartedAt)

	return &a, nil
}
