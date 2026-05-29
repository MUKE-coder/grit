package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"gritdemo/internal/config"
)

// SecObsBridge talks to the locally-mounted Sentinel + Pulse APIs over
// loopback. Tokens are cached and silently refreshed on 401.
type SecObsBridge struct {
	cfg         *config.Config
	httpClient  *http.Client
	mu          sync.Mutex
	sentinelTok string
	pulseTok    string
}

var ErrUpstream = errors.New("secobs upstream")

func NewSecObsBridge(cfg *config.Config) *SecObsBridge {
	return &SecObsBridge{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (b *SecObsBridge) baseURL() string {
	port := b.cfg.Port
	if port == "" {
		port = "8080"
	}
	return "http://127.0.0.1:" + port
}

func (b *SecObsBridge) SentinelGet(ctx context.Context, path string, out interface{}) error {
	return b.do(ctx, "sentinel", path, out)
}

func (b *SecObsBridge) PulseGet(ctx context.Context, path string, out interface{}) error {
	return b.do(ctx, "pulse", path, out)
}

func (b *SecObsBridge) do(ctx context.Context, kind, path string, out interface{}) error {
	tok, err := b.token(ctx, kind)
	if err != nil {
		return fmt.Errorf("%w: %s auth: %v", ErrUpstream, kind, err)
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, b.baseURL()+path, nil)
	req.Header.Set("Authorization", "Bearer "+tok)

	res, err := b.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %s GET %s: %v", ErrUpstream, kind, path, err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		b.invalidate(kind)
		tok2, err := b.token(ctx, kind)
		if err != nil {
			return fmt.Errorf("%w: %s re-auth: %v", ErrUpstream, kind, err)
		}
		req2, _ := http.NewRequestWithContext(ctx, http.MethodGet, b.baseURL()+path, nil)
		req2.Header.Set("Authorization", "Bearer "+tok2)
		res, err = b.httpClient.Do(req2)
		if err != nil {
			return fmt.Errorf("%w: %s retry: %v", ErrUpstream, kind, err)
		}
		defer res.Body.Close()
	}

	if res.StatusCode >= 400 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("%w: %s %d: %s", ErrUpstream, path, res.StatusCode, strings.TrimSpace(string(body)))
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(res.Body).Decode(out)
}

func (b *SecObsBridge) token(ctx context.Context, kind string) (string, error) {
	b.mu.Lock()
	cached := b.sentinelTok
	if kind == "pulse" {
		cached = b.pulseTok
	}
	b.mu.Unlock()
	if cached != "" {
		return cached, nil
	}

	var loginPath, user, pass string
	switch kind {
	case "sentinel":
		loginPath = "/sentinel/api/auth/login"
		user = b.cfg.SentinelUsername
		pass = b.cfg.SentinelPassword
	case "pulse":
		loginPath = "/pulse/api/auth/login"
		user = b.cfg.PulseUsername
		pass = b.cfg.PulsePassword
	default:
		return "", fmt.Errorf("unknown kind %q", kind)
	}

	payload := fmt.Sprintf(`{"username":%q,"password":%q}`, user, pass)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, b.baseURL()+loginPath, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	res, err := b.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("login %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var envelope struct {
		Token       string `json:"token"`
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&envelope); err != nil {
		return "", err
	}
	tok := envelope.Token
	if tok == "" {
		tok = envelope.AccessToken
	}
	if tok == "" {
		return "", errors.New("no token in login response")
	}

	b.mu.Lock()
	if kind == "pulse" {
		b.pulseTok = tok
	} else {
		b.sentinelTok = tok
	}
	b.mu.Unlock()

	return tok, nil
}

func (b *SecObsBridge) invalidate(kind string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if kind == "pulse" {
		b.pulseTok = ""
	} else {
		b.sentinelTok = ""
	}
}
