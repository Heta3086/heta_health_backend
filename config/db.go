package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	DB = db
	ensureSchema()
	fmt.Println("✅ DB Connected")
}

func ensureSchema() {
	if DB == nil {
		return
	}

	queries := []string{
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS auth_user_id INT`,
		`CREATE UNIQUE INDEX IF NOT EXISTS users_auth_user_id_unique_idx ON users(auth_user_id) WHERE auth_user_id IS NOT NULL`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			log.Printf("schema migration warning: %v", err)
		}
	}
}
