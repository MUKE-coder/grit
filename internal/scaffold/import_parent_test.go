package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// grit#129: the importer created a missing belongs_to parent from its natural
// key alone, so a parent model requiring anything else failed on the
// constraint once per row. In the reporter's CRM that was 4 of 4 rows, and in
// a real contact list it is nearly every row.
//
// What decides it is the related model on disk, so that is what these read.

const companyNeedsAnOwner = `package models

type Company struct {
	ID        string    ` + "`gorm:\"primaryKey\" json:\"id\"`" + `
	Name      string    ` + "`gorm:\"size:200;not null\" json:\"name\" binding:\"required\"`" + `
	OwnerID   string    ` + "`gorm:\"size:36;not null\" json:\"owner_id\" binding:\"required\"`" + `
	Notes     string    ` + "`json:\"notes\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
}
`

const groupNeedsOnlyAName = `package models

type Group struct {
	ID   string ` + "`gorm:\"primaryKey\" json:\"id\"`" + `
	Name string ` + "`gorm:\"size:200;not null\" json:\"name\" binding:\"required\"`" + `
}
`

const statusHasADefault = `package models

type Status struct {
	ID     string ` + "`gorm:\"primaryKey\" json:\"id\"`" + `
	Name   string ` + "`json:\"name\" binding:\"required\"`" + `
	Colour string ` + "`gorm:\"default:'grey'\" json:\"colour\" binding:\"required\"`" + `
}
`

func TestUnfillableRequiredColumns(t *testing.T) {
	for _, tc := range []struct {
		name string
		src  string
		rel  string
		key  string
		want []string
	}{
		{
			// The reporter's case. owner_id is required with a foreign key, so
			// creating the company from its name alone cannot succeed.
			name: "a required column the CSV has no place for",
			src:  companyNeedsAnOwner, rel: "Company", key: "name",
			want: []string{"owner_id"},
		},
		{
			// The case auto-create exists for, and must keep working.
			name: "a lookup table whose only column is a name",
			src:  groupNeedsOnlyAName, rel: "Group", key: "name",
			want: nil,
		},
		{
			// Required to the API, supplied by the database.
			name: "a required column with a default",
			src:  statusHasADefault, rel: "Status", key: "name",
			want: nil,
		},
		{
			// Without the natural key excluded, every parent looks unfillable
			// and nothing would auto-create at all.
			name: "the natural key does not count against itself",
			src:  groupNeedsOnlyAName, rel: "Group", key: "",
			want: []string{"name"},
		},
		{
			name: "a model this file does not declare",
			src:  groupNeedsOnlyAName, rel: "Company", key: "name",
			want: nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := unfillableRequiredColumns(tc.src, tc.rel, tc.key)
			if strings.Join(got, ",") != strings.Join(tc.want, ",") {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

// The reader works off the project, which is what both callers pass it.
func TestImportUnfillableColumnsReadsTheProject(t *testing.T) {
	api := t.TempDir()
	models := filepath.Join(api, "internal", "models")
	if err := os.MkdirAll(models, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(models, "company.go"), []byte(companyNeedsAnOwner), 0o644); err != nil {
		t.Fatal(err)
	}

	if got := ImportUnfillableColumns(api, "Company", "name"); strings.Join(got, ",") != "owner_id" {
		t.Errorf("got %v, want [owner_id]", got)
	}
	// A model that is not there at all must not make every parent unfillable:
	// that would turn a missing file into a refused import.
	if got := ImportUnfillableColumns(api, "Nowhere", "name"); got != nil {
		t.Errorf("a missing model file produced %v", got)
	}
}

// One function picks the resolver, so the generator and the upgrade cannot
// answer differently. That they could is how two copies of generated code
// start to disagree, which happened three times this week.
func TestImportNameResolverForPicksByWhatTheParentNeeds(t *testing.T) {
	creating := ImportNameResolverFor("group", "Group", "name", "Name", nil)
	if !strings.Contains(creating, "rel = models.Group{Name: v}") {
		t.Error("a name-only parent is no longer created, so a lookup table import " +
			"that worked before now fails every row (grit#129)")
	}

	refusing := ImportNameResolverFor("company", "Company", "name", "Name", []string{"owner_id"})
	if strings.Contains(refusing, "db.Create(&rel)") {
		t.Error("a parent the importer cannot fill in is still created, which is a " +
			"constraint violation per row (grit#129)")
	}
	for _, want := range []string{
		"does not exist",
		"Company requires owner_id",
		"Add the Company first",
	} {
		if !strings.Contains(refusing, want) {
			t.Errorf("the refusal does not say %q, so the person importing cannot act on it", want)
		}
	}
}

// And the upgrade reaches a project that already has the resolver that cannot
// work, because that resolver does not become correct by being old.
func TestUpgradeReplacesACreatingResolverItCannotFillIn(t *testing.T) {
	api := t.TempDir()
	models := filepath.Join(api, "internal", "models")
	if err := os.MkdirAll(models, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, src := range map[string]string{
		"company.go": companyNeedsAnOwner,
		"group.go":   groupNeedsOnlyAName,
	} {
		if err := os.WriteFile(filepath.Join(models, name), []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	src := ImportNameResolver("company", "Company", "name", "Name") +
		ImportNameResolver("group", "Group", "name", "Name")

	out := src
	n, names := replaceCreatingResolvers(&out, api)
	if n != 1 {
		t.Fatalf("replaced %d resolvers, want 1 (Company only): %v", n, names)
	}
	if strings.Contains(out, "rel = models.Company{Name: v}") {
		t.Error("the Company resolver still creates a company it cannot fill in")
	}
	if !strings.Contains(out, "rel = models.Group{Name: v}") {
		t.Error("the Group resolver was replaced too, so a lookup table that " +
			"imported correctly before now refuses every new value")
	}
	if len(names) != 1 || !strings.Contains(names[0], "owner_id") {
		t.Errorf("the upgrade report does not name the column: %v", names)
	}

	// Upgrading twice must not change it again.
	again := out
	if n, _ := replaceCreatingResolvers(&again, api); n != 0 || again != out {
		t.Error("a second upgrade rewrote the resolver again")
	}
}
