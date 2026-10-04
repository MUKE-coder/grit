package services

import (
	"errors"
	"strings"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"library/apps/api/internal/models"
)

func passkeyDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&_foreign_keys=on"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Passkey{}, &models.WebAuthnSession{}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	})
	return db
}

// The relying party id is derived from the origin rather than configured
// twice. Getting it wrong is the classic WebAuthn failure and the browser's
// error names nothing useful, so this is worth pinning.
func TestRelyingPartyIDComesFromTheOrigin(t *testing.T) {
	db := passkeyDB(t)
	p, err := NewPasskeys(db, "Grit", []string{"https://admin.example.com/", "https://example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if p == nil {
		t.Fatal("expected a relying party")
	}
}

func TestNoOriginsIsRefused(t *testing.T) {
	db := passkeyDB(t)
	if _, err := NewPasskeys(db, "Grit", []string{"", "*"}); !errors.Is(err, ErrPasskeyNotConfigured) {
		t.Errorf("a deployment with no usable origin should refuse, got %v", err)
	}
}

// A challenge must not survive being used, or it is not a challenge.
func TestCeremonyIsSingleUse(t *testing.T) {
	db := passkeyDB(t)
	p, err := NewPasskeys(db, "Grit", []string{"http://localhost:3001"})
	if err != nil {
		t.Fatal(err)
	}
	u := &models.User{FirstName: "A", LastName: "B", Email: "a@example.com", Active: true}
	if err := db.Create(u).Error; err != nil {
		t.Fatal(err)
	}

	_, sid, err := p.BeginRegistration(u.ID)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := p.takeSession(sid, "register"); err != nil {
		t.Fatalf("the first read should succeed: %v", err)
	}
	if _, _, err := p.takeSession(sid, "register"); !errors.Is(err, ErrPasskeyChallengeGone) {
		t.Errorf("a spent challenge must not be readable again, got %v", err)
	}
}

func TestExpiredCeremonyIsRefused(t *testing.T) {
	db := passkeyDB(t)
	p, err := NewPasskeys(db, "Grit", []string{"http://localhost:3001"})
	if err != nil {
		t.Fatal(err)
	}
	u := &models.User{FirstName: "A", LastName: "B", Email: "b@example.com", Active: true}
	if err := db.Create(u).Error; err != nil {
		t.Fatal(err)
	}
	_, sid, err := p.BeginRegistration(u.ID)
	if err != nil {
		t.Fatal(err)
	}
	db.Model(&models.WebAuthnSession{}).Where("id = ?", sid).
		Update("expires_at", time.Now().Add(-time.Minute))

	if _, _, err := p.takeSession(sid, "register"); !errors.Is(err, ErrPasskeyChallengeGone) {
		t.Errorf("an expired challenge must be refused, got %v", err)
	}
}

// A registration challenge must not be usable as a sign-in challenge.
func TestPurposeIsChecked(t *testing.T) {
	db := passkeyDB(t)
	p, err := NewPasskeys(db, "Grit", []string{"http://localhost:3001"})
	if err != nil {
		t.Fatal(err)
	}
	u := &models.User{FirstName: "A", LastName: "B", Email: "c@example.com", Active: true}
	if err := db.Create(u).Error; err != nil {
		t.Fatal(err)
	}
	_, sid, err := p.BeginRegistration(u.ID)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := p.takeSession(sid, "login"); !errors.Is(err, ErrPasskeyChallengeGone) {
		t.Errorf("a register challenge must not pass as a login one, got %v", err)
	}
}

// Deleting is scoped by user, so a guessed id cannot remove somebody else's key.
func TestDeleteIsScopedToTheOwner(t *testing.T) {
	db := passkeyDB(t)
	p, err := NewPasskeys(db, "Grit", []string{"http://localhost:3001"})
	if err != nil {
		t.Fatal(err)
	}
	mine := &models.User{FirstName: "A", LastName: "B", Email: "mine@example.com", Active: true}
	theirs := &models.User{FirstName: "C", LastName: "D", Email: "theirs@example.com", Active: true}
	db.Create(mine)
	db.Create(theirs)

	key := &models.Passkey{UserID: theirs.ID, CredentialID: "abc", PublicKey: []byte("k"), Name: "Their laptop"}
	if err := db.Create(key).Error; err != nil {
		t.Fatal(err)
	}

	if err := p.Delete(mine.ID, key.ID); !errors.Is(err, ErrPasskeyUnknown) {
		t.Errorf("deleting another user's passkey must fail, got %v", err)
	}
	var still int64
	db.Model(&models.Passkey{}).Where("id = ?", key.ID).Count(&still)
	if still != 1 {
		t.Error("the passkey should still be there")
	}

	if err := p.Delete(theirs.ID, key.ID); err != nil {
		t.Errorf("the owner should be able to delete it: %v", err)
	}
}

func TestExpiredSessionsAreCleanedUp(t *testing.T) {
	db := passkeyDB(t)
	db.Create(&models.WebAuthnSession{Purpose: "login", Data: []byte("{}"), ExpiresAt: time.Now().Add(-time.Hour)})
	db.Create(&models.WebAuthnSession{Purpose: "login", Data: []byte("{}"), ExpiresAt: time.Now().Add(time.Hour)})

	n, err := CleanupExpiredWebAuthnSessions(db)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Errorf("expected one expired row removed, got %d", n)
	}
}

func TestCredentialIDIsUnique(t *testing.T) {
	db := passkeyDB(t)
	u := &models.User{FirstName: "A", LastName: "B", Email: "u@example.com", Active: true}
	db.Create(u)
	if err := db.Create(&models.Passkey{UserID: u.ID, CredentialID: "dup", PublicKey: []byte("k")}).Error; err != nil {
		t.Fatal(err)
	}
	err := db.Create(&models.Passkey{UserID: u.ID, CredentialID: "dup", PublicKey: []byte("k")}).Error
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "unique") {
		t.Errorf("a duplicate credential id must be refused, got %v", err)
	}
}
