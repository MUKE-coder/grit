package scaffold

import (
	"strings"
	"testing"
)

// freshSidebar is the sidebar a new project gets.
func freshSidebar() string {
	return adminCollapsibleSidebarComponent(Options{})
}

// TestContentNavAnchorsExistInTheTemplate is the guard that matters.
//
// repairContentNav works by finding four strings in an existing sidebar and
// editing around them. If one of those strings drifts out of the template,
// the repair does not fail loudly: it prints a warning and returns, so every
// project that upgrades keeps the old sidebar and plugins keep filing their
// boards under System. Nothing else would notice.
//
// So the anchors are checked against the template itself, here, where a drift
// fails CI instead of quietly disabling an upgrade path.
func TestContentNavAnchorsExistInTheTemplate(t *testing.T) {
	src := freshSidebar()

	anchors := map[string]string{
		"SYSTEM_NAV array":       systemNavClose,
		"resource grouping loop": contentNavGroupsLoop,
		"Content heading":        contentNavHeadingNew,
		"resource link list":     contentNavRootMapEnd,
	}
	for name, anchor := range anchors {
		if n := strings.Count(src, anchor); n != 1 {
			t.Errorf("the %s anchor appears %d times in the scaffolded sidebar, want exactly 1.\n"+
				"repairContentNav edits around it, so a drift here silently disables the upgrade.\n"+
				"anchor:\n%s", name, n, anchor)
		}
	}

	// And the template itself must carry the marker a plugin injects into.
	if n := strings.Count(src, contentNavMarker); n != 1 {
		t.Errorf("a new project's sidebar has %d copies of %s, want exactly 1: "+
			"without it `grit plugin add` cannot wire a Content nav link",
			n, contentNavMarker)
	}
	if n := strings.Count(src, "// grit:nav:system"); n != 1 {
		t.Errorf("a new project's sidebar has %d copies of the system nav marker, want exactly 1", n)
	}
}

// TestContentNavRepairUpgradesAnOldSidebar builds the sidebar as it was before
// CONTENT_NAV existed and checks the repair produces the current one.
func TestContentNavRepairUpgradesAnOldSidebar(t *testing.T) {
	current := freshSidebar()
	old := oldSidebarWithoutContentNav(t, current)

	out, fixed, warnings := repairContentNavSource(old)
	if len(warnings) > 0 {
		t.Fatalf("the repair warned instead of applying: %v", warnings)
	}
	if len(fixed) == 0 {
		t.Fatal("the repair reported no change")
	}
	if !strings.Contains(out, contentNavMarker) {
		t.Fatal("the repaired sidebar has no grit:nav:content marker, so a plugin still cannot inject into it")
	}
	if !strings.Contains(out, "const visibleContentNav = CONTENT_NAV.filter(") {
		t.Fatal("the repaired sidebar does not compute visibleContentNav")
	}
	if !strings.Contains(out, "{visibleContentNav.map((entry) => {") {
		t.Fatal("the repaired sidebar never renders the plugin entries")
	}
	if !strings.Contains(out, contentNavHeadingNew) {
		t.Fatal("the repaired sidebar still hides the Content heading when only a plugin entry is present")
	}

	// The whole point: the repaired file matches what a new project gets.
	if out != current {
		t.Errorf("a repaired sidebar differs from a freshly scaffolded one.\n"+
			"repaired length %d, fresh length %d", len(out), len(current))
	}
}

// TestContentNavRepairIsIdempotent: upgrade runs on every release, and a
// repair that applies twice would add a second CONTENT_NAV.
func TestContentNavRepairIsIdempotent(t *testing.T) {
	current := freshSidebar()
	out, fixed, warnings := repairContentNavSource(current)
	if out != current {
		t.Error("the repair changed a sidebar that already has CONTENT_NAV")
	}
	if len(fixed) > 0 || len(warnings) > 0 {
		t.Errorf("the repair reported work on an up-to-date sidebar: fixed=%v warnings=%v", fixed, warnings)
	}
}

// TestContentNavRepairLeavesAnEditedSidebarAlone: somebody who has rewritten
// their sidebar gets advice, not a broken file.
func TestContentNavRepairLeavesAnEditedSidebarAlone(t *testing.T) {
	theirs := "export const SYSTEM_NAV = [];\n// a sidebar somebody rewrote\n"
	out, fixed, warnings := repairContentNavSource(theirs)
	if out != theirs {
		t.Error("the repair edited a sidebar it did not recognise")
	}
	if len(fixed) > 0 {
		t.Errorf("the repair claimed to fix an unrecognised sidebar: %v", fixed)
	}
	if len(warnings) != 1 || !strings.Contains(warnings[0], "CONTENT_NAV") {
		t.Errorf("the warning should say what to copy across; got %v", warnings)
	}
}

// oldSidebarWithoutContentNav reverses the four edits, producing the file as
// it stood before this change. Deriving it from the current template rather
// than pasting a copy means this test cannot rot into testing a sidebar the
// scaffold no longer writes.
func oldSidebarWithoutContentNav(t *testing.T, current string) string {
	t.Helper()

	out := current
	for _, cut := range []string{contentNavDeclaration, contentNavFilter, contentNavRender} {
		if !strings.Contains(out, cut) {
			t.Fatalf("the current sidebar does not contain a block this repair inserts, "+
				"so the before-and-after cannot be built:\n%s", cut)
		}
		out = strings.Replace(out, cut, "", 1)
	}
	if !strings.Contains(out, contentNavHeadingNew) {
		t.Fatal("the current sidebar does not contain the new heading condition")
	}
	out = strings.Replace(out, contentNavHeadingNew, contentNavHeadingOld, 1)

	if strings.Contains(out, contentNavMarker) {
		t.Fatal("the reconstructed old sidebar still has the content marker")
	}
	return out
}
