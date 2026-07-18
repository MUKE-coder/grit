package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ── removeLinesContaining ─────────────────────────────────────────────────────

func TestRemoveLinesContaining(t *testing.T) {
	t.Run("removes matching lines", func(t *testing.T) {
		f := writeTempFile(t, "models.go", `package models

func Models() []interface{} {
	return []interface{}{
		&User{},
		&Post{},
		// grit:models
	}
}
`)
		if err := removeLinesContaining(f, "&Post{}"); err != nil {
			t.Fatalf("removeLinesContaining: %v", err)
		}
		got := readFile(t, f)
		if strings.Contains(got, "&Post{}") {
			t.Error("&Post{} should have been removed")
		}
		if !strings.Contains(got, "&User{}") {
			t.Error("&User{} should still be present")
		}
		if !strings.Contains(got, "// grit:models") {
			t.Error("// grit:models marker should still be present")
		}
	})

	t.Run("returns error when pattern not found", func(t *testing.T) {
		f := writeTempFile(t, "models.go", "package models\n")
		err := removeLinesContaining(f, "&Missing{}")
		if err == nil {
			t.Error("expected error when pattern not found, got nil")
		}
	})

	t.Run("file not found returns error", func(t *testing.T) {
		err := removeLinesContaining("/tmp/does-not-exist-remove-test.go", "&Post{}")
		if err == nil {
			t.Error("expected error for missing file, got nil")
		}
	})
}

// ── removeInlineText ──────────────────────────────────────────────────────────

func TestRemoveInlineText(t *testing.T) {
	t.Run("removes inline text", func(t *testing.T) {
		f := writeTempFile(t, "routes.go",
			"studio.Mount(r, db, []interface{}{&models.User{}, &models.Post{}, /* grit:studio */}, cfg)\n")
		if err := removeInlineText(f, "&models.Post{}, "); err != nil {
			t.Fatalf("removeInlineText: %v", err)
		}
		got := readFile(t, f)
		if strings.Contains(got, "&models.Post{}") {
			t.Error("&models.Post{} should have been removed")
		}
		if !strings.Contains(got, "&models.User{}") {
			t.Error("&models.User{} should still be present")
		}
		if !strings.Contains(got, "/* grit:studio */") {
			t.Error("studio marker should still be present")
		}
	})

	t.Run("returns error when text not found", func(t *testing.T) {
		f := writeTempFile(t, "routes.go", "package routes\n")
		err := removeInlineText(f, "&models.Missing{}, ")
		if err == nil {
			t.Error("expected error when text not found, got nil")
		}
	})
}

// ── removeLineBlock ───────────────────────────────────────────────────────────

func TestRemoveLineBlock(t *testing.T) {
	t.Run("removes handler initialization block", func(t *testing.T) {
		f := writeTempFile(t, "routes.go", `package routes

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	userHandler := &handlers.UserHandler{
		DB: db,
	}
	postHandler := &handlers.PostHandler{
		DB: db,
	}
	// grit:handlers
}
`)
		if err := removeLineBlock(f,
			"postHandler := &handlers.PostHandler{",
			"}"); err != nil {
			t.Fatalf("removeLineBlock: %v", err)
		}
		got := readFile(t, f)
		if strings.Contains(got, "postHandler") {
			t.Error("postHandler block should have been removed")
		}
		if !strings.Contains(got, "userHandler") {
			t.Error("userHandler block should still be present")
		}
		if !strings.Contains(got, "// grit:handlers") {
			t.Error("// grit:handlers marker should still be present")
		}
	})

	t.Run("returns error when block not found", func(t *testing.T) {
		f := writeTempFile(t, "routes.go", "package routes\n")
		err := removeLineBlock(f, "missingHandler := &handlers.MissingHandler{", "}")
		if err == nil {
			t.Error("expected error when block not found, got nil")
		}
	})
}

// ── removeSchemaExportBlock ───────────────────────────────────────────────────

