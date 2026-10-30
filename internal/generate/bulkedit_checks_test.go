package generate

import (
	"regexp"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/codefmt"
	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// The grid's save, repaired, has to be the grid's save as the generator writes
// it. Two things that must agree, with nothing else comparing them:
// TestGeneratedWritesDoNotReadTheRowBack cannot, because its fixtures predate
// the grid and have no BulkEdit at all.
//
// They disagreed. The repair added the policy pass, which returns the rows it
// refused, and did not report them, so `denied` was declared and unused and
// every upgraded project with a grid stopped compiling. The generator had the
// other half of the same bug: readableType was a package-level function, so a
// second generated resource redeclared it.
//
// BulkEdit's body does not depend on the resource's fields, only on its name,
// which is why one fixture from a real project can be compared against a
// Contact rendered from two fields here.
func TestRepairedBulkEditIsTheGeneratedBulkEdit(t *testing.T) {
	const module = "cors107/apps/api"

	old := fixture(t, "grid_service.go.txt")
	if !strings.Contains(old, "func (s *ContactService) BulkEdit(ctx context.Context") {
		t.Fatal("the fixture has no BulkEdit, so this test compares nothing")
	}

	out, changed := scaffold.RepairGeneratedServiceChecks(old, module)
	if !changed {
		t.Fatal("the grid service was not repaired")
	}
	out = codefmt.Go(out)

	want := bulkEditOf(t, renderService(t, module, "Contact", "name:string,phone:string", nil))
	got := bulkEditOf(t, out)
	if got != want {
		t.Errorf("the repaired BulkEdit is not the generated one.\n--- repaired\n%s\n--- generated\n%s", got, want)
	}

	// The specific failure, named, so a future half-repair says why rather
	// than only that two blocks of text differ.
	if strings.Contains(got, ":= s.permitted(ctx, allowed,") && !strings.Contains(got, "deniedSet[w.id]") {
		t.Error("BulkEdit asks the policy and throws away the rows it refused: that is a declared and unused variable, and a project that does not compile")
	}

	// Run twice changes nothing.
	if again, changed := scaffold.RepairGeneratedServiceChecks(out, module); changed || again != out {
		t.Error("repairing the grid service twice changed it again")
	}
}

// Nothing a generated service declares at package level may be shared between
// two resources, because the generator writes one copy of it per resource into
// the same package. readableType was, and a project with two generated
// resources did not compile.
func TestGeneratedServiceSharesNoPackageLevelNames(t *testing.T) {
	const module = "l19x-old/apps/api"
	src := renderService(t, module, "Contact", "name:string,phone:string", nil)

	decl := regexp.MustCompile(`(?m)^(?:func|type|var|const) (\w+)`)
	for _, m := range decl.FindAllStringSubmatch(src, -1) {
		name := m[1]
		if !strings.Contains(strings.ToLower(name), "contact") {
			t.Errorf("the service declares %q at package level, which the next generated resource will declare too; make it a method or name it after the resource", name)
		}
	}
}

// bulkEditOf is the BulkEdit method, from its signature to the closing brace
// in the first column.
func bulkEditOf(t *testing.T, src string) string {
	t.Helper()
	const sig = "func (s *ContactService) BulkEdit(ctx context.Context"
	at := strings.Index(src, sig)
	if at < 0 {
		t.Fatal("no BulkEdit in this source")
	}
	rest := src[at:]
	end := strings.Index(rest, "\n}\n")
	if end < 0 {
		t.Fatal("BulkEdit has no closing brace")
	}
	return rest[:end+3]
}
