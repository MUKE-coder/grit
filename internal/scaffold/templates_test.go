package scaffold

import (
	"go/parser"
	"go/token"
	"regexp"
	"strings"
	"testing"
)

// placeholderRe is a scaffold placeholder and not Go syntax. A composite
// literal of structs opens with two braces ([]Attr{{Key: "id"}}), which is why
// this asks for an upper-case identifier between them rather than for "{{".
var placeholderRe = regexp.MustCompile(`\{\{[A-Z_][A-Z0-9_]*\}\}`)

// The point of moving templates into files: they can be checked here, at the
// point of editing, instead of only in CI after a project has been scaffolded.
//
// What these catch is syntax and formatting. A full type check is not possible
// in this repository, because the templates import packages that exist only in
// a generated project; that stays the live suite's job. Syntax and formatting
// are where the bugs came from anyway: a hook after an early return, a stray
// brace, a respond call indented into the wrong block.

// moduleForChecking stands in for {{MODULE}}. Any valid import path works; this
// one is obviously not a real module, so a failure message cannot be mistaken
// for a problem with a user's project.
const moduleForChecking = "example.test/app"

func templateSource(t *testing.T, name string) string {
	t.Helper()
	return strings.ReplaceAll(tmpl(name), "{{MODULE}}", moduleForChecking)
}

// Every template the package asks for by name has to exist, or tmpl panics at
// the moment somebody scaffolds rather than at the moment somebody typos.
func TestEveryTemplateResolves(t *testing.T) {
	names, err := templateNames()
	if err != nil {
		t.Fatalf("walking the templates: %v", err)
	}
	if len(names) == 0 {
		t.Fatal("no templates are embedded, so the embed pattern is wrong")
	}
	for _, name := range names {
		if got := tmpl(name); got == "" {
			t.Errorf("%s is empty", name)
		}
	}
	t.Logf("%d templates embedded", len(names))
}

// A template that produces Go has to parse as Go. This is the check that was
// not possible while the same text lived inside a Go string literal, where a
// missing brace was just a longer string.
func TestGoTemplatesParse(t *testing.T) {
	names, err := templateNames()
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range names {
		if !isGoTemplate(name) {
			continue
		}
		t.Run(name, func(t *testing.T) {
			src := templateSource(t, name)
			fset := token.NewFileSet()
			if _, err := parser.ParseFile(fset, name, src, parser.AllErrors); err != nil {
				t.Errorf("does not parse as Go:\n%v", err)
			}
		})
	}
}

// Deliberately not a gofmt check on the templates.
//
// The first version of this file had one, and it was wrong twice over. The
// scaffolder already runs Go through codefmt on the way out, so a generated
// project is formatted whatever the template looks like, which is why the
// writer's comment says the templates only have to be correct and not
// formatted. And formatting a template means substituting a stand-in for
// {{MODULE}} first, which sorts to a different place in the import block than
// a real module path does: four handlers came out with their imports in a
// different order, so a check meant to tidy the source was quietly changing
// what every project gets.
//
// Syntax is module-independent and is checked above. Formatting is the
// scaffolder's job and it already does it.

// {{MODULE}} is the only placeholder these templates carry. A second one that
// nothing substitutes would ship to a user's project as literal braces, which
// has happened before in this codebase and compiles fine until somebody reads
// it.
func TestTemplatesCarryOnlyTheModulePlaceholder(t *testing.T) {
	names, err := templateNames()
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range names {
		for _, m := range placeholderRe.FindAllString(tmpl(name), -1) {
			if m != "{{MODULE}}" {
				t.Errorf("%s has an unknown placeholder: %s", name, m)
			}
		}
	}
}
