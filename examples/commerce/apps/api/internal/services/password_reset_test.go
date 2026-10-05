package services_test

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/services"
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
