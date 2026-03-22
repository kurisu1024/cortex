package job

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kutidu2048/cortex/internal/audio"
	"github.com/kutidu2048/cortex/internal/config"
	"github.com/kutidu2048/cortex/internal/models"
	"github.com/kutidu2048/cortex/internal/script"
	"github.com/kutidu2048/cortex/internal/ui"
	"github.com/kutidu2048/cortex/internal/video"
)

// Status represents job status
type Status string

const (
	StatusPending    Status = "pending"
	StatusRunning    Status = "running"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

// Job represents a generation job
type Job struct {
	ID             string
	Topic          string
	Status         Status
	Progress       *Progress
	OutputDir      string
	Voice          string
	Background     string
	HighVoicesOnly bool
	Duration       int // Target duration in minutes
	CreatedAt      time.Time
	StartedAt      *time.Time
	CompletedAt    *time.Time
	Error          error
}

// Manager handles job lifecycle
type Manager struct {
	modelManager *models.Manager
	jobs         map[string]*Job
	ui           *ui.Terminal
}

// NewManager creates a new job manager
func NewManager() *Manager {
	return &Manager{
		modelManager: models.NewManager(),
		jobs:         make(map[string]*Job),
		ui:           ui.NewTerminal(),
	}
}

// CreateJob creates a new job
func (m *Manager) CreateJob(topic, outputDir, voice, background string, highVoicesOnly bool, duration int) (string, error) {
	jobID := fmt.Sprintf("job_%d", time.Now().Unix())

	job := &Job{
		ID:             jobID,
		Topic:          topic,
		Status:         StatusPending,
		Progress:       NewProgress(),
		OutputDir:      outputDir,
		Voice:          voice,
		Background:     background,
		HighVoicesOnly: highVoicesOnly,
		Duration:       duration,
		CreatedAt:      time.Now(),
	}

	m.jobs[jobID] = job

	return jobID, nil
}

// RunJob executes a job
func (m *Manager) RunJob(jobID string) error {
	job, exists := m.jobs[jobID]
	if !exists {
		return fmt.Errorf("job not found: %s", jobID)
	}

	job.Status = StatusRunning
	now := time.Now()
	job.StartedAt = &now

	m.ui.ShowHeader()
	m.ui.ShowJobStart(job.Topic)

	// Create output directory
	if err := os.MkdirAll(job.OutputDir, 0755); err != nil {
		return m.failJob(job, fmt.Errorf("failed to create output directory: %w", err))
	}

	// Step 1: Generate script
	job.Progress.SetStep("Generating script", 1, 5)
	m.ui.ShowProgress(1, 5, "Generating script", "", job.Progress.StartTime)

	scriptGen := script.NewGenerator(m.modelManager.GetLLM())

	// Set target duration for the script
	scriptGen.SetDuration(job.Duration)

	// Load config to get multiple voices if available
	cfg, err := config.Load()
	if err == nil && len(cfg.Models.TTS.Voices) > 0 {
		voices := cfg.Models.TTS.Voices

		// Filter for high-quality voices only if requested
		if job.HighVoicesOnly {
			voices = m.filterHighQualityVoices(voices)
		}

		scriptGen.SetVoices(voices)
	}

	scr, err := scriptGen.Generate(job.Topic)
	if err != nil {
		return m.failJob(job, fmt.Errorf("script generation failed: %w", err))
	}

	// Step 2: Generate audio segments
	job.Progress.SetStep("Generating audio segments", 2, 5)
	m.ui.ShowProgress(2, 5, "Generating audio segments", "", job.Progress.StartTime)

	audioGen := audio.NewGenerator(m.modelManager.GetTTS())
	audioPaths, err := audioGen.GenerateFromScript(scr, job.OutputDir)
	if err != nil {
		return m.failJob(job, fmt.Errorf("audio generation failed: %w", err))
	}

	// Step 3: Combine audio
	job.Progress.SetStep("Combining audio segments", 3, 5)
	m.ui.ShowProgress(3, 5, "Combining audio segments", "", job.Progress.StartTime)

	audioCombiner := audio.NewCombiner()
	finalAudioPath := filepath.Join(job.OutputDir, "final_audio.wav")
	if err := audioCombiner.Combine(audioPaths, finalAudioPath); err != nil {
		return m.failJob(job, fmt.Errorf("audio combination failed: %w", err))
	}

	// Step 4: Generate video
	job.Progress.SetStep("Generating video", 4, 5)
	m.ui.ShowProgress(4, 5, "Generating video", "", job.Progress.StartTime)

	videoGen := video.NewGenerator()
	videoPath := filepath.Join(job.OutputDir, "output.mp4")
	if err := videoGen.GenerateFromAudio(finalAudioPath, videoPath, job.Background, true); err != nil {
		return m.failJob(job, fmt.Errorf("video generation failed: %w", err))
	}

	// Step 5: Complete
	job.Progress.SetStep("Finalizing", 5, 5)
	m.ui.ShowProgress(5, 5, "Finalizing", "", job.Progress.StartTime)

	job.Status = StatusCompleted
	completedAt := time.Now()
	job.CompletedAt = &completedAt

	m.ui.ShowSuccess(videoPath)

	return nil
}

// failJob marks a job as failed
func (m *Manager) failJob(job *Job, err error) error {
	job.Status = StatusFailed
	job.Error = err
	completedAt := time.Now()
	job.CompletedAt = &completedAt

	m.ui.ShowError(err.Error())

	return err
}

// GetJob retrieves a job by ID
func (m *Manager) GetJob(jobID string) (*Job, error) {
	job, exists := m.jobs[jobID]
	if !exists {
		return nil, fmt.Errorf("job not found: %s", jobID)
	}

	return job, nil
}

// ListJobs returns all jobs
func (m *Manager) ListJobs() []*Job {
	jobs := make([]*Job, 0, len(m.jobs))
	for _, job := range m.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// filterHighQualityVoices filters voices to only include high-quality ones
func (m *Manager) filterHighQualityVoices(voices map[string]string) map[string]string {
	filtered := make(map[string]string)

	for speaker, voicePath := range voices {
		// Check if voice path contains "-high" indicating high quality
		if strings.Contains(voicePath, "-high") {
			filtered[speaker] = voicePath
		}
	}

	// If no high-quality voices found, return all voices as fallback
	if len(filtered) == 0 {
		return voices
	}

	return filtered
}
