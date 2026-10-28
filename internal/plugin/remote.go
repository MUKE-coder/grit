package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// IsRepoRef reports whether a `grit plugin add` argument names a git
// repository rather than a built-in plugin or a local directory.
//
// A host is required. "user/repo" would be friendlier to type and is not
// accepted, because it is ambiguous with a plugin name and the ambiguity
// resolves in the worst direction: a typo that silently fetches and runs
// somebody else's code is a bad trade for four saved keystrokes.
func IsRepoRef(arg string) bool {
	switch {
	case strings.HasPrefix(arg, "https://"), strings.HasPrefix(arg, "http://"):
		return true
	case strings.HasPrefix(arg, "git@"):
		return true
	case strings.HasPrefix(arg, "github.com/"),
		strings.HasPrefix(arg, "gitlab.com/"),
		strings.HasPrefix(arg, "codeberg.org/"),
		strings.HasPrefix(arg, "bitbucket.org/"):
		return true
	}
	return false
}

// RepoURL turns an argument into something git can clone.
func RepoURL(arg string) string {
	if strings.HasPrefix(arg, "http") || strings.HasPrefix(arg, "git@") {
		return arg
	}
	return "https://" + arg
}

// CloneAndLoad fetches a plugin repository into dir and reads its manifest.
//
// The clone is shallow and the directory is the caller's to remove. Nothing
// from the repository is executed: LoadDir reads plugin.json and the files it
// names, and Install writes them. That is the whole trust surface, and it is
// the same one as a local directory plugin, which is why this can exist at
// all. It is still somebody else's code landing in your project, so the
// caller is expected to have said so out loud first.
func CloneAndLoad(url, dir string) (Plugin, error) {
	cmd := exec.Command("git", "clone", "--depth", "1", "--quiet", url, dir)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return Plugin{}, fmt.Errorf("cloning %s: %w", url, err)
	}

	// A repository may hold the plugin at its root or under plugin/, which is
	// the layout the submission guide asks for when a repo is a whole project
	// that happens to ship a plugin.
	candidates := []string{dir, filepath.Join(dir, "plugin")}
	var lastErr error
	for _, candidate := range candidates {
		p, err := LoadDir(candidate)
		if err == nil {
			return p, nil
		}
		lastErr = err
	}
	return Plugin{}, fmt.Errorf("%s has no plugin.json at its root or in plugin/: %w", url, lastErr)
}

// TempCloneDir makes a directory for a clone, named after the plugin so the
// path in any error message says what is being installed.
func TempCloneDir(url string) (string, error) {
	name := filepath.Base(strings.TrimSuffix(url, ".git"))
	if name == "" || name == "." || name == string(filepath.Separator) {
		name = "plugin"
	}
	return os.MkdirTemp("", "grit-plugin-"+name+"-")
}
