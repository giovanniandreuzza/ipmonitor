# GitHub Actions CI/CD

This document provides an overview of all automated workflows in this repository.

## Workflows Overview

### 🧪 CI Workflow

**File**: `.github/workflows/ci.yml`  
**Triggers**: Push/PR to `main` or `develop`  
**Purpose**: Continuous Integration

**Jobs**:

- **test**: Run unit tests with race detection and coverage reporting to Codecov
- **build**: Cross-platform build verification (Linux/macOS × AMD64/ARM64)
- **build-rpi**: Raspberry Pi ARMv7 build verification
- **lint**: Code quality checks with golangci-lint
- **mod-tidy**: Verify go.mod/go.sum are up to date

### 🚀 Release Workflow

**File**: `.github/workflows/release.yml`  
**Triggers**: Tag push matching `v*.*.*` (e.g., `v1.0.0`)  
**Purpose**: Automated releases with multi-platform binaries

**Jobs**:

- Run all tests
- Build binaries for:
  - Linux: AMD64, ARM64, ARMv7 (Raspberry Pi)
  - macOS: AMD64, ARM64 (M1/M2)
- Generate SHA256 checksums
- Create GitHub release with auto-generated changelog
- Upload all binaries and checksums as release artifacts

**Creating a release**:

```bash
git tag v1.2.3
git push origin v1.2.3
```

### 🐳 Docker Workflow

**File**: `.github/workflows/docker.yml`  
**Triggers**:

- Push to `main` branch
- Tag push matching `v*.*.*`
- Pull requests to `main`

**Purpose**: Build and publish multi-architecture Docker images

**Jobs**:

- Build images for: linux/amd64, linux/arm64, linux/arm/v7
- Push to GitHub Container Registry (ghcr.io)
- Tag with: version, branch, commit SHA, `latest`
- Use layer caching for faster builds

**Image locations**:

```
ghcr.io/giovanniandreuzza/ipmonitor:latest
ghcr.io/giovanniandreuzza/ipmonitor:main
ghcr.io/giovanniandreuzza/ipmonitor:v1.2.3
ghcr.io/giovanniandreuzza/ipmonitor:sha-abc1234
```

### 🔒 Security Workflow

**File**: `.github/workflows/security.yml`  
**Triggers**:

- Push/PR to `main` or `develop`
- Weekly schedule (Sundays at midnight)

**Purpose**: Automated security scanning

**Jobs**:

- **gosec**: Static security analysis, uploads to GitHub Security tab
- **govulncheck**: Scan for known vulnerabilities in dependencies

### 🤖 Dependabot

**File**: `.github/dependabot.yml`  
**Schedule**: Weekly updates  
**Purpose**: Automated dependency updates

**Updates**:

- Go modules (all grouped together)
- GitHub Actions versions

## Badges

Add these badges to your README to show workflow status:

```markdown
[![CI](https://github.com/giovanniandreuzza/ipmonitor/actions/workflows/ci.yml/badge.svg)](https://github.com/giovanniandreuzza/ipmonitor/actions/workflows/ci.yml)
[![Release](https://github.com/giovanniandreuzza/ipmonitor/actions/workflows/release.yml/badge.svg)](https://github.com/giovanniandreuzza/ipmonitor/actions/workflows/release.yml)
[![Security](https://github.com/giovanniandreuzza/ipmonitor/actions/workflows/security.yml/badge.svg)](https://github.com/giovanniandreuzza/ipmonitor/actions/workflows/security.yml)
[![Docker](https://github.com/giovanniandreuzza/ipmonitor/actions/workflows/docker.yml/badge.svg)](https://github.com/giovanniandreuzza/ipmonitor/actions/workflows/docker.yml)
```

## Local Testing

Before pushing, test your changes locally:

```bash
# Run tests
go test -v -race ./test/...

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build ./cmd/ipmonitor
GOOS=linux GOARCH=arm64 go build ./cmd/ipmonitor
GOOS=linux GOARCH=arm GOARM=7 go build ./cmd/ipmonitor

# Run linter
golangci-lint run

# Security scan
gosec ./...

# Vulnerability check
govulncheck ./...

# Format check
gofmt -s -l .

# Mod tidy check
go mod tidy
git diff --exit-code go.mod go.sum
```

## Secrets Required

For workflows to function properly, ensure these secrets are configured in GitHub:

| Secret | Purpose | Where to Get |
|--------|---------|--------------|
| `GITHUB_TOKEN` | Automatic releases, Docker push | Auto-provided by GitHub |
| `CODECOV_TOKEN` | (Optional) Upload coverage | [codecov.io](https://codecov.io) |

No manual secrets needed - `GITHUB_TOKEN` is automatically provided!

## Workflow Permissions

Current permissions configured:

- **release.yml**: `contents: write` (create releases)
- **docker.yml**: `packages: write` (push to GHCR)
- **security.yml**: `security-events: write` (upload SARIF)

## Tips

1. **Fast Feedback**: The CI workflow provides the fastest feedback on code changes
2. **Release Process**: Just tag and push - everything else is automated
3. **Docker Images**: Multi-arch images work on Raspberry Pi, servers, and dev machines
4. **Security**: Weekly scans keep dependencies secure
5. **Code Quality**: golangci-lint enforces best practices

## Troubleshooting

**Build fails on push?**

- Check the CI workflow logs in the Actions tab
- Run tests locally: `go test ./test/...`
- Run linter locally: `golangci-lint run`

**Release not created?**

- Ensure tag matches `v*.*.*` format (e.g., `v1.0.0`)
- Check release workflow logs
- Verify tests passed

**Docker image not found?**

- Check if workflow completed successfully
- Verify you're using the correct image path: `ghcr.io/giovanniandreuzza/ipmonitor`
- Ensure you're authenticated: `docker login ghcr.io`

**Security alerts?**

- Check the Security tab for vulnerability details
- Update dependencies: Review Dependabot PRs
- Run `govulncheck ./...` locally
