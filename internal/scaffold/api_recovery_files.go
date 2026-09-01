package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeRecoveryFiles writes the recovery-contact feature: a verified secondary
// email (and a phone, when an SMS provider is wired) that can get somebody back
// into an account when the primary address is gone.
func writeRecoveryFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "sms", "sms.go"):                    smsGo(),
		filepath.Join(apiRoot, "internal", "sms", "sms_test.go"):               smsTestGo(),
		filepath.Join(apiRoot, "internal", "models", "recovery_contact.go"):    recoveryModelGo(),
		filepath.Join(apiRoot, "internal", "services", "recovery.go"):          recoveryServiceGo(),
		filepath.Join(apiRoot, "internal", "services", "recovery_test.go"):     recoveryServiceTestGo(),
		filepath.Join(apiRoot, "internal", "handlers", "recovery.go"):          recoveryHandlerGo(),
	}
	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

func smsGo() string {
	return `// Package sms is the seam an SMS provider plugs into.
//
// There is no default implementation, deliberately. Every provider is a paid
// account with its own credentials, and the right one depends on where your
// users are: Twilio is not the sensible choice in Kampala, and Africa's Talking
// is not the sensible choice in Berlin. Baking one in would make everybody
// carry a dependency most of them cannot use.
//
// So the framework defines the interface and the features that need SMS ask
// whether anything is registered. Nothing that needs a text message is offered
// in the UI until something is.
//
// Wire one in main.go:
//
//	sms.Register(sms.SenderFunc(func(ctx context.Context, to, body string) error {
//	    // call your provider here
//	    return nil
//	}))
package sms

import (
	"context"
	"errors"
	"sync"
)

// ErrNoProvider is returned by Send when nothing has been registered. Callers
// should check Configured() first and not offer the feature at all, so this is
// the backstop rather than the expected path.
var ErrNoProvider = errors.New("no SMS provider is configured")

// Sender delivers one text message.
type Sender interface {
	Send(ctx context.Context, to, body string) error
}

// SenderFunc adapts a plain function to Sender.
type SenderFunc func(ctx context.Context, to, body string) error

func (f SenderFunc) Send(ctx context.Context, to, body string) error { return f(ctx, to, body) }

var (
	mu      sync.RWMutex
	current Sender
)

// Register installs the provider. Call it once, from main, before serving.
func Register(s Sender) {
	mu.Lock()
	defer mu.Unlock()
	current = s
}

// Configured reports whether a provider is installed.
//
// The security overview endpoint reports this to the client, so the admin can
// hide phone recovery rather than offering a button that cannot work. A
// disabled control with no explanation is worse than no control.
func Configured() bool {
	mu.RLock()
	defer mu.RUnlock()
	return current != nil
}

// Send delivers a message through the registered provider.
func Send(ctx context.Context, to, body string) error {
	mu.RLock()
	s := current
	mu.RUnlock()
	if s == nil {
		return ErrNoProvider
	}
	return s.Send(ctx, to, body)
}

// Reset clears the provider. Tests only.
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	current = nil
}
`
}

func smsTestGo() string {
	return `package sms

import (
	"context"
	"errors"
	"testing"
)

func TestUnconfiguredReportsItselfAndRefuses(t *testing.T) {
	Reset()
	if Configured() {
		t.Fatal("nothing is registered, so Configured must be false")
	}
	// The backstop: a caller that skipped the Configured check gets a named
	// error rather than a nil dereference.
	if err := Send(context.Background(), "+256700000000", "hi"); !errors.Is(err, ErrNoProvider) {
		t.Errorf("expected ErrNoProvider, got %v", err)
	}
}

func TestRegisteredProviderReceivesTheMessage(t *testing.T) {
	Reset()
	defer Reset()

	var gotTo, gotBody string
	Register(SenderFunc(func(_ context.Context, to, body string) error {
		gotTo, gotBody = to, body
		return nil
	}))

	if !Configured() {
		t.Fatal("Configured must be true once a provider is registered")
	}
	if err := Send(context.Background(), "+256700000000", "code 123456"); err != nil {
		t.Fatal(err)
	}
	if gotTo != "+256700000000" || gotBody != "code 123456" {
		t.Errorf("provider got %q / %q", gotTo, gotBody)
	}
}

func TestProviderErrorIsReturned(t *testing.T) {
	Reset()
	defer Reset()

	boom := errors.New("provider rejected the number")
	Register(SenderFunc(func(_ context.Context, _, _ string) error { return boom }))

	if err := Send(context.Background(), "x", "y"); !errors.Is(err, boom) {
		t.Errorf("the provider's error should reach the caller, got %v", err)
	}
}
`
}

