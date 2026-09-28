package scaffold

import (
	"strings"
	"testing"
)

// Nothing a Grit app served said it was one, so a technology scanner could
// identify Next.js, React and Tailwind and stop there. The generator meta is
// the convention every static site generator and CMS uses for exactly this.
func TestFrontendsDeclareTheirGenerator(t *testing.T) {
	opts := Options{ProjectName: "demo"}

	next := map[string]string{
		"apps/web/app/layout.tsx":   webRootLayout(opts),
		"apps/admin/app/layout.tsx": adminRootLayout(opts),
	}
	for name, src := range next {
		if !strings.Contains(src, `generator: "Grit",`) {
			t.Errorf("%s does not declare its generator", name)
		}
	}

	vite := map[string]string{
		"apps/web/index.html":   webTanStackIndexHTML(opts),
		"apps/admin/index.html": adminTanStackIndexHTML(opts),
	}
	for name, src := range vite {
		if !strings.Contains(src, `<meta name="generator" content="Grit" />`) {
			t.Errorf("%s does not declare its generator", name)
		}
	}
}

// The name and not the version. Telling an unauthenticated visitor which
// release is running hands them that release's advisories, which is the same
// reason poweredByHeader is off in the Next config two files away.
func TestTheGeneratorMetaCarriesNoVersion(t *testing.T) {
	for name, src := range map[string]string{
		"web layout":   webRootLayout(Options{ProjectName: "demo"}),
		"admin layout": adminRootLayout(Options{ProjectName: "demo"}),
		"web html":     webTanStackIndexHTML(Options{ProjectName: "demo"}),
		"admin html":   adminTanStackIndexHTML(Options{ProjectName: "demo"}),
	} {
		for _, leak := range []string{"Grit " + DefaultVersion, "content=\"Grit v", "Grit/" + DefaultVersion} {
			if strings.Contains(src, leak) {
				t.Errorf("%s leaks the version through the generator meta: %s", name, leak)
			}
		}
	}
	// And the API still names nothing: a header that says Grit would
	// contradict the reason the Next one is switched off.
	if strings.Contains(apiRoutesGo(), `"X-Powered-By"`) {
		t.Error("the API advertises itself in a header")
	}
}
