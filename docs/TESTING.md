# Cortex Testing Guide

Comprehensive guide for Test-Driven Development with Cortex using Ginkgo v2, Gomega, and Playwright.

**Last Updated:** November 2025

---

## Table of Contents

1. [Testing Philosophy](#testing-philosophy)
2. [Test Structure](#test-structure)
3. [Quick Start](#quick-start)
4. [Running Tests](#running-tests)
5. [Writing Tests](#writing-tests)
6. [TDD Workflow](#tdd-workflow)
7. [Coverage](#coverage)
8. [CI/CD Integration](#cicd-integration)

---

## Testing Philosophy

Cortex follows **Test-Driven Development (TDD)** with both outer-loop and inner-loop cycles:

### Outer Loop TDD (Acceptance Tests)
1. Write failing acceptance test for user-facing feature (RED)
2. Write unit tests and implementation (inner loop)
3. Acceptance test passes (GREEN)
4. Refactor

### Inner Loop TDD (Unit Tests)
1. Write failing unit test (RED)
2. Write minimum code to pass (GREEN)
3. Refactor
4. Repeat

---

## Test Structure

```
cortex/
├── acceptance/              # Outer loop TDD
│   ├── cli/                # CLI acceptance tests (Ginkgo)
│   ├── web-ui/             # Web UI E2E tests (Playwright)
│   ├── api/                # API acceptance tests (future)
│   └── README.md
├── internal/               # Inner loop TDD
│   ├── neuron/
│   │   ├── neuron.go
│   │   ├── neuron_test.go  # Unit tests
│   │   └── neuron_suite_test.go
│   └── synapse/
├── coverage/               # Test coverage reports
└── Makefile               # Test automation
```

---

## Quick Start

### 1. Install Dependencies

```bash
make install-deps
```

This installs:
- **Ginkgo v2** - BDD test framework for Go
- **Gomega** - Matcher library
- **Playwright** - Browser automation for E2E tests

### 2. Run All Tests

```bash
make test-all
```

### 3. Watch Mode (TDD)

```bash
make watch
```

Tests auto-run on file changes - perfect for TDD!

---

## Running Tests

### All Tests

```bash
make test-all               # Run everything
make test                   # Alias for test-all
```

### Acceptance Tests (Outer Loop)

```bash
make test-acceptance        # All acceptance tests
make test-acceptance-cli    # CLI tests only
make test-acceptance-web    # Web UI tests only
```

### Unit Tests (Inner Loop)

```bash
make test-unit              # All unit tests
```

### Coverage

```bash
make coverage               # Generate HTML coverage report
open coverage/coverage.html # View in browser
```

### Specific Labels

```bash
# Run only AI-related tests
ginkgo -r --label-filter="ai"

# Run only neuron tests
ginkgo -r --label-filter="neuron"

# Run acceptance + AI tests
ginkgo -r --label-filter="acceptance && ai"
```

---

## Writing Tests

### CLI Acceptance Test (Ginkgo)

```go
// acceptance/cli/my_feature_test.go
package acceptance_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "github.com/onsi/gomega/gexec"
)

var _ = Describe("My Feature", Label("acceptance", "cli"), func() {
    It("should work from user's perspective", func() {
        // Start with failing test (RED)
        Skip("Not yet implemented - TDD RED phase")

        session := RunCortex("my-command", "--flag", "value")
        Eventually(session).Should(gexec.Exit(0))
        Eventually(session.Out).Should(gbytes.Say("expected output"))
    })
})
```

### Web UI E2E Test (Playwright)

```typescript
// acceptance/web-ui/tests/my-feature.spec.ts
import { test, expect } from '@playwright/test';

test.describe('My Feature', () => {
  test('should work from user's perspective', async ({ page }) => {
    // Start with failing test (RED)
    test.skip(true, 'Not yet implemented - TDD RED phase');

    await page.goto('/my-feature');
    await page.click('button:has-text("Action")');
    await expect(page.locator('.result')).toBeVisible();
  });
});
```

### Unit Test (Ginkgo)

```go
// internal/mypackage/mycode_test.go
package mypackage_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var _ = Describe("MyCode", Label("unit"), func() {
    It("should do something specific", func() {
        // Start with failing test (RED)
        Skip("Not yet implemented - TDD RED phase")

        result := MyFunction()
        Expect(result).To(Equal("expected"))
    })
})
```

---

## TDD Workflow

### Example: Adding AI Neuron Generation

#### 1. Write Acceptance Test (Outer Loop - RED)

```go
// acceptance/cli/ai_neuron_generation_test.go
It("should generate neuron from prompt", func() {
    Skip("Not yet implemented - TDD RED phase")

    session := RunCortex("generate-neuron",
        "--prompt", "Find which process is using port 8080 and show the full command with PID")

    Eventually(session).Should(gexec.Exit(0))
    Eventually(session.Out).Should(gbytes.Say("Neuron generated"))
    Eventually(session.Out).Should(gbytes.Say("lsof"))
})
```

```bash
make test-acceptance-cli  # FAILS (RED ✗)
```

#### 2. Write Unit Tests (Inner Loop - RED)

```go
// internal/ai/generator_test.go
It("should call OpenAI API", func() {
    Skip("Not yet implemented - TDD RED phase")

    generator := NewGenerator("openai")
    result := generator.Generate("Find which process is using port 8080 and show the full command with PID")

    Expect(result.Name).To(ContainSubstring("port"))
    Expect(result.Script).To(ContainSubstring("lsof"))
})
```

```bash
make test-unit  # FAILS (RED ✗)
```

#### 3. Implement Minimum Code (Inner Loop - GREEN)

```go
// internal/ai/generator.go
func (g *Generator) Generate(prompt string) (*Neuron, error) {
    // Minimum implementation
    return &Neuron{Name: "check-nginx"}, nil
}
```

```bash
make test-unit  # PASSES (GREEN ✓)
```

#### 4. Refactor Unit Code

```go
// Improve implementation, add error handling, etc.
```

```bash
make test-unit  # Still PASSES (GREEN ✓)
```

#### 5. Wire Up Command (Outer Loop - GREEN)

```go
// cmd/cortex/generate.go
// Implement CLI command using Generator
```

```bash
make test-acceptance-cli  # PASSES (GREEN ✓)
```

#### 6. Refactor Everything

```bash
make test-all  # Everything PASSES (GREEN ✓)
```

---

## Coverage

### Requirements

- **Minimum threshold:** 90%
- **Enforced in CI:** Yes
- **Focus areas:**
  - User workflows (acceptance)
  - Business logic (unit)
  - Error handling
  - Edge cases

### Checking Coverage

```bash
make coverage
```

Output shows:
- Coverage percentage per package
- Coverage report (HTML)
- Fails if below 90% threshold

### Improving Coverage

1. Identify uncovered code:
   ```bash
   go tool cover -func=coverage/coverage.out | grep -v "100.0%"
   ```

2. Write tests for uncovered areas

3. Re-run coverage:
   ```bash
   make coverage
   ```

---

## CI/CD Integration

### GitHub Actions

Tests run automatically on:
- Every pull request
- Before merge to main
- Nightly full test suite

### Local Pre-Commit

Add to `.git/hooks/pre-commit`:

```bash
#!/bin/bash
make test-all
if [ $? -ne 0 ]; then
    echo "Tests failed! Commit aborted."
    exit 1
fi
```

Make executable:
```bash
chmod +x .git/hooks/pre-commit
```

---

## Test Labels

All tests use Ginkgo labels for filtering:

| Label | Description |
|-------|-------------|
| `acceptance` | Outer loop acceptance tests |
| `unit` | Inner loop unit tests |
| `cli` | CLI-specific tests |
| `web-ui` | Web UI tests |
| `api` | API tests |
| `neuron` | Neuron-related tests |
| `synapse` | Synapse-related tests |
| `ai` | AI generation tests |

### Filter Examples

```bash
# Only acceptance tests
ginkgo -r --label-filter="acceptance"

# Only unit tests
ginkgo -r --label-filter="unit"

# AI neuron acceptance tests
ginkgo -r --label-filter="acceptance && ai && neuron"

# Everything except web-ui
ginkgo -r --label-filter="!web-ui"
```

---

## Troubleshooting

### Ginkgo not found

```bash
make install-deps
```

### Tests timeout

Increase timeout in test:
```go
Eventually(session, "30s").Should(gexec.Exit(0))
```

### Playwright browsers missing

```bash
cd acceptance/web-ui
npx playwright install --with-deps
```

### Coverage threshold failures

Current coverage is below 90%. Write more tests!

```bash
# See what's not covered
go tool cover -func=coverage/coverage.out | grep -v "100.0%"
```

---

## Performance Requirements

Based on specifications:

| Test Type | Requirement |
|-----------|-------------|
| CLI tests | < 10s per test |
| Web UI load | < 2s on 3G |
| WebSocket latency | < 100ms |
| Full test suite | < 5 minutes |

---

## Best Practices

### ✅ DO

- Start with failing test (RED)
- Write minimum code to pass (GREEN)
- Refactor with passing tests
- Use descriptive test names
- Test user workflows, not implementation
- Keep tests fast and focused
- Use `Skip()` for not-yet-implemented features

### ❌ DON'T

- Skip writing tests
- Test implementation details
- Write tests after code
- Have slow-running tests
- Ignore test failures
- Commit with failing tests

---

## Resources

- **Ginkgo v2 Docs:** https://onsi.github.io/ginkgo/
- **Gomega Matchers:** https://onsi.github.io/gomega/
- **Playwright Docs:** https://playwright.dev/
- **TDD Guide:** https://martinfowler.com/bliki/TestDrivenDevelopment.html

---

## Getting Help

- **Issues:** https://github.com/anoop2811/cortex/issues
- **Discussions:** Use `test` label

---

**Maintained By:** Cortex Core Team
**Last Updated:** November 2025
