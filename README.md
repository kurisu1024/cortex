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
- 🎙️ **High-Quality TTS** - Edge TTS with neural voices (default) or offline Piper TTS
- 🎨 **AI Image Generation** - SDXL-Turbo for fast, high-quality visuals (GPU-accelerated on Apple Silicon)
- 🎬 **Animated Videos** - AnimateDiff for fully animated cartoon characters and scenes
- 🖼️ **Flexible Backgrounds** - Gradient, solid color, custom images, or AI-generated scenes
- ⚡ **Optimized Performance** - fp16 precision, CPU offloading, memory-efficient generation
- 📊 **Progress Tracking** - Real-time progress with hacker-style terminal UI
- ⚙️ **Configurable** - Viper + Cobra for flexible configuration and CLI

## Quick Start

### Prerequisites

1. **Ollama** - For LLM script generation
   ```bash
   # Install Ollama
   curl -fsSL https://ollama.com/install.sh | sh

   # Pull a model
   ollama pull llama3.2

   # Start Ollama
   ollama serve
   ```

2. **Python 3.9+** - For AI image/animation generation
   ```bash
   # Check Python version
   python3 --version

   # Install required packages
   pip3 install torch torchvision torchaudio
   pip3 install diffusers transformers accelerate
   pip3 install opencv-python
   pip3 install edge-tts  # For text-to-speech
   ```

3. **FFmpeg** - For audio/video processing
   ```bash
   # macOS
   brew install ffmpeg

   # Linux
   sudo apt-get install ffmpeg
   ```

4. **Piper TTS** (Optional) - For offline text-to-speech
   ```bash
   # Only needed if you want offline TTS instead of Edge TTS
   # Download Piper binary from https://github.com/rhasspy/piper/releases

   # macOS example:
   mkdir -p ~/.local/bin
   curl -L https://github.com/rhasspy/piper/releases/download/2023.11.14-2/piper_macos_x64.tar.gz -o piper.tar.gz
   tar -xzf piper.tar.gz
   cp piper/piper ~/.local/bin/
   export PATH="$HOME/.local/bin:$PATH"

   # Download voice models
   make download-voices
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
# Download Ollama LLM model (required for script generation)
ollama pull llama3.2

# Start Ollama service
ollama serve
```

**Note:** Edge TTS works out of the box (requires internet). AI image models are downloaded automatically on first use.

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
  --background ai-generated

# AI-generated animated backgrounds
cortex generate "History of Animation" \
  --background ai-generated \
  --duration 10

# Static background with gradient
cortex generate "Quick tutorial" \
  --background gradient \
  -d 5
```

### Configuration

Create a `.cortex.yaml` file in your home directory or project directory:

```yaml
models:
  ollama:
    host: http://localhost:11434
    model: llama3.2

  image:
    # AI Image/Animation Generation
    model_id: stabilityai/sdxl-turbo  # Fast, high-quality image generation
    art_style: cinematic, high quality, 4k, detailed

  tts:
    # TTS Engine: "edgetts" (default, requires internet) or "piper" (offline)
    engine: edgetts

    # Default voice
    voice: en-US-AriaNeural

    # Multiple voices for dynamic narration
    voices:
      narrator: en-US-AriaNeural      # Female, professional
      host: en-US-GuyNeural           # Male, energetic
      expert: en-US-JennyNeural       # Female, friendly
      interviewer: en-GB-RyanNeural   # Male, British

    # Alternative: Piper TTS (offline, requires setup)
    # engine: piper
    # voice: ~/.local/share/piper/voices/en_US-lessac-medium.onnx

output:
  directory: ./output
  duration: 10  # Default video duration in minutes

  video:
    format: mp4
    background: ai-generated  # Options: gradient, solid, image, ai-generated
    waveform: false

    # AI Animation Settings
    animated: true           # Use AnimateDiff for animated scenes
    animation_frames: 10     # Frames per clip (8=1s, 10=1.25s, 12=1.5s, 16=2s)
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

1. **Script Generation** - Uses Ollama LLM to create an engaging script with scenes
2. **Text-to-Speech** - Converts script segments to audio using Edge TTS or Piper
3. **Audio Combination** - Merges audio segments with transitions
4. **Image/Animation Generation** - Creates AI visuals for each scene:
   - **Static mode**: SDXL-Turbo images with Ken Burns zoom effects (~5-10s per image)
   - **Animated mode**: AnimateDiff cartoon clips with character motion (~10-15s per clip)
5. **Video Creation** - Combines images/animations with audio into final MP4

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

Cortex supports multiple background modes:

