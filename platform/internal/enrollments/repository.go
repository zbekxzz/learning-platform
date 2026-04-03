package enrollments

import (
	"context"
	"errors"
	"platform/internal/courses"
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

func GetUserCourses(userID int64) ([]courses.Course, error) {

	rows, err := database.DB.Query(context.Background(), `
		SELECT c.id, c.title, c.description, c.is_published, c.created_at
		FROM enrollments e
		JOIN courses c ON c.id = e.course_id
		WHERE e.user_id = $1 AND c.deleted_at IS NULL
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]courses.Course, 0)

	for rows.Next() {
		var c courses.Course
		rows.Scan(&c.ID, &c.Title, &c.Description, &c.IsPublished, &c.CreatedAt)
		result = append(result, c)
	}

	return result, nil
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
