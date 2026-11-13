package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// adminOnlyOptions is `grit new my-app --api --admin`: the single architecture
// with the Vite SPA, no public site, and the panel served at /.
func adminOnlyOptions() Options {
	return Options{
		ProjectName:  "console",
		Architecture: ArchSingle,
		Frontend:     FrontendTanStack,
		Theme:        "atlas",
		AdminOnly:    true,
	}
}

// The panel's routes are the SPA's own routes, with no /admin segment above
// them: that segment is what makes room for a public site, and there isn't one.
func TestAdminOnlyServesThePanelAtTheRoot(t *testing.T) {
	root := t.TempDir()
	files := embeddedSingleAdminFileMap(root, adminOnlyOptions())
	if len(files) < 100 {
		t.Fatalf("the panel has only %d files; the transform has stopped finding them", len(files))
	}

	for path := range files {
		rel := filepath.ToSlash(strings.TrimPrefix(path, root))
		if strings.Contains(rel, "/src/routes/admin/") {
			t.Errorf("%s is under /admin, which in this shape is a path nothing serves", rel)
		}
	}

	for _, want := range []string{
		"src/routes/index.tsx",
		"src/routes/_dashboard/dashboard.tsx",
		"src/admin-panel/admin.css",
		"src/admin-panel/pages/dashboard.tsx",
	} {
		if _, ok := files[filepath.Join(root, filepath.FromSlash(want))]; !ok {
			t.Errorf("%s is missing: the panel is not reachable at the root", want)
		}
	}

	// TanStack's file convention has no routes/route.tsx: __root.tsx is the
	// root layout. Written anyway, the generated route tree does not compile.
	if _, ok := files[filepath.Join(root, "src", "routes", "route.tsx")]; ok {
		t.Error("routes/route.tsx is written; TanStack's root layout is __root.tsx, and this one is a build error")
	}

	// And the entry point sends you to the dashboard without a prefix.
	index := files[filepath.Join(root, "src", "routes", "index.tsx")]
	if !strings.Contains(index, "createFileRoute('/')") {
		t.Errorf("the SPA's index route does not claim /:\n%s", index)
	}
	if !strings.Contains(index, "to: '/dashboard'") {
		t.Errorf("the SPA's index route does not open the dashboard:\n%s", index)
	}
}

// A file route's id has to equal its path or TanStack refuses it. At the root
// the panel's own ids are already right, so the rewrite that adds /admin has to
// stand down rather than run.
func TestAdminOnlyRouteIDsMatchTheirPaths(t *testing.T) {
	root := t.TempDir()
	routeID := regexp.MustCompile(`createFileRoute\('([^']*)'\)`)
	checked := 0

	for path, content := range embeddedSingleAdminFileMap(root, adminOnlyOptions()) {
		rel := filepath.ToSlash(strings.TrimPrefix(path, root))
		marker := "/src/routes/"
		at := strings.Index(rel, marker)
		if at < 0 || !strings.HasSuffix(rel, ".tsx") {
			continue
		}
		// writeFile applies the rewrite, so that is what a test has to check.
		final := embeddedSingleAdminContent(path, content)
		match := routeID.FindStringSubmatch(final)
		if match == nil {
			continue // a layout written as a plain component
		}
		checked++

		want := "/" + strings.TrimSuffix(rel[at+len(marker):], ".tsx")
		if strings.HasSuffix(want, "/index") {
			want = strings.TrimSuffix(want, "index")
		}
		if match[1] != want {
			t.Errorf("%s declares route id %q but its path is %q: the router will refuse it",
				rel, match[1], want)
		}
	}
	if checked < 20 {
		t.Fatalf("only %d route ids checked; the test is not looking at the routes", checked)
	}
}

