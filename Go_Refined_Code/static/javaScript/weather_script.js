import { checkIfLoggedIn } from "./reuseable_functions.js"
import { showError } from "./reuseable_functions.js"

checkIfLoggedIn()

const WMO_DESCRIPTION = {
    0: "Clear sky",
    1: "Mainly clear", 2: "Partly cloudy", 3: "Overcast",
    45: "Foggy", 48: "Icy fog",
    51: "Light drizzle", 53: "Moderate drizzle", 55: "Dense drizzle",
    61: "Slight rain", 63: "Moderate rain", 65: "Heavy rain",
    71: "Slight snow", 73: "Moderate snow", 75: "Heavy snow",
    77: "Snow grains",
    80: "Slight showers", 81: "Moderate showers", 82: "Violent showers",
    85: "Slight snow showers", 86: "Heavy snow showers",
    95: "Thunderstorm",
    96: "Thunderstorm with hail", 99: "Thunderstorm with heavy hail",
}

const WMO_ICON = {
    0: "☀️",
    1: "🌤️", 2: "⛅", 3: "☁️",
    45: "🌫️", 48: "🌫️",
    51: "🌦️", 53: "🌦️", 55: "🌧️",
    61: "🌧️", 63: "🌧️", 65: "🌧️",
    71: "🌨️", 73: "🌨️", 75: "❄️",
    77: "❄️",
    80: "🌦️", 81: "🌧️", 82: "⛈️",
    85: "🌨️", 86: "❄️",
    95: "⛈️",
    96: "⛈️", 99: "⛈️",
}

function wmoDescription(code) {
    return WMO_DESCRIPTION[code] ?? "Unknown"
}

function wmoIcon(code) {
    return WMO_ICON[code] ?? "🌡️"
}

function formatDate(dateStr) {
    const d = new Date(dateStr + "T12:00:00")
    return d.toLocaleDateString("en-GB", { weekday: "short", month: "short", day: "numeric" })
}

function renderCurrent(current) {
    const container = document.getElementById("weather-current")
    const code = current.weathercode

    container.innerHTML = `
        <div class="weather-now-card">
            <div class="weather-now-icon">${wmoIcon(code)}</div>
            <div class="weather-now-info">
                <div class="weather-now-temp">${Math.round(current.temperature_2m)}<span class="weather-unit">°C</span></div>
                <div class="weather-now-desc">${wmoDescription(code)}</div>
                <div class="weather-now-wind">Wind: ${Math.round(current.windspeed_10m)} km/h</div>
            </div>
        </div>
    `
}

function renderForecast(daily) {
    const container = document.getElementById("weather-forecast")
    const { time, temperature_2m_max, temperature_2m_min, weathercode } = daily

    container.innerHTML = `<h2 class="weather-forecast-title">7-Day Forecast</h2>`

    const grid = document.createElement("div")
    grid.className = "weather-grid"

    time.forEach((date, i) => {
        const card = document.createElement("div")
        card.className = "weather-day-card"
        card.innerHTML = `
            <div class="weather-day-date">${formatDate(date)}</div>
            <div class="weather-day-icon">${wmoIcon(weathercode[i])}</div>
            <div class="weather-day-desc">${wmoDescription(weathercode[i])}</div>
            <div class="weather-day-temps">
                <span class="weather-temp-high">${Math.round(temperature_2m_max[i])}°</span>
                <span class="weather-temp-low">${Math.round(temperature_2m_min[i])}°</span>
            </div>
        `
        grid.appendChild(card)
    })

    container.appendChild(grid)
}

fetch("/api/weather")
    .then(res => {
        if (!res.ok) throw new Error(`Weather API returned ${res.status}`)
        return res.json()
    })
    .then(json => {
        const weather = json.data
        renderCurrent(weather.current)
        renderForecast(weather.daily)
    })
    .catch(err => {
        console.error("Weather fetch failed:", err)
        showError("Could not load weather data. Please try again later.")
        document.getElementById("weather-current").innerHTML =
            `<p class="weather-loading">Weather unavailable.</p>`
    })
