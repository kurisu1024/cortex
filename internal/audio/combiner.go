package audio

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Combiner handles combining multiple audio files
type Combiner struct{}

// NewCombiner creates a new audio combiner
func NewCombiner() *Combiner {
	return &Combiner{}
}

// Combine merges multiple audio files into one
func (c *Combiner) Combine(audioPaths []string, outputPath string) error {
	if len(audioPaths) == 0 {
		return fmt.Errorf("no audio files to combine")
	}

	// Check if ffmpeg is installed
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg is not installed. Please install it first")
	}

	fmt.Println("\n🔊 Combining audio segments...")

	// If only one file, just copy it
	if len(audioPaths) == 1 {
		return c.copyFile(audioPaths[0], outputPath)
	}

	// Create a concat file for ffmpeg
	concatFile := filepath.Join(filepath.Dir(outputPath), "concat_list.txt")
	if err := c.createConcatFile(audioPaths, concatFile); err != nil {
		return err
	}
	defer os.Remove(concatFile)

	// Use ffmpeg to concatenate audio files
	cmd := exec.Command("ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", concatFile,
		"-c", "copy",
		"-y", // Overwrite output file
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("✅ Combined audio saved to: %s\n", outputPath)

	return nil
}

// createConcatFile creates a concat file for ffmpeg
func (c *Combiner) createConcatFile(audioPaths []string, outputPath string) error {
	var content strings.Builder

	for _, path := range audioPaths {
		// Convert to absolute path
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		content.WriteString(fmt.Sprintf("file '%s'\n", absPath))
	}

	return os.WriteFile(outputPath, []byte(content.String()), 0644)
}

// copyFile copies a file from src to dst
func (c *Combiner) copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0644)
}

// GetDuration returns the duration of an audio file in seconds
func (c *Combiner) GetDuration(audioPath string) (float64, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		audioPath,
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe failed: %w", err)
	}

	var duration float64
	_, err = fmt.Sscanf(strings.TrimSpace(string(output)), "%f", &duration)
	if err != nil {
		return 0, err
	}

	return duration, nil
}
