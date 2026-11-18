package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/project"
)

// bindFailureOld is what a handler scaffolded before v3.546.0 answers a
// binding failure with.
//
// An exact string, and only ever written by Grit's own templates: nobody
// composes that line by hand. So replacing it is a substitution rather than a
// rewrite, which is what makes this safe to run over a developer's handler.
const bindFailureOld = "respond.Fail(c, respond.CodeValidationError, err.Error())"

// bindFailureNew fills error.details instead of sending the validator's own
// sentence. See respond.Invalid, and grit#225.
const bindFailureNew = "respond.Invalid(c, err)"

// repairBindFailures moves a project's handlers off err.Error() for a binding
// failure.
//
// Why a repair at all: handlers are the developer's after scaffolding, and
// upgrade preserves them, which is right. So converting the 53 templates
// reached new projects and nobody else. Every existing project went on
// answering a failed validation with
//
//	Key: 'RegisterRequest.first_name' Error:Field validation for 'first_name'
//	failed on the 'required' tag
//
// and an empty error.details, which is the field the error documentation tells
// clients to read.
//
// It runs last in the upgrade and has its own gate, because a check folded
// into an existing repair only reaches projects that repair has not already
// visited.
func RepairBindFailures(root string) error {
	apiRoot := project.APIRoot(root)

	// The gate: respond.Invalid has to exist before anything calls it.
	// writeRespondFiles delivers it earlier in the same upgrade, so this is a
	// check on the order rather than on the version.
	if !fileContains(filepath.Join(apiRoot, "internal", "respond", "respond.go"), "func Invalid(") {
		return nil
	}

	// Handlers hold almost all of them, and a project may have put one
	// elsewhere, so the whole Go module is walked rather than a fixed list.
	// Excluding internal/respond itself, whose own doc comment quotes the old
	// line while explaining why not to write it.
	var converted []string
	walkErr := filepath.Walk(apiRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// vendor and node_modules are somebody else's code and large.
			if name := info.Name(); name == "vendor" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".go") || strings.HasSuffix(info.Name(), "_test.go") {
			return nil
		}
		if strings.Contains(filepath.ToSlash(path), "/internal/respond/") {
			return nil
		}

		body, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		src := string(body)
		if !strings.Contains(src, bindFailureOld) {
			return nil
		}

		n := strings.Count(src, bindFailureOld)
		if writeErr := writeEdit(path, strings.ReplaceAll(src, bindFailureOld, bindFailureNew)); writeErr != nil {
			return writeErr
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			rel = path
		}
		converted = append(converted, fmt.Sprintf("%s (%d)", filepath.ToSlash(rel), n))
		return nil
	})
	if walkErr != nil {
		return walkErr
	}

	for _, shown := range converted {
		fmt.Printf("  ✓ %s: a failed validation fills error.details now, instead of "+
			"answering with the validator's own sentence\n", shown)
	}
	return nil
}
