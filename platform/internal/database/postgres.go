package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect() {
	dsn := os.Getenv("DB_DSN")

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	DB = pool

	_, err = DB.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS certificates (
			id SERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL,
			course_id BIGINT NOT NULL,
			certificate_url VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT unique_user_course_certificate UNIQUE (user_id, course_id)
		);
	`)
	if err != nil {
		log.Println("Warning: failed to create certificates table:", err)
	}
}
