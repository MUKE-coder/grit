package ui

import (
	"math"
	"os"
	"path/filepath"
	"regexp"
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

// TestTextRolesTrackTheWebsite is the comparator that stops the CLI and the
// website drifting apart, and it reads the website rather than a copy of it.
//
// The version before this held two hex strings transcribed from globals.css.
// It covered the brand and nothing else, and by the time anybody looked, five
// of the other six roles had drifted: close enough to look deliberate, far
// enough that the CLI and the site were two palettes.
//
// Where the CLI cannot take the site's value it says so here, with the
// measured ratio that forced it, so an exception has to be argued for rather
// than just existing.
func TestTextRolesTrackTheWebsite(t *testing.T) {
	css := websiteTokens(t)

	// role -> (css token, cli pair, and whether each mode must match)
	cases := []struct {
		role       string
		token      string
		pair       colourPair
		matchLight bool
		matchDark  bool
		why        string
	}{
		{"brand", "accent", colourPair{cBrand.Light.TrueColor, cBrand.Dark.TrueColor}, true, true, ""},
		{"muted", "fg-muted", colourPair{cMuted.Light.TrueColor, cMuted.Dark.TrueColor}, true, true, ""},
		{
			"success", "success",
			colourPair{cSuccess.Light.TrueColor, cSuccess.Dark.TrueColor},
			false, true,
			"the site's light --success is 4.66:1 on a light grey terminal, below the 5.1 floor",
		},
		{
			"warning", "warning",
			colourPair{cWarning.Light.TrueColor, cWarning.Dark.TrueColor},
			false, true,
			"the site's light --warning is 4.32:1 on a light grey terminal",
		},
		{
			"error", "danger",
			colourPair{cError.Light.TrueColor, cError.Dark.TrueColor},
			false, false,
			"the site's --danger misses in both modes: 4.45:1 on light grey, 5.04:1 on One Dark",
		},
	}

	for _, c := range cases {
		light, dark := css.light[c.token], css.dark[c.token]
		if light == "" || dark == "" {
			t.Errorf("the website has no --%s in one of its modes (light=%q dark=%q); "+
				"if it was renamed, this test is how the CLI finds out",
				c.token, light, dark)
			continue
		}

		if c.matchLight && !strings.EqualFold(c.pair.light, light) {
			t.Errorf("%s on a light terminal is %s; the website's --%s is %s",
				c.role, c.pair.light, c.token, light)
		}
		if c.matchDark && !strings.EqualFold(c.pair.dark, dark) {
			t.Errorf("%s on a dark terminal is %s; the website's --%s is %s",
				c.role, c.pair.dark, c.token, dark)
		}

		// An exception has to stay an exception. If the site moves its token
		// to something the CLI could now use, this says so rather than
		// letting the deviation outlive its reason.
		if !c.matchLight && strings.EqualFold(c.pair.light, light) {
			t.Errorf("%s light now equals the website's --%s (%s), so the exception "+
				"is stale: set matchLight and delete the note (%s)",
				c.role, c.token, light, c.why)
		}
		if !c.matchDark && strings.EqualFold(c.pair.dark, dark) {
			t.Errorf("%s dark now equals the website's --%s (%s), so the exception "+
				"is stale: set matchDark and delete the note (%s)",
				c.role, c.token, dark, c.why)
		}
	}
}

type colourPair struct{ light, dark string }

type tokenSet struct {
	light map[string]string
	dark  map[string]string
}

// websiteTokens reads docs/app/globals.css and returns the two palettes.
//
// The file declares the dark palette first, under :root, and the light one in
// a later block. Which is which is decided by --bg rather than by order or by
// a selector name, because that is the one token whose value says outright
// which mode a block is: a near-black background is the dark theme.
func websiteTokens(t *testing.T) tokenSet {
	t.Helper()

	path := filepath.Join("..", "..", "docs", "app", "globals.css")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("the website stylesheet is not in this checkout (%v), so there "+
			"is nothing to compare against", err)
	}

	blocks := splitCSSBlocks(string(body))
	out := tokenSet{light: map[string]string{}, dark: map[string]string{}}
	for _, block := range blocks {
		bg := block["bg"]
		if bg == "" {
			continue
		}
		target := out.light
		if isDarkHex(t, bg) {
			target = out.dark
		}
		for k, v := range block {
			// First declaration wins: the base palette is declared before the
			// per-theme overrides that follow it.
			if _, seen := target[k]; !seen {
				target[k] = v
			}
		}
	}

	if len(out.light) == 0 || len(out.dark) == 0 {
		t.Fatalf("could not read a light and a dark palette out of %s "+
			"(light=%d tokens, dark=%d)", path, len(out.light), len(out.dark))
	}
	return out
}

// isDarkHex reports whether a background colour is the dark theme's.
func isDarkHex(t *testing.T, hex string) bool {
	t.Helper()
	if !strings.HasPrefix(hex, "#") || len(hex) != 7 {
		return false
	}
	return luminance(t, hex) < 0.5
}

var (
	blockRe = regexp.MustCompile(`(?s)\{([^{}]*)\}`)
	varRe   = regexp.MustCompile(`--([a-z0-9-]+)\s*:\s*(#[0-9a-fA-F]{6})\s*;`)
)

// splitCSSBlocks returns one map of custom properties per brace-delimited
// block. Only hex values: a var() reference or an rgba() is not a colour this
// comparator can compare.
func splitCSSBlocks(css string) []map[string]string {
	var out []map[string]string
	for _, m := range blockRe.FindAllStringSubmatch(css, -1) {
		vars := map[string]string{}
		for _, v := range varRe.FindAllStringSubmatch(m[1], -1) {
			vars[v[1]] = v[2]
		}
		if len(vars) > 0 {
			out = append(out, vars)
		}
	}
	return out
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
