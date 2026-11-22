package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A model file with one field, written into a minimal project.
func writeCheckModel(t *testing.T, root, extra string) string {
	t.Helper()
	path := filepath.Join(root, "apps", "api", "internal", "models", "widget.go")
	writeTestFile(t, path, `package models

import "time"

type Widget struct {
	ID        string    `+"`"+`gorm:"primarykey;size:36" json:"id"`+"`"+`
	Name      string    `+"`"+`json:"name"`+"`"+`
`+extra+`	CreatedAt time.Time `+"`"+`json:"created_at"`+"`"+`
}
`)
	return path
}

// writeCurrentShared writes what the model implies, so the project starts clean.
func writeCurrentShared(t *testing.T, root, modelPath string) {
	t.Helper()
	structs, err := parseGoStructs(modelPath)
	if err != nil {
		t.Fatalf("parsing the model: %v", err)
	}
	shared := sharedRootFrom(root)
	for _, s := range structs {
		if skipSharedType(s.Name) {
			continue
		}
		writeTestFile(t, sharedTypePath(shared, s.Name), buildTSType(s))
		writeTestFile(t, sharedSchemaPath(shared, s.Name), buildZodSchema(s))
	}
}

func checkProject(t *testing.T) (root, modelPath string) {
	t.Helper()
	root = setupMinimalProject(t, "myapp/apps/api")
	// The shared package has to exist for sync to run at all.
	writeTestFile(t, filepath.Join(root, "packages", "shared", "package.json"),
		`{"name":"@myapp/shared"}`+"\n")
	modelPath = writeCheckModel(t, root, "")
	writeCurrentShared(t, root, modelPath)
	return root, modelPath
}

// A project whose shared files were written from its models passes.
func TestSyncCheckPassesAProjectInStep(t *testing.T) {
	root, _ := checkProject(t)

	drift, err := syncCheckIn(root)
	if err != nil {
		t.Fatalf("syncCheckIn: %v", err)
	}
	if len(drift) != 0 {
		for _, d := range drift {
			t.Errorf("unexpected drift: %s (%s)", d.Rel, d.Kind)
		}
	}
}

// The case the issue is about: a field added to a model, nobody ran sync.
//
// Everything else in a generated project's CI is happy with this. The Go
// compiles, the TypeScript type-checks, and the field is simply absent from
// the client, which shows up as undefined in whichever screen reads it first.
func TestSyncCheckCatchesAFieldAddedWithoutSyncing(t *testing.T) {
	root, _ := checkProject(t)

	// Add a field and do not sync.
	writeCheckModel(t, root, "\tColour    string    `json:\"colour\"`\n")

	drift, err := syncCheckIn(root)
	if err != nil {
		t.Fatalf("syncCheckIn: %v", err)
	}
	if len(drift) != 2 {
		t.Fatalf("got %d drifting files, want the type and the schema", len(drift))
	}
	for _, d := range drift {
		if d.Kind != DriftStale {
			t.Errorf("%s reported as %q, want stale", d.Rel, d.Kind)
		}
		if d.Model != "Widget" {
			t.Errorf("%s blamed on %q, want Widget", d.Rel, d.Model)
		}
		if !strings.Contains(d.Want, "colour") {
			t.Errorf("%s: what the model implies does not mention the new field", d.Rel)
		}
		line, want, _ := d.FirstDifference()
		if line == 0 {
			t.Errorf("%s: reported as stale with no differing line", d.Rel)
		}
		if !strings.Contains(want, "colour") {
			t.Errorf("%s: the first difference reported is line %d, %q, which is "+
				"not the field that was added", d.Rel, line, want)
		}
	}
}

// A model with no shared file is not drift, and that is the decision this
// records rather than an omission.
//
// The literal rule, every model implies a type and a schema, was measured
// against a freshly scaffolded triple project: 78 missing files on a project
// nobody had touched, because packages/shared does not mirror every framework
// model and never did. Nothing on a client validates an OAuthExchangeCode. A
// check that fails on day one is one nobody runs.
//
// What is still covered: every resource grit generate writes, which gets both
// files, and every framework model that has one.
func TestSyncCheckTreatsAnAbsentFileAsADecisionNotDrift(t *testing.T) {
	root, _ := checkProject(t)
	if err := os.Remove(sharedTypePath(sharedRootFrom(root), "Widget")); err != nil {
		t.Fatalf("removing the type: %v", err)
	}

	drift, err := syncCheckIn(root)
	if err != nil {
		t.Fatalf("syncCheckIn: %v", err)
	}
	for _, d := range drift {
		if strings.HasSuffix(filepath.ToSlash(d.Rel), "types/widget.ts") {
			t.Errorf("a model with no shared type was reported as drift: %v", d)
		}
	}

	// And the schema, which is still there, is still checked. Otherwise this
	// would pass for the wrong reason.
	writeCheckModel(t, root, "\tColour    string    `json:\"colour\"`\n")
	drift, err = syncCheckIn(root)
	if err != nil {
		t.Fatalf("syncCheckIn: %v", err)
	}
	if len(drift) != 1 {
		t.Fatalf("got %d drifting files, want only the schema that is there", len(drift))
	}
	if !strings.HasSuffix(filepath.ToSlash(drift[0].Rel), "schemas/widget.ts") {
		t.Errorf("the drift reported is %q, want the schema", drift[0].Rel)
	}
}

// User is skipped by sync because its schemas are hand-written, and a check
// that did not skip it would fail every project on its first run.
func TestSyncCheckSkipsTheSameModelsSyncDoes(t *testing.T) {
	root, _ := checkProject(t)

	// setupMinimalProject writes a User model and no user.ts.
	drift, err := syncCheckIn(root)
	if err != nil {
		t.Fatalf("syncCheckIn: %v", err)
	}
	for _, d := range drift {
		if d.Model == "User" {
			t.Errorf("the check wants to write %s for User, which sync never "+
				"writes: the two no longer skip the same models", d.Rel)
		}
	}
}

// CRLF is not drift.
//
// git checks these files out with CRLF on Windows, and a byte comparison would
// report every shared file as stale on every Windows machine while passing in
// CI, which is the worst shape a check can have.
func TestSyncCheckIgnoresLineEndings(t *testing.T) {
	root, _ := checkProject(t)

	path := sharedTypePath(sharedRootFrom(root), "Widget")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading the type: %v", err)
	}
	crlf := strings.ReplaceAll(strings.ReplaceAll(string(raw), "\r\n", "\n"), "\n", "\r\n")
	if err := os.WriteFile(path, []byte(crlf), 0o644); err != nil {
		t.Fatalf("rewriting with CRLF: %v", err)
	}

	drift, err := syncCheckIn(root)
	if err != nil {
		t.Fatalf("syncCheckIn: %v", err)
	}
	if len(drift) != 0 {
		t.Errorf("a CRLF checkout reported as drift: %v", drift)
	}
}

// A project with no shared package is not a project this can check, and
// saying so is better than reporting every model as missing a file.
func TestSyncCheckRefusesAProjectWithNoSharedPackage(t *testing.T) {
	root := setupMinimalProject(t, "myapp/apps/api")
	writeCheckModel(t, root, "")
	// The fixture writes one, because almost every other test needs it.
	if err := os.RemoveAll(filepath.Join(root, "packages", "shared")); err != nil {
		t.Fatalf("removing the shared package: %v", err)
	}

	if _, err := syncCheckIn(root); err == nil {
		t.Error("a project with no packages/shared passed the check")
	}
}
