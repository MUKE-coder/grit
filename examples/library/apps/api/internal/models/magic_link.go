package models

import (
	"time"

	"gorm.io/gorm"
)

// MagicLinkToken is one emailed sign-in link.
//
// Only the hash is stored. A table of live sign-in tokens in the clear is a
// table of passwords: anyone who can read it can sign in as anybody who has
// asked for a link in the last quarter of an hour.
type MagicLinkToken struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	UserID    string `gorm:"size:36;index;not null" json:"user_id"`
	TokenHash string `gorm:"size:64;uniqueIndex;not null" json:"-"`

	// UsedAt is set the moment the token is spent, and the row is kept rather
	// than deleted: a second click should be able to say "this link has already
	// been used" instead of "this link is invalid", which is what somebody sees
	// when they click the same email twice and needs a different answer.
	UsedAt *time.Time `json:"used_at,omitempty"`

	// What asked for it, for the admin's activity log and for anybody looking
	// at a link they did not request.
	IPAddress string `gorm:"size:45" json:"ip_address,omitempty"`
	UserAgent string `gorm:"size:500" json:"user_agent,omitempty"`

	ExpiresAt time.Time `gorm:"index;not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// CleanupExpiredMagicLinks removes links nobody can use any more.
func CleanupExpiredMagicLinks(db *gorm.DB) error {
	return db.Where("expires_at < ?", time.Now().Add(-24*time.Hour)).
		Delete(&MagicLinkToken{}).Error
}
