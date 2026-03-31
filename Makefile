.PHONY: help build run test clean install-deps lint fmt

# Variables
BINARY_NAME=mcp-server
BINARY_PATH=./bin/$(BINARY_NAME)
GO=go
MAIN_PKG=./cmd/server

# Default target
help:
	@echo "MCP Go Server - Build Commands"
	@echo "=============================="
	@echo ""
	@echo "Available targets:"
	@echo "  make build          - Build the server binary"
	@echo "  make run            - Build and run the server"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make install-deps   - Download Go dependencies"
	@echo "  make lint           - Run linter"
	@echo "  make fmt            - Format code"
	@echo "  make help           - Show this help message"
	@echo ""

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	@$(GO) build -o $(BINARY_PATH) $(MAIN_PKG)
	@echo "✓ Built: $(BINARY_PATH)"

# Run the server
run: build
	@echo "Starting server..."
	@$(BINARY_PATH)

# Run with custom config
run-config:
	@echo "Starting server with config..."
	@$(BINARY_PATH) -config config/default.yaml

# Run tests
test:
	@echo "Running tests..."
	@$(GO) test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@$(GO) clean
	@echo "✓ Cleaned"

# Install dependencies
install-deps:
	@echo "Downloading dependencies..."
	@$(GO) mod download
	@$(GO) mod tidy
	@echo "✓ Dependencies installed"

# Run linter (if golangci-lint is installed)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint &> /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, skipping"; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	@$(GO) fmt ./...
	@goimports -w . 2>/dev/null || $(GO) fmt ./...
	@echo "✓ Code formatted"

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	@$(GO) test -v -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

# Docker targets
docker-build:
	@echo "Building Docker image..."
	@docker build -t mcp-go:latest .
	@echo "✓ Docker image built"

docker-run:
	@echo "Running Docker container..."
	@docker run -p 9090:9090 mcp-go:latest

# Example client targets
run-py-example:
	@echo "Running Python example..."
	@python3 examples/python_client.py

run-js-example:
	@echo "Running JavaScript example..."
	@node examples/javascript_client.js

# Install dev tools
install-tools:
	@echo "Installing development tools..."
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@$(GO) install golang.org/x/tools/cmd/goimports@latest
	@echo "✓ Tools installed"
