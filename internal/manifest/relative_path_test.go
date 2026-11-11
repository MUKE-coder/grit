package manifest

import (
	"os"
	"path/filepath"
	"testing"
)

// Every repair grit upgrade runs after scaffold.Upgrade returns passes a
// relative root: filepath.Join("apps", "api"). So the paths reaching Refresh
// are relative, filepath.Rel was handed one absolute and one relative path,
// and the error read here as "outside the project".
//
// The effect was not a missing record but a wrong one: the file kept the hash
// it had before the repair, so from then on `grit upgrade` treated it as
// hand-edited and `grit generate resource` refused to regenerate it. Found
// when a repair fixed a service and the generator then declined to touch it.
//
// TestRepairsDoNotWriteBehindTheManifest could not catch this. It checks that
// a repair calls Refresh, which they all did.
func TestRelResolvesARelativePathAgainstTheWorkingDirectory(t *testing.T) {
	root := t.TempDir()
	before, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// Registered before the chdir, so the test cannot leave the process inside
	// a directory t.TempDir is about to remove: on Windows that fails the test
	// in cleanup, which is what the first version of this did.
	t.Cleanup(func() { _ = os.Chdir(before) })
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	// What the OS reports, not root: t.TempDir can hand back a symlinked path
	// whose resolved form differs.
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	relative := filepath.Join("apps", "api", "internal", "services", "memo.go")
	key, inside := Rel(cwd, relative)
	if !inside {
		t.Fatal("a relative path inside the project was reported as outside it, so every " +
			"repair that passes one refreshes nothing and poisons the file it fixed")
	}
	if key != "apps/api/internal/services/memo.go" {
		t.Errorf("key = %q, want apps/api/internal/services/memo.go", key)
	}

	// An absolute path still behaves exactly as before.
	absolute := filepath.Join(cwd, "apps", "api", "internal", "services", "memo.go")
	if key, inside := Rel(cwd, absolute); !inside || key != "apps/api/internal/services/memo.go" {
		t.Errorf("absolute: key=%q inside=%v", key, inside)
	}

	// And something genuinely outside is still outside.
	if _, inside := Rel(filepath.Join(cwd, "apps", "api"), filepath.Join(cwd, "elsewhere.go")); inside {
		t.Error("a path above the root was reported as inside it")
	}
}

// The other direction, which the first fix broke.
//
// `grit new app` starts the recording with the relative directory it was given
// and then writes absolute paths, so here it is the ROOT that is relative.
// v3.505.1 resolved only the path, so this pair failed instead and a freshly
// scaffolded project recorded nothing at all: 441 entries became 0. One
// direction was traded for the other, which is why this test takes all four
// combinations rather than the one that was broken at the time.
func TestRelTakesEitherSideRelative(t *testing.T) {
	root := t.TempDir()
	before, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(before) })
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	const want = "apps/api/internal/services/memo.go"
	relRoot := "app"
	absRoot := filepath.Join(cwd, "app")
	relPath := filepath.Join("app", "apps", "api", "internal", "services", "memo.go")
	absPath := filepath.Join(absRoot, "apps", "api", "internal", "services", "memo.go")

	for _, c := range []struct {
		name, root, path string
	}{
		{"both absolute", absRoot, absPath},
		{"both relative", relRoot, relPath},
		{"relative path, absolute root: every grit upgrade repair", absRoot, relPath},
		{"absolute path, relative root: grit new", relRoot, absPath},
	} {
		t.Run(c.name, func(t *testing.T) {
			key, inside := Rel(c.root, c.path)
			if !inside {
				t.Fatalf("reported as outside the project, so the record is dropped silently")
			}
			if key != want {
				t.Errorf("key = %q, want %q", key, want)
			}
		})
	}
}

// Refresh over a relative path updates the recorded hash, which is the whole
// point of the call.
func TestRefreshWorksOnTheRelativePathsRepairsUse(t *testing.T) {
	root := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	here, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	rel := filepath.Join("apps", "api", "internal", "services", "memo.go")
	if err := os.MkdirAll(filepath.Dir(rel), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(rel, []byte("package services // as generated\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	release, err := Start(here, "3.505.0", "test")
	if err != nil {
		t.Fatal(err)
	}
	Note(rel, "package services // as generated\n")
	if err := release(); err != nil {
		t.Fatal(err)
	}

	// A repair rewrites it and refreshes, exactly as writeRepairedGo does.
	release, err = Start(here, "3.505.0", "repair")
	if err != nil {
		t.Fatal(err)
	}
	repaired := "package services // repaired\n"
	if err := os.WriteFile(rel, []byte(repaired), 0o644); err != nil {
		t.Fatal(err)
	}
	Refresh(rel)
	if err := release(); err != nil {
		t.Fatal(err)
	}

	m, err := Load(here)
	if err != nil {
		t.Fatal(err)
	}
	if got := m.StatusOf(here, "apps/api/internal/services/memo.go"); got != Unchanged {
		t.Errorf("after a repair the file reads as %v; it has to read as Grit's, or the next "+
			"upgrade skips it and grit generate refuses to regenerate it", got)
	}
}
