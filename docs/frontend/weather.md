# Weather Page — weather.html & weather_script.js

**Files:**
- `Go_Refined_Code/static/html/weather.html`
- `Go_Refined_Code/static/javaScript/weather_script.js`

The weather page displays the current conditions and a 7-day forecast for Copenhagen, fetched from the Go backend which proxies and caches Open-Meteo data.

---

## weather.html

**URL:** `/weather`

A static HTML page that loads `weather_script.js` to populate two sections:

- **`#weather-current`** — current conditions card (temperature, wind, condition).
- **`#weather-forecast`** — 7-day forecast grid, one card per day.

The page shows a "Loading weather…" placeholder until the API responds.

---

## weather_script.js

Calls `GET /api/weather` on page load and renders the response into the DOM.

### WMO Weather Code Interpretation

Open-Meteo returns WMO weather interpretation codes (integers) rather than text descriptions. The script maps these to human-readable strings and emoji icons:

| Code range | Condition | Icon |
|---|---|---|
| 0 | Clear sky | ☀️ |
| 1–3 | Clear / partly cloudy / overcast | 🌤️ ⛅ ☁️ |
| 45, 48 | Fog | 🌫️ |
| 51–55 | Drizzle | 🌦️ |
| 61–65 | Rain | 🌧️ |
| 71–77 | Snow | 🌨️ ❄️ |
| 80–82 | Rain showers | 🌦️ ⛈️ |
| 95–99 | Thunderstorm | ⛈️ |

### `renderCurrent(current)`

Reads `current.temperature_2m`, `current.windspeed_10m`, and `current.weathercode` from the API response and builds the current-conditions card.

### `renderForecast(daily)`

Reads `daily.time`, `daily.temperature_2m_max`, `daily.temperature_2m_min`, and `daily.weathercode` (all arrays with one entry per day) and builds one card per day in a CSS grid.

### `formatDate(dateStr)`

Converts an ISO date string (`"2026-05-19"`) to a short, readable label (`"Mon, 19 May"`) using the browser's `Intl`-based `toLocaleDateString`.

### Error handling

If the fetch fails or the API returns an error status, `showError()` from `reuseable_functions.js` shows a toast notification and a fallback "Weather unavailable." message replaces the loading placeholder.

---

## Data Flow

```
weather.html loads
    └── weather_script.js runs
        └── fetch("/api/weather")
            └── Nginx proxies to Go backend
                └── WeatherAPIHandler
                    ├── Cache hit  → return cached data
                    └── Cache miss → fetch Open-Meteo → cache → return
            └── { "data": { current: {...}, daily: {...} } }
        └── renderCurrent(data.current)
        └── renderForecast(data.daily)
```
