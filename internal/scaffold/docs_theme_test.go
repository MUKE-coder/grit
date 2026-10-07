package scaffold

import (
	"strings"
	"testing"
)

// A --full project has three frontends, and THEME in .env paints all three.
//
// The docs site ignored it: fumadocs' neutral palette, a hardcoded dark class,
// and no reader of THEME anywhere. A project generated with --theme emerald had
// a green web app, a green admin panel and a black-and-white docs site.
func TestDocsSiteWearsTheProjectTheme(t *testing.T) {
	opts := Options{ProjectName: "contacts", Architecture: ArchTriple, Frontend: FrontendNext}

	css := docsGlobalCSS()
	if !strings.Contains(css, `[data-theme="emerald"]`) {
		t.Error("the docs stylesheet does not carry the palettes")
	}
	// Mapped onto fumadocs' own tokens, or the palette is loaded and unused.
	for _, token := range []string{"--color-fd-background: var(--bg-primary)", "--color-fd-primary: var(--accent)", "--color-fd-border: var(--border)"} {
		if !strings.Contains(css, token) {
			t.Errorf("fumadocs' %s is not mapped onto the palette", token)
		}
	}
	// After the imports: equal specificity, so the last declaration wins.
	if strings.Index(css, "--color-fd-background: var(--bg-primary)") < strings.Index(css, "fumadocs-ui/css/preset.css") {
		t.Error("the mapping comes before fumadocs' own defaults, which overwrite it")
	}

	layout := docsRootLayout(opts)
	if !strings.Contains(layout, `data-theme={theme}`) || !strings.Contains(layout, "NEXT_PUBLIC_THEME") {
		t.Error("the docs layout does not put the theme on <html>")
	}
	if strings.Contains(layout, `className="dark"`) {
		t.Error("the docs layout still forces dark, whatever the project's theme is")
	}
	if !strings.Contains(layout, "enabled: false") {
		t.Error("fumadocs' light/dark switch is still on, and it changes nothing")
	}

	cfg := docsNextConfig()
	if !strings.Contains(cfg, "NEXT_PUBLIC_THEME: process.env.THEME") {
		t.Error("the docs build never reads THEME")
	}
	if !strings.Contains(cfg, "rootEnv") {
		t.Error("THEME lives in the monorepo's root .env, which the docs app does not load on its own")
	}

	// And the site is the project's, not the framework's.
	home := docsHomePage(opts)
	if !strings.Contains(home, "opts.ProjectName") && !strings.Contains(home, "contacts") {
		t.Error("the docs home page never names the project")
	}
	if strings.Contains(home, "build with Grit") {
		t.Error("the docs home page is still Grit's own marketing")
	}
	if strings.Contains(home, "—") {
		t.Error("em dash in generated copy")
	}
}
