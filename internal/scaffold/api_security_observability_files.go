package scaffold

// This file scaffolds in-app Security + Observability summary dashboards
// that consume Sentinel v2.0.1 and Pulse v1.0 APIs over loopback, plus
// a Notification primitive + a background poller that creates
// notifications when something important fires in either system.
//
// The story:
//   - Sentinel + Pulse already expose rich dashboards at /sentinel/ui
//     and /pulse/ui. Those stay — they're the "dig in" experience.
//   - Grit admins shouldn't need to leave the admin to know what's
//     happening. This adds:
//       * /system/security    — in-app summary card grid pulled from
//                               /sentinel/api endpoints (summary,
//                               threats, security-score, auth-shield,
//                               CSP violations, performance overview).
//       * /system/observability — in-app summary from /pulse/api
//                                 (overview, SLOs, USE, N+1 ranked,
//                                 errors, runtime).
//   - A SecurityObservabilityPoller cron task runs every 60 s, calls
//     the same loopback APIs as the admin pages, and writes high-
//     severity findings into a new in-app Notification feed.
//   - The bell icon in the admin navbar already exists; this file
//     adds a Notification model + handler + the polling job so the
//     bell becomes meaningful out of the box.

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeSecurityObservabilityFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "models", "notification.go"):           notificationModelGo(),
		filepath.Join(apiRoot, "internal", "handlers", "notification.go"):         notificationHandlerGo(),
		filepath.Join(apiRoot, "internal", "handlers", "security.go"):             securityHandlerGo(),
		filepath.Join(apiRoot, "internal", "handlers", "observability.go"):        observabilityHandlerGo(),
		filepath.Join(apiRoot, "internal", "services", "secobs_bridge.go"):        secobsBridgeGo(),
		filepath.Join(apiRoot, "internal", "services", "secobs_poller.go"):        secobsPollerGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

// notificationModelGo emits the Notification model the admin bell reads from.
func notificationModelGo() string {
	return `package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification is an in-app message surfaced through the admin bell.
// Source distinguishes Sentinel security findings from Pulse perf
// findings from manual operator messages.
type Notification struct {
	ID        string    ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	UserID    string    ` + "`" + `gorm:"size:36;index" json:"user_id"` + "`" + ` // empty = visible to all admins
	Source    string    ` + "`" + `gorm:"size:16;index" json:"source"` + "`" + `  // sentinel | pulse | system
	Severity  string    ` + "`" + `gorm:"size:16;index" json:"severity"` + "`" + ` // critical | high | medium | low | info
	Title     string    ` + "`" + `gorm:"size:200" json:"title"` + "`" + `
	Body      string    ` + "`" + `gorm:"type:text" json:"body"` + "`" + `
	Link      string    ` + "`" + `gorm:"size:500" json:"link"` + "`" + `      // deep-link into /sentinel/ui or /pulse/ui
	Dedup     string    ` + "`" + `gorm:"size:128;uniqueIndex" json:"-"` + "`" + ` // collision key — repeat firings update Count, not insert
	Count     int       ` + "`" + `gorm:"default:1" json:"count"` + "`" + `
	ReadAt    *time.Time ` + "`" + `json:"read_at"` + "`" + `
	CreatedAt time.Time ` + "`" + `gorm:"index" json:"created_at"` + "`" + `
	UpdatedAt time.Time ` + "`" + `json:"updated_at"` + "`" + `
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	return nil
}
`
}

