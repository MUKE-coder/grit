package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

// A fresh project and an upgraded one must end up with the same code.
func TestFreshTemplatesNeedNoGrantCeilingRepair(t *testing.T) {
	for name, c := range map[string]struct {
		src string
		fn  func(string) (string, []string, []string)
	}{
		"handlers/role.go": {roleHandlerGo(), repairRoleCeilingSource},
		// The template that replaced it is checked the same way, since the
		// ceiling has to survive the queries moving onto the service.
		"handlers/role_test.go": {roleHandlerTestGo(), repairRoleTestCallerSource},
		"handlers/user.go":      {apiUserHandlerGo(), repairUserCeilingSource},
	} {
		out, fixed, warn := c.fn(c.src)
		if out != c.src || len(fixed) > 0 || len(warn) > 0 {
			t.Errorf("%s: the scaffold template still needs the repair (fixed %v, warn %v)", name, fixed, warn)
		}
	}
}

// The repair applied to handlers as v3.246.0 generated them: the ceiling taken
// back out of the handler each release wrote before this repair existed, which is
// the fixture, because the blocks the repair writes build their own queries and
// the templates do not any more.
func TestRepairAddsTheGrantCeiling(t *testing.T) {
	roleFresh := roleHandlerFixture(t, "v3.340.0.go.txt")
	role := roleFresh
	for _, s := range []string{roleCeilingCheck, roleSystemCheck, roleAssignCheck} {
		role = strings.ReplaceAll(role, s, "")
	}
	if strings.Contains(role, "authz.CannotGrant(") {
		t.Fatal("could not reconstruct the old role handler")
	}
	out, fixed, warn := repairRoleCeilingSource(role)
	if len(warn) > 0 || len(fixed) != 1 {
		t.Fatalf("role.go: fixed %v, warned %v", fixed, warn)
	}
	if out != roleFresh {
		t.Error("the repaired role handler differs from the one that release wrote")
	}

	// From the fixture, not from the template: the checks this repair writes name
	// &user, and the template holds a pointer since the queries moved onto the
	// service. What the repair meets is the older file, which is what testdata
	// holds.
	userFresh := userHandlerFixture(t, "v3.339.0.go.txt")
	user := userFresh
	for _, s := range []string{
		"\t// Only an ADMIN makes an ADMIN, and nobody hands out a role that grants more\n\t// than they hold.\n" + userRoleCheck + "\n",
		"\t// Only an ADMIN changes an ADMIN account or makes one. Before, a users.edit\n\t// holder could PUT {\"role\":\"ADMIN\"} on themselves, or reset an\n\t// administrator's password or email and sign in as them.\n" + userAdminCheck + userRoleCheck + "\n",
		"\t// Deleting an administrator is an ADMIN's call too.\n" + userAdminCheck + "\n",
	} {
		user = strings.ReplaceAll(user, s, "")
	}
	if strings.Contains(user, "authz.RoleBeyondCaller(") {
		t.Fatal("could not reconstruct the old user handler")
	}
	out, fixed, warn = repairUserCeilingSource(user)
	if len(warn) > 0 || len(fixed) != 1 {
		t.Fatalf("user.go: fixed %v, warned %v", fixed, warn)
	}
	if out != userFresh {
		t.Error("the repaired user handler differs from the one that release wrote")
	}

	test := strings.Replace(roleHandlerTestGo(), roleTestRouterAdmin, roleTestRouterAnchor, 1)
	if out, fixed, _ := repairRoleTestCallerSource(test); len(fixed) != 1 || !strings.Contains(out, `c.Set("user_role", "ADMIN")`) {
		t.Error("the role tests were not made to act as an ADMIN")
	}
}

func TestCeilingTemplatesAreGo(t *testing.T) {
	for name, src := range map[string]string{"ceiling.go": authzCeilingGo(), "ceiling_test.go": authzCeilingTestGo()} {
		if _, err := format.Source([]byte(src)); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}
