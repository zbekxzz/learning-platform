package enrollments

import "time"

type Enrollment struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	CourseID   int64     `json:"course_id"`
	EnrolledAt time.Time `json:"enrolled_at"`
}