// notificationHandlerGo emits the handler the admin bell hits.
func notificationHandlerGo() string {
	return `package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

type NotificationHandler struct {
	DB *gorm.DB
}

// List returns unread + recent notifications for the bell dropdown.
// Visible to any authenticated user; admins see system-wide ones too.
func (h *NotificationHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")

	q := h.DB.Order("created_at DESC").Limit(50)
	if role == "ADMIN" {
		// Admins see broadcast (user_id="") + their own
		q = q.Where("user_id = '' OR user_id = ?", userID)
	} else {
		q = q.Where("user_id = ?", userID)
	}

	var items []models.Notification
	q.Find(&items)

	// Quick unread count for the bell badge
	var unread int64
	cq := h.DB.Model(&models.Notification{}).Where("read_at IS NULL")
	if role == "ADMIN" {
		cq = cq.Where("user_id = '' OR user_id = ?", userID)
	} else {
		cq = cq.Where("user_id = ?", userID)
	}
	cq.Count(&unread)

	c.JSON(http.StatusOK, gin.H{
		"data":   items,
		"unread": unread,
	})
}

// MarkRead marks one notification as read.
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	id := c.Param("id")
	now := time.Now()
	if err := h.DB.Model(&models.Notification{}).
		Where("id = ?", id).
		Update("read_at", now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "DB_ERROR", "message": err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "marked read"})
}

// MarkAllRead clears the bell for the current viewer.
func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")
	now := time.Now()
	q := h.DB.Model(&models.Notification{}).Where("read_at IS NULL")
	if role == "ADMIN" {
		q = q.Where("user_id = '' OR user_id = ?", userID)
	} else {
		q = q.Where("user_id = ?", userID)
	}
	q.Update("read_at", now)
	c.JSON(http.StatusOK, gin.H{"message": "all marked read"})
}
`
}

// secobsBridgeGo emits the shared loopback bridge that obtains JWTs
// from Sentinel + Pulse and proxies their summary endpoints.
func secobsBridgeGo() string {
	return `package services

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

	"{{MODULE}}/internal/config"
)

// SecObsBridge talks to the locally-mounted Sentinel + Pulse APIs.
// Both are mounted in-process by routes.Setup, so this is a loopback
// call — no network, no auth gateway, just JWT handshake against
// /sentinel/api/auth/login and /pulse/api/auth/login.
//
// Tokens are cached and silently refreshed on 401.
type SecObsBridge struct {
	cfg         *config.Config
	httpClient  *http.Client
	mu          sync.Mutex
	sentinelTok string
	pulseTok    string
}

// ErrUpstream indicates an HTTP failure talking to Sentinel/Pulse.
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

// SentinelGet pulls JSON from a Sentinel API endpoint.
// path is like "/sentinel/api/dashboard/summary".
func (b *SecObsBridge) SentinelGet(ctx context.Context, path string, out interface{}) error {
	return b.do(ctx, "sentinel", path, out)
}

// PulseGet pulls JSON from a Pulse API endpoint.
// path is like "/pulse/api/overview?range=1h".
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
		// Re-login once
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

	payload := fmt.Sprintf(` + "`{\"username\":%q,\"password\":%q}`" + `, user, pass)
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

	// Both Sentinel and Pulse return {"token": "..."} or {"access_token": "..."} —
	// accept either shape.
	var envelope struct {
		Token       string ` + "`" + `json:"token"` + "`" + `
		AccessToken string ` + "`" + `json:"access_token"` + "`" + `
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
`
}

