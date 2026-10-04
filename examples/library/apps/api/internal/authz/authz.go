// Package authz contains the ownership-check helpers Grit uses to
// prevent IDOR (Insecure Direct Object Reference) — OWASP Top 10:2025
// A01 Broken Access Control's most common concrete form.
//
// The cardinal rule (from PHASE 2 §4.3 of the security course): every
// object access must be authorised against the current user, server-side.
// authz.MustOwn enforces that with a single call.
//
// Usage:
//
//	var invoice models.Invoice
//	if err := authz.MustOwn(c, db, &invoice, c.Param("id")); err != nil {
//	    return // helper has already written 404 / 401
//	}
//	// invoice belongs to c.MustGet("user_id"). Safe to use.
package authz

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"library/apps/api/internal/models"
	"library/apps/api/internal/respond"
)

// Ownable is implemented by domain models whose ownership is identified
// by a single user-id column. For models with team/tenant scoping use
// CheckScope instead.
type Ownable interface {
	GetOwnerID() string
}

// ErrNotFound and ErrForbidden are returned by the helpers below so
// callers can branch (e.g. log differently) — but the HTTP responses
// they produce are deliberately identical (404 NOT_FOUND) to avoid
// leaking the existence of rows the caller doesn't own.
var (
	ErrNotFound  = errors.New("authz: not found")
	ErrForbidden = errors.New("authz: forbidden")
)

// MustOwn loads the row by id and verifies that the authenticated user
// is its owner. On any failure it writes a 404 response and returns a
// non-nil error so the caller can return immediately.
//
// The 404 (not 403) is intentional. Returning 403 confirms the row
// exists, which lets attackers enumerate IDs.
func MustOwn(c *gin.Context, db *gorm.DB, dest Ownable, id string) error {
	userID, ok := c.Get("user_id")
	if !ok {
		respond.Fail(c, respond.CodeUnauthorized, "Authentication required")
		return ErrForbidden
	}

	if err := db.Where("id = ?", id).First(dest).Error; err != nil {
		writeNotFound(c)
		return ErrNotFound
	}

	if dest.GetOwnerID() != userID {
		writeNotFound(c) // 404, not 403 — see comment above
		return ErrForbidden
	}
	return nil
}

// IsAdmin reports whether the caller holds the ADMIN role.
//
// Ownership scoping has to let an admin through or the admin panel, which
// calls the same endpoints, shows an administrator only the rows they happen
// to have created themselves.
func IsAdmin(c *gin.Context) bool {
	role, _ := c.Get("user_role")
	return asString(role) == models.RoleAdmin
}

// MustOwnUnlessAdmin is MustOwn with that exemption. It is what a generated
// handler calls; MustOwn stays strict for code that wants no exemption at all.
func MustOwnUnlessAdmin(c *gin.Context, db *gorm.DB, dest Ownable, id string) error {
	if IsAdmin(c) {
		if err := db.Where("id = ?", id).First(dest).Error; err != nil {
			writeNotFound(c)
			return ErrNotFound
		}
		return nil
	}
	return MustOwn(c, db, dest, id)
}

// OwnsOr404 checks a row the handler has already loaded, writing 404 and
// returning false when the caller does not own it. ADMIN always passes.
//
// Same 404-not-403 reasoning as MustOwn: 403 confirms the row exists, which
// turns a guessed id into an enumeration oracle. This variant exists so a
// handler that already fetched the row does not fetch it a second time.
func OwnsOr404(c *gin.Context, row Ownable) bool {
	if IsAdmin(c) {
		return true
	}
	if row.GetOwnerID() != CurrentUserID(c) {
		writeNotFound(c)
		return false
	}
	return true
}

// ScopeToOwner narrows a list query to rows the caller owns. An ADMIN gets the
// query back untouched.
//
// This is the half of IDOR prevention that MustOwn cannot do. MustOwn protects
// a row fetched by id; without this, the list endpoint hands over every row in
// the table and the id is no longer a secret worth guessing.
//
// column is a fixed string from the generator, never user input.
func ScopeToOwner(c *gin.Context, q *gorm.DB, column string) *gorm.DB {
	if IsAdmin(c) {
		return q
	}
	userID, ok := c.Get("user_id")
	if !ok {
		// No authenticated user: match nothing rather than everything. A
		// scoping helper that opens up when it cannot identify the caller is
		// worse than no scoping, because it reads as if it is protecting you.
		return q.Where("1 = 0")
	}
	return q.Where(column+" = ?", userID)
}

// CurrentUserID returns the authenticated user's id, or "".
func CurrentUserID(c *gin.Context) string {
	userID, _ := c.Get("user_id")
	return asString(userID)
}

// CheckScope verifies a (column, value) pair matches the current user's
// authoritative scope (e.g. team_id, tenant_id). Use this when ownership
// is by membership rather than a single user_id column.
func CheckScope(c *gin.Context, scopeKey, expectedValue string) bool {
	got, ok := c.Get(scopeKey)
	return ok && got == expectedValue
}

// RequireRoles returns a gin middleware that 403s when the authenticated
// user's role isn't in the allowlist. This is a stricter sibling of the
// generic Auth middleware — use it on admin-only routes.
func RequireRoles(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		role, _ := c.Get("user_role")
		if _, ok := allowed[asString(role)]; !ok {
			respond.Fail(c, respond.CodeForbidden, "Insufficient role")
			return
		}
		c.Next()
	}
}

func asString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func writeNotFound(c *gin.Context) {
	respond.Fail(c, respond.CodeNotFound, "Resource not found")
}
