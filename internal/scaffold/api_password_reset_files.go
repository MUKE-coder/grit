package scaffold

// Password reset.
//
// The reset flow used to be a placeholder: ForgotPassword generated a token and
// logged it without storing it, and ResetPassword hashed the new password,
// discarded it, and returned "Password reset successfully". Anyone who used
// forgot-password believed they had locked an attacker out and had in fact
// changed nothing — the most dangerous kind of bug, because it reports success.
//
// This is the real thing:
//   - the token is random 256-bit, stored only as its SHA-256
//   - single use, enforced by a conditional UPDATE so two concurrent requests
//     cannot both consume it
//   - short TTL, and issuing a new one invalidates any earlier unused token
//   - consuming it revokes every session, because the whole point of a reset is
//     to evict whoever you think is in your account

func apiPasswordResetModelGo() string {
	return `package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"{{MODULE}}/internal/ids"
	"gorm.io/gorm"
)

// PasswordResetToken is a single-use, expiring grant to change one user's
// password without knowing the old one. Only the hash is stored, so a leak of
// this table cannot be used to reset anyone's password.
type PasswordResetToken struct {
	ID     string ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	UserID string ` + "`" + `gorm:"size:36;index;not null" json:"user_id"` + "`" + `

	TokenHash string ` + "`" + `gorm:"size:64;uniqueIndex;not null" json:"-"` + "`" + `

	// RequestIP is kept for audit: "who asked for this reset?" is the first
	// question after an account takeover.
	RequestIP string ` + "`" + `gorm:"size:64" json:"request_ip"` + "`" + `

	ExpiresAt time.Time  ` + "`" + `gorm:"index" json:"expires_at"` + "`" + `
	UsedAt    *time.Time ` + "`" + `gorm:"index" json:"used_at,omitempty"` + "`" + `
	CreatedAt time.Time  ` + "`" + `json:"created_at"` + "`" + `
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
`
}

func apiPasswordResetServiceGo() string {
	return `package services

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
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
`
}

func apiPasswordResetTestGo() string {
	return `package services_test

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
)

func newResetDB(tb testing.TB) *gorm.DB {
	tb.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(tb, err)
	require.NoError(tb, db.AutoMigrate(&models.PasswordResetToken{}))
	return db
}

func TestResetTokenStoresOnlyTheHash(t *testing.T) {
	db := newResetDB(t)
	row, err := services.CreatePasswordResetToken(db, "user-1", "secret-token", "1.2.3.4")
	require.NoError(t, err)

	assert.Equal(t, models.HashResetToken("secret-token"), row.TokenHash)
	assert.NotEqual(t, "secret-token", row.TokenHash)
	assert.Equal(t, "1.2.3.4", row.RequestIP)
}

func TestConsumeResetTokenReturnsTheOwner(t *testing.T) {
	db := newResetDB(t)
	_, err := services.CreatePasswordResetToken(db, "user-1", "tok", "")
	require.NoError(t, err)

	userID, err := services.ConsumePasswordResetToken(db, "tok")
	require.NoError(t, err)
	assert.Equal(t, "user-1", userID)
}

func TestResetTokenIsSingleUse(t *testing.T) {
	db := newResetDB(t)
	_, err := services.CreatePasswordResetToken(db, "user-1", "tok", "")
	require.NoError(t, err)

	_, err = services.ConsumePasswordResetToken(db, "tok")
	require.NoError(t, err)

	_, err = services.ConsumePasswordResetToken(db, "tok")
	assert.ErrorIs(t, err, services.ErrResetTokenInvalid, "a reset token must not work twice")
}

func TestResetTokenExpires(t *testing.T) {
	db := newResetDB(t)
	row, err := services.CreatePasswordResetToken(db, "user-1", "tok", "")
	require.NoError(t, err)

	require.NoError(t, db.Model(row).Update("expires_at", time.Now().Add(-time.Minute)).Error)

	_, err = services.ConsumePasswordResetToken(db, "tok")
	assert.ErrorIs(t, err, services.ErrResetTokenInvalid, "an expired token must not work")
}

func TestUnknownResetTokenIsRejected(t *testing.T) {
	db := newResetDB(t)
	_, err := services.ConsumePasswordResetToken(db, "never-issued")
	assert.ErrorIs(t, err, services.ErrResetTokenInvalid)
}

// Requesting a second link must retire the first, or every request widens the
// window of usable tokens.
func TestIssuingANewTokenBurnsTheOldOne(t *testing.T) {
	db := newResetDB(t)
	_, err := services.CreatePasswordResetToken(db, "user-1", "first", "")
	require.NoError(t, err)
	_, err = services.CreatePasswordResetToken(db, "user-1", "second", "")
	require.NoError(t, err)

	_, err = services.ConsumePasswordResetToken(db, "first")
	assert.ErrorIs(t, err, services.ErrResetTokenInvalid, "the superseded token must be dead")

	userID, err := services.ConsumePasswordResetToken(db, "second")
	require.NoError(t, err)
	assert.Equal(t, "user-1", userID)
}

// One user's request must not affect another's outstanding token.
func TestResetTokensAreScopedPerUser(t *testing.T) {
	db := newResetDB(t)
	_, err := services.CreatePasswordResetToken(db, "user-1", "mine", "")
	require.NoError(t, err)
	_, err = services.CreatePasswordResetToken(db, "user-2", "theirs", "")
	require.NoError(t, err)

	userID, err := services.ConsumePasswordResetToken(db, "mine")
	require.NoError(t, err)
	assert.Equal(t, "user-1", userID)
}

func TestPurgeExpiredResetTokens(t *testing.T) {
	db := newResetDB(t)
	row, err := services.CreatePasswordResetToken(db, "user-1", "old", "")
	require.NoError(t, err)
	require.NoError(t, db.Model(row).Update("expires_at", time.Now().Add(-48*time.Hour)).Error)

	_, err = services.CreatePasswordResetToken(db, "user-2", "fresh", "")
	require.NoError(t, err)

	n, err := services.PurgeExpiredPasswordResetTokens(db, 24*time.Hour)
	require.NoError(t, err)
	assert.Equal(t, int64(1), n)

	userID, err := services.ConsumePasswordResetToken(db, "fresh")
	require.NoError(t, err)
	assert.Equal(t, "user-2", userID)
}
`
}
