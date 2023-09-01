MAKEFLAGS += --no-print-directory

.DEFAULT_GOAL := help

.PHONY: build clean fmt help quality test vet

BINARY_NAME := grc
VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS    := -s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)

help: ## Show available targets
	@echo "grc - Available targets"
	@echo ""
	@grep -hE '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*## "} {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build binary
	@./.make/build.sh

test: ## Run all tests
	@./.make/test.sh

quality: ## Run all quality checks (fmt, vet, test)
	@./.make/quality.sh

fmt: ## Format Go source code
	@go fmt ./...

vet: ## Run go vet
	@go vet ./...

clean: ## Remove build artifacts
	@./.make/clean.sh
