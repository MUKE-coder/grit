package scaffold

import (
	"strings"
	"testing"
)

// The logo mark on an auth screen has to be visible on that screen.
//
// It defaulted to white text on fifteen percent white, which works on a dark
// hero panel. Six of the eight themes have a light one: aurora, coral, amber,
// sky, mono and emerald all showed an empty square where the letter was.
func TestAuthShellBrandMarksAreLegible(t *testing.T) {
	shells := map[string]string{
		"admin atlas":    adminAtlasAuthShell(),
		"admin aurora":   adminAuroraAuthShell(),
		"desktop atlas":  desktopAtlasAuthShell(),
		"desktop aurora": desktopAuroraAuthShell(),
		"desktop pulse":  desktopPulseAuthShell(),
		"desktop mark":   desktopClientBrandMark(),
	}
	for name, src := range shells {
		if strings.Contains(src, `background: "rgba(255,255,255,0.15)"`) {
			t.Errorf("%s still paints the mark with a wash of white", name)
		}
		if strings.Contains(src, "<BrandMark />") {
			t.Errorf("%s renders a mark with no colour", name)
		}
	}
	// Every mark that takes a background takes the letter's colour with it.
	for name, src := range shells {
		for _, call := range []string{"<BrandMark color={t.primary} />", "<BrandMark tint={t.primary} />", "<BrandMark accent={t.accent} />"} {
			if strings.Contains(src, call) {
				t.Errorf("%s: %s sets a background and leaves the letter hardcoded", name, call)
			}
		}
	}
}
