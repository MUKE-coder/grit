package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func singleOptions() Options {
	return Options{ProjectName: "app", Architecture: ArchSingle, Theme: "atlas"}
}

// The panel lands in the SPA: its code behind its own alias, its routes in the
// SPA's route tree, and nothing in apps/, which a single project does not have.
func TestSingleAdminLandsInsideTheSPA(t *testing.T) {
	root := t.TempDir()
	files := embeddedSingleAdminFileMap(root, singleOptions())
	if len(files) < 100 {
		t.Fatalf("the SPA panel has only %d files; the transform has stopped finding them", len(files))
	}

	routes, panel := 0, 0
	for path := range files {
		rel, err := filepath.Rel(root, path)
		if err != nil {
			t.Fatal(err)
		}
		rel = filepath.ToSlash(rel)
		switch {
		case strings.HasPrefix(rel, "frontend/src/routes/admin/"):
			routes++
		case strings.HasPrefix(rel, "frontend/src/admin-panel/"):
			panel++
		default:
			t.Errorf("%s is outside the SPA's panel and its routes", rel)
		}
	}
	if routes < 20 {
		t.Errorf("only %d route files: the panel's screens are not reachable", routes)
	}
	if panel < 100 {
		t.Errorf("only %d panel files: its components are not all there", panel)
	}

	for _, want := range []string{
		// The section's own entry point, layout and stylesheet.
		"frontend/src/routes/admin/index.tsx",
		"frontend/src/routes/admin/route.tsx",
		"frontend/src/admin-panel/admin.css",
		"frontend/src/routes/admin/_dashboard/dashboard.tsx",
		"frontend/src/admin-panel/pages/dashboard.tsx",
		"frontend/src/admin-panel/lib/api-client.ts",
		"frontend/src/admin-panel/resources/index.ts",
	} {
		if _, ok := files[filepath.Join(root, filepath.FromSlash(want))]; !ok {
			t.Errorf("%s is missing from the SPA panel", want)
		}
	}

	// The app shell the SPA already has. A second main.tsx or index.html inside it
	// is either ignored or harmful, and a second root route is a routing conflict.
	for _, unwanted := range []string{
		"frontend/src/admin-panel/main.tsx",
		"frontend/src/admin-panel/index.html",
		"frontend/src/admin-panel/package.json",
		"frontend/src/routes/admin/__root.tsx",
	} {
		if _, ok := files[filepath.Join(root, filepath.FromSlash(unwanted))]; ok {
			t.Errorf("%s should not be written: the SPA owns it", unwanted)
		}
	}
}

