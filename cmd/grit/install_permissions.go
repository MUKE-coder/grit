package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Whether this install can replace itself.
//
// `grit update` rewrites the binary in place, which needs write access to the
// directory it lives in, and on Windows needs to move the running .exe aside
// first. Both fail in ways nobody can act on: an unwritable directory surfaces
// as "Access is denied" from a rename several steps later, naming neither the
// directory nor the reason.
//
// So the check happens up front, once, and says what is wrong and where. The
// installers run the same check before they write anything, which is the point
// at which somebody can still choose a different directory.

// checkInstallDirWritable reports whether dir can take a new binary.
//
// It writes and removes a probe file rather than reading permission bits,
// because the bits are not the whole story: a directory can be writable and
// still refuse, under Windows ACLs, a read-only volume, or a container mount.
// The only honest test is to try.
func checkInstallDirWritable(dir string) error {
	probe, err := os.CreateTemp(dir, ".grit-write-probe-*")
	if err != nil {
		return fmt.Errorf("%s is not writable, so grit cannot update itself in place: %w\n\n%s",
			dir, err, installDirAdvice(dir))
	}
	name := probe.Name()
	_ = probe.Close()
	if err := os.Remove(name); err != nil {
		return fmt.Errorf("%s allows writes but not deletes, so grit cannot replace its binary: %w\n\n%s",
			dir, err, installDirAdvice(dir))
	}
	return nil
}

// installDirAdvice says what to do about a directory grit cannot write.
//
// The advice is to move the install, not to run as administrator. A CLI that
// needs elevation to update is a CLI that stops being updated.
func installDirAdvice(dir string) string {
	home, _ := os.UserHomeDir()
	suggested := filepath.Join(home, ".grit", "bin")

	var b strings.Builder
	if runtime.GOOS == "windows" {
		if isSystemDir(dir) {
			b.WriteString("  That is a system directory, so every update would need an administrator.\n")
		}
		fmt.Fprintf(&b, "  Reinstall somewhere you own, which needs no elevation:\n\n")
		fmt.Fprintf(&b, "    $env:GRIT_INSTALL_DIR = \"%s\"\n", suggested)
		b.WriteString("    irm https://gritframework.dev/install.ps1 | iex\n")
		return b.String()
	}
	if isSystemDir(dir) {
		b.WriteString("  That is a system directory, so every update would need sudo.\n")
	}
	fmt.Fprintf(&b, "  Reinstall somewhere you own, which needs no sudo:\n\n")
	fmt.Fprintf(&b, "    GRIT_INSTALL_DIR=%s curl -fsSL https://gritframework.dev/install.sh | sh\n", suggested)
	return b.String()
}

// isSystemDir reports whether dir is one of the locations that needs elevation
// to write. Not exhaustive, and does not need to be: it only decides whether
// to add a sentence explaining why the directory is refusing.
func isSystemDir(dir string) bool {
	lower := strings.ToLower(filepath.Clean(dir))
	for _, prefix := range []string{
		`c:\program files`, `c:\program files (x86)`, `c:\windows`,
		"/usr/bin", "/usr/local/bin", "/usr/sbin", "/bin", "/sbin", "/opt",
	} {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	return false
}

// asideName picks a name to move the running binary to on Windows.
//
// The obvious name is binary + ".old", and it is why `grit update` could work
// once and then fail forever: when something still holds a handle on the old
// file, Windows marks it delete-pending rather than removing it, so the name
// stays reserved and the next rename onto it fails with "Access is denied"
// with nothing to say why. A name carrying this process's id cannot collide
// with a file another process is holding.
func asideName(binPath string) string {
	if err := os.Remove(binPath + ".old"); err == nil || os.IsNotExist(err) {
		return binPath + ".old"
	}
	return fmt.Sprintf("%s.old-%d", binPath, os.Getpid())
}

// sweepAsideBinaries removes the leftovers earlier updates could not.
//
// Best effort: a file Windows is still holding stays, and the next update will
// pick a fresh name rather than trip over it.
func sweepAsideBinaries(binPath string) {
	matches, err := filepath.Glob(binPath + ".old*")
	if err != nil {
		return
	}
	for _, m := range matches {
		_ = os.Remove(m)
	}
}
