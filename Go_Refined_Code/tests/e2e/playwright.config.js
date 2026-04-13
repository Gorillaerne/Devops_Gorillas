// @ts-check
const { defineConfig } = require('@playwright/test');

module.exports = defineConfig({
  testDir: './tests',
  use: {
    baseURL: 'http://localhost:4321',
  },
  webServer: {
    command: 'npx serve ../../static -p 4321 --no-clipboard',
    port: 4321,
    reuseExistingServer: !process.env.CI,
  },
});
