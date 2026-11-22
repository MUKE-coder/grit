package explain

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// The editing advice must come from the manifest, not from the path.
//
// That is the whole design: a hardcoded list of framework-owned paths would
// give the same answer for a file somebody has already edited, and would then
// be telling them the opposite of what grit upgrade is going to do.
func TestAdviceFollowsTheManifestRatherThanThePath(t *testing.T) {
	root := t.TempDir()
	rel := "apps/api/internal/models/product.go"
	abs := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		t.Fatal(err)
	}

	const written = "package models\n\ntype Product struct{}\n"
	if err := os.WriteFile(abs, []byte(written), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &manifest.Manifest{Files: map[string]manifest.Entry{}}
	m.Record(rel, "resource:Product", "3.561.0", written)
	if err := m.Save(root, "3.561.0"); err != nil {
		t.Fatalf("save manifest: %v", err)
	}

	// Untouched: Grit wrote it, so an upgrade may rewrite it.
	e, err := Explain(root, abs)
	if err != nil {
		t.Fatalf("explain: %v", err)
	}
	if e.Kind != Generated {
		t.Fatalf("an unedited generated file is %v, want Generated", e.Kind)
	}
	if e.Generator != "resource:Product" {
		t.Errorf("generator = %q, want resource:Product", e.Generator)
	}
	if e.Version != "3.561.0" {
		t.Errorf("version = %q, want 3.561.0", e.Version)
	}
	if !strings.Contains(e.Advice(), "may rewrite it") {
		t.Errorf("advice for an unedited file does not mention rewriting: %q", e.Advice())
	}

	// The same path, edited. The advice must flip, and the path has not
	// changed, which is the point.
	if err := os.WriteFile(abs, []byte(written+"\n// mine\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	e, err = Explain(root, abs)
	if err != nil {
		t.Fatalf("explain after edit: %v", err)
	}
	if e.Kind != Edited {
		t.Fatalf("an edited generated file is %v, want Edited", e.Kind)
	}
	if !strings.Contains(e.Advice(), "leave it alone") {
		t.Errorf("advice for an edited file does not say it is left alone: %q", e.Advice())
	}

	// And gone.
	if err := os.Remove(abs); err != nil {
		t.Fatal(err)
	}
	e, err = Explain(root, abs)
	if err != nil {
		t.Fatalf("explain after delete: %v", err)
	}
	if e.Kind != Gone {
		t.Fatalf("a deleted generated file is %v, want Gone", e.Kind)
	}
}

// A file Grit never wrote is the reader's own, whatever its path looks like.
func TestUntrackedFileIsYoursEvenInAGeneratedDirectory(t *testing.T) {
	root := t.TempDir()
	rel := "apps/api/internal/services/mine.go"
	abs := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(abs, []byte("package services\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	e, err := Explain(root, abs)
	if err != nil {
		t.Fatalf("explain: %v", err)
	}
	if e.Kind != Yours {
		t.Fatalf("an untracked file is %v, want Yours", e.Kind)
	}
	// It still gets the description for the directory, because what the
	// directory is for is true regardless of who wrote this file.
	if !strings.Contains(e.Purpose, "business rules") {
		t.Errorf("a file under internal/services should be described as a service: %q", e.Purpose)
	}
	if !strings.Contains(e.Advice(), "yours") {
		t.Errorf("advice = %q, want it to say the file is yours", e.Advice())
	}
}

// A path outside the project is an error rather than a confident wrong answer.
func TestPathOutsideTheProjectIsRefused(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "elsewhere.go")

	if _, err := Explain(root, outside); err == nil {
		t.Fatal("a path outside the project was accepted")
	}
}

// Every pattern must say something in all three fields, or the output has a
// gap where a sentence should be.
func TestEveryPatternIsComplete(t *testing.T) {
	for i, p := range patterns {
		if strings.TrimSpace(p.purpose) == "" {
			t.Errorf("pattern %d has no purpose", i)
		}
		if strings.TrimSpace(p.edit) == "" {
			t.Errorf("pattern %d (%s) has no editing advice", i, firstWords(p.purpose))
		}
		if p.read != "" && !strings.HasPrefix(p.read, "/") {
			t.Errorf("pattern %d (%s) read path %q should start with /",
				i, firstWords(p.purpose), p.read)
		}
	}
	if len(patterns) < 15 {
		t.Fatalf("only %d patterns, which is fewer than the file map has shapes", len(patterns))
	}
}

func firstWords(s string) string {
	if len(s) > 40 {
		return s[:40]
	}
	return s
}
