package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/services"
)

func magicLinkDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db handle: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&models.MagicLinkToken{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func TestAMagicLinkSignsInOnce(t *testing.T) {
	db := magicLinkDB(t)
	ctx := context.Background()

	token, err := services.IssueMagicLink(ctx, db, "user-1", "127.0.0.1", "test")
	if err != nil {
		t.Fatal(err)
	}
	if len(token) < 32 {
		t.Errorf("token %q is too short to be unguessable", token)
	}

	// The plaintext is never stored: a table of live tokens is a table of
	// passwords.
	var row models.MagicLinkToken
	if err := db.First(&row).Error; err != nil {
		t.Fatal(err)
	}
	if row.TokenHash == token {
		t.Error("the token is stored in the clear")
	}

	userID, err := services.ConsumeMagicLink(ctx, db, token)
	if err != nil || userID != "user-1" {
		t.Fatalf("consume = %q, %v", userID, err)
	}

	if _, err := services.ConsumeMagicLink(ctx, db, token); !errors.Is(err, services.ErrMagicLinkUsed) {
		t.Errorf("a spent link signed somebody in again: %v", err)
	}
}

// Clicking the same email twice is an ordinary mistake, and deserves an answer
// that says so rather than "invalid link".
func TestASpentLinkSaysItIsSpent(t *testing.T) {
	db := magicLinkDB(t)
	ctx := context.Background()

	token, _ := services.IssueMagicLink(ctx, db, "user-1", "", "")
	if _, err := services.ConsumeMagicLink(ctx, db, token); err != nil {
		t.Fatal(err)
	}
	_, err := services.ConsumeMagicLink(ctx, db, token)
	if !errors.Is(err, services.ErrMagicLinkUsed) {
		t.Errorf("err = %v, want ErrMagicLinkUsed", err)
	}
	if _, err := services.ConsumeMagicLink(ctx, db, "never-existed"); !errors.Is(err, services.ErrMagicLinkInvalid) {
		t.Errorf("an unknown token should be invalid, got %v", err)
	}
}

func TestAnExpiredLinkIsRefused(t *testing.T) {
	db := magicLinkDB(t)
	ctx := context.Background()

	token, _ := services.IssueMagicLink(ctx, db, "user-1", "", "")
	if err := db.Model(&models.MagicLinkToken{}).Where("1 = 1").
		Update("expires_at", time.Now().Add(-time.Minute)).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := services.ConsumeMagicLink(ctx, db, token); !errors.Is(err, services.ErrMagicLinkExpired) {
		t.Errorf("an expired link was accepted: %v", err)
	}
}

// Without a rate limit the form is a way to send somebody a hundred emails.
func TestLinksAreRateLimitedPerAccount(t *testing.T) {
	db := magicLinkDB(t)
	ctx := context.Background()

	if _, err := services.IssueMagicLink(ctx, db, "user-1", "", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := services.IssueMagicLink(ctx, db, "user-1", "", ""); !errors.Is(err, services.ErrMagicLinkTooSoon) {
		t.Errorf("a second link was issued immediately: %v", err)
	}
	// A different account is unaffected.
	if _, err := services.IssueMagicLink(ctx, db, "user-2", "", ""); err != nil {
		t.Errorf("one account's link blocked another's: %v", err)
	}
}
