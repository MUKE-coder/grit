// Package docs is the documentation, searchable, embedded in the binary and
// pinned to the version that built it.
//
// An agent writing Grit code has to find the right API for the version the
// project is on. Without this it reads whatever the docs site currently says,
// which is a different version: Grit ships several releases a week, so an
// agent working in a project pinned to v3.350 and reading v3.393 docs is
// confidently wrong, and the symptom is generated code calling a helper that
// did not exist yet or a flag that has since been renamed.
//
// Embedded rather than fetched for two reasons. An installed binary cannot read
// the Next.js app the docs are written in, and a developer on a train has no
// network. Embedding also makes the index inseparable from the version that
// produced it, which is the whole point: Version() is a fact about these pages
// rather than a guess.
//
// The index is built by scripts/docs-index.py and committed. Stale is the
// failure worth guarding against, so there is a test that the embedded version
// matches the CLI's.
package docs

import (
	"embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"unicode"
)

//go:embed index.json
var embedded embed.FS

// Page is one documentation page.
type Page struct {
	// URL is the path on the docs site, e.g. /docs/admin/custom-pages. It is
	// also the identifier Read takes, because a URL is the thing an agent can
	// hand back to a human.
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Headings    []string `json:"headings"`
	Text        string   `json:"text"`
	// Kind is docs, blog or lesson. A blog post explains why and a doc page
	// says how, and an agent choosing between them benefits from knowing.
	Kind string `json:"kind"`
}

type index struct {
	Version string `json:"version"`
	Pages   []Page `json:"pages"`
}

var (
	once    sync.Once
	loaded  index
	loadErr error
)

func load() (*index, error) {
	once.Do(func() {
		raw, err := embedded.ReadFile("index.json")
		if err != nil {
			loadErr = fmt.Errorf("docs: the index is missing from the binary: %w", err)
			return
		}
		if err := json.Unmarshal(raw, &loaded); err != nil {
			loadErr = fmt.Errorf("docs: the index will not parse: %w", err)
		}
	})
	return &loaded, loadErr
}

// Version is the Grit version these pages document.
func Version() string {
	idx, err := load()
	if err != nil {
		return ""
	}
	return idx.Version
}

// Count is how many pages are indexed, for `grit doctor` and for a test that
// notices the index emptying itself.
func Count() int {
	idx, err := load()
	if err != nil {
		return 0
	}
	return len(idx.Pages)
}

// Hit is one search result.
type Hit struct {
	Page
	// Score is only meaningful against the other hits in the same search.
	Score int `json:"score"`
	// Snippet is the prose around the first match, which is what tells a
	// reader whether this is the page they wanted without opening it.
	Snippet string `json:"snippet"`
}

// Search finds the pages most likely to answer a question.
//
// Deliberately a keyword score and not an embedding. The index has to fit in a
// binary and run with no model and no network, and for "which page covers
// field types" a title and heading match is the right answer anyway. The thing
// that matters is that the ranking is stable and the snippet is honest.
func Search(query string, limit int) ([]Hit, error) {
	idx, err := load()
	if err != nil {
		return nil, err
	}
	terms := terms(query)
	if len(terms) == 0 {
		return nil, fmt.Errorf("docs: nothing to search for")
	}
	if limit <= 0 {
		limit = 5
	}

	hits := make([]Hit, 0, 16)
	for _, page := range idx.Pages {
		score, where := score(page, terms)
		if score == 0 {
			continue
		}
		hits = append(hits, Hit{Page: page, Score: score, Snippet: where})
	}

	sort.Slice(hits, func(i, j int) bool {
		if hits[i].Score != hits[j].Score {
			return hits[i].Score > hits[j].Score
		}
		// A stable tiebreak, so two runs of the same query do not reorder and
		// make a diff of agent output look like a change.
		return hits[i].URL < hits[j].URL
	})

	if len(hits) > limit {
		hits = hits[:limit]
	}
	return hits, nil
}

// Read returns one page by URL, with or without the leading slash.
func Read(url string) (*Page, error) {
	idx, err := load()
	if err != nil {
		return nil, err
	}
	want := normaliseURL(url)
	for i := range idx.Pages {
		if normaliseURL(idx.Pages[i].URL) == want {
			return &idx.Pages[i], nil
		}
	}
	// Suggest, rather than just refusing. A near miss on a URL is the common
	// case and the alternative is an agent concluding the page does not exist.
	near := make([]string, 0, 4)
	for i := range idx.Pages {
		if strings.Contains(normaliseURL(idx.Pages[i].URL), want) || strings.Contains(want, normaliseURL(idx.Pages[i].URL)) {
			near = append(near, idx.Pages[i].URL)
		}
	}
	if len(near) > 0 {
		sort.Strings(near)
		if len(near) > 5 {
			near = near[:5]
		}
		return nil, fmt.Errorf("docs: no page at %q; did you mean %s", url, strings.Join(near, ", "))
	}
	return nil, fmt.Errorf("docs: no page at %q", url)
}

// URLs is every indexed page, for listing.
func URLs() []string {
	idx, err := load()
	if err != nil {
		return nil
	}
	out := make([]string, 0, len(idx.Pages))
	for i := range idx.Pages {
		out = append(out, idx.Pages[i].URL)
	}
	sort.Strings(out)
	return out
}

func normaliseURL(u string) string {
	return strings.ToLower(strings.Trim(u, "/"))
}

// terms splits a query into lowercase words, dropping the ones that match
// everything.
func terms(query string) []string {
	var out []string
	for _, word := range strings.FieldsFunc(strings.ToLower(query), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-'
	}) {
		if len(word) < 2 || stop[word] {
			continue
		}
		out = append(out, word)
	}
	return out
}

