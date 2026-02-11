package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeAIFiles(root string, opts Options) error {
	apiRoot := filepath.Join(root, "apps", "api")
	module := opts.ProjectName + "/apps/api"

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "ai", "ai.go"):         aiServiceGo(),
		filepath.Join(apiRoot, "internal", "handlers", "ai.go"):   aiHandlerGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func aiServiceGo() string {
	return `package ai

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
	Role    string ` + "`" + `json:"role"` + "`" + `
	Content string ` + "`" + `json:"content"` + "`" + `
}

// CompletionRequest holds the input for a completion.
type CompletionRequest struct {
	Prompt      string    ` + "`" + `json:"prompt"` + "`" + `
	Messages    []Message ` + "`" + `json:"messages,omitempty"` + "`" + `
	MaxTokens   int       ` + "`" + `json:"max_tokens,omitempty"` + "`" + `
	Temperature float64   ` + "`" + `json:"temperature,omitempty"` + "`" + `
}

// CompletionResponse holds the AI response.
type CompletionResponse struct {
	Content string ` + "`" + `json:"content"` + "`" + `
	Model   string ` + "`" + `json:"model"` + "`" + `
	Usage   *Usage ` + "`" + `json:"usage,omitempty"` + "`" + `
}

// Usage contains token usage information.
type Usage struct {
	InputTokens  int ` + "`" + `json:"input_tokens"` + "`" + `
	OutputTokens int ` + "`" + `json:"output_tokens"` + "`" + `
}

// StreamHandler is called for each chunk of a streamed response.
type StreamHandler func(chunk string) error

// AI provides text generation via Claude or OpenAI APIs.
type AI struct {
	provider string
	apiKey   string
	model    string
	client   *http.Client
}

// New creates a new AI service instance.
func New(provider, apiKey, model string) *AI {
	return &AI{
		provider: strings.ToLower(provider),
		apiKey:   apiKey,
		model:    model,
		client:   &http.Client{Timeout: 120 * time.Second},
	}
}

// Complete generates a response from a single prompt.
func (a *AI) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	messages := req.Messages
	if len(messages) == 0 && req.Prompt != "" {
		messages = []Message{{Role: "user", Content: req.Prompt}}
	}

	if req.MaxTokens == 0 {
		req.MaxTokens = 1024
	}

	switch a.provider {
	case "openai":
		return a.openaiComplete(ctx, messages, req.MaxTokens, req.Temperature)
	default:
		return a.claudeComplete(ctx, messages, req.MaxTokens, req.Temperature)
	}
}

// Stream generates a streaming response, calling handler for each chunk.
func (a *AI) Stream(ctx context.Context, req CompletionRequest, handler StreamHandler) error {
	messages := req.Messages
	if len(messages) == 0 && req.Prompt != "" {
		messages = []Message{{Role: "user", Content: req.Prompt}}
	}

	if req.MaxTokens == 0 {
		req.MaxTokens = 1024
	}

	switch a.provider {
	case "openai":
		return a.openaiStream(ctx, messages, req.MaxTokens, req.Temperature, handler)
	default:
		return a.claudeStream(ctx, messages, req.MaxTokens, req.Temperature, handler)
	}
}

// ── Claude (Anthropic) ──────────────────────────────────────────

func (a *AI) claudeComplete(ctx context.Context, messages []Message, maxTokens int, temperature float64) (*CompletionResponse, error) {
	body := map[string]interface{}{
		"model":      a.model,
		"max_tokens": maxTokens,
		"messages":   messages,
	}
	if temperature > 0 {
		body["temperature"] = temperature
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.anthropic.com/v1/messages", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("x-api-key", a.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling Claude API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Claude API error (%d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Content []struct {
			Text string ` + "`" + `json:"text"` + "`" + `
		} ` + "`" + `json:"content"` + "`" + `
		Model string ` + "`" + `json:"model"` + "`" + `
		Usage struct {
			InputTokens  int ` + "`" + `json:"input_tokens"` + "`" + `
			OutputTokens int ` + "`" + `json:"output_tokens"` + "`" + `
		} ` + "`" + `json:"usage"` + "`" + `
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	content := ""
	if len(result.Content) > 0 {
		content = result.Content[0].Text
	}

	return &CompletionResponse{
		Content: content,
		Model:   result.Model,
		Usage: &Usage{
			InputTokens:  result.Usage.InputTokens,
			OutputTokens: result.Usage.OutputTokens,
		},
	}, nil
}

func (a *AI) claudeStream(ctx context.Context, messages []Message, maxTokens int, temperature float64, handler StreamHandler) error {
	body := map[string]interface{}{
		"model":      a.model,
		"max_tokens": maxTokens,
		"messages":   messages,
		"stream":     true,
	}
	if temperature > 0 {
		body["temperature"] = temperature
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.anthropic.com/v1/messages", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("x-api-key", a.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("calling Claude API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Claude API error (%d): %s", resp.StatusCode, string(respBody))
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
			Type  string ` + "`" + `json:"type"` + "`" + `
			Delta struct {
				Text string ` + "`" + `json:"text"` + "`" + `
			} ` + "`" + `json:"delta"` + "`" + `
		}

		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}

		if event.Type == "content_block_delta" && event.Delta.Text != "" {
			if err := handler(event.Delta.Text); err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}

// ── OpenAI ──────────────────────────────────────────────────────

func (a *AI) openaiComplete(ctx context.Context, messages []Message, maxTokens int, temperature float64) (*CompletionResponse, error) {
	body := map[string]interface{}{
		"model":      a.model,
		"max_tokens": maxTokens,
		"messages":   messages,
	}
	if temperature > 0 {
		body["temperature"] = temperature
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling OpenAI API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OpenAI API error (%d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string ` + "`" + `json:"content"` + "`" + `
			} ` + "`" + `json:"message"` + "`" + `
		} ` + "`" + `json:"choices"` + "`" + `
		Model string ` + "`" + `json:"model"` + "`" + `
		Usage struct {
			PromptTokens     int ` + "`" + `json:"prompt_tokens"` + "`" + `
			CompletionTokens int ` + "`" + `json:"completion_tokens"` + "`" + `
		} ` + "`" + `json:"usage"` + "`" + `
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

func (a *AI) openaiStream(ctx context.Context, messages []Message, maxTokens int, temperature float64, handler StreamHandler) error {
	body := map[string]interface{}{
		"model":      a.model,
		"max_tokens": maxTokens,
		"messages":   messages,
		"stream":     true,
	}
	if temperature > 0 {
		body["temperature"] = temperature
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("calling OpenAI API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("OpenAI API error (%d): %s", resp.StatusCode, string(respBody))
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
					Content string ` + "`" + `json:"content"` + "`" + `
				} ` + "`" + `json:"delta"` + "`" + `
			} ` + "`" + `json:"choices"` + "`" + `
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
`
}

