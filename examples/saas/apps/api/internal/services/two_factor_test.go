package services

import (
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"saas/apps/api/internal/models"
	"saas/apps/api/internal/totp"
)

func twoFactorDB(t *testing.T) (*gorm.DB, *TwoFactorService) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}, &models.TwoFactorConfig{}, &models.TOTPPendingToken{}, &models.TrustedDevice{}))
	return db, &TwoFactorService{DB: db}
}

func enrolled(t *testing.T, db *gorm.DB, userID string, codes ...string) *models.TwoFactorConfig {
	t.Helper()
	config := models.TwoFactorConfig{
		UserID:      userID,
		Enabled:     true,
		Secret:      "SECRET",
		BackupCodes: datatypes.JSONSlice[string](codes),
	}
	require.NoError(t, db.Create(&config).Error)
	return &config
}

// A code is good for its whole time step, so the step has to be spent once. Two
// requests arriving with the same code both validate it; only one may sign in.
func TestSpendTOTPStepIsWonByOneRequest(t *testing.T) {
	db, twoFactor := twoFactorDB(t)
	ctx := context.Background()
	config := enrolled(t, db, "u1")
	step := time.Now().Unix() / totp.Period

	spent, err := twoFactor.SpendTOTPStep(ctx, config, step)
	require.NoError(t, err)
	assert.True(t, spent, "the first request did not spend the code")

	// The second request has read the same row, so it carries the same step.
	spent, err = twoFactor.SpendTOTPStep(ctx, config, step)
	require.NoError(t, err)
	assert.False(t, spent, "the same code was spent twice, so a replay signs in")

	var stored models.TwoFactorConfig
	require.NoError(t, db.First(&stored, config.ID).Error)
	assert.Equal(t, step, stored.LastUsedStep)
}

// The same for a backup code: both requests find it in the list, and only the
// one that writes may use it.
func TestSpendBackupCodeIsWonByOneRequest(t *testing.T) {
	db, twoFactor := twoFactorDB(t)
	ctx := context.Background()
	config := enrolled(t, db, "u2", "hash-a", "hash-b", "hash-c")
	read := *config // what a second request would hold

	remaining, spent, err := twoFactor.SpendBackupCode(ctx, config, 1)
	require.NoError(t, err)
	require.True(t, spent)
	assert.Equal(t, []string{"hash-a", "hash-c"}, remaining)

	_, spent, err = twoFactor.SpendBackupCode(ctx, &read, 1)
	require.NoError(t, err)
	assert.False(t, spent, "a backup code was spent twice")

	var stored models.TwoFactorConfig
	require.NoError(t, db.First(&stored, config.ID).Error)
	assert.Equal(t, []string{"hash-a", "hash-c"}, []string(stored.BackupCodes))
}

// And for the pending token, which is what makes two requests carrying one
// half-finished sign-in resolve to one session.
func TestSpendPendingTokenIsWonByOneRequest(t *testing.T) {
	db, twoFactor := twoFactorDB(t)
	ctx := context.Background()
	pending := models.TOTPPendingToken{UserID: "u3", TokenHash: totp.HashToken("tok"), ExpiresAt: time.Now().Add(time.Minute)}
	require.NoError(t, db.Create(&pending).Error)

	spent, err := twoFactor.SpendPendingToken(ctx, pending.ID)
	require.NoError(t, err)
	assert.True(t, spent)

	spent, err = twoFactor.SpendPendingToken(ctx, pending.ID)
	require.NoError(t, err)
	assert.False(t, spent, "one pending token signed in twice")
}

// An expired pending token is not found, which is what stops a five-minute-old
// challenge from being answered an hour later.
func TestPendingByTokenIgnoresAnExpiredOne(t *testing.T) {
	db, twoFactor := twoFactorDB(t)
	ctx := context.Background()
	live := models.TOTPPendingToken{UserID: "u4", TokenHash: totp.HashToken("live"), ExpiresAt: time.Now().Add(time.Minute)}
	dead := models.TOTPPendingToken{UserID: "u4", TokenHash: totp.HashToken("dead"), ExpiresAt: time.Now().Add(-time.Minute)}
	require.NoError(t, db.Create(&live).Error)
	require.NoError(t, db.Create(&dead).Error)

	found, err := twoFactor.PendingByToken(ctx, "live")
	require.NoError(t, err)
	assert.Equal(t, live.ID, found.ID)

	_, err = twoFactor.PendingByToken(ctx, "dead")
	assert.Error(t, err, "an expired pending token was accepted")
}

