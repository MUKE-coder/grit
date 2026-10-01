package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeTOTPFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "totp", "totp.go"):         totpServiceGo(),
		filepath.Join(apiRoot, "internal", "models", "two_factor.go"): twoFactorModelsGo(),
		filepath.Join(apiRoot, "internal", "handlers", "totp.go"):     totpHandlerGo(),
		// The queries that handler used to run itself, since v3.342.0.
		filepath.Join(apiRoot, "internal", "services", "two_factor.go"):      apiTwoFactorServiceGo(),
		filepath.Join(apiRoot, "internal", "services", "two_factor_test.go"): apiTwoFactorServiceTestGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

// totpServiceGo generates the zero-dependency TOTP service (RFC 6238).
// Includes: secret generation, code validation, backup codes, trusted devices.
func totpServiceGo() string {
	return `package totp

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1" // #nosec G505 -- RFC 6238 TOTP is HMAC-SHA1, and every authenticator app requires it
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base32"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"net/url"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	// SecretSize is the number of random bytes for the TOTP secret.
	SecretSize = 20
	// CodeDigits is the number of digits in a TOTP code.
	CodeDigits = 6
	// Period is the time step in seconds.
	Period = 30
	// Window is the number of periods to check before/after current (clock skew tolerance).
	Window = 1
	// SkewSearch is how far out a code is searched when enrolling, and how far
	// a refused code is searched to explain the refusal. Five minutes either
	// way covers the drift a phone or a virtual machine accumulates; beyond
	// that the code is simply wrong.
	//
	// This is NOT the window a sign-in accepts. Sign-in stays at Window, one
	// step either side of the device's own recorded offset.
	SkewSearch = 10
	// BackupCodeLength is the character length of each backup code.
	BackupCodeLength = 8
	// BackupCodeCount is the default number of backup codes generated.
	BackupCodeCount = 10
	// PendingTokenExpiry is how long a TOTP pending token is valid.
	PendingTokenExpiry = 5 * time.Minute
	// MaxPendingAttempts is how many wrong codes one pending token takes before
	// it is spent and the sign-in has to start again from the password.
	MaxPendingAttempts = 5
	// TrustedDeviceDuration is how long a trusted device cookie lasts.
	TrustedDeviceDuration = 30 * 24 * time.Hour // 30 days
)

// MaxFailedAttempts is how many wrong codes, across sign-ins, lock the account
// for LockoutDuration. Wrong codes count on the same counter as wrong passwords.
var (
	MaxFailedAttempts = 10
	LockoutDuration   = 15 * time.Minute
)

// GenerateSecret creates a new random TOTP secret, base32-encoded.
func GenerateSecret() (string, error) {
	secret := make([]byte, SecretSize)
	if _, err := rand.Read(secret); err != nil {
		return "", fmt.Errorf("generating TOTP secret: %w", err)
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret), nil
}

// GenerateURI builds an otpauth:// URI for QR code generation.
func GenerateURI(secret, email, issuer string) string {
	v := url.Values{}
	v.Set("secret", secret)
	v.Set("issuer", issuer)
	v.Set("algorithm", "SHA1")
	v.Set("digits", fmt.Sprintf("%d", CodeDigits))
	v.Set("period", fmt.Sprintf("%d", Period))

	label := url.PathEscape(fmt.Sprintf("%s:%s", issuer, email))
	return fmt.Sprintf("otpauth://totp/%s?%s", label, v.Encode())
}

// EnrolStep finds the step a code belongs to, searching SkewSearch either way.
//
// Only for enrolment, where the person is already signed in, holds the secret
// they were just shown, and is proving they can read their own authenticator.
// The offset it returns is what makes the tight sign-in window work afterwards
// for a device whose clock nobody can fix.
func EnrolStep(secret, code string) (int64, bool, error) {
	counter := time.Now().Unix() / Period
	matched, found := int64(0), false
	for i := -int64(SkewSearch); i <= int64(SkewSearch); i++ {
		expected, err := generateCode(secret, counter+i)
		if err != nil {
			return 0, false, err
		}
		if hmac.Equal([]byte(expected), []byte(code)) {
			matched, found = counter+i, true
		}
	}
	return matched, found, nil
}

// SkewSeconds reports how far a code's time step is from this server's.
//
// Only called when a code has already been refused. A code from an
// authenticator whose clock has drifted is indistinguishable from a wrong code
// unless somebody measures it, and "Invalid verification code" sends people to
// re-scan a QR that was never the problem. Drift is the single most common
// reason a correct app is rejected, and it is the one thing the server can
// work out on the user's behalf.
//
// Returns the offset in seconds and whether the code matched at all within
// SkewSearch steps. A match here is NOT an acceptance: the caller still
// refuses the code. It only changes what the refusal says.
func SkewSeconds(secret, code string) (int, bool) {
	counter := time.Now().Unix() / Period
	for i := -int64(SkewSearch); i <= int64(SkewSearch); i++ {
		if i >= -int64(Window) && i <= int64(Window) {
			continue // already tried and refused by ValidateCodeStep
		}
		expected, err := generateCode(secret, counter+i)
		if err != nil {
			return 0, false
		}
		if hmac.Equal([]byte(expected), []byte(code)) {
			return int(i * Period), true
		}
	}
	return 0, false
}

// ValidateCode checks if the given TOTP code is valid for the secret.
// Accepts codes within ±Window periods for clock skew tolerance.
func ValidateCode(secret, code string) (bool, error) {
	_, ok, err := ValidateCodeStep(secret, code)
	return ok, err
}

// ValidateCodeStep is ValidateCode that also returns the time step the code
// belongs to. A caller records it and refuses any step not later than the last
// one used, because a code is otherwise good for its whole window, and anyone
// who saw it could use it again.
func ValidateCodeStep(secret, code string) (int64, bool, error) {
	return ValidateCodeOffset(secret, code, 0)
}

// ValidateCodeOffset is ValidateCodeStep for a device whose clock is known to
// be offset steps away from this server's.
//
// The window is unchanged: one step either side, of the device's time rather
// than the server's. A device that was 60 seconds behind at enrolment is not
// given a wider window, it is given the right one, which is the difference
// between accommodating a wrong clock and accepting older codes from everybody.
func ValidateCodeOffset(secret, code string, offset int64) (int64, bool, error) {
	counter := time.Now().Unix()/Period + offset
	matched, found := int64(0), false
	for i := -int64(Window); i <= int64(Window); i++ {
		expected, err := generateCode(secret, counter+i)
		if err != nil {
			return 0, false, err
		}
		// Every candidate is compared, in constant time, so the response time
		// says nothing about which one matched.
		if hmac.Equal([]byte(expected), []byte(code)) {
			matched, found = counter+i, true
		}
	}
	return matched, found, nil
}

// generateCode computes the HOTP code for a given counter (RFC 4226).
func generateCode(secret string, counter int64) (string, error) {
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", fmt.Errorf("decoding secret: %w", err)
	}

	// Counter to big-endian 8 bytes
	buf := make([]byte, 8)
	// #nosec G115 -- counter is a time step, always positive.
	binary.BigEndian.PutUint64(buf, uint64(counter))

	// HMAC-SHA1
	mac := hmac.New(sha1.New, key)
	mac.Write(buf)
	hash := mac.Sum(nil)

	// Dynamic truncation (RFC 4226 section 5.4)
	offset := hash[len(hash)-1] & 0x0f
	truncated := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7fffffff

	// Modulo 10^digits
	code := truncated % uint32(math.Pow10(CodeDigits))
	return fmt.Sprintf("%0*d", CodeDigits, code), nil
}

// GenerateBackupCodes creates a set of one-time-use recovery codes.
// Returns the plaintext codes (show once to user) and their bcrypt hashes (store in DB).
func GenerateBackupCodes(count int) ([]string, []string, error) {
	if count == 0 {
		count = BackupCodeCount
	}

	codes := make([]string, count)
	hashes := make([]string, count)

	for i := 0; i < count; i++ {
		raw := make([]byte, BackupCodeLength)
		if _, err := rand.Read(raw); err != nil {
			return nil, nil, fmt.Errorf("generating backup code: %w", err)
		}
		// Hex-encode and take first BackupCodeLength characters, uppercase for readability
		code := strings.ToUpper(hex.EncodeToString(raw)[:BackupCodeLength])
		codes[i] = code

		hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
		if err != nil {
			return nil, nil, fmt.Errorf("hashing backup code: %w", err)
		}
		hashes[i] = string(hash)
	}

	return codes, hashes, nil
}

// VerifyBackupCode checks if a code matches any of the stored hashes.
// Returns the index of the matched code, or -1 if no match.
func VerifyBackupCode(code string, hashes []string) int {
	code = strings.ToUpper(strings.TrimSpace(code))
	for i, h := range hashes {
		if bcrypt.CompareHashAndPassword([]byte(h), []byte(code)) == nil {
			return i
		}
	}
	return -1
}

// GeneratePendingToken creates a random token for the TOTP verification step.
func GeneratePendingToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// GenerateEmailCode returns the six digits that go in the email.
//
// crypto/rand, not math/rand: this is the whole second factor, and a code an
// attacker can predict from the clock is not one. Six digits with five attempts
// and a five-minute life is a one-in-two-hundred-thousand guess.
func GenerateEmailCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// EmailCodeMatches compares a typed code with the stored hash, in constant
// time, so the comparison cannot be timed a digit at a time.
func EmailCodeMatches(hash, code string) bool {
	if hash == "" {
		return false
	}
	want, err := hex.DecodeString(hash)
	if err != nil {
		return false
	}
	got := sha256.Sum256([]byte(strings.TrimSpace(code)))
	return subtle.ConstantTimeCompare(want, got[:]) == 1
}

// HashToken returns the SHA-256 hash of a token (for DB storage).
func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// GenerateDeviceToken creates a random token for trusted device cookies.
func GenerateDeviceToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
`
}

