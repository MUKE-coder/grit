package scaffold

// v3.30 — UserActivity (semantic event log, separate from the existing
// tamper-evident ActivityLog hash chain).
//
// Why a separate model:
//   - ActivityLog: per-HTTP-request log, hash chained, append-only,
//     compliance use case. Carries Method + Path + Status + a SHA digest.
//   - UserActivity: semantic events that operators read in the
//     /system/activity dashboard. Carries Action (e.g. "ticket.create"),
//     Severity, Summary, optional ResourceType + ResourceID. No hash
//     chain — soft-delete + retention via cron is fine.
//
// Handlers call activity.Log(db, c, "user.delete", "critical", summary,
// "user", userID) to emit. Auth events use convenience helpers.

// userActivityModelGo emits internal/models/user_activity.go.
func userActivityModelGo() string {
	return `package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserActivity is a semantic event written by handlers via
// activity.Log(...). Surfaced in the /system/activity dashboard so
// operators can see who did what, when, from where, and how severe.
//
// Separate from ActivityLog (the hash-chained HTTP audit) on purpose:
//   - This is human-readable, queryable, filter-friendly.
//   - ActivityLog is tamper-evident, immutable, mostly compliance fodder.
//
// Use the Severity column as the operator-facing impact level:
//   - info     → routine action (login, view, list)
//   - warn     → something noteworthy (failed login, role change)
//   - critical → high-impact action (user delete, billing change, mass
//                update). Surface these prominently in the dashboard.
type UserActivity struct {
	ID           string         ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	UserID       string         ` + "`" + `gorm:"size:36;index" json:"user_id"` + "`" + `      // actor (empty = system)
	Action       string         ` + "`" + `gorm:"size:64;index" json:"action"` + "`" + `       // dotted: "ticket.create", "user.delete"
	Severity     string         ` + "`" + `gorm:"size:16;index" json:"severity"` + "`" + `     // info | warn | critical
	Summary      string         ` + "`" + `gorm:"size:500" json:"summary"` + "`" + `           // single sentence for the dashboard row
	ResourceType string         ` + "`" + `gorm:"size:64;index" json:"resource_type"` + "`" + ` // e.g. "user", "ticket"
	ResourceID   string         ` + "`" + `gorm:"size:64;index" json:"resource_id"` + "`" + `   // primary key of the target
	IPAddress    string         ` + "`" + `gorm:"size:45" json:"ip_address"` + "`" + `
	UserAgent    string         ` + "`" + `gorm:"size:500" json:"user_agent"` + "`" + `
	Metadata     string         ` + "`" + `gorm:"type:text" json:"metadata"` + "`" + `          // optional JSON blob for extra context
	CreatedAt    time.Time      ` + "`" + `gorm:"index" json:"created_at"` + "`" + `
	DeletedAt    gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}

func (a *UserActivity) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	if a.Severity == "" {
		a.Severity = "info"
	}
	return nil
}
`
}

