# IP Monitor Architecture

## Overview

IP Monitor implements **Domain-Driven Design (DDD)**, **Clean Architecture**, and **Hexagonal Architecture** (Ports & Adapters) with a 4-layer structure:

```
┌──────────────────────────────────────┐
│    Presentation Layer (CLI, HTTP)    │  ← User interfaces
├──────────────────────────────────────┤
│    Application Layer (Use Cases)     │  ← Business orchestration
├──────────────────────────────────────┤
│    Domain Layer (Business Logic)     │  ← Core, no dependencies
├──────────────────────────────────────┤
│    Infrastructure (Adapters, Config) │  ← Technical details
└──────────────────────────────────────┘
```

## Directory Structure

```
internal/
├── presentation/                      # Presentation Layer
│   └── cli/
│       └── cli_adapter.go            # CLI interface (structured logging)
│
├── application/                       # Application Layer
│   ├── usecases/
│   │   └── monitorip/
│   │       ├── monitor_ip_usecase.go # Orchestrates IP monitoring
│   │       └── dto/
│   │           └── monitor_ip_dto.go # Command & result DTOs (with validation)
│   ├── ports/
│   │   ├── public_ip_provider.go     # Fetch IP from external API
│   │   └── notification.go           # Send notifications
│   └── events/
│       └── event_bus.go              # Domain event publishing
│
├── domain/                            # Domain Layer (Pure business logic)
│   └── ip/
│       ├── valueobjects/
│       │   └── ipv4_address.go      # Self-validating IPv4 (immutable)
│       ├── entities/
│       │   └── ip_monitor_state.go  # Aggregate root
│       ├── services/
│       │   └── ip_change_detector.go # Domain service
│       ├── events/
│       │   └── ip_events.go         # Domain events
│       ├── repositories/
│       │   └── repository.go        # Persistence contract
│       └── errors/
│           └── errors.go            # Domain-specific errors
│
└── infrastructure/                    # Infrastructure Layer
    ├── adapters/
    │   ├── persistence/file/
    │   │   └── file_adapter.go  # IPRepository impl (retry logic)
    │   ├── http/ipify/
    │   │   └── ipify_adapter.go # PublicIPProvider impl (retry + timeout)
    │   └── notification/telegram/
    │       └── telegram_adapter.go # NotificationPort impl (retry + timeout)
    ├── config/
    │   └── config.go                # Env var loading + validation
    └── logging/
        └── logger.go                # Structured logging (slog)

cmd/
└── ipmonitor/
    └── main.go                      # Bootstrap & DI

test/
├── domain/
│   ├── valueobjects/
│   │   └── ipv4_address_test.go    # Value object tests
│   └── services/
│       └── ip_change_detector_test.go # Domain service tests
└── application/
    └── monitorip_usecase_test.go    # Use case tests (with mocks)
```

## Layer Responsibilities

### 1. Presentation Layer

User-facing interfaces. Currently CLI; easily extended with HTTP/gRPC.

- **CLI Adapter** (`presentation/cli/`): Command-line interaction
  - Calls use cases
  - Logs via structured slog
  - Returns success/error to user
  - Replaceable without affecting other layers

### 2. Application Layer

Orchestration, DTOs, and application-level ports (external service contracts).

- **Use Cases** (`application/usecases/`):
  - `MonitorIPUseCase`: Validates input, orchestrates domain logic, publishes events
  - Single Responsibility: IP monitoring workflow

- **Application Ports** (`application/ports/`):
  - `PublicIPProvider`: External IP API contract
  - `NotificationPort`: Notification service contract
  - Decouples use cases from infrastructure

- **DTOs** (`application/usecases/*/dto/`):
  - `MonitorIPCommand`: Input (with validation hook)
  - `MonitorIPResult`: Output (success/failure)
  - Cross-boundary data transfer

- **Event Bus** (`application/events/`):
  - Publishes domain events to subscribers
  - Supports sync and noop implementations
  - Decouples side effects from core logic

### 3. Domain Layer

Pure business logic with zero infrastructure dependencies. Defines what the business is.

- **Value Objects** (`domain/ip/valueobjects/`):
  - `IPv4Address`: Immutable, self-validating
  - Equality by value, not identity
  - Raises domain errors for invalid addresses

- **Entities** (`domain/ip/entities/`):
  - `IPMonitorState`: Aggregate root (identity + lifecycle)
  - Encapsulates business invariants
  - Methods: `CurrentIP()`, `UpdateIP()`, `HasChanged()`

- **Domain Services** (`domain/ip/services/`):
  - `IPChangeDetector`: Detects changes, emits events
  - Stateless, pure business logic
  - No infrastructure dependencies

- **Domain Events** (`domain/ip/events/`):
  - `IPChangedEvent`: Old → New IP notification
  - `IPDetectedEvent`: Initial detection notification
  - Business-meaningful, published via event bus

- **Repositories** (`domain/ip/repositories/`):
  - `IPRepository`: Persistence contract (domain language)
  - Defines operations for aggregate persistence
  - Implemented by infrastructure adapters

- **Domain Errors** (`domain/ip/errors/`):
  - `DomainError`: Typed errors with codes
  - `InvalidIPv4Error`, `EmptyIPv4Error`: Validation
  - `RepositoryLoadError`, `NotificationError`: Domain failures

### 4. Infrastructure Layer

Technical implementations: adapters, config, logging.

- **Adapters** (`infrastructure/adapters/`):
  - **File Adapter** (IPRepository):
    - Stores IP in `/tmp/public_ip.txt`
    - Load error handling
  - **Ipify Adapter** (PublicIPProvider):
    - Fetches IP from ipify.org API
    - 3-attempt retry with backoff (200ms × attempt)
    - 5-second HTTP timeout
    - Structured logging on failures
  - **Telegram Adapter** (NotificationPort):
    - Sends Telegram messages with Markdown
    - 3-attempt retry with backoff
    - 5-second HTTP timeout
    - Structured logging on failures

