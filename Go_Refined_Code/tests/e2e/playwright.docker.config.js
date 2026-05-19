// @ts-check
const { defineConfig } = require('@playwright/test');

/**
 * Docker Compose integration test config.
 *
 * Run against the full running stack (Nginx + Go + MySQL):
 *   docker compose -f compose.dev.yml up -d
 *   npx playwright test --config=playwright.docker.config.js
 *
 * Unlike the default config, there is no webServer here — the stack must
 * already be running. Tests navigate through Nginx (port 8081) so routing,
 * API proxying, and the Go backend are all exercised for real.
 */
module.exports = defineConfig({
    testDir: './docker-tests',
    use: {
        baseURL: 'http://localhost:8081',
    },
});
