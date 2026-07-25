package scaffold

// OCSF audit export.
//
// The semantic activity log (models.UserActivity) records who did what —
// auth.login, user.delete, session.revoke_all — with actor, severity, resource
// and IP. That is exactly what a SIEM ingests, but no SIEM speaks Grit's schema.
// The Open Cybersecurity Schema Framework (ocsf.io) is the vendor-neutral shape
// Splunk, Elastic, Sentinel, Chronicle and Amazon Security Lake all understand.
//
// This maps each event to its OCSF class and exposes them as newline-delimited
// JSON on an admin endpoint a collector can poll incrementally. Pull rather than
// push: no credentials to store, no queue to babysit, and the collector owns its
// own cursor — the pattern every one of those SIEMs already ships a connector
// for.

func apiOCSFServiceGo() string {
	return `package services

import (
	"strings"
	"time"

	"{{MODULE}}/internal/models"
)

// OCSF schema version this exporter conforms to. Bump deliberately — a SIEM's
// parser is pinned to a major line.
const OCSFVersion = "1.3.0"

// ocsfClass names the OCSF class an event maps to. category/class/activity are
// the OCSF taxonomy; typeUID is the derived class_uid*100 + activity_id that
// OCSF uses as the stable event-type key.
type ocsfClass struct {
	CategoryUID  int
	CategoryName string
	ClassUID     int
	ClassName    string
	ActivityID   int
	ActivityName string
}

func (c ocsfClass) typeUID() int { return c.ClassUID*100 + c.ActivityID }

// OCSF category + class constants (subset actually emitted).
const (
	catIAM = 3 // Identity & Access Management
	catApp = 6 // Application Activity

	classAuthentication = 3002
	classAccountChange  = 3001
	classAPIActivity    = 6003
)

// classifyAction maps a Grit dotted action to an OCSF class. The rule is:
// auth/session/password/user/role events are IAM (a security team filters on
// these); everything else is generic API Activity keyed off the verb. Unknown
// verbs still map — to API Activity with activity_id 0 (Unknown) — so a new
// action name is never silently dropped from the SIEM.
func classifyAction(action string) ocsfClass {
	switch action {
	case "auth.login":
		return ocsfClass{catIAM, "Identity & Access Management", classAuthentication, "Authentication", 1, "Logon"}
	case "auth.login_failed":
		return ocsfClass{catIAM, "Identity & Access Management", classAuthentication, "Authentication", 1, "Logon"}
	case "auth.logout", "session.revoke", "session.revoke_all":
		return ocsfClass{catIAM, "Identity & Access Management", classAuthentication, "Authentication", 2, "Logoff"}
	case "auth.register", "user.create":
		return ocsfClass{catIAM, "Identity & Access Management", classAccountChange, "Account Change", 1, "Create"}
	case "password.change":
		return ocsfClass{catIAM, "Identity & Access Management", classAccountChange, "Account Change", 3, "Password Change"}
	case "password.reset", "auth.reset_password":
		return ocsfClass{catIAM, "Identity & Access Management", classAccountChange, "Account Change", 4, "Password Reset"}
	case "user.delete":
		return ocsfClass{catIAM, "Identity & Access Management", classAccountChange, "Account Change", 6, "Delete"}
	case "user.update", "role.assign", "role.revoke", "permission.grant", "permission.revoke":
		return ocsfClass{catIAM, "Identity & Access Management", classAccountChange, "Account Change", 13, "Change"}
	}

	// Generic <entity>.<verb> — map the verb onto API Activity.
	verb := action
	if i := strings.LastIndex(action, "."); i >= 0 {
		verb = action[i+1:]
	}
	switch verb {
	case "create":
		return ocsfClass{catApp, "Application Activity", classAPIActivity, "API Activity", 1, "Create"}
	case "read", "view", "list", "export":
		return ocsfClass{catApp, "Application Activity", classAPIActivity, "API Activity", 2, "Read"}
	case "update", "edit":
		return ocsfClass{catApp, "Application Activity", classAPIActivity, "API Activity", 3, "Update"}
	case "delete", "remove":
		return ocsfClass{catApp, "Application Activity", classAPIActivity, "API Activity", 4, "Delete"}
	default:
		return ocsfClass{catApp, "Application Activity", classAPIActivity, "API Activity", 0, "Unknown"}
	}
}

// severityID maps Grit's severity to the OCSF severity_id scale
// (1 Informational, 3 Medium, 5 Critical).
func severityID(sev string) (int, string) {
	switch sev {
	case "critical":
		return 5, "Critical"
	case "warn":
		return 3, "Medium"
	default:
		return 1, "Informational"
	}
}

// ToOCSF renders one activity row as an OCSF event object, ready to marshal.
// The map form (rather than a typed struct) keeps the OCSF shape — deeply
// nested and evolving — readable and lets optional fields be omitted cleanly.
func ToOCSF(a models.UserActivity, productName string) map[string]interface{} {
	class := classifyAction(a.Action)

	// status: a failed sign-in is the one action whose name encodes failure.
	statusID, status := 1, "Success"
	if a.Action == "auth.login_failed" {
		statusID, status = 2, "Failure"
	}

	sevID, sev := severityID(a.Severity)

	ev := map[string]interface{}{
		"activity_id":   class.ActivityID,
		"activity_name": class.ActivityName,
		"category_uid":  class.CategoryUID,
		"category_name": class.CategoryName,
		"class_uid":     class.ClassUID,
		"class_name":    class.ClassName,
		"type_uid":      class.typeUID(),
		"time":          a.CreatedAt.UnixMilli(),
		"severity_id":   sevID,
		"severity":      sev,
		"status_id":     statusID,
		"status":        status,
		"message":       a.Summary,
		"metadata": map[string]interface{}{
			"version": OCSFVersion,
			"uid":     a.ID,
			"product": map[string]interface{}{
				"name":        productName,
				"vendor_name": "Grit",
			},
		},
		// Grit's own action string is preserved so a SIEM can pivot back to the
		// native taxonomy without reverse-engineering the OCSF mapping.
		"unmapped": map[string]interface{}{
			"grit_action":        a.Action,
			"grit_resource_type": a.ResourceType,
			"grit_resource_id":   a.ResourceID,
		},
	}

	// actor.user — omit entirely for system events rather than emit a blank uid,
	// which a SIEM would read as "an account whose id is empty string".
	if a.UserID != "" {
		ev["actor"] = map[string]interface{}{
			"user": map[string]interface{}{"uid": a.UserID},
		}
	}
	if a.IPAddress != "" {
		ev["src_endpoint"] = map[string]interface{}{"ip": a.IPAddress}
	}
	if a.UserAgent != "" {
		ev["http_request"] = map[string]interface{}{"user_agent": a.UserAgent}
	}
	if a.ResourceType != "" || a.ResourceID != "" {
		ev["resources"] = []map[string]interface{}{
			{"type": a.ResourceType, "uid": a.ResourceID},
		}
	}
	return ev
}

// ExportCursor points at the last row a collector has already consumed. OCSF
// events carry millisecond time, but two rows can share a millisecond, so the
// cursor is (created_at, id) — the same total order VerifyChain uses — to
// guarantee no row is skipped or repeated across polls.
type ExportCursor struct {
	After   time.Time
	AfterID string
}
`
}

