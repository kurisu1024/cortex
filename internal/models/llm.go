package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// LLMClient handles interactions with Ollama
type LLMClient struct {
	host   string
	model  string
	client *http.Client
}

// NewLLMClient creates a new LLM client
func NewLLMClient(host, model string) *LLMClient {
	return &LLMClient{
		host:  host,
		model: model,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// Start ensures Ollama is running
func (l *LLMClient) Start() error {
	// Check if Ollama is already running
	if l.IsHealthy() {
		return nil
	}

	// Note: Ollama should be started manually or as a system service
	// This is just a health check
	return fmt.Errorf("Ollama is not running. Please start it with: ollama serve")
}

// Stop is a no-op since Ollama runs as a separate service
func (l *LLMClient) Stop() error {
	return nil
}

// IsHealthy checks if Ollama is responsive
func (l *LLMClient) IsHealthy() bool {
	resp, err := l.client.Get(l.host + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// GenerateRequest represents a generation request to Ollama
type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// GenerateResponse represents a response from Ollama
type GenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// Generate generates text using the LLM
func (l *LLMClient) Generate(prompt string) (string, error) {
	reqBody := GenerateRequest{
		Model:  l.model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := l.client.Post(
		l.host+"/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusNotFound {
			return "", fmt.Errorf("ollama returned status %d: %s\n\nℹ️  The model '%s' is not installed. Download it with:\n   ollama pull %s\n\nOr use the make command:\n   make download-model", resp.StatusCode, string(body), l.model, l.model)
		}
		return "", fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	var result GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Response, nil
}

// GenerateStream generates text with streaming support
func (l *LLMClient) GenerateStream(prompt string, callback func(string)) error {
	reqBody := GenerateRequest{
		Model:  l.model,
		Prompt: prompt,
		Stream: true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := l.client.Post(
		l.host+"/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	for {
		var result GenerateResponse
		if err := decoder.Decode(&result); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if result.Response != "" {
			callback(result.Response)
		}

		if result.Done {
			break
		}
	}

	return nil
}
