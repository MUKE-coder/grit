package scaffold

import (
	"path/filepath"
	"strings"
	"testing"
)

// The framework's test factories, and the list the generator reads.
//
// `grit generate resource X --fields "...:belongs_to:User"` emits
// `factory.CreateUser(t, db)`. Factories are generated per resource, so User's
// and Upload's belonged to nothing and were never written: every project with
// such a resource had a handlers test package that refused to build, with
//
//	vet: internal/handlers/note_test.go:62:18: undefined: factory.CreateUser
//
// and that is what kept CI red from v3.464.0 to v3.487.0. `go build` is clean,
// because it is a _test.go file; `go test` runs a subset of vet as an error
// first, so the package does not compile at all.
//
// Three things have to agree: the files the scaffold writes,
// FrameworkFactoryModels, and the names the generator emits calls to. A test on
// any one of them passes while the set is broken.

// TestFrameworkFactoriesMatchTheList reads both sides of the list.
func TestFrameworkFactoriesMatchTheList(t *testing.T) {
	sources := map[string]string{
		"User":   apiUserFactoryGo(),
		"Upload": apiUploadFactoryGo(),
	}

	for _, model := range FrameworkFactoryModels {
		src, ok := sources[model]
		if !ok {
			t.Errorf("FrameworkFactoryModels names %s and no factory is written for it: "+
				"the generator will emit factory.Create%s and the project will not compile",
				model, model)
			continue
		}
		for _, want := range []string{
			"func New" + model + "(",
			"func Create" + model + "(",
			"func Create" + model + "s(",
		} {
			if !strings.Contains(src, want) {
				t.Errorf("the %s factory has no %s", model, want)
			}
		}
	}

	// And the other direction: a factory written but not listed is one the
	// generator will not count on, so a belongs_to to it goes unwarned.
	for model := range sources {
		var listed bool
		for _, m := range FrameworkFactoryModels {
			if m == model {
				listed = true
			}
		}
		if !listed {
			t.Errorf("a factory is written for %s and FrameworkFactoryModels does not name it", model)
		}
	}
}

// TestFrameworkFactoriesAreDelivered: in the framework-owned map, so they reach
// a project that already exists. Written once at scaffold time, the fix would
// only ever have reached new projects.
func TestFrameworkFactoriesAreDelivered(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "demo", Architecture: ArchAPI}

	if err := writeFrameworkOwnedFiles(root, opts); err != nil {
		t.Fatalf("writeFrameworkOwnedFiles: %v", err)
	}

	for _, name := range []string{"factory.go", "user.go", "upload.go"} {
		path := filepath.Join(opts.APIRoot(root), "internal", "factory", name)
		if !fileExists(path) {
			t.Errorf("internal/factory/%s was not written, so `grit upgrade` does not deliver it", name)
		}
	}
}

// TestFactoryCounterHasOneWriter: internal/generate writes the counter file
// too, for a project old enough not to have it. It reads this text to do it, so
// the two copies cannot drift.
func TestFactoryCounterHasOneWriter(t *testing.T) {
	if FactoryCounterGo() != apiFactoryCounterGo() {
		t.Error("FactoryCounterGo is not the text the scaffold writes")
	}
	if !strings.Contains(FactoryCounterGo(), "func next() int {") {
		t.Error("the counter file has no next(), which every generated factory calls")
	}
	if !strings.Contains(FactoryCounterGo(), "package factory") {
		t.Error("the counter file is not in package factory")
	}
}
