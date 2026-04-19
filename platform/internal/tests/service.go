package tests

import (
	"context"
	"fmt"
	"platform/internal/database"
	"strconv"
	"strings"
)

func StartTest(userID, moduleID int64) (*Test, []Question, map[int64]interface{}, error) {

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

			answerMap[q.ID] = []interface{}{} // можно и nil

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

	score := 0

	questions, err := GetQuestions(testID)
	if err != nil {
		return 0, err
	}

	for _, q := range questions {
		fmt.Println("QUESTION:", q.ID, "TYPE:", q.Type)

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
