# Helper targets for building and using the Gmail Rules Creator

APP_NAME := grc
BIN_DIR := build
BIN := $(BIN_DIR)/$(APP_NAME)
GO ?= go
YAML_SAMPLE := resources/example.yaml
XML_SAMPLE := $(YAML_SAMPLE:.yaml=.xml)
UNAME_S := $(shell uname -s 2>/dev/null)
USER_ID := $(shell id -u 2>/dev/null)

# Use local caches so builds also work in sandboxed environments
GOCACHE ?= $(CURDIR)/.cache
GOMODCACHE ?= $(CURDIR)/.modcache
export GOCACHE
export GOMODCACHE

.PHONY: help all build run run-sample generate-sample fmt tidy clean install

help:
	@echo "Available targets:"
	@echo "  make build            Build the binary into $(BIN)"
	@echo "  make run              Run the compiled binary"
	@echo "  make run-sample       Generate XML using $(YAML_SAMPLE)"
	@echo "  make generate-sample  Build and write $(XML_SAMPLE)"
	@echo "  make fmt              Run go fmt"
	@echo "  make tidy             Run go mod tidy"
	@echo "  make clean            Remove build artifacts and caches"
	@echo "  make install          Install the binary (Linux only)"

all: build

build: $(BIN)

$(BIN): $(shell find . -name '*.go' -print)
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN) ./...

run: build
	$(BIN)

run-sample: build
	$(BIN) $(YAML_SAMPLE)

# Generate XML for the bundled sample file to simplify manual testing
generate-sample: build
	$(BIN) -output $(XML_SAMPLE) $(YAML_SAMPLE)

fmt:
	$(GO) fmt ./...

tidy:
	$(GO) mod tidy

clean:
	rm -rf $(BIN_DIR) $(XML_SAMPLE) .cache .modcache

install: build
ifeq ($(OS),Windows_NT)
	@echo "Install skipped: not required on Windows."
else ifeq ($(UNAME_S),Linux)
	@if [ "$(USER_ID)" = "0" ]; then \
		prefix=/usr/local/bin; \
	else \
		prefix="$$HOME/.local/bin"; \
	fi; \
	install -d "$$prefix"; \
	install -m 0755 "$(BIN)" "$$prefix/$(APP_NAME)"; \
	echo "Installed $(APP_NAME) to $$prefix";
else
	@echo "Install skipped: only Linux installation is supported."
endif
