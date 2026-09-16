package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// repairGrantCeiling brings a project scaffolded before v3.247.0 up to the fix
// for H2 in the contact-app review: a users.edit or roles.edit holder could make
// themselves ADMIN, reset an administrator's password, or write "*" into their
// own role.
//
// internal/authz/ceiling.go is new framework code and arrives whole. The user
// and role handlers are the developer's, so each check is inserted beside text
// Grit generated, and a handler that no longer has that text is named.
func repairGrantCeiling(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	authz := filepath.Join(apiRoot, "internal", "authz")
	if !fileExists(filepath.Join(authz, "permissions.go")) {
		return nil
	}
	for path, content := range map[string]string{
		filepath.Join(authz, "ceiling.go"):      authzCeilingGo(),
		filepath.Join(authz, "ceiling_test.go"): authzCeilingTestGo(),
	} {
		if err := writeFile(path, strings.ReplaceAll(content, "{{MODULE}}", opts.Module())); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	if !fileContains(filepath.Join(authz, "ceiling.go"), "func CannotGrant(") {
		return nil
	}

	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	handlers := filepath.Join(apiRoot, "internal", "handlers")
	for _, s := range []struct {
		path string
		fn   func(string) (string, []string, []string)
	}{
		{filepath.Join(handlers, "role.go"), repairRoleCeilingSource},
		{filepath.Join(handlers, "role_test.go"), repairRoleTestCallerSource},
		{filepath.Join(handlers, "user.go"), repairUserCeilingSource},
	} {
		if !fileExists(s.path) {
			continue
		}
		if err := repairSourceFile(root, m, s.path, s.fn); err != nil {
			return err
		}
	}
	return nil
}

// insertion puts text beside an anchor that must appear exactly count times.
type insertion struct {
	anchor, text string
	before       bool
	count        int
}

// applyInsertions makes every insertion or none of them.
func applyInsertions(src string, ins []insertion) (string, bool) {
	out := src
	for _, in := range ins {
		n := in.count
		if n == 0 {
			n = 1
		}
		if strings.Count(out, in.anchor) != n {
			return src, false
		}
		replacement := in.anchor + in.text
		if in.before {
			replacement = in.text + in.anchor
		}
		out = strings.ReplaceAll(out, in.anchor, replacement)
	}
	return out, true
}

const (
	roleValidateAnchor = "\tif bad := validateGrants(in.Grants); len(bad) > 0 {\n\t\trespond.BadRequest(c, \"Unknown permissions: \"+strings.Join(bad, \", \"))\n\t\treturn\n\t}\n"
	roleCeilingCheck   = "\t// A role hands out at most what its author holds, or a roles.create or\n\t// roles.edit holder could write \"*\" into a role and take it.\n\tif beyond := authz.CannotGrant(c, in.Grants); len(beyond) > 0 {\n\t\trespond.Forbidden(c, \"You cannot grant permissions you do not hold: \"+strings.Join(beyond, \", \"))\n\t\treturn\n\t}\n"
	roleSystemAnchor   = "\t// System roles keep their name."
	roleSystemCheck    = "\t// Only an ADMIN changes a built-in role: taking \"*\" out of ADMIN locks out\n\t// every administrator.\n\tif role.IsSystem && !authz.IsAdmin(c) {\n\t\trespond.Forbidden(c, \"Only an ADMIN can change a built-in role\")\n\t\treturn\n\t}\n\n"
	roleAssignAnchor   = "\t\tif int(n) != len(in.RoleIDs) {\n\t\t\trespond.BadRequest(c, \"One or more roles do not exist\")\n\t\t\treturn\n\t\t}\n\t}\n"
	roleAssignCheck    = `
	// The ceiling applies to assignment too: only an ADMIN changes an ADMIN's
	// roles or hands out ADMIN, and nobody else hands out a role that grants more
	// than they hold. Before, a users.edit holder could assign themselves ADMIN.
	if !authz.IsAdmin(c) {
		if authz.IsAdminAccount(h.DB, &user) {
			respond.Forbidden(c, "Only an ADMIN can change an ADMIN's roles")
			return
		}
		var roles []models.Role
		if err := h.DB.WithContext(c.Request.Context()).Where("id IN ?", in.RoleIDs).Find(&roles).Error; err != nil {
			respond.Internal(c, err)
			return
		}
		for i := range roles {
			if strings.EqualFold(roles[i].Name, models.RoleAdmin) {
				respond.Forbidden(c, "Only an ADMIN can assign the ADMIN role")
				return
			}
			if beyond := authz.CannotGrant(c, roles[i].GrantsList()); len(beyond) > 0 {
				respond.Forbidden(c, fmt.Sprintf("The %s role grants permissions you do not hold: %s", roles[i].Name, strings.Join(beyond, ", ")))
				return
			}
		}
	}
`
)

func repairRoleCeilingSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "func (h *RoleHandler) AssignUserRoles(") || strings.Contains(src, "authz.CannotGrant(") {
		return src, nil, nil
	}
	out, ok := applyInsertions(src, []insertion{
		{anchor: roleValidateAnchor, text: roleCeilingCheck, count: 2},
		{anchor: roleSystemAnchor, text: roleSystemCheck, before: true},
		{anchor: roleAssignAnchor, text: roleAssignCheck},
	})
	if !ok {
		return src, nil, []string{"is not the role handler Grit wrote, so a roles.edit or users.edit holder can still grant themselves more: check authz.CannotGrant in Create, Update and AssignUserRoles"}
	}
	return out, []string{"roles cannot grant, and assignments cannot hand out, more than the caller holds"}, nil
}

