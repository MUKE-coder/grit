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

	got, _ := os.ReadFile(card)
	if !strings.Contains(string(got), `resource.endpoint.split("/").filter(Boolean).pop()`) {
		t.Error("upgrade left the stat card asking for stats by the admin slug")
	}
	got, _ = os.ReadFile(chart)
	if !strings.Contains(string(got), "const apiName =") {
		t.Error("upgrade left the chart card asking by the saved slug")
	}
}
