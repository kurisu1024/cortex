package script

import (
	"fmt"
	"strings"

	"github.com/kutidu2048/cortex/internal/models"
)

// Segment represents a section of the script
type Segment struct {
	Index      int
	Speaker    string
	Text       string
	VoicePath  string // Path to the voice model for this segment
}

// Script represents a generated script with segments
type Script struct {
	Title    string
	Segments []Segment
	RawText  string
}

// Generator handles script generation
type Generator struct {
	llm             *models.LLMClient
	availableVoices map[string]string // speaker name -> voice path
	speakers        []string          // list of speaker names for rotation
	duration        int               // target duration in minutes (default: 10)
}

// NewGenerator creates a new script generator
func NewGenerator(llm *models.LLMClient) *Generator {
	return &Generator{
		llm:             llm,
		availableVoices: make(map[string]string),
		speakers:        []string{},
		duration:        10, // Default 10 minutes
	}
}

// SetVoices sets multiple voices for different speakers
func (g *Generator) SetVoices(voices map[string]string) {
	g.availableVoices = voices
	g.speakers = make([]string, 0, len(voices))
	for speaker := range voices {
		g.speakers = append(g.speakers, speaker)
	}
}

// SetDuration sets the target duration for the script
func (g *Generator) SetDuration(minutes int) {
	if minutes > 0 {
		g.duration = minutes
	}
}

// Generate creates a script for the given topic
func (g *Generator) Generate(topic string) (*Script, error) {
	// Calculate approximate word count (150 words per minute of speech)
	targetWords := g.duration * 150

	prompt := fmt.Sprintf(`Create an engaging, informative script about: %s

The script should be:
- Conversational and engaging
- Well-structured with clear segments
- Approximately %d minutes when spoken (around %d words total)
- Educational but entertaining
- Comprehensive and detailed to meet the target length

Format the script with clear segments like this:
[SEGMENT 1]
Text for first segment...

[SEGMENT 2]
Text for second segment...

Begin:`, topic, g.duration, targetWords)

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
	// Calculate approximate word count (150 words per minute of speech)
	targetWords := g.duration * 150

	prompt := fmt.Sprintf(`Create an engaging, informative script about: %s

The script should be:
- Conversational and engaging
- Well-structured with clear segments
- Approximately %d minutes when spoken (around %d words total)
- Educational but entertaining
- Comprehensive and detailed to meet the target length

Format the script with clear segments like this:
[SEGMENT 1]
Text for first segment...

[SEGMENT 2]
Text for second segment...

Begin:`, topic, g.duration, targetWords)

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

	// Try to parse [SEGMENT N] format by splitting on [SEGMENT markers
	if strings.Contains(text, "[SEGMENT") {
		// Split by [SEGMENT markers
		parts := strings.Split(text, "[SEGMENT")

		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}

			// Extract segment number and content
			// Format: "1]\nContent here..."
			lines := strings.SplitN(part, "\n", 2)
			if len(lines) >= 2 {
				content := strings.TrimSpace(lines[1])
				if content != "" {
					segment := Segment{
						Index:   len(segments),
						Speaker: g.assignSpeaker(len(segments)),
						Text:    content,
					}
					segment.VoicePath = g.getVoicePathForSpeaker(segment.Speaker)
					segments = append(segments, segment)
				}
			}
		}
	}

	// If no segments found, split by paragraphs
	if len(segments) == 0 {
		paragraphs := strings.Split(text, "\n\n")
		for _, para := range paragraphs {
			para = strings.TrimSpace(para)
			if para != "" && len(para) > 20 {
				segment := Segment{
					Index:   len(segments),
					Speaker: g.assignSpeaker(len(segments)),
					Text:    para,
				}
				segment.VoicePath = g.getVoicePathForSpeaker(segment.Speaker)
				segments = append(segments, segment)
			}
		}
	}

	// If still no segments, create one big segment
	if len(segments) == 0 {
		segment := Segment{
			Index:   0,
			Speaker: g.assignSpeaker(0),
			Text:    strings.TrimSpace(text),
		}
		segment.VoicePath = g.getVoicePathForSpeaker(segment.Speaker)
		segments = append(segments, segment)
	}

	return segments
}

// assignSpeaker assigns a speaker to a segment (rotates through available speakers)
func (g *Generator) assignSpeaker(segmentIndex int) string {
	if len(g.speakers) == 0 {
		return "Narrator"
	}
	// Rotate through available speakers
	return g.speakers[segmentIndex%len(g.speakers)]
}

// getVoicePathForSpeaker returns the voice path for a given speaker
func (g *Generator) getVoicePathForSpeaker(speaker string) string {
	if voicePath, exists := g.availableVoices[speaker]; exists {
		return voicePath
	}
	return "" // Will use default voice
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
