package scaffold

import (
	"regexp"
	"strings"
	"testing"
)

// excelize is pinned in two scaffolds and floored in a third place, and the
// three drifted: the API sat on v2.11.0 while the desktop app sat on v2.8.1,
// which is old enough to carry five advisories including the four that kept
// the Scan workflow red from 20 October. A desktop project was getting a
// vulnerable Excel export from the same framework that had already fixed the
// API one.
//
// Both scaffolds use it for the same thing, so they pin the same version and
// that version is the floor grit upgrade raises an existing project to.
func TestExcelizePinsAgree(t *testing.T) {
	pin := regexp.MustCompile(`github\.com/xuri/excelize/v2 (v[^\s]+)`)

	find := func(name, src string) string {
		m := pin.FindStringSubmatch(src)
		if m == nil {
			t.Fatalf("%s does not pin excelize at all", name)
		}
		return m[1]
	}

	api := find("the API go.mod", apiGoMod(Options{ProjectName: "app"}))
	desktop := find("the desktop go.mod", desktopGoMod(DesktopOptions{ProjectName: "app"}))

	if api != desktop {
		t.Errorf("excelize pins disagree: API %s, desktop %s.\n"+
			"Both scaffolds import it for an Excel export, so a fix that reaches "+
			"one and not the other leaves half of the projects vulnerable.", api, desktop)
	}

	var floor string
	for _, d := range frameworkDeps {
		if d.Path == "github.com/xuri/excelize/v2" {
			floor = d.Floor
		}
	}
	if floor == "" {
		t.Fatal("excelize has no frameworkDeps floor, so grit upgrade leaves " +
			"existing projects on whatever they were scaffolded with")
	}
	if floor != api {
		t.Errorf("the frameworkDeps floor is %s and new projects get %s: "+
			"an upgraded project and a new one would not match", floor, api)
	}

	// The advisories this pin exists for were fixed in a commit, not a tag.
	// When v2.11.1 is released this stops being true and the comment beside
	// the pin says to move to it; the assertion below is a reminder, not a
	// rule, so it only fires if somebody moves backwards to a plain v2.11.0.
	if strings.HasPrefix(api, "v2.11.0") || strings.HasPrefix(api, "v2.10") || strings.HasPrefix(api, "v2.8") {
		t.Errorf("excelize is pinned to %s, which is affected by GO-2026-6684, "+
			"6685, 6688 and 6689, all reachable from the spreadsheet import", api)
	}
}
