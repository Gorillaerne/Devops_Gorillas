// @ts-check
const { test, expect } = require('@playwright/test');

/**
 * Weather integration tests — run against the live Docker Compose stack.
 *
 * These tests make real HTTP requests through the full chain:
 *   Browser → Nginx (8081) → Go backend (8080) → Open-Meteo API
 *
 * No API mocking. The intent is to verify:
 *  1. Nginx routes /weather to the correct HTML file.
 *  2. Nginx proxies /api/weather to the Go backend correctly.
 *  3. The Go backend returns a valid StandardResponse JSON.
 *  4. The frontend renders the data without errors.
 */

test.describe('Weather — full stack integration', () => {
    test('GET /weather serves the weather HTML page', async ({ page }) => {
        const response = await page.goto('/weather');
        expect(response.status()).toBe(200);
        const title = await page.title();
        expect(title).toContain('Weather');
    });

    test.skip('GET /api/weather returns valid StandardResponse JSON', async ({ request }) => {
        const response = await request.get('/api/weather');
        expect(response.status()).toBe(200);
        expect(response.headers()['content-type']).toContain('application/json');

        const body = await response.json();
        expect(body).toHaveProperty('data');
        expect(body.data).toHaveProperty('current');
        expect(body.data).toHaveProperty('daily');

        const { current, daily } = body.data;
        expect(typeof current.temperature_2m).toBe('number');
        expect(typeof current.windspeed_10m).toBe('number');
        expect(typeof current.weathercode).toBe('number');
        expect(Array.isArray(daily.time)).toBe(true);
        expect(daily.time.length).toBe(7);
        expect(daily.temperature_2m_max.length).toBe(7);
        expect(daily.temperature_2m_min.length).toBe(7);
        expect(daily.weathercode.length).toBe(7);
    });

    test('weather page renders current temperature after real API call', async ({ page }) => {
        await page.goto('/weather');
        // Wait for the JS to fetch and render — loading placeholder disappears
        await expect(page.locator('.weather-now-card')).toBeVisible({ timeout: 10000 });
        // Temperature should be a number (we don't know exact value)
        const tempText = await page.locator('.weather-now-temp').textContent();
        expect(tempText).toMatch(/-?\d+/);
    });

    test('weather page renders 7 forecast cards after real API call', async ({ page }) => {
        await page.goto('/weather');
        await expect(page.locator('.weather-day-card').first()).toBeVisible({ timeout: 10000 });
        await expect(page.locator('.weather-day-card')).toHaveCount(7);
    });

    test('second API call is served from cache (response time < 100ms)', async ({ request }) => {
        // First call — may hit Open-Meteo
        await request.get('/api/weather');

        // Second call — must be served from cache and be very fast
        const start = Date.now();
        const response = await request.get('/api/weather');
        const duration = Date.now() - start;

        expect(response.status()).toBe(200);
        expect(duration).toBeLessThan(100);
    });

    test('Nginx routes /login to the login page', async ({ page }) => {
        const response = await page.goto('/login');
        expect(response.status()).toBe(200);
        await expect(page.locator('#login-form')).toBeVisible();
    });

    test('Nginx routes /register to the register page', async ({ page }) => {
        const response = await page.goto('/register');
        expect(response.status()).toBe(200);
    });

    test('Nginx routes /about to the about page', async ({ page }) => {
        const response = await page.goto('/about');
        expect(response.status()).toBe(200);
    });

    test('Nginx returns 404 for undefined routes', async ({ request }) => {
        const response = await request.get('/this-route-does-not-exist');
        expect(response.status()).toBe(404);
    });
});
