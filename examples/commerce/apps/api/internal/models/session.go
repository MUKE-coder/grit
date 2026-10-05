package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"commerce/apps/api/internal/ids"
	"gorm.io/gorm"
)

// Session is a server-side record of one issued refresh token — i.e. one
// logged-in device. It exists so a refresh token can be revoked; a bare JWT
// cannot be.
type Session struct {
	ID     string `gorm:"primarykey;size:36" json:"id"`
	UserID string `gorm:"size:36;index;not null" json:"user_id"`

	// TokenHash is the SHA-256 of the current refresh token. The token itself is
	// never persisted, so a dump of this table cannot be replayed as a session.
	TokenHash string `gorm:"size:64;uniqueIndex;not null" json:"-"`

	// PrevTokenHash is the hash of the token this one replaced. If a client ever
	// presents THAT token again, the token was captured and replayed after
	// rotation — the session is killed rather than refreshed.
	PrevTokenHash string `gorm:"size:64;index" json:"-"`

	UserAgent string `gorm:"size:512" json:"user_agent"`
	IP        string `gorm:"size:64" json:"ip"`

	CreatedAt  time.Time  `json:"created_at"`
	LastSeenAt time.Time  `gorm:"index" json:"last_seen_at"`
	ExpiresAt  time.Time  `gorm:"index" json:"expires_at"`
	RevokedAt  *time.Time `gorm:"index" json:"revoked_at,omitempty"`
}

func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = ids.New()
	}
	return nil
}

// HashSessionToken returns the storage form of a refresh token. Exported so the
// service and any test can derive the same value.
func HashSessionToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// Active reports whether the session may still be used to refresh.
func (s *Session) Active(now time.Time, idleTimeout time.Duration) bool {
	if s.RevokedAt != nil {
		return false
	}
	if now.After(s.ExpiresAt) {
		return false // absolute timeout
	}
	if idleTimeout > 0 && now.Sub(s.LastSeenAt) > idleTimeout {
		return false // idle timeout
	}
	return true
}
