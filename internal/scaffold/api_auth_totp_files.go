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
	"encoding/base32"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
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
	counter := time.Now().Unix() / Period
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
)

// TwoFactorConfig stores TOTP settings for a user.
type TwoFactorConfig struct {
	ID          uint                       ` + "`" + `gorm:"primarykey" json:"id"` + "`" + `
	UserID      string                     ` + "`" + `gorm:"size:36;uniqueIndex;not null" json:"user_id"` + "`" + `
	Secret      string                     ` + "`" + `gorm:"size:255;not null" json:"-"` + "`" + `
	Enabled     bool                       ` + "`" + `gorm:"default:false" json:"enabled"` + "`" + `
	BackupCodes datatypes.JSONSlice[string] ` + "`" + `gorm:"type:text" json:"-"` + "`" + `
	// LastUsedStep is the time step of the last code accepted. Only a later one
	// is accepted next, so a code cannot be used twice.
	LastUsedStep int64 ` + "`" + `gorm:"not null;default:0" json:"-"` + "`" + `
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
	Attempts  int       ` + "`" + `gorm:"not null;default:0" json:"-"` + "`" + `
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
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
	"{{MODULE}}/internal/totp"
)

// TOTPHandler handles two-factor authentication endpoints.
type TOTPHandler struct {
	DB          *gorm.DB
	AuthService *services.AuthService
	Issuer      string // App name for authenticator display
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
	u := user.(models.User)

	// Check if TOTP is already enabled
	var existing models.TwoFactorConfig
	if err := h.DB.WithContext(c.Request.Context()).Where("user_id = ?", userID).First(&existing).Error; err == nil && existing.Enabled {
		c.JSON(http.StatusConflict, gin.H{
			"error": gin.H{
				"code":    "TOTP_ALREADY_ENABLED",
				"message": "Two-factor authentication is already enabled",
			},
		})
		return
	}

	secret, err := totp.GenerateSecret()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "TOTP_ERROR", "message": "Failed to generate secret"},
		})
		return
	}

	uri := totp.GenerateURI(secret, u.Email, h.Issuer)

	// Medium recovery level: a phone camera reads it reliably at the size a
	// setup dialog shows it, without the density High produces.
	qrPNG, err := qrcode.Encode(uri, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "TOTP_ERROR", "message": "Failed to render the QR code"},
		})
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

// Enable verifies the initial TOTP code and activates 2FA for the user.
// Returns backup codes that the user should save.
func (h *TOTPHandler) Enable(c *gin.Context) {
	userID := c.GetString("user_id")

	var req EnableTOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	// Enabling replaces the secret, so it is refused while 2FA is on: a stolen
	// session could otherwise re-enrol the account to an authenticator the
	// attacker holds. Disabling first asks for the password.
	var existing models.TwoFactorConfig
	if err := h.DB.WithContext(c.Request.Context()).Where("user_id = ?", userID).First(&existing).Error; err == nil && existing.Enabled {
		c.JSON(http.StatusConflict, gin.H{
			"error": gin.H{
				"code":    "TOTP_ALREADY_ENABLED",
				"message": "Two-factor authentication is already enabled. Disable it first.",
			},
		})
		return
	}

	// Verify the code matches the secret
	step, valid, err := totp.ValidateCodeStep(req.Secret, req.Code)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_TOTP_CODE",
				"message": "Invalid verification code. Make sure your authenticator app is synced.",
			},
		})
		return
	}

	// Generate backup codes
	codes, hashes, err := totp.GenerateBackupCodes(0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "TOTP_ERROR", "message": "Failed to generate backup codes"},
		})
		return
	}

	// Upsert the TwoFactorConfig
	var config models.TwoFactorConfig
	if err := h.DB.WithContext(c.Request.Context()).Where("user_id = ?", userID).FirstOrCreate(&config, models.TwoFactorConfig{UserID: userID}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "TOTP_ERROR", "message": "Failed to enable two-factor authentication"},
		})
		return
	}

	config.Secret = req.Secret
	config.Enabled = true
	config.BackupCodes = hashes
	// The code that enabled 2FA is spent; it cannot also sign in.
	config.LastUsedStep = step

	if err := h.DB.WithContext(c.Request.Context()).Save(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "TOTP_ERROR", "message": "Failed to enable two-factor authentication"},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"enabled":      true,
			"backup_codes": codes,
		},
		"message": "Two-factor authentication enabled. Save your backup codes in a safe place.",
	})
}

