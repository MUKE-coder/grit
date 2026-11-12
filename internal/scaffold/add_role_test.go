package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// `grit add role` has to reach the admin, not just the Go and TypeScript
// halves.
//
// It looked for apps/admin/resources/users.ts. Resources moved into a folder
// each, so the real path is apps/admin/resources/users/users.ts, fileExists
// answered no, and the three admin injections were skipped in silence. The
// command printed "added to 4 location(s)" and looked like it had worked.
//
// The effect is the one that matters: the role reached the model, the schemas,
// the types and the constants, and in the admin it had no badge, no filter and
// no option in the form's select. So the new role could not be given to
// anybody from the screen that exists to give it to somebody.
//
// The count is asserted as well as the file, because a cheerful success
// message was the whole reason this went unnoticed.
func TestAddRoleReachesTheAdminResource(t *testing.T) {
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })

	if err := Run(Options{
		ProjectName:  "rolecheck",
		Architecture: ArchTriple,
		Frontend:     "next",
		Version:      DefaultVersion,
	}); err != nil {
		t.Fatalf("scaffolding: %v", err)
	}
	root := filepath.Join(dir, "rolecheck")
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}

	// The definition is in a folder of its own, which is the layout the old
	// path could not see. If this ever moves again, the test says so here
	// rather than through a role that quietly does nothing.
	def := FindResourceDef(filepath.Join(root, "apps", "admin", "resources"), "users")
	if def == "" {
		t.Fatal("no users resource definition found, so this test cannot check what it is for")
	}

	if err := AddRole("MANAGER"); err != nil {
		t.Fatalf("adding a role: %v", err)
	}

	raw, err := os.ReadFile(def)
	if err != nil {
		t.Fatal(err)
	}
	body := string(raw)

	for _, want := range []string{
		`MANAGER: { color: "success", label: "Manager" },`, // the table badge
		`{ label: "Manager", value: "MANAGER" },`,          // the filter and the form option
	} {
		if !strings.Contains(body, want) {
			t.Errorf("the admin users resource is missing %s", want)
		}
	}
	// Twice: once in the filter options, once in the form's select. A role
	// you can filter by and cannot assign is half a role.
	if n := strings.Count(body, `{ label: "Manager", value: "MANAGER" },`); n != 2 {
		t.Errorf("the role appears %d times in the admin resource, want 2 (filter and form option)", n)
	}

	// And the four places that always worked, so a regression in the fix
	// cannot trade one half for the other.
	for _, f := range []struct{ path, want string }{
		{filepath.Join(root, "apps", "api", "internal", "models", "user.go"), `RoleManager = "MANAGER"`},
		{filepath.Join(root, "packages", "shared", "schemas", "user.ts"), `"MANAGER"`},
		{filepath.Join(root, "packages", "shared", "types", "user.ts"), `"MANAGER"`},
		{filepath.Join(root, "packages", "shared", "constants", "index.ts"), `MANAGER`},
	} {
		raw, err := os.ReadFile(f.path)
		if err != nil {
			t.Errorf("reading %s: %v", filepath.Base(f.path), err)
			continue
		}
		if !strings.Contains(string(raw), f.want) {
			t.Errorf("%s is missing %s", filepath.Base(f.path), f.want)
		}
	}
}

// Adding a role and removing it again leaves the project exactly as it was.
//
// That is the property worth asserting, because the failure mode of a partial
// removal is quiet: a value left in the Zod enum that the Go model no longer
// accepts means the admin form offers a role the API refuses, and nothing
// about that looks wrong until somebody picks it.
func TestAddingAndRemovingARoleIsARoundTrip(t *testing.T) {
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })

	if err := Run(Options{
		ProjectName:  "roundtrip",
		Architecture: ArchTriple,
		Frontend:     "next",
		Version:      DefaultVersion,
	}); err != nil {
		t.Fatalf("scaffolding: %v", err)
	}
	root := filepath.Join(dir, "roundtrip")
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}

	paths := roleFilesFor(root)
	files := map[string]string{
		"model":     paths.Model,
		"schema":    paths.Schema,
		"types":     paths.Types,
		"constants": paths.Constants,
		"admin":     paths.Admin,
	}
	before := map[string]string{}
	for name, path := range files {
		if path == "" {
			t.Fatalf("a triple project has no %s file for roles", name)
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		before[name] = string(raw)
	}

	if err := AddRole("MANAGER"); err != nil {
		t.Fatalf("adding: %v", err)
	}
	// It really did change them, or the round trip proves nothing.
	for name, path := range files {
		raw, _ := os.ReadFile(path)
		if string(raw) == before[name] {
			t.Errorf("adding the role did not change %s", name)
		}
	}

	if err := RemoveRole("MANAGER", false); err != nil {
		t.Fatalf("removing: %v", err)
	}
	for name, path := range files {
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if string(raw) != before[name] {
			t.Errorf("%s did not come back to what it was:\n--- before ---\n%s\n--- after ---\n%s",
				name, before[name], string(raw))
		}
	}
}

// The two the framework itself depends on cannot be removed by typing one
// command. A project without them compiles and cannot sign anybody in.
func TestTheFrameworksOwnRolesAreProtected(t *testing.T) {
	prev, _ := os.Getwd()
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })

	if err := Run(Options{
		ProjectName:  "protected",
		Architecture: ArchTriple,
		Frontend:     "next",
		Version:      DefaultVersion,
	}); err != nil {
		t.Fatalf("scaffolding: %v", err)
	}
	if err := os.Chdir(filepath.Join(dir, "protected")); err != nil {
		t.Fatal(err)
	}

	for _, role := range []string{"ADMIN", "USER", "admin"} {
		err := RemoveRole(role, true) // even with --force
		if err == nil {
			t.Errorf("%s was removable", role)
			continue
		}
		if !strings.Contains(err.Error(), "cannot be removed") {
			t.Errorf("removing %s failed with %q, which does not say why", role, err)
		}
	}

	// And a role the project does not have says so rather than reporting a
	// successful no-op.
	if err := RemoveRole("NOBODY", false); err == nil {
		t.Error("removing a role the project does not have returned no error")
	}
}
