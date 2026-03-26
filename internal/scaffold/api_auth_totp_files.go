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
	"crypto/sha1"
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
	// TrustedDeviceDuration is how long a trusted device cookie lasts.
	TrustedDeviceDuration = 30 * 24 * time.Hour // 30 days
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
	now := time.Now().Unix()
	counter := now / Period

	for i := -int64(Window); i <= int64(Window); i++ {
		expected, err := generateCode(secret, counter+i)
		if err != nil {
			return false, err
		}
		if expected == code {
			return true, nil
		}
	}
	return false, nil
}

// generateCode computes the HOTP code for a given counter (RFC 4226).
func generateCode(secret string, counter int64) (string, error) {
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", fmt.Errorf("decoding secret: %w", err)
	}

	// Counter to big-endian 8 bytes
	buf := make([]byte, 8)
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
	UserID      uint                       ` + "`" + `gorm:"uniqueIndex;not null" json:"user_id"` + "`" + `
	Secret      string                     ` + "`" + `gorm:"size:255;not null" json:"-"` + "`" + `
	Enabled     bool                       ` + "`" + `gorm:"default:false" json:"enabled"` + "`" + `
	BackupCodes datatypes.JSONSlice[string] ` + "`" + `gorm:"type:text" json:"-"` + "`" + `
	CreatedAt   time.Time                  ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt   time.Time                  ` + "`" + `json:"updated_at"` + "`" + `
}

// TrustedDevice stores a remembered device that can skip TOTP.
type TrustedDevice struct {
	ID        uint           ` + "`" + `gorm:"primarykey" json:"id"` + "`" + `
	UserID    uint           ` + "`" + `gorm:"index;not null" json:"user_id"` + "`" + `
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
	UserID    uint      ` + "`" + `gorm:"index;not null" json:"user_id"` + "`" + `
	TokenHash string    ` + "`" + `gorm:"size:64;uniqueIndex;not null" json:"-"` + "`" + `
	ExpiresAt time.Time ` + "`" + `gorm:"not null" json:"expires_at"` + "`" + `
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
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
}

type totpEnableRequest struct {
	Secret string ` + "`" + `json:"secret" binding:"required"` + "`" + `
	Code   string ` + "`" + `json:"code" binding:"required"` + "`" + `
}

type totpVerifyRequest struct {
	PendingToken string ` + "`" + `json:"pending_token" binding:"required"` + "`" + `
	Code         string ` + "`" + `json:"code" binding:"required"` + "`" + `
	TrustDevice  bool   ` + "`" + `json:"trust_device"` + "`" + `
}

type totpDisableRequest struct {
	Password string ` + "`" + `json:"password" binding:"required"` + "`" + `
}

type backupCodeVerifyRequest struct {
	PendingToken string ` + "`" + `json:"pending_token" binding:"required"` + "`" + `
	Code         string ` + "`" + `json:"code" binding:"required"` + "`" + `
	TrustDevice  bool   ` + "`" + `json:"trust_device"` + "`" + `
}

// Setup generates a new TOTP secret and QR code URI for the user.
// The secret is NOT stored yet — the user must verify it first via Enable.
func (h *TOTPHandler) Setup(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, _ := c.Get("user")
	u := user.(models.User)

	// Check if TOTP is already enabled
	var existing models.TwoFactorConfig
	if err := h.DB.Where("user_id = ?", userID).First(&existing).Error; err == nil && existing.Enabled {
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

	c.JSON(http.StatusOK, gin.H{
		"data": totpSetupResponse{
			Secret: secret,
			URI:    uri,
		},
		"message": "Scan the QR code with your authenticator app, then verify with a code",
	})
}

// Enable verifies the initial TOTP code and activates 2FA for the user.
// Returns backup codes that the user should save.
func (h *TOTPHandler) Enable(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req totpEnableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	// Verify the code matches the secret
	valid, err := totp.ValidateCode(req.Secret, req.Code)
	if err != nil || !valid {
		c.JSON(http.StatusBadRequest, gin.H{
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
	h.DB.Where("user_id = ?", userID).FirstOrCreate(&config, models.TwoFactorConfig{UserID: userID})

	config.Secret = req.Secret
	config.Enabled = true
	config.BackupCodes = hashes

	if err := h.DB.Save(&config).Error; err != nil {
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
	var req totpVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	// Look up the pending token
	tokenHash := totp.HashToken(req.PendingToken)
	var pending models.TOTPPendingToken
	if err := h.DB.Where("token_hash = ? AND expires_at > ?", tokenHash, time.Now()).First(&pending).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_PENDING_TOKEN",
				"message": "Invalid or expired verification session. Please log in again.",
			},
		})
		return
	}

	// Get the user's TOTP config
	var config models.TwoFactorConfig
	if err := h.DB.Where("user_id = ? AND enabled = ?", pending.UserID, true).First(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "TOTP_ERROR", "message": "Two-factor configuration not found"},
		})
		return
	}

	// Validate the TOTP code
	valid, err := totp.ValidateCode(config.Secret, req.Code)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_TOTP_CODE",
				"message": "Invalid verification code",
			},
		})
		return
	}

	// Delete the pending token (one-time use)
	h.DB.Delete(&pending)

	// Load the user for token generation
	var user models.User
	if err := h.DB.First(&user, pending.UserID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "USER_ERROR", "message": "User not found"},
		})
		return
	}

	// Generate real JWT tokens
	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "TOKEN_ERROR", "message": "Failed to generate tokens"},
		})
		return
	}

	// If user wants to trust this device, create a trusted device cookie
	if req.TrustDevice {
		h.createTrustedDevice(c, user.ID)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user":   user,
			"tokens": tokens,
		},
		"message": "Logged in successfully",
	})
}

