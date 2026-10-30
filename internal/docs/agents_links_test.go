package docs

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Review 2.25: AGENTS.md linked to three files that are not in this
// repository.
//
// AGENTS.md is tracked. GRIT.md, PHASES.md and CLAUDE.md are in .gitignore,
// deliberately, alongside the plans and the social calendar. So the first table
// in the file whose whole job is to orient somebody who just cloned Grit sent
// them to three dead links, and a reviewer working from a clone hit exactly
// that.
//
// Checked against what git tracks rather than what is on disk, which is the
// only check that matters: every one of those files exists here, which is why
// nobody noticed.
func TestEveryAgentsLinkIsInTheRepository(t *testing.T) {
	root := repoRoot(t)
	tracked := trackedFiles(t, root)

	agents := make([]string, 0, 8)
	for path := range tracked {
		if filepath.Base(path) == "AGENTS.md" {
			agents = append(agents, path)
		}
	}
	if len(agents) < 5 {
		t.Fatalf("found %d AGENTS.md files in the repository, so this test is looking in the wrong place", len(agents))
	}

	// Markdown links, minus anchors, absolute URLs and bare fragments.
	link := regexp.MustCompile(`\]\(([^)]+)\)`)
	for _, rel := range agents {
		raw, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatal(err)
		}
		dir := filepath.ToSlash(filepath.Dir(rel))
		for _, m := range link.FindAllStringSubmatch(string(raw), -1) {
			target := m[1]
			if strings.HasPrefix(target, "http") || strings.HasPrefix(target, "#") || strings.Contains(target, "://") {
				continue
			}
			if i := strings.IndexAny(target, "#?"); i >= 0 {
				target = target[:i]
			}
			if target == "" {
				continue
			}
			joined := target
			if dir != "." {
				joined = dir + "/" + target
			}
			joined = strings.TrimSuffix(filepath.ToSlash(filepath.Clean(joined)), "/")

			if tracked[joined] || tracked[joined+"/"] {
				continue
			}
			// A directory link resolves if the repository tracks anything
			// inside it.
			found := false
			for p := range tracked {
				if strings.HasPrefix(p, joined+"/") {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("%s links to %q, which this repository does not contain. "+
					"Anybody who clones Grit gets a dead link in the file that is "+
					"meant to orient them (review 2.25)", rel, m[1])
			}
		}
	}
}

// trackedFiles is every path git tracks, as a set of forward-slash paths.
func trackedFiles(t *testing.T, root string) map[string]bool {
	t.Helper()
	cmd := exec.Command("git", "ls-files")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		t.Skipf("git ls-files is not available here: %v", err)
	}
	set := make(map[string]bool)
	for _, line := range strings.Split(string(out), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			set[line] = true
		}
	}
	if len(set) == 0 {
		t.Skip("git tracks no files here")
	}
	return set
}

// repoRoot walks up from the test's directory to the module root.
func repoRoot(t *testing.T) string {
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
