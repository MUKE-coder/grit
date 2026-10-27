package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

// An upgrade put back what a removal had taken out.
//
// Three of the blog example's files are written from the scaffold's own maps
// rather than by the generator, so `grit remove resource Blog` deleted them
// and the next `grit upgrade` wrote them back. The seeder refers to
// models.Blog, which was gone, so a project that built before the upgrade did
// not build after it, in files the developer had never touched (grit#111).
func TestUpgradeDoesNotResurrectARemovedResource(t *testing.T) {
	root := t.TempDir()
	api := filepath.Join(root, "apps", "api", "internal", "models")
	if err := os.MkdirAll(api, 0o755); err != nil {
		t.Fatal(err)
	}
	opts := Options{Architecture: ArchTriple, Upgrading: true}

	files := map[string]string{
		filepath.Join(root, "apps", "api", "internal", "database", "blogs_seeder.go"): "seeder",
		filepath.Join(root, "apps", "web", "hooks", "use-blogs.ts"):                   "hook",
		filepath.Join(root, "apps", "web", "lib", "blog-api.ts"):                      "lib",
		filepath.Join(root, "apps", "api", "internal", "database", "seed.go"):         "keep me",
	}

	// No blog.go: the resource was removed.
	kept := withoutRemoved(copyOf(files), root, opts)
	if len(kept) != 1 {
		t.Errorf("expected only the non-blog file to survive, got %d: %v", len(kept), keys(kept))
	}
	if _, ok := kept[filepath.Join(root, "apps", "api", "internal", "database", "seed.go")]; !ok {
		t.Error("seed.go was dropped; only the removed resource's files should be")
	}

	// And with the model present, every file is written as before. This is the
	// direction that matters: a guard that is slightly too eager stops every
	// project in the world getting its blog.
	if err := os.WriteFile(filepath.Join(api, "blog.go"), []byte("package models"), 0o644); err != nil {
		t.Fatal(err)
	}
	all := withoutRemoved(copyOf(files), root, opts)
	if len(all) != len(files) {
		t.Errorf("a project that still has its blog lost files: %v", keys(all))
	}
}

// On a new project the models have not been written when the seeder map is
// built, so asking "does blog.go exist" there answers no for every project
// Grit creates. Upgrading is what tells the two apart.
func TestANewProjectStillGetsTheBlog(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		filepath.Join(root, "apps", "api", "internal", "database", "blogs_seeder.go"): "seeder",
	}
	kept := withoutRemoved(copyOf(files), root, Options{Architecture: ArchTriple})
	if len(kept) != 1 {
		t.Error("a new project was refused the blog example because its model " +
			"had not been written yet")
	}
}

// The same file lands in a different place in each architecture, so the match
// is on the name rather than the path.
func TestBlogFilesAreRecognisedWhereverTheyLand(t *testing.T) {
	for _, path := range []string{
		filepath.Join("apps", "web", "hooks", "use-blogs.ts"),
		filepath.Join("frontend", "src", "hooks", "use-blogs.ts"),
		filepath.Join("apps", "expo", "hooks", "use-blogs.ts"),
		filepath.Join("src", "hooks", "use-blogs.ts"),
		filepath.Join("api", "internal", "database", "blogs_seeder.go"),
		filepath.Join("apps", "web", "lib", "blog-api.ts"),
	} {
		if !isBlogExampleFile(path) {
			t.Errorf("%s was not recognised as part of the blog example", path)
		}
	}
	for _, path := range []string{
		filepath.Join("apps", "api", "internal", "database", "seed.go"),
		filepath.Join("apps", "web", "hooks", "use-users.ts"),
		filepath.Join("apps", "api", "internal", "models", "blog.go"),
	} {
		if isBlogExampleFile(path) {
			t.Errorf("%s was treated as a blog example file and would be dropped", path)
		}
	}
}

func copyOf(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func keys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, filepath.Base(k))
	}
	return out
}
