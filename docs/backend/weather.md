# Weather — weatherApi.go

**File:** `Go_Refined_Code/handlers/weatherApi.go`

Handles `GET /api/weather`. Fetches a 7-day weather forecast from [Open-Meteo](https://open-meteo.com) and returns it wrapped in the `StandardResponse` shape defined by the OpenAPI spec (`Legacy/openApiSpec`).

---

## Why Open-Meteo?

- **Free** — no API key required, no billing setup.
- **Generous limits** — 10,000 requests per day on the free tier.
- **Open data** — based on publicly funded meteorological data.
- **Stable API** — has a published OpenAPI spec and versioned endpoints.

---

## Response Format

Matches the `StandardResponse` schema from the legacy OpenAPI spec:

```json
{
  "data": {
    "current": {
      "temperature_2m": 14.5,
      "windspeed_10m": 9.2,
      "weathercode": 2,
      "is_day": 1
    },
    "daily": {
      "time": ["2026-05-19", "2026-05-20", ...],
      "temperature_2m_max": [17.1, 15.8, ...],
      "temperature_2m_min": [10.2, 9.7, ...],
      "weathercode": [2, 61, ...]
    }
  }
}
```

The `data` object is the raw Open-Meteo response. The frontend interprets the WMO weather codes into human-readable descriptions and emoji icons.

---

## In-Memory Cache

The most important part of this handler is the 30-minute cache.

```
Request 1 → cache miss → fetch Open-Meteo → store in cache → respond
Request 2 → cache hit  → respond from cache (no upstream call)
...
Request N → cache hit  → respond from cache (no upstream call)
30 min later → cache expired
Request N+1 → cache miss → fetch Open-Meteo → store in cache → respond
```

### How it works

```go
type weatherCache struct {
    mu        sync.Mutex
    data      map[string]any
    fetchedAt time.Time
}

var globalWeatherCache weatherCache
```

The mutex ensures that concurrent requests don't trigger multiple simultaneous upstream calls. When the cache has valid data (`fetchedAt` is less than 30 minutes ago), the handler returns immediately without touching Open-Meteo.

### Scalability impact

Without cache: 1,000 users/minute = 1,000 Open-Meteo requests/minute (~1.4M/day — far over the free tier limit).

With 30-minute cache: 1,000 users/minute = 2 Open-Meteo requests/hour = 48/day. Well within limits at any traffic level.

---

## Open-Meteo Parameters

The handler fetches weather for **Copenhagen** (hardcoded coordinates):

| Parameter | Value |
|---|---|
| Latitude | 55.6761 |
| Longitude | 12.5683 |
| Timezone | Europe/Copenhagen |
| Forecast days | 7 |
| Current fields | `temperature_2m`, `windspeed_10m`, `weathercode`, `is_day` |
| Daily fields | `temperature_2m_max`, `temperature_2m_min`, `weathercode` |

---

## Functions

### `WeatherAPIHandler(w, r)`

The exported handler registered in `main.go` for `GET /api/weather`.

1. Acquires the cache mutex.
2. If the cache is valid, calls `writeWeatherJSON` and returns immediately.
3. Otherwise, builds an HTTP request to Open-Meteo using the request's context (so cancellations propagate correctly).
4. Returns `502 Bad Gateway` if Open-Meteo is unreachable or returns a non-200 status.
5. Decodes the JSON response, stores it in the cache with the current timestamp, and responds.

### `writeWeatherJSON(w, data)`

Private helper that wraps `data` in `{ "data": ... }` and writes it as JSON. Matches the `StandardResponse` schema.

---

## Testing — weather_test.go

Four tests using a local `httptest.Server` as a mock for Open-Meteo. The `openMeteoBaseURL` variable is overridden to point to the mock server, and `globalWeatherCache` is reset between tests.

| Test | What it verifies |
|---|---|
| `TestWeatherAPIHandler_Success` | Returns 200 with a `{ "data": ... }` response body |
| `TestWeatherAPIHandler_UsesCache` | Three requests trigger only one upstream call |
| `TestWeatherAPIHandler_UpstreamNon200` | A 503 from Open-Meteo results in a 502 to the client |
| `TestWeatherAPIHandler_DoesNotCacheOnError` | A failed fetch does not poison the cache — the next request retries |

---

## Error Handling

| Situation | Response |
|---|---|
| Failed to build request | `500 Internal Server Error` |
| Network error to Open-Meteo | `502 Bad Gateway` |
| Open-Meteo returns non-200 | `502 Bad Gateway` |
| Failed to decode JSON | `500 Internal Server Error` |
| Cache hit | `200 OK` (no upstream call) |
