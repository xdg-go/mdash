.PHONY: all build test lint format install clean run docker-build docker-run help

# Variables
BINARY_NAME=mdash
DOCKER_IMAGE=mdash
GO_FILES=$(shell find . -name '*.go' -type f)
GO_PACKAGES=$(shell go list ./...)

# Default target
all: lint test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) .

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
test-coverage: test
	@echo "Generating coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run .

# Format code
fmt:
	@echo "Formatting code..."
	@gofumpt -w $(GO_FILES)

# Install the binary
install:
	@echo "Installing $(BINARY_NAME)..."
	@go install .

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME) coverage.out coverage.html
	@go clean

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Help target
help:
	@echo "Available targets:"
	@echo "  make build         - Build the binary"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make lint          - Run golangci-lint"
	@echo "  make format        - Format code with gofmt"
	@echo "  make install       - Install the binary using go install"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make deps          - Download and tidy dependencies"
	@echo "  make help          - Show this help message"