const (
	roleTestRouterAnchor = "\tr := gin.New()\n\tr.GET(\"/permissions\", h.Catalog)\n"
	roleTestRouterAdmin  = "\tr := gin.New()\n\t// These tests act as an ADMIN. The grant ceiling has tests of its own.\n\tr.Use(func(c *gin.Context) { c.Set(\"user_role\", \"ADMIN\"); c.Next() })\n\tr.GET(\"/permissions\", h.Catalog)\n"
)

// repairRoleTestCallerSource makes the shipped role tests act as an ADMIN, which
// they always meant to: with the ceiling, a caller with no role is refused.
func repairRoleTestCallerSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "func roleRouter(") || strings.Contains(src, `c.Set("user_role", "ADMIN")`) {
		return src, nil, nil
	}
	if strings.Count(src, roleTestRouterAnchor) != 1 {
		return src, nil, []string{"has a different test router; its requests need user_role ADMIN now that roles have a grant ceiling"}
	}
	return strings.Replace(src, roleTestRouterAnchor, roleTestRouterAdmin, 1), []string{"the role tests act as an ADMIN"}, nil
}

const (
	userRoleCheck  = "\tif reason := authz.RoleBeyondCaller(c, h.DB, req.Role); reason != \"\" {\n\t\trespond.Fail(c, respond.CodeForbidden, \"You cannot give that role: \"+reason)\n\t\treturn\n\t}\n"
	userAdminCheck = "\tif !authz.IsAdmin(c) && authz.IsAdminAccount(h.DB, &user) {\n\t\trespond.Fail(c, respond.CodeForbidden, \"Only an ADMIN can change an ADMIN account\")\n\t\treturn\n\t}\n"
)

func repairUserCeilingSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "func (h *UserHandler) Update(") || strings.Contains(src, "authz.RoleBeyondCaller(") {
		return src, nil, nil
	}
	// The delete call is bound to the request's context in handlers written since
	// v3.270.0, and plain before it; the anchor is whichever this file has.
	deleteCall := "h.DB.Delete(&user)"
	if strings.Contains(src, "h.DB.WithContext(c.Request.Context()).Delete(&user)") {
		deleteCall = "h.DB.WithContext(c.Request.Context()).Delete(&user)"
	}
	// Create read the email before inserting it until v3.281.0, and builds the
	// row straight away since; the ceiling check goes before whichever of the
	// two this file has.
	createAnchor := "\t// Check email uniqueness\n\tvar existing models.User"
	if !strings.Contains(src, createAnchor) {
		createAnchor = "\tuser := models.User{\n\t\tFirstName: req.FirstName,"
	}
	// And the failure it answers with is respond.Fail in handlers written since
	// the error envelopes moved there, and a hand-built gin.H before it.
	deleteFail := "\tif err := " + deleteCall + ".Error; err != nil {\n\t\trespond.Fail(c, respond.CodeInternalError, \"Failed to delete user\")"
	if !strings.Contains(src, deleteFail) {
		deleteFail = "\tif err := " + deleteCall + ".Error; err != nil {\n\t\tc.JSON(http.StatusInternalServerError, gin.H{\n\t\t\t\"error\": gin.H{\n\t\t\t\t\"code\":    \"INTERNAL_ERROR\",\n\t\t\t\t\"message\": \"Failed to delete user\","
	}
	out, ok := applyInsertions(src, []insertion{
		{anchor: createAnchor, before: true,
			text: "\t// Only an ADMIN makes an ADMIN, and nobody hands out a role that grants more\n\t// than they hold.\n" + userRoleCheck + "\n"},
		{anchor: "\tupdates := map[string]interface{}{}\n\tif req.FirstName != \"\" {", before: true,
			text: "\t// Only an ADMIN changes an ADMIN account or makes one. Before, a users.edit\n\t// holder could PUT {\"role\":\"ADMIN\"} on themselves, or reset an\n\t// administrator's password or email and sign in as them.\n" + userAdminCheck + userRoleCheck + "\n"},
		{anchor: deleteFail, before: true,
			text: "\t// Deleting an administrator is an ADMIN's call too.\n" + userAdminCheck + "\n"},
	})
	if !ok {
		return src, nil, []string{"is not the user handler Grit wrote, so a users.edit holder can still make themselves ADMIN: call authz.RoleBeyondCaller and authz.IsAdminAccount in Create, Update and Delete"}
	}
	return out, []string{"only an ADMIN makes or changes an ADMIN account, and roles handed out stay within the caller's grants"}, nil
}
