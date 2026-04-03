package modules

import (
	"context"
	"platform/internal/database"
)

func CreateModule(m *Module) error {
	query := `
	INSERT INTO modules (course_id, title, order_index)
	VALUES ($1, $2, $3)
	RETURNING id, created_at`

	return database.DB.QueryRow(context.Background(),
		query,
		m.CourseID,
		m.Title,
		m.OrderIndex,
	).Scan(&m.ID, &m.CreatedAt)
}

func GetModulesByCourse(courseID int64) ([]Module, error) {

	rows, err := database.DB.Query(context.Background(),
		`SELECT id, course_id, title, order_index, created_at
		 FROM modules
		 WHERE course_id = $1
		 ORDER BY order_index ASC`, courseID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	modules := make([]Module, 0)

	for rows.Next() {
		var m Module
		err := rows.Scan(&m.ID, &m.CourseID, &m.Title, &m.OrderIndex, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		modules = append(modules, m)
	}

	return modules, nil
}

func CreateMaterial(mat *Material) error {

	query := `
	INSERT INTO materials 
	(module_id, title, type, file_url, external_url, content, order_index)
	VALUES ($1,$2,$3,$4,$5,$6,$7)
	RETURNING id, created_at`

	return database.DB.QueryRow(context.Background(),
		query,
		mat.ModuleID,
		mat.Title,
		mat.Type,
		mat.FileURL,
		mat.ExternalURL,
		mat.Content,
		mat.OrderIndex,
	).Scan(&mat.ID, &mat.CreatedAt)
}

func GetMaterialsByModule(moduleID int64) ([]Material, error) {

	rows, err := database.DB.Query(context.Background(),
		`SELECT id, module_id, title, type, file_url, external_url, content, order_index, created_at
		 FROM materials
		 WHERE module_id = $1
		 ORDER BY order_index ASC`, moduleID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	materials := make([]Material, 0)

	for rows.Next() {
		var m Material
		err := rows.Scan(
			&m.ID,
			&m.ModuleID,
			&m.Title,
			&m.Type,
			&m.FileURL,
			&m.ExternalURL,
			&m.Content,
			&m.OrderIndex,
			&m.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		materials = append(materials, m)
	}

	return materials, nil
}
