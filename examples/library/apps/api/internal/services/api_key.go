package services

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"

	"library/apps/api/internal/models"
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
	RateLimit int
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
		RateLimit:  opts.RateLimit,
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
// verifiedKeyTTL is how long a verified key is trusted without asking the
// database again. A key revoked on another replica can keep working that long.
const verifiedKeyTTL = 30 * time.Second

// touchInterval is the least time between two last_used_at writes for a key.
const touchInterval = time.Minute

type verifiedKey struct {
	key   models.APIKey
	until time.Time
}

var (
	verifiedKeys sync.Map // prefix + ":" + secret hash -> verifiedKey
	lastTouched  sync.Map // key id -> time.Time
)

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

	// A key verified a moment ago is not a SELECT per request. The cache key holds
	// the hash of the secret presented, so a wrong secret never matches an entry.
	cacheKey := prefix + ":" + models.HashAPIKeySecret(secret)
	if cached, ok := verifiedKeys.Load(cacheKey); ok {
		entry := cached.(verifiedKey)
		if time.Now().Before(entry.until) && entry.key.Active() {
			key := entry.key
			return &key, nil
		}
		verifiedKeys.Delete(cacheKey)
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

	verifiedKeys.Store(cacheKey, verifiedKey{key: key, until: time.Now().Add(verifiedKeyTTL)})
	return &key, nil
}

// TouchAPIKey records use. Best-effort and deliberately not in the request's
// transaction — "when was this key last used" is useful, and never worth
// failing a request over.
func TouchAPIKey(db *gorm.DB, id string) {
	now := time.Now()
	// last_used_at is read by a person, to the minute. Writing it on every request
	// was an UPDATE per call on the same row, which contends under load.
	if prev, ok := lastTouched.Load(id); ok && now.Sub(prev.(time.Time)) < touchInterval {
		return
	}
	lastTouched.Store(id, now)
	if err := db.Model(&models.APIKey{}).Where("id = ?", id).UpdateColumn("last_used_at", now).Error; err != nil {
		log.Printf("api keys: recording use of key %s: %v", id, err)
	}
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
	// Stop trusting any cached verification of it on this replica now.
	verifiedKeys.Range(func(k, v any) bool {
		if v.(verifiedKey).key.ID == id {
			verifiedKeys.Delete(k)
		}
		return true
	})
	return nil
}
