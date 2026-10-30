package main

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// Four places say what version this is, and nothing compared them.
//
// Enterprise review 2.26 found the fourth: release.yml passed
// -X main.Version, and the variable is `version`, lowercase. The ldflag set
// nothing, silently, so every released binary reported whatever the source
// constant said. That happened to be right, because the release ritual bumps
// it by hand, which is exactly the kind of correctness that holds until the
// one time somebody forgets.
//
// The ldflag is gone and this test is the replacement. A test rather than a CI
// step because the release ritual runs locally: finding out from a red build
// after the tag is pushed is finding out too late, and a tag is the one thing
// that is awkward to take back.
func TestEveryVersionSpotAgrees(t *testing.T) {
	root := repoRoot(t)

	spots := []struct {
		name    string
		path    string
		pattern string
	}{
		{"internal/scaffold DefaultVersion", filepath.Join("internal", "scaffold", "scaffold.go"),
			`(?m)^const DefaultVersion = "([0-9][^"]*)"`},
		{"docs/config/site.ts GRIT_VERSION", filepath.Join("docs", "config", "site.ts"),
			`(?m)^export const GRIT_VERSION = '([0-9][^']*)'`},
	}

	for _, spot := range spots {
		raw, err := os.ReadFile(filepath.Join(root, spot.path))
		if err != nil {
			t.Fatalf("reading %s: %v", spot.path, err)
		}
		m := regexp.MustCompile(spot.pattern).FindSubmatch(raw)
		if m == nil {
			t.Fatalf("%s: could not find the version in %s. The pattern this "+
				"test matches on has moved, so it has stopped checking anything",
				spot.name, spot.path)
		}
		if got := string(m[1]); got != version {
			t.Errorf("%s says %s, cmd/grit says %s.\n\n"+
				"A release bumps all of these together. Whichever is wrong, the "+
				"three have to match before the tag is pushed",
				spot.name, got, version)
		}
	}
}

// The changelog has an entry for the version being released, because a release
// whose changelog stops one version short is one nobody can read.
func TestTheChangelogMentionsThisVersion(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join(repoRoot(t), "docs", "app", "docs", "changelog", "page.tsx"))
	if err != nil {
		t.Fatalf("reading the changelog: %v", err)
	}
	want := `id="v` + version + `"`
	if !regexp.MustCompile(regexp.QuoteMeta(want)).Match(raw) {
		t.Errorf("the changelog has no entry for v%s (looking for %s)", version, want)
	}
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
	t.Fatal("could not find the repository root from the test's directory")
	return ""
}
