# API Reference

All API endpoints are served under the `/api` prefix by the Go backend on port 8080. In production, Nginx forwards these requests from port 8081.

All responses use `Content-Type: application/json`.

---

## Authentication

Endpoints that require authentication expect a JWT token in the `Authorization` header:

```
Authorization: Bearer <token>
```

The token is issued by `POST /api/login` and expires after **24 hours**.

---

## Endpoints

### Search

#### `GET /api/search`

Searches the `pages` table for matching titles or content.

**Query parameters**

| Parameter | Required | Default | Description |
|---|---|---|---|
| `q` | Yes | — | Search term |
| `language` | No | `en` | Language filter (`en`, `da`, etc.) |

**Success response — 200 OK**

```json
{
  "data": [
    {
      "title": "Example Page",
      "content": "Some content about the example topic...",
      "URL": "https://example.com/page"
    }
  ]
}
```

Returns at most 20 results.

**Error response — 422 Unprocessable Entity**

Returned when `q` is missing.

```json
{
  "statusCode": 422,
  "message": "`q` query parameter is required"
}
```

---

### Authentication

#### `POST /api/register`

Creates a new user account.

**Request body**

```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "mypassword",
  "password2": "mypassword"
}
```

| Field | Required | Description |
|---|---|---|
| `username` | Yes | Must be unique |
| `email` | Yes | Used for breach notifications |
| `password` | Yes | Plain text — hashed server-side with bcrypt |
| `password2` | Yes | Must match `password` |

**Responses**

| Status | Message | Condition |
|---|---|---|
| 201 Created | `"User created successfully"` | Account created |
| 400 Bad Request | `"Invalid JSON body"` | Malformed request |
| 409 Conflict | `"User already exists or DB error"` | Username taken |
| 422 Unprocessable Entity | `"Passwords do not match"` | `password` ≠ `password2` |

---

#### `POST /api/login`

Authenticates a user and returns a JWT token.

**Request body**

```json
{
  "username": "johndoe",
  "password": "mypassword"
}
```

**Success response — 200 OK**

```json
{
  "statusCode": 200,
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "breached": false
}
```

When `breached` is `true`, the credentials were found in the known breach list. The frontend forces these users to the profile page to change their password.

**Error responses**

| Status | Message | Condition |
|---|---|---|
| 400 Bad Request | `"Invalid JSON body"` | Malformed request |
| 401 Unauthorized | `"Invalid credentials"` | Wrong username or password |

---

### Profile

#### `POST /api/change-password`

Changes the authenticated user's password. Requires a valid JWT token.

**Headers**

```
Authorization: Bearer <token>
```

**Request body**

```json
{
  "current_password": "oldpassword",
  "new_password": "newpassword",
  "new_password2": "newpassword"
}
```

**Responses**

| Status | Message | Condition |
|---|---|---|
| 200 OK | `"Password updated successfully"` | Password changed |
| 400 Bad Request | `"Invalid JSON body"` | Malformed request |
| 401 Unauthorized | `"Invalid or missing token"` | Bad/expired JWT |
| 401 Unauthorized | `"User not found"` | User deleted since token issued |
| 401 Unauthorized | `"Current password is incorrect"` | Wrong current password |
| 422 Unprocessable Entity | `"Passwords do not match"` | `new_password` ≠ `new_password2` |
| 500 Internal Server Error | `"Internal server error"` | Database error |

---

### Monitoring

#### `GET /metrics`

Exposes Prometheus metrics for scraping. This endpoint is served directly by the Go backend and is **not** proxied through Nginx.

See [Middleware & Metrics](backend/middleware.md) for the full list of metrics.

---

## Common Response Shape

Most non-search endpoints return this structure:

```json
{
  "statusCode": 200,
  "message": "Human readable message",
  "token": "...",
  "breached": false
}
```

`token` and `breached` are only present on a successful login response.