// VerifyBackupCode validates a backup code during login (alternative to TOTP).
func (h *TOTPHandler) VerifyBackupCode(c *gin.Context) {
	var req backupCodeVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	// Look up the pending token
	tokenHash := totp.HashToken(req.PendingToken)
	var pending models.TOTPPendingToken
	if err := h.DB.Where("token_hash = ? AND expires_at > ?", tokenHash, time.Now()).First(&pending).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_PENDING_TOKEN",
				"message": "Invalid or expired verification session. Please log in again.",
			},
		})
		return
	}

	// Get the user's TOTP config
	var config models.TwoFactorConfig
	if err := h.DB.Where("user_id = ? AND enabled = ?", pending.UserID, true).First(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "TOTP_ERROR", "message": "Two-factor configuration not found"},
		})
		return
	}

	// Verify the backup code
	idx := totp.VerifyBackupCode(req.Code, config.BackupCodes)
	if idx < 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_BACKUP_CODE",
				"message": "Invalid backup code",
			},
		})
		return
	}

	// Remove the used backup code (one-time use)
	config.BackupCodes = append(config.BackupCodes[:idx], config.BackupCodes[idx+1:]...)
	h.DB.Save(&config)

	// Delete the pending token
	h.DB.Delete(&pending)

	// Load user and generate tokens
	var user models.User
	if err := h.DB.First(&user, pending.UserID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "USER_ERROR", "message": "User not found"},
		})
		return
	}

	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "TOKEN_ERROR", "message": "Failed to generate tokens"},
		})
		return
	}

	if req.TrustDevice {
		h.createTrustedDevice(c, user.ID)
	}

	remaining := len(config.BackupCodes)
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user":                  user,
			"tokens":               tokens,
			"backup_codes_remaining": remaining,
		},
		"message": "Logged in successfully with backup code",
	})
}

// Disable turns off 2FA for the user (requires password confirmation).
func (h *TOTPHandler) Disable(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req totpDisableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	// Verify password
	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
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

	// Delete TOTP config, trusted devices, and pending tokens
	h.DB.Where("user_id = ?", userID).Delete(&models.TwoFactorConfig{})
	h.DB.Where("user_id = ?", userID).Delete(&models.TrustedDevice{})
	h.DB.Where("user_id = ?", userID).Delete(&models.TOTPPendingToken{})

	c.JSON(http.StatusOK, gin.H{
		"message": "Two-factor authentication disabled",
	})
}

// Status returns whether TOTP is enabled for the current user.
func (h *TOTPHandler) Status(c *gin.Context) {
	userID := c.GetUint("user_id")

	var config models.TwoFactorConfig
	enabled := false
	backupCodesRemaining := 0

	if err := h.DB.Where("user_id = ?", userID).First(&config).Error; err == nil {
		enabled = config.Enabled
		backupCodesRemaining = len(config.BackupCodes)
	}

	// Count trusted devices
	var deviceCount int64
	h.DB.Model(&models.TrustedDevice{}).Where("user_id = ? AND expires_at > ?", userID, time.Now()).Count(&deviceCount)

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
	userID := c.GetUint("user_id")

	var config models.TwoFactorConfig
	if err := h.DB.Where("user_id = ? AND enabled = ?", userID, true).First(&config).Error; err != nil {
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
	if err := h.DB.Save(&config).Error; err != nil {
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

// RevokeTrustedDevices removes all trusted devices for the current user.
func (h *TOTPHandler) RevokeTrustedDevices(c *gin.Context) {
	userID := c.GetUint("user_id")

	h.DB.Where("user_id = ?", userID).Delete(&models.TrustedDevice{})

	c.JSON(http.StatusOK, gin.H{
		"message": "All trusted devices revoked. You will need to enter a TOTP code on next login.",
	})
}

// createTrustedDevice stores a trusted device and sets the cookie.
func (h *TOTPHandler) createTrustedDevice(c *gin.Context, userID uint) {
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

	if err := h.DB.Create(&device).Error; err != nil {
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
func IsTrustedDevice(c *gin.Context, db *gorm.DB, userID uint) bool {
	token, err := c.Cookie("totp_trusted")
	if err != nil || token == "" {
		return false
	}

	tokenHash := totp.HashToken(token)
	var device models.TrustedDevice
	if err := db.Where("user_id = ? AND token_hash = ? AND expires_at > ?", userID, tokenHash, time.Now()).First(&device).Error; err != nil {
		return false
	}

	// Refresh the device expiry (sliding window)
	device.ExpiresAt = time.Now().Add(totp.TrustedDeviceDuration)
	db.Save(&device)

	return true
}
`
}
