package docs

import (
	"strings"
	"testing"
)

// The command's own help prints `grit docs show /docs/concepts/field-types`,
// and in Git Bash on Windows that argument never reaches the process intact:
// MSYS rewrites anything that looks like an absolute POSIX path, so the CLI
// receives "C:/Program Files/Git/docs/concepts/field-types" and answered
// "no page at ...", which reads as a missing page rather than a mangled
// argument (grit#96).
func TestReadRecoversAMangledAbsolutePath(t *testing.T) {
	for _, arg := range []string{
		"/docs/concepts/field-types",
		"docs/concepts/field-types",
		"C:/Program Files/Git/docs/concepts/field-types",
		`C:\Program Files\Git\docs\concepts\field-types`,
		"/c/Program Files/Git/docs/concepts/field-types",
		"//docs/concepts/field-types",
	} {
		page, err := Read(arg)
		if err != nil {
			t.Errorf("Read(%q): %v", arg, err)
			continue
		}
		if page.URL != "/docs/concepts/field-types" {
			t.Errorf("Read(%q) returned %q", arg, page.URL)
		}
	}
}

// Twenty indexed pages are under /blog, so the recovery cannot be anchored on
// a "docs" segment.
func TestReadRecoversABlogPage(t *testing.T) {
	var blog string
	for _, u := range URLs() {
		if strings.HasPrefix(u, "/blog/") {
			blog = u
			break
		}
	}
	if blog == "" {
		t.Skip("no blog page in the index")
	}
	page, err := Read("C:/Program Files/Git" + blog)
	if err != nil {
		t.Fatalf("Read of a mangled blog path: %v", err)
	}
	if page.URL != blog {
		t.Errorf("got %q, want %q", page.URL, blog)
	}
}

// A real miss must still get the suggestion list. The recovery is allowed to
// rescue a mangled path and not to invent a match for a wrong one.
func TestReadStillRefusesAPageThatDoesNotExist(t *testing.T) {
	if _, err := Read("/docs/concepts/there-is-no-such-page"); err == nil {
		t.Fatal("expected an error for a page that does not exist")
	}

	_, err := Read("/docs/concepts/field-typos")
	if err == nil {
		t.Fatal("expected an error for a near miss")
	}
	if !strings.Contains(err.Error(), "did you mean") {
		t.Errorf("a near miss should suggest; got %v", err)
	}
}

// The suffix match takes the single page it is sure about. An argument that
// could be two pages gets the suggestion list instead, because a confident
// wrong page is worse than a question.
func TestSuffixMatchRefusesToGuessBetweenTwo(t *testing.T) {
	pages := []Page{
		{URL: "/docs/a/thing"},
		{URL: "/docs/b/thing"},
	}
	if got := bySuffix(&pages, "thing"); got != nil {
		t.Errorf("ambiguous suffix should match nothing, got %q", got.URL)
	}
	if got := bySuffix(&pages, "docs/a/thing"); got == nil || got.URL != "/docs/a/thing" {
		t.Error("an unambiguous suffix should match")
	}
}

// A suffix has to land on a segment boundary, or "field-types" would match
// "my-field-types".
func TestSuffixMatchRespectsSegmentBoundaries(t *testing.T) {
	pages := []Page{{URL: "/docs/concepts/types"}}
	if got := bySuffix(&pages, "docs/concepts/field-types"); got != nil {
		t.Errorf("a partial segment should not match, got %q", got.URL)
	}
}
