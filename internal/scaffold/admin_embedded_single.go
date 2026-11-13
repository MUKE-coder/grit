package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// The admin panel inside a single project's embedded SPA.
//
// A single project is one Go binary with a Vite SPA compiled into it. There is no
// admin app, and until now no admin panel either, so the same complaint applied as
// to a double: the shape offers a web UI and no way to administer anything.
//
// The source here is the Vite admin, the one a --triple --vite project gets, not
// the Next one. That is the whole reason this is tractable: a single project's
// frontend IS a TanStack Router app, so the panel's 203 files are already in the
// right dialect. They move into frontend/src, the routes gain an /admin segment,
// and the route ids inside the shims are rewritten to match, because TanStack
// checks that a file route's id is its path.
//
// What is dropped is the app shell: the SPA has its own index.html, vite config,
// main.tsx, root route and stylesheet.

// ShouldEmbedAdminInSPA reports whether the panel lives inside a Vite app.
//
// Two shapes qualify, and they are the same problem: a single project, whose one
// SPA is a TanStack Router app, and a double whose web app was scaffolded with
// --vite. Both host the panel as a section of themselves. A double on Next.js is
// the other transform, in admin_embedded.go.
//
// Before this, a --double --vite project got the Next.js panel: 41 files under an
// app/ directory a Vite app does not route, importing next/link, which it cannot
// resolve. The panel was unreachable and would not have built.
func (o Options) ShouldEmbedAdminInSPA() bool {
	// Not every single: a Next single embeds the panel as a route group like a
	// double does, because it has a Next router to put it in.
	return (o.Architecture == ArchSingle && !o.SingleUsesNext()) ||
		(o.ShouldEmbedAdmin() && o.UseTanStack())
}

// webAppRoot is the frontend application the panel lives inside.
//
// A single project's is frontend/, whether that is the Vite SPA or a Next app;
// everything else keeps apps/web. Writers that spelled out apps/web wrote into a
// directory a single project does not have, and nothing failed: the file was
// simply absent from the app, which is the same silence that shipped a Vite
// admin with no security screen.
func webAppRoot(root string, opts Options) string {
	if opts.Architecture == ArchSingle {
		// The root. A single project is a web app with a Go API in api/, so its
		// package.json, next.config.ts and app/ sit where anyone opening the
		// repository expects to find them. A project from before that layout keeps
		// its frontend where it was put.
		if opts.LegacySingleFlat {
			return filepath.Join(root, "frontend")
		}
		return root
	}
	return filepath.Join(root, "apps", "web")
}

// spaHostRoot is webAppRoot under its old name, kept for the SPA writers that
// read as "the Vite app this panel is a section of".
func spaHostRoot(root string, opts Options) string {
	return webAppRoot(root, opts)
}

// singleAdminShellFiles are the Vite admin's own application scaffolding, which
// the SPA already has. Relative to apps/admin.
var singleAdminShellFiles = map[string]bool{
	"package.json":           true,
	"vite.config.ts":         true,
	"index.html":             true,
	"tsconfig.json":          true,
	"src/main.tsx":           true,
	"src/vite-env.d.ts":      true,
	"src/globals.css":        true, // the SPA's own, plus admin.css for the screens
	"src/routes/__root.tsx":  true, // the SPA owns the root route
	"src/routes/index.tsx":   true, // and its own home page
	"public/.gitkeep":        true,
	"components.json":        true,
	"vitest.config.ts":       true,
	"vitest.setup.ts":        true,
	"src/routeTree.gen.ts":   true, // the router plugin regenerates it
	"src/router.tsx":         true,
	"src/styles/globals.css": true,
}

