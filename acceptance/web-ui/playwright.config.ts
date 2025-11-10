import { defineConfig, devices } from '@playwright/test';

/**
 * Playwright E2E Test Configuration for Cortex Web UI
 * Updated: November 2025
 *
 * See https://playwright.dev/docs/test-configuration
 */
export default defineConfig({
  testDir: './tests',

  /* Run tests in files in parallel */
  fullyParallel: true,

  /* Fail the build on CI if you accidentally left test.only */
  forbidOnly: !!process.env.CI,

  /* Retry on CI only */
  retries: process.env.CI ? 2 : 0,

  /* Opt out of parallel tests on CI */
  workers: process.env.CI ? 1 : undefined,

  /* Reporter to use */
  reporter: [
    ['html'],
    ['json', { outputFile: 'test-results/results.json' }],
    ['junit', { outputFile: 'test-results/junit.xml' }],
  ],

  /* Shared settings for all projects */
  use: {
    /* Base URL for tests */
    baseURL: process.env.BASE_URL || 'http://localhost:8080',

    /* Collect trace on first retry */
    trace: 'on-first-retry',

    /* Screenshot on failure */
    screenshot: 'only-on-failure',

    /* Video on failure */
    video: 'retain-on-failure',
  },

  /* Configure projects for major browsers */
  /* Note: Only Chromium installed to optimize test speed and disk space */
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },

    /* Disabled - Firefox not installed (uncomment and run: npx playwright install firefox)
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    */

    /* Disabled - WebKit not installed (uncomment and run: npx playwright install webkit)
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
    */

    /* Mobile viewports - Chromium engine */
    {
      name: 'Mobile Chrome',
      use: { ...devices['Pixel 5'] },
    },

    /* Disabled - Mobile Safari requires WebKit (uncomment and run: npx playwright install webkit)
    {
      name: 'Mobile Safari',
      use: { ...devices['iPhone 12'] },
    },
    */
  ],

  /* Run local dev server before tests (when available) */
  // webServer: {
  //   command: 'cortex ui --port 8080',
  //   url: 'http://localhost:8080',
  //   reuseExistingServer: !process.env.CI,
  //   timeout: 120 * 1000,
  // },
});
