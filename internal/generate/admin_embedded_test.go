package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Three shapes, three answers. Getting this wrong writes a resource's admin
// screens into a directory the project does not have, and the resource exists
// everywhere except the admin.
func TestAdminRootsPerArchitecture(t *testing.T) {
	for _, tc := range []struct {
		arch string
		code string
	}{
		{"triple", filepath.Join("apps", "admin")},
		{"double", filepath.Join("apps", "web", "admin-panel")},
		{"single", filepath.Join("frontend", "src", "admin-panel")},
	} {
		g := &Generator{Root: "r", Architecture: tc.arch}
		if got := g.AdminCodeRoot(); got != filepath.Join("r", tc.code) {
			t.Errorf("%s: code root is %s, want %s", tc.arch, got, filepath.Join("r", tc.code))
		}
		// The registry injection reads AdminRoot. It used to be a second copy of
		// this logic, and a resource generated into a single was written and never
		// registered, so the screens existed and the sidebar did not know.
		if g.AdminRoot() != g.AdminCodeRoot() {
			t.Errorf("%s: AdminRoot and AdminCodeRoot disagree", tc.arch)
		}
	}
}

// In the SPA the router is the SPA's, so a generated route file goes into its tree
// and its id carries the /admin segment.
func TestGeneratedSPARoutesAreUnderAdmin(t *testing.T) {
	g := &Generator{Root: "r", Architecture: "single"}
	if !g.AdminIsTanStack() {
		t.Fatal("a single project's panel is a TanStack app")
	}
	want := filepath.Join("r", "frontend", "src", "routes", "admin")
	if got := g.adminTanStackRoutesRoot(); got != want {
		t.Errorf("routes root is %s, want %s", got, want)
	}
	if got := g.tanStackResourcesRoot(); got != filepath.Join("r", "frontend", "src", "admin-panel", "resources") {
		t.Errorf("resources root is %s", got)
	}

	route := filepath.Join("r", "frontend", "src", "routes", "admin", "_dashboard", "resources", "invoices", "index.tsx")
	content := `import { createFileRoute } from '@tanstack/react-router'
import { ResourcePage } from '@/components/resource/resource-page'
import { invoiceResource } from '@/resources/invoices/invoices'

export const Route = createFileRoute('/_dashboard/resources/invoices/')({
  component: () => <ResourcePage resource={invoiceResource} />,
})
`
	final := embeddedAdminFileContent(route, content)
	if !strings.Contains(final, "createFileRoute('/admin/_dashboard/resources/invoices/')") {
		t.Errorf("the route id was not rewritten to match its path:\n%s", final)
	}
	if strings.Contains(final, `'@/components/`) || strings.Contains(final, `'@/resources/`) {
		t.Errorf("the panel's imports still point at the SPA's own code:\n%s", final)
	}

	// A page written for the panel itself, reached through the same alias.
	page := filepath.Join("r", "frontend", "src", "admin-panel", "pages", "reports.tsx")
	if got := embeddedAdminFileContent(page, `import { Thing } from "@/components/thing";`); !strings.Contains(got, `"@admin/components/thing"`) {
		t.Errorf("a panel page was not repointed: %s", got)
	}

	// And a file that is not the panel's is untouched, in every shape.
	web := filepath.Join("r", "frontend", "src", "routes", "blog", "index.tsx")
	plain := `import { Card } from "@/components/card";`
	if got := embeddedAdminFileContent(web, plain); got != plain {
		t.Errorf("a file outside the panel was rewritten: %s", got)
	}
}

// The shared schemas and types go where the project keeps them. A single has no
// workspace, so writing them to packages/shared writes them nowhere, and the
// generated customisation overlay imports the type it never got.
func TestSharedRootFollowsTheProject(t *testing.T) {
	monorepo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(monorepo, "packages", "shared"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got := sharedRootFrom(monorepo); got != filepath.Join(monorepo, "packages", "shared") {
		t.Errorf("a monorepo's shared root is %s", got)
	}

	single := t.TempDir()
	if err := os.MkdirAll(filepath.Join(single, "frontend", "src", "shared"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got := sharedRootFrom(single); got != filepath.Join(single, "frontend", "src", "shared") {
		t.Errorf("a single's shared root is %s", got)
	}
}

// The commands that take a root and no architecture still have to find the panel.
func TestAdminRootFromDisk(t *testing.T) {
	for _, dir := range []string{
		filepath.Join("apps", "admin"),
		filepath.Join("apps", "web", "admin-panel"),
		filepath.Join("frontend", "src", "admin-panel"),
	} {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
		if got := adminRootFrom(root); got != filepath.Join(root, dir) {
			t.Errorf("with %s on disk the panel was looked for at %s", dir, got)
		}
	}
}
