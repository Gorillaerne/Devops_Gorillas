# Database — database/sqlite.go

**File:** `Go_Refined_Code/database/sqlite.go`

Manages the MySQL database connection and runs schema migrations on startup. Despite the filename (a remnant from an earlier SQLite phase), this file connects to MySQL in production and exposes a package-level `DB` variable used by all handlers.

---

## Package-Level Variable

```go
var DB *sql.DB
```

`DB` is the single shared database connection pool. All handler functions receive it as a parameter rather than accessing this global directly — this makes handlers independently testable with a separate in-memory SQLite connection during tests.

---

## `Connect() error`

Opens and verifies the MySQL connection, then runs any pending schema migrations. Called once from `main.go` before any routes are registered.

**Steps:**

1. Loads environment variables from a `.env` file if one exists. If no file is found, the function continues using system environment variables (so it works in Docker without a `.env` file).
2. Reads the `DATABASE_PATH` environment variable. This is a MySQL DSN string, e.g.:
   ```
   user:password@tcp(mysql:3306)/whoknowsdb?parseTime=true
   ```
3. Opens the connection pool with `sql.Open`. Note: this does not actually connect yet.
4. Calls `DB.Ping()` to verify the database is reachable. Returns an error if not.
5. Runs all pending Goose migrations (see [Migrations](#migrations) below).
6. Sets connection pool limits:
   - `MaxOpenConns`: 25 — maximum number of open connections to the database.
   - `MaxIdleConns`: 5 — connections kept open in the pool when idle.
   - `ConnMaxLifetime`: 5 minutes — connections are recycled after this time to avoid stale connections.
7. Lists all tables in the connected database and logs their names. This is a startup sanity check.

---

## Migrations

Schema changes are managed with [Goose v3](https://github.com/pressly/goose). Migration files live in `database/migrations/` and are embedded directly into the compiled binary using Go's `embed.FS` — no separate migration files are needed on the server.

```go
//go:embed migrations/*.sql
var migrationsFS embed.FS
```

On every startup, `goose.Up()` checks which migrations have already run (tracked in the `goose_db_version` table) and applies any new ones. Migrations that have already run are never re-applied.

### Migration Files

| File | Description |
|---|---|
| `00001_create_users.sql` | Creates the `users` table |
| `00002_create_pages.sql` | Creates the `pages` table |
| `00003_create_search_queries.sql` | Creates the `search_queries` table |
| `00004_seed_admin.go` | Inserts the admin user from environment variables |
| `00005_fix_collations.sql` | Aligns both tables to `utf8mb4_unicode_ci` |
| `00006_fulltext_search.sql` | Adds FULLTEXT indexes on `pages.title` and `pages.content` |

### Adding a New Migration

Create a new numbered SQL file in `database/migrations/`:

```sql
-- +goose Up
ALTER TABLE pages ADD COLUMN new_col VARCHAR(255);

-- +goose Down
ALTER TABLE pages DROP COLUMN new_col;
```

The `-- +goose Up` and `-- +goose Down` comments are required. The `Down` block is used if you ever need to roll back a migration manually.

### Admin Seed (`00004_seed_admin.go`)

The admin user is seeded via a Go migration (rather than SQL) so the password can be read from environment variables and bcrypt-hashed at migration time. The following variables must be set:

| Variable | Description |
|---|---|
| `ADMIN_USERNAME` | Admin account username |
| `ADMIN_EMAIL` | Admin account email |
| `ADMIN_PASSWORD` | Admin account plaintext password (hashed before storage) |

`INSERT IGNORE` is used so re-running the migration on an existing database does not overwrite the admin account.

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
