package users

import (
	"context"
	"platform/internal/database"
)

func GetAllUsers() ([]User, error) {

	rows, err := database.DB.Query(context.Background(), `
		SELECT id, email, full_name, role, created_at
		FROM users
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)

	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Email, &u.FullName, &u.Role, &u.CreatedAt)
		users = append(users, u)
	}

	return users, nil
}

func CreateUser(email, fullName, passwordHash, role string) (*User, error) {

	row := database.DB.QueryRow(context.Background(), `
		INSERT INTO users (email, full_name, password_hash, role)
		VALUES ($1,$2,$3,$4)
		RETURNING id, email, full_name, role, created_at
	`, email, fullName, passwordHash, role)

	var u User
	err := row.Scan(&u.ID, &u.Email, &u.FullName, &u.Role, &u.CreatedAt)

	return &u, err
}

func UpdateUser(id int64, email, fullName, role string) error {

	_, err := database.DB.Exec(context.Background(), `
		UPDATE users
		SET email=$1, full_name=$2, role=$3
		WHERE id=$4
	`, email, fullName, role, id)

	return err
}

func DeleteUser(id int64) error {

	_, err := database.DB.Exec(context.Background(), `
		DELETE FROM users WHERE id=$1
	`, id)

	return err
}
