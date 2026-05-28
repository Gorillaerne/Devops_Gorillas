package handlers

import (
	"bytes"
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// AddPageHandler handles POST /api/pages.
// The DO scraper function calls this to push fetched Wikipedia content.
// Authenticated via X-Scraper-Key header (must match SCRAPER_KEY env var).
func AddPageHandler(db *sql.DB) http.HandlerFunc {
	scraperKey := os.Getenv("SCRAPER_KEY")
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkAPIKey(r.Header.Get("X-Scraper-Key"), scraperKey) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 2<<20)

		var page struct {
			Title    string `json:"title"`
			URL      string `json:"url"`
			Language string `json:"language"`
			Content  string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&page); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if page.Title == "" || page.URL == "" || page.Content == "" || page.Language == "" {
			http.Error(w, "title, url, content and language are required", http.StatusBadRequest)
			return
		}

		_, err := db.Exec(`
			INSERT INTO pages (title, url, language, content, last_updated)
			VALUES (?, ?, ?, ?, NOW())
			ON DUPLICATE KEY UPDATE
				title        = VALUES(title),
				content      = VALUES(content),
				last_updated = NOW()
		`, page.Title, page.URL, page.Language, page.Content)
		if err != nil {
			slog.Error("add page: db upsert failed", slog.Any("error", err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		slog.Info("page indexed via API", slog.String("title", page.Title), slog.String("language", page.Language))
		w.WriteHeader(http.StatusCreated)
	}
}

// TriggerScrapeHandler handles POST /api/scrape.
// Allows admins to trigger an on-demand scrape for a specific query by calling
// the DO scrape function directly. Authenticated via X-Scrape-Key header.
func TriggerScrapeHandler() http.HandlerFunc {
	scrapeKey := os.Getenv("SCRAPE_KEY")
	scrapeURL := os.Getenv("SCRAPE_FUNCTION_URL")
	functionKey := os.Getenv("FUNCTION_KEY")

	return func(w http.ResponseWriter, r *http.Request) {
		if !checkAPIKey(r.Header.Get("X-Scrape-Key"), scrapeKey) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var body struct {
			Query    string `json:"query"`
			Language string `json:"language"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if body.Query == "" {
			http.Error(w, "query is required", http.StatusBadRequest)
			return
		}
		if body.Language == "" {
			body.Language = "en"
		}
		if body.Language != "en" && body.Language != "da" {
			http.Error(w, "language must be 'en' or 'da'", http.StatusBadRequest)
			return
		}
		if scrapeURL == "" {
			http.Error(w, "scraper not configured", http.StatusServiceUnavailable)
			return
		}

		payload, _ := json.Marshal(map[string]interface{}{
			"query":    body.Query,
			"language": body.Language,
			"depth":    0,
		})

		req, err := http.NewRequest(http.MethodPost, scrapeURL, bytes.NewReader(payload))
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		if functionKey != "" {
			req.Header.Set("X-Function-Key", functionKey)
		}

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			slog.Error("trigger scrape: call failed", slog.Any("error", err))
			http.Error(w, "failed to trigger scrape", http.StatusInternalServerError)
			return
		}
		defer func() { _ = resp.Body.Close() }()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "scrape queued"})
	}
}

func checkAPIKey(incoming, expected string) bool {
	if expected == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(incoming), []byte(expected)) == 1
}
