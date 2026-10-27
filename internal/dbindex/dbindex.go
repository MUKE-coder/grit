// Package dbindex works out which indexes a project's list queries need.
//
// Every generated list is "WHERE <scope> ORDER BY <column>", and nothing was
// creating an index that serves it. On 500,000 rows that is a sort of the
// whole table on disk for every page: 2.3 seconds to return one page, and 16.6
// seconds at page 20,000, measured on Postgres 16 with the database also
// taking writes (grit#112).
//
// Three things decide the index, and all three are already written down in the
// project:
//
//   - The scope comes first, because it is an equality and the planner can
//     seek straight to it. A tenant-owned model is scoped by org_id, an
//     --owned-by one by its owner column, anything else by nothing.
//   - The sort column comes second.
//   - id comes last, because every list now orders by the sort column and
//     then by id to break ties, and an index without that trailing column
//     cannot serve the whole ORDER BY.
//
// It emits SQL rather than creating anything. The statements use CREATE INDEX
// CONCURRENTLY, which does not lock the table and therefore cannot run inside
// a transaction, which is exactly why this is a file you apply rather than
// something AutoMigrate does while somebody is using the product.
package dbindex

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Resource is one list-bearing table, as read out of the project.
type Resource struct {
	Model       string   // Deal
	Table       string   // deals
	Sortable    []string // the columns a caller may sort by
	DefaultSort string   // what a request with no sort_by uses
	Scope       string   // org_id, user_id, or empty
	ScopeReason string   // why, for the comment in the SQL
}

// Index is one statement to create.
type Index struct {
	Table   string
	Name    string
	Columns []string
	Reason  string
}

var (
	sortableRe    = regexp.MustCompile(`Sortable:\s*map\[string\]bool\{([^}]*)\}`)
	keyRe         = regexp.MustCompile(`"([^"]+)"\s*:\s*true`)
	defaultSortRe = regexp.MustCompile(`DefaultSort:\s*"([^"]+)"`)
	tableNameRe   = regexp.MustCompile(`func \((\w+)\) TableName\(\) string \{\s*return "([^"]+)"`)
	ownerColRe    = regexp.MustCompile(`(?m)^\s*(\w+)\s+\*?string\s+` + "`" + `[^` + "`" + `]*json:"([^",]+)"[^` + "`" + `]*owner:"true"`)
)

// Scan reads the API's services and models.
func Scan(apiRoot string) ([]Resource, error) {
	serviceDir := filepath.Join(apiRoot, "internal", "services")
	entries, err := os.ReadDir(serviceDir)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", serviceDir, err)
	}

	var out []Resource
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(serviceDir, e.Name()))
		if err != nil {
			continue
		}
		src := string(data)

		m := sortableRe.FindStringSubmatch(src)
		if m == nil {
			continue // not a list-bearing service
		}
		var sortable []string
		for _, k := range keyRe.FindAllStringSubmatch(m[1], -1) {
			sortable = append(sortable, k[1])
		}
		if len(sortable) == 0 {
			continue
		}

		base := strings.TrimSuffix(e.Name(), ".go")
		r := Resource{
			Model:       pascal(base),
			Table:       plural(base),
			Sortable:    sortable,
			DefaultSort: "created_at",
		}
		if d := defaultSortRe.FindStringSubmatch(src); d != nil {
			r.DefaultSort = d[1]
		}

		// The model says how the table is scoped, and what it is really called.
		modelPath := filepath.Join(apiRoot, "internal", "models", base+".go")
		if data, err := os.ReadFile(modelPath); err == nil {
			model := string(data)
			if t := tableNameRe.FindStringSubmatch(model); t != nil {
				r.Table = t[2]
			}
			switch {
			case strings.Contains(model, "tenant.Owned"):
				r.Scope = "org_id"
				r.ScopeReason = "tenant-owned: every list is filtered by the active organization"
			default:
				if o := ownerColRe.FindStringSubmatch(model); o != nil {
					r.Scope = o[2]
					r.ScopeReason = "--owned-by: every list is filtered by the signed-in user"
				}
			}
		}
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Table < out[j].Table })
	return out, nil
}

