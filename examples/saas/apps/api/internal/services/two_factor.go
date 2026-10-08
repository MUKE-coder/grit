package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"saas/apps/api/internal/crypto"
	"saas/apps/api/internal/models"
	"saas/apps/api/internal/totp"
)

// TwoFactorService owns the two-factor tables: the config a user enrols, the
// pending token a half-finished sign-in carries, and the devices allowed to skip
// the prompt.
//
// The handler decides what the answer is: which refusal a wrong code gets, when
// an account is locked, what a setup screen shows. This decides what the data
// does, and three of these methods are the reason the split matters, because
// each one is a compare-and-set that two concurrent requests must not both win:
// a code spent once, a backup code spent once, a pending token used once. Those
// are not rules a handler can hold.
type TwoFactorService struct {
	DB *gorm.DB
}

func (s *TwoFactorService) db(ctx context.Context) *gorm.DB {
	return s.DB.WithContext(ctx)
}

// Config reads a user's two-factor row, enabled or not.
func (s *TwoFactorService) Config(ctx context.Context, userID string) (*models.TwoFactorConfig, error) {
	var config models.TwoFactorConfig
	if err := s.db(ctx).Where("user_id = ?", userID).First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

// ConfigOrNew reads the row, or hands back an unsaved one for this user.
//
// Email setup writes to whichever it gets: an account with an authenticator
// keeps it working until the new method is confirmed, so there is nothing to
// create up front.
func (s *TwoFactorService) ConfigOrNew(ctx context.Context, userID string) (*models.TwoFactorConfig, error) {
	config, err := s.Config(ctx, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &models.TwoFactorConfig{UserID: userID}, nil
	}
	if err != nil {
		return nil, err
	}
	return config, nil
}

// EnsureConfig reads the row or creates it, for enrolment.
func (s *TwoFactorService) EnsureConfig(ctx context.Context, userID string) (*models.TwoFactorConfig, error) {
	var config models.TwoFactorConfig
	if err := s.db(ctx).Where("user_id = ?", userID).
		FirstOrCreate(&config, models.TwoFactorConfig{UserID: userID}).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

// EnabledConfig reads the row only when the account really has a second factor.
func (s *TwoFactorService) EnabledConfig(ctx context.Context, userID string) (*models.TwoFactorConfig, error) {
	var config models.TwoFactorConfig
	if err := s.db(ctx).Where("user_id = ? AND enabled = ?", userID, true).First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

// IsEnabled reports whether the account has a second factor, for a caller that
// wants an answer rather than a row.
func (s *TwoFactorService) IsEnabled(ctx context.Context, userID string) bool {
	config, err := s.Config(ctx, userID)
	return err == nil && config.Enabled
}

// SaveConfig writes a config the caller has changed.
func (s *TwoFactorService) SaveConfig(ctx context.Context, config *models.TwoFactorConfig) error {
	return s.db(ctx).Save(config).Error
}

// DisableFor removes the second factor and everything that depends on it.
//
// The config, the trusted devices and the pending tokens go together, or none of
// them do. The three deletes used to run unchecked, and the answer was
// "disabled" whether or not anything had been deleted.
func (s *TwoFactorService) DisableFor(ctx context.Context, userID string) error {
	return s.db(ctx).Transaction(func(tx *gorm.DB) error {
		for _, table := range []interface{}{&models.TwoFactorConfig{}, &models.TrustedDevice{}, &models.TOTPPendingToken{}} {
			if err := tx.Where("user_id = ?", userID).Delete(table).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// SealSecret stores a secret that is still in the clear encrypted, now that it
// has verified a code.
//
// EncryptedString reads a value without its enc:v1: prefix as it is, so a secret
// enabled before FIELD_ENCRYPTION_KEY was set keeps working. grit migrate
// encrypts all of them at once; this covers a project that has not run it. The
// update matches only a row that is still plaintext, and the caller treats a
// failure as nothing worse than a secret that stays readable: it verifies either
// way.
func (s *TwoFactorService) SealSecret(ctx context.Context, config *models.TwoFactorConfig) error {
	if !crypto.EncryptionEnabled() {
		return nil
	}
	return s.db(ctx).Model(&models.TwoFactorConfig{}).
		Where("id = ? AND secret NOT LIKE ?", config.ID, "enc:v1:%").
		UpdateColumn("secret", config.Secret).Error
}

// PendingByToken reads the half-finished sign-in a token names, if it is still
// live. The token is hashed here, because the plaintext is never stored.
func (s *TwoFactorService) PendingByToken(ctx context.Context, token string) (*models.TOTPPendingToken, error) {
	var pending models.TOTPPendingToken
	if err := s.db(ctx).Where("token_hash = ? AND expires_at > ?", totp.HashToken(token), time.Now()).
		First(&pending).Error; err != nil {
		return nil, err
	}
	return &pending, nil
}

// SpendEmailChallengeCode clears the code on a pending token, and reports
// whether this request was the one that cleared it.
//
// The row itself goes when the session is issued. Clearing the hash first, and
// only while it still matches what this request read, means a replay in the same
// instant finds nothing to match.
func (s *TwoFactorService) SpendEmailChallengeCode(ctx context.Context, pending *models.TOTPPendingToken) (bool, error) {
	res := s.db(ctx).Model(&models.TOTPPendingToken{}).
		Where("id = ? AND code_hash = ?", pending.ID, pending.CodeHash).
		Update("code_hash", "")
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected == 1, nil
}

// SpendTOTPStep records the time step a code belongs to, and reports whether
// this request was the one that recorded it.
//
// A code is good for its whole window, so without this the same code signs in
// again for as long as it lasts. The step only moves forward and the update is
// conditional, so two requests cannot both spend one code.
//
// The offset is re-recorded from the step that actually matched, which is RFC
// 6238 resynchronisation: a clock losing a second a day is followed instead of
// eventually locking the account out.
func (s *TwoFactorService) SpendTOTPStep(ctx context.Context, config *models.TwoFactorConfig, step int64) (bool, error) {
	res := s.db(ctx).Model(&models.TwoFactorConfig{}).
		Where("id = ? AND last_used_step < ?", config.ID, step).
		UpdateColumns(map[string]interface{}{
			"last_used_step": step,
			"step_offset":    step - time.Now().Unix()/totp.Period,
		})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected == 1, nil
}

// SpendBackupCode removes the code at idx from the account's list and reports
// the codes that remain, and whether this request was the one that spent it.
//
// The update matches only while the list is still the one the caller read. Two
// requests carrying the same code both find it in the list; only the first to
// write can spend it, and the other is refused.
func (s *TwoFactorService) SpendBackupCode(ctx context.Context, config *models.TwoFactorConfig, idx int) ([]string, bool, error) {
	if idx < 0 || idx >= len(config.BackupCodes) {
		return nil, false, errors.New("backup code index out of range")
	}
	remaining := make([]string, 0, len(config.BackupCodes)-1)
	remaining = append(remaining, config.BackupCodes[:idx]...)
	remaining = append(remaining, config.BackupCodes[idx+1:]...)

	read, err := json.Marshal([]string(config.BackupCodes))
	if err != nil {
		return nil, false, err
	}
	res := s.db(ctx).Model(&models.TwoFactorConfig{}).
		Where("id = ? AND backup_codes = ?", config.ID, string(read)).
		Update("backup_codes", datatypes.JSONSlice[string](remaining))
	if res.Error != nil {
		return nil, false, res.Error
	}
	return remaining, res.RowsAffected == 1, nil
}

// SpendPendingToken deletes the pending token, and reports whether this request
// was the one that deleted it, so two requests carrying one token cannot both
// sign in.
func (s *TwoFactorService) SpendPendingToken(ctx context.Context, pendingID uint) (bool, error) {
	res := s.db(ctx).Delete(&models.TOTPPendingToken{}, pendingID)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected == 1, nil
}

// CountPendingAttempt counts a wrong code against one pending token.
//
// Its own UPDATE, with the increment in SQL, so concurrent guesses cannot share
// one increment and a token cannot be guessed more times than it allows.
func (s *TwoFactorService) CountPendingAttempt(ctx context.Context, pendingID uint) error {
	return s.db(ctx).Model(&models.TOTPPendingToken{}).Where("id = ?", pendingID).
		UpdateColumn("attempts", gorm.Expr("attempts + 1")).Error
}

// The account's own lockout counters are AuthService's, not this service's: a
// wrong password and a wrong code are counted against one account by one rule,
// and there used to be two copies of that rule. See AuthService.CountLoginFailure
// and AuthService.LockAccount.

// TrustedDevices lists the live devices allowed to skip the prompt, newest
// first.
//
// Expired rows are filtered rather than deleted: the nightly cleanup owns
// deletion, and a read should not write.
func (s *TwoFactorService) TrustedDevices(ctx context.Context, userID string) ([]models.TrustedDevice, error) {
	var devices []models.TrustedDevice
	if err := s.db(ctx).Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Order("created_at desc").Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

// CountTrustedDevices counts them, for the status endpoint.
func (s *TwoFactorService) CountTrustedDevices(ctx context.Context, userID string) (int64, error) {
	var n int64
	if err := s.db(ctx).Model(&models.TrustedDevice{}).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).Count(&n).Error; err != nil {
		return 0, err
	}
	return n, nil
}

// TrustDevice records a device allowed to skip the prompt.
func (s *TwoFactorService) TrustDevice(ctx context.Context, device *models.TrustedDevice) error {
	return s.db(ctx).Create(device).Error
}

// TrustedDeviceByToken reads the live device a cookie names.
func (s *TwoFactorService) TrustedDeviceByToken(ctx context.Context, userID, tokenHash string) (*models.TrustedDevice, error) {
	var device models.TrustedDevice
	if err := s.db(ctx).Where("user_id = ? AND token_hash = ? AND expires_at > ?", userID, tokenHash, time.Now()).
		First(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

// ExtendTrustedDevice pushes a device's expiry out, which is the sliding window
// that keeps a device people use from expiring on them.
func (s *TwoFactorService) ExtendTrustedDevice(ctx context.Context, device *models.TrustedDevice, until time.Time) error {
	device.ExpiresAt = until
	return s.db(ctx).Model(device).Update("expires_at", until).Error
}

// RevokeTrustedDevice removes one device, and reports whether there was one to
// remove.
//
// Scoped to the user: without that predicate any signed-in account could revoke
// anyone's device by guessing an id.
func (s *TwoFactorService) RevokeTrustedDevice(ctx context.Context, userID, id string) (bool, error) {
	res := s.db(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&models.TrustedDevice{})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

// RevokeTrustedDevices removes all of a user's devices.
func (s *TwoFactorService) RevokeTrustedDevices(ctx context.Context, userID string) error {
	return s.db(ctx).Where("user_id = ?", userID).Delete(&models.TrustedDevice{}).Error
}
