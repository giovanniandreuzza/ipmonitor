#!/bin/bash
set -e

# Install script for Raspberry Pi

BINARY_NAME="ipmonitor"
INSTALL_DIR="/usr/local/bin"
SERVICE_DIR="/etc/systemd/system"

# Detect architecture
ARCH=$(uname -m)
case $ARCH in
    aarch64)
        BINARY_FILE="ipmonitor-linux-arm64"
        ;;
    armv7l)
        BINARY_FILE="ipmonitor-linux-armv7"
        ;;
    x86_64)
        BINARY_FILE="ipmonitor-linux-amd64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

echo "Installing IP Monitor for $ARCH..."

# Check if binary exists
if [ ! -f "bin/$BINARY_FILE" ]; then
    echo "Error: Binary not found. Please run 'make build-all' first."
    exit 1
fi

# Install binary
echo "Installing binary to $INSTALL_DIR..."
sudo cp bin/$BINARY_FILE $INSTALL_DIR/$BINARY_NAME
sudo chmod +x $INSTALL_DIR/$BINARY_NAME

# Install systemd files
echo "Installing systemd service files..."
sudo cp deployments/systemd/ipmonitor.service $SERVICE_DIR/
sudo cp deployments/systemd/ipmonitor.timer $SERVICE_DIR/

# Prompt for credentials
echo ""
echo "Please enter your Telegram credentials:"
read -p "Telegram Bot Token: " BOT_TOKEN
read -p "Telegram Chat ID: " CHAT_ID

# Update service file with credentials
sudo sed -i "s/YOUR_BOT_TOKEN_HERE/$BOT_TOKEN/g" $SERVICE_DIR/ipmonitor.service
sudo sed -i "s/YOUR_CHAT_ID_HERE/$CHAT_ID/g" $SERVICE_DIR/ipmonitor.service

# Reload systemd and enable timer
echo "Enabling and starting systemd timer..."
sudo systemctl daemon-reload
sudo systemctl enable ipmonitor.timer
sudo systemctl start ipmonitor.timer

echo ""
echo "✅ Installation complete!"
echo ""
echo "Check status with: sudo systemctl status ipmonitor.timer"
echo "View logs with: sudo journalctl -u ipmonitor -f"
