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
		filepath.Join(apiRoot, "internal", "models", "notification.go"):    notificationModelGo(),
		filepath.Join(apiRoot, "internal", "handlers", "notification.go"):  notificationHandlerGo(),
		filepath.Join(apiRoot, "internal", "handlers", "security.go"):      securityHandlerGo(),
		filepath.Join(apiRoot, "internal", "handlers", "observability.go"): observabilityHandlerGo(),
		filepath.Join(apiRoot, "internal", "services", "secobs_bridge.go"): secobsBridgeGo(),
		filepath.Join(apiRoot, "internal", "services", "secobs_poller.go"): secobsPollerGo(),
		filepath.Join(apiRoot, "internal", "services", "security_autoban.go"): securityAutoBanGo(),
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

	"{{MODULE}}/internal/ids"
	"gorm.io/gorm"
)

// Notification is an in-app message surfaced through the admin bell.
// Source distinguishes Sentinel security findings from Pulse perf
// findings from manual operator messages.
` + notificationIndexNote + notificationStructLine + `	ID        string    ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	UserID    string    ` + notificationUserIDTagNew + ` // empty = visible to all admins
	Source    string    ` + "`" + `gorm:"size:16;index" json:"source"` + "`" + `  // sentinel | pulse | system
	Severity  string    ` + "`" + `gorm:"size:16;index" json:"severity"` + "`" + ` // critical | high | medium | low | info
	Title     string    ` + "`" + `gorm:"size:200" json:"title"` + "`" + `
	Body      string    ` + "`" + `gorm:"type:text" json:"body"` + "`" + `
	Link      string    ` + "`" + `gorm:"size:500" json:"link"` + "`" + `      // deep-link into /sentinel/ui or /pulse/ui
	Dedup     string    ` + "`" + `gorm:"size:128;uniqueIndex" json:"-"` + "`" + ` // collision key — repeat firings update Count, not insert
	Count     int       ` + "`" + `gorm:"default:1" json:"count"` + "`" + `
	ReadAt    *time.Time ` + notificationReadAtTagNew + `
	CreatedAt time.Time ` + notificationCreatedAtTagNew + `
	UpdatedAt time.Time ` + "`" + `json:"updated_at"` + "`" + `
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = ids.New()
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
	"{{MODULE}}/internal/respond"
)

type NotificationHandler struct {
	DB *gorm.DB
}

// notificationListLimit is how many rows the bell's dropdown shows. The unread
// count beside it is exact; the list is only the most recent.
const notificationListLimit = 50

// scopeNotifications limits a query to the rows one viewer may see: their own,
// plus the broadcast rows (user_id "") when they are an ADMIN.
//
// One definition of the rule, spelled out three times before this and left out
// of MarkRead altogether, where any signed-in account could clear any
// notification by id.
func scopeNotifications(q *gorm.DB, c *gin.Context) *gorm.DB {
	userID, _ := c.Get("user_id")
	if role, _ := c.Get("user_role"); role == models.RoleAdmin {
		return q.Where("user_id = '' OR user_id = ?", userID)
	}
	return q.Where("user_id = ?", userID)
}

// List returns unread + recent notifications for the bell dropdown.
// Visible to any authenticated user; admins see system-wide ones too.
func (h *NotificationHandler) List(c *gin.Context) {
	var items []models.Notification
	q := scopeNotifications(h.DB.WithContext(c.Request.Context()), c).
		Order("created_at DESC").Limit(notificationListLimit)
	if err := q.Find(&items).Error; err != nil {
		respond.ServerError(c, "DB_ERROR", err, "Internal server error")
		return
	}

	// Quick unread count for the bell badge
	var unread int64
	cq := scopeNotifications(h.DB.WithContext(c.Request.Context()).Model(&models.Notification{}), c).
		Where("read_at IS NULL")
	if err := cq.Count(&unread).Error; err != nil {
		respond.ServerError(c, "DB_ERROR", err, "Internal server error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   items,
		"unread": unread,
	})
}

