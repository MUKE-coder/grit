package scaffold

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Contact-app review L34: about 1,150 lines of admin code that nothing imported.
//
// components/layout/sidebar.tsx has not been rendered since the chrome moved to
// components/chrome/CollapsibleSidebar.tsx. components/resource/view-modal.tsx
// stopped being used when records got their own detail page. The four files in
// components/widgets/ are used by the modern, minimal and glass dashboards only,
// and were written for every style. And a panel inside the web app carried its
// own copy of the realtime client beside the web app's, 476 lines that would
// have opened a second socket for anyone who imported them.
//
// A new project no longer gets any of them. Upgrade stops writing them, so an
// existing project's copies would sit in the tree forever; this removes each
// one when Grit wrote it, nobody has edited it, and nothing imports it. A file
// somebody changed, or that their own code imports, stays where it is.

// deadPanelFile is a file the panel no longer uses, relative to its code root.
type deadPanelFile struct {
	rel    string
	reason string
}

// deadPanelFiles are removed in this order. widget-grid.tsx goes before the
// other widgets because it is the file that imports them.
var deadPanelFiles = []deadPanelFile{
	{"components/layout/sidebar.tsx", "the sidebar is components/chrome/CollapsibleSidebar.tsx"},
	{"components/resource/view-modal.tsx", "a record opens on its own detail page"},
	{"components/widgets/widget-grid.tsx", "only the modern, minimal and glass dashboards use the widgets"},
	{"components/widgets/stats-card.tsx", "only the modern, minimal and glass dashboards use the widgets"},
	{"components/widgets/chart-widget.tsx", "only the modern, minimal and glass dashboards use the widgets"},
	{"components/widgets/activity-widget.tsx", "only the modern, minimal and glass dashboards use the widgets"},
}

// sharedRealtimeFiles are the panel's copy of the realtime client, removed only
// where the app the panel lives in has its own. use-realtime.ts goes first: it
// imports lib/realtime.ts.
var sharedRealtimeFiles = []deadPanelFile{
	{"hooks/use-realtime.ts", "the web app's hooks/use-realtime.ts serves the panel"},
	{"lib/realtime.ts", "the web app's lib/realtime.ts serves the panel"},
}

// panelTree is one admin panel as an import graph sees it.
type panelTree struct {
	// code is the panel's own directory.
	code string
	// app is the whole app, every file of which could import from the panel.
	app string
	// aliases maps an import prefix ("@/", "@admin/") to the directory it names.
	aliases map[string]string
}

// hostHasRealtime reports whether the app around an embedded panel has a
// realtime client of its own.
func (p panelTree) hostHasRealtime() bool {
	host := p.aliases["@/"]
	return host != "" && !sameFilePath(host, p.code) && fileExists(filepath.Join(host, "lib", "realtime.ts"))
}

// deadCodePanels finds the panels upgrade maintains: the Next.js admin app, the
// panel in a Next.js web app, and the panel in a SPA. A standalone Vite admin
// is not upgraded, so nothing is removed from it either.
func deadCodePanels(root string) []panelTree {
	var panels []panelTree
	for _, shape := range adminPanelShapes(root) {
		if shape.alias == "@" {
			panels = append(panels, panelTree{
				code:    shape.code,
				app:     shape.code,
				aliases: map[string]string{"@/": shape.code},
			})
			continue
		}
		host := filepath.Dir(shape.code)
		panels = append(panels, panelTree{
			code:    shape.code,
			app:     host,
			aliases: map[string]string{"@/": host, shape.alias + "/": shape.code},
		})
	}
	for _, host := range []string{filepath.Join(root, "frontend"), filepath.Join(root, "apps", "web")} {
		src := filepath.Join(host, "src")
		code := filepath.Join(src, "admin-panel")
		if dirExists(code) {
			panels = append(panels, panelTree{
				code:    code,
				app:     src,
				aliases: map[string]string{"@/": src, "@admin/": code},
			})
		}
	}
	return panels
}

// pruneDeadFrontendFiles applies L34 to every panel in the project.
func pruneDeadFrontendFiles(root string) error {
	// The recording, not the manifest on disk: the panel was written earlier in
	// this same upgrade, and the file on disk is not saved until it ends.
	return pruneDeadPanelFiles(root, deadCodePanels(root), manifest.IsUnchanged)
}

