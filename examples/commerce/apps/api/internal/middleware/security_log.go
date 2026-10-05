package middleware

import (
	"context"
	"log"

	"gorm.io/gorm"

	"commerce/apps/api/internal/audit"
	"commerce/apps/api/internal/models"
)

// SecurityEvent is the discriminator that goes into ActivityLog.Action
// for OWASP A09 "security event" records. Use these constants — never
// raw strings — so the audit dashboard can filter on a stable enum.
const (
	EventLoginSuccess      = "security.login.success"
	EventLoginFailure      = "security.login.failure"
	EventLogout            = "security.logout"
	EventPasswordChanged   = "security.password.changed"
	EventPasswordResetReq  = "security.password.reset_requested"
	EventTOTPEnabled       = "security.totp.enabled"
	EventTOTPDisabled      = "security.totp.disabled"
	EventTOTPChallengeFail = "security.totp.challenge_failed"
	EventRoleChanged       = "security.role.changed"
	EventAccountLocked     = "security.account.locked"
	EventAuthZDenied       = "security.authz.denied"
	EventSuspiciousRequest = "security.suspicious_request"
)

// LogSecurityEvent records a security-relevant event to the activity log.
//
// We piggyback on the existing tamper-evident ActivityLog table — Method
// "SECURITY" + Path = event-constant is the convention. This means every
// security event automatically picks up the hash-chain integrity the
// activity log already provides (a deleted security event breaks the
// chain at verify time), without a separate model to maintain.
//
// Written through the chain, so it is linked and stamped like every other
// entry. It used to be inserted directly with no hash at all: the first event
// broke verification, and the unique index on hash refused every event after.
//
// Database errors are logged but never returned — failing a login or
// logout because the audit log had a hiccup turns audit infrastructure
// into a DoS amplifier. The request itself was usually fine.
func LogSecurityEvent(ctx context.Context, db *gorm.DB, userID, event, ip, userAgent string) {
	if db == nil {
		return
	}
	entry := models.ActivityLog{
		UserID:    userID,
		Method:    "SECURITY",
		Path:      event,
		IPAddress: ip,
		UserAgent: userAgent,
	}
	if err := audit.AppendChained(db.WithContext(ctx), &entry); err != nil {
		log.Printf("security_log: failed to record %s: %v", event, err)
	}
}
