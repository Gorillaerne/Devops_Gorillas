// Package handlers profileApi
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// parseJWT extracts and validates the JWT from the Authorization header.
func parseJWT(r *http.Request) (*Claims, error) {
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(_ *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	return claims, err
}

// HandleAPIChangePassword POST /api/change-password
func HandleAPIChangePassword(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := parseJWT(r)
		if err != nil {
			sendJSON(w, http.StatusUnauthorized, "Invalid or missing token")
			return
		}

		var req struct {
			CurrentPassword string `json:"current_password"` //nolint:gosec
			NewPassword     string `json:"new_password"`     //nolint:gosec
			NewPassword2    string `json:"new_password2"`    //nolint:gosec
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendJSON(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		if req.NewPassword != req.NewPassword2 {
			sendJSON(w, http.StatusUnprocessableEntity, "Passwords do not match")
			return
		}

		var hashedPw string
		err = db.QueryRow("SELECT password FROM users WHERE id = ?", claims.UserID).Scan(&hashedPw)
		if err == sql.ErrNoRows {
			sendJSON(w, http.StatusUnauthorized, "User not found")
			return
		}
		if err != nil {
			sendJSON(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		if !verifyPassword(req.CurrentPassword, hashedPw) {
			sendJSON(w, http.StatusUnauthorized, "Current password is incorrect")
			return
		}

		newHash, err := hashPassword(req.NewPassword)
		if err != nil {
			sendJSON(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		_, err = db.Exec("UPDATE users SET password = ? WHERE id = ?", newHash, claims.UserID)
		if err != nil {
			sendJSON(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		sendJSON(w, http.StatusOK, "Password updated successfully")
	}
}
