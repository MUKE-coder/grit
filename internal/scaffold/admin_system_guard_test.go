package scaffold

import (
	"strings"
	"testing"
)

// /system is the platform's own console: jobs, backups, audit, health, roles.
// A tenant user must never reach it, whatever role they hold inside a tenant.
//
// The guard is its own component, rendered by the dashboard layouts, and not a
// block in AdminLayout. Plugins inject into admin-layout.tsx, so an upgrade
// treats that file as edited and leaves it alone: a guard written there would
// never reach a project with a plugin installed. This test pins the component
// and both places that render it.
func TestSystemPagesArePlatformStaffOnly(t *testing.T) {
	src := adminSystemGuardComponent()

	if !strings.Contains(src, `pathname.startsWith("/system")`) {
		t.Fatal("/system has no route guard: a tenant user can open the operations " +
			"console even though the API refuses its data")
	}
	for _, key := range []string{"system.view", "jobs.view", "backups.view", "audit.view", "roles.view"} {
		if !strings.Contains(src, `"`+key+`"`) {
			t.Errorf("the /system guard does not accept %q, so staff holding it are "+
				"sent away from a page they are entitled to", key)
		}
	}
	if !strings.Contains(src, `router.replace("/account")`) {
		t.Error("the /system guard does not send a non-staff user to their account")
	}
	if !strings.Contains(src, "if (userLoading || permsLoading || !user) return;") {
		t.Error("the guard runs before the user and permissions load, so it " +
			"redirects staff during the first paint")
	}

	if !strings.Contains(adminLayoutComponent(), "components/chrome/") {
		t.Fatal("the guard was written into AdminLayout again: plugins edit that " +
			"file and an upgrade would skip it")
	}
	if strings.Contains(adminLayoutComponent(), "system.view") {
		t.Error("AdminLayout still carries the /system check; a project with any " +
			"plugin installed will never receive it through grit upgrade")
	}

	next := adminDashboardLayout()
	if !strings.Contains(next, `import { SystemGuard } from "@/components/chrome/system-guard";`) ||
		!strings.Contains(next, "<SystemGuard />") {
		t.Error("the Next dashboard layout does not render the /system guard")
	}

	tanstack := adminTanStackDashboardLayout()
	if !strings.Contains(tanstack, `import { SystemGuard } from '@/components/chrome/system-guard'`) ||
		!strings.Contains(tanstack, "<SystemGuard />") {
		t.Error("the TanStack dashboard layout does not render the /system guard")
	}
}
