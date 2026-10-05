package models

import "commerce/apps/api/internal/erasure"

// The framework's own tables whose rows exist only to serve a user and carry no
// independent compliance value: hard-deleted when that user is erased.
//
// Resources generated with --owned-by register themselves from their own model
// files, so erasing a user deletes the records they own too. The activity log
// is deliberately absent: its rows hold only a UUID, so anonymizing the user
// anonymizes them, and deleting them would break the audit hash chain.
func init() {
	for _, m := range []interface{}{
		&Upload{}, &Session{}, &PasswordResetToken{}, &UserRole{},
		&TwoFactorConfig{}, &TrustedDevice{}, &TOTPPendingToken{},
		&DashboardLayout{}, &Notification{},
	} {
		erasure.Register(m, "user_id")
	}
}
