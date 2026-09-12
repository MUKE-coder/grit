package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fatih/color"
)

// Three kinds of page, three layouts. The root layout is the document and
// nothing else: when it drew the site's chrome, the admin panel rendered inside
// the marketing navbar because a client-side prefix list had never heard of it.
func TestWebSectionsHaveTheirOwnLayouts(t *testing.T) {
	root := t.TempDir()
	files := webFileMap(root, doubleOptions())

	for _, want := range []string{
		"apps/web/app/layout.tsx",
		"apps/web/app/(marketing)/layout.tsx",
		"apps/web/app/(marketing)/page.tsx",
		"apps/web/app/(marketing)/blog/page.tsx",
		"apps/web/app/(marketing)/blog/[slug]/page.tsx",
	} {
		if _, ok := files[filepath.Join(root, filepath.FromSlash(want))]; !ok {
			t.Errorf("%s is missing from the web app", want)
		}
	}

	// The old shape: pages at the app root, and a component that took the chrome
	// away again by matching path prefixes.
	for _, gone := range []string{
		"apps/web/app/page.tsx",
		"apps/web/app/blog/page.tsx",
		"apps/web/components/AppChrome.tsx",
	} {
		if _, ok := files[filepath.Join(root, filepath.FromSlash(gone))]; ok {
			t.Errorf("%s is still written: the section layouts replaced it", gone)
		}
	}

	rootLayout := files[filepath.Join(root, filepath.FromSlash("apps/web/app/layout.tsx"))]
	if strings.Contains(rootLayout, "AppChrome") {
		t.Error("the root layout still wraps every page in the site chrome")
	}
	if !strings.Contains(rootLayout, "<Providers>{children}</Providers>") {
		t.Error("the root layout should render the providers and nothing else")
	}

	marketing := files[filepath.Join(root, filepath.FromSlash("apps/web/app/(marketing)/layout.tsx"))]
	if !strings.Contains(marketing, "<Navbar />") || !strings.Contains(marketing, "<Footer />") {
		t.Error("the marketing layout is what draws the navbar and footer now")
	}

	// The public form share is outside every group, so it is bare by construction
	// rather than by being on a list.
	if _, ok := files[filepath.Join(root, filepath.FromSlash("apps/web/app/forms/[token]/page.tsx"))]; !ok {
		t.Error("the public form page moved: check it is still outside (marketing)")
	}
}

// The 404 is rendered from the app root, outside every group, so it has no
// section layout to inherit chrome from.
func TestWebNotFoundBringsItsOwnChrome(t *testing.T) {
	page := webNotFoundPage()
	if !strings.Contains(page, "<Navbar />") || !strings.Contains(page, "<Footer />") {
		t.Error("the 404 page renders with no header or footer at all")
	}
}

// NEXT_PUBLIC_ADMIN_URL is a URL only where the panel is its own application.
// Setting it to localhost:3001 everywhere is what kept a double's navbar
// pointing at a port with nothing behind it, even after the default was fixed.
func TestAdminURLEnvFollowsTheArchitecture(t *testing.T) {
	triple := envFile(tripleOptions())
	if !strings.Contains(triple, "\nNEXT_PUBLIC_ADMIN_URL=http://localhost:3001") {
		t.Error("a triple's admin app runs on its own port and .env should say so")
	}

	double := envFile(doubleOptions())
	if strings.Contains(double, "\nNEXT_PUBLIC_ADMIN_URL=http://localhost:3001") {
		t.Error("a double has no admin app on :3001, and this overrides the right default")
	}
	if !strings.Contains(double, "# NEXT_PUBLIC_ADMIN_URL=") {
		t.Error("the override should still be documented, commented out")
	}

	single := envFile(singleOptions())
	if strings.Contains(single, "NEXT_PUBLIC_ADMIN_URL") {
		t.Error("a single project has no Next.js app to read this")
	}
}

