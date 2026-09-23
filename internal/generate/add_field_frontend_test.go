package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// An --api project has no packages/shared and no apps/admin. grit generate field
// wrote to both without looking, so on an API-only project it failed on the first
// one: the model had already gained the field, the command reported an error, and
// what it had done was left half applied.
func TestAddFieldSkipsAFrontendThatIsNotThere(t *testing.T) {
	root := t.TempDir()
	def, err := ParseInlineFields("Widget", "name:string,colour:text")
	if err != nil {
		t.Fatal(err)
	}
	g := newTestGenerator(root, "app/apps/api", def)
	names := BuildNames(def)
	field := def.Fields[len(def.Fields)-1]

	if err := g.injectFrontendField(names, field); err != nil {
		t.Fatalf("an API-only project should need no frontend injection: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "packages")); !os.IsNotExist(err) {
		t.Error("skipping the frontend should not create packages/")
	}
	if _, err := os.Stat(filepath.Join(root, "apps", "admin")); !os.IsNotExist(err) {
		t.Error("skipping the frontend should not create apps/admin/")
	}
}

// And the guard is a check, not a blanket skip: with the directory there, the
// injection is attempted and its failure is reported rather than swallowed.
func TestAddFieldStillReportsAFrontendFailure(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "packages", "shared", "schemas"), 0o755); err != nil {
		t.Fatal(err)
	}

	def, err := ParseInlineFields("Widget", "name:string,colour:text")
	if err != nil {
		t.Fatal(err)
	}
	g := newTestGenerator(root, "app/apps/api", def)

	// The directory exists but widget.ts does not, which is a real problem in a
	// monorepo: the resource's schema file should be there.
	if err := g.injectFrontendField(BuildNames(def), def.Fields[len(def.Fields)-1]); err == nil {
		t.Error("a missing schema file in a monorepo should be reported, not skipped")
	}
}

// The built-in User does not follow the generated shape: its create-shaped
// schema is RegisterSchema, so adding a field to User wrote the Go column and
// then reported that CreateUserSchema was not there, leaving the model and the
// shared schema disagreeing.
func TestAddFieldOnUserFindsTheRegisterSchema(t *testing.T) {
	root := t.TempDir()
	schemas := filepath.Join(root, "packages", "shared", "schemas")
	types := filepath.Join(root, "packages", "shared", "types")
	for _, dir := range []string{schemas, types} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	// The shapes the scaffold actually writes for User.
	if err := os.WriteFile(filepath.Join(schemas, "user.ts"),
		[]byte("import { z } from \"zod\";\n\nexport const RegisterSchema = z.object({\n  email: z.string(),\n});\n\nexport const UpdateUserSchema = z.object({\n  email: z.string().optional(),\n});\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(types, "user.ts"),
		[]byte("export interface User {\n  id: string;\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	def, err := ParseInlineFields("User", "github_username:string:optional")
	if err != nil {
		t.Fatal(err)
	}
	g := newTestGenerator(root, "app/apps/api", def)
	if err := g.injectFrontendField(BuildNames(def), def.Fields[len(def.Fields)-1]); err != nil {
		t.Fatalf("adding a field to User: %v", err)
	}

	schema, err := os.ReadFile(filepath.Join(schemas, "user.ts"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(schema), "github_username") {
		t.Errorf("the field never reached the schema:\n%s", schema)
	}
	if strings.Count(string(schema), "github_username") != 2 {
		t.Errorf("want the field in both the register and update schemas:\n%s", schema)
	}
	interfaceFile, err := os.ReadFile(filepath.Join(types, "user.ts"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(interfaceFile), "github_username") {
		t.Errorf("the field never reached the type:\n%s", interfaceFile)
	}
}