// Verify validates a TOTP code during the login flow (after password check).
// Exchanges a pending token + valid TOTP code for real JWT tokens.
func (h *TOTPHandler) Verify(c *gin.Context) {
	var req VerifyTOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	pending, user, config, ok := h.beginSecondFactor(c, req.PendingToken)
	if !ok {
		return
	}

	step, valid, err := totp.ValidateCodeStep(config.Secret, req.Code)
	if err == nil && valid {
		// A code is good for its window, so without this the same code signs in
		// again for as long as it lasts. The step only moves forward, and the
		// update is conditional, so two requests cannot both spend one code.
		res := h.DB.WithContext(c.Request.Context()).Model(&models.TwoFactorConfig{}).
			Where("id = ? AND last_used_step < ?", config.ID, step).
			UpdateColumn("last_used_step", step)
		valid = res.Error == nil && res.RowsAffected == 1
	}
	if err != nil || !valid {
		h.failSecondFactor(pending, user)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_TOTP_CODE",
				"message": "Invalid verification code",
			},
		})
		return
	}

	h.completeSecondFactor(c, pending, user, req.TrustDevice, nil, "Logged in successfully")
}

// VerifyBackupCode validates a backup code during login (alternative to TOTP).
func (h *TOTPHandler) VerifyBackupCode(c *gin.Context) {
	var req VerifyBackupCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	pending, user, config, ok := h.beginSecondFactor(c, req.PendingToken)
	if !ok {
		return
	}

	idx := totp.VerifyBackupCode(req.Code, config.BackupCodes)
	if idx < 0 {
		h.failSecondFactor(pending, user)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_BACKUP_CODE",
				"message": "Invalid backup code",
			},
		})
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "TOTP_ERROR", "message": "Failed to use the backup code"},
		})
		return
	}
	res := h.DB.WithContext(c.Request.Context()).Model(&models.TwoFactorConfig{}).
		Where("id = ? AND backup_codes = ?", config.ID, string(read)).
		Update("backup_codes", datatypes.JSONSlice[string](remaining))
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "TOTP_ERROR", "message": "Failed to use the backup code"},
		})
		return
	}
	if res.RowsAffected != 1 {
		h.failSecondFactor(pending, user)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{"code": "INVALID_BACKUP_CODE", "message": "Invalid backup code"},
		})
		return
	}
	config.BackupCodes = remaining

	h.completeSecondFactor(c, pending, user, req.TrustDevice, gin.H{"backup_codes_remaining": len(remaining)},
		"Logged in successfully with backup code")
}

// beginSecondFactor loads what a second-factor attempt needs, and refuses one
// that must not go ahead: a pending token that is unknown, expired or out of
// attempts, or an account that is disabled or locked. Password sign-in makes
// the same account checks; before, the second step skipped them.
func (h *TOTPHandler) beginSecondFactor(c *gin.Context, pendingToken string) (*models.TOTPPendingToken, *models.User, *models.TwoFactorConfig, bool) {
	invalid := func() {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_PENDING_TOKEN",
				"message": "Invalid or expired verification session. Please log in again.",
			},
		})
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
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{"code": "ACCOUNT_DISABLED", "message": "Your account has been disabled"},
		})
		return nil, nil, nil, false
	}
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": gin.H{"code": "ACCOUNT_LOCKED", "message": "Too many failed attempts. Try again later, or reset your password."},
		})
		return nil, nil, nil, false
	}

	var config models.TwoFactorConfig
	if err := h.DB.WithContext(c.Request.Context()).Where("user_id = ? AND enabled = ?", pending.UserID, true).First(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "TOTP_ERROR", "message": "Two-factor configuration not found"},
		})
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
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_PENDING_TOKEN",
				"message": "Invalid or expired verification session. Please log in again.",
			},
		})
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "TOKEN_ERROR", "message": "Failed to generate tokens"},
		})
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
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	// Verify password
	var user models.User
	if err := h.DB.WithContext(c.Request.Context()).Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "User not found"},
		})
		return
	}

	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_PASSWORD",
				"message": "Incorrect password",
			},
		})
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "TOTP_ERROR", "message": "Failed to disable two-factor authentication"},
		})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "TOTP_NOT_ENABLED",
				"message": "Two-factor authentication is not enabled",
			},
		})
		return
	}

	codes, hashes, err := totp.GenerateBackupCodes(0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "TOTP_ERROR", "message": "Failed to generate backup codes"},
		})
		return
	}

	config.BackupCodes = hashes
	if err := h.DB.WithContext(c.Request.Context()).Save(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "TOTP_ERROR", "message": "Failed to save backup codes"},
		})
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to load trusted devices"},
		})
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to revoke the device"},
		})
		return
	}
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "Trusted device not found"},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device revoked. It will need a code at the next sign-in."})
}

// RevokeTrustedDevices removes all trusted devices for the current user.
func (h *TOTPHandler) RevokeTrustedDevices(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := h.DB.WithContext(c.Request.Context()).Where("user_id = ?", userID).Delete(&models.TrustedDevice{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "TOTP_ERROR", "message": "Failed to revoke trusted devices"},
		})
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

	c.SetCookie(
		"totp_trusted",
		deviceToken,
		int(totp.TrustedDeviceDuration.Seconds()),
		"/",
		"",    // domain
		false, // secure (set true in production)
		true,  // httpOnly
	)
}

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
`
}
