package scaffold

import (
	"strings"
	"testing"
)

// A money column has always rendered through formatMoney, which is
// Intl.NumberFormat and therefore groups: a table reads $1,234,567.89. The
// input beside it read 1234567.89, and telling six figures from seven by
// counting digits is the mistake the field exists to prevent.
func TestMoneyFieldGroupsItsDigits(t *testing.T) {
	src := adminMoneyField()

	for _, want := range []string{
		"function groupMajor(",
		"function parseDraft(",
		"toLocaleString(undefined, {",
		// the four places a draft is written from a value
		"useState(() =>\n    value ? groupMajor(toMajor(value), digits) : \"\"",
		"setDraft(value ? groupMajor(toMajor(value), digits) : \"\")",
		"setDraft(groupMajor(toMajor(value ?? { amount: 0, currency }), digits))",
		"setDraft(groupMajor(major, currencyExponent(next)))",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the money field is missing %q", want)
		}
	}

	// The guard that decides whether an outside change should rewrite the box
	// has to read the draft back through parseDraft. Number.parseFloat of a
	// grouped "1,234.56" is 1, so a bare parseFloat here would see every
	// grouped draft as wrong and rewrite the box on every render.
	if strings.Contains(src, "const shown = Number.parseFloat(draft)") {
		t.Error("the sync guard still reads the draft with Number.parseFloat, " +
			"which returns 1 for \"1,234.56\"")
	}
	if !strings.Contains(src, "const shown = parseDraft(draft)") {
		t.Error("the sync guard does not read the draft through parseDraft")
	}

	// Grouping must not happen on a keystroke: that is what moves the caret,
	// and the component holds a separate draft precisely to avoid it.
	if strings.Contains(src, "setDraft(groupMajor") && strings.Contains(src, "const handleAmount") {
		amount := src[strings.Index(src, "const handleAmount"):]
		if end := strings.Index(amount, "\n  };"); end > 0 {
			if strings.Contains(amount[:end], "groupMajor") {
				t.Error("handleAmount groups while typing, which moves the caret")
			}
		}
	}
}
