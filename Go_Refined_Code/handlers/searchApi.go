// Package handlers searchApi
package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
)

// SearchResult Struct
type SearchResult struct {
	Title   string `json:"title"`
	URL     string
	Content string `json:"content"`
}

// SearchResponse Struct
type SearchResponse struct {
	Data []SearchResult `json:"data"`
}

// ErrorResponse Struct
type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

// SearchAPIHandler handles searching
func SearchAPIHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")

		if q == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)

			err := json.NewEncoder(w).Encode(ErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "`q` query parameter is required",
			})

			if err != nil {
				slog.Error("searchApi: failed to send error response", slog.Any("error", err))
			}
			return
		}

		language := r.URL.Query().Get("language")
		if language == "" {
			language = "en"
		}

		var results []SearchResult

		rows, err := db.Query(`
SELECT title, content, url
FROM pages
WHERE (title LIKE ? OR content LIKE ?)
  AND language = ?
LIMIT 20
`, "%"+q+"%", "%"+q+"%", language)
		if err != nil {
			slog.Error("searchApi: database query failed",
				slog.String("query", q),
				slog.String("language", language),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer func() {
			if closeErr := rows.Close(); closeErr != nil {
				slog.Error("searchApi: error closing rows", slog.Any("error", closeErr))
			}
		}()

		for rows.Next() {
			var result SearchResult
			if err := rows.Scan(&result.Title, &result.Content, &result.URL); err != nil {
				slog.Error("searchApi: error scanning row", slog.Any("error", err))
				continue
			}
			results = append(results, result)
		}

		if err := rows.Err(); err != nil {
			slog.Error("searchApi: error during row iteration", slog.Any("error", err))
		}

		// Log the user search event — structured so it can be queried later.
		// We intentionally do not log personal data (no user ID, no IP here).
		slog.Info("user_search",
			slog.String("query", q), //nolint:gosec
			slog.String("language", language),
			slog.Int("result_count", len(results)),
		)

		response := SearchResponse{Data: results}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err = json.NewEncoder(w).Encode(response); err != nil {
			slog.Error("searchApi: failed to encode response", slog.Any("error", err))
		}
	}
}
