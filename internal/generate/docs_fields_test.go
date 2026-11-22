package generate

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Every --fields string printed in the documentation parses.
//
// Three published guides told a reader to run
//
//	--fields "...,status:enum(draft|sent|paid),..."
//
// and `enum(...)` is not a field type. The command fails immediately, on the
// one page a beginner is most likely to copy from, and nothing caught it:
// prose is not compiled, and the docs site only had to render.
//
// So the docs are read here, and every --fields value in them goes through the
// same parser the CLI uses. A syntax that works only in a sentence is now a
// failing test.
//
// The readable forms are covered too, because they go through the same
// entry point: select[a, b], a trailing `if other == value`, and the
// positional spellings.

// --fields "<value>" or --fields '<value>', with the backslash-newline
// continuations the docs wrap long commands with already removed.
var docsFieldsRe = regexp.MustCompile(`--fields\s+["']([^"']+)["']`)

// Pages that show a deliberately broken command, so that the text can explain
// what the error means. None today; the list exists so that adding one is a
// decision somebody writes down rather than a test somebody deletes.
var docsFieldsExempt = map[string]bool{}

func TestDocsFieldStringsParse(t *testing.T) {
	root := repoRoot(t)
	// app/ is the pages. components/ is where a tool or a showcase writes an
	// example command, and those are read by somebody too.
	roots := []string{
		filepath.Join(root, "docs", "app"),
		filepath.Join(root, "docs", "components"),
	}
	if _, err := os.Stat(roots[0]); err != nil {
		t.Skipf("docs not present: %v", err)
	}

	checked := 0
	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Generated lesson content and the changelog record history,
			// including syntax that has since changed.
			if name := info.Name(); name == "node_modules" || name == "changelog" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".tsx") && !strings.HasSuffix(path, ".ts") {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		rel = filepath.ToSlash(rel)
		if docsFieldsExempt[rel] {
			return nil
		}

		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		// A long command is wrapped with a backslash and a newline in the
		// page source; join those before matching so the whole value is seen.
		text := strings.NewReplacer("\\\\\n", " ", "\\\n", " ").Replace(string(body))

		for _, m := range docsFieldsRe.FindAllStringSubmatch(text, -1) {
			spec := strings.TrimSpace(m[1])
			if isDocsPlaceholder(spec) {
				continue
			}
			checked++
			if _, err := ParseInlineFields("DocsProbe", spec); err != nil {
				t.Errorf("%s: --fields %q does not parse:\n    %v", rel, spec, err)
			}
		}
		return nil
	}
	for _, r := range roots {
		if _, err := os.Stat(r); err != nil {
			continue
		}
		if err := filepath.Walk(r, walk); err != nil {
			t.Fatalf("walk %s: %v", r, err)
		}
	}

	// A walk that silently matched nothing would pass forever. The docs have
	// had far more than this for years; the floor only catches a broken walk.
	if checked < 20 {
		t.Fatalf("only %d --fields strings found in the docs, expected at least 20: "+
			"the walk or the pattern is broken, and this test is no longer checking anything",
			checked)
	}
	t.Logf("%d --fields strings in the docs, all parse", checked)
}

// isDocsPlaceholder reports whether a --fields value stands for a real one
// rather than being one. The LLM-facing pages show the SHAPE of the flag,
// which is their job, and a shape cannot parse.
//
// Deliberately narrow: "type" is not a field type and never will be, and an
// ellipsis is not a field name. Anything else has to parse.
func isDocsPlaceholder(spec string) bool {
	switch {
	case strings.ContainsAny(spec, "<>"):
		return true
	case strings.Contains(spec, "${"):
		return true
	case strings.Contains(spec, "..."):
		return true
	}
	// field:type, field1:type,field2:type: every segment typed "type".
	for _, part := range strings.Split(spec, ",") {
		name, typ, ok := strings.Cut(strings.TrimSpace(part), ":")
		if !ok || typ != "type" || name == "" {
			return false
		}
	}
	return spec != ""
}

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatal("could not find the repository root (no go.mod above the test)")
	return ""
}
