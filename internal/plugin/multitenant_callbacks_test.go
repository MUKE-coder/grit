package plugin

import (
	"strings"
	"testing"
)

/*
 * Scoping has to cover every processor that can READ.
 *
 * It did not. RegisterScoping registered Query, Update, Delete and Create,
 * and GORM has a fifth: Row. Scan goes through Row, not Query:
 *
 *	func (db *DB) Scan(dest any) *DB      { ... tx.Rows() ... }
 *	func (db *DB) Rows() (*sql.Rows, err) { tx.callbacks.Row().Execute(tx) }
 *
 * So db.Model(&Order{}).Scan(&out) on a tenant-owned model returned every
 * organization's rows while Find and First were scoped correctly. Reported
 * privately by a user; fixed in plugin 1.0.1.
 *
 * The generated tenant_test.go proves the behaviour against a real database.
 * This proves the registration still exists, which is the thing a refactor
 * silently drops: the behaviour test only runs in somebody's project, and
 * this runs here, on every build.
 */
func TestTenantScopingCoversEveryReadPath(t *testing.T) {
	src := mtTenantPackage(Context{Module: "example.com/app"})

	// processor -> why it has to be scoped
	required := map[string]string{
		`db.Callback().Query()`:  "Find, First, Take and Pluck",
		`db.Callback().Row()`:    "Scan, Rows and Row. Scan calls Rows, which executes THIS processor and not Query",
		`db.Callback().Update()`: "an update must not reach another organization's row",
		`db.Callback().Delete()`: "a delete must not reach another organization's row",
		`db.Callback().Create()`: "a new row is stamped with the organization",
	}

	for call, why := range required {
		if !strings.Contains(src, call) {
			t.Errorf("RegisterScoping does not register %s.\n\nIt covers %s.\n\n"+
				"Dropping it is a cross-tenant data leak, not a regression in "+
				"a feature: one organization reads another's rows and nothing "+
				"errors.", call, why)
		}
	}
}

// The fix is one line and the kind of line a later edit removes without
// noticing, so the generated project carries tests that would catch it
// against a real database.
func TestTenantTestsCoverScanAndRows(t *testing.T) {
	src := mtTenantTests(Context{Module: "example.com/app"})

	for _, name := range []string{
		"TestScan_IsScoped",
		"TestRows_IsScoped",
		"TestAggregateScan_IsScoped",
		"TestRaw_IsNotScoped",
	} {
		if !strings.Contains(src, name) {
			t.Errorf("the generated tenant tests no longer include %s", name)
		}
	}
}