- **Configuration** (`infrastructure/config/`):
  - Loads env vars: `TELEGRAM_BOT_TOKEN`, `TELEGRAM_CHAT_ID` (required)
  - Optional: `LOG_LEVEL` (DEBUG|INFO|WARN|ERROR; default: INFO)
  - Validates on startup; returns error if missing

- **Logging** (`infrastructure/logging/`):
  - Uses Go's slog package (1.21+)
  - Text handler with configurable level
  - Key-value logging throughout
  - Used in adapters and CLI

## Repositories vs Ports (DDD Distinction)

| Aspect | Repositories | Ports |
|--------|--------------|-------|
| **Layer** | Domain | Application |
| **Purpose** | Persist aggregates | Integrate external services |
| **Example** | IPRepository | PublicIPProvider, NotificationPort |
| **Impl** | `file_adapter.go` | `ipify_adapter.go`, `telegram_adapter.go` |
| **Scope** | Bounded context | Application-level services |

## Dependency Rules

**All dependencies point inward:**

```
Infrastructure → Application → Domain (only)
```

1. **Domain**: Zero dependencies (purest business)
2. **Application**: Depends only on Domain
3. **Infrastructure**: Depends on Application & Domain
4. **Presentation**: Depends on Application

## Data Flow

```
CLI Adapter
    ↓
MonitorIPUseCase.Execute(cmd)
    ↓ (validate)
IPChangeDetector.DetectChange()
    ↓ (get current IP)
Ipify Adapter (3 retries, timeout)
    ↓ (load stored state)
File Adapter (handle missing)
    ↓ (domain service generates event)
IPChangedEvent / IPDetectedEvent
    ↓ (publish event)
EventBus → Telegram Adapter (3 retries, timeout)
    ↓ (save new state)
File Adapter
    ↓
Return MonitorIPResult
    ↓
CLI outputs result
```

## Key Patterns

### Value Object (IPv4Address)

- Immutable after construction
- Self-validating via `NewIPv4Address()`
- Raises `InvalidIPv4Error` on validation failure
- Equality by value

### Entity (IPMonitorState)

- Has identity (implicit: current state)
- Mutable state (UpdateIP)
- Business logic (HasChanged)
- Aggregate root for IP domain

### Domain Service (IPChangeDetector)

- Stateless
- Operates on domain objects
- Emits domain events
- No infrastructure calls

### Domain Event

- Immutable snapshot of what happened
- Carries business context (old/new IP, timestamp)
- Published via event bus
- Decouples components

### Adapter Pattern

- Implements application port
- Contains infrastructure concerns (HTTP, file I/O)
- Retry logic, timeouts, error handling
- Adapter errors wrapped in domain errors where appropriate

### Event Bus

- Simple in-memory pub/sub
- Publishes to all subscribers synchronously
- Noop variant for testing
- Decouples use case from side effects

## Testing Strategy

### Unit Tests (test/)

1. **Value Objects** (`test/domain/valueobjects/`):
   - Valid construction
   - Invalid input rejection
   - Equality checks

2. **Domain Services** (`test/domain/services/`):
   - First detection (no stored state)
   - IP change detection (state changed)
   - No change (same IP)

3. **Use Cases** (`test/application/`):
   - Full flow with mock adapters
   - IP changed and notified
   - No change, no notification
   - Error handling

### Mock Adapters

- `mockIPProvider`: Returns configurable IP
- `mockRepo`: In-memory state storage
- `mockNotifier`: Tracks if called

## Configuration

Required environment variables:

```bash
export TELEGRAM_BOT_TOKEN="YOUR_BOT_TOKEN"
export TELEGRAM_CHAT_ID="YOUR_CHAT_ID"
```

Optional:

```bash
export LOG_LEVEL="DEBUG"  # DEBUG, INFO, WARN, ERROR (default: INFO)
```

Usage:

```bash
make run      # Build and run locally
make build-rpi # Build for Raspberry Pi
make test     # Run all tests
```

## Extending the Architecture

### Add Email Notification

1. Create `infrastructure/adapters/secondary/notification/email/email_adapter.go`
2. Implement `application/ports.NotificationPort`
3. Wire in `cmd/ipmonitor/main.go`
4. **No changes to domain/application core**

### Add HTTP API

1. Create `internal/presentation/http/http_adapter.go`
2. Call `monitorIPUseCase.Execute()`
3. Return JSON responses
4. **No changes to domain/application core**

### Swap Storage (File → Database)

1. Create `infrastructure/adapters/secondary/persistence/postgres/postgres_adapter.go`
2. Implement `domain/ip/repositories.IPRepository`
3. Wire in `cmd/ipmonitor/main.go`
4. **No changes to domain/application core**

## Architecture Strengths

✅ **Testability**: Domain/application logic testable without infrastructure  
✅ **Flexibility**: Swap implementations easily  
✅ **Maintainability**: Clear concerns, obvious where code belongs  
✅ **Business-Focused**: Domain layer speaks business language  
✅ **Framework-Independent**: Pure Go, no tight framework coupling  
✅ **Resilience**: Retry logic and error handling at infrastructure layer  
✅ **Observability**: Structured logging throughout  
✅ **Scalable**: Ready for complex business logic growth  

## Learning Resources

- **Domain-Driven Design** - Eric Evans
- **Clean Architecture** - Robert C. Martin
- **Implementing DDD** - Vaughn Vernon
- **Hexagonal Architecture** - Alistair Cockburn
