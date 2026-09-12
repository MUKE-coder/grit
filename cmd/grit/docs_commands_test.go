package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/docscheck"
)

// docsRoot is the docs site, relative to this package.
func docsRoot(t *testing.T) string {
	t.Helper()
	root := filepath.Join("..", "..", "docs", "app")
	if _, err := os.Stat(root); err != nil {
		t.Skipf("docs site not present: %v", err)
	}
	return root
}

// Every grit command the docs show has to exist, with the flags the docs give it.
//
// "The docs said to do X and X did not work" is the most repeated line in this
// project's changelog. Each time it was found by a person following a page, which
// is the most expensive way to find it. This asks the command tree instead.
func TestDocsCommandsExist(t *testing.T) {
	samples, err := docscheck.Extract(docsRoot(t))
	if err != nil {
		t.Fatalf("reading the docs: %v", err)
	}
	if len(samples) < 300 {
		// The extractor silently reading nothing, or less than it used to, is the
		// failure mode that would make this test pass forever while checking
		// nothing. The floor is below today's count and above the count before the
		// extractor learned to read the first line of a block.
		t.Fatalf("only %d commands found in the docs; the extractor has probably stopped matching", len(samples))
	}
	t.Logf("checked %d commands from the docs", len(samples))

	for _, problem := range docscheck.Check(samples, rootCommand()) {
		t.Errorf("%s", problem)
	}
}

// And every field spec the docs show has to parse, because a tutorial whose
// generate command the generator refuses is worse than a missing page.
func TestDocsFieldSpecsParse(t *testing.T) {
	samples, err := docscheck.Extract(docsRoot(t))
	if err != nil {
		t.Fatalf("reading the docs: %v", err)
	}
	for _, problem := range docscheck.CheckFieldSpecs(samples) {
		t.Errorf("%s", problem)
	}
}
