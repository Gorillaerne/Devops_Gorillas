// Package database to handle datasbase
package database

import (
	"database/sql"
	"log"
	"os"

	/* SQL import */
	_ "github.com/glebarez/go-sqlite"
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

	log.Println("Connecting to SQLite database...")

	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "data/Gorilla_whoknows.db"
	}

	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	if err := DB.Ping(); err != nil {
		return err
	}

	log.Println("✅ SQLite database connected")

	// 🔍 LIST TABELLER
	rows, err := DB.Query(`SELECT name FROM sqlite_master WHERE type='table'`)
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
