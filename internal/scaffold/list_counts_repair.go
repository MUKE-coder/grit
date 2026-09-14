package scaffold

import (
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to count users",
			},
		})
		return
	}
`

// blogListAnchor is the blog admin list's fetch and its error response.
const blogListAnchor = `	blogs, total, pages, err := h.Service.List(page, pageSize, search, sortBy, sortOrder)
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
	// same search as the list.
	countQuery := h.DB.WithContext(c.Request.Context()).Model(&models.Blog{})
	if search != "" {
		countQuery = countQuery.Where("LOWER(title) LIKE LOWER(?) OR LOWER(content) LIKE LOWER(?)", "%"+search+"%", "%"+search+"%")
	}
	counts, err := paginate.Counts(c, countQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to count blogs",
			},
		})
		return
	}
`

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
	return insertListCounts(src, module, "user.go", userListAnchor, userListCounts,
		"after the users total, call paginate.Counts(c, query) and return the result as meta.counts, or the users page's This Week, This Month and Updated Recently cards show a dash")
}

func repairBlogListCountsSource(src, module string) (string, []string, []string) {
	if !strings.Contains(src, "func (h *BlogHandler) List(") {
		return src, nil, nil
	}
	return insertListCounts(src, module, "blog_handler.go", blogListAnchor, blogListCounts,
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
	return out, []string{file + " returns the counts the admin's stat cards ask for on the list request"}, nil
}
