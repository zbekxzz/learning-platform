package courses

import "errors"

func CreateCourse(title, description string, createdBy int64, role string) (*Course, error) {

	if role != "admin" && role != "teacher" {
		return nil, errors.New("forbidden")
	}

	course := &Course{
		Title:       title,
		Description: description,
		CreatedBy:   createdBy,
		IsPublished: false,
	}

	err := Create(course)
	if err != nil {
		return nil, err
	}

	return course, nil
}

func PublishCourse(id int64, role string) error {

	if role != "admin" && role != "teacher" {
		return errors.New("forbidden")
	}

	return Publish(id)
}
