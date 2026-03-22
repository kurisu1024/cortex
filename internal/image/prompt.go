package image

import (
	"fmt"
	"strings"

	"github.com/kutidu2048/cortex/internal/script"
)

// PromptGenerator generates image prompts from script segments
type PromptGenerator struct {
	style string // Art style suffix to add to prompts
}

// NewPromptGenerator creates a new prompt generator
func NewPromptGenerator(style string) *PromptGenerator {
	if style == "" {
		style = "cinematic lighting, highly detailed, professional photography, 8k resolution"
	}

	return &PromptGenerator{
		style: style,
	}
}

// GeneratePrompts creates image prompts from script segments
func (p *PromptGenerator) GeneratePrompts(segments []script.Segment) []string {
	prompts := make([]string, 0)
	seenScenes := make(map[string]bool)

	for _, segment := range segments {
		// Use scene action if available, otherwise generate from text
		var prompt string
		if segment.SceneAction != "" {
			// Use the scene description directly - it's already a visual description
			prompt = fmt.Sprintf("%s, %s", segment.SceneAction, p.style)

			// Only add each unique scene once (avoid duplicates from multiple dialogue lines in same scene)
			if !seenScenes[segment.SceneAction] {
				prompts = append(prompts, prompt)
				seenScenes[segment.SceneAction] = true
			}
		} else {
			// Fallback to keyword-based generation for old format
			prompt = p.generatePromptFromText(segment.Text)
			prompts = append(prompts, prompt)
		}
	}

	return prompts
}

// generatePromptFromText extracts key concepts and creates an image prompt
func (p *PromptGenerator) generatePromptFromText(text string) string {
	// Simple keyword extraction and prompt generation
	// This is a basic implementation - could be enhanced with NLP

	text = strings.ToLower(text)

	// Determine subject matter based on keywords
	var subject string

	// Programming/Tech topics
	if containsAny(text, []string{"python", "programming", "code", "software", "developer"}) {
		subject = "modern programming workspace with code on screens"
	} else if containsAny(text, []string{"go", "golang", "concurrency"}) {
		subject = "futuristic technology laboratory with holographic displays"
	} else if containsAny(text, []string{"ai", "artificial intelligence", "machine learning", "neural"}) {
		subject = "futuristic AI neural network visualization with glowing nodes"
	} else if containsAny(text, []string{"data", "database", "storage"}) {
		subject = "abstract data visualization with colorful flowing information"
	} else if containsAny(text, []string{"function", "variable", "syntax"}) {
		subject = "clean modern code editor with colorful syntax highlighting"
	} else if containsAny(text, []string{"web", "internet", "browser", "http"}) {
		subject = "modern web development environment with multiple screens"
	} else if containsAny(text, []string{"server", "cloud", "infrastructure"}) {
		subject = "futuristic server room with glowing network connections"
	} else if containsAny(text, []string{"security", "encryption", "password"}) {
		subject = "cyber security concept with digital locks and shields"
	} else if containsAny(text, []string{"mobile", "app", "ios", "android"}) {
		subject = "sleek mobile app interface design with modern UI elements"
	} else if containsAny(text, []string{"design", "ui", "interface", "user experience"}) {
		subject = "beautiful user interface design with modern aesthetic"
	// Science topics
	} else if containsAny(text, []string{"science", "experiment", "research"}) {
		subject = "modern scientific laboratory with advanced equipment"
	} else if containsAny(text, []string{"space", "astronomy", "cosmos", "universe"}) {
		subject = "beautiful space scene with stars and nebulae"
	} else if containsAny(text, []string{"biology", "cell", "dna", "genetic"}) {
		subject = "microscopic view of cells and DNA structures"
	// Business topics
	} else if containsAny(text, []string{"business", "corporate", "company", "startup"}) {
		subject = "modern corporate office with glass and steel architecture"
	} else if containsAny(text, []string{"finance", "money", "investment", "stock"}) {
		subject = "financial charts and graphs with upward trends"
	} else if containsAny(text, []string{"marketing", "advertising", "brand"}) {
		subject = "creative marketing brainstorming session with mood boards"
	// Education topics
	} else if containsAny(text, []string{"learn", "education", "study", "teach"}) {
		subject = "modern educational environment with interactive displays"
	} else if containsAny(text, []string{"book", "reading", "literature"}) {
		subject = "cozy library with books and warm lighting"
	// General/Abstract
	} else {
		subject = "abstract concept visualization with flowing colors and light"
	}

	// Combine subject with style
	return fmt.Sprintf("%s, %s", subject, p.style)
}

// containsAny checks if text contains any of the keywords
func containsAny(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}
