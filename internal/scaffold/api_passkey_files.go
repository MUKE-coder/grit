package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writePasskeyFiles writes the WebAuthn passkey feature.
func writePasskeyFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "models", "passkey.go"):          passkeyModelGo(),
		filepath.Join(apiRoot, "internal", "services", "passkey.go"):        passkeyServiceGo(),
		filepath.Join(apiRoot, "internal", "services", "passkey_test.go"):   passkeyServiceTestGo(),
		filepath.Join(apiRoot, "internal", "handlers", "passkey.go"):        passkeyHandlerGo(),
	}
	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

func passkeyModelGo() string {
	return `package models

import (
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"

	"{{MODULE}}/internal/ids"
)

// Passkey is one registered WebAuthn credential: a phone, a laptop's biometric
// sensor, a hardware key.
//
// One row per authenticator, not one per user. Somebody registers their laptop
// and their phone and expects both to work, and losing the laptop should not
// mean losing the account.
type Passkey struct {
	ID     string ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	UserID string ` + "`" + `gorm:"size:36;index;not null" json:"user_id"` + "`" + `

	// CredentialID is what the authenticator calls itself, base64url encoded.
	// Unique because it is the lookup key on sign-in, and because registering
	// the same authenticator twice should update rather than duplicate.
	CredentialID string ` + "`" + `gorm:"size:512;uniqueIndex;not null" json:"credential_id"` + "`" + `

	PublicKey       []byte ` + "`" + `gorm:"type:blob;not null" json:"-"` + "`" + `
	AttestationType string ` + "`" + `gorm:"size:64" json:"-"` + "`" + `
	AAGUID          []byte ` + "`" + `gorm:"type:blob" json:"-"` + "`" + `
	Transports      string ` + "`" + `gorm:"size:255" json:"transports,omitempty"` + "`" + `

	// SignCount is the authenticator's own counter.
	//
	// Kept because a counter that goes backwards means two authenticators are
	// answering for one credential, which is the signature of a cloned key. Not
	// every authenticator implements it (many report zero forever), so it is a
	// signal rather than a gate.
	SignCount uint32 ` + "`" + `gorm:"default:0" json:"-"` + "`" + `

	BackupEligible bool ` + "`" + `gorm:"default:false" json:"backup_eligible"` + "`" + `
	BackupState    bool ` + "`" + `gorm:"default:false" json:"synced"` + "`" + `

	// Name is what the owner sees in the list. Defaulted from the user agent at
	// registration, because "Passkey 1" tells nobody which device to go and get.
	Name string ` + "`" + `gorm:"size:120" json:"name"` + "`" + `

	LastUsedAt *time.Time ` + "`" + `json:"last_used_at,omitempty"` + "`" + `
	CreatedAt  time.Time  ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt  time.Time  ` + "`" + `json:"updated_at"` + "`" + `
}

func (p *Passkey) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = ids.New()
	}
	return nil
}

// WebAuthnSession holds one in-flight ceremony.
//
// In a table rather than in memory, for the same reason refresh sessions are:
// the moment there are two API instances, an in-memory challenge is a coin
// flip on whether sign-in works. Single use and short lived, because a replayed
// challenge is the thing the challenge exists to prevent.
type WebAuthnSession struct {
	ID string ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `

	// UserID is empty for a usernameless sign-in, where the authenticator tells
	// us who it is only at the end of the ceremony.
	UserID string ` + "`" + `gorm:"size:36;index" json:"user_id,omitempty"` + "`" + `

	Purpose   string    ` + "`" + `gorm:"size:16;not null" json:"purpose"` + "`" + `
	Data      []byte    ` + "`" + `gorm:"type:blob;not null" json:"-"` + "`" + `
	ExpiresAt time.Time ` + "`" + `gorm:"index" json:"expires_at"` + "`" + `
	CreatedAt time.Time ` + "`" + `json:"created_at"` + "`" + `
}

func (s *WebAuthnSession) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = ids.New()
	}
	return nil
}

// PasskeyUser adapts a User to what the webauthn library expects.
//
// It carries the credentials with it because the library asks for them during
// both ceremonies: to exclude already-registered authenticators when
// registering, and to know what to verify against when signing in.
type PasskeyUser struct {
	User        *User
	Credentials []webauthn.Credential
}

func (u PasskeyUser) WebAuthnID() []byte { return []byte(u.User.ID) }

func (u PasskeyUser) WebAuthnName() string { return u.User.Email }

func (u PasskeyUser) WebAuthnDisplayName() string {
	name := u.User.FirstName + " " + u.User.LastName
	if name == " " {
		return u.User.Email
	}
	return name
}

func (u PasskeyUser) WebAuthnCredentials() []webauthn.Credential { return u.Credentials }
`
}