// embeddedSingleAdminFileMap is the Vite admin, keyed by where it goes in the SPA.
func embeddedSingleAdminFileMap(root string, opts Options) map[string]string {
	// Where the panel is served, for the content transforms writeFile applies
	// to these files. Run and RunSingle set it too, early enough for files
	// this map does not hold; setting it here as well is what lets a test
	// build the map and get the same answer the scaffolder does.
	adminMountsAtRoot = opts.AdminOnly

	source := adminTanStackFileMap(root, opts)
	adminRoot := filepath.Join(root, "apps", "admin")
	spa := filepath.Join(spaHostRoot(root, opts), "src")

	out := map[string]string{}
	for path, content := range source {
		rel, err := filepath.Rel(adminRoot, path)
		if err != nil {
			continue
		}
		rel = filepath.ToSlash(rel)
		if singleAdminShellFiles[rel] || strings.HasPrefix(rel, "public/") ||
			strings.HasPrefix(rel, "__tests__/") {
			continue
		}
		if !strings.HasPrefix(rel, "src/") {
			continue // anything outside src/ is app scaffolding
		}
		inner := strings.TrimPrefix(rel, "src/")

		var dest string
		if strings.HasPrefix(inner, "routes/") {
			// The routes, under the SPA's /admin, or at its root when there is
			// no public site to share it with (grit#89).
			dest = filepath.Join(spa, "routes", adminRouteSegment(opts),
				filepath.FromSlash(strings.TrimPrefix(inner, "routes/")))
		} else {
			// Pages, components, lib, hooks and the resource definitions.
			dest = filepath.Join(spa, "admin-panel", filepath.FromSlash(inner))
		}
		out[dest] = content
	}

	// The section's own entry point: /admin sends you to the dashboard, and
	// so does / when the admin is the whole app.
	out[filepath.Join(spa, "routes", adminRouteSegment(opts), "index.tsx")] =
		singleAdminIndexRoute(opts)

	// A layout for the whole section, which exists to own the stylesheet. The
	// panel's screens need the print, toaster and dark-mode rules, and in a Vite
	// app a stylesheet is pulled in by importing it from a module that renders.
	// At the root there is no routes/route.tsx in TanStack's file convention:
	// __root.tsx is the layout, and it already exists, so the stylesheet is
	// imported there instead. See singleAdminRootStylesheet.
	if !opts.AdminOnly {
		out[filepath.Join(spa, "routes", "admin", "route.tsx")] = singleAdminSectionRoute()
	}

	// And the half of the admin stylesheet that is about screens. The SPA's
	// globals.css already declares the same tokens with the same values.
	out[filepath.Join(spa, "admin-panel", "admin.css")] = adminScreenCSS()
	return out
}

// singleAdminSectionRoute is the layout for everything under /admin.
//
// It exists for the stylesheet: the panel's screens need the print isolation, the
// toaster styling and the dark-mode rules, and in a Vite app a stylesheet reaches
// the bundle by being imported from a module that renders. The SPA's own
// globals.css already declares the tokens, with the same names and values, so this
// is the half that is about admin screens.
func singleAdminSectionRoute() string {
	return `import { createFileRoute, Outlet } from '@tanstack/react-router'

import '@admin/admin.css'

export const Route = createFileRoute('/admin')({
  component: AdminSection,
})

function AdminSection() {
  return <Outlet />
}
`
}

// adminRouteSegment is the directory the panel's routes live in, which is the
// path they are served at: "admin" beside a public site, and "." when the
// admin is the whole app (grit#89).
func adminRouteSegment(opts Options) string {
	if opts.AdminOnly {
		return "."
	}
	return "admin"
}

// singleAdminIndexRoute sends the panel's entry point to the dashboard.
func singleAdminIndexRoute(opts Options) string {
	prefix := "/admin"
	if opts.AdminOnly {
		prefix = ""
	}
	return `import { createFileRoute, redirect } from '@tanstack/react-router'

// ` + orSlash(prefix) + ` has nothing of its own to show: the panel starts at the dashboard.
export const Route = createFileRoute('` + prefix + `/')({
  beforeLoad: () => {
    throw redirect({ to: '` + prefix + `/dashboard' })
  },
})
`
}

// orSlash renders an empty path prefix as / for a sentence.
func orSlash(prefix string) string {
	if prefix == "" {
		return "/"
	}
	return prefix
}

// routeIDPattern finds the id a file route declares.
var routeIDPattern = regexp.MustCompile(`createFileRoute\('([^']*)'\)`)

// adminMountsAtRoot says the panel is the whole SPA, so its routes are served
// at / rather than under /admin (grit#89).
//
// Package state, which is not how the rest of this file works, and the reason
// is that embeddedSingleAdminContent is keyed on a destination path and called
// from writeFile, which has no Options. Two of its rewrites are wrong at the
// root: the route id prefix, and the link prefix that turns href="/dashboard"
// into "/admin/dashboard". The path cannot distinguish them, because in a
// normal single project src/routes/ also holds the public site, and giving
// that the admin's alias rewrite would repoint its own components.
//
// Set by Run, RunSingle and embeddedSingleAdminFileMap, each of which has the
// Options and runs before any of these files are written, and set on every
// call rather than only when true so a second scaffold in one process cannot
// inherit the first one's answer.
var adminMountsAtRoot bool

