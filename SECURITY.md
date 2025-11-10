# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Currently supported versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of Cortex seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### Where to Report

**Please do NOT report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to:
- **anoop2811** (repository owner)

You can find contact information in the GitHub profile.

### What to Include

Please include the following information in your report:

- Type of issue (e.g., buffer overflow, SQL injection, cross-site scripting, etc.)
- Full paths of source file(s) related to the manifestation of the issue
- The location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit it

### Response Timeline

- **Initial Response**: Within 48 hours
- **Confirmation**: Within 7 days
- **Fix Timeline**: Depends on severity
  - **Critical**: 1-7 days
  - **High**: 7-30 days
  - **Medium**: 30-90 days
  - **Low**: Next release cycle

## Security Update Process

1. Security vulnerability is reported and confirmed
2. A fix is developed in a private branch
3. Security advisory is prepared
4. Fix is released with security advisory
5. Public disclosure after users have had time to update

## Automated Security Scanning

This project uses multiple automated security tools:

### Dependabot
- **Go modules**: Weekly updates for dependencies
- **npm packages**: Weekly updates for Playwright tests
- **GitHub Actions**: Weekly updates for workflow dependencies

### Continuous Security Scanning
- **govulncheck**: Official Go vulnerability scanner
- **npm audit**: Node.js dependency vulnerability scanner
- **CodeQL**: Semantic code analysis
- **Trivy**: Container and filesystem vulnerability scanner
- **Gosec**: Go security checker

### Running Security Scans Locally

#### Go Vulnerability Check
```bash
# Install govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest

# Run scan
govulncheck ./...
```

#### NPM Audit
```bash
cd acceptance/web-ui
npm audit
npm audit fix  # Apply automatic fixes
```

#### Go Security Scanner (gosec)
```bash
# Install gosec
go install github.com/securego/gosec/v2/cmd/gosec@latest

# Run scan
gosec ./...
```

#### Trivy Scanner
```bash
# Install trivy
brew install aquasecurity/trivy/trivy  # macOS
# or
curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh

# Run scan
trivy fs .
```

## Security Best Practices

### For Contributors

1. **Never commit secrets**
   - No API keys, passwords, or tokens
   - Use environment variables
   - Review commits before pushing

2. **Validate input**
   - Sanitize all user input
   - Use parameterized queries
   - Validate file paths

3. **Handle errors securely**
   - Don't expose sensitive information in errors
   - Log security events
   - Fail securely

4. **Keep dependencies updated**
   - Review Dependabot PRs promptly
   - Test updates thoroughly
   - Monitor security advisories

5. **Use security tools**
   - Run `govulncheck` before commits
   - Enable IDE security linters
   - Review CodeQL findings

### For Users

1. **Keep Cortex updated**
   - Use the latest stable release
   - Subscribe to security advisories
   - Review changelogs

2. **Secure your environment**
   - Use environment variables for secrets
   - Limit execution permissions
   - Review generated neurons before execution

3. **Report issues**
   - Report suspicious behavior
   - Share security concerns
   - Help improve security

## Known Security Considerations

### Neuron Execution
- Neurons execute shell scripts with system permissions
- Review neuron code before execution
- Use least-privilege principles
- Avoid running as root/administrator

### AI Generation (Future)
- AI-generated code should be reviewed
- Validate generated scripts
- Use security scanning on generated code
- Don't execute untrusted generated neurons

### API Keys (Future)
- Store API keys in environment variables
- Never log API keys
- Rotate keys regularly
- Use key management systems in production

## Security Contacts

For security-related discussions:
- **GitHub Issues**: https://github.com/anoop2811/cortex/issues (for non-sensitive security improvements)
- **Private Contact**: See repository owner's GitHub profile

## Acknowledgments

We appreciate security researchers who responsibly disclose vulnerabilities. Contributors will be acknowledged in:
- Security advisories
- Release notes
- This document (with permission)

---

**Last Updated**: November 2025
**Policy Version**: 1.0
