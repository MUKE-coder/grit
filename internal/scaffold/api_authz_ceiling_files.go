package scaffold

// authzCeilingGo emits internal/authz/ceiling.go: the grant ceiling (H2 in the
// contact-app review). Nobody but an ADMIN hands out more than they hold, and
// only an ADMIN changes an ADMIN account.
func authzCeilingGo() string {
	return `package authz

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// The grant ceiling.
//
// A role or user editor could hand out more than they held. A users.edit holder
// could PUT {"role":"ADMIN"} on themselves, assign themselves the ADMIN role, or
// reset an administrator's password and sign in as them, and a roles.edit holder
// could add "*" to their own role. An ADMIN hands out anything; everyone else
// hands out at most what they hold, and only an ADMIN changes an ADMIN account.

// CallerGrants are the grants the auth middleware resolved for this request.
func CallerGrants(c *gin.Context) []string {
	if v, ok := c.Get("user_grants"); ok {
		if grants, ok := v.([]string); ok {
			return grants
		}
	}
	return nil
}

// Beyond lists the requested grants that held does not cover. A wildcard is
// covered only when every permission it expands to is.
func Beyond(held, requested []string) []string {
	if HasAll(held) {
		return nil
	}
	var out []string
	for _, g := range requested {
		switch {
		case g == "*":
			out = append(out, g)
		case strings.Contains(g, "*"):
			keys := Expand([]string{g})
			covered := len(keys) > 0
			for _, k := range keys {
				if !Granted(held, k) {
					covered = false
					break
				}
			}
			if !covered {
				out = append(out, g)
			}
		default:
			if !Granted(held, g) {
				out = append(out, g)
			}
		}
	}
	return out
}

// CannotGrant lists what the caller may not hand out among requested. An ADMIN
// may hand out anything.
func CannotGrant(c *gin.Context, requested []string) []string {
	if IsAdmin(c) {
		return nil
	}
	return Beyond(CallerGrants(c), requested)
}

// IsAdminAccount reports whether user is an administrator: the ADMIN role, or
// roles that together grant everything. When the grants cannot be read, the
// account is treated as one, which is the strict answer.
func IsAdminAccount(db *gorm.DB, user *models.User) bool {
	if strings.EqualFold(user.Role, models.RoleAdmin) {
		return true
	}
	grants, err := GrantsFor(db, user.ID)
	if err != nil {
		return true
	}
	return HasAll(grants)
}

// RoleBeyondCaller says why the caller may not give someone the role named
// roleName, or returns "" when they may. A name with no role row grants nothing
// through the roles table, so only ADMIN is refused among those.
func RoleBeyondCaller(c *gin.Context, db *gorm.DB, roleName string) string {
	if roleName == "" || IsAdmin(c) {
		return ""
	}
	if strings.EqualFold(roleName, models.RoleAdmin) {
		return "only an ADMIN can make an ADMIN"
	}
	var role models.Role
	if err := db.Where("name = ?", roleName).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ""
		}
		return "the role could not be checked"
	}
	if beyond := CannotGrant(c, role.GrantsList()); len(beyond) > 0 {
		return "the " + role.Name + " role grants permissions you do not hold: " + strings.Join(beyond, ", ")
	}
	return ""
}
`
}

func authzCeilingTestGo() string {
	return `package authz

import (
	"reflect"
	"testing"
)

// The ceiling on handing out permissions. Before it, a users.edit or roles.edit
// holder could make themselves ADMIN.
func TestBeyondIsTheGrantCeiling(t *testing.T) {
	cases := []struct {
		name            string
		held, requested []string
		want            []string
	}{
		{"a holder of everything may grant anything", []string{"*"}, []string{"*", "users.edit"}, nil},
		{"nobody else may grant everything", []string{"users.edit"}, []string{"*"}, []string{"*"}},
		{"a held permission may be granted", []string{"users.view", "users.edit"}, []string{"users.edit"}, nil},
		{"one not held may not", []string{"users.view"}, []string{"users.delete"}, []string{"users.delete"}},
		{"a wildcard needs everything it expands to", []string{"users.view"}, []string{"users.*"}, []string{"users.*"}},
		{"a held wildcard covers its permissions", []string{"users.*"}, []string{"users.edit", "users.*"}, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Beyond(tc.held, tc.requested); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Beyond(%v, %v) = %v, want %v", tc.held, tc.requested, got, tc.want)
			}
		})
	}
}
`
}
