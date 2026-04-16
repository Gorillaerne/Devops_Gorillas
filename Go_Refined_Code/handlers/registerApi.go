// Package handlers register
package handlers

/* RegisterRequest struct to handle registrations */
type RegisterRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"` //nolint:gosec
	Password2 string `json:"password2"`
}

/* LoginRequest struct to handle logins */
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"` //nolint:gosec
}

/* AuthResponse struct to handle responses */
type AuthResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Token      string `json:"token"`
	Breached   bool   `json:"breached,omitempty"`
}
