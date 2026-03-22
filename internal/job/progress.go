package job

import (
	"fmt"
	"time"
)

// Progress tracks job progress
type Progress struct {
	CurrentStep int
	TotalSteps  int
	StepName    string
	StartTime   time.Time
	Message     string
}

// NewProgress creates a new progress tracker
func NewProgress() *Progress {
	return &Progress{
		StartTime: time.Now(),
	}
}

// SetStep sets the current step
func (p *Progress) SetStep(name string, current, total int) {
	p.StepName = name
	p.CurrentStep = current
	p.TotalSteps = total
}

// SetMessage sets a progress message
func (p *Progress) SetMessage(msg string) {
	p.Message = msg
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

// String returns a string representation of progress
func (p *Progress) String() string {
	return fmt.Sprintf("[%d/%d] %s (%.1f%%)",
		p.CurrentStep,
		p.TotalSteps,
		p.StepName,
		p.Percentage(),
	)
}
