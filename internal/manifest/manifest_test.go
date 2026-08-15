package manifest

import (
	"os"
	"path/filepath"
	"testing"
)

func write(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestLoadMissingIsEmptyNotAnError(t *testing.T) {
	m, err := Load(t.TempDir())
	if err != nil {
		t.Fatalf("a project with no manifest must load cleanly, got %v", err)
	}
	if len(m.Files) != 0 {
		t.Fatalf("expected no entries, got %d", len(m.Files))
	}
}

func TestLoadCorruptFallsBackRatherThanFailing(t *testing.T) {
	root := t.TempDir()
	write(t, Path(root), "{not json at all")

	m, err := Load(root)
	if err != nil {
		t.Fatalf("a corrupt manifest must not stop the command, got %v", err)
	}
	if len(m.Files) != 0 {
		t.Fatalf("expected an empty manifest, got %d entries", len(m.Files))
	}
}

func TestRoundTrip(t *testing.T) {
	root := t.TempDir()

	m, _ := Load(root)
	m.Record("apps/admin/components/data-table.tsx", "scaffold", "3.147.0", "export const x = 1\n")
	if err := m.Save(root, "3.147.0"); err != nil {
		t.Fatalf("save: %v", err)
	}

	back, err := Load(root)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	entry, ok := back.Files["apps/admin/components/data-table.tsx"]
	if !ok {
		t.Fatal("entry did not survive the round trip")
	}
	if entry.Generator != "scaffold" || entry.Grit != "3.147.0" {
		t.Fatalf("provenance lost: %+v", entry)
	}
	if back.Format != FormatVersion {
		t.Fatalf("format version = %d, want %d", back.Format, FormatVersion)
	}
}

func TestStatusOf(t *testing.T) {
	root := t.TempDir()
	body := "line one\nline two\n"
	rel := "apps/web/app/page.tsx"

	m, _ := Load(root)
	m.Record(rel, "scaffold", "3.147.0", body)

	// Not on disk at all.
	if got := m.StatusOf(root, rel); got != Missing {
		t.Fatalf("deleted file: got %v, want missing", got)
	}

	// Byte-for-byte what we wrote.
	write(t, filepath.Join(root, filepath.FromSlash(rel)), body)
	if got := m.StatusOf(root, rel); got != Unchanged {
		t.Fatalf("untouched file: got %v, want unchanged", got)
	}

	// Edited since.
	write(t, filepath.Join(root, filepath.FromSlash(rel)), body+"// mine\n")
	if got := m.StatusOf(root, rel); got != Modified {
		t.Fatalf("edited file: got %v, want modified", got)
	}

	// Never recorded.
	if got := m.StatusOf(root, "apps/web/app/mine.tsx"); got != Untracked {
		t.Fatalf("unknown file: got %v, want untracked", got)
	}
}

// A file recorded with LF and checked out by git with CRLF is the same file.
// Hashing raw bytes would call every file in a fresh Windows clone modified,
// and an upgrade trusting that would refuse to update anything at all.
func TestCRLFIsNotAnEdit(t *testing.T) {
	root := t.TempDir()
	rel := "turbo.json"
	lf := "{\n  \"tasks\": {}\n}\n"
	crlf := "{\r\n  \"tasks\": {}\r\n}\r\n"

	m, _ := Load(root)
	m.Record(rel, "scaffold", "3.147.0", lf)
	write(t, filepath.Join(root, rel), crlf)

	if got := m.StatusOf(root, rel); got != Unchanged {
		t.Fatalf("CRLF checkout read as %v, want unchanged", got)
	}
}

func TestHashDistinguishesRealEdits(t *testing.T) {
	if Hash("a\n") == Hash("b\n") {
		t.Fatal("different content hashed the same")
	}
	if Hash("a\nb\n") != Hash("a\r\nb\r\n") {
		t.Fatal("line endings must not change the hash")
	}
}

func TestKeyNormalisation(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"apps/web/page.tsx", "apps/web/page.tsx"},
		{"./apps/web/page.tsx", "apps/web/page.tsx"},
		{filepath.Join("apps", "web", "page.tsx"), "apps/web/page.tsx"},
	} {
		if got := Key(tc.in); got != tc.want {
			t.Errorf("Key(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestRelRejectsPathsOutsideTheProject(t *testing.T) {
	root := filepath.Join(t.TempDir(), "project")

	if _, ok := Rel(root, filepath.Join(root, "apps", "api", "main.go")); !ok {
		t.Fatal("a path inside the project should be recordable")
	}
	if got, ok := Rel(root, filepath.Join(filepath.Dir(root), "elsewhere.txt")); ok {
		t.Fatalf("a path outside the project should not be recorded, got %q", got)
	}
}

func TestForgetAndEntries(t *testing.T) {
	m := &Manifest{Files: map[string]Entry{}}
	m.Record("b.ts", "scaffold", "1", "x")
	m.Record("a.ts", "scaffold", "1", "x")

	got := m.Entries()
	if len(got) != 2 || got[0] != "a.ts" || got[1] != "b.ts" {
		t.Fatalf("Entries not sorted: %v", got)
	}

	m.Forget("a.ts")
	if _, still := m.Files["a.ts"]; still {
		t.Fatal("Forget did not remove the entry")
	}
}
