package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// grit#140: the model needs the unique index for the importer's
// OnConflict{DoNothing} to mean anything.

func uniqueModelSource(t *testing.T, unique []string, fields []Field) string {
	t.Helper()
	root := t.TempDir()
	g := &Generator{
		Root:   root,
		Module: "example",
		Definition: &ResourceDefinition{
			Name:   "Contact",
			Fields: fields,
			Unique: unique,
		},
	}
	if err := g.Run(); err != nil {
		t.Fatalf("generate: %v", err)
	}
	body, err := os.ReadFile(filepath.Join(g.APIRoot(), "internal", "models", "contact.go"))
	if err != nil {
		t.Fatalf("reading the model: %v", err)
	}
	return string(body)
}

// One column, the simple case: a contact unique by email.
func TestUniqueIndexOnOneColumn(t *testing.T) {
	src := uniqueModelSource(t, []string{"email"}, []Field{
		{Name: "name", Type: "string"},
		{Name: "email", Type: "email"},
	})

	if !strings.Contains(src, "uniqueIndex:idx_contact_unique,priority:1") {
		t.Error("the email column has no unique index, so the importer's " +
			"OnConflict has nothing to conflict with and a re-imported CSV " +
			"duplicates every row (grit#140)")
	}
	// And only that column.
	if n := strings.Count(src, "idx_contact_unique"); n != 1 {
		t.Errorf("the index covers %d columns, want 1", n)
	}
}

// A composite, in the order named: unique by name WITHIN a company. The
// priorities are what decide that, and a foreign key can carry one.
func TestUniqueIndexAcrossAColumnAndARelation(t *testing.T) {
	src := uniqueModelSource(t, []string{"name", "company"}, []Field{
		{Name: "name", Type: "string"},
		{Name: "company", Type: "belongs_to", RelatedModel: "Company"},
	})

	for _, want := range []string{
		"uniqueIndex:idx_contact_unique,priority:1",
		"uniqueIndex:idx_contact_unique,priority:2",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the model is missing %q", want)
		}
	}
	// The foreign key keeps its ordinary index as well: the unique one is
	// composite, so it does not serve a lookup by company alone.
	// Checked as a fragment, not a whole tag: an optional relation adds
	// default:null and a required one does not, and neither is what this is
	// about.
	if !strings.Contains(src, "size:36;index;uniqueIndex:idx_contact_unique,priority:2") {
		t.Error("the relation lost its own index, so every lookup by company " +
			"is now a table scan")
	}
}

// Without the flag, nothing changes. This is the half that keeps every
// existing project and every test of one working.
func TestNoUniqueFlagMeansNoIndex(t *testing.T) {
	src := uniqueModelSource(t, nil, []Field{
		{Name: "name", Type: "string"},
		{Name: "email", Type: "email"},
	})
	if strings.Contains(src, "idx_contact_unique") {
		t.Error("a resource generated without --unique got a unique index")
	}
}

// A column with no other gorm tag must not come out with a leading
// semicolon: gorm reads that as an empty setting and the tag is malformed.
func TestUniqueIndexOnAColumnWithNoOtherTag(t *testing.T) {
	src := uniqueModelSource(t, []string{"code"}, []Field{
		{Name: "code", Type: "int"},
	})
	if strings.Contains(src, `gorm:";`) {
		t.Error("the gorm tag starts with a semicolon, which gorm reads as an " +
			"empty setting")
	}
	if !strings.Contains(src, "uniqueIndex:idx_contact_unique,priority:1") {
		t.Error("the column did not get the index")
	}
}
