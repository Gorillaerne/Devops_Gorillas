# Authentication — authApi.go & registerApi.go

**Files:**
- `Go_Refined_Code/handlers/authApi.go`
- `Go_Refined_Code/handlers/registerApi.go`

These files handle user registration, login, password hashing, and JWT token generation.

---

## Data Structures

Defined in `registerApi.go`:

```go
type RegisterRequest struct { ... }  // POST /api/register body
type LoginRequest    struct { ... }  // POST /api/login body
type AuthResponse    struct {
    StatusCode int    `json:"statusCode"`
    Message    string `json:"message"`
    Token      string `json:"token,omitempty"`
    Breached   bool   `json:"breached,omitempty"`
}
```

`Claims` is defined in `authApi.go` and embeds the standard JWT claims plus a user ID:

```go
type Claims struct {
    UserID int `json:"user_id"`
    jwt.RegisteredClaims
}
```

---

## Functions

### `HandleAPIRegister(db) http.HandlerFunc`

**Route:** `POST /api/register`

1. Decodes the JSON request body (username, email, password, password2).
2. Returns `422` if `password` and `password2` do not match.
3. Hashes the password with bcrypt (cost 12).
4. Inserts the user into the `users` table.
5. Returns `409` if the username already exists (MySQL unique constraint violation).
6. Returns `201` on success.

---

### `HandleAPILogin(db) http.HandlerFunc`

**Route:** `POST /api/login`

1. Decodes the JSON request body (username, password).
2. Looks up the user in the database by username.
3. Returns `401` if the user is not found or the password does not match the stored bcrypt hash.
4. Creates a JWT token that expires in 24 hours, signed with the `JWT_SECRET` environment variable.
5. Checks whether the credentials appear in the [breach list](breach-detection.md) and sets `breached: true` in the response if so.
6. Returns `200` with the token and breach flag.

---

### `hashPassword(password string) (string, error)`

Wraps `bcrypt.GenerateFromPassword` with a configurable cost factor. The cost is set to `12` in production and overridden to `bcrypt.MinCost` in tests to keep tests fast.

### `verifyPassword(password, hash string) bool`

Wraps `bcrypt.CompareHashAndPassword`. Returns `true` if the plain-text password matches the stored hash.

### `sendJSON(w, code, message)`

A helper used throughout the handlers package to send a consistent JSON response. Sets `Content-Type: application/json`, writes the status code, and encodes an `AuthResponse`.

---

## Security Notes

- **bcrypt cost 12** — deliberately slow to make brute-force attacks expensive.
- **JWT secret from environment** — never hard-coded. The `jwtKey` variable reads from `JWT_SECRET` at startup.
- **No timing side-channel** — bcrypt comparison is constant-time; the handler returns the same `401` for both "user not found" and "wrong password", preventing username enumeration.
- **Passwords never logged** — the logging middleware records method, path, status, and duration, but never the request body.
