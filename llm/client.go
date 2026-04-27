// Package llm provides a client for interacting with LLM APIs
// to generate git commit messages and related content.
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/HexmosTech/git-lrc/config"
)

const (
	defaultTimeout    = 30 * time.Second
	openAIAPIEndpoint = "https://api.openai.com/v1/chat/completions"
)

// Message represents a single message in a chat completion request.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest is the payload sent to the LLM API.
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
}

// ChatResponse is the response received from the LLM API.
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// Client handles communication with the LLM API.
type Client struct {
	httpClient *http.Client
	cfg        *config.Config
}

// NewClient creates a new LLM client using the provided configuration.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		cfg: cfg,
	}
}

// Complete sends a prompt to the LLM and returns the generated text.
func (c *Client) Complete(ctx context.Context, prompt string) (string, error) {
	reqBody := ChatRequest{
		Model: c.cfg.LLMModel,
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
		Temperature: 0.3,
		MaxTokens:   512,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshalling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIAPIEndpoint, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.LLMAPIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body: %w", err)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("unmarshalling response: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("API error (%s): %s", chatResp.Error.Type, chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from API")
	}

	return chatResp.Choices[0].Message.Content, nil
}
