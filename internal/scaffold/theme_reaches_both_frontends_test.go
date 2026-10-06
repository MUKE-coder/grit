package scaffold

import (
	"strings"
	"testing"
)

// The theme flag reaches a Vite app as well as a Next one.
//
// Vite exposes only VITE_-prefixed variables to client code. The env file wrote
// THEME, which the Go binary and the Next apps read, and nothing a TanStack app
// could see: import.meta.env.VITE_THEME was undefined, so getTheme fell back to
// atlas. The dashboard still took its colours from the stylesheet baked at
// scaffold time, and the sign-in screen took its layout from the fallback, so
// --theme emerald produced an emerald dashboard behind an indigo atlas login.
//
// The same flag and the same theme name gave two different answers depending on
// which frontend was asked, which is the kind of difference nobody reports as a
// bug: they assume that is what the theme looks like.
func TestTheThemeReachesBothFrontends(t *testing.T) {
	for _, theme := range []string{"emerald", "coral", "atlas"} {
		env := envFile(Options{
			ProjectName:  "app",
			Architecture: ArchSingle,
			Theme:        theme,
		})

		if !strings.Contains(env, "\nTHEME="+theme) {
			t.Errorf("%s: THEME is not set, so the Go side and the Next apps cannot see it", theme)
		}
		if !strings.Contains(env, "\nVITE_THEME="+theme) {
			t.Errorf("%s: VITE_THEME is not set, so a Vite app falls back to atlas", theme)
		}
	}

	// And the two never disagree, which would be worse than either being wrong.
	env := envFile(Options{ProjectName: "app", Architecture: ArchTriple, Frontend: FrontendNext, Theme: "pulse"})
	if strings.Count(env, "=pulse") != 2 {
		t.Errorf("the two theme variables disagree:\n%s", themeLines(env))
	}
}

// themeLines pulls just the theme settings out of an env file, for an error
// message that fits on a screen.
func themeLines(env string) string {
	var out []string
	for _, line := range strings.Split(env, "\n") {
		if strings.HasPrefix(line, "THEME=") || strings.HasPrefix(line, "VITE_THEME=") {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n")
}
