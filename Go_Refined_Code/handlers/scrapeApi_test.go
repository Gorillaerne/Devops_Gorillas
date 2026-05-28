package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// --- checkAPIKey ---

func TestCheckAPIKey_EmptyExpected(t *testing.T) {
	if checkAPIKey("anything", "") {
		t.Error("expected false when expected key is empty")
	}
}

func TestCheckAPIKey_Match(t *testing.T) {
	if !checkAPIKey("secret", "secret") {
		t.Error("expected true for matching keys")
	}
}

func TestCheckAPIKey_Mismatch(t *testing.T) {
	if checkAPIKey("wrong", "secret") {
		t.Error("expected false for mismatched keys")
	}
}

// --- validateScrapeLanguage ---

func TestValidateScrapeLanguage_Empty(t *testing.T) {
	lang, ok := validateScrapeLanguage("")
	if !ok || lang != "en" {
		t.Errorf("expected ('en', true), got (%q, %v)", lang, ok)
	}
}

func TestValidateScrapeLanguage_En(t *testing.T) {
	lang, ok := validateScrapeLanguage("en")
	if !ok || lang != "en" {
		t.Errorf("expected ('en', true), got (%q, %v)", lang, ok)
	}
}

func TestValidateScrapeLanguage_Da(t *testing.T) {
	lang, ok := validateScrapeLanguage("da")
	if !ok || lang != "da" {
		t.Errorf("expected ('da', true), got (%q, %v)", lang, ok)
	}
}

func TestValidateScrapeLanguage_Invalid(t *testing.T) {
	_, ok := validateScrapeLanguage("fr")
	if ok {
		t.Error("expected false for unsupported language")
	}
}

// --- AddPageHandler ---

func TestAddPageHandler_Unauthorized(t *testing.T) {
	t.Setenv("SCRAPER_KEY", "secret")
	h := AddPageHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/pages", nil)
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAddPageHandler_InvalidJSON(t *testing.T) {
	t.Setenv("SCRAPER_KEY", "secret")
	h := AddPageHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/pages", strings.NewReader("not-json"))
	req.Header.Set("X-Scraper-Key", "secret")
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestAddPageHandler_MissingFields(t *testing.T) {
	t.Setenv("SCRAPER_KEY", "secret")
	h := AddPageHandler(nil)

	body := `{"title":"Only Title"}`
	req := httptest.NewRequest(http.MethodPost, "/api/pages", strings.NewReader(body))
	req.Header.Set("X-Scraper-Key", "secret")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- TriggerScrapeHandler ---

func TestTriggerScrapeHandler_Unauthorized(t *testing.T) {
	t.Setenv("SCRAPE_KEY", "secret")
	h := TriggerScrapeHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/scrape", nil)
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestTriggerScrapeHandler_InvalidJSON(t *testing.T) {
	t.Setenv("SCRAPE_KEY", "secret")
	h := TriggerScrapeHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/scrape", strings.NewReader("bad"))
	req.Header.Set("X-Scrape-Key", "secret")
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestTriggerScrapeHandler_MissingQuery(t *testing.T) {
	t.Setenv("SCRAPE_KEY", "secret")
	h := TriggerScrapeHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/scrape", strings.NewReader(`{"language":"en"}`))
	req.Header.Set("X-Scrape-Key", "secret")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestTriggerScrapeHandler_InvalidLanguage(t *testing.T) {
	t.Setenv("SCRAPE_KEY", "secret")
	h := TriggerScrapeHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/scrape", strings.NewReader(`{"query":"Python","language":"fr"}`))
	req.Header.Set("X-Scrape-Key", "secret")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestTriggerScrapeHandler_ScraperNotConfigured(t *testing.T) {
	t.Setenv("SCRAPE_KEY", "secret")
	t.Setenv("SCRAPE_FUNCTION_URL", "")
	h := TriggerScrapeHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/scrape", strings.NewReader(`{"query":"Python","language":"en"}`))
	req.Header.Set("X-Scrape-Key", "secret")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", w.Code)
	}
}

func TestTriggerScrapeHandler_Success(t *testing.T) {
	mockScraper := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockScraper.Close()

	t.Setenv("SCRAPE_KEY", "secret")
	t.Setenv("SCRAPE_FUNCTION_URL", mockScraper.URL)
	t.Setenv("FUNCTION_KEY", "fn-key")
	h := TriggerScrapeHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/scrape", strings.NewReader(`{"query":"Python","language":"en"}`))
	req.Header.Set("X-Scrape-Key", "secret")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusAccepted {
		t.Errorf("expected 202, got %d", w.Code)
	}
}

func TestTriggerScrapeHandler_DefaultLanguage(t *testing.T) {
	mockScraper := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockScraper.Close()

	t.Setenv("SCRAPE_KEY", "secret")
	t.Setenv("SCRAPE_FUNCTION_URL", mockScraper.URL)
	h := TriggerScrapeHandler()

	// No language field — should default to "en" and succeed
	req := httptest.NewRequest(http.MethodPost, "/api/scrape", strings.NewReader(`{"query":"Python"}`))
	req.Header.Set("X-Scrape-Key", "secret")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusAccepted {
		t.Errorf("expected 202, got %d", w.Code)
	}
}

func TestTriggerScrapeHandler_CallFailed(t *testing.T) {
	// Close the server before using its URL so client.Do returns a connection error.
	ts := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	ts.Close()

	t.Setenv("SCRAPE_KEY", "secret")
	t.Setenv("SCRAPE_FUNCTION_URL", ts.URL)
	h := TriggerScrapeHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/scrape", strings.NewReader(`{"query":"test","language":"en"}`))
	req.Header.Set("X-Scrape-Key", "secret")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestAddPageHandler_UpsertError(t *testing.T) {
	// SQLite does not support MySQL's ON DUPLICATE KEY UPDATE syntax,
	// so the upsert will fail and the handler should return 500.
	db := newPagesDB(t)
	t.Setenv("SCRAPER_KEY", "secret")
	h := AddPageHandler(db)

	body := `{"title":"T","url":"https://example.com","language":"en","content":"some content"}`
	req := httptest.NewRequest(http.MethodPost, "/api/pages", strings.NewReader(body))
	req.Header.Set("X-Scraper-Key", "secret")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestCallScrapeFunction_InvalidURL(t *testing.T) {
	// "http://[" is an invalid URL — http.NewRequest will fail to parse it.
	err := callScrapeFunction("http://[", "", "test", "en")
	if err == nil {
		t.Error("expected error for invalid URL")
	}
}
