package scaffold

import (
	"strings"
	"testing"
)

// Read-only and computed admin fields are shown and never submitted.
//
// A disabled input was the only way to show a value, and it was still sent:
// over a column the server computes and the PATCH allow-list names, saving the
// form wrote the displayed figure back over the server's. Every path a form
// submits through strips display-only values, and these pin each one, because
// the one that is missed is the one that writes.
func TestDisplayOnlyFieldsAreNeverSubmitted(t *testing.T) {
	for name, c := range map[string]struct {
		src   string
		wants []string
	}{
		"form-builder": {adminFormBuilder(), []string{
			"handleSubmit((data) => onSubmit(writableValues(formDef.fields, data)))",
			"if (isDisplayOnly(field)) {",
			"useWatch({ control })",
			`from "@/lib/form-values"`,
		}},
		"form-stepper": {adminFormStepper(), []string{
			"handleSubmit((data) => onSubmit(writableValues(formDef.fields, data)))",
			"await onStepSave(writableValues(formDef.fields, patch));",
		}},
		"update-groups": {adminUpdateGroups(), []string{
			"body: writableValues(groupFields, values)",
		}},
		"lib/resource.ts": {adminResourceTypes(), []string{
			"readOnly?: boolean;",
			"compute?: (values: Record<string, unknown>) => unknown;",
		}},
	} {
		for _, want := range c.wants {
			if !strings.Contains(c.src, want) {
				t.Errorf("%s is missing %s", name, want)
			}
		}
	}

	// No submit path may hand raw form data to its caller any more.
	if strings.Contains(adminFormBuilder(), "handleSubmit(onSubmit)") ||
		strings.Contains(adminFormStepper(), "handleSubmit(onSubmit)") {
		t.Error("a form still submits its raw values, display-only fields included")
	}
}

// The rules live in a module with no runtime imports, so they can be tested
// without rendering anything, and the test ships with every project.
func TestFormValuesModuleStandsAlone(t *testing.T) {
	src := adminFormValues()
	if strings.Contains(src, "import {") || strings.Contains(src, "from \"react") {
		t.Error("lib/form-values.ts gained a runtime import")
	}
	if !strings.Contains(adminFormValuesTest("@/lib"), `from "@/lib/form-values"`) {
		t.Error("the shipped test does not exercise the shipped module")
	}
	// The Vite admin keeps lib under src/.
	if !strings.Contains(adminFormValuesTest("@/src/lib"), `from "@/src/lib/form-values"`) {
		t.Error("the Vite admin's copy of the test imports the wrong path")
	}
}
