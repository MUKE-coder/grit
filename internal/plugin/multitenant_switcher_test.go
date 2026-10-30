package plugin

import (
	"strings"
	"testing"
)

// grit#141: the plugin resolved the active organization from X-Organization-ID
// and nothing in a generated client ever set it. A user who belongs to one
// organization is fine, because the middleware activates it automatically; a
// user who belongs to two had no way to choose, which is the case the plugin
// exists for.

func switcherFiles(t *testing.T, arch, frontend string) map[string]string {
	t.Helper()
	return multitenantFiles(Context{
		Module:       "example/apps/api",
		Architecture: arch,
		Frontend:     frontend,
	})
}

// A project with an admin gets the switcher and the hook.
func TestTheSwitcherIsWrittenWhereThereIsAnAdmin(t *testing.T) {
	for _, arch := range []string{"triple", "full"} {
		files := switcherFiles(t, arch, "next")
		for _, want := range []string{
			"apps/admin/hooks/use-organizations.ts",
			"apps/admin/components/chrome/OrganizationSwitcher.tsx",
		} {
			if _, ok := files[want]; !ok {
				t.Errorf("%s: %s was not written, so somebody in two organizations "+
					"has no way to choose between them (grit#141)", arch, want)
			}
		}
	}
}

// And a project with no admin gets neither, rather than files nothing compiles.
func TestNoSwitcherWhereThereIsNoAdmin(t *testing.T) {
	for _, arch := range []string{"api", "single", "double"} {
		for path := range switcherFiles(t, arch, "next") {
			if strings.Contains(path, "apps/admin") {
				t.Errorf("%s: wrote %s into a project with no admin app", arch, path)
			}
		}
	}
}

// The TanStack admin keeps its source under src/, and its components cannot
// use "use client" or next/navigation.
func TestTheSwitcherIsConvertedForTanStack(t *testing.T) {
	files := switcherFiles(t, "triple", "tanstack")

	src, ok := files["apps/admin/src/components/chrome/OrganizationSwitcher.tsx"]
	if !ok {
		t.Fatalf("the switcher is not under src/ for TanStack; got %v", switcherKeys(files))
	}
	if strings.Contains(src, `"use client"`) {
		t.Error(`the TanStack switcher still has "use client", which Vite does not understand`)
	}
	hook, ok := files["apps/admin/src/hooks/use-organizations.ts"]
	if !ok {
		t.Fatal("the hook is not under src/ for TanStack")
	}
	if strings.Contains(hook, `"use client"`) {
		t.Error(`the TanStack hook still has "use client"`)
	}
}

func switcherKeys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// Switching organization must clear the cache, not invalidate it.
//
// An invalidated query keeps its old data while it refetches, and that data
// belongs to the organization being left: React Query would go on rendering
// one tenant's rows under another tenant's name. That is the bug that would
// make this feature worse than not having it.
func TestSwitchingClearsTheQueryCache(t *testing.T) {
	hook := switcherFiles(t, "triple", "next")["apps/admin/hooks/use-organizations.ts"]

	if !strings.Contains(hook, "queryClient.clear()") {
		t.Error("switching organization does not clear the query cache, so every " +
			"list on screen keeps showing the previous organization's rows under " +
			"the new organization's name (grit#141)")
	}
	if strings.Contains(hook, "invalidateQueries") {
		t.Error("switching organization invalidates rather than clears. An " +
			"invalidated query serves its old data while it refetches, and that " +
			"data is the organization being left")
	}
}

// A stored id for an organization the reader no longer belongs to must be
// dropped, or every request carries an organization the API refuses and the
// app looks broken rather than moving them back.
func TestARevokedMembershipClearsTheStoredChoice(t *testing.T) {
	hook := switcherFiles(t, "triple", "next")["apps/admin/hooks/use-organizations.ts"]
	if !strings.Contains(hook, "!organizations.some((o) => o.id === activeID)") {
		t.Error("a stored organization id is never checked against the list, so a " +
			"revoked membership leaves every request carrying an organization the " +
			"API refuses (grit#141)")
	}
}

// The switcher is injected where the layout has a marker, not into the navbar.
// The navbar is rewritten by every upgrade and has no marker, so a switcher
// put there would vanish from a project that upgraded.
func TestTheSwitcherIsInjectedAtTheLayoutMarker(t *testing.T) {
	var found int
	for _, inj := range multitenantInjections(Context{
		Module: "example/apps/api", Architecture: "triple", Frontend: "next",
	}) {
		if !strings.Contains(inj.File, "admin-layout.tsx") {
			continue
		}
		found++
		if !strings.HasPrefix(inj.Marker, "// grit:layout:imports") &&
			!strings.HasPrefix(inj.Marker, "{/* grit:layout:banner") {
			t.Errorf("injected at %q, which is not a marker the upgrade respects", inj.Marker)
		}
	}
	if found != 2 {
		t.Errorf("%d layout injections, want 2 (the import and the element)", found)
	}
}

// The hook must call "/api/organizations", not "/organizations".
//
// versionedPath in api-core only rewrites a path that starts with /api/, so a
// bare path goes to the host root and 404s. Every other generated hook writes
// "/api/...". The first version of this hook did not, and nothing caught it:
// it compiled, tsc passed, and these tests passed. The browser 404'd.
func TestTheHookCallsTheVersionedPath(t *testing.T) {
	hook := switcherFiles(t, "triple", "next")["apps/admin/hooks/use-organizations.ts"]

	if strings.Contains(hook, `apiClient.get("/organizations")`) {
		t.Error(`the hook calls "/organizations". versionedPath only rewrites a ` +
			`path starting with /api/, so this goes to the host root and 404s ` +
			`(grit#141)`)
	}
	if !strings.Contains(hook, `apiClient.get("/api/organizations")`) {
		t.Error(`the hook does not call "/api/organizations", which is how every ` +
			`other generated hook reaches the API`)
	}
}
