# IP Monitor - Documentation Index

Welcome to IP Monitor's documentation! Start here to understand the project.

## Quick Navigation

- **[README.md](../README.md)** - Quick start, installation, setup
- **[ARCHITECTURE.md](../ARCHITECTURE.md)** - Full architecture guide (DDD, Clean, Hexagonal)
- **[CONTRIBUTING.md](../CONTRIBUTING.md)** - Development setup and guidelines
- **[CHANGELOG.md](../CHANGELOG.md)** - Version history and release notes
- **[LICENSE](../LICENSE)** - MIT License

## What to Read First?

**New to the project?**

1. Read [README.md](../README.md) for quick start
2. Read [ARCHITECTURE.md](../ARCHITECTURE.md) to understand how it's built
3. Read [CONTRIBUTING.md](../CONTRIBUTING.md) if you want to contribute

**Architecture deep dive?**

- [ARCHITECTURE.md](../ARCHITECTURE.md) covers:
  - 4-layer architecture (Presentation, Application, Domain, Infrastructure)
  - DDD concepts (value objects, entities, domain services, events)
  - Hexagonal architecture (ports & adapters)
  - Repositories vs Ports distinction
  - How to extend (add new adapters, interfaces)
  - Testing strategy
  - Data flow diagrams

**Want to contribute?**

- [CONTRIBUTING.md](../CONTRIBUTING.md) covers:
  - Development setup (Go 1.21+)
  - Building and testing
  - Code style
  - Pull request process
  - Project structure

**Deploying to Raspberry Pi?**

- See [README.md](../README.md) §"Installation" for:
  - Download pre-built binaries
  - systemd service and timer setup
  - Environment variable configuration

## Architecture at a Glance

```
┌──────────────────────────────────────┐
│    Presentation Layer (CLI)          │
├──────────────────────────────────────┤
│    Application Layer (Use Cases)     │
├──────────────────────────────────────┤
│    Domain Layer (Business Logic)     │
├──────────────────────────────────────┤
│    Infrastructure (Adapters, Config) │
└──────────────────────────────────────┘
```

See [ARCHITECTURE.md](../ARCHITECTURE.md) for full details.

## Development Workflow

```bash
# Clone and build
git clone https://github.com/giovanniandreuzza/ipmonitor.git
cd ipmonitor
make build

# Test
make test

# Run locally
export TELEGRAM_BOT_TOKEN="your_token"
export TELEGRAM_CHAT_ID="your_chat_id"
make run

# Build for Raspberry Pi
make build-rpi
```

## Learning Path

**Beginner**: Read README → Run locally → Try make test  
**Intermediate**: Read ARCHITECTURE → Explore code → Run tests  
**Advanced**: Review domain services → Try extending with new adapter  

---

For more details, see each documentation file!