func TestRemoveSchemaExportBlock(t *testing.T) {
	t.Run("removes multi-line schema export block", func(t *testing.T) {
		f := writeTempFile(t, "index.ts", `export {
  CreateUserSchema,
  UpdateUserSchema,
  type CreateUserInput,
  type UpdateUserInput,
} from "./user";
export {
  CreatePostSchema,
  UpdatePostSchema,
  type CreatePostInput,
  type UpdatePostInput,
} from "./post";
// grit:schemas
`)
		if err := removeSchemaExportBlock(f, "Post", "post"); err != nil {
			t.Fatalf("removeSchemaExportBlock: %v", err)
		}
		got := readFile(t, f)
		if strings.Contains(got, "CreatePostSchema") {
			t.Error("Post schema export should have been removed")
		}
		if !strings.Contains(got, "CreateUserSchema") {
			t.Error("User schema export should still be present")
		}
		if !strings.Contains(got, "// grit:schemas") {
			t.Error("// grit:schemas marker should still be present")
		}
	})

	t.Run("returns error when schema not found", func(t *testing.T) {
		f := writeTempFile(t, "index.ts", "// grit:schemas\n")
		err := removeSchemaExportBlock(f, "Missing", "missing")
		if err == nil {
			t.Error("expected error when schema not found, got nil")
		}
	})
}

// ── RemoveResource (via direct file manipulation, not findProjectRoot) ────────

// setupProjectWithResource generates a resource into a temp project,
// returning the root path for removal testing.
func setupProjectWithResource(t *testing.T, resourceName, fields, module string) string {
	t.Helper()
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields(resourceName, fields)
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}

	g := newTestGenerator(root, module, def)
	if err := g.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	return root
}

func TestRemoveResource_Files(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupProjectWithResource(t, "Post", "title:string,content:text", module)

	// Verify files exist before removal
	modelPath := filepath.Join(root, "apps", "api", "internal", "models", "post.go")
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Fatal("model file must exist before removal")
	}

	// We can't call RemoveResource() directly (uses findProjectRoot),
	// so test the sub-functions manually:

	// Test file deletion
	if err := os.Remove(modelPath); err != nil {
		t.Fatalf("removing model file: %v", err)
	}
	if _, err := os.Stat(modelPath); !os.IsNotExist(err) {
		t.Error("model file should be gone after removal")
	}

	// Verify that other generated files still exist (service, handler, etc.)
	svcPath := filepath.Join(root, "apps", "api", "internal", "services", "post.go")
	if _, err := os.Stat(svcPath); os.IsNotExist(err) {
		t.Error("service file should still exist (we only removed model)")
	}
}

func TestRemoveResource_Injections(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupProjectWithResource(t, "Post", "title:string,content:text", module)

	// Reverse: remove model injection from user.go
	userGoPath := filepath.Join(root, "apps", "api", "internal", "models", "user.go")
	if err := removeLinesContaining(userGoPath, "&Post{}"); err != nil {
		t.Fatalf("removeLinesContaining Post: %v", err)
	}
	got := readFile(t, userGoPath)
	if strings.Contains(got, "&Post{}") {
		t.Error("&Post{} should have been removed from user.go")
	}
	if !strings.Contains(got, "// grit:models") {
		t.Error("// grit:models marker must remain after removal")
	}

	// Reverse: remove studio injection
	routesPath := filepath.Join(root, "apps", "api", "internal", "routes", "routes.go")
	if err := removeInlineText(routesPath, "&models.Post{}, "); err != nil {
		t.Fatalf("removeInlineText: %v", err)
	}
	routesContent := readFile(t, routesPath)
	if strings.Contains(routesContent, "&models.Post{}") {
		t.Error("&models.Post{} should be removed from routes.go")
	}

	// Reverse: remove handler initialization block
	if err := removeLineBlock(routesPath,
		"postHandler := &handlers.PostHandler{",
		"}"); err != nil {
		t.Fatalf("removeLineBlock: %v", err)
	}
	// Also remove route lines that reference postHandler (protected + admin routes)
	removeLinesContaining(routesPath, "postHandler.")

	routesContent = readFile(t, routesPath)
	if strings.Contains(routesContent, "postHandler") {
		t.Error("postHandler references should be fully removed from routes.go")
	}

	// Reverse: remove schema export
	schemaIndexPath := filepath.Join(root, "packages", "shared", "schemas", "index.ts")
	if err := removeSchemaExportBlock(schemaIndexPath, "Post", "post"); err != nil {
		t.Fatalf("removeSchemaExportBlock: %v", err)
	}
	schemaIndex := readFile(t, schemaIndexPath)
	if strings.Contains(schemaIndex, "CreatePostSchema") {
		t.Error("CreatePostSchema should be removed from schemas/index.ts")
	}
	if !strings.Contains(schemaIndex, "// grit:schemas") {
		t.Error("// grit:schemas marker must remain")
	}

	// Reverse: remove type export
	typesIndexPath := filepath.Join(root, "packages", "shared", "types", "index.ts")
	if err := removeLinesContaining(typesIndexPath, `from "./post"`); err != nil {
		t.Fatalf("removeLinesContaining type: %v", err)
	}
	typesIndex := readFile(t, typesIndexPath)
	if strings.Contains(typesIndex, "Post") {
		t.Error("Post type export should be removed from types/index.ts")
	}
	if !strings.Contains(typesIndex, "// grit:types") {
		t.Error("// grit:types marker must remain")
	}

	// Reverse: remove API route constants
	constantsPath := filepath.Join(root, "packages", "shared", "constants", "index.ts")
	if err := removeLineBlock(constantsPath, "POSTS: {", "},"); err != nil {
		t.Fatalf("removeLineBlock constants: %v", err)
	}
	constants := readFile(t, constantsPath)
	if strings.Contains(constants, "POSTS") {
		t.Error("POSTS route block should be removed from constants/index.ts")
	}
}

