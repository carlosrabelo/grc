# Core Makefile

BIN ?= grc
VERSION ?= dev
BUILD_TIME ?= unknown
LD_FLAGS := -s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)
GOCACHE_DIR ?= $(abspath ../.cache/go-build)
export GOCACHE := $(GOCACHE_DIR)

.DEFAULT_GOAL := help

.PHONY: build clean help lint run test

build: ## Build selected binary (BIN=<name>) into ../bin
	@mkdir -p "$(GOCACHE_DIR)"
	GO111MODULE=on go build -ldflags "$(LD_FLAGS)" -o ../bin/$(BIN) ./cmd/$(BIN)

clean: ## Remove build artifacts while keeping ../bin/.gitkeep
	@if [ -d ../bin ]; then \
		find ../bin -type f ! -name ".gitkeep" -delete; \
	else \
		mkdir -p ../bin; \
	fi

help: ## Show this structured help for core/
	@echo "Build & Install (core):"
	@printf " %-15s %s\n" "build" "Build selected binary (BIN=<name>)"
	@printf " %-15s %s\n" "clean" "Remove artifacts keeping ../bin/.gitkeep"
	@echo ""
	@echo "Quality (core):"
	@printf " %-15s %s\n" "lint" "Run golangci-lint when available"
	@printf " %-15s %s\n" "test" "Run Go tests for the core module"
	@echo ""
	@echo "Runtime (core):"
	@printf " %-15s %s\n" "run" "Run the binary (BIN=<name>) via go run"
	@printf " %-15s %s\n" "help" "Show this structured help"

lint: ## Run golangci-lint when available
	@if command -v golangci-lint >/dev/null 2>&1; then \
		mkdir -p "$(GOCACHE_DIR)"; \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint nao encontrado"; \
	fi

run: ## Run the selected binary (BIN=<name>) via go run
	@mkdir -p "$(GOCACHE_DIR)"
	GO111MODULE=on go run -ldflags "$(LD_FLAGS)" ./cmd/$(BIN)

test: ## Run Go tests for the core module
	@mkdir -p "$(GOCACHE_DIR)"
	GO111MODULE=on go test ./...
