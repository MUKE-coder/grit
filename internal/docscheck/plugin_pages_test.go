package docscheck

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/plugin"
)

/*
 * Every plugin a reader can install should have a page telling them what it
 * does. `grit plugin add webhooks` worked for months with nothing on the
 * site about it: findable only by running `grit plugin list`, which nobody
 * does before deciding whether a thing exists.
 *
 * The marketplace page lists them, but a line in a list is not documentation,
 * and the other nine plugins each have a page.
 *
 * Two lists that have to agree, with nothing comparing them, so: this.
 */

// Pages under /docs/plugins that are about plugins as a subject rather than
// about one plugin.
var notAPlugin = map[string]bool{
	"overview":    true,
	"marketplace": true,
	"authoring":   true,
	"go-modules":  true,
}

func TestEveryInstallablePluginHasAPage(t *testing.T) {
	documented := map[string]bool{}
	root := filepath.Join("..", "..", "docs", "app", "docs", "plugins")

	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("reading %s: %v", root, err)
	}
	for _, e := range entries {
		if !e.IsDir() || notAPlugin[e.Name()] {
			continue
		}
		if _, err := os.Stat(filepath.Join(root, e.Name(), "page.tsx")); err == nil {
			documented[e.Name()] = true
		}
	}

	var missing []string
	for _, p := range plugin.All() {
		if !documented[p.Name] {
			missing = append(missing, p.Name)
		}
	}
	sort.Strings(missing)

	if len(missing) > 0 {
		t.Errorf("installable with `grit plugin add` and no page under "+
			"/docs/plugins: %s\n\nA plugin nobody can read about is a plugin "+
			"nobody installs.", strings.Join(missing, ", "))
	}
}
