package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// No model Grit writes may carry a bool with a database default of true.
//
// GORM omits a zero-valued field from an INSERT when the column has a default,
// so the database fills in true and GORM writes it back: the row is stored
// ticked and nothing errors. This project has now found it three times in its
// own models, on User.Active, on the SSO flags, and on a product variant that
// the first two fixes missed (grit#179).
//
// A sweep rather than a note in a comment, because the two earlier fixes were
// notes in comments and the third one still shipped.
//
// It reads the TEMPLATE SOURCE rather than a generated project, and that is
// not laziness. The first version of this test scaffolded a project and walked
// internal/models, and it passed while the bug was still in the tree: the
// variant model is written by `grit generate resource --variants`, through
// internal/generate, so the one file the fault was actually in was the one the
// sweep could not see. Reading the templates covers every writer there is.
func TestNoScaffoldedModelHasABoolDefaultingToTrue(t *testing.T) {
	// A bool struct field whose gorm tag carries default:true, as it appears
	// inside a template string. Not *bool: nil and false are different values
	// to GORM, so a pointer has no zero-value problem and a default there is
	// legitimate.
	//
	// The tag is spliced as ` + "`" + ` in a Go raw string, so the pattern
	// allows anything between the type and the gorm key rather than a literal
	// backtick.
	offending := regexp.MustCompile(`(?m)^\s*(\w+)\s+bool\s+[^\n]*gorm:\\?"[^"\n]*default:(?i:true)`)

	roots := []string{".", filepath.Join("..", "generate")}
	checked, hits := 0, 0

	for _, root := range roots {
		entries, err := os.ReadDir(root)
		if err != nil {
			t.Fatalf("reading %s: %v", root, err)
		}
		for _, e := range entries {
			name := e.Name()
			if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
				continue
			}
			data, err := os.ReadFile(filepath.Join(root, name))
			if err != nil {
				t.Fatal(err)
			}
			checked++
			for _, m := range offending.FindAllStringSubmatch(string(data), -1) {
				hits++
				t.Errorf("%s: a template declares %s as a bool with gorm default:true, "+
					"so it can never be saved as false.\n\n"+
					"GORM leaves a zero value out of the INSERT, the database fills in "+
					"true, and GORM writes that back into the struct: saving \"unticked\" "+
					"stores \"ticked\" with no error. Drop the default and set it on every "+
					"create path, as User.Active does, or make it *bool.",
					filepath.Join(root, name), m[1])
			}
		}
	}

	if checked < 100 {
		t.Fatalf("only %d template files read; this sweep is not covering what it "+
			"thinks it is", checked)
	}
	if hits == 0 {
		t.Logf("swept %d template files", checked)
	}
}
