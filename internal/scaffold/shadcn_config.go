package scaffold

import "fmt"

// components.json, the file the shadcn CLI reads before it will install
// anything.
//
// Why every scaffolded frontend ships one: Grit UI is distributed as a shadcn
// registry, and the documented way to install a block is
//
//	npx shadcn@latest add https://ui.gritframework.dev/r/<block>.json
//
// Without this file that command does not fail, which would at least be clear.
// It drops into an interactive setup instead:
//
//	? You need to create a components.json file to add components. Proceed?
//	? Select a component library >  Base UI (Recommended) / React Aria / Radix UI
//	? Which color would you like to use as the base color?
//	...
//
// So the first thing anybody does with Grit UI is answer four questions about a
// project that already has all the answers, and one of them (the component
// library) has a wrong option that looks recommended. Grit is built on Radix
// through shadcn/ui, and picking Base UI there produces components that import
// packages the project does not have.
//
// Every value below is read from the project rather than guessed:
//
//   - the tailwind config and css paths are the ones the scaffold wrote
//   - lib/utils.ts already exports cn(), which is what "utils" points at
//   - lucide-react is already a dependency, so iconLibrary matches
//   - cssVariables is true because the scaffold's globals.css defines the
//     shadcn variable set
//
// baseColor is slate, which only affects components the CLI generates from
// scratch. Grit UI blocks carry their own palette, so it changes nothing for
// them and stays a sensible default for a plain `shadcn add button`.
func shadcnComponentsJSON(cfg shadcnPaths) string {
	return fmt.Sprintf(`{
  "$schema": "https://ui.shadcn.com/schema.json",
  "style": "new-york",
  "rsc": %t,
  "tsx": true,
  "tailwind": {
    "config": "tailwind.config.ts",
    "css": %q,
    "baseColor": "slate",
    "cssVariables": true,
    "prefix": ""
  },
  "iconLibrary": "lucide",
  "aliases": {
    "components": "%s/components",
    "utils": "%s/lib/utils",
    "ui": "%s/components/ui",
    "lib": "%s/lib",
    "hooks": "%s/hooks"
  }
}
`, cfg.RSC, cfg.CSS, cfg.Alias, cfg.Alias, cfg.Alias, cfg.Alias, cfg.Alias)
}

// shadcnPaths is what differs between a Next.js app and a Vite one.
type shadcnPaths struct {
	// RSC is true for the App Router, false for Vite. The CLI uses it to decide
	// whether a generated component needs "use client".
	RSC bool
	// CSS is where the Tailwind entry point lives, relative to the app root.
	CSS string
	// Alias is the import prefix the app's tsconfig maps, without a trailing
	// slash. Both layouts use "@", pointing at different roots.
	Alias string
}

// nextComponentsJSON is the config for a Next.js app (web, admin).
func nextComponentsJSON() string {
	return shadcnComponentsJSON(shadcnPaths{
		RSC:   true,
		CSS:   "app/globals.css",
		Alias: "@",
	})
}

// viteComponentsJSON is the config for a Vite app (--frontend vite, single
// mode). Same aliases, different css path, and no server components.
func viteComponentsJSON() string {
	return shadcnComponentsJSON(shadcnPaths{
		RSC:   false,
		CSS:   "src/globals.css",
		Alias: "@",
	})
}
