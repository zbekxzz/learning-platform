package auth

import (
	"context"
	"platform/internal/database"
)

func CreateUser(user *User) error {
	query := `
	INSERT INTO users (full_name, email, password_hash, role)
	VALUES ($1, $2, $3, $4)
	RETURNING id`

	return database.DB.QueryRow(context.Background(),
		query,
		user.FullName,
		user.Email,
		user.PasswordHash,
		user.Role,
	).Scan(&user.ID)
}

func GetUserByEmail(email string) (*User, error) {
	query := `
	SELECT id, full_name, email, password_hash, role
	FROM users WHERE email=$1`

	row := database.DB.QueryRow(context.Background(), query, email)

	var user User
	err := row.Scan(&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(id int64) (*User, error) {
	query := `
	SELECT id, full_name, email, password_hash, role
	FROM users WHERE id=$1`

	row := database.DB.QueryRow(context.Background(), query, id)

	var user User
	err := row.Scan(&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