// userActivityServiceGo emits internal/services/activity.go — the helper
// every handler calls to emit a semantic event.
func userActivityServiceGo() string {
	return `package services

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// ActivityService centralises the "log a semantic event" call so
// handlers don't repeat IP-pull / user-agent / actor lookup boilerplate.
//
// Usage from a handler:
//
//	services.LogActivity(h.DB, c, services.ActivityArgs{
//	    Action:       "ticket.create",
//	    Severity:     "info",
//	    Summary:      fmt.Sprintf("Opened ticket %q", ticket.Subject),
//	    ResourceType: "ticket",
//	    ResourceID:   ticket.ID,
//	})
//
// Errors are logged, not returned — losing an audit row should never
// fail a real request. If you care about guaranteed delivery, queue
// these via asynq instead of writing inline.
type ActivityArgs struct {
	// UserID overrides the actor when set. Auth handlers (login, register)
	// pass it explicitly because the auth middleware hasn't yet populated
	// the gin context with "user_id" — without this override every auth
	// event would be attributed to the empty system actor.
	UserID       string
	Action       string
	Severity     string                 // info | warn | critical
	Summary      string
	ResourceType string
	ResourceID   string
	Metadata     map[string]interface{} // optional JSON-encodable extras
}

// LogActivity writes a UserActivity row. Picks actor + IP + user-agent
// from the request context automatically, falling back to args.UserID
// when the caller is in an unauthenticated handler (auth flows).
func LogActivity(db *gorm.DB, c *gin.Context, args ActivityArgs) {
	userID := args.UserID
	if userID == "" {
		if v, ok := c.Get("user_id"); ok {
			if s, ok := v.(string); ok {
				userID = s
			}
		}
	}

	var metaJSON string
	if args.Metadata != nil {
		if b, err := json.Marshal(args.Metadata); err == nil {
			metaJSON = string(b)
		}
	}

	row := models.UserActivity{
		UserID:       userID,
		Action:       args.Action,
		Severity:     args.Severity,
		Summary:      args.Summary,
		ResourceType: args.ResourceType,
		ResourceID:   args.ResourceID,
		IPAddress:    ResolveClientIP(c),
		UserAgent:    c.GetHeader("User-Agent"),
		Metadata:     metaJSON,
	}

	if err := db.Create(&row).Error; err != nil {
		// Audit failures are non-fatal but worth knowing about — log and
		// keep moving. In production, wire a metric to alert on a sudden
		// surge in these (suggests DB write pressure).
		log.Printf("activity: failed to write %s: %v", args.Action, err)
	}
}

// Convenience helpers for the most common events. Use these in auth
// handlers + middleware so the dotted action names stay consistent.

func LogLogin(db *gorm.DB, c *gin.Context, userID, email string) {
	LogActivity(db, c, ActivityArgs{
		UserID:       userID,
		Action:       "auth.login",
		Severity:     "info",
		Summary:      email + " signed in",
		ResourceType: "user",
		ResourceID:   userID,
	})
}

// LogLoginFailed intentionally leaves UserID empty — the failed attempt
// might be an unknown email or a real account being brute-forced; either
// way, attributing it to a specific actor is misleading. The summary
// captures the attempted email for audit.
func LogLoginFailed(db *gorm.DB, c *gin.Context, email string) {
	LogActivity(db, c, ActivityArgs{
		Action:   "auth.login_failed",
		Severity: "warn",
		Summary:  "Failed sign-in attempt for " + email,
	})
}

func LogLogout(db *gorm.DB, c *gin.Context, userID, email string) {
	LogActivity(db, c, ActivityArgs{
		UserID:   userID,
		Action:   "auth.logout",
		Severity: "info",
		Summary:  email + " signed out",
		ResourceType: "user", ResourceID: userID,
	})
}

func LogRegister(db *gorm.DB, c *gin.Context, userID, email string) {
	LogActivity(db, c, ActivityArgs{
		UserID:       userID,
		Action:       "auth.register",
		Severity:     "info",
		Summary:      email + " created an account",
		ResourceType: "user",
		ResourceID:   userID,
	})
}

// v3.31.39 — CUD helpers. Format convention going forward:
//
//	{verb} {entityType} {identifier}[: {detail}]
//
// ` + "`identifier`" + ` is the human-readable label (name, title, slug, sku);
// it must never be blank -- callers should fall back to a snippet of
// the ID rather than emit "Created Product ". ` + "`detail`" + ` is optional
// extra context (price for Create, diff for Update) and only renders
// when non-empty.
//
// Example summaries:
//
//	Created Product Desktop: KES 340,000
//	Updated Product Desktop: name "Desktop" → "Desktop Pro"
//	Updated Category Phones: image changed
//	Deleted Blog "Welcome to the new site"
//
// Severity is fixed at "info" for routine Creates/Updates/Deletes.
// If a particular delete should pop louder (user accounts, billing
// rows), call LogActivity directly with Severity: "critical".

func formatCUDSummary(verb, entityType, identifier, detail string) string {
	if identifier == "" {
		identifier = "(unnamed)"
	}
	summary := verb + " " + entityType + " " + identifier
	if detail != "" {
		summary += ": " + detail
	}
	return summary
}

// LogCreate writes a "Created X Y[: detail]" row.
func LogCreate(db *gorm.DB, c *gin.Context, entityType, identifier, resourceID, detail string) {
	LogActivity(db, c, ActivityArgs{
		Action:       strings.ToLower(entityType) + ".create",
		Severity:     "info",
		Summary:      formatCUDSummary("Created", entityType, identifier, detail),
		ResourceType: strings.ToLower(entityType),
		ResourceID:   resourceID,
	})
}

// LogUpdate writes an "Updated X Y[: detail]" row. ` + "`detail`" + ` is the
// caller-built diff string (e.g. ` + "`name \"old\" → \"new\"`" + `); pass "" if
// you only want to record that a touch happened.
func LogUpdate(db *gorm.DB, c *gin.Context, entityType, identifier, resourceID, detail string) {
	LogActivity(db, c, ActivityArgs{
		Action:       strings.ToLower(entityType) + ".update",
		Severity:     "info",
		Summary:      formatCUDSummary("Updated", entityType, identifier, detail),
		ResourceType: strings.ToLower(entityType),
		ResourceID:   resourceID,
	})
}

// LogDelete writes a "Deleted X Y" row. No detail -- by the time you
// log a delete the snippet is the only thing left.
func LogDelete(db *gorm.DB, c *gin.Context, entityType, identifier, resourceID string) {
	LogActivity(db, c, ActivityArgs{
		Action:       strings.ToLower(entityType) + ".delete",
		Severity:     "info",
		Summary:      formatCUDSummary("Deleted", entityType, identifier, ""),
		ResourceType: strings.ToLower(entityType),
		ResourceID:   resourceID,
	})
}

// DiffSummary renders an Updates() map as a human-readable
// "what changed" string for the UserActivity Summary field:
//
//	1 field   → ` + "`image changed`" + `
//	2-3 fields → ` + "`changed name, image, price`" + `
//	4+ fields → ` + "`5 fields changed (name, image, price, ...)`" + `
//
// Keys are sorted so the output is deterministic across runs (handy
// for tests + log grep). Empty map returns "".
func DiffSummary(updates map[string]interface{}) string {
	if len(updates) == 0 {
		return ""
	}
	keys := make([]string, 0, len(updates))
	for k := range updates {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if len(keys) == 1 {
		return keys[0] + " changed"
	}
	if len(keys) <= 3 {
		return "changed " + strings.Join(keys, ", ")
	}
	return fmt.Sprintf("%d fields changed (%s, ...)", len(keys), strings.Join(keys[:3], ", "))
}
`
}

