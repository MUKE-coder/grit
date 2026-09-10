package scaffold

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// grit upgrade treated every project as Next.js, and on a TanStack project
// wrote the Next admin and web apps, package.json included, over the Vite ones.
func TestReadProjectFrontend(t *testing.T) {
	cases := map[string]struct {
		gritJSON string
		viteApp  string
		want     Frontend
	}{
		"recorded tanstack":           {gritJSON: `{"frontend":"tanstack"}`, want: FrontendTanStack},
		"recorded next":               {gritJSON: `{"frontend":"next"}`, want: FrontendNext},
		"unrecorded, vite admin":      {viteApp: "admin", want: FrontendTanStack},
		"unrecorded, vite web":        {viteApp: "web", want: FrontendTanStack},
		"unrecorded, nothing to tell": {want: FrontendNext},
	}
	for name, c := range cases {
		root := t.TempDir()
		if c.gritJSON != "" {
			if err := os.WriteFile(filepath.Join(root, "grit.json"), []byte(c.gritJSON), 0644); err != nil {
				t.Fatal(err)
			}
		}
		if c.viteApp != "" {
			dir := filepath.Join(root, "apps", c.viteApp)
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, "vite.config.ts"), []byte("export default {}\n"), 0644); err != nil {
				t.Fatal(err)
			}
		}
		if got := readProjectFrontend(root); got != c.want {
			t.Errorf("%s: got %q, want %q", name, got, c.want)
		}
	}
}

// The Vite apps' tests could not start: the shared setup file imports
// jest-dom, the login test imports user-event, the environment is jsdom, and
// none of the three was declared.
func TestViteAppsDeclareTheirTestLibraries(t *testing.T) {
	opts := Options{ProjectName: "demo", Frontend: FrontendTanStack}
	for name, src := range map[string]string{
		"vite admin": adminTanStackPackageJSON(opts),
		"vite web":   webTanStackPackageJSON(opts),
	} {
		var pkg struct {
			DevDependencies map[string]string `json:"devDependencies"`
		}
		if err := json.Unmarshal([]byte(src), &pkg); err != nil {
			t.Fatalf("%s: package.json does not parse: %v", name, err)
		}
		for _, dep := range []string{"vitest", "jsdom", "@testing-library/react", "@testing-library/jest-dom", "@testing-library/user-event"} {
			if pkg.DevDependencies[dep] == "" {
				t.Errorf("%s does not declare %s", name, dep)
			}
		}
	}
}

// A Vite app's tests walked up the tree for a PostCSS config, found a stray
// root one in Tailwind v3 style, and failed on a plugin nobody had installed.
func TestVitestConfigsPinPostCSS(t *testing.T) {
	for name, src := range map[string]string{
		"admin": adminVitestConfig(),
		"web":   webVitestConfig(),
	} {
		if !strings.Contains(src, "css: { postcss: { plugins: [] } }") {
			t.Errorf("the %s vitest config does not pin PostCSS", name)
		}
	}
}
