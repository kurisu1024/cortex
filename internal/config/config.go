package config

import (
	"github.com/spf13/viper"
)

// Config holds all configuration for Cortex
type Config struct {
	Models ModelsConfig `mapstructure:"models"`
	Output OutputConfig `mapstructure:"output"`
}

// ModelsConfig holds model-specific configuration
type ModelsConfig struct {
	Ollama OllamaConfig `mapstructure:"ollama"`
	TTS    TTSConfig    `mapstructure:"tts"`
}

// OllamaConfig holds Ollama LLM configuration
type OllamaConfig struct {
	Host  string `mapstructure:"host"`
	Model string `mapstructure:"model"`
}

// TTSConfig holds TTS configuration
type TTSConfig struct {
	Engine string            `mapstructure:"engine"`
	Voice  string            `mapstructure:"voice"`  // Default voice
	Voices map[string]string `mapstructure:"voices"` // Multiple voices: speaker_name -> voice_path
}

// OutputConfig holds output-related configuration
type OutputConfig struct {
	Directory string      `mapstructure:"directory"`
	Duration  int         `mapstructure:"duration"` // Target video duration in minutes
	Video     VideoConfig `mapstructure:"video"`
}

// VideoConfig holds video generation configuration
type VideoConfig struct {
	Format     string `mapstructure:"format"`
	Background string `mapstructure:"background"`
	Waveform   bool   `mapstructure:"waveform"`
}

// Load loads the configuration from viper
func Load() (*Config, error) {
	var cfg Config

	// Set defaults
	viper.SetDefault("models.ollama.host", "http://localhost:11434")
	viper.SetDefault("models.ollama.model", "llama3.2")
	viper.SetDefault("models.tts.engine", "say")       // Use macOS say by default
	viper.SetDefault("models.tts.voice", "Samantha")   // Default macOS voice
	viper.SetDefault("models.tts.voices", map[string]string{
		"narrator":    "Samantha", // Female, natural
		"host":        "Alex",     // Male, clear
		"expert":      "Victoria", // Female, professional
		"interviewer": "Daniel",   // Male, British
		"analyst":     "Karen",    // Female, Australian
	})
	viper.SetDefault("output.directory", "./output")
	viper.SetDefault("output.duration", 10) // 10 minutes default
	viper.SetDefault("output.video.format", "mp4")
	viper.SetDefault("output.video.background", "gradient")
	viper.SetDefault("output.video.waveform", true)

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}