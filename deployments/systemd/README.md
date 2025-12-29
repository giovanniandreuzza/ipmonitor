# Systemd Deployment

This directory contains systemd service and timer files for running IP Monitor on Linux systems.

## Files

- `ipmonitor.service` - Systemd service unit (runs the monitor once)
- `ipmonitor.timer` - Systemd timer unit (runs every 5 minutes)
- `ipmonitor.env.example` - Environment variables template

## Installation

### 1. Install the binary

Download the latest release or build from source:

```bash
# Download for Raspberry Pi (64-bit)
wget https://github.com/giovanniandreuzza/ipmonitor/releases/latest/download/ipmonitor-linux-arm64
chmod +x ipmonitor-linux-arm64
sudo mv ipmonitor-linux-arm64 /usr/local/bin/ipmonitor

# Or for Raspberry Pi (32-bit ARMv7)
wget https://github.com/giovanniandreuzza/ipmonitor/releases/latest/download/ipmonitor-linux-armv7
chmod +x ipmonitor-linux-armv7
sudo mv ipmonitor-linux-armv7 /usr/local/bin/ipmonitor

# Or for Linux AMD64
wget https://github.com/giovanniandreuzza/ipmonitor/releases/latest/download/ipmonitor-linux-amd64
chmod +x ipmonitor-linux-amd64
sudo mv ipmonitor-linux-amd64 /usr/local/bin/ipmonitor
```

### 2. Configure environment variables

```bash
# Copy the example environment file
sudo cp ipmonitor.env.example /etc/default/ipmonitor

# Edit with your credentials
sudo nano /etc/default/ipmonitor
```

Update these values:
- `TELEGRAM_BOT_TOKEN` - Get from @BotFather on Telegram
- `TELEGRAM_CHAT_ID` - Get from @userinfobot on Telegram
- `LOG_LEVEL` - Optional (DEBUG, INFO, WARN, ERROR)

**Important:** Secure the file (contains secrets):
```bash
sudo chmod 600 /etc/default/ipmonitor
sudo chown root:root /etc/default/ipmonitor
```

### 3. Install systemd units

```bash
# Copy service and timer files
sudo cp ipmonitor.service /etc/systemd/system/
sudo cp ipmonitor.timer /etc/systemd/system/

# Set correct permissions
sudo chmod 644 /etc/systemd/system/ipmonitor.service
sudo chmod 644 /etc/systemd/system/ipmonitor.timer

# Reload systemd
sudo systemctl daemon-reload
```

### 4. Enable and start

```bash
# Enable the timer (starts on boot)
sudo systemctl enable ipmonitor.timer

# Start the timer immediately
sudo systemctl start ipmonitor.timer

# Optionally, run once now to test
sudo systemctl start ipmonitor.service
```

## Verification

### Check timer status

```bash
# View timer status
sudo systemctl status ipmonitor.timer

# List all timers
sudo systemctl list-timers ipmonitor.timer
```

### Check service logs

```bash
# View recent logs
sudo journalctl -u ipmonitor -n 50

# Follow logs in real-time
sudo journalctl -u ipmonitor -f

# View logs since boot
sudo journalctl -u ipmonitor -b
```

### Test the service manually

```bash
# Run the service once
sudo systemctl start ipmonitor.service

# Check if it succeeded
sudo systemctl status ipmonitor.service
```

## Customization

### Change execution interval

Edit `ipmonitor.timer` and modify the `OnUnitActiveSec` value:

```ini
[Timer]
OnBootSec=1min
OnUnitActiveSec=5min  # Change this (e.g., 10min, 1h, 30s)
Persistent=true
```

Then reload:
```bash
sudo systemctl daemon-reload
sudo systemctl restart ipmonitor.timer
```

### Change user

By default, the service runs as user `pi`. To change:

Edit `ipmonitor.service`:
```ini
[Service]
User=youruser  # Change this
```

Then:
```bash
sudo systemctl daemon-reload
sudo systemctl restart ipmonitor.timer
```

### Update environment variables

```bash
# Edit the environment file
sudo nano /etc/default/ipmonitor

# No need to reload systemd, changes take effect on next run
```

## Troubleshooting

### Service fails to start

```bash
# Check detailed logs
sudo journalctl -u ipmonitor -n 100 --no-pager

# Verify binary exists and is executable
ls -l /usr/local/bin/ipmonitor

# Verify environment file exists
sudo cat /etc/default/ipmonitor

# Test manually
sudo -u pi /usr/local/bin/ipmonitor
```

### Timer not running

```bash
# Check if timer is enabled
sudo systemctl is-enabled ipmonitor.timer

# Check if timer is active
sudo systemctl is-active ipmonitor.timer

# View timer details
sudo systemctl status ipmonitor.timer
```

### Environment variables not loaded

```bash
# Verify file path is correct in service file
grep EnvironmentFile /etc/systemd/system/ipmonitor.service

# Should show: EnvironmentFile=/etc/default/ipmonitor

# Check file permissions
ls -l /etc/default/ipmonitor
```

## Uninstallation

```bash
# Stop and disable the timer
sudo systemctl stop ipmonitor.timer
sudo systemctl disable ipmonitor.timer

# Remove systemd files
sudo rm /etc/systemd/system/ipmonitor.service
sudo rm /etc/systemd/system/ipmonitor.timer
sudo rm /etc/default/ipmonitor

# Remove binary
sudo rm /usr/local/bin/ipmonitor

# Reload systemd
sudo systemctl daemon-reload
```

## Alternative: Docker Deployment

For containerized deployment, see the Docker files in the project root:
- `Dockerfile`
- `docker-compose.yml`

Docker is recommended for easier updates and isolation.
