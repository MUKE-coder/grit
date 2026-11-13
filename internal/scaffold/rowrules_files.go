package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// writeRowRules writes internal/rowrules/rules.go, and only when it is not
// already there.
//
// The engine beside it travels with every upgrade, because nothing in it is
// meant to be edited. This file is the opposite: it holds the project's own
// access rules, and `grit generate resource` injects into it. Overwriting it
// on upgrade would delete somebody's policies, and the manifest guard would
// report a conflict on every upgrade of every project that had written one,
// which teaches people to ignore conflicts.
//
// So: created once, and left alone afterwards. A project that has never had
// it gets it on the next upgrade, which is how the feature reaches a project
// scaffolded before it existed.
func writeRowRules(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	path := filepath.Join(apiRoot, "internal", "rowrules", "rules.go")

	if _, err := os.Stat(path); err == nil {
		return nil // the project owns it
	}

	content := strings.ReplaceAll(tmpl("api/rowrules/rules.go"), "{{MODULE}}", opts.Module())
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating the rowrules directory: %w", err)
	}
	return writeFile(path, content)
}