// Plan is the indexes for these resources.
//
// allSortable decides the ambition. The default sort is the one every list
// request uses when nobody has clicked a column, so its index pays for itself
// on every page load of every list in the product. The rest are only used when
// somebody sorts by them, and each one costs write throughput on a table that
// is mostly written, so they are opt in and the developer picks them with a
// measurement in hand rather than from a guess made at generation time.
func Plan(resources []Resource, allSortable bool, only map[string][]string) []Index {
	var out []Index
	for _, r := range resources {
		wanted := []string{r.DefaultSort}
		if cols, ok := only[strings.ToLower(r.Model)]; ok {
			wanted = cols
		} else if allSortable {
			wanted = r.Sortable
		}

		for _, col := range wanted {
			if col == "id" {
				// id is the primary key, already unique and already indexed.
				continue
			}
			cols := []string{}
			if r.Scope != "" {
				cols = append(cols, r.Scope)
			}
			cols = append(cols, col, "id")

			reason := fmt.Sprintf("serves ORDER BY %s, id", col)
			if col == r.DefaultSort {
				reason = fmt.Sprintf("the default sort: every list request uses it, ORDER BY %s, id", col)
			}
			if r.ScopeReason != "" {
				reason += ". " + r.ScopeReason
			}
			out = append(out, Index{
				Table:   r.Table,
				Name:    indexName(r.Table, cols),
				Columns: cols,
				Reason:  reason,
			})
		}
	}
	return out
}

// SQL renders the statements.
func SQL(indexes []Index, version string) string {
	var b strings.Builder
	b.WriteString("-- Indexes for the list queries, written by grit index " + version + ".\n")
	b.WriteString("--\n")
	b.WriteString("-- Every generated list is WHERE <scope> ORDER BY <column>, id. Without an\n")
	b.WriteString("-- index in that shape the database sorts the whole table on disk for each\n")
	b.WriteString("-- page, which is invisible until the table is large and then is the slowest\n")
	b.WriteString("-- thing in the product.\n")
	b.WriteString("--\n")
	b.WriteString("-- CONCURRENTLY, so building these does not lock the table. That is also why\n")
	b.WriteString("-- they cannot run inside a transaction: apply this file with psql directly,\n")
	b.WriteString("-- not through a migration runner that wraps everything in BEGIN.\n")
	b.WriteString("--\n")
	b.WriteString("--     psql \"$DATABASE_URL\" -f this-file.sql\n")
	b.WriteString("--\n")
	b.WriteString("-- Postgres only. On MySQL drop CONCURRENTLY; on SQLite drop it and expect\n")
	b.WriteString("-- the table to be locked while each one builds.\n\n")

	if len(indexes) == 0 {
		b.WriteString("-- Nothing to create: no list-bearing resource was found.\n")
		return b.String()
	}

	for _, idx := range indexes {
		fmt.Fprintf(&b, "-- %s\n", idx.Reason)
		fmt.Fprintf(&b, "CREATE INDEX CONCURRENTLY IF NOT EXISTS %s\n  ON %s (%s);\n\n",
			idx.Name, idx.Table, strings.Join(idx.Columns, ", "))
	}
	return b.String()
}

// indexName is deterministic and inside Postgres's 63 byte identifier limit,
// which a long table name plus three columns passes easily.
func indexName(table string, cols []string) string {
	parts := append([]string{"idx", table}, cols...)
	name := strings.Join(parts, "_")
	if len(name) <= 63 {
		return name
	}
	// Trimmed from the middle columns rather than the end, so the table and
	// the trailing id stay readable.
	short := "idx_" + table
	for _, c := range cols {
		if len(short)+1+len(c) > 60 {
			short += "_" + c[:3]
			continue
		}
		short += "_" + c
	}
	if len(short) > 63 {
		short = short[:63]
	}
	return short
}

func plural(snake string) string {
	switch {
	case strings.HasSuffix(snake, "s"), strings.HasSuffix(snake, "x"),
		strings.HasSuffix(snake, "ch"), strings.HasSuffix(snake, "sh"):
		return snake + "es"
	case strings.HasSuffix(snake, "y") && len(snake) > 1 && !isVowel(snake[len(snake)-2]):
		return snake[:len(snake)-1] + "ies"
	default:
		return snake + "s"
	}
}

func isVowel(b byte) bool { return strings.ContainsRune("aeiou", rune(b)) }

func pascal(snake string) string {
	parts := strings.Split(snake, "_")
	for i, p := range parts {
		if p != "" {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}
