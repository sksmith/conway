# Conway Polyhedron Notation Library Makefile

# Variables
GO := go
GOFMT := gofmt
GOLINT := golangci-lint
GOTEST := $(GO) test
GOBUILD := $(GO) build
GOMOD := $(GO) mod
GOVET := $(GO) vet

# Directories
SRC_DIR := ./conway
EXAMPLES_DIR := ./examples
COVERAGE_DIR := ./coverage

# Default target
.PHONY: all
all: fmt lint test build

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting Go code..."
	@$(GOFMT) -s -w .
	@$(GO) mod tidy

# Lint code
.PHONY: lint
lint:
	@echo "Running linter..."
	@$(GOLINT) run ./...

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@$(GOTEST) -v -race ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	@$(GOTEST) -v -race -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	@$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated at $(COVERAGE_DIR)/coverage.html"

# Run benchmarks
.PHONY: bench
bench:
	@echo "Running benchmarks..."
	@$(GOTEST) -bench=. -benchmem ./...

# Vet code
.PHONY: vet
vet:
	@echo "Running go vet..."
	@$(GOVET) ./...

# Build examples
.PHONY: build
build:
	@echo "Building examples..."
	@cd $(EXAMPLES_DIR)/basic && $(GOBUILD) -o basic .
	@cd $(EXAMPLES_DIR)/advanced && $(GOBUILD) -o advanced .

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(COVERAGE_DIR)
	@find . -name "*.test" -delete
	@find . -name "*.out" -delete
	@cd $(EXAMPLES_DIR)/basic && rm -f basic basic.exe
	@cd $(EXAMPLES_DIR)/advanced && rm -f advanced advanced.exe

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@$(GOMOD) download
	@$(GOMOD) tidy

# Install development tools
.PHONY: dev-deps
dev-deps:
	@echo "Installing development dependencies..."
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run property-based tests
.PHONY: property-test
property-test:
	@echo "Running property-based tests..."
	@$(GOTEST) -v -run TestProperty ./...

# Run concurrency tests
.PHONY: concurrency-test
concurrency-test:
	@echo "Running concurrency tests..."
	@$(GOTEST) -v -race -run TestConcurrency ./...

# Run all checks (used by CI)
.PHONY: ci
ci: deps fmt vet lint test-coverage

# Check if code is properly formatted
.PHONY: check-fmt
check-fmt:
	@echo "Checking if code is formatted..."
	@test -z "$$($(GOFMT) -l .)" || (echo "Code is not formatted. Run 'make fmt'"; exit 1)

# Security scan
.PHONY: security
security:
	@echo "Running security scan..."
	@$(GO) list -json -deps ./... | nancy sleuth

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Run fmt, lint, test, and build"
	@echo "  fmt          - Format Go code"
	@echo "  lint         - Run linter"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  bench        - Run benchmarks"
	@echo "  vet          - Run go vet"
	@echo "  build        - Build examples"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Install dependencies"
	@echo "  dev-deps     - Install development dependencies"
	@echo "  ci           - Run all checks (CI pipeline)"
	@echo "  check-fmt    - Check if code is formatted"
	@echo "  security     - Run security scan"
	@echo "  help         - Show this help message"