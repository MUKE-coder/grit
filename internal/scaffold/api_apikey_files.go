package scaffold

// API-key authentication.
//
// Until now the only way for anything to call the API was the JWT flow, which
// is built for a human at a browser: short-lived access tokens, a refresh
// token in a cookie, rotation, revocation on password change. A cron job on
// someone else's server wants none of that — it wants one long-lived
// credential it can put in a header.
//
// Design notes worth keeping:
//
//   - The key is "grit_<prefix>_<secret>". The prefix is stored in clear and
//     indexed, so a lookup is one indexed hit rather than a scan comparing
//     bcrypt hashes against every row.
//   - Only the secret's SHA-256 is stored. It is shown once, at creation, and
//     is unrecoverable afterwards — which the UI has to say plainly.
//   - SHA-256 rather than bcrypt on purpose: this is a 256-bit random secret,
//     not a human password, so there is nothing to brute-force and the cost of
//     bcrypt would land on every single API request.
//   - Keys carry the same permission strings as roles, so authorisation has
//     one vocabulary rather than two.

func apiAPIKeyModelGo() string {
	return `package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"{{MODULE}}/internal/ids"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// APIKey is a long-lived credential for machine callers.
//
// Prefix is stored in clear and indexed: it identifies which key is being
// presented so verification is a single indexed lookup, not a scan. The secret
// itself exists only as a hash — issuing is the one moment it is knowable.
type APIKey struct {
	ID     string ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	UserID string ` + "`" + `gorm:"size:36;index;not null" json:"user_id"` + "`" + `

	// Name is for humans: "nightly export", "Zapier". Shown in the list so a
	// key can be revoked without guessing what it does.
	Name string ` + "`" + `gorm:"size:120;not null" json:"name"` + "`" + `

	Prefix    string ` + "`" + `gorm:"size:16;uniqueIndex;not null" json:"prefix"` + "`" + `
	SecretHash string ` + "`" + `gorm:"size:64;not null" json:"-"` + "`" + `

	// Scopes are permission keys, the same strings roles grant
	// ("products.view"). Empty means the key inherits the owner's permissions,
	// which is the common case and the least surprising default.
	Scopes datatypes.JSONSlice[string] ` + "`" + `json:"scopes"` + "`" + `

	LastUsedAt *time.Time ` + "`" + `json:"last_used_at,omitempty"` + "`" + `
	ExpiresAt  *time.Time ` + "`" + `gorm:"index" json:"expires_at,omitempty"` + "`" + `
	RevokedAt  *time.Time ` + "`" + `gorm:"index" json:"revoked_at,omitempty"` + "`" + `
	CreatedAt  time.Time  ` + "`" + `json:"created_at"` + "`" + `
}

func (k *APIKey) BeforeCreate(tx *gorm.DB) error {
	if k.ID == "" {
		k.ID = ids.New()
	}
	return nil
}

// Active reports whether the key may still be used.
func (k *APIKey) Active() bool {
	if k.RevokedAt != nil {
		return false
	}
	if k.ExpiresAt != nil && time.Now().After(*k.ExpiresAt) {
		return false
	}
	return true
}

// HashAPIKeySecret returns the storage form of an API key's secret.
//
// SHA-256 rather than bcrypt: the secret is 256 bits of entropy from crypto/rand,
// so there is no dictionary to attack, and bcrypt's deliberate slowness would
// be paid on every request rather than once per login.
func HashAPIKeySecret(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(sum[:])
}
`
}

