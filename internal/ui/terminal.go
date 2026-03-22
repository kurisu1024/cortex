package ui

import (
	"fmt"
	"strings"
	"time"
)

// Terminal handles terminal UI output with hacker-style aesthetics
type Terminal struct {
	width int
}

// NewTerminal creates a new terminal UI
func NewTerminal() *Terminal {
	return &Terminal{
		width: 80,
	}
}

// Color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	Bold   = "\033[1m"
	Dim    = "\033[2m"
)

// ShowHeader displays the Cortex ASCII header
func (t *Terminal) ShowHeader() {
	banner := `
    ╔═══════════════════════════════════════════════════════════╗
    ║                                                           ║
    ║   ██████╗ ██████╗ ██████╗ ████████╗███████╗██╗  ██╗    ║
    ║  ██╔════╝██╔═══██╗██╔══██╗╚══██╔══╝██╔════╝╚██╗██╔╝    ║
    ║  ██║     ██║   ██║██████╔╝   ██║   █████╗   ╚███╔╝     ║
    ║  ██║     ██║   ██║██╔══██╗   ██║   ██╔══╝   ██╔██╗     ║
    ║  ╚██████╗╚██████╔╝██║  ██║   ██║   ███████╗██╔╝ ██╗    ║
    ║   ╚═════╝ ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚══════╝╚═╝  ╚═╝    ║
    ║                                                           ║
    ║            AI-Powered Script to Video Generator          ║
    ║                     [LOCAL MODELS]                       ║
    ╚═══════════════════════════════════════════════════════════╝
	`
	fmt.Printf("%s%s%s\n", Cyan, banner, Reset)
}

// ShowJobStart displays job initiation message
func (t *Terminal) ShowJobStart(topic string) {
	fmt.Printf("\n%s%s╔══════════════════════════════════════════════════════════════╗%s\n", Bold, Green, Reset)
	fmt.Printf("%s%s║  INITIALIZING JOB                                             ║%s\n", Bold, Green, Reset)
	fmt.Printf("%s%s╠══════════════════════════════════════════════════════════════╣%s\n", Bold, Green, Reset)
	fmt.Printf("%s%s║  Topic: %-53s ║%s\n", Bold, Green, topic, Reset)
	fmt.Printf("%s%s╚══════════════════════════════════════════════════════════════╝%s\n\n", Bold, Green, Reset)
}

// Progress represents job progress data
type Progress struct {
	CurrentStep int
	TotalSteps  int
	StepName    string
	StartTime   time.Time
	Message     string
}

// Percentage returns the progress percentage
func (p *Progress) Percentage() float64 {
	if p.TotalSteps == 0 {
		return 0
	}
	return float64(p.CurrentStep) / float64(p.TotalSteps) * 100
}

// Elapsed returns time elapsed since start
func (p *Progress) Elapsed() time.Duration {
	return time.Since(p.StartTime)
}

