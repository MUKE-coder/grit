package generate

import (
	"fmt"
	"strings"
)

// columnPack describes a single packed-column substitution: a primary
// field whose row in the columns array carries a `cell:` callback
// rendering two stacked values, plus one or more "absorbed" fields the
// emitter should skip.
type columnPack struct {
	// primary is the snake_case key of the field whose original column
	// slot is replaced by the packed cell. The other absorbed fields
	// are silently skipped.
	primary string
	// line is the full `{ key: ..., label: ..., cell: (row) => ... },`
	// row, ready to be appended into the columns array.
	line string
}

// detectColumnPacks scans the resource's field list for known
// pack-eligible patterns and returns:
//   - a map from each absorbed field's snake_case key to the pack it
//     belongs to (so the emitter can both detect "skip this row" and
//     "this row is the pack's primary — emit the packed line here"),
//   - whether ANY pack fired, so the caller knows to inject the
//     StackedCell import.
//
// Currently recognises two patterns:
//
//	name + email                  → packed "Contact" column
//	first_name + last_name        → packed "Name" column
//
// Both are exclusive — name/email wins over first_name/last_name if
// somehow both are present (rare in practice). Other heuristics
// (price + currency badge, status + relative date) are roadmap items.
func detectColumnPacks(fields []Field) (map[string]columnPack, bool) {
	out := map[string]columnPack{}
	have := map[string]bool{}
	for _, f := range fields {
		have[toSnakeCase(f.Name)] = true
	}

	switch {
	case have["name"] && have["email"]:
		pack := columnPack{
			primary: "name",
			line: fmt.Sprintf(`{ key: "name", label: "Contact", sortable: true, searchable: true, cell: (row) => StackedCell({ top: String(row.name ?? %q), bottom: String(row.email ?? %q) }) },`,
				"", ""),
		}
		out["name"] = pack
		out["email"] = pack

	case have["first_name"] && have["last_name"]:
		pack := columnPack{
			primary: "first_name",
			line: `{ key: "first_name", label: "Name", sortable: true, searchable: true, cell: (row) => StackedCell({ top: String(row.first_name ?? "") + " " + String(row.last_name ?? "") }) },`,
		}
		out["first_name"] = pack
		out["last_name"] = pack
	}

	return out, len(out) > 0
}

// ensure strings package is referenced (toSnakeCase lives in the
// generate package and uses strings).
var _ = strings.ToLower