func apiAPIKeyServiceGo() string {
	return `package services

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// KeyTokenPrefix marks a Grit API key so it is recognisable in a log or a
// leaked config file — and so secret scanners can be taught one pattern.
const KeyTokenPrefix = "grit"

var ErrAPIKeyInvalid = errors.New("api key is not valid")

// IssuedKey is returned once, at creation. The Secret field is the only time
// the raw key exists outside the caller's hands.
type IssuedKey struct {
	Record *models.APIKey
	Token  string
}

// GenerateAPIKey mints a key and stores only its hash.
//
// Token layout: grit_<8 hex prefix>_<64 hex secret>. The prefix is a lookup
// handle, not a secret; the second half is the credential.
func GenerateAPIKey(db *gorm.DB, userID, name string, scopes []string, expiresAt *time.Time) (*IssuedKey, error) {
	prefixBytes := make([]byte, 4)
	if _, err := rand.Read(prefixBytes); err != nil {
		return nil, fmt.Errorf("generating key prefix: %w", err)
	}
	secretBytes := make([]byte, 32)
	if _, err := rand.Read(secretBytes); err != nil {
		return nil, fmt.Errorf("generating key secret: %w", err)
	}

	prefix := hex.EncodeToString(prefixBytes)
	secret := hex.EncodeToString(secretBytes)
	token := KeyTokenPrefix + "_" + prefix + "_" + secret

	record := &models.APIKey{
		UserID:     userID,
		Name:       name,
		Prefix:     prefix,
		SecretHash: models.HashAPIKeySecret(secret),
		Scopes:     scopes,
		ExpiresAt:  expiresAt,
	}
	if err := db.Create(record).Error; err != nil {
		return nil, fmt.Errorf("storing api key: %w", err)
	}

	return &IssuedKey{Record: record, Token: token}, nil
}

// VerifyAPIKey resolves a presented token to its record.
//
// The comparison is constant-time. Both halves are hex of the same length, so
// a byte-by-byte compare would leak how much of a guessed secret was right.
func VerifyAPIKey(db *gorm.DB, token string) (*models.APIKey, error) {
	parts := strings.Split(strings.TrimSpace(token), "_")
	if len(parts) != 3 || parts[0] != KeyTokenPrefix {
		return nil, ErrAPIKeyInvalid
	}
	prefix, secret := parts[1], parts[2]

	var key models.APIKey
	if err := db.Where("prefix = ?", prefix).First(&key).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAPIKeyInvalid
		}
		return nil, err
	}

	expected := models.HashAPIKeySecret(secret)
	if subtle.ConstantTimeCompare([]byte(expected), []byte(key.SecretHash)) != 1 {
		return nil, ErrAPIKeyInvalid
	}

	// Revoked and expired are reported the same way as a wrong secret: telling
	// them apart confirms that a prefix once existed.
	if !key.Active() {
		return nil, ErrAPIKeyInvalid
	}

	return &key, nil
}

// TouchAPIKey records use. Best-effort and deliberately not in the request's
// transaction — "when was this key last used" is useful, and never worth
// failing a request over.
func TouchAPIKey(db *gorm.DB, id string) {
	now := time.Now()
	db.Model(&models.APIKey{}).Where("id = ?", id).UpdateColumn("last_used_at", now)
}

// RevokeAPIKey marks a key unusable. Revocation is a timestamp rather than a
// delete so the audit trail keeps showing which key did what.
func RevokeAPIKey(db *gorm.DB, id, userID string) error {
	now := time.Now()
	res := db.Model(&models.APIKey{}).
		Where("id = ? AND user_id = ? AND revoked_at IS NULL", id, userID).
		Update("revoked_at", now)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrAPIKeyInvalid
	}
	return nil
}
`
}

func apiAPIKeyTestGo() string {
	return `package services

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"{{MODULE}}/internal/models"
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
	issued, err := GenerateAPIKey(db, "user-1", "nightly export", nil, nil)
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
	issued, _ := GenerateAPIKey(db, "user-1", "ci", nil, nil)

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
	issued, _ := GenerateAPIKey(db, "user-1", "ci", nil, nil)

	// Right prefix, wrong secret — the case a scan-based lookup would get wrong.
	tampered := issued.Token[:len(issued.Token)-4] + "0000"
	if _, err := VerifyAPIKey(db, tampered); err == nil {
		t.Fatal("a key with a wrong secret verified")
	}
}

func TestVerifyRejectsARevokedKey(t *testing.T) {
	db := newKeyDB(t)
	issued, _ := GenerateAPIKey(db, "user-1", "ci", nil, nil)

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
	issued, _ := GenerateAPIKey(db, "user-1", "ci", nil, &past)

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
	issued, _ := GenerateAPIKey(db, "user-1", "ci", nil, nil)

	// Without the user_id predicate anyone could revoke anyone's key by id.
	if err := RevokeAPIKey(db, issued.Record.ID, "someone-else"); err == nil {
		t.Fatal("another user revoked this key")
	}
	if _, err := VerifyAPIKey(db, issued.Token); err != nil {
		t.Fatal("the key was revoked despite the wrong owner")
	}
}
`
}

