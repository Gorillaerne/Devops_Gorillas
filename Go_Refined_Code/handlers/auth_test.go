package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

// TestMain sets a low bcrypt cost so register/login tests run fast.
func TestMain(m *testing.M) {
	bcryptCost = bcrypt.MinCost
	os.Exit(m.Run())
}

func newUsersDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE users (
		id       INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT    NOT NULL UNIQUE,
		email    TEXT    NOT NULL,
		password TEXT    NOT NULL
	)`)
	if err != nil {
		t.Fatalf("create users table: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// --- Unit tests for password helpers ---

func TestHashAndVerifyPassword_Bcrypt(t *testing.T) {
	hash, err := hashPassword("secret")
	if err != nil {
		t.Fatalf("hashPassword: %v", err)
	}
	if !verifyPassword("secret", hash) {
		t.Error("bcrypt verify should succeed with correct password")
	}
	if verifyPassword("wrong", hash) {
		t.Error("bcrypt verify should fail with wrong password")
	}
}

// --- Integration tests for HandleAPIRegister ---

func TestHandleAPIRegister_Success(t *testing.T) {
	db := newUsersDB(t)

	body := `{"username":"alice","email":"alice@example.com","password":"pass123","password2":"pass123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandleAPIRegister(db)(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandleAPIRegister_PasswordMismatch(t *testing.T) {
	db := newUsersDB(t)

	body := `{"username":"bob","email":"bob@example.com","password":"abc","password2":"xyz"}`
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandleAPIRegister(db)(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}

func TestHandleAPIRegister_DuplicateUser(t *testing.T) {
	db := newUsersDB(t)

	body := `{"username":"carol","email":"carol@example.com","password":"pass","password2":"pass"}`

	// First registration should succeed
	req1 := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBufferString(body))
	req1.Header.Set("Content-Type", "application/json")
	HandleAPIRegister(db)(httptest.NewRecorder(), req1)

	// Second registration of the same username should conflict
	req2 := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBufferString(body))
	req2.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	HandleAPIRegister(db)(w, req2)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409 on duplicate, got %d", w.Code)
	}
}

func TestHandleAPIRegister_InvalidJSON(t *testing.T) {
	db := newUsersDB(t)

	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()

	HandleAPIRegister(db)(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- Integration tests for HandleAPILogin ---

// seedUser inserts a user with a bcrypt-hashed password directly into the DB.
func seedUser(t *testing.T, db *sql.DB, username, email, password string) {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}
	_, err = db.Exec(
		`INSERT INTO users (username, email, password) VALUES (?, ?, ?)`,
		username, email, string(hash),
	)
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
}

func TestHandleAPILogin_Success(t *testing.T) {
	db := newUsersDB(t)
	seedUser(t, db, "dave", "dave@example.com", "mypass")

	body := `{"username":"dave","password":"mypass"}`
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandleAPILogin(db)(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp AuthResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Token == "" {
		t.Error("expected non-empty JWT token in response")
	}
}

func TestHandleAPILogin_WrongPassword(t *testing.T) {
	db := newUsersDB(t)
	seedUser(t, db, "eve", "eve@example.com", "correct")

	body := `{"username":"eve","password":"wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandleAPILogin(db)(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestHandleAPILogin_UnknownUser(t *testing.T) {
	db := newUsersDB(t)

	body := `{"username":"nobody","password":"pass"}`
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandleAPILogin(db)(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestHandleAPILogin_InvalidJSON(t *testing.T) {
	db := newUsersDB(t)

	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString("bad"))
	w := httptest.NewRecorder()

	HandleAPILogin(db)(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
