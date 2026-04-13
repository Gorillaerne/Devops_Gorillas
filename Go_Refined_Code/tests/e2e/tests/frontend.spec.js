// @ts-check
const { test, expect } = require('@playwright/test');

test.describe('Search page', () => {
  test('search button triggers /api/search call', async ({ page }) => {
    await page.route('/api/search*', (route) =>
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ data: [] }),
      })
    );

    await page.goto('/html/index.html');

    const requestPromise = page.waitForRequest((req) =>
      req.url().includes('/api/search')
    );
    await page.fill('#search-input', 'golang');
    await page.click('#search-button');

    const request = await requestPromise;
    expect(request.url()).toContain('q=golang');
  });

  test('nav login link href is /login', async ({ page }) => {
    await page.goto('/html/index.html');
    expect(await page.getAttribute('#nav-login', 'href')).toBe('/login');
  });

  test('nav register link href is /register', async ({ page }) => {
    await page.goto('/html/index.html');
    expect(await page.getAttribute('#nav-register', 'href')).toBe('/register');
  });

  test('footer about link href is /about', async ({ page }) => {
    await page.goto('/html/index.html');
    const href = await page.locator('.footer a[href="/about"]').getAttribute('href');
    expect(href).toBe('/about');
  });
});

test.describe('Login page', () => {
  test('login form submit triggers /api/login call', async ({ page }) => {
    await page.route('/api/login', (route) =>
      route.fulfill({
        status: 401,
        contentType: 'application/json',
        body: JSON.stringify({ statusCode: 401, message: 'Invalid credentials' }),
      })
    );

    await page.goto('/html/login.html');

    const requestPromise = page.waitForRequest((req) =>
      req.url().includes('/api/login')
    );
    await page.fill('#username', 'testuser');
    await page.fill('#password', 'testpass');
    await page.click('#login-button');

    const request = await requestPromise;
    expect(request.method()).toBe('POST');
  });

  test('nav login link href is /login', async ({ page }) => {
    await page.goto('/html/login.html');
    expect(await page.getAttribute('#nav-login', 'href')).toBe('/login');
  });

  test('nav register link href is /register', async ({ page }) => {
    await page.goto('/html/login.html');
    expect(await page.getAttribute('#nav-register', 'href')).toBe('/register');
  });
});

test.describe('Register page', () => {
  test('register form submit triggers /api/register call', async ({ page }) => {
    await page.route('/api/register', (route) =>
      route.fulfill({
        status: 409,
        contentType: 'application/json',
        body: JSON.stringify({ statusCode: 409, message: 'User already exists' }),
      })
    );

    await page.goto('/html/register.html');

    const requestPromise = page.waitForRequest((req) =>
      req.url().includes('/api/register')
    );
    await page.fill('#reg-username', 'newuser');
    await page.fill('#reg-email', 'new@example.com');
    await page.fill('#reg-password', 'pass123');
    await page.fill('#reg-password-confirm', 'pass123');
    await page.click('#register-button');

    const request = await requestPromise;
    expect(request.method()).toBe('POST');
  });

  test('nav login link href is /login', async ({ page }) => {
    await page.goto('/html/register.html');
    expect(await page.getAttribute('#nav-login', 'href')).toBe('/login');
  });

  test('nav register link href is /register', async ({ page }) => {
    await page.goto('/html/register.html');
    expect(await page.getAttribute('#nav-register', 'href')).toBe('/register');
  });
});

test.describe('About page', () => {
  test('nav login link href is /login', async ({ page }) => {
    await page.goto('/html/about.html');
    expect(await page.getAttribute('#nav-login', 'href')).toBe('/login');
  });

  test('nav register link href is /register', async ({ page }) => {
    await page.goto('/html/about.html');
    expect(await page.getAttribute('#nav-register', 'href')).toBe('/register');
  });

  test('footer home link href is /', async ({ page }) => {
    await page.goto('/html/about.html');
    const href = await page.locator('.footer a').getAttribute('href');
    expect(href).toBe('/');
  });
});
