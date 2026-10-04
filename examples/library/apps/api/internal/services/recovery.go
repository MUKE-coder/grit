package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"gorm.io/gorm"

	"library/apps/api/internal/models"
)

// RecoveryCodeTTL is how long a code is good for.
//
// Fifteen minutes: long enough to switch to another device and read a message,
// short enough that a code left on a screen is not a standing key to the
// account.
const RecoveryCodeTTL = 15 * time.Minute

// MaxRecoveryAttempts caps guesses against one code.
//
// A six-digit code is a million possibilities, which sounds like a lot until
// you can try all of them. Five is the difference between a secret and a
// formality.
const MaxRecoveryAttempts = 5

var (
	ErrRecoveryCodeInvalid   = errors.New("that code is not valid")
	ErrRecoverySameAsPrimary = errors.New("a recovery address must be different from your sign-in address")
	ErrRecoveryInUse         = errors.New("that address is already in use on another account")
)

// NewRecoveryCode returns a six-digit code from crypto/rand.
//
// Not math/rand: a predictable recovery code is a way into every account at
// once, and the difference in effort here is nil.
func NewRecoveryCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", fmt.Errorf("generating recovery code: %w", err)
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// CreateRecoveryToken issues a code for a destination, burning any earlier
// unused ones of the same kind.
//
// Burning matters: two live codes for one account means the older one still
// works after the user re-requests, which is exactly what somebody who
// intercepted the first message wants.
func CreateRecoveryToken(db *gorm.DB, userID string, kind models.RecoveryContactKind, destination, code string) (*models.RecoveryContactToken, error) {
	row := &models.RecoveryContactToken{
		UserID:      userID,
		Kind:        kind,
		Destination: destination,
		CodeHash:    models.HashRecoveryCode(code),
		ExpiresAt:   time.Now().Add(RecoveryCodeTTL),
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		if err := tx.Model(&models.RecoveryContactToken{}).
			Where("user_id = ? AND kind = ? AND used_at IS NULL", userID, kind).
			Update("used_at", now).Error; err != nil {
			return err
		}
		return tx.Create(row).Error
	})
	if err != nil {
		return nil, fmt.Errorf("creating recovery token: %w", err)
	}
	return row, nil
}

// ConsumeRecoveryToken checks a code and, if it is right, marks the matching
// contact verified on the user.
//
// The lookup is by user and kind rather than by code hash. Looking up by hash
// would mean a wrong guess never touches a row, so the attempt counter would
// never increment and the cap would do nothing.
func ConsumeRecoveryToken(db *gorm.DB, userID string, kind models.RecoveryContactKind, code string) (string, error) {
	var row models.RecoveryContactToken
	if err := db.Where("user_id = ? AND kind = ? AND used_at IS NULL", userID, kind).
		Order("created_at desc").First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrRecoveryCodeInvalid
		}
		return "", err
	}

	if time.Now().After(row.ExpiresAt) || row.Attempts >= MaxRecoveryAttempts {
		return "", ErrRecoveryCodeInvalid
	}

	if row.CodeHash != models.HashRecoveryCode(code) {
		// Recorded outside a transaction, deliberately.
		//
		// The first version did the read, the increment and the rejection
		// inside one transaction, so returning the error rolled the increment
		// back with everything else. The counter never moved, the cap never
		// applied, and a six-digit code had unlimited guesses. A test that
		// spent the cap and then tried the right code is what caught it.
		if err := db.Model(&models.RecoveryContactToken{}).
			Where("id = ?", row.ID).
			UpdateColumn("attempts", gorm.Expr("attempts + 1")).Error; err != nil {
			return "", err
		}
		return "", ErrRecoveryCodeInvalid
	}

	// Only the success path is transactional: marking the token spent and
	// writing the contact onto the user have to happen together or not at all.
	var destination string
	err := db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		// Single-use enforced with a conditional UPDATE, so two concurrent
		// requests cannot both succeed.
		res := tx.Model(&models.RecoveryContactToken{}).
			Where("id = ? AND used_at IS NULL", row.ID).
			Update("used_at", now)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrRecoveryCodeInvalid
		}

		// Upsert on (user_id, kind): confirming a new address replaces the old
		// one rather than leaving two rows where only one can be current.
		contact := models.RecoveryContact{
			UserID:      userID,
			Kind:        kind,
			Destination: row.Destination,
			VerifiedAt:  now,
		}
		var existing models.RecoveryContact
		err := tx.Where("user_id = ? AND kind = ?", userID, kind).First(&existing).Error
		if err == nil {
			if err := tx.Model(&existing).
				Updates(map[string]interface{}{"destination": row.Destination, "verified_at": now}).Error; err != nil {
				return err
			}
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := tx.Create(&contact).Error; err != nil {
				return err
			}
		} else {
			return err
		}

		destination = row.Destination
		return nil
	})

	return destination, err
}

// ValidateRecoveryEmail refuses the two addresses that are not recoveries.
func ValidateRecoveryEmail(db *gorm.DB, userID, primary, candidate string) error {
	candidate = strings.ToLower(strings.TrimSpace(candidate))
	// Your own sign-in address is not a recovery path. If you have lost access
	// to it, sending the code there helps nobody.
	if candidate == strings.ToLower(strings.TrimSpace(primary)) {
		return ErrRecoverySameAsPrimary
	}
	// Somebody else's sign-in address is worse: it would let them reset into
	// this account.
	var count int64
	if err := db.Model(&models.User{}).
		Where("LOWER(email) = ? AND id <> ?", candidate, userID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrRecoveryInUse
	}
	return nil
}

// ClearRecoveryContact removes a verified contact and burns its outstanding
// codes.
func ClearRecoveryContact(db *gorm.DB, userID string, kind models.RecoveryContactKind) error {
	return db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		// Outstanding codes die with the contact, or a code sent moments before
		// removal could re-add it.
		if err := tx.Model(&models.RecoveryContactToken{}).
			Where("user_id = ? AND kind = ? AND used_at IS NULL", userID, kind).
			Update("used_at", now).Error; err != nil {
			return err
		}
		return tx.Where("user_id = ? AND kind = ?", userID, kind).
			Delete(&models.RecoveryContact{}).Error
	})
}

// LoadRecoveryContacts returns the verified contacts for a user, keyed by kind.
func LoadRecoveryContacts(db *gorm.DB, userID string) (map[models.RecoveryContactKind]models.RecoveryContact, error) {
	var rows []models.RecoveryContact
	if err := db.Where("user_id = ?", userID).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[models.RecoveryContactKind]models.RecoveryContact, len(rows))
	for _, r := range rows {
		out[r.Kind] = r
	}
	return out, nil
}
