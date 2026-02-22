package handlers

type RegisterRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Password2 string `json:"password2"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
    StatusCode int    `json:"statusCode"`
    Message    string `json:"message"`
    Token      string `json:"token"`
}

