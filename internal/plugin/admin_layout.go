package plugin

import (
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// adminDir is the admin app's source root for this project.
//
// The Next admin keeps components, hooks and lib at the app root. The TanStack
// admin puts all of them under src/, and its router reads src/routes/ rather
// than app/(dashboard)/.
//
// Every plugin here wrote the Next shape unconditionally, so on a TanStack
// project the files landed in apps/admin/components/ where the app's "@/" alias
// does not point and nothing imports them. That is not a build error, which is
// what made it hard to notice: the plugin reports success, the admin compiles,
// and the feature is simply not there. The injections missed for the same
// reason, so even the import line never appeared.
//
// It also made the plugin uninstallable a second time: the stray file counts as
// "already exists", and the install refuses rather than overwrite.
func adminDir(ctx Context) string {
	if ctx.Architecture == "double" {
		if ctx.Frontend == "tanstack" {
			return "apps/web/src/admin-panel"
		}
		return "apps/web/admin-panel"
	}
	if ctx.Frontend == "tanstack" {
		return "apps/admin/src"
	}
	return "apps/admin"
}

// adminFile returns a path under the admin source root.
//
//	adminFile(ctx, "components/command-palette.tsx")
func adminFile(ctx Context, rel string) string {
	return adminDir(ctx) + "/" + strings.TrimPrefix(rel, "/")
}

// adminPageFile returns where a dashboard page lives, and the route shim that
// reaches it on TanStack.
//
// routePath is the URL under the dashboard ("system/impersonate"); name is the
// file stem used under src/pages/ on TanStack. The second return is empty on
// Next, where the page file is its own route.
func adminPageFile(ctx Context, routePath, name string) (page string, route string) {
	// A double has no admin app: its panel is embedded in the web app, under
	// /admin. The route group keeps its name, so the URL stays /admin/<path>
	// and the same layout wraps it.
	if ctx.Architecture == "double" {
		if ctx.Frontend == "tanstack" {
			return "apps/web/src/admin-panel/pages/" + name + ".tsx",
				"apps/web/src/routes/admin/_dashboard/" + routePath + ".tsx"
		}
		return "apps/web/app/admin/(dashboard)/" + routePath + "/page.tsx", ""
	}
	if ctx.Frontend != "tanstack" {
		return adminDir(ctx) + "/app/(dashboard)/" + routePath + "/page.tsx", ""
	}
	return adminDir(ctx) + "/pages/" + name + ".tsx",
		adminDir(ctx) + "/routes/_dashboard/" + routePath + ".tsx"
}

// tanStackRouteShim is the file TanStack Router needs to put a component on a
// URL: the routing lives in the file tree, so the component cannot also be the
// route the way a Next page is.
func tanStackRouteShim(routePath, importPath string) string {
	return `import { createFileRoute } from '@tanstack/react-router'
import Page from '` + importPath + `'

export const Route = createFileRoute('/_dashboard/` + routePath + `')({
  component: Page,
})
`
}

// adminSource converts a component written in the Next dialect for whichever
// admin this project has.
//
// The plugin components use "use client" and next/navigation, neither of which
// exists in a Vite app. scaffold.NextToTanStack is the same converter the
// TanStack admin scaffold uses, so a rule added there reaches plugins without
// anybody remembering to add it twice.
//
// An embedded panel (double) also imports through @admin/, not @/: both aliases
// exist in that web app, and @/ resolves to the web source, so a component
// written for the standalone admin would import the wrong tree or nothing.
func adminSource(ctx Context, code string) string {
	if ctx.Frontend == "tanstack" {
		code = scaffold.NextToTanStack(code)
	}
	if ctx.Architecture == "double" {
		code = strings.ReplaceAll(code, `from "@/`, `from "@admin/`)
	}
	return code
}
