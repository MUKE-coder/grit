package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// repairCIWorkflows brings a project's CI up to the fix for H21 in the
// contact-app review: Dependabot and the security scans pointed at the
// single-app layout whatever the project was, nothing ran the tests on a pull
// request, and the release workflow was rejected by GitHub on the first tag.
//
// Grit's own workflow files are rewritten where they are still what Grit wrote;
// the upgrade guard keeps an edited one. ci.yml is new, so it is written only
// when the project has no file by that name: an untracked file is one the guard
// would otherwise replace.
func repairCIWorkflows(root string, opts Options) error {
	workflows := filepath.Join(root, ".github", "workflows")
	if !dirExists(workflows) {
		return nil
	}
	ci := opts
	// Upgrade reads the shape from the directories, and a desktop app is one of
	// them it does not record.
	ci.IncludeDesktop = ci.IncludeDesktop || dirExists(filepath.Join(root, "apps", "desktop"))

	files := map[string]string{
		filepath.Join(root, ".github", "dependabot.yml"): dependabotYAML(ci),
		filepath.Join(workflows, "security.yml"):         securityCIYAML(ci),
		filepath.Join(workflows, "release.yml"):          releaseCIYAML(ci),
		filepath.Join(workflows, "lint.yml"):             lintCIYAML(ci),
	}
	if path := filepath.Join(workflows, "ci.yml"); !fileExists(path) {
		files[path] = ciYAML(ci)
	}
	for path, content := range files {
		if !fileExists(path) && !strings.HasSuffix(path, "ci.yml") {
			// A workflow the project deleted stays deleted.
			continue
		}
		if err := writeFile(path, strings.ReplaceAll(content, "{{PROJECT}}", opts.ProjectName)); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}