// userActivityHandlerGo emits internal/handlers/user_activity_handler.go.
func userActivityHandlerGo() string {
	return `package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/paginate"
)

type UserActivityHandler struct {
	DB *gorm.DB
}

// List returns paginated activity events. Filters: user_id, action,
// severity, resource_type, q (substring against summary). Default sort
// is newest first.
//
//	GET /api/user-activity?severity=critical&page=1&page_size=25
//	GET /api/user-activity?q=login&user_id=...
func (h *UserActivityHandler) List(c *gin.Context) {
	q := h.DB.Model(&models.UserActivity{}).Order("created_at desc")

	params := paginate.Bind(c).
		With("user_id", c.Query("user_id")).
		With("action", c.Query("action")).
		With("severity", c.Query("severity")).
		With("resource_type", c.Query("resource_type"))

	if needle := c.Query("q"); needle != "" {
		q = q.Where("summary LIKE ?", "%"+needle+"%")
	}

	res, err := paginate.List[models.UserActivity](q, params, paginate.Config{
		Sortable: map[string]bool{
			"created_at": true,
			"severity":   true,
			"action":     true,
		},
		DefaultSort:  "created_at",
		DefaultOrder: "desc",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

// Stats returns severity bucket counts for the last 24h. Powers the
// header chips on the activity dashboard.
//
//	GET /api/user-activity/stats
//	→ { "data": { "info": 142, "warn": 8, "critical": 1, "total": 151 } }
func (h *UserActivityHandler) Stats(c *gin.Context) {
	type bucket struct {
		Severity string ` + "`" + `json:"severity"` + "`" + `
		Count    int64  ` + "`" + `json:"count"` + "`" + `
	}
	var rows []bucket
	h.DB.Model(&models.UserActivity{}).
		Select("severity, COUNT(*) AS count").
		Where("created_at > NOW() - INTERVAL '24 hours'").
		Group("severity").
		Scan(&rows)

	out := map[string]int64{"info": 0, "warn": 0, "critical": 0, "total": 0}
	for _, r := range rows {
		out[r.Severity] = r.Count
		out["total"] += r.Count
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}
`
}

// v3.31.49 -- clientIPHelperGo emits internal/services/clientip.go.
// ResolveClientIP wraps gin.Context.ClientIP and honours the
// X-Public-IP-Hint header sent by the admin / web axios clients
// when the TCP peer is loopback (::1 / 127.0.0.1). That way local
// dev activity logs show the operator's actual public IP instead of
// "::1" on every row. Production traffic from real proxies goes
// through gin's trusted X-Forwarded-For path and never honours the
// hint -- so it can't be used to spoof audit records.
func clientIPHelperGo() string {
	return `package services

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// v3.31.49 -- ResolveClientIP returns the IP we should log for an
// activity event. It honours the browser-supplied X-Public-IP-Hint
// header, but only when the TCP peer is loopback -- so production
// traffic from real proxies (which sets X-Forwarded-For for
// gin.Context.ClientIP to consume) keeps using the trusted path
// and can't be spoofed by a client header.
//
// The hint exists for one reason: in local dev the operator runs
// admin + API on the same machine, so c.ClientIP() returns "::1"
// for every action. The activity feed showing "::1" on every row
// is technically correct but useless. The admin client looks up
// the operator's public IP once per session (via api.ipify.org)
// and attaches it as X-Public-IP-Hint; this function reads it.
func ResolveClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	if isLoopback(ip) {
		if hint := strings.TrimSpace(c.GetHeader("X-Public-IP-Hint")); hint != "" {
			// Cap at 64 chars matching the column width so a
			// malformed header can't bloat the audit row.
			if len(hint) > 64 {
				hint = hint[:64]
			}
			return hint
		}
	}
	return ip
}

// isLoopback covers the two loopback forms gin.Context.ClientIP can
// return on a local-only deployment: "::1" (IPv6) and "127.0.0.1"
// (IPv4). "0.0.0.0" shouldn't appear as a client IP but we include
// it for symmetry with the frontend prettyIP helper.
func isLoopback(ip string) bool {
	return ip == "::1" || ip == "127.0.0.1" || ip == "0.0.0.0"
}
`
}
