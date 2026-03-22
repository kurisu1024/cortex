package video

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	case "ai-generated":
		// This will be called with image paths from job manager
		err = fmt.Errorf("ai-generated background requires GenerateFromImages method")
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

// GenerateFromImages creates a video from AI-generated images with Ken Burns effects
func (g *Generator) GenerateFromImages(imagePaths []string, audioPath, outputPath string) error {
	if len(imagePaths) == 0 {
		return fmt.Errorf("no images provided")
	}

	fmt.Println("\n🎬 Generating video from AI images...")

	// Get audio duration to calculate timing
	audioDuration, err := g.getAudioDuration(audioPath)
	if err != nil {
		return fmt.Errorf("failed to get audio duration: %w", err)
	}

	// Calculate duration per image (in seconds)
	durationPerImage := audioDuration / float64(len(imagePaths))

	// Build ffmpeg command with inputs
	var args []string

	// Add all image inputs
	for _, imgPath := range imagePaths {
		args = append(args, "-loop", "1", "-t", fmt.Sprintf("%.2f", durationPerImage), "-i", imgPath)
	}

	// Add audio input
	args = append(args, "-i", audioPath)

	// Build filter complex for Ken Burns effects
	var filterParts []string
	for i := range imagePaths {
		var zoomFilter string
		if i%2 == 0 {
			// Zoom in effect
			zoomFilter = fmt.Sprintf("[%d:v]scale=1920:1080:force_original_aspect_ratio=increase,crop=1920:1080,"+
				"zoompan=z='min(zoom+0.0015,1.5)':d=%d:s=1920x1080:fps=25[v%d]",
				i, int(durationPerImage*25), i)
		} else {
			// Zoom out effect
			zoomFilter = fmt.Sprintf("[%d:v]scale=1920:1080:force_original_aspect_ratio=increase,crop=1920:1080,"+
				"zoompan=z='if(lte(zoom,1.0),1.5,max(1.0,zoom-0.0015))':d=%d:s=1920x1080:fps=25[v%d]",
				i, int(durationPerImage*25), i)
		}
		filterParts = append(filterParts, zoomFilter)
	}

	// Concatenate all video segments
	var concatInputs string
	for i := range imagePaths {
		concatInputs += fmt.Sprintf("[v%d]", i)
	}
	concatFilter := fmt.Sprintf("%sconcat=n=%d:v=1:a=0[outv]", concatInputs, len(imagePaths))
	filterParts = append(filterParts, concatFilter)

	filterComplex := strings.Join(filterParts, ";")

	args = append(args,
		"-filter_complex", filterComplex,
		"-map", "[outv]",
		"-map", fmt.Sprintf("%d:a", len(imagePaths)), // Audio is the last input
		"-c:v", "libx264",
		"-c:a", "aac",
		"-b:a", "192k",
		"-shortest",
		"-y",
		outputPath,
	)

	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("✅ Video saved to: %s\n", outputPath)
	return nil
}

// getAudioDuration returns the duration of an audio file in seconds
func (g *Generator) getAudioDuration(audioPath string) (float64, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		audioPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("ffprobe failed: %w", err)
	}

	var duration float64
	_, err = fmt.Sscanf(strings.TrimSpace(string(output)), "%f", &duration)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %w", err)
	}

	return duration, nil
}
