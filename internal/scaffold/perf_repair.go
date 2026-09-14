package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// paginateConnectHook installs the list-total cache when the database connects.
const paginateConnectHook = `	// A list's total is reused across its pages until a write to the table, so
	// paging does not count the whole match every time. See internal/paginate.
	if err := paginate.Install(db); err != nil {
		return nil, fmt.Errorf("installing the list total cache: %w", err)
	}

`

// migrateSearchIndexHook follows models.Migrate in cmd/migrate, inside the
// recorded run, so the history lists the indexes and a rollback drops them.
const migrateSearchIndexHook = `	// Postgres gets a trigram index on every column tagged search:"trigram", the
	// ones list search reads.
	if err := paginate.EnsureSearchIndexes(db, models.Models()...); err != nil {
		log.Printf("Search indexes were not created, so search reads the whole table: %v", err)
	}
`

const generatedArchivedMarker = "\tArchivedAt *time.Time  `gorm:\"index\" json:\"archived_at,omitempty\"`\n"

var (
	// generatedArchivedRe finds the field grit generate gives every resource,
	// with the spacing gofmt leaves it.
	generatedArchivedRe  = regexp.MustCompile("(?m)^\\tArchivedAt[ \\t]+\\*time\\.Time[ \\t]+`gorm:\"index\" json:\"archived_at,omitempty\"`$")
	unindexedCreatedAtRe = regexp.MustCompile("(?m)^(\\tCreatedAt[ \\t]+time\\.Time[ \\t]+)`json:\"created_at\"`$")
	modelStructRe        = regexp.MustCompile(`(?m)^type (\w+) struct \{$`)
	searchableListRe     = regexp.MustCompile(`(?s)var (\w+)ListConfig = paginate\.Config\{\s*Searchable:\s*\[\]string\{([^}]*)\}`)
)

// repairListPerformance brings a project up to the fix for H18 in the
// contact-app review: list queries sorted and filtered on columns without an
// index, counted the whole match on every page, and searched with a LIKE no
// index could serve. On a million Postgres rows a page took 77 ms, its count
// 73 ms more, and a search 270 ms.
//
// internal/paginate and cmd/migrate are framework code and arrive whole. The
// rest is the developer's: database.go gets the cache installed, and models get
// an index on created_at and a search tag on the columns their service searches,
// each only where the line is still what Grit wrote. The indexes are built by
// the next grit migrate.
func repairListPerformance(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	if !fileContains(filepath.Join(apiRoot, "internal", "paginate", "count.go"), "func Install(") {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	database := filepath.Join(apiRoot, "internal", "database", "database.go")
	if fileExists(database) {
		module := opts.Module()
		if err := repairSourceFile(root, m, database, func(src string) (string, []string, []string) {
			return repairPaginateWiringSource(src, module)
		}); err != nil {
			return err
		}
	}

	searchable, err := searchableColumns(filepath.Join(apiRoot, "internal", "services"))
	if err != nil {
		return err
	}
	modelFiles, err := filepath.Glob(filepath.Join(apiRoot, "internal", "models", "*.go"))
	if err != nil {
		return err
	}
	for _, path := range modelFiles {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		base := filepath.Base(path)
		if err := repairSourceFile(root, m, path, func(src string) (string, []string, []string) {
			return repairModelIndexesSource(src, base, searchable)
		}); err != nil {
			return err
		}
	}
	return nil
}

func repairPaginateWiringSource(src, module string) (string, []string, []string) {
	if strings.Contains(src, "paginate.Install(db)") || !strings.Contains(src, "func Connect(") {
		return src, nil, nil
	}
	const anchor = "\tsqlDB, err := db.DB()\n"
	if strings.Count(src, anchor) != 1 {
		return src, nil, []string{"database.go is not the file Grit wrote: call paginate.Install(db) after connecting, or every list page counts its whole table"}
	}
	out := strings.Replace(src, anchor, paginateConnectHook+anchor, 1)
	var ok bool
	if out, ok = addImportGroup(out, module+"/internal/paginate"); !ok {
		return src, nil, []string{"could not add the paginate import to database.go: call paginate.Install(db) after connecting"}
	}
	return out, []string{"list totals are reused across pages until the table is written to"}, nil
}

// searchableColumns reads each generated service's Searchable list, keyed by
// its model: contactListConfig belongs to Contact.
func searchableColumns(dir string) (map[string][]string, error) {
	out := map[string][]string{}
	files, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return nil, err
	}
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		for _, match := range searchableListRe.FindAllStringSubmatch(string(data), -1) {
			model := strings.ToUpper(match[1][:1]) + match[1][1:]
			for _, col := range strings.Split(match[2], ",") {
				if col = strings.Trim(strings.TrimSpace(col), `"`); col != "" {
					out[model] = append(out[model], col)
				}
			}
		}
	}
	return out, nil
}

func repairModelIndexesSource(src, base string, searchable map[string][]string) (string, []string, []string) {
	generated := generatedArchivedRe.MatchString(src)
	if !generated && base != "upload.go" && base != "user.go" {
		return src, nil, nil
	}
	out := src
	var fixed []string
	if loc := unindexedCreatedAtRe.FindStringSubmatchIndex(out); loc != nil && (generated || base == "upload.go" || strings.Contains(out, "type User struct {")) {
		out = unindexedCreatedAtRe.ReplaceAllString(out, "${1}`gorm:\"index\" json:\"created_at\"`")
		fixed = append(fixed, "created_at is indexed, so a list sorted by it reads the newest rows first")
	}
	if generated {
		tagged := 0
		for _, match := range modelStructRe.FindAllStringSubmatch(out, -1) {
			for _, col := range searchable[match[1]] {
				fieldRe := regexp.MustCompile("(?m)^(\\t\\w+[ \\t]+\\*?string[ \\t]+`[^`]*json:\"" + regexp.QuoteMeta(col) + "\"[^`]*)`$")
				loc := fieldRe.FindStringSubmatchIndex(out)
				if loc == nil || strings.Contains(out[loc[2]:loc[3]], "search:") {
					continue
				}
				out = out[:loc[3]] + ` search:"trigram"` + out[loc[3]:]
				tagged++
			}
		}
		if tagged > 0 {
			fixed = append(fixed, "the columns its list searches are tagged for a trigram index on Postgres")
		}
	}
	if len(fixed) == 0 {
		return src, nil, nil
	}
	return out, fixed, nil
}
