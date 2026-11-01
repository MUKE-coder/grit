package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

// grit#172: `grit upgrade` re-created a removed resource's shared type.
//
// v3.460.0 put blog and upload on SharedModelTypes, so the comparator test
// would cover them. The list is also what repairSharedTypes walks, and that is
// createIfMissing over every entry, so a blog.ts absent because somebody
// removed the resource looked exactly like a blog.ts absent because the
// project predates the list. The Blog's type and schema came back on the next
// upgrade, and types/index.ts re-exported them.
//
// The model is what tells the two apart: if it is not there, grit sync would
// not write the type either.
func TestASharedTypeIsNotWrittenForAModelTheProjectRemoved(t *testing.T) {
	root := t.TempDir()
	models := filepath.Join(root, "apps", "api", "internal", "models")
	if err := os.MkdirAll(models, 0o755); err != nil {
		t.Fatal(err)
	}
	// Every model except blog, which is what removing the Blog leaves behind.
	for _, typ := range SharedModelTypes {
		if typ.Kebab == "blog" {
			continue
		}
		name := filepath.Join(models, snakeOf(typ.Kebab)+".go")
		if err := os.WriteFile(name, []byte("package models\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	shared := filepath.Join(root, "packages", "shared")
	files := sharedModelTypeFiles(root, shared)

	for _, gone := range []string{
		filepath.Join(shared, "types", "blog.ts"),
		filepath.Join(shared, "schemas", "blog.ts"),
	} {
		if _, ok := files[gone]; ok {
			t.Errorf("would write %s for a project whose blog model is gone, which undoes the removal", gone)
		}
	}
	// And the rest still arrive, or the gate is doing too much.
	if _, ok := files[filepath.Join(shared, "types", "api-key.ts")]; !ok {
		t.Error("api-key.ts was withheld although the project has that model")
	}
}

// The scaffold passes no root, because it is writing the models in the same
// run and asking whether they exist yet would skip all of them.
func TestTheScaffoldStillGetsEverySharedType(t *testing.T) {
	shared := filepath.Join("packages", "shared")
	files := sharedModelTypeFiles("", shared)
	for _, typ := range SharedModelTypes {
		if _, ok := files[filepath.Join(shared, "types", typ.Kebab+".ts")]; !ok {
			t.Errorf("a new project would not get types/%s.ts", typ.Kebab)
		}
	}
}

// A tree with no models directory is not a project that removed something, it
// is a question we cannot answer. Withholding there would mean a frontend-only
// checkout, or a test fixture, silently losing its shared types.
func TestNoModelsDirectoryMeansWriteThemAll(t *testing.T) {
	root := t.TempDir()
	shared := filepath.Join(root, "packages", "shared")
	files := sharedModelTypeFiles(root, shared)
	if len(files) == 0 {
		t.Error("no shared types at all for a tree with no models directory")
	}
}

func snakeOf(kebab string) string {
	out := make([]rune, 0, len(kebab))
	for _, r := range kebab {
		if r == '-' {
			out = append(out, '_')
			continue
		}
		out = append(out, r)
	}
	return string(out)
}
