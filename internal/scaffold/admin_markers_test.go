package scaffold

import (
	"strings"
	"testing"
)

// Every admin extension point a plugin can inject at has a marker.
//
// grit#228, found reviewing the first community plugin submission:
// grit-calendar anchored 10 of its 18 injections on arbitrary lines of Grit's
// generated source, because the markers for what it extends did not exist. One
// anchor was a version comment ("v3.31.44 -- per-resource widgets"). Another
// was a whole seeded catalogue row. A third was a JSX condition.
//
// An injection whose marker is absent is not refused: it does not happen, and
// the install still reports success. So a Grit release that reformatted any of
// those lines would have left a plugin installing "cleanly" with its features
// missing, and nothing anywhere connecting the two.
func TestEveryAdminExtensionPointHasAMarker(t *testing.T) {
	next := Options{Architecture: ArchTriple, Frontend: FrontendNext}

	for _, c := range []struct {
		marker string
		file   string
		body   string
		why    string
	}{
		{
			"// grit:nav:imports",
			"components/chrome/CollapsibleSidebar.tsx",
			adminCollapsibleSidebarComponent(next),
			"a plugin adding a sidebar entry needs somewhere to import its config from",
		},
		{
			"// grit:dashboard:catalog",
			"lib/dashboard-catalog.ts",
			adminDashboardCatalogTS(),
			"a plugin's widget has to appear in the dashboard settings list or nobody can turn it on",
		},
		{
			"{/* grit:resource:slot */}",
			"components/resource/resource-page.tsx",
			adminResourcePage(),
			"a plugin mounting something page-level on every resource screen has nowhere else to go",
		},
	} {
		if !strings.Contains(c.body, c.marker) {
			t.Errorf("%s has no %s marker.\n  %s\n  Without it a plugin anchors on "+
				"a line of generated source, and the day that line moves the "+
				"injection silently stops happening.", c.file, c.marker, c.why)
		}
	}
}

// All four dashboard styles carry both dashboard markers.
//
// This is the one that matters, and the reason it is a loop. The default
// dashboard is the only style with a ResourceWidgetsRow, and grit-calendar
// anchored its widget on that import and on the comment above it. So on
// --style modern, minimal and glass the plugin installed successfully and its
// dashboard widget simply never appeared.
//
// A test over one style would have passed while three were broken, which is
// exactly how the bug got there.
func TestEveryDashboardStyleTakesAPluginWidget(t *testing.T) {
	for _, style := range []string{"", "centered", "modern", "minimal", "glass"} {
		body := adminDashboardPageForStyle(style)
		name := style
		if name == "" {
			name = "(default)"
		}
		for _, marker := range []string{
			"// grit:dashboard:imports",
			"{/* grit:dashboard:widgets */}",
		} {
			if !strings.Contains(body, marker) {
				t.Errorf("the %s dashboard has no %s marker, so a plugin's widget "+
					"is silently absent on a project scaffolded with that style",
					name, marker)
			}
		}
	}
}

// The widget slot is inside the page's JSX, not after it.
//
// A marker in the wrong place is worse than no marker: the injection happens,
// the file stops compiling, and the message points at the plugin rather than
// at this.
func TestTheWidgetSlotIsInsideTheJSX(t *testing.T) {
	for _, style := range []string{"", "modern", "minimal", "glass"} {
		body := adminDashboardPageForStyle(style)
		at := strings.Index(body, "{/* grit:dashboard:widgets */}")
		if at < 0 {
			continue // the test above reports this
		}
		rest := body[at:]
		// What follows is the closing of the page element and the return, in
		// that order. Anything else means the slot is outside the tree.
		closing := strings.Index(rest, "</div>")
		returning := strings.Index(rest, ");")
		if closing < 0 || returning < 0 || closing > returning {
			name := style
			if name == "" {
				name = "(default)"
			}
			t.Errorf("the %s dashboard's widget slot is not inside the page element: "+
				"an injection there would not compile", name)
		}
	}
}

// The markers survive the Next to TanStack transform.
//
// The Vite admin's components are the Next ones put through nextToTanStack, so
// a marker that the transform rewrites or drops exists in one admin and not the
// other, and a plugin supporting both would work on one shape only.
func TestTheMarkersSurviveTheTanStackTransform(t *testing.T) {
	for _, c := range []struct {
		marker string
		body   string
	}{
		{"// grit:nav:imports", nextToTanStack(adminCollapsibleSidebarComponent(Options{Architecture: ArchTriple, Frontend: FrontendTanStack}))},
		{"// grit:dashboard:catalog", nextToTanStack(adminDashboardCatalogTS())},
		{"{/* grit:resource:slot */}", nextToTanStack(adminResourcePage())},
		{"// grit:dashboard:imports", nextToTanStack(adminDashboardPageForStyle(""))},
		{"{/* grit:dashboard:widgets */}", nextToTanStack(adminDashboardPageForStyle(""))},
	} {
		if !strings.Contains(c.body, c.marker) {
			t.Errorf("%s does not survive nextToTanStack, so it exists in the Next "+
				"admin and not the Vite one", c.marker)
		}
	}
}
