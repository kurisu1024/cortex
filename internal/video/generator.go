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
func (g *Generator) GenerateFromImages(imagePaths []string, audioPath, outputPath string, segmentDurations []float64) error {
	if len(imagePaths) == 0 {
		return fmt.Errorf("no images provided")
	}

	if len(imagePaths) != len(segmentDurations) {
		return fmt.Errorf("mismatch: %d images but %d segment durations", len(imagePaths), len(segmentDurations))
	}

	fmt.Printf("\n🎬 Generating video from %d AI images with Ken Burns effects...\n", len(imagePaths))
	for i, imgPath := range imagePaths {
		fmt.Printf("  Image %d: %s (%.2fs)\n", i+1, filepath.Base(imgPath), segmentDurations[i])
	}

	// Create temp directory for individual video clips
	tempClipsDir := filepath.Join(filepath.Dir(outputPath), "temp_image_clips")
	if err := os.MkdirAll(tempClipsDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp clips directory: %w", err)
	}
	defer os.RemoveAll(tempClipsDir)

	// Step 1: Create individual video clips from each image with Ken Burns effect
	var clipPaths []string
	for i, imagePath := range imagePaths {
		clipPath := filepath.Join(tempClipsDir, fmt.Sprintf("clip_%03d.mp4", i))
		duration := segmentDurations[i]
		durationFrames := int(duration * 25) // Convert to frames at 25fps

		fmt.Printf("  [%d/%d] Creating clip from %s...\n", i+1, len(imagePaths), filepath.Base(imagePath))

		var zoomFilter string
		if i%2 == 0 {
			// Zoom in effect
			zoomFilter = fmt.Sprintf("scale=1920:1080:force_original_aspect_ratio=increase,crop=1920:1080,"+
				"zoompan=z='min(zoom+0.0015,1.5)':d=%d:s=1920x1080:fps=25", durationFrames)
		} else {
			// Zoom out effect
			zoomFilter = fmt.Sprintf("scale=1920:1080:force_original_aspect_ratio=increase,crop=1920:1080,"+
				"zoompan=z='if(lte(zoom,1.0),1.5,max(1.0,zoom-0.0015))':d=%d:s=1920x1080:fps=25", durationFrames)
		}

		args := []string{
			"-loop", "1",
			"-i", imagePath,
			"-vf", zoomFilter,
			"-t", fmt.Sprintf("%.2f", duration),
			"-c:v", "libx264",
			"-pix_fmt", "yuv420p",
			"-y",
			clipPath,
		}

		cmd := exec.Command("ffmpeg", args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to create clip %d: %w\nOutput: %s", i, err, string(output))
		}

		clipPaths = append(clipPaths, clipPath)
	}

	fmt.Printf("\n✅ Created %d video clips\n", len(clipPaths))

	// Step 2: Create concat file
	concatFile := filepath.Join(tempClipsDir, "concat_list.txt")
	var concatContent strings.Builder
	for _, clipPath := range clipPaths {
		absPath, _ := filepath.Abs(clipPath)
		concatContent.WriteString(fmt.Sprintf("file '%s'\n", absPath))
	}

	if err := os.WriteFile(concatFile, []byte(concatContent.String()), 0644); err != nil {
		return fmt.Errorf("failed to write concat file: %w", err)
	}

	// Step 3: Concatenate all clips and add audio
	fmt.Printf("\n🎬 Concatenating %d clips and adding audio...\n", len(clipPaths))

	args := []string{
		"-f", "concat",
		"-safe", "0",
		"-i", concatFile,
		"-i", audioPath,
		"-c:v", "copy", // Copy video since clips are already encoded
		"-c:a", "aac",
		"-b:a", "192k",
		"-shortest",
		"-y",
		outputPath,
	}

	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg concat failed: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("✅ Video saved to: %s\n", outputPath)
	return nil
}

// GenerateFromAnimatedClips creates a final video by concatenating animated video clips
func (g *Generator) GenerateFromAnimatedClips(clipPaths []string, audioPath, outputPath string, segmentDurations []float64) error {
	if len(clipPaths) == 0 {
		return fmt.Errorf("no clips provided")
	}

	if len(clipPaths) != len(segmentDurations) {
		return fmt.Errorf("mismatch: %d clips but %d segment durations", len(clipPaths), len(segmentDurations))
	}

	fmt.Printf("\n🎬 Generating final video from %d animated clips...\n", len(clipPaths))
	for i, clipPath := range clipPaths {
		fmt.Printf("  Clip %d: %s (%.2fs)\n", i, filepath.Base(clipPath), segmentDurations[i])
	}

	// Create a concat file for ffmpeg
	concatFile := filepath.Join(filepath.Dir(outputPath), "concat_clips.txt")
	var concatContent strings.Builder

	// For each clip, we need to potentially loop or trim to match audio duration
	// AnimateDiff clips are typically 2 seconds (16 frames @ 8fps)
	tempClipsDir := filepath.Join(filepath.Dir(outputPath), "temp_clips")
	if err := os.MkdirAll(tempClipsDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp clips directory: %w", err)
	}
	defer os.RemoveAll(tempClipsDir)

	for i, clipPath := range clipPaths {
		// Get the clip duration
		clipDuration, err := g.getAudioDuration(clipPath)
		if err != nil {
			// If can't get duration, assume 2 seconds (typical AnimateDiff output)
			clipDuration = 2.0
		}

		targetDuration := segmentDurations[i]
		processedClip := filepath.Join(tempClipsDir, fmt.Sprintf("processed_%03d.mp4", i))

		// If clip is shorter than target, loop it; if longer, trim it
		var args []string
		if targetDuration > clipDuration {
			// Loop the clip
			loops := int(targetDuration/clipDuration) + 1
			args = []string{
				"-stream_loop", fmt.Sprintf("%d", loops),
				"-i", clipPath,
				"-t", fmt.Sprintf("%.2f", targetDuration),
				"-c", "copy",
				"-y",
				processedClip,
			}
		} else {
			// Trim the clip
			args = []string{
				"-i", clipPath,
				"-t", fmt.Sprintf("%.2f", targetDuration),
				"-c", "copy",
				"-y",
				processedClip,
			}
		}

		cmd := exec.Command("ffmpeg", args...)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to process clip %d: %w\nOutput: %s", i, err, string(output))
		}

		// Add to concat list
		absPath, _ := filepath.Abs(processedClip)
		concatContent.WriteString(fmt.Sprintf("file '%s'\n", absPath))
	}

	// Write concat file
	if err := os.WriteFile(concatFile, []byte(concatContent.String()), 0644); err != nil {
		return fmt.Errorf("failed to write concat file: %w", err)
	}
	defer os.Remove(concatFile)

	// Concatenate all clips and add audio
	args := []string{
		"-f", "concat",
		"-safe", "0",
		"-i", concatFile,
		"-i", audioPath,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-b:a", "192k",
		"-shortest",
		"-y",
		outputPath,
	}

	fmt.Printf("\nFFmpeg concat command: ffmpeg %s\n\n", strings.Join(args, " "))

	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg concat failed: %w\nOutput: %s", err, string(output))
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
