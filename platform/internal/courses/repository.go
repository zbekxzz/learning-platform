package courses

import (
	"context"
	"platform/internal/database"
)

func Create(course *Course) error {
	query := `
	INSERT INTO courses (title, description, created_by, is_published)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at`

	return database.DB.QueryRow(
		context.Background(),
		query,
		course.Title,
		course.Description,
		course.CreatedBy,
		course.IsPublished,
	).Scan(&course.ID, &course.CreatedAt)
}

func GetAllPublished() ([]Course, error) {
	query := `
	SELECT id, title, description, created_by, is_published, created_at
	FROM courses
	WHERE is_published = TRUE
	ORDER BY created_at DESC`

	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []Course

	for rows.Next() {
		var c Course
		err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.CreatedBy, &c.IsPublished, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}

	return courses, nil
}

func GetByID(id int64) (*Course, error) {
	query := `
	SELECT id, title, description, created_by, is_published, created_at
	FROM courses
	WHERE id = $1`

	row := database.DB.QueryRow(context.Background(), query, id)

	var c Course
	err := row.Scan(&c.ID, &c.Title, &c.Description, &c.CreatedBy, &c.IsPublished, &c.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func Publish(id int64) error {
	query := `UPDATE courses SET is_published = TRUE WHERE id = $1`
	_, err := database.DB.Exec(context.Background(), query, id)
	return err
}
