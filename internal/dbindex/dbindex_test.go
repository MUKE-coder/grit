package dbindex

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func project(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range files {
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

const plainService = `package services
var dealConfig = paginate.Config{
	Sortable:   map[string]bool{"id": true, "created_at": true, "title": true},
	Filterable: map[string]bool{"title": true},
}
`

// The scope goes first in the index because it is an equality the planner can
// seek to, and it is different for each ownership model. Getting this wrong
// produces an index the query cannot use, which looks like having solved the
// problem.
func TestTheScopeComesFromTheModel(t *testing.T) {
	cases := []struct {
		name  string
		model string
		scope string
	}{
		{
			name: "tenant-owned",
			model: `package models
type Deal struct {
	tenant.Owned
	ID string ` + "`gorm:\"primarykey\"`" + `
}`,
			scope: "org_id",
		},
		{
			name: "owned-by",
			model: `package models
type Deal struct {
	ID     string ` + "`gorm:\"primarykey\" json:\"id\"`" + `
	UserID string ` + "`gorm:\"size:36;index\" json:\"user_id\" binding:\"required\" owner:\"true\"`" + `
}`,
			scope: "user_id",
		},
		{
			name: "neither",
			model: `package models
type Deal struct {
	ID    string ` + "`gorm:\"primarykey\" json:\"id\"`" + `
	Title string ` + "`json:\"title\"`" + `
}`,
			scope: "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			root := project(t, map[string]string{
				"internal/services/deal.go": plainService,
				"internal/models/deal.go":   c.model,
			})
			got, err := Scan(root)
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != 1 {
				t.Fatalf("expected one resource, got %d", len(got))
			}
			if got[0].Scope != c.scope {
				t.Errorf("scope is %q, want %q", got[0].Scope, c.scope)
			}

			idx := Plan(got, false, nil)
			if len(idx) != 1 {
				t.Fatalf("expected one index, got %d", len(idx))
			}
			want := []string{"created_at", "id"}
			if c.scope != "" {
				want = append([]string{c.scope}, want...)
			}
			if strings.Join(idx[0].Columns, ",") != strings.Join(want, ",") {
				t.Errorf("columns %v, want %v", idx[0].Columns, want)
			}
		})
	}
}

// id last, always. Every list orders by the column and then by id to break
// ties (grit#112), and an index without that trailing column cannot serve the
// whole ORDER BY, which makes it an index that costs writes and buys nothing.
func TestEveryIndexEndsWithID(t *testing.T) {
	root := project(t, map[string]string{
		"internal/services/deal.go": plainService,
		"internal/models/deal.go":   "package models\ntype Deal struct{}",
	})
	res, _ := Scan(root)
	for _, idx := range Plan(res, true, nil) {
		if idx.Columns[len(idx.Columns)-1] != "id" {
			t.Errorf("%s ends with %s", idx.Name, idx.Columns[len(idx.Columns)-1])
		}
	}
}

// Only the default sort by default. An index per sortable column is write
// amplification on a table that is mostly written, and nothing at generation
// time knows which columns anybody actually sorts by.
func TestOnlyTheDefaultSortUnlessAsked(t *testing.T) {
	root := project(t, map[string]string{
		"internal/services/deal.go": plainService,
		"internal/models/deal.go":   "package models\ntype Deal struct{}",
	})
	res, _ := Scan(root)

	if n := len(Plan(res, false, nil)); n != 1 {
		t.Errorf("default plan has %d indexes, want 1 (created_at)", n)
	}
	// id is already the primary key, so it never gets one of its own.
	all := Plan(res, true, nil)
	if n := len(all); n != 2 {
		t.Errorf("--all-sortable has %d indexes, want 2 (created_at, title)", n)
	}
	for _, idx := range all {
		if strings.HasSuffix(idx.Name, "_id_id") {
			t.Errorf("%s indexes the primary key", idx.Name)
		}
	}

	picked := Plan(res, false, map[string][]string{"deal": {"title"}})
	if len(picked) != 1 || picked[0].Columns[0] != "title" {
		t.Errorf("--only did not pick the column: %v", picked)
	}
}

// A real table name wins over the guess, because TableName() is the one the
// database actually uses and an index on a table that does not exist fails at
// the worst moment.
func TestTableNameOverridesTheGuess(t *testing.T) {
	root := project(t, map[string]string{
		"internal/services/person.go": strings.ReplaceAll(plainService, "dealConfig", "personConfig"),
		"internal/models/person.go": `package models
type Person struct{}
func (Person) TableName() string { return "people" }`,
	})
	res, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	if res[0].Table != "people" {
		t.Errorf("table is %q, want people", res[0].Table)
	}
}

// The SQL has to be applicable. CONCURRENTLY cannot run inside a transaction,
// which is the whole reason this is a file rather than something AutoMigrate
// does, so the file says so.
func TestTheSQLSaysHowToApplyIt(t *testing.T) {
	root := project(t, map[string]string{
		"internal/services/deal.go": plainService,
		"internal/models/deal.go":   "package models\ntype Deal struct{}",
	})
	res, _ := Scan(root)
	sql := SQL(Plan(res, false, nil), "v1")

	for _, want := range []string{
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS",
		"cannot run inside a transaction",
		"psql",
		"Postgres only",
	} {
		if !strings.Contains(sql, want) {
			t.Errorf("the SQL file does not mention %q", want)
		}
	}
}

// A service with no Sortable map is not a list, and an empty project is not
// an error.
func TestNonListServicesAreIgnored(t *testing.T) {
	root := project(t, map[string]string{
		"internal/services/mail.go": "package services\nfunc Send() {}",
	})
	res, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 0 {
		t.Errorf("found %d resources in a project with no lists", len(res))
	}
	if !strings.Contains(SQL(nil, "v1"), "Nothing to create") {
		t.Error("an empty plan should say so rather than emit a bare header")
	}
}
