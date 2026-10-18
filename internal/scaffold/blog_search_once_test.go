package scaffold

import (
	"strings"
	"testing"
)

// The blog list and the admin's stat cards narrow by the same search. They
// used to build that WHERE separately, in two files, so adding a searchable
// column to one left the counts answering a different question from the table
// above them. These tests hold the clause to one definition.

// TestBlogSearchDefinedOnce is the property itself: whatever else changes, the
// LIKE belongs to exactly one function in the service.
func TestBlogSearchDefinedOnce(t *testing.T) {
	const clause = `LOWER(title) LIKE LOWER(?)`

	service := blogServiceGo()
	if got := strings.Count(service, clause); got != 1 {
		t.Errorf("blog_service.go states the search %d times, want 1", got)
	}
	if !strings.Contains(service, "func (s *BlogService) searchScope(") {
		t.Error("blog_service.go has no searchScope for the list and the counts to share")
	}
	if !strings.Contains(service, "func (s *BlogService) CountQuery(") {
		t.Error("blog_service.go has no CountQuery for the handler to ask for")
	}

	handler := blogHandlerGo()
	if strings.Contains(handler, clause) {
		t.Error("blog_handler.go rebuilds the search, which is the copy that drifts")
	}
	if !strings.Contains(handler, "h.Service.CountQuery(") {
		t.Error("blog_handler.go does not get its count query from the service")
	}
	// Reaching past the service to the database is how the second copy got
	// there in the first place.
	if strings.Contains(handler, "h.DB.WithContext") {
		t.Error("blog_handler.go queries the database directly; that is the service's job")
	}
}

// TestBlogServiceCountQueryRepair covers the upgrade path. The blog files are
// written for a new project and do not travel on upgrade, so an existing
// project gets CountQuery from this repair or not at all.
func TestBlogServiceCountQueryRepair(t *testing.T) {
	fresh := blogServiceGo()

	if out, changes, warnings := repairBlogServiceCountQuery(fresh, "{{MODULE}}"); out != fresh || len(changes)+len(warnings) != 0 {
		t.Errorf("a fresh blog_service.go needs no repair, got changes %v warnings %v", changes, warnings)
	}

	old := blogServiceBeforeSearchScope(t, fresh)
	out, changes, warnings := repairBlogServiceCountQuery(old, "{{MODULE}}")
	if len(warnings) != 0 || len(changes) != 1 {
		t.Fatalf("repairing the previous blog_service.go: changes %v warnings %v", changes, warnings)
	}
	if !strings.Contains(out, "func (s *BlogService) CountQuery(") {
		t.Error("the repair did not add CountQuery")
	}
	if !strings.Contains(out, `"context"`) {
		t.Error("the repair did not add the context import CountQuery needs")
	}
	// The point of the repair is one definition, so List has to move onto it
	// rather than keep its own inline copy beside the new one.
	if got := strings.Count(out, `LOWER(title) LIKE LOWER(?)`); got != 1 {
		t.Errorf("after the repair the search is stated %d times, want 1", got)
	}
	if !strings.Contains(out, "s.searchScope(s.DB.Model(&models.Blog{}), search)") {
		t.Error("List did not move onto searchScope")
	}

	// Running it twice must not add a second copy.
	again, changes, warnings := repairBlogServiceCountQuery(out, "{{MODULE}}")
	if again != out || len(changes)+len(warnings) != 0 {
		t.Errorf("the repair is not idempotent: changes %v warnings %v", changes, warnings)
	}
}

// TestBlogServiceRepairWarnsOnEditedConstructor is the safety property. When
// the anchor is gone the repair must say so rather than guess, because the
// handler repair is skipped on that warning and a handler calling a method
// that is not there would not compile.
func TestBlogServiceRepairWarnsOnEditedConstructor(t *testing.T) {
	old := blogServiceBeforeSearchScope(t, blogServiceGo())
	edited := strings.Replace(old,
		"func NewBlogService(db *gorm.DB) *BlogService {\n\treturn &BlogService{DB: db}\n}",
		"func NewBlogService(db *gorm.DB) *BlogService {\n\tsvc := &BlogService{DB: db}\n\treturn svc\n}", 1)
	if edited == old {
		t.Fatal("could not edit the constructor")
	}

	out, changes, warnings := repairBlogServiceCountQuery(edited, "{{MODULE}}")
	if out != edited {
		t.Error("the repair changed a file it did not recognise")
	}
	if len(changes) != 0 || len(warnings) != 1 {
		t.Fatalf("want one warning and no changes, got changes %v warnings %v", changes, warnings)
	}
	if !strings.Contains(warnings[0], "CountQuery") {
		t.Errorf("the warning does not say what to add: %q", warnings[0])
	}
}

// blogServiceBeforeSearchScope rebuilds the service as it was before the
// search had one home, which is what an upgrading project still has on disk.
func blogServiceBeforeSearchScope(t *testing.T, fresh string) string {
	t.Helper()

	start := strings.Index(fresh, "// searchScope is what searching a blog means")
	end := strings.Index(fresh, "// List returns a paginated list of all blogs (admin).")
	if start < 0 || end < 0 || end < start {
		t.Fatal("could not find the search pair to remove")
	}
	old := fresh[:start] + fresh[end:]

	old = strings.Replace(old,
		"\tquery := s.searchScope(s.DB.Model(&models.Blog{}), search)\n",
		"\tquery := s.DB.Model(&models.Blog{})\n\n\tif search != \"\" {\n"+
			"\t\t// Portable case-insensitive search — ILIKE is Postgres-only.\n"+
			"\t\tquery = query.Where(\"LOWER(title) LIKE LOWER(?) OR LOWER(content) LIKE LOWER(?)\", \"%\"+search+\"%\", \"%\"+search+\"%\")\n\t}\n", 1)
	old = strings.Replace(old, "\t\"context\"\n", "", 1)

	if strings.Contains(old, "CountQuery") || strings.Contains(old, "searchScope") {
		t.Fatal("could not rebuild the service as it was before")
	}
	return old
}
