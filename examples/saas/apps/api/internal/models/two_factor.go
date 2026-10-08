package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"saas/apps/api/internal/crypto"
)

// TwoFactorConfig stores TOTP settings for a user.
// How a second factor is delivered.
const (
	// TwoFactorMethodApp is an authenticator app: the stronger of the two,
	// because the code never leaves the device.
	TwoFactorMethodApp = "app"
	// TwoFactorMethodEmail sends a code to the address on the account. Only as
	// safe as that mailbox, which is also where a password reset goes, so it
	// adds less than an authenticator does. It is here because the alternative
	// for somebody who will not install an app is no second factor at all.
	TwoFactorMethodEmail = "email"
)

type TwoFactorConfig struct {
	ID     uint   `gorm:"primarykey" json:"id"`
	UserID string `gorm:"size:36;uniqueIndex;not null" json:"user_id"`
	// Secret is encrypted with FIELD_ENCRYPTION_KEY when one is set. It was
	// stored in the clear, so anybody who could read the table, a backup or a
	// GORM Studio session could generate every user's codes. A secret written
	// before the key reads back as it is and is encrypted by grit migrate, or
	// the next time its code verifies. Without a key it is stored as before.
	Secret  crypto.EncryptedString `gorm:"type:text;not null" json:"-"`
	Enabled bool                   `gorm:"default:false" json:"enabled"`
	// Method is how the second factor is delivered: "app" for an authenticator
	// and "email" for a code sent to the address on the account.
	//
	// An authenticator is the stronger of the two, because an emailed code is
	// only as safe as the mailbox, and for most people the mailbox is what the
	// password reset goes to anyway. Email is here because the alternative for
	// somebody who will not install an app is no second factor at all, and that
	// is worse than a mailbox.
	//
	// Empty means "app": rows written before this column existed were all
	// authenticators, and defaulting it keeps them working without a backfill.
	Method      string                      `gorm:"size:16;not null;default:app" json:"method"`
	BackupCodes datatypes.JSONSlice[string] `gorm:"type:text" json:"-"`
	// LastUsedStep is the time step of the last code accepted. Only a later one
	// is accepted next, so a code cannot be used twice.
	LastUsedStep int64 `gorm:"not null;default:0" json:"-"`
	// StepOffset is how many 30 second steps this device's clock is from the
	// server's, measured when the authenticator was enrolled and re-measured on
	// every accepted code.
	//
	// Without it, a laptop a minute behind produces correct codes that are
	// always one step too old, and the only advice the server can give is "fix
	// your clock", which on a managed machine is not advice. It does not widen
	// what is accepted: the window stays one step, around this device's time
	// rather than the server's. Zero for every device whose clock is right,
	// which is almost all of them.
	StepOffset int64 `gorm:"not null;default:0" json:"-"`
	// SetupCodeHash and SetupCodeExpiresAt hold the code that proves somebody
	// can read the mailbox before email becomes their second factor. Turning on
	// a factor you cannot receive is how an account locks itself out, so it is
	// deliberately two steps.
	SetupCodeHash      string     `gorm:"size:64" json:"-"`
	SetupCodeExpiresAt *time.Time `json:"-"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// TrustedDevice stores a remembered device that can skip TOTP.
type TrustedDevice struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    string         `gorm:"size:36;index;not null" json:"user_id"`
	TokenHash string         `gorm:"size:64;uniqueIndex;not null" json:"-"`
	UserAgent string         `gorm:"size:500" json:"user_agent"`
	IPAddress string         `gorm:"size:45" json:"ip_address"`
	ExpiresAt time.Time      `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TOTPPendingToken stores short-lived tokens for the TOTP verification step.
type TOTPPendingToken struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    string    `gorm:"size:36;index;not null" json:"user_id"`
	TokenHash string    `gorm:"size:64;uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	// Attempts counts wrong codes against this token; at totp.MaxPendingAttempts
	// it is refused.
	Attempts int `gorm:"not null;default:0" json:"-"`
	// CodeHash is set when the second factor is a code sent by email, and holds
	// the SHA-256 of that code. Hashed for the same reason the token is: a
	// six-digit code sitting in a table is a six-digit code an attacker with a
	// read of that table can type.
	//
	// Empty means the challenge is an authenticator, and the code is computed
	// rather than stored.
	CodeHash  string    `gorm:"size:64" json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

// CleanupExpiredTOTPTokens removes expired pending tokens and trusted devices.
func CleanupExpiredTOTPTokens(db *gorm.DB) error {
	now := time.Now()
	if err := db.Where("expires_at < ?", now).Delete(&TOTPPendingToken{}).Error; err != nil {
		return err
	}
	return db.Where("expires_at < ?", now).Delete(&TrustedDevice{}).Error
}
