package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// The components from section 6 of the style guide.
//
// Each returns a string rather than printing one, so a caller can put it in a
// table, a test can read it, and nothing writes to stdout from inside a helper.

// Wordmark is shown by three commands and no others: grit new, grit version,
// and grit on its own. Everywhere else opens with Header.
//
// The grain is the identity: loose data on the left, a solid app on the right.
func Wordmark(version string) string {
	if !unicode {
		return "\n" + Gutter + Brand.Render("grit") + "\n\n" +
			Gutter + "Describe your data. Get the whole app.  " + Muted.Render("v"+version) + "\n\n"
	}
	rows := []string{
		"████████  ███████   ██  ████████",
		"██        ██    ██  ██     ██",
		"██  ████  ███████   ██     ██",
		"██    ██  ██  ██    ██     ██",
		"████████  ██    ██  ██     ██",
	}
	var b strings.Builder
	b.WriteString("\n")
	for _, row := range rows {
		b.WriteString(Gutter + Brand.Render("░▒▓"+row) + "\n")
	}
	b.WriteString("\n" + Gutter + "Describe your data. Get the whole app.  " +
		Muted.Render("v"+version) + "\n\n")
	return b.String()
}

// Header is the one-line opening every other command uses:
//
//	░▒▓█ grit migrate  v3.387.0 · sqlite
//
// meta is whatever is worth knowing before the output starts: the version, the
// database, the architecture. It may be empty.
func Header(command, meta string) string {
	line := Gutter + Brand.Render(GlyphMark) + " " + Bold.Render("grit "+command)
	if meta != "" {
		line += "  " + Muted.Render(meta)
	}
	return line
}

// Step is a finished step: a check and a noun.
//
// "Go API", not "Scaffolding Go API...". The list reports results; only the
// step that is running animates.
func Step(label, detail string) string {
	line := Gutter + Success.Render(GlyphOK) + " " + label
	if detail != "" {
		line += "  " + Muted.Render(detail)
	}
	return line
}

// StepTimed is Step with the duration on the right, aligned to column.
func StepTimed(label string, d time.Duration, column int) string {
	line := Gutter + Success.Render(GlyphOK) + " " + label
	took := Duration(d)
	if took == "" {
		return line
	}
	return line + pad(label, column) + Muted.Render(took)
}

