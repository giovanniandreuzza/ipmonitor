# Contributing to IP Monitor

Thank you for your interest in contributing to IP Monitor! This document provides guidelines for contributing to this project.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/ipmonitor.git`
3. Create a new branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test your changes
6. Commit your changes: `git commit -am 'Add some feature'`
7. Push to the branch: `git push origin feature/your-feature-name`
8. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Make (optional, but recommended)

### Building the Project

```bash
# Build for current platform
make build

# Build for Raspberry Pi
make build-rpi

# Build for all platforms
make build-all
```

### Running Tests

```bash
make test
```

### Running Locally

```bash
# Set environment variables
export TELEGRAM_BOT_TOKEN="your_token"
export TELEGRAM_CHAT_ID="your_chat_id"

# Run the application
make run
# or
go run ./cmd/ipmonitor
```

## Code Style

- Follow standard Go formatting guidelines
- Run `make fmt` to format your code
- Run `make vet` to check for common mistakes
- Keep functions small and focused
- Add comments for exported functions and types
- Write tests for new functionality

## Project Structure

The project follows Domain-Driven Design (DDD) and Hexagonal Architecture:

```
internal/
├── presentation/               # User interfaces (CLI, future: HTTP)
│   └── cli/
│       └── cli_adapter.go     # CLI interface
│
├── application/                # Use cases and application services
│   ├── usecases/
│   │   └── monitorip/
│   │       ├── monitor_ip_usecase.go
│   │       └── dto/
│   ├── ports/                  # External service contracts
│   │   ├── public_ip_provider.go
│   │   └── notification.go
│   └── events/
│       └── event_bus.go       # Domain event publisher
│
├── domain/                     # Pure business logic (no dependencies)
│   └── ip/
│       ├── valueobjects/      # IPv4Address (immutable)
│       ├── entities/          # IPMonitorState (aggregate root)
│       ├── services/          # IPChangeDetector (domain logic)
│       ├── events/            # Domain events
│       ├── repositories/      # Persistence contracts
│       └── errors/            # Domain-specific errors
│
└── infrastructure/             # Technical implementations
    ├── adapters/
    │   └── secondary/
    │       ├── persistence/file/
    │       ├── http/ipify/
    │       └── notification/telegram/
    ├── config/                 # Configuration loading
    └── logging/                # Structured logging

cmd/
└── ipmonitor/
    └── main.go                # Bootstrap & dependency injection

test/
├── domain/
│   ├── valueobjects/
│   └── services/
└── application/
```

## Development Workflow

1. **Understand DDD**: Read [ARCHITECTURE.md](../ARCHITECTURE.md)
2. **Make changes**: Edit files in `internal/`
3. **Add tests**: Add corresponding tests in `test/`
4. **Run tests**: `make test`
5. **Format code**: `make fmt`
6. **Build**: `make build`

## Code Style

- Follow standard Go formatting guidelines
- Run `make fmt` to format your code
- Run `make vet` to check for common mistakes
- Keep functions small and focused
- Add comments for exported functions and types
- Write tests for new functionality
- Use structured logging (slog) instead of fmt.Printf or log.Print

## Testing

```bash
# Run all tests
make test

# Run specific test file
go test ./test/domain/valueobjects/

