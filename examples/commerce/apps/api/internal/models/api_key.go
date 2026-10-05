package models

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"commerce/apps/api/internal/ids"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// APIKey is a long-lived credential for machine callers.
//
// Prefix is stored in clear and indexed: it identifies which key is being
// presented so verification is a single indexed lookup, not a scan. The secret
// itself exists only as a hash — issuing is the one moment it is knowable.
type APIKey struct {
	ID     string `gorm:"primarykey;size:36" json:"id"`
	UserID string `gorm:"size:36;index;not null" json:"user_id"`

	// Name is for humans: "nightly export", "Zapier". Shown in the list so a
	// key can be revoked without guessing what it does.
	Name string `gorm:"size:120;not null" json:"name"`

	Prefix     string `gorm:"size:16;uniqueIndex;not null" json:"prefix"`
	SecretHash string `gorm:"size:64;not null" json:"-"`

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
	Kind string `gorm:"size:16;not null;default:secret;index" json:"kind"`

	// Token is the full key, stored in clear, and ONLY for publishable keys.
	//
	// That is not an oversight. A publishable key is already public: it is in
	// every copy of your app. Hashing it would buy nothing and would cost the
	// one thing that makes it pleasant to work with, which is being able to
	// read it again from the admin when somebody sets up a new environment.
	// Secret keys keep only their hash and are shown exactly once.
	Token string `gorm:"size:120" json:"token,omitempty"`

	// Scopes are permission keys, the same strings roles grant
	// ("products.view"). Empty means the key inherits the owner's permissions.
	//
	// Ignored entirely for publishable keys. A key in a browser inheriting an
	// admin's permissions is the exact failure this model exists to prevent.
	Scopes datatypes.JSONSlice[string] `json:"scopes"`

	// Endpoints narrows a key to specific routes, as method plus path with an
	// optional trailing wildcard:
	//
	//	["GET /api/v1/shop/products", "GET /api/v1/shop/products/*"]
	//
	// Empty means every route the kind already allows. This is a second axis
	// to Scopes, not a replacement: scopes say what you may do, endpoints say
	// where you may go, and a partner integration usually wants both narrowed.
	Endpoints datatypes.JSONSlice[string] `json:"endpoints"`

	// Origins restricts browser use to specific sites, checked against the
	// request's Origin header.
	//
	// Worth having and worth not overestimating. It stops another site's page
	// using this key from a customer's browser. It stops nothing that is not a
	// browser, because curl does not send an Origin it does not like. Leave it
	// empty for a mobile app: native clients send no Origin at all, so an
	// allowlist would reject every request they make.
	Origins datatypes.JSONSlice[string] `json:"origins"`

	// RateLimit is requests per minute for this key alone. Zero means no
	// per-key limit, and the global Sentinel IP limit still applies either way.
	//
	// Per key rather than per IP, because those answer different questions. An
	// IP limit protects the server from a flood. A key limit protects you from
	// one client: a partner integration polling every second, or a storefront
	// with a render loop, throttled without touching the limit that applies to
	// everybody else.
	RateLimit int `gorm:"default:0" json:"rate_limit"`

	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt  *time.Time `gorm:"index" json:"expires_at,omitempty"`
	RevokedAt  *time.Time `gorm:"index" json:"revoked_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
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
