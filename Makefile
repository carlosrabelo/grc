# Makefile for compiling and running the Go program

# Variables
APP_NAME := grc
BUILD_DIR := build
SRC_FILE := main.go
RESOURCES_DIR=resources
YAML_FILE=$(RESOURCES_DIR)/example.yaml
XML_FILE=$(RESOURCES_DIR)/example.xml


.PHONY: all build run clean

# Default target: build the program
all: build

# Build the binary
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC_FILE)
	@echo "Build completed. Binary is located in $(BUILD_DIR)/$(APP_NAME)"

# Run the program
run: build
	@echo "Generating XML from $(YAML_FILE)..."
	$(BUILD_DIR)/$(APP_NAME) $(YAML_FILE)

# Clean the build directory
clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean completed."
