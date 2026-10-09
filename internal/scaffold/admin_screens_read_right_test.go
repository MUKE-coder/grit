package scaffold

import (
	"strings"
	"testing"
)

// Four things a generated admin panel said or did that reading the screens
// caught and no unit test could.
//
// Every one of them was in a shared component, so every resource in every
// generated project inherited it. Found by generating a contacts app, standing
// it up, and looking at it.

// A boolean column renders every boolean there is, so it cannot borrow the
// words of one of them.
func TestABooleanColumnSaysYesOrNo(t *testing.T) {
	for name, src := range map[string]string{
		"admin":   adminCellRenderers(),
		"desktop": desktopClientDataTable(),
	} {
		// The rendered words, not the comment explaining why they changed.
		if strings.Contains(src, ">Active</span>") || strings.Contains(src, "/> Active<") ||
			strings.Contains(src, ">Inactive</span>") || strings.Contains(src, "/> Inactive<") {
			t.Errorf("%s: the boolean cell still says Active or Inactive. It renders every "+
				"boolean column, so a contact's `starred` column reads \"Inactive\" and an "+
				"order's `paid` column reads \"Active\".", name)
		}
		if !strings.Contains(src, "Yes") || !strings.Contains(src, "No") {
			t.Errorf("%s: the boolean cell says neither Yes nor No", name)
		}
	}
}

// The tabs are about archiving, and most resources have nothing to publish.
func TestTheArchiveTabsSayActiveNotPublished(t *testing.T) {
	src := adminResourcePage()
	if strings.Contains(src, `{ label: "Published", archived: false }`) {
		t.Error("the not-archived tab is still labelled Published, on every resource. " +
			"A Contact is not published and a Group is not published.")
	}
	if !strings.Contains(src, `{ label: "Active", archived: false }`) {
		t.Error("the not-archived tab is not labelled Active")
	}
	if !strings.Contains(src, `{ label: "Archived", archived: true }`) {
		t.Error("the Archived tab went missing")
	}
}

// Advice about filters is only useful to somebody who set one.
func TestTheEmptyStateDistinguishesFilteredFromEmpty(t *testing.T) {
	empty := adminTableEmptyState()
	if !strings.Contains(empty, "filtered") {
		t.Fatal("the empty state does not take a filtered prop, so it gives the same " +
			"advice to a resource with no rows as to a search that matched nothing")
	}
	for _, want := range []string{"Try adjusting your search or filters", "Nothing here yet"} {
		if !strings.Contains(empty, want) {
			t.Errorf("the empty state never says %q", want)
		}
	}

	table := adminDataTable()
	if !strings.Contains(table, "filtered={Boolean(isFiltered)}") {
		t.Error("the table does not pass isFiltered to the empty state")
	}

	page := adminResourcePage()
	if !strings.Contains(page, "isFiltered={") {
		t.Error("the resource page never tells the table whether anything is narrowing " +
			"the rows, so the table cannot know which empty state is true")
	}
}

// A portalled panel that opens past the bottom edge cannot be scrolled to.
// The date picker has flipped above since it was written; the relationship
// pickers did not, and the owner field on an --owned-by resource is the last
// field on the form.
func TestEveryPortalledPanelFlipsWhenThereIsNoRoomBelow(t *testing.T) {
	src := adminDateField() + adminRelationshipSelectField() + adminMultiRelationshipSelectField()

	// One per picker, plus the date picker that always had it.
	if got := strings.Count(src, "window.innerHeight - r"); got < 3 {
		t.Errorf("only %d panels measure the room below them; the single and multi "+
			"relationship pickers and the date picker all need to. A picker at the "+
			"bottom of a drawer opens its list off the window, and the panel is "+
			"position: fixed, so nothing can scroll it back.", got)
	}

	for _, gone := range []string{
		`setPos({ top: rect.bottom + 4, left: rect.left, width: rect.width });`,
		`setPosition({ top: r.bottom + window.scrollY + 4, left: r.left + window.scrollX, width: r.width });`,
	} {
		if strings.Contains(src, gone) {
			t.Errorf("a picker still places itself below unconditionally:\n  %s", gone)
		}
	}

	// And the list gives up its height to the panel rather than fixing it, so
	// the search box and the create row survive a tight fit.
	if strings.Contains(src, `className="max-h-60 overflow-y-auto p-1"`) {
		t.Error("a picker's list is still a fixed 240px tall, which overflows the panel " +
			"when the panel is clamped to the space available")
	}
}
