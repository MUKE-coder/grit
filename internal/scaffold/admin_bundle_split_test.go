package scaffold

import (
	"strings"
	"testing"
)

// SheetJS (about 500 KB) and recharts (about 400 KB) were static imports of
// modules on every resource list page and the dashboard, so both rode in those
// pages' first load. They must only ever be imported dynamically there.

func TestExcelUtilsLoadsSheetJSOnDemand(t *testing.T) {
	src := adminExcelUtils()
	if strings.Contains(src, `from "xlsx"`) {
		t.Error(`excel-utils imports "xlsx" statically, which puts SheetJS in the first load of every resource list page`)
	}
	if !strings.Contains(src, `import("xlsx")`) {
		t.Error("excel-utils no longer loads SheetJS at all")
	}
	for _, fn := range []string{"export async function exportToFile(", "export async function downloadImportTemplate("} {
		if !strings.Contains(src, fn) {
			t.Errorf("%s is not async, so it cannot wait for SheetJS", strings.TrimSuffix(fn, "("))
		}
	}
	if !strings.Contains(adminExportMenu(), "await exportToFile(") {
		t.Error("the export menu reports success before the file is written")
	}
}

func TestDashboardLoadsRechartsOnDemand(t *testing.T) {
	for name, src := range map[string]string{
		"dashboard page":     adminCaptivatingDashboard(),
		"resource stat card": adminResourceStatCardTSX(),
	} {
		if strings.Contains(src, `from "recharts"`) {
			t.Errorf(`the %s imports "recharts" statically, which puts it in the dashboard's first load`, name)
		}
		if !strings.Contains(src, "<Suspense") {
			t.Errorf("the %s renders a lazy chart without Suspense; the Vite admin's dynamic() is React.lazy", name)
		}
	}
	for name, src := range map[string]string{
		"DashboardCharts":   adminDashboardChartsTSX(),
		"ResourceSparkline": adminResourceSparklineTSX(),
	} {
		if !strings.Contains(src, `from "recharts"`) {
			t.Errorf("%s no longer holds the recharts code", name)
		}
	}
	if !strings.Contains(adminCaptivatingDashboard(), `import("@/components/dashboard/DashboardCharts")`) {
		t.Error("the dashboard page does not load its charts module")
	}
	if !strings.Contains(adminResourceStatCardTSX(), `import("@/components/dashboard/ResourceSparkline")`) {
		t.Error("the stat card does not load its sparkline module")
	}
}
