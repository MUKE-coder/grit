package clierr_test

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/clierr"
)

// The catalogue and the pages are two things that must agree.
//
// The CLI prints https://gritframework.dev/e/<slug> when something fails. A
// slug with no page sends a stuck person to a 404, which is worse than
// printing nothing, and nothing else would notice: the Go side compiles, the
// docs site builds, and the only symptom is on somebody else's screen.
//
// So both sides are read here, in one test. A test per side would pass while
// the pair was broken.
func TestEveryProblemHasAPage(t *testing.T) {
	dir := pagesDir(t)

	onDisk := map[string]bool{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(dir, e.Name(), "page.tsx")); err == nil {
			onDisk[e.Name()] = true
		}
	}

	inCatalogue := map[string]bool{}
	for _, slug := range clierr.Slugs() {
		inCatalogue[slug] = true
	}

	var missingPage, orphanPage []string
	for slug := range inCatalogue {
		if !onDisk[slug] {
			missingPage = append(missingPage, slug)
		}
	}
	for slug := range onDisk {
		if !inCatalogue[slug] {
			orphanPage = append(orphanPage, slug)
		}
	}
	sort.Strings(missingPage)
	sort.Strings(orphanPage)

	if len(missingPage) > 0 {
		t.Errorf("the CLI can print a link to these, and there is no page:\n  %s\n"+
			"Add docs/app/e/<slug>/page.tsx, or remove the entry from internal/clierr.",
			strings.Join(missingPage, "\n  "))
	}
	if len(orphanPage) > 0 {
		t.Errorf("these pages exist and nothing links to them:\n  %s\n"+
			"Add an entry to internal/clierr, or delete the page.",
			strings.Join(orphanPage, "\n  "))
	}

	// A walk that found nothing would pass both checks above by agreeing
	// with an empty catalogue. It cannot be empty.
	if len(onDisk) == 0 {
		t.Fatal("no pages found under docs/app/e: the path is wrong and this test is checking nothing")
	}
}

// Each page must state the same slug it lives at, because the breadcrumb and
// the canonical URL are built from it.
func TestEachPageDeclaresItsOwnSlug(t *testing.T) {
	dir := pagesDir(t)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}
	checked := 0
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		path := filepath.Join(dir, e.Name(), "page.tsx")
		body, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		checked++
		if want := `slug="` + e.Name() + `"`; !strings.Contains(string(body), want) {
			t.Errorf("%s does not contain %s, so its breadcrumb points elsewhere", path, want)
		}
		if want := "/e/" + e.Name() + "'"; !strings.Contains(string(body), want) {
			t.Errorf("%s has no canonical ending %s", path, want)
		}
	}
	if checked == 0 {
		t.Fatal("no pages checked")
	}
}

func pagesDir(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for i := 0; i < 6; i++ {
		candidate := filepath.Join(dir, "docs", "app", "e")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Skip("docs/app/e not found; running outside the repository")
	return ""
}
