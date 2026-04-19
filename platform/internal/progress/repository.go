package progress

import (
	"context"
	"platform/internal/database"
)

func MarkCompleted(userID, moduleID int64) error {

	_, err := database.DB.Exec(context.Background(),
		`INSERT INTO progress (user_id, module_id, is_completed)
		 VALUES ($1,$2,TRUE)
		 ON CONFLICT (user_id, module_id)
		 DO UPDATE SET is_completed=TRUE`,
		userID, moduleID)

	return err
}

func GetCompletedModules(userID int64) ([]int64, error) {

	rows, _ := database.DB.Query(context.Background(),
		`SELECT module_id FROM progress WHERE user_id=$1 AND is_completed=TRUE`, userID)

	defer rows.Close()

	var result []int64

	for rows.Next() {
		var id int64
		rows.Scan(&id)
		result = append(result, id)
	}

	return result, nil
}
