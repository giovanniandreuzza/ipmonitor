.PHONY: build build-linux build-rpi test clean run install help

# Variables
BINARY_NAME=ipmonitor
VERSION=$(shell cat VERSION)
BUILD_DIR=bin
MAIN_PATH=./cmd/ipmonitor

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build for current platform
	@echo "Building $(BINARY_NAME) v$(VERSION) for current platform..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

build-linux: ## Build for Linux AMD64
	@echo "Building $(BINARY_NAME) v$(VERSION) for Linux AMD64..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64"

build-rpi: ## Build for Raspberry Pi (ARM64)
	@echo "Building $(BINARY_NAME) v$(VERSION) for Raspberry Pi (ARM64)..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64"

build-all: build-linux build-rpi ## Build for all platforms

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)
	@rm -f cmd/ipmonitor/$(BINARY_NAME)
	@echo "Clean complete"

run: ## Run the application locally
	@echo "Running $(BINARY_NAME)..."
	@go run $(MAIN_PATH)

install: build ## Install the binary to /usr/local/bin
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete"

fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

lint: ## Run golangci-lint (requires golangci-lint installed)
	@echo "Running golangci-lint..."
	@golangci-lint run

release: clean test build-all ## Build release binaries
	@echo "Creating release v$(VERSION)..."
	@mkdir -p dist
	@cp $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 dist/
	@cp $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 dist/
	@cd dist && sha256sum * > checksums.txt
	@echo "Release v$(VERSION) ready in dist/"
