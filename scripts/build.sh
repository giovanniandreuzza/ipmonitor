#!/bin/bash
set -e

echo "Building IP Monitor..."

# Get version from VERSION file
VERSION=$(cat VERSION)

# Build directory
BUILD_DIR="bin"
mkdir -p $BUILD_DIR

# Build for multiple platforms
echo "Building for Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$VERSION" -o $BUILD_DIR/ipmonitor-linux-amd64 ./cmd/ipmonitor

echo "Building for Linux ARM64 (Raspberry Pi)..."
GOOS=linux GOARCH=arm64 go build -ldflags "-X main.Version=$VERSION" -o $BUILD_DIR/ipmonitor-linux-arm64 ./cmd/ipmonitor

echo "Building for Linux ARMv7 (Raspberry Pi 32-bit)..."
GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-X main.Version=$VERSION" -o $BUILD_DIR/ipmonitor-linux-armv7 ./cmd/ipmonitor

echo "Build complete! Binaries available in $BUILD_DIR/"
