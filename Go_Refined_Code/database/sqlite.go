package database

import (
	"database/sql"
	_ "github.com/glebarez/go-sqlite"
	"log"
	"os"
	"github.com/joho/godotenv"
)


var DB *sql.DB

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
	defer rows.Close()

	for rows.Next() {
		var table string
		rows.Scan(&table)
		log.Println("Found table:", table)
	}

	return nil
}
