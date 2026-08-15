package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// record scaffolds one file the way a real run would: written through
// writeFile with the recorder on, so the manifest holds its hash.
func record(t *testing.T, root, rel, body string) string {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	release, err := manifest.Start(root, "3.147.0", "scaffold")
	if err != nil {
		t.Fatalf("start recording: %v", err)
	}
	if err := writeFile(path, body); err != nil {
		t.Fatalf("write %s: %v", rel, err)
	}
	if err := release(); err != nil {
		t.Fatalf("save manifest: %v", err)
	}
	return path
}

func readBack(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

// The whole point of the manifest: an upgrade may replace what it wrote, and
// may not replace what you wrote. Before v3.147.0 upgrade did the second thing
// to everybody, every time.
func TestGuardOverwritesUntouchedAndSpareEdited(t *testing.T) {
	root := t.TempDir()

	untouched := record(t, root, "components/data-table.tsx", "export const version = 1\n")
	edited := record(t, root, "components/sidebar.tsx", "export const nav = []\n")

	// The developer edits one of them.
	mine := "export const nav = [{ href: '/reports' }]\n"
	if err := os.WriteFile(edited, []byte(mine), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := startGuard(root, false); err != nil {
		t.Fatalf("start guard: %v", err)
	}
	release, _ := manifest.Start(root, "3.148.0", "scaffold")
	if err := writeFile(untouched, "export const version = 2\n"); err != nil {
		t.Fatal(err)
	}
	if err := writeFile(edited, "export const nav = ['upstream']\n"); err != nil {
		t.Fatal(err)
	}
	written, skipped := stopGuard()
	_ = release()

	if got := readBack(t, untouched); got != "export const version = 2\n" {
		t.Errorf("untouched file was not updated, got %q", got)
	}
	if got := readBack(t, edited); got != mine {
		t.Errorf("an edited file was overwritten: got %q, want %q", got, mine)
	}
	if written != 1 {
		t.Errorf("written = %d, want 1", written)
	}
	if len(skipped) != 1 || skipped[0].Rel != "components/sidebar.tsx" {
		t.Fatalf("skipped = %+v, want just components/sidebar.tsx", skipped)
	}
	if !strings.Contains(skipped[0].Diff(), "upstream") {
		t.Errorf("the diff should show what the upgrade wanted to write:\n%s", skipped[0].Diff())
	}
}

// --force is the escape hatch, and it has to still work: it is exactly what
// upgrade did for every file before the guard existed.
func TestForceOverwritesEvenEditedFiles(t *testing.T) {
	root := t.TempDir()
	path := record(t, root, "app/page.tsx", "original\n")
	if err := os.WriteFile(path, []byte("mine\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := startGuard(root, true); err != nil {
		t.Fatal(err)
	}
	release, _ := manifest.Start(root, "3.148.0", "scaffold")
	if err := writeFile(path, "upstream\n"); err != nil {
		t.Fatal(err)
	}
	_, skipped := stopGuard()
	_ = release()

	if got := readBack(t, path); got != "upstream\n" {
		t.Errorf("--force did not overwrite: got %q", got)
	}
	if len(skipped) != 0 {
		t.Errorf("--force should skip nothing, got %+v", skipped)
	}
}

// Every project generated before the manifest existed has no provenance at
// all. Refusing to update those would mean no existing project ever receives
// another framework fix, which is worse than the problem being solved.
func TestProjectWithoutAManifestUpgradesAsBefore(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "components", "data-table.tsx")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("old\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := startGuard(root, false); err != nil {
		t.Fatal(err)
	}
	if err := writeFile(path, "new\n"); err != nil {
		t.Fatal(err)
	}
	written, skipped := stopGuard()

	if got := readBack(t, path); got != "new\n" {
		t.Errorf("a pre-manifest project should upgrade as it always did, got %q", got)
	}
	if written != 1 || len(skipped) != 0 {
		t.Errorf("written=%d skipped=%d, want 1 and 0", written, len(skipped))
	}
}

// An edit that happens to match the new template is not a conflict, and
// reporting it as one would train people to ignore the warning.
func TestEditThatMatchesUpstreamIsNotAConflict(t *testing.T) {
	root := t.TempDir()
	path := record(t, root, "lib/utils.ts", "old\n")
	if err := os.WriteFile(path, []byte("new\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := startGuard(root, false); err != nil {
		t.Fatal(err)
	}
	if err := writeFile(path, "new\n"); err != nil {
		t.Fatal(err)
	}
	_, skipped := stopGuard()

	if len(skipped) != 0 {
		t.Errorf("an edit matching upstream was reported as a conflict: %+v", skipped)
	}
}

// A file Grit wrote and the developer deleted comes back. Deleting a framework
// component is how people ask for the stock one again.
func TestDeletedFileIsRestored(t *testing.T) {
	root := t.TempDir()
	path := record(t, root, "components/badge.tsx", "stock\n")
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}

	if err := startGuard(root, false); err != nil {
		t.Fatal(err)
	}
	if err := writeFile(path, "stock v2\n"); err != nil {
		t.Fatal(err)
	}
	written, skipped := stopGuard()

	if got := readBack(t, path); got != "stock v2\n" {
		t.Errorf("deleted file not restored, got %q", got)
	}
	if written != 1 || len(skipped) != 0 {
		t.Errorf("written=%d skipped=%d, want 1 and 0", written, len(skipped))
	}
}

// With the guard off, writeFile behaves exactly as it always did. Every other
// command in the CLI depends on that.
func TestWriteFileIsUnaffectedWhenTheGuardIsOff(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "a", "b", "c.ts")
	if err := writeFile(path, "hello\n"); err != nil {
		t.Fatal(err)
	}
	if got := readBack(t, path); got != "hello\n" {
		t.Fatalf("got %q", got)
	}
	if err := writeFile(path, "goodbye\n"); err != nil {
		t.Fatal(err)
	}
	if got := readBack(t, path); got != "goodbye\n" {
		t.Fatalf("second write did not land: %q", got)
	}
}
