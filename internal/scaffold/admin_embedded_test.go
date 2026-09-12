package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func doubleOptions() Options {
	return Options{ProjectName: "app", Architecture: ArchDouble, Frontend: FrontendNext, Theme: "atlas"}
}

func tripleOptions() Options {
	return Options{ProjectName: "app", Architecture: ArchTriple, Frontend: FrontendNext, Theme: "atlas"}
}

// A double's admin panel lives inside the web app: routes under app/admin, code
// under admin-panel, and none of the standalone app's shell.
func TestEmbeddedAdminLandsInsideTheWebApp(t *testing.T) {
	root := t.TempDir()
	files := embeddedAdminFileMap(root, doubleOptions())
	if len(files) < 100 {
		t.Fatalf("the embedded panel has only %d files; the transform has stopped finding them", len(files))
	}

	for path := range files {
		rel, err := filepath.Rel(root, path)
		if err != nil {
			t.Fatal(err)
		}
		rel = filepath.ToSlash(rel)
		if !strings.HasPrefix(rel, "apps/web/") {
			t.Errorf("%s is outside the web app", rel)
		}
		// apps/admin must not be written at all: a double has no admin app, and a
		// file there is one nothing compiles.
		if strings.HasPrefix(rel, "apps/admin/") {
			t.Errorf("%s was written to the admin app, which a double does not have", rel)
		}
	}

	// The pieces that make it a section of the web app rather than a copy of an
	// app: its own layout, and the half of the stylesheet that is about screens.
	for _, want := range []string{
		"apps/web/app/admin/layout.tsx",
		"apps/web/app/admin/admin.css",
		"apps/web/app/admin/(dashboard)/dashboard/page.tsx",
		"apps/web/admin-panel/components/layout/admin-layout.tsx",
		"apps/web/admin-panel/lib/api-client.ts",
		"apps/web/admin-panel/resources/index.ts",
	} {
		if _, ok := files[filepath.Join(root, filepath.FromSlash(want))]; !ok {
			t.Errorf("%s is missing from the embedded panel", want)
		}
	}

	// And the app shell it must NOT bring: the web app has these already, and a
	// second package.json or globals.css inside it is either ignored or harmful.
	for _, unwanted := range []string{
		"apps/web/admin-panel/package.json",
		"apps/web/admin-panel/next.config.ts",
		"apps/web/admin-panel/tsconfig.json",
		"apps/web/app/admin/globals.css",
	} {
		if _, ok := files[filepath.Join(root, filepath.FromSlash(unwanted))]; ok {
			t.Errorf("%s should not be written: the web app owns it", unwanted)
		}
	}
}

// Imports and links are repointed, or the panel compiles against the web app's own
// files and links to pages that do not exist.
func TestEmbeddedAdminContentIsRepointed(t *testing.T) {
	root := t.TempDir()
	opts := doubleOptions()

	for path, content := range embeddedAdminFileMap(root, opts) {
		// writeFile applies the rewrite, so that is what a test has to check.
		final := embeddedAdminContent(path, content)
		rel := filepath.ToSlash(strings.TrimPrefix(path, root))

		if leftoverAliasPattern.MatchString(final) {
			t.Errorf("%s still imports @/ from the admin's own tree: inside the web app that is the web app's code",
				rel)
		}
		for _, segment := range adminRoutePrefixes {
			if strings.Contains(final, `href="/`+segment) {
				t.Errorf("%s links to /%s, which inside the web app is not a page", rel, segment)
			}
		}
		if strings.Contains(final, "/admin/admin/") {
			t.Errorf("%s has a doubled prefix, so the rewrite ran twice", rel)
		}
	}
}

// Every route the panel writes has its first segment in adminRoutePrefixes.
//
// The rewrite sees one file at a time and cannot derive the list, so the list is
// written down. This is what keeps it honest: add a screen at /admin/reports and
// this fails until /reports is listed, rather than shipping a sidebar link to a
// page that 404s.
func TestEmbeddedAdminRoutesAreAllPrefixed(t *testing.T) {
	root := t.TempDir()
	known := map[string]bool{}
	for _, segment := range adminRoutePrefixes {
		known[segment] = true
	}

	for path := range embeddedAdminFileMap(root, doubleOptions()) {
		rel := filepath.ToSlash(strings.TrimPrefix(path, root))
		marker := "/apps/web/app/admin/"
		at := strings.Index(rel, marker)
		if at < 0 {
			continue
		}
		for _, segment := range strings.Split(rel[at+len(marker):], "/") {
			if segment == "" || strings.HasPrefix(segment, "(") {
				continue // a route group is not part of the URL
			}
			if strings.Contains(segment, ".") || strings.HasPrefix(segment, "[") {
				break // a file, or a dynamic segment: no first segment of its own
			}
			if !known[segment] {
				t.Errorf("the panel has a route at /admin/%s and %q is not in adminRoutePrefixes, so links to it are not rewritten",
					segment, segment)
			}
			break
		}
	}
}

// A triple is untouched by any of this.
func TestTripleAdminIsUnchanged(t *testing.T) {
	triple := tripleOptions()
	if triple.ShouldEmbedAdmin() {
		t.Fatal("a triple must not embed the panel: it has an admin app")
	}
	if !triple.HasAdminPanel() || !doubleOptions().HasAdminPanel() {
		t.Error("both shapes have an admin panel, wherever it lives")
	}

	// The web app's dependencies are the admin's only where the admin is inside it.
	if deps := webAdminDependencies(triple); deps != "" {
		t.Error("a triple's web app must not carry the admin's dependencies")
	}
	if deps := webAdminDependencies(doubleOptions()); !strings.Contains(deps, "@tiptap/react") {
		t.Error("a double's web app needs the admin's editor dependency")
	}

	// And the link: a triple points at the admin app, a double at its own route.
	if got := adminHref(`const ADMIN_URL = "{{ADMIN_HREF}}";`, triple); !strings.Contains(got, "localhost:3001") {
		t.Errorf("a triple should link to the admin app, got %q", got)
	}
	if got := adminHref(`const ADMIN_URL = "{{ADMIN_HREF}}";`, doubleOptions()); !strings.Contains(got, "/admin/dashboard") {
		t.Errorf("a double should link to its own admin route, got %q", got)
	}
}

// The whole panel reaches disk, including the screens written by feature writers
// rather than by the map: that is how the first build failed, on a dashboard widget.
func TestEmbeddedAdminWritesTheFeatureScreensToo(t *testing.T) {
	root := t.TempDir()
	opts := doubleOptions()
	if err := writeEmbeddedAdminFiles(root, opts); err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{
		// Written by adminExtraWriters, which the embedded path has to run.
		"apps/web/admin-panel/components/dashboard/ResourceWidgetsRow.tsx",
		"apps/web/admin-panel/components/dashboard/CustomChartCard.tsx",
	} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(want))); err != nil {
			t.Errorf("%s was not written: %v", want, err)
		}
	}

	// The account-security screen comes from a feature writer that asks the path
	// helpers, so it only lands if those know about the embedded layout.
	if err := writeAdminSecurityFiles(root, opts); err != nil {
		t.Fatal(err)
	}
	security := filepath.Join(root, filepath.FromSlash("apps/web/app/admin/(dashboard)/account/security/page.tsx"))
	if _, err := os.Stat(security); err != nil {
		t.Errorf("the account-security screen is missing from the embedded panel: %v", err)
	}
}
