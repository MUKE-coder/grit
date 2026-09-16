package scaffold

import (
	"path/filepath"
	"strings"
	"testing"
)

// Contact-app review L24 to L28: work the admin and web frontends did on every
// render or every route for nothing, and a dashboard chart of random numbers.

func TestDataTableSelectionIsConstantTimePerRow(t *testing.T) {
	src := adminDataTable()
	for _, want := range []string{
		"useMemo(() => new Set(selectedRows), [selectedRows])",
		"isSelected={selected.has(id)}",
		"const DataTableRow = memo(function DataTableRow(",
		"const toggleRow = useCallback(",
		"const value = getNestedValue(row, col.key);",
		"{renderCell(col, value, row)}",
		"selectedRows = noSelection,",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("data-table.tsx is missing %q", want)
		}
	}
	if strings.Contains(src, "selectedRows.includes(id)") {
		t.Error("data-table.tsx still scans the selection array once per row")
	}
	if n := strings.Count(src, "getNestedValue(row, col.key)"); n != 1 {
		t.Errorf("each cell should read its value once; found %d reads", n)
	}
	// Hooks before the loading and empty early returns, or React throws when a
	// table goes from loading to loaded.
	if strings.Index(src, "const toggleRow = useCallback(") > strings.Index(src, "if (isLoading) {") {
		t.Error("a hook runs after an early return")
	}
}

func TestUsePermissionsBuildsItsSetOncePerFetch(t *testing.T) {
	src := adminUsePermissions()
	for _, want := range []string{
		"select: toPermissionState,",
		"function toPermissionState(data: MyPermissions): PermissionState {",
		"const can = useCallback(",
		"[granted, isSuper],",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("use-permissions.ts is missing %q", want)
		}
	}
	if strings.Contains(src, "new Set(data?.permissions") || strings.Contains(src, "function can(permission") {
		t.Error("use-permissions.ts still builds a Set and a can() on every render")
	}
}

func TestDropzoneRevokesItsPreviews(t *testing.T) {
	src := adminDropzone()
	if strings.Count(src, "URL.createObjectURL(") != 1 {
		t.Fatal("expected exactly one object URL to be created")
	}
	for _, want := range []string{
		"previewUrls.current.add(url);",
		"URL.revokeObjectURL(f.url);",
		"for (const url of urls) URL.revokeObjectURL(url);",
		"releasePreviews(updated, files);",
		"releasePreviews(updated, [...files, ...newFiles]);",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("dropzone.tsx is missing %q", want)
		}
	}
	// Removing a file releases its preview.
	remove := src[strings.Index(src, "const removeFile = (index: number) => {"):]
	remove = remove[:strings.Index(remove, "};")]
	if !strings.Contains(remove, "releasePreviews(") {
		t.Error("removing a file does not revoke its preview")
	}
}

func TestLayoutsDoNotPreloadTheMonoFont(t *testing.T) {
	for _, theme := range []string{"atlas", "aurora", "pulse"} {
		opts := Options{ProjectName: "app", Theme: theme}
		for name, src := range map[string]string{
			"admin": adminRootLayout(opts),
			"web":   webRootLayout(opts),
		} {
			if strings.Contains(src, "weight: [") {
				t.Errorf("%s layout (%s) pins static weights on a variable font", name, theme)
			}
			for _, line := range strings.Split(src, "\n") {
				if strings.Contains(line, `variable: "--font-mono"`) && !strings.Contains(line, "preload: false") {
					t.Errorf("%s layout (%s) preloads the mono font on every route: %s", name, theme, strings.TrimSpace(line))
				}
				if strings.Contains(line, `variable: "--font-display"`) && strings.Contains(line, "preload: false") {
					t.Errorf("%s layout (%s) stopped preloading the body font", name, theme)
				}
			}
		}
	}
}

func TestDashboardPlotsNoRandomNumbers(t *testing.T) {
	src := adminCaptivatingDashboard()
	for _, gone := range []string{"Math.random", "ActivityAreaChart", "weekSeries", "Activity, past 7 days"} {
		if strings.Contains(src, gone) {
			t.Errorf("the dashboard still contains %q", gone)
		}
	}
	if !strings.Contains(src, "<SeverityPieChart data={severitySeries} />") || !strings.Contains(src, "Recent activity") {
		t.Error("the dashboard lost a widget backed by real data")
	}
	if strings.Contains(adminDashboardCatalogTS(), "system:activity-7d") {
		t.Error("dashboard settings still offers the removed chart")
	}
	// A project whose dashboard page was edited keeps its import of the chart,
	// so the charts module must still export it.
	if !strings.Contains(adminDashboardChartsTSX(), "export function ActivityAreaChart(") {
		t.Error("DashboardCharts dropped an export an edited dashboard may still import")
	}
}

// grit upgrade writes these files whole in every panel shape, so the fixes
// reach existing projects without a repair. The test fails if one of them
// becomes user-owned or leaves the maps upgrade reads.
func TestRenderWorkFixesReachEveryPanelShape(t *testing.T) {
	root := t.TempDir()
	fixed := map[string]string{
		"components/tables/data-table.tsx": "const DataTableRow = memo(",
		"hooks/use-permissions.ts":         "select: toPermissionState,",
		"components/ui/dropzone.tsx":       "URL.revokeObjectURL(f.url);",
		"lib/dashboard-catalog.ts":         "system:severity-mix",
	}
	dashboard := "app/(dashboard)/dashboard/page.tsx"

	adminRoot := filepath.Join(root, "apps", "admin")
	triple := adminFileMap(root, tripleOptions())
	for rel, want := range fixed {
		path := filepath.Join(adminRoot, filepath.FromSlash(rel))
		if isUserOwnedAdminFile(adminRoot, path) {
			t.Errorf("%s is user-owned, so upgrade would not deliver the fix", rel)
		}
		if !strings.Contains(triple[path], want) {
			t.Errorf("admin app %s does not carry the fix", rel)
		}
	}
	if body := triple[filepath.Join(adminRoot, filepath.FromSlash(dashboard))]; body == "" || strings.Contains(body, "Math.random") {
		t.Error("the admin app's dashboard is missing or still random")
	}
	if !strings.Contains(triple[filepath.Join(adminRoot, "app", "layout.tsx")], "preload: false") {
		t.Error("the admin app's layout still preloads the mono font")
	}
	if !strings.Contains(webFileMap(root, tripleOptions())[filepath.Join(root, "apps", "web", "app", "layout.tsx")], "preload: false") {
		t.Error("the web app's layout still preloads the mono font")
	}

	for name, files := range map[string]map[string]string{
		"embedded in apps/web": embeddedAdminFileMap(root, doubleOptions()),
		"embedded in the SPA":  embeddedSingleAdminFileMap(root, singleOptions()),
	} {
		for rel, want := range fixed {
			found := false
			for path, body := range files {
				if strings.HasSuffix(filepath.ToSlash(path), "/"+rel) && strings.Contains(body, want) {
					found = true
				}
			}
			if !found {
				t.Errorf("panel %s: %s does not carry the fix", name, rel)
			}
		}
		for path, body := range files {
			if strings.HasSuffix(filepath.ToSlash(path), "/dashboard/page.tsx") && strings.Contains(body, "Math.random") {
				t.Errorf("panel %s: %s still plots random numbers", name, path)
			}
		}
	}
}
