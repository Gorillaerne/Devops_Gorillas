package database

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE users (
		id       INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT    NOT NULL,
		email    TEXT    NOT NULL,
		password TEXT    NOT NULL
	)`)
	if err != nil {
		t.Fatalf("create users table: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func TestPurgeMD5Users_RemovesLegacyPasswords(t *testing.T) {
	db := newTestDB(t)
	DB = db

	_, err := db.Exec(`INSERT INTO users (username, email, password) VALUES
		('alice', 'alice@example.com', 'abc123def456789012345678901234'),
		('bob',   'bob@example.com',   '$2a$10$validbcrypthashXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX')`)
	if err != nil {
		t.Fatalf("seed users: %v", err)
	}

	PurgeMD5Users()

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 user remaining (bcrypt), got %d", count)
	}
}

func TestPurgeMD5Users_DBError(t *testing.T) {
	// DB without a users table triggers a SQL error — PurgeMD5Users must not panic.
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer func() { _ = db.Close() }()
	DB = db
	PurgeMD5Users() // should log error and return gracefully
}

func TestPurgeMD5Users_NoLegacyUsers(t *testing.T) {
	db := newTestDB(t)
	DB = db

	_, err := db.Exec(`INSERT INTO users (username, email, password) VALUES
		('carol', 'carol@example.com', '$2b$12$validhashXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX')`)
	if err != nil {
		t.Fatalf("seed users: %v", err)
	}

	PurgeMD5Users()

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 user remaining, got %d", count)
	}
}
