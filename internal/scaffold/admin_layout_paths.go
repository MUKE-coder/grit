package scaffold

import "path/filepath"

// The two admin layouts, and the one function that knows the difference.
//
// The Next admin keeps components, hooks and lib at the app root and its pages
// under app/(dashboard)/. The TanStack admin puts all of those under src/, and
// its router reads src/routes/, where each route is a shim pointing at a
// component in src/pages/.
//
// Every writer that hardcoded the Next shape produced files a TanStack project
// never imports. They are not a build error, which is the problem: the feature
// is simply absent, the app compiles, and nobody finds out until somebody looks
// for the page. That is how the TanStack admin shipped with no account-security
// screen, and how three plugins wrote themselves into a directory nothing reads.

// adminCodeRoot is where this project's admin code lives.
//
// apps/admin for a triple, apps/web/admin-panel for a double whose panel is a
// route group inside the web app. Writers that hardcoded the first one wrote into
// a directory nothing compiles, which is how the embedded panel's first build
// failed on a dashboard widget that had been written to apps/admin/components.
func adminCodeRoot(root string, opts Options) string {
	if opts.ShouldEmbedAdmin() {
		return filepath.Join(root, "apps", "web", "admin-panel")
	}
	return filepath.Join(root, "apps", "admin")
}

// adminComponent returns the path for a component file.
//
//	adminComponent(root, opts, "security", "passkeys.tsx")
func adminComponent(root string, opts Options, parts ...string) string {
	return adminPath(root, opts, append([]string{"components"}, parts...)...)
}

// adminHook returns the path for a hook file.
func adminHook(root string, opts Options, name string) string {
	return adminPath(root, opts, "hooks", name)
}

// adminLib returns the path for a lib file.
func adminLib(root string, opts Options, name string) string {
	return adminPath(root, opts, "lib", name)
}

// adminPath joins parts under the admin root, adding src/ for TanStack.
//
// Three layouts now, not two: a double project has no admin app, so the panel's
// components, hooks and lib live inside the web app under admin-panel/ and its
// routes under app/admin/. See admin_embedded.go.
func adminPath(root string, opts Options, parts ...string) string {
	base := filepath.Join(root, "apps", "admin")
	switch {
	case opts.ShouldEmbedAdmin():
		base = filepath.Join(root, "apps", "web", "admin-panel")
	case opts.UseTanStack():
		base = filepath.Join(base, "src")
	}
	return filepath.Join(base, filepath.Join(parts...))
}

// adminPageFiles returns the file(s) that put a dashboard page on screen.
//
// One file on Next, which is the page itself at its route path. Two on
// TanStack, because the router wants a route shim separate from the component:
// the page lands in src/pages/ and src/routes/ points at it.
//
//	adminPageFiles(root, opts, "account/security", "security", body)
//
// routePath is the URL under the dashboard, slash-separated. name is the file
// stem used under src/pages/ for TanStack.
func adminPageFiles(root string, opts Options, routePath, name, body string) map[string]string {
	if !opts.UseTanStack() {
		parts := append([]string{"app", "(dashboard)"}, splitSlash(routePath)...)
		parts = append(parts, "page.tsx")
		return map[string]string{
			adminPathNext(root, opts, parts...): body,
		}
	}

	pageFile := adminPath(root, opts, "pages", name+".tsx")
	routeFile := adminPath(root, opts, append([]string{"routes", "_dashboard"}, splitSlash(routePath+".tsx")...)...)
	return map[string]string{
		pageFile:  nextToTanStack(body),
		routeFile: adminTanStackPageRoute("/_dashboard/"+routePath, "@/pages/"+name),
	}
}

// adminPathNext joins parts under the Next admin root, ignoring the frontend
// choice but not the architecture.
//
// When the panel is embedded in the web app, app/(dashboard)/x/page.tsx belongs at
// apps/web/app/admin/(dashboard)/x/page.tsx: the route group keeps its name, so
// the URL is /admin/x, and the layout that draws the chrome still wraps it.
func adminPathNext(root string, opts Options, parts ...string) string {
	if opts.ShouldEmbedAdmin() && len(parts) > 0 && parts[0] == "app" {
		return filepath.Join(append([]string{root, "apps", "web", "app", "admin"}, parts[1:]...)...)
	}
	if opts.ShouldEmbedAdmin() {
		return filepath.Join(append([]string{root, "apps", "web", "admin-panel"}, parts...)...)
	}
	return filepath.Join(append([]string{root, "apps", "admin"}, parts...)...)
}

func splitSlash(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '/' {
			if i > start {
				out = append(out, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}
