package generate

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// marshalJSONHead is the first line of the generated method's doc comment, and
// the marker that finds it again. A rebuild has to replace the method rather
// than write a second one, and the comment is the only thing on it that no
// developer has a reason to change.
const marshalJSONHead = "// MarshalJSON leaves out the columns declared sensitive:"

// hiddenColumn is one column the method drops: the Go field name for the
// shadow field, and the JSON name it takes over.
type hiddenColumn struct {
	Go   string
	JSON string
}

// sensitiveFieldRe reads a model line that declares a sensitive column:
// the Go name, and the json name from its tag.
var sensitiveFieldRe = regexp.MustCompile(
	"^\\s*(\\w+)\\s+[\\w\\.\\*\\[\\]]+\\s+`[^`]*json:\"([a-zA-Z0-9_]+)[^`]*sensitive:\"true\"")

// marshalJSONFor is the method's source, or "" when there is nothing to hide.
//
// This is the only place the promise is kept without somebody choosing to keep
// it. Keeping a column out of the admin table, out of the export and out of a
// filter are three decisions in three generators, and each of them is a
// decision somebody can later make differently: a hand-added column, a custom
// handler returning the row, a webhook payload, the body that answers the POST
// that just set the value. All of those marshal the model, and all of them go
// through here.
//
// The shape is the standard shadowing trick rather than a marshal into a map
// and back. The embedded *alias has no methods, so there is no recursion; the
// shadow field sits at depth 0 and wins over the real one at depth 1; and
// being a nil pointer with omitempty, encoding/json leaves it out entirely
// rather than writing null. One marshal, no reflection of our own, and the
// field order of everything else is untouched.
func marshalJSONFor(pascal string, hidden []hiddenColumn) string {
	if len(hidden) == 0 {
		return ""
	}
	var shadows strings.Builder
	var listed []string
	for _, h := range hidden {
		listed = append(listed, h.JSON)
		shadows.WriteString(fmt.Sprintf(
			"\t\t%s *struct{} `json:%q`\n", h.Go, h.JSON+",omitempty"))
	}
	return fmt.Sprintf(`%s %s.
//
// Those may be written and are never read back. Not in a list, not in the
// body that answers the request that set them, not in a webhook payload, not
// in an audit event. Every one of those marshals this struct, so this is the
// one place that has to be right.
//
// The shadow fields below are nil pointers with omitempty, so encoding/json
// drops them; sitting at depth 0 they take the JSON name from the embedded
// originals at depth 1. The alias type has no methods of its own, which is
// what keeps this from calling itself.
func (m %s) MarshalJSON() ([]byte, error) {
	type alias %s
	return json.Marshal(&struct {
		*alias
%s	}{alias: (*alias)(&m)})
}
`, marshalJSONHead, strings.Join(listed, ", "), pascal, pascal, shadows.String())
}

// sensitiveMarshalJSON writes the method for a model the generator is building
// from a definition, with a leading blank line because it is appended to a
// file already holding the struct and its hooks.
func sensitiveMarshalJSON(names Names, fields []Field) string {
	var hidden []hiddenColumn
	for _, f := range fields {
		if f.Sensitive {
			hidden = append(hidden, hiddenColumn{
				Go:   toPascalCase(f.Name),
				JSON: toSnakeCase(f.Name),
			})
		}
	}
	body := marshalJSONFor(names.Pascal, hidden)
	if body == "" {
		return ""
	}
	return "\n" + body
}

// ensureSensitiveMarshalJSON rebuilds the method on a model that already
// exists, from the columns that file currently declares.
//
// Rebuilt rather than appended to, because `grit add field` runs on a model
// that may already hide a column: a method naming only the field being added
// would start returning the one that was already hidden. Reading the file is
// also the only way to be right about a column somebody tagged by hand.
func ensureSensitiveMarshalJSON(path, pascal string) error {
	l, err := readGoLines(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", filepath.Base(path), err)
	}

	var hidden []hiddenColumn
	for _, line := range l.lines {
		if m := sensitiveFieldRe.FindStringSubmatch(line); m != nil {
			hidden = append(hidden, hiddenColumn{Go: m[1], JSON: m[2]})
		}
	}

	// Where the existing method is, if there is one: the doc comment's first
	// line through the closing brace of the func that follows it.
	start := l.findPrefix(0, marshalJSONHead)
	end := -1
	if start >= 0 {
		if open := l.findPrefix(start, "func (m "+pascal+") MarshalJSON("); open >= 0 {
			end = l.closeOf(open)
		}
		if end < 0 {
			// A method whose shape this cannot read is a method somebody has
			// edited. Leaving it alone and saying so is better than writing a
			// second one beside it.
			return fmt.Errorf("models/%s has a MarshalJSON this could not re-read: "+
				"check by hand that it drops every column tagged sensitive", filepath.Base(path))
		}
	}

	body := marshalJSONFor(pascal, hidden)
	if body == "" {
		return nil // nothing to hide, and nothing was written before
	}

	if start >= 0 {
		// Replace in place, and keep the blank line that precedes it.
		l.lines = append(l.lines[:start], append(strings.Split(strings.TrimRight(body, "\n"), "\n"), l.lines[end+1:]...)...)
	} else {
		l.lines = append(l.lines, "", strings.TrimRight(body, "\n"))
	}

	// encoding/json sorts before every other standard-library import a model
	// has, so it goes at the top of the block. Added as a line in the list
	// rather than by rewriting an anchor: an anchor that has moved fails
	// silently, and codefmt sorts the block afterwards anyway.
	if !l.contains(`"encoding/json"`) {
		if open := l.find(0, "import ("); open >= 0 {
			l.insert(open+1, "\t\"encoding/json\"")
		} else {
			return fmt.Errorf("models/%s has no import block: add \"encoding/json\" by hand",
				filepath.Base(path))
		}
	}
	return l.write(path)
}
