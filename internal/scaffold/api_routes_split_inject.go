package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ensureRouteRegistry gives an existing project the per-resource routes split.
//
// Writing resources.go is free, because it is a new file. The other half is
// not: routes.go has to call mountResources, and routes.go is exactly the kind
// of file `grit upgrade` refuses to rewrite, because it is the file people
// edit most. So the call is injected, once, against a marker that has been in
// every scaffolded routes.go for many versions.
//
// Until this runs, `grit generate resource` in an upgraded project falls back
// to the old marker injections and keeps working. After it runs, new resources
// get their own file and the ones already inline stay where they are, which is
// the only safe answer: moving somebody's routes out from under them would
// silently drop any edit they had made to those lines.
//
// A no-op on a project that already has it, so running upgrade twice is safe.
func ensureRouteRegistry(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	routesDir := filepath.Join(apiRoot, "internal", "routes")
	if _, err := os.Stat(routesDir); err != nil {
		return nil
	}

	registry := strings.ReplaceAll(apiRoutesRegistryGo(), "{{MODULE}}", opts.Module())
	if err := writeFile(filepath.Join(routesDir, "resources.go"), registry); err != nil {
		return fmt.Errorf("writing resources.go: %w", err)
	}

	return ensureMountResourcesCall(filepath.Join(routesDir, "routes.go"))
}

func ensureMountResourcesCall(path string) error {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(body)
	if strings.Contains(content, "mountResources(") {
		return nil
	}

	// Anchored on the legacy-alias call rather than on grit:routes:custom.
	// Both sit at the end of Setup where every group is still in scope, but
	// this one is a statement rather than a comment, so a project whose
	// markers have been tidied away still matches.
	const anchor = "\tmountLegacyAPIAlias(r)"
	if !strings.Contains(content, anchor) {
		fmt.Println("  ⚠ routes.go looks hand-edited; add the mountResources call by hand")
		fmt.Println("    See internal/routes/resources.go for the Mount fields it takes.")
		return nil
	}

	// Only the groups this routes.go actually declares.
	//
	// A project from before the public API group has no publicAPI variable,
	// and naming it here is an undefined-symbol error, not a nil field. A
	// resource that then mounts on a group its project does not have fails its
	// own nil dereference at boot, loudly, which beats the whole API refusing
	// to compile over a group nothing uses.
	fields := []string{
		"\t\tEngine:    r,",
		"\t\tDB:        db,",
	}
	for _, g := range []struct{ field, decl string }{
		{"Cfg:       cfg,", "cfg *config.Config"},
		{"Svc:       svc,", "svc *Services"},
		{"V1:        v1,", "v1 :="},
		{"Public:    publicAPI,", "publicAPI :="},
		{"Protected: protected,", "protected :="},
		{"Admin:     admin,", "admin :="},
	} {
		if strings.Contains(content, g.decl) {
			fields = append(fields, "\t\t"+g.field)
		}
	}

	call := "\t// Every generated resource, each from its own <resource>_routes.go.\n" +
		"\t//\n" +
		"\t// A resource file registers itself from an init(), so this loop is the\n" +
		"\t// only place routes.go mentions them. Adding a resource does not edit\n" +
		"\t// this file, and neither does removing one.\n" +
		"\tmountResources(&Mount{\n" +
		strings.Join(fields, "\n") + "\n" +
		"\t})\n\n"

	content = strings.Replace(content, anchor, call+anchor, 1)
	return os.WriteFile(path, []byte(content), 0o644)
}
