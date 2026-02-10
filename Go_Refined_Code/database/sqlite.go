package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var DB *sql.DB

func Connect() error {
	log.Println("Connecting to SQLite database...")

	var err error
	DB, err = sql.Open("sqlite3", "data/Gorilla_whoknows.db")
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
