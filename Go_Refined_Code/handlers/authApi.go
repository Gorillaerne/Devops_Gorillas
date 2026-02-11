package handlers

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

// Helper: Hash password (MD5 to match your old DB)
func hashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

// Helper: Verify password
func verifyPassword(storedHash, password string) bool {
	return storedHash == hashPassword(password)
}

// POST /api/register
func HandleApiRegister(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		pw1 := r.FormValue("password")
		pw2 := r.FormValue("password2")

		if pw1 != pw2 {
			tmpl.ExecuteTemplate(w, "layout.html", map[string]interface{}{"Error": "Passwords do not match"})
			return
		}

		// Fixed: hashPassword only returns one value
		hashed := hashPassword(pw1)

		_, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, hashed)
		if err != nil {
			log.Printf("Register error: %v", err)
			tmpl.ExecuteTemplate(w, "layout.html", map[string]interface{}{"Error": "User already exists or DB error"})
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// POST /api/login
func HandleApiLogin(db *sql.DB, store sessions.Store, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var userID int
		var hashedPw string

		err := db.QueryRow("SELECT id, password FROM users WHERE username = ?", username).Scan(&userID, &hashedPw)

		if err == sql.ErrNoRows || !verifyPassword(hashedPw, password) {
			tmpl.ExecuteTemplate(w, "layout.html", map[string]interface{}{"Error": "Invalid credentials"})
			return
		} else if err != nil {
			log.Printf("DB error: %v", err)
			return
		}

		session, _ := store.Get(r, "session-name")
		session.Values["user_id"] = userID
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// GET /api/logout
func HandleApiLogout(store sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		session.Options.MaxAge = -1
		session.Save(r, w)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}