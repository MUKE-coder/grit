package services

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
)

// PasswordResetTTL is how long a reset link stays usable. Short on purpose —
// the link sits in an inbox, which is exactly where an attacker with mailbox
// access goes looking.
var PasswordResetTTL = time.Hour

// ErrResetTokenInvalid covers every failure mode deliberately: unknown,
// already used, and expired all look identical to the caller, so the endpoint
// cannot be used to probe which tokens once existed.
var ErrResetTokenInvalid = errors.New("password reset token is not valid")

// CreatePasswordResetToken records a freshly issued reset token.
//
// Any earlier unused token for the same user is burned first. Otherwise
// requesting a second link would leave the first one live, and a reset flow
// that accumulates valid tokens is a reset flow with a widening attack window.
func CreatePasswordResetToken(db *gorm.DB, userID, token, requestIP string) (*models.PasswordResetToken, error) {
	now := time.Now()

	if err := db.Model(&models.PasswordResetToken{}).
		Where("user_id = ? AND used_at IS NULL", userID).
		Update("used_at", &now).Error; err != nil {
		return nil, err
	}

	row := &models.PasswordResetToken{
		UserID:    userID,
		TokenHash: models.HashResetToken(token),
		RequestIP: requestIP,
		ExpiresAt: now.Add(PasswordResetTTL),
	}
	if err := db.Create(row).Error; err != nil {
		return nil, err
	}
	return row, nil
}

// ConsumePasswordResetToken validates a token and marks it used in one
// conditional UPDATE, returning the user it belongs to.
//
// The single statement is the point: checking "is it unused?" and then writing
// "it is now used" as two steps lets two concurrent requests both pass the
// check and both reset the password. Here the database decides — exactly one
// UPDATE reports a row affected.
func ConsumePasswordResetToken(db *gorm.DB, token string) (string, error) {
	hash := models.HashResetToken(token)
	now := time.Now()

	res := db.Model(&models.PasswordResetToken{}).
		Where("token_hash = ? AND used_at IS NULL AND expires_at > ?", hash, now).
		Update("used_at", &now)
	if res.Error != nil {
		return "", res.Error
	}
	if res.RowsAffected == 0 {
		return "", ErrResetTokenInvalid
	}

	var row models.PasswordResetToken
	if err := db.Where("token_hash = ?", hash).First(&row).Error; err != nil {
		return "", err
	}
	return row.UserID, nil
}

// PurgeExpiredPasswordResetTokens drops spent and expired rows. Safe to call
// from a cron job; nothing depends on the history beyond the audit window.
func PurgeExpiredPasswordResetTokens(db *gorm.DB, olderThan time.Duration) (int64, error) {
	cutoff := time.Now().Add(-olderThan)
	res := db.Where("expires_at < ? OR (used_at IS NOT NULL AND used_at < ?)", cutoff, cutoff).
		Delete(&models.PasswordResetToken{})
	return res.RowsAffected, res.Error
}
