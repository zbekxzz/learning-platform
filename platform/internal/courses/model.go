package courses

import "time"

type Course struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedBy   int64     `json:"created_by"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
}
