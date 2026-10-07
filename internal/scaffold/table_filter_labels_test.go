package scaffold

import (
	"strings"
	"testing"
)

// Every control in the table toolbar has an accessible name.
//
// The filters had none. A screen reader announced each one as "combobox, All
// Role": the current value, with nothing saying what it filters. The date-range
// inputs had nothing at all, and the number pair had only a placeholder, which
// is the weakest source of a name and gone the moment you type.
func TestTableFilterControlsAreNamed(t *testing.T) {
	src := adminTableFilters()

	// The two selects, the number pair and the date pair: five names.
	if got := strings.Count(src, "aria-label="); got < 6 {
		t.Errorf("the filter controls carry %d aria-labels, want one per control", got)
	}
	for _, want := range []string{
		"Filter by ${filter.label}",
		"${filter.label}, from",
		"${filter.label}, to",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("no filter control is named %q", want)
		}
	}
	// The name comes from the filter, so it stays right when a column is renamed.
	if strings.Contains(src, `aria-label="Filter"`) {
		t.Error("a filter is named by a constant, which will be wrong for every column but one")
	}

	// And the page-size select, which is a control like any other.
	if !strings.Contains(adminTablePagination(), `aria-label={t("table.perPageLabel", "Rows per page")}`) {
		t.Error("the page-size select has no name")
	}
}
