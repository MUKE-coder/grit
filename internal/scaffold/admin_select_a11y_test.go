package scaffold

import (
	"strings"
	"testing"
)

// The single-choice pickers, held to the standard the multi-choice one already
// met (grit#179, proposal 5).
//
// The proposal asked for a searchable combobox for every belongs_to field and
// for an enum over about seven options. Most of that had shipped: both
// relationship pickers and SelectField all have a search box, and the
// relationship ones search on the server. What had not shipped was everything
// around the search box.
//
// SelectField rendered its list as an absolutely positioned child of the form.
// Inside a dialog or any scroll container with overflow hidden it clipped, and
// near the bottom of a long form it opened past the edge of the window, which
// is the same as having no list. The date picker and both relationship pickers
// are portalled and flip. This one was the odd one out.
//
// And neither single picker was a listbox. aria-haspopup on the trigger and
// nothing else: no role="listbox", no role="option", no aria-selected, no
// aria-activedescendant. A screen reader got a pile of unlabelled buttons.
// That is the failure mode that goes unnoticed for months, because it is
// invisible to everybody who is in the room when the feature is demonstrated.
func TestTheSingleChoicePickersArePortalledListboxes(t *testing.T) {
	pickers := map[string]string{
		"select-field":              adminSelectField(),
		"relationship-select-field": adminRelationshipSelectField(),
	}

	for name, src := range pickers {
		// Portalled and fixed, so the panel is not clipped by an ancestor.
		if !strings.Contains(src, "createPortal") {
			t.Errorf("%s does not portal its panel, so it clips inside a dialog", name)
		}
		if !strings.Contains(src, `className="fixed z-[9999]`) &&
			!strings.Contains(src, `"fixed z-[9999]`) {
			t.Errorf("%s does not position its panel fixed above the page", name)
		}
		// Opens upward when there is more room above. Without this a picker on
		// the last field of a form opens off screen.
		if !strings.Contains(src, "window.innerHeight") {
			t.Errorf("%s never measures the space it has, so it cannot flip", name)
		}

		// Real listbox semantics.
		for _, want := range []string{
			`role="listbox"`,
			`role="option"`,
			"aria-selected",
			"aria-activedescendant",
		} {
			if !strings.Contains(src, want) {
				t.Errorf("%s is missing %s, so a screen reader announces unlabelled buttons", name, want)
			}
		}

		// Keyboard movement through the options, and focus handed back on
		// close. Focus left on an unmounted node drops the user at the top of
		// the document.
		for _, want := range []string{`"ArrowDown"`, `"ArrowUp"`, `"Home"`, `"End"`, `"Escape"`} {
			if !strings.Contains(src, want) {
				t.Errorf("%s does not handle %s", name, want)
			}
		}
		if !strings.Contains(src, "triggerRef.current?.focus()") {
			t.Errorf("%s does not return focus to its trigger when it closes", name)
		}
	}
}

// The threshold the proposal names. A search box over three choices is
// clutter; over forty it is the only way to use the control.
func TestTheEnumPickerSearchesOnlyWhenThereIsSomethingToSearch(t *testing.T) {
	src := adminSelectField()
	if !strings.Contains(src, "const SEARCH_FROM = 8;") {
		t.Error("select-field does not declare a threshold for showing the search box")
	}
	if !strings.Contains(src, "options.length >= SEARCH_FROM") {
		t.Error("the threshold is declared and not applied")
	}
	// And the box is actually conditional on it, rather than the constant
	// being computed and ignored.
	if !strings.Contains(src, "{searchable && (") {
		t.Error("the search box is rendered unconditionally")
	}
	// With no search box, the arrow keys have to work from the trigger or the
	// control cannot be used without a mouse at all.
	if !strings.Contains(src, "onTriggerKeyDown") {
		t.Error("with no search box there is no keyboard way into the list")
	}
}

// The TanStack admin gets the same components through nextToTanStack, so a fix
// to one has to survive that translation rather than reaching one admin only.
func TestTheTanStackAdminGetsThePickerFixesToo(t *testing.T) {
	for name, src := range map[string]string{
		"select-field":              nextToTanStack(adminSelectField()),
		"relationship-select-field": nextToTanStack(adminRelationshipSelectField()),
	} {
		for _, want := range []string{"createPortal", `role="listbox"`, "aria-activedescendant"} {
			if !strings.Contains(src, want) {
				t.Errorf("the TanStack %s lost %s in translation", name, want)
			}
		}
	}
}
