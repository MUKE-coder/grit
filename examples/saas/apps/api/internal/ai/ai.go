package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Message represents a chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionRequest holds the input for a completion.
type CompletionRequest struct {
	Prompt      string    `json:"prompt"`
	Messages    []Message `json:"messages,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

// CompletionResponse holds the AI response.
type CompletionResponse struct {
	Content string `json:"content"`
	Model   string `json:"model"`
	Usage   *Usage `json:"usage,omitempty"`
}

// Usage contains token usage information.
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// StreamHandler is called for each chunk of a streamed response.
type StreamHandler func(chunk string) error

// AI provides text generation via Vercel AI Gateway.
// One API key, hundreds of models — use "provider/model" format
// (e.g. "anthropic/claude-sonnet-4.6", "openai/gpt-5.4", "google/gemini-2.5-pro").
type AI struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

// New creates a new AI service backed by Vercel AI Gateway.
func New(apiKey, model, baseURL string) *AI {
	if baseURL == "" {
		baseURL = "https://ai-gateway.vercel.sh/v1"
	}
	// Ensure no trailing slash
	baseURL = strings.TrimRight(baseURL, "/")

	return &AI{
		apiKey:  apiKey,
		model:   model,
		baseURL: baseURL,
		client:  &http.Client{Timeout: 120 * time.Second},
	}
}

// Complete generates a response from a single prompt or conversation.
func (a *AI) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	messages := req.Messages
	if len(messages) == 0 && req.Prompt != "" {
		messages = []Message{{Role: "user", Content: req.Prompt}}
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 1024
	}

	body := map[string]interface{}{
		"model":      a.model,
		"max_tokens": maxTokens,
		"messages":   messages,
	}
	if req.Temperature > 0 {
		body["temperature"] = req.Temperature
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+"/chat/completions", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+a.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("calling AI Gateway: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AI Gateway error (%d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Model string `json:"model"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	content := ""
	if len(result.Choices) > 0 {
		content = result.Choices[0].Message.Content
	}

	return &CompletionResponse{
		Content: content,
		Model:   result.Model,
		Usage: &Usage{
			InputTokens:  result.Usage.PromptTokens,
			OutputTokens: result.Usage.CompletionTokens,
		},
	}, nil
}

// Stream generates a streaming response, calling handler for each chunk.
func (a *AI) Stream(ctx context.Context, req CompletionRequest, handler StreamHandler) error {
	messages := req.Messages
	if len(messages) == 0 && req.Prompt != "" {
		messages = []Message{{Role: "user", Content: req.Prompt}}
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 1024
	}

	body := map[string]interface{}{
		"model":      a.model,
		"max_tokens": maxTokens,
		"messages":   messages,
		"stream":     true,
	}
	if req.Temperature > 0 {
		body["temperature"] = req.Temperature
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+"/chat/completions", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+a.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("calling AI Gateway: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("AI Gateway error (%d): %s", resp.StatusCode, string(respBody))
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var event struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}

		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}

		if len(event.Choices) > 0 && event.Choices[0].Delta.Content != "" {
			if err := handler(event.Choices[0].Delta.Content); err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}
