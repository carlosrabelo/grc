# Root Makefile

.DEFAULT_GOAL := help

VERSION ?= dev
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

.PHONY: build clean help install lint run test uninstall

build: ## Build binaries into ./bin using ldflags
	@$(MAKE) -C core build VERSION=$(VERSION) BUILD_TIME=$(BUILD_TIME)

clean: ## Remove artifacts while keeping bin/.gitkeep
	@$(MAKE) -C core clean

install: build ## Install grc via scripts/install.sh
	@./scripts/install.sh grc

lint: ## Execute golangci-lint when available
	@$(MAKE) -C core lint

run: ## Run the main application via go run
	@$(MAKE) -C core run VERSION=$(VERSION) BUILD_TIME=$(BUILD_TIME)

test: ## Run Go tests (go test ./...)
	@$(MAKE) -C core test

uninstall: ## Uninstall grc via scripts/uninstall.sh
	@./scripts/uninstall.sh grc

help: ## Show this structured help
	@echo "Build & Install:"
	@printf " %-15s %s\n" "build" "Build binaries into ./bin using ldflags"
	@printf " %-15s %s\n" "install" "Install grc via scripts/install.sh"
	@printf " %-15s %s\n" "uninstall" "Uninstall grc via scripts/uninstall.sh"
	@echo ""
	@echo "Quality:"
	@printf " %-15s %s\n" "lint" "Execute golangci-lint when available"
	@printf " %-15s %s\n" "test" "Run Go tests (go test ./...)"
	@echo ""
	@echo "Runtime:"
	@printf " %-15s %s\n" "run" "Run the main application via go run"
	@printf " %-15s %s\n" "clean" "Remove artifacts while keeping bin/.gitkeep"
	@printf " %-15s %s\n" "help" "Show this structured help"

%:
	@$(MAKE) -C core $@ VERSION=$(VERSION) BUILD_TIME=$(BUILD_TIME)
