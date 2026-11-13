package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// spaWithPanelAt builds the two directories the shape is read from: a project
// whose panel owns /admin has src/routes/admin, and one whose panel is the
// whole app has its screens at src/routes/_dashboard.
func spaWithPanelAt(t *testing.T, segment string) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "src", "routes", segment, "_dashboard")
	if segment == "" {
		dir = filepath.Join(root, "src", "routes", "_dashboard")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "src", "admin-panel", "pages"), 0o755); err != nil {
		t.Fatal(err)
	}
	return root
}

// `grit new --api --admin` serves the panel at /, so every prefix the embed
// rewrites add points at a path the project does not have. The shape is read
// from the SPA's route tree, because these rewrites are keyed on the
// destination path and run from a writer with no Generator (grit#89).
func TestAdminServedAtRootReadsTheRouteTree(t *testing.T) {
	underAdmin := spaWithPanelAt(t, "admin")
	atRoot := spaWithPanelAt(t, "")

	for _, c := range []struct {
		name string
		root string
		rel  string
		want bool
	}{
		{"panel code beside a site", underAdmin, "src/admin-panel/pages/dashboard.tsx", false},
		{"route shim beside a site", underAdmin, "src/routes/admin/_dashboard/x.tsx", false},
		{"panel code with no site", atRoot, "src/admin-panel/pages/dashboard.tsx", true},
		{"route shim with no site", atRoot, "src/routes/_dashboard/x.tsx", true},
	} {
		path := filepath.ToSlash(filepath.Join(c.root, filepath.FromSlash(c.rel)))
		if got := adminServedAtRoot(path); got != c.want {
			t.Errorf("%s: adminServedAtRoot = %v, want %v", c.name, got, c.want)
		}
	}

	// A path that belongs to no SPA, and a project with neither directory:
	// both answer the old way, so a layout nobody has seen cannot silently
	// lose its prefixes.
	if adminServedAtRoot("/some/apps/admin/src/app/page.tsx") {
		t.Error("a standalone admin app is not served at an SPA's root")
	}

	// And a TanStack admin that IS its own app, which has routes/_dashboard
	// like the embedded one and must still be left alone: it owns @/, so
	// repointing its imports aims them at an @admin alias it never declares.
	// The first version of this probe claimed it, and every generated resource
	// in a triple project stopped compiling.
	standalone := t.TempDir()
	if err := os.MkdirAll(filepath.Join(standalone, "apps", "admin", "src", "routes", "_dashboard"), 0o755); err != nil {
		t.Fatal(err)
	}
	shim := filepath.Join(standalone, "apps", "admin", "src", "routes", "_dashboard", "x.tsx")
	if adminServedAtRoot(filepath.ToSlash(shim)) {
		t.Error("a standalone TanStack admin is its own app, not a panel inside an SPA")
	}
	if got := embeddedAdminFileContent(shim, "import { X } from '@/resources/posts'\n"); !strings.Contains(got, "'@/resources/posts'") {
		t.Errorf("a standalone admin's imports were repointed: %s", got)
	}
	bare := t.TempDir()
	if adminServedAtRoot(filepath.ToSlash(filepath.Join(bare, "src", "admin-panel", "x.tsx"))) {
		t.Error("with no route tree to read, the answer should be the old one")
	}
}

// And the routes root follows it, or a generated resource gets its route shims
// in a directory nothing serves.
func TestGeneratedRoutesFollowThePanelsSegment(t *testing.T) {
	atRoot := spaWithPanelAt(t, "")
	g := &Generator{Root: atRoot, Architecture: "single", Frontend: "tanstack"}
	if got, want := g.adminTanStackRoutesRoot(), filepath.Join(atRoot, "src", "routes"); got != want {
		t.Errorf("routes root is %s, want %s", got, want)
	}

	underAdmin := spaWithPanelAt(t, "admin")
	g = &Generator{Root: underAdmin, Architecture: "single", Frontend: "tanstack"}
	if got, want := g.adminTanStackRoutesRoot(), filepath.Join(underAdmin, "src", "routes", "admin"); got != want {
		t.Errorf("routes root is %s, want %s", got, want)
	}
}

// The content rewrite, which is where the two shapes differ in what they must
// NOT do. The aliases are repointed either way: the panel's code sits behind
// @admin wherever it is served, and @/components in the SPA is the SPA's own.
func TestGeneratedAdminContentIsPrefixedOnlyWhenThereIsAPrefix(t *testing.T) {
	const shim = `import { ResourcePage } from '@/components/resource/resource-page'
import { machineResource } from '@/resources/machines/machines'

export const Route = createFileRoute('/_dashboard/resources/machines/')({
  component: () => <ResourcePage resource={machineResource} href="/resources/machines" />,
})
`

	atRoot := spaWithPanelAt(t, "")
	final := embeddedAdminFileContent(
		filepath.Join(atRoot, "src", "routes", "_dashboard", "resources", "machines", "index.tsx"), shim)
	if !strings.Contains(final, "@admin/components/resource/resource-page") {
		t.Error("the alias was not repointed: @/components in the SPA is the SPA's own code")
	}
	if strings.Contains(final, "/admin/") {
		t.Errorf("the panel is the whole app, so nothing should be prefixed:\n%s", final)
	}
	if !strings.Contains(final, "createFileRoute('/_dashboard/resources/machines/')") {
		t.Errorf("the route id was rewritten away from the path it is at:\n%s", final)
	}

	underAdmin := spaWithPanelAt(t, "admin")
	final = embeddedAdminFileContent(
		filepath.Join(underAdmin, "src", "routes", "admin", "_dashboard", "resources", "machines", "index.tsx"), shim)
	if !strings.Contains(final, "createFileRoute('/admin/_dashboard/resources/machines/')") {
		t.Errorf("a panel under /admin still needs its route id prefixed:\n%s", final)
	}
	if !strings.Contains(final, `href="/admin/resources/machines"`) {
		t.Errorf("a panel under /admin still needs its links prefixed:\n%s", final)
	}
}
