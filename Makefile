MAKEFLAGS += --no-print-directory

.DEFAULT_GOAL := help

.PHONY: build build-static clean deb fmt help install lint package quality test uninstall vet

BINARY_NAME := grc
VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS    := -s -w -X github.com/carlosrabelo/grc/grc/internal/app.Version=$(VERSION) -X github.com/carlosrabelo/grc/grc/internal/app.BuildTime=$(BUILD_TIME)

export BINARY_NAME

help: ## Show available targets
	@echo "grc - Available targets"
	@echo ""
	@grep -hE '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*## "} {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the binary
	@./.make/build.sh

build-static: ## Build fully static binary (CGO_ENABLED=0, netgo)
	@./.make/build-static.sh

test: ## Run tests
	@./.make/test.sh

fmt: ## Format code
	@go fmt ./...

vet: ## Run go vet
	@go vet ./...

lint: vet ## Run linter (go vet)

quality: ## Run all quality checks
	@./.make/quality.sh

clean: ## Remove build artifacts
	@./.make/clean.sh

install: ## Install to ~/.local/bin (requires prior make build)
	@./.make/install.sh

uninstall: ## Remove from ~/.local/bin (sudo only if needed)
	@./.make/uninstall.sh

package: ## Build a .deb into dist/ (requires prior make build)
	@./.make/package.sh

deb: package ## Alias for package
