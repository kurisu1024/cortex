package audio

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/topher/cortex/internal/models"
	"github.com/topher/cortex/internal/script"
)

// Generator handles audio generation from script segments
type Generator struct {
	tts *models.TTSClient
}

// NewGenerator creates a new audio generator
func NewGenerator(tts *models.TTSClient) *Generator {
	return &Generator{tts: tts}
}

// GenerateFromScript generates audio files for each script segment
func (g *Generator) GenerateFromScript(scr *script.Script, outputDir string) ([]string, error) {
	// Create output directory if it doesn't exist
	audioDir := filepath.Join(outputDir, "audio_segments")
	if err := os.MkdirAll(audioDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create audio directory: %w", err)
	}

	var audioPaths []string

	fmt.Printf("\n🎙️  Generating audio for %d segments...\n", len(scr.Segments))

	for i, segment := range scr.Segments {
		outputPath := filepath.Join(audioDir, fmt.Sprintf("segment_%03d.wav", i))

		fmt.Printf("  [%d/%d] Generating audio...\n", i+1, len(scr.Segments))

		if err := g.tts.GenerateAudio(segment.Text, outputPath); err != nil {
			return nil, fmt.Errorf("failed to generate audio for segment %d: %w", i, err)
		}

		audioPaths = append(audioPaths, outputPath)
	}

	fmt.Printf("✅ Generated %d audio segments\n", len(audioPaths))

	return audioPaths, nil
}

// GenerateSegment generates audio for a single segment
func (g *Generator) GenerateSegment(text, outputPath string) error {
	return g.tts.GenerateAudio(text, outputPath)
}
