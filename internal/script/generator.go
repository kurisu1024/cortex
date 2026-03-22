package script

import (
	"fmt"
	"strings"

	"github.com/kutidu2048/cortex/internal/models"
)

// Segment represents a section of the script
type Segment struct {
	Index         int
	Speaker       string
	Text          string    // The spoken dialogue (without speaker name)
	SceneAction   string    // Visual description for video generation
	VoicePath     string    // Path to the voice model for this segment
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

	prompt := fmt.Sprintf(`Create an engaging, animated video script about: %s

The script should be:
- Conversational with character dialogue
- Well-structured with visual scene descriptions
- Approximately %d minutes when spoken (around %d words total)
- Educational but entertaining with personality
- Include vivid scene descriptions for animation

Format EXACTLY like this example:

[SCENE 1: Wide shot of futuristic Mars colony with orange sky and domed habitats]
NARRATOR: Welcome to the year 2045, where humanity has taken its greatest leap.

[SCENE 2: Inside high-tech research lab with holographic displays and scientists working]
SCIENTIST: The data we're seeing is unprecedented. Life may exist beneath the Martian surface.
NARRATOR: This discovery could change everything we know about the universe.

[SCENE 3: Close-up of rover on Mars surface, red dust swirling around it]
ENGINEER: The rover's sensors detected organic compounds in this crater.

IMPORTANT RULES:
- Start each scene with [SCENE N: visual description of what we see]
- Use CHARACTER: before each line of dialogue
- Make scenes visually descriptive and cinematic
- Characters should have distinct personalities
- Scene descriptions should guide animation/visuals

Create approximately %d scenes for a %d minute video.

Begin:`, topic, targetWords/50, targetWords, g.duration)

	fmt.Println("🧠 Generating script with AI...")

	rawScript, err := g.llm.Generate(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate script: %w", err)
	}

	script := &Script{
		Title:   topic,
		RawText: rawScript,
	}

	// DEBUG: Print raw script to see what LLM generated
	fmt.Printf("\n--- RAW SCRIPT ---\n%s\n--- END RAW SCRIPT ---\n\n", rawScript)

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

	prompt := fmt.Sprintf(`Create an engaging, animated video script about: %s

The script should be:
- Conversational with character dialogue
- Well-structured with visual scene descriptions
- Approximately %d minutes when spoken (around %d words total)
- Educational but entertaining with personality
- Include vivid scene descriptions for animation

Format EXACTLY like this example:

[SCENE 1: Wide shot of futuristic Mars colony with orange sky and domed habitats]
NARRATOR: Welcome to the year 2045, where humanity has taken its greatest leap.

[SCENE 2: Inside high-tech research lab with holographic displays and scientists working]
SCIENTIST: The data we're seeing is unprecedented. Life may exist beneath the Martian surface.
NARRATOR: This discovery could change everything we know about the universe.

[SCENE 3: Close-up of rover on Mars surface, red dust swirling around it]
ENGINEER: The rover's sensors detected organic compounds in this crater.

IMPORTANT RULES:
- Start each scene with [SCENE N: visual description of what we see]
- Use CHARACTER: before each line of dialogue
- Make scenes visually descriptive and cinematic
- Characters should have distinct personalities
- Scene descriptions should guide animation/visuals

Create approximately %d scenes for a %d minute video.

Begin:`, topic, targetWords/50, targetWords, g.duration)

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

// parseSegments parses raw script text into segments with scene actions and dialogue
func (g *Generator) parseSegments(text string) []Segment {
	var segments []Segment

	// Try to parse [SCENE N: description] format
	if strings.Contains(text, "[SCENE") {
		// Split by [SCENE markers
		parts := strings.Split(text, "[SCENE")

		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}

			// Extract scene description and dialogue
			// Format: "1: Scene description]\nCHARACTER: dialogue\n..."
			lines := strings.Split(part, "\n")
			if len(lines) < 2 {
				continue
			}

			// Parse scene header: "1: Scene description]"
			sceneHeader := lines[0]
			sceneAction := ""
			if idx := strings.Index(sceneHeader, ":"); idx != -1 {
				// Extract text between ":" and "]"
				actionText := sceneHeader[idx+1:]
				if endIdx := strings.Index(actionText, "]"); endIdx != -1 {
					sceneAction = strings.TrimSpace(actionText[:endIdx])
				}
			}

			// Parse dialogue lines
			for _, line := range lines[1:] {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}

				// Parse "CHARACTER: dialogue" format
				if idx := strings.Index(line, ":"); idx != -1 {
					speaker := strings.TrimSpace(line[:idx])
					dialogue := strings.TrimSpace(line[idx+1:])

					if dialogue != "" {
						// Clean up speaker name - remove parenthetical notes and extra info
						// e.g. "NARRATOR (in an adventurous tone)" -> "narrator"
						// e.g. "DR. PATEL, Lead Researcher" -> "dr. patel"
						speaker = cleanSpeakerName(speaker)

						segment := Segment{
							Index:       len(segments),
							Speaker:     speaker,
							Text:        dialogue, // Just the dialogue, no speaker name
							SceneAction: sceneAction,
							VoicePath:   g.getVoicePathForSpeaker(speaker),
						}
						segments = append(segments, segment)
					}
				}
			}
		}
	}

	// Fallback: If no SCENE format found, try old SEGMENT format
	if len(segments) == 0 && strings.Contains(text, "[SEGMENT") {
		parts := strings.Split(text, "[SEGMENT")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			lines := strings.SplitN(part, "\n", 2)
			if len(lines) >= 2 {
				content := strings.TrimSpace(lines[1])
				if content != "" {
					segment := Segment{
						Index:       len(segments),
						Speaker:     g.assignSpeaker(len(segments)),
						Text:        content,
						SceneAction: "",
						VoicePath:   g.getVoicePathForSpeaker(g.assignSpeaker(len(segments))),
					}
					segments = append(segments, segment)
				}
			}
		}
	}

	// Last fallback: split by paragraphs
	if len(segments) == 0 {
		paragraphs := strings.Split(text, "\n\n")
		for _, para := range paragraphs {
			para = strings.TrimSpace(para)
			if para != "" && len(para) > 20 {
				segment := Segment{
					Index:       len(segments),
					Speaker:     g.assignSpeaker(len(segments)),
					Text:        para,
					SceneAction: "",
					VoicePath:   g.getVoicePathForSpeaker(g.assignSpeaker(len(segments))),
				}
				segments = append(segments, segment)
			}
		}
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

// cleanSpeakerName removes parenthetical notes and extra descriptions from speaker names
func cleanSpeakerName(speaker string) string {
	// Remove parenthetical notes: "NARRATOR (in an adventurous tone)" -> "NARRATOR"
	if idx := strings.Index(speaker, "("); idx != -1 {
		speaker = strings.TrimSpace(speaker[:idx])
	}

	// Remove comma and everything after: "DR. PATEL, Lead Researcher" -> "DR. PATEL"
	if idx := strings.Index(speaker, ","); idx != -1 {
		speaker = strings.TrimSpace(speaker[:idx])
	}

	// Convert to lowercase
	return strings.ToLower(speaker)
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