# Run with verbose output
go test -v ./...
```

Tests include:

- Domain layer (value objects, services)
- Application layer (use cases with mocks)
- Infrastructure adapters (with integration tests)

## Adding New Features

### To add a new notification channel (e.g., Email)

1. Create adapter: `internal/infrastructure/adapters/secondary/notification/email/email_adapter.go`
2. Implement `application/ports.NotificationPort`
3. Wire in `cmd/ipmonitor/main.go`
4. Add tests in `test/application/`

### To add an HTTP API

1. Create adapter: `internal/presentation/http/http_adapter.go`
2. Call `monitorIPUseCase.Execute()`
3. Return JSON responses
4. Wire in `cmd/ipmonitor/main.go`

## Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Make your changes
4. Add tests
5. Run `make test` and `make fmt`
6. Commit with meaningful messages
7. Push to your fork
8. Open a pull request with a clear description

## Questions?

See [ARCHITECTURE.md](../ARCHITECTURE.md) for detailed architecture explanation.

## Commit Message Guidelines

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

- `feat:` - A new feature
- `fix:` - A bug fix
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks

Examples:

```
feat: add support for IPv6 monitoring
fix: handle network timeout errors gracefully
docs: update installation instructions
```

## CI/CD and Automation

### GitHub Actions Workflows

The project uses GitHub Actions for continuous integration and deployment:

#### **CI Workflow** (`.github/workflows/ci.yml`)

Runs on every push and pull request to `main` or `develop`:

- **Test**: Runs all unit tests with race detection and coverage
- **Build**: Cross-platform build verification (Linux, macOS, ARM)
- **Build for Raspberry Pi**: Specific ARMv7 build
- **Lint**: Runs golangci-lint for code quality checks
- **Mod Tidy**: Ensures go.mod and go.sum are clean

#### **Release Workflow** (`.github/workflows/release.yml`)

Triggered when pushing a tag like `v1.0.0`:

- Runs all tests
- Builds binaries for multiple platforms:
  - Linux: AMD64, ARM64, ARMv7
  - macOS: AMD64, ARM64 (M1/M2)
- Generates checksums
- Creates GitHub release with changelog
- Uploads all binaries as release artifacts

#### **Docker Workflow** (`.github/workflows/docker.yml`)

Builds and publishes Docker images:

- Multi-architecture builds (AMD64, ARM64, ARMv7)
- Pushes to GitHub Container Registry (ghcr.io)
- Tags with version, branch, and commit SHA
- Runs on push to `main` or new tags

#### **Security Workflow** (`.github/workflows/security.yml`)

Security scanning that runs on push/PR and weekly:

- **gosec**: Static security analysis
- **govulncheck**: Vulnerability scanning for dependencies
- Uploads results to GitHub Security tab

#### **Dependabot** (`.github/dependabot.yml`)

Automated dependency updates:

- Weekly checks for Go module updates
- Weekly checks for GitHub Actions updates
- Auto-labels PRs for easy review

### Creating a Release

To create a new release:

```bash
# Ensure all changes are committed and pushed
git add .
git commit -m "Prepare release v1.2.3"
git push

# Create and push a tag
git tag v1.2.3
git push origin v1.2.3
```

The Release workflow will automatically:

1. Run all tests
2. Build binaries for all platforms
3. Create a GitHub release with changelog
4. Upload binaries and checksums

### Docker Images

Docker images are automatically built and published to:

```
ghcr.io/giovanniandreuzza/ipmonitor:latest
ghcr.io/giovanniandreuzza/ipmonitor:main
ghcr.io/giovanniandreuzza/ipmonitor:v1.2.3
ghcr.io/giovanniandreuzza/ipmonitor:sha-abc1234
```

Images support multiple architectures:

- linux/amd64
- linux/arm64
- linux/arm/v7 (Raspberry Pi)

### Code Quality

Before submitting a PR, ensure code quality:

```bash
# Run linter locally
golangci-lint run

# Format code
gofmt -s -w .

# Run security scanner
gosec ./...

# Check for vulnerabilities
govulncheck ./...
```

Install tools:

```bash
# golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# gosec
go install github.com/securego/gosec/v2/cmd/gosec@latest

# govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest
```

## Pull Request Guidelines

- Keep pull requests focused on a single feature or fix
- Update documentation as needed
- Add tests for new functionality
- Ensure all tests pass
- Update the CHANGELOG.md if applicable
- Reference any related issues in the PR description

## Reporting Issues

When reporting issues, please include:

- A clear and descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Environment details (OS, Go version, etc.)
- Relevant logs or error messages

## Feature Requests

We welcome feature requests! Please:

- Check if the feature has already been requested
- Clearly describe the feature and its use case
- Explain why this feature would be useful
- Provide examples if possible

## Questions?

If you have questions, feel free to:

- Open an issue with the `question` label
- Start a discussion in the GitHub Discussions tab

## License

By contributing to IP Monitor, you agree that your contributions will be licensed under the MIT License.

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers and help them get started
- Focus on constructive feedback
- Respect differing viewpoints and experiences

Thank you for contributing! 🎉
