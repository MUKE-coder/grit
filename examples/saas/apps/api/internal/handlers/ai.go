package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"saas/apps/api/internal/ai"
	"saas/apps/api/internal/respond"
)

// aiFailure turns an error from the ai package into something the caller can
// act on, and logs the rest.
//
// These cases need different responses from whoever hits them and used to be
// one message: a key the gateway rejected, a rate limit, a model that does not
// exist, and everything else. The gateway's own reply is logged rather than
// returned, since it is somebody else's error text and may quote the request
// back.
//
// Worth having because the generic version cost real time: the API said
// "Failed to generate completion" while the gateway had said "Authentication
// failed. Create an API key and set in AI_GATEWAY_API_KEY", which is the whole
// answer.
func aiFailure(c *gin.Context, op string, err error) {
	log.Printf("[ai] %s: %v", op, err)

	msg := err.Error()
	switch {
	case strings.Contains(msg, "(401)"):
		respond.Fail(c, respond.CodeAIUnauthorized, "The AI gateway rejected the API key. Check AI_GATEWAY_API_KEY "+"in your .env, and that it is a gateway key rather than a provider key.")
	case strings.Contains(msg, "(403)"):
		// 403 from the gateway is usually entitlement rather than identity:
		// the key is real and the plan does not cover the model. The default
		// AI_GATEWAY_MODEL is a large one, so this is what a free-tier key
		// meets first, and "rejected the key" sends people to look in the
		// wrong place.
		respond.Fail(c, respond.CodeAIForbidden, "The AI gateway accepted the key but refused the request. This is "+"usually the plan not covering AI_GATEWAY_MODEL. The server log has "+"the gateway's own words.")
	case strings.Contains(msg, "(429)"):
		respond.Fail(c, respond.CodeAIRateLimited, "The AI gateway is rate limiting this key. Try again shortly.")
	case strings.Contains(msg, "(404)"):
		respond.Fail(c, respond.CodeAIModelNotFound, "The AI gateway does not know that model. Check AI_GATEWAY_MODEL "+"in your .env.")
	default:
		respond.Fail(c, respond.CodeAIError, "The AI gateway did not answer. The server log has the detail.")
	}
}

// AIHandler handles AI completion endpoints.
type AIHandler struct {
	AI *ai.AI
}

type CompleteRequest struct {
	Prompt      string  `json:"prompt" binding:"required"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

type ChatRequest struct {
	Messages    []ai.Message `json:"messages" binding:"required"`
	MaxTokens   int          `json:"max_tokens"`
	Temperature float64      `json:"temperature"`
}

// Complete handles a single prompt completion.
func (h *AIHandler) Complete(c *gin.Context) {
	if h.AI == nil {
		respond.Fail(c, respond.CodeAIUnavailable, "AI service is not configured")
		return
	}

	var req CompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	resp, err := h.AI.Complete(c.Request.Context(), ai.CompletionRequest{
		Prompt:      req.Prompt,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	})
	if err != nil {
		aiFailure(c, "complete", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}

// Chat handles a multi-turn conversation.
func (h *AIHandler) Chat(c *gin.Context) {
	if h.AI == nil {
		respond.Fail(c, respond.CodeAIUnavailable, "AI service is not configured")
		return
	}

	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	resp, err := h.AI.Complete(c.Request.Context(), ai.CompletionRequest{
		Messages:    req.Messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	})
	if err != nil {
		aiFailure(c, "chat", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}

// Stream handles a streaming completion via SSE.
func (h *AIHandler) Stream(c *gin.Context) {
	if h.AI == nil {
		respond.Fail(c, respond.CodeAIUnavailable, "AI service is not configured")
		return
	}

	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
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
