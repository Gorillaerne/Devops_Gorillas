package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "modernc.org/sqlite"
)

func newPagesDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE pages (
		id       INTEGER PRIMARY KEY AUTOINCREMENT,
		title    TEXT    NOT NULL,
		content  TEXT    NOT NULL,
		url      TEXT    NOT NULL,
		language TEXT    NOT NULL DEFAULT 'en'
	)`)
	if err != nil {
		t.Fatalf("create pages table: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func seedPage(t *testing.T, db *sql.DB, title, content, url, language string) {
	t.Helper()
	_, err := db.Exec(
		`INSERT INTO pages (title, content, url, language) VALUES (?, ?, ?, ?)`,
		title, content, url, language,
	)
	if err != nil {
		t.Fatalf("seed page: %v", err)
	}
}

func TestSearchAPIHandler_MissingQuery(t *testing.T) {
	db := newPagesDB(t)

	req := httptest.NewRequest(http.MethodGet, "/api/search", nil)
	w := httptest.NewRecorder()

	SearchAPIHandler(db)(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}

func TestSearchAPIHandler_ReturnsResults(t *testing.T) {
	db := newPagesDB(t)
	seedPage(t, db, "Go Programming", "Go is a statically typed language", "https://go.dev", "en")

	req := httptest.NewRequest(http.MethodGet, "/api/search?q=Go", nil)
	w := httptest.NewRecorder()

	SearchAPIHandler(db)(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp SearchResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Data) == 0 {
		t.Error("expected at least one search result")
	}
}

func TestSearchAPIHandler_EmptyResults(t *testing.T) {
	db := newPagesDB(t)

	req := httptest.NewRequest(http.MethodGet, "/api/search?q=nonexistent", nil)
	w := httptest.NewRecorder()

	SearchAPIHandler(db)(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 with empty results, got %d", w.Code)
	}
}

func TestSearchAPIHandler_DefaultLanguageIsEn(t *testing.T) {
	db := newPagesDB(t)
	seedPage(t, db, "English page", "about something", "https://example.com", "en")
	seedPage(t, db, "Danish page", "about something", "https://example.dk", "da")

	// No language param → defaults to "en", so only the English page should match
	req := httptest.NewRequest(http.MethodGet, "/api/search?q=something", nil)
	w := httptest.NewRecorder()

	SearchAPIHandler(db)(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp SearchResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Errorf("expected 1 result (only 'en'), got %d", len(resp.Data))
	}
	if len(resp.Data) > 0 && resp.Data[0].Title != "English page" {
		t.Errorf("expected English page, got %q", resp.Data[0].Title)
	}
}

func TestSearchAPIHandler_LanguageFilter(t *testing.T) {
	db := newPagesDB(t)
	seedPage(t, db, "Danish page", "Danish content", "https://example.dk", "da")
	seedPage(t, db, "English page", "Danish content", "https://example.com", "en")

	req := httptest.NewRequest(http.MethodGet, "/api/search?q=Danish&language=da", nil)
	w := httptest.NewRecorder()

	SearchAPIHandler(db)(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp SearchResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Errorf("expected 1 result (only 'da'), got %d", len(resp.Data))
	}
}
