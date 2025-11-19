CORE_DIR := core

.DEFAULT_GOAL := help

VERSION ?= dev
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

.PHONY: build clean help install lint run test uninstall

build: ## Build binaries into ./bin using ldflags
	@$(MAKE) -C $(CORE_DIR) build VERSION=$(VERSION) BUILD_TIME=$(BUILD_TIME)

clean: ## Remove artifacts while keeping bin/.gitkeep
	@$(MAKE) -C $(CORE_DIR) clean

help: ## Show available targets
	@echo "GRC - Gmail Rules Creator"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*## "} {printf "  %-15s %s\n", $$1, $$2}'
	@echo ""
	@echo "For more targets, run 'make -C core help'"

install: build ## Install grc via scripts/install.sh
	@./scripts/install.sh grc

lint: ## Execute golangci-lint when available
	@$(MAKE) -C $(CORE_DIR) lint

run: ## Run the main application via go run
	@$(MAKE) -C $(CORE_DIR) run VERSION=$(VERSION) BUILD_TIME=$(BUILD_TIME)

test: ## Run Go tests (go test ./...)
	@$(MAKE) -C $(CORE_DIR) test

uninstall: ## Uninstall grc via scripts/uninstall.sh
	@./scripts/uninstall.sh grc
