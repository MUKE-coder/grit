package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"gorm.io/gorm"

	"saas/apps/api/internal/ids"
)

// RecoveryContactKind distinguishes the two destinations a code can go to.
type RecoveryContactKind string

const (
	RecoveryEmail RecoveryContactKind = "email"
	RecoveryPhone RecoveryContactKind = "phone"
)

// RecoveryContactToken is a single-use, expiring proof that somebody can read
// mail or texts at a destination.
//
// Modelled on EmailVerificationToken and for the same reasons: only the hash is
// stored, so leaking this table does not let an attacker confirm a contact they
// do not control, and the destination is recorded as it was when the code was
// issued. Without that last part, changing the pending address after requesting
// a code would let you verify an address you never owned.
type RecoveryContactToken struct {
	ID     string `gorm:"primarykey;size:36" json:"id"`
	UserID string `gorm:"size:36;index;not null" json:"user_id"`

	Kind        RecoveryContactKind `gorm:"size:16;not null" json:"kind"`
	Destination string              `gorm:"size:255;not null" json:"destination"`

	CodeHash string `gorm:"size:64;index;not null" json:"-"`

	// Attempts is capped so a six-digit code cannot be brute-forced inside its
	// fifteen-minute life. A million codes and unlimited guesses is not a
	// secret, it is a formality.
	Attempts int `gorm:"default:0" json:"-"`

	ExpiresAt time.Time  `gorm:"index" json:"expires_at"`
	UsedAt    *time.Time `gorm:"index" json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

func (t *RecoveryContactToken) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = ids.New()
	}
	return nil
}

// RecoveryContact is a verified way back into an account.
//
// Its own table rather than columns on User, and that is not a modelling
// preference. grit upgrade does not regenerate the User model, so shipping
// this as user columns meant an upgraded project got the handler that reads
// them and a model without them: a half-feature that does not compile. A table
// of its own arrives complete through AutoMigrate, and it leaves room for more
// than one contact per kind later without another migration.
type RecoveryContact struct {
	ID     string `gorm:"primarykey;size:36" json:"id"`
	UserID string `gorm:"size:36;index:idx_recovery_user_kind,unique;not null" json:"user_id"`

	Kind        RecoveryContactKind `gorm:"size:16;index:idx_recovery_user_kind,unique;not null" json:"kind"`
	Destination string              `gorm:"size:255;not null" json:"destination"`

	// Only ever written once a code has been confirmed. An unverified row
	// would look like a way back in and would not be one.
	VerifiedAt time.Time `json:"verified_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (c *RecoveryContact) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = ids.New()
	}
	return nil
}

// HashRecoveryCode returns the storage form of a recovery code.
func HashRecoveryCode(code string) string {
	sum := sha256.Sum256([]byte(code))
	return hex.EncodeToString(sum[:])
}