// twoFactorModelsGo generates the TwoFactorConfig and TrustedDevice models.
func twoFactorModelsGo() string {
	return `package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"{{MODULE}}/internal/crypto"
)

// TwoFactorConfig stores TOTP settings for a user.
// How a second factor is delivered.
const (
	// TwoFactorMethodApp is an authenticator app: the stronger of the two,
	// because the code never leaves the device.
	TwoFactorMethodApp = "app"
	// TwoFactorMethodEmail sends a code to the address on the account. Only as
	// safe as that mailbox, which is also where a password reset goes, so it
	// adds less than an authenticator does. It is here because the alternative
	// for somebody who will not install an app is no second factor at all.
	TwoFactorMethodEmail = "email"
)

type TwoFactorConfig struct {
	ID     uint   ` + "`" + `gorm:"primarykey" json:"id"` + "`" + `
	UserID string ` + "`" + `gorm:"size:36;uniqueIndex;not null" json:"user_id"` + "`" + `
` + twoFactorSecretField + `	Enabled     bool                       ` + "`" + `gorm:"default:false" json:"enabled"` + "`" + `
	// Method is how the second factor is delivered: "app" for an authenticator
	// and "email" for a code sent to the address on the account.
	//
	// An authenticator is the stronger of the two, because an emailed code is
	// only as safe as the mailbox, and for most people the mailbox is what the
	// password reset goes to anyway. Email is here because the alternative for
	// somebody who will not install an app is no second factor at all, and that
	// is worse than a mailbox.
	//
	// Empty means "app": rows written before this column existed were all
	// authenticators, and defaulting it keeps them working without a backfill.
	Method      string                     ` + "`" + `gorm:"size:16;not null;default:app" json:"method"` + "`" + `
	BackupCodes datatypes.JSONSlice[string] ` + "`" + `gorm:"type:text" json:"-"` + "`" + `
	// LastUsedStep is the time step of the last code accepted. Only a later one
	// is accepted next, so a code cannot be used twice.
	LastUsedStep int64 ` + "`" + `gorm:"not null;default:0" json:"-"` + "`" + `
	// StepOffset is how many 30 second steps this device's clock is from the
	// server's, measured when the authenticator was enrolled and re-measured on
	// every accepted code.
	//
	// Without it, a laptop a minute behind produces correct codes that are
	// always one step too old, and the only advice the server can give is "fix
	// your clock", which on a managed machine is not advice. It does not widen
	// what is accepted: the window stays one step, around this device's time
	// rather than the server's. Zero for every device whose clock is right,
	// which is almost all of them.
	StepOffset int64 ` + "`" + `gorm:"not null;default:0" json:"-"` + "`" + `
	// SetupCodeHash and SetupCodeExpiresAt hold the code that proves somebody
	// can read the mailbox before email becomes their second factor. Turning on
	// a factor you cannot receive is how an account locks itself out, so it is
	// deliberately two steps.
	SetupCodeHash      string     ` + "`" + `gorm:"size:64" json:"-"` + "`" + `
	SetupCodeExpiresAt *time.Time ` + "`" + `json:"-"` + "`" + `
	CreatedAt   time.Time                  ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt   time.Time                  ` + "`" + `json:"updated_at"` + "`" + `
}

// TrustedDevice stores a remembered device that can skip TOTP.
type TrustedDevice struct {
	ID        uint           ` + "`" + `gorm:"primarykey" json:"id"` + "`" + `
	UserID    string         ` + "`" + `gorm:"size:36;index;not null" json:"user_id"` + "`" + `
	TokenHash string         ` + "`" + `gorm:"size:64;uniqueIndex;not null" json:"-"` + "`" + `
	UserAgent string         ` + "`" + `gorm:"size:500" json:"user_agent"` + "`" + `
	IPAddress string         ` + "`" + `gorm:"size:45" json:"ip_address"` + "`" + `
	ExpiresAt time.Time      ` + "`" + `gorm:"not null;index" json:"expires_at"` + "`" + `
	CreatedAt time.Time      ` + "`" + `json:"created_at"` + "`" + `
	DeletedAt gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}

// TOTPPendingToken stores short-lived tokens for the TOTP verification step.
type TOTPPendingToken struct {
	ID        uint      ` + "`" + `gorm:"primarykey" json:"id"` + "`" + `
	UserID    string    ` + "`" + `gorm:"size:36;index;not null" json:"user_id"` + "`" + `
	TokenHash string    ` + "`" + `gorm:"size:64;uniqueIndex;not null" json:"-"` + "`" + `
	ExpiresAt time.Time ` + "`" + `gorm:"not null" json:"expires_at"` + "`" + `
	// Attempts counts wrong codes against this token; at totp.MaxPendingAttempts
	// it is refused.
	Attempts int ` + "`" + `gorm:"not null;default:0" json:"-"` + "`" + `
	// CodeHash is set when the second factor is a code sent by email, and holds
	// the SHA-256 of that code. Hashed for the same reason the token is: a
	// six-digit code sitting in a table is a six-digit code an attacker with a
	// read of that table can type.
	//
	// Empty means the challenge is an authenticator, and the code is computed
	// rather than stored.
	CodeHash  string    ` + "`" + `gorm:"size:64" json:"-"` + "`" + `
	CreatedAt time.Time ` + "`" + `json:"created_at"` + "`" + `
}

// CleanupExpiredTOTPTokens removes expired pending tokens and trusted devices.
func CleanupExpiredTOTPTokens(db *gorm.DB) error {
	now := time.Now()
	if err := db.Where("expires_at < ?", now).Delete(&TOTPPendingToken{}).Error; err != nil {
		return err
	}
	return db.Where("expires_at < ?", now).Delete(&TrustedDevice{}).Error
}
`
}

// totpHandlerGo generates the HTTP handlers for TOTP setup, verification, and management.
func totpHandlerGo() string { return tmpl("api/handlers/totp.go") }
