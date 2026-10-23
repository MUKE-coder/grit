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

	// Grouping happens on the keystroke, because reading the number you are
	// halfway through typing is the whole point of the separators. What must
	// not happen is the caret being left at the end afterwards.
	amount := src[strings.Index(src, "const handleAmount"):]
	end := strings.Index(amount, "\n  };")
	if end < 0 {
		t.Fatal("could not find the end of handleAmount")
	}
	body := amount[:end]
	if !strings.Contains(body, "groupDraft(cleaned)") {
		t.Error("handleAmount does not group while typing, so the number stays " +
			"unreadable until the field is left")
	}
	if !strings.Contains(body, "caret.current = caretAfter(") {
		t.Error("handleAmount regroups without restoring the caret, which leaves " +
			"it at the end of the box on every keystroke")
	}
	// The currency's decimals must not be forced mid-type: that is what turns
	// "19." into "19.00" with the point taken out from under the caret.
	if strings.Contains(body, "groupMajor") {
		t.Error("handleAmount applies the currency's decimals while typing")
	}

	if !strings.Contains(src, "useLayoutEffect(() => {") {
		t.Error("the caret is restored in an effect rather than a layout effect, " +
			"so it would appear at the end of the box and then jump")
	}
	if !strings.Contains(src, "ref={inputRef}") {
		t.Error("the input takes no ref, so the caret cannot be put back")
	}
	if !strings.Contains(src, "e.target.selectionStart") {
		t.Error("the change handler never reads where the caret was")
	}
}

// Only the integer part is grouped while typing. Grouping the whole number
// means deciding how many decimals it has, and mid-typing there is no answer:
// "19." is on its way to "19.99".
func TestMoneyFieldLeavesTheFractionAloneWhileTyping(t *testing.T) {
	src := adminMoneyField()
	for _, want := range []string{
		"function groupDraft(",
		"function caretAfter(",
		"const fraction = dot === -1 ?",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the money field is missing %q", want)
		}
	}
}
