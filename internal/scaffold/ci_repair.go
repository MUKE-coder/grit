package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
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
	// ci.yml is delivered when the project has none, or when Grit wrote the one
	// it has. A ci.yml the project wrote itself is untracked, and the upgrade
	// guard would replace it, so it is left alone.
	if path := filepath.Join(workflows, "ci.yml"); !fileExists(path) || gritWroteFile(root, path) {
		files[path] = ciYAML(ci)
	}
	// The allowlist security.yml reads. New, like ci.yml, and kept when the
	// project has one of its own.
	allow := filepath.Join(root, ".github", "govulncheck-allow.txt")
	if fileExists(filepath.Join(workflows, "security.yml")) && (!fileExists(allow) || gritWroteFile(root, allow)) {
		files[allow] = govulncheckAllowTXT()
	}
	for path, content := range files {
		if !fileExists(path) && !strings.HasSuffix(path, "ci.yml") && path != allow {
			// A workflow the project deleted stays deleted.
			continue
		}
		if err := writeFile(path, strings.ReplaceAll(content, "{{PROJECT}}", opts.ProjectName)); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

// gritWroteFile reports whether the manifest records Grit writing path. The
// upgrade guard then keeps the file if it has been edited since.
func gritWroteFile(root, path string) bool {
	m, err := manifest.Load(root)
	if err != nil {
		return false
	}
	rel, ok := manifest.Rel(root, path)
	return ok && m.StatusOf(root, rel) != manifest.Untracked
}
