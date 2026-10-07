package scaffold

import (
	"strings"
	"testing"
)

// No em dashes in the docs a project ships.
//
// The house standard bans them in anything a reader sees. The generated docs
// had seventy-odd, mostly in bulleted definitions where a colon reads better.
func TestGeneratedDocsHaveNoEmDashes(t *testing.T) {
	opts := Options{ProjectName: "contacts", Architecture: ArchTriple, Frontend: FrontendNext, Version: "3.0.0"}
	pages := map[string]string{
		"index":           docsContentIndex(opts),
		"getting started": docsContentGettingStarted(opts),
		"auth":            docsContentAuth(),
		"cli":             docsContentCLI(),
		"migrations":      docsContentMigrationsSeeding(),
		"admin overview":  docsContentAdminOverview(),
		"admin resources": docsContentAdminResources(),
		"batteries":       docsContentBatteries(),
		"about":           docsContentAbout(),
		"home page":       docsHomePage(opts),
	}
	for name, src := range pages {
		if i := strings.Index(src, "—"); i >= 0 {
			line := src[maxInt(0, i-60):minInt(len(src), i+60)]
			t.Errorf("%s has an em dash: ...%s...", name, strings.ReplaceAll(line, "\n", " "))
		}
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
