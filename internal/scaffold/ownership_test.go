package scaffold

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// OwnershipOf mirrors guardAllows. A test for either one alone would pass
// while the pair disagreed, which is the shape of bug this project keeps
// finding, so this drives both over the same files and compares the answers.
//
// Only the cases where the verdict does not depend on the new template's bytes
// can be compared this way, which is exactly what Ownership.Certain marks. The
// uncertain ones are checked separately below.
func TestOwnershipAgreesWithTheUpgradeGuard(t *testing.T) {
	root := t.TempDir()

	write := func(rel, content string) string {
		abs := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		return abs
	}

	// A project with a manifest, and four files in the four states.
	m, err := manifest.Load(root)
	if err != nil {
		t.Fatal(err)
	}

	pristine := write("apps/admin/pristine.ts", "// as grit wrote it\n")
	edited := write("apps/admin/edited.ts", "// as grit wrote it\n")
	deleted := write("apps/admin/deleted.ts", "// as grit wrote it\n")
	gone := write("apps/removed/gone.ts", "// as grit wrote it\n")
	mine := write("apps/admin/mine.ts", "// mine\n")

	for _, rel := range []string{
		"apps/admin/pristine.ts", "apps/admin/edited.ts",
		"apps/admin/deleted.ts", "apps/removed/gone.ts",
	} {
		m.Record(rel, "test", "0.0.0", "// as grit wrote it\n")
	}
	if err := m.Save(root, "0.0.0"); err != nil {
		t.Fatal(err)
	}

	// Now make the states real.
	if err := os.WriteFile(edited, []byte("// mine now\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(deleted); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(filepath.Join(root, "apps", "removed")); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name string
		path string
		// what a template would write there, which is what the guard sees
		proposed string
	}{
		{"pristine", pristine, "// a newer grit wrote this\n"},
		{"deleted, directory still there", deleted, "// a newer grit wrote this\n"},
		{"deleted with its directory", gone, "// a newer grit wrote this\n"},
		{"edited", edited, "// a newer grit wrote this\n"},
		{"never grit's", mine, "// a newer grit wrote this\n"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o, err := OwnershipOf(root, tc.path)
			if err != nil {
				t.Fatalf("OwnershipOf: %v", err)
			}
			if o.Headline == "" || o.Detail == "" {
				t.Errorf("no answer for %s: %+v", tc.path, o)
			}

			// Run the real guard over the same file.
			if err := startGuard(root, false); err != nil {
				t.Fatal(err)
			}
			allowed := guardAllows(tc.path, tc.proposed)
			stopGuardFull()

			if o.Certain && allowed != o.WillReplace {
				t.Errorf("%s: grit own says WillReplace=%v, the upgrade guard says %v.\n"+
					"These two have to agree: `grit own` is only useful if it predicts what\n"+
					"the upgrade actually does. Headline was %q.",
					o.Rel, o.WillReplace, allowed, o.Headline)
			}
			if !o.Certain && allowed {
				t.Errorf("%s: grit own says the upgrade cannot be sure to replace it, and the\n"+
					"guard replaced it anyway with content that differs. Headline was %q.",
					o.Rel, o.Headline)
			}
		})
	}
}

// The one case the manifest cannot distinguish by status alone, and the one
// that cost somebody their work: a project with no manifest is not a project
// full of edited files.
func TestOwnershipSaysWhenThereIsNoManifest(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "apps", "admin", "thing.ts")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("// whatever\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	o, err := OwnershipOf(root, path)
	if err != nil {
		t.Fatal(err)
	}
	if o.Tracked {
		t.Error("reported a manifest in a project that has none")
	}
	if !o.WillReplace || !o.Certain {
		t.Errorf("a project with no manifest has every file overwritten by an upgrade; "+
			"got WillReplace=%v Certain=%v", o.WillReplace, o.Certain)
	}

	if err := startGuard(root, false); err != nil {
		t.Fatal(err)
	}
	allowed := guardAllows(path, "// a newer grit wrote this\n")
	stopGuardFull()
	if !allowed {
		t.Error("the guard held back a file in a project with no manifest, so this answer is wrong")
	}
}

// Outside the project is an error rather than a guess.
func TestOwnershipRefusesAPathOutsideTheProject(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "elsewhere.ts")
	if _, err := OwnershipOf(root, outside); err == nil {
		t.Error("answered for a path outside the project instead of saying so")
	}
}
