package manifest

import (
	"os"
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

	recording.mu.Lock()
	recording.manifest = m
	recording.root = root
	recording.grit = gritVersion
	recording.generator = generator
	recording.active = true
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
	data, err := os.ReadFile(absPath)
	if err != nil {
		return
	}
	recording.manifest.Record(rel, existing.Generator, recording.grit, string(data))
}

// Drop removes a path from the recording, for a generator that deletes a file
// it previously wrote.
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
