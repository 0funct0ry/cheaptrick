.PHONY: all build clean test lint vet fmt run install help

# Binary name
BINARY_NAME=cheaptrick

# Build flags
LDFLAGS=-ldflags="-s -w"

# Directories
BIN_DIR=bin
FIXTURES_DIR=test_fixtures

all: clean build lint test

## Build:
build: ## Build the cheaptrick binary
	@echo "==> Building ${BINARY_NAME}..."
	@mkdir -p ${BIN_DIR}
	@go build ${LDFLAGS} -o ${BIN_DIR}/${BINARY_NAME} .
	@echo "==> Done. Binary is in ${BIN_DIR}/${BINARY_NAME}"

install: ## Install the binary to GOPATH/bin
	@echo "==> Installing ${BINARY_NAME}..."
	@go install ${LDFLAGS} .
	@echo "==> Done. Installed to $$GOPATH/bin"

## Development:
run: build ## Run the server with default test settings
	@echo "==> Running ${BINARY_NAME}..."
	@./${BIN_DIR}/${BINARY_NAME} --fixtures=${FIXTURES_DIR} --log=mock_log.jsonl

test: ## Run unit tests
	@echo "==> Running tests..."
	@go test -v ./...

lint: ## Run golangci-lint (if installed)
	@echo "==> Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping..."; \
	fi

vet: ## Run go vet
	@echo "==> Running go vet..."
	@go vet ./...

fmt: ## Run go fmt
	@echo "==> Formatting code..."
	@go fmt ./...

clean: ## Clean generated binaries and logs
	@echo "==> Cleaning..."
	@go clean
	@rm -rf ${BIN_DIR}
	@rm -f mock_log.jsonl
	@rm -rf ${FIXTURES_DIR}/*.json

## Help:
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2} /^## / {printf "\n\033[1m%s\033[0m\n", substr($$0, 4)}' $(MAKEFILE_LIST)
