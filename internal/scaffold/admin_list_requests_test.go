package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

// A resource list page sent five requests on every view (the list and a
// request per stat card), one per keystroke in the search box with the table
// blanking to a skeleton each time, and five more after every save.

func TestResourceListAsksOncePerSearchAndKeepsRows(t *testing.T) {
	hook := adminUseResource()
	for _, want := range []string{
		"placeholderData: keepPreviousData",
		"queryFn: async ({ signal })",
		", { signal });",
		"export function useDebouncedValue<T>(",
		`searchParams.set("counts", counts.join(","))`,
	} {
		if !strings.Contains(hook, want) {
			t.Errorf("use-resource is missing %q", want)
		}
	}
	if n := strings.Count(hook, "refreshAfterSave(queryClient, endpoint);"); n != 2 {
		t.Errorf("update and patch should both refresh without refetching the stat cards; found %d", n)
	}

	controller := adminUseResourceController()
	for _, want := range []string{
		"useDebouncedValue(search, 300)",
		"search: debouncedSearch,",
		"counts: statsEnabled && !customStatCards ? DEFAULT_STAT_COUNTS : undefined,",
		"isFetching: isFetching && !isLoading,",
	} {
		if !strings.Contains(controller, want) {
			t.Errorf("the resource controller is missing %q", want)
		}
	}
	if strings.Contains(controller, "page_size=1&created_since") || strings.Contains(controller, `endpoint: ep + "?page_size=1"`) {
		t.Error("the default stat cards still send a request each")
	}

	if !strings.Contains(adminDataTable(), "{isFetching && (") {
		t.Error("the table has no refetch indicator, so a slow search looks like nothing happened")
	}
	if !strings.Contains(adminPageHeader(), "|| stat.loading ?") {
		t.Error("a stat card with a static value has no loading state")
	}
}

func TestPaginateAnswersCounts(t *testing.T) {
	src := apiPaginateGo()
	for _, want := range []string{"func Counts(c *gin.Context, query *gorm.DB)", `"counts": true,`, "result.Meta.Counts = counts"} {
		if !strings.Contains(src, want) {
			t.Errorf("paginate.go is missing %q", want)
		}
	}
}

func TestHandWrittenListsReturnCounts(t *testing.T) {
	cases := []struct {
		name   string
		fresh  string
		block  string
		repair func(string, string) (string, []string, []string)
	}{
		{"blog_handler.go", blogHandlerGo(), blogListCounts, repairBlogListCountsSource},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if strings.Count(tc.fresh, "paginate.Counts(") != 1 || strings.Count(tc.fresh, listCountsMeta) != 1 {
				t.Fatal("the template does not return the counts")
			}
			if out, changes, warnings := tc.repair(tc.fresh, "{{MODULE}}"); out != tc.fresh || len(changes)+len(warnings) != 0 {
				t.Errorf("a fresh %s needs no repair, got changes %v warnings %v", tc.name, changes, warnings)
			}

			old := strings.Replace(tc.fresh, tc.block, "", 1)
			old = strings.Replace(old, listCountsMeta, "", 1)
			old = strings.Replace(old, "\t\"{{MODULE}}/internal/paginate\"\n", "", 1)
			if strings.Contains(old, "paginate.Counts(") || strings.Contains(old, "/internal/paginate\"") {
				t.Fatal("could not rebuild the file as it was before")
			}
			out, changes, warnings := tc.repair(old, "{{MODULE}}")
			if len(warnings) != 0 || len(changes) != 1 {
				t.Fatalf("repairing the previous %s: changes %v warnings %v", tc.name, changes, warnings)
			}
			if strings.Count(out, "paginate.Counts(") != 1 || strings.Count(out, listCountsMeta) != 1 {
				t.Error("the repaired file does not return the counts")
			}
			if _, err := format.Source([]byte(strings.ReplaceAll(out, "{{MODULE}}", "example.com/app"))); err != nil {
				t.Errorf("the repaired file is not valid Go: %v", err)
			}

			edited := strings.Replace(old, "Failed to fetch blogs", "Could not fetch blogs", 1)
			if _, _, warnings := tc.repair(edited, "{{MODULE}}"); len(warnings) != 1 {
				t.Errorf("an edited %s should be left alone with a warning, got %v", tc.name, warnings)
			}
		})
	}
}

// The users list stopped splicing paginate.Counts in v3.281.0 because it stopped
// counting by hand: paginate.List answers ?counts= for it, the way it does for
// every generated resource. It also reads the filters the admin's users page
// sends, which the hand-rolled version did not.
func TestUsersListGoesThroughPaginate(t *testing.T) {
	src := apiUserHandlerGo()
	for _, want := range []string{
		"paginate.List[models.User](",
		"userListConfig",
		`Filterable:   map[string]bool{"role": true, "active": true, "provider": true}`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the users list is missing %q", want)
		}
	}
	for _, gone := range []string{"allowedSorts", "offset := (page - 1) * pageSize", "query.Count(&total)"} {
		if strings.Contains(src, gone) {
			t.Errorf("the users list still hand-rolls %q", gone)
		}
	}
	// repairListCounts must leave it alone rather than warn about a missing anchor.
	if out, fixed, warnings := repairUserListCountsSource(src, "{{MODULE}}"); out != src || len(fixed)+len(warnings) != 0 {
		t.Errorf("a paginate-backed users list needs no counts repair, got fixed %v warnings %v", fixed, warnings)
	}
}
