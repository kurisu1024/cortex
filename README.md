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
git clone https://github.com/kutidu2048/cortex.git
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

### Download Models

After installation, download the required AI models:

```bash
# Download Piper TTS voices (required for text-to-speech)
make download-voices

# Download Ollama LLM model (required for script generation)
make download-model

# Start Ollama service
ollama serve
```

## Usage

### Check Model Status

```bash
cortex status
```

### Generate a Video

```bash
# Basic usage (10 minute video by default)
cortex generate "The history of artificial intelligence"

# Custom duration (5 minute video)
cortex generate "Quick intro to Python" --duration 5
cortex generate "Quick intro to Python" -d 5  # Short form

# Longer video (20 minutes)
cortex generate "Deep dive into machine learning" -d 20

# With custom options
cortex generate "Space exploration" \
  --output ./videos \
  --duration 15 \
  --voice en_US-lessac-medium \
  --background gradient

# Use only high-quality voices (better audio quality)
cortex generate "Quantum computing explained" \
  --high-voices-only \
  --duration 8

# Combine multiple flags
cortex generate "AI Ethics" -H -d 12 -o ./videos
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
    # Default voice
    voice: ~/.local/share/piper/voices/en_US-lessac-medium.onnx

    # Multiple voices (optional) - rotates through voices for different segments
    voices:
      narrator: ~/.local/share/piper/voices/en_US-lessac-medium.onnx
      host: ~/.local/share/piper/voices/en_US-ryan-high.onnx
      expert: ~/.local/share/piper/voices/en_US-amy-medium.onnx

output:
  directory: ./output
  video:
    format: mp4
    background: gradient  # Options: gradient, solid, image
    waveform: true
```

### Commands

```bash
# Model management
cortex start         # Start local AI models
cortex stop          # Stop local AI models
cortex status        # Check model health

# Content generation
cortex generate      # Generate script, audio, and video

# Voice management
make download-voices         # Download all medium + high quality voices
make download-voice-medium   # Download only medium quality voices
make download-voice-high     # Download only high quality voices

# Help
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

# Download Piper voices
make download-voices          # All voices (medium + high)
make download-voice-medium    # Only medium quality
make download-voice-high      # Only high quality

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

## Multiple Voice Support

Cortex supports using multiple voices that automatically rotate between segments, creating a more dynamic and engaging video.

### Configuration

Add multiple voices to your `.cortex.yaml`:

```yaml
models:
  tts:
    engine: piper
    voice: ~/.local/share/piper/voices/en_US-lessac-medium.onnx  # Default fallback
    voices:
      narrator: ~/.local/share/piper/voices/en_US-lessac-medium.onnx  # Female voice
      host: ~/.local/share/piper/voices/en_US-ryan-high.onnx         # Male voice
      expert: ~/.local/share/piper/voices/en_US-amy-medium.onnx      # Alternative female
```

### How It Works

1. **Automatic Rotation**: Each script segment is assigned a different voice in rotation
2. **Speaker Names**: The names (narrator, host, expert) are used for logging during generation
3. **Fallback**: If voices aren't configured, uses the default `voice` setting
4. **Quality Filtering**: Use `--high-voices-only` or `-H` flag to use only high-quality voices (filters voices with "-high" in the path)

### Example Output

```
🎙️  Generating audio for 6 segments...
  [1/6] narrator: Generating audio...
  [2/6] host: Generating audio...
  [3/6] expert: Generating audio...
  [4/6] narrator: Generating audio...
  ...
```

### Downloading Multiple Voices

**Easy Way (Recommended):**

```bash
# Download all medium and high quality voices
make download-voices

# Or download specific quality levels
make download-voice-medium   # ~100MB - Faster, smaller files
make download-voice-high     # ~150MB - Better quality, larger files
```

This will download:
- **Medium Quality (5 voices)**: lessac, amy, joe, kristin, ryan
- **High Quality (3 voices)**: ryan, ljspeech, libritts_r

**Manual Download (if needed):**

```bash
cd ~/.local/share/piper/voices

# Female voice (lessac)
wget https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/lessac/medium/en_US-lessac-medium.onnx
wget https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/lessac/medium/en_US-lessac-medium.onnx.json

# Male voice (ryan)
wget https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/ryan/high/en_US-ryan-high.onnx
wget https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/ryan/high/en_US-ryan-high.onnx.json

# More voices available at: https://huggingface.co/rhasspy/piper-voices/tree/main/en/en_US
```

## Command Line Flags

### Generate Command

- `-o, --output <dir>` - Output directory for generated files (default: `./output`)
- `-d, --duration <minutes>` - Target video duration in minutes (default: `10`, max: `60`)
- `-v, --voice <path>` - Specific TTS voice to use (overrides config)
- `-b, --background <style>` - Video background style: `gradient`, `solid`, or `image` (default: `gradient`)
- `-H, --high-voices-only` - Use only high-quality voices from configuration
- `--config <file>` - Custom config file path
- `--verbose` - Enable verbose output

### Examples

```bash
# Use high-quality voices only
cortex generate "AI Ethics" -H

# Custom duration (15 minute video)
cortex generate "Machine Learning" -d 15

# Custom output directory
cortex generate "Quantum Computing" -o ~/videos

# Combine multiple flags
cortex generate "Deep Learning" -H -d 20 -o ./output -b solid
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

**No high-quality voices available:**
- Download high-quality voices: `make download-voice-high`
- Or download all voices: `make download-voices`

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
