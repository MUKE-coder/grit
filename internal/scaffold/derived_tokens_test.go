package scaffold

import (
	"regexp"
	"strings"
	"testing"
)

// Every token the Tailwind map points at must actually be declared.
//
// The map and the palette are two lists that have to agree, and nothing makes
// them. A --color-x pointing at a --y nobody declares compiles without
// complaint and renders as nothing: a pill with no fill, a border that is not
// there. The failure is silent and it is invisible in exactly the themes
// nobody opened.
func TestEveryMappedTokenIsDeclared(t *testing.T) {
	palette := themePaletteCSS()

	declared := map[string]bool{}
	for _, m := range regexp.MustCompile(`(?m)^\s*(--[a-z0-9-]+)\s*:`).FindAllStringSubmatch(palette, -1) {
		declared[m[1]] = true
	}

	refs := regexp.MustCompile(`--color-[a-z0-9-]+\s*:\s*var\((--[a-z0-9-]+)\)`)
	seen := 0
	for _, m := range refs.FindAllStringSubmatch(themeTokenMap, -1) {
		seen++
		if !declared[m[1]] {
			t.Errorf("the Tailwind map points at %s, which no theme and no derivation declares, "+
				"so every utility built on it renders nothing", m[1])
		}
	}
	if seen < 20 {
		t.Fatalf("only %d mapped colours found; the map's shape changed and this "+
			"test is no longer reading it", seen)
	}
}

// The derived tokens have to resolve against the theme they are used in.
//
// Declared on :root alone they resolve once, against whatever theme the root
// carries, and a section with its own data-theme inherits that resolved value
// while the --accent it was mixed from changes underneath it. That is grit#98,
// and it looks like a pill keeping the old theme's tint on a screen where
// everything else changed.
func TestDerivedTokensResolveAgainstTheirOwnTheme(t *testing.T) {
	css := derivedTokenCSS()

	if !strings.Contains(css, "[data-theme] {") {
		t.Error("the derived block does not select [data-theme], so a section with " +
			"its own theme inherits tints mixed from the root's accent")
	}
	if !strings.Contains(css, ":root,") {
		t.Error("the derived block does not cover :root, so a surface with no " +
			"data-theme attribute gets none of these tokens")
	}
}

// A tinted fill mixes into transparent, not into a surface.
//
// These sit on a card on one screen and on the page on the next. Mixed into
// --bg-elevated a pill carries the card's colour with it, which is visible the
// moment one lands on the page background.
func TestTintedFillsAreTransparent(t *testing.T) {
	css := derivedTokenCSS()

	for _, name := range []string{
		"--accent-soft", "--accent-line",
		"--success-soft", "--success-line",
		"--warning-soft", "--warning-line",
		"--danger-soft", "--danger-line",
		"--info-soft", "--info-line",
	} {
		line := tokenLine(css, name)
		if line == "" {
			t.Errorf("%s is not declared", name)
			continue
		}
		if !strings.Contains(line, "transparent") {
			t.Errorf("%s mixes into a surface (%s), so it carries that surface's "+
				"colour onto any other background it lands on", name, strings.TrimSpace(line))
		}
	}
}

// Every theme declares what the derivations read.
//
// The derived block mixes --accent, --success, --warning, --danger, --info,
// --border, --text-muted, --text-primary, --bg-hover and --bg-primary. A theme
// missing one of them produces an invalid color-mix, which a browser drops: the
// token is then unset and whatever used it renders nothing.
func TestEveryThemeDeclaresWhatTheDerivationsRead(t *testing.T) {
	needed := []string{
		"accent", "success", "warning", "danger", "info",
		"border", "text-muted", "text-primary", "bg-hover", "bg-primary",
	}
	themes := append([]string{}, ValidThemes...)
	themes = append(themes, "midnight")

	for _, theme := range themes {
		vars := themeVars(theme)
		if len(vars) == 0 {
			t.Errorf("theme %q has no palette block", theme)
			continue
		}
		for _, name := range needed {
			if vars[name] == "" {
				t.Errorf("theme %q declares no --%s, so every token derived from it "+
					"is an invalid color-mix the browser drops", theme, name)
			}
		}
	}
}

// tokenLine returns the declaration of one custom property.
func tokenLine(css, name string) string {
	for _, line := range strings.Split(css, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), name+":") {
			return line
		}
	}
	return ""
}
