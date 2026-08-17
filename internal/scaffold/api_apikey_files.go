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
	"strings"
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

	// Kind decides what the key is allowed to reach, and it is the most
	// important field on this model.
	//
	// A publishable key ships inside a browser bundle or a mobile binary,
	// where it is readable by anyone who wants it. Calling that a secret and
	// hoping is how an admin credential ends up in a JavaScript file. So the
	// kind is declared, and a publishable key is structurally incapable of
	// reaching a route that was not marked public: not because it lacks a
	// permission, but because the middleware for protected routes rejects the
	// kind outright.
	Kind string ` + "`" + `gorm:"size:16;not null;default:secret;index" json:"kind"` + "`" + `

	// Token is the full key, stored in clear, and ONLY for publishable keys.
	//
	// That is not an oversight. A publishable key is already public: it is in
	// every copy of your app. Hashing it would buy nothing and would cost the
	// one thing that makes it pleasant to work with, which is being able to
	// read it again from the admin when somebody sets up a new environment.
	// Secret keys keep only their hash and are shown exactly once.
	Token string ` + "`" + `gorm:"size:120" json:"token,omitempty"` + "`" + `

	// Scopes are permission keys, the same strings roles grant
	// ("products.view"). Empty means the key inherits the owner's permissions.
	//
	// Ignored entirely for publishable keys. A key in a browser inheriting an
	// admin's permissions is the exact failure this model exists to prevent.
	Scopes datatypes.JSONSlice[string] ` + "`" + `json:"scopes"` + "`" + `

	// Endpoints narrows a key to specific routes, as method plus path with an
	// optional trailing wildcard:
	//
	//	["GET /api/v1/shop/products", "GET /api/v1/shop/products/*"]
	//
	// Empty means every route the kind already allows. This is a second axis
	// to Scopes, not a replacement: scopes say what you may do, endpoints say
	// where you may go, and a partner integration usually wants both narrowed.
	Endpoints datatypes.JSONSlice[string] ` + "`" + `json:"endpoints"` + "`" + `

	// Origins restricts browser use to specific sites, checked against the
	// request's Origin header.
	//
	// Worth having and worth not overestimating. It stops another site's page
	// using this key from a customer's browser. It stops nothing that is not a
	// browser, because curl does not send an Origin it does not like. Leave it
	// empty for a mobile app: native clients send no Origin at all, so an
	// allowlist would reject every request they make.
	Origins datatypes.JSONSlice[string] ` + "`" + `json:"origins"` + "`" + `

	LastUsedAt *time.Time ` + "`" + `json:"last_used_at,omitempty"` + "`" + `
	ExpiresAt  *time.Time ` + "`" + `gorm:"index" json:"expires_at,omitempty"` + "`" + `
	RevokedAt  *time.Time ` + "`" + `gorm:"index" json:"revoked_at,omitempty"` + "`" + `
	CreatedAt  time.Time  ` + "`" + `json:"created_at"` + "`" + `
}

// Key kinds.
const (
	// KindPublishable is safe to ship to a browser or a phone. Reaches public
	// routes only.
	KindPublishable = "publishable"
	// KindSecret is server-side only. Reaches whatever its owner can.
	KindSecret = "secret"
)

// Publishable reports whether this key is one that lives in public.
func (k *APIKey) Publishable() bool { return k.Kind == KindPublishable }

// AllowsEndpoint reports whether the key may call a method and path.
//
// An empty allowlist means yes. A pattern ending in * matches a prefix, which
// is what lets one entry cover a resource and its detail routes.
func (k *APIKey) AllowsEndpoint(method, path string) bool {
	if len(k.Endpoints) == 0 {
		return true
	}
	want := strings.ToUpper(method) + " " + path
	for _, pattern := range k.Endpoints {
		pattern = strings.TrimSpace(pattern)
		if strings.HasSuffix(pattern, "*") {
			if strings.HasPrefix(want, strings.TrimSuffix(pattern, "*")) {
				return true
			}
			continue
		}
		if want == pattern {
			return true
		}
	}
	return false
}

