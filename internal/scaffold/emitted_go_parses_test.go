package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/codefmt"
)

// Does every Go file the scaffold writes parse?
//
// The same gap as grit#159 on the generator side, over far more templates: the
// API is hundreds of raw Go strings, asserted on with strings.Contains, and a
// Contains that passes says nothing about whether the file is valid Go. CI does
// compile a scaffolded project, which catches these, ten minutes later and with
// a failure pointing at a generated project rather than at the template.
//
// Four architectures, because a template only one shape reaches is exactly the
// one nothing else checks. A single keeps the Go module at the root, which
// changes every import path the templates emit.
func TestEveryGoFileTheScaffoldWritesParses(t *testing.T) {
	shapes := []struct {
		name string
		opts Options
	}{
		{"triple-next", Options{ProjectName: "tripleapp", Architecture: ArchTriple, Frontend: FrontendNext}},
		{"double-next", Options{ProjectName: "doubleapp", Architecture: ArchDouble, Frontend: FrontendNext}},
		{"single-vite", Options{ProjectName: "singleapp", Architecture: ArchSingle, Frontend: FrontendTanStack}},
		{"api-only", Options{ProjectName: "apiapp", Architecture: ArchAPI}},
	}

	for _, shape := range shapes {
		t.Run(shape.name, func(t *testing.T) {
			root := t.TempDir()
			if err := createDirectories(root, shape.opts); err != nil {
				t.Fatalf("createDirectories: %v", err)
			}
			if err := writeAPIFiles(root, shape.opts); err != nil {
				t.Fatalf("writeAPIFiles: %v", err)
			}

			checked := parseEveryGoFile(t, root)

			// A sweep that found nothing is indistinguishable from a sweep that
			// ran over nothing, and this repository has shipped the second kind.
			if checked < 50 {
				t.Fatalf("only %d Go files under %s; this test is not reading what it "+
					"thinks it is", checked, root)
			}
			t.Logf("parsed %d Go files", checked)
		})
	}
}

// parseEveryGoFile parses every .go file under root, reporting each failure, and
// returns how many it read.
func parseEveryGoFile(t *testing.T, root string) int {
	t.Helper()
	checked := 0
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		rel, _ := filepath.Rel(root, path)
		mustParseGoFile(t, filepath.ToSlash(rel), string(data))
		checked++
		return nil
	})
	if err != nil {
		t.Fatalf("walking %s: %v", root, err)
	}
	return checked
}

// mustParseGoFile fails the test when src is not valid Go, with the line shown.
//
// Through codefmt.GoChecked rather than parser.ParseFile, so this and the
// formatter agree on what parsing means. codefmt.Go, which every write goes
// through, discards exactly this error: right for a write, because a file on
// disk somebody can read beats a scaffolding failure with nothing written, and
// wrong for a test.
func mustParseGoFile(t *testing.T, filename, src string) {
	t.Helper()
	if _, err := codefmt.GoChecked(src); err != nil {
		t.Errorf("%s is not valid Go: %v\n\n%s", filename, err, goErrorWindow(src, err))
	}
}

// goErrorWindow prints the lines around the error, numbered. The difference
// between "syntax error at 41:9" and knowing what is on line 41.
func goErrorWindow(src string, err error) string {
	line := 0
	if parts := strings.SplitN(err.Error(), ":", 2); len(parts) == 2 {
		_, _ = fmt.Sscanf(parts[0], "%d", &line)
	}
	lines := strings.Split(src, "\n")
	from, to := 0, len(lines)
	if line > 0 {
		from, to = max(0, line-8), min(len(lines), line+4)
	}
	var b strings.Builder
	for i := from; i < to; i++ {
		marker := "  "
		if i+1 == line {
			marker = "->"
		}
		fmt.Fprintf(&b, "%s %4d | %s\n", marker, i+1, lines[i])
	}
	return b.String()
}