// Duration formats a timing: 8.4s, 120ms, and nothing at all below 100ms,
// where it is noise.
func Duration(d time.Duration) string {
	switch {
	case d < 100*time.Millisecond:
		return ""
	case d < time.Second:
		return fmt.Sprintf("%dms", d.Milliseconds())
	default:
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
}

// Addition is a thing that was created: a table, a file, a route.
func Addition(name string) string {
	return Gutter + Success.Render(GlyphAdded) + " " + name
}

// PendingStep is a step that has not started yet.
func PendingStep(label string) string {
	return Gutter + Faint.Render(GlyphPending) + " " + Muted.Render(label)
}

// Running is the step that is in flight, given a spinner frame.
func Running(frame, label string) string {
	return Gutter + Brand.Render(frame) + " " + label
}

// Result is the last line a command prints: what happened, then the detail.
//
//	Created my-app  Single · Atlas · 14.2s
func Result(headline, name, detail string) string {
	line := Gutter + Success.Bold(true).Render(headline)
	if name != "" {
		line += " " + Bold.Render(name)
	}
	if detail != "" {
		line += "  " + Muted.Render(detail)
	}
	return line
}

// WarningLine is a warning's first line: the glyph carries the colour, the
// title is bold in the terminal's own foreground so it still reads under
// NO_COLOR.
func WarningLine(title, detail string) string {
	line := Gutter + Warning.Render(GlyphWarn) + " " + Bold.Render(title)
	if detail != "" {
		line += "  " + Muted.Render(detail)
	}
	return line
}

// WarningBlock says what is wrong, what still works, and how to fix it, in
// that order, because the middle line is what stops somebody abandoning a run
// that was going to be fine.
func WarningBlock(title, detail string, body ...string) string {
	var b strings.Builder
	b.WriteString(WarningLine(title, detail))
	for _, line := range body {
		b.WriteString("\n" + Gutter + "  " + line)
	}
	return b.String()
}

// FailureBlock is a warning that stopped the command: same shape, a Try line
// instead of a Fix, and the caller exits non-zero.
func FailureBlock(title, cause, try string) string {
	var b strings.Builder
	b.WriteString(Gutter + Error.Render(GlyphFail) + " " + Bold.Render(title))
	if cause != "" {
		b.WriteString("\n" + Gutter + "  " + cause)
	}
	if try != "" {
		b.WriteString("\n" + Gutter + "  " + Muted.Render("Try") + "  " + Brand.Render(try))
	}
	return b.String()
}

// Fix is the line under a warning that says what to run.
func Fix(label, command string) string {
	return Muted.Render(label) + "  " + Brand.Render(command)
}

// Note is information that is not a problem: a baseline migration, a
// convention worth knowing. A brand bar down the left and no glyph, because
// nothing is wrong.
func Note(title string, lines ...string) string {
	var b strings.Builder
	b.WriteString(Gutter + Brand.Render(GlyphNoteBar) + " " + Bold.Render(title))
	for _, line := range lines {
		b.WriteString("\n" + Gutter + Brand.Render(GlyphNoteBar) + " " + line)
	}
	return b.String()
}

// Command is something to type, with its muted comment aligned at column.
func Command(cmd, comment string, column int) string {
	line := Gutter + "  " + Brand.Render(cmd)
	if comment == "" {
		return line
	}
	return line + pad(cmd, column) + Muted.Render(comment)
}

// URL is a row of the table a command ends with: a muted label, then the
// address underlined so terminals make it clickable.
//
// The underline is written directly rather than through lipgloss. Its v1.1.0
// Underline emits the pair around every single character, so a twenty-character
// address became four hundred bytes of escape codes: fine on screen, unreadable
// in a log and absurd in a pipe.
func URL(label, href string, column int) string {
	return Gutter + Muted.Render(label) + pad(label, column) + underline(href)
}

// underline wraps a string in the SGR pair, or returns it untouched when the
// output is not styled at all.
func underline(s string) string {
	if lipgloss.ColorProfile() == termenv.Ascii {
		return s
	}
	return "[4m" + s + "[24m"
}

// Rule is the single faint line allowed above the URL table, and nowhere else.
func Rule(width int) string {
	if width <= 0 {
		width = 52
	}
	ch := "─"
	if !unicode {
		ch = "-"
	}
	return Gutter + Faint.Render(strings.Repeat(ch, width))
}

// Section is a bold label introducing a group, like "Next steps".
func Section(title string) string {
	return Gutter + Bold.Render(title)
}

// More is how a list longer than eight rows ends.
func More(n int, how string) string {
	return Gutter + Muted.Render(fmt.Sprintf("  … %d more, see them all with %s", n, how))
}

// LogLine is a runtime log: the time, a fixed-width source, then the message.
// No date, no banner, no rule made of equals signs.
func LogLine(at time.Time, source, message, detail string) string {
	line := Gutter + Muted.Render(at.Format("15:04:05")+"  "+padRight(source, 8)) + message
	if detail != "" {
		line += "  " + Muted.Render(detail)
	}
	return line
}

// Answered is a prompt that has been answered, collapsed to one line so the
// question being asked is always the last thing on screen.
func Answered(label, value, detail string) string {
	line := Gutter + Success.Render(GlyphOK) + " " + Muted.Render(label) + "  " + value
	if detail != "" {
		line += " " + Muted.Render("· "+detail)
	}
	return line
}

// pad returns the spaces between a left column of width column and whatever
// follows it, never fewer than two.
func pad(s string, column int) string {
	n := column - len([]rune(s))
	if n < 2 {
		n = 2
	}
	return strings.Repeat(" ", n)
}

func padRight(s string, n int) string {
	if len([]rune(s)) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len([]rune(s)))
}