// The links too. The panel's pages say href="/dashboard" and here that is the
// address, so prefixing points every one of them at an /admin that does not
// exist. Both transforms writeFile applies have to decline: the auth shells were
// rewritten by the general one before the SPA's own could say no.
func TestAdminOnlyLeavesTheLinksAlone(t *testing.T) {
	root := t.TempDir()
	checked := 0

	for path, content := range embeddedSingleAdminFileMap(root, adminOnlyOptions()) {
		rel := filepath.ToSlash(strings.TrimPrefix(path, root))
		// Both, in the order writeFile applies them.
		final := embeddedSingleAdminContent(path, embeddedAdminContent(path, content))
		checked++

		for _, quote := range []string{`"`, "`", `'`} {
			if strings.Contains(final, quote+"/admin/") || strings.Contains(final, quote+"/admin"+quote) {
				t.Errorf("%s links to /admin, which this project does not serve", rel)
				break
			}
		}
		// The alias rewrite still has to run: these files live in the SPA, and
		// @/pages is a directory it does not have. The guard that skips the
		// link prefixing must not skip this as well.
		if strings.Contains(final, `from "@/pages/`) || strings.Contains(final, `from '@/pages/`) {
			t.Errorf("%s imports @/pages, which the SPA does not have", rel)
		}
		if leftoverAliasPattern.MatchString(final) {
			t.Errorf("%s still imports @/ from the panel's own tree", rel)
		}
	}
	if checked < 100 {
		t.Fatalf("only %d files checked; the map is not being read", checked)
	}
}

// The site chrome goes with the site. navbar, footer, mega-menu and nav-menu
// exist to be rendered by _site.tsx; with no site nothing imports them but each
// other, and nav-menu carries a link to /admin/dashboard.
func TestAdminOnlyLeavesOutTheSiteChrome(t *testing.T) {
	chrome := []string{
		"src/components/navbar.tsx",
		"src/components/footer.tsx",
		"src/components/mega-menu.tsx",
		"src/lib/nav-menu.ts",
		"src/routes/_site.tsx",
		"src/routes/_site/index.tsx",
		"src/routes/_site/blog/index.tsx",
	}

	root := t.TempDir()
	if err := writeSingleFrontendFiles(root, adminOnlyOptions()); err != nil {
		t.Fatal(err)
	}
	for _, rel := range chrome {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); err == nil {
			t.Errorf("%s is written: this project has no public site to render it", rel)
		}
	}

	// The stylesheet the /admin section's layout used to import. There is no
	// such layout here, and in a Vite app a stylesheet reaches the bundle only
	// by being imported from a module that renders, so the root route does it.
	rootRoute, err := os.ReadFile(filepath.Join(root, "src", "routes", "__root.tsx"))
	if err != nil {
		t.Fatalf("the SPA has no root route: %v", err)
	}
	if !strings.Contains(string(rootRoute), "@admin/admin.css") {
		t.Error("the root route does not import the admin stylesheet, so the panel renders unstyled")
	}

	// And the other direction, because a guard that is always on is the same
	// bug wearing the opposite sign: a normal single project keeps all of it.
	plain := t.TempDir()
	if err := writeSingleFrontendFiles(plain, singleOptions()); err != nil {
		t.Fatal(err)
	}
	for _, rel := range chrome {
		if _, err := os.Stat(filepath.Join(plain, filepath.FromSlash(rel))); err != nil {
			t.Errorf("a --single project lost %s: %v", rel, err)
		}
	}
}

// adminMountsAtRoot is package state, which is how a transform keyed on a path
// learns something only the Options know. Two scaffolds in one process is a
// real case: the test binary, and `grit new` twice from a script.
func TestAdminAtRootIsNotStickyBetweenProjects(t *testing.T) {
	_ = embeddedSingleAdminFileMap(t.TempDir(), adminOnlyOptions())

	root := t.TempDir()
	files := embeddedSingleAdminFileMap(root, singleOptions())
	dashboard := filepath.Join(root, "src", "routes", "admin", "_dashboard", "dashboard.tsx")
	content, ok := files[dashboard]
	if !ok {
		t.Fatal("the --single project has no dashboard route: the previous project's answer leaked")
	}
	if !strings.Contains(embeddedSingleAdminContent(dashboard, content), "createFileRoute('/admin/_dashboard/dashboard')") {
		t.Error("the --single project's route id lost its /admin prefix to the previous project")
	}
}
