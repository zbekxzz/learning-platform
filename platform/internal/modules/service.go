package modules

import (
	"context"
	"errors"
	"platform/internal/database"
)

func CreateModuleService(chapterID int64, title string, order int, role string) (*Module, error) {

	if role != "admin" && role != "teacher" {
		return nil, errors.New("forbidden")
	}

	module := &Module{
		ChapterID:  chapterID,
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

func GetPreviousModule(moduleID int64) (int64, error) {

	var chapterID int64
	var orderIndex int

	err := database.DB.QueryRow(context.Background(),
		`SELECT chapter_id, order_index 
		 FROM modules 
		 WHERE id=$1`,
		moduleID,
	).Scan(&chapterID, &orderIndex)

	if err != nil {
		return 0, err
	}

	var prevID int64

	err = database.DB.QueryRow(context.Background(),
		`SELECT id FROM modules
		 WHERE chapter_id=$1 AND order_index=$2`,
		chapterID, orderIndex-1,
	).Scan(&prevID)

	if err != nil {
		// первый модуль
		return 0, nil
	}

	return prevID, nil
}

func UpdateMaterialsOrderService(moduleID int64, updates []OrderUpdate, role string) error {
	if role != "admin" && role != "teacher" {
		return errors.New("forbidden")
	}
	return UpdateMaterialsOrder(moduleID, updates)
}
