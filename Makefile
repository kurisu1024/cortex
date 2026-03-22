.PHONY: build install clean test run help

# Binary name
BINARY_NAME=cortex

# Build directory
BUILD_DIR=./bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	@echo "Building Cortex..."
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/cortex
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

install: build ## Install binary to GOPATH/bin
	@echo "Installing Cortex..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
	@echo "✅ Installed to $(GOPATH)/bin/$(BINARY_NAME)"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

test: ## Run tests
	@echo "Running tests..."
	@$(GOTEST) -v ./...

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@$(GOMOD) download
	@$(GOMOD) tidy
	@echo "✅ Dependencies updated"

run: build ## Build and run the application
	@$(BUILD_DIR)/$(BINARY_NAME)

dev: ## Run in development mode (rebuild on changes)
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	@air

# Development helpers
fmt: ## Format code
	@echo "Formatting code..."
	@$(GOCMD) fmt ./...
	@echo "✅ Code formatted"

lint: ## Run linter
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@golangci-lint run ./...

# Quick build and test
quick: fmt build test ## Format, build, and test

.DEFAULT_GOAL := help
