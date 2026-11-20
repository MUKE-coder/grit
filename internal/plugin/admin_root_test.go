package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// projectWithAdminAt builds a project root whose admin panel is at rel.
func projectWithAdminAt(t *testing.T, rel string) string {
	t.Helper()
	root := t.TempDir()
	if rel != "" {
		if err := os.MkdirAll(filepath.Join(root, filepath.FromSlash(rel)), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// {{ADMIN_ROOT}} resolves to wherever the panel actually is.
//
// grit#228: there was no placeholder for this, so grit-calendar's manifest
// needed 21 file entries for 11 files, every one written twice with its own
// when clause, and double, single and api projects could not be targeted at
// all.
func TestAdminRootResolvesPerShape(t *testing.T) {
	for _, c := range []struct {
		shape string
		made  string
		want  string
	}{
		{"triple + next", "apps/admin/components", "apps/admin"},
		{"triple + vite", "apps/admin/src/components", "apps/admin/src"},
		{"double + next", "apps/web/admin-panel/components", "apps/web/admin-panel"},
		{"double + vite", "apps/web/src/admin-panel/components", "apps/web/src/admin-panel"},
		{"single", "src/admin-panel/components", "src/admin-panel"},
		{"single, older layout", "frontend/src/admin-panel/components", "frontend/src/admin-panel"},
	} {
		root := projectWithAdminAt(t, c.made)
		got, ok := adminRootFor(Context{Root: root})
		if !ok {
			t.Errorf("%s: no admin panel found, though %s exists", c.shape, c.made)
			continue
		}
		if got != c.want {
			t.Errorf("%s: admin root = %q, want %q", c.shape, got, c.want)
		}
	}
}

// A Vite admin resolves to src/, not to the application directory.
//
// Both exist in that shape, so the order of the candidates is the whole
// answer: get it the other way round and every plugin writes its components
// one directory too high, where nothing imports them.
func TestAViteAdminResolvesToItsSrc(t *testing.T) {
	root := t.TempDir()
	// apps/admin AND apps/admin/src, which is what a Vite project has.
	for _, d := range []string{"apps/admin", "apps/admin/src/components"} {
		if err := os.MkdirAll(filepath.Join(root, filepath.FromSlash(d)), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	got, ok := adminRootFor(Context{Root: root})
	if !ok || got != "apps/admin/src" {
		t.Errorf("admin root = %q (%v), want apps/admin/src", got, ok)
	}
}

// A project with no admin panel says so.
func TestNoAdminPanelIsReported(t *testing.T) {
	root := projectWithAdminAt(t, "apps/api/internal")
	if got, ok := adminRootFor(Context{Root: root}); ok {
		t.Errorf("found an admin panel at %q in an api-only project", got)
	}
}

// An admin file is skipped on a project with no admin, not written to a path
// with a hole in it.
//
// The same rule a node dependency for a missing workspace already gets. The
// alternative is a file at "/components/x.tsx" relative to the project root,
// which is a directory nobody asked for and which removal would then prune.
func TestAdminFilesAreSkippedWithoutAnAdmin(t *testing.T) {
	dir := writePluginDir(t, map[string]string{
		"plugin.json": `{
			"name": "calendar", "version": "0.5.0", "summary": "A calendar",
			"files": [
				{"from": "files/widget.tsx", "to": "{{ADMIN_ROOT}}/components/widget.tsx"},
				{"from": "files/service.go", "to": "{{API_ROOT}}/internal/calendar/calendar.go"}
			],
			"directories": [
				{"from": "files/pages", "to": "{{ADMIN_ROOT}}/pages/calendar"}
			],
			"injections": [
				{"file": "{{ADMIN_ROOT}}/lib/icons.ts", "marker": "// grit:icons:map", "code": "x"},
				{"file": "{{API_ROOT}}/internal/routes/routes.go", "marker": "// grit:imports", "code": "y"}
			]
		}`,
		"files/widget.tsx":      "widget",
		"files/service.go":      "package calendar",
		"files/pages/index.tsx": "page",
	})

	p, err := LoadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	// An api-only project: no admin anywhere.
	apiOnly := Context{
		Root:         projectWithAdminAt(t, "apps/api/internal"),
		Architecture: "api",
		Module:       "example.com/app",
	}
	files := p.Files(apiOnly)
	for path := range files {
		if strings.Contains(path, "components") || strings.Contains(path, "pages") {
			t.Errorf("an admin file was written into a project with no admin: %s", path)
		}
	}
	if files["apps/api/internal/calendar/calendar.go"] != "package calendar" {
		t.Errorf("the Go file did not install: %v", keysOf(files))
	}

	injections := p.Injections(apiOnly)
	for _, in := range injections {
		if strings.Contains(in.File, "icons.ts") {
			t.Error("an admin injection was kept on a project with no admin")
		}
	}
	if len(injections) != 1 {
		t.Errorf("expected only the API injection, got %d", len(injections))
	}
}

// And with an admin, every one of them resolves.
func TestAdminFilesResolveWithAnAdmin(t *testing.T) {
	dir := writePluginDir(t, map[string]string{
		"plugin.json": `{
			"name": "calendar", "version": "0.5.0", "summary": "A calendar",
			"files": [{"from": "files/widget.tsx", "to": "{{ADMIN_ROOT}}/components/widget.tsx"}],
			"directories": [{"from": "files/pages", "to": "{{ADMIN_ROOT}}/pages/calendar"}],
			"injections": [{"file": "{{ADMIN_ROOT}}/lib/icons.ts", "marker": "// grit:icons:map", "code": "x"}]
		}`,
		"files/widget.tsx":      "widget",
		"files/pages/index.tsx": "page",
	})

	p, err := LoadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	// One manifest, two shapes, no when clauses: the point of the placeholder.
	for _, c := range []struct{ made, want string }{
		{"apps/admin/components", "apps/admin"},
		{"apps/admin/src/components", "apps/admin/src"},
	} {
		ctx := Context{
			Root:         projectWithAdminAt(t, c.made),
			Architecture: "triple",
			Module:       "example.com/app",
		}
		files := p.Files(ctx)
		if files[c.want+"/components/widget.tsx"] != "widget" {
			t.Errorf("%s: the widget did not land: %v", c.want, keysOf(files))
		}
		if files[c.want+"/pages/calendar/index.tsx"] != "page" {
			t.Errorf("%s: the tree did not land: %v", c.want, keysOf(files))
		}
		injections := p.Injections(ctx)
		if len(injections) != 1 || injections[0].File != c.want+"/lib/icons.ts" {
			t.Errorf("%s: injection targets %v", c.want, injections)
		}
	}
}
