package scaffold

import (
	"strings"
	"testing"
)

// The customer area's layout carries the two markers the generator writes at.
//
// This is a test on the template rather than on the generator, because the
// failure it guards is an edit here. internal/generate injects a sidebar entry
// before // grit:nav and an icon import before // grit:nav-icons; take either
// out and every generated customer screen is still written, still reachable by
// typing the URL, and invisible. Nothing else fails.
//
// Both, not one. An entry naming an icon the file does not import is a layout
// that will not compile, which takes the whole signed-in area down rather than
// one page.
func TestTheCustomerLayoutKeepsItsInjectionMarkers(t *testing.T) {
	layout := webAccountLayout(Options{ProjectName: "keep"})

	for _, marker := range []string{"// grit:nav", "// grit:nav-icons"} {
		if !strings.Contains(layout, marker) {
			t.Errorf("app/(app)/layout.tsx has no %s marker, so grit generate resource "+
				"--portal screens writes a screen nothing links to", marker)
		}
	}

	// Inside the array, not after it. An entry injected below the closing
	// bracket is a syntax error in the file every signed-in page renders
	// through.
	nav := strings.Index(layout, "const NAV = [")
	end := strings.Index(layout[nav:], "];")
	if nav < 0 || end < 0 {
		t.Fatal("the NAV array is not in the layout any more, so the marker has nothing to be inside")
	}
	if !strings.Contains(layout[nav:nav+end], "// grit:nav") {
		t.Error("the // grit:nav marker is outside the NAV array, so an injected entry lands after it")
	}
}
