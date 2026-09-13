package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func viteDoubleOptions() Options {
	return Options{ProjectName: "app", Architecture: ArchDouble, Frontend: FrontendTanStack, Theme: "atlas"}
}

func viteTripleOptions() Options {
	return Options{ProjectName: "app", Architecture: ArchTriple, Frontend: FrontendTanStack, Theme: "atlas"}
}

// Two hosts, one transform. A double built with --vite hosts the panel the same
// way a single does, because both are TanStack Router apps.
func TestViteDoubleHostsThePanelLikeASingle(t *testing.T) {
	vite := viteDoubleOptions()
	if !vite.ShouldEmbedAdminInSPA() {
		t.Fatal("a --double --vite project hosts the panel in its Vite app")
	}
	if !singleOptions().ShouldEmbedAdminInSPA() {
		t.Error("a single still qualifies")
	}
	// A double on Next.js must be untouched by any of this.
	if doubleOptions().ShouldEmbedAdminInSPA() {
		t.Error("a Next.js double takes the route-group transform, not this one")
	}
	if viteTripleOptions().ShouldEmbedAdminInSPA() {
		t.Error("a triple has an admin application; nothing is embedded")
	}

	root := t.TempDir()
	if got := spaHostRoot(root, vite); got != filepath.Join(root, "apps", "web") {
		t.Errorf("the host is %s", got)
	}
	if got := spaHostRoot(root, singleOptions()); got != filepath.Join(root, "frontend") {
		t.Errorf("a single's host is %s", got)
	}

	// Both predicates are true for this shape, so the order they are asked in is
	// what decides whether the panel lands in a directory the app routes.
	if got := adminCodeRoot(root, vite); got != filepath.Join(root, "apps", "web", "src", "admin-panel") {
		t.Errorf("the panel's code root is %s", got)
	}
	if got := adminPath(root, vite, "routes", "_dashboard", "x.tsx"); got != filepath.Join(root, "apps", "web", "src", "routes", "admin", "_dashboard", "x.tsx") {
		t.Errorf("a route file went to %s, where the app's router will not find it", got)
	}
}

// The files land inside the Vite app, and none of them in an app/ directory it
// does not route.
func TestViteDoublePanelFileMap(t *testing.T) {
	root := t.TempDir()
	files := embeddedSingleAdminFileMap(root, viteDoubleOptions())
	if len(files) < 100 {
		t.Fatalf("only %d files", len(files))
	}

	for path := range files {
		rel, err := filepath.Rel(root, path)
		if err != nil {
			t.Fatal(err)
		}
		rel = filepath.ToSlash(rel)
		if !strings.HasPrefix(rel, "apps/web/src/") {
			t.Errorf("%s is outside the Vite app's source", rel)
		}
		if strings.HasPrefix(rel, "apps/web/app/") {
			t.Errorf("%s is a Next.js route group: this app has no app directory", rel)
		}
	}

	for _, want := range []string{
		"apps/web/src/routes/admin/route.tsx",
		"apps/web/src/routes/admin/_dashboard/dashboard.tsx",
		"apps/web/src/admin-panel/pages/dashboard.tsx",
		"apps/web/src/admin-panel/lib/api-client.ts",
	} {
		if _, ok := files[filepath.Join(root, filepath.FromSlash(want))]; !ok {
			t.Errorf("%s is missing", want)
		}
	}

	// And the rewrite still applies: the destination is what decides, and the
	// destination moved.
	client := filepath.Join(root, filepath.FromSlash("apps/web/src/admin-panel/lib/api-client.ts"))
	final := embeddedSingleAdminContent(client, files[client])
	if strings.Contains(final, "localhost:8080") {
		t.Error("the panel should talk to the origin serving it")
	}
	if leftoverAliasPattern.MatchString(final) {
		t.Error("the panel's imports were not repointed at @admin")
	}
}

