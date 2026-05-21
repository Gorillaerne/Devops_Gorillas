# Profile — profileApi.go

**File:** `Go_Refined_Code/handlers/profileApi.go`

Handles the `POST /api/change-password` endpoint. Requires the user to be authenticated with a valid JWT token.

---

## Functions

### `parseJWT(r *http.Request) (*Claims, error)`

A private helper that:
1. Reads the `Authorization` header from the request.
2. Strips the `Bearer ` prefix to extract the raw token string.
3. Parses and validates the JWT using the `JWT_SECRET` key.
4. Returns the decoded `Claims` (which contain the `UserID`) or an error if the token is invalid or expired.

---

### `HandleAPIChangePassword(db) http.HandlerFunc`

**Route:** `POST /api/change-password`

**Requires:** `Authorization: Bearer <token>` header.

### Flow

1. Calls `parseJWT` to validate the token. Returns `401` if invalid.
2. Decodes the JSON request body:
   - `current_password`
   - `new_password`
   - `new_password2`
3. Returns `422` if `new_password` and `new_password2` do not match.
4. Fetches the current bcrypt hash from the database for the user identified by the JWT's `UserID`.
5. Verifies that `current_password` matches the stored hash. Returns `401` if it does not.
6. Hashes the new password with bcrypt.
7. Updates the `password` column in the database.
8. Returns `200` on success.

---

## Why the Current Password Is Required

Requiring the current password prevents an attacker who has stolen a JWT token from immediately locking out the real user. Even with a valid token, they cannot change the password without also knowing what it currently is.

---

## JWT Claims

The `Claims` type (defined in `authApi.go`) carries the `UserID`, which is the primary key in the `users` table. The handler uses this ID directly in the database query rather than looking up the user by username, making it safe against username changes.
