package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The generator names the files it actually wrote.
//
// Resource definitions moved to a folder each, so a new project gets
// apps/admin/resources/posts/posts.ts. The success output kept announcing the
// flat layout it replaced:
//
//	✓ apps/admin/resources/posts.ts        <- never written
//
// Which is the first thing a newcomer clicks after a generate, and the first
// thing that does not open. The overlay beside it, the one file that is never
// regenerated and therefore the only safe place to customise, was written and
// not mentioned at all.
func TestGenerateAnnouncesPathsThatExist(t *testing.T) {
	const module = "notes/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Post", "title:string,published:bool")
	if err != nil {
		t.Fatalf("fields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	names := g.Names()
	resourcesRoot := filepath.Join(root, "apps", "admin", "resources")

	if err := g.writeResourceDefinition(names); err != nil {
		t.Fatalf("resource definition: %v", err)
	}

	out := captureStdout(t, func() {
		g.announceResourceDef(resourcesRoot, names, false)
	})

	lines := announcedPaths(out)
	if len(lines) != 2 {
		t.Fatalf("want the definition and its overlay announced, got %d line(s):\n%s", len(lines), out)
	}
	for _, rel := range lines {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); err != nil {
			t.Errorf("announced a file that is not there: %s", rel)
		}
	}

	// And specifically the folder layout, not the flat path that was printed
	// before: a test that only checks existence passes on either.
	if got := lines[0]; got != "apps/admin/resources/posts/posts.ts" {
		t.Errorf("definition announced as %q", got)
	}
	if got := lines[1]; !strings.HasPrefix(got, "apps/admin/resources/posts/posts.custom.tsx") {
		t.Errorf("overlay announced as %q", got)
	}
}

// An overlay that was already there is not announced as new. It holds work
// this run did not do, and writeResourceCustomStub deliberately left it alone.
func TestExistingOverlayIsNotAnnouncedAsCreated(t *testing.T) {
	const module = "notes/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Post", "title:string")
	if err != nil {
		t.Fatalf("fields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	names := g.Names()
	resourcesRoot := filepath.Join(root, "apps", "admin", "resources")

	if err := g.writeResourceDefinition(names); err != nil {
		t.Fatalf("resource definition: %v", err)
	}

	out := captureStdout(t, func() {
		g.announceResourceDef(resourcesRoot, names, true)
	})
	if lines := announcedPaths(out); len(lines) != 1 {
		t.Fatalf("want only the definition, got %d line(s):\n%s", len(lines), out)
	}
}

// announcedPaths pulls the path out of each "  ✓ <path>[  (note)]" line.
func announcedPaths(out string) []string {
	var paths []string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "✓ ") {
			continue
		}
		rest := strings.TrimPrefix(line, "✓ ")
		if i := strings.Index(rest, "  ("); i >= 0 {
			rest = rest[:i]
		}
		paths = append(paths, strings.TrimSpace(rest))
	}
	return paths
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	saved := os.Stdout
	os.Stdout = w
	fn()
	os.Stdout = saved
	if err := w.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	var sb strings.Builder
	buf := make([]byte, 4096)
	for {
		n, readErr := r.Read(buf)
		sb.Write(buf[:n])
		if readErr != nil {
			break
		}
	}
	return sb.String()
}
