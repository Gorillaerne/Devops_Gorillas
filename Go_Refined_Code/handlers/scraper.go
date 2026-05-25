// Package handlers scraper
package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

const (
	scraperTopN    = 20
	scraperTimeout = 10 * time.Second
)

type wikipediaSummary struct {
	Title   string `json:"title"`
	Extract string `json:"extract"`
	URL     struct {
		Desktop struct {
			Page string `json:"page"`
		} `json:"desktop"`
	} `json:"content_urls"`
}

// StartScraper starts a background goroutine that periodically fetches Wikipedia
// pages for the most searched queries and upserts them into the pages table.
func StartScraper(db *sql.DB, interval time.Duration) {
	go func() {
		slog.Info("scraper: started", slog.Duration("interval", interval))
		runScrape(db)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			runScrape(db)
		}
	}()
}

func runScrape(db *sql.DB) {
	queries, err := topQueries(db)
	if err != nil {
		slog.Error("scraper: failed to fetch top queries", slog.Any("error", err))
		return
	}

	client := &http.Client{Timeout: scraperTimeout}

	for _, q := range queries {
		summary, err := fetchWikipedia(client, q.query, q.language)
		if err != nil {
			slog.Warn("scraper: wikipedia fetch failed",
				slog.String("query", q.query),
				slog.String("language", q.language),
				slog.Any("error", err),
			)
			continue
		}

		if err := upsertPage(db, summary, q.language); err != nil {
			slog.Error("scraper: failed to upsert page",
				slog.String("url", summary.URL.Desktop.Page),
				slog.Any("error", err),
			)
			continue
		}

		slog.Info("scraper: upserted page",
			slog.String("title", summary.Title),
			slog.String("language", q.language),
		)
	}
}

type queryRow struct {
	query    string
	language string
}

func topQueries(db *sql.DB) ([]queryRow, error) {
	rows, err := db.Query(`
		SELECT query, language
		FROM search_queries
		ORDER BY count DESC
		LIMIT ?
	`, scraperTopN)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []queryRow
	for rows.Next() {
		var q queryRow
		if err := rows.Scan(&q.query, &q.language); err != nil {
			return nil, err
		}
		results = append(results, q)
	}
	return results, rows.Err()
}

func fetchWikipedia(client *http.Client, query, language string) (*wikipediaSummary, error) {
	endpoint := fmt.Sprintf(
		"https://%s.wikipedia.org/api/rest_v1/page/summary/%s",
		language,
		url.PathEscape(query),
	)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "WhoKnows/1.0 (https://github.com/Gorillaerne/Devops_Gorillas)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no wikipedia article for %q", query)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wikipedia returned status %d", resp.StatusCode)
	}

	var summary wikipediaSummary
	if err := json.NewDecoder(resp.Body).Decode(&summary); err != nil {
		return nil, err
	}

	if summary.Extract == "" || summary.URL.Desktop.Page == "" {
		return nil, fmt.Errorf("wikipedia article for %q has no usable content", query)
	}

	return &summary, nil
}

func upsertPage(db *sql.DB, s *wikipediaSummary, language string) error {
	_, err := db.Exec(`
		INSERT INTO pages (title, url, language, content, last_updated)
		VALUES (?, ?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE
			title        = VALUES(title),
			content      = VALUES(content),
			last_updated = NOW()
	`, s.Title, s.URL.Desktop.Page, language, s.Extract)
	return err
}
