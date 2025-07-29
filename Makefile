.PHONY: build run clean test deps

# Variables
BINARY_NAME=mcp-server
BUILD_DIR=bin
MAIN_PATH=./cmd/mcp-server

# Default target
all: build

# Install dependencies
deps:
	go mod tidy
	go mod download

# Build the server
build: deps
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Build for different platforms
build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)

build-windows:
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)

# Run the server
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run with debug logging
run-debug: build
	./$(BUILD_DIR)/$(BINARY_NAME) -debug

# Run on different port
run-port: build
	./$(BUILD_DIR)/$(BINARY_NAME) -port 9090

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Development server with auto-reload (requires air)
dev:
	air -c .air.toml

# Docker build
docker-build:
	docker build -t play-mcp-server .

# Docker run
docker-run:
	docker run -p 8080:8080 play-mcp-server

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the server binary"
	@echo "  run           - Build and run the server"
	@echo "  run-debug     - Run with debug logging"
	@echo "  run-port      - Run on port 9090"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean         - Clean build artifacts"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  deps          - Install dependencies"
	@echo "  dev           - Run development server with auto-reload"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