func apiOCSFHandlerGo() string {
	return `package handlers

import (
	"bufio"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
)

// OCSFHandler streams the semantic activity log as OCSF events for a SIEM to
// ingest. Admin-only — an audit trail readable by the people it audits is not
// an audit trail.
type OCSFHandler struct {
	DB          *gorm.DB
	ProductName string
}

func NewOCSFHandler(db *gorm.DB, productName string) *OCSFHandler {
	if productName == "" {
		productName = "grit-app"
	}
	return &OCSFHandler{DB: db, ProductName: productName}
}

const (
	ocsfDefaultLimit = 500
	ocsfMaxLimit     = 5000
)

// Export streams OCSF events as newline-delimited JSON (application/x-ndjson),
// oldest first, so a collector ingests chronologically and never has to hold
// the whole response in memory.
//
//	GET /api/audit/ocsf?since=2026-07-01T00:00:00Z&after=<id>&limit=1000
//
// Pagination is a cursor, not an offset: pass the response's X-Grit-Next-Since
// and X-Grit-Next-After headers back on the next poll to resume exactly where
// this response stopped. When fewer than limit rows come back, the collector is
// caught up.
func (h *OCSFHandler) Export(c *gin.Context) {
	limit := ocsfDefaultLimit
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	if limit > ocsfMaxLimit {
		limit = ocsfMaxLimit
	}

	q := h.DB.Model(&models.UserActivity{}).Order("created_at asc, id asc").Limit(limit)

	// since is a wall-clock floor (a collector's first poll); after is the exact
	// cursor for every poll thereafter. Both may be present — after wins ties
	// within the same millisecond as since.
	var since time.Time
	if v := c.Query("since"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{"code": "VALIDATION_ERROR", "message": "since must be RFC3339, e.g. 2026-07-01T00:00:00Z"},
			})
			return
		}
		since = t
		q = q.Where("created_at >= ?", t)
	}
	if afterID := c.Query("after"); afterID != "" {
		var cursor models.UserActivity
		if err := h.DB.Select("created_at", "id").First(&cursor, "id = ?", afterID).Error; err == nil {
			q = q.Where("(created_at, id) > (?, ?)", cursor.CreatedAt, cursor.ID)
		}
		// An unknown after id falls through to since/start rather than erroring:
		// a collector that lost its place still makes progress instead of
		// wedging.
		_ = since
	}

	var rows []models.UserActivity
	if err := q.Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "failed to read audit log"},
		})
		return
	}

	c.Header("Content-Type", "application/x-ndjson")
	// Advertise where the next poll should resume. Empty when this page was
	// empty — the collector then just retries the same cursor later.
	if n := len(rows); n > 0 {
		last := rows[n-1]
		c.Header("X-Grit-Next-Since", last.CreatedAt.UTC().Format(time.RFC3339Nano))
		c.Header("X-Grit-Next-After", last.ID)
	}
	c.Header("X-Grit-Count", strconv.Itoa(len(rows)))

	w := bufio.NewWriter(c.Writer)
	defer w.Flush()
	enc := json.NewEncoder(w)
	for i := range rows {
		// json.Encoder writes a trailing newline per value, which is exactly the
		// NDJSON record separator.
		if err := enc.Encode(services.ToOCSF(rows[i], h.ProductName)); err != nil {
			return
		}
	}
}
`
}

