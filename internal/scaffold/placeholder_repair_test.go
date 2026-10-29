package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// A repair that writes a template without substituting its placeholders
// writes a file that does not compile.
//
// repairBenchTestGo did exactly that: it replaced {{PROJECT}} and not
// {{MODULE}}, so an upgrade left a literal "{{MODULE}}/internal/config"
// import in internal/handlers/bench_test.go and the project stopped
// building. Caught by running an upgrade on a real project, which is the
// only place it shows.
//
// This reads the repair sources for the shape: a call that renders a
// template function and replaces some placeholders and not the others. It
// cannot see every way of getting this wrong, and it sees the one that
// happened.
var (
	repairPlaceholderRe = regexp.MustCompile(`\{\{[A-Z_]+\}\}`)
)

func TestRepairsSubstituteEveryPlaceholderTheyWrite(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}

	checked := 0
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, "_repair.go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(".", name))
		if err != nil {
			t.Fatal(err)
		}
		src := string(raw)

		// Which template functions this repair renders.
		for _, call := range regexp.MustCompile(`\b(api\w+Go|k6\w+)\(\)`).FindAllString(src, -1) {
			fn := strings.TrimSuffix(call, "()")
			body, ok := templateBodyFor(t, fn)
			if !ok {
				continue
			}
			want := repairPlaceholderRe.FindAllString(body, -1)
			if len(want) == 0 {
				continue
			}
			checked++
			seen := map[string]bool{}
			for _, p := range want {
				seen[p] = true
			}
			for p := range seen {
				if !strings.Contains(src, `"`+p+`"`) {
					t.Errorf("%s renders %s, whose template contains %s, and never "+
						"substitutes it. The file it writes will carry the literal "+
						"placeholder and the project will not compile.", name, call, p)
				}
			}
		}
	}

	if checked == 0 {
		t.Error("no repair was found rendering a template with placeholders; " +
			"the reader stopped matching and this test checks nothing")
	}
}

// templateBodyFor renders the named template function, when it is one this
// test knows how to call. Deliberately a small list: a template needing
// Options is rendered with an empty one, which is enough to see its
// placeholders.
func templateBodyFor(t *testing.T, fn string) (string, bool) {
	t.Helper()
	switch fn {
	case "apiBenchTestGo":
		return apiBenchTestGo(), true
	case "apiAuthTestGo":
		return apiAuthTestGo(), true
	case "apiSSOHooksGo":
		return apiSSOHooksGo(), true
	case "apiOAuthHooksGo":
		return apiOAuthHooksGo(), true
	case "k6CommonJS":
		return k6CommonJS(), true
	case "k6SmokeJS":
		return k6SmokeJS(), true
	}
	return "", false
}
