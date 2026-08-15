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
	written  int
	skipped  []SkippedFile
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
	guard.written = 0
	guard.skipped = nil
	return nil
}

// stopGuard leaves protective mode and reports what happened.
func stopGuard() (written int, skipped []SkippedFile) {
	guard.mu.Lock()
	defer guard.mu.Unlock()
	written, skipped = guard.written, guard.skipped
	guard.active = false
	guard.manifest = nil
	guard.skipped = nil
	return written, skipped
}

// guardAllows decides whether a write may proceed, and records the decision.
//
// The rule is narrow on purpose: only a file Grit is recorded as having
// written, whose bytes have changed since, and which would change again, is
// held back. Everything else is written as before.
//
// A project with no manifest has every file Untracked, so it upgrades exactly
// as it did before v3.147.0. That is the honest outcome: with no provenance
// there is no basis for claiming a file was edited, and silently refusing to
// update an entire pre-manifest project would be worse than the problem.
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
	if guard.manifest.StatusOf(guard.root, rel) != manifest.Modified {
		guard.written++
		return true
	}

	current, err := os.ReadFile(path)
	if err != nil {
		return true
	}
	// An edit that happens to leave the file identical to the new template is
	// not a conflict, and reporting it as one would be noise.
	if manifest.Hash(string(current)) == manifest.Hash(proposed) {
		return true
	}

	guard.skipped = append(guard.skipped, SkippedFile{
		Rel:      rel,
		Current:  string(current),
		Proposed: proposed,
	})
	return false
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
