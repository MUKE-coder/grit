package scaffold

import "strings"

// roleModelGo emits internal/models/role.go — the Role model, the user↔role
// join, and the default-role seeding.
//
// Shape decisions:
//
//   - Grants are a JSON array in a text column, not a normalised
//     role_permissions table. One row, one query, no joins — and because the
//     stored strings keep their wildcards, a role granted "products.*"
//     automatically covers actions added to the catalog later.
//
//   - user_roles is many-to-many even though most apps assign one role. It costs
//     nothing now, and the multi-tenant plugin later adds org_id to this same
//     table rather than migrating off a single users.role_id column.
//
//   - The join table deliberately has NO DeletedAt. Soft-delete tombstones
//     collide with the composite primary key when a role is re-granted, which
//     surfaces as a unique-constraint violation on an operation the user
//     reasonably expects to be idempotent.
func roleModelGo() string {
	src := `package models

import (
	"encoding/json"
	"time"

	"{{MODULE}}/internal/ids"
	"gorm.io/gorm"
)

// Role is a named bag of permission grants.
//
// Grants hold whatever the operator authored, wildcards included — see
// internal/authz for the key format and matcher. Storing "products.*" rather
// than the expanded leaves is deliberate: the role keeps working when a new
// action shows up in the catalog.
type Role struct {
	ID          string ~gorm:"primarykey;size:36" json:"id"~
	Name        string ~gorm:"size:80;uniqueIndex;not null" json:"name" binding:"required"~
	Description string ~gorm:"size:500" json:"description"~

	// Grants is a JSON array of permission keys. Use GrantsList/SetGrants
	// rather than touching it directly.
	Grants string ~gorm:"type:text" json:"-"~

	// IsSystem marks a role the app seeds and depends on. System roles can be
	// edited but not deleted or renamed — enforced in the handler, not just the
	// UI, so a direct API call can't rename ADMIN out from under the guards.
	IsSystem bool ~gorm:"default:false;index" json:"is_system"~

	Version   int            ~gorm:"not null;default:1" json:"version"~
	CreatedAt time.Time      ~json:"created_at"~
	UpdatedAt time.Time      ~json:"updated_at"~
	DeletedAt gorm.DeletedAt ~gorm:"index" json:"-"~
}

func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = ids.New()
	}
	return nil
}

func (r *Role) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}

// GrantsList decodes the stored grants.
//
// A malformed value yields NO permissions rather than an error: this is an
// authorization path, and the safe failure mode is to deny.
func (r *Role) GrantsList() []string {
	if r.Grants == "" {
		return nil
	}
	var out []string
	if err := json.Unmarshal([]byte(r.Grants), &out); err != nil {
		return nil
	}
	return out
}

func (r *Role) SetGrants(grants []string) error {
	if grants == nil {
		grants = []string{}
	}
	b, err := json.Marshal(grants)
	if err != nil {
		return err
	}
	r.Grants = string(b)
	return nil
}

// MarshalJSON exposes grants as a real array to API clients while keeping the
// column a string internally.
func (r Role) MarshalJSON() ([]byte, error) {
	type alias Role
	return json.Marshal(struct {
		alias
		Grants []string ~json:"grants"~
	}{alias(r), r.GrantsList()})
}

// UserRole assigns a Role to a User.
//
// No DeletedAt on purpose — see the note on roleModelGo. Re-granting a
// previously removed role must be idempotent, and a tombstone row makes the
// composite PK reject it.
type UserRole struct {
	UserID    string    ~gorm:"primaryKey;size:36;index" json:"user_id"~
	RoleID    string    ~gorm:"primaryKey;size:36;index" json:"role_id"~
	CreatedAt time.Time ~json:"created_at"~
}

func (UserRole) TableName() string { return "user_roles" }

// DefaultRoles are seeded on first boot.
//
// The grants mirror what the routes enforced BEFORE permissions existed, so
// upgrading an app doesn't change who can do what:
//
//	ADMIN  — everything, via the "*" superuser grant
//	EDITOR — content, matching the role's existing meaning
//	USER   — nothing privileged; ordinary users hold no admin permissions
//
// Note this is NOT "every role gets every permission". The route guard passes if
// the caller matches the legacy role name OR holds the permission, so granting
// USER the full catalog would hand every signed-up user the admin panel.
func DefaultRoles() []Role {
	return []Role{
		{
			Name:        RoleAdmin,
			Description: "Full access to every feature and action.",
			IsSystem:    true,
			Grants:      ~["*"]~,
		},
		{
			Name:        RoleEditor,
			Description: "Manages content and uploads. No user or system administration.",
			IsSystem:    true,
			Grants:      ~["uploads.create","uploads.view","uploads.delete","users.view"]~,
		},
		{
			Name:        RoleUser,
			Description: "Standard account. No administrative permissions.",
			IsSystem:    true,
			Grants:      ~[]~,
		},
	}
}

// SeedRoles inserts the default roles if they are missing.
//
// Idempotent, and deliberately NON-destructive: an existing role's grants are
// left alone. Re-applying defaults on every boot would silently revert an
// operator's edits — a real bug in the system this was modelled on, where
// seeding ran on every roles-page load and quietly undid their changes.
// "Reset to defaults" is an explicit action in the UI instead.
func SeedRoles(db *gorm.DB) error {
	for _, want := range DefaultRoles() {
		var existing Role
		// Find, not First: on a fresh database every role is missing, and First
		// logs each miss as "record not found" — three scary-looking lines on
		// the first migrate even though creating them is exactly the point.
		// Find treats zero rows as normal and stays quiet.
		if err := db.Where("name = ?", want.Name).Limit(1).Find(&existing).Error; err != nil {
			return err
		}
		if existing.ID != "" {
			continue // already present — leave the operator's grants intact
		}
		role := want
		if err := db.Create(&role).Error; err != nil {
			return err
		}
	}
	return nil
}
`
	return strings.ReplaceAll(src, "~", "`")
}

