package ui

import (
	"math"
	"strconv"
	"strings"
	"testing"
)

/*
 * The contrast claim, as a test.
 *
 * ui.go has always said "every text role clears 5.1:1 on black, #0D1117,
 * #1E1E1E, One Dark, Solarized Dark, white, #FAFAFA, #EEEEEE and Solarized
 * Light". Nothing checked it, so it was a promise about whichever hues
 * happened to be in the file when somebody last did the arithmetic by hand.
 *
 * It matters more here than on a web page. A terminal has no zoom, no user
 * stylesheet and no high-contrast mode to fall back on: if the colour is
 * unreadable on somebody's background, the output is unreadable, and they
 * will not know whether it is their terminal or our CLI.
 *
 * 5.1 rather than WCAG's 4.5 because terminal fonts are thin at small sizes
 * and the margin is cheap.
 */

// The nine backgrounds ui.go names. Five dark, four light.
var darkBackgrounds = map[string]string{
	"black":          "#000000",
	"GitHub dark":    "#0D1117",
	"VS Code dark":   "#1E1E1E",
	"One Dark":       "#282C34",
	"Solarized Dark": "#002B36",
}

var lightBackgrounds = map[string]string{
	"white":           "#FFFFFF",
	"near-white":      "#FAFAFA",
	"light grey":      "#EEEEEE",
	"Solarized Light": "#FDF6E3",
}

// minRatio is the floor every text role has to clear.
const minRatio = 5.1

func channel(v float64) float64 {
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

// luminance is WCAG relative luminance for a #rrggbb string.
func luminance(t *testing.T, hex string) float64 {
	t.Helper()
	h := strings.TrimPrefix(hex, "#")
	if len(h) != 6 {
		t.Fatalf("%q is not a #rrggbb colour", hex)
	}
	var c [3]float64
	for i := 0; i < 3; i++ {
		n, err := strconv.ParseUint(h[i*2:i*2+2], 16, 8)
		if err != nil {
			t.Fatalf("%q is not a #rrggbb colour: %v", hex, err)
		}
		c[i] = channel(float64(n) / 255)
	}
	return 0.2126*c[0] + 0.7152*c[1] + 0.0722*c[2]
}

func contrast(t *testing.T, a, b string) float64 {
	t.Helper()
	la, lb := luminance(t, a), luminance(t, b)
	hi, lo := math.Max(la, lb), math.Min(la, lb)
	return (hi + 0.05) / (lo + 0.05)
}

// textRoles are the roles that render words. faint is excluded on purpose:
// ui.go says it is only ever a rule, a pending dot or a progress track, and
// it is deliberately below the floor.
func textRoles() map[string]struct{ light, dark string } {
	return map[string]struct{ light, dark string }{
		"brand":   {cBrand.Light.TrueColor, cBrand.Dark.TrueColor},
		"muted":   {cMuted.Light.TrueColor, cMuted.Dark.TrueColor},
		"success": {cSuccess.Light.TrueColor, cSuccess.Dark.TrueColor},
		"warning": {cWarning.Light.TrueColor, cWarning.Dark.TrueColor},
		"error":   {cError.Light.TrueColor, cError.Dark.TrueColor},
	}
}

func TestEveryTextRoleIsReadableOnEveryTerminalBackground(t *testing.T) {
	for role, pair := range textRoles() {
		for name, bg := range lightBackgrounds {
			if got := contrast(t, pair.light, bg); got < minRatio {
				t.Errorf("%s (%s) on a %s background is %.2f:1, below %.2f",
					role, pair.light, name, got, minRatio)
			}
		}
		for name, bg := range darkBackgrounds {
			if got := contrast(t, pair.dark, bg); got < minRatio {
				t.Errorf("%s (%s) on a %s background is %.2f:1, below %.2f",
					role, pair.dark, name, got, minRatio)
			}
		}
	}
}

// TestBrandMatchesTheWebsite is the comparator that stops the two identities
// drifting apart again. The values are the design system's --accent, and if
// the website's accent changes this is where the CLI finds out.
func TestBrandMatchesTheWebsite(t *testing.T) {
	const (
		webAccentLight = "#5B47E0"
		webAccentDark  = "#A99BFF"
	)
	if !strings.EqualFold(cBrand.Light.TrueColor, webAccentLight) {
		t.Errorf("the CLI brand on a light terminal is %s; the website's --accent is %s",
			cBrand.Light.TrueColor, webAccentLight)
	}
	if !strings.EqualFold(cBrand.Dark.TrueColor, webAccentDark) {
		t.Errorf("the CLI brand on a dark terminal is %s; the website's --accent is %s",
			cBrand.Dark.TrueColor, webAccentDark)
	}
}

// TestFaintIsBelowTheFloorOnPurpose documents the one exclusion, so that a
// future reader does not "fix" faint into a readable colour and quietly turn
// every rule and progress track into something that competes with the words.
func TestFaintIsBelowTheFloorOnPurpose(t *testing.T) {
	worst := math.Inf(1)
	for _, bg := range lightBackgrounds {
		worst = math.Min(worst, contrast(t, cFaint.Light.TrueColor, bg))
	}
	for _, bg := range darkBackgrounds {
		worst = math.Min(worst, contrast(t, cFaint.Dark.TrueColor, bg))
	}
	if worst >= minRatio {
		t.Errorf("faint now clears %.2f:1, which makes it a text colour. "+
			"It is meant to be a rule, a pending dot or a progress track, and "+
			"nothing that has to be read.", worst)
	}
}
