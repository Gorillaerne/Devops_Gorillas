package handlers

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// SearchQueriesTotal counts search queries by term, used for Grafana topk dashboards.
var SearchQueriesTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "search_queries_total",
		Help: "Total number of search queries, labelled by the search term.",
	},
	[]string{"query"},
)

var registeredUsersTotal = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "registered_users_total",
	Help: "Current total number of registered users.",
})

var httpRequestsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests by method, path, and status.",
	},
	[]string{"method", "path", "status"},
)

var httpRequestDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration in seconds by method, path, and status.",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"method", "path", "status"},
)

// PrometheusMiddleware records http_requests_total and http_request_duration_seconds
// for every request. Uses the Gorilla Mux route pattern as the path label to avoid
// high-cardinality labels from dynamic URL segments.
func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := newResponseWriter(w)

		next.ServeHTTP(wrapped, r)

		// Use the matched route pattern (e.g. /api/search) not the raw URL.
		route := "unknown"
		if match := mux.CurrentRoute(r); match != nil {
			if pattern, err := match.GetPathTemplate(); err == nil {
				route = pattern
			}
		}

		status := fmt.Sprintf("%d", wrapped.statusCode)
		duration := time.Since(start).Seconds()

		httpRequestsTotal.WithLabelValues(r.Method, route, status).Inc()
		httpRequestDuration.WithLabelValues(r.Method, route, status).Observe(duration)
	})
}

// StartUserMetricsCollector runs a background goroutine that refreshes
// DB-backed gauges every interval. Call once from main after DB is ready.
func StartUserMetricsCollector(db *sql.DB, interval time.Duration) {
	go func() {
		for {
			var count float64
			if err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
				slog.Error("metrics: failed to count users", slog.Any("error", err))
			} else {
				registeredUsersTotal.Set(count)
			}
			time.Sleep(interval)
		}
	}()
}
