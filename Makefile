# Gmail Rules Creator - Build automation for Go project

# Variables
GO		= go
BIN		= grc
SRC		= ./cmd/grc
BUILD_DIR	= ./bin
VERSION		= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
BUILD_TIME	= $(shell date +%Y-%m-%dT%H:%M:%S%z)
LDFLAGS		= -s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)

# Default target - show help
.DEFAULT_GOAL := help

.PHONY: help all build clean run test fmt vet lint deps mod-tidy info

help:	## Show this help
	@echo "Gmail Rules Creator - Build automation"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-12s %s\n", $$1, $$2}'

all: clean fmt vet build	## Clean, format, vet and build

build:	## Build the project
	@echo "Building $(BIN)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 $(GO) build -trimpath -tags netgo -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BIN) $(SRC)

clean:	## Remove build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@$(GO) clean -cache -testcache -modcache 2>/dev/null || true

run:	## Run the binary
	@$(BUILD_DIR)/$(BIN)

test:	## Run Go test suite
	$(GO) test ./...

fmt:	## Format Go sources
	$(GO) fmt ./...

vet:	## Vet Go code
	$(GO) vet ./...

lint:	## Lint with golangci-lint
	@command -v golangci-lint >/dev/null 2>&1 || { \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BUILD_DIR) latest; \
	}
	PATH=$(BUILD_DIR):$$PATH golangci-lint run ./...

deps:	## Download Go module dependencies
	$(GO) mod download

mod-tidy:	## Tidy Go modules
	$(GO) mod tidy

info:	## Show project information
	@echo "Project: Gmail Rules Creator"
	@echo "Binary: $(BIN)"
	@echo "Source: $(SRC)"
	@echo "Build: $(BUILD_DIR)"
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go version: $$($(GO) version)"
