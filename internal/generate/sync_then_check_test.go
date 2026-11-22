package generate

import (
	"os"
	"path/filepath"
	"testing"
)

// Sync writes it, and the check has to agree it is right.
//
// They are two readers of one rule: which models get a shared file, where it
// goes, and what is in it. The content comes from buildTSType and
// buildZodSchema in both, and the path and the skip from skipSharedType and
// sharedTypePath, so there is one answer by construction. This is what holds
// them to it: run the real Sync over a project and then the real check, and
// require it to find nothing.
//
// Without it the pair can come apart in the way that is worst: the check
// compares a path sync does not write, finds the file absent or matching an
// older rule, and reports a clean project forever.
func TestSyncThenCheckFindsNothing(t *testing.T) {
	root := setupMinimalProject(t, "myapp/apps/api")
	writeCheckModel(t, root, "\tColour    string    `json:\"colour\"`\n\tRank      int       `json:\"rank\"`\n")

	// Sync() finds the project from the working directory.
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("working directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(wd); err != nil {
			t.Fatalf("restoring the working directory: %v", err)
		}
	})
	if err := os.Chdir(root); err != nil {
		t.Fatalf("entering the project: %v", err)
	}

	if err := Sync(); err != nil {
		t.Fatalf("Sync: %v", err)
	}

	// Both files exist, so the check is comparing something.
	for _, path := range []string{
		sharedTypePath(sharedRootFrom(root), "Widget"),
		sharedSchemaPath(sharedRootFrom(root), "Widget"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("sync did not write %s: %v", filepath.Base(path), err)
		}
	}

	drift, err := syncCheckIn(root)
	if err != nil {
		t.Fatalf("syncCheckIn: %v", err)
	}
	for _, d := range drift {
		line, want, got := d.FirstDifference()
		t.Errorf("the check disagrees with what sync just wrote: %s (%s)\n"+
			"  line %d\n  on disk:   %s\n  the model: %s",
			d.Rel, d.Kind, line, got, want)
	}
}