// TestGenerateAndRemove_RoundTrip verifies that after a full generate+remove
// cycle the marker files are back to their original (empty-injection) state.
func TestGenerateAndRemove_RoundTrip(t *testing.T) {
	const module = "myapp/apps/api"

	// Capture original content
	root := setupMinimalProject(t, module)
	origUserGo := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "user.go"))
	origTypesIndex := readTestFile(t, filepath.Join(root, "packages", "shared", "types", "index.ts"))

	// Generate
	def, err := ParseInlineFields("Widget", "name:string")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	if err := g.Run(); err != nil {
		t.Fatalf("Run(): %v", err)
	}

	// Verify injection happened
	after := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "user.go"))
	if !strings.Contains(after, "&Widget{}") {
		t.Fatal("Widget should have been injected")
	}

	// Reverse injections manually
	userGoPath := filepath.Join(root, "apps", "api", "internal", "models", "user.go")
	removeLinesContaining(userGoPath, "&Widget{}")

	routesPath := filepath.Join(root, "apps", "api", "internal", "routes", "routes.go")
	removeInlineText(routesPath, "&models.Widget{}, ")
	removeLineBlock(routesPath, "widgetHandler := &handlers.WidgetHandler{", "}")
	removeLinesContaining(routesPath, "widgetHandler.")

	schemaPath := filepath.Join(root, "packages", "shared", "schemas", "index.ts")
	removeSchemaExportBlock(schemaPath, "Widget", "widget")

	typesPath := filepath.Join(root, "packages", "shared", "types", "index.ts")
	removeLinesContaining(typesPath, `from "./widget"`)

	constantsPath := filepath.Join(root, "packages", "shared", "constants", "index.ts")
	removeLineBlock(constantsPath, "WIDGETS: {", "},")

	// After removal: user.go should match original
	restoredUserGo := readTestFile(t, userGoPath)
	if restoredUserGo != origUserGo {
		t.Errorf("user.go not restored to original after removal.\nWant:\n%s\nGot:\n%s",
			origUserGo, restoredUserGo)
	}

	// After removal: types/index.ts should match original
	restoredTypes := readTestFile(t, typesPath)
	if restoredTypes != origTypesIndex {
		t.Errorf("types/index.ts not restored to original.\nWant:\n%s\nGot:\n%s",
			origTypesIndex, restoredTypes)
	}

	// Marker must still be present
	if !strings.Contains(restoredUserGo, "// grit:models") {
		t.Error("// grit:models marker missing after round-trip")
	}
}

