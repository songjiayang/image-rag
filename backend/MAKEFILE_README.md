# Backend Makefile Usage Guide

This Makefile provides comprehensive tooling for the Image RAG Service backend development.

## Quick Start

```bash
# Set up development environment
make setup

# Run all checks (format, vet, lint, test)
make check-all

# Run tests
make test

# Run linting
make lint

# Fix import ordering automatically
make fix-imports
```

## Available Commands

### Development
- `make run` - Start the backend server
- `make dev` - Start with hot reload using air
- `make build` - Build production binary

### Code Quality
- `make fmt` - Format code with go fmt
- `make vet` - Run go vet for suspicious constructs
- `make lint` - Run comprehensive golangci-lint
- `make lint-fast` - Run fast linting checks only
- `make check-imports` - Check import ordering standards
- `make fix-imports` - Fix import ordering automatically

### Testing
- `make test` - Run tests with race detection and coverage
- `make test-short` - Run tests without race detection (faster)

### Utilities
- `make tidy` - Clean up go.mod and go.sum
- `make clean` - Clean build artifacts
- `make deps` - Install/update dependencies
- `make info` - Show project information
- `make install-tools` - Install development tools

## Import Ordering Standards

The project follows strict import ordering:

1. **Standard library packages** (e.g., `fmt`, `net/http`)
2. **Third-party packages** (e.g., `github.com/gin-gonic/gin`)
3. **Project packages** (e.g., `image-rag-backend/internal/logger`)

Each group is separated by a blank line.

### Example:

```go
import (
    "fmt"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    "image-rag-backend/internal/logger"
    "image-rag-backend/internal/models"
)
```

## Configuration

### Environment Variables
- `GIN_MODE` - Set to "release" for production mode
- All other configuration is loaded from `.env` file

### Linting Configuration
- `.golangci.yml` - Comprehensive linting rules
- Local import prefix: `image-rag-backend`

## Development Workflow

1. **Before committing:**
   ```bash
   make pre-commit
   ```

2. **Daily development:**
   ```bash
   make dev  # Start hot reload server
   ```

3. **Before pushing:**
   ```bash
   make check-all
   ```

## Troubleshooting

### Common Issues

1. **Import ordering issues:**
   ```bash
   make fix-imports
   ```

2. **Missing development tools:**
   ```bash
   make install-tools
   ```

3. **Linting errors:**
   ```bash
   # Check specific issues
   make lint
   
   # Fix automatically where possible
   make fix-imports
   make fmt
   ```

### Dependencies

Required tools (automatically installed with `make install-tools`):
- golangci-lint
- goimports
- air (for hot reload)

All tools are installed in your `$GOPATH/bin` directory.