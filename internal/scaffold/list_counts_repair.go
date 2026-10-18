package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// The admin's default stat cards sent a request each: four on every list page
// view and four more after every save. On the users and blog pages three of
// the four were also wrong, since those handlers never read created_since or
// updated_since and answered every window with the grand total. The cards now
// ask for their counts on the list request (?counts=), which paginate.List
// answers for every generated resource. The users and blog handlers build
// their own queries, so they call paginate.Counts, spliced into their
// templates from the text below and added to existing projects by
// repairListCounts.

// userListAnchor is the users list's total count, which the counts follow.
const userListAnchor = "\t// Count total\n\tvar total int64\n\tquery.Count(&total)\n"

const userListCounts = `
	// The admin's stat cards ask for their counts on this request, over the
	// same search as the total.
	counts, err := paginate.Counts(c, query)
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Failed to count users")
		return
	}
`

// blogListAnchor is the blog admin list's fetch and its error response.
const blogListAnchor = `	blogs, total, pages, err := h.Service.List(page, pageSize, search, sortBy, sortOrder)
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Failed to fetch blogs")
		return
	}
`

// blogListAnchorLegacy is the same fetch in a project written before the error
// envelopes moved to respond.Fail. Kept so a project that skipped several
// releases still gets the counts rather than a warning.
const blogListAnchorLegacy = `	blogs, total, pages, err := h.Service.List(page, pageSize, search, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to fetch blogs",
			},
		})
		return
	}
`

const blogListCounts = `
	// The admin's stat cards ask for their counts on this request, over the
	// same search as the list. The query comes from the service so that the
	// search is defined in one place rather than rebuilt here.
	counts, err := paginate.Counts(c, h.Service.CountQuery(c.Request.Context(), search))
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Failed to count blogs")
		return
	}
`

// blogServiceAnchor is the blog service's constructor, which the search pair
// follows.
const blogServiceAnchor = `func NewBlogService(db *gorm.DB) *BlogService {
	return &BlogService{DB: db}
}
`

// blogServiceCountQuery gives an existing project's blog service the single
// definition of its search.
//
// It is here rather than only in the template because the blog files are
// written for a new project and do not travel on upgrade, and the handler
// block below calls CountQuery. Delivering one without the other would be a
// project that no longer builds.
const blogServiceCountQuery = `
// searchScope is what searching a blog means, written once.
//
// The list and the admin's stat cards both narrow by the same text, and when
// each built its own WHERE they could answer different questions: add a column
// to one clause and the table filtered by it while the counts above it did
// not. One definition, two callers, so that cannot happen.
func (s *BlogService) searchScope(query *gorm.DB, search string) *gorm.DB {
	if search == "" {
		return query
	}
	// Portable case-insensitive search, ILIKE is Postgres-only.
	return query.Where("LOWER(title) LIKE LOWER(?) OR LOWER(content) LIKE LOWER(?)", "%"+search+"%", "%"+search+"%")
}

// CountQuery is the list's query for a search, without paging or ordering.
//
// It exists for the handler's stat cards, which need to count the same rows
// the list is about to return. Exported so the handler can ask the service for
// it rather than reaching past the service to the database and rebuilding the
// search by hand, which is the version that drifts.
func (s *BlogService) CountQuery(ctx context.Context, search string) *gorm.DB {
	return s.searchScope(s.DB.WithContext(ctx).Model(&models.Blog{}), search)
}
`

// repairBlogServiceCountQuery adds that pair, and moves List onto it so the
// clause really does exist once rather than twice.
func repairBlogServiceCountQuery(src, module string) (string, []string, []string) {
	if !strings.Contains(src, "func NewBlogService(") {
		return src, nil, nil
	}
	if strings.Contains(src, "func (s *BlogService) CountQuery(") {
		return src, nil, nil
	}
	if strings.Count(src, blogServiceAnchor) != 1 {
		return src, nil, []string{"blog_service.go is not the file Grit wrote: add CountQuery(ctx, search) returning the list's query, or the blog stat cards keep their own copy of the search"}
	}
	at := strings.Index(src, blogServiceAnchor) + len(blogServiceAnchor)
	out := src[:at] + blogServiceCountQuery + src[at:]

	// List keeping its own inline clause would leave two definitions, which is
	// the thing being fixed. Only rewrite the one this file shipped with.
	const listClause = `	query := s.DB.Model(&models.Blog{})

	if search != "" {
		// Portable case-insensitive search — ILIKE is Postgres-only.
		query = query.Where("LOWER(title) LIKE LOWER(?) OR LOWER(content) LIKE LOWER(?)", "%"+search+"%", "%"+search+"%")
	}
`
	if strings.Count(out, listClause) == 1 {
		out = strings.Replace(out, listClause, "\tquery := s.searchScope(s.DB.Model(&models.Blog{}), search)\n", 1)
	}

	var ok bool
	if out, ok = addImportGroup(out, "context"); !ok {
		return src, nil, []string{"could not add the context import to blog_service.go"}
	}
	return out, []string{"blog_service.go defines the blog search once, for the list and the stat cards"}, nil
}

