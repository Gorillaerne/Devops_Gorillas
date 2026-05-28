// Package handlers searchApi
package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
)

const descriptionMaxLen = 160

// SearchResult Struct
type SearchResult struct {
	Title       string `json:"title"`
	URL         string
	Content     string `json:"content"`
	Description string `json:"description"`
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

		results, err := executeSearch(db, q, language)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.Info("user_search", //nolint:gosec
			slog.String("query", q),
			slog.String("language", language),
			slog.Int("result_count", len(results)),
		)
		SearchQueriesTotal.WithLabelValues(q).Inc()

		_, dbErr := db.Exec(`
			INSERT INTO search_queries (query, language, count)
			VALUES (?, ?, 1)
			ON DUPLICATE KEY UPDATE count = count + 1
		`, q, language)
		if dbErr != nil {
			slog.Error("searchApi: failed to upsert search_queries", slog.Any("error", dbErr))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(SearchResponse{Data: results}); err != nil {
			slog.Error("searchApi: failed to encode response", slog.Any("error", err))
		}
	}
}

func executeSearch(db *sql.DB, q, language string) ([]SearchResult, error) {
	results, err := runFullTextSearch(db, q, language)
	if err != nil {
		slog.Info("searchApi: FULLTEXT unavailable, falling back to LIKE", slog.String("reason", err.Error()))
		return runLikeSearch(db, q, language)
	}
	if len(results) == 0 {
		return runLikeSearch(db, q, language)
	}
	return results, nil
}

func runFullTextSearch(db *sql.DB, q, language string) ([]SearchResult, error) {
	rows, err := db.Query(`
SELECT title, content, url
FROM pages
WHERE language = ?
  AND (MATCH(title) AGAINST(? IN NATURAL LANGUAGE MODE) > 0
       OR MATCH(content) AGAINST(? IN NATURAL LANGUAGE MODE) > 0)
ORDER BY MATCH(title) AGAINST(? IN NATURAL LANGUAGE MODE) * 3
       + MATCH(content) AGAINST(? IN NATURAL LANGUAGE MODE) DESC
LIMIT 20
`, language, q, q, q, q)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	return scanSearchRows(rows)
}

func runLikeSearch(db *sql.DB, q, language string) ([]SearchResult, error) {
	rows, err := db.Query(`
SELECT title, content, url
FROM pages
WHERE (title LIKE ? OR content LIKE ?)
  AND language = ?
LIMIT 20
`, "%"+q+"%", "%"+q+"%", language)
	if err != nil {
		slog.Error("searchApi: fallback query failed", //nolint:gosec
			slog.String("query", q),
			slog.String("language", language),
			slog.Any("error", err),
		)
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	return scanSearchRows(rows)
}

func scanSearchRows(rows *sql.Rows) ([]SearchResult, error) {
	var results []SearchResult
	for rows.Next() {
		var r SearchResult
		if err := rows.Scan(&r.Title, &r.Content, &r.URL); err != nil {
			slog.Error("searchApi: error scanning row", slog.Any("error", err))
			continue
		}
		r.Description = truncateToDescription(r.Content)
		results = append(results, r)
	}
	return results, rows.Err()
}

func truncateToDescription(content string) string {
	runes := []rune(content)
	if len(runes) > descriptionMaxLen {
		return string(runes[:descriptionMaxLen]) + "..."
	}
	return content
}
