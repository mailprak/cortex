# Cortex Acceptance Tests

This directory contains all acceptance tests (outer-loop TDD) for Cortex, organized by user interface.

## Directory Structure

```
acceptance/
├── cli/          # Command-line interface tests (Ginkgo/Gomega)
├── web-ui/       # Web UI end-to-end tests (Playwright)
├── api/          # REST API acceptance tests (future)
└── README.md     # This file
```

## Test Organization Philosophy

**Acceptance tests** validate the system from the **user's perspective**, testing complete features end-to-end. They follow **outer-loop TDD**:

1. Write failing acceptance test (RED)
2. Write unit tests and implementation (GREEN - inner loop)
3. Refactor (REFACTOR)
4. Acceptance test passes (GREEN - outer loop)

## CLI Acceptance Tests (`cli/`)

Tests the Cortex command-line interface using **Ginkgo v2** and **Gomega**.

**Technology**: Go + Ginkgo/Gomega
**Target**: Cortex binary execution

**Example tests**:
- `cortex --help` displays help text
- `cortex create-synapse <name>` creates synapse structure
- `cortex generate-neuron --prompt "..."` generates AI neurons
- `cortex execute-synapse <path>` runs workflows

**Run CLI tests**:
```bash
make test-acceptance-cli
# or
ginkgo -r -v --label-filter="cli" ./acceptance/cli/
```

## Web UI Acceptance Tests (`web-ui/`)

Tests the Cortex web dashboard using **Playwright** for browser automation.

**Technology**: TypeScript + Playwright
**Target**: Web UI in real browsers (Chrome, Firefox, Safari)

**Example tests**:
- Dashboard loads in < 2 seconds
- Real-time log streaming via WebSocket
- Visual synapse builder drag-and-drop
- Mobile responsiveness
- Accessibility (WCAG 2.1 AAA)

**Run Web UI tests**:
```bash
make test-acceptance-web
# or
cd acceptance/web-ui && npx playwright test
```

## API Acceptance Tests (`api/`)

Tests the REST API endpoints (when implemented).

**Technology**: Go + Ginkgo/Gomega (or similar)
**Target**: HTTP API endpoints

**Future tests**:
- POST /api/neurons - Create neuron
- GET /api/neurons - List neurons
- POST /api/synapses/execute - Execute workflow
- WebSocket /ws/logs - Stream logs

## Running All Acceptance Tests

```bash
# Run all acceptance tests (CLI + Web UI + API)
make test-acceptance

# Run with coverage
make coverage

# Watch mode (re-run on changes)
make watch
```

## Test Labels

All acceptance tests use Ginkgo labels for filtering:

- `Label("acceptance")` - All acceptance tests
- `Label("cli")` - CLI-specific tests
- `Label("web-ui")` - Web UI tests
- `Label("api")` - API tests
- `Label("neuron")` - Neuron-related tests
- `Label("synapse")` - Synapse-related tests
- `Label("ai")` - AI generation tests

**Filter examples**:
```bash
# Run only AI-related acceptance tests
ginkgo -r --label-filter="acceptance && ai"

# Run only CLI neuron tests
ginkgo -r --label-filter="cli && neuron"
```

## Writing New Acceptance Tests

### CLI Tests (Ginkgo)

```go
// acceptance/cli/my_feature_test.go
package cli_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "github.com/onsi/gomega/gexec"
)

var _ = Describe("My Feature", Label("acceptance", "cli"), func() {
    It("should do something from user's perspective", func() {
        session := RunCortex("my-command", "--flag", "value")
        Eventually(session).Should(gexec.Exit(0))
        Eventually(session.Out).Should(gbytes.Say("expected output"))
    })
})
```

### Web UI Tests (Playwright)

```typescript
// acceptance/web-ui/my-feature.spec.ts
import { test, expect } from '@playwright/test';

test.describe('My Feature', () => {
  test('should do something from user's perspective', async ({ page }) => {
    await page.goto('/my-feature');
    await page.click('button:has-text("Action")');
    await expect(page.locator('.result')).toBeVisible();
  });
});
```

## Test Data & Fixtures

Shared test fixtures are located in:
- `acceptance/fixtures/` - Shared test data (YAML configs, scripts, etc.)
- `acceptance/helpers/` - Shared test utilities

## Coverage Goals

- **Acceptance tests**: Cover all user-facing features
- **Coverage threshold**: 90% (enforced in CI)
- **Focus**: User workflows, not implementation details

## CI/CD Integration

Acceptance tests run in GitHub Actions on:
- Every pull request
- Before merge to main
- Nightly full test suite

See `.github/workflows/test.yml` for CI configuration.

## Performance Requirements

Based on spec requirements:

- **CLI tests**: Each test < 10s
- **Web UI load**: < 2s on 3G
- **WebSocket latency**: < 100ms
- **Total suite**: < 5 minutes

## Troubleshooting

### Ginkgo not found
```bash
make install-deps
```

### Playwright browsers missing
```bash
cd acceptance/web-ui && npx playwright install --with-deps
```

### Tests timing out
Increase timeout in test:
```go
Eventually(session, "30s").Should(gexec.Exit(0))
```

---

**Last Updated**: November 2025
**Maintained By**: Cortex Core Team
