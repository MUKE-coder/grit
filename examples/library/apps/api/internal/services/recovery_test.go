package services

import (
	"errors"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"library/apps/api/internal/models"
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

	if err := ValidateRecoveryEmail(db, u.ID, u.Email, "Primary@Example.com"); !errors.Is(err, ErrRecoverySameAsPrimary) {
		t.Errorf("the primary address must be refused, case-insensitively: %v", err)
	}
}

// Somebody else's sign-in address is worse: it would let them reset into this
// account.
func TestAnotherUsersAddressIsRefused(t *testing.T) {
	db := recoveryDB(t)
	mine := seedUser(t, db, "mine@example.com")
	seedUser(t, db, "theirs@example.com")

	if err := ValidateRecoveryEmail(db, mine.ID, mine.Email, "theirs@example.com"); !errors.Is(err, ErrRecoveryInUse) {
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
