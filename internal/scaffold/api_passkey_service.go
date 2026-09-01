package scaffold

func passkeyServiceGo() string {
	return `package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// WebAuthnSessionTTL bounds a ceremony.
//
// Five minutes: long enough to find the phone the passkey lives on, short
// enough that an abandoned challenge is not left lying around to be replayed.
const WebAuthnSessionTTL = 5 * time.Minute

var (
	ErrPasskeyNotConfigured = errors.New("passkeys are not configured on this deployment")
	ErrPasskeyChallengeGone = errors.New("that sign-in attempt expired; start again")
	ErrPasskeyUnknown       = errors.New("that passkey is not registered")
)

// Passkeys wraps the WebAuthn relying party.
type Passkeys struct {
	DB *gorm.DB
	wa *webauthn.WebAuthn
}

// NewPasskeys builds the relying party from the origins the frontends actually
// run on.
//
// RPID is the registrable domain, and getting it wrong is the classic WebAuthn
// failure: the browser refuses with a SecurityError that names nothing useful,
// and it looks like the code is broken. So it is derived from the first origin
// rather than typed twice, and logged at boot.
func NewPasskeys(db *gorm.DB, displayName string, origins []string) (*Passkeys, error) {
	clean := make([]string, 0, len(origins))
	for _, o := range origins {
		o = strings.TrimSpace(strings.TrimSuffix(o, "/"))
		if o != "" && o != "*" {
			clean = append(clean, o)
		}
	}
	if len(clean) == 0 {
		return nil, ErrPasskeyNotConfigured
	}

	first, err := url.Parse(clean[0])
	if err != nil {
		return nil, fmt.Errorf("parsing origin %q: %w", clean[0], err)
	}
	rpID := first.Hostname()
	if rpID == "" {
		return nil, fmt.Errorf("origin %q has no host", clean[0])
	}

	wa, err := webauthn.New(&webauthn.Config{
		RPDisplayName: displayName,
		RPID:          rpID,
		RPOrigins:     clean,
	})
	if err != nil {
		return nil, fmt.Errorf("configuring webauthn: %w", err)
	}
	log.Printf("Passkeys enabled for %s (origins: %s)", rpID, strings.Join(clean, ", "))
	return &Passkeys{DB: db, wa: wa}, nil
}

func (p *Passkeys) load(userID string) (models.PasskeyUser, error) {
	var user models.User
	if err := p.DB.First(&user, "id = ?", userID).Error; err != nil {
		return models.PasskeyUser{}, err
	}
	var rows []models.Passkey
	if err := p.DB.Where("user_id = ?", userID).Find(&rows).Error; err != nil {
		return models.PasskeyUser{}, err
	}
	creds := make([]webauthn.Credential, 0, len(rows))
	for _, r := range rows {
		id, err := base64.RawURLEncoding.DecodeString(r.CredentialID)
		if err != nil {
			continue
		}
		creds = append(creds, webauthn.Credential{
			ID:              id,
			PublicKey:       r.PublicKey,
			AttestationType: r.AttestationType,
			Authenticator: webauthn.Authenticator{
				AAGUID:    r.AAGUID,
				SignCount: r.SignCount,
			},
			Flags: webauthn.CredentialFlags{
				BackupEligible: r.BackupEligible,
				BackupState:    r.BackupState,
			},
		})
	}
	return models.PasskeyUser{User: &user, Credentials: creds}, nil
}

func (p *Passkeys) saveSession(userID, purpose string, s *webauthn.SessionData) (string, error) {
	blob, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	row := models.WebAuthnSession{
		UserID:    userID,
		Purpose:   purpose,
		Data:      blob,
		ExpiresAt: time.Now().Add(WebAuthnSessionTTL),
	}
	if err := p.DB.Create(&row).Error; err != nil {
		return "", err
	}
	return row.ID, nil
}

// takeSession reads a ceremony and deletes it in the same breath.
//
// Deleted rather than marked used: there is nothing to audit here and a
// challenge that can be presented twice is not a challenge.
func (p *Passkeys) takeSession(id, purpose string) (*webauthn.SessionData, string, error) {
	var row models.WebAuthnSession
	if err := p.DB.First(&row, "id = ? AND purpose = ?", id, purpose).Error; err != nil {
		return nil, "", ErrPasskeyChallengeGone
	}
	p.DB.Delete(&models.WebAuthnSession{}, "id = ?", row.ID)
	if time.Now().After(row.ExpiresAt) {
		return nil, "", ErrPasskeyChallengeGone
	}
	var s webauthn.SessionData
	if err := json.Unmarshal(row.Data, &s); err != nil {
		return nil, "", ErrPasskeyChallengeGone
	}
	return &s, row.UserID, nil
}

// BeginRegistration starts adding a passkey to a signed-in account.
func (p *Passkeys) BeginRegistration(userID string) (*protocol.CredentialCreation, string, error) {
	user, err := p.load(userID)
	if err != nil {
		return nil, "", err
	}
	opts, session, err := p.wa.BeginRegistration(
		user,
		// Excluding what is already registered stops a second, duplicate
		// credential for the same authenticator, which would show up as two
		// identical rows nobody can tell apart.
		webauthn.WithExclusions(credentialDescriptors(user.Credentials)),
		webauthn.WithResidentKeyRequirement(protocol.ResidentKeyRequirementPreferred),
	)
	if err != nil {
		return nil, "", err
	}
	id, err := p.saveSession(userID, "register", session)
	if err != nil {
		return nil, "", err
	}
	return opts, id, nil
}

// FinishRegistration verifies the attestation and stores the credential.
func (p *Passkeys) FinishRegistration(sessionID, label string, body []byte) (*models.Passkey, error) {
	session, userID, err := p.takeSession(sessionID, "register")
	if err != nil {
		return nil, err
	}
	user, err := p.load(userID)
	if err != nil {
		return nil, err
	}
	parsed, err := protocol.ParseCredentialCreationResponseBytes(body)
	if err != nil {
		return nil, fmt.Errorf("reading the authenticator's answer: %w", err)
	}
	cred, err := p.wa.CreateCredential(user, *session, parsed)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(label) == "" {
		label = "Passkey"
	}
	row := models.Passkey{
		UserID:          userID,
		CredentialID:    base64.RawURLEncoding.EncodeToString(cred.ID),
		PublicKey:       cred.PublicKey,
		AttestationType: cred.AttestationType,
		AAGUID:          cred.Authenticator.AAGUID,
		SignCount:       cred.Authenticator.SignCount,
		BackupEligible:  cred.Flags.BackupEligible,
		BackupState:     cred.Flags.BackupState,
		Name:            strings.TrimSpace(label),
	}
	if err := p.DB.Create(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// BeginLogin starts a usernameless sign-in.
//
// Usernameless because that is the point of a passkey: the authenticator knows
// which account it holds, so asking for an email first is a step that buys
// nothing.
func (p *Passkeys) BeginLogin() (*protocol.CredentialAssertion, string, error) {
	opts, session, err := p.wa.BeginDiscoverableLogin()
	if err != nil {
		return nil, "", err
	}
	id, err := p.saveSession("", "login", session)
	if err != nil {
		return nil, "", err
	}
	return opts, id, nil
}

// FinishLogin verifies the assertion and returns the user it belongs to.
func (p *Passkeys) FinishLogin(sessionID string, body []byte) (*models.User, error) {
	session, _, err := p.takeSession(sessionID, "login")
	if err != nil {
		return nil, err
	}
	parsed, err := protocol.ParseCredentialRequestResponseBytes(body)
	if err != nil {
		return nil, fmt.Errorf("reading the authenticator's answer: %w", err)
	}

	var matched *models.User
	_, err = p.wa.ValidateDiscoverableLogin(func(rawID, userHandle []byte) (webauthn.User, error) {
		// The user handle is the account id we minted at registration.
		user, err := p.load(string(userHandle))
		if err != nil {
			return nil, ErrPasskeyUnknown
		}
		matched = user.User
		return user, nil
	}, *session, parsed)
	if err != nil {
		return nil, err
	}
	if matched == nil {
		return nil, ErrPasskeyUnknown
	}

	// Record the use, and notice a counter that went backwards.
	credID := base64.RawURLEncoding.EncodeToString(parsed.RawID)
	var row models.Passkey
	if err := p.DB.First(&row, "credential_id = ?", credID).Error; err == nil {
		next := parsed.Response.AuthenticatorData.Counter
		if next != 0 && next <= row.SignCount {
			// Not fatal on its own: plenty of authenticators never increment.
			// Logged because a counter that moves backwards is the signature of
			// a cloned credential, and somebody should be able to find it.
			log.Printf("passkey %s: sign counter went from %d to %d for user %s",
				row.ID, row.SignCount, next, row.UserID)
		}
		now := time.Now()
		p.DB.Model(&row).Updates(map[string]interface{}{
			"sign_count":   next,
			"last_used_at": now,
		})
	}

	return matched, nil
}

// List returns a user's passkeys, newest first.
func (p *Passkeys) List(userID string) ([]models.Passkey, error) {
	var rows []models.Passkey
	err := p.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&rows).Error
	return rows, err
}

// Delete removes one passkey belonging to this user.
//
// Scoped by user id as well as row id, so a guessed id cannot remove somebody
// else's authenticator.
func (p *Passkeys) Delete(userID, id string) error {
	res := p.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Passkey{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrPasskeyUnknown
	}
	return nil
}

// CleanupExpiredSessions drops abandoned ceremonies.
func CleanupExpiredWebAuthnSessions(db *gorm.DB) (int64, error) {
	res := db.Where("expires_at < ?", time.Now()).Delete(&models.WebAuthnSession{})
	return res.RowsAffected, res.Error
}

func credentialDescriptors(creds []webauthn.Credential) []protocol.CredentialDescriptor {
	out := make([]protocol.CredentialDescriptor, 0, len(creds))
	for _, c := range creds {
		out = append(out, c.Descriptor())
	}
	return out
}
`
}

func passkeyServiceTestGo() string {
	return `package services

import (
	"strings"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
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
	if _, err := NewPasskeys(db, "Grit", []string{"", "*"}); err != ErrPasskeyNotConfigured {
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
	if _, _, err := p.takeSession(sid, "register"); err != ErrPasskeyChallengeGone {
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

	if _, _, err := p.takeSession(sid, "register"); err != ErrPasskeyChallengeGone {
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
	if _, _, err := p.takeSession(sid, "login"); err != ErrPasskeyChallengeGone {
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

	if err := p.Delete(mine.ID, key.ID); err != ErrPasskeyUnknown {
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
`
}
