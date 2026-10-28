package plugin

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// marketplacePage is the docs page that lists every plugin for a browser.
const marketplacePage = "docs/app/docs/plugins/marketplace/page.tsx"

// TestMarketplacePageListsEveryPlugin compares the registry with the page.
//
// The page cannot import the Go registry, so its list of official plugins is
// hand-maintained. Two lists with nothing comparing them drift, and the way
// this one drifts is the quiet way: somebody adds a plugin, `grit plugin
// list` shows it, the marketplace does not, and nobody notices because the
// page still looks complete. This is the thing comparing them.
func TestMarketplacePageListsEveryPlugin(t *testing.T) {
	page, ok := readRepoFile(t, marketplacePage)
	if !ok {
		t.Skipf("%s not found; this test only runs inside the Grit repository", marketplacePage)
	}

	for _, p := range All() {
		// Quoted, so a plugin called "push" is not satisfied by the word
		// "push" appearing in somebody's prose.
		needle := `name: '` + p.Name + `'`
		if !strings.Contains(page, needle) {
			t.Errorf("the %q plugin is registered but missing from %s.\n"+
				"Add it to the OFFICIAL list there, with a summary and a docs link.",
				p.Name, marketplacePage)
		}
	}
}

// TestMarketplacePageHasNoPluginsThatDoNotExist is the other direction: a
// page advertising `grit plugin add something` for a plugin that was renamed
// or removed sends people to an error.
func TestMarketplacePageHasNoPluginsThatDoNotExist(t *testing.T) {
	page, ok := readRepoFile(t, marketplacePage)
	if !ok {
		t.Skipf("%s not found; this test only runs inside the Grit repository", marketplacePage)
	}

	registered := map[string]bool{}
	for _, p := range All() {
		registered[p.Name] = true
	}
	// Go module plugins are listed on the page too, and are not in the
	// registry because they install with `go get`. They live in their own
	// MODULES array, so only the OFFICIAL block is checked.
	// Anchored on the name rather than the whole declaration: the page types
	// the array now ("const OFFICIAL: PluginEntry[] = ["), and an anchor that
	// pins the exact spelling makes this test refuse to run every time the
	// page is refactored.
	official, found := sliceBetween(page, "const OFFICIAL", "\n]")
	if !found {
		t.Fatalf("%s no longer has an OFFICIAL list; this test cannot check it", marketplacePage)
	}

	for _, name := range namesIn(official) {
		if !registered[name] {
			t.Errorf("%s lists %q under OFFICIAL, but no plugin by that name is registered.\n"+
				"`grit plugin add %s` would fail for anybody who followed the page.",
				marketplacePage, name, name)
		}
	}
}

func namesIn(block string) []string {
	var out []string
	rest := block
	for {
		at := strings.Index(rest, "name: '")
		if at < 0 {
			return out
		}
		rest = rest[at+len("name: '"):]
		end := strings.Index(rest, "'")
		if end < 0 {
			return out
		}
		out = append(out, rest[:end])
		rest = rest[end:]
	}
}

func sliceBetween(s, open, close string) (string, bool) {
	start := strings.Index(s, open)
	if start < 0 {
		return "", false
	}
	start += len(open)
	end := strings.Index(s[start:], close)
	if end < 0 {
		return "", false
	}
	return s[start : start+end], true
}

// readRepoFile reads a path relative to the repository root, found by walking
// up from the test's directory until go.mod appears.
func readRepoFile(t *testing.T, rel string) (string, bool) {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		return "", false
	}
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			raw, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(rel)))
			if err != nil {
				return "", false
			}
			return string(raw), true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", false
}

// TestEveryPluginPageHasAWorkedExample.
//
// The marketplace sends people to a plugin's page to find out how to use it.
// A page that only says what the plugin is leaves them to work the rest out
// from the CLI, which is the thing the page existed to save them.
func TestEveryPluginPageHasAWorkedExample(t *testing.T) {
	for _, p := range All() {
		page, ok := readRepoFile(t, "docs/app/docs/plugins/"+p.Name+"/page.tsx")
		if !ok {
			// webhooks documents itself under /docs/backend; the marketplace
			// test above is what checks every plugin is reachable.
			continue
		}
		if !strings.Contains(page, "const STEPS: PluginStep[]") {
			t.Errorf("the %s page has no STEPS array, so it has no worked example", p.Name)
		}
		if !strings.Contains(page, "<PluginSteps steps={STEPS}") {
			t.Errorf("the %s page declares STEPS but never renders them", p.Name)
		}
		if !strings.Contains(page, "const RELATED: RelatedPlugin[]") {
			t.Errorf("the %s page has no RELATED array", p.Name)
		}
		if !strings.Contains(page, "<RelatedPlugins items={RELATED}") {
			t.Errorf("the %s page declares RELATED but never renders it", p.Name)
		}
	}
}

// TestRelatedPluginLinksResolve: a "see also" list quietly rots into dead
// links, and a 404 from a page recommending somewhere to go next is worse
// than not recommending anything.
func TestRelatedPluginLinksResolve(t *testing.T) {
	hrefRe := regexp.MustCompile(`href: '(/docs/[^']+)'`)

	for _, p := range All() {
		rel := "docs/app/docs/plugins/" + p.Name + "/page.tsx"
		page, ok := readRepoFile(t, rel)
		if !ok {
			continue
		}
		block, found := sliceBetween(page, "const RELATED: RelatedPlugin[] = [", "\n]")
		if !found {
			continue
		}
		for _, m := range hrefRe.FindAllStringSubmatch(block, -1) {
			href := m[1]
			target := "docs/app" + href + "/page.tsx"
			if _, exists := readRepoFile(t, target); !exists {
				t.Errorf("the %s page links to %s under Related plugins, "+
					"but %s does not exist", p.Name, href, target)
			}
			if href == "/docs/plugins/"+p.Name {
				t.Errorf("the %s page lists itself as a related plugin", p.Name)
			}
		}
	}
}
