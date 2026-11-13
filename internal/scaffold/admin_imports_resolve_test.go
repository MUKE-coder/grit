package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// Every "@/..." an admin template imports has to be a file the admin ships.
//
// The sidebar imported the command palette from "@/components/command-palette"
// while the palette was registered at "components/chrome/CommandPalette.tsx".
// Nothing in Go noticed: both are strings. The project scaffolded cleanly, the
// upgrade reported 192 files updated, and every single page then answered 500,
// because the import is resolved by Next and not by anything that runs here.
//
// This is cheap to check and it covers the whole class: a file moved, a file
// renamed, a file registered in one admin shape and not the other.
func TestEveryAdminImportResolves(t *testing.T) {
	opts := Options{ProjectName: "shop", Version: "9.9.9", Architecture: ArchTriple, Frontend: FrontendNext}
	root := t.TempDir()
	adminRoot := filepath.Join(root, "apps", "admin")
	files := adminFileMap(root, opts)

	// The map is not the whole admin: four writers of their own put files
	// beside it, and an import may point at one of those. Running them into
	// the same temp root is how the set here matches what a project gets.
	for _, write := range []func(string, Options) error{
		writeAdminResourceDashboardWidgets,
		writeAdminCustomChartFiles,
		writeAdminPasskeyFiles,
		writePasswordStrengthFiles,
	} {
		if err := write(root, opts); err != nil {
			t.Fatalf("running an admin writer: %v", err)
		}
	}
	_ = filepath.Walk(adminRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if _, already := files[path]; !already {
			files[path] = ""
		}
		return nil
	})

	// What the admin ships, as module specifiers.
	have := map[string]bool{}
	for path := range files {
		rel, err := filepath.Rel(adminRoot, path)
		if err != nil || strings.HasPrefix(rel, "..") {
			continue
		}
		rel = filepath.ToSlash(rel)
		rel = strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(
			rel, ".tsx"), ".ts"), ".jsx"), ".js")
		have["@/"+rel] = true
		// A directory import resolves to its index file.
		if base := strings.TrimSuffix(rel, "/index"); base != rel {
			have["@/"+base] = true
		}
	}

	// Specifiers that are written by something other than this map: the
	// per-resource registry, the generator's own output, and the i18n
	// messages a project may or may not have wired.
	skip := func(spec string) bool {
		for _, p := range []string{
			"@/resources",          // the registry and every generated resource under it
			"@/messages",           // i18n, written only when i18n is on
			"@/app",                // route files reached by Next's router, not by import
			"@/i18n",               // next-intl's own config, written by the i18n wiring
			"@/components/plugins", // plugin-injected components
		} {
			if spec == p || strings.HasPrefix(spec, p+"/") {
				return true
			}
		}
		return false
	}

	importRe := regexp.MustCompile(`from\s+"(@/[^"]+)"`)

	var bad []string
	for path, body := range files {
		rel, err := filepath.Rel(adminRoot, path)
		if err != nil || strings.HasPrefix(rel, "..") {
			continue
		}
		for _, m := range importRe.FindAllStringSubmatch(body, -1) {
			spec := m[1]
			if skip(spec) || have[spec] {
				continue
			}
			bad = append(bad, filepath.ToSlash(rel)+"  imports  "+spec)
		}
	}

	if len(bad) > 0 {
		sort.Strings(bad)
		t.Errorf("%d import(s) point at a file the admin does not ship, so every page "+
			"that reaches them answers 500:\n  %s", len(bad), strings.Join(bad, "\n  "))
	}
}

// The two admin shapes ship the same components.
//
// A component registered in the Next map and not the Vite one is missing from
// every --single and --admin project, and CI does not build the Vite admin, so
// nothing says so until somebody opens it.
func TestBothAdminShapesShipTheSameComponents(t *testing.T) {
	root := t.TempDir()
	next := adminFileMap(root, Options{
		ProjectName: "shop", Version: "9.9.9",
		Architecture: ArchTriple, Frontend: FrontendNext,
	})
	vite := adminTanStackFileMap(root, Options{
		ProjectName: "shop", Version: "9.9.9",
		Architecture: ArchSingle, Frontend: FrontendTanStack,
	})

	// Compare by the part that matters: components/chrome and hooks, which is
	// where this release put things. Pages differ by router and are not
	// expected to line up.
	names := func(files map[string]string, under ...string) map[string]bool {
		out := map[string]bool{}
		for path := range files {
			p := filepath.ToSlash(path)
			for _, dir := range under {
				if i := strings.Index(p, dir); i >= 0 {
					out[p[i+len(dir):]] = true
				}
			}
		}
		return out
	}

	nextNames := names(next, "/components/chrome/", "/hooks/")
	viteNames := names(vite, "/components/chrome/", "/hooks/")

	var missing []string
	for name := range nextNames {
		if !viteNames[name] {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		t.Errorf("the Vite admin is missing %d file(s) the Next admin ships, and CI never "+
			"builds the Vite admin:\n  %s", len(missing), strings.Join(missing, "\n  "))
	}
}
