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
}
