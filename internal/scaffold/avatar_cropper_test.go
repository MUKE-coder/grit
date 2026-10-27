package scaffold

import (
	"strings"
	"testing"
)

// The cropper is the one generated component whose whole job is direct
// manipulation, so it is the one that keeps being given a keyboard path that
// does not work.
//
// It shipped with tabIndex={0} on a div whose role said it was not
// interactive, which Biome refuses, so a freshly scaffolded project failed its
// own `pnpm lint` on a file the developer had never opened (grit#116). The
// suggested fix, deleting the tab stop, would have left the cropper unusable
// without a pointer.
//
// The canary runs the generated project's own lint nightly, which is the real
// check. This is the same property asserted where it is cheap, so a regression
// is caught at `go test` rather than a day later.
func TestTheCropperKeyboardPathIsRealControls(t *testing.T) {
	src := adminAvatarCropperTSX()

	// A tab stop on an element whose role is not interactive is the thing
	// Biome refuses, and deleting the tab stop is not the fix.
	if strings.Contains(src, "tabIndex={0}") {
		t.Error("a tabIndex={0} is back. If it is on an element with an " +
			"interactive role that is fine, and this test should say so; if it " +
			"is on a div, a fresh scaffold fails its own pnpm lint")
	}

	// What replaced it: native range inputs, which announce their name, value
	// and bounds without being told to, and move by step on an arrow key.
	if n := strings.Count(src, `type="range"`); n != 1 {
		t.Errorf("found %d range inputs, want the one shared control", n)
	}
	for _, label := range []string{`label="Size"`, `label="Horizontal"`, `label="Vertical"`} {
		if !strings.Contains(src, label) {
			t.Errorf("the %s control is missing: the crop cannot be reached "+
				"from the keyboard without it", label)
		}
	}

	// The step is the nudge, so an arrow key moves what it always moved.
	if !strings.Contains(src, "step={NUDGE}") {
		t.Error("the range step is not NUDGE, so arrow keys no longer move by " +
			"the amount the rest of this file is written in")
	}

	// The resize handles are a pointer affordance. They carried an aria-label
	// and a tab stop while having no key handler at all, so the keyboard tour
	// of this dialog was four buttons that announced a job they could not do.
	if !strings.Contains(src, "tabIndex={-1}") || !strings.Contains(src, `aria-hidden="true"`) {
		t.Error("the corner handles are focusable again, which puts four " +
			"buttons that do nothing back in the tab order")
	}
	if strings.Contains(src, "aria-label={`Resize from the") {
		t.Error("the handles announce themselves again while still having no " +
			"keyboard handler")
	}
}

// The dialog tells people how to use it, and that sentence has to describe the
// thing that is actually there.
func TestTheCropperSaysHowToUseIt(t *testing.T) {
	src := adminAvatarCropperTSX()
	if strings.Contains(src, "drag a corner to resize it") {
		t.Error("the instructions still offer only dragging, after the " +
			"keyboard path became the sliders")
	}
	if !strings.Contains(src, "use the sliders below") {
		t.Error("the instructions do not mention the controls that are there")
	}
}
