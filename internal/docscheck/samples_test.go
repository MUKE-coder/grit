package docscheck

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// docsRoot is the docs site, found from this file rather than from the working
// directory, so the test runs the same from anywhere.
func docsRoot(t *testing.T) string {
	t.Helper()
	_, here, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate this test file")
	}
	return filepath.Join(filepath.Dir(here), "..", "..", "docs")
}

// The homepage is the first thing an agent reads, before the skill. For months
// it showed a handler running its own query, a float64 price, a dropped bind
// error and a hand-built error envelope, which are four of the things the skill
// tells an agent never to do. Found by a product reviewer in October 2026, not
// by the build, because nothing checked the code the way the commands were
// already checked.
func TestDocsSamplesFollowTheFrameworksOwnRules(t *testing.T) {
	problems, err := CheckSamples(docsRoot(t))
	if err != nil {
		t.Fatalf("scanning the docs: %v", err)
	}
	if len(problems) == 0 {
		return
	}
	var b strings.Builder
	b.WriteString("the docs show Go that grit generate would never write:\n\n")
	for _, p := range problems {
		b.WriteString("  " + p.String() + "\n\n")
	}
	b.WriteString("Every sample on the site has to be something the generator would produce.\n")
	b.WriteString("An agent reads the homepage before it reads the skill.")
	t.Error(b.String())
}

// Each rule catches what it is for, and the fixed form passes, or the check is
// either decoration or a nuisance.
func TestSampleRulesCatchTheRealShapes(t *testing.T) {
	cases := []struct {
		name string
		rule string
		bad  string
		good string
	}{
		{
			name: "query in a handler",
			rule: "handler-query",
			bad:  "    h.DB.Where(\"featured = ?\", true).Find(&products)",
			good: "    products, err := h.service().Featured(h.ctx(c))",
		},
		{
			name: "float money in a model",
			rule: "float-money",
			bad:  "    Price      float64   `json:\"price\"`",
			good: "    Price      money.Money `json:\"price\"`",
		},
		{
			name: "float money on the command line",
			rule: "float-money",
			bad:  "  grit generate resource Product --fields \"name:string,price:float\"",
			good: "  grit generate resource Product --fields \"name:string,price:money\"",
		},
		{
			name: "dropped bind error",
			rule: "dropped-bind",
			bad:  "    c.ShouldBindJSON(&req)",
			good: "    if err := c.ShouldBindJSON(&req); err != nil {",
		},
		{
			name: "hand-built error envelope",
			rule: "hand-built-error",
			bad:  "        c.JSON(401, gin.H{\"error\": \"Invalid credentials\"})",
			good: "        respond.Error(c, err)",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			page := "<CodeBlock language=\"go\" code={`\n" + tc.bad + "\n`} />"
			got := checkSampleFile("page.tsx", page)
			if len(got) != 1 || got[0].Rule != tc.rule {
				t.Fatalf("the bad sample was not caught by %s: %v", tc.rule, got)
			}
			fixed := "<CodeBlock language=\"go\" code={`\n" + tc.good + "\n`} />"
			if got := checkSampleFile("page.tsx", fixed); len(got) != 0 {
				t.Errorf("the fixed sample was reported anyway: %v", got)
			}
		})
	}
}

// A rule that fires on prose, or on a TypeScript block that happens to contain
// the word price, is a rule somebody switches off within a week.
func TestSampleRulesDoNotFireOutsideGo(t *testing.T) {
	for _, page := range []string{
		// Prose explaining the rule.
		"<p>Never use <code>Price float64</code>: money is integer minor units.</p>",
		// A TypeScript block after a Go one.
		"<CodeBlock language=\"go\" code={`type P struct{}`} />\n" +
			"<CodeBlock language=\"ts\" code={`const price: number = 0`} />\n" +
			"c.ShouldBindJSON(&req)",
		// A shell block showing the money form.
		"<CodeBlock language=\"bash\" code={`grit generate resource P --fields price:money`} />",
	} {
		if got := checkSampleFile("page.tsx", page); len(got) != 0 {
			t.Errorf("fired outside a Go sample: %v\n  in: %s", got, page)
		}
	}
}