func apiOCSFTestGo() string {
	return `package services_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
)

func TestToOCSFAuthentication(t *testing.T) {
	ev := services.ToOCSF(models.UserActivity{
		ID:        "act-1",
		UserID:    "user-9",
		Action:    "auth.login",
		Severity:  "info",
		Summary:   "ada@example.com signed in",
		IPAddress: "203.0.113.7",
		CreatedAt: time.Unix(1_700_000_000, 0),
	}, "test-app")

	assert.Equal(t, 3002, ev["class_uid"], "auth.login -> Authentication class")
	assert.Equal(t, 1, ev["activity_id"], "Logon activity")
	assert.Equal(t, 300201, ev["type_uid"], "class_uid*100 + activity_id")
	assert.Equal(t, 1, ev["status_id"], "success")
	assert.Equal(t, int64(1_700_000_000_000), ev["time"], "epoch millis")

	actor := ev["actor"].(map[string]interface{})
	assert.Equal(t, "user-9", actor["user"].(map[string]interface{})["uid"])
	src := ev["src_endpoint"].(map[string]interface{})
	assert.Equal(t, "203.0.113.7", src["ip"])
}

// A failed sign-in is the one action whose OCSF status must be Failure — that
// is what a SIEM alerts on.
func TestToOCSFFailedLogin(t *testing.T) {
	ev := services.ToOCSF(models.UserActivity{
		Action:   "auth.login_failed",
		Severity: "warn",
	}, "test-app")

	assert.Equal(t, 3002, ev["class_uid"])
	assert.Equal(t, 2, ev["status_id"], "Failure")
	assert.Equal(t, "Failure", ev["status"])
	assert.Equal(t, 3, ev["severity_id"], "warn -> Medium")
}

func TestToOCSFAccountChange(t *testing.T) {
	cases := map[string]int{
		"auth.register":  1, // Create
		"password.reset": 4, // Password Reset
		"user.delete":    6, // Delete
		"user.update":    13,
	}
	for action, wantActivity := range cases {
		ev := services.ToOCSF(models.UserActivity{Action: action}, "app")
		assert.Equal(t, 3001, ev["class_uid"], action+" -> Account Change")
		assert.Equal(t, wantActivity, ev["activity_id"], action)
	}
}

// An unknown action must still map — dropping it would blind the SIEM to
// exactly the novel event worth seeing.
func TestToOCSFUnknownActionStillMaps(t *testing.T) {
	ev := services.ToOCSF(models.UserActivity{Action: "invoice.approve"}, "app")
	assert.Equal(t, 6003, ev["class_uid"], "falls back to API Activity")
	assert.Equal(t, 0, ev["activity_id"], "Unknown activity, never dropped")

	ev2 := services.ToOCSF(models.UserActivity{Action: "invoice.create"}, "app")
	assert.Equal(t, 6003, ev2["class_uid"])
	assert.Equal(t, 1, ev2["activity_id"], "verb create -> API Activity Create")
}

// System events have no actor; emitting a blank uid would read as a real
// account to a SIEM.
func TestToOCSFSystemEventHasNoActor(t *testing.T) {
	ev := services.ToOCSF(models.UserActivity{Action: "cron.cleanup", Summary: "nightly purge"}, "app")
	_, hasActor := ev["actor"]
	assert.False(t, hasActor, "no actor key when UserID is empty")
}

func TestToOCSFPreservesNativeAction(t *testing.T) {
	ev := services.ToOCSF(models.UserActivity{
		Action:       "ticket.close",
		ResourceType: "ticket",
		ResourceID:   "t-42",
	}, "app")
	um := ev["unmapped"].(map[string]interface{})
	assert.Equal(t, "ticket.close", um["grit_action"])
	assert.Equal(t, "ticket", um["grit_resource_type"])

	md := ev["metadata"].(map[string]interface{})
	require.Equal(t, services.OCSFVersion, md["version"])
}
`
}
