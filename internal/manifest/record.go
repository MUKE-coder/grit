package manifest

import (
	"os"
	"strings"
	"sync"
)

// Recording is process-wide state, and deliberately so.
//
// Every file Grit writes goes through one of two helpers, but those helpers are
// called from roughly 150 template functions that thread no context of any
// kind. Threading a recorder through all of them would be a change larger than
// the feature. The trade is that recording is a mode the process is in, not a
// value it passes around, so it is off unless a command turns it on and a test
// that does not opt in is unaffected.
var recording struct {
	mu        sync.Mutex
	manifest  *Manifest
	root      string
	grit      string
	generator string
	active    bool

	// edited is every tracked file that was already hand-edited when this
	// session began, by project-relative key.
	//
	// Taken once, at the start, because that is the only moment the answer is
	// knowable: Refresh runs after the write, when the bytes on disk are the
	// ones Grit just put there and an edit is indistinguishable from an
	// injection. Without it, injecting into a file somebody had edited
	// re-hashed it as pristine and the next upgrade overwrote their work
	// (review 2.6).
	edited map[string]bool
}

// Start begins recording every write under root as having been made by
// generator, and returns the function that ends it. It loads any existing
// manifest first, so a resource generated today does not erase the provenance
// of a project scaffolded last month.
//
// Nested calls relabel rather than restart. Generator.Run calls itself for an
// inline child resource, and a nested Start that reloaded from disk would
// throw away everything the parent had recorded and not yet saved.
//
// Usage: defer the returned release, and let it report the save error.
func Start(root, gritVersion, generator string) (release func() error, err error) {
	recording.mu.Lock()
	if recording.active {
		previous := recording.generator
		recording.generator = generator
		recording.mu.Unlock()
		return func() error {
			recording.mu.Lock()
			defer recording.mu.Unlock()
			recording.generator = previous
			return nil
		}, nil
	}
	recording.mu.Unlock()

	m, loadErr := Load(root)
	if loadErr != nil {
		return func() error { return nil }, loadErr
	}

	// Before anything is written, while "has somebody edited this" still has
	// an answer. See recording.edited.
	edited := editedFiles(root, m)

	recording.mu.Lock()
	recording.manifest = m
	recording.root = root
	recording.grit = gritVersion
	recording.generator = generator
	recording.active = true
	recording.edited = edited
	recording.mu.Unlock()

	return Stop, nil
}

// Generator relabels subsequent writes. Used when one command runs several
// generators in sequence and each batch of files has a different owner.
func Generator(name string) {
	recording.mu.Lock()
	defer recording.mu.Unlock()
	if recording.active {
		recording.generator = name
	}
}

// Note records that content was written to absPath. It is called from the file
// writers and does nothing at all unless Start has been called, which is what
// keeps it out of the way of tests and of library use.
func Note(absPath, content string) {
	recording.mu.Lock()
	defer recording.mu.Unlock()
	if !recording.active {
		return
	}
	rel, inside := Rel(recording.root, absPath)
	if !inside {
		return
	}
	recording.manifest.Record(rel, recording.generator, recording.grit, content)
}

// Refresh re-reads a file from disk and updates its hash, keeping whichever
// generator originally claimed it. This is for injection: adding a route to
// routes.go or a resource to the registry is an edit Grit made, and without
// this every injected file would read as hand-edited from then on.
//
// Files Grit does not already track are left alone. An injection into
// something the user wrote does not make it Grit's.
func Refresh(absPath string) {
	recording.mu.Lock()
	defer recording.mu.Unlock()
	if !recording.active {
		return
	}
	rel, inside := Rel(recording.root, absPath)
	if !inside {
		return
	}
	existing, tracked := recording.manifest.Files[rel]
	if !tracked {
		return
	}
	// Already edited when this command started, so its hash stays as it was
	// and the file keeps reading as modified. The injection still happened and
	// the project still works; what does not happen is the next upgrade
	// deciding the file is Grit's and overwriting it (review 2.6).
	if recording.edited[Key(rel)] {
		return
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		return
	}
	recording.manifest.Record(rel, existing.Generator, recording.grit, string(data))
}