// A file route's id has to equal its path or TanStack refuses it, and the path now
// starts with /admin. This is the check that the ids were rewritten with the files.
func TestSingleAdminRouteIDsMatchTheirPaths(t *testing.T) {
	root := t.TempDir()
	routeID := regexp.MustCompile(`createFileRoute\('([^']*)'\)`)
	checked := 0

	for path, content := range embeddedSingleAdminFileMap(root, singleOptions()) {
		rel := filepath.ToSlash(strings.TrimPrefix(path, root))
		marker := "/frontend/src/routes/"
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
		want = strings.TrimSuffix(want, "/route")
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

// Imports and links are repointed, or the panel compiles against the SPA's own
// files and links to pages that do not exist.
func TestSingleAdminContentIsRepointed(t *testing.T) {
	root := t.TempDir()

	for path, content := range embeddedSingleAdminFileMap(root, singleOptions()) {
		final := embeddedSingleAdminContent(path, content)
		rel := filepath.ToSlash(strings.TrimPrefix(path, root))

		if leftoverAliasPattern.MatchString(final) {
			t.Errorf("%s still imports @/ from the panel's own tree: in the SPA that is the SPA's code", rel)
		}
		if strings.Contains(final, `from "@/pages/`) || strings.Contains(final, `from '@/pages/`) {
			t.Errorf("%s imports a page the SPA does not have", rel)
		}
		for _, segment := range adminRoutePrefixes {
			if strings.Contains(final, `href="/`+segment) {
				t.Errorf("%s links to /%s, which in the SPA is not a page", rel, segment)
			}
		}
		if strings.Contains(final, "/admin/admin/") {
			t.Errorf("%s has a doubled prefix, so the rewrite ran twice", rel)
		}
		// Vite does not polyfill process, and esbuild does not typecheck, so this
		// reaches the browser as "process is not defined" rather than a build error.
		if strings.Contains(final, "process.env") {
			t.Errorf("%s reads process.env, which does not exist in a Vite app", rel)
		}
		// next/link and next/navigation are Next.js. The compat shim is the one
		// file allowed to name them, in its comments.
		if strings.Contains(rel, "next-compat") {
			continue
		}
		for _, next := range []string{`from "next/`, `from 'next/`} {
			if strings.Contains(final, next) {
				t.Errorf("%s imports from next/, which a Vite app cannot resolve", rel)
			}
		}
	}
}

// The API base. A triple admin is on its own port and talks across origins; here
// the binary serves the page and the API, so the origin of the page is the API.
func TestSingleAdminTalksToItsOwnOrigin(t *testing.T) {
	root := t.TempDir()
	client := filepath.Join(root, "frontend", "src", "admin-panel", "lib", "api-client.ts")
	content, ok := embeddedSingleAdminFileMap(root, singleOptions())[client]
	if !ok {
		t.Fatal("the panel has no api-client")
	}
	final := embeddedSingleAdminContent(client, content)
	if strings.Contains(final, "localhost:8080") {
		t.Error("the panel points at localhost:8080, which from the dev server is a cross-origin request the API does not allow")
	}
	if !strings.Contains(final, "window.location.origin") {
		t.Error("the panel should fall back to the origin serving it")
	}
}

// Everything reaches disk, including the screens written by feature writers rather
// than by the map, and nothing reaches an apps/ directory.
func TestSingleAdminWritesTheFeatureScreensAndNothingElse(t *testing.T) {
	root := t.TempDir()
	opts := singleOptions()
	if err := writeEmbeddedSingleAdminFiles(root, opts); err != nil {
		t.Fatal(err)
	}
	if err := writeAdminSecurityFiles(root, opts); err != nil {
		t.Fatal(err)
	}
	if err := writeAdminPasskeyFiles(root, opts); err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{
		"frontend/src/admin-panel/components/dashboard/ResourceWidgetsRow.tsx",
		"frontend/src/admin-panel/components/dashboard/CustomChartCard.tsx",
		// The security screen is a page plus a route shim in a TanStack app. A
		// single project is one whatever Frontend says, which is what
		// AdminIsTanStack exists to answer.
		"frontend/src/admin-panel/pages/account-security.tsx",
		"frontend/src/routes/admin/_dashboard/account/security.tsx",
		"frontend/src/admin-panel/components/security/passkeys.tsx",
	} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(want))); err != nil {
			t.Errorf("%s was not written: %v", want, err)
		}
	}

	// An apps/ directory in a single project holds files nothing compiles: no
	// package.json, no build, nothing that reads them.
	if _, err := os.Stat(filepath.Join(root, "apps")); err == nil {
		t.Error("a single project grew an apps/ directory")
	}
}

// The SPA knows how to resolve the panel: aliases in both tsconfig and the Vite
// config, because Vite does not read tsconfig paths.
func TestSingleSPAResolvesThePanel(t *testing.T) {
	ts := singleFrontendTSConfig()
	for _, want := range []string{
		`"@admin/*": ["./src/admin-panel/*"]`,
		`"@repo/shared/brand": ["./src/shared/brand.config.ts"]`,
		`"@repo/upload": ["../packages/upload/src/index.ts"]`,
	} {
		if !strings.Contains(ts, want) {
			t.Errorf("the SPA tsconfig is missing %s", want)
		}
	}

	vite := singleFrontendViteConfig()
	for _, want := range []string{
		"'@admin': path.resolve",
		"'@repo/shared/brand': path.resolve",
		"'@repo/shared': path.resolve",
		"'@repo/upload/web': path.resolve",
	} {
		if !strings.Contains(vite, want) {
			t.Errorf("the SPA vite config is missing %s", want)
		}
	}
	// The subpath aliases have to come first: Vite matches a prefix, so
	// @repo/shared would swallow @repo/shared/brand and resolve it to a file that
	// does not exist.
	if strings.Index(vite, "'@repo/shared/brand'") > strings.Index(vite, "'@repo/shared':") {
		t.Error("@repo/shared is aliased before its own subpath, so the subpath never matches")
	}
	if strings.Index(vite, "'@repo/upload/web'") > strings.Index(vite, "'@repo/upload':") {
		t.Error("@repo/upload is aliased before its own subpath")
	}
}

