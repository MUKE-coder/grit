package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"commerce/apps/api/internal/ids"
	"gorm.io/gorm"
)

// PasswordResetToken is a single-use, expiring grant to change one user's
// password without knowing the old one. Only the hash is stored, so a leak of
// this table cannot be used to reset anyone's password.
type PasswordResetToken struct {
	ID     string `gorm:"primarykey;size:36" json:"id"`
	UserID string `gorm:"size:36;index;not null" json:"user_id"`

	TokenHash string `gorm:"size:64;uniqueIndex;not null" json:"-"`

	// RequestIP is kept for audit: "who asked for this reset?" is the first
	// question after an account takeover.
	RequestIP string `gorm:"size:64" json:"request_ip"`

	ExpiresAt time.Time  `gorm:"index" json:"expires_at"`
	UsedAt    *time.Time `gorm:"index" json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

func (t *PasswordResetToken) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = ids.New()
	}
	return nil
}

// HashResetToken returns the storage form of a reset token.
func HashResetToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