// authzGrantsGo emits internal/authz/grants.go — resolving a user to their
// permission grants.
//
// This is the ONE place that answers "what can this user do". Everything else
// (middleware, handlers, the /me endpoint) goes through GrantsFor, which is what
// makes the multi-tenant plugin a contained change later: it swaps the
// resolution for a per-organization lookup without touching call sites.
//
// Resolution order:
//  1. explicit user_roles assignments (the real model)
//  2. failing that, the legacy users.role string — so apps upgrading from
//     role-only authorization keep working before anyone is assigned a Role row
func authzGrantsGo() string {
	src := `package authz

import (
	"sync"
	"sync/atomic"

	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// grantCache memoises user -> grants for the lifetime of a generation.
//
// Authorization runs on every request, so hitting the database twice per call
// is a real cost. Rather than expire by time (which makes a revoked permission
// linger for up to the TTL), the whole cache is dropped whenever roles change —
// see Invalidate. Revocation is therefore immediate.
var (
	grantCache sync.Map // userID -> []string
	generation atomic.Uint64
)

// Invalidate drops every cached grant set. Call it after any write that could
// change authorization: role grants edited, role deleted, user's roles changed.
func Invalidate() {
	generation.Add(1)
	grantCache.Range(func(k, _ any) bool {
		grantCache.Delete(k)
		return true
	})
}

type cachedGrants struct {
	gen    uint64
	grants []string
}

// GrantsFor returns every permission grant the user holds, unioned across their
// roles. Wildcards are preserved — pass the result to Granted, which understands
// them; do not compare strings directly.
func GrantsFor(db *gorm.DB, userID string) ([]string, error) {
	if userID == "" {
		return nil, nil
	}

	gen := generation.Load()
	if v, ok := grantCache.Load(userID); ok {
		if c, ok := v.(cachedGrants); ok && c.gen == gen {
			return c.grants, nil
		}
	}

	grants, err := resolveGrants(db, userID)
	if err != nil {
		return nil, err
	}
	grantCache.Store(userID, cachedGrants{gen: gen, grants: grants})
	return grants, nil
}

func resolveGrants(db *gorm.DB, userID string) ([]string, error) {
	var roles []models.Role
	err := db.
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}

	// Backwards compatibility: an app upgrading from role-string authorization
	// has no user_roles rows yet. Fall back to the role named on the user so
	// existing admins don't lose access the moment permissions ship.
	if len(roles) == 0 {
		var user models.User
		if err := db.Select("id", "role").Where("id = ?", userID).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, nil
			}
			return nil, err
		}
		if user.Role == "" {
			return nil, nil
		}
		if err := db.Where("name = ?", user.Role).Find(&roles).Error; err != nil {
			return nil, err
		}
	}

	seen := map[string]bool{}
	var out []string
	for _, r := range roles {
		for _, g := range r.GrantsList() {
			if seen[g] {
				continue
			}
			seen[g] = true
			out = append(out, g)
		}
	}
	return out, nil
}

// Can reports whether the user holds the permission.
// Prefer the middleware guard for routes; use this for conditional logic inside
// a handler (e.g. hiding fields).
func Can(db *gorm.DB, userID, permission string) bool {
	grants, err := GrantsFor(db, userID)
	if err != nil {
		return false // fail closed
	}
	return Granted(grants, permission)
}
`
	return strings.ReplaceAll(src, "~", "`")
}