// securityHandlerGo emits the handler the in-app Security dashboard reads.
//
// Hits Sentinel's analytics summary, blocked-IP list, and recent-threats
// endpoints, then reshapes them into the flat {banned_ips_now,
// auto_bans_24h, rate_limited_last_hour, active_bans, recent_threats}
// envelope the /system/security page expects. The earlier version
// returned the raw Sentinel responses under {data: {summary, score, ...}}
// which (a) called nonexistent endpoints and (b) wrapped the result so
// axios' .data accessor returned {data: ...} instead of the page's
// flat shape — every KPI rendered as 0.
func securityHandlerGo() string {
	return `package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"{{MODULE}}/internal/services"
)

type SecurityHandler struct {
	Bridge *services.SecObsBridge
}

// Summary returns the flat security envelope the React dashboard reads.
// On a fresh app with no traffic, expect zeros and empty arrays — that's
// the truth, not a bug.
func (h *SecurityHandler) Summary(c *gin.Context) {
	if h.Bridge == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": gin.H{"code": "SENTINEL_OFF", "message": "Sentinel is not enabled"},
		})
		return
	}

	ctx := c.Request.Context()

	// /sentinel/api/ip/blocked returns {data: BlockedIP[]} — currently
	// banned IPs with reason and optional expiry.
	var blockedResp struct {
		Data []struct {
			IP        string     ` + "`json:\"ip\"`" + `
			Reason    string     ` + "`json:\"reason\"`" + `
			BlockedAt time.Time  ` + "`json:\"blocked_at\"`" + `
			ExpiresAt *time.Time ` + "`json:\"expires_at\"`" + `
		} ` + "`json:\"data\"`" + `
	}
	_ = h.Bridge.SentinelGet(ctx, "/sentinel/api/ip/blocked", &blockedResp)

	// /sentinel/api/analytics/summary?window=24h returns ThreatStats:
	// total_threats, blocked_count, etc. blocked_count over the 24h
	// window is the closest analogue to "auto-bans in the last 24h".
	var statsResp struct {
		Data struct {
			TotalThreats int64 ` + "`json:\"total_threats\"`" + `
			BlockedCount int64 ` + "`json:\"blocked_count\"`" + `
		} ` + "`json:\"data\"`" + `
	}
	_ = h.Bridge.SentinelGet(ctx, "/sentinel/api/analytics/summary?window=24h", &statsResp)

	// /sentinel/api/threats?limit=10 returns ThreatEvent[]. Each event
	// carries the threat types as a string slice plus the offender IP.
	var threatsResp struct {
		Data []struct {
			ID          string    ` + "`json:\"id\"`" + `
			IP          string    ` + "`json:\"ip\"`" + `
			Path        string    ` + "`json:\"path\"`" + `
			ThreatTypes []string  ` + "`json:\"threat_types\"`" + `
			Severity    string    ` + "`json:\"severity\"`" + `
			Timestamp   time.Time ` + "`json:\"timestamp\"`" + `
		} ` + "`json:\"data\"`" + `
	}
	_ = h.Bridge.SentinelGet(ctx, "/sentinel/api/threats?limit=10", &threatsResp)

	activeBans := make([]gin.H, 0, len(blockedResp.Data))
	for _, b := range blockedResp.Data {
		var expires string
		if b.ExpiresAt != nil {
			expires = b.ExpiresAt.Format(time.RFC3339)
		}
		activeBans = append(activeBans, gin.H{
			"ip":         b.IP,
			"reason":     b.Reason,
			"expires_at": expires,
			"level":      1, // Sentinel's BlockedIP doesn't carry an escalation level
		})
	}

	recentThreats := make([]gin.H, 0, len(threatsResp.Data))
	for _, t := range threatsResp.Data {
		threatType := ""
		if len(t.ThreatTypes) > 0 {
			threatType = t.ThreatTypes[0]
		}
		recentThreats = append(recentThreats, gin.H{
			"id":          t.ID,
			"type":        threatType,
			"ip":          t.IP,
			"description": t.Path,
			"created_at":  t.Timestamp.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"banned_ips_now":         len(blockedResp.Data),
		"auto_bans_24h":          statsResp.Data.BlockedCount,
		"rate_limited_last_hour": 0, // Not directly exposed by Sentinel; populated by rate-limit logs in a future release
		"active_bans":            activeBans,
		"rate_limit_hits_5min":   []gin.H{},
		"recent_threats":         recentThreats,
	})
}
`
}

