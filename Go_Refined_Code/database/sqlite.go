package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Connect() error {
	var err error
	DB, err = sql.Open("sqlite3", "./data/Gorilla_whoknows.db")
	if err != nil {
		return err
	}
	return DB.Ping()
}
