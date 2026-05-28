# main.go — Application Entry Point

**File:** `Go_Refined_Code/main.go`

This file is the starting point of the Go application. It wires together the database, middleware, routes, and HTTP server.

---

## What It Does

The `main()` function runs four setup steps in order:

### 1. Logging

Configures the global logger to emit structured JSON lines to stdout:

```go
slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
})))
```

Every log event (from any handler or middleware) is a JSON object with `time`, `level`, `msg`, and any extra fields. This makes logs easy to parse and query in production.

### 2. Database

```go
database.Connect()
database.PurgeMD5Users()
apiHandlers.StartUserMetricsCollector(database.DB, 60*time.Second)
apiHandlers.StartUserMetricsCollector(database.DB, 60*time.Second)
```

- Connects to MySQL using the `DATABASE_PATH` environment variable. Goose migrations are run automatically during `Connect()` (see [Database](database.md)).
- Deletes any legacy users whose password is not a bcrypt hash (see [Database](database.md)).
- Starts a background goroutine that counts registered users for Prometheus every 60 seconds.

Wikipedia content is populated by a separate DigitalOcean serverless function on a cron schedule — the server is not involved (see [Crawler & Scraper](scraper.md)).

If `SEND_BREACH_EMAILS=true` is set, it also kicks off a goroutine that emails all breached users (see [Email Service](email-service.md)).

### 3. Router

All HTTP routes are registered on a [Gorilla Mux](https://github.com/gorilla/mux) router:

| Method | Path | Handler |
|---|---|---|
| GET | `/metrics` | Prometheus metrics |
| GET | `/api/search` | `SearchAPIHandler` |
| GET | `/api/weather` | Placeholder stub |
| POST | `/api/register` | `HandleAPIRegister` |
| POST | `/api/login` | `HandleAPILogin` |
| GET | `/api/logout` | Placeholder stub |
| POST | `/api/change-password` | `HandleAPIChangePassword` |

All `/api/*` routes pass through two middleware layers (applied in order):
1. `PrometheusMiddleware` — records request count and duration metrics.
2. `LoggingMiddleware` — emits a structured log line per request.

CORS is applied at the outermost layer and allows all origins, methods, and headers.

### 4. HTTP Server

```go
srv := &http.Server{
    Addr:         ":8080",
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  120 * time.Second,
}
```

The server listens on port `8080`. Timeouts are set to protect against slow clients and idle connections.

---

## Middleware Order

```
Incoming request
    └── CORS middleware
        └── Gorilla Mux router
            └── /api/* subrouter
                └── PrometheusMiddleware
                    └── LoggingMiddleware
                        └── Handler function
```

---

## Key Design Choices

- **No global state for handlers.** Each handler receives the `*sql.DB` connection as a closure argument. This makes handlers testable without needing to mock globals.
- **Structured logging from the start.** Using `slog.NewJSONHandler` means every log line is already machine-readable, with no extra work for log aggregation tools.
- **Timeouts on the server.** Without `ReadTimeout` and `WriteTimeout`, a slow or malicious client could hold connections open indefinitely.