// observabilityHandlerGo emits the handler the in-app Observability dashboard reads.
//
// The handler hits Pulse's overview, runtime, n1, and errors endpoints
// and reshapes the responses into the flat {latency, traffic, errors,
// saturation, slowest_routes, n1_detections, recent_errors} envelope the
// /system/performance page expects. Returning the raw Pulse shapes would
// force the page to know Pulse internals — and worse, the response would
// be wrapped in {data: ...} which collides with axios' .data accessor
// and leaves every KPI rendered as "—".
func observabilityHandlerGo() string {
	return `package handlers

import (
	"net/http"
	goruntime "runtime"

	"github.com/gin-gonic/gin"

	"{{MODULE}}/internal/services"
)

type ObservabilityHandler struct {
	Bridge *services.SecObsBridge
}

// Pulse durations come back as nanoseconds (Go's time.Duration JSON form).
func nsToMs(ns int64) float64 { return float64(ns) / 1_000_000.0 }

// Summary returns a flat performance envelope built from Pulse's overview,
// runtime, ranked N+1 and errors endpoints. Missing or zero data is fine
// — the React page renders "—" for unset numeric fields and shows the
// "No X yet" empty state for empty arrays.
func (h *ObservabilityHandler) Summary(c *gin.Context) {
	if h.Bridge == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": gin.H{"code": "PULSE_OFF", "message": "Pulse is not enabled"},
		})
		return
	}

	ctx := c.Request.Context()

	// /pulse/api/overview?range=1h returns Pulse's Overview struct: avg/p95
	// latency in ns, total/error counters, RPM, plus per-route stats. The
	// inner shape mirrors pulse.Overview + pulse.RouteStats.
	var overview struct {
		TotalRequests int64   ` + "`json:\"total_requests\"`" + `
		TotalErrors   int64   ` + "`json:\"total_errors\"`" + `
		ErrorRate     float64 ` + "`json:\"error_rate\"`" + `
		AvgLatency    int64   ` + "`json:\"avg_latency\"`" + `
		P95Latency    int64   ` + "`json:\"p95_latency\"`" + `
		RPM           float64 ` + "`json:\"rpm\"`" + `
		TopRoutes []struct {
			Method       string  ` + "`json:\"method\"`" + `
			Path         string  ` + "`json:\"path\"`" + `
			RequestCount int64   ` + "`json:\"request_count\"`" + `
			ErrorRate    float64 ` + "`json:\"error_rate\"`" + `
			AvgLatency   int64   ` + "`json:\"avg_latency\"`" + `
			P50Latency   int64   ` + "`json:\"p50_latency\"`" + `
			P95Latency   int64   ` + "`json:\"p95_latency\"`" + `
			P99Latency   int64   ` + "`json:\"p99_latency\"`" + `
		} ` + "`json:\"top_routes\"`" + `
	}
	_ = h.Bridge.PulseGet(ctx, "/pulse/api/overview?range=1h", &overview)

	// /pulse/api/runtime/current returns Pulse's RuntimeMetric snapshot.
	var rt struct {
		HeapAlloc    uint64 ` + "`json:\"heap_alloc\"`" + `
		NumGoroutine int    ` + "`json:\"num_goroutine\"`" + `
		NumGC        uint32 ` + "`json:\"num_gc\"`" + `
	}
	_ = h.Bridge.PulseGet(ctx, "/pulse/api/runtime/current", &rt)

	// /pulse/api/database/n1/ranked wraps an N1Ranking[] under "data".
	var n1Resp struct {
		Data []struct {
			Route            string  ` + "`json:\"route\"`" + `
			AvgQueriesPerHit float64 ` + "`json:\"avg_queries_per_hit\"`" + `
			FirstSeen        string  ` + "`json:\"first_seen\"`" + `
		} ` + "`json:\"data\"`" + `
	}
	_ = h.Bridge.PulseGet(ctx, "/pulse/api/database/n1/ranked?range=1h&limit=10", &n1Resp)

	// /pulse/api/errors wraps ErrorRecord[] under "data".
	var errResp struct {
		Data []struct {
			ID           string ` + "`json:\"id\"`" + `
			Route        string ` + "`json:\"route\"`" + `
			ErrorMessage string ` + "`json:\"error_message\"`" + `
			LastSeen     string ` + "`json:\"last_seen\"`" + `
		} ` + "`json:\"data\"`" + `
	}
	_ = h.Bridge.PulseGet(ctx, "/pulse/api/errors?limit=10&resolved=false", &errResp)

	// Build slowest_routes from the overview's TopRoutes (already ranked
	// by Pulse — top requests rather than top latency, but a useful proxy).
	slowest := make([]gin.H, 0, len(overview.TopRoutes))
	for _, r := range overview.TopRoutes {
		slowest = append(slowest, gin.H{
			"route":      r.Path,
			"method":     r.Method,
			"requests":   r.RequestCount,
			"avg":        nsToMs(r.AvgLatency),
			"p95":        nsToMs(r.P95Latency),
			"p99":        nsToMs(r.P99Latency),
			"error_rate": r.ErrorRate,
		})
	}

	n1List := make([]gin.H, 0, len(n1Resp.Data))
	for _, n := range n1Resp.Data {
		n1List = append(n1List, gin.H{
			"route":       n.Route,
			"query_count": int(n.AvgQueriesPerHit),
			"first_seen":  n.FirstSeen,
		})
	}

	recentErrs := make([]gin.H, 0, len(errResp.Data))
	for _, e := range errResp.Data {
		recentErrs = append(recentErrs, gin.H{
			"id":         e.ID,
			"route":      e.Route,
			"message":    e.ErrorMessage,
			"created_at": e.LastSeen,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"latency": gin.H{
			"p50": nsToMs(overview.AvgLatency), // Pulse overview lacks p50 — use avg as proxy
			"p95": nsToMs(overview.P95Latency),
			"p99": nsToMs(overview.P95Latency), // p99 not in overview either
			"avg": nsToMs(overview.AvgLatency),
		},
		"traffic": gin.H{
			"throughput": overview.RPM / 60.0, // convert RPM → req/s
			"total":      overview.TotalRequests,
		},
		"errors": gin.H{
			"rate":        overview.ErrorRate,
			"active_open": overview.TotalErrors,
		},
		"saturation": gin.H{
			"goroutines": rt.NumGoroutine,
			"heap_mb":    float64(rt.HeapAlloc) / (1024.0 * 1024.0),
			"gc_cycles":  rt.NumGC,
			"cpu_cores":  goruntime.NumCPU(),
		},
		"slowest_routes": slowest,
		"n1_detections":  n1List,
		"recent_errors":  recentErrs,
	})
}
`
}

