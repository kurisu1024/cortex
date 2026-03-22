package video

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Generator handles video generation from audio
type Generator struct {
	visualizer *Visualizer
}

// NewGenerator creates a new video generator
func NewGenerator() *Generator {
	return &Generator{
		visualizer: NewVisualizer(),
	}
}

// GenerateFromAudio creates a video from an audio file
func (g *Generator) GenerateFromAudio(audioPath, outputPath, background string, showWaveform bool) error {
	// Check if ffmpeg is installed
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg is not installed. Please install it first")
	}

	fmt.Println("\n🎬 Generating video...")

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	var err error

	switch background {
	case "gradient":
		err = g.generateWithGradient(audioPath, outputPath, showWaveform)
	case "solid":
		err = g.generateWithSolidColor(audioPath, outputPath, showWaveform)
	case "image":
		err = g.generateWithImage(audioPath, outputPath, "", showWaveform)
	default:
		err = g.generateWithGradient(audioPath, outputPath, showWaveform)
	}

	if err != nil {
		return err
	}

	fmt.Printf("✅ Video saved to: %s\n", outputPath)

	return nil
}

// generateWithGradient creates a video with a gradient background
func (g *Generator) generateWithGradient(audioPath, outputPath string, showWaveform bool) error {
	var filterComplex string

	if showWaveform {
		// Create gradient background with waveform overlay
		filterComplex = "gradients=s=1920x1080:c0=0x0f0c29:c1=0x302b63:c2=0x24243e[gradient];" +
			"[0:a]showwaves=s=1920x1080:mode=line:rate=25:colors=0x00FF00[waves];" +
			"[gradient][waves]blend=all_mode=screen:all_opacity=0.6"
	} else {
		// Just gradient background
		filterComplex = "gradients=s=1920x1080:c0=0x0f0c29:c1=0x302b63:c2=0x24243e"
	}

	cmd := exec.Command("ffmpeg",
		"-i", audioPath,
		"-filter_complex", filterComplex,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-b:a", "192k",
		"-shortest",
		"-y",
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// generateWithSolidColor creates a video with a solid color background
func (g *Generator) generateWithSolidColor(audioPath, outputPath string, showWaveform bool) error {
	var filterComplex string

	if showWaveform {
		filterComplex = "color=c=0x1a1a2e:s=1920x1080:d=10[bg];" +
			"[0:a]showwaves=s=1920x1080:mode=cline:rate=25:colors=0x00FF00|0x0080FF[waves];" +
			"[bg][waves]overlay=0:0"
	} else {
		filterComplex = "color=c=0x1a1a2e:s=1920x1080"
	}

	cmd := exec.Command("ffmpeg",
		"-i", audioPath,
		"-filter_complex", filterComplex,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-b:a", "192k",
		"-shortest",
		"-y",
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// generateWithImage creates a video with an image background
func (g *Generator) generateWithImage(audioPath, outputPath, imagePath string, showWaveform bool) error {
	if imagePath == "" {
		// Fall back to solid color if no image provided
		return g.generateWithSolidColor(audioPath, outputPath, showWaveform)
	}

	var args []string

	if showWaveform {
		args = []string{
			"-loop", "1",
			"-i", imagePath,
			"-i", audioPath,
			"-filter_complex",
			"[1:a]showwaves=s=1920x1080:mode=cline:rate=25:colors=0x00FF00|0x0080FF:scale=sqrt[waves];" +
				"[0:v]scale=1920:1080[bg];" +
				"[bg][waves]overlay=0:0",
			"-c:v", "libx264",
			"-c:a", "aac",
			"-b:a", "192k",
			"-shortest",
			"-y",
			outputPath,
		}
	} else {
		args = []string{
			"-loop", "1",
			"-i", imagePath,
			"-i", audioPath,
			"-c:v", "libx264",
			"-c:a", "aac",
			"-b:a", "192k",
			"-shortest",
			"-y",
			outputPath,
		}
	}

	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// AddSubtitles adds subtitles to an existing video
func (g *Generator) AddSubtitles(videoPath, subtitlesPath, outputPath string) error {
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-vf", fmt.Sprintf("subtitles=%s", subtitlesPath),
		"-c:a", "copy",
		"-y",
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}
