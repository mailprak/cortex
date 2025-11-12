import { test, expect } from '@playwright/test';

/**
 * E2E Tests for Cortex Dashboard
 * Testing the web UI from the user's perspective (outer loop TDD)
 */

test.describe('Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to dashboard before each test
    await page.goto('/');
  });

  test('should load in under 2 seconds', async ({ page }) => {
    const startTime = Date.now();
    await page.goto('/');
    await page.waitForLoadState('networkidle');
    const loadTime = Date.now() - startTime;

    expect(loadTime).toBeLessThan(2000);
  });

  test('should display neuron library', async ({ page }) => {
    await expect(page.getByText('Neuron Library')).toBeVisible();

    const neuronCards = page.locator('[data-testid="neuron-card"]');
    expect(await neuronCards.count()).toBeGreaterThan(0);
  });

  test('should display system metrics', async ({ page }) => {
    await expect(page.getByText('CPU')).toBeVisible();
    await expect(page.getByText('Memory')).toBeVisible();
    await expect(page.getByText('Disk')).toBeVisible();
  });

  test('should be responsive on mobile devices', async ({ page }) => {
    // Set viewport to iPhone dimensions
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto('/');

    // Navigation should be hamburger menu on mobile
    await expect(page.getByRole('button', { name: /menu/i })).toBeVisible();
  });
});

test.describe('Neuron Execution', () => {
  test.skip('should execute neuron and show real-time logs', async ({ page }) => {
    // TODO: Requires backend WebSocket implementation for real-time log streaming
    await page.goto('/');

    // Click neuron to execute
    await page.click('[data-testid="neuron-card"]:first-child');
    await page.click('button:has-text("Execute")');

    // Wait for WebSocket logs
    await expect(page.locator('[data-testid="log-stream"]')).toBeVisible();

    // Verify logs updating
    const logContent = page.locator('[data-testid="log-stream"]');
    await expect(logContent).toContainText(/executing|running|complete/i);
  });

  test.skip('should display execution status updates', async ({ page }) => {
    // TODO: Requires backend execution API and status tracking
    await page.goto('/');
    await page.click('[data-testid="neuron-card"]:first-child');
    await page.click('button:has-text("Execute")');

    // Verify status changes
    await expect(page.locator('[data-testid="execution-status"]'))
      .toHaveText('Running');
    await expect(page.locator('[data-testid="execution-status"]'))
      .toHaveText('Completed', { timeout: 10000 });
  });
});

test.describe('Accessibility', () => {
  test('should have proper ARIA labels', async ({ page }) => {
    await page.goto('/');

    await expect(page.locator('[aria-label="Main navigation"]')).toBeVisible();
    await expect(page.locator('[aria-label="Neuron library"]')).toBeVisible();
  });

  test('should support keyboard navigation', async ({ page }) => {
    await page.goto('/');

    // Tab through elements
    await page.keyboard.press('Tab');
    const firstFocused = await page.evaluate(() => document.activeElement?.tagName);
    expect(firstFocused).toBeTruthy();
  });
});

test.describe('WebSocket Performance', () => {
  test.skip('should maintain latency under 100ms', async ({ page }) => {
    // TODO: Requires backend WebSocket server implementation
    await page.goto('/');

    // Measure WebSocket round-trip time
    const startTime = Date.now();

    await page.evaluate(() => {
      // @ts-ignore - WebSocket will be available when implemented
      window.ws.send(JSON.stringify({ type: 'ping' }));
    });

    await page.waitForFunction(() => {
      // @ts-ignore
      return window.lastPongReceived;
    });

    const latency = Date.now() - startTime;
    expect(latency).toBeLessThan(100);
  });
});