// The panel's dependencies, and the mirrored shared package it imports as values.
func TestSinglePanelHasWhatItImports(t *testing.T) {
	pkg := singleFrontendPackageJSON(singleOptions())
	for _, dep := range []string{"@tiptap/react", "recharts", "sonner", "xlsx", "react-hook-form"} {
		if !strings.Contains(pkg, `"`+dep+`"`) {
			t.Errorf("the SPA does not depend on %s, which the panel imports", dep)
		}
	}

	root := t.TempDir()
	mirror := singleSharedMirrorFiles(root, singleOptions())
	for _, want := range []string{
		// Re-exported by the barrels, so a missing one fails the build on an
		// import of "./money" rather than on anything the panel wrote.
		"frontend/src/shared/schemas/money.ts",
		"frontend/src/shared/types/money.ts",
		"frontend/src/shared/types/errors.ts",
		"frontend/src/shared/brand.config.ts",
	} {
		if _, ok := mirror[filepath.Join(root, filepath.FromSlash(want))]; !ok {
			t.Errorf("%s is not mirrored into the SPA", want)
		}
	}
}

// The site's chrome steps aside for the panel, and only for the panel.
func TestSingleRootRouteYieldsToThePanel(t *testing.T) {
	spa := webTanStackRootRoute(singleOptions())
	if !strings.Contains(spa, "startsWith('/admin/')") {
		t.Error("the SPA's root route wraps the panel in the site navbar")
	}
	if !strings.Contains(spa, "<Navbar />") {
		t.Error("the SPA's root route lost the site navbar entirely")
	}

	// A monorepo web app has no panel inside it, so its root route is unchanged.
	web := webTanStackRootRoute(Options{ProjectName: "app", Architecture: ArchTriple, Frontend: FrontendTanStack})
	if strings.Contains(web, "/admin") {
		t.Error("a web app with no embedded panel should not test for /admin")
	}
}

// The navbar links to the panel, so it can be found without knowing the URL.
func TestSingleNavbarLinksToThePanel(t *testing.T) {
	if !strings.Contains(singleViteNavbar(singleOptions()), `"/admin/dashboard"`) {
		t.Error("the SPA navbar has no link to the panel")
	}
}

// The wiring upgrade applies to an existing project: idempotent, and it never
// replaces a value the project chose for itself.
func TestInsertMissingAfterIsIdempotentAndRespectsWhatIsThere(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tsconfig.json")
	original := `{
  "compilerOptions": {
    "paths": {
      "@/*": ["./src/*"],
      "@admin/*": ["./somewhere/else/*"]
    }
  }
}
`
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}
	lines := []string{
		`      "@admin/*": ["./src/admin-panel/*"],`,
		`      "@repo/shared/*": ["./src/shared/*"],`,
	}

	changed, err := insertMissingAfter(path, `"paths": {`, lines)
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("the missing path was not added")
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(after), `"@admin/*"`) != 1 {
		t.Error("the project's own @admin alias was duplicated rather than left alone")
	}
	if !strings.Contains(string(after), "./somewhere/else") {
		t.Error("the project's own choice was overwritten")
	}
	if !strings.Contains(string(after), `"@repo/shared/*"`) {
		t.Error("the missing alias was not inserted")
	}

	again, err := insertMissingAfter(path, `"paths": {`, lines)
	if err != nil {
		t.Fatal(err)
	}
	if again {
		t.Error("a second run changed the file again, so upgrade is not idempotent")
	}
}

// A single has a panel, and it is not the triple's separate app.
func TestSinglePanelPredicates(t *testing.T) {
	single := singleOptions()
	if !single.ShouldEmbedAdminInSPA() || !single.HasAdminPanel() {
		t.Error("a single project has an admin panel, inside its SPA")
	}
	if single.ShouldIncludeAdmin() || single.ShouldEmbedAdmin() {
		t.Error("a single must not get an admin app or a web app route group")
	}
	if !single.AdminIsTanStack() {
		t.Error("a single's SPA is a TanStack app whatever Frontend says")
	}
	if adminCodeRoot("r", single) != filepath.Join("r", "frontend", "src", "admin-panel") {
		t.Errorf("the panel's code root is wrong: %s", adminCodeRoot("r", single))
	}
	if got := adminPath("r", single, "routes", "_dashboard", "x.tsx"); got != filepath.Join("r", "frontend", "src", "routes", "admin", "_dashboard", "x.tsx") {
		t.Errorf("a route file went to %s, where the SPA's router will not find it", got)
	}
}
