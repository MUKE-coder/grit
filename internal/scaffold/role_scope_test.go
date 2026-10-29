package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Every lookup of a role BY NAME has to say whose role it means.
//
// A role's name is unique per owner now, not across the platform, so "the
// role called ADMIN" is a question with as many answers as there are
// organizations. A query that does not narrow to an owner can pick up a
// customer's own role as if the platform had defined it, which is how a
// membership ends up holding permissions nobody granted (grit#130, and the
// shape GHSA-4j45-98gm-mv9v is about).
//
// The report listed four such places. There were six. That gap is why this
// is a test and not a checklist: the next one will be added by somebody who
// has never read the issue.
var roleNameLookup = regexp.MustCompile(
	`(?m)^.*\bWhere\("name = \?"[^)]*\).*$`)

// A line that already narrows, is not about roles, or carries the marker.
//
// grit:role-lookup-ok is the escape hatch, and it has to carry a reason:
// one place needs it, the repair that SEARCHES for the old unnarrowed query
// in a project's source. Adding the marker to silence a real lookup would be
// visible in review, which is the point of making it a comment rather than a
// list in this file.
var roleLookupScoped = regexp.MustCompile(
	`PlatformRoles\(|OrgRoles\(|owner_org_id|grit:role-lookup-ok`)

// Files whose "name = ?" is about something else entirely.
var notRoles = []string{
	"cluster", // Generation rows
	"variant", // product option names
	"option",
}

func TestEveryRoleLookupByNameNarrowsToAnOwner(t *testing.T) {
	roots := []string{".", filepath.Join("templates", "api")}

	checked, flagged := 0, 0
	for _, root := range roots {
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return err
			}
			name := info.Name()
			if !strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, ".tmpl") {
				return nil
			}
			if strings.HasSuffix(name, "_test.go") {
				return nil
			}
			for _, skip := range notRoles {
				if strings.Contains(strings.ToLower(name), skip) {
					return nil
				}
			}

			raw, readErr := os.ReadFile(path)
			if readErr != nil {
				return nil
			}
			src := string(raw)

			for _, line := range roleNameLookup.FindAllString(src, -1) {
				// Only the ones querying a Role.
				if !strings.Contains(src, "models.Role") && !strings.Contains(src, "Role{") &&
					!strings.Contains(line, "Role") && !strings.Contains(line, "role") {
					continue
				}
				if !mentionsARole(src, line) {
					continue
				}
				checked++
				if roleLookupScoped.MatchString(line) {
					continue
				}
				flagged++
				t.Errorf("%s: a role is looked up by name without narrowing to an "+
					"owner, so an organization's own role of that name can answer "+
					"(grit#130):\n    %s", path, strings.TrimSpace(line))
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walking %s: %v", root, err)
		}
	}

	if checked == 0 {
		t.Error("no role lookup by name was found anywhere, so this test is " +
			"checking nothing: the reader has stopped matching")
	}
	t.Logf("checked %d role lookups by name, %d unnarrowed", checked, flagged)
}

// mentionsARole keeps the reader off queries that happen to filter a name
// column on some other model.
func mentionsARole(src, line string) bool {
	if strings.Contains(line, "RoleAdmin") || strings.Contains(line, "roleName") ||
		strings.Contains(line, "user.Role") || strings.Contains(line, "want.Name") {
		return true
	}
	// A line inside a file that is plainly about roles.
	return strings.Contains(src, "func SeedRoles(") ||
		strings.Contains(src, "type RoleService ") ||
		strings.Contains(src, "RoleBeyondCaller(")
}
