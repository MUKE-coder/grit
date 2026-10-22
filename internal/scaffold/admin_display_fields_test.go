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
}

// A field that only applies when another field holds a particular value.
//
// Three layouts render a form and each one has its own field loop, so a
// condition honoured by the plain form and not by the wizard is a feature that
// breaks the moment somebody adds steps. Every layout is pinned here, along
// with the clearing rule: a hidden field is submitted empty rather than left
// alone, or a project moved off recurring keeps billing monthly on a column
// nobody can see.
func TestConditionalFieldsAreHonouredByEveryLayout(t *testing.T) {
	for name, c := range map[string]struct {
		src   string
		wants []string
	}{
		"lib/resource.ts": {adminResourceTypes(), []string{
			"visibleWhen?: (values: Record<string, unknown>) => boolean;",
			"requiredWhen?: (values: Record<string, unknown>) => boolean;",
		}},
		"lib/form-values.ts": {adminFormValues(), []string{
			"export function isVisible(",
			"export function isRequired(",
			"export function emptyValueFor(",
			"if (!isVisible(field, values)) return false;",
			"out[key] = clear.has(key) ? clear.get(key) : value;",
		}},
		"form-builder": {adminFormBuilder(), []string{
			"export function useConditionalFields(",
			"useWatch({ control, disabled: !conditional })",
			"{visibleFields.map((field) => (",
		}},
		"form-stepper": {adminFormStepper(), []string{
			"useConditionalFields(",
			"{visibleIn(step.fields).map((field) => (",
		}},
		"update-groups": {adminUpdateGroups(), []string{
			"useConditionalFields(groupFields, control)",
			"{visibleGroupFields.map((field) => (",
		}},
		// The suite every new project gets. It does not reach a project that
		// already exists: frontend test files are written by scaffold and not
		// by upgrade, so an existing project gets the behaviour and keeps its
		// older copy of this file.
		"__tests__/conditional-fields.test.tsx": {adminConditionalFieldsTest(), []string{
			"visibleWhen: (v) => v.kind === \"recurring\"",
			"expect(onSubmit.mock.calls[0][0].cycle).toBe(\"\")",
		}},
		"__tests__/form-values.test.ts": {adminFormValuesTest("@/lib"), []string{
			"describe(\"conditional fields\"",
			"clears a hidden field instead of keeping the stale value",
			"never requires a hidden field, whatever required says",
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
