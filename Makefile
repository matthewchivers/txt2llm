# txt2llm Makefile

BINARY_NAME=txt2llm
BIN_DIR=build/bin
COVERAGE_DIR=build/coverage

.PHONY: build build-all install clean help test test-verbose test-coverage test-short lint lint-fix format

# Run tests
test: lint
	go test ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Run tests with coverage
test-coverage:
	mkdir -p $(COVERAGE_DIR)
	go test -v -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated: $(COVERAGE_DIR)/coverage.html"

# Run short tests (skip integration tests)
test-short:
	go test -short ./...

# Run linting
lint: format
	$(shell go env GOPATH)/bin/golangci-lint run

# Run linting with auto-fix
lint-fix: format
	$(shell go env GOPATH)/bin/golangci-lint run --fix

# Format code
format:
	gofmt -w .

# Build the binary for current platform
build: test
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY_NAME) .

# Build binaries for all platforms
build-all: test clean
	mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o $(BIN_DIR)/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME)-windows-amd64.exe .

# Install using go install (recommended)
install: test
	go install .

# Clean build artifacts
clean: 
	go clean
	rm -rf build/

# Show help
help:
	@echo "Available targets:"
	@echo "  test         - Run all tests"
	@echo "  test-verbose - Run tests with verbose output"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  test-short   - Run short tests (skip integration tests)"
	@echo "  lint         - Run linting"
	@echo "  lint-fix     - Run linting with auto-fix"
	@echo "  format       - Format code with gofmt"
	@echo "  build        - Build the binary for current platform"
	@echo "  build-all    - Build binaries for all platforms"
	@echo "  install      - Install using 'go install' (to GOPATH/bin)"
	@echo "  clean        - Clean build artifacts"
	@echo "  help         - Show this help"