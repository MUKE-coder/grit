package scaffold

import (
	"strings"
	"testing"
)

// The mobile app wears the project's theme like the other three frontends.
//
// It carried #6c5ce7 as a literal in its Tailwind config and its JavaScript
// palette, so --theme emerald shipped a green web app, a green admin, a green
// desktop client and a purple phone app. NativeWind resolves classes at build
// time and has no CSS variables, so the values are written at generation time.
func TestExpoAppWearsTheProjectTheme(t *testing.T) {
	emerald := Options{ProjectName: "contacts", Theme: "emerald", IncludeExpo: true}
	accent := themeColor("emerald", "accent", "")
	if accent == "" {
		t.Fatal("the emerald palette has no accent")
	}

	tw := expoTailwindConfig(emerald)
	if !strings.Contains(tw, accent) {
		t.Errorf("the mobile Tailwind config does not carry %s", accent)
	}
	if strings.Contains(tw, "#6c5ce7") {
		t.Error("the mobile Tailwind config still carries the old purple")
	}

	provider := expoThemeProvider(emerald)
	if strings.Count(provider, accent) < 4 {
		t.Errorf("the mobile palette uses %s in %d places, want the logo, the gradients and both refresh spinners", accent, strings.Count(provider, accent))
	}
	if strings.Contains(provider, "#6c5ce7") {
		t.Error("the mobile palette still carries the old purple")
	}

	// A project with no --theme is atlas, and atlas is not purple either.
	none := expoTailwindConfig(Options{ProjectName: "contacts"})
	if !strings.Contains(none, themeColor("atlas", "accent", "")) {
		t.Error("an unthemed project does not get the default palette")
	}
}

// One palette, read by the surfaces that cannot share a stylesheet.
func TestThemeVarsReadsTheSharedPalette(t *testing.T) {
	for _, name := range ValidThemes {
		vars := themeVars(name)
		for _, want := range []string{"accent", "accent-hover", "bg-primary", "text-primary", "border", "success", "danger"} {
			if vars[want] == "" {
				t.Errorf("%s has no --%s", name, want)
			}
		}
	}
	// An unknown name falls back rather than returning nothing: a colour that
	// is the empty string is a Tailwind build error and an invisible element.
	if themeVars("no-such-theme")["accent"] != themeVars("atlas")["accent"] {
		t.Error("an unknown theme does not fall back to atlas")
	}
	if got := rgba("#059669", 0.1); got != "rgba(5, 150, 105, 0.1)" {
		t.Errorf("rgba: %s", got)
	}
}
