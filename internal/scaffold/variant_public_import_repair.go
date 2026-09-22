package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// grit add variants wrote a public variants handler that called respond.Fail
// without importing respond, so the API stopped compiling the moment variants
// were added to a resource. Nothing compiled that file: the fresh-project
// checks never run grit add variants. Found building the storefront blueprint.

// modelsImportLine is the models import in a generated handler, whose module
// path the respond import copies.
var modelsImportLine = regexp.MustCompile(`(?m)^\t"([^"]+)/internal/models"\n`)

func repairVariantPublicImportSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "respond.Fail(") || strings.Contains(src, `/internal/respond"`) {
		return src, nil, nil
	}
	m := modelsImportLine.FindStringSubmatchIndex(src)
	if m == nil {
		return src, nil, []string{"the public variants handler calls respond without importing it: add the internal/respond import, or the API does not compile"}
	}
	module := src[m[2]:m[3]]
	out := src[:m[1]] + "\t\"" + module + "/internal/respond\"\n" + src[m[1]:]
	return out, []string{"the public variants handler imports respond, so the API compiles"}, nil
}

// repairVariantPublicImport applies it to every resource that has variants.
func repairVariantPublicImport(root string) error {
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	for _, apiRoot := range []string{filepath.Join(root, "apps", "api"), filepath.Join(root, "api"), root} {
		matches, err := filepath.Glob(filepath.Join(apiRoot, "internal", "handlers", "*_variant_public.go"))
		if err != nil {
			return err
		}
		for _, path := range matches {
			if _, err := os.Stat(path); err != nil {
				continue
			}
			if err := repairTextFile(root, m, path, repairVariantPublicImportSource); err != nil {
				return err
			}
		}
	}
	return nil
}
