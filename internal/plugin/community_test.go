package plugin

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"
)

var dateOnly = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// TestCommunityIndexParses is the one that matters, because the index is
// edited by hand in a pull request from somebody who is not running the
// tests. A trailing comma ships a binary whose community list is silently
// empty: Community() swallows the unmarshal error on purpose, since an
// empty list is the honest answer at runtime, which means this test is the
// only thing standing between a typo and a feature that quietly disappears.
func TestCommunityIndexParses(t *testing.T) {
	var doc struct {
		Comment []string          `json:"$comment"`
		Plugins []CommunityPlugin `json:"plugins"`
	}
	if err := json.Unmarshal(CommunityJSON(), &doc); err != nil {
		t.Fatalf("internal/plugin/community.json does not parse: %v", err)
	}
	if len(doc.Comment) == 0 {
		t.Error("community.json has lost its $comment, which is where the submission rules are written")
	}
}

func TestCommunityEntriesAreComplete(t *testing.T) {
	seen := map[string]bool{}
	for _, p := range Community() {
		if p.Name == "" {
			t.Error("a community entry has no name")
			continue
		}
		if seen[p.Name] {
			t.Errorf("%s is listed twice", p.Name)
		}
		seen[p.Name] = true

		if _, builtIn := Get(p.Name); builtIn == nil {
			t.Errorf("the community entry %q has the same name as a built-in plugin, "+
				"so `grit plugin add %s` would install the built-in and the author's "+
				"plugin would be unreachable", p.Name, p.Name)
		}
		if p.Author == "" {
			t.Errorf("%s has no author", p.Name)
		}
		if p.Summary == "" {
			t.Errorf("%s has no summary, so `grit plugin list` would print a blank line", p.Name)
		}
		if !IsRepoRef(p.Repo) {
			t.Errorf("%s has repo %q, which `grit plugin add` would not recognise as a repository",
				p.Name, p.Repo)
		}
		if !dateOnly.MatchString(p.Added) {
			t.Errorf("%s has added %q, want YYYY-MM-DD", p.Name, p.Added)
		}
	}
}

func TestIsRepoRef(t *testing.T) {
	repos := []string{
		"https://github.com/someone/grit-calendar",
		"http://example.com/x.git",
		"git@github.com:someone/grit-calendar.git",
		"github.com/someone/grit-calendar",
		"gitlab.com/someone/thing",
		"codeberg.org/someone/thing",
	}
	for _, r := range repos {
		if !IsRepoRef(r) {
			t.Errorf("IsRepoRef(%q) = false, want true", r)
		}
	}

	// A bare user/repo is deliberately NOT a repository reference: it is
	// ambiguous with a plugin name, and resolving the ambiguity toward "fetch
	// and run somebody's code" is the wrong default.
	notRepos := []string{"kanban", "saved-views", "./plugins/mine", "/abs/path", "someone/grit-calendar"}
	for _, r := range notRepos {
		if IsRepoRef(r) {
			t.Errorf("IsRepoRef(%q) = true, want false", r)
		}
	}
}

func TestRepoURLAddsAScheme(t *testing.T) {
	if got := RepoURL("github.com/a/b"); got != "https://github.com/a/b" {
		t.Errorf("RepoURL(github.com/a/b) = %q", got)
	}
	for _, already := range []string{"https://github.com/a/b", "git@github.com:a/b.git"} {
		if got := RepoURL(already); got != already {
			t.Errorf("RepoURL(%q) = %q, want it unchanged", already, got)
		}
	}
}

// TestSubmitAndMarketplaceURLsAgree: the CLI points people at the docs page
// and the docs page is generated from this file, so a rename that breaks one
// should break the test rather than the link.
func TestSubmitAndMarketplaceURLsAgree(t *testing.T) {
	if !strings.HasPrefix(MarketplaceURL, "https://gritframework.dev/docs/plugins/") {
		t.Errorf("MarketplaceURL = %q, which is not under the plugins docs", MarketplaceURL)
	}
	if !strings.Contains(SubmitURL, "plugin-submission") {
		t.Errorf("SubmitURL = %q, which does not name the submission template", SubmitURL)
	}
}
