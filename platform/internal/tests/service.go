package tests

import (
	"context"
	"platform/internal/database"
)

func StartTest(userID, moduleID int64) (*Test, []Question, map[int64][]Answer, error) {

	test, err := GetTestByModule(moduleID)
	if err != nil {
		return nil, nil, nil, err
	}

	questions, _ := GetQuestions(test.ID)

	answerMap := make(map[int64][]Answer)

	for _, q := range questions {
		ans, _ := GetAnswers(q.ID)

		// скрываем правильные ответы
		for i := range ans {
			ans[i].IsCorrect = false
		}

		answerMap[q.ID] = ans
	}

	CreateAttempt(userID, test.ID)

	return test, questions, answerMap, nil
}

func SubmitTest(userAnswers map[int64]int64, testID int64, userID int64) (int, error) {

	score := 0

	for qID, answerID := range userAnswers {

		var correct bool

		database.DB.QueryRow(context.Background(),
			`SELECT is_correct FROM answers WHERE id=$1 AND question_id=$2`,
			answerID, qID,
		).Scan(&correct)

		if correct {
			score++
		}
	}

	// обновляем attempt
	database.DB.Exec(context.Background(),
		`UPDATE attempts SET score=$1, finished_at=NOW()
		 WHERE user_id=$2 AND test_id=$3`,
		score, userID, testID)

	return score, nil
}
