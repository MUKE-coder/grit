package scaffold

import (
	"regexp"
	"strings"
	"testing"
)

// The History tab reads one record's activity, and that takes two things that
// have to agree.
//
// paginate refuses a filter the whitelist does not name, and the whitelist
// cannot see a parameter the handler never binds. Either one alone looks
// right in review and produces a tab that silently shows every record's
// history instead of this record's, or none at all. So both files are read
// here, in one test (grit#261).
func TestRecordHistoryFilterIsBoundAndWhitelisted(t *testing.T) {
	handler := tmpl("api/handlers/user_activity.go")
	reads := tmpl("api/services/activity_reads.go")

	if !strings.Contains(handler, `With("resource_id", c.Query("resource_id"))`) {
		t.Error("the user-activity handler does not bind resource_id, so the History tab cannot ask for one record")
	}

	// The whitelist, read as the one line it is, so a resource_id appearing in
	// some other config in this file cannot satisfy the assertion.
	line := regexp.MustCompile(`(?m)^\s*Filterable:\s*map\[string\]bool\{[^}]*\},`).FindString(
		reads[strings.Index(reads, "userActivityListConfig"):],
	)
	if line == "" {
		t.Fatal("no Filterable map found in userActivityListConfig")
	}
	for _, want := range []string{"resource_id", "resource_type"} {
		if !strings.Contains(line, `"`+want+`"`) {
			t.Errorf("userActivityListConfig.Filterable does not whitelist %q: %s", want, strings.TrimSpace(line))
		}
	}
}

// The tab itself, in the file that is delivered as one unit.
//
// RecordHistory lives in resource-detail-page.tsx rather than a file of its
// own on purpose: a component in a second file is a second file `grit
// upgrade` has to write in the same pass, and a project that received one
// without the other would not build.
func TestHistoryTabShipsInsideTheDetailPage(t *testing.T) {
	page := adminResourceDetailPage()

	for _, want := range []string{
		"function RecordHistory(",
		`useState<"details" | "history">("details")`,
		`{tab === "history" && <RecordHistory resource={resource} id={id} />}`,
		"/api/user-activity?",
		"resource_type",
		"resource_id",
	} {
		if !strings.Contains(page, want) {
			t.Errorf("the detail page is missing %q", want)
		}
	}

	// The strip is no-print: the printed record and the PDF are the record,
	// not the audit trail of it.
	if !strings.Contains(page, `<div className="no-print flex items-center gap-1 border-b border-border">`) {
		t.Error("the tab strip is not marked no-print, so it would appear in the PDF")
	}

	// An empty list and a failed request must not look alike, because one of
	// them is a bug and the other is a new record.
	if !strings.Contains(page, "The history could not be loaded") {
		t.Error("a failed history request renders as an empty history")
	}
	if !strings.Contains(page, "Nothing has been recorded against this record yet") {
		t.Error("an empty history has no empty state")
	}
}

// The feed stores the entity name lowercased, which is what LogCreate and its
// siblings write, and the tab has to ask for it in the same spelling.
//
// Two readers of one convention: services/activity.go lowercases the entity
// when it writes the row, and the tab lowercases resource.name when it asks
// for it. A test over either alone passes while the pair disagrees.
func TestHistoryTabAsksForTheSpellingTheFeedStores(t *testing.T) {
	if !strings.Contains(userActivityServiceGo(), `ResourceType: strings.ToLower(entityType)`) {
		t.Error("the activity service no longer lowercases the entity when it writes a row")
	}
	if !strings.Contains(adminResourceDetailPage(), `(resource.name ?? "").toLowerCase()`) {
		t.Error("the History tab does not lowercase the resource name, so it would ask for a type the feed never stored")
	}
}

// A summary on a record page is a sentence somebody reads, so it must not
// name columns they did not touch and cannot see.
func TestDiffSummaryLeavesOutBookkeepingColumns(t *testing.T) {
	src := userActivityServiceGo()

	if !strings.Contains(src, "bookkeepingColumns") {
		t.Fatal("DiffSummary no longer filters bookkeeping columns")
	}
	for _, col := range []string{"version", "updated_at", "created_at"} {
		if !regexp.MustCompile(`"` + col + `":\s*true`).MatchString(src) {
			t.Errorf("%q is not treated as a bookkeeping column, so it appears in a record's history", col)
		}
	}
	// And an update that wrote nothing else says nothing, rather than "0
	// fields changed".
	if !strings.Contains(src, "if len(keys) == 0 {") {
		t.Error("DiffSummary has no empty case after filtering, so a bookkeeping-only update renders a summary about nothing")
	}
}
