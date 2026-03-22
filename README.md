# Cortex

**AI-Powered Script to Video Generator** - Generate engaging videos from text topics using local AI models. No subscriptions required.

```
╔═══════════════════════════════════════════════════════════╗
║   ██████╗ ██████╗ ██████╗ ████████╗███████╗██╗  ██╗    ║
║  ██╔════╝██╔═══██╗██╔══██╗╚══██╔══╝██╔════╝╚██╗██╔╝    ║
║  ██║     ██║   ██║██████╔╝   ██║   █████╗   ╚███╔╝     ║
║  ██║     ██║   ██║██╔══██╗   ██║   ██╔══╝   ██╔██╗     ║
║  ╚██████╗╚██████╔╝██║  ██║   ██║   ███████╗██╔╝ ██╗    ║
║   ╚═════╝ ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚══════╝╚═╝  ╚═╝    ║
╚═══════════════════════════════════════════════════════════╝
```

## Features

- 🧠 **Local LLM Integration** - Uses Ollama for script generation (no API keys needed)
- 🎙️ **Local TTS** - Piper TTS for high-quality speech synthesis
- 🎬 **Video Generation** - Automatic video creation with audio waveform visualizations
- 🎨 **Customizable Backgrounds** - Gradient, solid color, or image backgrounds
- 📊 **Progress Tracking** - Real-time progress with hacker-style terminal UI
- ⚙️ **Configurable** - Viper + Cobra for flexible configuration and CLI

## Quick Start

### Prerequisites

1. **Ollama** - For LLM script generation
   ```bash
   # Install Ollama
   curl -fsSL https://ollama.com/install.sh | sh

   # Pull a model
   ollama pull llama3

   # Start Ollama
   ollama serve
   ```

2. **Piper TTS** - For text-to-speech
   ```bash
   # Download Piper binary from GitHub releases
   # https://github.com/rhasspy/piper/releases

   # macOS example:
   mkdir -p ~/.local/bin
   curl -L https://github.com/rhasspy/piper/releases/download/2023.11.14-2/piper_macos_x64.tar.gz -o piper.tar.gz
   tar -xzf piper.tar.gz
   cp piper/piper ~/.local/bin/
   export PATH="$HOME/.local/bin:$PATH"

   # Download a voice model
   mkdir -p ~/.local/share/piper/voices
   cd ~/.local/share/piper/voices
   wget https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/lessac/medium/en_US-lessac-medium.onnx
   wget https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/lessac/medium/en_US-lessac-medium.onnx.json
   ```

3. **FFmpeg** - For audio/video processing
   ```bash
   # macOS
   brew install ffmpeg

   # Linux
   sudo apt-get install ffmpeg
   ```

### Installation

```bash
# Clone the repository
git clone https://github.com/topher/cortex.git
cd cortex

# Check prerequisites
make check-prereqs

# Install everything (prerequisites + binary)
make install

# Or just build (without installing prerequisites)
make build
./bin/cortex
```

The `make install` command will:
1. Check if prerequisites are installed
2. Automatically install missing prerequisites (FFmpeg, Ollama, Piper)
3. Build the Cortex binary
4. Install to `$GOPATH/bin`

## Usage

### Check Model Status

```bash
cortex status
```

### Generate a Video

```bash
# Basic usage
cortex generate "The history of artificial intelligence"

# With custom options
cortex generate "Space exploration" \
  --output ./videos \
  --voice en_US-lessac-medium \
  --background gradient
```

### Configuration

Create a `.cortex.yaml` file in your home directory or project directory:

```yaml
models:
  ollama:
    host: http://localhost:11434
    model: llama3
  tts:
    engine: piper
    voice: ~/.local/share/piper/voices/en_US-lessac-medium.onnx  # Path to .onnx model

output:
  directory: ./output
  video:
    format: mp4
    background: gradient  # Options: gradient, solid, image
    waveform: true
```

### Commands

```bash
cortex start         # Start local AI models
cortex stop          # Stop local AI models
cortex status        # Check model health
cortex generate      # Generate script, audio, and video
cortex --help        # Show all commands
```

## Pipeline

Cortex follows a multi-step pipeline:

1. **Script Generation** - Uses Ollama LLM to create an engaging script
2. **Text-to-Speech** - Converts script segments to audio using Piper
3. **Audio Combination** - Merges audio segments with transitions
4. **Video Creation** - Combines audio with visualizations into MP4

## Project Structure

```
cortex/
├── cmd/cortex/              # CLI entry point
├── pkg/commands/            # Cobra CLI commands
├── internal/
│   ├── config/             # Viper configuration
│   ├── models/             # LLM and TTS clients
│   ├── script/             # Script generation
│   ├── audio/              # Audio generation & combination
│   ├── video/              # Video generation
│   ├── job/                # Job management
│   └── ui/                 # Terminal UI
├── .cortex.yaml            # Configuration file
├── Makefile                # Build commands
└── README.md
```

## Development

```bash
# Check prerequisites
make check-prereqs

# Install missing prerequisites only
make install-prereqs

# Install Go dependencies
make deps

# Format code
make fmt

# Run tests
make test

# Build
make build

# Run linter
make lint

# Quick: format, build, and test
make quick
```

## Environment Variables

Override config with environment variables (prefix: `CORTEX_`):

```bash
export CORTEX_MODELS_OLLAMA_HOST=http://localhost:11434
export CORTEX_MODELS_OLLAMA_MODEL=mistral
export CORTEX_OUTPUT_DIRECTORY=./my-videos
```

## Video Customization

### Background Styles

- **gradient** - Smooth color gradients (default)
- **solid** - Solid dark background
- **image** - Use a custom background image

### Waveform Visualization

Enable/disable audio waveform overlay in config:

```yaml
output:
  video:
    waveform: true  # or false
```

## Troubleshooting

**Ollama not running:**
```bash
ollama serve
```

**Piper not found:**
- Ensure Piper is in your PATH
- Install from: https://github.com/rhasspy/piper

**FFmpeg errors:**
- Install/update FFmpeg: `brew install ffmpeg`

## Roadmap

- [ ] Web UI for monitoring jobs (React dashboard)
- [ ] Multiple voice support with different speakers
- [ ] Custom video templates and transitions
- [ ] Subtitle generation from script
- [ ] Export to multiple formats
- [ ] Batch processing

## License

MIT License - See LICENSE file

## Contributing

Contributions welcome! Please open an issue or PR.

---

**Built with** [Cobra](https://github.com/spf13/cobra) • [Viper](https://github.com/spf13/viper) • [Ollama](https://ollama.com) • [Piper TTS](https://github.com/rhasspy/piper) • [FFmpeg](https://ffmpeg.org)
