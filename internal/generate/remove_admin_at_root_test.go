package generate

import (
	"os"
	"path/filepath"
	"testing"
)

// Every single-project entry in the removal's path list assumed the panel's
// routes sit under src/routes/admin/. `grit new --api --admin` serves the panel
// at the root, so they sit at src/routes/, and a removed resource left its
// route files behind: still in the route tree, still importing a page whose
// resource definition had just been deleted, which is a build error rather than
// a stale screen (grit#89).
func TestRemoveDeletesRootMountedAdminRoutes(t *testing.T) {
	root := t.TempDir()
	write := func(rel string) string {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("export const Route = {}\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		return path
	}

	// A single project, which is what --admin always is.
	if err := os.WriteFile(filepath.Join(root, "grit.json"),
		[]byte(`{"architecture":"single","frontend":"tanstack"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	routes := []string{
		"src/routes/_dashboard/resources/blogs/index.tsx",
		"src/routes/_dashboard/resources/blogs/$id.tsx",
		"src/routes/_dashboard/resources/blogs.tsx",
	}
	for _, rel := range routes {
		write(rel)
	}
	// A sibling the removal has no business touching, so a test that passes
	// because everything was deleted fails here instead.
	keep := write("src/routes/_dashboard/resources/users/index.tsx")

	if err := RemoveResourceAt(root, "Blog"); err != nil {
		t.Fatalf("RemoveResourceAt: %v", err)
	}

	for _, rel := range routes {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); err == nil {
			t.Errorf("%s survived the removal: it imports a page that is now deleted", rel)
		}
	}
	// And the folder, which left behind still reads as "the blog is sort of
	// here".
	dir := filepath.Join(root, "src", "routes", "_dashboard", "resources", "blogs")
	if _, err := os.Stat(dir); err == nil {
		t.Error("the resource's empty route folder is still there")
	}
	if _, err := os.Stat(keep); err != nil {
		t.Errorf("the removal took another resource's routes with it: %v", err)
	}
}
