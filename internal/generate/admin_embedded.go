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
	if g.Architecture == "single" {
		return true
	}
	// A double whose web app was scaffolded with --vite hosts the panel the same
	// way, in src/admin-panel rather than in a Next.js route group. Read off the
	// disk rather than from the architecture, which does not record the frontend.
	return dirExists(filepath.Join(g.Root, "apps", "web", "src", "admin-panel"))
}

// spaHost is the Vite app hosting the panel.
func (g *Generator) spaHost() string {
	if g.Architecture == "single" {
		// The frontend owns the project root, with the Go module in api/ beside
		// it. A project scaffolded before that keeps its app in frontend/, and
		// `grit generate resource` has to keep working in it, so the directory on
		// disk decides rather than the version that wrote it: a package.json in
		// frontend/ is the old shape, and anything else is the new one.
		if legacy := filepath.Join(g.Root, "frontend"); fileExists(filepath.Join(legacy, "package.json")) {
			return legacy
		}
		return g.Root
	}
	return filepath.Join(g.Root, "apps", "web")
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
		return filepath.Join(g.spaHost(), "src", "admin-panel")
	}
	return filepath.Join(g.Root, "apps", "admin", "src")
}

// adminTanStackRoutesRoot is the directory its route files live under. In the SPA
// that is the SPA's own route tree, under the segment the panel owns, which is
// no segment at all when the panel is the whole app.
func (g *Generator) adminTanStackRoutesRoot() string {
	if g.EmbedsAdminInSPA() {
		routes := filepath.Join(g.spaHost(), "src", "routes")
		probe := filepath.ToSlash(filepath.Join(g.spaHost(), "src", "admin-panel")) + "/"
		if adminServedAtRoot(probe) {
			return routes
		}
		return filepath.Join(routes, "admin")
	}
	return filepath.Join(g.Root, "apps", "admin", "src", "routes")
}

// adminServedAtRoot answers, from the destination path alone, whether the panel
// a file belongs to is served at the SPA's root rather than under /admin
// (grit#89).
//
// `grit new --api --admin` builds an app that is nothing but the panel, so its
// screens are at /dashboard rather than /admin/dashboard, and every prefix this
// file adds would point at a path the project does not serve. The file's own
// path says which SPA it is going into, and that SPA's route tree says which
// shape it is: routes/admin is the segment a panel beside a public site owns.
//
// Read from disk rather than carried in a flag because these rewrites are keyed
// on the destination and run from writeFileWithDirs, which has no Generator.
// Positive evidence either way: with neither directory present the answer is
// the old one, so a layout nobody has seen cannot silently lose its prefixes.
func adminServedAtRoot(slashed string) bool {
	var spa string
	switch {
	case strings.Contains(slashed, "/src/admin-panel/"):
		spa = slashed[:strings.Index(slashed, "/src/admin-panel/")]
	case strings.Contains(slashed, "/src/routes/"):
		spa = slashed[:strings.Index(slashed, "/src/routes/")]
	default:
		return false
	}
	host := filepath.FromSlash(spa)
	// Embedded at all. A standalone admin app has routes/_dashboard too, and
	// it owns @/: rewriting its imports points them at an @admin alias it does
	// not declare. What makes a panel embedded is that its code sits in a host
	// SPA's src/admin-panel, so that is the question asked first.
	if !dirExists(filepath.Join(host, "src", "admin-panel")) {
		return false
	}
	routes := filepath.Join(host, "src", "routes")
	if dirExists(filepath.Join(routes, "admin")) {
		return false
	}
	// _dashboard is the panel's in both shapes: a public site lives under
	// _site. At the SPA's root it is the panel saying it is the whole app.
	return dirExists(filepath.Join(routes, "_dashboard"))
}

// AdminCodeRoot is where the admin's components, lib, hooks and resource
// definitions live.
func (g *Generator) AdminCodeRoot() string {
	switch {
	// The Vite case first: a double built with --vite satisfies both, and its
	// panel is a section of a Vite app rather than a Next.js route group. Asking
	// in the other order wrote a resource's definition to a directory that does
	// not exist, so the sidebar never learned about it.
	case g.EmbedsAdminInSPA():
		return filepath.Join(g.spaHost(), "src", "admin-panel")
	case g.EmbedsAdmin():
		return filepath.Join(g.Root, "apps", "web", "admin-panel")
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
	// A Vite host: frontend/src for a single, apps/web/src for a double built
	// with --vite. The Next.js panel has no src segment, so the two transforms
	// cannot both claim a file.
	if strings.Contains(slashed, "/src/admin-panel/") ||
		strings.Contains(slashed, "/src/routes/admin/") {
		return repointSPAAdminContent(content, !adminServedAtRoot(slashed))
	}
	// The panel's route files when it is served at the SPA's root: there is no
	// /admin segment to key on, so the host SPA beside them is asked instead.
	// Nothing is prefixed here, because the panel is the whole app.
	if adminServedAtRoot(slashed) {
		return repointSPAAdminContent(content, false)
	}
	if !strings.Contains(slashed, "/apps/web/admin-panel/") &&
		!strings.Contains(slashed, "/apps/web/app/admin/") {
		return content
	}
	return repointAdminContent(content)
}

// repointSPAAdminContent is the same rewrite for a panel inside the SPA, plus
// the one thing only that shape needs: the id a file route declares has to equal
// its path, and its path starts with /admin.
//
// underAdmin is false when the panel is the whole app. The aliases still have
// to be repointed, because the panel's code sits behind @admin wherever it is
// served, but nothing is prefixed: its own links and route ids are already the
// addresses it answers on.
func repointSPAAdminContent(content string, underAdmin bool) string {
	content = repointAdminAliases(content)
	content = strings.ReplaceAll(content, `"@/pages/`, `"@admin/pages/`)
	content = strings.ReplaceAll(content, `'@/pages/`, `'@admin/pages/`)
	if !underAdmin {
		return content
	}
	content = prefixAdminLinks(content)
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
	return prefixAdminLinks(repointAdminAliases(content))
}

// repointAdminAliases puts the panel's own code behind @admin. True wherever
// the panel is served: these are imports, not addresses.
func repointAdminAliases(content string) string {
	for _, dir := range []string{"components", "lib", "hooks", "resources", "config"} {
		content = strings.ReplaceAll(content, `"@/`+dir+`/`, `"@admin/`+dir+`/`)
		content = strings.ReplaceAll(content, `'@/`+dir+`/`, `'@admin/`+dir+`/`)
		content = strings.ReplaceAll(content, `"@/`+dir+`"`, `"@admin/`+dir+`"`)
	}
	return content
}

// prefixAdminLinks moves the links a generated screen carries under /admin: a
// resource page links to its own list and to other resources, and inside a host
// app those pages are one segment deeper than the panel thinks.
func prefixAdminLinks(content string) string {
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
		filepath.Join(root, "apps", "web", "src", "admin-panel"),
		// --single: the frontend owns the project root, so the panel is a
		// section of the SPA at src/admin-panel.
		//
		// This was missing, and it is the layout `grit new --single` has
		// written since the SPA moved out of frontend/. Three commands read
		// this helper, so on every single project `grit remove resource` left
		// the admin definition and its page behind, still registered and
		// pointing at an endpoint that no longer existed, and `grit sync`
		// silently never added a new field to the admin table or form.
		// AdminCodeRoot, which `grit generate` uses, had the case right, so
		// generation worked and nothing else did.
		filepath.Join(root, "src", "admin-panel"),
		// The older single shape, with the app in frontend/.
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
