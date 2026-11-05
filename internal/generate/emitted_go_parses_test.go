package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/codefmt"
)

// Does the Go this generator writes actually parse?
//
// The templates are raw Go strings and the tests assert on them with around
// 1,600 strings.Contains calls. A Contains that passes says the substring is
// there. It says nothing about whether the file is valid Go, so a template that
// emits a syntax error passes every test in this package and is found by CI
// compiling a scaffolded project, ten minutes later, with a failure pointing at
// a generated project rather than at the template (grit#159).
//
// Two resources, not one, and that is the point of the second. Every test in
// this package generated a single resource, so a template that declares a
// package-level function without the resource name in it looked fine: the
// redeclaration only appears on the second resource. That shipped in v3.451.0
// and nothing here could see it.

// mustParseGo fails the test when src is not valid Go, with the line numbered.
//
// Through codefmt.GoChecked rather than parser.ParseFile directly, so that this
// and the formatter agree on what parsing means: if gofmt can format it, it
// compiles far enough for the template to be the generator's problem rather
// than a syntax error.
func mustParseGo(t *testing.T, filename, src string) {
	t.Helper()
	if _, err := codefmt.GoChecked(src); err != nil {
		t.Errorf("%s is not valid Go: %v\n\n%s", filename, err, numberedSource(src, err))
	}
}

// numberedSource prints the source with line numbers, which is the difference
// between "syntax error at 41:9" and knowing what is on line 41.
func numberedSource(src string, err error) string {
	// A format.Source error reads "41:9: expected ';', found ...". Show a
	// window around it rather than the whole file, which for a service is 600
	// lines of noise around one bad one.
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

func TestEveryGoFileTheGeneratorWritesParses(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	// Two resources. The second is what catches a package-level declaration
	// that is not named after its resource.
	for _, spec := range []struct{ name, fields string }{
		{"Post", "title:string,content:text,published:bool,price:float"},
		{"Comment", "body:text,author:string,approved:bool"},
	} {
		def, err := ParseInlineFields(spec.name, spec.fields)
		if err != nil {
			t.Fatalf("ParseInlineFields(%s): %v", spec.name, err)
		}
		if err := newTestGenerator(root, module, def).Run(); err != nil {
			t.Fatalf("Generator.Run() for %s: %v", spec.name, err)
		}
	}

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
		mustParseGo(t, filepath.ToSlash(rel), string(data))
		checked++
		return nil
	})
	if err != nil {
		t.Fatalf("walking the project: %v", err)
	}

	// A sweep that found nothing is indistinguishable from a sweep that ran
	// over nothing, and this repository has shipped the second kind.
	if checked < 6 {
		t.Fatalf("only %d Go files found under %s; this test is not reading what it "+
			"thinks it is", checked, root)
	}
	t.Logf("parsed %d Go files", checked)
}

// And the helper has to fail on bad Go, or the sweep above is decoration.
func TestMustParseGoRejectsBrokenGo(t *testing.T) {
	broken := "package models\n\nfunc Thing( {\n"
	if _, err := codefmt.GoChecked(broken); err == nil {
		t.Fatal("codefmt.GoChecked accepted source that does not parse, so every " +
			"test built on it is passing for the wrong reason")
	}

	// And the window around the error names the offending line.
	_, err := codefmt.GoChecked(broken)
	out := numberedSource(broken, err)
	if !strings.Contains(out, "func Thing( {") {
		t.Errorf("the error window does not show the broken line:\n%s", out)
	}
}
