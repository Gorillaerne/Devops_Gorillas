package handlers

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET")) // In production, use os.Getenv("JWT_SECRET")

type Claims struct {
    UserID int `json:"user_id"`
    jwt.RegisteredClaims
}






// Helper: Hash password (MD5 to match your old DB)
func hashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

// Helper: Verify password
func verifyPassword(storedHash, password string) bool {
	return storedHash == hashPassword(password)
}

// Helper to send JSON responses consistently
func sendJSON(w http.ResponseWriter, code int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(AuthResponse{
        StatusCode: code,
        Message:    message,
    })
}

// --- Handlers ---

// POST /api/register
func HandleApiRegister(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Username  string `json:"username"`
            Email     string `json:"email"`
            Password  string `json:"password"`
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

        hashed := hashPassword(req.Password)
        _, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", 
            req.Username, req.Email, hashed)
        
        if err != nil {
            sendJSON(w, http.StatusConflict, "User already exists or DB error")
            return
        }

        sendJSON(w, http.StatusCreated, "User created successfully")
    }
}

// POST /api/login
func HandleApiLogin(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }

        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            sendJSON(w, http.StatusBadRequest, "Invalid JSON")
            return
        }

        var userID int
        var hashedPw string
        err := db.QueryRow("SELECT id, password FROM users WHERE username = ?", req.Username).Scan(&userID, &hashedPw)

        if err == sql.ErrNoRows || !verifyPassword(hashedPw, req.Password) {
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

        // Return token in the JSON body
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(AuthResponse{
            StatusCode: 200,
            Message:    "Login successful",
            Token:      tokenString,
        })
    }
}
