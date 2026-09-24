package scaffold

import (
	"regexp"
	"strings"
	"testing"
)

// What an audit of the generated admin found, written down so it stays found.
//
// Every one of these was measured in a browser against a real generated
// project, not reasoned about: a screen reader read twenty identical "Select"
// checkboxes, the search box was a placeholder with no name, the muted text
// sat at 3.1:1 against the dark canvas, and the primary button was white on
// #60a5fa at 2.5:1, which is a blue rectangle with a rumour of text on it.

// A label that is not attached to its input is decoration. The admin shipped
// twelve of thirteen form fields that way once; these are the ones an audit
// found still detached.
func TestAdminControlsAreLabelled(t *testing.T) {
	table := adminDataTable()

	if !strings.Contains(table, `aria-label={"Select " + rowName(row, columns)}`) {
		t.Error("a row checkbox does not say which row it selects, so a screen reader reads a page of identical checkboxes")
	}
	if !strings.Contains(table, "aria-label={allSelected ?") {
		t.Error("the select-all checkbox has no name")
	}
	if !strings.Contains(adminTableToolbar(), `aria-label={t("table.searchLabel"`) {
		t.Error("the search box is labelled by its placeholder, which is announced inconsistently and vanishes when you type")
	}
	if !strings.Contains(table, "function rowName(") {
		t.Error("rowName is gone, so the checkboxes have nothing to name themselves after")
	}

	options := AdminOptionLibraryTSX()
	for _, id := range []string{`htmlFor={"add-value-" + option.id}`, `id={"add-value-" + option.id}`} {
		if !strings.Contains(options, id) {
			t.Errorf("the option library's value field is not associated with its label (%s)", id)
		}
	}
}

// 2.4.1: the sidebar is eight links deep on every page.
func TestAdminHasASkipLink(t *testing.T) {
	layout := adminLayoutComponent()
	if !strings.Contains(layout, `href="#main"`) || !strings.Contains(layout, "Skip to content") {
		t.Error("no skip link: reaching the content by keyboard means tabbing the sidebar on every page")
	}
	if !strings.Contains(layout, "focus:not-sr-only") {
		t.Error("the skip link never becomes visible, so a sighted keyboard user cannot see where they are")
	}
}

// 1.3.1 and 2.4.1: the auth pages are a hero panel and a form, and the form is
// the part somebody came for.
func TestAuthPagesHaveAMainLandmark(t *testing.T) {
	shells := adminAtlasAuthShell()
	if !strings.Contains(shells, "<main") || !strings.Contains(shells, "</main>") {
		t.Error("the auth shell has no main landmark, so there is no way to skip the hero panel")
	}
}

// 1.4.3: measured, not guessed. Each of these is the ratio the browser
// reported against the background the token actually sits on.
func TestAdminTextMeetsAAContrast(t *testing.T) {
	css := adminGlobalCSS() + adminDarkModeCSSAddon()

	if strings.Contains(css, "--text-muted: #606078") || strings.Contains(css, "--text-muted: #94a3b8") {
		t.Error("the muted token is back to the value that measured 3.1:1 on dark and 2.5:1 on light")
	}

	// Every accent needs a label colour that reads on it.
	accents := regexp.MustCompile(`--accent: (#[0-9a-fA-F]{6});`).FindAllStringSubmatch(css, -1)
	if len(accents) < 4 {
		t.Fatalf("found %d accent definitions, expected the themes to still be there", len(accents))
	}
	if !strings.Contains(css, "--accent-fg:") {
		t.Fatal("no --accent-fg token: the button label is hard-coded white again")
	}
	for _, theme := range []string{"#60a5fa", "#3b82f6"} {
		block := css[strings.Index(css, "--accent: "+theme):]
		cut := strings.Index(block, "}")
		if cut > 0 {
			block = block[:cut]
		}
		if !strings.Contains(block, "--accent-fg: #0a0a0f") {
			t.Errorf("a white label on %s measures under 4.5:1, so that theme needs a dark --accent-fg", theme)
		}
	}
}

// 2.5.8: a 16px checkbox and a 20px eye toggle are under the minimum target.
func TestAdminTargetsAreBigEnough(t *testing.T) {
	table := adminDataTable()
	if !strings.Contains(table, `<label className="inline-flex cursor-pointer p-1">`) {
		t.Error("the table checkboxes lost the padded label that makes them a 24px target")
	}
	if strings.Contains(table, `className="text-xs text-text-secondary hover:text-accent transition-colors"`) {
		t.Error("the row actions are back to a 16px-tall target")
	}
}