// authzGrantsTestGo emits internal/authz/grants_test.go.
//
// These ship with the generated app on purpose. The permission catalog is
// application data that teams edit, and grant resolution is a security path —
// these pin the behaviour that matters:
//   - default roles seed automatically, and re-seeding never clobbers an
//     operator's edits (the failure mode in the system this was modelled on,
//     where seeding ran on every roles-page load and silently reverted changes)
//   - the legacy users.role string still authorises, so upgrading from
//     role-only auth doesn't lock existing admins out
//   - an explicit user_roles assignment wins, and doesn't leak other grants
//   - revocation takes effect immediately rather than after a cache TTL
func authzGrantsTestGo() string {
	return `package authz_test

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"{{MODULE}}/internal/authz"
	"{{MODULE}}/internal/models"
)

func newDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Role{}, &models.UserRole{}); err != nil {
		t.Fatal(err)
	}
	if err := models.SeedRoles(db); err != nil {
		t.Fatal(err)
	}
	return db
}

// Roles must exist automatically, with sane grants.
func TestSeedRoles(t *testing.T) {
	db := newDB(t)

	var roles []models.Role
	db.Find(&roles)
	if len(roles) != 3 {
		t.Fatalf("seeded %d roles, want 3", len(roles))
	}

	var admin models.Role
	db.Where("name = ?", models.RoleAdmin).First(&admin)
	if got := admin.GrantsList(); len(got) != 1 || got[0] != "*" {
		t.Errorf("ADMIN grants = %v, want [*]", got)
	}
	if !admin.IsSystem {
		t.Error("ADMIN must be a system role")
	}

	// Re-seeding must not duplicate or clobber.
	admin.SetGrants([]string{"users.view"})
	db.Save(&admin)
	if err := models.SeedRoles(db); err != nil {
		t.Fatal(err)
	}
	db.Where("name = ?", models.RoleAdmin).First(&admin)
	if got := admin.GrantsList(); len(got) != 1 || got[0] != "users.view" {
		t.Errorf("re-seeding clobbered an operator edit: %v", got)
	}
}

// The legacy users.role string must still authorise, so apps upgrading from
// role-only auth don't lose access before anyone is assigned a Role row.
func TestGrantsFor_LegacyRoleFallback(t *testing.T) {
	db := newDB(t)
	admin := models.User{ID: "u-admin", Email: "a@x.com", FirstName: "A", LastName: "B", Role: models.RoleAdmin}
	plain := models.User{ID: "u-plain", Email: "p@x.com", FirstName: "P", LastName: "B", Role: models.RoleUser}
	db.Create(&admin)
	db.Create(&plain)

	g, err := authz.GrantsFor(db, admin.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !authz.Granted(g, "users.delete") {
		t.Errorf("legacy ADMIN lost access; grants=%v", g)
	}

	g, err = authz.GrantsFor(db, plain.ID)
	if err != nil {
		t.Fatal(err)
	}
	if authz.Granted(g, "users.delete") {
		t.Errorf("plain USER must not hold users.delete; grants=%v", g)
	}
}

// An explicit user_roles assignment takes precedence over the legacy string.
func TestGrantsFor_ExplicitAssignmentWins(t *testing.T) {
	db := newDB(t)
	u := models.User{ID: "u1", Email: "u1@x.com", FirstName: "U", LastName: "1", Role: models.RoleUser}
	db.Create(&u)

	custom := models.Role{Name: "Support"}
	custom.SetGrants([]string{"uploads.*"})
	db.Create(&custom)
	db.Create(&models.UserRole{UserID: u.ID, RoleID: custom.ID})
	authz.Invalidate()

	g, _ := authz.GrantsFor(db, u.ID)
	if !authz.Granted(g, "uploads.delete") {
		t.Errorf("uploads.* should cover uploads.delete; grants=%v", g)
	}
	if authz.Granted(g, "users.delete") {
		t.Errorf("assignment must not leak other permissions; grants=%v", g)
	}
}

// Revoking must take effect immediately, not after a cache TTL.
func TestInvalidate(t *testing.T) {
	db := newDB(t)
	u := models.User{ID: "u2", Email: "u2@x.com", FirstName: "U", LastName: "2", Role: models.RoleAdmin}
	db.Create(&u)

	if g, _ := authz.GrantsFor(db, u.ID); !authz.Granted(g, "users.delete") {
		t.Fatal("expected initial admin access")
	}

	db.Model(&models.User{}).Where("id = ?", u.ID).Update("role", models.RoleUser)
	authz.Invalidate()

	if g, _ := authz.GrantsFor(db, u.ID); authz.Granted(g, "users.delete") {
		t.Error("demotion did not take effect after Invalidate")
	}
}
`
}
