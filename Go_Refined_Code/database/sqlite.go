// Package database to handle datasbase
package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	/* SQL import */
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

// DB connection to DB
var DB *sql.DB

// Connect - connection to DB
func Connect() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using defaults or system environment")
	}

	log.Println("Connecting to Postgres database...")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_DB"),
	)

	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	if err := DB.Ping(); err != nil {
		return err
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("✅ Postgres database connected")

	// List tables
	rows, err := DB.Query(`SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'`)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			log.Printf("error scanning row: %v", err)
			continue
		}
		log.Println("Found table:", table)
	}

	return nil
}

// PurgeMD5Users deletes all users whose password is not a bcrypt hash.
// Bcrypt hashes always start with "$2", so anything else is a legacy MD5 hash.
func PurgeMD5Users() {
	result, err := DB.Exec("DELETE FROM users WHERE password NOT LIKE '$2%'")
	if err != nil {
		log.Printf("PurgeMD5Users: failed to delete legacy users: %v", err)
		return
	}
	n, _ := result.RowsAffected()
	if n > 0 {
		log.Printf("PurgeMD5Users: removed %d legacy MD5 user(s)", n)
	}
}