// The Vite app needs telling about the panel in two files, because Vite does not
// read tsconfig paths and tsc does not read Vite aliases.
func TestViteHostResolvesThePanel(t *testing.T) {
	cfg := webTanStackViteConfig(viteDoubleOptions())
	for _, want := range []string{
		"'@admin': path.resolve(__dirname, './src/admin-panel')",
		"'@repo/upload/web': path.resolve(__dirname, '../../packages/upload/src/web.ts')",
	} {
		if !strings.Contains(cfg, want) {
			t.Errorf("the vite config is missing %s", want)
		}
	}
	ts := webTanStackTSConfig(viteDoubleOptions())
	if !strings.Contains(ts, `"@admin/*": ["./src/admin-panel/*"]`) {
		t.Error("the tsconfig is missing the @admin path")
	}

	// A Vite app with no panel inside it gets neither.
	if strings.Contains(webTanStackViteConfig(viteTripleOptions()), "@admin") {
		t.Error("a triple's web app has no panel and should not alias one")
	}
	if strings.Contains(webTanStackTSConfig(viteTripleOptions()), "@admin") {
		t.Error("same for its tsconfig")
	}

	// The panel's dependencies, including the two the Next.js web app carried in
	// its own list and this one did not.
	pkg := webTanStackPackageJSON(viteDoubleOptions())
	for _, dep := range []string{"@tiptap/react", "recharts", "xlsx", "react-hook-form", "@hookform/resolvers"} {
		if !strings.Contains(pkg, `"`+dep+`"`) {
			t.Errorf("the web app does not depend on %s, which the panel imports", dep)
		}
	}
	if strings.Contains(webTanStackPackageJSON(viteTripleOptions()), "@tiptap/react") {
		t.Error("a triple's web app should not carry the panel's dependencies")
	}
}

// The admin link, which is a path when the panel is inside this app. Every Vite
// web app shipped the literal placeholder before this.
func TestViteWebAppGetsARealAdminLink(t *testing.T) {
	root := t.TempDir()
	if err := writeWebTanStackFiles(root, viteDoubleOptions()); err != nil {
		t.Fatal(err)
	}
	navbar, err := os.ReadFile(filepath.Join(root, "apps", "web", "src", "components", "navbar.tsx"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(navbar), "{{ADMIN_HREF}}") {
		t.Error("the navbar links to a page called {{ADMIN_HREF}}")
	}
	if !strings.Contains(string(navbar), "/admin/dashboard") {
		t.Error("a double's navbar should link to the panel inside it")
	}
}

// The SPA's API client is a bearer client, and still has to carry the CSRF token:
// the admin panel in the same app signs in with cookies, and after that every
// mutation the public site makes arrives with a session cookie and no token.
func TestSPAClientEchoesTheCSRFToken(t *testing.T) {
	client := viteAPIClientWithAuth()
	if !strings.Contains(client, "X-CSRF-Token") || !strings.Contains(client, "grit_csrf") {
		t.Error("a sign-up in a browser that has used the admin panel is refused with CSRF_INVALID")
	}
}

// A Vite app keeps its code under src/. Writing to apps/web/hooks created a
// directory nothing imports, which the resource generator then followed.
func TestRealtimeClientFollowsTheAppLayout(t *testing.T) {
	root := t.TempDir()
	if err := writeRealtimeClientFiles(root, viteDoubleOptions()); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "apps", "web", "src", "hooks", "use-realtime.ts")); err != nil {
		t.Errorf("the hook is not where a Vite app keeps hooks: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "apps", "web", "hooks")); err == nil {
		t.Error("apps/web/hooks was created in a Vite app")
	}
	// The panel's copy goes inside the panel, not into a src/ of its own.
	if _, err := os.Stat(filepath.Join(root, "apps", "web", "src", "admin-panel", "hooks", "use-realtime.ts")); err != nil {
		t.Errorf("the panel's copy is missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "apps", "web", "src", "admin-panel", "src")); err == nil {
		t.Error("the panel grew a second src/ directory")
	}
}

// grit add web-auth writes the customer area into whichever Vite app it is.
func TestWebAuthTargetsTheViteApp(t *testing.T) {
	root := t.TempDir()
	webRoot := filepath.Join(root, "apps", "web")
	paths := map[string]bool{}
	for _, f := range singleAuthFiles(webRoot, viteDoubleOptions()) {
		rel, err := filepath.Rel(root, f.path)
		if err != nil {
			t.Fatal(err)
		}
		paths[filepath.ToSlash(rel)] = true
	}
	for _, want := range []string{
		"apps/web/src/routes/_auth/login.tsx",
		"apps/web/src/routes/account/route.tsx",
		// The library and a client that carries its token: a Vite web app shipped
		// neither, so the screens had nothing to call.
		"apps/web/src/lib/auth.ts",
		"apps/web/src/lib/api.ts",
	} {
		if !paths[want] {
			t.Errorf("%s is missing", want)
		}
	}

	// The navbar is built from the host's own, not from the single's.
	web := viteHostNavbar(viteDoubleOptions())
	if !strings.Contains(web, "ADMIN_URL") {
		t.Error("a double's auth navbar should be the web app's navbar, which links to the panel")
	}
	if !strings.Contains(singleViteNavbarWithAuth(viteDoubleOptions()), "<UserMenu />") {
		t.Error("the account menu was not injected into the web app's navbar")
	}
	if !strings.Contains(singleViteNavbarWithAuth(singleOptions()), "<UserMenu />") {
		t.Error("nor into the single's")
	}
}
