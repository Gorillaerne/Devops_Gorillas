// Package migrations registers Goose migration functions.
package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	goose.AddMigrationContext(upSeedAdmin, downSeedAdmin)
}

func upSeedAdmin(_ context.Context, tx *sql.Tx) error {
	username := os.Getenv("ADMIN_USERNAME")
	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")

	if username == "" || email == "" || password == "" {
		return fmt.Errorf("seed admin: ADMIN_USERNAME, ADMIN_EMAIL and ADMIN_PASSWORD must be set")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("seed admin: failed to hash password: %w", err)
	}

	_, err = tx.Exec(`
		INSERT IGNORE INTO users (username, email, password)
		VALUES (?, ?, ?)
	`, username, email, string(hash))
	return err
}

func downSeedAdmin(_ context.Context, tx *sql.Tx) error {
	username := os.Getenv("ADMIN_USERNAME")
	if username == "" {
		return nil
	}
	_, err := tx.Exec(`DELETE FROM users WHERE username = ?`, username)
	return err
}
