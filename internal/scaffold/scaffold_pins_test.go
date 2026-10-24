package scaffold

import (
	"regexp"
	"strings"
	"testing"
)

// requireLine matches one line of a go.mod require block.
var requireLine = regexp.MustCompile(`(?m)^\s+([a-zA-Z0-9][^\s/]*(?:/[^\s]+)*)\s+(v[^\s]+)`)

// pinsIn reads the module versions out of a go.mod template.
func pinsIn(src string) map[string]string {
	out := map[string]string{}
	for _, m := range requireLine.FindAllStringSubmatch(src, -1) {
		mod, version := m[1], m[2]
		// "// indirect" lines and comment text are not pins we own.
		if strings.HasPrefix(mod, "//") {
			continue
		}
		if _, seen := out[mod]; !seen {
			out[mod] = version
		}
	}
	return out
}

// The API and the desktop app are scaffolded from two separate go.mod
// templates that share most of their dependencies, and nothing compared them.
// They drifted, quietly, in the direction that matters:
//
//   excelize        API v2.11.0, desktop v2.8.1 — three years and five
//                   advisories apart, for the same Excel export
//   golang.org/x/crypto  API v0.57.0, desktop v0.31.0 — below the floor that
//                   exists because govulncheck flags anything older
//
// Neither was going to fail anything. The generated-app scan builds an API
// project and never builds a desktop one, so a desktop-only vulnerability is
// invisible to it. This compares the two templates directly instead.
func TestScaffoldPinsAgreeAcrossTemplates(t *testing.T) {
	api := pinsIn(apiGoMod(Options{ProjectName: "app"}))
	desktop := pinsIn(desktopGoMod(DesktopOptions{ProjectName: "app"}))

	if len(api) < 10 || len(desktop) < 5 {
		t.Fatalf("parsed too few pins to be reading the templates properly: "+
			"api=%d desktop=%d", len(api), len(desktop))
	}

	for mod, apiVersion := range api {
		desktopVersion, shared := desktop[mod]
		if !shared {
			continue
		}
		if apiVersion != desktopVersion {
			t.Errorf("%s: the API scaffolds %s and the desktop one scaffolds %s.\n"+
				"Both are framework-owned pins for the same library; a fix that "+
				"reaches one and not the other leaves half the projects on the bug.",
				mod, apiVersion, desktopVersion)
		}
	}
}

// Every library with a grit upgrade floor is scaffolded at exactly that floor.
//
// Below it, a new project is born needing an upgrade. Above it, an upgraded
// project and a new one differ, and the floor stops describing what anybody
// is actually running.
func TestScaffoldPinsMatchTheUpgradeFloors(t *testing.T) {
	templates := map[string]map[string]string{
		"the API go.mod":     pinsIn(apiGoMod(Options{ProjectName: "app"})),
		"the desktop go.mod": pinsIn(desktopGoMod(DesktopOptions{ProjectName: "app"})),
	}

	checked := 0
	for _, dep := range frameworkDeps {
		for name, pins := range templates {
			version, present := pins[dep.Path]
			if !present {
				continue
			}
			checked++
			if version != dep.Floor {
				t.Errorf("%s pins %s at %s, and the grit upgrade floor is %s.\n"+
					"Why the floor exists: %s",
					name, dep.Path, version, dep.Floor, dep.Why)
			}
		}
	}
	if checked == 0 {
		t.Fatal("no frameworkDeps module was found in either template, so this " +
			"test is asserting nothing")
	}
}

// excelize specifically, because the four advisories it was behind are
// reachable from the spreadsheet import in both the API and the desktop app,
// and the fix is a commit rather than a tag, which is easy to "tidy up" into a
// plain version that reintroduces them.
func TestExcelizeIsPinnedPastItsAdvisories(t *testing.T) {
	const mod = "github.com/xuri/excelize/v2"
	for name, src := range map[string]string{
		"the API go.mod":     apiGoMod(Options{ProjectName: "app"}),
		"the desktop go.mod": desktopGoMod(DesktopOptions{ProjectName: "app"}),
	} {
		version, ok := pinsIn(src)[mod]
		if !ok {
			t.Errorf("%s does not pin excelize at all", name)
			continue
		}
		for _, bad := range []string{"v2.11.0", "v2.10", "v2.9", "v2.8", "v2.11.1-0.202609"} {
			if strings.HasPrefix(version, bad) {
				t.Errorf("%s pins excelize %s, which is affected by GO-2026-6684, "+
					"6685, 6688 and 6689: panics and unbounded allocation on a "+
					"crafted workbook, all reachable from the import.", name, version)
			}
		}
	}
}
