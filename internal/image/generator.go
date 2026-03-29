package image

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Generator handles AI image generation and animated video generation
type Generator struct {
	modelID         string
	scriptPath      string
	animatedScript  string
	useAnimation    bool
	animationFrames int
}

// NewGenerator creates a new image generator
func NewGenerator(modelID string) *Generator {
	if modelID == "" {
		modelID = "stabilityai/sdxl-turbo"
	}

	return &Generator{
		modelID:         modelID,
		scriptPath:      filepath.Join("scripts", "sd_image_gen.py"),
		animatedScript:  filepath.Join("scripts", "animatediff_gen.py"),
		useAnimation:    false,
		animationFrames: 16,
	}
}

// SetAnimationMode enables or disables animated video generation
func (g *Generator) SetAnimationMode(enabled bool, frames int) {
	g.useAnimation = enabled
	if frames > 0 {
		g.animationFrames = frames
	}
}

// GenerateImage generates an AI image from a text prompt
func (g *Generator) GenerateImage(prompt, outputPath string) error {
	// Run the Python script to generate the image
	cmd := exec.Command("python3", g.scriptPath, prompt, outputPath, g.modelID)

	// Pipe stderr to see progress output from Python script
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("image generation failed: %w", err)
	}

	// Verify the image was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return fmt.Errorf("image file was not created: %s", outputPath)
	}

	return nil
}

// GenerateAnimatedVideo generates an animated video clip from a text prompt
func (g *Generator) GenerateAnimatedVideo(prompt, outputPath string) error {
	// Run the AnimateDiff Python script
	cmd := exec.Command("python3", g.animatedScript, prompt, outputPath, fmt.Sprintf("%d", g.animationFrames))

	// Pipe stderr to see progress output
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("animated video generation failed: %w", err)
	}

	// Verify the video was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return fmt.Errorf("video file was not created: %s", outputPath)
	}

	return nil
}

// GenerateImagesForSegments generates one image per script segment
func (g *Generator) GenerateImagesForSegments(prompts []string, outputDir string) ([]string, error) {
	paths := make([]string, 0, len(prompts))

	// ANSI color codes
	const (
		Green = "\033[1;32m"
		Reset = "\033[0m"
	)

	if g.useAnimation {
		fmt.Println("🎬 Using AnimateDiff for animated video clips")
		for i, prompt := range prompts {
			outputPath := filepath.Join(outputDir, fmt.Sprintf("clip_%03d.mp4", i))

			fmt.Printf("%s  [%d/%d] Generating animated clip for: %s...%s\n", Green, i+1, len(prompts), truncate(prompt, 50), Reset)

			if err := g.GenerateAnimatedVideo(prompt, outputPath); err != nil {
				return nil, fmt.Errorf("failed to generate animated video %d: %w", i, err)
			}

			paths = append(paths, outputPath)
		}
	} else {
		for i, prompt := range prompts {
			outputPath := filepath.Join(outputDir, fmt.Sprintf("image_%03d.png", i))

			fmt.Printf("%s  [%d/%d] Generating image for: %s...%s\n", Green, i+1, len(prompts), truncate(prompt, 50), Reset)

			if err := g.GenerateImage(prompt, outputPath); err != nil {
				return nil, fmt.Errorf("failed to generate image %d: %w", i, err)
			}

			paths = append(paths, outputPath)
		}
	}

	return paths, nil
}

// truncate truncates a string to a maximum length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
