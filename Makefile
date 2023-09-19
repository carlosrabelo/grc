MAKEFLAGS += --no-print-directory

.DEFAULT_GOAL := help

.PHONY: build build-static clean fmt help install lint quality test vet

BINARY_NAME := grc
VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS    := -s -w -X github.com/carlosrabelo/grc/grc/internal/app.Version=$(VERSION) -X github.com/carlosrabelo/grc/grc/internal/app.BuildTime=$(BUILD_TIME)

help: ## Show available targets
	@echo "grc - Available targets"
	@echo ""
	@grep -hE '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*## "} {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build binary
	@./.make/build.sh

build-static: ## Build fully static binary (CGO_ENABLED=0, netgo)
	@./.make/build-static.sh

install: ## Build and install to /usr/local/bin (root) or ~/.local/bin (user)
	@./.make/install.sh

test: ## Run all tests
	@./.make/test.sh

lint: ## Run golangci-lint
	@golangci-lint run ./...

quality: ## Run all quality checks (fmt, vet, lint, test)
	@./.make/quality.sh

fmt: ## Format Go source code
	@go fmt ./...

vet: ## Run go vet
	@go vet ./...

clean: ## Remove build artifacts
	@./.make/clean.sh
