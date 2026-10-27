package scaffold

import (
	"strings"
	"testing"
)

// The public form asks for whatever the model does not tell it to leave out,
// and the model had stopped telling it two things.
//
// A slug column is skipped by name, so one called "handle" was offered as a
// question nobody can answer. Required was read from binding:"required" alone,
// which the generator only puts on strings because gin's required rejects 0
// and false, so every required number and money column said "optional". And
// archived_at was never in the skip list at all: a stranger filling in a public
// form was asked to pick the date their own submission was put away, which is
// the fault grit#101 fixed in the admin's forms and nobody looked for here
// (found while fixing grit#108).
func TestPublicFormSkipsServerFilledColumns(t *testing.T) {
	src := formShareDispatchGo()

	for _, want := range []string{
		// Any slug column, by its tag rather than by being called "slug".
		`if f.Tag.Get("slug") != "" {`,
		// Both halves of "required".
		`required := f.Tag.Get("required") == "true"`,
		`f.Tag.Get("binding")`,
		// And the archive bookkeeping.
		`"archived_at": true,`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the public form reader is missing %q", want)
		}
	}
}
