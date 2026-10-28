package blueprint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The website's copy of the list. Two pages read it (/blueprints, the
// gallery, and /docs/blueprints, the explanation), so comparing this one
// file with the registry covers both.
const blueprintsPage = "docs/config/blueprints.ts"

// TestBlueprintsPageListsEveryBlueprint compares the registry with the docs
// page.
//
// The website cannot import this package, so it holds its own copy of the list.
// Two lists with nothing comparing them drift, and this one drifts quietly: a
// blueprint ships, `grit blueprint list` shows it, the docs page does not, and
// the page still looks complete.
func TestBlueprintsPageListsEveryBlueprint(t *testing.T) {
	page, ok := readRepoFile(t, blueprintsPage)
	if !ok {
		t.Skipf("%s not found; this test only runs inside the Grit repository", blueprintsPage)
	}

	for _, b := range All() {
		needle := `name: '` + b.Name + `'`
		if !strings.Contains(page, needle) {
			t.Errorf("the %q blueprint is in the registry but missing from %s", b.Name, blueprintsPage)
		}
		if !strings.Contains(page, b.Summary) {
			t.Errorf("%s describes %q differently from the registry.\nregistry: %s",
				blueprintsPage, b.Name, b.Summary)
		}
	}
}

// TestBlueprintsPageHasNoBlueprintsThatDoNotExist is the other direction.
func TestBlueprintsPageHasNoBlueprintsThatDoNotExist(t *testing.T) {
	page, ok := readRepoFile(t, blueprintsPage)
	if !ok {
		t.Skipf("%s not found; this test only runs inside the Grit repository", blueprintsPage)
	}

	known := map[string]bool{}
	for _, n := range Names() {
		known[n] = true
	}

	rest := page
	for {
		at := strings.Index(rest, "name: '")
		if at < 0 {
			return
		}
		rest = rest[at+len("name: '"):]
		end := strings.Index(rest, "'")
		if end < 0 {
			return
		}
		name := rest[:end]
		rest = rest[end:]
		if !known[name] {
			t.Errorf("%s lists a blueprint %q that the registry does not have, so "+
				"`grit blueprint info %s` would fail for anybody who followed the page",
				blueprintsPage, name, name)
		}
	}
}

// TestEveryBlueprintSaysWhatItTeaches: a catalogue entry whose only
// description is "an app built with Grit" helps nobody choose.
func TestEveryBlueprintSaysWhatItTeaches(t *testing.T) {
	for _, b := range All() {
		if len(b.Teaches) == 0 {
			t.Errorf("%s lists nothing under Teaches", b.Name)
		}
		if b.Summary == "" {
			t.Errorf("%s has no summary", b.Name)
		}
		if b.Access != Planned && b.Repo == "" {
			t.Errorf("%s is listed as available but names no repository", b.Name)
		}
		if b.Access == Planned && b.Repo != "" {
			t.Errorf("%s is planned but names a repository, which would 404", b.Name)
		}
	}
}

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
