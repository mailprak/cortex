# Security Guide for Cortex

This guide explains how Cortex maintains security and how you can contribute to keeping the project secure.

## Table of Contents

- [Security Overview](#security-overview)
- [Automated Security Scanning](#automated-security-scanning)
- [Dependency Management](#dependency-management)
- [Running Security Scans Locally](#running-security-scans-locally)
- [Security Best Practices](#security-best-practices)
- [Reporting Vulnerabilities](#reporting-vulnerabilities)

## Security Overview

Cortex implements multiple layers of security:

1. **Automated Dependency Updates** - Dependabot monitors and updates dependencies weekly
2. **Continuous Security Scanning** - Multiple tools scan for vulnerabilities on every commit
3. **Code Analysis** - CodeQL and Gosec analyze code for security issues
4. **Container Scanning** - Trivy scans for vulnerabilities in dependencies and filesystems
5. **Manual Review** - Security-sensitive changes require review

## Automated Security Scanning

### GitHub Actions Workflows

Security scans run automatically on:
- Every push to `main` and `develop` branches
- Every pull request
- Daily at 2 AM UTC
- On-demand via workflow dispatch

### Scanning Tools

#### 1. govulncheck (Go Vulnerabilities)
**What it does**: Scans Go code for known vulnerabilities in dependencies

**When it runs**: On every commit and daily

**Configuration**: `.github/workflows/security-scan.yml`

**Example output**:
```
No vulnerabilities found.
```

#### 2. npm audit (Node.js Vulnerabilities)
**What it does**: Scans npm dependencies for known vulnerabilities

**When it runs**: On every commit and daily

**Severity levels**:
- **Critical**: Immediate action required
- **High**: Fix within 7 days
- **Moderate**: Fix within 30 days
- **Low**: Fix in next release

#### 3. CodeQL (Semantic Analysis)
**What it does**: Analyzes code for security vulnerabilities and coding errors

**When it runs**: On every commit and pull request

**Languages**: Go (currently), JavaScript (future)

**Queries**: `security-extended`, `security-and-quality`

#### 4. Trivy (Container & Filesystem Scanning)
**What it does**: Scans for vulnerabilities in:
- Dependencies (Go modules, npm packages)
- Container images
- Filesystem

**When it runs**: On every commit and daily

**Severity levels**: CRITICAL, HIGH, MEDIUM

#### 5. Gosec (Go Security)
**What it does**: Analyzes Go source code for security problems

**When it runs**: On every commit and daily

**Checks for**:
- SQL injection
- Command injection
- Path traversal
- Insecure cryptography
- Hardcoded credentials

#### 6. Dependency Review (Pull Requests)
**What it does**: Reviews dependency changes in pull requests

**When it runs**: On every pull request

**Blocks PRs with**:
- Moderate or higher severity vulnerabilities
- AGPL-3.0 or GPL-3.0 licensed dependencies

## Dependency Management

### Dependabot Configuration

Dependabot automatically creates pull requests for:

**Go modules** (`go.mod`):
- **Schedule**: Weekly on Monday at 9 AM UTC
- **Grouping**: Related packages grouped together (e.g., Ginkgo/Gomega)
- **Auto-update**: Security patches applied automatically

**npm packages** (`acceptance/web-ui/package.json`):
- **Schedule**: Weekly on Monday at 9 AM UTC
- **Grouping**: Playwright and type definitions grouped
- **Versioning**: Uses `increase` strategy for compatible updates

**GitHub Actions** (`.github/workflows/*.yml`):
- **Schedule**: Weekly on Monday at 9 AM UTC
- **Updates**: Action versions to latest stable

### Reviewing Dependabot PRs

1. **Check PR title**: Should indicate package and version change
2. **Review changelog**: Click on release notes link
3. **Check CI status**: All tests must pass
4. **Review diff**: Look for breaking changes
5. **Merge or close**: Merge if safe, close if unnecessary

**Example Dependabot PR**:
```
deps(go): bump github.com/onsi/ginkgo/v2 from 2.27.2 to 2.28.0
```

## Running Security Scans Locally

### Before Committing

```bash
# 1. Run Go vulnerability check
govulncheck ./...

# 2. Run Go security scanner
gosec ./...

# 3. Verify Go modules
go mod verify

# 4. Run tests
make test-all
```

### Installation

#### govulncheck
```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
```

#### gosec
```bash
go install github.com/securego/gosec/v2/cmd/gosec@latest
```

#### Trivy (macOS)
```bash
brew install aquasecurity/trivy/trivy
```

#### Trivy (Linux)
```bash
curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh
```

### Running Scans

#### Full Security Audit
```bash
# Run all security checks
make security-scan

# Or manually:
govulncheck ./...
gosec ./...
trivy fs .
go mod verify
```

#### Go Vulnerabilities Only
```bash
govulncheck ./...
```

#### npm Audit (if package-lock.json exists)
```bash
cd acceptance/web-ui
npm audit

# Fix automatically
npm audit fix

# Fix including breaking changes
npm audit fix --force
```

#### Container Security
```bash
# Scan filesystem
trivy fs .

# Scan Docker image (future)
trivy image cortex:latest
```

## Security Best Practices

### For Developers

#### 1. Never Commit Secrets
```bash
# ❌ BAD
apiKey := "sk-1234567890abcdef"

# ✅ GOOD
apiKey := os.Getenv("OPENAI_API_KEY")
```

**Prevention**:
- Use `.gitignore` for sensitive files
- Use environment variables
- Review diffs before committing

#### 2. Validate All Input
```go
// ❌ BAD
func executeCommand(userInput string) {
    cmd := exec.Command("sh", "-c", userInput)
    cmd.Run()
}

// ✅ GOOD
func executeCommand(userInput string) error {
    // Validate input
    if !isValidCommand(userInput) {
        return errors.New("invalid command")
    }

    // Use explicit arguments
    cmd := exec.Command("ls", "-la", userInput)
    return cmd.Run()
}
```

#### 3. Handle Errors Securely
```go
// ❌ BAD
if err != nil {
    log.Printf("Database error: %v", err) // May expose connection string
}

// ✅ GOOD
if err != nil {
    log.Printf("Database connection failed")
    log.Debug().Err(err).Msg("Database error details") // Debug logs only
}
```

#### 4. Use Least Privilege
```yaml
# neuron.yaml
---
name: check_disk_space
type: check  # Read-only, not mutate
# ... no sudo/elevated permissions needed
```

#### 5. Review AI-Generated Code
```bash
# When AI generates neurons, always review before executing
cortex generate-neuron "check disk space"

# Review generated files
cat neurons/check_disk_space/run.sh

# Look for:
# - Suspicious commands (rm -rf, curl | bash)
# - Hardcoded credentials
# - Unnecessary sudo
```

### For Users

#### 1. Keep Cortex Updated
```bash
# Check for updates
cortex --version

# Update to latest
# (depends on installation method)
```

#### 2. Review Neurons Before Execution
```bash
# Always review neuron code
cat path/to/neuron/run.sh

# Understand what it does
# Check for destructive commands
```

#### 3. Use Environment Variables
```bash
# ❌ BAD
export OPENAI_API_KEY="sk-1234567890abcdef"
cortex generate-neuron "..."

# ✅ GOOD (use .env file or secure vault)
# Store in ~/.cortex/.env (not committed to git)
```

## Reporting Vulnerabilities

### Where to Report

**DO NOT** create public GitHub issues for security vulnerabilities.

Instead:
1. Email the repository owner (see GitHub profile)
2. Include detailed information (see [SECURITY.md](../../SECURITY.md))
3. Wait for acknowledgment (within 48 hours)

### What to Include

- Type of vulnerability
- Affected component/file
- Steps to reproduce
- Impact assessment
- Suggested fix (if any)

### Response Timeline

- **Critical**: 1-7 days
- **High**: 7-30 days
- **Medium**: 30-90 days
- **Low**: Next release

## Security Checklist

### Before Committing
- [ ] Run `govulncheck ./...`
- [ ] Run `gosec ./...`
- [ ] No hardcoded secrets
- [ ] All tests pass
- [ ] Code reviewed

### Before Releasing
- [ ] All security scans pass
- [ ] Dependencies updated
- [ ] Changelog includes security fixes
- [ ] Security advisory if needed
- [ ] CVE assigned if applicable

### Monthly Review
- [ ] Review Dependabot PRs
- [ ] Check security advisories
- [ ] Update dependencies
- [ ] Review access logs
- [ ] Rotate secrets

## Resources

### Documentation
- [Security Policy (SECURITY.md)](../../SECURITY.md)
- [Contributing Guide](contributing.md)
- [Testing Guide](../TESTING.md)

### Tools
- [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)
- [gosec](https://github.com/securego/gosec)
- [Trivy](https://github.com/aquasecurity/trivy)
- [CodeQL](https://codeql.github.com/)
- [Dependabot](https://docs.github.com/en/code-security/dependabot)

### External Resources
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [CWE Top 25](https://cwe.mitre.org/top25/)
- [Go Security Policy](https://go.dev/security/policy)

---

**Last Updated**: November 2025
**Maintained By**: Cortex Security Team
