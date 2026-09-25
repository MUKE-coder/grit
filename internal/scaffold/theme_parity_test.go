package scaffold

import (
	"regexp"
	"sort"
	"strings"
	"testing"
)

// THEME=<name> paints every surface, or it paints none of them.
//
// The five themes v3.322.0 added went into the admin's stylesheet only, so
// `grit new shop --theme emerald` produced an emerald admin behind an Atlas
// marketing site. Nothing failed: a [data-theme] block that does not exist
// leaves :root's values in place, which is the default theme, which looks
// deliberate.
func TestEveryThemePaintsTheWebAppToo(t *testing.T) {
	block := regexp.MustCompile(`\[data-theme="([a-z]+)"\]`)

	names := func(src string) []string {
		seen := map[string]bool{}
		var out []string
		for _, m := range block.FindAllStringSubmatch(src, -1) {
			if !seen[m[1]] {
				seen[m[1]] = true
				out = append(out, m[1])
			}
		}
		sort.Strings(out)
		return out
	}

	admin := names(adminGlobalCSS())
	web := names(webGlobalCSS())

	if len(admin) == 0 {
		t.Fatal("the admin stylesheet defines no themes at all")
	}
	if strings.Join(admin, ",") != strings.Join(web, ",") {
		t.Errorf("the web app is missing themes the admin has.\n admin: %v\n   web: %v", admin, web)
	}

	// Every theme the CLI offers has to be one of them, or --theme accepts a
	// name that paints nothing.
	for _, name := range ValidThemes {
		found := false
		for _, have := range web {
			if have == name || name == "atlas" {
				found = true
			}
		}
		if !found {
			t.Errorf("--theme %s is accepted but the web app has no block for it", name)
		}
	}

	// A button's label is a token, not text-white: on a light accent the two
	// are the same colour.
	if !strings.Contains(webGlobalCSS(), "--color-accent-fg: var(--accent-fg);") {
		t.Error("the web app has no accent-fg utility, so an accent button's label cannot be themed")
	}
	for _, name := range web {
		i := strings.Index(webGlobalCSS(), `[data-theme="`+name+`"]`)
		if i < 0 {
			continue
		}
		rest := webGlobalCSS()[i:]
		rest = rest[:strings.Index(rest, "}")]
		if !strings.Contains(rest, "--accent-fg:") {
			t.Errorf("web theme %q never sets --accent-fg, so its buttons fall back to the body colour", name)
		}
	}
}
