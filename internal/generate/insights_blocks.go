package generate

import (
	"fmt"
	"strings"
)

// What a list page can say about itself without being configured.
//
// Every resource page shows four cards: the total and three date windows. They
// are the same four on every page, because they are the only questions that can
// be asked of a table whose columns are unknown. The columns are not unknown
// here: the generator knows which are booleans and which are a short list of
// options, and those are exactly the columns worth counting. A page that says
// "Active 2, Inactive 237, Draft 14, Live 3" the moment it is generated is a
// different page from one that says "Total 239", and neither needed a line of
// configuration.
//
// Only booleans and option lists. A free-text column has as many values as it
// has rows, and a chart of ten thousand one-row bars is not an insight; the
// API caps what it returns at twelve values either way, so the cap would hide
// the problem rather than fix it.

// maxCountedColumns matches the API's cap on ?breakdown=. Emitting more would
// write config that the server drops, which reads like a bug in the admin.
const maxCountedColumns = 3

// insightsBlocks returns the stats and insights properties for a generated
// resource definition, or two empty strings when the model has nothing worth
// counting.
func (g *Generator) insightsBlocks() (string, string) {
	counted := g.countableFields()
	if len(counted) == 0 {
		return "", ""
	}

	var cards []string
	var columns []string
	var labels []string
	for _, f := range counted {
		column := toSnakeCase(f.Name)
		columns = append(columns, fmt.Sprintf("%q", column))

		switch {
		case FieldType(f.Type) == FieldBool || f.IsToggle():
			// "Active" and "Inactive" say what they are, so no column prefix.
			// Both spellings of each value, because Postgres answers true and
			// false where SQLite and MySQL answer 1 and 0, and the same admin
			// is pointed at all three over a project's life.
			yes, no := boolLabels(f.Name)
			cards = append(cards, fmt.Sprintf(`      {
        field: %q,
        labels: { "true": %q, "1": %q, "false": %q, "0": %q },
        only: ["true", "false"],
        icon: "ToggleLeft",
      },`, column, yes, yes, no, no))
			labels = append(labels, fmt.Sprintf(`      %s: { "true": %q, "1": %q, "false": %q, "0": %q },`, tsKey(column), yes, yes, no, no))

		default:
			// A select or radio: one card per option, in the order the model
			// declares them rather than by size, so the cards do not reorder
			// themselves as the data changes.
			var only []string
			var optionLabels []string
			for _, option := range f.Options {
				only = append(only, fmt.Sprintf("%q", option.Value))
				optionLabels = append(optionLabels, fmt.Sprintf("%q: %q", option.Value, option.Label))
			}
			cards = append(cards, fmt.Sprintf(`      {
        field: %q,
        label: %q,
        labels: { %s },
        only: [%s],
      },`, column, fieldLabel(f.Name), strings.Join(optionLabels, ", "), strings.Join(only, ", ")))
			labels = append(labels, fmt.Sprintf(`      %s: { %s },`, tsKey(column), strings.Join(optionLabels, ", ")))
		}
	}

	stats := fmt.Sprintf(`
  // One card per value, counted over the rows the table is showing. Generated
  // from the model: %s. Delete a line to drop its cards, or set
  // stats: { countBy: [] } to keep only the four defaults.
  stats: {
    countBy: [
%s
    ],
  },`, strings.Join(countedNames(counted), ", "), strings.Join(cards, "\n"))

	insights := fmt.Sprintf(`
  // The collapsible charts above the table: created over time, and these
  // columns by value. Nothing is requested until somebody opens the panel.
  insights: {
    breakdown: [%s],
    labels: {
%s
    },
  },`, strings.Join(columns, ", "), strings.Join(labels, "\n"))

	return stats, insights
}

// countableFields picks the columns whose values are a short, known list.
func (g *Generator) countableFields() []Field {
	var out []Field
	for _, f := range g.Definition.Fields {
		if f.Encrypted || f.IsManyToMany() {
			continue
		}
		switch {
		case FieldType(f.Type) == FieldBool || f.IsToggle():
		case (f.IsSelect() || f.IsRadio()) && f.HasOptions():
			// A select with more options than the chart can show is still
			// worth counting: the API returns the largest twelve and the
			// cards follow the model's own order.
		default:
			continue
		}
		out = append(out, f)
		if len(out) == maxCountedColumns {
			break
		}
	}
	return out
}

func countedNames(fields []Field) []string {
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		out = append(out, toSnakeCase(f.Name))
	}
	return out
}

// boolLabels names the two sides of a boolean column.
//
// "Active" and "Inactive" for active, "Published" and "Unpublished" for
// published: the column is already an adjective, so the negative is the
// adjective with a prefix. The prefixes are the ones English actually uses for
// the words that turn up as boolean columns, and anything else falls back to
// the column's own name with "Not".
func boolLabels(name string) (string, string) {
	label := fieldLabel(name)
	lower := strings.ToLower(label)
	// A slice rather than a map keyed by the prefix: a map is iterated in a
	// random order, and a generator that writes different bytes on each run
	// cannot be diffed, tested or repaired.
	for _, group := range []struct {
		prefix string
		words  []string
	}{
		{"In", []string{"active", "complete", "valid", "visible", "correct"}},
		{"Un", []string{"published", "read", "paid", "verified", "confirmed", "approved", "available", "locked", "archived", "subscribed"}},
	} {
		for _, word := range group.words {
			if lower == word {
				return label, group.prefix + lower
			}
		}
	}
	return label, "Not " + label
}

// tsKey quotes an object key only when it has to be quoted. A snake_case
// column never does, but this is what writes the key, so it should not be the
// thing that produces invalid TypeScript for an odd column name.
func tsKey(name string) string {
	for i, r := range name {
		if r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (i > 0 && r >= '0' && r <= '9') {
			continue
		}
		return fmt.Sprintf("%q", name)
	}
	return name
}

// fieldLabel is the column as a person reads it: is_active -> Is active,
// payment_status -> Payment status.
func fieldLabel(name string) string {
	words := strings.ReplaceAll(toSnakeCase(name), "_", " ")
	if words == "" {
		return name
	}
	return strings.ToUpper(words[:1]) + words[1:]
}
