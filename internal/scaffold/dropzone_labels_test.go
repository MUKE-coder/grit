package scaffold

import (
	"strings"
	"testing"
)

// grit#151: every control in the dropzone is an icon and nothing else.
//
// The component is operable: react-dropzone's getRootProps gives the drop zone
// its tabIndex and its Enter/Space, and every control is a real <button>, so
// Tab reaches it and Enter fires it. What it was not is legible. A screen
// reader announced "button", once per file, with nothing to say which file each
// one removed.
//
// That is the same fault as the thirteen form fields that shipped with a label
// attached to nothing: fine with a mouse, unusable otherwise, and invisible in
// a demo because the person demonstrating it can see.
//
// title does not count. It is a tooltip, announced inconsistently and not at
// all on a touch device.
func TestEveryIconOnlyButtonInTheDropzoneIsNamed(t *testing.T) {
	source := adminDropzone()

	// A <button ...> opening tag, ending at the ">" that closes it rather than
	// at the first one found. These buttons carry onClick={(e) => ...}, and an
	// arrow function puts a ">" inside the tag, so a [^>]* regex stops before
	// the attributes that follow and reads a labelled button as unnamed. That
	// is exactly what the first version of this test did.
	openings := buttonTags(source)
	if len(openings) < 5 {
		t.Fatalf("found %d buttons in the dropzone; this test is reading the wrong thing", len(openings))
	}

	for _, tag := range openings {
		if strings.Contains(tag, "aria-label") || strings.Contains(tag, "aria-labelledby") {
			continue
		}
		// A button with text in it names itself. Only the icon-only ones are
		// the problem, and those are the ones this component is made of.
		t.Errorf("a dropzone button has no accessible name:\n  %s\n\n"+
			"If it contains visible text this is a false alarm and the test should "+
			"learn to tell them apart. If it contains only an icon, a screen reader "+
			"announces \"button\" and nothing else (grit#151).", collapse(tag))
	}
}

// And the name has to say which file, because "Remove" five times in a list of
// five files is not a name, it is a count.
func TestTheDropzoneLabelsNameTheirFile(t *testing.T) {
	source := adminDropzone()
	for _, want := range []string{
		"Remove ${file.name}",
		"Remove ${fileLabel}",
		"Move ${fileLabel} up",
		"Move ${fileLabel} down",
	} {
		if !strings.Contains(source, want) {
			t.Errorf("no label reading %q, so the control does not say what it acts on", want)
		}
	}
	// The reorder buttons also say where in the list they are, which is the
	// other thing somebody who cannot see the list has no way to know.
	if !strings.Contains(source, "position ${index + 1} of ${total}") {
		t.Error("the reorder buttons do not announce their position in the list")
	}
}

// buttonTags returns each <button> opening tag, stepping over anything inside
// braces so an arrow function's ">" does not end the tag early.
func buttonTags(source string) []string {
	var out []string
	i := 0
	for {
		start := strings.Index(source[i:], "<button")
		if start < 0 {
			return out
		}
		start += i
		depth, closed := 0, false
		for j := start; j < len(source); j++ {
			switch source[j] {
			case '{':
				depth++
			case '}':
				depth--
			case '>':
				if depth == 0 {
					out = append(out, source[start:j+1])
					i, closed = j+1, true
				}
			}
			if closed {
				break
			}
		}
		if !closed {
			return out
		}
	}
}

func collapse(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
