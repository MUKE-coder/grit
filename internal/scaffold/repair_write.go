package scaffold

import (
	"fmt"
	"os"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Review 2.5: repairs were not uniformly guarded by the manifest.
//
// grit upgrade runs around a hundred repair passes. Most go through
// repairSourceFile, which checks the manifest before writing and refreshes it
// after. Eleven wrote with os.WriteFile directly and told the manifest
// nothing, which has a cost in both directions:
//
//   - a file the repair edited read as hand-edited afterwards, so every later
//     framework fix to it was skipped. The repair poisoned the file it fixed.
//   - nothing consulted the upgrade guard, so a repair could write over
//     something the guard would have refused.
//
// Two helpers, because there are two cases, and using the wrong one is the
// mistake worth making hard:
//
//   - writeRepaired, for editing a file that is already there
//   - writeFile, for a whole framework-owned file arriving for the first time
//
// TestRepairsDoNotWriteBehindTheManifest keeps new repairs on one of them.

// writeRepaired writes an edited file and tells the manifest what happened.
//
// Refresh keeps the file reading as Grit's when it was Grit's, so later fixes
// still reach it, and leaves it reading as edited when somebody had edited it,
// so the next upgrade does not overwrite their work. Which of those two it
// does is decided inside the manifest, from a snapshot taken before this
// command wrote anything (review 2.6).
func writeRepaired(path, content string) error {
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	manifest.Refresh(path)
	return nil
}

// writeRepairedMode is writeRepaired for a file whose permissions matter, such
// as a shell script a repair appends to.
func writeRepairedMode(path, content string, mode os.FileMode) error {
	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	manifest.Refresh(path)
	return nil
}
