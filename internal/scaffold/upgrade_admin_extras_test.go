package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The dashboard widgets and chart builder were written by writers upgrade
// never ran, so v3.222.0's fix to them (ask the API by its resource name, not
// the admin slug) reached new projects only. Every multi-word resource's stat
// card kept failing in a project that had upgraded.
func TestUpgradeRefreshesTheDashboardWidgets(t *testing.T) {
	root := t.TempDir()
	admin := filepath.Join(root, "apps", "admin")
	card := filepath.Join(admin, "components", "dashboard", "ResourceStatCard.tsx")
	chart := filepath.Join(admin, "components", "dashboard", "CustomChartCard.tsx")

	// The widgets as a project from before the fix has them.
	writeTestFile(t, card, `"/api/admin/dashboard/resource-stats/" + resource.slug +`)
	writeTestFile(t, chart, `"/api/admin/dashboard/chart/" + chart.resource +`)

	opts := Options{ProjectName: "app", Architecture: ArchTriple}
	opts.Normalize()
	if _, err := upgradeAdminFiles(root, opts, UpgradeOptions{}); err != nil {
		t.Fatalf("upgradeAdminFiles: %v", err)
	}

	// The card asks through the shared hook, by endpoint, and the hook takes
	// the API's name for the resource from the endpoint's last segment.
	got, _ := os.ReadFile(card)
	if !strings.Contains(string(got), `useResourceDashboardStats(resource.endpoint,`) {
		t.Error("upgrade left the stat card asking for stats by the admin slug")
	}
	if !strings.Contains(adminUseResource(), `endpoint.split("/").filter(Boolean).pop()`) {
		t.Error("the dashboard stats hook no longer asks by the API's resource name")
	}
	got, _ = os.ReadFile(chart)
	if !strings.Contains(string(got), "const apiName =") {
		t.Error("upgrade left the chart card asking by the saved slug")
	}
}
