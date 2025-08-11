# Conway Polyhedron Notation Library Makefile

# Variables
GO := go
GOFMT := gofmt
GOFUMPT := gofumpt
GCI := gci
GOLINT := golangci-lint
STATICCHECK := staticcheck
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
	@$(GOFUMPT) -w .
	@$(GCI) write --skip-generated -s standard -s default -s "prefix(github.com/sksmith/conway)" .
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

# Run staticcheck
.PHONY: staticcheck
staticcheck:
	@echo "Running staticcheck..."
	@$(STATICCHECK) ./...

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
	@$(GO) install honnef.co/go/tools/cmd/staticcheck@latest
	@$(GO) install mvdan.cc/gofumpt@latest
	@$(GO) install github.com/daixiang0/gci@latest

# Run property-based tests
.PHONY: property-test
property-test:
	@echo "Running property-based tests..."
	@$(GOTEST) -v -run TestProperty ./...

# Run concurrency tests
.PHONY: concurrency-test
concurrency-test:
	@echo "Running concurrency tests with race detector..."
	@$(GOTEST) -v -race -run "TestConcurrent|TestCentroidCachingRace|TestBoundingBoxCalculationRace|TestAtomicIDGeneration" ./...

# Run all pre-commit checks (essential only)
.PHONY: pre-commit
pre-commit: check-fmt vet staticcheck lint-critical test concurrency-test build
	@echo "✅ All pre-commit checks passed!"

# Run all checks (used by CI)
.PHONY: ci
ci: deps fmt vet staticcheck lint test-coverage

# Quick check (fast checks only)
.PHONY: quick-check
quick-check: check-fmt vet staticcheck
	@echo "✅ Quick checks passed!"

# Run critical linting only (for pre-commit)
.PHONY: lint-critical
lint-critical:
	@echo "Running critical linter checks..."
	@$(GOLINT) run --config .golangci-critical.yml ./...

# Check if code is properly formatted
.PHONY: check-fmt
check-fmt:
	@echo "Checking if code is formatted..."
	@test -z "$$($(GOFUMPT) -l .)" || (echo "Code is not formatted. Run 'make fmt'"; exit 1)

# Security scan
.PHONY: security
security:
	@echo "Running security scan..."
	@$(GO) list -json -deps ./... | nancy sleuth

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo ""
	@echo "🚀 Main targets:"
	@echo "  pre-commit   - Run all pre-commit checks (recommended before push)"
	@echo "  quick-check  - Run fast checks only (fmt, vet, staticcheck)"
	@echo "  all          - Run fmt, lint, test, and build"
	@echo "  ci           - Run all checks (CI pipeline)"
	@echo ""
	@echo "🔧 Development:"
	@echo "  fmt          - Format Go code (gofumpt + gci + go mod tidy)"
	@echo "  lint         - Run golangci-lint"
	@echo "  vet          - Run go vet"
	@echo "  staticcheck  - Run staticcheck"
	@echo "  test         - Run tests"
	@echo "  build        - Build examples"
	@echo ""
	@echo "📊 Analysis:"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  concurrency-test- Run concurrency tests with race detector"
	@echo "  property-test- Run property-based tests"
	@echo "  bench        - Run benchmarks"
	@echo "  security     - Run security scan"
	@echo ""
	@echo "🛠️  Setup:"
	@echo "  deps         - Install dependencies"
	@echo "  dev-deps     - Install development tools"
	@echo "  clean        - Clean build artifacts"
	@echo ""
	@echo "✅ Validation:"
	@echo "  check-fmt    - Check if code is formatted"
	@echo ""
	@echo "  help         - Show this help message"