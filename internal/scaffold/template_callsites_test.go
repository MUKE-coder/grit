package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// tmpl panics on a name it cannot find, and that panic lands in `grit new`:
// the scaffolder stops partway, leaving a half-written project and a stack
// trace. It is the worst failure this package has, because it happens to a
// user on their first command.
//
// TestEveryTemplateResolves walks the templates directory and checks each file
// loads. That is the wrong direction. Every file resolving says nothing about
// whether the calls name those files: `tmpl("admin/pages/security")` for a
// file called security.tsx.tmpl passes that test and panics on the first
// `grit new --triple`, which is exactly what it did.
//
// This checks the other direction: every name passed to tmpl() is a template
// that exists.
func TestEveryTmplCallNamesATemplateThatExists(t *testing.T) {
	calls, err := tmplCallSites()
	if err != nil {
		t.Fatalf("reading the package source: %v", err)
	}
	if len(calls) == 0 {
		t.Fatal("no tmpl() calls found, so this test is checking nothing")
	}

	have := map[string]bool{}
	names, err := templateNames()
	if err != nil {
		t.Fatalf("walking the templates: %v", err)
	}
	for _, name := range names {
		have[name] = true
	}

	for _, call := range calls {
		if have[call.name] {
			continue
		}
		hint := ""
		// The mistake is nearly always a missing extension, because the file
		// is foo.tsx.tmpl and the obvious name is "foo".
		for candidate := range have {
			if strings.HasPrefix(candidate, call.name+".") {
				hint = ", did you mean " + candidate + "?"
				break
			}
		}
		t.Errorf("%s:%d: tmpl(%q) names no template%s. This panics in grit new, "+
			"part way through writing the project", call.file, call.line, call.name, hint)
	}
	t.Logf("%d tmpl() call sites, all resolved", len(calls))
}

// And the reverse, as a tidiness check rather than a correctness one: a
// template nothing asks for is either a file somebody forgot to wire up or
// dead weight in the binary.
func TestEveryTemplateIsAskedFor(t *testing.T) {
	calls, err := tmplCallSites()
	if err != nil {
		t.Fatal(err)
	}
	asked := map[string]bool{}
	for _, c := range calls {
		asked[c.name] = true
	}
	names, err := templateNames()
	if err != nil {
		t.Fatal(err)
	}
	var orphans []string
	for _, name := range names {
		if !asked[name] {
			orphans = append(orphans, name)
		}
	}
	sort.Strings(orphans)
	if len(orphans) > 0 {
		t.Errorf("%d template(s) that nothing calls tmpl() for, so they are "+
			"written to no project:\n  %s", len(orphans), strings.Join(orphans, "\n  "))
	}
}

type tmplCall struct {
	file string
	line int
	name string
}

// Only a plain string literal is read. A name built at runtime cannot be
// checked here, and there are none today; if one appears, this finds the
// others and leaves it alone rather than failing on it.
var tmplCallRe = regexp.MustCompile(`\btmpl\(\s*"([^"]+)"\s*\)`)

func tmplCallSites() ([]tmplCall, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}
	var calls []tmplCall
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(".", name))
		if err != nil {
			return nil, err
		}
		for i, line := range strings.Split(string(raw), "\n") {
			for _, m := range tmplCallRe.FindAllStringSubmatch(line, -1) {
				calls = append(calls, tmplCall{file: name, line: i + 1, name: m[1]})
			}
		}
	}
	return calls, nil
}
