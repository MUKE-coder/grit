package scaffold

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// stampProjectVersion records the version an upgrade just applied, in
// grit.json.
//
// The field was written once at scaffold time and never touched again, so a
// project upgraded four times still reported the version it was created with.
// That is the field people read to answer "what version am I on", and it was
// the wrong answer every time after the first upgrade.
//
// The real record lives in .grit/manifest.json, which the manifest keeps
// current. This is not a second source of truth so much as the human-facing
// copy of the same fact, and the two disagreeing is worse than either.
//
// Only the version key changes. The file is decoded into a map and re-encoded
// rather than rewritten from the template, because architecture, frontend and
// apps describe how this project was built and are not the upgrade's business.
func stampProjectVersion(root, version string) error {
	path := filepath.Join(root, "grit.json")
	body, err := os.ReadFile(path)
	if err != nil {
		// No grit.json means this is not a project the upgrade can label, and
		// findProjectRoot would have failed long before here.
		return nil
	}

	var project map[string]any
	if err := json.Unmarshal(body, &project); err != nil {
		// Hand-edited into something that is no longer JSON. Leaving it alone
		// beats replacing it with a guess.
		return nil
	}
	if project["version"] == version {
		return nil
	}
	project["version"] = version

	// Indented to match what the scaffold writes, so an upgrade does not turn
	// the file into one long line in somebody's diff.
	out, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding grit.json: %w", err)
	}
	return os.WriteFile(path, append(out, '\n'), 0o644)
}
