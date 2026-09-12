package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/fatih/color"
)

// The SPA has four sections and the root is not one of them: it renders an
// outlet and nothing else, so no section can inherit another's chrome.
func TestSingleRootRouteDrawsNothing(t *testing.T) {
	root := webTanStackRootRoute(singleOptions())
	if strings.Contains(root, "<Navbar />") || strings.Contains(root, "useRouterState") {
		t.Error("the root route is deciding chrome again; the sections own their layouts")
	}
	if !strings.Contains(root, "component: () => <Outlet />") {
		t.Error("the root route should render an outlet")
	}

	// A monorepo web app with no panel inside it keeps the old root, chrome and all.
	web := webTanStackRootRoute(Options{ProjectName: "app", Architecture: ArchTriple, Frontend: FrontendTanStack})
	if !strings.Contains(web, "<Navbar />") {
		t.Error("a web app with no embedded panel still draws its chrome from the root")
	}
}

// The public pages live under a pathless layout route, which is what gives them
// the navbar and the footer. The URLs do not change.
func TestSinglePublicPagesAreASection(t *testing.T) {
	root := t.TempDir()
	if err := writeSingleFrontendFiles(root, singleOptions()); err != nil {
		t.Fatal(err)
	}

	layout := filepath.Join(root, "frontend", "src", "routes", "_site.tsx")
	body, err := os.ReadFile(layout)
	if err != nil {
		t.Fatalf("the public section has no layout: %v", err)
	}
	if !strings.Contains(string(body), "<Navbar />") || !strings.Contains(string(body), "<Footer />") {
		t.Error("_site.tsx is what draws the site chrome now")
	}

	routeID := regexp.MustCompile(`createFileRoute\('([^']*)'\)`)
	for rel, want := range map[string]string{
		"_site/index.tsx":      "/_site/",
		"_site/blog/index.tsx": "/_site/blog/",
		"_site/blog/$slug.tsx": "/_site/blog/$slug",
	} {
		path := filepath.Join(root, "frontend", "src", "routes", filepath.FromSlash(rel))
		content, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s is missing: %v", rel, err)
			continue
		}
		match := routeID.FindStringSubmatch(string(content))
		if match == nil {
			t.Errorf("%s declares no route id", rel)
			continue
		}
		if match[1] != want {
			t.Errorf("%s declares %q, want %q: the router refuses an id that is not its path", rel, match[1], want)
		}
	}

	// The old flat pages must be gone, or two routes resolve to "/".
	for _, gone := range []string{"index.tsx", "blog/index.tsx"} {
		if _, err := os.Stat(filepath.Join(root, "frontend", "src", "routes", filepath.FromSlash(gone))); err == nil {
			t.Errorf("routes/%s is still written alongside its copy in _site", gone)
		}
	}
}

