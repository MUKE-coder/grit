package errorcodes

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The docs table is generated, so the only way it can be wrong is by being stale.
//
// A catalogue entry added without re-running the generator would leave the page
// documenting an API that returns one more code than it lists, which is the exact
// failure the catalogue exists to prevent.
func TestDocsTableIsCurrent(t *testing.T) {
	path := filepath.Join("..", "..", "docs", "components", "error-codes.ts")
	have, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("docs site not present: %v", err)
	}
	if string(have) != DocsTypeScript() {
		t.Errorf("%s is out of date: run go run ./tools/errorcodes", filepath.ToSlash(path))
	}
}

// And the page that renders it has to name every area, or a whole group of codes
// is generated into the data file and shown nowhere.
func TestDocsPageRendersEveryArea(t *testing.T) {
	path := filepath.Join("..", "..", "docs", "app", "docs", "backend", "errors", "page.tsx")
	page, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("errors page not present: %v", err)
	}
	// The page renders from errorCodeAreas rather than listing areas itself, which
	// is what makes it complete by construction. This asserts that it still does.
	for _, want := range []string{"errorCodeAreas", "errorCodeCount"} {
		if !strings.Contains(string(page), want) {
			t.Errorf("the errors page does not use %s, so it can fall behind the catalogue", want)
		}
	}
}