// embeddedSingleAdminContent repoints a Vite admin file living in the SPA.
//
// Three things beyond what the Next panel needs: the route id inside each shim has
// to match its new path or the router refuses it, the lazy page import has to
// point at admin-panel, and the pages alias is @admin/pages rather than a
// directory the SPA has.
func embeddedSingleAdminContent(path, content string) string {
	slashed := filepath.ToSlash(path)
	// Keyed on the destination, and the destination is inside a Vite app's src:
	// frontend/src for a single, apps/web/src for a double built with --vite. The
	// Next.js panel lives at apps/web/admin-panel, with no src, so the two
	// transforms cannot both claim a file.
	// Where the panel's route files are. Under --admin they are the SPA's own
	// routes at the root, and matching only /src/routes/admin/ skipped them
	// entirely: forty shims kept their "@/pages/..." imports, which resolve to
	// a directory the SPA does not have.
	//
	// Widening to /src/routes/ is correct only in that shape, and only because
	// of what defines it: with no public site, nothing else lives there. In a
	// normal single project this would hand _site.tsx the admin's alias
	// rewrite and repoint the site's own components.
	adminRoutes := "/src/routes/admin/"
	if adminMountsAtRoot {
		adminRoutes = "/src/routes/"
	}
	if !strings.Contains(slashed, "/src/admin-panel/") &&
		!strings.Contains(slashed, adminRoutes) {
		return content
	}

	// The admin's own code, behind its own alias.
	for _, dir := range []string{"pages", "components", "lib", "hooks", "resources", "config"} {
		content = strings.ReplaceAll(content, `"@/`+dir+`/`, `"@admin/`+dir+`/`)
		content = strings.ReplaceAll(content, `'@/`+dir+`/`, `'@admin/`+dir+`/`)
		content = strings.ReplaceAll(content, `"@/`+dir+`"`, `"@admin/`+dir+`"`)
		content = strings.ReplaceAll(content, `'@/`+dir+`'`, `'@admin/`+dir+`'`)
	}

	// The route id, which TanStack requires to equal the file's path.
	//
	// Untouched when the panel is the whole app: its own ids are already the
	// paths it will be served at, and prefixing them would mount an admin with
	// no site at /admin anyway.
	if !adminMountsAtRoot {
		content = routeIDPattern.ReplaceAllStringFunc(content, func(match string) string {
			id := routeIDPattern.FindStringSubmatch(match)[1]
			if strings.HasPrefix(id, "/admin") {
				return match
			}
			return "createFileRoute('/admin" + id + "')"
		})
	}

	// The API base. A triple admin runs on its own port and talks to the API
	// across origins, so it defaults to http://localhost:8080. Here the Go binary
	// serves this page and the API, and in dev Vite proxies to it: the origin of
	// the page is the API, in development and in production both.
	content = strings.ReplaceAll(content,
		`|| "http://localhost:8080"`, `|| window.location.origin`)

	// And the links, as in the Next panel: the pages say href="/dashboard", and in
	// the SPA that screen is at /admin/dashboard.
	//
	// Not when the panel is the whole app: there the pages are already right,
	// and prefixing would point every link at an /admin that does not exist.
	if !adminMountsAtRoot {
		for _, segment := range adminRoutePrefixes {
			for _, quote := range []string{`"`, "`", `'`} {
				content = strings.ReplaceAll(content, quote+"/"+segment, quote+"/admin/"+segment)
			}
		}
		content = strings.ReplaceAll(content, `href="/"`, `href="/admin"`)
	}

	// Vite has no process. It does not polyfill one and esbuild does not
	// typecheck, so a process.env read reaches the browser as "process is not
	// defined" rather than failing the build. Every NEXT_PUBLIC_ the panel
	// reads needs its Vite name here, and TestSingleAdminContentIsRepointed is
	// what catches a new one: it caught this very entry being added without it.
	for from, to := range map[string]string{
		"process.env.NEXT_PUBLIC_IP_HINT": "import.meta.env.VITE_IP_HINT",
		"process.env.NODE_ENV":            "import.meta.env.MODE",
	} {
		content = strings.ReplaceAll(content, from, to)
	}

	return strings.ReplaceAll(content, "/admin/admin/", "/admin/")
}

