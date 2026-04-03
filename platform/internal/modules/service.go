package modules

import (
	"errors"
)

func CreateModuleService(courseID int64, title string, order int, role string) (*Module, error) {

	if role != "admin" && role != "teacher" {
		return nil, errors.New("forbidden")
	}

	module := &Module{
		CourseID:   courseID,
		Title:      title,
		OrderIndex: order,
	}

	err := CreateModule(module)
	if err != nil {
		return nil, err
	}

	return module, nil
}

func CreateMaterialService(mat *Material, role string) error {

	if role != "admin" && role != "teacher" {
		return errors.New("forbidden")
	}

	return CreateMaterial(mat)
}
