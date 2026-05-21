// @ts-check
const { test, expect } = require('@playwright/test');

const MOCK_WEATHER = {
    data: {
        current: {
            temperature_2m: 18.5,
            windspeed_10m: 12.3,
            weathercode: 2,
            is_day: 1,
        },
        daily: {
            time: ['2026-05-19', '2026-05-20', '2026-05-21', '2026-05-22', '2026-05-23', '2026-05-24', '2026-05-25'],
            temperature_2m_max: [20, 18, 15, 17, 19, 21, 16],
            temperature_2m_min: [12, 10, 8, 11, 13, 14, 9],
            weathercode: [2, 61, 3, 1, 0, 2, 80],
        },
    },
};

test.describe('Weather page', () => {
    test.beforeEach(async ({ page }) => {
        await page.route('/api/weather', (route) =>
            route.fulfill({
                status: 200,
                contentType: 'application/json',
                body: JSON.stringify(MOCK_WEATHER),
            })
        );
        await page.goto('/html/weather.html');
    });

    test('shows current temperature', async ({ page }) => {
        await expect(page.locator('.weather-now-temp')).toBeVisible();
        // 18.5 rounds to 19 (Math.round)
        await expect(page.locator('.weather-now-temp')).toContainText('19');
    });

    test('shows current weather condition', async ({ page }) => {
        await expect(page.locator('.weather-now-desc')).toContainText('Partly cloudy');
    });

    test('shows wind speed', async ({ page }) => {
        await expect(page.locator('.weather-now-wind')).toContainText('12 km/h');
    });

    test('renders 7 forecast day cards', async ({ page }) => {
        await expect(page.locator('.weather-day-card')).toHaveCount(7);
    });

    test('first forecast card shows correct high and low temperatures', async ({ page }) => {
        const firstCard = page.locator('.weather-day-card').first();
        await expect(firstCard.locator('.weather-temp-high')).toContainText('20°');
        await expect(firstCard.locator('.weather-temp-low')).toContainText('12°');
    });

    test('second forecast card shows rain condition', async ({ page }) => {
        // weathercode 61 = "Slight rain"
        const secondCard = page.locator('.weather-day-card').nth(1);
        await expect(secondCard.locator('.weather-day-desc')).toContainText('Slight rain');
    });

    test('weather nav link is present and points to /weather', async ({ page }) => {
        await expect(page.locator('nav a[href="/weather"]')).toBeVisible();
    });
});

test.describe('Weather page — API error handling', () => {
    test('shows fallback message when API returns 500', async ({ page }) => {
        await page.route('/api/weather', (route) =>
            route.fulfill({ status: 500, body: 'Internal Server Error' })
        );
        await page.goto('/html/weather.html');

        await expect(page.locator('.weather-loading')).toContainText('Weather unavailable');
        await expect(page.locator('.toast--error')).toBeVisible();
    });

    test('shows error toast with message when API returns 502', async ({ page }) => {
        await page.route('/api/weather', (route) =>
            route.fulfill({ status: 502, body: 'Bad Gateway' })
        );
        await page.goto('/html/weather.html');

        const toast = page.locator('.toast--error');
        await expect(toast).toBeVisible();
        await expect(toast).toContainText('Could not load weather data');
    });
});

test.describe('Weather nav link on all pages', () => {
    test('index page has weather nav link', async ({ page }) => {
        await page.goto('/html/index.html');
        await expect(page.locator('nav a[href="/weather"]')).toBeVisible();
    });

    test('login page has weather nav link', async ({ page }) => {
        await page.goto('/html/login.html');
        await expect(page.locator('nav a[href="/weather"]')).toBeVisible();
    });

    test('register page has weather nav link', async ({ page }) => {
        await page.goto('/html/register.html');
        await expect(page.locator('nav a[href="/weather"]')).toBeVisible();
    });

    test('about page has weather nav link', async ({ page }) => {
        await page.goto('/html/about.html');
        await expect(page.locator('nav a[href="/weather"]')).toBeVisible();
    });

    test('profile page has weather nav link', async ({ page }) => {
        await page.addInitScript(() => localStorage.setItem('token', 'fake-jwt'));
        await page.goto('/html/profile.html');
        await expect(page.locator('nav a[href="/weather"]')).toBeVisible();
    });
});
