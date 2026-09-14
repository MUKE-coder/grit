package scaffold

import (
	"strings"
	"testing"
)

// The admin layout used to return a spinner in place of the page until
// /auth/me answered, so no page query could start before it did. The page now
// renders on the first browser pass behind an overlay. Not on the server:
// admin pages read browser state while rendering, and prerendering them
// failed the build.
func TestAdminLayoutRendersThePageBeforeTheUserIsKnown(t *testing.T) {
	layout := adminLayoutComponent()

	if strings.Contains(layout, "if (isLoading) {") || strings.Contains(layout, "if (!user) {") {
		t.Error("the layout still returns early while /auth/me is in flight, which holds back every page query")
	}
	main := strings.Index(layout, "<main")
	if main < 0 || !strings.Contains(layout[main:], "{inBrowser ? children : null}</main>") {
		t.Fatal("the layout no longer renders the page in <main>, in the browser only")
	}
	if !strings.Contains(layout, "useSyncExternalStore(subscribeNever, () => true, () => false)") {
		t.Error("the page is no longer kept out of the server render")
	}
	if !strings.Contains(layout, `{!user && (`) || !strings.Contains(layout, `aria-busy="true"`) {
		t.Error("the loading overlay is missing, so a signed-out visitor would see the page before the redirect")
	}
	if !strings.Contains(layout, "usePermissions();") {
		t.Error("permissions are not requested beside /auth/me")
	}
	// Plugins inject beside these markers; losing one drops their component.
	for _, marker := range []string{"// grit:layout:imports", "{/* grit:layout:banner */}"} {
		if !strings.Contains(layout, marker) {
			t.Errorf("marker %s is gone", marker)
		}
	}
	if !strings.Contains(nextToTanStack(layout), "usePermissions();") {
		t.Error("the TanStack admin lost the parallel permissions request")
	}
}

func TestAdminDashboardLayoutIsAServerComponent(t *testing.T) {
	if strings.Contains(adminDashboardLayout(), "use client") {
		t.Error(`the (dashboard) layout is marked "use client"; AdminLayout is already the client boundary`)
	}
}
