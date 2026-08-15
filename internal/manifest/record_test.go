package manifest

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNoteIsANoOpWhenNotRecording(t *testing.T) {
	if Active() {
		t.Fatal("recording should be off by default")
	}
	// The point: a library caller or a test that never opted in must not have
	// its writes recorded into some other project's manifest.
	Note(filepath.Join(t.TempDir(), "x.ts"), "hello")
	if err := Stop(); err != nil {
		t.Fatalf("Stop without Start should be harmless, got %v", err)
	}
}

func TestStartNoteStop(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "apps", "web", "page.tsx")

	release, err := Start(root, "3.147.0", "scaffold")
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	if !Active() {
		t.Fatal("recording should be on after Start")
	}
	write(t, path, "body\n")
	Note(path, "body\n")
	if err := release(); err != nil {
		t.Fatalf("release: %v", err)
	}
	if Active() {
		t.Fatal("recording should be off after release")
	}

	m, _ := Load(root)
	if got := m.StatusOf(root, "apps/web/page.tsx"); got != Unchanged {
		t.Fatalf("recorded file reads as %v, want unchanged", got)
	}
	if m.Files["apps/web/page.tsx"].Generator != "scaffold" {
		t.Fatalf("wrong generator: %+v", m.Files["apps/web/page.tsx"])
	}
}

// Generator.Run calls itself for an inline child resource. A nested Start that
// reloaded the manifest from disk would discard everything the parent had
// recorded and not yet saved, so the parent's own files would silently vanish
// from the manifest.
func TestNestedStartKeepsTheParentsRecords(t *testing.T) {
	root := t.TempDir()

	release, err := Start(root, "3.147.0", "resource:Invoice")
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	write(t, filepath.Join(root, "invoice.go"), "parent\n")
	Note(filepath.Join(root, "invoice.go"), "parent\n")

	// The inline child resource, generated mid-run.
	inner, err := Start(root, "3.147.0", "resource:InvoiceItem")
	if err != nil {
		t.Fatalf("nested start: %v", err)
	}
	write(t, filepath.Join(root, "invoice_item.go"), "child\n")
	Note(filepath.Join(root, "invoice_item.go"), "child\n")
	if err := inner(); err != nil {
		t.Fatalf("nested release: %v", err)
	}

	// Back to the parent, which must still be recording under its own label.
	write(t, filepath.Join(root, "invoice_handler.go"), "parent2\n")
	Note(filepath.Join(root, "invoice_handler.go"), "parent2\n")
	if err := release(); err != nil {
		t.Fatalf("release: %v", err)
	}

	m, _ := Load(root)
	for path, wantGenerator := range map[string]string{
		"invoice.go":         "resource:Invoice",
		"invoice_item.go":    "resource:InvoiceItem",
		"invoice_handler.go": "resource:Invoice",
	} {
		entry, ok := m.Files[path]
		if !ok {
			t.Errorf("%s was not recorded", path)
			continue
		}
		if entry.Generator != wantGenerator {
			t.Errorf("%s recorded as %q, want %q", path, entry.Generator, wantGenerator)
		}
	}
}

func TestNoteIgnoresPathsOutsideTheProject(t *testing.T) {
	base := t.TempDir()
	root := filepath.Join(base, "project")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}

	release, _ := Start(root, "3.147.0", "scaffold")
	Note(filepath.Join(base, "outside.txt"), "x")
	_ = release()

	m, _ := Load(root)
	if len(m.Files) != 0 {
		t.Fatalf("a write outside the project was recorded: %v", m.Files)
	}
}

// Injection edits a file Grit already owns. Without Refresh, routes.go would
// read as hand-edited the moment the first resource was generated.
func TestRefreshKeepsTheOriginalOwner(t *testing.T) {
	root := t.TempDir()
	routes := filepath.Join(root, "routes.go")

	release, _ := Start(root, "3.147.0", "scaffold")
	write(t, routes, "// grit:routes\n")
	Note(routes, "// grit:routes\n")
	_ = release()

	release, _ = Start(root, "3.147.0", "resource:Product")
	write(t, routes, "protected.GET(\"/products\")\n// grit:routes\n")
	Refresh(routes)
	_ = release()

	m, _ := Load(root)
	entry := m.Files["routes.go"]
	if entry.Generator != "scaffold" {
		t.Errorf("Refresh changed the owner to %q; injection does not transfer ownership", entry.Generator)
	}
	if got := m.StatusOf(root, "routes.go"); got != Unchanged {
		t.Errorf("injected file reads as %v, want unchanged", got)
	}
}

// Refreshing something Grit never wrote must not adopt it.
func TestRefreshIgnoresUntrackedFiles(t *testing.T) {
	root := t.TempDir()
	mine := filepath.Join(root, "mine.ts")
	write(t, mine, "handwritten\n")

	release, _ := Start(root, "3.147.0", "resource:Product")
	Refresh(mine)
	_ = release()

	m, _ := Load(root)
	if _, adopted := m.Files["mine.ts"]; adopted {
		t.Fatal("Refresh adopted a file Grit never wrote")
	}
}
