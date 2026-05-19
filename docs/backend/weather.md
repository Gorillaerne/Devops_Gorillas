# Weather — weatherApi.go

**File:** `Go_Refined_Code/handlers/weatherApi.go`

This file is a placeholder for a planned weather feature. The handler is currently empty — the route `GET /api/weather` is registered in `main.go` and responds, but does not return any data.

---

## Current State

The route is mapped to `homeHandler` in `main.go`:

```go
api.HandleFunc("/weather", homeHandler).Methods("GET")
```

`homeHandler` simply writes the text `"test endpoints"` to the response. No actual weather data is fetched or returned.

---

## Planned Use

The intent is for this endpoint to return weather information relevant to the user's location or a searched location. The file exists to reserve the route and signal that the feature is planned.

---

## What Needs to Be Done

To implement the weather feature:

1. Choose a weather API (e.g. Open-Meteo, OpenWeatherMap).
2. Implement `WeatherAPIHandler` in `weatherApi.go`.
3. Register it in `main.go` instead of `homeHandler`.
4. Add the frontend UI (likely on the index or a dedicated weather page).
5. Add tests in a corresponding `weather_test.go` file.
