# Backend Makefile for Image RAG Service
# Includes test, lint, and import ordering standards

.PHONY: help build test lint fmt vet tidy clean run dev deps check-imports fix-imports check-all

# Default target
help:
	@echo "Backend Makefile for Image RAG Service"
	@echo ""
	@echo "Available targets:"
	@echo "  build        - Build the backend binary"
	@echo "  test         - Run all tests with race detection and coverage"
	@echo "  test-short   - Run tests without race detection (faster)"
	@echo "  lint         - Run golangci-lint with full configuration"
	@echo "  lint-fast    - Run golangci-lint with fast checks only"
	@echo "  fmt          - Format code with go fmt"
	@echo "  vet          - Run go vet for suspicious constructs"
	@echo "  tidy         - Clean up go.mod and go.sum"
	@echo "  clean        - Clean build artifacts"
	@echo "  run          - Run the backend server"
	@echo "  dev          - Run with hot reload (using air)"
	@echo "  deps         - Install/update dependencies"
	@echo "  check-imports - Check import ordering standards"
	@echo "  fix-imports  - Fix import ordering automatically"
	@echo "  check-all    - Run all checks (fmt, vet, lint, test)"
	@echo ""
	@echo "Import Ordering Standards:"
	@echo "  1. Standard library packages"
	@echo "  2. Third-party packages"
	@echo "  3. Project packages (internal/)"
	@echo ""
	@echo "Each group separated by a blank line"

# Build configuration
BINARY_NAME=image-rag-backend
MAIN_PATH=./cmd/server/main.go
BUILD_DIR=./bin

# Build the backend binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Run all tests with comprehensive coverage
test:
	@echo "Running tests with race detection and coverage..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests without race detection for faster execution
test-short:
	@echo "Running tests (fast mode)..."
	go test -v ./...

# Run golangci-lint with full configuration
lint:
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
	else \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run --timeout=5m; \
	fi

# Run golangci-lint with fast checks only
lint-fast:
	@echo "Running golangci-lint (fast)..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --fast --timeout=2m; \
	else \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run --fast --timeout=2m; \
	fi

# Format code with go fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run go vet for suspicious constructs
vet:
	@echo "Running go vet..."
	go vet ./...

# Clean up go.mod and go.sum
tidy:
	@echo "Tidying go.mod..."
	go mod tidy
	go mod verify

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)/$(BINARY_NAME)
	rm -rf coverage.out coverage.html

# Run the backend server
run:
	@echo "Starting backend server..."
	go run $(MAIN_PATH)

# Run with hot reload using air (if installed)
dev:
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "air not found. Installing..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

# Install/update dependencies
deps:
	@echo "Installing/updating dependencies..."
	go mod download
	go mod verify

# Check import ordering standards
check-imports:
	@echo "Checking import ordering..."
	@failed=0; \
	for file in $$(find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*"); do \
		if ! goimports -l -local image-rag-backend $$file | grep -q "^"; then \
			echo "✓ $$file"; \
		else \
			echo "✗ $$file - import ordering issues found"; \
			failed=1; \
		fi; \
	done; \
	if [ $$failed -eq 1 ]; then \
		echo ""; \
		echo "Import ordering issues found. Run 'make fix-imports' to fix automatically."; \
		exit 1; \
	else \
		echo "All imports are properly ordered!"; \
	fi

# Fix import ordering automatically
fix-imports:
	@echo "Fixing import ordering..."
	@if command -v goimports >/dev/null 2>&1; then \
		find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" -exec goimports -w -local image-rag-backend {} \;; \
		echo "Import ordering fixed!"; \
	else \
		echo "goimports not found. Installing..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
		find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" -exec goimports -w -local image-rag-backend {} \;; \
		echo "Import ordering fixed!"; \
	fi

# Run all checks
check-all: fmt vet lint test
	@echo "All checks passed! ✓"

# Pre-commit hook
pre-commit: check-imports fmt vet lint test
	@echo "Pre-commit checks passed! ✓"

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "Development tools installed!"

# Show project info
info:
	@echo "Backend Info:"
	@echo "  Go version: $$(go version)"
	@echo "  Module: $$(go list -m)"
	@echo "  Binary: $(BINARY_NAME)"
	@echo "  Build dir: $(BUILD_DIR)"
	@echo ""
	@echo "Dependencies:"
	@go list -m all | head -20

# Create .golangci.yml config file for consistent linting
create-lint-config:
	@if [ ! -f .golangci.yml ]; then \
		echo "Creating .golangci.yml configuration..."; \
		cat > .golangci.yml << 'EOF'; \
	linters-settings:
	  goimports:
	    local-prefixes: image-rag-backend
	  gofmt:
	    simplify: true
	  govet:
	    check-shadowing: true
	  misspell:
	    locale: US
	  lll:
	    line-length: 120
	  gci:
	    sections:
	      - standard
	      - default
	      - prefix(image-rag-backend)

	linters:
	  enable:
	    - bodyclose
	    - dogsled
	    - dupl
	    - errcheck
	    - exhaustive
	    - funlen
	    - gochecknoinits
	    - goconst
	    - gocritic
	    - gocyclo
	    - gofmt
	    - goimports
	    - gomnd
	    - goprintffuncname
	    - gosec
	    - gosimple
	    - govet
	    - ineffassign
	    - lll
	    - misspell
	    - nakedret
	    - noctx
	    - nolintlint
	    - rowserrcheck
	    - staticcheck
	    - structcheck
	    - stylecheck
	    - typecheck
	    - unconvert
	    - unparam
	    - unused
	    - varcheck
	    - whitespace

	issues:
	  exclude-use-default: false
	  exclude-rules:
	    - path: _test\.go
	      linters:
	        - gomnd
	        - funlen
	    - path: internal/models/
	      linters:
	        - lll
	
	EOF
		echo ".golangci.yml created!"; \
	else \
		echo ".golangci.yml already exists"; \
	fi

# Initialize project for development
setup: create-lint-config install-tools
	@echo "Setting up development environment..."
	go mod tidy
	@echo "Setup complete! Run 'make check-all' to verify everything works."