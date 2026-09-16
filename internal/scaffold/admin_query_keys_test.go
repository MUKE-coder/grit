package scaffold

import (
	"path/filepath"
	"strings"
	"testing"
)

// The admin asked for the same data under different keys: the dashboard's stat
// card and latest table fetched one URL twice per resource, the dashboard and
// the bell polled /api/notifications separately, the multi relationship picker
// cached its options under a key no save invalidated, and tab counts lived in
// component state where nothing could refresh them. Every read of a resource
// now keys off resourceKeys in use-resource.ts.

func TestResourceKeysCoverEveryReadOfAResource(t *testing.T) {
	src := adminUseResource()
	for _, want := range []string{
		`dashboardStats: (endpoint: string, params: object) => [endpoint, "dashboard-stats", params]`,
		`options: (endpoint: string, search = "") => [endpoint, "options", search]`,
		`tabCount: (endpoint: string, tab: string, filters: object) => [endpoint, "tab-count", tab, filters]`,
		`tree: (endpoint: string) => [endpoint, "tree"]`,
		"export function useResourceDashboardStats",
		"export function useRelationshipOptions",
		"queryKey: resourceKeys.dashboardStats(endpoint, params)",
		"queryKey: resourceKeys.options(endpoint, term)",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("use-resource.ts is missing %q", want)
		}
	}
}

func TestDashboardWidgetsShareOneQuery(t *testing.T) {
	for name, src := range map[string]string{
		"ResourceStatCard":    adminResourceStatCardTSX(),
		"ResourceLatestTable": adminResourceLatestTableTSX(),
	} {
		if !strings.Contains(src, "useResourceDashboardStats(resource.endpoint, dateRangeToQueryParams(dateRange), ") {
			t.Errorf("%s does not read the shared dashboard stats query", name)
		}
		// A key or URL of its own is how the two drifted apart.
		for _, banned := range []string{"queryKey:", "/api/admin/dashboard/resource-stats/", "limit: String(limit)"} {
			if strings.Contains(src, banned) {
				t.Errorf("%s still has %q", name, banned)
			}
		}
	}
}

func TestNotificationListIsOneQuery(t *testing.T) {
	for name, src := range map[string]string{
		"NotificationBell":   adminNotificationBellComponent(),
		"dashboard":          adminCaptivatingDashboard(),
		"notifications page": adminNotificationsPage(),
	} {
		if !strings.Contains(src, "useNotificationList(") {
			t.Errorf("%s does not read the shared notification list", name)
		}
		if strings.Contains(src, `("/api/notifications")`) || strings.Contains(src, "notifications-unread") {
			t.Errorf("%s still fetches the notification list itself", name)
		}
	}
	if !strings.Contains(adminCaptivatingDashboard(), "useNotificationList(unreadCount)") {
		t.Error("the dashboard's tile should select the unread count from the shared list")
	}
}

func TestRelationshipPickersShareTheirOptions(t *testing.T) {
	for name, src := range map[string]string{
		"relationship-select-field":       adminRelationshipSelectField(),
		"multi-relationship-select-field": adminMultiRelationshipSelectField(),
	} {
		for _, want := range []string{
			"useRelationshipOptions(endpoint, { enabled: open ||",
			"search: debouncedSearch,",
			`enabled: open && debouncedSearch !== "",`,
		} {
			if !strings.Contains(src, want) {
				t.Errorf("%s is missing %q", name, want)
			}
		}
		for _, banned := range []string{"queryKey:", "relationship-options", "apiClient.get("} {
			if strings.Contains(src, banned) {
				t.Errorf("%s still has %q", name, banned)
			}
		}
	}
}

func TestTableTabCountsAreQueries(t *testing.T) {
	src := adminTableTabs()
	for _, want := range []string{"useQueries({", "queryKey: resourceKeys.tabCount(endpoint, tab.key, filters)"} {
		if !strings.Contains(src, want) {
			t.Errorf("table-tabs is missing %q", want)
		}
	}
	for _, banned := range []string{"useEffect", "setCounts"} {
		if strings.Contains(src, banned) {
			t.Errorf("table-tabs still fetches counts in an effect (%q)", banned)
		}
	}
}

func TestTreeAndStatCardsKeyUnderTheEndpoint(t *testing.T) {
	tree := adminResourceTree()
	if strings.Contains(tree, "[resource.slug") {
		t.Error("the tree view still keys on the slug, which no save invalidates")
	}
	if !strings.Contains(tree, "resourceKeys.tree(resource.endpoint)") ||
		!strings.Contains(tree, "resourceKeys.all(resource.endpoint)") {
		t.Error("the tree view does not key through resourceKeys")
	}
	if !strings.Contains(adminStatCards(), "queryKey: [...resourceKeys.stats(baseEndpoint), stat.endpoint, stat.field]") {
		t.Error("PageHeader's stat cards do not key through resourceKeys")
	}
}

// Old keys must not survive anywhere in the panel, and the two fetches that
// used to be copied must each live in one file.
func TestNoAdminFileKeepsADriftedKey(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "app", Architecture: ArchTriple}
	opts.Normalize()

	files := adminFileMap(root, opts)
	files["ResourceStatCard.tsx"] = adminResourceStatCardTSX()
	files["ResourceLatestTable.tsx"] = adminResourceLatestTableTSX()

	for path, src := range files {
		for _, banned := range []string{
			`"relationship-options"`,
			`"resource-latest"`,
			`"dashboard", "resource-stats"`,
			`"dashboard", "notifications-unread"`,
		} {
			if strings.Contains(src, banned) {
				t.Errorf("%s still keys a query as %s", path, banned)
			}
		}
		base := filepath.Base(path)
		if strings.Contains(src, "/api/admin/dashboard/resource-stats/") && base != "use-resource.ts" {
			t.Errorf("%s fetches resource stats outside the shared hook", path)
		}
		if strings.Contains(src, `("/api/notifications")`) && base != "use-notifications.ts" {
			t.Errorf("%s fetches the notification list outside the shared hook", path)
		}
	}
}

func TestUseNotificationsShipsInEveryAdminShape(t *testing.T) {
	root := t.TempDir()

	triple := Options{ProjectName: "app", Architecture: ArchTriple}
	triple.Normalize()
	if _, ok := adminFileMap(root, triple)[filepath.Join(root, "apps", "admin", "hooks", "use-notifications.ts")]; !ok {
		t.Error("the Next.js admin does not get hooks/use-notifications.ts")
	}

	double := Options{ProjectName: "app", Architecture: ArchDouble}
	double.Normalize()
	if _, ok := embeddedAdminFileMap(root, double)[filepath.Join(root, "apps", "web", "admin-panel", "hooks", "use-notifications.ts")]; !ok {
		t.Error("the admin inside the web app does not get hooks/use-notifications.ts")
	}

	vite := Options{ProjectName: "app", Architecture: ArchTriple, Frontend: FrontendTanStack}
	vite.Normalize()
	if _, ok := adminTanStackFileMap(root, vite)[filepath.Join(root, "apps", "admin", "src", "hooks", "use-notifications.ts")]; !ok {
		t.Error("the Vite admin does not get src/hooks/use-notifications.ts")
	}
}