// listCountsMetaAnchor is the last key of a hand-written list's meta map, and
// listCountsMeta the key that follows it.
const (
	listCountsMetaAnchor = "\t\t\t\"pages\":     pages,\n"
	listCountsMeta       = "\t\t\t\"counts\":    counts,\n"
)

// repairListCounts gives the users and blog list handlers of an existing
// project the counts the admin's stat cards now ask for. It runs after the
// codegen runtime, which delivers paginate.Counts.
func repairListCounts(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	if !fileContains(filepath.Join(apiRoot, "internal", "paginate", "paginate.go"), "func Counts(") {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	module := opts.Module()

	// The service first, because the blog handler's block calls CountQuery.
	// If this leaves the service alone the handler is skipped below, rather
	// than written with a call to a method that is not there.
	servicePath := filepath.Join(apiRoot, "internal", "services", "blog_service.go")
	if fileExists(servicePath) {
		if err := repairSourceFile(root, m, servicePath, func(src string) (string, []string, []string) {
			return repairBlogServiceCountQuery(src, module)
		}); err != nil {
			return err
		}
	}
	blogServiceReady := fileContains(servicePath, "func (s *BlogService) CountQuery(")

	handlers := []struct {
		file string
		fix  func(src, module string) (string, []string, []string)
	}{
		{"user.go", repairUserListCountsSource},
		{"blog_handler.go", repairBlogListCountsSource},
	}
	for _, h := range handlers {
		path := filepath.Join(apiRoot, "internal", "handlers", h.file)
		if !fileExists(path) {
			continue
		}
		if h.file == "blog_handler.go" && !blogServiceReady {
			fmt.Printf("  ⚠ %s: left alone, because blog_service.go has no CountQuery for it to call. Add CountQuery(ctx, search) there and run grit upgrade again, or the blog page's stat cards show a dash\n", h.file)
			continue
		}
		fix := h.fix
		if err := repairSourceFile(root, m, path, func(src string) (string, []string, []string) {
			return fix(src, module)
		}); err != nil {
			return err
		}
	}
	return nil
}

func repairUserListCountsSource(src, module string) (string, []string, []string) {
	if !strings.Contains(src, "func (h *UserHandler) List(") {
		return src, nil, nil
	}
	// The users list went through paginate.List in v3.281.0, and onto
	// UserService in v3.340.0. paginate.List answers ?counts= itself either way:
	// nothing to splice, and nothing to warn about.
	if strings.Contains(src, "paginate.List[models.User]") || strings.Contains(src, "h.users().List(") {
		return src, nil, nil
	}
	return insertListCounts(src, module, "user.go", userListAnchor, userListCounts,
		"after the users total, call paginate.Counts(c, query) and return the result as meta.counts, or the users page's This Week, This Month and Updated Recently cards show a dash")
}

func repairBlogListCountsSource(src, module string) (string, []string, []string) {
	if !strings.Contains(src, "func (h *BlogHandler) List(") {
		return src, nil, nil
	}
	anchor := blogListAnchor
	if !strings.Contains(src, anchor) && strings.Contains(src, blogListAnchorLegacy) {
		anchor = blogListAnchorLegacy
	}
	return insertListCounts(src, module, "blog_handler.go", anchor, blogListCounts,
		"in the admin List, call paginate.Counts over the blog query and return the result as meta.counts, or the blog page's This Week, This Month and Updated Recently cards show a dash")
}

// insertListCounts adds block after anchor, and the counts key to the first
// meta map after it.
func insertListCounts(src, module, file, anchor, block, what string) (string, []string, []string) {
	if strings.Contains(src, "paginate.Counts(") {
		return src, nil, nil
	}
	if strings.Count(src, anchor) != 1 {
		return src, nil, []string{file + " is not the file Grit wrote: " + what}
	}
	cut := strings.Index(src, anchor) + len(anchor)
	rest := src[cut:]
	meta := strings.Index(rest, listCountsMetaAnchor)
	if meta < 0 {
		return src, nil, []string{file + " is not the file Grit wrote: " + what}
	}
	meta += len(listCountsMetaAnchor)
	out := src[:cut] + block + rest[:meta] + listCountsMeta + rest[meta:]
	var ok bool
	if out, ok = addImportGroup(out, module+"/internal/paginate"); !ok {
		return src, nil, []string{"could not add the paginate import to " + file}
	}
	// The block answers a failed count with respond.Fail, so a handler written
	// before the envelopes moved there needs the import as well.
	if out, ok = addImportGroup(out, module+"/internal/respond"); !ok {
		return src, nil, []string{"could not add the respond import to " + file}
	}
	return out, []string{file + " returns the counts the admin's stat cards ask for on the list request"}, nil
}
