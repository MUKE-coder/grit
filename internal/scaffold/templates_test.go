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

// Every template is LF, on every machine.
//
// Git converts text files on checkout unless told not to, so a clone on Windows
// turned all 650 line endings in auth.go.tmpl into CRLF. The CLI then embedded
// CRLF, the repair constants stopped matching the templates they mirror, and
// every generated project differed from the one CI verified. CI runs on Linux
// and saw none of it.
//
// .gitattributes pins these to LF. This is the check that says so, because an
// attributes file is easy to lose in a merge and the failure it prevents is
// invisible on the machine most of the testing happens on.
func TestTemplatesAreLF(t *testing.T) {
	names, err := templateNames()
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range names {
		if strings.Contains(tmpl(name), "\r\n") {
			t.Errorf("%s has CRLF line endings: check .gitattributes, then "+
				"git add --renormalize internal/scaffold/templates", name)
		}
	}
}

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

// Every repair block a repair asks for still exists in its template, and still
// holds the code the repair means to write.
//
// The markers were lost once already, when the template was restored from a
// copy taken before they were placed. repairBlock panics at init, so that was
// caught in seconds, but a block whose markers survive in the wrong place would
// not be: it would hand an upgraded project a slice of somebody else's
// function. So this checks each block's shape and not only that it resolves.
func TestEveryRepairBlockResolves(t *testing.T) {
	blocks := map[string]struct {
		text        string
		startsWith  string
		endsWith    string
		mustContain []string
	}{
		"login": {
			text:       loginBlockNew,
			startsWith: "// invalidCredentials is the one answer",
			// The blank line after Login is part of it: the text it replaces
			// carried one, and gofmt will not put it back.
			endsWith:    "}\n\n",
			mustContain: []string{"func (h *AuthHandler) Login(", "func (h *AuthHandler) loginRefusal(", "func (h *AuthHandler) startTOTPChallenge("},
		},
		"two-factor-lookup": {
			text:        emailChallengeNew,
			startsWith:  "\t// Which columns are safe to read here",
			endsWith:    "\t}",
			mustContain: []string{"h.AuthService.EnabledTwoFactor("},
		},
		"email-challenge": {
			text:        emailChallengeTokenNew,
			startsWith:  "\tpending := models.TOTPPendingToken{",
			endsWith:    "\treturn true\n}",
			mustContain: []string{"h.AuthService.StartPendingTOTP(", "Two-factor authentication required"},
		},
		"auth-service-queries": {
			text:       repairBlock("api/services/auth.go", "auth-service-queries"),
			startsWith: "// The reads and writes behind signing in",
			endsWith:   "}\n\n",
			mustContain: []string{
				"func (s *AuthService) UserByEmail(",
				"func (s *AuthService) UserByID(",
				"func (s *AuthService) EnabledTwoFactor(",
				"func (s *AuthService) StartPendingTOTP(",
				"func (s *AuthService) ClearLoginFailures(",
			},
		},
	}
	for id, want := range blocks {
		t.Run(id, func(t *testing.T) {
			if !strings.HasPrefix(want.text, want.startsWith) {
				t.Errorf("starts at the wrong line:\n%.80q", want.text)
			}
			if !strings.HasSuffix(want.text, want.endsWith) {
				t.Errorf("ends at the wrong line:\n%.80q", want.text[max(0, len(want.text)-80):])
			}
			for _, s := range want.mustContain {
				if !strings.Contains(want.text, s) {
					t.Errorf("does not contain %q", s)
				}
			}
			if strings.Contains(want.text, repairMarkerPrefix) {
				t.Error("carries a marker line, which would end up in a user's file")
			}
		})
	}
}

// No generated project ever sees a marker. They are a fact about this
// repository and about nobody's app.
func TestNoTemplateLeaksARepairMarker(t *testing.T) {
	names, err := templateNames()
	if err != nil {
		t.Fatal(err)
	}
	marked := 0
	for _, name := range names {
		if strings.Contains(tmpl(name), repairMarkerPrefix) {
			t.Errorf("%s leaks a repair marker into generated code", name)
		}
		raw, err := templateFS.ReadFile("templates/" + name + ".tmpl")
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(raw), repairMarkerPrefix) {
			marked++
		}
	}
	if marked == 0 {
		t.Error("no template carries a marker, so every repair block must have been lost")
	}
}
