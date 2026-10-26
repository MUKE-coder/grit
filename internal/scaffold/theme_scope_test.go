package scaffold

import (
	"regexp"
	"strings"
	"testing"
)

// A theme scoped to a surface rather than to <html> did not repaint.
//
// Tailwind declares the @theme mappings on :root, and a custom property that
// references another one resolves where it is declared, so a section with its
// own data-theme redefined --bg-primary and still inherited the --color-background
// that had already been resolved at the root. The tokens changed and
// bg-background did not (grit#98).
//
// Both stylesheets, because the mapping was written out twice and the fix had
// to land in both. It is one function now, and this is the test that says so.
func TestScopedThemeReDeclaresTheTokenMappings(t *testing.T) {
	for _, sheet := range []struct {
		name string
		css  string
	}{
		{"admin", adminGlobalCSS()},
		{"web", webGlobalCSS()},
	} {
		root := declarationsIn(t, sheet.css, "@theme {")
		if len(root) == 0 {
			t.Fatalf("%s: no @theme block", sheet.name)
		}
		scoped := declarationsIn(t, sheet.css, "[data-theme] [data-theme],")
		if len(scoped) == 0 {
			t.Fatalf("%s: the mappings are not re-declared for a scoped theme, so "+
				"data-theme on anything but <html> changes the tokens and not the colours", sheet.name)
		}
		for name, value := range root {
			if scoped[name] != value {
				t.Errorf("%s: %s is mapped at the root but not for a scoped theme (%q vs %q)",
					sheet.name, name, value, scoped[name])
			}
		}
		for name := range scoped {
			if _, ok := root[name]; !ok {
				t.Errorf("%s: %s is re-declared for a scoped theme but never mapped at the root",
					sheet.name, name)
			}
		}
	}
}

// And the scoped selector covers the dark-mode signals the toggle writes, not
// just data-theme: a nested .dark has the identical problem.
func TestScopedThemeCoversTheDarkSignals(t *testing.T) {
	for _, sheet := range []struct {
		name string
		css  string
	}{
		{"admin", adminGlobalCSS()},
		{"web", webGlobalCSS()},
	} {
		for _, want := range []string{
			"[data-theme] [data-theme],",
			"[data-theme] .dark,",
			`[data-theme] [data-theme-mode="dark"] {`,
		} {
			if !strings.Contains(sheet.css, want) {
				t.Errorf("%s: the scoped selector is missing %q", sheet.name, want)
			}
		}
	}
}

var cssDeclRe = regexp.MustCompile(`^\s*(--[\w-]+)\s*:\s*(.+?);\s*$`)

// declarationsIn reads the custom properties out of the block that opens with
// the given line.
func declarationsIn(t *testing.T, css, opener string) map[string]string {
	t.Helper()
	i := strings.Index(css, opener)
	if i < 0 {
		return nil
	}
	rest := css[i:]
	end := strings.Index(rest, "\n}")
	if end < 0 {
		t.Fatalf("block %q is never closed", opener)
	}
	out := map[string]string{}
	for _, line := range strings.Split(rest[:end], "\n") {
		if m := cssDeclRe.FindStringSubmatch(line); m != nil {
			out[m[1]] = m[2]
		}
	}
	return out
}
