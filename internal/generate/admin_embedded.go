package generate

import (
	"path/filepath"
	"regexp"
	"strings"
)

// Where a generated resource's admin screens go.
//
// A triple project has an admin app, and they go in it. A double has no admin
// app: the panel is a route group inside the web app, so the definition belongs
// under apps/web/admin-panel/resources and the pages under
// apps/web/app/admin/(dashboard)/resources. Writing them to apps/admin in a double
// produces files nothing compiles and a resource that exists in the API, in the
// types and in the hooks, and nowhere in the admin.
//
// The content needs the same treatment as the scaffolded screens: @/ is the web
// app's own code inside the web app, and an absolute link to /resources/invoices
// is /admin/resources/invoices there. See internal/scaffold/admin_embedded.go,
// which does this for the panel the scaffold writes.

// EmbedsAdmin reports whether this project's admin panel lives inside the web app.
func (g *Generator) EmbedsAdmin() bool {
	return g.Architecture == "double"
}

// EmbedsAdminInSPA reports whether the panel lives inside the single binary's SPA.
//
// The third shape. Its panel is a TanStack Router app like the --vite admin, so
// the screens are the same; what differs is where they go and that the router is
// the SPA's, so every route id carries the /admin segment.
func (g *Generator) EmbedsAdminInSPA() bool {
	return g.Architecture == "single"
}

// AdminIsTanStack reports whether the panel is a TanStack Router app, which
// decides whether a screen is one page file or a page plus a route shim.
func (g *Generator) AdminIsTanStack() bool {
	if g.EmbedsAdminInSPA() {
		return true
	}
	return dirExists(filepath.Join(g.Root, "apps", "admin", "src", "routes"))
}

// adminTanStackRoot is where a TanStack panel keeps its own code.
func (g *Generator) adminTanStackRoot() string {
	if g.EmbedsAdminInSPA() {
		return filepath.Join(g.Root, "frontend", "src", "admin-panel")
	}
	return filepath.Join(g.Root, "apps", "admin", "src")
}

// adminTanStackRoutesRoot is the directory its route files live under. In the SPA
// that is the SPA's own route tree, under the segment the panel owns.
func (g *Generator) adminTanStackRoutesRoot() string {
	if g.EmbedsAdminInSPA() {
		return filepath.Join(g.Root, "frontend", "src", "routes", "admin")
	}
	return filepath.Join(g.Root, "apps", "admin", "src", "routes")
}

// AdminCodeRoot is where the admin's components, lib, hooks and resource
// definitions live.
func (g *Generator) AdminCodeRoot() string {
	switch {
	case g.EmbedsAdmin():
		return filepath.Join(g.Root, "apps", "web", "admin-panel")
	case g.EmbedsAdminInSPA():
		return filepath.Join(g.Root, "frontend", "src", "admin-panel")
	default:
		return filepath.Join(g.Root, "apps", "admin")
	}
}

// AdminRoutesRoot is the directory the admin's routes live under.
func (g *Generator) AdminRoutesRoot() string {
	if g.EmbedsAdmin() {
		return filepath.Join(g.Root, "apps", "web", "app", "admin")
	}
	return filepath.Join(g.Root, "apps", "admin", "app")
}

// embeddedAdminFileContent repoints a generated admin file, decided by where it
// is being written rather than by who is writing it.
//
// Anything under the web app's admin-panel directory or its app/admin route group
// is the panel's code. Everything else is returned untouched, so this is inert for
// a triple project and for every file that is not an admin screen.
func embeddedAdminFileContent(path, content string) string {
	slashed := filepath.ToSlash(path)
	if strings.Contains(slashed, "/frontend/src/admin-panel/") ||
		strings.Contains(slashed, "/frontend/src/routes/admin/") {
		return repointSPAAdminContent(content)
	}
	if !strings.Contains(slashed, "/apps/web/admin-panel/") &&
		!strings.Contains(slashed, "/apps/web/app/admin/") {
		return content
	}
	return repointAdminContent(content)
}

// repointSPAAdminContent is the same rewrite for a panel inside the SPA, plus
// the one thing only that shape needs: the id a file route declares has to equal
// its path, and its path now starts with /admin.
func repointSPAAdminContent(content string) string {
	content = repointAdminContent(content)
	content = strings.ReplaceAll(content, `"@/pages/`, `"@admin/pages/`)
	content = strings.ReplaceAll(content, `'@/pages/`, `'@admin/pages/`)
	content = routeIDPattern.ReplaceAllStringFunc(content, func(match string) string {
		id := routeIDPattern.FindStringSubmatch(match)[1]
		if strings.HasPrefix(id, "/admin") {
			return match
		}
		return "createFileRoute('/admin" + id + "')"
	})
	return content
}

// routeIDPattern finds the id a TanStack file route declares.
var routeIDPattern = regexp.MustCompile(`createFileRoute\('([^']*)'\)`)

// adminContent repoints a generated admin file for a project whose panel is
// embedded, and returns it unchanged for one whose admin is its own app.
func (g *Generator) adminContent(content string) string {
	if !g.EmbedsAdmin() {
		return content
	}
	return repointAdminContent(content)
}

// repointAdminContent rewrites the imports and links of one admin file.
func repointAdminContent(content string) string {
	for _, dir := range []string{"components", "lib", "hooks", "resources", "config"} {
		content = strings.ReplaceAll(content, `"@/`+dir+`/`, `"@admin/`+dir+`/`)
		content = strings.ReplaceAll(content, `'@/`+dir+`/`, `'@admin/`+dir+`/`)
		content = strings.ReplaceAll(content, `"@/`+dir+`"`, `"@admin/`+dir+`"`)
	}
	// The links a generated screen carries: a resource page links to its own list
	// and to other resources.
	for _, segment := range []string{"resources", "dashboard", "profile", "system", "account", "settings"} {
		for _, quote := range []string{`"`, "`", `'`} {
			content = strings.ReplaceAll(content, quote+"/"+segment, quote+"/admin/"+segment)
		}
	}
	return strings.ReplaceAll(content, "/admin/admin/", "/admin/")
}

// adminRootFrom finds a project's admin panel with nothing but its root.
//
// The commands that take a root and no Options, like sync and remove, cannot ask
// which shape the project is. Reading the layout off the disk answers it: exactly
// one of these three directories exists in any project that has a panel.
func adminRootFrom(root string) string {
	for _, candidate := range []string{
		filepath.Join(root, "apps", "admin"),
		filepath.Join(root, "apps", "web", "admin-panel"),
		filepath.Join(root, "frontend", "src", "admin-panel"),
	} {
		if dirExists(candidate) {
			return candidate
		}
	}
	return filepath.Join(root, "apps", "admin")
}

// relativeToRoot renders a path for the screen, relative to the project root.
func relativeToRoot(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(rel)
}

// writeAdminFile writes a generated admin file, repointed if the panel is
// embedded.
func (g *Generator) writeAdminFile(path, content string) error {
	return writeFileWithDirs(path, g.adminContent(content))
}
