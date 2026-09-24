package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestAWritableDirPasses(t *testing.T) {
	if err := checkInstallDirWritable(t.TempDir()); err != nil {
		t.Errorf("a temp directory should be writable: %v", err)
	}
}

func TestAMissingDirFailsWithSomethingToDoAboutIt(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "no", "such", "place")
	err := checkInstallDirWritable(missing)
	if err == nil {
		t.Fatal("a directory that does not exist was reported as writable")
	}
	// The whole point of checking early is that the message can name the
	// directory and say what to do. "Access is denied" from a rename does
	// neither.
	msg := err.Error()
	if !strings.Contains(msg, missing) {
		t.Errorf("the error does not name the directory: %s", msg)
	}
	if !strings.Contains(msg, "GRIT_INSTALL_DIR") {
		t.Errorf("the error does not say how to install somewhere else: %s", msg)
	}
	if strings.Contains(strings.ToLower(msg), "administrator") ||
		strings.Contains(strings.ToLower(msg), "sudo") {
		// Only a system directory should suggest elevation, and only to
		// explain why it is refusing. A temp path is not one.
		t.Errorf("a non-system directory should not mention elevation: %s", msg)
	}
}

func TestASystemDirectoryExplainsWhyItRefuses(t *testing.T) {
	dir := "/usr/local/bin"
	if runtime.GOOS == "windows" {
		dir = `C:\Program Files\grit`
	}
	if !isSystemDir(dir) {
		t.Fatalf("%s should be recognised as a system directory", dir)
	}
	advice := installDirAdvice(dir)
	if !strings.Contains(strings.ToLower(advice), "administrator") &&
		!strings.Contains(strings.ToLower(advice), "sudo") {
		t.Errorf("the advice does not explain why a system directory refuses: %s", advice)
	}
	// And it still points somewhere that needs neither.
	if !strings.Contains(advice, "GRIT_INSTALL_DIR") {
		t.Errorf("the advice does not offer a directory the user owns: %s", advice)
	}
}

// The bug this was written for: grit update worked once, then failed forever.
//
// Windows marks a file delete-pending when something still holds a handle, so
// the name stays reserved and renaming onto it fails with "Access is denied".
// A name that cannot collide is the fix.
func TestTheAsideNameAvoidsALockedLeftover(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "grit.exe")
	if err := os.WriteFile(bin, []byte("binary"), 0o755); err != nil {
		t.Fatal(err)
	}

	// No leftover: the plain name is fine, and nothing is left behind.
	if got := asideName(bin); got != bin+".old" {
		t.Errorf("with no leftover the aside name should be %s.old, got %s", bin, got)
	}

	// A removable leftover is removed and the plain name reused.
	if err := os.WriteFile(bin+".old", []byte("stale"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := asideName(bin); got != bin+".old" {
		t.Errorf("a removable leftover should be cleared, got %s", got)
	}
	if _, err := os.Stat(bin + ".old"); !os.IsNotExist(err) {
		t.Error("the removable leftover was not cleared")
	}
}

func TestTheSweepClearsEveryLeftover(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "grit.exe")
	for _, name := range []string{bin + ".old", bin + ".old-123", bin + ".old-456"} {
		if err := os.WriteFile(name, []byte("stale"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	sweepAsideBinaries(bin)
	left, _ := filepath.Glob(bin + ".old*")
	if len(left) != 0 {
		t.Errorf("leftovers survived the sweep: %v", left)
	}
}
