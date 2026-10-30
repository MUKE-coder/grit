package generate

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Review 2.24: the README's five-minute tutorial generated Ticket, which the
// generator refuses because the support desk ships a Ticket model.
//
// The refusal message is good: it names the feature that holds the name. The
// problem is that the tutorial is the first thing anybody runs, and it stopped
// on step 2 of six.
//
// Two things that must agree, with nothing comparing them: the reserved list
// in this package, and every `grit generate resource` in the documentation. A
// name added to the list tomorrow could break a tutorial written today, and
// nobody would find out until somebody followed it.
func TestNoDocumentedResourceUsesAReservedName(t *testing.T) {
	root := repoRootForNouns(t)

	// Every place a reader is told to run the generator. Not the changelog:
	// it records what commands were, including ones that have since become
	// reserved, and rewriting history to suit a new reservation would be
	// worse than leaving it.
	files := []string{
		"README.md",
		filepath.Join("docs", "public", "skill.md"),
	}
	// Plus every page of the docs site, which is the specification a user
	// reads and the place a tutorial actually lives.
	docsRoot := filepath.Join(root, "docs", "app", "docs")
	_ = filepath.WalkDir(docsRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || d.Name() != "page.tsx" {
			return nil
		}
		if strings.Contains(filepath.ToSlash(path), "/changelog/") {
			return nil
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr == nil {
			files = append(files, rel)
		}
		return nil
	})
	if len(files) < 20 {
		t.Fatalf("found %d files to check, so this test is looking in the wrong place", len(files))
	}

	// `grit generate resource Name` and `grit g r Name`, with the name captured.
	cmd := regexp.MustCompile(`grit (?:generate resource|g r) ([A-Z][A-Za-z0-9]*)`)

	var found int
	for _, rel := range files {
		raw, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			continue
		}
		for _, m := range cmd.FindAllStringSubmatch(string(raw), -1) {
			found++
			name := m[1]
			if err := CheckReservedName(name); err != nil {
				t.Errorf("%s tells the reader to run `grit generate resource %s`, "+
					"and the generator refuses it: %v. A tutorial that stops on its "+
					"own second step is worse than no tutorial (review 2.24)",
					filepath.ToSlash(rel), name, err)
			}
		}
	}
	if found < 5 {
		t.Errorf("found only %d documented generate commands, so this test is probably "+
			"matching nothing: check the pattern against the docs", found)
	}
}

func repoRootForNouns(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatal("could not find the repository root")
	return ""
}
