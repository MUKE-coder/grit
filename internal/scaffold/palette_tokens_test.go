package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// A colour utility Tailwind does not know about compiles to nothing.
//
// Tailwind v4 builds one utility per --color-* in the @theme block, and the
// admin defines sixteen. Anything else, text-text-primary or bg-card or
// text-muted-foreground, produces no rule at all: no build error, no warning in
// the console, just an element keeping whatever it inherited. The page looks
// almost right, which is worse than looking broken, and it is why the account
// screen reads as "off design and inconsistent" rather than as a bug anybody
// could point at.
//
// Sixty-seven of these had accumulated across the admin. This walks every file
// the scaffolder writes and checks each colour utility against the tokens that
// app's own stylesheet defines, so the next one fails here instead of shipping.

var (
	// A colour utility, with any variant prefixes already consumed by the
	// preceding boundary: bg-, text-, border-, ring- and friends.
	colorUtility = regexp.MustCompile(`(?:^|[\s"'` + "`" + `:{(\[])(bg|text|border|ring|divide|placeholder|caret|fill|stroke|outline|decoration|from|via|to|shadow|accent)-([a-z][a-z0-9]*(?:-[a-z0-9]+)*)`)

	// The --color-* names in an @theme block.
	themeToken = regexp.MustCompile(`--color-([a-z0-9-]+)\s*:`)
)

// Values that are not palette lookups: Tailwind's own keywords, its numbered
// scales, and the non-colour words that share these prefixes (text-sm,
// border-b, bg-none, from-left).
var notAPaletteColor = map[string]bool{
	"transparent": true, "current": true, "inherit": true, "white": true, "black": true,
	"none": true, "auto": true, "clip": true, "ellipsis": true, "wrap": true, "nowrap": true,
	"balance": true, "pretty": true, "left": true, "center": true, "right": true,
	"justify": true, "start": true, "end": true, "top": true, "bottom": true,
	"solid": true, "dashed": true, "dotted": true, "double": true, "hidden": true,
	"collapse": true, "separate": true, "fixed": true, "local": true, "scroll": true,
	"cover": true, "contain": true, "repeat": true, "opacity": true, "origin": true,
	"box": true, "size": true, "position": true, "image": true, "gradient": true,
	"blend": true, "clip-text": true, "inner": true, "b": true, "t": true, "l": true,
	"r": true, "x": true, "y": true, "s": true, "e": true, "se": true, "sm": true,
	"base": true, "lg": true, "xl": true, "px": true, "line": true, "through": true,
	"underline": true, "overline": true, "offset": true, "dense": true, "reverse": true,
}

// Tailwind's built-in numbered scales (red-500, slate-800). Allowed by this
// test, though a themed app should rarely need one.
var numberedScale = regexp.MustCompile(`^[a-z]+-(50|[1-9]00|950)$`)

// unknownColorTokens returns the colour utilities in src that name something
// the theme does not define.
func unknownColorTokens(src string, known map[string]bool) []string {
	var bad []string
	seen := map[string]bool{}
	for _, m := range colorUtility.FindAllStringSubmatch(src, -1) {
		prefix, value := m[1], m[2]
		value, ok := colorPart(prefix, value)
		if !ok {
			continue
		}
		if known[value] || notAPaletteColor[value] || numberedScale.MatchString(value) {
			continue
		}
		// A bare word that is not one of the names an author reaches for is a
		// size or a keyword (text-sm, shadow-lg), not a missing colour.
		if !strings.Contains(value, "-") && !looksLikeAColorWord(value) {
			continue
		}
		key := prefix + "-" + value
		if !seen[key] {
			seen[key] = true
			bad = append(bad, key)
		}
	}
	sort.Strings(bad)
	return bad
}

// looksLikeAColorWord keeps the check from flagging text-sm and border-b while
// still catching bg-card and text-primary.
func looksLikeAColorWord(value string) bool {
	switch value {
	case "primary", "secondary", "muted", "card", "popover", "destructive",
		"accent", "background", "foreground", "border", "input", "ring",
		"success", "danger", "warning", "info":
		return true
	}
	return false
}

