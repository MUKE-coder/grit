// Package uiregistry talks to the hosted Grit UI registry.
//
// The registry is the same shadcn-format endpoint anyone can use with
// `npx shadcn add`; this package exists so `grit ui add` can put the file in the
// right place for the project's architecture without the user having to know
// whether their frontend lives at the root, in apps/web, or in apps/admin.
package uiregistry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// DefaultBaseURL is the hosted registry. Overridable with GRIT_UI_REGISTRY,
// which is what the registry's own developers use to point at localhost.
const DefaultBaseURL = "https://ui.gritframework.dev"

// BaseURL resolves the registry origin.
func BaseURL() string {
	if v := strings.TrimSpace(os.Getenv("GRIT_UI_REGISTRY")); v != "" {
		return strings.TrimRight(v, "/")
	}
	return DefaultBaseURL
}

// Item is one entry in the registry index.
type Item struct {
	Name         string   `json:"name"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Categories   []string `json:"categories"`
	Dependencies []string `json:"dependencies"`
}

// Category returns the item's primary category, or "misc".
func (i Item) Category() string {
	if len(i.Categories) > 0 {
		return i.Categories[0]
	}
	return "misc"
}

type index struct {
	Name  string `json:"name"`
	Items []Item `json:"items"`
}

// File is one file carried by a registry item.
type File struct {
	Path    string `json:"path"`
	Target  string `json:"target"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

// Component is a full registry item, source included.
type Component struct {
	Name         string   `json:"name"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Dependencies []string `json:"dependencies"`
	Files        []File   `json:"files"`
	CSSVars      struct {
		Theme map[string]string `json:"theme"`
	} `json:"cssVars"`
}

func httpClient() *http.Client {
	// A component fetch is a few KB; anything past this is a hung connection,
	// not a slow one.
	return &http.Client{Timeout: 20 * time.Second}
}

func get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry returned %s for %s", resp.Status, url)
	}
	return body, nil
}

// List fetches every component in the registry, sorted by name.
func List(ctx context.Context) ([]Item, error) {
	body, err := get(ctx, BaseURL()+"/r/registry.json")
	if err != nil {
		return nil, err
	}
	var idx index
	if err := json.Unmarshal(body, &idx); err != nil {
		return nil, fmt.Errorf("parsing registry index: %w", err)
	}
	sort.Slice(idx.Items, func(i, j int) bool { return idx.Items[i].Name < idx.Items[j].Name })
	return idx.Items, nil
}

// Fetch retrieves one component with its source.
func Fetch(ctx context.Context, name string) (*Component, error) {
	body, err := get(ctx, fmt.Sprintf("%s/r/%s.json", BaseURL(), name))
	if err != nil {
		return nil, err
	}
	var c Component
	if err := json.Unmarshal(body, &c); err != nil {
		return nil, fmt.Errorf("parsing component %s: %w", name, err)
	}
	if len(c.Files) == 0 || strings.TrimSpace(c.Files[0].Content) == "" {
		// Guarding this explicitly because the original registry shipped items
		// whose files carried a path but no content — writing that would create
		// an empty component and look like a successful install.
		return nil, fmt.Errorf("component %s came back with no source", name)
	}
	return &c, nil
}

// TargetDir picks where a component should be written.
//
// Preference order is deliberate: apps/web is the user-facing frontend and the
// likeliest destination, apps/admin next, then a flat single-app frontend, then
// the project root. Returns the directory and a human label for the report.
func TargetDir(root string) (dir string, label string, err error) {
	candidates := []struct{ rel, label string }{
		{filepath.Join("apps", "web"), "apps/web"},
		{filepath.Join("apps", "admin"), "apps/admin"},
		{"frontend", "frontend"},
		{"web", "web"},
	}
	for _, c := range candidates {
		p := filepath.Join(root, c.rel)
		if st, statErr := os.Stat(p); statErr == nil && st.IsDir() {
			return p, c.label, nil
		}
	}
	return "", "", fmt.Errorf(
		"no frontend app found (looked for apps/web, apps/admin, frontend/ and web/)\n\n" +
			"grit ui targets a React app. An --api project has no frontend to install into")
}

// Write saves the component's source at the path the registry item asks for.
//
// The item's target is honoured rather than a name derived here, so that
// `grit ui add` and `npx shadcn add` put the file in exactly the same place.
// Deriving it locally is how the two drift apart.
//
// It refuses to clobber an existing file unless force is set, because a
// component someone has already edited is their code, not ours.
func Write(appDir string, c *Component, force bool) (string, error) {
	rel := strings.TrimSpace(c.Files[0].Target)
	if rel == "" {
		rel = strings.TrimSpace(c.Files[0].Path)
	}
	if rel == "" {
		rel = filepath.Join("components", "grit-ui", c.Name+".tsx")
	}
	// A registry is remote input: a target of "../../.ssh/authorized_keys"
	// must not escape the app directory.
	rel = filepath.Clean(filepath.FromSlash(rel))
	if filepath.IsAbs(rel) || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("component %s has an unsafe target path %q", c.Name, c.Files[0].Target)
	}

	out := filepath.Join(appDir, rel)
	outDir := filepath.Dir(out)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", fmt.Errorf("creating %s: %w", outDir, err)
	}
	if !force {
		if _, err := os.Stat(out); err == nil {
			return "", fmt.Errorf("%s already exists — pass --force to overwrite", out)
		}
	}
	if err := os.WriteFile(out, []byte(c.Files[0].Content), 0o644); err != nil {
		return "", fmt.Errorf("writing %s: %w", out, err)
	}
	return out, nil
}
