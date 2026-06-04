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

func GetTestByChapter(chapterID int64) (*Test, error) {
	row := database.DB.QueryRow(context.Background(),
		`SELECT id, chapter_id, title, time_limit, max_attempts, type
		 FROM tests WHERE chapter_id=$1 AND type='chapter'`, chapterID)

	var t Test
	err := row.Scan(&t.ID, &t.ChapterID, &t.Title, &t.TimeLimit, &t.MaxAttempts, &t.Type)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func GetFinalTest(courseID int64) (*Test, error) {
	row := database.DB.QueryRow(context.Background(),
		`SELECT id, course_id, title, time_limit, max_attempts, type
		 FROM tests WHERE course_id=$1 AND type='final'`, courseID)

	var t Test
	err := row.Scan(&t.ID, &t.CourseID, &t.Title, &t.TimeLimit, &t.MaxAttempts, &t.Type)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func GetQuestions(testID int64) ([]Question, error) {
	rows, err := database.DB.Query(context.Background(),
		`SELECT id, test_id, type, question_text, order_index
		 FROM questions 
		 WHERE test_id=$1 
		 ORDER BY order_index`, testID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	q := make([]Question, 0)

	for rows.Next() {
		var item Question
		err := rows.Scan(
			&item.ID,
			&item.TestID,
			&item.Type,
			&item.QuestionText,
			&item.OrderIndex,
		)
		if err != nil {
			return nil, err
		}
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

func GetMatchingPairs(questionID int64) ([]map[string]string, []string, error) {

	rows, err := database.DB.Query(context.Background(),
		`SELECT left_text, right_text FROM matching_pairs WHERE question_id=$1`,
		questionID,
	)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	pairs := []map[string]string{}
	optionsMap := map[string]bool{}

	for rows.Next() {
		var left, right string
		rows.Scan(&left, &right)

		pairs = append(pairs, map[string]string{
			"left": left,
		})

		optionsMap[right] = true
	}

	// options (перемешанный список)
	options := []string{}
	for k := range optionsMap {
		options = append(options, k)
	}

	return pairs, options, nil
}

func DeleteTest(testID int64) error {
	ctx := context.Background()
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Delete attempts first since it has a foreign key to tests without cascade
	_, err = tx.Exec(ctx, `DELETE FROM attempts WHERE test_id = $1`, testID)
	if err != nil {
		return err
	}

	// due to cascade rules on database questions, answers and matching pairs will delete automatically
	_, err = tx.Exec(ctx, `DELETE FROM tests WHERE id = $1`, testID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

