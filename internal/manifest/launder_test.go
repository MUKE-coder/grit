package manifest

import (
	"path/filepath"
	"testing"
)

// Review 2.6: an injection laundered a hand-edit into "framework-owned".
//
// Refresh re-hashes a file after Grit injects into it, which is right: adding
// a route to routes.go is an edit Grit made, and without it every injected
// file would read as hand-edited from then on. It was wrong when the file had
// already been edited: the re-hash recorded somebody's own work as Grit's, and
// the next upgrade overwrote it without a word.
//
// Two files, one edited and one not, through the same injection, so the test
// shows both halves of the rule rather than only the new one.
func TestAnInjectionDoesNotLaunderAHandEdit(t *testing.T) {
	root := t.TempDir()

	mine := filepath.Join(root, "routes.go")
	theirs := filepath.Join(root, "registry.ts")
	write(t, mine, "// grit wrote this\n")
	write(t, theirs, "// grit wrote this too\n")

	// Grit records both as its own.
	release, err := Start(root, "3.0.0", "scaffold")
	if err != nil {
		t.Fatal(err)
	}
	Note(mine, "// grit wrote this\n")
	Note(theirs, "// grit wrote this too\n")
	if err := release(); err != nil {
		t.Fatal(err)
	}

	// Somebody edits one of them.
	write(t, mine, "// grit wrote this\n// and I added this line\n")

	// A later command injects into both, as grit generate resource does.
	release, err = Start(root, "3.0.0", "resource:Product")
	if err != nil {
		t.Fatal(err)
	}
	if !WasEditedBeforeThisCommand(mine) {
		t.Error("the edited file is not reported as edited, so nothing can protect it")
	}
	if WasEditedBeforeThisCommand(theirs) {
		t.Error("an untouched file is reported as edited")
	}
	write(t, mine, "// grit wrote this\n// and I added this line\nproducts()\n")
	Refresh(mine)
	write(t, theirs, "// grit wrote this too\nproducts()\n")
	Refresh(theirs)
	if err := release(); err != nil {
		t.Fatal(err)
	}

	// The edited file still reads as edited, so an upgrade leaves it alone.
	m, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if got := m.StatusOf(root, "routes.go"); got != Modified {
		t.Errorf("the edited file reads as %s: an upgrade will overwrite somebody's work", got)
	}
	// And the untouched one still reads as Grit's, so it keeps receiving fixes.
	if got := m.StatusOf(root, "registry.ts"); got != Unchanged {
		t.Errorf("the untouched file reads as %s: it will stop receiving framework fixes", got)
	}
}

// The snapshot belongs to the session that took it. Two commands in one
// process, which is a test or a command running two generators, must not read
// one session's edits against another's files.
func TestTheEditSnapshotDoesNotOutliveItsSession(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "a.go")
	write(t, path, "one\n")

	release, err := Start(root, "3.0.0", "scaffold")
	if err != nil {
		t.Fatal(err)
	}
	Note(path, "one\n")
	if err := release(); err != nil {
		t.Fatal(err)
	}

	write(t, path, "one\nedited\n")

	// Session two sees the edit.
	release, err = Start(root, "3.0.0", "scaffold")
	if err != nil {
		t.Fatal(err)
	}
	if !WasEditedBeforeThisCommand(path) {
		t.Fatal("session two did not see the edit")
	}
	// Grit adopts the file deliberately, as --force does.
	Note(path, "one\nedited\n")
	if err := release(); err != nil {
		t.Fatal(err)
	}

	// Session three must not still think it is edited: it is now exactly what
	// Grit recorded.
	release, err = Start(root, "3.0.0", "scaffold")
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	if WasEditedBeforeThisCommand(path) {
		t.Error("the snapshot from an earlier session is still being consulted")
	}
}

// Nothing happens outside a recording session, which is what keeps the
// manifest out of the way of library use and of tests that do not opt in.
func TestWasEditedIsFalseWhenNotRecording(t *testing.T) {
	if WasEditedBeforeThisCommand(filepath.Join(t.TempDir(), "x.go")) {
		t.Error("reported an edit with no session running")
	}
}
