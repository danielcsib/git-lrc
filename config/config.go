// Package config handles loading and validation of application configuration
// from environment variables and .env files.
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration values required by git-lrc.
type Config struct {
	// LLMAPIKey is the API key used to authenticate with the LLM provider.
	LLMAPIKey string

	// LLMModel specifies which model to use for generating commit messages.
	LLMModel string

	// LLMAPIEndpoint is the base URL for the LLM API (supports OpenAI-compatible endpoints).
	LLMAPIEndpoint string

	// MaxTokens is the maximum number of tokens to request in LLM responses.
	MaxTokens int
}

// Load reads configuration from the environment, optionally loading a .env file
// from the current working directory if it exists.
func Load() (*Config, error) {
	// Attempt to load .env file; ignore error if it doesn't exist.
	_ = godotenv.Load()

	cfg := &Config{
		LLMAPIKey: os.Getenv("LLM_API_KEY"),
		// Switched to gpt-4.1-nano — even cheaper and still great for short commit messages.
		LLMModel:       getEnvWithDefault("LLM_MODEL", "gpt-4.1-nano"),
		LLMAPIEndpoint: getEnvWithDefault("LLM_API_ENDPOINT", "https://api.openai.com/v1"),
		// Bumped to 512 tokens to give the model a bit more room for detailed messages.
		MaxTokens: 512,
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks that all required configuration fields are present.
func (c *Config) validate() error {
	if c.LLMAPIKey == "" {
		return fmt.Errorf("config: LLM_API_KEY environment variable is required")
	}
	if c.LLMModel == "" {
		return fmt.Errorf("config: LLM_MODEL must not be empty")
	}
	if c.LLMAPIEndpoint == "" {
		return fmt.Errorf("config: LLM_API_ENDPOINT must not be empty")
	}
	return nil
}

// getEnvWithDefault returns the value of the environment variable named by key,
// or defaultValue if the variable is not set or is empty.
func getEnvWithDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
