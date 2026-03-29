package audio

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kutidu2048/cortex/internal/models"
	"github.com/kutidu2048/cortex/internal/script"
)

// Generator handles audio generation from script segments
type Generator struct {
	tts *models.TTSClient
}

// NewGenerator creates a new audio generator
func NewGenerator(tts *models.TTSClient) *Generator {
	return &Generator{tts: tts}
}

// SegmentInfo holds audio path and duration for a segment
type SegmentInfo struct {
	Path     string
	Duration float64
}

// GenerateFromScript generates audio files for each script segment and returns paths with durations
func (g *Generator) GenerateFromScript(scr *script.Script, outputDir string) ([]SegmentInfo, error) {
	// Create output directory if it doesn't exist
	audioDir := filepath.Join(outputDir, "audio_segments")
	if err := os.MkdirAll(audioDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create audio directory: %w", err)
	}

	segments := make([]SegmentInfo, 0, len(scr.Segments))
	combiner := NewCombiner()

	fmt.Printf("\n🎙️  Generating audio for %d segments...\n", len(scr.Segments))

	for i, segment := range scr.Segments {
		outputPath := filepath.Join(audioDir, fmt.Sprintf("segment_%03d.wav", i))

		// Show which speaker/voice is being used
		speaker := segment.Speaker
		if speaker == "" {
			speaker = "Narrator"
		}
		fmt.Printf("  [%d/%d] %s: Generating audio...\n", i+1, len(scr.Segments), speaker)

		// Use segment-specific voice if available
		if segment.VoicePath != "" {
			if err := g.tts.GenerateAudioWithVoice(segment.Text, outputPath, segment.VoicePath); err != nil {
				return nil, fmt.Errorf("failed to generate audio for segment %d: %w", i, err)
			}
		} else {
			// Fall back to default voice
			if err := g.tts.GenerateAudio(segment.Text, outputPath); err != nil {
				return nil, fmt.Errorf("failed to generate audio for segment %d: %w", i, err)
			}
		}

		// Get duration of the generated audio
		duration, err := combiner.GetDuration(outputPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get duration for segment %d: %w", i, err)
		}

		segments = append(segments, SegmentInfo{
			Path:     outputPath,
			Duration: duration,
		})
	}

	fmt.Printf("✅ Generated %d audio segments\n", len(segments))

	return segments, nil
}

// GenerateSegment generates audio for a single segment
func (g *Generator) GenerateSegment(text, outputPath string) error {
	return g.tts.GenerateAudio(text, outputPath)
}