// sides are the edge selectors that sit between a prefix and its value:
// border-t-accent is a colour, border-t-2 is a width.
var sides = map[string]bool{
	"t": true, "r": true, "b": true, "l": true, "x": true, "y": true,
	"s": true, "e": true, "ss": true, "se": true, "ee": true, "es": true,
}

// colorPart strips the parts of a utility that are not the colour, and reports
// whether what is left is a palette lookup at all.
func colorPart(prefix, value string) (string, bool) {
	parts := strings.Split(value, "-")

	switch prefix {
	case "border", "divide":
		if len(parts) > 1 && sides[parts[0]] {
			parts = parts[1:]
		}
	case "outline", "ring", "shadow":
		// outline-offset-2 and ring-offset-2 are distances; ring-offset-<color>
		// is a colour, but of the ring offset rather than the ring.
		if len(parts) > 1 && parts[0] == "offset" {
			parts = parts[1:]
		}
	}
	if len(parts) == 0 {
		return "", false
	}
	// bg-gradient-to-r names a direction, not a colour; the colour is in the
	// from-/via-/to- utilities beside it.
	if parts[0] == "gradient" {
		return "", false
	}
	// A width, a radius, a size: border-2, ring-4, shadow-2xl.
	if isNumber(parts[0]) && len(parts) == 1 {
		return "", false
	}
	return strings.Join(parts, "-"), true
}

func isNumber(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// knownTokensFrom pulls the palette out of a generated stylesheet.
func knownTokensFrom(t *testing.T, css string) map[string]bool {
	t.Helper()
	known := map[string]bool{}
	for _, m := range themeToken.FindAllStringSubmatch(css, -1) {
		known[m[1]] = true
	}
	if len(known) < 8 {
		t.Fatalf("found only %d --color-* tokens; the stylesheet is not the one this test expects", len(known))
	}
	// bg-background is how --color-background is used, and border/ring/divide
	// all fall back to --color-border's name.
	return known
}

func TestAdminUsesOnlyColoursItDefines(t *testing.T) {
	for _, frontend := range []Frontend{FrontendNext, FrontendTanStack} {
		t.Run(string(frontend), func(t *testing.T) {
			files := adminTSX(t, frontend)

			var css string
			for path, src := range adminAllFiles(t, frontend) {
				if strings.HasSuffix(path, "globals.css") || strings.HasSuffix(path, "index.css") {
					css = src
					break
				}
			}
			if css == "" {
				t.Skip("no stylesheet emitted for this frontend")
			}
			known := knownTokensFrom(t, css)

			total := 0
			for path, src := range files {
				for _, bad := range unknownColorTokens(src, known) {
					t.Errorf("%s: %s names a colour the theme does not define, so it compiles to nothing", path, bad)
					total++
					if total > 12 {
						t.Fatalf("stopping after %d; fix these first", total)
					}
				}
			}
		})
	}
}

// adminAllFiles is adminTSX including the stylesheet, so the test reads the
// palette from the same tree it checks rather than a copy that can drift.
func adminAllFiles(t *testing.T, frontend Frontend) map[string]string {
	t.Helper()

	root := t.TempDir()
	opts := Options{ProjectName: "palette-test", Architecture: ArchTriple, Frontend: frontend}
	if err := createDirectories(root, opts); err != nil {
		t.Fatalf("createDirectories: %v", err)
	}
	write := writeAdminFiles
	if frontend == FrontendTanStack {
		write = writeAdminTanStackFiles
	}
	if err := write(root, opts); err != nil {
		t.Fatalf("write admin files: %v", err)
	}

	out := map[string]string{}
	adminDir := filepath.Join(root, "apps", "admin")
	err := filepath.Walk(adminDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".css") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(adminDir, path)
		out[filepath.ToSlash(rel)] = string(data)
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", adminDir, err)
	}
	return out
}