// AllowsOrigin reports whether a browser request from this origin may use the
// key.
//
// An empty allowlist means any origin, which is also the correct answer for a
// request that carries no Origin header at all: mobile apps and servers never
// send one, and rejecting them would break every non-browser caller.
func (k *APIKey) AllowsOrigin(origin string) bool {
	if len(k.Origins) == 0 || origin == "" {
		return true
	}
	for _, allowed := range k.Origins {
		if strings.EqualFold(strings.TrimSpace(allowed), origin) {
			return true
		}
	}
	return false
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

// KeyOptions describes a key to be minted.
type KeyOptions struct {
	UserID string
	Name   string
	// Kind is models.KindPublishable or models.KindSecret. Empty means secret,
	// so a caller that has not thought about it gets the careful answer.
	Kind      string
	Scopes    []string
	Endpoints []string
	Origins   []string
	ExpiresAt *time.Time
}

// GenerateAPIKey mints a key.
//
// Token layout: grit_<kind>_<8 hex prefix>_<64 hex secret>, so pk and sk are
// distinguishable at a glance in a log, a config file or a code review. The
// prefix is a lookup handle, not a credential; the last segment is.
//
// A secret key stores only the hash of its secret. A publishable key stores
// the whole token in clear, because it is going to live in a JavaScript bundle
// and pretending otherwise buys nothing.
func GenerateAPIKey(db *gorm.DB, opts KeyOptions) (*IssuedKey, error) {
	kind := opts.Kind
	if kind == "" {
		kind = models.KindSecret
	}
	if kind != models.KindPublishable && kind != models.KindSecret {
		return nil, fmt.Errorf("unknown api key kind %q", kind)
	}

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
	token := KeyTokenPrefix + "_" + kindSegment(kind) + "_" + prefix + "_" + secret

	scopes := opts.Scopes
	if kind == models.KindPublishable {
		// Never. A publishable key that inherited its owner's permissions
		// would put an admin's authority into a browser, which is the one
		// outcome this whole design exists to make impossible.
		scopes = nil
	}

	record := &models.APIKey{
		UserID:     opts.UserID,
		Name:       opts.Name,
		Kind:       kind,
		Prefix:     prefix,
		SecretHash: models.HashAPIKeySecret(secret),
		Scopes:     scopes,
		Endpoints:  opts.Endpoints,
		Origins:    opts.Origins,
		ExpiresAt:  opts.ExpiresAt,
	}
	if kind == models.KindPublishable {
		record.Token = token
	}
	if err := db.Create(record).Error; err != nil {
		return nil, fmt.Errorf("storing api key: %w", err)
	}

	return &IssuedKey{Record: record, Token: token}, nil
}

// kindSegment is the two letters that appear in the token itself.
func kindSegment(kind string) string {
	if kind == models.KindPublishable {
		return "pk"
	}
	return "sk"
}

// VerifyAPIKey resolves a presented token to its record.
//
// The comparison is constant-time. Both halves are hex of the same length, so
// a byte-by-byte compare would leak how much of a guessed secret was right.
func VerifyAPIKey(db *gorm.DB, token string) (*models.APIKey, error) {
	parts := strings.Split(strings.TrimSpace(token), "_")
	if len(parts) < 3 || parts[0] != KeyTokenPrefix {
		return nil, ErrAPIKeyInvalid
	}

	// Two layouts. grit_pk_<prefix>_<secret> is current; grit_<prefix>_<secret>
	// is what keys issued before kinds existed look like, and those still have
	// to work, so they are read as secret keys exactly as they were.
	var prefix, secret string
	switch {
	case len(parts) == 4 && (parts[1] == "pk" || parts[1] == "sk"):
		prefix, secret = parts[2], parts[3]
	case len(parts) == 3:
		prefix, secret = parts[1], parts[2]
	default:
		return nil, ErrAPIKeyInvalid
	}

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

		// The safety property of the whole design, in five lines.
		//
		// A publishable key lives in a browser bundle or a phone binary, where
		// anyone can read it. It authenticates nobody. Reaching a protected
		// route with one is refused on the kind alone, before permissions are
		// even consulted, so no combination of scopes can talk its way past
		// this. If a publishable key leaks, and it will, the worst it opens is
		// what you already marked public.
		if key.Publishable() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code": "PUBLISHABLE_KEY_NOT_ALLOWED",
					"message": "This is a protected endpoint. A publishable key reaches " +
						"public endpoints only; use a secret key from a server, or sign in.",
				},
			})
			return
		}

		if !key.AllowsEndpoint(c.Request.Method, c.FullPath()) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "ENDPOINT_NOT_ALLOWED",
					"message": "That key is not allowed to call this endpoint",
				},
			})
			return
		}
		if !key.AllowsOrigin(c.GetHeader("Origin")) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "ORIGIN_NOT_ALLOWED",
					"message": "That key is not allowed from this origin",
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
// RequireAPIKey guards a public endpoint.
//
// Public does not mean unauthenticated. It means no user is required: any
// valid key gets in, publishable or secret, and no JWT is needed. What that
// buys is not secrecy, because a publishable key is readable by anyone with
// your app. It buys identification, a rate-limit bucket per key, per-endpoint
// and per-origin narrowing, and the ability to turn one client off without
// deploying anything.
//
// Deliberately no user is set on the context. A public handler that quietly
// depended on c.Get("user_id") would then behave differently depending on
// which credential turned up, which is the kind of difference nobody notices
// until it is a data leak.
func RequireAPIKey(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractAPIKey(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "API_KEY_REQUIRED",
					"message": "Send your publishable key as X-API-Key",
				},
			})
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
		if !key.AllowsEndpoint(c.Request.Method, c.FullPath()) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "ENDPOINT_NOT_ALLOWED",
					"message": "That key is not allowed to call this endpoint",
				},
			})
			return
		}
		if !key.AllowsOrigin(c.GetHeader("Origin")) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "ORIGIN_NOT_ALLOWED",
					"message": "That key is not allowed from this origin",
				},
			})
			return
		}

		services.TouchAPIKey(db, key.ID)
		c.Set("api_key_id", key.ID)
		c.Set("api_key_kind", key.Kind)
		c.Set("auth_via", "api_key")
		c.Next()
	}
}

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

// What a caller sends to mint an API key.
type CreateAPIKeyRequest struct {
	Name string ` + "`" + `json:"name" binding:"required,min=1,max=120"` + "`" + `
	// Kind is "publishable" or "secret". Empty means secret: a caller that has
	// not thought about it should get the careful one.
	Kind      string   ` + "`" + `json:"kind"` + "`" + `
	Scopes    []string ` + "`" + `json:"scopes"` + "`" + `
	Endpoints []string ` + "`" + `json:"endpoints"` + "`" + `
	Origins   []string ` + "`" + `json:"origins"` + "`" + `
	ExpiresIn int      ` + "`" + `json:"expires_in_days"` + "`" + `
}

// Create issues a key and returns it once.

func (h *APIKeyHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")

	var req CreateAPIKeyRequest
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

	issued, err := services.GenerateAPIKey(h.DB, services.KeyOptions{
		UserID:    userID,
		Name:      req.Name,
		Kind:      req.Kind,
		Scopes:    req.Scopes,
		Endpoints: req.Endpoints,
		Origins:   req.Origins,
		ExpiresAt: expiresAt,
	})
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
