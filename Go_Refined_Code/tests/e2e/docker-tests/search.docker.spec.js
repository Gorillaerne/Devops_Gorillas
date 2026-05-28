// @ts-check
const { test, expect } = require('@playwright/test');

/**
 * Search integration tests — run against the live Docker Compose stack.
 *
 * These tests make real HTTP requests through the full chain:
 *   Browser → Nginx (8081) → Go backend (8080) → MySQL
 *
 * No API mocking. The intent is to verify:
 *  1. The search API returns real results from the database.
 *  2. The frontend renders results correctly end-to-end.
 *  3. The register → login → search flow works as a whole.
 */

test.describe('Search API — full stack integration', () => {
    test('GET /api/search returns JSON with data array', async ({ request }) => {
        const response = await request.get('/api/search?q=python&language=en');
        expect(response.status()).toBe(200);
        expect(response.headers()['content-type']).toContain('application/json');

        const body = await response.json();
        expect(body).toHaveProperty('data');
        expect(Array.isArray(body.data)).toBe(true);
    });

    test('GET /api/search results contain title, description and URL', async ({ request }) => {
        const response = await request.get('/api/search?q=python&language=en');
        const body = await response.json();

        expect(body.data.length).toBeGreaterThan(0);

        const first = body.data[0];
        expect(typeof first.title).toBe('string');
        expect(first.title.length).toBeGreaterThan(0);
        expect(typeof first.description).toBe('string');
        expect(first.description.length).toBeGreaterThan(0);
        expect(typeof first.URL).toBe('string');
        expect(first.URL).toMatch(/^https?:\/\//);
    });

    test('GET /api/search top result for "python" is Python', async ({ request }) => {
        const response = await request.get('/api/search?q=python&language=en');
        const body = await response.json();

        expect(body.data.length).toBeGreaterThan(0);
        expect(body.data[0].title.toLowerCase()).toContain('python');
    });

    test('GET /api/search returns 422 when q is missing', async ({ request }) => {
        const response = await request.get('/api/search');
        expect(response.status()).toBe(422);

        const body = await response.json();
        expect(body).toHaveProperty('message');
    });

    test('GET /api/search returns empty data for gibberish query', async ({ request }) => {
        const response = await request.get('/api/search?q=xqzwjfkplm&language=en');
        expect(response.status()).toBe(200);
        const body = await response.json();
        expect(Array.isArray(body.data)).toBe(true);
        expect(body.data.length).toBe(0);
    });
});

test.describe('Search UI — full stack integration', () => {
    test('searching for "python" renders results on the page', async ({ page }) => {
        await page.goto('/');
        await page.fill('#search-input', 'python');
        await page.click('#search-button');

        await expect(page.locator('.result-item').first()).toBeVisible({ timeout: 10000 });
        const count = await page.locator('.result-item').count();
        expect(count).toBeGreaterThan(0);
    });

    test('each result shows a title and description', async ({ page }) => {
        await page.goto('/');
        await page.fill('#search-input', 'python');
        await page.click('#search-button');

        await expect(page.locator('.result-item').first()).toBeVisible({ timeout: 10000 });

        const title = await page.locator('.result-title').first().textContent();
        const snippet = await page.locator('.result-snippet').first().textContent();

        expect(title.length).toBeGreaterThan(0);
        expect(snippet).not.toBe('Ingen beskrivelse tilgængelig');
    });

    test('pressing Enter triggers search', async ({ page }) => {
        await page.goto('/');
        await page.fill('#search-input', 'python');
        await page.press('#search-input', 'Enter');

        await expect(page.locator('.result-item').first()).toBeVisible({ timeout: 10000 });
    });

    test('gibberish query shows no-results message', async ({ page }) => {
        await page.goto('/');
        await page.fill('#search-input', 'xqzwjfkplm');
        await page.click('#search-button');

        await expect(page.locator('.no-results')).toBeVisible({ timeout: 10000 });
    });
});

test.describe('Register → Login → Search — full stack flow', () => {
    const timestamp = Date.now();
    const username = `e2euser_${timestamp}`;
    const email = `e2e_${timestamp}@test.com`;
    const password = 'TestPassword123!';

    test('can register a new user', async ({ page }) => {
        await page.goto('/register');
        await page.fill('#reg-username', username);
        await page.fill('#reg-email', email);
        await page.fill('#reg-password', password);
        await page.fill('#reg-password-confirm', password);
        await page.click('#register-button');

        // Should redirect to login or show success
        await expect(page).toHaveURL(/login/, { timeout: 10000 });
    });

    test('can login and then search', async ({ page }) => {
        await page.goto('/login');
        await page.fill('#username', username);
        await page.fill('#password', password);
        await page.click('#login-button');

        // Wait for redirect to home after login
        await expect(page).toHaveURL(/\/$/, { timeout: 10000 });

        // Search for something
        await page.fill('#search-input', 'python');
        await page.click('#search-button');

        await expect(page.locator('.result-item').first()).toBeVisible({ timeout: 10000 });
    });
});
