package scaffold

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Where a resource's files live inside the admin.
//
// One folder per resource, since v3.143.0:
//
//	resources/
//	  index.ts
//	  products/
//	    products.ts          generated, rewritten freely
//	    products.custom.tsx  yours, written once
//
// The flat layout that came before put both files directly in resources/, which
// reads fine with three resources and badly with twenty: the overlay doubled
// the file count and the two halves of one resource sorted apart from each
// other whenever a name fell between them alphabetically.
//
// The overlay import inside the definition is unaffected. It was always
// "./<name>.custom" and the two files are still siblings.
//
// Both layouts are read, because a project upgrades on its own schedule and
// `grit generate` must not break on one that has not moved yet. Writes always
// use the folder layout.

// resourceDefPath returns the path to write a resource definition to, and the
// directory that has to exist first.
func ResourceDefPath(resourcesRoot, kebab string) (dir string, file string) {
	dir = filepath.Join(resourcesRoot, kebab)
	return dir, filepath.Join(dir, kebab+".ts")
}

// resourceCustomPath returns the overlay path, beside its definition.
func ResourceCustomPath(resourcesRoot, kebab string) string {
	return filepath.Join(resourcesRoot, kebab, kebab+".custom.tsx")
}

// findResourceDef locates an existing definition in either layout, preferring
// the folder one. Returns "" when neither is present.
func FindResourceDef(resourcesRoot, kebab string) string {
	nested := filepath.Join(resourcesRoot, kebab, kebab+".ts")
	if isFile(nested) {
		return nested
	}
	flat := filepath.Join(resourcesRoot, kebab+".ts")
	if isFile(flat) {
		return flat
	}
	return ""
}

// findResourceCustom locates an existing overlay in either layout.
func FindResourceCustom(resourcesRoot, kebab string) string {
	nested := filepath.Join(resourcesRoot, kebab, kebab+".custom.tsx")
	if isFile(nested) {
		return nested
	}
	flat := filepath.Join(resourcesRoot, kebab+".custom.tsx")
	if isFile(flat) {
		return flat
	}
	return ""
}

// registryImportPath is what resources/index.ts imports the definition from.
// Relative, extensionless, forward slashes: this goes into TypeScript, not into
// a filesystem call, so filepath.Join would produce backslashes on Windows and
// a module nobody can resolve.
func RegistryImportPath(kebab string) string {
	return "./" + kebab + "/" + kebab
}

// MigrateResourceLayout moves a flat resources/ directory into one folder per
// resource and rewrites the registry's imports to match.
//
// Returns the number of resources moved. Safe to run repeatedly: a resource
// already in a folder is skipped, and nothing is deleted, only moved.
func MigrateResourceLayout(resourcesRoot string) (int, error) {
	entries, err := os.ReadDir(resourcesRoot)
	if err != nil {
		return 0, nil // no resources directory: nothing to do, not an error
	}

	moved := 0
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".ts") || strings.HasSuffix(name, ".custom.tsx") {
			continue
		}
		kebab := strings.TrimSuffix(name, ".ts")
		// index.ts is the registry, not a resource, and registry.ts was an old
		// name for the same thing.
		if kebab == "index" || kebab == "registry" {
			continue
		}

		dir := filepath.Join(resourcesRoot, kebab)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return moved, err
		}
		if err := os.Rename(filepath.Join(resourcesRoot, name), filepath.Join(dir, name)); err != nil {
			return moved, err
		}
		// The overlay travels with its definition. Moving one without the other
		// would break the "./<name>.custom" import that assumes they are
		// siblings.
		overlay := kebab + ".custom.tsx"
		if isFile(filepath.Join(resourcesRoot, overlay)) {
			if err := os.Rename(
				filepath.Join(resourcesRoot, overlay),
				filepath.Join(dir, overlay),
			); err != nil {
				return moved, err
			}
		}
		// A .bak left by `grit remove resource` goes too, so it does not sit in
		// the root looking like a resource that failed to move.
		if isFile(filepath.Join(resourcesRoot, overlay+".bak")) {
			_ = os.Rename(
				filepath.Join(resourcesRoot, overlay+".bak"),
				filepath.Join(dir, overlay+".bak"),
			)
		}
		moved++
	}

	// Both rewrites run every time, not only when something moved. A project
	// that moved under an earlier build has the folders but not the corrected
	// imports, and gating on `moved > 0` would never fix it. Both patterns
	// refuse an already-nested path, so repeating them changes nothing.
	if err := rewriteRegistryImports(filepath.Join(resourcesRoot, "index.ts")); err != nil {
		return moved, err
	}
	// The route pages import the definition by alias rather than by relative
	// path. Missing them leaves a project that moved cleanly and then fails to
	// compile, which is worse than not moving at all.
	if err := rewriteAliasImports(filepath.Dir(resourcesRoot)); err != nil {
		return moved, err
	}
	return moved, nil
}

// aliasImportRe matches @/resources/<name> that is not already nested. The
// trailing group refuses a following slash, so a path this has already fixed
// is left alone and running the migration twice is safe.
var aliasImportRe = regexp.MustCompile(`(['"])@/resources/([a-z0-9-]+)(['"])`)

// rewriteAliasImports turns @/resources/products into
// @/resources/products/products across the app, leaving @/resources itself
// (the registry) untouched.
func rewriteAliasImports(appRoot string) error {
	return filepath.WalkDir(appRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if d.Name() == "node_modules" || d.Name() == ".next" || d.Name() == "dist" {
				return filepath.SkipDir
			}
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".ts" && ext != ".tsx" {
			return nil
		}
		body, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		out := aliasImportRe.ReplaceAllString(string(body), "${1}@/resources/${2}/${2}${3}")
		if out == string(body) {
			return nil
		}
		return os.WriteFile(path, []byte(out), 0644)
	})
}

// rewriteRegistryImports turns from "./products" into from "./products/products".
//
// Anchored on the quote so it cannot match a path that is already nested, and
// so "./products" is never rewritten inside "./products/products".
var registryImportRe = regexp.MustCompile(`from "\./([a-z0-9-]+)";`)

func rewriteRegistryImports(indexPath string) error {
	body, err := os.ReadFile(indexPath)
	if err != nil {
		return nil // no registry: a single-resource project may not have one
	}
	out := registryImportRe.ReplaceAllStringFunc(string(body), func(match string) string {
		sub := registryImportRe.FindStringSubmatch(match)
		kebab := sub[1]
		return `from "` + RegistryImportPath(kebab) + `";`
	})
	if out == string(body) {
		return nil
	}
	return os.WriteFile(indexPath, []byte(out), 0644)
}

// isFile answers "is there a file here", under a name that does not collide
// with the identically-shaped helper in internal/generate. That package
// imports this one, so the dependency cannot run the other way and the two
// cannot share.
func isFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
