package chapters

import "time"

type Chapter struct {
	ID         int64     `json:"id"`
	CourseID   int64     `json:"course_id"`
	Title      string    `json:"title"`
	OrderIndex int       `json:"order_index"`
	CreatedAt  time.Time `json:"created_at"`
}
