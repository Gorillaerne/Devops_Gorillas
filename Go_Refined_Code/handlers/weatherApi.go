// Package handlers weatherApi
package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

const (
	weatherCacheTTL = 30 * time.Minute
	openMeteoPath   = "/v1/forecast" +
		"?latitude=55.6761&longitude=12.5683" +
		"&current=temperature_2m,windspeed_10m,weathercode,is_day" +
		"&daily=temperature_2m_max,temperature_2m_min,weathercode" +
		"&timezone=Europe%2FCopenhagen" +
		"&forecast_days=7"
)

// openMeteoBaseURL is the Open-Meteo API base. Overridden in tests.
var openMeteoBaseURL = "https://api.open-meteo.com" //nolint:gochecknoglobals

type weatherCache struct {
	mu        sync.Mutex
	data      map[string]any
	fetchedAt time.Time
}

var globalWeatherCache weatherCache //nolint:gochecknoglobals

// WeatherAPIHandler handles GET /api/weather.
// It fetches a 7-day forecast from Open-Meteo (free, no API key required)
// and caches the result for 30 minutes so all users share a single upstream call.
func WeatherAPIHandler(w http.ResponseWriter, r *http.Request) {
	globalWeatherCache.mu.Lock()
	defer globalWeatherCache.mu.Unlock()

	if globalWeatherCache.data != nil && time.Since(globalWeatherCache.fetchedAt) < weatherCacheTTL {
		writeWeatherJSON(w, globalWeatherCache.data)
		return
	}

	url := openMeteoBaseURL + openMeteoPath //nolint:gosec

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, url, nil)
	if err != nil {
		slog.Error("weatherApi: failed to build request", slog.Any("error", err))
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("weatherApi: failed to fetch from Open-Meteo", slog.Any("error", err))
		http.Error(w, "Failed to fetch weather data", http.StatusBadGateway)
		return
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			slog.Error("weatherApi: failed to close response body", slog.Any("error", closeErr))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		slog.Error("weatherApi: Open-Meteo returned non-200", slog.Int("status", resp.StatusCode))
		http.Error(w, "Weather service unavailable", http.StatusBadGateway)
		return
	}

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		slog.Error("weatherApi: failed to decode response", slog.Any("error", err))
		http.Error(w, "Failed to parse weather data", http.StatusInternalServerError)
		return
	}

	globalWeatherCache.data = data
	globalWeatherCache.fetchedAt = time.Now()

	writeWeatherJSON(w, data)
}

func writeWeatherJSON(w http.ResponseWriter, data map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]any{"data": data}); err != nil {
		slog.Error("weatherApi: failed to encode response", slog.Any("error", err))
	}
}
