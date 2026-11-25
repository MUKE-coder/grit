package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Per-resource retention (grit#260): internal/retention is the registry a
// model declares its window in, and services/retention.go is the sweep that
// enforces it.
//
// grit new and grit upgrade deliver all four files through
// writeFrameworkOwnedFiles, beside the Trash the sweep is built on.
// WriteRetentionPackage is the generator's: the registry is all a generated
// model needs to compile, so it writes it into a project that predates it, the
// way it writes internal/erasure.

func apiRetentionGo() string            { return tmpl("api/retention/retention.go") }
func apiRetentionTestGo() string        { return tmpl("api/retention/retention_test.go") }
func apiRetentionServiceGo() string     { return tmpl("api/services/retention.go") }
func apiRetentionServiceTestGo() string { return tmpl("api/services/retention_test.go") }

// WriteRetentionPackage writes internal/retention under apiRoot. overwrite=false
// leaves an existing copy alone, for the generator on an older project.
func WriteRetentionPackage(apiRoot, module string, overwrite bool) error {
	files := map[string]string{
		filepath.Join(apiRoot, "internal", "retention", "retention.go"):      apiRetentionGo(),
		filepath.Join(apiRoot, "internal", "retention", "retention_test.go"): apiRetentionTestGo(),
	}
	for path, content := range files {
		if !overwrite {
			if _, err := os.Stat(path); err == nil {
				continue
			}
		}
		if err := writeFile(path, strings.ReplaceAll(content, "{{MODULE}}", module)); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}