// TestRemoveCaseBlocks covers the switch-dispatch arms that `generate resource`
// injects. Missing these was why a removed resource still failed to compile:
// the arm kept referencing models.<Name> after the model file was deleted.
func TestRemoveCaseBlocks(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "dispatch.go")

	src := `package services

func dispatch(name string) error {
	switch name {
	case "users":
		return handle(&models.User{})
	case "gadgets":
		return handle(&models.Gadget{})

	// grit:resource-stats:dispatch
	default:
		return nil
	}
}
`
	if err := os.WriteFile(path, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	if n := removeCaseBlocks(path, `case "gadgets":`); n != 1 {
		t.Fatalf("removeCaseBlocks = %d, want 1", n)
	}

	got := readTestFile(t, path)
	if strings.Contains(got, "Gadget") {
		t.Errorf("gadgets arm still present:\n%s", got)
	}
	// The neighbouring arm, the marker and the default must survive.
	for _, keep := range []string{`case "users":`, "models.User{}", "// grit:resource-stats:dispatch", "default:"} {
		if !strings.Contains(got, keep) {
			t.Errorf("removal ate %q:\n%s", keep, got)
		}
	}

	// Removing an absent arm must be a no-op, not a mangle.
	if n := removeCaseBlocks(path, `case "absent":`); n != 0 {
		t.Errorf("removing an absent case reported %d removals", n)
	}
}

// TestPruneUnusedImports guards the specific trap that made the first fix fail:
// these dispatch files describe the codegen contract in prose ("via
// json.Marshal(fields)"), so a naive usage check matches the COMMENT and keeps
// an import that is genuinely unused — trading one compile error for another.
func TestPruneUnusedImports(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "svc.go")

	src := `package services

import (
	"fmt"

	"encoding/json"

	"example.com/app/internal/models"
)

// Each case re-marshals fields via json.Marshal(fields) — prose only.
func run() error {
	return fmt.Errorf("nothing uses json or models now")
}
`
	if err := os.WriteFile(path, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	pruneUnusedImports(path)
	got := readTestFile(t, path)

	if strings.Contains(got, `"encoding/json"`) {
		t.Errorf("encoding/json kept despite only appearing in a comment:\n%s", got)
	}
	if strings.Contains(got, `"example.com/app/internal/models"`) {
		t.Errorf("unused models import kept:\n%s", got)
	}
	if !strings.Contains(got, `"fmt"`) {
		t.Errorf("fmt is used and must be kept:\n%s", got)
	}
}

