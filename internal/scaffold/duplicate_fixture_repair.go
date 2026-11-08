package scaffold

// A generated handlers test that cannot compile, for a duplicated fixture.
//
// v3.483.0 added a bulk-edit stale-row test and interpolated the parent
// fixtures once per row, so a resource with a belongs_to got:
//
//	file := factory.CreateUpload(t, db)
//	first := factory.CreateAttachment(...)
//	file := factory.CreateUpload(t, db)      <- the second one
//	second := factory.CreateAttachment(...)
//
//	vet: internal/handlers/attachment_test.go:147:7: no new variables on left side of :=
//
// v3.488.0 generates it once. This is for the resources already generated
// between those releases, because `grit upgrade` does not rewrite a generated
// test file: the reader may have added assertions to it, and taking those away
// would be worse than the bug.
//
// Deleting the one duplicated line is safe in a way that rewriting the file is
// not. The file does not compile, so there is nothing in it that currently
// works; both rows want the same parent anyway, which is what the generator now
// emits; and the edit is confined to a line this generator wrote and nobody
// would have written twice on purpose.

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// duplicateFixtureLine is a parent fixture declaration, as the generator emits
// it. Anchored to the tab indent and the factory call so a line somebody wrote
// themselves, with different spacing or a different helper, is not touched.
var duplicateFixtureLine = regexp.MustCompile(`^\t(\w+) := factory\.Create\w+\(t, db\)$`)

// repairDuplicateFixtures removes a redeclared parent fixture from every
// generated handlers test that has one.
func repairDuplicateFixtures(root string, opts Options) error {
	dir := filepath.Join(opts.APIRoot(root), "internal", "handlers")
	entries, err := os.ReadDir(dir)
	if err != nil {
		// No handlers directory is a project shape this does not apply to.
		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		raw, readErr := os.ReadFile(path)
		if readErr != nil {
			continue
		}
		crlf := strings.Contains(string(raw), "\r\n")
		src := strings.ReplaceAll(string(raw), "\r\n", "\n")

		out, removed := dropRedeclaredFixtures(src)
		if removed == 0 {
			continue
		}
		if crlf {
			out = strings.ReplaceAll(out, "\n", "\r\n")
		}
		if err := writeRepaired(path, out); err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		fmt.Printf("  ✓ %s: removed %d %s that redeclared a fixture the test above it "+
			"already made; the package could not compile with %s\n",
			filepath.ToSlash(rel), removed, plural(removed, "line", "lines"),
			plural(removed, "it", "them"))
	}
	return nil
}

// dropRedeclaredFixtures removes a fixture line whose variable is already
// declared in the same function, and reports how many it took out.
//
// Function by function, because the same fixture name in two test functions is
// correct and common.
func dropRedeclaredFixtures(src string) (string, int) {
	lines := strings.Split(src, "\n")
	seen := map[string]bool{}
	keep := make([]string, 0, len(lines))
	removed := 0

	for _, line := range lines {
		// A new top-level declaration ends the previous function's scope.
		if strings.HasPrefix(line, "func ") {
			seen = map[string]bool{}
		}
		if m := duplicateFixtureLine.FindStringSubmatch(line); m != nil {
			if seen[m[1]] {
				removed++
				continue
			}
			seen[m[1]] = true
		}
		keep = append(keep, line)
	}
	if removed == 0 {
		return src, 0
	}
	return strings.Join(keep, "\n"), removed
}
