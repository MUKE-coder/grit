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
func totpHandlerGo() string {
	return `package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"{{MODULE}}/internal/config"
	"{{MODULE}}/internal/crypto"
	"{{MODULE}}/internal/jobs"
	"{{MODULE}}/internal/mail"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
	"{{MODULE}}/internal/totp"
	"{{MODULE}}/internal/respond"
)

// TOTPHandler handles two-factor authentication endpoints.
type TOTPHandler struct {
	DB          *gorm.DB
	AuthService *services.AuthService
	Issuer      string // App name for authenticator display

	// Set when the second factor can be a code by email.
	Config *config.Config
	Mailer *mail.Mailer
	Jobs   *jobs.Client
}

type totpSetupResponse struct {
	Secret string ` + "`" + `json:"secret"` + "`" + `
	URI    string ` + "`" + `json:"uri"` + "`" + `
	// QRCode is a data: URI holding a PNG of URI, ready to drop straight into
	// an <img src>. Rendered here rather than in each client because there are
	// four of them (admin, web, desktop, Expo) and none should have to ship a
	// QR encoder — or handle the raw secret — to show a setup screen.
	QRCode string ` + "`" + `json:"qr_code"` + "`" + `
}

type EnableTOTPRequest struct {
	Secret string ` + "`" + `json:"secret" binding:"required"` + "`" + `
	Code   string ` + "`" + `json:"code" binding:"required"` + "`" + `
}

type VerifyTOTPRequest struct {
	PendingToken string ` + "`" + `json:"pending_token" binding:"required"` + "`" + `
	Code         string ` + "`" + `json:"code" binding:"required"` + "`" + `
	TrustDevice  bool   ` + "`" + `json:"trust_device"` + "`" + `
}

type DisableTOTPRequest struct {
	Password string ` + "`" + `json:"password" binding:"required"` + "`" + `
}

type VerifyBackupCodeRequest struct {
	PendingToken string ` + "`" + `json:"pending_token" binding:"required"` + "`" + `
	Code         string ` + "`" + `json:"code" binding:"required"` + "`" + `
	TrustDevice  bool   ` + "`" + `json:"trust_device"` + "`" + `
}

// Setup generates a new TOTP secret and QR code URI for the user.
// The secret is NOT stored yet — the user must verify it first via Enable.
func (h *TOTPHandler) Setup(c *gin.Context) {
	userID := c.GetString("user_id")
	user, _ := c.Get("user")
	// Set by the auth middleware. Checked rather than asserted, so the route
	// mounted without it answers 401 instead of panicking.
	u, ok := user.(models.User)
	if !ok {
		respond.Fail(c, respond.CodeUnauthorized, "Not signed in")
		return
	}

	// Check if TOTP is already enabled
	var existing models.TwoFactorConfig
	if err := h.DB.WithContext(c.Request.Context()).Where("user_id = ?", userID).First(&existing).Error; err == nil && existing.Enabled {
		respond.Fail(c, respond.CodeTOTPAlreadyEnabled, "Two-factor authentication is already enabled")
		return
	}

	secret, err := totp.GenerateSecret()
	if err != nil {
		respond.Fail(c, respond.CodeTOTPError, "Failed to generate secret")
		return
	}

	uri := totp.GenerateURI(secret, u.Email, h.Issuer)

	// Medium recovery level: a phone camera reads it reliably at the size a
	// setup dialog shows it, without the density High produces.
	qrPNG, err := qrCodePNG(uri, 256)
	if err != nil {
		respond.Fail(c, respond.CodeTOTPError, "Failed to render the QR code")
		return
	}
	qrDataURI := "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrPNG)

	c.JSON(http.StatusOK, gin.H{
		"data": totpSetupResponse{
			Secret: secret,
			URI:    uri,
			QRCode: qrDataURI,
		},
		"message": "Scan the QR code with your authenticator app, then verify with a code",
	})
}

// totpRefusal explains a refused code.
//
// The code is refused either way. This only decides what to tell the person,
// and the difference matters: "your clock is 90 seconds behind" is something
// they can fix in a minute, while "invalid code" sends them to re-scan a QR
// that was never wrong.
//
// It says nothing an attacker can use. Learning that a guess landed near the
// window requires guessing a valid six-digit code for a nearby step, which is
// the same work as guessing the current one, and the attempt is counted the
// same way.
func totpRefusal(secret, code string) string {
	if off, found := totp.SkewSeconds(secret, code); found {
		direction := "behind"
		if off > 0 {
			direction = "ahead of"
		}
		if off < 0 {
			off = -off
		}
		return fmt.Sprintf(
			"That code was right, but your device's clock is about %d seconds %s this server's, so the code had already changed. Turn on automatic time on your phone and try the next code.",
			off, direction)
	}
	return "Invalid verification code. Check you are reading the code for this account, and that it has not just expired."
}

// Enable verifies the initial TOTP code and activates 2FA for the user.
// Returns backup codes that the user should save.
func (h *TOTPHandler) Enable(c *gin.Context) {
	userID := c.GetString("user_id")

	var req EnableTOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	// Enabling replaces the secret, so it is refused while 2FA is on: a stolen
	// session could otherwise re-enrol the account to an authenticator the
	// attacker holds. Disabling first asks for the password.
	var existing models.TwoFactorConfig
	if err := h.DB.WithContext(c.Request.Context()).Where("user_id = ?", userID).First(&existing).Error; err == nil && existing.Enabled {
		respond.Fail(c, respond.CodeTOTPAlreadyEnabled, "Two-factor authentication is already enabled. Disable it first.")
		return
	}

	// Enrolment searches wider than a sign-in does, and keeps the offset it
	// finds. Refusing a correct code because the device's clock is a minute out
	// leaves somebody who cannot change that clock with no second factor at all.
	step, valid, err := totp.EnrolStep(req.Secret, req.Code)
	if err != nil || !valid {
		respond.Fail(c, respond.CodeInvalidTOTPCode, totpRefusal(req.Secret, req.Code))
		return
	}
	offset := step - time.Now().Unix()/totp.Period

	// Generate backup codes
	codes, hashes, err := totp.GenerateBackupCodes(0)
	if err != nil {
		respond.Fail(c, respond.CodeTOTPError, "Failed to generate backup codes")
		return
	}

	// Upsert the TwoFactorConfig
	var config models.TwoFactorConfig
	if err := h.DB.WithContext(c.Request.Context()).Where("user_id = ?", userID).FirstOrCreate(&config, models.TwoFactorConfig{UserID: userID}).Error; err != nil {
		respond.Fail(c, respond.CodeTOTPError, "Failed to enable two-factor authentication")
		return
	}

	// Encrypted at rest when FIELD_ENCRYPTION_KEY is set: see models.TwoFactorConfig.
	config.Secret = crypto.EncryptedString(req.Secret)
	config.Enabled = true
	config.BackupCodes = hashes
	// The code that enabled 2FA is spent; it cannot also sign in.
	config.LastUsedStep = step
	config.StepOffset = offset

	if err := h.DB.WithContext(c.Request.Context()).Save(&config).Error; err != nil {
		respond.Fail(c, respond.CodeTOTPError, "Failed to enable two-factor authentication")
		return
	}

	message := "Two-factor authentication enabled. Save your backup codes in a safe place."
	if offset != 0 {
		// Said out loud, because the account now depends on a clock the person
		// may not know is wrong, and because a clock that has drifted once
		// usually keeps drifting.
		message += fmt.Sprintf(
			" Your device's clock is about %d seconds %s this server's; that has been allowed for, but it is worth fixing.",
			abs64(offset*totp.Period), aheadOrBehind(offset))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"enabled":      true,
			"backup_codes": codes,
			"clock_offset": offset * totp.Period,
		},
		"message": message,
	})
}

// abs64 is the absolute value, for turning an offset into a distance.
func abs64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

// aheadOrBehind names the direction of a clock offset.
func aheadOrBehind(offset int64) string {
	if offset > 0 {
		return "ahead of"
	}
	return "behind"
}

// SendEmailSetupCode emails a code to prove the address works, before email
// becomes this account's second factor.
//
//	POST /api/v1/auth/totp/email/send
func (h *TOTPHandler) SendEmailSetupCode(c *gin.Context) {
	userID := c.GetString("user_id")

	var user models.User
	if err := h.DB.WithContext(c.Request.Context()).Where("id = ?", userID).First(&user).Error; err != nil {
		respond.Fail(c, respond.CodeNotFound, "User not found")
		return
	}

	// Nothing to send with means nothing to sign in with later. Refused here
	// rather than discovered at the next sign-in, when the person would be
	// holding a factor they can never satisfy.
	if h.Mailer == nil && h.Jobs == nil {
		respond.Fail(c, respond.CodeMailFailed, "This deployment cannot send email yet, so codes by email cannot be turned on. Set MAIL_MAILER first.")
		return
	}

	code, err := totp.GenerateEmailCode()
	if err != nil {
		respond.Fail(c, respond.CodeTOTPError, "Failed to create a code")
		return
	}
	expires := time.Now().Add(totp.PendingTokenExpiry)

	var config models.TwoFactorConfig
	err = h.DB.WithContext(c.Request.Context()).Where("user_id = ?", userID).First(&config).Error
	if err != nil {
		config = models.TwoFactorConfig{UserID: userID}
	}
	config.Method = models.TwoFactorMethodEmail
	config.SetupCodeHash = totp.HashToken(code)
	config.SetupCodeExpiresAt = &expires
	// Enabled is untouched: an account already using an authenticator keeps it
	// working until the new method is confirmed.
	if err := h.DB.WithContext(c.Request.Context()).Save(&config).Error; err != nil {
		respond.Fail(c, respond.CodeTOTPError, "Failed to start email two-factor setup")
		return
	}

	if err := dispatchMail(c.Request.Context(), h.Mailer, h.Jobs, "two-factor-setup:"+userID+":"+config.SetupCodeHash, mail.SendOptions{
		To:       user.Email,
		Subject:  "Your sign-in code",
		Template: "two-factor-code",
		Data: map[string]interface{}{
			"AppName": h.Config.AppName,
			"Title":   "Your sign-in code",
			"Code":    code,
			"Minutes": int(totp.PendingTokenExpiry.Minutes()),
			"Year":    time.Now().Year(),
		},
	}); err != nil {
		log.Printf("two-factor setup: emailing a code to %s: %v", userID, err)
		respond.Fail(c, respond.CodeMailFailed, "We could not email your code. Check the mail settings and try again.")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    gin.H{"sent_to": user.Email},
		"message": "We sent a code to " + user.Email + ".",
	})
}

// EnableEmail turns on the email second factor once its code comes back.
//
//	POST /api/v1/auth/totp/email/enable
func (h *TOTPHandler) EnableEmail(c *gin.Context) {
	var req struct {
		Code string ` + "`" + `json:"code" binding:"required"` + "`" + `
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	userID := c.GetString("user_id")
	var config models.TwoFactorConfig
	if err := h.DB.WithContext(c.Request.Context()).Where("user_id = ?", userID).First(&config).Error; err != nil {
		respond.Fail(c, respond.CodeNotFound, "Ask for a code first.")
		return
	}
	if config.SetupCodeHash == "" || config.SetupCodeExpiresAt == nil || time.Now().After(*config.SetupCodeExpiresAt) {
		respond.Fail(c, respond.CodeInvalidCode, "That code has expired. Ask for a new one.")
		return
	}
	if !totp.EmailCodeMatches(config.SetupCodeHash, req.Code) {
		respond.Fail(c, respond.CodeInvalidCode, "That code is wrong. Check the email and try again.")
		return
	}

	// Spent on use, and the method is only committed here: until this point an
	// account with an authenticator still has one.
	codes, hashes, err := totp.GenerateBackupCodes(0)
	if err != nil {
		respond.Fail(c, respond.CodeTOTPError, "Failed to create backup codes")
		return
	}
	config.Enabled = true
	config.Method = models.TwoFactorMethodEmail
	config.SetupCodeHash = ""
	config.SetupCodeExpiresAt = nil
	config.BackupCodes = hashes
	if err := h.DB.WithContext(c.Request.Context()).Save(&config).Error; err != nil {
		respond.Fail(c, respond.CodeTOTPError, "Failed to enable two-factor authentication")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"enabled":      true,
			"method":       models.TwoFactorMethodEmail,
			"backup_codes": codes,
		},
		"message": "Codes by email are on. Save your backup codes: they are how you get in if you lose the mailbox.",
	})
}

// Verify validates a TOTP code during the login flow (after password check).
// Exchanges a pending token + valid TOTP code for real JWT tokens.
func (h *TOTPHandler) Verify(c *gin.Context) {
	var req VerifyTOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	pending, user, config, ok := h.beginSecondFactor(c, req.PendingToken)
	if !ok {
		return
	}

	// A challenge sent by email carries its own code. Nothing else about the
	// flow changes: same pending token, same attempt counting, same trusted
	// device, same session at the end.
	if pending.CodeHash != "" {
		if !totp.EmailCodeMatches(pending.CodeHash, req.Code) {
			h.failSecondFactor(pending, user)
			respond.Fail(c, respond.CodeInvalidCode, "That code is wrong or has expired. Check the email, or sign in again for a new code.")
			return
		}
		// Spent: the row goes when the session is issued, and clearing the hash
		// first means a replay in the same instant finds nothing to match.
		if err := h.DB.WithContext(c.Request.Context()).Model(&models.TOTPPendingToken{}).
			Where("id = ? AND code_hash = ?", pending.ID, pending.CodeHash).
			Update("code_hash", "").Error; err != nil {
			respond.Fail(c, respond.CodeTOTPError, "Failed to complete sign-in")
			return
		}
		h.completeSecondFactor(c, pending, user, req.TrustDevice, nil, "Signed in")
		return
	}

	step, valid, err := totp.ValidateCodeOffset(string(config.Secret), req.Code, config.StepOffset)
	if err == nil && valid {
		// A code is good for its window, so without this the same code signs in
		// again for as long as it lasts. The step only moves forward, and the
		// update is conditional, so two requests cannot both spend one code.
		//
		// The offset is re-recorded from the step that actually matched, which
		// is RFC 6238 resynchronisation: a clock that keeps losing a second a
		// day is followed instead of eventually locking the account out.
		res := h.DB.WithContext(c.Request.Context()).Model(&models.TwoFactorConfig{}).
			Where("id = ? AND last_used_step < ?", config.ID, step).
			UpdateColumns(map[string]interface{}{
				"last_used_step": step,
				"step_offset":    step - time.Now().Unix()/totp.Period,
			})
		valid = res.Error == nil && res.RowsAffected == 1
	}
	if err != nil || !valid {
		h.failSecondFactor(pending, user)
		// Same diagnosis as enrolment: a clock that jumped since the device was
		// enrolled is the usual reason a correct-looking code stops working,
		// and "invalid code" sends people to their backup codes for nothing.
		respond.Fail(c, respond.CodeInvalidTOTPCode, totpRefusal(string(config.Secret), req.Code))
		return
	}

	h.sealSecret(c, config)
	h.completeSecondFactor(c, pending, user, req.TrustDevice, nil, "Logged in successfully")
}

` + totpSealSecretFunc + `
// VerifyBackupCode validates a backup code during login (alternative to TOTP).
func (h *TOTPHandler) VerifyBackupCode(c *gin.Context) {
	var req VerifyBackupCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	pending, user, config, ok := h.beginSecondFactor(c, req.PendingToken)
	if !ok {
		return
	}

	idx := totp.VerifyBackupCode(req.Code, config.BackupCodes)
	if idx < 0 {
		h.failSecondFactor(pending, user)
		respond.Fail(c, respond.CodeInvalidBackupCode, "Invalid backup code")
		return
	}

	// Spend the backup code before signing in, and fail if that cannot be saved:
	// a code that stays in the list stays usable.
	remaining := make([]string, 0, len(config.BackupCodes)-1)
	remaining = append(remaining, config.BackupCodes[:idx]...)
	remaining = append(remaining, config.BackupCodes[idx+1:]...)
	// The update matches only while the list is still the one this request
	// read. Two requests carrying the same code both find it in the list; only
	// the first to write can spend it, and the other is refused.
	read, err := json.Marshal([]string(config.BackupCodes))
	if err != nil {
		respond.Fail(c, respond.CodeTOTPError, "Failed to use the backup code")
		return
	}
	res := h.DB.WithContext(c.Request.Context()).Model(&models.TwoFactorConfig{}).
		Where("id = ? AND backup_codes = ?", config.ID, string(read)).
		Update("backup_codes", datatypes.JSONSlice[string](remaining))
	if res.Error != nil {
		respond.Fail(c, respond.CodeTOTPError, "Failed to use the backup code")
		return
	}
	if res.RowsAffected != 1 {
		h.failSecondFactor(pending, user)
		respond.Fail(c, respond.CodeInvalidBackupCode, "Invalid backup code")
		return
	}
	config.BackupCodes = remaining
	h.sealSecret(c, config)

	h.completeSecondFactor(c, pending, user, req.TrustDevice, gin.H{"backup_codes_remaining": len(remaining)},
		"Logged in successfully with backup code")
}

// beginSecondFactor loads what a second-factor attempt needs, and refuses one
// that must not go ahead: a pending token that is unknown, expired or out of
// attempts, or an account that is disabled or locked. Password sign-in makes
// the same account checks; before, the second step skipped them.
func (h *TOTPHandler) beginSecondFactor(c *gin.Context, pendingToken string) (*models.TOTPPendingToken, *models.User, *models.TwoFactorConfig, bool) {
	invalid := func() {
		respond.Fail(c, respond.CodeInvalidPendingToken, "Invalid or expired verification session. Please log in again.")
	}

	var pending models.TOTPPendingToken
	if err := h.DB.WithContext(c.Request.Context()).Where("token_hash = ? AND expires_at > ?", totp.HashToken(pendingToken), time.Now()).First(&pending).Error; err != nil {
		invalid()
		return nil, nil, nil, false
	}
	if pending.Attempts >= totp.MaxPendingAttempts {
		invalid()
		return nil, nil, nil, false
	}

	var user models.User
	if err := h.DB.WithContext(c.Request.Context()).Where("id = ?", pending.UserID).First(&user).Error; err != nil {
		invalid()
		return nil, nil, nil, false
	}
	if !user.Active {
		respond.Fail(c, respond.CodeAccountDisabled, "Your account has been disabled")
		return nil, nil, nil, false
	}
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		respond.Fail(c, respond.CodeAccountLocked, "Too many failed attempts. Try again later, or reset your password.")
		return nil, nil, nil, false
	}

	var config models.TwoFactorConfig
	if err := h.DB.WithContext(c.Request.Context()).Where("user_id = ? AND enabled = ?", pending.UserID, true).First(&config).Error; err != nil {
		respond.Fail(c, respond.CodeTOTPError, "Two-factor configuration not found")
		return nil, nil, nil, false
	}
	return &pending, &user, &config, true
}

// failSecondFactor counts a wrong code against the pending token and against
// the account. Each count is its own UPDATE, so concurrent guesses cannot share
// one increment. Before, a pending token took unlimited guesses, and signing in
// again with the password cleared the account's count.
func (h *TOTPHandler) failSecondFactor(pending *models.TOTPPendingToken, user *models.User) {
	if err := h.DB.Model(&models.TOTPPendingToken{}).Where("id = ?", pending.ID).
		UpdateColumn("attempts", gorm.Expr("attempts + 1")).Error; err != nil {
		log.Printf("totp: counting an attempt for %s: %v", user.ID, err)
	}
	if err := h.DB.Model(&models.User{}).Where("id = ?", user.ID).
		UpdateColumn("failed_login_count", gorm.Expr("failed_login_count + 1")).Error; err != nil {
		log.Printf("totp: counting a failure for %s: %v", user.ID, err)
		return
	}

	var fresh models.User
	if err := h.DB.Select("id", "failed_login_count").First(&fresh, "id = ?", user.ID).Error; err != nil {
		return
	}
	if totp.MaxFailedAttempts <= 0 || fresh.FailedLoginCount < totp.MaxFailedAttempts {
		return
	}
	until := time.Now().Add(totp.LockoutDuration)
	if err := h.DB.Model(&models.User{}).Where("id = ?", user.ID).
		Updates(map[string]interface{}{"locked_until": until, "failed_login_count": 0}).Error; err != nil {
		log.Printf("totp: locking %s: %v", user.ID, err)
		return
	}
	if err := h.DB.Where("user_id = ?", user.ID).Delete(&models.TOTPPendingToken{}).Error; err != nil {
		log.Printf("totp: clearing pending tokens for %s: %v", user.ID, err)
	}
	log.Printf("lockout: %s locked until %s after %d wrong 2FA codes", user.Email, until.Format(time.RFC3339), totp.MaxFailedAttempts)
}

// completeSecondFactor spends the pending token, clears the failure count and
// signs the user in.
func (h *TOTPHandler) completeSecondFactor(c *gin.Context, pending *models.TOTPPendingToken, user *models.User, trustDevice bool, extra gin.H, message string) {
	// Only the request that deletes the pending token may use it, so two
	// requests carrying one token cannot both sign in.
	res := h.DB.WithContext(c.Request.Context()).Delete(&models.TOTPPendingToken{}, pending.ID)
	if res.Error != nil || res.RowsAffected != 1 {
		respond.Fail(c, respond.CodeInvalidPendingToken, "Invalid or expired verification session. Please log in again.")
		return
	}
	if user.FailedLoginCount > 0 || user.LockedUntil != nil {
		if err := h.DB.WithContext(c.Request.Context()).Model(&models.User{}).Where("id = ?", user.ID).
			Updates(map[string]interface{}{"failed_login_count": 0, "locked_until": nil}).Error; err != nil {
			log.Printf("totp: clearing the failure count for %s: %v", user.ID, err)
		}
	}

	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		respond.Fail(c, respond.CodeTokenError, "Failed to generate tokens")
		return
	}
	// Record the session. An access token names its session, and one whose
	// session was never recorded is refused on its first request.
	if _, err := services.CreateSession(h.DB, c, user.ID, tokens.RefreshToken); err != nil {
		log.Printf("totp: failed to record session for %s: %v", user.ID, err)
	}

	if trustDevice {
		h.createTrustedDevice(c, user.ID)
	}

	// Mirror tokens into HttpOnly cookies so the browser client doesn't
	// need to handle them in JS. Native bearer clients use the JSON body.
	h.AuthService.SetAuthCookies(c, tokens)

	data := gin.H{"user": user, "tokens": tokens}
	for k, v := range extra {
		data[k] = v
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "message": message})
}

// Disable turns off 2FA for the user (requires password confirmation).
func (h *TOTPHandler) Disable(c *gin.Context) {
	userID := c.GetString("user_id")

	var req DisableTOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	// Verify password
	var user models.User
	if err := h.DB.WithContext(c.Request.Context()).Where("id = ?", userID).First(&user).Error; err != nil {
		respond.Fail(c, respond.CodeNotFound, "User not found")
		return
	}

	if !user.CheckPassword(req.Password) {
		respond.Fail(c, respond.CodeInvalidPassword, "Incorrect password")
		return
	}

	// The config, the trusted devices and the pending tokens go together, or
	// none of them do. The three deletes ran unchecked, and the answer was
	// "disabled" whether or not anything had been deleted.
	if err := h.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		for _, table := range []interface{}{&models.TwoFactorConfig{}, &models.TrustedDevice{}, &models.TOTPPendingToken{}} {
			if err := tx.Where("user_id = ?", userID).Delete(table).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		respond.Fail(c, respond.CodeTOTPError, "Failed to disable two-factor authentication")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Two-factor authentication disabled",
	})
}

// Status returns whether TOTP is enabled for the current user.
func (h *TOTPHandler) Status(c *gin.Context) {
	userID := c.GetString("user_id")

	var config models.TwoFactorConfig
	enabled := false
	backupCodesRemaining := 0

	if err := h.DB.WithContext(c.Request.Context()).Where("user_id = ?", userID).First(&config).Error; err == nil {
		enabled = config.Enabled
		backupCodesRemaining = len(config.BackupCodes)
	}

	// Count trusted devices
	var deviceCount int64
	h.DB.WithContext(c.Request.Context()).Model(&models.TrustedDevice{}).Where("user_id = ? AND expires_at > ?", userID, time.Now()).Count(&deviceCount)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"enabled":                enabled,
			"backup_codes_remaining": backupCodesRemaining,
			"trusted_devices":        deviceCount,
		},
	})
}

// RegenerateBackupCodes creates a new set of backup codes, replacing the old ones.
func (h *TOTPHandler) RegenerateBackupCodes(c *gin.Context) {
	userID := c.GetString("user_id")

	var config models.TwoFactorConfig
	if err := h.DB.WithContext(c.Request.Context()).Where("user_id = ? AND enabled = ?", userID, true).First(&config).Error; err != nil {
		respond.Fail(c, respond.CodeTOTPNotEnabled, "Two-factor authentication is not enabled")
		return
	}

	codes, hashes, err := totp.GenerateBackupCodes(0)
	if err != nil {
		respond.Fail(c, respond.CodeTOTPError, "Failed to generate backup codes")
		return
	}

	config.BackupCodes = hashes
	if err := h.DB.WithContext(c.Request.Context()).Save(&config).Error; err != nil {
		respond.Fail(c, respond.CodeTOTPError, "Failed to save backup codes")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"backup_codes": codes,
		},
		"message": "New backup codes generated. Previous codes are now invalid.",
	})
}

// ListTrustedDevices returns the devices allowed to skip the TOTP prompt.
//
// The status endpoint only reports a count, which is enough to nag with and
// useless to act on: "3 trusted devices" tells a user nothing about whether
// one of them is a laptop they sold. Expired rows are filtered rather than
// deleted here — the nightly cleanup owns deletion, and a read should not
// mutate.
func (h *TOTPHandler) ListTrustedDevices(c *gin.Context) {
	userID := c.GetString("user_id")

	var devices []models.TrustedDevice
	if err := h.DB.WithContext(c.Request.Context()).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Order("created_at desc").
		Find(&devices).Error; err != nil {
		respond.Fail(c, respond.CodeInternalError, "Failed to load trusted devices")
		return
	}

	// The caller's own device is marked so the UI can label it rather than
	// inviting someone to revoke the browser they are sitting in.
	currentHash := ""
	if cookie, err := c.Cookie("totp_trusted"); err == nil && cookie != "" {
		currentHash = totp.HashToken(cookie)
	}

	out := make([]gin.H, 0, len(devices))
	for _, d := range devices {
		out = append(out, gin.H{
			"id":         d.ID,
			"user_agent": d.UserAgent,
			"ip_address": d.IPAddress,
			"created_at": d.CreatedAt,
			"expires_at": d.ExpiresAt,
			"current":    currentHash != "" && d.TokenHash == currentHash,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": out})
}

// RevokeTrustedDevice removes one trusted device by id.
func (h *TOTPHandler) RevokeTrustedDevice(c *gin.Context) {
	userID := c.GetString("user_id")

	// Scoped to the caller: without the user_id predicate this would let any
	// authenticated user revoke anyone's device by guessing an id.
	res := h.DB.WithContext(c.Request.Context()).Where("id = ? AND user_id = ?", c.Param("id"), userID).
		Delete(&models.TrustedDevice{})
	if res.Error != nil {
		respond.Fail(c, respond.CodeInternalError, "Failed to revoke the device")
		return
	}
	if res.RowsAffected == 0 {
		respond.Fail(c, respond.CodeNotFound, "Trusted device not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device revoked. It will need a code at the next sign-in."})
}

// RevokeTrustedDevices removes all trusted devices for the current user.
func (h *TOTPHandler) RevokeTrustedDevices(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := h.DB.WithContext(c.Request.Context()).Where("user_id = ?", userID).Delete(&models.TrustedDevice{}).Error; err != nil {
		respond.Fail(c, respond.CodeTOTPError, "Failed to revoke trusted devices")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All trusted devices revoked. You will need to enter a TOTP code on next login.",
	})
}

// createTrustedDevice stores a trusted device and sets the cookie.
func (h *TOTPHandler) createTrustedDevice(c *gin.Context, userID string) {
	deviceToken, err := totp.GenerateDeviceToken()
	if err != nil {
		return // Non-critical, skip silently
	}

	device := models.TrustedDevice{
		UserID:    userID,
		TokenHash: totp.HashToken(deviceToken),
		UserAgent: c.Request.UserAgent(),
		IPAddress: c.ClientIP(),
		ExpiresAt: time.Now().Add(totp.TrustedDeviceDuration),
	}

	if err := h.DB.WithContext(c.Request.Context()).Create(&device).Error; err != nil {
		return // Non-critical
	}

` + trustedDeviceCookieNew + `}

// IsTrustedDevice checks if the current request has a valid trusted device cookie.
func IsTrustedDevice(c *gin.Context, db *gorm.DB, userID string) bool {
	token, err := c.Cookie("totp_trusted")
	if err != nil || token == "" {
		return false
	}

	tokenHash := totp.HashToken(token)
	var device models.TrustedDevice
	if err := db.Where("user_id = ? AND token_hash = ? AND expires_at > ?", userID, tokenHash, time.Now()).First(&device).Error; err != nil {
		return false
	}

	// Refresh the device expiry (sliding window). The device is trusted either
	// way; a refresh that fails only means the window is not extended.
	device.ExpiresAt = time.Now().Add(totp.TrustedDeviceDuration)
	if err := db.Model(&device).Update("expires_at", device.ExpiresAt).Error; err != nil {
		log.Printf("totp: extending trusted device %d: %v", device.ID, err)
	}

	return true
}
` + totpQRHelper
}

// totpQRHelper draws the 2FA setup QR code with boombuler/barcode, which
// replaced skip2/go-qrcode (last released in 2020).
const totpQRHelper = `
// qrCodePNG renders content as a PNG QR code at most size pixels square, at the
// Medium recovery level, with the four-module quiet zone a phone camera needs
// to find the code on a dark page. Whole pixels per module keep the edges sharp.
func qrCodePNG(content string, size int) ([]byte, error) {
	code, err := qr.Encode(content, qr.M, qr.Auto)
	if err != nil {
		return nil, err
	}
	bounds := code.Bounds()
	modules := bounds.Dx() + 8
	scale := size / modules
	if scale < 1 {
		scale = 1
	}
	img := image.NewGray(image.Rect(0, 0, modules*scale, modules*scale))
	draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if color.GrayModel.Convert(code.At(x, y)).(color.Gray).Y >= 128 {
				continue
			}
			left, top := (x-bounds.Min.X+4)*scale, (y-bounds.Min.Y+4)*scale
			draw.Draw(img, image.Rect(left, top, left+scale, top+scale), image.Black, image.Point{}, draw.Src)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
`
