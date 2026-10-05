package generate

import (
	"os"
	"path/filepath"
	"regexp"
)

// What the admin calls a related row.
//
// Every relationship field carries a displayField: the column the dropdown
// shows instead of a UUID. The generator used to write "name" for all of them
// without looking at the related model, and plenty of models have no `name`: a
// Collection has a `title`, an Invoice has a `number`, a Tag has a `label`.
//
// The form field survived it on a fallback chain (displayField || name || title
// || id), so the mistake was invisible until the bulk-create grid's dropdown
// rendered raw UUIDs. The fallback is still there, and worth keeping for a
// model this guesses wrong about, but the generator should not be guessing when
// the answer is in the file next door.

// labelColumns are the names a label goes by, best first.
//
// Order is the point. A model with both a `name` and a `title` is showing the
// name; one with a `title` and a `label` is showing the title. Only when none
// of them is there does the first text column win.
var labelColumns = []string{"name", "title", "label", "subject", "number", "code", "reference"}

// modelStringField matches a model's exported string column and its json name:
//
//	Title string `gorm:"size:255" json:"title" ...`
var modelStringField = regexp.MustCompile("(?m)^\\s*([A-Z][A-Za-z0-9_]*)\\s+string\\b[^`\\n]*`([^`]*)`")

// modelJSONName pulls the json tag's name out of a struct tag.
var modelJSONName = regexp.MustCompile(`json:"([^",]+)`)

// displayFieldFor returns the column the admin should show for a row of the
// related model: the first of labelColumns the model has, else its first text
// column, else "name".
//
// "name" is the fallback rather than an error because the related model may not
// exist yet: `grit generate resource Order --fields "customer:belongs_to:Customer"`
// run before Customer is a reasonable thing to do, and the form's own chain
// covers it until the resource is regenerated.
func (g *Generator) displayFieldFor(relatedModel string) string {
	if relatedModel == "" {
		return "name"
	}
	path := filepath.Join(g.APIRoot(), "internal", "models", toSnakeCase(relatedModel)+".go")
	raw, err := os.ReadFile(path)
	if err != nil {
		return "name"
	}

	// Only the model's own struct, so a nested type declared below it does not
	// contribute columns.
	src := string(raw)
	start := regexp.MustCompile(`(?m)^type ` + regexp.QuoteMeta(relatedModel) + ` struct \{`).FindStringIndex(src)
	if start == nil {
		return "name"
	}
	end := regexp.MustCompile(`(?m)^\}`).FindStringIndex(src[start[1]:])
	body := src[start[1]:]
	if end != nil {
		body = body[:end[0]]
	}

	var columns []string
	for _, m := range modelStringField.FindAllStringSubmatch(body, -1) {
		tag := modelJSONName.FindStringSubmatch(m[2])
		if tag == nil || tag[1] == "-" {
			continue
		}
		columns = append(columns, tag[1])
	}

	have := make(map[string]bool, len(columns))
	for _, c := range columns {
		have[c] = true
	}
	for _, want := range labelColumns {
		if have[want] {
			return want
		}
	}
	// No column with a label's name. The first text column is a better guess
	// than a column that is not there: a dropdown of UUIDs is useless, and a
	// dropdown of the wrong text at least names the row.
	for _, c := range columns {
		if c != "id" {
			return c
		}
	}
	return "name"
}
