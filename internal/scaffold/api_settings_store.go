package scaffold

// APISettingsRegistryGo and friends are exported so grit generate and grit add
// can install the settings module into a project that predates it.
func APISettingsRegistryGo() string { return apiSettingsRegistryGo() }

// APISettingsStoreGo emits the store: persistence, resolution order and cache.
func APISettingsStoreGo() string { return apiSettingsStoreGo() }

// APISettingsModelGo emits the settings table.
func APISettingsModelGo() string { return apiSettingsModelGo() }

// APISettingsHandlerGo emits the read/write endpoints.
func APISettingsHandlerGo() string { return apiSettingsHandlerGo() }

// APISettingsDefaultsGo emits the starter declarations.
func APISettingsDefaultsGo() string { return apiSettingsDefaultsGo() }

// apiSettingsModelGo emits internal/models/setting.go.
func apiSettingsModelGo() string {
	return `package models

import "time"

// Setting is one stored override.
//
// Keyed by (key, scope, scope_id) rather than by key alone, because the same
// setting can hold a global value, a per-tenant value and a per-user value at
// once, and resolution picks the most specific one that exists.
//
// ScopeID is empty for a global row. A partial unique index would express
// that better on Postgres, but a plain composite unique works on all three
// dialects Grit supports and the empty string is a perfectly good sentinel
// for "no scope".
type Setting struct {
	ID        string    ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	Key       string    ` + "`" + `gorm:"size:100;not null;uniqueIndex:idx_setting_scope" json:"key"` + "`" + `
	Scope     string    ` + "`" + `gorm:"size:16;not null;default:global;uniqueIndex:idx_setting_scope" json:"scope"` + "`" + `
	ScopeID   string    ` + "`" + `gorm:"size:36;not null;default:'';uniqueIndex:idx_setting_scope" json:"scope_id"` + "`" + `
	Value     string    ` + "`" + `gorm:"type:text" json:"value"` + "`" + `
	UpdatedBy string    ` + "`" + `gorm:"size:36" json:"updated_by,omitempty"` + "`" + `
	CreatedAt time.Time ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time ` + "`" + `json:"updated_at"` + "`" + `
}

func (Setting) TableName() string { return "settings" }
`
}

