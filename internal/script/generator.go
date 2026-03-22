package script

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/topher/cortex/internal/models"
)

// Segment represents a section of the script
type Segment struct {
	Index   int
	Speaker string
	Text    string
}

// Script represents a generated script with segments
type Script struct {
	Title    string
	Segments []Segment
	RawText  string
}

// Generator handles script generation
type Generator struct {
	llm *models.LLMClient
}

// NewGenerator creates a new script generator
func NewGenerator(llm *models.LLMClient) *Generator {
	return &Generator{llm: llm}
}

// Generate creates a script for the given topic
func (g *Generator) Generate(topic string) (*Script, error) {
	prompt := fmt.Sprintf(`Create an engaging, informative script about: %s

The script should be:
- Conversational and engaging
- Well-structured with clear segments
- About 2-3 minutes when spoken
- Educational but entertaining

Format the script with clear segments like this:
[SEGMENT 1]
Text for first segment...

[SEGMENT 2]
Text for second segment...

Begin:`, topic)

	fmt.Println("🧠 Generating script with AI...")

	rawScript, err := g.llm.Generate(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate script: %w", err)
	}

	script := &Script{
		Title:   topic,
		RawText: rawScript,
	}

	// Parse script into segments
	segments := g.parseSegments(rawScript)
	script.Segments = segments

	fmt.Printf("✅ Generated script with %d segments\n", len(segments))

	return script, nil
}

// GenerateStream creates a script with streaming output
func (g *Generator) GenerateStream(topic string, callback func(string)) (*Script, error) {
	prompt := fmt.Sprintf(`Create an engaging, informative script about: %s

The script should be:
- Conversational and engaging
- Well-structured with clear segments
- About 2-3 minutes when spoken
- Educational but entertaining

Format the script with clear segments like this:
[SEGMENT 1]
Text for first segment...

[SEGMENT 2]
Text for second segment...

Begin:`, topic)

	fmt.Println("🧠 Generating script with AI...")

	var fullText strings.Builder
	err := g.llm.GenerateStream(prompt, func(chunk string) {
		fullText.WriteString(chunk)
		if callback != nil {
			callback(chunk)
		}
	})

	if err != nil {
		return nil, fmt.Errorf("failed to generate script: %w", err)
	}

	rawScript := fullText.String()
	script := &Script{
		Title:   topic,
		RawText: rawScript,
	}

	segments := g.parseSegments(rawScript)
	script.Segments = segments

	fmt.Printf("\n✅ Generated script with %d segments\n", len(segments))

	return script, nil
}

// parseSegments parses raw script text into segments
func (g *Generator) parseSegments(text string) []Segment {
	var segments []Segment

	// Try to parse [SEGMENT N] format
	segmentRegex := regexp.MustCompile(`\[SEGMENT\s+(\d+)\]\s*\n([\s\S]*?)(?=\[SEGMENT|\z)`)
	matches := segmentRegex.FindAllStringSubmatch(text, -1)

	if len(matches) > 0 {
		for _, match := range matches {
			if len(match) >= 3 {
				segments = append(segments, Segment{
					Index:   len(segments),
					Speaker: "Narrator",
					Text:    strings.TrimSpace(match[2]),
				})
			}
		}
	} else {
		// If no segments found, split by paragraphs
		paragraphs := strings.Split(text, "\n\n")
		for _, para := range paragraphs {
			para = strings.TrimSpace(para)
			if para != "" && len(para) > 20 {
				segments = append(segments, Segment{
					Index:   len(segments),
					Speaker: "Narrator",
					Text:    para,
				})
			}
		}
	}

	// If still no segments, create one big segment
	if len(segments) == 0 {
		segments = append(segments, Segment{
			Index:   0,
			Speaker: "Narrator",
			Text:    strings.TrimSpace(text),
		})
	}

	return segments
}

// SaveToFile saves the script to a text file
func (s *Script) SaveToFile(filepath string) error {
	var content strings.Builder

	content.WriteString(fmt.Sprintf("Title: %s\n\n", s.Title))
	content.WriteString("=" + strings.Repeat("=", len(s.Title)+7) + "\n\n")

	for _, seg := range s.Segments {
		content.WriteString(fmt.Sprintf("[SEGMENT %d]\n", seg.Index+1))
		content.WriteString(seg.Text + "\n\n")
	}

	// Write to file
	return nil // Implementation would use os.WriteFile
}
