package services

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
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
	Severity     string // info | warn | critical
	Summary      string
	ResourceType string
	ResourceID   string
	Metadata     map[string]interface{} // optional JSON-encodable extras
}

// LogActivity writes a UserActivity row before it returns. Picks actor + IP +
// user-agent from the request context automatically, falling back to
// args.UserID when the caller is in an unauthenticated handler (auth flows).
func LogActivity(db *gorm.DB, c *gin.Context, args ActivityArgs) {
	if err := LogActivityErr(db, c, args); err != nil {
		// Audit failures are non-fatal but worth knowing about: log and keep
		// moving. In production, wire a metric to alert on a sudden surge in
		// these (suggests DB write pressure).
		log.Printf("activity: failed to write %s: %v", args.Action, err)
	}
}

// LogActivityErr is LogActivity for the handful of events where losing the row
// matters more than the noise: a GDPR erasure, a key revocation, anything an
// auditor will later ask you to produce. It writes the same row and hands back
// the error instead of logging it, so the caller can tell the client the audit
// trail is incomplete.
func LogActivityErr(db *gorm.DB, c *gin.Context, args ActivityArgs) error {
	row := activityRow(c, args)
	return db.Create(&row).Error
}

// activityRow builds the row for args from the request: the actor, the client
// IP and the user agent. It reads the request, so it runs before the handler
// returns even when the row is written later.
//
// The building itself is activityRowCtx, in request_meta.go, which takes those
// three facts from a context.Context instead. This is the gin-shaped front door
// to it, and the only reason it still exists is that handlers hold a gin
// context and jobs do not.
func activityRow(c *gin.Context, args ActivityArgs) models.UserActivity {
	return activityRowCtx(ContextOf(c), args)
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
		UserID:       userID,
		Action:       "auth.logout",
		Severity:     "info",
		Summary:      email + " signed out",
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
// `identifier` is the human-readable label (name, title, slug, sku);
// it must never be blank -- callers should fall back to a snippet of
// the ID rather than emit "Created Product ". `detail` is optional
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

// hasActor reports whether the request carries an authenticated identity.
//
// The resource helpers below record "who changed what". With no actor there is
// no who, so the row answers nothing anyone can audit — and on a public
// endpoint it is worse than useless: every anonymous write turns into two
// inserts in two transactions, which is write amplification an attacker
// controls for free.
//
// Auth events are deliberately different. LogLoginFailed records an empty
// actor on purpose, because a failed sign-in matters precisely when the
// account is unknown. Only the CRUD helpers guard.
func hasActor(c *gin.Context) bool {
	v, ok := c.Get("user_id")
	if !ok {
		return false
	}
	s, ok := v.(string)
	return ok && s != ""
}

// LogCreate writes a "Created X Y[: detail]" row.
func LogCreate(db *gorm.DB, c *gin.Context, entityType, identifier, resourceID, detail string) {
	if !hasActor(c) {
		return
	}
	queueActivity(db, c, ActivityArgs{
		Action:       strings.ToLower(entityType) + ".create",
		Severity:     "info",
		Summary:      formatCUDSummary("Created", entityType, identifier, detail),
		ResourceType: strings.ToLower(entityType),
		ResourceID:   resourceID,
	})
}

// LogUpdate writes an "Updated X Y[: detail]" row. `detail` is the
// caller-built diff string (e.g. `name "old" → "new"`); pass "" if
// you only want to record that a touch happened.
func LogUpdate(db *gorm.DB, c *gin.Context, entityType, identifier, resourceID, detail string) {
	if !hasActor(c) {
		return
	}
	queueActivity(db, c, ActivityArgs{
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
	if !hasActor(c) {
		return
	}
	queueActivity(db, c, ActivityArgs{
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
//	1 field   → `image changed`
//	2-3 fields → `changed name, image, price`
//	4+ fields → `5 fields changed (name, image, price, ...)`
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
