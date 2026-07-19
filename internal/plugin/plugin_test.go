package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// newTestProject builds a minimal project with one injectable file.
func newTestProject(t *testing.T) Context {
	t.Helper()
	root := t.TempDir()

	routes := filepath.Join(root, "apps", "api", "internal", "routes")
	if err := os.MkdirAll(routes, 0755); err != nil {
		t.Fatal(err)
	}
	src := `package routes

func Setup() {
	// grit:models
	// grit:routes
}
`
	if err := os.WriteFile(filepath.Join(routes, "routes.go"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	return Context{
		Root:         root,
		Module:       "testapp/apps/api",
		APIRoot:      filepath.Join(root, "apps", "api"),
		Architecture: "triple",
		Frontend:     "next",
	}
}

func testPlugin() Plugin {
	return Plugin{
		Name:    "demo",
		Version: "1.0.0",
		Summary: "test fixture",
		Files: func(ctx Context) map[string]string {
			return map[string]string{
				"apps/api/internal/demo/demo.go": "package demo\n\nfunc Hello() string { return \"hi\" }\n",
			}
		},
		Injections: func(ctx Context) []Injection {
			return []Injection{{
				File:   "apps/api/internal/routes/routes.go",
				Marker: "// grit:routes",
				Code:   "\tdemo.Mount(r)",
			}}
		},
	}
}

func read(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestInstall_WritesFilesAndInjects(t *testing.T) {
	ctx := newTestProject(t)

	rec, err := Install(ctx, testPlugin())
	if err != nil {
		t.Fatal(err)
	}

	if got := filepath.Join(ctx.Root, "apps/api/internal/demo/demo.go"); !fileExists(got) {
		t.Error("plugin file was not written")
	}
	routes := read(t, filepath.Join(ctx.Root, "apps/api/internal/routes/routes.go"))
	if !strings.Contains(routes, "demo.Mount(r)") {
		t.Errorf("injection missing:\n%s", routes)
	}
	// The marker must survive so a later install can use it too.
	if !strings.Contains(routes, "// grit:routes") {
		t.Error("injection consumed the marker")
	}

	// The record must describe everything, since removal is driven by it.
	if len(rec.Files) != 1 || len(rec.Injections) != 1 {
		t.Fatalf("record incomplete: %d files, %d injections", len(rec.Files), len(rec.Injections))
	}
	if rec.Injections[0].Code != "\tdemo.Mount(r)" {
		t.Error("record must store the EXACT injected text")
	}
}

// Removal must restore the project to its pre-install state.
func TestRemove_RestoresOriginal(t *testing.T) {
	ctx := newTestProject(t)
	routesPath := filepath.Join(ctx.Root, "apps/api/internal/routes/routes.go")
	before := read(t, routesPath)

	if _, err := Install(ctx, testPlugin()); err != nil {
		t.Fatal(err)
	}
	warnings, err := Remove(ctx.Root, "demo")
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 0 {
		t.Errorf("unexpected warnings: %v", warnings)
	}

	if after := read(t, routesPath); after != before {
		t.Errorf("routes.go not restored.\nwant:\n%s\ngot:\n%s", before, after)
	}
	if fileExists(filepath.Join(ctx.Root, "apps/api/internal/demo/demo.go")) {
		t.Error("plugin file survived removal")
	}
	// The emptied directory should be gone too.
	if _, err := os.Stat(filepath.Join(ctx.Root, "apps/api/internal/demo")); !os.IsNotExist(err) {
		t.Error("empty plugin directory was left behind")
	}

	lock, _ := LoadLock(ctx.Root)
	if len(lock.Plugins) != 0 {
		t.Error("lockfile still lists the plugin")
	}
}

// Editing an injected block must be reported, never silently overwritten —
// removal guessing at a changed block would delete the user's work.
func TestRemove_WarnsOnEditedInjection(t *testing.T) {
	ctx := newTestProject(t)
	if _, err := Install(ctx, testPlugin()); err != nil {
		t.Fatal(err)
	}

	routesPath := filepath.Join(ctx.Root, "apps/api/internal/routes/routes.go")
	edited := strings.Replace(read(t, routesPath), "demo.Mount(r)", "demo.Mount(r, customOpts)", 1)
	if err := os.WriteFile(routesPath, []byte(edited), 0644); err != nil {
		t.Fatal(err)
	}

	warnings, err := Remove(ctx.Root, "demo")
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) == 0 {
		t.Fatal("expected a warning about the modified block")
	}
	if !strings.Contains(read(t, routesPath), "customOpts") {
		t.Error("the user's edit was destroyed")
	}
}

func TestInstall_RefusesDuplicate(t *testing.T) {
	ctx := newTestProject(t)
	if _, err := Install(ctx, testPlugin()); err != nil {
		t.Fatal(err)
	}
	if _, err := Install(ctx, testPlugin()); err == nil {
		t.Error("installing twice should fail")
	}
}

// A plugin must never clobber a file that already exists — removal would then
// delete something the plugin didn't create.
func TestInstall_RefusesToOverwrite(t *testing.T) {
	ctx := newTestProject(t)
	target := filepath.Join(ctx.Root, "apps/api/internal/demo/demo.go")
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte("// mine\n"), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := Install(ctx, testPlugin()); err == nil {
		t.Fatal("expected refusal to overwrite")
	}
	if read(t, target) != "// mine\n" {
		t.Error("existing file was overwritten")
	}
}

func TestInstall_MissingMarkerFails(t *testing.T) {
	ctx := newTestProject(t)
	p := testPlugin()
	p.Injections = func(ctx Context) []Injection {
		return []Injection{{
			File:   "apps/api/internal/routes/routes.go",
			Marker: "// grit:does-not-exist",
			Code:   "\tnope()",
		}}
	}
	if _, err := Install(ctx, p); err == nil {
		t.Error("a missing marker should fail the install, not pass silently")
	}
}

// An optional injection is for surfaces that legitimately don't exist (a
// frontend file in an --api project) and must not fail the install.
func TestInstall_OptionalInjectionSkips(t *testing.T) {
	ctx := newTestProject(t)
	p := testPlugin()
	p.Injections = func(ctx Context) []Injection {
		return []Injection{{
			File:     "apps/admin/nope.ts",
			Marker:   "// grit:x",
			Code:     "x",
			Optional: true,
		}}
	}
	if _, err := Install(ctx, p); err != nil {
		t.Errorf("optional injection should not fail the install: %v", err)
	}
}

func TestRequires(t *testing.T) {
	ctx := newTestProject(t)
	dependent := testPlugin()
	dependent.Name = "dependent"
	dependent.Requires = []string{"demo"}
	dependent.Files = func(ctx Context) map[string]string {
		return map[string]string{"apps/api/internal/dependent/x.go": "package dependent\n"}
	}
	dependent.Injections = nil

	if _, err := Install(ctx, dependent); err == nil {
		t.Fatal("expected a missing-requirement error")
	}

	if _, err := Install(ctx, testPlugin()); err != nil {
		t.Fatal(err)
	}
	if _, err := Install(ctx, dependent); err != nil {
		t.Fatalf("install should succeed once the requirement is met: %v", err)
	}

	// Removing something still depended on must be refused.
	if _, err := Remove(ctx.Root, "demo"); err == nil {
		t.Error("expected refusal: demo is required by dependent")
	}
}
