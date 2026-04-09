// Package handlers authApi
package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET")) // In production, use os.Getenv("JWT_SECRET")

// Claims struct that contains the userToken
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// Helper: Hash password (MD5 to match your old DB)
func hashPassword(password string) (string, error) {
	// GenerateFromPassword handles salting and hashing automatically
	// Cost of 10-14 is usually a good balance of speed vs security
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Helper: Verify password using bcrypt
func verifyPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Helper to send JSON responses consistently
func sendJSON(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(AuthResponse{
		StatusCode: code,
		Message:    message,
	})
	if err != nil {
		// Log the error because the client might not have received the JSON
		log.Printf("sendJSON: failed to encode response: %v", err)
	}
}

// --- Handlers ---

// HandleAPIRegister POST /api/register
func HandleAPIRegister(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username  string `json:"username"`
			Email     string `json:"email"`
			Password  string `json:"password"` //nolint:gosec
			Password2 string `json:"password2"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendJSON(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		if req.Password != req.Password2 {
			sendJSON(w, http.StatusUnprocessableEntity, "Passwords do not match")
			return
		}

		hashed, err := hashPassword(req.Password)
		if err != nil {
			sendJSON(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		_, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
			req.Username, req.Email, hashed)

		if err != nil {
			sendJSON(w, http.StatusConflict, "User already exists or DB error")
			return
		}

		sendJSON(w, http.StatusCreated, "User created successfully")
	}
}

// HandleAPILogin POST /api/login
func HandleAPILogin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req struct {
			Username string `json:"username"`
			Password string `json:"password"` //nolint:gosec
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendJSON(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		var userID int
		var hashedPw string
		err := db.QueryRow("SELECT id, password FROM users WHERE username = ?", req.Username).Scan(&userID, &hashedPw)

		if err == sql.ErrNoRows ||
			!verifyPassword(req.Password, hashedPw) {
			sendJSON(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		// Create JWT
		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &Claims{
			UserID: userID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			sendJSON(w, http.StatusInternalServerError, "Error generating token")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// It is good practice to set the status code explicitly before encoding
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(AuthResponse{
			StatusCode: 200,
			Message:    "Login successful",
			Token:      tokenString,
		})

		if err != nil {
			// We log it because we can't send a new HTTP error once encoding starts
			log.Printf("authApi: failed to encode login response: %v", err)
		}
	}
}