// The customer area: a route group with its own layout, and the two pages the
// user menu and the middleware have been pointing at with nothing behind them.
func TestWebAuthWritesTheCustomerArea(t *testing.T) {
	root := t.TempDir()
	webRoot := filepath.Join(root, "apps", "web")
	byPath := map[string]string{}
	for _, f := range webAuthFiles(webRoot, doubleOptions()) {
		rel, err := filepath.Rel(root, f.path)
		if err != nil {
			t.Fatal(err)
		}
		byPath[filepath.ToSlash(rel)] = f.content
	}

	for _, want := range []string{
		"apps/web/app/(app)/layout.tsx",
		"apps/web/app/(app)/account/page.tsx",
		"apps/web/app/(app)/account/profile/page.tsx",
	} {
		if _, ok := byPath[want]; !ok {
			t.Errorf("%s is missing: signing in and opening your account is still a 404", want)
		}
	}
	if _, ok := byPath["apps/web/components/AppChrome.tsx"]; ok {
		t.Error("web-auth should not write AppChrome any more")
	}

	layout := byPath["apps/web/app/(app)/layout.tsx"]
	if !strings.Contains(layout, "ProtectedWebRoute") {
		t.Error("the customer area must be behind the guard")
	}
	if strings.Contains(layout, "<Navbar />") {
		t.Error("the customer area draws its own chrome, not the marketing navbar")
	}
	// The web app's palette has bg-secondary and text-secondary, not the admin's
	// muted/card/primary tokens. An unknown utility compiles to nothing, which is
	// how a page ends up unstyled with no error anywhere.
	for _, absent := range []string{"text-muted-foreground", "bg-card", "bg-primary ", "text-primary-foreground"} {
		for name, content := range byPath {
			if strings.Contains(name, "(app)") && strings.Contains(content, absent) {
				t.Errorf("%s uses %q, which is not a token this app defines", name, absent)
			}
		}
	}
}

// Upgrade moves the framework's own pages into the group and leaves everything
// else alone, because only the project knows which section its pages belong to.
func TestMigrateWebRouteGroupsMovesOnlyWhatItOwns(t *testing.T) {
	root := t.TempDir()
	app := filepath.Join(root, "apps", "web", "app")
	write := func(rel, body string) {
		path := filepath.Join(app, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("page.tsx", "// the landing page, with a local edit\n")
	write("blog/page.tsx", "// blog list\n")
	write("pricing/page.tsx", "// the project's own page\n")
	chrome := filepath.Join(root, "apps", "web", "components", "AppChrome.tsx")
	if err := os.MkdirAll(filepath.Dir(chrome), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(chrome, []byte("// prefix list\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	quiet := color.New()
	migrateWebRouteGroups(root, quiet, quiet, quiet)

	moved, err := os.ReadFile(filepath.Join(app, "(marketing)", "page.tsx"))
	if err != nil {
		t.Fatalf("the landing page was not moved into the group: %v", err)
	}
	if !strings.Contains(string(moved), "with a local edit") {
		t.Error("the page was rewritten rather than moved, losing the project's edits")
	}
	if _, err := os.Stat(filepath.Join(app, "page.tsx")); err == nil {
		t.Error("the old landing page is still there: two pages now resolve to /")
	}
	if _, err := os.Stat(filepath.Join(app, "(marketing)", "blog", "page.tsx")); err != nil {
		t.Error("the blog did not move with it")
	}
	if _, err := os.Stat(filepath.Join(app, "pricing", "page.tsx")); err != nil {
		t.Error("a page the project wrote itself was moved: only its author knows where it belongs")
	}
	if _, err := os.Stat(chrome); err == nil {
		t.Error("AppChrome is unused after the migration and should be removed")
	}

	// Running it twice must not move anything a second time.
	migrateWebRouteGroups(root, quiet, quiet, quiet)
	if _, err := os.Stat(filepath.Join(app, "(marketing)", "(marketing)")); err == nil {
		t.Error("a second upgrade nested the group inside itself")
	}
}

// The line .env carried is retired for a project that already exists, and only
// when it is the exact line the scaffold wrote.
func TestRetireAdminURLEnvLeavesAChoiceAlone(t *testing.T) {
	quiet := color.New()

	scaffolded := t.TempDir()
	path := filepath.Join(scaffolded, ".env")
	if err := os.WriteFile(path, []byte("A=1\nNEXT_PUBLIC_ADMIN_URL=http://localhost:3001\nB=2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	retireAdminURLEnv(scaffolded, quiet)
	out, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "\nNEXT_PUBLIC_ADMIN_URL=http://localhost:3001") {
		t.Error("the stale line survived, so the navbar still points at :3001")
	}
	if !strings.Contains(string(out), "A=1") || !strings.Contains(string(out), "B=2") {
		t.Error("the rest of .env was disturbed")
	}

	chosen := t.TempDir()
	path = filepath.Join(chosen, ".env")
	const own = "NEXT_PUBLIC_ADMIN_URL=https://admin.example.com\n"
	if err := os.WriteFile(path, []byte(own), 0o644); err != nil {
		t.Fatal(err)
	}
	retireAdminURLEnv(chosen, quiet)
	out, err = os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != own {
		t.Error("a value somebody chose was rewritten")
	}
}
