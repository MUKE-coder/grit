package scaffold

import (
	"strings"
	"testing"
)

// The marketing navbar is a mega menu, declared once and rendered by every
// frontend Grit ships.
func TestMegaMenuIsDeclaredOnceAndRenderedEverywhere(t *testing.T) {
	triple := Options{ProjectName: "contacts", Architecture: ArchTriple, Frontend: FrontendNext, Theme: "emerald"}

	menu := webMegaMenu()
	// Base UI's NavigationMenu, not a hand-rolled dropdown: CLAUDE.md asks for
	// it when a primitive needs behaviour, and this one needs a portalled
	// popup, hover intent, Escape, and arrow-key movement.
	if !strings.Contains(menu, `from "@base-ui/react/navigation-menu"`) {
		t.Error("the mega menu is not built on Base UI")
	}
	for _, part := range []string{"NavigationMenu.Root", "NavigationMenu.Trigger", "NavigationMenu.Content", "NavigationMenu.Portal", "NavigationMenu.Viewport"} {
		if !strings.Contains(menu, part) {
			t.Errorf("the mega menu does not use %s", part)
		}
	}
	// Portalled on purpose: the header is sticky, so a panel inside it is
	// clipped by the header's own height.
	if !strings.Contains(menu, "NavigationMenu.Portal") {
		t.Error("the panel is not portalled out of the sticky header")
	}
	// An external link says so to a screen reader.
	if !strings.Contains(menu, "(opens in a new tab)") {
		t.Error("an external link does not announce that it leaves the site")
	}
	// A whole row is the target, not the 14px title alone.
	if !strings.Contains(menu, "flex items-start gap-3 rounded-lg p-3") {
		t.Error("a panel link's target is smaller than its row")
	}

	// The entries are data, in one file.
	cfg := webNavMenuConfig(triple)
	for _, want := range []string{"export const navEntries", "MegaMenuColumn", "MegaMenuFeature", "lib/nav-menu.ts"} {
		if !strings.Contains(cfg, want) {
			t.Errorf("the nav config has no %s", want)
		}
	}
	// And it names the project, not the framework.
	if !strings.Contains(cfg, "contacts") {
		t.Error("the nav config never names the project")
	}

	// A triple's panel is another origin; a double's and a single's is a route.
	if !strings.Contains(cfg, `process.env.NEXT_PUBLIC_ADMIN_URL`) {
		t.Error("a triple's admin link is not overridable")
	}
	if strings.Contains(cfg, "{{ADMIN_HREF}}") {
		t.Error("the placeholder survived: the Vite writers do not substitute it")
	}
	for _, inApp := range []Options{
		{ProjectName: "contacts", Architecture: ArchDouble, Frontend: FrontendNext},
		{ProjectName: "contacts", Architecture: ArchSingle, Frontend: FrontendTanStack},
	} {
		got := webNavMenuConfig(inApp)
		if !strings.Contains(got, `const ADMIN_URL = "/admin/dashboard"`) {
			t.Errorf("%s does not link the panel as a route of this app", inApp.Architecture)
		}
		if strings.Contains(got, "external: true") {
			t.Errorf("%s opens its own panel in a new tab", inApp.Architecture)
		}
	}

	// Every navbar renders it, and none of them lists the entries itself.
	for name, nav := range map[string]string{
		"web":      webNavbar(triple),
		"web auth": webNavbarWithAuth(triple),
		"single":   singleViteNavbar(Options{ProjectName: "contacts", Architecture: ArchSingle}),
	} {
		if !strings.Contains(nav, "<MegaMenu />") || !strings.Contains(nav, "<MegaMenuMobile") {
			t.Errorf("the %s navbar does not render the menu on both widths", name)
		}
		if strings.Contains(nav, "const navLinks = [") {
			t.Errorf("the %s navbar still carries its own copy of the entries", name)
		}
	}
}
