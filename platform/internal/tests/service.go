package tests

import (
	"context"
	"fmt"
	"platform/internal/database"
	"strconv"
	"strings"
)

func StartTest(userID, moduleID int64) (*Test, []Question, map[int64]interface{}, error) {

	canAccess, err := CanAccessModule(userID, moduleID)
	if err != nil {
		return nil, nil, nil, err
	}

	if !canAccess {
		return nil, nil, nil, fmt.Errorf("module locked")
	}

	test, err := GetTestByModule(moduleID)
	if err != nil {
		return nil, nil, nil, err
	}

	questions, _ := GetQuestions(test.ID)

	answerMap := make(map[int64]interface{})

	for i, q := range questions {

		switch q.Type {
		case "mcq":

			ans, _ := GetAnswers(q.ID)

			for i := range ans {
				ans[i].IsCorrect = false
			}

			answerMap[q.ID] = ans

		case "matching":

			pairs, options, _ := GetMatchingPairs(q.ID)

			questions[i].Pairs = pairs
			questions[i].Options = options

			answerMap[q.ID] = []interface{}{}

		case "open":
			answerMap[q.ID] = []interface{}{}
		}
	}

	CreateAttempt(userID, test.ID)

	return test, questions, answerMap, nil
}

func StartChapterTest(userID, chapterID int64) (*Test, []Question, map[int64][]Answer, error) {

	test, err := GetTestByChapter(chapterID)
	if err != nil {
		return nil, nil, nil, err
	}

	questions, _ := GetQuestions(test.ID)

	answerMap := make(map[int64][]Answer)

	for _, q := range questions {
		ans, _ := GetAnswers(q.ID)

		for i := range ans {
			ans[i].IsCorrect = false
		}

		answerMap[q.ID] = ans
	}

	CreateAttempt(userID, test.ID)

	return test, questions, answerMap, nil
}

func StartFinalTest(userID, courseID int64) (*Test, []Question, map[int64][]Answer, error) {

	test, err := GetFinalTest(courseID)
	if err != nil {
		return nil, nil, nil, err
	}

	questions, _ := GetQuestions(test.ID)

	answerMap := make(map[int64][]Answer)

	for _, q := range questions {
		ans, _ := GetAnswers(q.ID)

		for i := range ans {
			ans[i].IsCorrect = false
		}

		answerMap[q.ID] = ans
	}

	CreateAttempt(userID, test.ID)

	return test, questions, answerMap, nil
}

func SubmitTest(userAnswers map[string]interface{}, testID int64, userID int64) (int, error) {
	var maxAttempts int

	err := database.DB.QueryRow(context.Background(),
		`SELECT max_attempts FROM tests WHERE id=$1`,
		testID,
	).Scan(&maxAttempts)

	if err != nil {
		return 0, err
	}

	totalQuestions, err := GetQuestionCount(testID)
	if err != nil {
		return 0, err
	}

	passed, err := IsTestPassed(userID, testID, totalQuestions)
	if err != nil {
		return 0, err
	}

	if passed {
		return 0, fmt.Errorf("test already passed")
	}

	attemptCount, err := GetAttemptCount(userID, testID)
	if err != nil {
		return 0, err
	}

	if attemptCount >= maxAttempts {
		return 0, fmt.Errorf("no attempts left")
	}

	score := 0

	questions, err := GetQuestions(testID)
	if err != nil {
		return 0, err
	}

	for _, q := range questions {

		key := strconv.FormatInt(q.ID, 10)

		userAnswer, exists := userAnswers[key]
		if !exists {
			continue
		}

		switch q.Type {

		case "mcq":

			var answerID int64

			switch v := userAnswer.(type) {
			case float64:
				answerID = int64(v)
			case int64:
				answerID = v
			default:
				continue
			}

			var correct bool

			err := database.DB.QueryRow(context.Background(),
				`SELECT is_correct FROM answers WHERE id=$1 AND question_id=$2`,
				answerID, q.ID,
			).Scan(&correct)

			if err == nil && correct {
				score++
			}

		case "open":

			userText, ok := userAnswer.(string)
			if !ok {
				continue
			}

			var correctText string

			err := database.DB.QueryRow(context.Background(),
				`SELECT correct_text FROM questions WHERE id=$1`,
				q.ID,
			).Scan(&correctText)

			if err != nil {
				continue
			}

			if normalize(userText) == normalize(correctText) {
				score++
			}

		case "matching":

			userMap, ok := userAnswer.(map[string]interface{})
			if !ok {
				continue
			}

			rows, err := database.DB.Query(context.Background(),
				`SELECT left_text, right_text FROM matching_pairs WHERE question_id=$1`,
				q.ID,
			)
			if err != nil {
				continue
			}

			correctPairs := make(map[string]string)

			for rows.Next() {
				var left, right string
				rows.Scan(&left, &right)
				correctPairs[left] = right
			}
			rows.Close()

			correctCount := 0

			for left, right := range correctPairs {

				userVal, ok := userMap[left]
				if !ok {
					continue
				}

				userStr, ok := userVal.(string)
				if !ok {
					continue
				}

				if normalize(userStr) == normalize(right) {
					correctCount++
				}
			}

			if correctCount == len(correctPairs) {
				score++
			}
		}
	}

	_, err = database.DB.Exec(context.Background(),
		`INSERT INTO attempts (user_id, test_id, score, finished_at)
		 VALUES ($1,$2,$3,NOW())`,
		userID, testID, score)

	if err != nil {
		return 0, err
	}

	return score, nil
}

