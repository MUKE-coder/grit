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
	"log"

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
		IPAddress:    c.ClientIP(),
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
