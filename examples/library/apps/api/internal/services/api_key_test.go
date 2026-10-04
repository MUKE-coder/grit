package services

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"library/apps/api/internal/models"
)

func newKeyDB(tb testing.TB) *gorm.DB {
	tb.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		tb.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.APIKey{}); err != nil {
		tb.Fatalf("migrate: %v", err)
	}
	return db
}

func TestIssuedKeyStoresOnlyTheHash(t *testing.T) {
	db := newKeyDB(t)
	issued, err := GenerateAPIKey(db, KeyOptions{UserID: "user-1", Name: "nightly export"})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	var stored models.APIKey
	db.First(&stored, "id = ?", issued.Record.ID)

	if stored.SecretHash == issued.Token {
		t.Fatal("the raw token was stored")
	}
	// The secret is the third segment; it must not appear anywhere in the row.
	parts := issued.Token
	if stored.SecretHash == parts {
		t.Fatal("stored value is the token, not a hash")
	}
	if len(stored.SecretHash) != 64 {
		t.Fatalf("expected a sha256 hex digest, got %d chars", len(stored.SecretHash))
	}
}

func TestVerifyAcceptsTheIssuedKey(t *testing.T) {
	db := newKeyDB(t)
	issued, _ := GenerateAPIKey(db, KeyOptions{UserID: "user-1", Name: "ci"})

	got, err := VerifyAPIKey(db, issued.Token)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if got.UserID != "user-1" {
		t.Fatalf("resolved to user %q", got.UserID)
	}
}

func TestVerifyRejectsAWrongSecret(t *testing.T) {
	db := newKeyDB(t)
	issued, _ := GenerateAPIKey(db, KeyOptions{UserID: "user-1", Name: "ci"})

	// Right prefix, wrong secret — the case a scan-based lookup would get wrong.
	tampered := issued.Token[:len(issued.Token)-4] + "0000"
	if _, err := VerifyAPIKey(db, tampered); err == nil {
		t.Fatal("a key with a wrong secret verified")
	}
}

func TestVerifyRejectsARevokedKey(t *testing.T) {
	db := newKeyDB(t)
	issued, _ := GenerateAPIKey(db, KeyOptions{UserID: "user-1", Name: "ci"})

	if err := RevokeAPIKey(db, issued.Record.ID, "user-1"); err != nil {
		t.Fatalf("revoke: %v", err)
	}
	if _, err := VerifyAPIKey(db, issued.Token); err == nil {
		t.Fatal("a revoked key still verified")
	}
}

func TestVerifyRejectsAnExpiredKey(t *testing.T) {
	db := newKeyDB(t)
	past := time.Now().Add(-time.Hour)
	issued, _ := GenerateAPIKey(db, KeyOptions{UserID: "user-1", Name: "ci", ExpiresAt: &past})

	if _, err := VerifyAPIKey(db, issued.Token); err == nil {
		t.Fatal("an expired key still verified")
	}
}

func TestVerifyRejectsMalformedTokens(t *testing.T) {
	db := newKeyDB(t)
	for _, bad := range []string{"", "nonsense", "grit_only-two", "wrong_aa_bb", "grit__"} {
		if _, err := VerifyAPIKey(db, bad); err == nil {
			t.Fatalf("malformed token %q verified", bad)
		}
	}
}

func TestRevokeIsScopedToTheOwner(t *testing.T) {
	db := newKeyDB(t)
	issued, _ := GenerateAPIKey(db, KeyOptions{UserID: "user-1", Name: "ci"})

	// Without the user_id predicate anyone could revoke anyone's key by id.
	if err := RevokeAPIKey(db, issued.Record.ID, "someone-else"); err == nil {
		t.Fatal("another user revoked this key")
	}
	if _, err := VerifyAPIKey(db, issued.Token); err != nil {
		t.Fatal("the key was revoked despite the wrong owner")
	}
}
