package generate

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

// grit#207: a --confidential resource with a belongs_to generated a test
// package that did not compile.
//
//	vet.exe: internal/handlers/record_test.go:99:2: declared and not used: folder
//
// confidentialOwnerOverride points the fixture row at the signed-in user by
// replacing the owner's assignment. Called a second time on its own output it
// could not find that assignment, fell into a branch that rebuilt the override
// from nothing, and so dropped every OTHER belongs_to assignment while leaving
// the lines that created those parents.
//
// Why the tests shipped with v3.505.0 did not catch it: every fixture there was
// a confidential resource whose only relation was the owner. With one parent
// there is nothing to drop. So these take the two-parent shape, which is what a
// real project has: the reporter had it on fifteen resources at once.

// confidentialWithAParent is a resource with an owner AND another belongs_to.
//
// Built through ParseInlineFields and ApplyFlags, which is what the CLI does
// with the same words the reporter typed. A hand-written ResourceDefinition is
// how the first version of this test went wrong: it left the parent relation
// optional, parentFixtures skipped it, and the test then asserted that an
// assignment the generator was never asked to make was missing.
func confidentialWithAParent(t *testing.T) *Generator {
	t.Helper()
	def, err := ParseInlineFields("Record", "folder:belongs_to[Folder],notes:text")
	if err != nil {
		t.Fatalf("parsing the fields: %v", err)
	}
	// grit generate resource Record --owned-by user --confidential
	def.ApplyFlags(false, false, "user", false, false, false, true)
	def.EnsureOwnerField()

	if owner := def.OwnerField(); owner == nil {
		t.Fatal("--owned-by user produced no owner field")
	}
	g := &Generator{Definition: def, Module: "app"}
	var parents int
	for _, f := range def.Fields {
		if f.IsBelongsTo() && !g.relOptional(f) {
			parents++
		}
	}
	if parents < 2 {
		t.Fatalf("the fixture has %d required belongs_to relations, want 2: with fewer than two "+
			"there is no other parent to drop and this test cannot see grit#207", parents)
	}
	return g
}

// The generated file parses. The compiler's own answer, not a substring.
func TestAConfidentialResourceWithAParentGeneratesValidGo(t *testing.T) {
	g := confidentialWithAParent(t)
	src := g.resourceTestSource(g.Names())

	if _, err := parser.ParseFile(token.NewFileSet(), "record_test.go", src, parser.AllErrors); err != nil {
		t.Fatalf("the generated test file does not parse: %v", err)
	}
}

// Every variable the file declares is read. Parsing cannot see this (unused
// locals are a type-check error, not a syntax one) and it is exactly the
// failure: a fixture created and never used.
func TestEveryFixtureInAGeneratedTestIsUsed(t *testing.T) {
	g := confidentialWithAParent(t)
	src := g.resourceTestSource(g.Names())

	for _, fn := range splitGoFuncs(src) {
		for _, name := range declaredFixtures(fn.body) {
			// Counting occurrences rather than Contains: the declaration
			// itself is one, so a variable that is only declared reads as 1.
			if strings.Count(fn.body, name) < 2 {
				t.Errorf("%s declares %q and never reads it, so the handlers package does not "+
					"compile: \"declared and not used\" (grit#207)", fn.name, name)
			}
		}
	}
}

// And the row is fully formed: the parent is set, not merely created. A test
// asserting an administrator sees nothing is weaker when the row it cannot see
// is missing a required foreign key, and on Postgres that row does not insert
// at all.
func TestTheHiddenRowPointsAtItsParent(t *testing.T) {
	g := confidentialWithAParent(t)
	src := g.resourceTestSource(g.Names())

	for _, want := range []struct{ test, owner string }{
		{"TestRecord_OwnerSeesWhatWasCreated", "owner"},
		{"TestRecord_AnAdminSeesNothing", "rowOwner"},
		{"TestRecord_AGrantHolderSeesEveryRow", "rowOwner"},
	} {
		body, ok := funcBody(src, want.test)
		if !ok {
			t.Errorf("%s was not generated", want.test)
			continue
		}
		if !strings.Contains(body, "m.FolderID = folder.ID") {
			t.Errorf("%s creates the folder and never points the record at it", want.test)
		}
		if !strings.Contains(body, "m.UserID = "+want.owner+".ID") {
			t.Errorf("%s does not point the row at %s, so the assertion holds either way",
				want.test, want.owner)
		}
	}
}

// An ordinary owned resource with a parent is untouched by any of this.
func TestAnOwnedResourceWithAParentIsUnchanged(t *testing.T) {
	g := confidentialWithAParent(t)
	g.Definition.Confidential = false
	src := g.resourceTestSource(g.Names())

	if _, err := parser.ParseFile(token.NewFileSet(), "record_test.go", src, parser.AllErrors); err != nil {
		t.Fatalf("does not parse: %v", err)
	}
	body, ok := funcBody(src, "TestRecord_AdminSeesWhatWasCreated")
	if !ok {
		t.Fatal("the admin round-trip test was not generated for a plain owned resource")
	}
	if !strings.Contains(body, "m.FolderID = folder.ID") || !strings.Contains(body, "m.UserID = user.ID") {
		t.Error("a plain --owned-by resource lost one of its belongs_to assignments")
	}
}

// ------------------------------------------------------------------ helpers

type goFunc struct{ name, body string }

// splitGoFuncs is every top-level func in src, by name and body. Text, because
// the subject is generated text and a real parse would hide the thing being
// looked for behind a type-check error.
func splitGoFuncs(src string) []goFunc {
	var out []goFunc
	lines := strings.Split(src, "\n")
	for i, line := range lines {
		if !strings.HasPrefix(line, "func ") {
			continue
		}
		name := strings.TrimPrefix(line, "func ")
		if cut := strings.IndexAny(name, "("); cut > 0 {
			name = name[:cut]
		}
		var body []string
		for _, l := range lines[i+1:] {
			if strings.HasPrefix(l, "}") {
				break
			}
			body = append(body, l)
		}
		out = append(out, goFunc{name: name, body: strings.Join(body, "\n")})
	}
	return out
}

func funcBody(src, name string) (string, bool) {
	for _, fn := range splitGoFuncs(src) {
		if fn.name == name {
			return fn.body, true
		}
	}
	return "", false
}

// declaredFixtures is the names declared with := at the top level of a body,
// which is what a fixture line looks like.
func declaredFixtures(body string) []string {
	var out []string
	for _, line := range strings.Split(body, "\n") {
		if !strings.HasPrefix(line, "\t") || strings.HasPrefix(line, "\t\t") {
			continue
		}
		trimmed := strings.TrimSpace(line)
		cut := strings.Index(trimmed, " := ")
		if cut <= 0 {
			continue
		}
		for _, name := range strings.Split(trimmed[:cut], ",") {
			name = strings.TrimSpace(name)
			if name == "" || name == "_" || strings.Contains(name, " ") {
				continue
			}
			out = append(out, name)
		}
	}
	return out
}
