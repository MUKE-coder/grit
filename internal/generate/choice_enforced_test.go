package generate

import (
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// Every field type that declares a closed set is one the API enforces.
//
// grit#226: the enum tag was written into the model for whoever read it and
// enforced nowhere, so the shared Zod schema was the only thing holding a
// select to its own options and the API stored status "whatever" with a 201.
// internal/fieldtypes reads the tag now, on every write that has a model.
//
// Two things that have to agree, so they are read in one test: which field
// types emit the tag, which is EnumTag's business, and which column shapes the
// checker understands, which is the template's. A test over either alone
// passes while the pair is broken. That is how the first version of the
// checker compared a check column's whole value, datatypes.JSONSlice[string],
// against one option and would have refused every multi-select write.
func TestEveryChoiceFieldTypeIsEnforceable(t *testing.T) {
	template := scaffold.FieldTypesSource()

	// The checker's two branches: a single value compared as text, and a list
	// walked element by element.
	for _, branch := range []string{
		"func offendingChoice(",
		"case reflect.Slice, reflect.Array:",
		`f.Tag.Get("enum")`,
	} {
		if !strings.Contains(template, branch) {
			t.Fatalf("internal/fieldtypes does not contain %q, so the enum tag is "+
				"back to being decoration", branch)
		}
	}

	// And it is wired into the callback rather than merely defined.
	for _, wiring := range []string{
		"choices := choiceFields(stmt.Schema)",
		"checkChoice(tx, choices,",
		"checkChoiceMap(tx, choices,",
	} {
		if !strings.Contains(template, wiring) {
			t.Errorf("the choice check is not wired into the write callback: %q missing", wiring)
		}
	}

	// Now the field types. Anything that emits an enum tag must land on a
	// column shape one of those two branches handles.
	handled := map[string]bool{
		"string":                      true, // select, radio
		"datatypes.JSONSlice[string]": true, // check
	}
	var withTag []string
	for _, typeName := range ValidFieldTypes() {
		f := Field{
			Name: "status",
			Type: typeName,
			Options: []FieldOption{
				{Value: "open", Label: "Open"},
				{Value: "closed", Label: "Closed"},
			},
		}
		if f.EnumTag() == "" {
			continue
		}
		withTag = append(withTag, typeName)
		if !handled[f.GoType()] {
			t.Errorf("field type %q declares enum:%q on a %s column, which the "+
				"checker in internal/fieldtypes does not understand: it compares a "+
				"single value as text and walks a slice element by element, and "+
				"anything else is silently never refused",
				typeName, f.EnumTag(), f.GoType())
		}
	}

	// A sanity floor, because a test that iterates an empty list passes.
	if len(withTag) < 3 {
		t.Errorf("only %v declare a closed set, expected at least select, radio "+
			"and check: if one stopped emitting the tag, its column is no longer "+
			"enforced", withTag)
	}
}
