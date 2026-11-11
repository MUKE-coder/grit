package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/errorcodes"
)

// WriteErrorCodeFiles regenerates the two files the error catalogue produces,
// from Grit's codes plus the project's own in .grit/error-codes.yaml.
//
// Its own entry point because the two files have to be rewritten after
// somebody edits that YAML, and the only command that did it was `grit
// upgrade`, which rewrites the whole framework. Nobody runs an upgrade to add
// an error code, so the extension point would have been one people tried once
// (grit#208). `grit sync` calls this.
//
// A hard error on a malformed file, unlike the scaffold and upgrade paths
// which warn and carry on: this is the command whose entire job is to apply
// that file, so finishing with a cheerful message having applied nothing would
// be the worst of the three behaviours.
func WriteErrorCodeFiles(root string) ([]string, error) {
	project, err := errorcodes.LoadProject(root)
	if err != nil {
		return nil, err
	}

	opts := projectOptions(root)
	files := map[string]string{
		filepath.Join(opts.APIRoot(root), "internal", "respond", "codes.go"): apiRespondCodesGo(project),
	}
	// The shared package is absent from an --api project, and its absence is
	// not a failure: the Go half is still worth regenerating.
	for _, shared := range []string{
		filepath.Join(root, "packages", "shared"),
		filepath.Join(root, "src", "shared"),
	} {
		if info, err := os.Stat(shared); err == nil && info.IsDir() {
			files[filepath.Join(shared, "types", "errors.ts")] = sharedErrorsTS(project)
			break
		}
	}

	var written []string
	for path, content := range files {
		rel, _ := filepath.Rel(root, path)
		if err := writeFile(path, strings.ReplaceAll(content, "{{MODULE}}", opts.Module())); err != nil {
			return written, fmt.Errorf("writing %s: %w", rel, err)
		}
		written = append(written, filepath.ToSlash(rel))
	}

	if len(project) > 0 {
		written = append(written, fmt.Sprintf("%d code(s) from %s", len(project), errorcodes.ProjectFile))
	}
	return written, nil
}