// apiSettingsStoreGo emits internal/settings/store.go.
func apiSettingsStoreGo() string {
	return `package settings

import (
	"context"
	"strconv"
	"sync"
	"time"

	"gorm.io/gorm"

	"{{MODULE}}/internal/ids"
	"{{MODULE}}/internal/models"
)

// Context keys the resolution order reads. Set by the auth middleware and,
// where the multitenant plugin is installed, by its scoping middleware.
type ctxKey string

const (
	// CtxUserID scopes a lookup to one person.
	CtxUserID ctxKey = "settings.user_id"
	// CtxTenantID scopes a lookup to one organisation.
	CtxTenantID ctxKey = "settings.tenant_id"
)

// WithUser returns a context that resolves user-scoped settings.
func WithUser(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, CtxUserID, userID)
}

// WithTenant returns a context that resolves tenant-scoped settings.
func WithTenant(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, CtxTenantID, tenantID)
}

// Store reads and writes settings, with a cache in front.
//
// A setting is read on nearly every request that touches the thing it
// configures, and going to Postgres for a string each time is silly. The
// cache is invalidated on write rather than expiring, so a change made in the
// admin takes effect on the next request rather than up to a TTL later.
type Store struct {
	db *gorm.DB

	mu     sync.RWMutex
	cache  map[string]string
	loaded bool
}

var defaultStore *Store

// Init creates the process store. Call once at boot, after Define.
func Init(db *gorm.DB) *Store {
	defaultStore = &Store{db: db, cache: map[string]string{}}
	return defaultStore
}

// Default returns the process store, or nil before Init.
func Default() *Store { return defaultStore }

// cacheKey identifies one stored row.
func cacheKey(key string, scope Scope, scopeID string) string {
	return string(scope) + ":" + scopeID + ":" + key
}

// load fills the cache from the database once.
//
// Every override in one query rather than one query per key. A settings table
// holds tens of rows, not millions, and the alternative is a lazy per-key read
// that turns a page rendering twelve settings into twelve round trips.
func (s *Store) load() {
	s.mu.RLock()
	if s.loaded {
		s.mu.RUnlock()
		return
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.loaded {
		return
	}
	var rows []models.Setting
	if err := s.db.Find(&rows).Error; err == nil {
		for _, r := range rows {
			s.cache[cacheKey(r.Key, Scope(r.Scope), r.ScopeID)] = r.Value
		}
	}
	s.loaded = true
}

// Raw resolves a setting to its string value.
//
// The order is: user override, then tenant override, then the stored global,
// then the environment, then the declared default. Most specific wins, which
// is what makes a per-user locale and a per-tenant currency work without
// either knowing the other exists.
//
// The environment sits below the stored value and above the default, so a
// container can supply what the application boots with and an admin can still
// change it afterwards. That order is what makes this a settings page rather
// than a read-only view of the deployment: something that genuinely must not
// be changeable belongs in config, not here.
func (s *Store) Raw(ctx context.Context, key string) string {
	declared, ok := Get(key)
	if !ok {
		return ""
	}
	s.load()

	s.mu.RLock()
	defer s.mu.RUnlock()

	if declared.Scope == User {
		if uid, _ := ctx.Value(CtxUserID).(string); uid != "" {
			if v, ok := s.cache[cacheKey(key, User, uid)]; ok {
				return v
			}
		}
	}
	if declared.Scope == User || declared.Scope == Tenant {
		if tid, _ := ctx.Value(CtxTenantID).(string); tid != "" {
			if v, ok := s.cache[cacheKey(key, Tenant, tid)]; ok {
				return v
			}
		}
	}
	if v, ok := s.cache[cacheKey(key, Global, "")]; ok {
		return v
	}
	if v, ok := declared.EnvOverride(); ok {
		return v
	}
	return declared.Default
}

// Set stores a value after validating it.
func (s *Store) Set(key string, scope Scope, scopeID, value, actor string) error {
	declared, ok := Get(key)
	if !ok {
		return ErrUnknownSetting{Key: key}
	}
	if err := declared.Parse(value); err != nil {
		return err
	}
	// A setting declared global cannot be overridden per tenant or per user.
	// Storing the row anyway would be worse than refusing: resolution would
	// never read it and nobody could tell why their change did nothing.
	if scope != Global && declared.Scope == Global {
		return ErrScopeNotAllowed{Key: key, Scope: scope, Allowed: declared.Scope}
	}
	if scope == User && declared.Scope == Tenant {
		return ErrScopeNotAllowed{Key: key, Scope: scope, Allowed: declared.Scope}
	}

	row := models.Setting{
		Key: key, Scope: string(scope), ScopeID: scopeID,
		Value: value, UpdatedBy: actor,
	}
	err := s.db.Where("key = ? AND scope = ? AND scope_id = ?", key, string(scope), scopeID).
		Assign(map[string]interface{}{"value": value, "updated_by": actor, "updated_at": time.Now()}).
		Attrs(map[string]interface{}{"id": ids.New()}).
		FirstOrCreate(&row).Error
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.cache[cacheKey(key, scope, scopeID)] = value
	s.mu.Unlock()
	return nil
}

// Unset removes an override so the next level down applies again.
func (s *Store) Unset(key string, scope Scope, scopeID string) error {
	if err := s.db.Where("key = ? AND scope = ? AND scope_id = ?", key, string(scope), scopeID).
		Delete(&models.Setting{}).Error; err != nil {
		return err
	}
	s.mu.Lock()
	delete(s.cache, cacheKey(key, scope, scopeID))
	s.mu.Unlock()
	return nil
}

// Reload drops the cache. For a process that did not make the change itself:
// a second API instance, or a value written by a migration.
func (s *Store) Reload() {
	s.mu.Lock()
	s.cache = map[string]string{}
	s.loaded = false
	s.mu.Unlock()
}

// ── Errors ───────────────────────────────────────────────────────────────

// ErrUnknownSetting is returned for a key nothing declared.
type ErrUnknownSetting struct{ Key string }

func (e ErrUnknownSetting) Error() string {
	return "no setting named " + e.Key + " is declared"
}

// ErrScopeNotAllowed is returned when a value is stored at a scope the
// declaration does not permit.
type ErrScopeNotAllowed struct {
	Key     string
	Scope   Scope
	Allowed Scope
}

func (e ErrScopeNotAllowed) Error() string {
	return e.Key + " cannot be set per " + string(e.Scope) +
		": it is declared " + string(e.Allowed)
}

// ── Typed accessors ──────────────────────────────────────────────────────
//
// Package-level, reading the process store, so calling code does not thread a
// store through. A read before Init returns the declared default rather than
// panicking, which keeps a unit test that never set one up working.

// String returns a setting's value.
func String(ctx context.Context, key string) string {
	if defaultStore == nil {
		if s, ok := Get(key); ok {
			if v, ok := s.EnvOverride(); ok {
				return v
			}
			return s.Default
		}
		return ""
	}
	return defaultStore.Raw(ctx, key)
}

// Bool returns a boolean setting, false if it does not parse.
func Bool(ctx context.Context, key string) bool {
	v, _ := strconv.ParseBool(String(ctx, key))
	return v
}

// Int returns a whole-number setting, 0 if it does not parse.
func Int(ctx context.Context, key string) int {
	n, _ := strconv.Atoi(String(ctx, key))
	return n
}

// Float returns a decimal setting, 0 if it does not parse.
func Float(ctx context.Context, key string) float64 {
	n, _ := strconv.ParseFloat(String(ctx, key), 64)
	return n
}

// Duration returns a duration setting, 0 if it does not parse.
func Duration(ctx context.Context, key string) time.Duration {
	d, _ := time.ParseDuration(String(ctx, key))
	return d
}
`
}
