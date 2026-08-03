package scaffold

// Email verification.
//
// The User model has carried EmailVerifiedAt since the beginning, but only an
// identity provider ever set it — there was no way for a password signup to
// prove an address. That left a field that looks like a feature and is not
// one, which is worse than not having it.
//
// Modelled on password reset deliberately: same single-use hashed-token shape,
// same "issuing a new one burns the old", same purge helper. Two flows that
// differ only in what they authorise should not be implemented two ways.
//
// Enforcement is opt-in via REQUIRE_EMAIL_VERIFICATION. Defaulting it on would
// break every existing project on upgrade — their users all have a NULL
// EmailVerifiedAt and would be locked out at once.

func apiEmailVerifyModelGo() string {
	return `package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"{{MODULE}}/internal/ids"
	"gorm.io/gorm"
)

// EmailVerificationToken is a single-use, expiring proof that someone can read
// mail at an address. Only the hash is stored, so leaking this table does not
// let an attacker verify addresses they do not control.
type EmailVerificationToken struct {
	ID     string ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	UserID string ` + "`" + `gorm:"size:36;index;not null" json:"user_id"` + "`" + `

	TokenHash string ` + "`" + `gorm:"size:64;uniqueIndex;not null" json:"-"` + "`" + `

	// Email is recorded as it was when the token was issued. If the user
	// changes their address before clicking the link, the token must not
	// verify the new one — that is how you verify an address you never owned.
	Email string ` + "`" + `gorm:"size:255;not null" json:"email"` + "`" + `

	ExpiresAt time.Time  ` + "`" + `gorm:"index" json:"expires_at"` + "`" + `
	UsedAt    *time.Time ` + "`" + `gorm:"index" json:"used_at,omitempty"` + "`" + `
	CreatedAt time.Time  ` + "`" + `json:"created_at"` + "`" + `
}

func (t *EmailVerificationToken) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = ids.New()
	}
	return nil
}

// HashVerificationToken returns the storage form of a verification token.
func HashVerificationToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
`
}

func apiEmailVerifyServiceGo() string {
	return `package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// EmailVerificationTTL is generous on purpose: unlike a password reset, a
// verification link is often opened on a different device hours later.
var EmailVerificationTTL = 48 * time.Hour

var (
	ErrVerificationTokenInvalid = errors.New("email verification token is not valid")
	ErrEmailAlreadyVerified     = errors.New("email is already verified")
)

// GenerateVerificationToken returns a 256-bit URL-safe token.
func GenerateVerificationToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generating verification token: %w", err)
	}
	return hex.EncodeToString(buf), nil
}

// CreateEmailVerificationToken issues a token for a user, invalidating any
// earlier unused ones. Storing only the hash means this function is the single
// place the raw token exists — the caller must mail it and then forget it.
func CreateEmailVerificationToken(db *gorm.DB, userID, email, token string) (*models.EmailVerificationToken, error) {
	row := &models.EmailVerificationToken{
		UserID:    userID,
		Email:     email,
		TokenHash: models.HashVerificationToken(token),
		ExpiresAt: time.Now().Add(EmailVerificationTTL),
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		// Burn outstanding tokens: two live links for one address means the
		// older one still works after the user re-requests, which is exactly
		// what someone who intercepted the first mail wants.
		now := time.Now()
		if err := tx.Model(&models.EmailVerificationToken{}).
			Where("user_id = ? AND used_at IS NULL", userID).
			Update("used_at", now).Error; err != nil {
			return err
		}
		return tx.Create(row).Error
	})
	if err != nil {
		return nil, fmt.Errorf("creating verification token: %w", err)
	}
	return row, nil
}

// ConsumeEmailVerificationToken marks the user's email verified and returns
// their id. Single-use is enforced with a conditional UPDATE, so two
// concurrent requests cannot both succeed.
func ConsumeEmailVerificationToken(db *gorm.DB, token string) (string, error) {
	var userID string

	err := db.Transaction(func(tx *gorm.DB) error {
		var row models.EmailVerificationToken
		if err := tx.Where("token_hash = ?", models.HashVerificationToken(token)).
			First(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrVerificationTokenInvalid
			}
			return err
		}

		if row.UsedAt != nil || time.Now().After(row.ExpiresAt) {
			return ErrVerificationTokenInvalid
		}

		now := time.Now()
		res := tx.Model(&models.EmailVerificationToken{}).
			Where("id = ? AND used_at IS NULL", row.ID).
			Update("used_at", now)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrVerificationTokenInvalid // lost the race
		}

		// Only verify the address the token was issued for. If the user
		// changed their email in the meantime, this matches nothing and the
		// stale link cannot verify the new address.
		upd := tx.Model(&models.User{}).
			Where("id = ? AND email = ?", row.UserID, row.Email).
			Update("email_verified_at", now)
		if upd.Error != nil {
			return upd.Error
		}
		if upd.RowsAffected == 0 {
			return ErrVerificationTokenInvalid
		}

		userID = row.UserID
		return nil
	})
	if err != nil {
		return "", err
	}
	return userID, nil
}

// PurgeExpiredVerificationTokens removes spent and expired rows.
func PurgeExpiredVerificationTokens(db *gorm.DB, olderThan time.Duration) (int64, error) {
	cutoff := time.Now().Add(-olderThan)
	res := db.Where("expires_at < ? OR (used_at IS NOT NULL AND used_at < ?)", cutoff, cutoff).
		Delete(&models.EmailVerificationToken{})
	return res.RowsAffected, res.Error
}
`
}

