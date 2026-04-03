package enrollments

import (
	"errors"
	"platform/internal/courses"
)

func EnrollUser(userID, courseID int64, role string) error {

	if role != "student" {
		return errors.New("only students can enroll")
	}

	course, err := courses.GetByID(courseID)
	if err != nil {
		return errors.New("course not found")
	}

	if !course.IsPublished {
		return errors.New("course not published")
	}

	exists, err := Exists(userID, courseID)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("already enrolled")
	}

	return Create(userID, courseID)
}

func UnenrollUser(userID, courseID int64, role string) error {

	if role != "student" {
		return errors.New("only students can unenroll")
	}

	return Delete(userID, courseID)
}
