package scaffold

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// The admin panel, living inside the web app.
//
// A triple project has three apps, and the admin is one of them: its own
// Next.js app on its own port, which the web app links to with an absolute URL.
// A double has two, and until now that meant no admin panel at all. The link was
// still in the navbar, pointing at http://localhost:3001, where nothing was
// listening.
//
// So for a double the admin becomes a route group inside the web app: the same
// 147 screens, at /admin/dashboard, with their own layout. Not a second copy of
// the code. The file map is the one the standalone app uses, put through a
// transform that moves the routes under app/admin, moves everything else under
// admin-panel/, repoints the imports at it, and prefixes the internal links.
// One source, one rewrite, so a fix to an admin screen reaches both shapes.
//
// What is deliberately dropped: the app shell. package.json, next.config.ts,
// tsconfig.json, postcss, globals.css and the root layout all belong to the web
// app, which already has them. The pieces of those the admin genuinely needs (its
// dark-mode tokens, its dependencies, a path alias) are merged into the web app's
// own copies instead.

// ShouldEmbedAdmin reports whether the admin panel lives inside the web app.
//
// Double only for now. Single is a Vite SPA rather than Next, so it needs its own
// port of the same idea and will get it separately.
func (o Options) ShouldEmbedAdmin() bool {
	return o.Architecture == ArchDouble
}

// adminShellFiles are the standalone app's own scaffolding, which the web app
// already provides. Relative to apps/admin.
var adminShellFiles = map[string]bool{
	"package.json":      true,
	"next.config.ts":    true,
	"tsconfig.json":     true,
	"postcss.config.js": true,
	"components.json":   true,
	"Dockerfile":        true,
	"vitest.config.ts":  true,
	"vitest.setup.ts":   true,
	"app/globals.css":   true, // web's own, plus the dark-mode blocks merged in
	"app/layout.tsx":    true, // web owns <html> and <body>; see embeddedAdminLayout
}

// embeddedAdminFileMap is the admin panel, keyed by where it goes in apps/web.
func embeddedAdminFileMap(root string, opts Options) map[string]string {
	source := adminFileMap(root, opts)
	adminRoot := filepath.Join(root, "apps", "admin")
	webRoot := filepath.Join(root, "apps", "web")

	out := map[string]string{}
	for path, content := range source {
		rel, err := filepath.Rel(adminRoot, path)
		if err != nil {
			continue
		}
		rel = filepath.ToSlash(rel)

		if adminShellFiles[rel] || strings.HasPrefix(rel, "public/") ||
			strings.HasPrefix(rel, "__tests__/") {
			continue
		}

		var dest string
		switch {
		case strings.HasPrefix(rel, "app/"):
			// The routes, under the web app's /admin.
			dest = filepath.Join(webRoot, "app", "admin", filepath.FromSlash(strings.TrimPrefix(rel, "app/")))
		default:
			// Components, lib, hooks and the resource definitions, in one place
			// that is obviously the admin's and collides with nothing the web app
			// has of its own.
			dest = filepath.Join(webRoot, "admin-panel", filepath.FromSlash(rel))
		}
		// The content rewrite happens in writeFile, keyed on the destination, so
		// the feature writers that bypass this map get it too.
		out[dest] = content
	}

	// The layout the standalone app's root layout becomes: the providers, without
	// the document.
	out[filepath.Join(webRoot, "app", "admin", "layout.tsx")] = embeddedAdminLayout(opts)

	// And the half of the admin stylesheet that is about admin screens: print
	// isolation, the toaster, dark mode. The tokens are already in the web app's
	// own stylesheet, with the same names and values, so they are not repeated.
	out[filepath.Join(webRoot, "app", "admin", "admin.css")] = adminScreenCSS()
	return out
}

// adminRouteSegments collects the first URL segment of every admin route.
//
// app/(dashboard)/system/activity/page.tsx gives "system"; app/(auth)/login/page.tsx
// gives "login". Route groups in parentheses are not part of the URL, so they are
// stepped over.
func adminRouteSegments(source map[string]string, adminRoot string) []string {
	seen := map[string]bool{}
	for path := range source {
		rel, err := filepath.Rel(adminRoot, path)
		if err != nil {
			continue
		}
		rel = filepath.ToSlash(rel)
		if !strings.HasPrefix(rel, "app/") {
			continue
		}
		for _, segment := range strings.Split(strings.TrimPrefix(rel, "app/"), "/") {
			if strings.HasPrefix(segment, "(") || segment == "" {
				continue // a route group: not in the URL
			}
			if strings.HasSuffix(segment, ".tsx") || strings.HasSuffix(segment, ".ts") ||
				strings.HasSuffix(segment, ".css") {
				break // a file, so this route has no segment of its own
			}
			if strings.HasPrefix(segment, "[") {
				break // a dynamic segment cannot be the first one
			}
			seen[segment] = true
			break
		}
	}

	out := make([]string, 0, len(seen))
	for segment := range seen {
		out = append(out, segment)
	}
	// Longest first, so /system/activity is rewritten before /system would match
	// its prefix.
	sort.Slice(out, func(i, j int) bool {
		if len(out[i]) != len(out[j]) {
			return len(out[i]) > len(out[j])
		}
		return out[i] < out[j]
	})
	return out
}

