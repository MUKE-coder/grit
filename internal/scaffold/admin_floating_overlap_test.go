package scaffold

import (
	"strings"
	"testing"
)

// The admin has two elements that float over everything, and both sat on top
// of content until v3.530.0.
//
// Measured in a browser at 375px rather than reasoned about:
//
//	the menu toggle   fixed top-3 left-3, 36px square, so it covers the first
//	                  48px of the viewport. The page header is sticky top-0,
//	                  which puts the two in exactly the same place at every
//	                  scroll position, and the toggle has the higher z-index:
//	                  the title sat underneath it the moment anybody scrolled.
//	the quick-access  fixed bottom-6, 48px tall, reserving the last 72px.
//	button            Nothing reserved it, so the final element on the page
//	                  sat beneath it: the "Last" pagination button at 1440,
//	                  the search field at 375.
//
// Both fixes are one utility each, which is exactly why they need a test:
// a utility is the easiest thing in the world to tidy away.
func TestTheFloatingElementsLeaveRoomForContent(t *testing.T) {
	header := adminPageHeaderComponent()

	// The toggle is 36px at left-3, so it ends at 48px. pl-14 is 56px.
	if !strings.Contains(header, "pl-14") {
		t.Error("the page header has no mobile left padding, so the fixed menu toggle " +
			"covers the page title whenever the sticky header is pinned")
	}
	// And the desktop padding still wins at md, or every page gains an indent.
	if !strings.Contains(header, "md:px-8") {
		t.Error("the page header lost its desktop padding")
	}

	layout := adminLayoutComponent()
	if !strings.Contains(layout, "pb-20") {
		t.Error("<main> reserves no room at the bottom, so the last element on the " +
			"page sits under the floating quick-access button")
	}

	// The thing that does not work, recorded so nobody tries it again:
	// padding on the scrolling wrapper cannot move a sticky header, because
	// a sticky element pins to the viewport rather than to its parent's box.
	if strings.Contains(layout, "pt-16") {
		t.Error("<main> has mobile top padding, which was the first attempt at the " +
			"toggle overlap and does nothing once the page is scrolled: the header " +
			"is sticky, so it pins to the viewport top regardless")
	}
}