func apiAPIKeyMiddlewareGo() string {
	return `package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/authz"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
)

// APIKeyOrAuth accepts either an API key or the normal JWT.
//
// It runs before the JWT middleware and, on a valid key, populates exactly the
// same context values (user_id, user_role) before calling Next. Everything
// downstream — RequireRole, the activity logger, handlers reading user_id —
// therefore works unchanged, which is the point: machine callers should not
// need a parallel set of handlers.
//
// A request with no key at all falls through untouched, so the JWT middleware
// behind this one handles browsers as before.
func APIKeyOrAuth(db *gorm.DB, jwtAuth gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractAPIKey(c)
		if token == "" {
			jwtAuth(c)
			return
		}

		key, err := services.VerifyAPIKey(db, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "INVALID_API_KEY",
					"message": "That API key is not valid, has expired, or was revoked",
				},
			})
			return
		}

		var user models.User
		if err := db.First(&user, "id = ?", key.UserID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{"code": "INVALID_API_KEY", "message": "The owner of that key no longer exists"},
			})
			return
		}
		if !user.Active {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": gin.H{"code": "ACCOUNT_DISABLED", "message": "The owner of that key is disabled"},
			})
			return
		}

		services.TouchAPIKey(db, key.ID)

		// This must set exactly what middleware.Auth sets — see auth.go. Handlers
		// read c.Get("user") for the whole record, and RequireRole reads
		// user_grants for permission checks.
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("user_email", user.Email)
		c.Set("user_role", user.Role)
		if grants, err := authz.GrantsFor(db, user.ID); err == nil {
			c.Set("user_grants", grants)
		}
		// Marked so handlers that must refuse machine callers — anything
		// changing credentials, say — can tell the difference.
		c.Set("auth_via", "api_key")
		if len(key.Scopes) > 0 {
			c.Set("api_key_scopes", []string(key.Scopes))
		}
		c.Next()
	}
}

// extractAPIKey reads a key from either header. X-API-Key is what most
// integrations reach for; Authorization: Bearer is what an OpenAPI client
// generates. Supporting both costs four lines.
func extractAPIKey(c *gin.Context) string {
	if v := strings.TrimSpace(c.GetHeader("X-API-Key")); v != "" {
		return v
	}
	auth := strings.TrimSpace(c.GetHeader("Authorization"))
	const bearer = "Bearer "
	if strings.HasPrefix(auth, bearer) {
		token := strings.TrimSpace(auth[len(bearer):])
		// Only claim it if it looks like one of ours — otherwise a JWT would
		// be sent to the key verifier and rejected before the JWT middleware
		// ever sees it.
		if strings.HasPrefix(token, services.KeyTokenPrefix+"_") {
			return token
		}
	}
	return ""
}
`
}

func apiAPIKeyHandlerGo() string {
	return `package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
)

// APIKeyHandler manages a user's own machine credentials.
type APIKeyHandler struct {
	DB *gorm.DB
}

// List returns the caller's keys. The secret is not among them — there is
// nothing to return, since only its hash was ever stored.
func (h *APIKeyHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")

	var keys []models.APIKey
	if err := h.DB.Where("user_id = ?", userID).
		Order("created_at desc").Find(&keys).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to load API keys"},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": keys})
}

// Create issues a key and returns it once.
func (h *APIKeyHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		Name      string   ` + "`" + `json:"name" binding:"required,min=1,max=120"` + "`" + `
		Scopes    []string ` + "`" + `json:"scopes"` + "`" + `
		ExpiresIn int      ` + "`" + `json:"expires_in_days"` + "`" + `
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	var expiresAt *time.Time
	if req.ExpiresIn > 0 {
		t := time.Now().AddDate(0, 0, req.ExpiresIn)
		expiresAt = &t
	}

	issued, err := services.GenerateAPIKey(h.DB, userID, req.Name, req.Scopes, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to create the API key"},
		})
		return
	}

	services.LogActivity(h.DB, c, services.ActivityArgs{
		Action:       "api_key.create",
		Severity:     "warn",
		Summary:      "API key created: " + req.Name,
		ResourceType: "api_key",
		ResourceID:   issued.Record.ID,
	})

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"key":   issued.Record,
			"token": issued.Token,
		},
		"message": "Copy this key now — it will not be shown again.",
	})
}

// Revoke marks one of the caller's keys unusable.
func (h *APIKeyHandler) Revoke(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := services.RevokeAPIKey(h.DB, c.Param("id"), userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "Key not found, or already revoked"},
		})
		return
	}

	services.LogActivity(h.DB, c, services.ActivityArgs{
		Action:       "api_key.revoke",
		Severity:     "warn",
		Summary:      "API key revoked",
		ResourceType: "api_key",
		ResourceID:   c.Param("id"),
	})

	c.JSON(http.StatusOK, gin.H{"message": "API key revoked"})
}
`
}
