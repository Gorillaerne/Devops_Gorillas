package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// makeToken creates a signed JWT for the given userID, using the package-level jwtKey.
func makeToken(t *testing.T, userID int) string {
	t.Helper()
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtKey)
	if err != nil {
		t.Fatalf("makeToken: %v", err)
	}
	return signed
}

func TestHandleAPIChangePassword_Success(t *testing.T) {
	db := newUsersDB(t)
	seedUser(t, db, "frank", "frank@example.com", "oldpass")

	var userID int
	_ = db.QueryRow("SELECT id FROM users WHERE username = ?", "frank").Scan(&userID)
	token := makeToken(t, userID)

	body := `{"current_password":"oldpass","new_password":"newpass","new_password2":"newpass"}`
	req := httptest.NewRequest(http.MethodPost, "/api/change-password", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w := httptest.NewRecorder()

	HandleAPIChangePassword(db)(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandleAPIChangePassword_WrongCurrentPassword(t *testing.T) {
	db := newUsersDB(t)
	seedUser(t, db, "grace", "grace@example.com", "correct")

	var userID int
	_ = db.QueryRow("SELECT id FROM users WHERE username = ?", "grace").Scan(&userID)
	token := makeToken(t, userID)

	body := `{"current_password":"wrong","new_password":"new","new_password2":"new"}`
	req := httptest.NewRequest(http.MethodPost, "/api/change-password", bytes.NewBufferString(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w := httptest.NewRecorder()

	HandleAPIChangePassword(db)(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestHandleAPIChangePassword_PasswordMismatch(t *testing.T) {
	db := newUsersDB(t)
	seedUser(t, db, "henry", "henry@example.com", "pass")

	var userID int
	_ = db.QueryRow("SELECT id FROM users WHERE username = ?", "henry").Scan(&userID)
	token := makeToken(t, userID)

	body := `{"current_password":"pass","new_password":"abc","new_password2":"xyz"}`
	req := httptest.NewRequest(http.MethodPost, "/api/change-password", bytes.NewBufferString(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w := httptest.NewRecorder()

	HandleAPIChangePassword(db)(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}

func TestHandleAPIChangePassword_NoToken(t *testing.T) {
	db := newUsersDB(t)

	body := `{"current_password":"a","new_password":"b","new_password2":"b"}`
	req := httptest.NewRequest(http.MethodPost, "/api/change-password", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	HandleAPIChangePassword(db)(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestHandleAPIChangePassword_InvalidJSON(t *testing.T) {
	db := newUsersDB(t)
	seedUser(t, db, "irma", "irma@example.com", "pass")

	var userID int
	_ = db.QueryRow("SELECT id FROM users WHERE username = ?", "irma").Scan(&userID)
	token := makeToken(t, userID)

	req := httptest.NewRequest(http.MethodPost, "/api/change-password", bytes.NewBufferString("bad-json"))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w := httptest.NewRecorder()

	HandleAPIChangePassword(db)(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleAPILogin_BreachedFlagSet(t *testing.T) {
	db := newUsersDB(t)
	// Seed a breached user with their leaked password
	seedUser(t, db, "Benthe1954", "benthe@example.com", "^Jt^pLkzW2")

	body := `{"username":"Benthe1954","password":"^Jt^pLkzW2"}`
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandleAPILogin(db)(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp AuthResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !resp.Breached {
		t.Error("expected breached=true for a known breached user")
	}
}

func TestHandleAPILogin_BreachedFlagNotSetForNormalUser(t *testing.T) {
	db := newUsersDB(t)
	seedUser(t, db, "normaluser", "normal@example.com", "safepassword")

	body := `{"username":"normaluser","password":"safepassword"}`
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandleAPILogin(db)(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp AuthResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Breached {
		t.Error("expected breached=false for a normal user")
	}
}
