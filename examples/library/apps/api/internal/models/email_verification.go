package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
	"library/apps/api/internal/ids"
)

// EmailVerificationToken is a single-use, expiring proof that someone can read
// mail at an address. Only the hash is stored, so leaking this table does not
// let an attacker verify addresses they do not control.
type EmailVerificationToken struct {
	ID     string `gorm:"primarykey;size:36" json:"id"`
	UserID string `gorm:"size:36;index;not null" json:"user_id"`

	TokenHash string `gorm:"size:64;uniqueIndex;not null" json:"-"`

	// Email is recorded as it was when the token was issued. If the user
	// changes their address before clicking the link, the token must not
	// verify the new one — that is how you verify an address you never owned.
	Email string `gorm:"size:255;not null" json:"email"`

	ExpiresAt time.Time  `gorm:"index" json:"expires_at"`
	UsedAt    *time.Time `gorm:"index" json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

func (t *EmailVerificationToken) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = ids.New()
	}
	return nil
}

// HashVerificationToken returns the storage form of a verification token.
func HashVerificationToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
