package plugin

import (
	"os"
	"path/filepath"
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
	official, found := sliceBetween(page, "const OFFICIAL = [", "\n]")
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
