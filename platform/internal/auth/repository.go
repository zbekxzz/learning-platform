package auth

import (
	"context"
	"time"
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

func SaveRestoreCode(email, code string, expiresAt time.Time) error {
	ctx := context.Background()
	
	// Delete any existing codes for this email first
	_, err := database.DB.Exec(ctx, "DELETE FROM restore_codes WHERE email = $1", email)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO restore_codes (email, code, expires_at)
	VALUES ($1, $2, $3)`

	_, err = database.DB.Exec(ctx, query, email, code, expiresAt)
	return err
}

func GetRestoreCode(email string) (string, time.Time, error) {
	query := `
	SELECT code, expires_at
	FROM restore_codes
	WHERE email = $1
	ORDER BY created_at DESC
	LIMIT 1`

	var code string
	var expiresAt time.Time

	err := database.DB.QueryRow(context.Background(), query, email).Scan(&code, &expiresAt)
	if err != nil {
		return "", time.Time{}, err
	}

	return code, expiresAt, nil
}

func DeleteRestoreCode(email string) error {
	_, err := database.DB.Exec(context.Background(), "DELETE FROM restore_codes WHERE email = $1", email)
	return err
}

func UpdateUserPassword(email string, passwordHash string) error {
	query := `
	UPDATE users
	SET password_hash = $1
	WHERE email = $2`

	_, err := database.DB.Exec(context.Background(), query, passwordHash, email)
	return err
}

