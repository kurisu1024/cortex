package job

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kutidu2048/cortex/internal/audio"
	"github.com/kutidu2048/cortex/internal/config"
	"github.com/kutidu2048/cortex/internal/image"
	"github.com/kutidu2048/cortex/internal/models"
	"github.com/kutidu2048/cortex/internal/script"
	"github.com/kutidu2048/cortex/internal/ui"
	"github.com/kutidu2048/cortex/internal/video"
)

// Status represents job status
type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
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

	// Show dancing Cortex while generating script
	stopCortex := make(chan bool)
	go m.animateCortex("🧠 Generating script with AI...", stopCortex)

	scr, err := scriptGen.Generate(job.Topic)

	// Stop animation
	stopCortex <- true
	time.Sleep(100 * time.Millisecond) // Give goroutine time to stop
	fmt.Print("\033[4A\033[J")         // Move up 4 lines and clear to end

	if err != nil {
		return m.failJob(job, fmt.Errorf("script generation failed: %w", err))
	}

	// Step 2: Generate audio segments
	job.Progress.SetStep("Generating audio segments", 2, 5)
	m.ui.ShowProgress(2, 5, "Generating audio segments", "", job.Progress.StartTime)

	audioGen := audio.NewGenerator(m.modelManager.GetTTS())
	audioSegments, err := audioGen.GenerateFromScript(scr, job.OutputDir)
	if err != nil {
		return m.failJob(job, fmt.Errorf("audio generation failed: %w", err))
	}

	// Extract paths for combining
	audioPaths := make([]string, len(audioSegments))
	segmentDurations := make([]float64, len(audioSegments))
	for i, seg := range audioSegments {
		audioPaths[i] = seg.Path
		segmentDurations[i] = seg.Duration
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

	// Check if we need to generate AI images
	if job.Background == "ai-generated" {
		// Load config for image/animation settings FIRST
		cfg, err := config.Load()
		if err != nil {
			return m.failJob(job, fmt.Errorf("failed to load config: %w", err))
		}

		// Check if we should use animation
		useAnimation := cfg.Output.Video.Animated
		animationFrames := cfg.Output.Video.AnimationFrames
		if animationFrames == 0 {
			animationFrames = 16 // Default 16 frames = 2 seconds
		}

		fmt.Printf("\n[DEBUG] Config loaded: animated=%v, frames=%d\n", useAnimation, animationFrames)

		// Generate prompts from script segments
		promptGen := image.NewPromptGenerator("cinematic, high quality, 4k, detailed")
		prompts := promptGen.GeneratePrompts(scr.Segments)

		if len(prompts) == 0 {
			return m.failJob(job, fmt.Errorf("failed to generate image prompts from script"))
		}

		modelID := cfg.Models.Image.ModelID
		if modelID == "" {
			modelID = "stabilityai/sdxl-turbo"
		}
		imageGen := image.NewGenerator(modelID)

		// Set animation mode BEFORE generating
		imageGen.SetAnimationMode(useAnimation, animationFrames)

		if useAnimation {
			fmt.Printf("\n🎬 Generating %d animated video clips...\n", len(prompts))
		} else {
			fmt.Printf("\n🎨 Generating %d AI images...\n", len(prompts))
		}
		// Create directory for images or clips
		mediaDir := filepath.Join(job.OutputDir, "images")
		if useAnimation {
			mediaDir = filepath.Join(job.OutputDir, "clips")
		}
		if err := os.MkdirAll(mediaDir, 0755); err != nil {
			return m.failJob(job, fmt.Errorf("failed to create media directory: %w", err))
		}

		mediaPaths, err := imageGen.GenerateImagesForSegments(prompts, mediaDir)
		if err != nil {
			return m.failJob(job, fmt.Errorf("media generation failed: %w", err))
		}

		if useAnimation {
			fmt.Printf("✅ Generated %d animated clips\n", len(mediaPaths))
			for _, path := range mediaPaths {
				fmt.Printf("Clip saved to: %s\n", path)
			}
			// Generate video by concatenating animated clips
			if err := videoGen.GenerateFromAnimatedClips(mediaPaths, finalAudioPath, videoPath, segmentDurations); err != nil {
				return m.failJob(job, fmt.Errorf("video generation failed: %w", err))
			}
		} else {
			fmt.Printf("✅ Generated %d images\n", len(mediaPaths))
			for _, path := range mediaPaths {
				fmt.Printf("Images saved to: %s\n", path)
			}
			// Generate video from images with Ken Burns effects and segment timing
			if err := videoGen.GenerateFromImages(mediaPaths, finalAudioPath, videoPath, segmentDurations); err != nil {
				return m.failJob(job, fmt.Errorf("video generation failed: %w", err))
			}
		}
	} else {
		// Generate video with standard backgrounds
		if err := videoGen.GenerateFromAudio(finalAudioPath, videoPath, job.Background, true); err != nil {
			return m.failJob(job, fmt.Errorf("video generation failed: %w", err))
		}
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

// animateCortex shows the dancing Cortex robot during long operations
func (m *Manager) animateCortex(message string, stop chan bool) {
	frame := 0
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			// Clear previous robot (4 lines) - always clear to maintain position
			if frame > 0 {
				fmt.Print("\033[4A\033[J") // Move up 4 lines and clear to end
			}
			m.ui.ShowCortexWithMessage(message, frame)
			frame++
		}
	}
}
