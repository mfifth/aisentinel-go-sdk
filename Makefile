.PHONY: all build test clean lint fmt vet mod-tidy install-tools help

# Default target
all: mod-tidy fmt vet lint test build

# Build the project
build:
	go build -v ./...

# Run tests
test:
	go test -v -race -cover ./...

# Run tests with coverage
test-cover:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	go clean ./...
	rm -f coverage.out coverage.html

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...
	gofmt -s -w .

# Vet code
vet:
	go vet ./...

# Tidy modules
mod-tidy:
	go mod tidy
	go mod verify

# Download dependencies
mod-download:
	go mod download

# Install development tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

# Run benchmarks
bench:
	go test -bench=. -benchmem ./...

# Run fuzz tests (if any)
fuzz:
	go test -fuzz=. -fuzztime=10s ./...

# Generate mocks (if using mockery)
mocks:
	@echo "Generating mocks..."
	@which mockery > /dev/null || (echo "Install mockery: go install github.com/vektra/mockery/v2@latest" && exit 1)
	mockery --all --output ./mocks

# Run security scan
security:
	gosec ./...

# Help
help:
	@echo "Available targets:"
	@echo "  all          - Run mod-tidy, fmt, vet, lint, test, and build"
	@echo "  build        - Build the project"
	@echo "  test         - Run tests"
	@echo "  test-cover   - Run tests with coverage report"
	@echo "  clean        - Clean build artifacts"
	@echo "  lint         - Run linter"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code"
	@echo "  mod-tidy     - Tidy and verify modules"
	@echo "  mod-download - Download dependencies"
	@echo "  install-tools- Install development tools"
	@echo "  bench        - Run benchmarks"
	@echo "  fuzz         - Run fuzz tests"
	@echo "  mocks        - Generate mocks"
	@echo "  security     - Run security scan"
	@echo "  help         - Show this help"