// adminRoutePrefixes are the first URL segments the admin panel owns.
//
// Every absolute link to one of these has to gain the /admin prefix when the panel
// lives inside the web app: the admin's own code says href="/dashboard", and in
// the web app that page is at /admin/dashboard.
//
// A list rather than something derived, because the rewrite happens one file at a
// time and cannot see the others. TestEmbeddedAdminRoutesAreAllPrefixed walks what
// the scaffold actually writes and fails if a route is missing from here, so a new
// screen cannot quietly link to a page that does not exist.
var adminRoutePrefixes = []string{
	"forgot-password", "reset-password", "verify-email", "sign-up",
	"resources", "dashboard", "settings", "account", "profile", "callback",
	"system", "login",
}

// embeddedAdminContent repoints an admin file written inside the web app.
//
// Keyed on the destination path: anything under the web app's admin-panel
// directory or its app/admin route group is the panel's code and needs the
// rewrite. Anything else is returned untouched, so this is inert for a triple
// project and for every non-admin file.
func embeddedAdminContent(path, content string) string {
	slashed := filepath.ToSlash(path)
	if !strings.Contains(slashed, "/apps/web/admin-panel/") &&
		!strings.Contains(slashed, "/apps/web/app/admin/") {
		return content
	}
	return embedAdminContent(content, adminRoutePrefixes)
}

// embedAdminContent repoints an admin file's imports and links.
func embedAdminContent(content string, routes []string) string {
	// Imports: @/components/... is the admin's own code, which now lives under
	// admin-panel/ behind its own alias. The web app keeps @/ for its own files,
	// so nothing has to be renamed on either side.
	for _, dir := range []string{"components", "lib", "hooks", "resources", "config"} {
		content = strings.ReplaceAll(content, `"@/`+dir+`/`, `"@admin/`+dir+`/`)
		content = strings.ReplaceAll(content, `'@/`+dir+`/`, `'@admin/`+dir+`/`)
		content = strings.ReplaceAll(content, `"@/`+dir+`"`, `"@admin/`+dir+`"`)
	}

	// Links: every absolute link to a route the admin owns gains the prefix. The
	// admin's own code says href="/dashboard", and inside the web app that page is
	// at /admin/dashboard.
	for _, segment := range routes {
		for _, quote := range []string{`"`, "`", `'`} {
			content = strings.ReplaceAll(content, quote+"/"+segment, quote+"/admin/"+segment)
		}
	}
	// Anything that asked for the admin's root now asks for /admin.
	content = strings.ReplaceAll(content, `href="/"`, `href="/admin"`)
	content = strings.ReplaceAll(content, `replace("/")`, `replace("/admin")`)
	content = strings.ReplaceAll(content, `push("/")`, `push("/admin")`)

	// The prefix must not be applied twice, which a second pass over an already
	// rewritten string would do.
	content = strings.ReplaceAll(content, "/admin/admin/", "/admin/")
	return content
}

// embeddedAdminLayout is the admin's layout inside the web app.
//
// The standalone app's root layout owns <html>, <body>, the fonts and the theme
// attribute. In the web app those belong to the web app's root layout, so what is
// left is the part that matters: the providers the admin screens need, and a
// metadata title for the section.
func embeddedAdminLayout(opts Options) string {
	return `import type { Metadata } from "next";

import "./admin.css";
import { Providers } from "@admin/components/shared/providers";

export const metadata: Metadata = {
  title: "` + opts.ProjectName + ` Admin",
  description: "Admin panel — Built with Grit",
};

// The admin panel, as a route group inside this app.
//
// Everything under /admin renders inside these providers: React Query for the
// resource hooks, the theme, and the toaster. <html> and <body> are the web app's
// root layout's job, so this adds nothing to the document.
//
// Access is the layout's own business, not this file's: components/layout/admin-layout
// redirects anybody who is not signed in to the login page, and sends a plain USER
// role to their profile rather than the dashboard.
export default function AdminSectionLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <Providers>{children}</Providers>;
}
`
}

// writeEmbeddedAdminFiles writes the admin panel into the web app.
func writeEmbeddedAdminFiles(root string, opts Options) error {
	if !opts.ShouldEmbedAdmin() {
		return nil
	}
	for path, content := range embeddedAdminFileMap(root, opts) {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	// And the screens written per-resource or in groups rather than as map
	// entries: the dashboard widgets and the chart builder. They ask adminCodeRoot
	// for their destination, so they land in admin-panel/ here. Leaving them out
	// is what made the embedded panel's first build fail, on a dashboard importing
	// a widget that had been written into the app that does not exist.
	for _, write := range adminExtraWriters {
		if err := write(root, opts); err != nil {
			return err
		}
	}
	return nil
}

// leftoverAliasPattern finds an import the transform failed to repoint.
//
// Used by a test rather than at run time: a surviving @/ in an admin file would
// resolve to the WEB app's file of that name, which either does not exist or is
// the wrong one, and the failure would be a type error in a generated project
// rather than anything visible here.
var leftoverAliasPattern = regexp.MustCompile(`["']@/(?:components|lib|hooks|resources|config)/`)