// The lockout reads the count the database holds rather than the one the request
// arrived with, because several wrong codes can be in flight at once.
//
// Through AuthService, which owns the account's counters: a wrong password and a
// wrong code are counted against one account by one rule, and there used to be a
// copy of that rule per factor.
func TestCountLoginFailureReturnsTheStoredCount(t *testing.T) {
	db, _ := twoFactorDB(t)
	ctx := context.Background()
	auth := &AuthService{DB: db}
	u := models.User{ID: "u5", Email: "u5@x.com", FirstName: "U", LastName: "5", Role: models.RoleUser}
	require.NoError(t, db.Create(&u).Error)

	for want := 1; want <= 3; want++ {
		got, err := auth.CountLoginFailure(ctx, u.ID)
		require.NoError(t, err)
		assert.Equal(t, want, got)
	}

	until := time.Now().Add(time.Hour)
	pending := models.TOTPPendingToken{UserID: u.ID, TokenHash: totp.HashToken("x"), ExpiresAt: time.Now().Add(time.Minute)}
	require.NoError(t, db.Create(&pending).Error)
	require.NoError(t, auth.LockAccount(ctx, u.ID, until))

	var after models.User
	require.NoError(t, db.First(&after, "id = ?", u.ID).Error)
	require.NotNil(t, after.LockedUntil)
	assert.Equal(t, 0, after.FailedLoginCount, "the count was not reset with the lock")
	var left int64
	db.Model(&models.TOTPPendingToken{}).Where("user_id = ?", u.ID).Count(&left)
	assert.Equal(t, int64(0), left, "a locked account kept its pending tokens, which are more guesses")
}

// Turning the second factor off takes the config, the devices and the pending
// tokens with it, or leaves all three alone.
func TestDisableForRemovesEverythingItDependsOn(t *testing.T) {
	db, twoFactor := twoFactorDB(t)
	ctx := context.Background()
	enrolled(t, db, "u6", "hash")
	require.NoError(t, db.Create(&models.TrustedDevice{UserID: "u6", TokenHash: "d", ExpiresAt: time.Now().Add(time.Hour)}).Error)
	require.NoError(t, db.Create(&models.TOTPPendingToken{UserID: "u6", TokenHash: "p", ExpiresAt: time.Now().Add(time.Minute)}).Error)

	require.NoError(t, twoFactor.DisableFor(ctx, "u6"))

	for _, m := range []interface{}{&models.TwoFactorConfig{}, &models.TrustedDevice{}, &models.TOTPPendingToken{}} {
		var n int64
		require.NoError(t, db.Model(m).Where("user_id = ?", "u6").Count(&n).Error)
		assert.Equal(t, int64(0), n, "%T outlived the second factor", m)
	}
}

// A trusted device is revoked by the person who holds it and nobody else: the id
// alone used to be enough.
func TestRevokeTrustedDeviceIsScopedToItsOwner(t *testing.T) {
	db, twoFactor := twoFactorDB(t)
	ctx := context.Background()
	mine := models.TrustedDevice{UserID: "u7", TokenHash: "mine", ExpiresAt: time.Now().Add(time.Hour)}
	require.NoError(t, db.Create(&mine).Error)

	revoked, err := twoFactor.RevokeTrustedDevice(ctx, "somebody-else", "1")
	require.NoError(t, err)
	assert.False(t, revoked, "another account revoked this device by guessing its id")

	devices, err := twoFactor.TrustedDevices(ctx, "u7")
	require.NoError(t, err)
	require.Len(t, devices, 1)

	revoked, err = twoFactor.RevokeTrustedDevice(ctx, "u7", "1")
	require.NoError(t, err)
	assert.True(t, revoked)
}

// An expired device is not trusted, and listing one does not delete it: the
// nightly cleanup owns deletion, and a read should not write.
func TestTrustedDevicesIgnoresExpiredRows(t *testing.T) {
	db, twoFactor := twoFactorDB(t)
	ctx := context.Background()
	require.NoError(t, db.Create(&models.TrustedDevice{UserID: "u8", TokenHash: "old", ExpiresAt: time.Now().Add(-time.Hour)}).Error)

	devices, err := twoFactor.TrustedDevices(ctx, "u8")
	require.NoError(t, err)
	assert.Empty(t, devices)
	n, err := twoFactor.CountTrustedDevices(ctx, "u8")
	require.NoError(t, err)
	assert.Equal(t, int64(0), n)

	var rows int64
	db.Model(&models.TrustedDevice{}).Where("user_id = ?", "u8").Count(&rows)
	assert.Equal(t, int64(1), rows, "a read deleted the row")
}