func apiEmailVerifyTestGo() string {
	return `package services

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"{{MODULE}}/internal/models"
)

func newVerifyDB(tb testing.TB) *gorm.DB {
	tb.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		tb.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.EmailVerificationToken{}); err != nil {
		tb.Fatalf("migrate: %v", err)
	}
	return db
}

func seedUnverifiedUser(tb testing.TB, db *gorm.DB, email string) *models.User {
	tb.Helper()
	u := &models.User{FirstName: "Ada", LastName: "L", Email: email, Password: "x"}
	if err := db.Create(u).Error; err != nil {
		tb.Fatalf("seed user: %v", err)
	}
	return u
}

func TestVerificationTokenStoresOnlyTheHash(t *testing.T) {
	db := newVerifyDB(t)
	u := seedUnverifiedUser(t, db, "ada@example.com")

	row, err := CreateEmailVerificationToken(db, u.ID, u.Email, "raw-token-value")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if row.TokenHash == "raw-token-value" {
		t.Fatal("the raw token was stored — a leak of this table would verify addresses")
	}
	if row.TokenHash != models.HashVerificationToken("raw-token-value") {
		t.Fatal("stored value is not the hash of the token")
	}
}

func TestConsumeMarksTheUserVerified(t *testing.T) {
	db := newVerifyDB(t)
	u := seedUnverifiedUser(t, db, "ada@example.com")
	if _, err := CreateEmailVerificationToken(db, u.ID, u.Email, "tok"); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := ConsumeEmailVerificationToken(db, "tok")
	if err != nil {
		t.Fatalf("consume: %v", err)
	}
	if got != u.ID {
		t.Fatalf("returned user %q, want %q", got, u.ID)
	}

	var after models.User
	db.First(&after, "id = ?", u.ID)
	if after.EmailVerifiedAt == nil {
		t.Fatal("email_verified_at was not set")
	}
}

func TestVerificationTokenIsSingleUse(t *testing.T) {
	db := newVerifyDB(t)
	u := seedUnverifiedUser(t, db, "ada@example.com")
	if _, err := CreateEmailVerificationToken(db, u.ID, u.Email, "tok"); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := ConsumeEmailVerificationToken(db, "tok"); err != nil {
		t.Fatalf("first consume: %v", err)
	}
	if _, err := ConsumeEmailVerificationToken(db, "tok"); err == nil {
		t.Fatal("the token verified twice")
	}
}

func TestIssuingANewVerificationTokenBurnsTheOld(t *testing.T) {
	db := newVerifyDB(t)
	u := seedUnverifiedUser(t, db, "ada@example.com")
	if _, err := CreateEmailVerificationToken(db, u.ID, u.Email, "first"); err != nil {
		t.Fatalf("create first: %v", err)
	}
	if _, err := CreateEmailVerificationToken(db, u.ID, u.Email, "second"); err != nil {
		t.Fatalf("create second: %v", err)
	}
	if _, err := ConsumeEmailVerificationToken(db, "first"); err == nil {
		t.Fatal("the superseded token still worked")
	}
	if _, err := ConsumeEmailVerificationToken(db, "second"); err != nil {
		t.Fatalf("the current token failed: %v", err)
	}
}

// The case the Email column exists for: a stale link must not verify an
// address the user only just switched to.
func TestTokenDoesNotVerifyAChangedAddress(t *testing.T) {
	db := newVerifyDB(t)
	u := seedUnverifiedUser(t, db, "ada@example.com")
	if _, err := CreateEmailVerificationToken(db, u.ID, u.Email, "tok"); err != nil {
		t.Fatalf("create: %v", err)
	}

	db.Model(&models.User{}).Where("id = ?", u.ID).Update("email", "someone-else@example.com")

	if _, err := ConsumeEmailVerificationToken(db, "tok"); err == nil {
		t.Fatal("a stale link verified an address it was never issued for")
	}
}

func TestVerificationTokenExpires(t *testing.T) {
	db := newVerifyDB(t)
	u := seedUnverifiedUser(t, db, "ada@example.com")
	row, err := CreateEmailVerificationToken(db, u.ID, u.Email, "tok")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	db.Model(&models.EmailVerificationToken{}).
		Where("id = ?", row.ID).
		Update("expires_at", time.Now().Add(-time.Minute))

	if _, err := ConsumeEmailVerificationToken(db, "tok"); err == nil {
		t.Fatal("an expired token was accepted")
	}
}

func TestUnknownVerificationTokenIsRejected(t *testing.T) {
	db := newVerifyDB(t)
	if _, err := ConsumeEmailVerificationToken(db, "never-issued"); err == nil {
		t.Fatal("an unknown token was accepted")
	}
}
`
}
