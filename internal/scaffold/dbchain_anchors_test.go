package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// Every GORM callback the framework registers is an anchor a project may name,
// and dbchain's list of them has to be the real spellings.
//
// The list is hand-written because GORM offers no way to ask what is
// registered. Two entries were wrong on the first attempt:
// appendonly:refuse_update for appendonly:update. An anchor spelled the way
// the list said would have passed dbchain's validation and then been ignored
// by GORM, because GORM accepts an After naming a callback it has never seen
// and simply drops the constraint. That is precisely the failure the list
// exists to prevent, inside the list.
//
// So the two sides are read in one test. One grep over everything the
// framework ships, one parse of the list, and a comparison, rather than a test
// per side that both pass while the pair is broken.
func TestDBChainKnowsEveryFrameworkCallback(t *testing.T) {
	registered := frameworkCallbackNames(t)
	if len(registered) < 8 {
		t.Fatalf("found only %d framework callbacks (%v); this test is reading for the "+
			"wrong thing, and a test that finds nothing passes", len(registered), registered)
	}

	known := dbchainKnownAnchors(t)
	if len(known) == 0 {
		t.Fatal("could not read knownAnchors out of the dbchain template")
	}

	for _, name := range registered {
		if !known[name] {
			t.Errorf("the framework registers %q and dbchain's knownAnchors does not have it, "+
				"so a link anchored to it is refused although it is real", name)
		}
	}
	// And the other direction, which is the one that bit: a name in the list
	// that nothing registers lets an anchor through and then does nothing.
	have := map[string]bool{}
	for _, name := range registered {
		have[name] = true
	}
	for name := range known {
		if strings.HasPrefix(name, "gorm:") {
			continue // GORM's own, not ours to find
		}
		if !have[name] {
			t.Errorf("dbchain's knownAnchors has %q and nothing in this framework registers it, "+
				"so an anchor spelled that way passes the check and GORM then ignores it", name)
		}
	}
}

var callbackRegisterRe = regexp.MustCompile(`Register\("([a-z_]+:[a-z_]+)"`)

// frameworkCallbackNames greps every Register("ns:name") out of the scaffold's
// Go sources and its templates. Test callbacks are left out: they are
// registered inside a test and removed again.
func frameworkCallbackNames(t *testing.T) []string {
	t.Helper()
	found := map[string]bool{}
	for _, dir := range []string{".", "templates"} {
		err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil //nolint:nilerr // an unreadable file is not this test's subject
			}
			name := d.Name()
			if !strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, ".tmpl") {
				return nil
			}
			if strings.HasSuffix(name, "_test.go") || strings.Contains(name, "_test.go.tmpl") {
				return nil
			}
			raw, err := os.ReadFile(path) // #nosec G304 -- this package's own files
			if err != nil {
				return nil //nolint:nilerr
			}
			for _, m := range callbackRegisterRe.FindAllStringSubmatch(string(raw), -1) {
				if strings.HasPrefix(m[1], "test:") || strings.HasPrefix(m[1], "dbchain:") {
					continue
				}
				found[m[1]] = true
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walking %s: %v", dir, err)
		}
	}
	out := make([]string, 0, len(found))
	for name := range found {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// dbchainKnownAnchors reads the names out of the var block in the template,
// which is the file a generated project compiles.
func dbchainKnownAnchors(t *testing.T) map[string]bool {
	t.Helper()
	src := tmpl("api/dbchain/dbchain.go")
	const head = "var knownAnchors = map[string]bool{"
	i := strings.Index(src, head)
	if i < 0 {
		t.Fatal("no knownAnchors in the dbchain template")
	}
	rest := src[i+len(head):]
	end := strings.Index(rest, "\n}")
	if end < 0 {
		t.Fatal("knownAnchors is not closed")
	}
	out := map[string]bool{}
	for _, m := range regexp.MustCompile(`"([a-z_]+:[a-z_]+)":\s*true`).FindAllStringSubmatch(rest[:end], -1) {
		out[m[1]] = true
	}
	return out
}