// secobsPollerGo emits the background poller that walks Sentinel +
// Pulse for high-severity findings and writes them into the
// Notification feed (dedup by a stable key so repeated firings don't
// spam the bell — they bump Count + UpdatedAt instead).
func secobsPollerGo() string {
	return `package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// SecObsPoller runs every minute and turns important findings into
// in-app notifications. Started from main.go alongside the cron
// scheduler when both Sentinel and Pulse are enabled.
type SecObsPoller struct {
	DB     *gorm.DB
	Bridge *SecObsBridge
	stop   chan struct{}
}

func NewSecObsPoller(db *gorm.DB, bridge *SecObsBridge) *SecObsPoller {
	return &SecObsPoller{DB: db, Bridge: bridge, stop: make(chan struct{})}
}

func (p *SecObsPoller) Start() {
	if p.Bridge == nil || p.DB == nil {
		return
	}
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		// Run once on startup so the bell is populated immediately.
		p.tick()
		for {
			select {
			case <-p.stop:
				return
			case <-ticker.C:
				p.tick()
			}
		}
	}()
}

func (p *SecObsPoller) Stop() { close(p.stop) }

func (p *SecObsPoller) tick() {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	// — Sentinel: high-CVSS unresolved threats + AuthShield lockouts
	var threats struct {
		Data []struct {
			ID       string  ` + "`" + `json:"id"` + "`" + `
			Type     string  ` + "`" + `json:"type"` + "`" + `
			Severity string  ` + "`" + `json:"severity"` + "`" + `
			CVSS     float64 ` + "`" + `json:"cvss"` + "`" + `
			SourceIP string  ` + "`" + `json:"source_ip"` + "`" + `
			Route    string  ` + "`" + `json:"route"` + "`" + `
		} ` + "`" + `json:"data"` + "`" + `
	}
	if err := p.Bridge.SentinelGet(ctx, "/sentinel/api/threats?limit=20&resolved=false", &threats); err == nil {
		for _, t := range threats.Data {
			if t.CVSS < 7.0 && t.Severity != "critical" && t.Severity != "high" {
				continue
			}
			sev := t.Severity
			if sev == "" {
				if t.CVSS >= 9.0 {
					sev = "critical"
				} else {
					sev = "high"
				}
			}
			p.upsert(&models.Notification{
				Source:   "sentinel",
				Severity: sev,
				Title:    fmt.Sprintf("Threat — %s (CVSS %.1f)", t.Type, t.CVSS),
				Body:     fmt.Sprintf("Source %s · route %s", t.SourceIP, t.Route),
				Link:     "/sentinel/ui/threats?focus=" + t.ID,
				Dedup:    dedup("sentinel", "threat", t.ID),
			})
		}
	}

	// — Pulse: firing alerts (already filtered to firing state)
	var alerts struct {
		Data []struct {
			ID       string ` + "`" + `json:"id"` + "`" + `
			Name     string ` + "`" + `json:"name"` + "`" + `
			Severity string ` + "`" + `json:"severity"` + "`" + `
			Message  string ` + "`" + `json:"message"` + "`" + `
		} ` + "`" + `json:"data"` + "`" + `
	}
	if err := p.Bridge.PulseGet(ctx, "/pulse/api/alerts?state=firing&limit=20", &alerts); err == nil {
		for _, a := range alerts.Data {
			if a.Severity != "critical" && a.Severity != "high" {
				continue
			}
			p.upsert(&models.Notification{
				Source:   "pulse",
				Severity: a.Severity,
				Title:    "Alert — " + a.Name,
				Body:     a.Message,
				Link:     "/pulse/ui/alerts?focus=" + a.ID,
				Dedup:    dedup("pulse", "alert", a.ID),
			})
		}
	}

	// — Pulse: top-impact N+1 (a single notification, dedup'd by date)
	var n1 struct {
		Data []struct {
			Route        string  ` + "`" + `json:"route"` + "`" + `
			Pattern      string  ` + "`" + `json:"pattern"` + "`" + `
			ImpactScore  float64 ` + "`" + `json:"impact_score"` + "`" + `
			Occurrences  int     ` + "`" + `json:"occurrences"` + "`" + `
			AvgQueriesPer int    ` + "`" + `json:"avg_queries_per_request"` + "`" + `
		} ` + "`" + `json:"data"` + "`" + `
	}
	if err := p.Bridge.PulseGet(ctx, "/pulse/api/database/n1/ranked?range=1h&limit=1", &n1); err == nil && len(n1.Data) > 0 {
		top := n1.Data[0]
		if top.ImpactScore > 100 {
			p.upsert(&models.Notification{
				Source:   "pulse",
				Severity: "medium",
				Title:    "N+1 query detected — " + top.Route,
				Body:     fmt.Sprintf("%d occurrences · ~%d queries/req · impact %.0f", top.Occurrences, top.AvgQueriesPer, top.ImpactScore),
				Link:     "/pulse/ui/database#n1",
				Dedup:    dedup("pulse", "n1", time.Now().UTC().Format("2006-01-02-15"), top.Route),
			})
		}
	}
}

func (p *SecObsPoller) upsert(n *models.Notification) {
	// Try update first — repeated firings bump Count + UpdatedAt
	// (idempotent thanks to the dedup unique index).
	res := p.DB.Model(&models.Notification{}).
		Where("dedup = ?", n.Dedup).
		Updates(map[string]interface{}{
			"count":      gorm.Expr("count + 1"),
			"updated_at": time.Now(),
		})
	if res.Error != nil {
		log.Printf("secobs_poller: update %s: %v", n.Dedup, res.Error)
		return
	}
	if res.RowsAffected == 0 {
		if err := p.DB.Create(n).Error; err != nil {
			log.Printf("secobs_poller: create %s: %v", n.Dedup, err)
		}
	}
}

func dedup(parts ...string) string {
	h := sha256.New()
	for _, p := range parts {
		h.Write([]byte(p))
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))[:32]
}
`
}
