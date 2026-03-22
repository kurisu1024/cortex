package models

import (
	"fmt"
	"strings"

	"github.com/topher/cortex/internal/config"
)

// Manager handles lifecycle of AI models
type Manager struct {
	llm *LLMClient
	tts *TTSClient
	cfg *config.Config
}

// NewManager creates a new model manager
func NewManager() *Manager {
	cfg, err := config.Load()
	if err != nil {
		// Use defaults if config fails to load
		cfg = &config.Config{}
	}

	return &Manager{
		llm: NewLLMClient(cfg.Models.Ollama.Host, cfg.Models.Ollama.Model),
		tts: NewTTSClient(cfg.Models.TTS.Engine, cfg.Models.TTS.Voice),
		cfg: cfg,
	}
}

// Start starts all required models
func (m *Manager) Start() error {
	if err := m.llm.Start(); err != nil {
		return fmt.Errorf("failed to start LLM: %w", err)
	}

	if err := m.tts.Start(); err != nil {
		return fmt.Errorf("failed to start TTS: %w", err)
	}

	return nil
}

// Stop stops all running models
func (m *Manager) Stop() error {
	var errors []string

	if err := m.llm.Stop(); err != nil {
		errors = append(errors, fmt.Sprintf("LLM: %v", err))
	}

	if err := m.tts.Stop(); err != nil {
		errors = append(errors, fmt.Sprintf("TTS: %v", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors stopping models: %s", strings.Join(errors, "; "))
	}

	return nil
}

// Status returns the status of all models
func (m *Manager) Status() string {
	var status strings.Builder

	llmStatus := "❌ Offline"
	if m.llm.IsHealthy() {
		llmStatus = "✅ Online"
	}

	ttsStatus := "❌ Offline"
	if m.tts.IsHealthy() {
		ttsStatus = "✅ Online"
	}

	status.WriteString(fmt.Sprintf("LLM (Ollama):  %s\n", llmStatus))
	status.WriteString(fmt.Sprintf("  Host:  %s\n", m.cfg.Models.Ollama.Host))
	status.WriteString(fmt.Sprintf("  Model: %s\n\n", m.cfg.Models.Ollama.Model))

	status.WriteString(fmt.Sprintf("TTS (%s): %s\n", m.cfg.Models.TTS.Engine, ttsStatus))
	status.WriteString(fmt.Sprintf("  Voice: %s\n", m.cfg.Models.TTS.Voice))

	return status.String()
}

// GetLLM returns the LLM client
func (m *Manager) GetLLM() *LLMClient {
	return m.llm
}

// GetTTS returns the TTS client
func (m *Manager) GetTTS() *TTSClient {
	return m.tts
}
