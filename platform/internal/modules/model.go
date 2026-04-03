package modules

import "time"

type Module struct {
	ID         int64     `json:"id"`
	CourseID   int64     `json:"course_id"`
	Title      string    `json:"title"`
	OrderIndex int       `json:"order_index"`
	CreatedAt  time.Time `json:"created_at"`
}

type Material struct {
	ID          int64     `json:"id"`
	ModuleID    int64     `json:"module_id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	FileURL     *string   `json:"file_url"`
	ExternalURL *string   `json:"external_url"`
	Content     *string   `json:"content"`
	OrderIndex  int       `json:"order_index"`
	CreatedAt   time.Time `json:"created_at"`
}
