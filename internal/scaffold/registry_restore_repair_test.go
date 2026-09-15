package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Upgrading a double wrote the template over the admin's resource registry, and
// every resource grit generate had registered vanished from the panel.
func TestEmbeddedAdminUpgradeLeavesUserOwnedFiles(t *testing.T) {
	web := filepath.Join("p", "apps", "web")
	for path, want := range map[string]bool{
		filepath.Join(web, "admin-panel", "resources", "index.ts"):                          true,
		filepath.Join(web, "admin-panel", "resources", "notes", "notes.ts"):                 true,
		filepath.Join(web, "admin-panel", "resources", "users", "users.ts"):                 false,
		filepath.Join(web, "admin-panel", "components", "shared", "providers.tsx"):          false,
		filepath.Join(web, "app", "admin", "(dashboard)", "resources", "notes", "page.tsx"): true,
		filepath.Join(web, "app", "admin", "(dashboard)", "dashboard", "page.tsx"):          false,
		filepath.Join(web, "app", "(marketing)", "page.tsx"):                                false,
	} {
		if got := isUserOwnedEmbeddedAdminFile(web, path); got != want {
			t.Errorf("isUserOwnedEmbeddedAdminFile(%s) = %v, want %v", path, got, want)
		}
	}
	spa := filepath.Join("p", "src")
	if !isUserOwnedEmbeddedAdminFile(spa, filepath.Join(spa, "admin-panel", "resources", "index.ts")) {
		t.Error("the SPA panel's registry is still rewritten on upgrade")
	}
}

func TestRestoreResourceRegistrations(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "apps", "web", "admin-panel", "resources")
	write := func(path, body string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	index := filepath.Join(dir, "index.ts")
	template := strings.ReplaceAll(adminResourceRegistry(), `"@/lib/resource"`, `"@admin/lib/resource"`)
	write(index, template)
	write(filepath.Join(dir, "users", "users.ts"), "export const usersResource = defineResource({});\n")
	write(filepath.Join(dir, "blogs", "blogs.ts"), "export const blogsResource = defineResource({});\n")
	write(filepath.Join(dir, "notes", "notes.ts"), "export const noteResource = defineResource({});\n")
	write(filepath.Join(dir, "tags", "tags.ts"), "export const tagResource = defineResource({});\n")

	restored, err := restoreResourceRegistrations(root, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(restored, ",") != "noteResource,tagResource" {
		t.Fatalf("restored %v, want noteResource and tagResource", restored)
	}
	got, _ := os.ReadFile(index)
	for _, want := range []string{
		"import { noteResource } from \"./notes/notes\";\nimport { tagResource } from \"./tags/tags\";\n// grit:resources",
		"  blogsResource,\n  noteResource,\n  tagResource,\n  // grit:resource-list",
	} {
		if !strings.Contains(string(got), want) {
			t.Errorf("the registry is missing %q:\n%s", want, got)
		}
	}
	if again, _ := restoreResourceRegistrations(root, Options{}); len(again) != 0 {
		t.Errorf("a second run registered %v again", again)
	}

	// A registry a person edited, with a resource of their own, is left alone.
	edited := strings.Replace(template, "// grit:resources", "import { productResource } from \"./products/products\";\n// grit:resources", 1)
	write(index, edited)
	if restored, _ := restoreResourceRegistrations(root, Options{}); len(restored) != 0 {
		t.Errorf("an edited registry was changed: %v", restored)
	}
	if got, _ := os.ReadFile(index); string(got) != edited {
		t.Error("an edited registry was rewritten")
	}
}

// The blog editor saved, published and deleted through /api/blogs/:id, which the
// API does not serve: every save returned 404.
func TestBlogEditorWritesThroughAdminRoutes(t *testing.T) {
	raw, err := os.ReadFile("admin_v3317_blog.go")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), blogEditorOldRoute) {
		t.Error("the blog editor template still writes to /api/blogs/:id")
	}
	old := "a(" + blogEditorOldRoute + ", patch)\nb(" + blogEditorOldRoute + ")\n"
	out, fixed := repairBlogEditorRoutesSource(old)
	if strings.Contains(out, blogEditorOldRoute) || strings.Count(out, `"/api/admin/blogs/" + params.id`) != 2 || len(fixed) != 1 {
		t.Errorf("the blog editor routes were not repaired: %q %v", out, fixed)
	}
	if again, fixed := repairBlogEditorRoutesSource(out); again != out || len(fixed) != 0 {
		t.Error("the blog editor route repair is not idempotent")
	}
}
