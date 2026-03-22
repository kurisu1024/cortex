.PHONY: build install clean test run help check-prereqs install-prereqs

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

# Detect OS
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	OS = macos
	PACKAGE_MANAGER = brew
endif
ifeq ($(UNAME_S),Linux)
	OS = linux
	# Detect package manager
	ifeq ($(shell command -v apt-get 2> /dev/null),)
		ifeq ($(shell command -v yum 2> /dev/null),)
			PACKAGE_MANAGER = unknown
		else
			PACKAGE_MANAGER = yum
		endif
	else
		PACKAGE_MANAGER = apt
	endif
endif

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

check-prereqs: ## Check if prerequisites are installed
	@echo "🔍 Checking prerequisites..."
	@echo ""
	@echo "Checking FFmpeg..."
	@command -v ffmpeg >/dev/null 2>&1 && echo "  ✅ FFmpeg is installed" || echo "  ❌ FFmpeg is NOT installed"
	@echo ""
	@echo "Checking Ollama..."
	@command -v ollama >/dev/null 2>&1 && echo "  ✅ Ollama is installed" || echo "  ❌ Ollama is NOT installed"
	@echo ""
	@echo "Checking Piper TTS..."
	@command -v piper >/dev/null 2>&1 && echo "  ✅ Piper is installed" || echo "  ❌ Piper is NOT installed"
	@echo ""

install-prereqs: ## Install missing prerequisites
	@echo "🚀 Installing prerequisites for $(OS)..."
	@echo ""
ifeq ($(OS),macos)
	@echo "📦 Installing with Homebrew..."
	@command -v brew >/dev/null 2>&1 || (echo "❌ Homebrew not found. Install from https://brew.sh" && exit 1)
	@command -v ffmpeg >/dev/null 2>&1 || (echo "Installing FFmpeg..." && brew install ffmpeg)
	@command -v ollama >/dev/null 2>&1 || (echo "Installing Ollama..." && curl -fsSL https://ollama.com/install.sh | sh)
	@if ! command -v piper >/dev/null 2>&1; then \
		echo "⚠️  Piper TTS is not available via Homebrew."; \
		echo "📥 Installing Piper TTS manually..."; \
		mkdir -p $(HOME)/.local/bin; \
		curl -L https://github.com/rhasspy/piper/releases/download/2023.11.14-2/piper_macos_x64.tar.gz -o /tmp/piper.tar.gz; \
		tar -xzf /tmp/piper.tar.gz -C /tmp; \
		cp /tmp/piper/piper $(HOME)/.local/bin/; \
		chmod +x $(HOME)/.local/bin/piper; \
		rm -rf /tmp/piper /tmp/piper.tar.gz; \
		echo "✅ Piper installed to $(HOME)/.local/bin/piper"; \
		echo "⚠️  Make sure $(HOME)/.local/bin is in your PATH"; \
		echo '   Add to ~/.zshrc or ~/.bashrc: export PATH="$$HOME/.local/bin:$$PATH"'; \
	fi
else ifeq ($(OS),linux)
ifeq ($(PACKAGE_MANAGER),apt)
	@echo "📦 Installing with apt..."
	@command -v ffmpeg >/dev/null 2>&1 || (echo "Installing FFmpeg..." && sudo apt-get update && sudo apt-get install -y ffmpeg)
	@command -v ollama >/dev/null 2>&1 || (echo "Installing Ollama..." && curl -fsSL https://ollama.com/install.sh | sh)
	@command -v piper >/dev/null 2>&1 || (echo "⚠️  Piper must be installed manually from https://github.com/rhasspy/piper/releases")
else ifeq ($(PACKAGE_MANAGER),yum)
	@echo "📦 Installing with yum..."
	@command -v ffmpeg >/dev/null 2>&1 || (echo "Installing FFmpeg..." && sudo yum install -y ffmpeg)
	@command -v ollama >/dev/null 2>&1 || (echo "Installing Ollama..." && curl -fsSL https://ollama.com/install.sh | sh)
	@command -v piper >/dev/null 2>&1 || (echo "⚠️  Piper must be installed manually from https://github.com/rhasspy/piper/releases")
else
	@echo "❌ Unknown package manager. Please install prerequisites manually:"
	@echo "  - FFmpeg: https://ffmpeg.org/download.html"
	@echo "  - Ollama: https://ollama.com"
	@echo "  - Piper TTS: https://github.com/rhasspy/piper/releases"
	@exit 1
endif
else
	@echo "❌ Unsupported OS. Please install prerequisites manually:"
	@echo "  - FFmpeg: https://ffmpeg.org/download.html"
	@echo "  - Ollama: https://ollama.com"
	@echo "  - Piper TTS: https://github.com/rhasspy/piper/releases"
	@exit 1
endif
	@echo ""
	@echo "✅ Prerequisites installation complete!"
	@echo ""
	@echo "📝 Next steps:"
	@echo "  1. Download a Piper voice model:"
	@echo "     mkdir -p $(HOME)/.local/share/piper/voices"
	@echo "     cd $(HOME)/.local/share/piper/voices"
	@echo "     wget https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/lessac/medium/en_US-lessac-medium.onnx"
	@echo "     wget https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/lessac/medium/en_US-lessac-medium.onnx.json"
	@echo "  2. Pull an Ollama model: ollama pull llama3"
	@echo "  3. Start Ollama: ollama serve"
	@echo ""

install: check-prereqs install-prereqs build ## Install binary and prerequisites to GOPATH/bin
	@echo ""
	@echo "Installing Cortex binary..."
	@mkdir -p $(GOPATH)/bin
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
	@echo "✅ Cortex installed to $(GOPATH)/bin/$(BINARY_NAME)"
	@echo ""
	@echo "🎉 Installation complete! Run 'cortex --help' to get started."

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