// stop words, kept short on purpose. A long list starts removing the words that
// carry the question: "how do I generate a resource" is mostly stop words and
// "generate" and "resource" are enough.
var stop = map[string]bool{
	"a": true, "an": true, "the": true, "is": true, "are": true, "was": true,
	"to": true, "of": true, "in": true, "on": true, "for": true, "and": true,
	"or": true, "do": true, "does": true, "how": true, "what": true, "my": true,
	"it": true, "that": true, "this": true, "with": true, "can": true, "i": true,
}

// score weighs where a term was found.
//
// A term in the title is a much better signal than the same term in the body,
// because a page's title is what it is about and its body mentions everything
// nearby. Without that weighting a search for "permissions" returns the thirty
// pages that mention permissions in passing above the one page about them.
func score(page Page, queryTerms []string) (int, string) {
	title := strings.ToLower(page.Title)
	desc := strings.ToLower(page.Description)
	text := strings.ToLower(page.Text)
	heads := strings.ToLower(strings.Join(page.Headings, " \n "))
	url := strings.ToLower(page.URL)

	total := 0
	matched := 0
	// strong counts the terms found somewhere that says what a page is about,
	// rather than somewhere it merely mentions.
	strong := 0
	snippetAt := -1

	for _, term := range queryTerms {
		found := false
		here := false
		if strings.Contains(title, term) {
			total += 40
			found, here = true, true
		}
		if strings.Contains(url, term) {
			total += 25
			found, here = true, true
		}
		if strings.Contains(heads, term) {
			// Weighted by how specific the heading set is.
			//
			// A term among five headings says what the page is about. Among
			// eight hundred it says nothing: the changelog has 807 headings
			// because it has an entry per release, so it matched the headings
			// of every query and sat at the top of all of them. An aggregate
			// page should be findable and should not win.
			total += headingWeight(len(page.Headings))
			found = true
			if len(page.Headings) <= specificHeadings {
				here = true
			}
		}
		if strings.Contains(desc, term) {
			total += 10
			found, here = true, true
		}
		if at := strings.Index(text, term); at >= 0 {
			// Term frequency, capped. An index page that lists a word forty
			// times is not forty times more relevant than the page about it.
			count := strings.Count(text, term)
			if count > 5 {
				count = 5
			}
			total += count
			found = true
			if snippetAt < 0 || at < snippetAt {
				snippetAt = at
			}
		}
		if found {
			matched++
		}
		if here {
			strong++
		}
	}

	if matched == 0 {
		return 0, ""
	}

	// Coverage is the strongest signal there is, and it has to be earned in a
	// position that says what the page is about.
	//
	// Counting a body mention towards it put the changelog top of "generate a
	// resource with owned by user": it is one enormous page that mentions
	// every word in the vocabulary, so it collected full coverage on any query
	// and beat the page actually about the thing. An aggregate page should be
	// findable and should not win.
	if strong == len(queryTerms) && len(queryTerms) > 1 {
		total += 30 * strong
	} else if strong > 0 && matched == len(queryTerms) && len(queryTerms) > 1 {
		// Partial credit: every term is present and some of them in a telling
		// place.
		total += 10 * strong
	}

	// An aggregate page is a catalogue and ranks below a page about the thing.
	//
	// The heading weighting above handles multi-word queries. A single-word one
	// has no coverage bonus to earn, so the changelog still won "owned-by" on
	// raw mentions: it has an entry per release and therefore mentions every
	// feature Grit has. It should be the second answer, not the first.
	if len(page.Headings) > specificHeadings {
		total = total / 2
	}

	// A reference page outranks a post about the same thing.
	//
	// Both belong in the index: "why does placement work this way" is answered
	// by a post and "what are the field types" by a page. But a post is dated
	// narrative and a doc page is the current contract, and an agent about to
	// write code wants the contract. Ten per cent is enough to break a tie
	// without burying the posts.
	if page.Kind != "docs" {
		total = total * 9 / 10
	}

	return total, snippet(page.Text, snippetAt)
}

// specificHeadings is where a heading list stops being a statement about the
// page and becomes a catalogue of everything in it.
const specificHeadings = 60

// headingWeight decays with the size of the heading set.
//
// Fifteen for a focused page, and approaching nothing for one with hundreds of
// sections, which is what an index or a changelog looks like.
func headingWeight(count int) int {
	if count <= 0 {
		return 0
	}
	w := 15 * 20 / (20 + count)
	if w < 1 {
		return 1
	}
	return w
}

// snippet is the prose around a match, cut at word boundaries.
func snippet(text string, at int) string {
	if text == "" {
		return ""
	}
	const width = 240
	if at < 0 {
		if len(text) <= width {
			return text
		}
		return trimToWords(text[:width]) + "..."
	}

	start := at - width/3
	if start < 0 {
		start = 0
	}
	end := start + width
	if end > len(text) {
		end = len(text)
	}
	out := text[start:end]
	if start > 0 {
		out = "..." + trimLeadingPartialWord(out)
	}
	if end < len(text) {
		out = trimToWords(out) + "..."
	}
	return out
}

func trimToWords(s string) string {
	if i := strings.LastIndexByte(s, ' '); i > 0 {
		return s[:i]
	}
	return s
}

func trimLeadingPartialWord(s string) string {
	if i := strings.IndexByte(s, ' '); i >= 0 && i < 20 {
		return s[i+1:]
	}
	return s
}