// writeEmbeddedSingleAdminFiles writes the admin panel into the SPA.
func writeEmbeddedSingleAdminFiles(root string, opts Options) error {
	if !opts.ShouldEmbedAdminInSPA() {
		return nil
	}
	for path, content := range embeddedSingleAdminFileMap(root, opts) {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	// The screens written per-resource or in groups rather than as map entries.
	for _, write := range adminExtraWriters {
		if err := write(root, opts); err != nil {
			return err
		}
	}
	return nil
}

// ensureSPAAdminWiring gives an existing SPA what the embedded panel needs.
//
// Three files, none of which can be rewritten from the template: frontend's
// package.json holds whatever dependencies the project added, vite.config.ts its
// own aliases and plugins, tsconfig.json its own paths. So each entry is checked
// and inserted only when it is missing, and a project that already has it is left
// exactly as it is.
//
// Returns the files it changed, for the upgrade to report.
func ensureSPAAdminWiring(feRoot string, opts Options) ([]string, error) {
	var changed []string

	// The dependencies the panel's screens import. Tiptap for the rich-text
	// field, recharts for the dashboard widgets, xlsx and @react-pdf/renderer for
	// the exports, and the rest for the form builder.
	pkg := filepath.Join(feRoot, "package.json")
	deps := [][2]string{
		{"@base-ui/react", baseUIVersion},
		{"libphonenumber-js", libphonenumberVersion},
		{"@hookform/resolvers", HookformResolversVersion},
		{"@react-pdf/renderer", "^4.1.5"},
		{"react-dropzone", "^14.2.0"},
		{"react-hook-form", "^7.55.0"},
		{"recharts", "^2.12.0"},
		{"sonner", "^1.3.0"},
		{"tw-animate-css", "^1.4.0"},
		// The CDN tarball, not the registry. SheetJS stopped publishing to npm
		// after 0.18.5, and that version carries two HIGH advisories: prototype
		// pollution, fixed in 0.19.3, and a ReDoS, fixed in 0.20.2. So "move it
		// to a registry version" means shipping both. The single architecture was
		// on ^0.18.5 while every other shape used the tarball (grit#148).
		//
		// The price is the pnpm pin in infra_hygiene_repair.go: 10.34 and later
		// refuse to install a tarball URL. That pin is the cost of not shipping
		// two known vulnerabilities, which is the right way round.
		{"xlsx", "https://cdn.sheetjs.com/xlsx-0.20.3/xlsx-0.20.3.tgz"},
	}
	for _, name := range tiptapPackages {
		deps = append(deps, [2]string{name, tiptapVersion})
	}
	var depLines []string
	for _, d := range deps {
		depLines = append(depLines, fmt.Sprintf(`    %q: %q,`, d[0], d[1]))
	}
	if ok, err := insertMissingAfter(pkg, `"dependencies": {`, depLines); err != nil {
		return changed, err
	} else if ok {
		changed = append(changed, "frontend/package.json")
	}

	// The aliases. Vite resolves imports itself and does not read tsconfig paths,
	// so both files need telling, and the subpath entries go in ahead of the
	// directory ones or the prefix match swallows them.
	vite := filepath.Join(feRoot, "vite.config.ts")
	// How far packages/ is from the file holding the alias.
	//
	// A single project's frontend is the project root, so packages/ is a child.
	// A monorepo web app is two levels below it. And a single scaffolded before
	// that layout keeps its frontend in frontend/, one level down, which is the
	// case that only bites a project missing the alias entirely: the insert below
	// skips one that is already there, so an upgraded project survived this by
	// luck rather than by being right.
	up := "./packages/upload/src"
	if opts.LegacySingleFlat {
		up = "../packages/upload/src"
	}
	aliases := []string{
		`      '@admin': path.resolve(__dirname, './src/admin-panel'),`,
		`      '@repo/shared/brand': path.resolve(__dirname, './src/shared/brand.config.ts'),`,
		`      '@repo/shared': path.resolve(__dirname, './src/shared'),`,
	}
	if opts.Architecture != ArchSingle {
		up = "../../packages/upload/src"
		// The shared package is a workspace dependency there, resolved by node.
		aliases = aliases[:1]
	}
	aliases = append(aliases,
		"      '@repo/upload/web': path.resolve(__dirname, '"+up+"/web.ts'),",
		"      '@repo/upload': path.resolve(__dirname, '"+up+"/index.ts'),")
	if ok, err := insertMissingAfter(vite, "alias: {", aliases); err != nil {
		return changed, err
	} else if ok {
		changed = append(changed, "frontend/vite.config.ts")
	}

	// The dev proxy, so the panel's Studio, Pulse and Sentinel links reach the Go
	// binary rather than the Vite dev server.
	if ok, err := insertMissingAfter(vite, "proxy: {", []string{
		`      '/studio': { target: 'http://localhost:8080', changeOrigin: true },`,
		`      '/pulse': { target: 'http://localhost:8080', changeOrigin: true },`,
		`      '/sentinel': { target: 'http://localhost:8080', changeOrigin: true },`,
		`      '/docs': { target: 'http://localhost:8080', changeOrigin: true },`,
	}); err != nil {
		return changed, err
	} else if ok && !contains(changed, "frontend/vite.config.ts") {
		changed = append(changed, "frontend/vite.config.ts")
	}

	ts := filepath.Join(feRoot, "tsconfig.json")
	paths := []string{
		`      "@admin/*": ["./src/admin-panel/*"],`,
		"      \"@repo/upload/web\": [\"" + up + "/web.ts\"],",
		"      \"@repo/upload\": [\"" + up + "/index.ts\"],",
	}
	if opts.Architecture == ArchSingle {
		paths = append(paths,
			`      "@repo/shared/brand": ["./src/shared/brand.config.ts"],`,
			`      "@repo/shared/*": ["./src/shared/*"],`)
	}
	if ok, err := insertMissingAfter(ts, `"paths": {`, paths); err != nil {
		return changed, err
	} else if ok {
		changed = append(changed, "frontend/tsconfig.json")
	}

	return changed, nil
}

// insertMissingAfter inserts each line after the anchor, skipping any whose key
// the file already has.
//
// The key is the part before the first colon, which is what identifies an entry
// in all three of these files. So a project that pinned its own version of a
// dependency, or aliased @repo/shared somewhere else, keeps what it chose.
func insertMissingAfter(path, anchor string, lines []string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		// Nothing to wire: a project without this file is not a single project.
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	content := string(data)
	idx := strings.Index(content, anchor)
	if idx < 0 {
		return false, nil
	}
	cut := idx + len(anchor)

	var missing []string
	for _, line := range lines {
		key := strings.TrimSpace(line)
		if c := strings.Index(key, ":"); c > 0 {
			key = key[:c]
		}
		if strings.Contains(content, key+":") {
			continue
		}
		missing = append(missing, line)
	}
	if len(missing) == 0 {
		return false, nil
	}
	out := content[:cut] + "\n" + strings.Join(missing, "\n") + content[cut:]
	if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
		return false, err
	}
	return true, nil
}

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

// singleRootRoute is the SPA's root layout.
//
// Normally that is the shared TanStack root, and the admin's own stylesheet is
// imported by routes/admin/route.tsx, the layout for everything under /admin.
//
// When the admin is the whole app there is no such layout: TanStack's file
// convention has no routes/route.tsx, the root IS __root.tsx. So the
// stylesheet import moves here, which is the only thing that layout carried
// (grit#89).
//
// Getting this wrong is quiet rather than loud: in a Vite app a stylesheet
// reaches the bundle by being imported from a module that renders, so a
// missing import is not an error, it is a panel with no print isolation, no
// toaster styling and no dark mode.
func singleRootRoute(opts Options) string {
	if !opts.AdminOnly {
		return webTanStackRootRoute(opts)
	}
	return `import { createRootRoute, Outlet } from '@tanstack/react-router'
// grit:web-org-imports

// The admin is the whole app here, so its stylesheet is imported at the root.
// In a plain --single project this lives in routes/admin/route.tsx, the layout
// for the /admin section, which this shape has no equivalent of.
import '@admin/admin.css'

export const Route = createRootRoute({
  component: () => (
    <>
      {/* grit:web-org-switcher */}
      <Outlet />
    </>
  ),
})
`
}
