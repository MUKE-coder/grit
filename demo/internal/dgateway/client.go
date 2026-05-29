// Package dgateway is a thin HTTP client for Desispay's DGateway mobile-money
// API. We use it to collect loan repayments via MTN MoMo / Airtel Money, and
// to verify the status of pending transactions.
//
// Auth: header `X-Api-Key`. Two endpoints we care about:
//   POST /v1/payments/collect    — initiate a request-to-pay
//   POST /v1/webhooks/verify     — poll for a transaction's final status
//
// All amounts are in major units (UGX). Phones must be `256XXXXXXXXX` or
// `0XXXXXXXXX`. DGateway charges an 8% platform fee per collect by default —
// the fee/net split appears in the verify response, not the initiate response.
package dgateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultBaseURL = "https://dgateway.desispay.com/api"
	defaultTimeout = 15 * time.Second

	// Provider constants (Uganda).
	ProviderIotec   = "iotec"   // MTN + Airtel via iotec
	ProviderRelworx = "relworx" // MTN + Airtel + Visa via Relworx (multi-region)
)

// Status values returned in the verify response.
const (
	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
)

// Client talks to DGateway over HTTPS. Construct with New.
type Client struct {
	baseURL string
	apiKey  string
	hc      *http.Client
}

// New constructs a Client. Pass an empty baseURL to use the production default.
func New(apiKey, baseURL string) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		hc:      &http.Client{Timeout: defaultTimeout},
	}
}

// CollectRequest is the body of POST /v1/payments/collect.
type CollectRequest struct {
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"` // e.g. "UGX"
	PhoneNumber string  `json:"phone_number"`
	Provider    string  `json:"provider,omitempty"` // iotec | relworx; omit = default
	Reference   string  `json:"reference,omitempty"` // optional client-side correlation id
	Description string  `json:"description,omitempty"`
	CallbackURL string  `json:"callback_url,omitempty"`
}

// CollectResponse is the response body from a successful collect request.
type CollectResponse struct {
	Reference string  `json:"reference"` // dgw_xxxxxxxxx
	Status    string  `json:"status"`    // typically "pending"
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	Provider  string  `json:"provider"`
}

// VerifyResponse is the verify endpoint's response — works for both webhooks and polling.
type VerifyResponse struct {
	Reference     string  `json:"reference"`
	Status        string  `json:"status"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Provider      string  `json:"provider"`
	Direction     string  `json:"direction"` // "collect" | "disburse"
	FailureReason string  `json:"failure_reason"`
	NetAmount     float64 `json:"net_amount"`
	Fee           float64 `json:"fee"`
}

// WebhookEvent is the shape of the inbound webhook payload.
type WebhookEvent struct {
	Event     string         `json:"event"` // "transaction.completed" etc.
	Data      VerifyResponse `json:"data"`
	Timestamp string         `json:"timestamp"`
}

// Collect initiates a mobile-money request-to-pay. The returned reference is
// what we store on the Repayment row to correlate the webhook + verify polls.
func (c *Client) Collect(ctx context.Context, req CollectRequest) (*CollectResponse, error) {
	var out CollectResponse
	if err := c.do(ctx, "POST", "/v1/payments/collect", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Verify queries the current status of a transaction by its DGateway reference.
// Safe to call repeatedly (poll until status != pending).
func (c *Client) Verify(ctx context.Context, reference string) (*VerifyResponse, error) {
	var out VerifyResponse
	body := map[string]string{"reference": reference}
	if err := c.do(ctx, "POST", "/v1/webhooks/verify", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) do(ctx context.Context, method, path string, in any, out any) error {
	if c.apiKey == "" {
		return fmt.Errorf("dgateway: API key not configured")
	}

	var body io.Reader
	if in != nil {
		buf, err := json.Marshal(in)
		if err != nil {
			return fmt.Errorf("dgateway: marshal request: %w", err)
		}
		body = bytes.NewReader(buf)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return err
	}
	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("dgateway: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("dgateway: read body: %w", err)
	}

	if resp.StatusCode >= 400 {
		// DGateway returns {"error": "...", "code": "..."} on failure.
		var apiErr struct {
			Error string `json:"error"`
			Code  string `json:"code"`
		}
		_ = json.Unmarshal(respBody, &apiErr)
		if apiErr.Error != "" {
			return fmt.Errorf("dgateway %d: %s (code=%s)", resp.StatusCode, apiErr.Error, apiErr.Code)
		}
		return fmt.Errorf("dgateway %d: %s", resp.StatusCode, string(respBody))
	}

	if out != nil {
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("dgateway: decode response: %w (body=%s)", err, string(respBody))
		}
	}
	return nil
}