func recoveryModelGo() string {
	return `package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"gorm.io/gorm"

	"{{MODULE}}/internal/ids"
)

// RecoveryContactKind distinguishes the two destinations a code can go to.
type RecoveryContactKind string

const (
	RecoveryEmail RecoveryContactKind = "email"
	RecoveryPhone RecoveryContactKind = "phone"
)

// RecoveryContactToken is a single-use, expiring proof that somebody can read
// mail or texts at a destination.
//
// Modelled on EmailVerificationToken and for the same reasons: only the hash is
// stored, so leaking this table does not let an attacker confirm a contact they
// do not control, and the destination is recorded as it was when the code was
// issued. Without that last part, changing the pending address after requesting
// a code would let you verify an address you never owned.
type RecoveryContactToken struct {
	ID     string ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	UserID string ` + "`" + `gorm:"size:36;index;not null" json:"user_id"` + "`" + `

	Kind        RecoveryContactKind ` + "`" + `gorm:"size:16;not null" json:"kind"` + "`" + `
	Destination string              ` + "`" + `gorm:"size:255;not null" json:"destination"` + "`" + `

	CodeHash string ` + "`" + `gorm:"size:64;index;not null" json:"-"` + "`" + `

	// Attempts is capped so a six-digit code cannot be brute-forced inside its
	// fifteen-minute life. A million codes and unlimited guesses is not a
	// secret, it is a formality.
	Attempts int ` + "`" + `gorm:"default:0" json:"-"` + "`" + `

	ExpiresAt time.Time  ` + "`" + `gorm:"index" json:"expires_at"` + "`" + `
	UsedAt    *time.Time ` + "`" + `gorm:"index" json:"used_at,omitempty"` + "`" + `
	CreatedAt time.Time  ` + "`" + `json:"created_at"` + "`" + `
}

func (t *RecoveryContactToken) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = ids.New()
	}
	return nil
}

// RecoveryContact is a verified way back into an account.
//
// Its own table rather than columns on User, and that is not a modelling
// preference. grit upgrade does not regenerate the User model, so shipping
// this as user columns meant an upgraded project got the handler that reads
// them and a model without them: a half-feature that does not compile. A table
// of its own arrives complete through AutoMigrate, and it leaves room for more
// than one contact per kind later without another migration.
type RecoveryContact struct {
	ID     string ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	UserID string ` + "`" + `gorm:"size:36;index:idx_recovery_user_kind,unique;not null" json:"user_id"` + "`" + `

	Kind        RecoveryContactKind ` + "`" + `gorm:"size:16;index:idx_recovery_user_kind,unique;not null" json:"kind"` + "`" + `
	Destination string              ` + "`" + `gorm:"size:255;not null" json:"destination"` + "`" + `

	// Only ever written once a code has been confirmed. An unverified row
	// would look like a way back in and would not be one.
	VerifiedAt time.Time ` + "`" + `json:"verified_at"` + "`" + `
	CreatedAt  time.Time ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt  time.Time ` + "`" + `json:"updated_at"` + "`" + `
}

func (c *RecoveryContact) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = ids.New()
	}
	return nil
}

// HashRecoveryCode returns the storage form of a recovery code.
func HashRecoveryCode(code string) string {
	sum := sha256.Sum256([]byte(code))
	return hex.EncodeToString(sum[:])
}
`
}