- **ai-generated** (default) - AI-generated scenes with two modes:
  - **Static** (`animated: false`): SDXL-Turbo images with Ken Burns zoom/pan effects
    - Fast generation: ~5-10 seconds per image on M4 GPU
    - Low memory usage: ~1-2GB
  - **Animated** (`animated: true`): AnimateDiff cartoon clips with character motion
    - Optimized generation: ~10-15 seconds per clip on M4 GPU
    - Memory efficient: ~2-3GB peak with fp16 + CPU offloading
    - Configure frames: 8 (1s), 10 (1.25s), 12 (1.5s), or 16 (2s)

- **gradient** - Smooth color gradients (simple, no AI required)
- **solid** - Solid dark background
- **image** - Use a custom background image

### AI Animation Configuration

Configure animation in `.cortex.yaml`:

```yaml
output:
  video:
    background: ai-generated
    animated: true           # false for static images, true for animated clips
    animation_frames: 10     # Recommended: 10 frames (1.25s clips)
    waveform: false         # Optional audio waveform overlay
```

**Performance Guide (Apple Silicon M4):**
- 8 frames (1.0s): ~8-10 seconds per clip, ~1.5GB memory
- 10 frames (1.25s): ~10-12 seconds per clip, ~2GB memory ⭐ **Recommended**
- 12 frames (1.5s): ~12-14 seconds per clip, ~2.5GB memory
- 16 frames (2.0s): ~15-18 seconds per clip, ~3GB memory

**Optimizations Applied:**
- fp16 precision for 50% memory reduction
- CPU offloading to prevent OOM errors
- Reduced inference steps (15) for 2-3x speed boost
- Memory cleanup between clips

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
- `-b, --background <style>` - Video background style: `gradient`, `solid`, `image`, or `ai-generated` (default: `ai-generated`)
- `-H, --high-voices-only` - Use only high-quality voices from configuration (for Piper TTS)
- `--config <file>` - Custom config file path
- `--verbose` - Enable verbose output

### Examples

```bash
# AI-generated animated video (default)
cortex generate "AI Ethics"

# AI-generated static images (faster)
# Set animated: false in .cortex.yaml

# Custom duration (15 minute video)
cortex generate "Machine Learning" -d 15

# Simple gradient background (no AI generation)
cortex generate "Quick Tutorial" -b gradient -d 5

# Custom output directory
cortex generate "Quantum Computing" -o ~/videos

# Combine multiple flags
cortex generate "Deep Learning" -d 20 -o ./output -b ai-generated
```

## Troubleshooting

**Ollama not running:**
```bash
ollama serve
```

**Python/PyTorch issues:**
```bash
# Check Python version (need 3.9+)
python3 --version

# Reinstall PyTorch with MPS support (Apple Silicon)
pip3 install --upgrade torch torchvision torchaudio

# Install missing packages
pip3 install diffusers transformers accelerate opencv-python edge-tts
```

**GPU/MPS errors:**
```bash
# "Invalid buffer size" or "Out of memory" errors:
# 1. Reduce animation_frames in .cortex.yaml (try 8 instead of 10-16)
# 2. Close other GPU-intensive applications
# 3. Set animated: false for static images (uses less memory)

# Black images or MPS errors:
# System will automatically fall back to CPU if MPS fails
```

**Edge TTS connection errors:**
```bash
# Requires internet connection
# Alternative: Use offline Piper TTS
# Set engine: piper in .cortex.yaml and install voices
```

**FFmpeg errors:**
```bash
# Install/update FFmpeg
brew install ffmpeg  # macOS
sudo apt-get install ffmpeg  # Linux
```

**Slow generation:**
```bash
# For faster AI generation:
# 1. Reduce animation_frames to 8 (1 second clips)
# 2. Use animated: false for static images (~5-10s vs 10-15s)
# 3. Ensure you're using GPU (check for MPS messages in output)
```

**Piper TTS issues (if using offline TTS):**
```bash
# Ensure Piper is in your PATH
# Install from: https://github.com/rhasspy/piper
# Download voices: make download-voices
```

## Roadmap

- [x] AI-generated images and animations
- [x] Multiple voice support with different speakers
- [x] GPU acceleration for Apple Silicon
- [ ] Web UI for monitoring jobs (React dashboard)
- [ ] Custom video templates and transitions
- [ ] Subtitle generation from script
- [ ] Export to multiple formats
- [ ] Batch processing
- [ ] Alternative animation models (AnimateDiff v2, v3)

## License

MIT License - See LICENSE file

## Contributing

Contributions welcome! Please open an issue or PR.

---

**Built with** [Cobra](https://github.com/spf13/cobra) • [Viper](https://github.com/spf13/viper) • [Ollama](https://ollama.com) • [SDXL-Turbo](https://huggingface.co/stabilityai/sdxl-turbo) • [AnimateDiff](https://github.com/guoyww/AnimateDiff) • [Edge TTS](https://github.com/rany2/edge-tts) • [FFmpeg](https://ffmpeg.org)
