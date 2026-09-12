package generate

import (
	"path/filepath"
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

// AdminCodeRoot is where the admin's components, lib, hooks and resource
// definitions live.
func (g *Generator) AdminCodeRoot() string {
	if g.EmbedsAdmin() {
		return filepath.Join(g.Root, "apps", "web", "admin-panel")
	}
	return filepath.Join(g.Root, "apps", "admin")
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
	if !strings.Contains(slashed, "/apps/web/admin-panel/") &&
		!strings.Contains(slashed, "/apps/web/app/admin/") {
		return content
	}
	return repointAdminContent(content)
}

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

// writeAdminFile writes a generated admin file, repointed if the panel is
// embedded.
func (g *Generator) writeAdminFile(path, content string) error {
	return writeFileWithDirs(path, g.adminContent(content))
}
