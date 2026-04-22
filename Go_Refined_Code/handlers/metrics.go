package handlers

import (
	"database/sql"
	"log/slog"
	"time"

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
