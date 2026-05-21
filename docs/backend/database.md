# Database — database/sqlite.go

**File:** `Go_Refined_Code/database/sqlite.go`

Manages the MySQL database connection. Despite the filename (a remnant from an earlier SQLite phase), this file connects to MySQL in production and exposes a package-level `DB` variable used by all handlers.

---

## Package-Level Variable

```go
var DB *sql.DB
```

`DB` is the single shared database connection pool. All handler functions receive it as a parameter rather than accessing this global directly — this makes handlers independently testable with a separate in-memory SQLite connection during tests.

---

## `Connect() error`

Opens and verifies the MySQL connection. Called once from `main.go` before any routes are registered.

**Steps:**

1. Loads environment variables from a `.env` file if one exists. If no file is found, the function continues using system environment variables (so it works in Docker without a `.env` file).
2. Reads the `DATABASE_PATH` environment variable. This is a MySQL DSN string, e.g.:
   ```
   user:password@tcp(mysql:3306)/whoknowsdb?parseTime=true
   ```
3. Opens the connection pool with `sql.Open`. Note: this does not actually connect yet.
4. Calls `DB.Ping()` to verify the database is reachable. Returns an error if not.
5. Sets connection pool limits:
   - `MaxOpenConns`: 25 — maximum number of open connections to the database.
   - `MaxIdleConns`: 5 — connections kept open in the pool when idle.
   - `ConnMaxLifetime`: 5 minutes — connections are recycled after this time to avoid stale connections.
6. Lists all tables in the connected database and logs their names. This is a startup sanity check.

---

## `PurgeMD5Users()`

Called once from `main.go` immediately after `Connect()`.

```sql
DELETE FROM users WHERE password NOT LIKE '$2%'
```

bcrypt hashes always start with `$2` (e.g. `$2a$12$...`). Any password that does not match this pattern is a legacy MD5 hash from the original Python application. This function removes those users so the system only contains accounts with bcrypt passwords.

If any rows are deleted, the count is logged:
```
PurgeMD5Users: removed 3 legacy MD5 user(s)
```

---

## Connection Pool Settings Explained

| Setting | Value | Reason |
|---|---|---|
| `MaxOpenConns` | 25 | Limits load on the MySQL server; requests queue rather than opening unlimited connections |
| `MaxIdleConns` | 5 | Keeps a small pool of warm connections ready; reduces connection overhead |
| `ConnMaxLifetime` | 5 minutes | Prevents connections from going stale if MySQL drops them server-side |

---

## Test Setup

In tests (`handlers/*_test.go`), the handlers are initialised with an in-memory SQLite database rather than MySQL. This means:
- Tests run without needing a running MySQL server.
- Each test creates a fresh, isolated database.
- The `database` package itself is not used in tests.
