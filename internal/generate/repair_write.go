package generate

import (
	"fmt"
	"go/format"
	"os"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Review 2.5, the half that lives in this package.
//
// The repairs here edit a project's generated services, handlers and routes
// after an upgrade, and they wrote with os.WriteFile and told the manifest
// nothing. The cost is the same as it was in internal/scaffold: a file a repair
// edited read as hand-edited afterwards, so every later framework fix to it was
// skipped. The repair poisoned the file it fixed.
//
// TestRepairsDoNotWriteBehindTheManifest, in internal/scaffold, scans this
// package too.

// writeRepairedGo formats an edited Go file, writes it, and tells the manifest.
//
// Formatted, because these repairs splice text in and a repair that leaves a
// file gofmt-dirty fails the check CI runs over a generated project. A file
// that does not parse is left alone and reported rather than written: a repair
// that produces broken Go should not be the reason a project stops compiling.
//
// Refresh keeps the file reading as Grit's when it was Grit's, so later fixes
// still reach it, and leaves it reading as edited when somebody had edited it,
// so the next upgrade does not overwrite their work. Which of the two it does
// is decided inside the manifest, from a snapshot taken before this command
// wrote anything (review 2.6).
func writeRepairedGo(path, content string) error {
	formatted, err := format.Source([]byte(content))
	if err != nil {
		return fmt.Errorf("the repair to %s did not produce valid Go, so it was left alone: %w", path, err)
	}
	if err := os.WriteFile(path, formatted, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	manifest.Refresh(path)
	return nil
}

// writeRepairedFile is writeRepairedGo for a file that is not Go.
func writeRepairedFile(path, content string) error {
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	manifest.Refresh(path)
	return nil
}
