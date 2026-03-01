// Package database to handle datasbase
package database

import (
	"database/sql"
	"log"
	"os"

	/* SQL import */
	_ "github.com/go-sql-driver/mysql"
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

	log.Println("Connecting to MySQL database...")

	dsn := os.Getenv("DATABASE_PATH")
	if dsn == "" {
		log.Fatal("DATABASE_PATH environment variable is not set")
	}

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	if err := DB.Ping(); err != nil {
		return err
	}

	log.Println("✅ MySQL database connected")

	// List tables
	rows, err := DB.Query(`SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE()`)
	if err != nil {
		return err
	}

	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			// Log the error and decide whether to continue or return
			log.Printf("error scanning row: %v", err)
			continue // or return err
		}
		log.Println("Found table:", table)
	}

	return nil
}
