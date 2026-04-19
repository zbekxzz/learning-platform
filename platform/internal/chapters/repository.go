package chapters

import (
	"context"
	"platform/internal/database"
)

func CreateChapter(c *Chapter) error {

	query := `
	INSERT INTO chapters (course_id, title, order_index)
	VALUES ($1,$2,$3)
	RETURNING id, created_at`

	return database.DB.QueryRow(context.Background(),
		query,
		c.CourseID,
		c.Title,
		c.OrderIndex,
	).Scan(&c.ID, &c.CreatedAt)
}

func GetChaptersByCourse(courseID int64) ([]Chapter, error) {

	rows, err := database.DB.Query(context.Background(),
		`SELECT id, course_id, title, order_index, created_at
		 FROM chapters
		 WHERE course_id=$1
		 ORDER BY order_index`, courseID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []Chapter{}

	for rows.Next() {
		var c Chapter
		rows.Scan(&c.ID, &c.CourseID, &c.Title, &c.OrderIndex, &c.CreatedAt)
		result = append(result, c)
	}

	return result, nil
}
