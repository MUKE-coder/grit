package scaffold

import (
	"strings"
	"testing"
)

// A generated project's home page is about that project.
//
// It shipped as Grit's own marketing instead: an <h1> reading "Grit", a
// paragraph describing the meta-framework, and a "Get Started" button pointing
// at gritframework.dev. Every --double, --triple and --full project advertised
// somebody else's framework on its front page, while the TanStack frontend's
// version of the same page used the project's name. These check the two agree.
func TestLandingPageIsAboutTheProject(t *testing.T) {
	for _, tc := range []struct {
		name string
		src  string
	}{
		{"next", webLandingPage(Options{ProjectName: "contacts", Architecture: ArchTriple, Frontend: FrontendNext})},
		{"tanstack", webTanStackIndexRoute(Options{ProjectName: "contacts", Architecture: ArchSingle, Frontend: FrontendTanStack})},
	} {
		if !strings.Contains(tc.src, "contacts") {
			t.Errorf("the %s home page never names the project", tc.name)
		}
		for _, grit := range []string{"gritframework.dev", "gritframework.com", "github.com/MUKE-coder/grit", "grit new "} {
			if strings.Contains(tc.src, grit) {
				t.Errorf("the %s home page sends a visitor to %s", tc.name, grit)
			}
		}
		if strings.Contains(tc.src, "meta-framework") {
			t.Errorf("the %s home page describes the framework, not the project", tc.name)
		}
		// On the amber theme --accent-fg is near-black, so a label hardcoded to
		// white on the accent button is unreadable there.
		if strings.Contains(tc.src, "bg-accent px-6 py-3 text-sm font-semibold text-white") {
			t.Errorf("the %s home page hardcodes a white label on the accent button", tc.name)
		}
	}
	// "Read the Blog" is a 404 once the Blog resource is removed, so it has to
	// be inside markers that grit remove can cut.
	next := webLandingPage(Options{ProjectName: "contacts", Architecture: ArchTriple, Frontend: FrontendNext})
	for _, marker := range []string{"{/* grit:home:blog-cta-start */}", "{/* grit:home:blog-cta-end */}"} {
		if !strings.Contains(next, marker) {
			t.Errorf("the home page's blog CTA is not wrapped in %s", marker)
		}
	}
}
