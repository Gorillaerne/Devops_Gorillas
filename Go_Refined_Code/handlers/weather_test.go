package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func resetWeatherCache() {
	globalWeatherCache.mu.Lock()
	globalWeatherCache.data = nil
	globalWeatherCache.fetchedAt = time.Time{}
	globalWeatherCache.mu.Unlock()
}

func TestWeatherAPIHandler_Success(t *testing.T) {
	resetWeatherCache()

	mockPayload := map[string]any{
		"current": map[string]any{
			"temperature_2m": 15.0,
			"windspeed_10m":  8.5,
			"weathercode":    1.0,
			"is_day":         1.0,
		},
		"daily": map[string]any{
			"time":               []any{"2026-05-19"},
			"temperature_2m_max": []any{17.0},
			"temperature_2m_min": []any{10.0},
			"weathercode":        []any{1.0},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(mockPayload); err != nil {
			t.Errorf("mock server encode: %v", err)
		}
	}))
	defer srv.Close()

	openMeteoBaseURL = srv.URL
	t.Cleanup(func() {
		openMeteoBaseURL = "https://api.open-meteo.com"
		resetWeatherCache()
	})

	req := httptest.NewRequest(http.MethodGet, "/api/weather", nil)
	rr := httptest.NewRecorder()
	WeatherAPIHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var resp map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if _, ok := resp["data"]; !ok {
		t.Error("expected 'data' key in response matching StandardResponse schema")
	}
}

func TestWeatherAPIHandler_UsesCache(t *testing.T) {
	resetWeatherCache()

	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		if err := json.NewEncoder(w).Encode(map[string]any{"cached": true}); err != nil {
			t.Errorf("mock server encode: %v", err)
		}
	}))
	defer srv.Close()

	openMeteoBaseURL = srv.URL
	t.Cleanup(func() {
		openMeteoBaseURL = "https://api.open-meteo.com"
		resetWeatherCache()
	})

	for i := range 3 {
		req := httptest.NewRequest(http.MethodGet, "/api/weather", nil)
		rr := httptest.NewRecorder()
		WeatherAPIHandler(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("request %d: expected 200, got %d", i, rr.Code)
		}
	}

	if callCount != 1 {
		t.Errorf("expected exactly 1 upstream call (cache should serve requests 2 and 3), got %d", callCount)
	}
}

func TestWeatherAPIHandler_UpstreamNon200(t *testing.T) {
	resetWeatherCache()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	openMeteoBaseURL = srv.URL
	t.Cleanup(func() {
		openMeteoBaseURL = "https://api.open-meteo.com"
		resetWeatherCache()
	})

	req := httptest.NewRequest(http.MethodGet, "/api/weather", nil)
	rr := httptest.NewRecorder()
	WeatherAPIHandler(rr, req)

	if rr.Code != http.StatusBadGateway {
		t.Errorf("expected 502 Bad Gateway on upstream error, got %d", rr.Code)
	}
}

func TestWeatherAPIHandler_DoesNotCacheOnError(t *testing.T) {
	resetWeatherCache()

	firstCall := true
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if firstCall {
			firstCall = false
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		if err := json.NewEncoder(w).Encode(map[string]any{"recovered": true}); err != nil {
			t.Errorf("mock server encode: %v", err)
		}
	}))
	defer srv.Close()

	openMeteoBaseURL = srv.URL
	t.Cleanup(func() {
		openMeteoBaseURL = "https://api.open-meteo.com"
		resetWeatherCache()
	})

	req1 := httptest.NewRequest(http.MethodGet, "/api/weather", nil)
	rr1 := httptest.NewRecorder()
	WeatherAPIHandler(rr1, req1)
	if rr1.Code != http.StatusBadGateway {
		t.Errorf("first request: expected 502, got %d", rr1.Code)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/api/weather", nil)
	rr2 := httptest.NewRecorder()
	WeatherAPIHandler(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Errorf("second request: expected 200 after recovery, got %d", rr2.Code)
	}
}