func aiHandlerGo() string {
	return `package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"{{MODULE}}/internal/ai"
)

// AIHandler handles AI completion endpoints.
type AIHandler struct {
	AI *ai.AI
}

type completionRequest struct {
	Prompt      string  ` + "`" + `json:"prompt" binding:"required"` + "`" + `
	MaxTokens   int     ` + "`" + `json:"max_tokens"` + "`" + `
	Temperature float64 ` + "`" + `json:"temperature"` + "`" + `
}

type chatRequest struct {
	Messages    []ai.Message ` + "`" + `json:"messages" binding:"required"` + "`" + `
	MaxTokens   int          ` + "`" + `json:"max_tokens"` + "`" + `
	Temperature float64      ` + "`" + `json:"temperature"` + "`" + `
}

// Complete handles a single prompt completion.
func (h *AIHandler) Complete(c *gin.Context) {
	if h.AI == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": gin.H{
				"code":    "AI_UNAVAILABLE",
				"message": "AI service is not configured",
			},
		})
		return
	}

	var req completionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	resp, err := h.AI.Complete(c.Request.Context(), ai.CompletionRequest{
		Prompt:      req.Prompt,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "AI_ERROR",
				"message": "Failed to generate completion",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}

// Chat handles a multi-turn conversation.
func (h *AIHandler) Chat(c *gin.Context) {
	if h.AI == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": gin.H{
				"code":    "AI_UNAVAILABLE",
				"message": "AI service is not configured",
			},
		})
		return
	}

	var req chatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	resp, err := h.AI.Complete(c.Request.Context(), ai.CompletionRequest{
		Messages:    req.Messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "AI_ERROR",
				"message": "Failed to generate response",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}

// Stream handles a streaming completion via SSE.
func (h *AIHandler) Stream(c *gin.Context) {
	if h.AI == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": gin.H{
				"code":    "AI_UNAVAILABLE",
				"message": "AI service is not configured",
			},
		})
		return
	}

	var req chatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	err := h.AI.Stream(c.Request.Context(), ai.CompletionRequest{
		Messages:    req.Messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	}, func(chunk string) error {
		c.SSEvent("message", chunk)
		c.Writer.Flush()
		return nil
	})

	if err != nil {
		c.SSEvent("error", fmt.Sprintf("Stream error: %v", err))
		c.Writer.Flush()
	}

	c.SSEvent("done", "[DONE]")
	c.Writer.Flush()
}
`
}