// grit add web-auth on a single writes the screens the SPA's auth library has
// never had, in the dialect that app speaks.
func TestSingleWebAuthWritesTheScreens(t *testing.T) {
	root := t.TempDir()
	feRoot := filepath.Join(root, "frontend")
	byPath := map[string]string{}
	for _, f := range singleAuthFiles(feRoot, singleOptions()) {
		rel, err := filepath.Rel(root, f.path)
		if err != nil {
			t.Fatal(err)
		}
		byPath[filepath.ToSlash(rel)] = f.content
	}

	for _, want := range []string{
		"frontend/src/hooks/use-auth.ts",
		"frontend/src/components/user-menu.tsx",
		"frontend/src/routes/_auth.tsx",
		"frontend/src/routes/_auth/login.tsx",
		"frontend/src/routes/_auth/register.tsx",
		"frontend/src/routes/_auth/forgot-password.tsx",
		"frontend/src/routes/_auth/reset-password.tsx",
		"frontend/src/routes/account/route.tsx",
		"frontend/src/routes/account/index.tsx",
		"frontend/src/routes/account/profile.tsx",
	} {
		if _, ok := byPath[want]; !ok {
			t.Errorf("%s is missing", want)
		}
	}

	// Every route file's id has to equal its path.
	routeID := regexp.MustCompile(`createFileRoute\('([^']*)'\)`)
	for rel, want := range map[string]string{
		"frontend/src/routes/_auth.tsx":           "/_auth",
		"frontend/src/routes/_auth/login.tsx":     "/_auth/login",
		"frontend/src/routes/_auth/register.tsx":  "/_auth/register",
		"frontend/src/routes/account/route.tsx":   "/account",
		"frontend/src/routes/account/index.tsx":   "/account/",
		"frontend/src/routes/account/profile.tsx": "/account/profile",
	} {
		match := routeID.FindStringSubmatch(byPath[rel])
		if match == nil {
			t.Errorf("%s declares no route id", rel)
			continue
		}
		if match[1] != want {
			t.Errorf("%s declares %q, want %q", rel, match[1], want)
		}
	}

	// The customer area is guarded before a screen renders, not after.
	account := byPath["frontend/src/routes/account/route.tsx"]
	if !strings.Contains(account, "beforeLoad") || !strings.Contains(account, "redirect({ to: '/login' })") {
		t.Error("the account section must turn away a visitor with no session")
	}
	if strings.Contains(account, "<Navbar />") {
		t.Error("the customer area draws its own chrome")
	}

	// The SPA's palette, not the admin's: an unknown utility renders nothing.
	for rel, content := range byPath {
		for _, absent := range []string{"text-muted-foreground", "bg-card", "text-primary-foreground"} {
			if strings.Contains(content, absent) {
				t.Errorf("%s uses %q, which this app does not define", rel, absent)
			}
		}
	}

	// The navbar gains the user menu and keeps everything else.
	navbar := byPath["frontend/src/components/navbar.tsx"]
	if !strings.Contains(navbar, "<UserMenu />") {
		t.Error("the navbar has no way into the account area")
	}
	if !strings.Contains(navbar, `{ href: "/blog", label: "Blog" }`) {
		t.Error("the auth-aware navbar lost the site's own links")
	}
}

// Upgrade moves an existing SPA's public pages into the section and corrects the
// id each one declares, or the router refuses them.
func TestMigrateSPARouteSections(t *testing.T) {
	root := t.TempDir()
	routes := filepath.Join(root, "frontend", "src", "routes")
	write := func(rel, body string) {
		path := filepath.Join(routes, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("index.tsx", "// the landing page, edited locally\nexport const Route = createFileRoute('/')({})\n")
	write("blog/index.tsx", "export const Route = createFileRoute('/blog/')({})\n")
	write("blog/$slug.tsx", "export const Route = createFileRoute('/blog/$slug')({})\n")

	quiet := color.New()
	migrateSPARouteSections(root, quiet, quiet)

	landing, err := os.ReadFile(filepath.Join(routes, "_site", "index.tsx"))
	if err != nil {
		t.Fatalf("the landing page did not move: %v", err)
	}
	if !strings.Contains(string(landing), "edited locally") {
		t.Error("the page was rewritten rather than moved")
	}
	if !strings.Contains(string(landing), "createFileRoute('/_site/')") {
		t.Errorf("the moved page still declares its old id:\n%s", landing)
	}
	slug, err := os.ReadFile(filepath.Join(routes, "_site", "blog", "$slug.tsx"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(slug), "createFileRoute('/_site/blog/$slug')") {
		t.Errorf("a nested page kept its old id:\n%s", slug)
	}
	if _, err := os.Stat(filepath.Join(routes, "index.tsx")); err == nil {
		t.Error("the old landing page is still there: two routes now resolve to /")
	}

	// A second upgrade is a no-op.
	migrateSPARouteSections(root, quiet, quiet)
	if _, err := os.Stat(filepath.Join(routes, "_site", "_site")); err == nil {
		t.Error("the section was nested inside itself")
	}
}

// siteRouteID is the rewrite the move depends on.
func TestSiteRouteID(t *testing.T) {
	for in, want := range map[string]string{
		"createFileRoute('/')":           "createFileRoute('/_site/')",
		"createFileRoute('/blog/')":      "createFileRoute('/_site/blog/')",
		"createFileRoute('/blog/$slug')": "createFileRoute('/_site/blog/$slug')",
		"createFileRoute('/_site/')":     "createFileRoute('/_site/')",
	} {
		if got := siteRouteID(in); got != want {
			t.Errorf("siteRouteID(%q) = %q, want %q", in, got, want)
		}
	}
}
