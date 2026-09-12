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

// ShouldEmbedAdminInSPA reports whether the panel lives inside the single
// binary's SPA.
func (o Options) ShouldEmbedAdminInSPA() bool {
	return o.Architecture == ArchSingle
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
	source := adminTanStackFileMap(root, opts)
	adminRoot := filepath.Join(root, "apps", "admin")
	spa := filepath.Join(root, "frontend", "src")

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
			// The routes, under the SPA's /admin.
			dest = filepath.Join(spa, "routes", "admin",
				filepath.FromSlash(strings.TrimPrefix(inner, "routes/")))
		} else {
			// Pages, components, lib, hooks and the resource definitions.
			dest = filepath.Join(spa, "admin-panel", filepath.FromSlash(inner))
		}
		out[dest] = content
	}

	// The section's own entry point: /admin sends you to the dashboard.
	out[filepath.Join(spa, "routes", "admin", "index.tsx")] = singleAdminIndexRoute()

	// A layout for the whole section, which exists to own the stylesheet. The
	// panel's screens need the print, toaster and dark-mode rules, and in a Vite
	// app a stylesheet is pulled in by importing it from a module that renders.
	out[filepath.Join(spa, "routes", "admin", "route.tsx")] = singleAdminSectionRoute()

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

// singleAdminIndexRoute sends /admin to the dashboard.
func singleAdminIndexRoute() string {
	return `import { createFileRoute, redirect } from '@tanstack/react-router'

// /admin has nothing of its own to show: the panel starts at the dashboard.
export const Route = createFileRoute('/admin/')({
  beforeLoad: () => {
    throw redirect({ to: '/admin/dashboard' })
  },
})
`
}

// routeIDPattern finds the id a file route declares.
var routeIDPattern = regexp.MustCompile(`createFileRoute\('([^']*)'\)`)

// embeddedSingleAdminContent repoints a Vite admin file living in the SPA.
//
// Three things beyond what the Next panel needs: the route id inside each shim has
// to match its new path or the router refuses it, the lazy page import has to
// point at admin-panel, and the pages alias is @admin/pages rather than a
// directory the SPA has.
func embeddedSingleAdminContent(path, content string) string {
	slashed := filepath.ToSlash(path)
	if !strings.Contains(slashed, "/frontend/src/admin-panel/") &&
		!strings.Contains(slashed, "/frontend/src/routes/admin/") {
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
	content = routeIDPattern.ReplaceAllStringFunc(content, func(match string) string {
		id := routeIDPattern.FindStringSubmatch(match)[1]
		if strings.HasPrefix(id, "/admin") {
			return match
		}
		return "createFileRoute('/admin" + id + "')"
	})

	// The API base. A triple admin runs on its own port and talks to the API
	// across origins, so it defaults to http://localhost:8080. Here the Go binary
	// serves this page and the API, and in dev Vite proxies to it: the origin of
	// the page is the API, in development and in production both.
	content = strings.ReplaceAll(content,
		`|| "http://localhost:8080"`, `|| window.location.origin`)

	// And the links, as in the Next panel: the pages say href="/dashboard", and in
	// the SPA that screen is at /admin/dashboard.
	for _, segment := range adminRoutePrefixes {
		for _, quote := range []string{`"`, "`", `'`} {
			content = strings.ReplaceAll(content, quote+"/"+segment, quote+"/admin/"+segment)
		}
	}
	content = strings.ReplaceAll(content, `href="/"`, `href="/admin"`)
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
func ensureSPAAdminWiring(root string, opts Options) ([]string, error) {
	feRoot := filepath.Join(root, "frontend")
	var changed []string

	// The dependencies the panel's screens import. Tiptap for the rich-text
	// field, recharts for the dashboard widgets, xlsx and @react-pdf/renderer for
	// the exports, and the rest for the form builder.
	pkg := filepath.Join(feRoot, "package.json")
	deps := [][2]string{
		{"@hookform/resolvers", "^3.3.0"},
		{"@react-pdf/renderer", "^4.1.5"},
		{"@tiptap/extension-color", "^2.1.0"},
		{"@tiptap/extension-highlight", "^2.1.0"},
		{"@tiptap/extension-image", "^2.1.0"},
		{"@tiptap/extension-link", "^2.1.0"},
		{"@tiptap/extension-placeholder", "^2.1.0"},
		{"@tiptap/extension-table", "^2.1.0"},
		{"@tiptap/extension-table-cell", "^2.1.0"},
		{"@tiptap/extension-table-header", "^2.1.0"},
		{"@tiptap/extension-table-row", "^2.1.0"},
		{"@tiptap/extension-text-align", "^2.1.0"},
		{"@tiptap/extension-text-style", "^2.1.0"},
		{"@tiptap/extension-underline", "^2.1.0"},
		{"@tiptap/pm", "^2.1.0"},
		{"@tiptap/react", "^2.1.0"},
		{"@tiptap/starter-kit", "^2.1.0"},
		{"react-dropzone", "^14.2.0"},
		{"react-hook-form", "^7.49.0"},
		{"recharts", "^2.12.0"},
		{"sonner", "^1.3.0"},
		{"tw-animate-css", "^1.4.0"},
		{"xlsx", "^0.18.5"},
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
	if ok, err := insertMissingAfter(vite, "alias: {", []string{
		`      '@admin': path.resolve(__dirname, './src/admin-panel'),`,
		`      '@repo/shared/brand': path.resolve(__dirname, './src/shared/brand.config.ts'),`,
		`      '@repo/shared': path.resolve(__dirname, './src/shared'),`,
		`      '@repo/upload/web': path.resolve(__dirname, '../packages/upload/src/web.ts'),`,
		`      '@repo/upload': path.resolve(__dirname, '../packages/upload/src/index.ts'),`,
	}); err != nil {
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
	if ok, err := insertMissingAfter(ts, `"paths": {`, []string{
		`      "@repo/upload/web": ["../packages/upload/src/web.ts"],`,
		`      "@repo/upload": ["../packages/upload/src/index.ts"],`,
		`      "@repo/shared/brand": ["./src/shared/brand.config.ts"],`,
		`      "@repo/shared/*": ["./src/shared/*"],`,
		`      "@admin/*": ["./src/admin-panel/*"],`,
	}); err != nil {
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
