# IP Monitor

[![CI](https://github.com/giovanniandreuzza/ipmonitor/actions/workflows/ci.yml/badge.svg)](https://github.com/giovanniandreuzza/ipmonitor/actions/workflows/ci.yml)
[![Release](https://github.com/giovanniandreuzza/ipmonitor/actions/workflows/release.yml/badge.svg)](https://github.com/giovanniandreuzza/ipmonitor/actions/workflows/release.yml)
[![Security](https://github.com/giovanniandreuzza/ipmonitor/actions/workflows/security.yml/badge.svg)](https://github.com/giovanniandreuzza/ipmonitor/actions/workflows/security.yml)
[![Docker](https://github.com/giovanniandreuzza/ipmonitor/actions/workflows/docker.yml/badge.svg)](https://github.com/giovanniandreuzza/ipmonitor/actions/workflows/docker.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/giovanniandreuzza/ipmonitor)](https://goreportcard.com/report/github.com/giovanniandreuzza/ipmonitor)
[![codecov](https://codecov.io/gh/giovanniandreuzza/ipmonitor/branch/main/graph/badge.svg)](https://codecov.io/gh/giovanniandreuzza/ipmonitor)
[![License](https://img.shields.io/github/license/giovanniandreuzza/ipmonitor)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/giovanniandreuzza/ipmonitor)](go.mod)

A lightweight tool that monitors your public IP address and sends Telegram notifications whenever it changes. Perfect for Raspberry Pi and other Linux systems.

## Features

- ✅ Monitors your public IPv4 address
- ✅ Sends formatted Telegram notifications when IP changes
- ✅ Runs automatically every 5 minutes (configurable)
- ✅ Persists across reboots with systemd
- ✅ Lightweight and efficient (~7MB binary)
- ✅ No external dependencies at runtime
- ✅ Clean architecture with domain-driven design

## Quick Start

### Download Pre-built Binary

Download the latest release for your platform:

```bash
# For Raspberry Pi (64-bit)
wget https://github.com/giovanniandreuzza/ipmonitor/releases/latest/download/ipmonitor-linux-arm64
chmod +x ipmonitor-linux-arm64
sudo mv ipmonitor-linux-arm64 /usr/local/bin/ipmonitor

# For Raspberry Pi (32-bit)
wget https://github.com/giovanniandreuzza/ipmonitor/releases/latest/download/ipmonitor-linux-armv7
chmod +x ipmonitor-linux-armv7
sudo mv ipmonitor-linux-armv7 /usr/local/bin/ipmonitor

# For Linux AMD64
wget https://github.com/giovanniandreuzza/ipmonitor/releases/latest/download/ipmonitor-linux-amd64
chmod +x ipmonitor-linux-amd64
sudo mv ipmonitor-linux-amd64 /usr/local/bin/ipmonitor
```

### Build from Source

```bash
git clone https://github.com/giovanniandreuzza/ipmonitor.git
cd ipmonitor
make build-rpi  # or make build for current platform
```

## Prerequisites

- Raspberry Pi or any Linux system
- Active internet connection
- Go 1.21+ (if building from source)
- Telegram Bot Token (from @BotFather on Telegram)
- Telegram Chat ID (your user ID or group ID)

## Installation

### Quick Install (Recommended)

Download and install the binary:

```bash
# For Raspberry Pi (64-bit)
wget https://github.com/giovanniandreuzza/ipmonitor/releases/latest/download/ipmonitor-linux-arm64
chmod +x ipmonitor-linux-arm64
sudo mv ipmonitor-linux-arm64 /usr/local/bin/ipmonitor

# For Raspberry Pi (32-bit ARMv7)
wget https://github.com/giovanniandreuzza/ipmonitor/releases/latest/download/ipmonitor-linux-armv7
chmod +x ipmonitor-linux-armv7
sudo mv ipmonitor-linux-armv7 /usr/local/bin/ipmonitor
```

### Systemd Setup (Linux/Raspberry Pi)

For full systemd installation instructions, see [deployments/systemd/README.md](deployments/systemd/README.md).

**Quick setup:**

```bash
# Copy systemd files
sudo cp deployments/systemd/ipmonitor.service /etc/systemd/system/
sudo cp deployments/systemd/ipmonitor.timer /etc/systemd/system/

# Create environment file
sudo cp deployments/systemd/ipmonitor.env.example /etc/default/ipmonitor
sudo nano /etc/default/ipmonitor  # Edit with your credentials
sudo chmod 600 /etc/default/ipmonitor

# Enable and start
sudo systemctl daemon-reload
sudo systemctl enable --now ipmonitor.timer
```

**Configuration:**

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `TELEGRAM_BOT_TOKEN` | Yes | — | Bot token from @BotFather |
| `TELEGRAM_CHAT_ID` | Yes | — | Chat ID from @userinfobot |
| `LOG_LEVEL` | No | `INFO` | `DEBUG`, `INFO`, `WARN`, `ERROR` |

### Docker Setup

See Docker documentation in the root directory or run:

```bash
docker pull ghcr.io/giovanniandreuzza/ipmonitor:latest
docker run --rm \
  -e TELEGRAM_BOT_TOKEN="your_token" \
  -e TELEGRAM_CHAT_ID="your_id" \
  ghcr.io/giovanniandreuzza/ipmonitor:latest
```

Full Docker setup with docker-compose available in root `docker-compose.yml`.

## Verification

### Check if the timer is running

```bash
sudo systemctl status ipmonitor.timer
```

You should see: `Active: active (waiting)`

### View the next scheduled run

```bash
sudo systemctl list-timers ipmonitor.timer
```

### Check the logs

```bash
sudo journalctl -u ipmonitor -f
```

The `-f` flag shows live logs (press `Ctrl+C` to exit).

## How It Works

1. **Initial Check**: When first run, it stores your current IP in `/tmp/public_ip.txt`
2. **Periodic Checks**: Every 5 minutes, the timer triggers the service
3. **Comparison**: The tool compares your current IP with the stored IP
4. **Notification**: If the IP has changed:
   - A Telegram message is sent
   - The new IP is stored
5. **No Change**: If the IP hasn't changed, nothing happens (no notification)

## Example Telegram Messages

**First run (Initial IP detection):**

```
🌐 IP Monitor Alert

✅ Initial IP detected: 203.0.113.45
```

**When IP changes:**

```
🌐 IP Monitor Alert

🔄 Your public IP has changed:

Old IP: 203.0.113.45
New IP: 203.0.113.99
```

## Troubleshooting

### Timer won't start

```bash
sudo systemctl daemon-reload
sudo systemctl start ipmonitor.timer
```

### Check if the binary exists

```bash
ls -lh ~/ipmonitor
```

### Test the binary manually

```bash
TELEGRAM_BOT_TOKEN="your_token" TELEGRAM_CHAT_ID="your_id" ~/ipmonitor
```

### No notifications received

- Verify your tokens are correct in the service file:

  ```bash
  sudo cat /etc/systemd/system/ipmonitor.service
  ```

- Check if Telegram bot is working by sending a message to your bot
- View logs for error messages:

  ```bash
  sudo journalctl -u ipmonitor -n 20
  ```

### View last 20 log entries

```bash
sudo journalctl -u ipmonitor -n 20
```

## Common Commands

| Command | Description |
|---------|-------------|
| `sudo systemctl start ipmonitor.timer` | Start the timer |
| `sudo systemctl stop ipmonitor.timer` | Stop the timer |
| `sudo systemctl restart ipmonitor.timer` | Restart the timer |
| `sudo systemctl status ipmonitor.timer` | Check timer status |
| `sudo systemctl list-timers ipmonitor.timer` | View next run time |
| `sudo journalctl -u ipmonitor -f` | Live logs (press Ctrl+C to exit) |
| `sudo journalctl -u ipmonitor -n 50` | Last 50 log entries |

## Changing the Check Interval

To change from 5 minutes to a different interval:

1. Edit the timer file:

   ```bash
   sudo nano /etc/systemd/system/ipmonitor.timer
   ```

2. Change this line:

   ```ini
   OnUnitActiveSec=5min
   ```

   Examples:
   - Every 1 minute: `OnUnitActiveSec=1min`
   - Every 10 minutes: `OnUnitActiveSec=10min`
   - Every hour: `OnUnitActiveSec=1h`

3. Reload and restart:

   ```bash
   sudo systemctl daemon-reload
   sudo systemctl restart ipmonitor.timer
   ```

## Updating the Tool

Download and replace with the latest version:

```bash
# Download latest release
wget https://github.com/giovanniandreuzza/ipmonitor/releases/latest/download/ipmonitor-linux-arm64 -O /tmp/ipmonitor
chmod +x /tmp/ipmonitor
sudo mv /tmp/ipmonitor /usr/local/bin/ipmonitor

# The service will automatically use the new binary on the next run
```

## Development

### Docker

Run with Docker:

```bash
# Pull from GitHub Container Registry
docker pull ghcr.io/giovanniandreuzza/ipmonitor:latest

# Run with environment variables
docker run --rm \
  -e TELEGRAM_BOT_TOKEN="your_bot_token" \
  -e TELEGRAM_CHAT_ID="your_chat_id" \
  -e LOG_LEVEL="info" \
  -v ipmonitor-data:/data \
  ghcr.io/giovanniandreuzza/ipmonitor:latest
```

Build locally:

```bash
# Build image
docker build -t ipmonitor .

# Run container
docker run --rm \
  -e TELEGRAM_BOT_TOKEN="your_bot_token" \
  -e TELEGRAM_CHAT_ID="your_chat_id" \
  ipmonitor
```

Docker Compose:

```yaml
version: '3.8'
services:
  ipmonitor:
    image: ghcr.io/giovanniandreuzza/ipmonitor:latest
    environment:
      - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
      - TELEGRAM_CHAT_ID=${TELEGRAM_CHAT_ID}
      - LOG_LEVEL=info
    volumes:
      - ipmonitor-data:/data
    restart: unless-stopped

volumes:
  ipmonitor-data:
```

### Building from Source

```bash
# Clone the repository
git clone https://github.com/giovanniandreuzza/ipmonitor.git
cd ipmonitor

# Build for current platform
make build

# Build for Raspberry Pi
make build-rpi

# Build for all platforms
make build-all

# Run tests
make test
```

## Architecture

IP Monitor is built with **Domain-Driven Design (DDD)**, **Clean Architecture**, and **Hexagonal Architecture** principles.

### Layers

```
Presentation Layer (CLI)
    ↓
Application Layer (Use Cases)
    ↓
Domain Layer (Business Logic)
    ↓
Infrastructure Layer (Adapters, Config, Logging)
```

### Key Features

- **Value Objects**: Immutable IPv4 address with validation
- **Entities**: IP monitoring state aggregate
- **Domain Services**: Change detection logic
- **Domain Events**: IP changed/detected events
- **Repositories**: Persistence abstraction
- **Ports & Adapters**: Dependency inversion pattern
- **Structured Logging**: slog with configurable levels
- **Resilience**: Retry logic with exponential backoff for external services

See [ARCHITECTURE.md](ARCHITECTURE.md) for a complete architecture guide.

### Project Structure

```
internal/
├── presentation/        # CLI interface
├── application/         # Use cases and ports
├── domain/             # Business logic (entities, services, events)
└── infrastructure/     # Adapters, config, logging

cmd/
└── ipmonitor/          # Application entry point

test/
├── domain/             # Value object and domain service tests
└── application/        # Use case tests with mocks
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.

## Testing

Run unit tests:

```bash
make test
```

Tests include:

- Value object validation (IPv4Address)
- Domain service logic (change detection)
- Use case workflows (with mock adapters)

All tests pass with 0 infrastructure dependencies.

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a list of changes in each version.

## Support

For issues or questions about the tool, check the logs first:

```bash
sudo journalctl -u ipmonitor -f
```

If you find a bug or have a feature request, please [open an issue](https://github.com/giovanniandreuzza/ipmonitor/issues).

## Acknowledgments

- Built with ❤️ for the Raspberry Pi community
- Uses [ipify](https://www.ipify.org/) API for IP detection
- Telegram Bot API for notifications

---

**Author**: Giovanni Andreuzza  
**Repository**: [github.com/giovanniandreuzza/ipmonitor](https://github.com/giovanniandreuzza/ipmonitor)
