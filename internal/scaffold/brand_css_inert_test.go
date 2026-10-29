package scaffold

import (
	"regexp"
	"strings"
	"testing"
)

// brand.css must declare nothing when it is written.
//
// The whole promise of the file is that an upgrade can create it in an
// existing project and change nothing about how that project looks. It was
// read as a repaint once (grit#135) because the values sit in a block whose
// comment markers are on their own lines, and then a note added inside that
// block ended the comment early: CSS comments do not nest, so the closing
// brace became live CSS. Four bytes, in a file that is not allowed to have
// any.
var cssComment = regexp.MustCompile(`(?s)/\*.*?\*/`)

func TestBrandCSSDeclaresNothing(t *testing.T) {
	cases := map[string]struct{ light, dark map[string]string }{
		"an app with light and dark palettes": {
			light: map[string]string{"accent": "#41608f", "bg-primary": "#ffffff"},
			dark:  map[string]string{"accent": "#a99bff", "bg-primary": "#08090b"},
		},
		"an app with no dark palette": {
			light: map[string]string{"accent": "#41608f"},
			dark:  map[string]string{},
		},
		"an app whose stylesheet could not be read": {
			light: map[string]string{},
			dark:  map[string]string{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			out := brandCSS("example", tc.light, tc.dark)

			live := strings.TrimSpace(cssComment.ReplaceAllString(out, ""))
			if live != "" {
				t.Errorf("brand.css declares something, so creating it changes how "+
					"the app looks (grit#135). Live after stripping comments:\n%s", live)
			}

			// Unbalanced markers are how the live bytes got there.
			if opens, closes := strings.Count(out, "/*"), strings.Count(out, "*/"); opens != closes {
				t.Errorf("%d /* and %d */: a comment is left open or closed twice", opens, closes)
			}
		})
	}
}

// Only colours the app actually uses. A value from Grit's palette in a file
// headed "your colours" is what made the report read as a repaint.
func TestBrandCSSShowsOnlyTheAppsOwnColours(t *testing.T) {
	out := brandCSS("example",
		map[string]string{"accent": "#41608f"},
		map[string]string{})

	if !strings.Contains(out, "#41608f") {
		t.Error("the example does not show the app's own accent")
	}
	for _, grit := range []string{"#5b47e0", "#a99bff", "#4a38c9"} {
		if strings.Contains(out, grit) {
			t.Errorf("the example shows Grit's %s, which this app does not use "+
				"(grit#135)", grit)
		}
	}
}
