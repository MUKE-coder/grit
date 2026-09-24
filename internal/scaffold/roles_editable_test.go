package scaffold

import (
	"strings"
	"testing"
)

// The roles that ship with a project are editable, ADMIN included.
//
// ADMIN holds the "*" grant. The editor read that off the role and, when it was
// set, replaced the whole permission grid with a paragraph telling the operator
// to "remove that grant" through a control that did not exist: the first role on
// the page could not be changed at all, and save always sent "*" straight back.
func TestTheWildcardGrantIsAControlNotAVerdict(t *testing.T) {
	page := adminRolesPage()

	if strings.Contains(page, `const isSuper = role ? role.grants.indexOf("*") >= 0 : false;`) {
		t.Fatal("the wildcard is still read-only state, so a role holding it cannot be edited")
	}
	if !strings.Contains(page, "const [superGrant, setSuperGrant] = useState(") {
		t.Fatal("no editable state for the wildcard")
	}
	// Seeded from the role, so opening ADMIN still shows it as holding everything.
	if !strings.Contains(page, `role ? role.grants.indexOf("*") >= 0 : false`) {
		t.Error("the toggle is not seeded from the role, so ADMIN would open as if it held nothing")
	}
	// There is a control, and it drives the same flag the save reads.
	if !strings.Contains(page, "onChange={setSuperGrant}") {
		t.Error("nothing can turn the wildcard off")
	}
	if !strings.Contains(page, `const grants = superGrant ? ["*"] : collapseGrants(selected, modules);`) {
		t.Error("save does not read the toggle, so unticking it would not be saved")
	}
	// And the paragraph points at that control rather than at nothing.
	if strings.Contains(page, "Remove that grant to") {
		t.Error("the explanation still asks for an action the page does not offer")
	}
	if !strings.Contains(page, "Untick the") {
		t.Error("the explanation does not say how to pick individual permissions")
	}
}

// A padlock on a card that opens says the wrong thing. The badge stays, because
// built-in is a fact worth knowing; the lock goes, because it is not true.
func TestBuiltInRolesDoNotRenderAsLocked(t *testing.T) {
	page := adminRolesPage()

	if strings.Contains(page, "<Lock ") {
		t.Error("built-in role cards still show a padlock, which reads as 'cannot be opened'")
	}
	if strings.Contains(page, `Lock,`) {
		t.Error("the padlock is still imported after being removed from the markup")
	}
	if !strings.Contains(page, "Built-in") {
		t.Error("the built-in badge was removed; it is a fact worth showing")
	}
	// The notice in the editor says what is fixed and what is not.
	if !strings.Contains(page, "Its description and\n					its permissions are yours to change") {
		t.Error("the built-in notice does not say the permissions are editable")
	}
	if !strings.Contains(page, "resolve it by name") {
		t.Error("the notice does not say why the name is fixed")
	}
}

// The standing rule against em dashes applies to the text a generated project
// shows its users, which is where it is easiest to miss.
func TestTheRolesScreensHaveNoEmDashesInProse(t *testing.T) {
	for name, src := range map[string]string{
		"admin roles page": adminRolesPage(),
		"desktop roles":    desktopClientSystemRolesPage(),
	} {
		for _, line := range strings.Split(src, "\n") {
			trimmed := strings.TrimSpace(line)
			// Code comments are not what the rule is about. What a person reads
			// on the screen is.
			if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "{/*") {
				continue
			}
			// A lone dash standing in for an empty cell is a glyph, not prose.
			if strings.Contains(line, "—") {
				t.Errorf("%s: em dash in %q", name, trimmed)
			}
			if strings.Contains(line, "&mdash;") && trimmed != "&mdash;" {
				t.Errorf("%s: em dash entity in prose: %q", name, trimmed)
			}
		}
	}
}
