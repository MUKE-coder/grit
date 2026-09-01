package generate

import (
	"fmt"
	"sort"
	"strings"
)

// Model names the scaffold already owns.
//
// This list exists because of a real, silent, build-breaking failure:
// `grit generate resource Ticket` wrote its own internal/models/ticket.go over
// the support-desk model the scaffold ships, taking TicketReply with it. The
// generator reported success. The next `go build` failed with an undefined
// symbol in an unrelated file, and nothing connected the two.
//
// A resource is not just its model — generating one also injects routes,
// handler wiring, permissions and a resource definition, all keyed on the same
// name. Colliding with a framework model corrupts every one of those, and
// `grit remove resource` then finishes the job by deleting the originals.
//
// Grouped by what they belong to, so the error message can say *why* a name is
// taken rather than only that it is.
var reservedModels = map[string]string{
	// Core identity and storage.
	"User":                   "authentication",
	"Session":                "authentication",
	"Setting":                "the settings registry",
	"UserIdentity":           "authentication",
	"UserRole":               "authorization",
	"Role":                   "authorization",
	"Upload":                 "file storage",
	"PasswordResetToken":     "authentication",
	"EmailVerificationToken": "authentication",
	"RecoveryContact":        "authentication",
	"RecoveryContactToken":   "authentication",

	// Two-factor and SSO.
	"TwoFactorConfig":  "two-factor auth",
	"TOTPPendingToken": "two-factor auth",
	"TrustedDevice":    "two-factor auth",
	"SSOConnection":    "SSO",
	"SAMLKeypair":      "SSO",
	"APIKey":           "API keys",

	// Built-in admin features.
	"Ticket":           "the support desk",
	"TicketReply":      "the support desk",
	"Notification":     "notifications",
	"ActivityLog":      "the audit log",
	"UserActivity":     "the audit log",
	"AccessReview":     "access reviews",
	"AccessReviewItem": "access reviews",
	"DeletionJournal":  "GDPR erasure",
	"FeatureFlag":      "feature flags",
	"FlagExposure":     "feature flags",
	"FlagRules":        "feature flags",
	"FormShare":        "public form sharing",
	"FormSubmission":   "public form sharing",
	"DashboardLayout":  "dashboard customisation",
	"ImportJob":        "CSV import",
	"Backup":           "backups",
	"BackupSchedule":   "backups",
	"WebhookEvent":     "webhooks",
}

// CheckReservedName reports whether name collides with a model the scaffold
// owns. Returns nil for a free name.
func CheckReservedName(name string) error {
	owner, taken := reservedModels[name]
	if !taken {
		return nil
	}
	return fmt.Errorf(
		"%q is a built-in model — it belongs to %s\n\n"+
			"Generating over it would overwrite apps/api/internal/models/%s.go and\n"+
			"break the build, and `grit remove resource %s` would then delete the\n"+
			"original. Pick another name:\n\n"+
			"  grit generate resource Support%s --fields \"...\"\n\n"+
			"If you really mean to replace the built-in, pass --force.",
		name, owner, toSnake(name), name, name,
	)
}

// ReservedNames returns the reserved model names, sorted. Used by the docs
// generator and the tests, so the published list cannot drift from the code.
func ReservedNames() []string {
	out := make([]string, 0, len(reservedModels))
	for name := range reservedModels {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// toSnake is the file-name form of a model name (TicketReply → ticket_reply).
// Names() does the same job, but building a Generator needs a project root and
// this check runs before that.
func toSnake(name string) string {
	var b strings.Builder
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('_')
		}
		if r >= 'A' && r <= 'Z' {
			r += 'a' - 'A'
		}
		b.WriteRune(r)
	}
	return b.String()
}
