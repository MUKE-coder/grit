package scaffold

import (
	"strings"
	"testing"
)

// grit#145, the second half: the README says "the API validates against the
// same Zod schema as the form", and that was true for the pages somebody
// writes by hand and not for the resources Grit generates. form-builder.tsx
// applied its own required plus format checks, and the schema next door was
// used by nothing.
func TestTheFormBuilderUsesTheResourcesSchemaWhenItHasOne(t *testing.T) {
	src := adminFormBuilder()

	if !strings.Contains(src, `import { standardSchemaResolver } from "@hookform/resolvers/standard-schema";`) {
		t.Fatal("the form builder does not import a resolver, so a form definition's " +
			"schema is ignored")
	}
	if !strings.Contains(src, "resolver: formDef.schema") {
		t.Error("the form builder does not pass the resource's schema to useForm")
	}

	// Opt-in. A form with no schema has to behave exactly as it did before,
	// which is what keeps this from being a behaviour change for every form in
	// every existing project: `grit upgrade` does not rewrite a resource
	// definition, so an existing resource keeps the field rules alone.
	if !strings.Contains(src, ": undefined,") {
		t.Error("a form with no schema must pass no resolver at all, so it keeps the " +
			"per-field rules alone")
	}

	// The field rules are not replaced. They know a rating's max and whether a
	// conditional field is even visible; the schema knows a refinement that
	// reads a sibling. Dropping either loses rules.
	for _, want := range []string{"formatRules(", "required:"} {
		if !strings.Contains(src, want) {
			t.Errorf("%q is gone from the form builder: the per-field rules are ANDed "+
				"with the schema, not replaced by it", want)
		}
	}
}

// The stepped form keeps the field rules alone, and that is a decision.
//
// A stepped form on edit PATCHes one step's fields at a time. Validating that
// payload against the whole create schema would fail on every field in another
// step, so the schema cannot express what a per-step save needs.
func TestTheSteppedFormDoesNotValidateAgainstTheWholeCreateSchema(t *testing.T) {
	// Its own template, not the form builder's.
	src := adminFormStepper()

	if !strings.Contains(src, "function FormStepper") {
		t.Fatal("FormStepper is gone, so this test is reading the wrong template")
	}
	if strings.Contains(src, "resolver:") {
		t.Error("FormStepper passes a resolver. On edit it PATCHes one step's fields, " +
			"and the create schema requires fields that are in other steps, so every " +
			"per-step save would fail validation.")
	}
}

// And the generator sets it, or none of the above does anything.
func TestAGeneratedResourceCarriesItsSchema(t *testing.T) {
	// The form definition's type has to allow it first.
	defs := adminResourceTypes()
	if !strings.Contains(defs, "schema?:") {
		t.Error("FormDefinition has no schema field, so a resource cannot name one")
	}
}
