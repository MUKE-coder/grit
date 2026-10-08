package generate

import (
	"strings"
	"testing"
)

// A column you can filter by is usually one you want to sort by.
//
// Sortable was built by testing GoType against "string", "int" and "uint". A
// date's GoType is *jsontime.Date and a timestamp's is *jsontime.DateTime, so
// neither matched and neither was sortable, while both were filterable. An
// invoices table that cannot be ordered by the date it was issued is the
// obvious case.
//
// It failed quietly, which is the part worth a test. paginate drops a sort
// column that is not on the whitelist rather than refusing the request, so
// ?sort=issued_on returned 200 with the rows in whatever order the database
// chose, and the page looked like it had a sorting bug rather than a missing
// permission.
func TestADateColumnCanBeSorted(t *testing.T) {
	gen := &Generator{Definition: &ResourceDefinition{
		Name: "Invoice",
		Fields: []Field{
			{Name: "Number", Type: "string"},
			{Name: "Amount", Type: "money"},
			{Name: "Rate", Type: "float"},
			{Name: "Paid", Type: "bool"},
			{Name: "IssuedOn", Type: "date"},
			{Name: "SettledAt", Type: "datetime"},
		},
	}}

	p := gen.crud(gen.Names())

	for _, column := range []string{"issued_on", "settled_at"} {
		if !strings.Contains(p.sortCols, `"`+column+`": true`) {
			t.Errorf("%s is not sortable, so ?sort=%s is silently ignored:\n  %s", column, column, p.sortCols)
		}
	}

	// A float covers float, percent and rating.
	if !strings.Contains(p.sortCols, `"rate": true`) {
		t.Errorf("a float column is not sortable:\n  %s", p.sortCols)
	}

	// Money sorts by the column that exists, not the struct.
	if !strings.Contains(p.sortCols, `"amount_amount": true`) {
		t.Errorf("money does not sort by its amount column:\n  %s", p.sortCols)
	}
	if strings.Contains(p.sortCols, `"amount": true`) {
		t.Errorf("money sorts by a column that does not exist:\n  %s", p.sortCols)
	}

	// Ordering by true before false is a sort nobody asks for, and every name
	// here widens the SQL surface.
	if strings.Contains(p.sortCols, `"paid": true`) {
		t.Errorf("a bool is sortable, which buys nothing:\n  %s", p.sortCols)
	}
}
