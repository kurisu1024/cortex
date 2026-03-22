package image

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Generator handles AI image generation
type Generator struct {
	modelID    string
	scriptPath string
}

// NewGenerator creates a new image generator
func NewGenerator(modelID string) *Generator {
	if modelID == "" {
		modelID = "runwayml/stable-diffusion-v1-5"
	}

	return &Generator{
		modelID:    modelID,
		scriptPath: filepath.Join("scripts", "sd_image_gen.py"),
	}
}

// GenerateImage generates an AI image from a text prompt
func (g *Generator) GenerateImage(prompt, outputPath string) error {
	// Run the Python script to generate the image
	cmd := exec.Command("python3", g.scriptPath, prompt, outputPath, g.modelID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("image generation failed: %w\nOutput: %s", err, string(output))
	}

	// Verify the image was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return fmt.Errorf("image file was not created: %s", outputPath)
	}

	return nil
}

// GenerateImagesForSegments generates one image per script segment
func (g *Generator) GenerateImagesForSegments(prompts []string, outputDir string) ([]string, error) {
	imagePaths := make([]string, 0, len(prompts))

	for i, prompt := range prompts {
		outputPath := filepath.Join(outputDir, fmt.Sprintf("image_%03d.png", i))

		fmt.Printf("  [%d/%d] Generating image for: %s...\n", i+1, len(prompts), truncate(prompt, 50))

		if err := g.GenerateImage(prompt, outputPath); err != nil {
			return nil, fmt.Errorf("failed to generate image %d: %w", i, err)
		}

		imagePaths = append(imagePaths, outputPath)
	}

	return imagePaths, nil
}

// truncate truncates a string to a maximum length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
