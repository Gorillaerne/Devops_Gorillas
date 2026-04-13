// Package handlers searchApi
package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
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
				// If the client disconnected or the network failed, we log it here.
				log.Printf("searchApi: failed to send error response: %v", err)
			}
			return
		}

		language := r.URL.Query().Get("language")

		if language == "" {
			language = "en"
		}

		log.Println("API search query:", q) //nolint:gosec

		var results []SearchResult

		rows, err := db.Query(`
SELECT title, content, url
FROM pages
WHERE (title LIKE ? OR content LIKE ?)
  AND language = ?
LIMIT 20
`, "%"+q+"%", "%"+q+"%", language)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer func() {
			if closeErr := rows.Close(); closeErr != nil {
				log.Printf("searchApi: error closing rows: %v", closeErr)
			}
		}()

		for rows.Next() {
			var result SearchResult
			// Check the error returned by Scan
			if err := rows.Scan(&result.Title, &result.Content, &result.URL); err != nil {
				log.Printf("searchApi: error scanning row: %v", err)
				continue // Skip this specific malformed row and try the next one
			}
			results = append(results, result)
		}

		// ALWAYS check for errors that may have occurred during the iteration
		if err := rows.Err(); err != nil {
			log.Printf("searchApi: error during row iteration: %v", err)
		}

		response := SearchResponse{
			Data: results,
		}
		log.Println("WRAPPED RESPONSE EXECUTED")
		// JSON response
		w.Header().Set("Content-Type", "application/json")
		// Explicitly set 200 OK if we've reached this point successfully
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			// If we fail here, the client likely won't see the error message,
			// but your server logs will explain why the connection was cut.
			log.Printf("searchApi: failed to encode response: %v", err)
		}
	}
}
