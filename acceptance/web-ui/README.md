# Cortex Web UI E2E Tests

End-to-end tests for the Cortex web dashboard using Playwright.

## Setup

```bash
# From project root
make install-deps

# Or manually
cd acceptance/web-ui
npm install
npx playwright install --with-deps
```

## Running Tests

```bash
# From project root
make test-acceptance-web

# Or from this directory
npm test                  # Run all tests
npm run test:ui          # Run in UI mode (interactive)
npm run test:debug       # Run in debug mode
npm run test:headed      # Run with browser visible
npm run report           # Show test report
```

## Test Structure

All tests follow **outer-loop TDD** - they test complete user workflows:

```
tests/
├── dashboard.spec.ts       # Dashboard UI tests
├── neuron-execution.spec.ts # Neuron execution tests (future)
├── synapse-builder.spec.ts  # Visual builder tests (future)
└── fleet-management.spec.ts # Fleet UI tests (future)
```

## Writing Tests

Tests should be skipped until features are implemented (RED phase):

```typescript
test('should do something', async ({ page }) => {
  test.skip(true, 'Feature not yet implemented - TDD RED phase');

  // Test code here...
});
```

Remove the skip when ready to implement (GREEN phase).

## Browser Coverage

Tests run on:
- Chrome/Chromium
- Firefox
- Safari/WebKit
- Mobile Chrome (Pixel 5)
- Mobile Safari (iPhone 12)

## Performance Requirements

From specification:
- Dashboard load: < 2s on 3G
- WebSocket latency: < 100ms
- Real-time updates: < 100ms

## CI/CD Integration

Tests run in GitHub Actions on every PR.

See `../../.github/workflows/test.yml` for configuration.
