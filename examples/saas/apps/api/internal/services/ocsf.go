package services

import (
	"strings"
	"time"

	"saas/apps/api/internal/models"
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
