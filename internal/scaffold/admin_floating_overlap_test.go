package scaffold

import (
	"strings"
	"testing"
)

// The admin used to have two elements that floated over everything, and both
// sat on top of content until v3.530.0.
//
// Measured in a browser at 375px rather than reasoned about:
//
//	the menu toggle   fixed top-3 left-3, 36px square, so it covered the first
//	                  48px of the viewport. The page header was sticky top-0,
//	                  which put the two in exactly the same place at every
//	                  scroll position, and the toggle had the higher z-index:
//	                  the title sat underneath it the moment anybody scrolled.
//	the quick-access  fixed bottom-6, 48px tall, reserving the last 72px.
//	button            Nothing reserved it, so the final element on the page
//	                  sat beneath it: the "Last" pagination button at 1440,
//	                  the search field at 375.
//
// v3.532.0 removed both rather than clearing room for them. The toggle is a
// normal button in the TopBar's own row, and the quick-access button is gone:
// it offered one action from on top of the last row of whatever was
// underneath it, and Cmd+K offers every action and covers nothing.
//
// So this test now guards the absence. A fixed element over the content is
// easy to reach for and the overlap it causes is only visible at one width,
// on one scroll position, which is how both of these shipped.
func TestNothingFloatsOverTheAdminsContent(t *testing.T) {
	layout := adminLayoutComponent()

	if strings.Contains(layout, "<QuickAccess") {
		t.Error("the floating quick-access button is back in the layout. It covers the " +
			"last row of the page, offers one action, and the command palette " +
			"already offers every action from no pixels at all")
	}
	if strings.Contains(layout, "fixed top-3 left-3") {
		t.Error("the menu toggle is fixed over the page again, which puts it on top of " +
			"the page title at every scroll position on a phone")
	}
	if !strings.Contains(layout, "<TopBar") {
		t.Error("the layout does not mount the TopBar, so there is nothing holding the " +
			"menu button, the breadcrumbs or the status pill")
	}

	// The thing that does not work, recorded so nobody tries it again:
	// padding on the scrolling wrapper cannot move a sticky header, because
	// a sticky element pins to the viewport rather than to its parent's box.
	if strings.Contains(layout, "pt-16") {
		t.Error("<main> has mobile top padding, which was the first attempt at the " +
			"toggle overlap and does nothing once the page is scrolled: a sticky " +
			"header pins to the viewport top regardless")
	}
}

// The page header carries the page, and the TopBar carries the panel.
//
// They used to be one thing. Every page's header held the theme toggle, the
// notification bell and the user menu on the same line as that page's title
// and its own primary action, so the one control the page existed for was the
// least prominent thing in the row, and three panel-wide controls were
// re-rendered on every screen.
func TestPanelChromeIsNotInThePageHeader(t *testing.T) {
	header := adminPageHeaderComponent()

	for _, chrome := range []string{"<DarkModeToggle", "<NotificationBell", "<UserMenu"} {
		if strings.Contains(header, chrome) {
			t.Errorf("PageHeader renders %s, which belongs to the panel rather than to "+
				"any page and already sits in the TopBar: two of them is two places "+
				"to press the same thing", chrome)
		}
	}

	if strings.Contains(header, "sticky top-0") {
		t.Error("PageHeader is sticky again. The TopBar is the sticky one; two sticky " +
			"headers stack and the second covers the first row of the content")
	}

	bar := adminTopBarComponent()
	for _, chrome := range []string{"<DarkModeToggle", "<NotificationBell"} {
		if !strings.Contains(bar, chrome) {
			t.Errorf("the TopBar does not render %s, so it is now in neither place", chrome)
		}
	}
}
