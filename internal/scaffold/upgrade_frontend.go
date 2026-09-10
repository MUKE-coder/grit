package scaffold

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// readProjectFrontend reads which frontend the project was scaffolded with.
//
// grit upgrade used to leave this unset, and Normalize fills an empty frontend
// with Next.js, so every project was upgraded as if its apps were Next apps.
// On a TanStack project that meant the Next admin and web apps, package.json
// included, were written over the Vite ones. Found upgrading a TanStack
// project: both apps came out depending on next, with vite and the router gone
// and two apps' worth of files in each directory.
//
// grit.json records the choice. A project from before it did falls back to
// what is on disk, where a Vite app is the one with a vite.config.ts.
func readProjectFrontend(root string) Frontend {
	if data, err := os.ReadFile(filepath.Join(root, "grit.json")); err == nil {
		var cfg struct {
			Frontend string `json:"frontend"`
		}
		if json.Unmarshal(data, &cfg) == nil {
			switch Frontend(cfg.Frontend) {
			case FrontendTanStack:
				return FrontendTanStack
			case FrontendNext:
				return FrontendNext
			}
		}
	}
	for _, app := range []string{"admin", "web"} {
		if fileExists(filepath.Join(root, "apps", app, "vite.config.ts")) {
			return FrontendTanStack
		}
	}
	return FrontendNext
}

// skipViteApp says a Vite app was left alone, and why.
//
// The refresh that upgrade performs is written for the Next.js layout. Until a
// Vite equivalent exists, leaving the app exactly as it is beats the one other
// option, which was replacing its package.json with a Next one.
func skipViteApp(app string) {
	fmt.Printf("  ⚠ apps/%s is a Vite (TanStack) app. grit upgrade does not refresh Vite apps yet,\n", app)
	fmt.Printf("    so it was left as it is rather than overwritten with the Next.js layout.\n")
}
