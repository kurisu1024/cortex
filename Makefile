.PHONY: build install clean test run help check-prereqs install-prereqs download-voices download-voice-medium download-voice-high download-model

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

# Piper voices configuration
VOICE_DIR=$(HOME)/.local/share/piper/voices
VOICE_BASE_URL=https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US

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
	@echo "  1. Download Piper voices: make download-voices"
	@echo "  2. Download Ollama model: make download-model"
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

# Voice download helpers
download-voice-medium: ## Download all medium quality en_US voices
	@echo "📥 Downloading medium quality Piper voices..."
	@mkdir -p $(VOICE_DIR)
	@echo ""
	@echo "Downloading lessac (female, medium)..."
	@curl -L -o $(VOICE_DIR)/en_US-lessac-medium.onnx $(VOICE_BASE_URL)/lessac/medium/en_US-lessac-medium.onnx
	@curl -L -o $(VOICE_DIR)/en_US-lessac-medium.onnx.json $(VOICE_BASE_URL)/lessac/medium/en_US-lessac-medium.onnx.json
	@echo "✅ lessac-medium downloaded"
	@echo ""
	@echo "Downloading amy (female, medium)..."
	@curl -L -o $(VOICE_DIR)/en_US-amy-medium.onnx $(VOICE_BASE_URL)/amy/medium/en_US-amy-medium.onnx
	@curl -L -o $(VOICE_DIR)/en_US-amy-medium.onnx.json $(VOICE_BASE_URL)/amy/medium/en_US-amy-medium.onnx.json
	@echo "✅ amy-medium downloaded"
	@echo ""
	@echo "Downloading joe (male, medium)..."
	@curl -L -o $(VOICE_DIR)/en_US-joe-medium.onnx $(VOICE_BASE_URL)/joe/medium/en_US-joe-medium.onnx
	@curl -L -o $(VOICE_DIR)/en_US-joe-medium.onnx.json $(VOICE_BASE_URL)/joe/medium/en_US-joe-medium.onnx.json
	@echo "✅ joe-medium downloaded"
	@echo ""
	@echo "Downloading kristin (female, medium)..."
	@curl -L -o $(VOICE_DIR)/en_US-kristin-medium.onnx $(VOICE_BASE_URL)/kristin/medium/en_US-kristin-medium.onnx
	@curl -L -o $(VOICE_DIR)/en_US-kristin-medium.onnx.json $(VOICE_BASE_URL)/kristin/medium/en_US-kristin-medium.onnx.json
	@echo "✅ kristin-medium downloaded"
	@echo ""
	@echo "Downloading ryan (male, medium)..."
	@curl -L -o $(VOICE_DIR)/en_US-ryan-medium.onnx $(VOICE_BASE_URL)/ryan/medium/en_US-ryan-medium.onnx
	@curl -L -o $(VOICE_DIR)/en_US-ryan-medium.onnx.json $(VOICE_BASE_URL)/ryan/medium/en_US-ryan-medium.onnx.json
	@echo "✅ ryan-medium downloaded"
	@echo ""
	@echo "✅ All medium quality voices downloaded to $(VOICE_DIR)"

download-voice-high: ## Download all high quality en_US voices
	@echo "📥 Downloading high quality Piper voices..."
	@mkdir -p $(VOICE_DIR)
	@echo ""
	@echo "Downloading ryan (male, high)..."
	@curl -L -o $(VOICE_DIR)/en_US-ryan-high.onnx $(VOICE_BASE_URL)/ryan/high/en_US-ryan-high.onnx
	@curl -L -o $(VOICE_DIR)/en_US-ryan-high.onnx.json $(VOICE_BASE_URL)/ryan/high/en_US-ryan-high.onnx.json
	@echo "✅ ryan-high downloaded"
	@echo ""
	@echo "Downloading ljspeech (female, high)..."
	@curl -L -o $(VOICE_DIR)/en_US-ljspeech-high.onnx $(VOICE_BASE_URL)/ljspeech/high/en_US-ljspeech-high.onnx
	@curl -L -o $(VOICE_DIR)/en_US-ljspeech-high.onnx.json $(VOICE_BASE_URL)/ljspeech/high/en_US-ljspeech-high.onnx.json
	@echo "✅ ljspeech-high downloaded"
	@echo ""
	@echo "Downloading libritts_r (multi-speaker, high)..."
	@curl -L -o $(VOICE_DIR)/en_US-libritts_r-high.onnx $(VOICE_BASE_URL)/libritts_r/high/en_US-libritts_r-high.onnx
	@curl -L -o $(VOICE_DIR)/en_US-libritts_r-high.onnx.json $(VOICE_BASE_URL)/libritts_r/high/en_US-libritts_r-high.onnx.json
	@echo "✅ libritts_r-high downloaded"
	@echo ""
	@echo "✅ All high quality voices downloaded to $(VOICE_DIR)"

download-voices: download-voice-medium download-voice-high ## Download all medium and high quality voices
	@echo ""
	@echo "🎉 All Piper voices downloaded successfully!"
	@echo ""
	@echo "📁 Voices installed in: $(VOICE_DIR)"
	@echo ""
	@echo "Available voices:"
	@echo "  Medium Quality:"
	@echo "    - lessac (female)"
	@echo "    - amy (female)"
	@echo "    - joe (male)"
	@echo "    - kristin (female)"
	@echo "    - ryan (male)"
	@echo ""
	@echo "  High Quality:"
	@echo "    - ryan (male)"
	@echo "    - ljspeech (female)"
	@echo "    - libritts_r (multi-speaker)"
	@echo ""
	@echo "💡 Update your .cortex.yaml to use these voices!"

download-model: ## Download the default Ollama model (llama3.2)
	@echo "📥 Downloading Ollama model llama3.2..."
	@ollama pull llama3.2
	@echo "✅ llama3.2 model downloaded successfully!"
	@echo ""
	@echo "💡 You can now use 'cortex generate' to create videos"

.DEFAULT_GOAL := help