// WasEditedBeforeThisCommand reports whether a tracked file had already been
// hand-edited when the current command started.
//
// For a caller that wants to say so rather than only avoid making it worse: an
// injection into somebody's edited file is worth a line of output, and after
// the write there is no other way to tell.
func WasEditedBeforeThisCommand(absPath string) bool {
	recording.mu.Lock()
	defer recording.mu.Unlock()
	if !recording.active {
		return false
	}
	rel, inside := Rel(recording.root, absPath)
	if !inside {
		return false
	}
	return recording.edited[Key(rel)]
}

// editedFiles is every tracked file whose bytes on disk are not the bytes Grit
// last wrote there.
//
// Missing and untracked files are not in it: a file Grit wrote and somebody
// deleted is not an edit to preserve, and a file Grit never wrote is not
// Refresh's business anyway.
func editedFiles(root string, m *Manifest) map[string]bool {
	edited := make(map[string]bool)
	for rel := range m.Files {
		if m.StatusOf(root, rel) == Modified {
			edited[Key(rel)] = true
		}
	}
	return edited
}

// IsUnchanged reports whether a tracked file still holds exactly what Grit last
// wrote. For an edit that should adopt its result only when nobody else has
// touched the file: check before editing, then Refresh after. False when not
// recording, and for files Grit does not track.
func IsUnchanged(absPath string) bool {
	recording.mu.Lock()
	defer recording.mu.Unlock()
	if !recording.active {
		return false
	}
	rel, inside := Rel(recording.root, absPath)
	if !inside {
		return false
	}
	return recording.manifest.StatusOf(recording.root, rel) == Unchanged
}

// StatusOfPath reports what has happened to an absolute path since Grit last
// wrote it: Untracked when Grit never wrote it or nothing is recording,
// Missing when it was written and has been deleted, and otherwise Unchanged
// or Modified.
//
// IsUnchanged answers one of those questions and collapses the rest into
// false, which is right for "may I edit this" and wrong for "has this been
// deleted on purpose".
func StatusOfPath(absPath string) Status {
	recording.mu.Lock()
	defer recording.mu.Unlock()
	if !recording.active {
		return Untracked
	}
	rel, inside := Rel(recording.root, absPath)
	if !inside {
		return Untracked
	}
	return recording.manifest.StatusOf(recording.root, rel)
}

// Drop removes a path from the recording, for a generator that deletes a file
// it previously wrote.
// DropTree forgets every recorded file under a directory, for the case where
// Grit removes a whole one.
//
// Needed for the same reason as Drop: the upgrade guard reads a recorded file
// that is absent, in a directory that is also absent, as a feature the project
// removed, and refuses to re-create it. A directory Grit removed itself is not
// that.
func DropTree(absDir string) {
	recording.mu.Lock()
	defer recording.mu.Unlock()
	if !recording.active || recording.manifest == nil {
		return
	}
	rel, inside := Rel(recording.root, absDir)
	if !inside {
		return
	}
	prefix := Key(rel) + "/"
	for key := range recording.manifest.Files {
		if strings.HasPrefix(key, prefix) {
			delete(recording.manifest.Files, key)
		}
	}
}

func Drop(absPath string) {
	recording.mu.Lock()
	defer recording.mu.Unlock()
	if !recording.active {
		return
	}
	if rel, inside := Rel(recording.root, absPath); inside {
		recording.manifest.Forget(rel)
	}
}

// Stop writes the manifest and leaves recording mode. Safe to call when Start
// was never called, so a command can defer it unconditionally.
func Stop() error {
	recording.mu.Lock()
	m, root, grit, active := recording.manifest, recording.root, recording.grit, recording.active
	recording.active = false
	recording.manifest = nil
	// The snapshot belongs to the session that took it. Left behind, the next
	// Start in the same process (a test, or a command that runs two
	// generators) would read one project's edits against another's files.
	recording.edited = nil
	recording.mu.Unlock()

	if !active || m == nil {
		return nil
	}
	return m.Save(root, grit)
}

// Active reports whether writes are currently being recorded.
func Active() bool {
	recording.mu.Lock()
	defer recording.mu.Unlock()
	return recording.active
}
