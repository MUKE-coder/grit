// Package ui is how every grit command prints.
//
// One palette, one set of glyphs, one gutter. Commands call the helpers here
// rather than building their own styles, because thirty commands each reaching
// for color.New is how the CLI ended up with a bright magenta wordmark and no
// two steps formatted the same way.
//
// The rules it enforces, from cli-design/cli-design-style-guide.md:
//
//   - Body text has no colour. The terminal's own foreground is the only one
//     readable on every theme anybody runs.
//   - Every colour is a pair, a dark-terminal value and a light-terminal one.
//     Never a single hex.
//   - Colour never carries meaning alone. Each status has a glyph, so the
//     output still reads with NO_COLOR set or piped to a file.
//   - Brand colour means "you can act on this": a command to type, the
//     selected option, the spinner. Emphasis is bold, not colour.
package ui

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// pair builds one role: a true-colour value and an ANSI-16 fallback for each
// background. The ANSI values matter more than they look, because a terminal
// with no 24-bit support still has to render this legibly.
func pair(light, dark, ansiLight, ansiDark string) lipgloss.CompleteAdaptiveColor {
	return lipgloss.CompleteAdaptiveColor{
		Light: lipgloss.CompleteColor{TrueColor: light, ANSI256: light, ANSI: ansiLight},
		Dark:  lipgloss.CompleteColor{TrueColor: dark, ANSI256: dark, ANSI: ansiDark},
	}
}

// The seven roles. Every text role clears 5.1:1 on black, #0D1117, #1E1E1E,
// One Dark, Solarized Dark, white, #FAFAFA, #EEEEEE and Solarized Light.
//
// faint is deliberately below that, which is why it is only ever a rule, a
// pending dot or a progress track, and never a word.
var (
	// The brand violet, the same pair the website uses: --accent is #5B47E0
	// on a light background and #A99BFF on a dark one.
	//
	// It was cobalt, which meant somebody ran `grit new`, watched a blue CLI
	// scaffold their project, then opened a violet website about it. Two
	// identities with nothing comparing them is the exact failure DESIGN.md
	// exists to stop, and the CLI was the half nobody had looked at.
	//
	// TestBrandIsReadableOnEveryTerminalBackground keeps the hue honest: the
	// nine backgrounds below are the claim this file has always made in a
	// comment, and now makes in a test.
	cBrand = pair("#5B47E0", "#A99BFF", "5", "13")

	// The rest of the roles are the website's tokens too, wherever the
	// website's value clears this file's contrast floor.
	//
	// It usually does not, and that is not a disagreement: the site holds its
	// colours to two backgrounds, its own light and dark. The CLI holds them
	// to nine, because a terminal has no zoom, no user stylesheet and no
	// high-contrast mode to fall back on. Where a site value would fail, the
	// CLI keeps a darker or lighter one and the comment carries the measured
	// ratio that forced it, so the next person does not "fix" it back.
	//
	// TestTextRolesTrackTheWebsite is the comparator. It reads
	// docs/app/globals.css rather than a copy of it, so a token changed on
	// the website fails the CLI's build until somebody has looked.

	// --fg-muted exactly, both modes.
	cMuted = pair("#555B64", "#A1A7B0", "8", "8")

	// Not a text role, so it has no site equivalent and no floor to clear:
	// faint is a rule, a pending dot or a progress track, never a word.
	cFaint = pair("#C4C9D2", "#4A5160", "7", "8")

	// Dark is --success exactly. Light is one step darker: the site's
	// #0b7a45 measures 4.66:1 on a light grey terminal, below the 5.1 floor.
	cSuccess = pair("#0A7343", "#3DD68C", "2", "10")

	// Dark is --warning exactly. Light is darker: #a35f00 measures 4.32:1 on
	// light grey, which is the worst miss of the three.
	cWarning = pair("#855300", "#F5A524", "3", "11")

	// Neither mode can take the site's value. --danger #ff6b6b is 5.04:1 on
	// One Dark, and #c93636 is 4.45:1 on light grey. Both miss by little, and
	// little is still unreadable for the person it is unreadable for.
	cError = pair("#BE1F2B", "#FF7B7B", "1", "9")
)

// The styles. Text sets no Foreground on purpose: that is the rule, not an
// omission.
var (
	Text    = lipgloss.NewStyle()
	Bold    = lipgloss.NewStyle().Bold(true)
	Brand   = lipgloss.NewStyle().Foreground(cBrand)
	Muted   = lipgloss.NewStyle().Foreground(cMuted)
	Faint   = lipgloss.NewStyle().Foreground(cFaint)
	Success = lipgloss.NewStyle().Foreground(cSuccess)
	Warning = lipgloss.NewStyle().Foreground(cWarning)
	Error   = lipgloss.NewStyle().Foreground(cError)
	Link    = lipgloss.NewStyle().Underline(true)
)

// Gutter is the two spaces every line starts with. Nothing begins at column 1
// except the shell's own prompt.
const Gutter = "  "

func init() {
	// lipgloss asks the terminal for its background and picks the right half of
	// each pair. When it cannot ask, it assumes dark, which is right more often
	// than not but not always, so the answer is overridable.
	switch strings.ToLower(strings.TrimSpace(os.Getenv("GRIT_THEME"))) {
	case "light":
		lipgloss.SetHasDarkBackground(false)
	case "dark":
		lipgloss.SetHasDarkBackground(true)
	}
}

// The glyphs, each with the ASCII the same line falls back to.
//
// A terminal that has not told us it speaks UTF-8 gets the second column: a
// box-drawing character rendered as a question mark is worse than a plain one
// that was never going to be pretty.
var (
	GlyphOK      = glyph("✓", "ok")
	GlyphAdded   = glyph("+", "+")
	GlyphWarn    = glyph("!", "!")
	GlyphFail    = glyph("✗", "x")
	GlyphPending = glyph("·", ".")
	GlyphAsk     = glyph("?", "?")
	GlyphPointer = glyph("❯", ">")
	GlyphNoteBar = glyph("│", "|")
	GlyphMark    = glyph("░▒▓█", "grit")
)

var (
	spinnerUTF8  = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinnerASCII = []string{"-", "\\", "|", "/"}
)

// Frames is the spinner this terminal can draw. Brand coloured by the caller,
// and only ever on the step that is actually running.
func Frames() []string {
	if unicode {
		return spinnerUTF8
	}
	return spinnerASCII
}

// unicode reports whether the terminal has told us it speaks UTF-8. Read once:
// the environment does not change under a running command.
var unicode = hasUTF8()

func hasUTF8() bool {
	for _, key := range []string{"LC_ALL", "LC_CTYPE", "LANG"} {
		if v := os.Getenv(key); v != "" {
			up := strings.ToUpper(v)
			return strings.Contains(up, "UTF-8") || strings.Contains(up, "UTF8")
		}
	}
	// Windows Terminal, VS Code and every modern Windows console draw these
	// fine and set none of the variables above. So does the legacy console,
	// which does not, hence the TERM check rather than a flat true.
	return os.Getenv("WT_SESSION") != "" || os.Getenv("TERM_PROGRAM") != "" ||
		os.Getenv("TERM") != ""
}

func glyph(utf8, ascii string) string {
	if unicode {
		return utf8
	}
	return ascii
}
