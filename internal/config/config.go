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
	Image  ImageConfig  `mapstructure:"image"`
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

// ImageConfig holds AI image generation configuration
type ImageConfig struct {
	ModelID  string  `mapstructure:"model_id"`  // Stable Diffusion model ID
	ArtStyle string  `mapstructure:"art_style"` // Default art style for prompts
	Steps    int     `mapstructure:"steps"`     // Inference steps (quality vs speed)
	Guidance float64 `mapstructure:"guidance"`  // Guidance scale (how closely to follow prompt)
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
	viper.SetDefault("models.tts.engine", "edgetts")         // Use Edge TTS by default (high quality)
	viper.SetDefault("models.tts.voice", "en-US-AriaNeural") // Default Edge TTS voice
	viper.SetDefault("models.tts.voices", map[string]string{
		"narrator":    "en-US-AriaNeural",  // Female, professional, news-style
		"host":        "en-US-GuyNeural",   // Male, passionate
		"expert":      "en-US-JennyNeural", // Female, friendly
		"interviewer": "en-GB-RyanNeural",  // Male, British
		"analyst":     "en-US-EmmaNeural",  // Female, cheerful
	})
	viper.SetDefault("models.image.model_id", "runwayml/stable-diffusion-v1-5")
	viper.SetDefault("models.image.art_style", "cinematic, high quality, 4k, detailed")
	viper.SetDefault("models.image.steps", 30)
	viper.SetDefault("models.image.guidance", 7.5)
	viper.SetDefault("output.directory", "./output")
	viper.SetDefault("output.duration", 10) // 10 minutes default
	viper.SetDefault("output.video.format", "mp4")
	viper.SetDefault("output.video.background", "ai-generated")
	viper.SetDefault("output.video.waveform", true)

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