// ShowProgress displays job progress
func (t *Terminal) ShowProgress(currentStep, totalSteps int, stepName, message string, startTime time.Time) {
	percentage := float64(currentStep) / float64(totalSteps) * 100
	bar := t.createProgressBar(int(percentage), 50)

	fmt.Printf("%s▶ [%d/%d] %s%s\n", Cyan, currentStep, totalSteps, stepName, Reset)
	fmt.Printf("  %s%s%s %.1f%%\n", Green, bar, Reset, percentage)

	if message != "" {
		fmt.Printf("  %s%s%s\n", Dim, message, Reset)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("  %sElapsed: %s%s\n\n", Gray, elapsed.Round(time.Second), Reset)
}

// createProgressBar creates a progress bar string
func (t *Terminal) createProgressBar(percentage, width int) string {
	filled := (percentage * width) / 100
	empty := width - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	return fmt.Sprintf("[%s]", bar)
}

// ShowSuccess displays success message
func (t *Terminal) ShowSuccess(outputPath string) {
	fmt.Printf("\n%s%s╔══════════════════════════════════════════════════════════════╗%s\n", Bold, Green, Reset)
	fmt.Printf("%s%s║  ✓ JOB COMPLETED SUCCESSFULLY                                 ║%s\n", Bold, Green, Reset)
	fmt.Printf("%s%s╠══════════════════════════════════════════════════════════════╣%s\n", Bold, Green, Reset)
	fmt.Printf("%s%s║  Output: %-52s ║%s\n", Bold, Green, outputPath, Reset)
	fmt.Printf("%s%s╚══════════════════════════════════════════════════════════════╝%s\n\n", Bold, Green, Reset)
}

// ShowError displays error message
func (t *Terminal) ShowError(message string) {
	boxWidth := 62 // Total width between ║ characters
	fmt.Printf("\n%s%s╔══════════════════════════════════════════════════════════════╗%s\n", Bold, Red, Reset)
	fmt.Printf("%s%s║  ✗ ERROR                                                      ║%s\n", Bold, Red, Reset)
	fmt.Printf("%s%s╠══════════════════════════════════════════════════════════════╣%s\n", Bold, Red, Reset)

	// Split long error messages
	words := strings.Split(message, " ")
	line := ""
	maxLineLen := boxWidth - 4 // Subtract 4 for "║  " and " ║"

	for _, word := range words {
		if len(line)+len(word)+1 > maxLineLen {
			// Print current line with proper padding
			padding := maxLineLen - len(line)
			fmt.Printf("%s%s║  %s%s ║%s\n", Bold, Red, line, strings.Repeat(" ", padding), Reset)
			line = word
		} else {
			if line != "" {
				line += " "
			}
			line += word
		}
	}

	// Print last line
	if line != "" {
		padding := maxLineLen - len(line)
		fmt.Printf("%s%s║  %s%s ║%s\n", Bold, Red, line, strings.Repeat(" ", padding), Reset)
	}

	fmt.Printf("%s%s╚══════════════════════════════════════════════════════════════╝%s\n\n", Bold, Red, Reset)
}

// ShowMatrix displays matrix-style loading animation (simplified)
func (t *Terminal) ShowMatrix(text string) {
	fmt.Printf("%s%s%s%s\n", Green, Dim, text, Reset)
}

// ShowInfo displays info message
func (t *Terminal) ShowInfo(message string) {
	fmt.Printf("%s▶ %s%s\n", Cyan, message, Reset)
}

// ShowWarning displays warning message
func (t *Terminal) ShowWarning(message string) {
	fmt.Printf("%s⚠ %s%s\n", Yellow, message, Reset)
}

// Cortex robot animation frames
var cortexFrames = []string{
	// Frame 1: Arms up
	`    🤖
   \║║/
    ║║
   /  \   `,
	// Frame 2: Arms out
	`    🤖
   -║║-
    ║║
   /  \   `,
	// Frame 3: Arms down
	`    🤖
   /║║\
    ║║
   \  /   `,
	// Frame 4: Arms wave
	`    🤖
   /║║~
    ║║
   /  \   `,
}

// ShowCortexRobot displays the dancing Cortex robot
func (t *Terminal) ShowCortexRobot(frame int) string {
	frameIndex := frame % len(cortexFrames)
	lines := strings.Split(cortexFrames[frameIndex], "\n")

	var result strings.Builder
	for _, line := range lines {
		result.WriteString(fmt.Sprintf("%s%s%s\n", Cyan, line, Reset))
	}
	return result.String()
}

// ClearLines clears N lines up from current position
func (t *Terminal) ClearLines(n int) {
	for i := 0; i < n; i++ {
		fmt.Print("\033[1A\033[2K") // Move up and clear line
	}
}

// ShowCortexWithMessage displays Cortex robot next to a message
func (t *Terminal) ShowCortexWithMessage(message string, frame int) {
	robot := strings.Split(cortexFrames[frame%len(cortexFrames)], "\n")

	// Print message with robot on the right
	fmt.Printf("%s%-50s%s%s%s\n", Cyan, message, Reset, Cyan, robot[0])
	for i := 1; i < len(robot); i++ {
		fmt.Printf("%s%-50s%s%s%s\n", "", "", Reset, Cyan, robot[i])
	}
	fmt.Print(Reset)
}

// ShowModelStatus displays model status in a styled box
func (t *Terminal) ShowModelStatus(llmStatus, ttsStatus, llmHost, llmModel, ttsVoice string) {
	fmt.Printf("\n%s╔══════════════════════════════════════════════════════════════╗%s\n", Cyan, Reset)
	fmt.Printf("%s║  MODEL STATUS                                                 ║%s\n", Cyan, Reset)
	fmt.Printf("%s╠══════════════════════════════════════════════════════════════╣%s\n", Cyan, Reset)
	fmt.Printf("%s║  LLM (Ollama):  %-45s ║%s\n", Cyan, llmStatus, Reset)
	fmt.Printf("%s║    Host:  %-51s ║%s\n", Cyan, llmHost, Reset)
	fmt.Printf("%s║    Model: %-51s ║%s\n", Cyan, llmModel, Reset)
	fmt.Printf("%s║                                                               ║%s\n", Cyan, Reset)
	fmt.Printf("%s║  TTS (Piper):   %-45s ║%s\n", Cyan, ttsStatus, Reset)
	fmt.Printf("%s║    Voice: %-51s ║%s\n", Cyan, ttsVoice, Reset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════════╝%s\n\n", Cyan, Reset)
}
