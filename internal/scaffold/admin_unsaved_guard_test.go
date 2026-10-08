package scaffold

import (
	"path/filepath"
	"strings"
	"testing"
)

// A form with edits in it warns before they are thrown away.
//
// The admin had no beforeunload anywhere, so a resource form with twenty
// fields and a rich text body was lost to a stray click on the sidebar with
// no prompt and nothing on screen saying there was anything to lose.
func TestAFormWithEditsWarnsBeforeLosingThem(t *testing.T) {
	guard := adminUnsavedGuard()

	// Both halves, because they catch different accidents. beforeunload does
	// not fire on a Next <Link>, which in an admin panel is the common way a
	// half-filled form gets abandoned.
	for _, needed := range []string{"beforeunload", `addEventListener("click"`} {
		if !strings.Contains(guard, needed) {
			t.Errorf("the guard does not listen for %s", needed)
		}
	}

	// Capture phase, or the router has already navigated by the time this runs.
	if !strings.Contains(guard, `document.addEventListener("click", onClick, true)`) {
		t.Error("the click listener is not on the capture phase, so the router wins the race")
	}

	// Both listeners come off when the form goes clean, or a saved form keeps
	// asking and people learn to click through the prompt.
	if !strings.Contains(guard, "removeEventListener") {
		t.Error("the guard never removes its listeners")
	}
	if !strings.Contains(guard, "if (!dirty) return;") {
		t.Error("the guard arms itself on a form nobody has typed in")
	}

	// A modified click opens a new tab and leaves this one alone, so stopping
	// it would be stopping a navigation that is not happening.
	for _, key := range []string{"metaKey", "ctrlKey", "shiftKey", "altKey"} {
		if !strings.Contains(guard, key) {
			t.Errorf("the guard does not let a %s click through", key)
		}
	}

	form := adminFormBuilder()
	if !strings.Contains(form, "useUnsavedGuard(isDirty && !isSubmitting)") {
		t.Error("the form builder does not arm the guard")
	}
	if !strings.Contains(form, "<UnsavedBadge") {
		t.Error("the form builder never says on screen that there is something to lose")
	}
}

// The guard reaches every admin, in both dialects.
func TestTheGuardIsDeliveredToBothAdmins(t *testing.T) {
	next := adminFileMap(t.TempDir(), tripleOptions())
	found := false
	for path, body := range next {
		if strings.HasSuffix(filepath.ToSlash(path), "components/forms/unsaved-guard.tsx") {
			found = true
			if strings.Contains(body, `"use client"`) == false {
				t.Error("the Next guard is not a client component, so it has no listeners")
			}
		}
	}
	if !found {
		t.Error("the Next.js admin does not get components/forms/unsaved-guard.tsx")
	}

	// The Vite admin cannot carry "use client" at the top of a module.
	vite := nextToTanStack(adminUnsavedGuard())
	if strings.Contains(vite, `"use client"`) {
		t.Error("the Vite guard still carries \"use client\"")
	}
	if !strings.Contains(vite, "export function useUnsavedGuard") {
		t.Error("the Vite guard does not export the hook")
	}
}

// A counted card narrows the table to the rows it counted.
//
// Every stat card was a div carrying hover:border-border/80, so the whole row
// looked clickable and none of it was. "Lead 18" is exactly the number
// somebody wants to click.
func TestACountedCardNarrowsTheTable(t *testing.T) {
	cards := adminStatCards()

	if !strings.Contains(cards, "filter?: { field: string; value: string };") {
		t.Error("a stat card cannot say which rows it counted")
	}
	// The branch is the point: a card with a subset is a button, one without
	// stays a div and loses the hover state with it.
	if !strings.Contains(cards, "if (!stat.filter || !onFilter) {") {
		t.Error("the card does not choose between a button and a div")
	}
	if !strings.Contains(cards, `aria-pressed={active}`) {
		t.Error("a filtering card does not tell a screen reader whether it is on")
	}

	// The generic cards must not become buttons: there is nothing to narrow a
	// total to, and a hover state on one is the fault being fixed.
	page := adminResourcePage()
	if !strings.Contains(page, "onStatFilter=") {
		t.Error("the resource page does not wire the cards to its filters")
	}
	if !strings.Contains(page, "isStatFilterActive=") {
		t.Error("the resource page cannot say which card is active")
	}

	// Clicking the active card again clears it. A filter you cannot take off
	// from where you put it on is one people reload the page to escape.
	if !strings.Contains(page, `c.filters[filter.field] === filter.value ? "" : filter.value`) {
		t.Error("clicking the active card does not clear the filter")
	}

	controller := adminUseResourceController()
	if !strings.Contains(controller, "filter: { field: column.field, value }") {
		t.Error("the counted cards do not carry the filter they stand for")
	}
}
