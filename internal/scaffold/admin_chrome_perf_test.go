package scaffold

import (
	"strings"
	"testing"
)

// Contact-app review L20 to L23. Every file here is framework-owned and written
// whole by grit upgrade, so the templates are the whole fix; these pin them in
// each shape the panel ships in.

// adminChromeShapes renders one admin template the way each panel shape gets it.
func adminChromeShapes(src string) map[string]string {
	return map[string]string{
		"next admin":   src,
		"panel in web": embedAdminContent(src, adminRoutePrefixes),
		"vite admin":   nextToTanStack(src),
	}
}

func TestRefreshLeavesChromeQueriesAlone(t *testing.T) {
	for shape, src := range adminChromeShapes(adminPageHeaderComponent()) {
		if strings.Contains(src, "queryClient.invalidateQueries();") {
			t.Errorf("%s: the refresh button still invalidates every query", shape)
		}
		for _, want := range []string{
			`const CHROME_QUERY_ROOTS = ["me", "my-permissions", "notifications"];`,
			"predicate: (query) => !CHROME_QUERY_ROOTS.includes(String(query.queryKey[0])),",
		} {
			if !strings.Contains(src, want) {
				t.Errorf("%s: PageHeader is missing %s", shape, want)
			}
		}
	}
}

func TestSessionWatchdogDoesNotRearmOnActivity(t *testing.T) {
	for shape, src := range adminChromeShapes(adminSessionWatchdogComponent()) {
		for _, gone := range []string{"clearTimeout", "setTimeout", "scheduleIdleCheck"} {
			if strings.Contains(src, gone) {
				t.Errorf("%s: the watchdog still uses %s", shape, gone)
			}
		}
		for _, want := range []string{
			"const IDLE_CHECK_EVERY_MS = 5 * 1000;",
			"if (!openRef.current) lastActivityRef.current = Date.now();",
			"}, IDLE_CHECK_EVERY_MS);",
		} {
			if !strings.Contains(src, want) {
				t.Errorf("%s: the watchdog is missing %s", shape, want)
			}
		}
		// One activity listener effect, registered once: it must not depend on
		// the modal state, or every open and close re-adds five listeners.
		if strings.Contains(src, "[onActivity") || strings.Contains(src, "[open, scheduleIdleCheck]") {
			t.Errorf("%s: the activity listeners are re-registered when the modal opens", shape)
		}
	}
}

func TestSidebarGroupExpansionIsDerived(t *testing.T) {
	src := adminCollapsibleSidebarComponent(Options{ProjectName: "app", Version: "1.0.0"})
	for shape, got := range adminChromeShapes(src) {
		for _, gone := range []string{"expandedGroups", "setExpandedGroups", `pathname.startsWith("/" + r.slug)`} {
			if strings.Contains(got, gone) {
				t.Errorf("%s: the sidebar still has %s", shape, gone)
			}
		}
		if strings.Count(got, "isGroupOpen(groupName)") != 2 {
			t.Errorf("%s: the group header and its items do not both read isGroupOpen", shape)
		}
	}
	// The active group is found from the same path the links use, which the
	// panel inside the web app prefixes with /admin.
	embedded := embedAdminContent(src, adminRoutePrefixes)
	if !strings.Contains(embedded, `pathname.startsWith("/admin/resources/" + r.slug)`) {
		t.Error("panel in web: the active group is not matched against /admin/resources/")
	}
	if strings.Count(embedded, `"/admin/resources/" + r.slug`) < 3 {
		t.Error("panel in web: the group match and the links disagree on the resource path")
	}
}

func TestThumbnailsLoadLazilyWithASize(t *testing.T) {
	cells := adminCellRenderers()
	if n := strings.Count(cells, `loading="lazy"`); n != 4 {
		t.Errorf("cell renderers: %d lazy images, want 4 (user, image, file, files)", n)
	}
	if n := strings.Count(cells, "width={32}"); n != 4 {
		t.Errorf("cell renderers: %d sized images, want 4", n)
	}
	if !strings.Contains(adminUserCellComponent(), `loading="lazy"`) {
		t.Error("UserCell avatar is not lazy")
	}
	// Account avatars are always on screen: sized, but never lazy.
	for name, src := range map[string]string{
		"UserMenu": adminUserMenuComponent(),
		"sidebar":  adminCollapsibleSidebarComponent(Options{ProjectName: "app"}),
	} {
		if !strings.Contains(src, "width={36} height={36}") || strings.Contains(src, `alt={fullName} width={36} height={36} loading="lazy"`) {
			t.Errorf("%s: the account avatar is not sized, or is lazy", name)
		}
	}
	// A plain <img>, not next/image: covers can be any URL an editor pasted,
	// and next/image throws on a host it was not configured for.
	for name, src := range map[string]string{"blog list": webBlogListPage(), "home": webLandingPage(Options{ProjectName: "demo", Architecture: ArchDouble, Frontend: FrontendNext})} {
		if strings.Contains(src, "next/image") {
			t.Errorf("%s: covers moved to next/image", name)
		}
		if !strings.Contains(src, `loading="lazy"`) || !strings.Contains(src, `decoding="async"`) || !strings.Contains(src, "width={640}") {
			t.Errorf("%s: cover images are not lazy, async and sized", name)
		}
	}
	post := webBlogDetailPage()
	if strings.Contains(post, `loading="lazy"`) || !strings.Contains(post, `fetchPriority="high"`) || !strings.Contains(post, "width={1200}") {
		t.Error("blog post: the hero cover must stay eager and high priority, with a size")
	}
}