// TestRemoveMarkedRegion covers demo content embedded in a hand-designed file
// (the web home page's blog section), which has no standalone file to delete.
func TestRemoveMarkedRegion(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "page.tsx")

	src := `<div>
  <Hero />
  {/* grit:home:blog-start */}
  <BlogSection />
  {/* grit:home:blog-end */}
  <Footer />
</div>
`
	if err := os.WriteFile(path, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	if !removeMarkedRegion(path, "grit:home:blog-start", "grit:home:blog-end") {
		t.Fatal("removeMarkedRegion reported no cut")
	}
	got := readTestFile(t, path)
	if strings.Contains(got, "BlogSection") {
		t.Errorf("marked region survived:\n%s", got)
	}
	for _, keep := range []string{"<Hero />", "<Footer />"} {
		if !strings.Contains(got, keep) {
			t.Errorf("removal ate %q:\n%s", keep, got)
		}
	}

	// An unmatched start marker must NOT truncate to end-of-file.
	partial := filepath.Join(dir, "partial.tsx")
	if err := os.WriteFile(partial, []byte("a\n{/* grit:x-start */}\nb\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if removeMarkedRegion(partial, "grit:x-start", "grit:x-end") {
		t.Error("unmatched markers should be left alone")
	}
	if !strings.Contains(readTestFile(t, partial), "b") {
		t.Error("unmatched markers truncated the file")
	}
}

// TestRemoveResource_SyncRegistryNotMangled pins the ordering trap found by
// running a real release: stripping ", &models.X{}" from inline model lists also
// matches
//
//	syncRegistry.Register("xs", &models.X{})
//
// turning it into Register("xs") — a call with too few arguments that the
// whole-line removal then no longer recognised. The registration must be deleted
// BEFORE the inline surgery runs.
func TestRemoveResource_SyncRegistryNotMangled(t *testing.T) {
	dir := t.TempDir()
	routesDir := filepath.Join(dir, "apps", "api", "internal", "routes")
	if err := os.MkdirAll(routesDir, 0755); err != nil {
		t.Fatal(err)
	}
	routes := filepath.Join(routesDir, "routes.go")

	src := `package routes

func Setup() {
	studio.Mount(r, db, []interface{}{&models.User{}, &models.Product{}, /* grit:studio */}, cfg)
	pulseCfg.Models = []interface{}{&models.User{}, &models.Product{}}
	syncRegistry.Register("users", &models.User{})
	syncRegistry.Register("products", &models.Product{})
}
`
	if err := os.WriteFile(routes, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	// Mirror RemoveResource's order: registration first, then inline lists.
	removeLinesContaining(routes, `syncRegistry.Register("products", &models.Product{})`)
	removeInlineText(routes, "&models.Product{}, ")
	removeInlineText(routes, ", &models.Product{}")

	got := readTestFile(t, routes)

	if strings.Contains(got, "Product") {
		t.Errorf("Product survived removal:\n%s", got)
	}
	// The surviving registration must keep BOTH arguments.
	if !strings.Contains(got, `syncRegistry.Register("users", &models.User{})`) {
		t.Errorf("user registration was mangled — this is the bug:\n%s", got)
	}
	if strings.Contains(got, `Register("users")`) {
		t.Errorf("registration lost its model argument:\n%s", got)
	}
	// The studio marker and the remaining model must survive.
	for _, keep := range []string{"grit:studio", "&models.User{}"} {
		if !strings.Contains(got, keep) {
			t.Errorf("removal ate %q:\n%s", keep, got)
		}
	}
}

// TestRemoveStructBlock covers the permission-catalog entry that
// `grit generate resource` injects. The literal nests three levels
// (Module > Group > Feature), so an indentation-based scan would stop at the
// first inner closing brace and leave a fragment that doesn't compile.
func TestRemoveStructBlock(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "permissions.go")

	src := `package authz

func generatedModules() []Module {
	return []Module{
		// grit:perms:auto-start
		{
			Key:  "orders",
			Name: "Orders",
			Groups: []Group{
				{
					Key:  "orders",
					Name: "Orders",
					Features: []Feature{
						{Key: "orders", Name: "Orders", Actions: AllActions},
					},
				},
			},
		},
		{
			Key:  "products",
			Name: "Products",
			Groups: []Group{
				{
					Key:  "products",
					Name: "Products",
					Features: []Feature{
						{Key: "products", Name: "Products", Actions: AllActions},
					},
				},
			},
		},
		// grit:perms:auto-end
	}
}
`
	if err := os.WriteFile(path, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	if !removeStructBlock(path, `Key:  "products",`) {
		t.Fatal("removeStructBlock reported no removal")
	}
	got := readTestFile(t, path)

	if strings.Contains(got, "products") {
		t.Errorf("products entry survived:\n%s", got)
	}
	// The sibling entry must be intact — this is what naive brace matching eats.
	if !strings.Contains(got, `Key:  "orders",`) || !strings.Contains(got, `{Key: "orders", Name: "Orders", Actions: AllActions}`) {
		t.Errorf("sibling entry was damaged:\n%s", got)
	}
	// Markers must survive so the next generate still has an anchor.
	for _, m := range []string{"grit:perms:auto-start", "grit:perms:auto-end"} {
		if !strings.Contains(got, m) {
			t.Errorf("marker %q was removed:\n%s", m, got)
		}
	}
	// Braces must still balance, or the file won't compile.
	if strings.Count(got, "{") != strings.Count(got, "}") {
		t.Errorf("unbalanced braces after removal:\n%s", got)
	}

	// Removing an absent entry is a no-op.
	if removeStructBlock(path, `Key:  "absent",`) {
		t.Error("removing an absent block reported success")
	}
}
