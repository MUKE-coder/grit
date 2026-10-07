package scaffold

import (
	"strings"
	"testing"
)

// The web app's error pages use the web app's palette.
//
// They were written with the admin's token names. bg-primary, text-primary and
// text-primary-foreground are not utilities in the web stylesheet, which maps
// --accent, --accent-fg and the bg-* scale, and Tailwind emits nothing for a
// class it cannot resolve. The 404's button rendered as bare text on no
// background, which is the kind of thing only a screenshot catches.
func TestWebErrorPagesUseTheWebPalette(t *testing.T) {
	opts := Options{ProjectName: "contacts", Theme: "emerald", Architecture: ArchTriple, Frontend: FrontendNext}
	pages := map[string]string{
		"error":     webErrorPage(),
		"not found": webNotFoundView(),
		"root 404":  webNotFoundPage(),
	}
	for name, src := range pages {
		for _, token := range []string{"bg-primary", "text-primary", "text-primary-foreground", "bg-card", "text-muted-foreground"} {
			if strings.Contains(src, token) {
				t.Errorf("the %s page uses %s, which the web stylesheet does not define", name, token)
			}
		}
	}

	// One 404, rendered in two places: the root for a URL that matched nothing,
	// the group for notFound() from a page that did match. The root brings the
	// chrome; the group's must not, or it appears twice.
	if !strings.Contains(webNotFoundPage(), "<Navbar />") || !strings.Contains(webNotFoundPage(), "<Footer />") {
		t.Error("the root 404 has no chrome, and nothing above it to provide any")
	}
	group := webMarketingNotFoundPage()
	if strings.Contains(group, "Navbar") || strings.Contains(group, "Footer") {
		t.Error("the marketing 404 draws chrome the marketing layout has already drawn")
	}
	if !strings.Contains(group, "NotFoundView") || !strings.Contains(webNotFoundPage(), "NotFoundView") {
		t.Error("the two 404s do not share one component, so they will drift")
	}

	// global-error replaces the document, so its colours are values, and they
	// have to be this project's values.
	ge := webGlobalErrorPage(opts)
	accent := themeColor("emerald", "accent", "")
	if !strings.Contains(ge, accent) {
		t.Errorf("the global error page does not carry the project's accent (%s)", accent)
	}
	for _, midnight := range []string{"#6c5ce7", "#0a0a0f", "#e8e8f0", `className="dark"`} {
		if strings.Contains(ge, midnight) {
			t.Errorf("the global error page still hardcodes %s from the midnight palette", midnight)
		}
	}

	// And all three are delivered.
	files := webFileMap("root", opts)
	for _, want := range []string{"not-found-view.tsx", "(marketing)"} {
		found := false
		for path := range files {
			if strings.Contains(path, want) && strings.Contains(path, "not-found") {
				found = true
			}
		}
		if !found {
			t.Errorf("no file matching %s + not-found is written", want)
		}
	}
}