func pruneDeadPanelFiles(root string, panels []panelTree, pristine func(string) bool) error {
	for _, panel := range panels {
		files := deadPanelFiles
		if panel.hostHasRealtime() {
			files = append(append([]deadPanelFile{}, deadPanelFiles...), sharedRealtimeFiles...)
		}
		for _, dead := range files {
			path := filepath.Join(panel.code, filepath.FromSlash(dead.rel))
			if !fileExists(path) {
				continue
			}
			shown := path
			if rel, err := filepath.Rel(root, path); err == nil {
				shown = filepath.ToSlash(rel)
			}
			if !pristine(path) {
				// Untracked files predate the manifest or are the reader's own;
				// only an edit Grit can see is worth a line.
				if gritWroteFile(root, path) {
					fmt.Printf("  ⚠ %s: Grit no longer uses it and left it in place because it was edited (%s)\n", shown, dead.reason)
				}
				continue
			}
			importer, err := firstImporter(panel, path)
			if err != nil {
				return err
			}
			if importer != "" {
				continue
			}
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("removing %s: %w", shown, err)
			}
			manifest.Drop(path)
			removeDirIfEmpty(filepath.Dir(path))
			fmt.Printf("  ✓ %s: removed, nothing imported it (%s)\n", shown, dead.reason)
		}
	}
	return nil
}

// importSpecifierRe finds module specifiers: import/export ... from "x",
// import("x"), a side-effect import "x", require("x") and vi.mock("x").
var importSpecifierRe = regexp.MustCompile(`(?:\bfrom\s*|\bimport\s*\(\s*|\bimport\s+|\brequire\(\s*|\bvi\.mock\(\s*)["']([^"'\n]+)["']`)

var sourceExts = map[string]bool{".ts": true, ".tsx": true, ".js": true, ".jsx": true, ".mjs": true}

// firstImporter returns a file in the panel's app that imports target, or ""
// when none does. It resolves each specifier rather than searching for a name,
// so a comment that mentions the file is not an import and a relative import
// written from next door is.
func firstImporter(panel panelTree, target string) (string, error) {
	want := moduleKey(target)
	found := ""
	err := filepath.WalkDir(panel.app, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case "node_modules", ".next", "dist", ".turbo", ".tanstack":
				return filepath.SkipDir
			}
			return nil
		}
		if found != "" || !sourceExts[filepath.Ext(p)] || sameFilePath(p, target) {
			return nil
		}
		body, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		for _, m := range importSpecifierRe.FindAllStringSubmatch(string(body), -1) {
			if resolved := resolveSpecifier(panel, p, m[1]); resolved != "" && strings.EqualFold(moduleKey(resolved), want) {
				found = p
				return filepath.SkipAll
			}
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("looking for imports of %s: %w", target, err)
	}
	return found, nil
}

// resolveSpecifier turns an import specifier into a path, or "" for a package.
func resolveSpecifier(panel panelTree, from, spec string) string {
	if strings.HasPrefix(spec, "./") || strings.HasPrefix(spec, "../") {
		return filepath.Join(filepath.Dir(from), filepath.FromSlash(spec))
	}
	// Longest alias first, so "@admin/" is not read as "@/"-something.
	best := ""
	for prefix := range panel.aliases {
		if strings.HasPrefix(spec, prefix) && len(prefix) > len(best) {
			best = prefix
		}
	}
	if best == "" {
		return ""
	}
	return filepath.Join(panel.aliases[best], filepath.FromSlash(strings.TrimPrefix(spec, best)))
}

// moduleKey is a path as a module name: clean, slashed, without its extension
// or a trailing /index.
func moduleKey(path string) string {
	key := filepath.ToSlash(filepath.Clean(path))
	if ext := filepath.Ext(key); sourceExts[ext] {
		key = strings.TrimSuffix(key, ext)
	}
	return strings.TrimSuffix(key, "/index")
}

// removeDirIfEmpty deletes dir when the last file in it was just removed. A
// directory with anything left in it is none of this repair's business.
func removeDirIfEmpty(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) > 0 {
		return
	}
	if err := os.Remove(dir); err != nil {
		fmt.Printf("  ⚠ %s: empty, but could not be removed: %v\n", filepath.ToSlash(dir), err)
	}
}