// MarkRead marks one notification as read.
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	now := time.Now()
	q := scopeNotifications(h.DB.WithContext(c.Request.Context()).Model(&models.Notification{}), c).
		Where("id = ?", c.Param("id"))
	res := q.Update("read_at", now)
	if res.Error != nil {
		respond.ServerError(c, "DB_ERROR", res.Error, "Internal server error")
		return
	}
	if res.RowsAffected == 0 {
		respond.Fail(c, respond.CodeNotFound, "notification not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "marked read"})
}

// MarkAllRead clears the bell for the current viewer.
func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	now := time.Now()
	q := scopeNotifications(h.DB.WithContext(c.Request.Context()).Model(&models.Notification{}), c).
		Where("read_at IS NULL")
	if err := q.Update("read_at", now).Error; err != nil {
		respond.ServerError(c, "DB_ERROR", err, "Internal server error")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "all marked read"})
}
`
}

// secobsBridgeGo emits the shared loopback bridge that obtains JWTs
// from Sentinel + Pulse and proxies their summary endpoints.
func secobsBridgeGo() string {
	return `package services

import (
	"bytes"
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

// SentinelPost sends JSON to a Sentinel endpoint: banning an address, lifting
// a ban, sweeping repeat offenders. body may be nil.
func (b *SecObsBridge) SentinelPost(ctx context.Context, path string, body, out interface{}) error {
	return b.send(ctx, "sentinel", http.MethodPost, path, body, out)
}

// SentinelDelete removes something from Sentinel, such as a ban.
func (b *SecObsBridge) SentinelDelete(ctx context.Context, path string, out interface{}) error {
	return b.send(ctx, "sentinel", http.MethodDelete, path, nil, out)
}

func (b *SecObsBridge) do(ctx context.Context, kind, path string, out interface{}) error {
	return b.send(ctx, kind, http.MethodGet, path, nil, out)
}

// send is the one request. It was a GET and nothing else, which was fine while
// these dashboards only read; banning an address is a write, and a second
// copy of the token handling is how the retry stops matching the first call.
func (b *SecObsBridge) send(ctx context.Context, kind, method, path string, body, out interface{}) error {
	// Encoded once, because a retry needs to send the same bytes and a
	// consumed reader sends none.
	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("%w: %s %s: encoding body: %v", ErrUpstream, kind, path, err)
		}
	}
	build := func(tok string) (*http.Request, error) {
		var reader io.Reader
		if payload != nil {
			reader = bytes.NewReader(payload)
		}
		req, err := http.NewRequestWithContext(ctx, method, b.baseURL()+path, reader)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+tok)
		if payload != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		return req, nil
	}

	tok, err := b.token(ctx, kind)
	if err != nil {
		return fmt.Errorf("%w: %s auth: %v", ErrUpstream, kind, err)
	}
	req, err := build(tok)
	if err != nil {
		return fmt.Errorf("%w: %s %s %s: %v", ErrUpstream, kind, method, path, err)
	}

	res, err := b.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %s %s %s: %v", ErrUpstream, kind, method, path, err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		// Re-login once
		b.invalidate(kind)
		tok2, err := b.token(ctx, kind)
		if err != nil {
			return fmt.Errorf("%w: %s re-auth: %v", ErrUpstream, kind, err)
		}
		req2, err := build(tok2)
		if err != nil {
			return fmt.Errorf("%w: %s retry: %v", ErrUpstream, kind, err)
		}
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
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"{{MODULE}}/internal/respond"
	"{{MODULE}}/internal/services"
)

type SecurityHandler struct {
	Bridge *services.SecObsBridge
	// AutoBan is the sweep, shared with the scheduled job so the button and
	// the schedule cannot drift apart on policy.
	AutoBan *services.SecurityAutoBan
	// AutoBanEvery is how often the schedule runs it, for the page to say.
	AutoBanEvery time.Duration
}

// upstream runs one call to Sentinel or Pulse and records the ones that fail.
//
// Every call on these two dashboards used to throw its error away with a blank
// assignment. With the upstream down, each struct kept its zero value and the
// page rendered a 200 full of noughts: no banned IPs, no threats, no errors and
// no latency. That reads as "you are safe and fast", which is the opposite of
// what had happened.
//
// The names collected here go back in the response as "degraded", so the page
// can say which panels are guesses.
type upstream struct {
	degraded []string
}

