package main

import (
	"bytes"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	scrapeMaxDepth = 1
	scrapeMaxLinks = 5
	scrapeTimeout  = 10 * time.Second
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

type mediaWikiLinksResponse struct {
	Query struct {
		Pages map[string]struct {
			Links []struct {
				Title string `json:"title"`
			} `json:"links"`
		} `json:"pages"`
	} `json:"query"`
}

// Main is the DigitalOcean Functions entry point for the HTTP-triggered scraper.
// It receives {query, language, depth}, fetches the Wikipedia page, posts the
// content to POST /api/pages on the Go server, then follows links by calling
// itself recursively (up to scrapeMaxDepth hops). It never touches the database.
func Main(args map[string]interface{}) map[string]interface{} {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))

	functionKey := os.Getenv("FUNCTION_KEY")
	if functionKey != "" {
		headers, _ := args["__ow_headers"].(map[string]interface{})
		incoming, _ := headers["x-function-key"].(string)
		if subtle.ConstantTimeCompare([]byte(incoming), []byte(functionKey)) != 1 {
			return webRespond(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
	}

	query, _ := args["query"].(string)
	language, _ := args["language"].(string)
	depth := 0
	if d, ok := args["depth"].(float64); ok {
		depth = int(d)
	}

	if query == "" || language == "" {
		return webRespond(http.StatusBadRequest, map[string]string{"error": "query and language are required"})
	}

	apiURL := os.Getenv("API_URL")
	scraperKey := os.Getenv("SCRAPER_KEY")
	scrapeURL := os.Getenv("SCRAPE_FUNCTION_URL")

	if apiURL == "" || scraperKey == "" {
		return webRespond(http.StatusInternalServerError, map[string]string{"error": "API_URL or SCRAPER_KEY not configured"})
	}

	client := &http.Client{Timeout: scrapeTimeout}

	summary, err := fetchWikipedia(client, query, language)
	if err != nil {
		slog.Warn("wikipedia fetch failed", slog.String("query", query), slog.Any("error", err))
		return webRespond(http.StatusOK, map[string]string{"status": "skipped", "reason": err.Error()})
	}

	if err := postPage(client, apiURL, scraperKey, summary, language); err != nil {
		slog.Error("post page failed", slog.Any("error", err))
		return webRespond(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	slog.Info("indexed page", slog.String("title", summary.Title), slog.String("language", language), slog.Int("depth", depth))

	if depth < scrapeMaxDepth && scrapeURL != "" {
		links, err := fetchWikipediaLinks(client, query, language)
		if err != nil {
			slog.Warn("links fetch failed", slog.String("query", query), slog.Any("error", err))
		} else {
			for _, link := range links {
				callScrape(client, scrapeURL, functionKey, link, language, depth+1)
			}
		}
	}

	return webRespond(http.StatusOK, map[string]string{"status": "ok", "title": summary.Title})
}

func postPage(client *http.Client, apiURL, scraperKey string, s *wikipediaSummary, language string) error {
	payload, err := json.Marshal(map[string]string{
		"title":    s.Title,
		"url":      s.URL.Desktop.Page,
		"language": language,
		"content":  s.Extract,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, apiURL+"/api/pages", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Scraper-Key", scraperKey)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("api/pages returned status %d", resp.StatusCode)
	}
	return nil
}

func callScrape(client *http.Client, scrapeURL, functionKey, query, language string, depth int) {
	payload, _ := json.Marshal(map[string]interface{}{
		"query":    query,
		"language": language,
		"depth":    depth,
	})

	req, err := http.NewRequest(http.MethodPost, scrapeURL, bytes.NewReader(payload))
	if err != nil {
		slog.Warn("callScrape: build request failed", slog.String("query", query), slog.Any("error", err))
		return
	}
	req.Header.Set("Content-Type", "application/json")
	if functionKey != "" {
		req.Header.Set("X-Function-Key", functionKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		slog.Warn("callScrape: call failed", slog.String("query", query), slog.Any("error", err))
		return
	}
	defer func() { _ = resp.Body.Close() }()
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

func fetchWikipediaLinks(client *http.Client, title, language string) ([]string, error) {
	endpoint := fmt.Sprintf(
		"https://%s.wikipedia.org/w/api.php?action=query&titles=%s&prop=links&pllimit=%d&plnamespace=0&format=json",
		language,
		url.QueryEscape(title),
		scrapeMaxLinks,
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("mediawiki api returned status %d", resp.StatusCode)
	}

	var result mediaWikiLinksResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var links []string
	for _, page := range result.Query.Pages {
		for _, l := range page.Links {
			links = append(links, l.Title)
		}
	}
	return links, nil
}

// webRespond formats a response for a DO Functions web action.
func webRespond(status int, body interface{}) map[string]interface{} {
	return map[string]interface{}{
		"statusCode": status,
		"headers":    map[string]string{"Content-Type": "application/json"},
		"body":       body,
	}
}
