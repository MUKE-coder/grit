package manifest

import (
	"os"
	"path/filepath"
	"testing"
)

// A recording session records what was written, and the test that says so is
// about the COUNT as well as the keys.
//
// v3.505.1 changed Rel to resolve a relative path against the working
// directory, which fixed `grit upgrade`'s repairs and broke `grit new`: that
// command starts the session with the relative directory it was handed and
// then writes absolute paths, so the pair failed the other way round and a
// freshly scaffolded project recorded nothing. 441 entries became 0, and the
// whole suite stayed green, because every manifest test built its manifest by
// hand instead of going through a session the way the commands do.
//
// This is deliberately the cheap version of that: a session opened the way
// `grit new` opens one, writes the way the scaffold writes, and an assertion
// that the entries are there. It does not need the scaffold itself to be
// right, only the recording, which is what broke.
func TestASessionOpenedWithARelativeRootStillRecords(t *testing.T) {
	tmp := t.TempDir()
	before, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(before) })
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Exactly `grit new app`: the root is the relative directory the user
	// typed, and the files are written by absolute path.
	const relativeRoot = "app"
	if err := os.MkdirAll(filepath.Join(relativeRoot, "apps", "api", "internal", "services"), 0o755); err != nil {
		t.Fatal(err)
	}

	release, err := Start(relativeRoot, "3.505.2", "scaffold")
	if err != nil {
		t.Fatal(err)
	}

	written := []string{
		filepath.Join("apps", "api", "internal", "services", "memo.go"),
		filepath.Join("apps", "api", "internal", "services", "user.go"),
	}
	for _, rel := range written {
		abs := filepath.Join(cwd, relativeRoot, rel)
		if err := os.WriteFile(abs, []byte("package services\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		Note(abs, "package services\n")
	}
	if err := release(); err != nil {
		t.Fatal(err)
	}

	m, err := Load(relativeRoot)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(m.Entries()); got != len(written) {
		t.Fatalf("the session recorded %d of %d files it was told about; a scaffolded project "+
			"with an empty manifest gets no framework fixes ever again, and nothing says so",
			got, len(written))
	}
	for _, rel := range written {
		key := Key(rel)
		if m.StatusOf(relativeRoot, key) != Unchanged {
			t.Errorf("%s reads as %v rather than Grit's own", key, m.StatusOf(relativeRoot, key))
		}
	}
}
