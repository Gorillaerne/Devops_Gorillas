# Middleware & Metrics — loggingMiddleware.go & metrics.go

**Files:**
- `Go_Refined_Code/handlers/loggingMiddleware.go`
- `Go_Refined_Code/handlers/metrics.go`

Both files implement HTTP middleware — code that wraps every API request to add cross-cutting behaviour (logging, metrics) without touching the handler logic.

---

## responseWriter — Shared Helper

Both middleware files rely on a thin wrapper around `http.ResponseWriter` that captures the HTTP status code written by the handler:

```go
type responseWriter struct {
    http.ResponseWriter
    statusCode int
}
```

The default status is `200`. When the handler calls `WriteHeader(code)`, the wrapper stores the code before passing it through. This allows middleware to read the final status code after the handler has run.

---

## Logging Middleware — `LoggingMiddleware`

**File:** `loggingMiddleware.go`

Wraps every `/api/*` request and emits a single structured JSON log line after the handler completes.

### Logged Fields

| Field | Type | Example |
|---|---|---|
| `method` | string | `"POST"` |
| `path` | string | `"/api/login"` |
| `status` | int | `200` |
| `duration_ms` | int64 | `12` |
| `remote_ip` | string | `"172.18.0.1:54312"` |

### Log Level by Status Code

| Status range | Log level |
|---|---|
| 2xx / 3xx | `INFO` |
| 4xx | `WARN` |
| 5xx | `ERROR` |

### What Is Never Logged

Request bodies, passwords, and tokens are never included. The middleware only logs metadata about the request, not its content.

---

## Prometheus Middleware — `PrometheusMiddleware`

**File:** `metrics.go`

Records two metrics for every `/api/*` request. Uses the Gorilla Mux **route pattern** (e.g. `/api/search`) rather than the raw URL as the `path` label, so dynamic segments (like query strings) do not create unbounded label cardinality.

### Metrics Recorded Per Request

| Metric | Type | Labels | Description |
|---|---|---|---|
| `http_requests_total` | Counter | `method`, `path`, `status` | Total number of requests |
| `http_request_duration_seconds` | Histogram | `method`, `path`, `status` | Request duration in seconds |

---

## Other Prometheus Metrics

### `search_queries_total`

- **Type:** Counter
- **Label:** `query` (the search term)
- **Incremented by:** `SearchAPIHandler` on every successful search.
- **Use case:** Identify the most popular search terms (e.g. in a Grafana `topk` panel).

### `registered_users_total`

- **Type:** Gauge
- **No labels**
- **Updated by:** `StartUserMetricsCollector` — a background goroutine that runs a `SELECT COUNT(*) FROM users` query every 60 seconds and sets the gauge.

---

## `StartUserMetricsCollector(db, interval)`

Called once from `main.go` after the database is ready. Starts a goroutine that periodically updates the `registered_users_total` gauge. Errors are logged but do not stop the loop.

---

## Middleware Registration

Both middleware layers are applied to the `/api` subrouter in `main.go`:

```go
api := r.PathPrefix("/api").Subrouter()
api.Use(apiHandlers.PrometheusMiddleware)
api.Use(apiHandlers.LoggingMiddleware)
```

Prometheus runs first so it starts timing before the logging middleware adds any overhead.
