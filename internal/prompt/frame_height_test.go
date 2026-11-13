package prompt

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
	"github.com/MUKE-coder/grit/v3/internal/ui"
)

const (
	testCols = 100
	testRows = 40
)

// frames returns the height of the first frame the form draws and the height
// of the frame after the terminal's size arrives.
//
// This is the real order of events. huh queues the size message as an Init
// command, so the first frame is already on screen when it lands; `presized`
// says whether the form was told the size beforehand, which is what
// ui.SizeForm does for real.
func frames(t *testing.T, presized bool) (first, after int) {
	t.Helper()
	opts := &scaffold.Options{}
	form := newProjectForm(opts, new(string), new(string), new(string), new(string), new(string)).
		WithTheme(ui.PromptTheme())

	size := tea.WindowSizeMsg{Width: testCols, Height: testRows}
	if presized {
		form.Update(size)
	}
	form.Init()
	first = lipgloss.Height(form.View())
	form.Update(size)
	return first, lipgloss.Height(form.View())
}

// The architecture list appeared twice on Git Bash, the stale copy a line
// short and missing Mobile.
//
// The cause is a frame that changes height. huh sizes every group to the
// tallest of them, and it learns the terminal's size from a message Init
// queues, so the first frame is the architecture picker's own height and the
// next is the theme picker's, three lines taller. Repainting at a new height
// means moving the cursor up and erasing, and on Git Bash the erase does not
// land: the short frame stays on screen under the tall one.
//
// There is nothing to get wrong if the size never changes.
func TestTheFirstFrameIsTheSizeEveryFrameAfterItIs(t *testing.T) {
	first, after := frames(t, true)
	if first != after {
		t.Errorf("the form is drawn at %d lines and then repainted at %d.\n"+
			"A repaint at a new height is what leaves a stale copy of the "+
			"picker on screen in Git Bash.", first, after)
	}
}

// The control: without the pre-sizing the heights differ, so the test above
// is testing the fix and not the weather.
func TestWithoutPreSizingTheFrameChangesHeight(t *testing.T) {
	first, after := frames(t, false)
	if first == after {
		t.Skipf("the pickers are all the same height now (%d lines), so the "+
			"reflow this guards against cannot arise; the guard stays for "+
			"when one of them grows again", first)
	}
}

// ui.SizeForm has to run before Init, because Run calls Init. If Init threw
// the size away the fix would do nothing, silently.
func TestInitKeepsTheSizeItWasGivenFirst(t *testing.T) {
	opts := &scaffold.Options{}
	form := newProjectForm(opts, new(string), new(string), new(string), new(string), new(string)).
		WithTheme(ui.PromptTheme())

	var sized *huh.Form = form
	sized.Update(tea.WindowSizeMsg{Width: testCols, Height: testRows})
	before := lipgloss.Height(sized.View())
	sized.Init()
	if got := lipgloss.Height(sized.View()); got != before {
		t.Errorf("Init changed the form from %d lines to %d, so sizing it "+
			"beforehand buys nothing", before, got)
	}
}
