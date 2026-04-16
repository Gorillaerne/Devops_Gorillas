// @ts-check
const { test, expect } = require('@playwright/test');

test.describe('Profile page', () => {
    test('shows breach warning when breachWarning is set in sessionStorage', async ({ page }) => {
        await page.addInitScript(() => sessionStorage.setItem('breachWarning', '1'));
        await page.goto('/html/profile.html');

        const banner = page.locator('#breach-warning');
        await expect(banner).toBeVisible();
        await expect(banner).toContainText('Security alert');
    });

    test('does not show breach warning without query param', async ({ page }) => {
        await page.goto('/html/profile.html');

        const banner = page.locator('#breach-warning');
        await expect(banner).toBeHidden();
    });

    test('change-password form submit triggers /api/change-password call', async ({ page }) => {
        await page.route('/api/change-password', (route) =>
            route.fulfill({
                status: 200,
                contentType: 'application/json',
                body: JSON.stringify({ statusCode: 200, message: 'Password updated successfully' }),
            })
        );

        await page.goto('/html/profile.html');

        const requestPromise = page.waitForRequest((req) =>
            req.url().includes('/api/change-password')
        );

        await page.fill('#current-password', 'oldpass');
        await page.fill('#new-password', 'newpass123');
        await page.fill('#new-password-confirm', 'newpass123');
        await page.click('#change-password-button');

        const request = await requestPromise;
        expect(request.method()).toBe('POST');

        const body = JSON.parse(request.postData());
        expect(body.current_password).toBe('oldpass');
        expect(body.new_password).toBe('newpass123');
        expect(body.new_password2).toBe('newpass123');
    });

    test('change-password form sends Authorization header with token', async ({ page }) => {
        await page.addInitScript(() => {
            localStorage.setItem('token', 'test-jwt-token');
        });

        await page.route('/api/change-password', (route) =>
            route.fulfill({
                status: 200,
                contentType: 'application/json',
                body: JSON.stringify({ statusCode: 200, message: 'Password updated successfully' }),
            })
        );

        await page.goto('/html/profile.html');

        const requestPromise = page.waitForRequest((req) =>
            req.url().includes('/api/change-password')
        );

        await page.fill('#current-password', 'old');
        await page.fill('#new-password', 'new123');
        await page.fill('#new-password-confirm', 'new123');
        await page.click('#change-password-button');

        const request = await requestPromise;
        expect(request.headers()['authorization']).toBe('Bearer test-jwt-token');
    });

    test('shows error toast on failed password change', async ({ page }) => {
        await page.route('/api/change-password', (route) =>
            route.fulfill({
                status: 401,
                contentType: 'application/json',
                body: JSON.stringify({ statusCode: 401, message: 'Current password is incorrect' }),
            })
        );

        await page.goto('/html/profile.html');

        await page.fill('#current-password', 'wrong');
        await page.fill('#new-password', 'new');
        await page.fill('#new-password-confirm', 'new');
        await page.click('#change-password-button');

        const toast = page.locator('.toast--error');
        await expect(toast).toBeVisible();
        await expect(toast).toContainText('Current password is incorrect');
    });

    test('breach warning is hidden after successful password change', async ({ page }) => {
        await page.addInitScript(() => sessionStorage.setItem('breachWarning', '1'));

        await page.route('/api/change-password', (route) =>
            route.fulfill({
                status: 200,
                contentType: 'application/json',
                body: JSON.stringify({ statusCode: 200, message: 'Password updated successfully' }),
            })
        );

        await page.goto('/html/profile.html');
        await expect(page.locator('#breach-warning')).toBeVisible();

        await page.fill('#current-password', 'old');
        await page.fill('#new-password', 'new123');
        await page.fill('#new-password-confirm', 'new123');
        await page.click('#change-password-button');

        await expect(page.locator('#breach-warning')).toBeHidden();
    });

    test('nav login link href is /login', async ({ page }) => {
        await page.goto('/html/profile.html');
        expect(await page.getAttribute('#nav-login', 'href')).toBe('/login');
    });

    test('nav register link href is /register', async ({ page }) => {
        await page.goto('/html/profile.html');
        expect(await page.getAttribute('#nav-register', 'href')).toBe('/register');
    });
});

test.describe('Login page — breach redirect', () => {
    test('redirects to /profile?breached=1 when breached flag is true', async ({ page }) => {
        await page.route('/api/login', (route) =>
            route.fulfill({
                status: 200,
                contentType: 'application/json',
                body: JSON.stringify({
                    statusCode: 200,
                    message: 'Login successful',
                    token: 'fake-jwt',
                    breached: true,
                }),
            })
        );

        await page.goto('/html/login.html');
        await page.fill('#username', 'Benthe1954');
        await page.fill('#password', '^Jt^pLkzW2');
        await page.click('#login-button');

        await page.waitForURL('**/profile', { timeout: 5000 });
        expect(page.url()).toContain('/profile');
    });

    test('redirects to / when breached flag is false', async ({ page }) => {
        await page.route('/api/login', (route) =>
            route.fulfill({
                status: 200,
                contentType: 'application/json',
                body: JSON.stringify({
                    statusCode: 200,
                    message: 'Login successful',
                    token: 'fake-jwt',
                    breached: false,
                }),
            })
        );

        await page.goto('/html/login.html');
        await page.fill('#username', 'normaluser');
        await page.fill('#password', 'safepass');
        await page.click('#login-button');

        await page.waitForURL('**/', { timeout: 5000 });
        expect(page.url()).toMatch(/\/$/);
    });
});
