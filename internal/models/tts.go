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
	_, err := exec.LookPath(t.engine)
	return err == nil
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
	default:
		return fmt.Errorf("unsupported TTS engine: %s", t.engine)
	}
}

// generateWithPiperVoice uses Piper TTS to generate audio with a specific voice
func (t *TTSClient) generateWithPiperVoice(text, outputPath, voicePath string) error {
	// Piper command: echo "text" | piper --model voice.onnx --output_file output.wav
	// Note: voicePath should be the path to the .onnx model file
	// Expand tilde in voice path
	if len(voicePath) > 0 && voicePath[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			voicePath = filepath.Join(homeDir, voicePath[1:])
		}
	}

	cmd := exec.Command("piper", "--model", voicePath, "--output_file", outputPath)

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

// GetVoice returns the current voice
func (t *TTSClient) GetVoice() string {
	return t.voice
}

// SetVoice sets a new voice
func (t *TTSClient) SetVoice(voice string) {
	t.voice = voice
}
