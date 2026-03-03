package enrollments

import (
	"context"
	"errors"
	"platform/internal/database"
)

func Create(userID, courseID int64) error {

	query := `
	INSERT INTO enrollments (user_id, course_id)
	VALUES ($1, $2)
	`

	_, err := database.DB.Exec(context.Background(), query, userID, courseID)
	return err
}

func Exists(userID, courseID int64) (bool, error) {

	var exists bool

	query := `
	SELECT EXISTS (
		SELECT 1 FROM enrollments
		WHERE user_id = $1 AND course_id = $2
	)
	`

	err := database.DB.QueryRow(context.Background(), query, userID, courseID).Scan(&exists)
	return exists, err
}

func GetUserEnrollments(userID int64) ([]Enrollment, error) {

	query := `
	SELECT id, user_id, course_id, enrolled_at
	FROM enrollments
	WHERE user_id = $1
	ORDER BY enrolled_at DESC
	`

	rows, err := database.DB.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	enrollments := make([]Enrollment, 0)

	for rows.Next() {
		var e Enrollment
		err := rows.Scan(&e.ID, &e.UserID, &e.CourseID, &e.EnrolledAt)
		if err != nil {
			return nil, err
		}
		enrollments = append(enrollments, e)
	}

	return enrollments, nil
}

func Delete(userID, courseID int64) error {

	query := `
	DELETE FROM enrollments
	WHERE user_id = $1 AND course_id = $2
	`

	result, err := database.DB.Exec(context.Background(), query, userID, courseID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("enrollment not found")
	}

	return nil
}
