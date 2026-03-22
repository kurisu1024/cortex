package models

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// TTSClient handles text-to-speech operations
type TTSClient struct {
	engine string
	voice  string
}

// NewTTSClient creates a new TTS client
func NewTTSClient(engine, voice string) *TTSClient {
	return &TTSClient{
		engine: engine,
		voice:  voice,
	}
}

// Start ensures TTS engine is available
func (t *TTSClient) Start() error {
	if !t.IsHealthy() {
		return fmt.Errorf("%s is not installed. Please install it first", t.engine)
	}
	return nil
}

// Stop is a no-op for TTS
func (t *TTSClient) Stop() error {
	return nil
}

// IsHealthy checks if TTS engine is available
func (t *TTSClient) IsHealthy() bool {
	return t.findPiperPath() != ""
}

// findPiperPath finds the piper executable in PATH or common locations
func (t *TTSClient) findPiperPath() string {
	// First try exec.LookPath (checks PATH)
	if path, err := exec.LookPath(t.engine); err == nil {
		return path
	}

	// Try common installation locations
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	commonPaths := []string{
		filepath.Join(homeDir, ".local", "bin", "piper", "piper"), // piper in extracted directory
		filepath.Join(homeDir, ".local", "bin", "piper"),          // standalone piper binary
		"/usr/local/bin/piper",
		"/opt/homebrew/bin/piper",
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// GenerateAudio generates audio from text using the default voice
func (t *TTSClient) GenerateAudio(text, outputPath string) error {
	return t.GenerateAudioWithVoice(text, outputPath, t.voice)
}

// GenerateAudioWithVoice generates audio from text using a specific voice
func (t *TTSClient) GenerateAudioWithVoice(text, outputPath, voicePath string) error {
	switch t.engine {
	case "piper":
		return t.generateWithPiperVoice(text, outputPath, voicePath)
	case "say":
		return t.generateWithSay(text, outputPath, voicePath)
	default:
		return fmt.Errorf("unsupported TTS engine: %s", t.engine)
	}
}

// generateWithPiperVoice uses Piper TTS to generate audio with a specific voice
func (t *TTSClient) generateWithPiperVoice(text, outputPath, voicePath string) error {
	// Piper command: echo "text" | piper --model voice.onnx --output_file output.wav
	// Note: voicePath should be the path to the .onnx model file

	// Find piper executable
	piperPath := t.findPiperPath()
	if piperPath == "" {
		return fmt.Errorf(`piper executable not found

Installation instructions for macOS:
1. Install espeak-ng: brew install espeak-ng
2. Download piper from: https://github.com/rhasspy/piper/releases
3. Or use the official Docker image

For more info: https://github.com/rhasspy/piper`)
	}

	// Expand tilde in voice path
	if len(voicePath) > 0 && voicePath[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			voicePath = filepath.Join(homeDir, voicePath[1:])
		}
	}

	cmd := exec.Command(piperPath, "--model", voicePath, "--output_file", outputPath)

	// Set library path for macOS (where espeak-ng is installed via Homebrew)
	cmd.Env = append(os.Environ(),
		"DYLD_LIBRARY_PATH=/opt/homebrew/lib:/usr/local/lib",
	)

	// Create a pipe for stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start piper: %w", err)
	}

	// Write text to stdin
	if _, err := stdin.Write([]byte(text)); err != nil {
		return fmt.Errorf("failed to write to piper stdin: %w", err)
	}
	stdin.Close()

	// Wait for completion
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("piper failed: %w", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return fmt.Errorf("audio file was not created: %s", outputPath)
	}

	return nil
}

// generateWithSay uses macOS native say command to generate audio
func (t *TTSClient) generateWithSay(text, outputPath, voiceName string) error {
	// macOS say command: say -v VoiceName -o output.aiff "text"
	// Note: voiceName should be a macOS voice name (e.g., "Samantha", "Alex")

	// If no voice specified, use default
	if voiceName == "" {
		voiceName = "Samantha" // Default to Samantha voice
	}

	// say outputs AIFF format, but we need WAV for consistency
	// Generate to temp AIFF first, then convert to WAV
	tempAiff := outputPath + ".aiff"

	// Build say command
	cmd := exec.Command("say", "-v", voiceName, "-o", tempAiff, text)

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("say command failed: %w", err)
	}

	// Verify temp AIFF was created
	if _, err := os.Stat(tempAiff); os.IsNotExist(err) {
		return fmt.Errorf("audio file was not created: %s", tempAiff)
	}

	// Convert AIFF to WAV using ffmpeg
	convertCmd := exec.Command("ffmpeg", "-i", tempAiff, "-y", outputPath)
	if err := convertCmd.Run(); err != nil {
		os.Remove(tempAiff) // Clean up temp file
		return fmt.Errorf("ffmpeg conversion failed: %w", err)
	}

	// Clean up temp AIFF file
	os.Remove(tempAiff)

	// Verify final WAV exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return fmt.Errorf("converted audio file was not created: %s", outputPath)
	}

	return nil
}

// GetVoice returns the current voice
func (t *TTSClient) GetVoice() string {
	return t.voice
}

// SetVoice sets a new voice
func (t *TTSClient) SetVoice(voice string) {
	t.voice = voice
}
