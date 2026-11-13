package ui

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/x/term"
)

// RunForm runs a huh form with Grit's prompt theme, after telling it how big
// the terminal is.
//
// huh learns the terminal size from a message, and Init queues that message as
// a command, which means it arrives after the first frame is already on
// screen. A form of more than one group sizes every group to the tallest of
// them, so the first frame is the current group's own height and every frame
// after it is taller. `grit new` showed that as a duplicated architecture
// list: nine lines went up, twelve came back, and the repaint at the new
// height left the shorter frame behind, stale and missing the options that
// had not fitted in it. The erase relies on cursor-up codes landing exactly,
// which is a thing Git Bash does not reliably do.
//
// Measuring the terminal here and handing huh the same message it was going to
// send itself means the first frame is already the final size, so there is no
// repaint at a different height to get wrong. The numbers are huh's own: this
// changes when the form is sized, not what it decides.
//
// Theme first, because the theme's padding and borders are part of what the
// sizing measures.
func RunForm(form *huh.Form) error {
	form = form.WithTheme(PromptTheme())
	SizeForm(form)
	return form.Run()
}

// SizeForm gives a form the terminal's size before it runs. Exported for the
// test that proves the first frame and the frames after it agree.
//
// huh renders to stderr, so that is the handle whose size decides. When it is
// not a terminal, as under `go test` or a pipe, there is nothing to measure
// and nothing to fix: huh falls back to its own defaults as before.
func SizeForm(form *huh.Form) {
	width, height, err := term.GetSize(os.Stderr.Fd())
	if err != nil || width <= 0 || height <= 0 {
		return
	}
	form.Update(tea.WindowSizeMsg{Width: width, Height: height})
}
