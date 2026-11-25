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
// Two writers because the two halves travel differently. The package is all a
// generated model needs to compile, so the generator writes it into a project
// that predates it, the way it writes internal/erasure. The sweep needs the
// job and the cron entry beside it, which only grit new and grit upgrade
// deliver.

func apiRetentionGo() string            { return tmpl("api/retention/retention.go") }
func apiRetentionTestGo() string        { return tmpl("api/retention/retention_test.go") }
func apiRetentionServiceGo() string     { return tmpl("api/services/retention.go") }
func apiRetentionServiceTestGo() string { return tmpl("api/services/retention_test.go") }

// WriteRetentionPackage writes internal/retention under apiRoot. overwrite=false
// leaves an existing copy alone, for the generator on an older project.
func WriteRetentionPackage(apiRoot, module string, overwrite bool) error {
	return writeModuleFiles(module, overwrite, map[string]string{
		filepath.Join(apiRoot, "internal", "retention", "retention.go"):      apiRetentionGo(),
		filepath.Join(apiRoot, "internal", "retention", "retention_test.go"): apiRetentionTestGo(),
	})
}

// writeRetentionService writes the sweep and its test. They rely on
// TrashService.purge, so they belong to the same delivery as trash.go.
func writeRetentionService(apiRoot, module string) error {
	return writeModuleFiles(module, true, map[string]string{
		filepath.Join(apiRoot, "internal", "services", "retention.go"):      apiRetentionServiceGo(),
		filepath.Join(apiRoot, "internal", "services", "retention_test.go"): apiRetentionServiceTestGo(),
	})
}

func writeModuleFiles(module string, overwrite bool, files map[string]string) error {
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
