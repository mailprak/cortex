# Contributing to Cortex

Thank you for your interest in contributing to Cortex! This guide will help you get started.

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
2. [How Can I Contribute?](#how-can-i-contribute)
3. [Development Setup](#development-setup)
4. [Development Workflow](#development-workflow)
5. [Testing](#testing)
6. [Submitting Changes](#submitting-changes)
7. [Code Style](#code-style)
8. [Documentation](#documentation)

## Code of Conduct

We are committed to providing a welcoming and inclusive experience for everyone. Please be respectful and constructive in all interactions.

Key principles:
- Be respectful and inclusive
- Welcome newcomers and help them get started
- Focus on what is best for the community
- Show empathy towards other community members

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates.

**Good bug reports** include:
- Clear, descriptive title
- Steps to reproduce
- Expected vs actual behavior
- Your environment (OS, Go version, etc.)
- Screenshots if applicable

**Example:**

```markdown
### Bug: Neuron execution fails with special characters

**Environment:**
- OS: macOS 14.0
- Cortex version: 1.0.0
- Go version: 1.25.4

**Steps to reproduce:**
1. Create neuron with name containing `@` symbol
2. Run `cortex exec -p neuron@test`
3. Observe error

**Expected:** Neuron executes successfully
**Actual:** Error: "invalid neuron name"
```

### Suggesting Features

We welcome feature suggestions! Please:
- Check if the feature already exists or is planned
- Describe the problem you're trying to solve
- Explain your proposed solution
- Consider alternatives

**Template:**

```markdown
### Feature: AI cost estimation before generation

**Problem:**
Users want to know costs before generating neurons with AI.

**Proposed Solution:**
Add `--estimate-cost` flag to show cost before generation.

**Alternatives Considered:**
- Show cost after generation (not helpful)
- Interactive prompt (breaks automation)
```

### Contributing Code

We follow **Test-Driven Development (TDD)**:
1. Write failing test (RED)
2. Write minimum code to pass (GREEN)
3. Refactor
4. Submit PR

See [Development Workflow](#development-workflow) below.

### Improving Documentation

Documentation is as important as code! You can:
- Fix typos or clarify confusing sections
- Add examples and use cases
- Translate documentation
- Write tutorials or guides

## Development Setup

### Prerequisites

- **Go 1.25.4+** - [Install Go](https://go.dev/doc/install)
- **Node.js 24.x LTS** - For web UI development
- **Make** - Build automation
- **Git** - Version control

### Clone and Build

```bash
# 1. Fork the repository on GitHub
# 2. Clone your fork
git clone https://github.com/YOUR-USERNAME/cortex.git
cd cortex

# 3. Add upstream remote
git remote add upstream https://github.com/anoop2811/cortex.git

# 4. Install dependencies
make install-deps

# 5. Build
make build

# 6. Run tests
make test-all

# 7. Verify
./cortex --version
```

### Directory Structure

```
cortex/
├── cmd/cortex/          # CLI entry point
├── internal/            # Private packages
│   ├── neuron/          # Neuron execution
│   ├── synapse/         # Synapse orchestration
│   ├── ai/              # AI generation (future)
│   ├── api/             # REST API (future)
│   └── db/              # Database layer (future)
├── acceptance/          # Acceptance tests
│   ├── cli/             # CLI tests (Ginkgo)
│   └── web-ui/          # Web UI tests (Playwright)
├── web/                 # Web UI (future)
├── docs/                # Documentation
├── examples/            # Example neurons/synapses
└── Makefile            # Build automation
```

## Development Workflow

We use **Test-Driven Development (TDD)** with outer-loop and inner-loop cycles.

### Adding a New Feature

#### 1. Create Feature Branch

```bash
git checkout -b feature/my-awesome-feature
```

#### 2. Write Acceptance Test (Outer Loop - RED)

```go
// acceptance/cli/my_feature_test.go
package acceptance_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var _ = Describe("My Feature", Label("acceptance", "cli"), func() {
    It("should work from user's perspective", func() {
        session := RunCortex("my-command", "--flag", "value")
        Eventually(session).Should(gexec.Exit(0))
        Eventually(session.Out).Should(gbytes.Say("expected output"))
    })
})
```

```bash
make test-acceptance-cli  # FAILS (RED ✗)
```

#### 3. Write Unit Tests (Inner Loop - RED)

```go
// internal/mypackage/mycode_test.go
var _ = Describe("MyCode", Label("unit"), func() {
    It("should do something specific", func() {
        result := MyFunction()
        Expect(result).To(Equal("expected"))
    })
})
```

```bash
make test-unit  # FAILS (RED ✗)
```

#### 4. Implement Code (Inner Loop - GREEN)

```go
// internal/mypackage/mycode.go
func MyFunction() string {
    return "expected"
}
```

```bash
make test-unit  # PASSES (GREEN ✓)
```

#### 5. Wire Up CLI Command (Outer Loop - GREEN)

```go
// cmd/cortex/my_command.go
// Implement CLI command using MyFunction
```

```bash
make test-acceptance-cli  # PASSES (GREEN ✓)
```

#### 6. Refactor

```bash
make test-all  # Everything PASSES ✓
```

#### 7. Commit

```bash
git add -A
git commit -m "Add my awesome feature

- Implements feature X for use case Y
- Includes acceptance and unit tests
- Updates documentation"
```

### TDD Watch Mode

```bash
make watch  # Tests auto-run on file changes
```

This is the fastest way to develop!

## Testing

### Running Tests

```bash
# All tests
make test-all

# Only unit tests
make test-unit

# Only acceptance tests
make test-acceptance

# Specific labels
ginkgo -r --label-filter="ai && neuron"

# With coverage
make coverage
```

### Writing Tests

See **[Testing Guide](../TESTING.md)** for complete details.

**Quick example:**

```go
var _ = Describe("Neuron Execution", Label("unit"), func() {
    var neuron *Neuron

    BeforeEach(func() {
        neuron = NewNeuron("test-neuron")
    })

    Context("when neuron is valid", func() {
        It("should execute successfully", func() {
            err := neuron.Execute()
            Expect(err).NotTo(HaveOccurred())
        })
    })

    Context("when neuron is invalid", func() {
        It("should return error", func() {
            neuron.Command = ""
            err := neuron.Execute()
            Expect(err).To(HaveOccurred())
        })
    })
})
```

### Coverage Requirements

- **Minimum**: 90% code coverage
- **Enforced**: In CI/CD
- **Check**: `make coverage`

## Submitting Changes

### Before Submitting

1. ✅ All tests pass: `make test-all`
2. ✅ Code is formatted: `make fmt`
3. ✅ No linter errors: `make lint`
4. ✅ Coverage meets threshold: `make coverage`
5. ✅ Documentation updated
6. ✅ Commit messages are clear

### Pull Request Process

#### 1. Update Your Fork

```bash
git fetch upstream
git rebase upstream/main
```

#### 2. Push Changes

```bash
git push origin feature/my-awesome-feature
```

#### 3. Create Pull Request

Go to GitHub and create a PR with:

**Title:** Clear, concise description
```
Add AI neuron generation with OpenAI support
```

**Description:**
```markdown
## Summary
Implements AI-powered neuron generation using OpenAI API.

## Changes
- Add AI generation package
- Implement OpenAI provider
- Add CLI command `generate-neuron`
- Include acceptance and unit tests

## Testing
- ✅ All tests pass
- ✅ Coverage: 92%
- ✅ Tested with OpenAI API

## Documentation
- Updated user guide
- Added architecture docs
- Included examples

## Related Issues
Closes #123
```

#### 4. Address Review Feedback

- Respond to comments
- Make requested changes
- Push updates (will auto-update PR)
- Re-request review when ready

#### 5. Merge

Once approved, a maintainer will merge your PR. Thank you!

## Code Style

### Go Code

We follow standard Go conventions:

```go
// Good
func GenerateNeuron(prompt string) (*Neuron, error) {
    if prompt == "" {
        return nil, errors.New("prompt cannot be empty")
    }

    neuron := &Neuron{
        Name: sanitizeName(prompt),
    }

    return neuron, nil
}

// Bad
func generate_neuron(prompt string)(*Neuron,error){
    neuron:=&Neuron{Name:prompt}
    return neuron,nil
}
```

**Rules:**
- Use `gofmt` (run `make fmt`)
- Use meaningful variable names
- Write comments for exported functions
- Keep functions small (< 50 lines)
- Handle all errors

### Test Code

```go
// Good - Descriptive, organized
var _ = Describe("Neuron Generator", func() {
    Context("when prompt is valid", func() {
        It("should generate neuron successfully", func() {
            neuron, err := GenerateNeuron("Find which process is using port 8080 and show full command with PID")
            Expect(err).NotTo(HaveOccurred())
            Expect(neuron.Name).To(ContainSubstring("port"))
            Expect(neuron.Script).To(ContainSubstring("lsof"))
        })
    })
})

// Bad - Unclear, hard to debug
var _ = Describe("Test", func() {
    It("works", func() {
        n, _ := GenerateNeuron("test")
        Expect(n).NotTo(BeNil())
    })
})
```

### Commit Messages

```
# Good
Add AI neuron generation with OpenAI support

Implements the AI generation feature specified in #123.
Users can now generate neurons from natural language prompts.

- Add OpenAI provider integration
- Implement prompt engineering with few-shot learning
- Add cost estimation before generation
- Include comprehensive tests (95% coverage)

Closes #123

# Bad
fixed stuff
```

**Format:**
```
<type>: <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `test`: Tests only
- `refactor`: Code refactoring
- `perf`: Performance improvement

## Documentation

### When to Update Docs

Always update docs when:
- Adding new features
- Changing existing behavior
- Fixing bugs that affect usage
- Improving error messages

### What to Document

1. **User-facing changes**:
   - Update `docs/guides/user-guide.md`
   - Add examples to `examples/`
   - Update `README.md` if needed

2. **Architecture changes**:
   - Update `docs/architecture/`
   - Add diagrams if helpful

3. **API changes**:
   - Update godoc comments
   - Update `docs/specs/`

### Documentation Style

```markdown
# Good - Clear, with examples

## Generate Neuron from Prompt

Generate a debugging script from natural language.

**Usage:**
```bash
cortex generate-neuron --prompt "Find which process is using port 8080 and show the full command with PID"
```

**Options:**
- `--prompt`: Natural language description (required)
- `--provider`: AI provider (openai, anthropic, ollama)
- `--output`: Output directory (default: current)

**Example:**
```bash
cortex generate-neuron \
  --prompt "Check if PostgreSQL on port 5432 is accepting connections and can execute a simple query" \
  --provider openai \
  --output postgres-health-check
```

# Bad - Unclear, no examples

## generate-neuron

Generates neurons.

Usage: cortex generate-neuron
```

## Getting Help

- **Questions**: [GitHub Discussions](https://github.com/anoop2811/cortex/discussions)
- **Bugs**: [GitHub Issues](https://github.com/anoop2811/cortex/issues)
- **Contributing**: This document

## Recognition

Contributors are recognized in:
- GitHub contributors page
- Release notes
- Project README

Thank you for contributing to Cortex! 🎉

---

**Next Steps:**
- Read [Testing Guide](../TESTING.md)
- Check [Architecture Docs](../architecture/)
- Join [GitHub Discussions](https://github.com/anoop2811/cortex/discussions)
