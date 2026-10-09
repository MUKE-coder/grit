package docs

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// The index ships inside the binary, so the failure worth guarding against is
// a stale one: it would answer confidently in the voice of a version it is not.

// The embedded index must be for this CLI's version.
//
// scripts/docs-index.py reads the same `var version` line, so a release that
// forgets to rebuild the index fails here rather than shipping documentation
// for the previous version to everybody who installs it.
func TestTheIndexIsForThisVersion(t *testing.T) {
	raw, err := os.ReadFile("../../cmd/grit/main.go")
	if err != nil {
		t.Fatalf("reading main.go: %v", err)
	}
	match := regexp.MustCompile(`var version = "([^"]+)"`).FindSubmatch(raw)
	if match == nil {
		t.Fatal("cannot find `var version` in cmd/grit/main.go")
	}
	want := string(match[1])

	if got := Version(); got != want {
		t.Fatalf("the embedded docs index is for v%s and this CLI is v%s.\n"+
			"Run: python scripts/docs-index.py", got, want)
	}
}

// An index that quietly emptied itself would make every search return nothing,
// which reads as "the docs do not cover this".
func TestTheIndexIsPopulated(t *testing.T) {
	if n := Count(); n < 100 {
		t.Fatalf("%d pages indexed; the docs site has far more, so the extractor or the "+
			"index is broken", n)
	}
}

// Every page has to be worth returning.
func TestEveryPageHasSomethingToShow(t *testing.T) {
	thin := 0
	for _, url := range URLs() {
		page, err := Read(url)
		if err != nil {
			t.Fatalf("%s: %v", url, err)
		}
		if page.Title == "" {
			t.Errorf("%s has no title", url)
		}
		if len(page.Text) < 200 && page.Description == "" {
			thin++
		}
	}
	// A couple of card-grid index pages have no prose of their own, which is
	// fine. A lot of them would mean the extractor stopped working.
	if thin > 5 {
		t.Fatalf("%d pages have neither prose nor a description", thin)
	}
}

// The ranking, as a set of questions with a right answer.
//
// These are the questions an agent actually asks, and the assertion is the page
// a person would send them to. It is a judgement encoded as a test, which is
// the only way ranking stays good: without it, a scoring change that helps one
// query silently ruins four.
func TestSearchFindsTheRightPage(t *testing.T) {
	cases := []struct {
		query string
		want  string
	}{
		{"field types", "/docs/concepts/field-types"},
		{"money field", "/docs/concepts/money"},
		{"admin custom table", "/docs/admin/custom-pages"},
		{"offline sync desktop", "/docs/desktop/offline"},
		{"roles and permissions", "/docs/admin/roles"},
		{"deploy to a vps", "/docs/deployment/vps"},
	}

	for _, tc := range cases {
		hits, err := Search(tc.query, 3)
		if err != nil {
			t.Errorf("%q: %v", tc.query, err)
			continue
		}
		if len(hits) == 0 {
			t.Errorf("%q found nothing", tc.query)
			continue
		}
		if hits[0].URL != tc.want {
			var got []string
			for _, h := range hits {
				got = append(got, h.URL)
			}
			t.Errorf("%q ranked %s first, want %s (top 3: %s)",
				tc.query, hits[0].URL, tc.want, strings.Join(got, ", "))
		}
	}
}

// The changelog must be findable and must not win.
//
// It has 807 headings, one per release, so it mentions every word in the
// vocabulary. Counting those as statements about the page put it top of every
// multi-word query, above the page actually about the thing. That is what
// heading weighting by specificity is for, and this is the test that keeps it.
func TestAnAggregatePageDoesNotWinEveryQuery(t *testing.T) {
	hits, err := Search("generate a resource owned by user", 5)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(hits) == 0 {
		t.Fatal("found nothing")
	}
	if hits[0].URL == "/docs/changelog" {
		t.Fatal("the changelog is top of a query about generating a resource, so heading " +
			"matches are being counted without regard to how many headings there are")
	}

	// And it is still reachable when it is the right answer.
	hits, err = Search("changelog", 3)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	found := false
	for _, h := range hits {
		if h.URL == "/docs/changelog" {
			found = true
		}
	}
	if !found {
		t.Error("the changelog is no longer findable by name")
	}
}

// A reference page beats a post about the same subject, because a post is
// dated narrative and a page is the current contract.
func TestReferenceOutranksNarrative(t *testing.T) {
	hits, err := Search("roles and permissions", 5)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(hits) == 0 {
		t.Fatal("found nothing")
	}
	if hits[0].Kind != "docs" {
		t.Fatalf("a %s page ranked above every doc page for a reference question: %s",
			hits[0].Kind, hits[0].URL)
	}
}

// A snippet is what tells a reader whether to open the page, so it has to
// contain the thing they searched for.
func TestTheSnippetShowsTheMatch(t *testing.T) {
	hits, err := Search("offline sync", 1)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(hits) == 0 {
		t.Fatal("found nothing")
	}
	snippet := strings.ToLower(hits[0].Snippet)
	if snippet == "" {
		t.Fatal("no snippet")
	}
	if !strings.Contains(snippet, "offline") && !strings.Contains(snippet, "sync") {
		t.Fatalf("the snippet does not contain either search term: %q", hits[0].Snippet)
	}
	if len(hits[0].Snippet) > 320 {
		t.Errorf("the snippet is %d characters, which is a page and not a snippet", len(hits[0].Snippet))
	}
}

// Reading a page by URL, forgivingly, because an agent will get the slash wrong.
func TestReadIsForgivingAboutSlashes(t *testing.T) {
	for _, url := range []string{
		"/docs/concepts/field-types",
		"docs/concepts/field-types",
		"/docs/concepts/field-types/",
		"/DOCS/concepts/field-types",
	} {
		if _, err := Read(url); err != nil {
			t.Errorf("Read(%q): %v", url, err)
		}
	}
}

// A near miss suggests, rather than refusing, because an agent told "no such
// page" concludes the documentation does not cover it.
func TestAWrongURLSuggests(t *testing.T) {
	_, err := Read("/docs/concepts")
	if err == nil {
		t.Skip("there is a page at /docs/concepts, so this is not a near miss")
	}
	if !strings.Contains(err.Error(), "did you mean") {
		t.Fatalf("a near miss gave no suggestions: %v", err)
	}

	_, err = Read("/docs/this-does-not-exist-at-all-xyzzy")
	if err == nil {
		t.Fatal("a page that does not exist was found")
	}
}

// An empty query is a mistake worth naming rather than a search that returns
// the whole index.
func TestAnEmptyQueryIsRefused(t *testing.T) {
	for _, q := range []string{"", "   ", "the a of"} {
		if _, err := Search(q, 5); err == nil {
			t.Errorf("Search(%q) was accepted", q)
		}
	}
}

// Two runs of the same query return the same order, so a diff of agent output
// does not churn.
func TestRankingIsStable(t *testing.T) {
	first, err := Search("admin table", 5)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	second, err := Search("admin table", 5)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(first) != len(second) {
		t.Fatalf("%d hits then %d", len(first), len(second))
	}
	for i := range first {
		if first[i].URL != second[i].URL {
			t.Fatalf("position %d was %s then %s", i, first[i].URL, second[i].URL)
		}
	}
}
