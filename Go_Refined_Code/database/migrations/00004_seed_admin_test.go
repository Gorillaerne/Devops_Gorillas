package migrations

import (
	"context"
	"testing"
)

func TestUpSeedAdmin_MissingEnvVars(t *testing.T) {
	t.Setenv("ADMIN_USERNAME", "")
	t.Setenv("ADMIN_EMAIL", "")
	t.Setenv("ADMIN_PASSWORD", "")

	// tx is nil but the function returns before using it when env vars are empty.
	err := upSeedAdmin(context.Background(), nil)
	if err == nil {
		t.Error("expected error when ADMIN_USERNAME/EMAIL/PASSWORD are unset")
	}
}

func TestUpSeedAdmin_MissingEmail(t *testing.T) {
	t.Setenv("ADMIN_USERNAME", "admin")
	t.Setenv("ADMIN_EMAIL", "")
	t.Setenv("ADMIN_PASSWORD", "pass")

	err := upSeedAdmin(context.Background(), nil)
	if err == nil {
		t.Error("expected error when ADMIN_EMAIL is unset")
	}
}

func TestDownSeedAdmin_EmptyUsername(t *testing.T) {
	t.Setenv("ADMIN_USERNAME", "")

	// Returns nil immediately without touching the (nil) transaction.
	err := downSeedAdmin(context.Background(), nil)
	if err != nil {
		t.Errorf("expected nil error for empty username, got %v", err)
	}
}