func CreateFullTest(req CreateTestRequest) error {

	tx, err := database.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	var testID int64

	err = tx.QueryRow(context.Background(),
		`INSERT INTO tests (module_id, chapter_id, course_id, type, title, time_limit, max_attempts)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)
		 RETURNING id`,
		req.ModuleID,
		req.ChapterID,
		req.CourseID,
		req.Type,
		req.Title,
		req.TimeLimit,
		req.MaxAttempts,
	).Scan(&testID)

	if err != nil {
		return err
	}

	for _, q := range req.Questions {

		var questionID int64

		err := tx.QueryRow(context.Background(),
			`INSERT INTO questions (test_id, type, question_text, order_index, correct_text)
			 VALUES ($1,$2,$3,$4,$5)
			 RETURNING id`,
			testID,
			q.Type,
			q.QuestionText,
			q.OrderIndex,
			q.CorrectText, // 👈 для open
		).Scan(&questionID)

		if err != nil {
			return err
		}

		if q.Type == "mcq" {
			for _, a := range q.Answers {

				_, err := tx.Exec(context.Background(),
					`INSERT INTO answers (question_id, text, is_correct)
					 VALUES ($1,$2,$3)`,
					questionID,
					a.Text,
					a.IsCorrect,
				)

				if err != nil {
					return err
				}
			}
		}

		if q.Type == "matching" {
			for _, p := range q.Pairs {

				_, err := tx.Exec(context.Background(),
					`INSERT INTO matching_pairs (question_id, left_text, right_text)
					 VALUES ($1,$2,$3)`,
					questionID,
					p.Left,
					p.Right,
				)

				if err != nil {
					return err
				}
			}
		}
	}

	return tx.Commit(context.Background())
}

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func GetQuestionCount(testID int64) (int, error) {
	var count int
	err := database.DB.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM questions WHERE test_id=$1`,
		testID,
	).Scan(&count)
	return count, err
}

func IsTestPassed(userID, testID int64, totalQuestions int) (bool, error) {

	var exists bool

	err := database.DB.QueryRow(context.Background(),
		`SELECT EXISTS (
			SELECT 1 FROM attempts
			WHERE user_id=$1 
			AND test_id=$2 
			AND score=$3
		)`,
		userID, testID, totalQuestions,
	).Scan(&exists)

	return exists, err
}

func GetAttemptCount(userID, testID int64) (int, error) {

	var count int

	err := database.DB.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM attempts 
		 WHERE user_id=$1 AND test_id=$2`,
		userID, testID,
	).Scan(&count)

	return count, err
}

func GetPreviousModule(moduleID int64) (int64, error) {

	var chapterID int64
	var orderIndex int

	err := database.DB.QueryRow(context.Background(),
		`SELECT chapter_id, order_index 
		 FROM modules 
		 WHERE id=$1`,
		moduleID,
	).Scan(&chapterID, &orderIndex)

	if err != nil {
		return 0, err
	}

	var prevID int64

	err = database.DB.QueryRow(context.Background(),
		`SELECT id FROM modules
		 WHERE chapter_id=$1 AND order_index=$2`,
		chapterID, orderIndex-1,
	).Scan(&prevID)

	if err != nil {
		// первый модуль
		return 0, nil
	}

	return prevID, nil
}

func IsModuleCompleted(userID, moduleID int64) (bool, error) {

	var total int

	err := database.DB.QueryRow(context.Background(),
		`SELECT COUNT(q.id)
		 FROM tests t
		 JOIN questions q ON q.test_id = t.id
		 WHERE t.module_id=$1`,
		moduleID,
	).Scan(&total)

	if err != nil {
		return false, err
	}

	var exists bool

	err = database.DB.QueryRow(context.Background(),
		`SELECT EXISTS (
			SELECT 1 FROM attempts a
			JOIN tests t ON t.id = a.test_id
			WHERE t.module_id=$1
			AND a.user_id=$2
			AND a.score=$3
		)`,
		moduleID, userID, total,
	).Scan(&exists)

	return exists, err
}

func CanAccessModule(userID, moduleID int64) (bool, error) {

	prevID, err := GetPreviousModule(moduleID)
	if err != nil {
		return false, err
	}

	// первый модуль
	if prevID == 0 {
		return true, nil
	}

	return IsModuleCompleted(userID, prevID)
}
