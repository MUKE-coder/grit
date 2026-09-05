package plugin

import (
	"strings"
	"testing"
)

// ctxFor builds a triple-tier context, which is the only architecture where
// the admin app exists and therefore the only one these checks apply to.
func ctxFor(frontend string) Context {
	return Context{
		Root:         "/tmp/project",
		Module:       "shop/apps/api",
		APIRoot:      "apps/api",
		Architecture: "triple",
		Frontend:     frontend,
	}
}

// adminPaths returns every file and injection target under apps/admin.
func adminPaths(p Plugin, ctx Context) []string {
	var out []string
	if p.Files != nil {
		for path := range p.Files(ctx) {
			if strings.HasPrefix(path, "apps/admin/") {
				out = append(out, path)
			}
		}
	}
	if p.Injections != nil {
		for _, inj := range p.Injections(ctx) {
			if strings.HasPrefix(inj.File, "apps/admin/") {
				out = append(out, inj.File)
			}
		}
	}
	return out
}

// A TanStack admin lives under apps/admin/src. Anything a plugin writes above
// that is somewhere the "@/" alias does not point, so nothing imports it.
//
// This is the failure that took a while to find, because it is not a build
// error: the plugin reports success, the admin compiles, and the feature is
// simply absent. It also makes the plugin uninstallable a second time, since
// the stray file counts as "already exists" and the install refuses.
func TestPluginsWriteIntoTheTanStackAdmin(t *testing.T) {
	ctx := ctxFor("tanstack")

	for name, p := range registry {
		paths := adminPaths(p, ctx)
		if len(paths) == 0 {
			continue // not an admin-touching plugin
		}
		for _, path := range paths {
			if !strings.HasPrefix(path, "apps/admin/src/") {
				t.Errorf("%s writes %s on a TanStack project; nothing under apps/admin/ "+
					"outside src/ is reachable from the app", name, path)
			}
			// The other direction: apps/admin/src/src/... from prefixing twice.
			if strings.Contains(path, "/src/src/") {
				t.Errorf("%s writes %s; the src/ prefix was applied twice", name, path)
			}
		}
	}
}

// And the Next admin keeps everything at the app root, so src/ must not appear.
func TestPluginsWriteIntoTheNextAdmin(t *testing.T) {
	ctx := ctxFor("next")

	for name, p := range registry {
		for _, path := range adminPaths(p, ctx) {
			if strings.HasPrefix(path, "apps/admin/src/") {
				t.Errorf("%s writes %s on a Next project, which has no src/ directory", name, path)
			}
		}
	}
}

// A plugin that ships a page needs a route for it on TanStack, because routing
// there is the file tree: a component under src/pages/ with nothing in
// src/routes/ pointing at it is a page with no URL.
func TestTanStackPagesHaveRoutes(t *testing.T) {
	ctx := ctxFor("tanstack")

	for name, p := range registry {
		if p.Files == nil {
			continue
		}
		files := p.Files(ctx)

		var pages, routes int
		for path := range files {
			if strings.HasPrefix(path, "apps/admin/src/pages/") {
				pages++
			}
			if strings.HasPrefix(path, "apps/admin/src/routes/") {
				routes++
			}
		}
		if pages > 0 && routes < pages {
			t.Errorf("%s ships %d page(s) under src/pages but only %d route(s); "+
				"the extra pages have no URL", name, pages, routes)
		}
	}
}

// The components are written in the Next dialect, and neither "use client" nor
// next/navigation exists in a Vite app.
func TestTanStackFilesAreConverted(t *testing.T) {
	ctx := ctxFor("tanstack")

	for name, p := range registry {
		if p.Files == nil {
			continue
		}
		for path, body := range p.Files(ctx) {
			if !strings.HasPrefix(path, "apps/admin/src/") {
				continue
			}
			if !strings.HasSuffix(path, ".tsx") && !strings.HasSuffix(path, ".ts") {
				continue
			}
			if strings.Contains(body, `"use client"`) {
				t.Errorf(`%s: %s still has "use client", which Vite does not understand`, name, path)
			}
			if strings.Contains(body, `"next/navigation"`) {
				t.Errorf("%s: %s still imports next/navigation, which is not installed", name, path)
			}
		}
	}
}
