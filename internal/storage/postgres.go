package storage

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewPostgresDB initializes the connection and ensures the schema exists.
func NewPostgresDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Automatic Schema Migration
	// Ensures the table exists regardless of which environment the app runs in.
	err = createSchema(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS crawled_urls (
		id SERIAL PRIMARY KEY,
		url TEXT UNIQUE NOT NULL,
		status_code INT,
		error_message TEXT,
		crawled_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(schema)
	if err != nil {
		log.Printf("Failed to create schema: %v", err)
		return err
	}
	return nil
}
