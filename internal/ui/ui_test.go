package ui

import (
	"strings"
	"testing"
	"time"
)

// Body text carries no colour of its own: the terminal's foreground is the one
// value readable on every theme anybody runs.
func TestBodyTextHasNoColour(t *testing.T) {
	if got := Text.Render("hello"); got != "hello" {
		t.Errorf("Text styled body copy: %q", got)
	}
}

// Every status has a glyph, so the output still reads with NO_COLOR set or
// piped to a file, where the colour is gone and the glyph is not.
func TestEveryStatusHasAGlyph(t *testing.T) {
	for name, line := range map[string]string{
		"step":    Step("Go API", ""),
		"added":   Addition("users"),
		"pending": PendingStep("Admin panel"),
		"warning": WarningLine("MinIO is not answering", ""),
		"failure": FailureBlock("Port 8080 is already in use", "", ""),
	} {
		plain := strip(line)
		if !strings.ContainsAny(plain, "✓+·!✗okx.") {
			t.Errorf("%s has no glyph: %q", name, plain)
		}
	}
}

// Timings: seconds with one decimal, milliseconds as whole numbers, and
// nothing at all under 100ms where it is noise.
func TestDurationFormatting(t *testing.T) {
	for d, want := range map[time.Duration]string{
		99 * time.Millisecond:   "",
		120 * time.Millisecond:  "120ms",
		8400 * time.Millisecond: "8.4s",
	} {
		if got := Duration(d); got != want {
			t.Errorf("%v: got %q, want %q", d, got, want)
		}
	}
}

// Every line starts at column 3. Nothing begins at column 1 except the shell's
// own prompt.
func TestEveryLineHasTheGutter(t *testing.T) {
	lines := []string{
		Header("migrate", "v1 · sqlite"),
		Step("Go API", ""),
		Addition("users"),
		PendingStep("Admin panel"),
		Result("Created", "my-app", "Single"),
		WarningLine("title", ""),
		Note("Baseline run", "detail"),
		Command("grit migrate", "create tables", 20),
		URL("API", "http://localhost:8080", 14),
		Section("Next steps"),
	}
	for _, line := range lines {
		for _, l := range strings.Split(line, "\n") {
			if !strings.HasPrefix(strip(l), "  ") {
				t.Errorf("no gutter: %q", strip(l))
			}
		}
	}
}

// A second column lines up: the comment beside a command, the address beside a
// label. Never fewer than two spaces, however long the left side gets.
func TestColumnsAlign(t *testing.T) {
	// The column is the longest command plus a gap, which is how callers
	// compute it. Anything at or past the column cannot align by definition,
	// and gets the two-space minimum instead.
	const column = 24
	short := strip(Command("cd app", "comment", column))
	long := strip(Command("docker compose up -d", "comment", column))
	if strings.Index(short, "comment") != strings.Index(long, "comment") {
		t.Errorf("comments do not line up:\n%q\n%q", short, long)
	}
	// A command longer than the column still gets a gap rather than running
	// into its comment.
	over := strip(Command("a-very-long-command-indeed", "comment", 10))
	if !strings.Contains(over, "d  comment") {
		t.Errorf("an over-long command lost its gap: %q", over)
	}
}

// strip removes ANSI escapes so a test reads what a person sees.
func strip(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == 0x1b {
			for i < len(s) && s[i] != 'm' {
				i++
			}
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}
