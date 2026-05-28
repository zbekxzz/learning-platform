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

func TogglePublishCourse(id int64, userID int64, role string) error {

	if role == "admin" {
		return TogglePublish(id)
	}

	if role == "teacher" {
		course, err := GetByID(id)
		if err != nil {
			return errors.New("course not found")
		}
		if course.CreatedBy != userID {
			return errors.New("forbidden: not your course")
		}
		return TogglePublish(id)
	}

	return errors.New("forbidden")
}

func GetTeacherCourses(userID int64) ([]Course, error) {
	return GetByCreator(userID)
}

func DeleteCourse(id int64, role string) error {

	if role != "admin" {
		return errors.New("forbidden")
	}

	return SoftDelete(id)
}
