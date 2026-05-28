package modules

import (
	"context"
	"platform/internal/database"
)

func CreateModule(m *Module) error {
	query := `
	INSERT INTO modules (chapter_id, title, order_index)
	VALUES ($1, $2, $3)
	RETURNING id, created_at`

	return database.DB.QueryRow(context.Background(),
		query,
		m.ChapterID,
		m.Title,
		m.OrderIndex,
	).Scan(&m.ID, &m.CreatedAt)
}

func GetModulesByChapter(chapterID int64) ([]Module, error) {
	rows, _ := database.DB.Query(context.Background(),
		`SELECT id, chapter_id, title, order_index
		 FROM modules
		 WHERE chapter_id=$1
		 ORDER BY order_index`, chapterID)

	defer rows.Close()

	var result []Module

	for rows.Next() {
		var m Module
		rows.Scan(&m.ID, &m.ChapterID, &m.Title, &m.OrderIndex)
		result = append(result, m)
	}

	return result, nil
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

// OrderUpdate represents an incoming update for a material's order index.
type OrderUpdate struct {
	ID         int64 `json:"id"`
	OrderIndex int   `json:"order_index"`
}

// UpdateMaterialsOrder updates the order_index for a slice of materials within a transaction.
func UpdateMaterialsOrder(moduleID int64, updates []OrderUpdate) error {
	tx, err := database.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	query := `
		UPDATE materials
		SET order_index = $1
		WHERE id = $2 AND module_id = $3
	`

	for _, u := range updates {
		_, err := tx.Exec(context.Background(), query, u.OrderIndex, u.ID, moduleID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}
