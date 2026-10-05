package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
)

// How long a sign-in link lives.
//
// Fifteen minutes is the compromise everybody lands on: long enough to switch
// to a phone and find the email, short enough that a link forwarded or left in
// a mailbox is not a standing key to the account.
const MagicLinkExpiry = 15 * time.Minute

// MagicLinkRate is how often one account may ask for a link. Without it the
// form is a way to send somebody a hundred emails.
const MagicLinkRate = 60 * time.Second

var (
	// ErrMagicLinkUsed is a link that has already signed somebody in. Told
	// apart from an invalid one on purpose: clicking the same email twice is
	// an ordinary mistake and deserves an answer that says so.
	ErrMagicLinkUsed = errors.New("this sign-in link has already been used")
	// ErrMagicLinkExpired is a link past its life.
	ErrMagicLinkExpired = errors.New("this sign-in link has expired")
	// ErrMagicLinkInvalid is anything else: a token that never existed, or one
	// whose row is gone.
	ErrMagicLinkInvalid = errors.New("this sign-in link is not valid")
	// ErrMagicLinkTooSoon is the rate limit.
	ErrMagicLinkTooSoon = errors.New("a link was sent a moment ago: check your email")
)

// IssueMagicLink records a new link for a user and returns the token to email.
//
// The plaintext is returned once, here, and never stored.
func IssueMagicLink(ctx context.Context, db *gorm.DB, userID, ip, agent string) (string, error) {
	var recent int64
	if err := db.WithContext(ctx).Model(&models.MagicLinkToken{}).
		Where("user_id = ? AND created_at > ?", userID, time.Now().Add(-MagicLinkRate)).
		Count(&recent).Error; err != nil {
		return "", fmt.Errorf("counting recent links: %w", err)
	}
	if recent > 0 {
		return "", ErrMagicLinkTooSoon
	}

	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generating a sign-in link: %w", err)
	}
	token := hex.EncodeToString(raw)

	row := models.MagicLinkToken{
		UserID:    userID,
		TokenHash: HashMagicLink(token),
		IPAddress: ip,
		UserAgent: agent,
		ExpiresAt: time.Now().Add(MagicLinkExpiry),
	}
	if err := db.WithContext(ctx).Create(&row).Error; err != nil {
		return "", fmt.Errorf("recording a sign-in link: %w", err)
	}
	return token, nil
}

// HashMagicLink is how a token is stored and looked up.
func HashMagicLink(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// ConsumeMagicLink spends a token and returns whose it was.
//
// The update is conditional on the token still being unused, so two clicks
// arriving together cannot both succeed: exactly one updates a row, and the
// other is told the link has been used.
func ConsumeMagicLink(ctx context.Context, db *gorm.DB, token string) (string, error) {
	if token == "" {
		return "", ErrMagicLinkInvalid
	}

	var row models.MagicLinkToken
	err := db.WithContext(ctx).Where("token_hash = ?", HashMagicLink(token)).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrMagicLinkInvalid
		}
		return "", fmt.Errorf("reading the sign-in link: %w", err)
	}
	if row.UsedAt != nil {
		return "", ErrMagicLinkUsed
	}
	if time.Now().After(row.ExpiresAt) {
		return "", ErrMagicLinkExpired
	}

	now := time.Now()
	res := db.WithContext(ctx).Model(&models.MagicLinkToken{}).
		Where("id = ? AND used_at IS NULL", row.ID).
		Update("used_at", now)
	if res.Error != nil {
		return "", fmt.Errorf("spending the sign-in link: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		// Somebody else spent it between the read and the update.
		return "", ErrMagicLinkUsed
	}
	return row.UserID, nil
}
