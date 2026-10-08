package services

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"saas/apps/api/internal/models"
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
