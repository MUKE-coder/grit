package scaffold

func recoveryServiceGo() string {
	return `package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
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
`
}

func recoveryServiceTestGo() string {
	return `package services

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

func recoveryDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&_foreign_keys=on"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.RecoveryContactToken{}, &models.RecoveryContact{}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	})
	return db
}

func seedUser(t *testing.T, db *gorm.DB, email string) *models.User {
	t.Helper()
	u := &models.User{FirstName: "A", LastName: "B", Email: email, Active: true}
	if err := db.Create(u).Error; err != nil {
		t.Fatal(err)
	}
	return u
}

func TestRecoveryCodeVerifiesAndStoresTheContact(t *testing.T) {
	db := recoveryDB(t)
	u := seedUser(t, db, "primary@example.com")

	code, err := NewRecoveryCode()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := CreateRecoveryToken(db, u.ID, models.RecoveryEmail, "backup@example.com", code); err != nil {
		t.Fatal(err)
	}

	got, err := ConsumeRecoveryToken(db, u.ID, models.RecoveryEmail, code)
	if err != nil {
		t.Fatalf("the right code should verify: %v", err)
	}
	if got != "backup@example.com" {
		t.Errorf("verified the wrong destination: %s", got)
	}

	contacts, err := LoadRecoveryContacts(db, u.ID)
	if err != nil {
		t.Fatal(err)
	}
	got2, ok := contacts[models.RecoveryEmail]
	if !ok || got2.Destination != "backup@example.com" || got2.VerifiedAt.IsZero() {
		t.Errorf("the contact was not stored as verified: %+v", contacts)
	}
}

// A code that has been used must not work twice.
func TestRecoveryCodeIsSingleUse(t *testing.T) {
	db := recoveryDB(t)
	u := seedUser(t, db, "primary@example.com")

	code, _ := NewRecoveryCode()
	if _, err := CreateRecoveryToken(db, u.ID, models.RecoveryEmail, "backup@example.com", code); err != nil {
		t.Fatal(err)
	}
	if _, err := ConsumeRecoveryToken(db, u.ID, models.RecoveryEmail, code); err != nil {
		t.Fatal(err)
	}
	if _, err := ConsumeRecoveryToken(db, u.ID, models.RecoveryEmail, code); err == nil {
		t.Error("a spent code must not verify again")
	}
}

// Requesting a second code must kill the first, or an intercepted message
// stays useful after the user re-requests.
func TestRequestingAgainBurnsTheOldCode(t *testing.T) {
	db := recoveryDB(t)
	u := seedUser(t, db, "primary@example.com")

	first, _ := NewRecoveryCode()
	if _, err := CreateRecoveryToken(db, u.ID, models.RecoveryEmail, "backup@example.com", first); err != nil {
		t.Fatal(err)
	}
	second, _ := NewRecoveryCode()
	if _, err := CreateRecoveryToken(db, u.ID, models.RecoveryEmail, "backup@example.com", second); err != nil {
		t.Fatal(err)
	}

	if _, err := ConsumeRecoveryToken(db, u.ID, models.RecoveryEmail, first); err == nil {
		t.Error("the superseded code must stop working")
	}
	if _, err := ConsumeRecoveryToken(db, u.ID, models.RecoveryEmail, second); err != nil {
		t.Errorf("the current code should still work: %v", err)
	}
}

// Six digits is a million options, which only helps if guesses are capped.
func TestGuessesAreCapped(t *testing.T) {
	db := recoveryDB(t)
	u := seedUser(t, db, "primary@example.com")

	code, _ := NewRecoveryCode()
	if _, err := CreateRecoveryToken(db, u.ID, models.RecoveryEmail, "backup@example.com", code); err != nil {
		t.Fatal(err)
	}

	wrong := "000000"
	if wrong == code {
		wrong = "111111"
	}
	for i := 0; i < MaxRecoveryAttempts; i++ {
		if _, err := ConsumeRecoveryToken(db, u.ID, models.RecoveryEmail, wrong); err == nil {
			t.Fatal("a wrong code must not verify")
		}
	}
	// The cap is spent, so even the right code is now refused. The user
	// requests a new one, which is the intended recovery from this state.
	if _, err := ConsumeRecoveryToken(db, u.ID, models.RecoveryEmail, code); err == nil {
		t.Error("the code should be dead once the attempt cap is reached")
	}
}

// Your own sign-in address is not a recovery path.
func TestPrimaryAddressIsRefused(t *testing.T) {
	db := recoveryDB(t)
	u := seedUser(t, db, "primary@example.com")

	if err := ValidateRecoveryEmail(db, u.ID, u.Email, "Primary@Example.com"); err != ErrRecoverySameAsPrimary {
		t.Errorf("the primary address must be refused, case-insensitively: %v", err)
	}
}

// Somebody else's sign-in address is worse: it would let them reset into this
// account.
func TestAnotherUsersAddressIsRefused(t *testing.T) {
	db := recoveryDB(t)
	mine := seedUser(t, db, "mine@example.com")
	seedUser(t, db, "theirs@example.com")

	if err := ValidateRecoveryEmail(db, mine.ID, mine.Email, "theirs@example.com"); err != ErrRecoveryInUse {
		t.Errorf("another account's address must be refused: %v", err)
	}
}

func TestClearingRemovesTheContact(t *testing.T) {
	db := recoveryDB(t)
	u := seedUser(t, db, "primary@example.com")

	code, _ := NewRecoveryCode()
	if _, err := CreateRecoveryToken(db, u.ID, models.RecoveryEmail, "backup@example.com", code); err != nil {
		t.Fatal(err)
	}
	if _, err := ConsumeRecoveryToken(db, u.ID, models.RecoveryEmail, code); err != nil {
		t.Fatal(err)
	}
	if err := ClearRecoveryContact(db, u.ID, models.RecoveryEmail); err != nil {
		t.Fatal(err)
	}

	contacts, err := LoadRecoveryContacts(db, u.ID)
	if err != nil {
		t.Fatal(err)
	}
	if _, still := contacts[models.RecoveryEmail]; still {
		t.Errorf("the contact should be gone: %+v", contacts)
	}
}
`
}
