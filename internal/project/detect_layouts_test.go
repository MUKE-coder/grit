package project

import (
	"os"
	"path/filepath"
	"testing"
)

// Every Grit project is a Grit project, whatever shape it is.
//
// Detection tested for turbo.json (a monorepo) and wails.json (a standalone
// desktop app) and nothing else. An --api project has no turbo.json, because it
// has no frontends to orchestrate, and a --single project has neither: both fell
// through to "not inside a Grit project", so `grit start` printed its own help
// and started nothing. That is the command the home page prints and the one
// every tutorial reaches at step four.
func TestEveryArchitectureIsDetected(t *testing.T) {
	for _, tc := range []struct {
		name   string
		layout map[string]string
		want   ProjectType
	}{
		{
			name: "monorepo",
			layout: map[string]string{
				"grit.json":             `{"architecture":"triple"}`,
				"turbo.json":            `{}`,
				"apps/api/go.mod":       "module contacts/apps/api\n",
				"apps/web/package.json": `{}`,
			},
			want: ProjectWeb,
		},
		{
			name: "api only, which has no turbo.json",
			layout: map[string]string{
				"grit.json":       `{"architecture":"api"}`,
				"apps/api/go.mod": "module contacts/apps/api\n",
			},
			want: ProjectWeb,
		},
		{
			name: "single, whose module is in api/",
			layout: map[string]string{
				"grit.json":    `{"architecture":"single"}`,
				"api/go.mod":   "module contacts\n",
				"package.json": `{}`,
			},
			want: ProjectWeb,
		},
		{
			name: "a single from before the api/ layout",
			layout: map[string]string{
				"grit.json":             `{"architecture":"single"}`,
				"go.mod":                "module contacts\n",
				"frontend/package.json": `{}`,
			},
			want: ProjectWeb,
		},
		{
			name: "standalone desktop",
			layout: map[string]string{
				"wails.json": `{}`,
				"go.mod":     "module contacts\n",
			},
			want: ProjectDesktop,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := writeLayout(t, tc.layout)

			info, err := DetectProjectFrom(root)
			if err != nil {
				t.Fatalf("not detected: %v", err)
			}
			if info.Type != tc.want {
				t.Errorf("type is %q, want %q", info.Type, tc.want)
			}
			if info.Root != root {
				t.Errorf("root is %s, want %s", info.Root, root)
			}
			if info.Module == "" {
				t.Error("no module was read, so anything keyed on it has nothing to work with")
			}
		})
	}
}

// And a directory that is not a Grit project is still not one: the detection
// has to stay a detection, or every command starts doing something in a
// directory the user did not mean.
func TestANonProjectIsStillNotAProject(t *testing.T) {
	root := writeLayout(t, map[string]string{
		"go.mod":  "module unrelated\n",
		"main.go": "package main\n",
	})
	if _, err := DetectProjectFrom(root); err == nil {
		t.Error("a plain Go module was taken for a Grit project")
	}
}

// Detection walks up, so a command run from inside the project finds its root.
func TestDetectionWalksUpFromASubdirectory(t *testing.T) {
	root := writeLayout(t, map[string]string{
		"grit.json":       `{"architecture":"api"}`,
		"apps/api/go.mod": "module contacts/apps/api\n",
	})
	info, err := DetectProjectFrom(filepath.Join(root, "apps", "api"))
	if err != nil {
		t.Fatalf("not detected from a subdirectory: %v", err)
	}
	if info.Root != root {
		t.Errorf("root is %s, want %s", info.Root, root)
	}
}

func writeLayout(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, body := range files {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}
