package scaffold

// Raising @hookform/resolvers in a project that cannot build its admin.
//
// FormBuilder imports @hookform/resolvers/standard-schema, a subpath that
// arrived in resolvers 4. The TanStack admin's package.json kept ^3.9.1 from
// v3.475.0 to v3.484.0, so every --vite project generated in that window fails
// to build its admin outright:
//
//	[vite]: Rollup failed to resolve import "@hookform/resolvers/standard-schema"
//	from src/components/forms/form-builder.tsx
//
// Rewriting package.json from the template is not enough. A project that added
// one dependency of its own has edited the file, the upgrade guard holds it back
// for exactly the right reason, and the admin stays unbuildable. So this edits
// the one line, which is a change no reader will mind and which they cannot make
// for themselves without knowing why.

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// resolversEntryRe is the dependency line, with its indentation and its comma
// kept so the rewrite changes nothing else.
var resolversEntryRe = regexp.MustCompile(`("@hookform/resolvers"\s*:\s*")([^"]+)(")`)

// raiseHookformResolvers brings every admin package.json that ships FormBuilder
// up to a resolvers version that has the subpath.
func raiseHookformResolvers(root string, opts Options) error {
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	want, err := majorOfPin(HookformResolversVersion)
	if err != nil {
		return nil
	}

	for _, path := range []string{
		filepath.Join(root, "apps", "admin", "package.json"),
		filepath.Join(webAppRoot(root, opts), "package.json"),
		filepath.Join(root, "frontend", "package.json"),
	} {
		raw, readErr := os.ReadFile(path)
		if readErr != nil {
			continue
		}
		src := string(raw)
		// Only an app that actually imports the subpath. Raising a major in a
		// package.json for a component it does not have is a change nobody
		// asked for.
		if !appUsesStandardSchemaResolver(filepath.Dir(path)) {
			continue
		}
		match := resolversEntryRe.FindStringSubmatchIndex(src)
		if match == nil {
			continue
		}
		have := src[match[4]:match[5]]
		if major, convErr := majorOfPin(have); convErr != nil || major >= want {
			continue
		}

		// A package.json nobody has edited stays Grit's after this. Without the
		// refresh the next upgrade reads this edit as the reader's and stops
		// updating the file.
		pristine := false
		if key, inside := manifest.Rel(root, path); inside {
			pristine = m.StatusOf(root, key) == manifest.Unchanged
		}
		out := src[:match[4]] + HookformResolversVersion + src[match[5]:]
		if err := writeRepaired(path, out); err != nil {
			return err
		}
		if pristine {
			manifest.Refresh(path)
		}
		rel, _ := filepath.Rel(root, path)
		fmt.Printf("  ✓ %s: @hookform/resolvers %s to %s, the version with the "+
			"standard-schema subpath the form builder imports; without it the app does not build\n",
			filepath.ToSlash(rel), have, HookformResolversVersion)
	}
	return nil
}

// appUsesStandardSchemaResolver reports whether this app has the component that
// needs the subpath.
func appUsesStandardSchemaResolver(appRoot string) bool {
	for _, candidate := range []string{
		filepath.Join(appRoot, "components", "forms", "form-builder.tsx"),
		filepath.Join(appRoot, "src", "components", "forms", "form-builder.tsx"),
	} {
		body, err := os.ReadFile(candidate)
		if err != nil {
			continue
		}
		if strings.Contains(string(body), "@hookform/resolvers/standard-schema") {
			return true
		}
	}
	return false
}

// majorOfPin reads the major out of a pin like ^5.9.0.
func majorOfPin(pin string) (int, error) {
	head, _, _ := strings.Cut(strings.TrimLeft(pin, "^~>=< v"), ".")
	return strconv.Atoi(head)
}