func (u *upstream) get(name string, call func() error) {
	if err := call(); err != nil {
		u.degraded = append(u.degraded, name)
		log.Printf("dashboard: %s is unavailable: %v", name, err)
	}
}

// names returns the failures, never nil, so the JSON field is [] and not null.
func (u *upstream) names() []string {
	if u.degraded == nil {
		return []string{}
	}
	return u.degraded
}

// allFailed reports whether not one of total calls came back.
func (u *upstream) allFailed(total int) bool { return len(u.degraded) == total }

// sentinelBlocked is /sentinel/api/ip/blocked: the currently banned IPs.
type sentinelBlocked struct {
	Data []struct {
		IP        string     ` + "`json:\"ip\"`" + `
		Reason    string     ` + "`json:\"reason\"`" + `
		BlockedAt time.Time  ` + "`json:\"blocked_at\"`" + `
		ExpiresAt *time.Time ` + "`json:\"expires_at\"`" + `
	} ` + "`json:\"data\"`" + `
}

// sentinelStats is /sentinel/api/analytics/summary: ThreatStats for a window.
type sentinelStats struct {
	Data struct {
		TotalThreats int64 ` + "`json:\"total_threats\"`" + `
		BlockedCount int64 ` + "`json:\"blocked_count\"`" + `
	} ` + "`json:\"data\"`" + `
}

// sentinelThreats is /sentinel/api/threats: the most recent ThreatEvents.
type sentinelThreats struct {
	Data []struct {
		ID          string    ` + "`json:\"id\"`" + `
		IP          string    ` + "`json:\"ip\"`" + `
		Path        string    ` + "`json:\"path\"`" + `
		ThreatTypes []string  ` + "`json:\"threat_types\"`" + `
		Severity    string    ` + "`json:\"severity\"`" + `
		Timestamp   time.Time ` + "`json:\"timestamp\"`" + `
	} ` + "`json:\"data\"`" + `
}

// sentinelScore is /sentinel/api/score: the graded posture, 0-100.
type sentinelScore struct {
	Data struct {
		Overall int    ` + "`json:\"overall\"`" + `
		Grade   string ` + "`json:\"grade\"`" + `
		Trend   string ` + "`json:\"trend\"`" + `
	} ` + "`json:\"data\"`" + `
}

// sentinelFullStats is the whole ThreatStats, not the two fields the first
// version of this dashboard read.
type sentinelFullStats struct {
	Data struct {
		TotalThreats   int64 ` + "`json:\"total_threats\"`" + `
		CriticalCount  int64 ` + "`json:\"critical_count\"`" + `
		HighCount      int64 ` + "`json:\"high_count\"`" + `
		MediumCount    int64 ` + "`json:\"medium_count\"`" + `
		LowCount       int64 ` + "`json:\"low_count\"`" + `
		BlockedCount   int64 ` + "`json:\"blocked_count\"`" + `
		UniqueIPs      int64 ` + "`json:\"unique_ips\"`" + `
		TopAttackTypes []struct {
			Type  string ` + "`json:\"type\"`" + `
			Count int64  ` + "`json:\"count\"`" + `
		} ` + "`json:\"top_attack_types\"`" + `
	} ` + "`json:\"data\"`" + `
}

// sentinelTrends is /sentinel/api/analytics/attack-trends: one bucket per
// interval.
type sentinelTrends struct {
	Data []struct {
		Period string           ` + "`json:\"period\"`" + `
		Total  int64            ` + "`json:\"total\"`" + `
		ByType map[string]int64 ` + "`json:\"by_type\"`" + `
	} ` + "`json:\"data\"`" + `
}

// sentinelTargets is /sentinel/api/analytics/top-targets: the routes being
// asked for.
type sentinelTargets struct {
	Data []struct {
		Route  string ` + "`json:\"route\"`" + `
		Method string ` + "`json:\"method\"`" + `
		Count  int64  ` + "`json:\"count\"`" + `
	} ` + "`json:\"data\"`" + `
}

// sentinelActors is /sentinel/api/actors: who is sending it.
type sentinelActors struct {
	Data []struct {
		IP          string    ` + "`json:\"ip\"`" + `
		ThreatCount int       ` + "`json:\"threat_count\"`" + `
		RiskScore   int       ` + "`json:\"risk_score\"`" + `
		Country     string    ` + "`json:\"country\"`" + `
		Status      string    ` + "`json:\"status\"`" + `
		LastSeen    time.Time ` + "`json:\"last_seen\"`" + `
	} ` + "`json:\"data\"`" + `
}

// sentinelThreatPage is the threat log with its total, which is what makes
// pagination possible.
type sentinelThreatPage struct {
	Data []struct {
		ID          string    ` + "`json:\"id\"`" + `
		IP          string    ` + "`json:\"ip\"`" + `
		Method      string    ` + "`json:\"method\"`" + `
		Path        string    ` + "`json:\"path\"`" + `
		ThreatTypes []string  ` + "`json:\"threat_types\"`" + `
		Severity    string    ` + "`json:\"severity\"`" + `
		Blocked     bool      ` + "`json:\"blocked\"`" + `
		Timestamp   time.Time ` + "`json:\"timestamp\"`" + `
	} ` + "`json:\"data\"`" + `
	Meta struct {
		Total    int ` + "`json:\"total\"`" + `
		Page     int ` + "`json:\"page\"`" + `
		PageSize int ` + "`json:\"page_size\"`" + `
	} ` + "`json:\"meta\"`" + `
}

// Summary returns everything the five cards and the score explainer read.
//
// It is one call rather than five because the page shows them together: five
// requests would paint the row in five steps, and any one of them failing
// would leave a card that is indistinguishable from a card that is zero. The
// degraded list is what tells those two apart.
//
// On a fresh app with no traffic, expect zeros and empty arrays: that is the
// truth, not a bug.
func (h *SecurityHandler) Summary(c *gin.Context) {
	if h.Bridge == nil {
		respond.Fail(c, respond.CodeSentinelOff, "Sentinel is not enabled")
		return
	}

	ctx := c.Request.Context()
	var up upstream

	get := func(name, path string, out interface{}) {
		up.get(name, func() error { return h.Bridge.SentinelGet(ctx, path, out) })
	}

	var blockedResp sentinelBlocked
	get("blocked_ips", "/sentinel/api/ip/blocked", &blockedResp)

	var statsResp sentinelFullStats
	get("threat_stats", "/sentinel/api/analytics/summary?window=24h", &statsResp)

	var scoreResp sentinelScore
	get("security_score", "/sentinel/api/score", &scoreResp)

	var threatsResp sentinelThreats
	get("recent_threats", "/sentinel/api/threats?limit=10", &threatsResp)

	// Nothing answered. Zeros here would be a claim that no IP is banned and
	// nothing has attacked this app, which is a worse answer than an error.
	if up.allFailed(4) {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": gin.H{
				"code":    "SENTINEL_UNAVAILABLE",
				"message": "Sentinel is not answering",
				"details": gin.H{"degraded": up.names()},
			},
		})
		return
	}

	// Sentinel's BlockedIP carries no flag saying who made the ban, so the
	// reason is the only evidence. services.IsAutoBan owns that reading, and
	// the sweep writes the prefix it looks for.
	automatic, manual := 0, 0
	activeBans := make([]gin.H, 0, len(blockedResp.Data))
	for _, b := range blockedResp.Data {
		var expires string
		if b.ExpiresAt != nil {
			expires = b.ExpiresAt.Format(time.RFC3339)
		}
		auto := services.IsAutoBan(b.Reason)
		if auto {
			automatic++
		} else {
			manual++
		}
		activeBans = append(activeBans, gin.H{
			"ip":         b.IP,
			"reason":     b.Reason,
			"blocked_at": b.BlockedAt.Format(time.RFC3339),
			"expires_at": expires,
			"permanent":  b.ExpiresAt == nil,
			"automatic":  auto,
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

	stats := statsResp.Data
	c.JSON(http.StatusOK, gin.H{
		"score":       scoreResp.Data.Overall,
		"grade":       scoreResp.Data.Grade,
		"score_trend": scoreResp.Data.Trend,

		"threats_24h":    stats.TotalThreats,
		"unique_ips_24h": stats.UniqueIPs,
		"critical_24h":   stats.CriticalCount,
		"high_24h":       stats.HighCount,
		"medium_24h":     stats.MediumCount,
		"low_24h":        stats.LowCount,
		"blocked_24h":    stats.BlockedCount,

		"banned_ips_now":   len(blockedResp.Data),
		"banned_automatic": automatic,
		"banned_manual":    manual,

		// Kept because the desktop system page and the older admin build read
		// them. Removing them would break a screen that nobody is changing.
		"auto_bans_24h":          stats.BlockedCount,
		"rate_limited_last_hour": 0,
		"rate_limit_hits_5min":   []gin.H{},

		"active_bans":    activeBans,
		"recent_threats": recentThreats,
		"auto_ban":       h.autoBanPolicy(),
		// The panels whose numbers are missing rather than zero.
		"degraded": up.names(),
	})
}

// autoBanPolicy is what the sweep will do, so the page can say it rather than
// make the reader go and find the configuration.
func (h *SecurityHandler) autoBanPolicy() gin.H {
	if h.AutoBan == nil {
		return gin.H{"enabled": false}
	}
	cfg := h.AutoBan.Config
	return gin.H{
		"enabled":     cfg.Enabled,
		"min_threats": cfg.MinThreats,
		"window":      cfg.Window.String(),
		"every":       h.AutoBanEvery.String(),
		"permanent":   cfg.Permanent,
	}
}

// Trends is the hourly shape of the last day: how many threats arrived, and
// how many of them were stopped rather than written down.
func (h *SecurityHandler) Trends(c *gin.Context) {
	if h.Bridge == nil {
		respond.Fail(c, respond.CodeSentinelOff, "Sentinel is not enabled")
		return
	}
	window := c.DefaultQuery("window", "24h")
	if _, err := time.ParseDuration(window); err != nil {
		respond.Fail(c, respond.CodeValidationError, "window must be a duration such as \"24h\"")
		return
	}
	var trends sentinelTrends
	if err := h.Bridge.SentinelGet(c.Request.Context(),
		"/sentinel/api/analytics/attack-trends?window="+url.QueryEscape(window)+"&interval=hour", &trends); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{
			"code": "SENTINEL_UNAVAILABLE", "message": "Sentinel is not answering"}})
		return
	}
	points := make([]gin.H, 0, len(trends.Data))
	for _, t := range trends.Data {
		points = append(points, gin.H{"period": t.Period, "total": t.Total, "by_type": t.ByType})
	}
	c.JSON(http.StatusOK, gin.H{"data": points, "window": window})
}

// Top is the three ranked lists under the chart: what is being tried, where,
// and by whom. One call, because they are read together and each is small.
func (h *SecurityHandler) Top(c *gin.Context) {
	if h.Bridge == nil {
		respond.Fail(c, respond.CodeSentinelOff, "Sentinel is not enabled")
		return
	}
	window := c.DefaultQuery("window", "168h")
	if _, err := time.ParseDuration(window); err != nil {
		respond.Fail(c, respond.CodeValidationError, "window must be a duration such as \"168h\"")
		return
	}
	ctx := c.Request.Context()
	var up upstream
	get := func(name, path string, out interface{}) {
		up.get(name, func() error { return h.Bridge.SentinelGet(ctx, path, out) })
	}

	// Attack types come from the trend buckets rather than from
	// analytics/summary's top_attack_types, which is filled by Sentinel's
	// in-memory store and by none of the persistent ones: on SQLite, Postgres
	// or MySQL, which is every real deployment, that field is always empty.
	// The per-bucket by_type map holds the same counts and is populated
	// everywhere, so it is summed here (sentinel#44).
	var trends sentinelTrends
	get("attack_types", "/sentinel/api/analytics/attack-trends?window="+url.QueryEscape(window)+"&interval=day", &trends)

	var targets sentinelTargets
	get("targeted_routes", "/sentinel/api/analytics/top-targets?window="+url.QueryEscape(window)+"&limit=10", &targets)

	// Sentinel's ListActors orders by last_seen and ignores sort_by in every
	// store (sentinel#44), so asking for threat_count returns the most recent
	// instead of the busiest. A page of 10 would therefore be "the last 10 addresses seen",
	// which is a different list wearing this one's title. Read a wide page and
	// rank here.
	var actors sentinelActors
	get("active_ips", "/sentinel/api/actors?page_size=200", &actors)

	if up.allFailed(3) {
		c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{
			"code": "SENTINEL_UNAVAILABLE", "message": "Sentinel is not answering"}})
		return
	}

	byType := map[string]int64{}
	for _, bucket := range trends.Data {
		for name, n := range bucket.ByType {
			byType[name] += n
		}
	}
	attackTypes := make([]gin.H, 0, len(byType))
	for name, n := range byType {
		attackTypes = append(attackTypes, gin.H{"type": name, "count": n})
	}
	sort.Slice(attackTypes, func(i, j int) bool {
		a, b := attackTypes[i]["count"].(int64), attackTypes[j]["count"].(int64)
		if a != b {
			return a > b
		}
		return attackTypes[i]["type"].(string) < attackTypes[j]["type"].(string)
	})
	routes := make([]gin.H, 0, len(targets.Data))
	for _, t := range targets.Data {
		routes = append(routes, gin.H{"route": t.Route, "method": t.Method, "count": t.Count})
	}
	ranked := actors.Data
	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].ThreatCount != ranked[j].ThreatCount {
			return ranked[i].ThreatCount > ranked[j].ThreatCount
		}
		return ranked[i].LastSeen.After(ranked[j].LastSeen)
	})
	if len(ranked) > 10 {
		ranked = ranked[:10]
	}
	ips := make([]gin.H, 0, len(ranked))
	for _, a := range ranked {
		ips = append(ips, gin.H{
			"ip": a.IP, "count": a.ThreatCount, "risk": a.RiskScore,
			"country": a.Country, "status": a.Status,
			"last_seen": a.LastSeen.Format(time.RFC3339),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"attack_types": attackTypes,
		"routes":       routes,
		"ips":          ips,
		"window":       window,
		"degraded":     up.names(),
	})
}

// Threats is the log itself: filtered by severity, searched by address or
// path, and paged, because 18,000 rows is a normal number and a page that
// shows the newest ten of them cannot be acted on.
func (h *SecurityHandler) Threats(c *gin.Context) {
	if h.Bridge == nil {
		respond.Fail(c, respond.CodeSentinelOff, "Sentinel is not enabled")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	query := fmt.Sprintf("/sentinel/api/threats?page=%d&page_size=%d", page, pageSize)
	if severity := strings.TrimSpace(c.Query("severity")); severity != "" && !strings.EqualFold(severity, "all") {
		// Sentinel filters on the exact stored word, which is capitalised.
		canonical, ok := services.CanonicalSeverity(severity)
		if !ok {
			respond.Fail(c, respond.CodeValidationError, "severity must be critical, high, medium, low or all")
			return
		}
		query += "&severity=" + canonical
	}
	if search := strings.TrimSpace(c.Query("search")); search != "" {
		query += "&search=" + url.QueryEscape(search)
	}

	var got sentinelThreatPage
	if err := h.Bridge.SentinelGet(c.Request.Context(), query, &got); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{
			"code": "SENTINEL_UNAVAILABLE", "message": "Sentinel is not answering"}})
		return
	}

	rows := make([]gin.H, 0, len(got.Data))
	for _, t := range got.Data {
		threatType := ""
		if len(t.ThreatTypes) > 0 {
			threatType = t.ThreatTypes[0]
		}
		outcome := "Logged"
		if t.Blocked {
			outcome = "Blocked"
		}
		rows = append(rows, gin.H{
			"id":        t.ID,
			"ip":        t.IP,
			"method":    t.Method,
			"path":      t.Path,
			"type":      threatType,
			"severity":  t.Severity,
			"blocked":   t.Blocked,
			"outcome":   outcome,
			"timestamp": t.Timestamp.Format(time.RFC3339),
		})
	}
	total := got.Meta.Total
	pages := 1
	if pageSize > 0 && total > 0 {
		pages = (total + pageSize - 1) / pageSize
	}
	c.JSON(http.StatusOK, gin.H{
		"data": rows,
		"meta": gin.H{"total": total, "page": page, "page_size": pageSize, "pages": pages},
	})
}

// banRequest is one address, by hand.
type banRequest struct {
	IP        string ` + "`json:\"ip\" binding:\"required\"`" + `
	Reason    string ` + "`json:\"reason\"`" + `
	Permanent bool   ` + "`json:\"permanent\"`" + `
	Duration  string ` + "`json:\"duration\"`" + `
}

// Ban blocks one address.
func (h *SecurityHandler) Ban(c *gin.Context) {
	if h.Bridge == nil {
		respond.Fail(c, respond.CodeSentinelOff, "Sentinel is not enabled")
		return
	}
	var req banRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, "An IP address is required")
		return
	}
	req.IP = strings.TrimSpace(req.IP)
	if net.ParseIP(req.IP) == nil {
		respond.Fail(c, respond.CodeValidationError, req.IP+" is not an IP address")
		return
	}
	// The app's own infrastructure speaks from these, and behind a proxy that
	// Sentinel has not been told to trust, so does every one of your users.
	if services.IsPrivateOrLoopbackIP(req.IP) {
		respond.Fail(c, respond.CodeValidationError,
			req.IP+" is a loopback or private address. Behind a reverse proxy every "+
				"request carries one, so banning it would ban everybody. Set "+
				"SENTINEL_TRUSTED_PROXIES so the real client address is read.")
		return
	}
	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		reason = "blocked from the admin dashboard"
	}
	body := map[string]interface{}{"ip": req.IP, "reason": reason}
	if req.Permanent {
		body["permanent"] = true
	} else if req.Duration != "" {
		if _, err := time.ParseDuration(req.Duration); err != nil {
			respond.Fail(c, respond.CodeValidationError, "duration must be a duration such as \"24h\"")
			return
		}
		body["duration"] = req.Duration
	}
	if err := h.Bridge.SentinelPost(c.Request.Context(), "/sentinel/api/ip/block", body, nil); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{
			"code": "SENTINEL_UNAVAILABLE", "message": "Sentinel did not accept the ban"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": req.IP + " is banned", "data": gin.H{"ip": req.IP}})
}

// Unban lifts a ban, whoever made it.
func (h *SecurityHandler) Unban(c *gin.Context) {
	if h.Bridge == nil {
		respond.Fail(c, respond.CodeSentinelOff, "Sentinel is not enabled")
		return
	}
	ip := strings.TrimSpace(c.Param("ip"))
	if net.ParseIP(ip) == nil {
		respond.Fail(c, respond.CodeValidationError, "That is not an IP address")
		return
	}
	if err := h.Bridge.SentinelDelete(c.Request.Context(), "/sentinel/api/ip/block/"+url.PathEscape(ip), nil); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{
			"code": "SENTINEL_UNAVAILABLE", "message": "Sentinel did not lift the ban"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": ip + " is no longer banned", "data": gin.H{"ip": ip}})
}

// BanOffenders runs the sweep now, on the same policy the schedule uses.
//
// The button exists because the schedule only sees what arrives after it is
// turned on, and the log it is turned on in front of is usually already full.
func (h *SecurityHandler) BanOffenders(c *gin.Context) {
	if h.AutoBan == nil {
		respond.Fail(c, respond.CodeSentinelOff, "Sentinel is not enabled")
		return
	}
	var req struct {
		DryRun bool ` + "`json:\"dry_run\"`" + `
	}
	_ = c.ShouldBindJSON(&req)

	result, err := h.AutoBan.Run(c.Request.Context(), req.DryRun)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{
			"code": "SENTINEL_UNAVAILABLE", "message": err.Error()}})
		return
	}
	message := "No address is over the threshold"
	switch {
	case req.DryRun && len(result.Banned) > 0:
		message = fmt.Sprintf("%d address(es) would be banned", len(result.Banned))
	case len(result.Banned) > 0:
		message = fmt.Sprintf("Banned %d high-threat IP(s). Enforcement starts within about 30 seconds.", len(result.Banned))
	}
	c.JSON(http.StatusOK, gin.H{"message": message, "data": result})
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
	"{{MODULE}}/internal/respond"
)

type ObservabilityHandler struct {
	Bridge *services.SecObsBridge
}

// pulseOverview is /pulse/api/overview: pulse.Overview plus pulse.RouteStats.
// Latencies are nanoseconds, Go's time.Duration JSON form.
type pulseOverview struct {
	TotalRequests int64   ` + "`json:\"total_requests\"`" + `
	TotalErrors   int64   ` + "`json:\"total_errors\"`" + `
	ErrorRate     float64 ` + "`json:\"error_rate\"`" + `
	AvgLatency    int64   ` + "`json:\"avg_latency\"`" + `
	P95Latency    int64   ` + "`json:\"p95_latency\"`" + `
	RPM           float64 ` + "`json:\"rpm\"`" + `
	TopRoutes     []struct {
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

// pulseRuntime is /pulse/api/runtime/current: a RuntimeMetric snapshot.
type pulseRuntime struct {
	HeapAlloc    uint64 ` + "`json:\"heap_alloc\"`" + `
	NumGoroutine int    ` + "`json:\"num_goroutine\"`" + `
	NumGC        uint32 ` + "`json:\"num_gc\"`" + `
}

// pulseN1 is /pulse/api/database/n1/ranked: N1Ranking[] under "data".
type pulseN1 struct {
	Data []struct {
		Route            string  ` + "`json:\"route\"`" + `
		AvgQueriesPerHit float64 ` + "`json:\"avg_queries_per_hit\"`" + `
		FirstSeen        string  ` + "`json:\"first_seen\"`" + `
	} ` + "`json:\"data\"`" + `
}

// pulseErrors is /pulse/api/errors: ErrorRecord[] under "data".
type pulseErrors struct {
	Data []struct {
		ID           string ` + "`json:\"id\"`" + `
		Route        string ` + "`json:\"route\"`" + `
		ErrorMessage string ` + "`json:\"error_message\"`" + `
		LastSeen     string ` + "`json:\"last_seen\"`" + `
	} ` + "`json:\"data\"`" + `
}

// Pulse durations come back as nanoseconds (Go's time.Duration JSON form).
func nsToMs(ns int64) float64 { return float64(ns) / 1_000_000.0 }

// Summary returns a flat performance envelope built from Pulse's overview,
// runtime, ranked N+1 and errors endpoints. Missing or zero data is fine
// — the React page renders "—" for unset numeric fields and shows the
// "No X yet" empty state for empty arrays.
func (h *ObservabilityHandler) Summary(c *gin.Context) {
	if h.Bridge == nil {
		respond.Fail(c, respond.CodePulseOff, "Pulse is not enabled")
		return
	}

	ctx := c.Request.Context()
	var up upstream

	get := func(name, path string, out interface{}) {
		up.get(name, func() error { return h.Bridge.PulseGet(ctx, path, out) })
	}

	var overview pulseOverview
	get("overview", "/pulse/api/overview?range=1h", &overview)

	var rt pulseRuntime
	get("runtime", "/pulse/api/runtime/current", &rt)

	var n1Resp pulseN1
	get("n1_detections", "/pulse/api/database/n1/ranked?range=1h&limit=10", &n1Resp)

	var errResp pulseErrors
	get("recent_errors", "/pulse/api/errors?limit=10&resolved=false", &errResp)

	// Nothing answered. A page of noughts would say this app served no
	// requests, had no errors and ran at nought milliseconds.
	if up.allFailed(4) {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": gin.H{
				"code":    "PULSE_UNAVAILABLE",
				"message": "Pulse is not answering",
				"details": gin.H{"degraded": up.names()},
			},
		})
		return
	}

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
		// The panels whose numbers are missing rather than zero.
		"degraded": up.names(),
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
				Title:    fmt.Sprintf("Threat: %s (CVSS %.1f)", t.Type, t.CVSS),
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
