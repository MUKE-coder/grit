package scaffold

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/aymanbagabas/go-udiff"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// The upgrade guard is what stops `grit upgrade` overwriting files you have
// edited.
//
// It lives at writeFile rather than at the callers because there is no single
// caller. Upgrade regenerates the web app through writeWebFiles, the admin
// through upgradeAdminFiles, docs through writeDocsFiles, and root config
// through writeUpgradeFiles, and each of those fans out to dozens of template
// functions. One check at the one place every byte passes through is both
// smaller and harder to leave a hole in.
//
// Like the manifest recorder it is process-wide mode rather than a value
// threaded through, and off unless Upgrade turns it on.
var guard struct {
	mu       sync.Mutex
	active   bool
	force    bool
	root     string
	manifest *manifest.Manifest
	// tracked is whether the project has a manifest at all, which is a
	// different question from whether a given file has an entry in it. See
	// manifest.Exists and guardAllows.
	tracked bool
	written int
	skipped []SkippedFile
	// unowned are files Grit has no record of writing, which differ from what
	// it would write, in a project that does have a manifest. Reported apart
	// from skipped because the sentence is different: you did not edit these,
	// they were never Grit's.
	unowned []SkippedFile
	// deleted are files Grit wrote and the project has since removed. Not put
	// back, and listed so that is a decision rather than a silence.
	deleted []string
}

// SkippedFile is one file left alone because it had been edited since Grit
// wrote it, along with what Grit would have written instead.
type SkippedFile struct {
	Rel      string
	Current  string
	Proposed string
}

// Diff renders what the upgrade would have changed, as a unified diff.
func (s SkippedFile) Diff() string {
	return udiff.Unified("yours/"+s.Rel, "grit/"+s.Rel, s.Current, s.Proposed)
}

// startGuard puts writeFile into protective mode for the duration of an
// upgrade. force keeps the old behaviour of overwriting everything.
func startGuard(root string, force bool) error {
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	guard.mu.Lock()
	defer guard.mu.Unlock()
	guard.active = true
	guard.force = force
	guard.root = root
	guard.manifest = m
	guard.tracked = manifest.Exists(root)
	guard.written = 0
	guard.skipped = nil
	guard.unowned = nil
	guard.deleted = nil
	return nil
}

// stopGuard leaves protective mode and reports what happened.
func stopGuard() (written int, skipped []SkippedFile) {
	written, skipped, _, _ = stopGuardFull()
	return written, skipped
}

// stopGuardFull is the same with the two kinds of held-back file the upgrade
// report words differently.
func stopGuardFull() (written int, skipped, unowned []SkippedFile, deleted []string) {
	guard.mu.Lock()
	defer guard.mu.Unlock()
	written, skipped = guard.written, guard.skipped
	unowned, deleted = guard.unowned, guard.deleted
	guard.active = false
	guard.manifest = nil
	guard.skipped = nil
	guard.unowned = nil
	guard.deleted = nil
	return written, skipped, unowned, deleted
}

// guardAdopt tells the guard a file is Grit's after all, as it stands now.
// The guard reads the manifest once, when the upgrade starts, so a file
// adopted during the upgrade would otherwise still be held back and reported
// as edited by you in the same run that adopted it.
func guardAdopt(path, content string) {
	guard.mu.Lock()
	defer guard.mu.Unlock()
	if !guard.active || guard.manifest == nil {
		return
	}
	rel, inside := manifest.Rel(guard.root, path)
	if !inside {
		return
	}
	if entry, tracked := guard.manifest.Files[rel]; tracked {
		guard.manifest.Record(rel, entry.Generator, "", content)
	}
}

// guardAllows decides whether a write may proceed, and records the decision.
//
// Four answers, by what the manifest says about the file:
//
//   - Grit wrote it and it has changed since: held back, as it has been since
//     v3.147.0. The upgrade reports it as yours.
//   - Grit wrote it and it is gone, along with the directory it lived in: the
//     project removed a feature, so it is NOT put back. "Missing" and "never
//     written" used to be the same answer here, which is how six deleted Expo
//     screens returned on every upgrade (grit#125).
//   - Grit wrote it and it is gone but its directory is still there: written,
//     because deleting one framework file out of a directory full of them is
//     how people ask for the stock version back, and has been the documented
//     way since v3.147.0. The directory is the discriminator: removing one is
//     a bigger statement than removing a file, and the only one of the two
//     that reads as "I do not want this".
//   - Grit has no record of it, in a project that HAS a manifest: the file is
//     not Grit's. If it differs from what Grit would write, held back and
//     reported. An Expo app copied in from another project is the whole of
//     this case, and overwriting it was the serious half of grit#125.
//   - Grit has no record of it and the project has no manifest: written,
//     exactly as before v3.147.0. With no provenance anywhere there is no
//     basis for claiming any file was edited, and refusing to update an
//     entire pre-manifest project would be worse than the problem. This is
//     why the Untracked case could not simply be protected.
func guardAllows(path, proposed string) bool {
	guard.mu.Lock()
	defer guard.mu.Unlock()
	if !guard.active {
		return true
	}
	if guard.force {
		guard.written++
		return true
	}

	rel, inside := manifest.Rel(guard.root, path)
	if !inside {
		return true
	}

	switch guard.manifest.StatusOf(guard.root, rel) {
	case manifest.Missing:
		// Grit wrote it; it is not there now. Whether that was a decision or a
		// request for the stock file back is read from the directory.
		if dirExists(filepath.Dir(path)) {
			guard.written++
			return true
		}
		guard.deleted = append(guard.deleted, rel)
		return false

	case manifest.Untracked:
		if !guard.tracked {
			guard.written++
			return true
		}
		current, err := os.ReadFile(path)
		if err != nil {
			// Not there, and never ours: nothing to lose by writing it.
			guard.written++
			return true
		}
		if manifest.Hash(string(current)) == manifest.Hash(proposed) {
			// Already what Grit would write. Adopt it rather than report it:
			// the bytes are identical, so there is nothing to decide, and a
			// record means the next upgrade knows whose it is.
			guard.manifest.Record(rel, "upgrade", "", proposed)
			guard.written++
			return true
		}
		guard.unowned = append(guard.unowned, SkippedFile{
			Rel:      rel,
			Current:  string(current),
			Proposed: proposed,
		})
		return false

	case manifest.Modified:
		current, err := os.ReadFile(path)
		if err != nil {
			return true
		}
		// An edit that happens to leave the file identical to the new template
		// is not a conflict, and reporting it as one would be noise.
		if manifest.Hash(string(current)) == manifest.Hash(proposed) {
			return true
		}
		guard.skipped = append(guard.skipped, SkippedFile{
			Rel:      rel,
			Current:  string(current),
			Proposed: proposed,
		})
		return false

	default: // Unchanged
		guard.written++
		return true
	}
}

// guardedWrite is writeFile's protected form: it asks the guard first and
// reports whether it wrote.
func guardedWrite(path, final string) (bool, error) {
	if !guardAllows(path, final) {
		return false, nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return false, err
	}
	if err := os.WriteFile(path, []byte(final), 0644); err != nil {
		return false, err
	}
	manifest.Note(path, final)
	return true, nil
}
