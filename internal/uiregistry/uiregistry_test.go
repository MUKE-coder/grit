package uiregistry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const itemJSON = `{
  "name": "hero-split-01",
  "title": "Split Hero",
  "description": "Side-by-side hero",
  "dependencies": ["lucide-react"],
  "files": [{
    "path": "components/grit-ui/hero-split-01.tsx",
    "target": "components/grit-ui/hero-split-01.tsx",
    "type": "registry:component",
    "content": "export default function HeroSplit01() { return null }\n"
  }],
  "cssVars": { "theme": { "--accent": "#6c5ce7" } }
}`

const indexJSON = `{
  "name": "grit-ui",
  "items": [
    {"name": "zebra-01", "title": "Zebra", "categories": ["layout"]},
    {"name": "alpha-01", "title": "Alpha", "categories": ["marketing"], "dependencies": ["lucide-react"]}
  ]
}`

func testServer(t *testing.T, handler http.HandlerFunc) string {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	t.Setenv("GRIT_UI_REGISTRY", srv.URL)
	return srv.URL
}

func TestBaseURLPrefersTheEnvOverride(t *testing.T) {
	if got := BaseURL(); got != DefaultBaseURL {
		t.Errorf("BaseURL() = %q, want the hosted default %q", got, DefaultBaseURL)
	}
	t.Setenv("GRIT_UI_REGISTRY", "http://localhost:3100/")
	if got := BaseURL(); got != "http://localhost:3100" {
		t.Errorf("BaseURL() = %q, want the override with its trailing slash trimmed", got)
	}
}

func TestListSortsByName(t *testing.T) {
	testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/r/registry.json" {
			t.Errorf("unexpected index path %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(indexJSON))
	})

	items, err := List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0].Name != "alpha-01" {
		t.Errorf("index not sorted: first item is %q", items[0].Name)
	}
	if got := items[0].Category(); got != "marketing" {
		t.Errorf("Category() = %q, want marketing", got)
	}
}

func TestCategoryFallsBackWhenAbsent(t *testing.T) {
	if got := (Item{Name: "x"}).Category(); got != "misc" {
		t.Errorf("Category() = %q, want misc", got)
	}
}

func TestFetchReturnsTheInlinedSource(t *testing.T) {
	testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/r/hero-split-01.json" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte(itemJSON))
	})

	c, err := Fetch(context.Background(), "hero-split-01")
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if !strings.Contains(c.Files[0].Content, "HeroSplit01") {
		t.Errorf("source not carried in the payload: %q", c.Files[0].Content)
	}
	if len(c.Dependencies) != 1 || c.Dependencies[0] != "lucide-react" {
		t.Errorf("dependencies = %v", c.Dependencies)
	}
}

// The original registry shipped items whose files carried a path but no
// content. Writing one of those produces an empty component and looks like a
// successful install, so an empty payload has to be an error.
func TestFetchRejectsAnItemWithNoSource(t *testing.T) {
	testServer(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"name":"empty-01","files":[{"path":"a.tsx","content":""}]}`))
	})

	if _, err := Fetch(context.Background(), "empty-01"); err == nil {
		t.Fatal("expected an error for an item with no source")
	} else if !strings.Contains(err.Error(), "no source") {
		t.Errorf("error should say the source was missing, got %v", err)
	}
}

func TestFetchSurfacesAnHTTPError(t *testing.T) {
	testServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"nope"}`))
	})

	if _, err := Fetch(context.Background(), "missing-01"); err == nil {
		t.Fatal("expected an error for a 404")
	}
}

func TestTargetDirPrefersTheUserFacingApp(t *testing.T) {
	root := t.TempDir()
	for _, d := range []string{"apps/admin", "apps/web", "frontend"} {
		if err := os.MkdirAll(filepath.Join(root, filepath.FromSlash(d)), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	dir, label, err := TargetDir(root)
	if err != nil {
		t.Fatalf("TargetDir: %v", err)
	}
	if label != "apps/web" {
		t.Errorf("label = %q, want apps/web to win over admin and frontend", label)
	}
	if dir != filepath.Join(root, "apps", "web") {
		t.Errorf("dir = %q", dir)
	}
}

func TestTargetDirFallsBackToAFlatFrontend(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "frontend"), 0o755); err != nil {
		t.Fatal(err)
	}
	_, label, err := TargetDir(root)
	if err != nil {
		t.Fatalf("TargetDir: %v", err)
	}
	if label != "frontend" {
		t.Errorf("label = %q, want frontend", label)
	}
}

// An --api project has nowhere to put a React component, and the message has to
// say so rather than inventing a directory.
func TestTargetDirExplainsWhenThereIsNoFrontend(t *testing.T) {
	_, _, err := TargetDir(t.TempDir())
	if err == nil {
		t.Fatal("expected an error when no frontend exists")
	}
	if !strings.Contains(err.Error(), "apps/web") || !strings.Contains(err.Error(), "--api") {
		t.Errorf("error should name what was looked for and why, got %v", err)
	}
}

func TestWriteCreatesTheFileAndRefusesToClobber(t *testing.T) {
	app := t.TempDir()
	c := &Component{
		Name:  "hero-split-01",
		Files: []File{{Content: "export default function X() { return null }\n"}},
	}

	path, err := Write(app, c, false)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}
	want := filepath.Join(app, "components", "grit-ui", "hero-split-01.tsx")
	if path != want {
		t.Errorf("path = %q, want %q", path, want)
	}
	data, err := os.ReadFile(path)
	if err != nil || !strings.Contains(string(data), "export default") {
		t.Fatalf("file not written correctly: %v %q", err, data)
	}

	// A component the user may have edited is their code now.
	if _, err := Write(app, c, false); err == nil {
		t.Error("expected Write to refuse to overwrite without --force")
	}
	if _, err := Write(app, c, true); err != nil {
		t.Errorf("--force should overwrite, got %v", err)
	}
}
