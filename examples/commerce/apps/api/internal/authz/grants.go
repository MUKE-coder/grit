package authz

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"gorm.io/gorm"

	"commerce/apps/api/internal/cluster"
	"commerce/apps/api/internal/models"
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
//
// With Share it also tells every other replica. The cache used to be per
// process, so a permission revoked through one API instance kept working on the
// others until they restarted.
func Invalidate() {
	dropLocal()
	if sharedDB != nil {
		if err := cluster.Bump(sharedDB, "authz"); err != nil {
			log.Printf("authz: could not tell the other replicas about a role change: %v", err)
		}
	}
}

func dropLocal() {
	generation.Add(1)
	grantCache.Range(func(k, _ any) bool {
		grantCache.Delete(k)
		return true
	})
}

var (
	sharedDB    *gorm.DB
	sharedWatch *cluster.Watch
)

// Share makes Invalidate reach every replica through the database they share,
// and makes each replica drop its cache within a second of another's change.
// Call it once at start-up, before serving.
func Share(db *gorm.DB) {
	sharedDB = db
	sharedWatch = cluster.NewWatch(db, "authz", time.Second)
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

	// Another replica changed a role: this one's cache is stale too.
	if sharedWatch != nil && sharedWatch.Changed() {
		dropLocal()
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

// GrantsForRole returns the grants of one role.
//
// For a per-organization role: the multitenant plugin's middleware reads the
// membership's role and adds these to the caller's grants for the duration of the
// request, so a role can mean "an administrator of this organization" without
// meaning it everywhere. Cached beside GrantsFor and invalidated by the same
// generation counter, so editing the role takes effect at once.
func GrantsForRole(db *gorm.DB, roleID string) ([]string, error) {
	if roleID == "" {
		return nil, nil
	}

	if sharedWatch != nil && sharedWatch.Changed() {
		dropLocal()
	}

	key := "role:" + roleID
	gen := generation.Load()
	if v, ok := grantCache.Load(key); ok {
		if c, ok := v.(cachedGrants); ok && c.gen == gen {
			return c.grants, nil
		}
	}

	var role models.Role
	if err := db.Where("id = ?", roleID).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	grants := role.GrantsList()
	grantCache.Store(key, cachedGrants{gen: gen, grants: grants})
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
			if errors.Is(err, gorm.ErrRecordNotFound) {
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
