package main

import (
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/generate"
)

// grit#140: the generated importer writes with clause.OnConflict{DoNothing},
// which does nothing unless the database has a constraint to conflict with. A
// generated resource had a unique index on its primary key alone and the
// importer assigns a fresh id per row, so nothing ever conflicted: the same
// CSV imported twice produced two of every row, and "skipped 0" was
// structural rather than a fact.
//
// --unique names the columns, because which ones is a decision per resource
// that no generator can make. These are the refusals: a unique index that
// fails to create leaves a resource whose importer silently duplicates every
// row, which is the bug the flag exists to fix, so the spellings that cannot
// work are refused at generation time rather than at migrate time.

func uniqueDef() *generate.ResourceDefinition {
	return &generate.ResourceDefinition{
		Name: "Contact",
		Fields: []generate.Field{
			{Name: "name", Type: "string"},
			{Name: "email", Type: "email"},
			{Name: "notes", Type: "richtext"},
			{Name: "price", Type: "money"},
			{Name: "photos", Type: "files"},
			{Name: "company", Type: "belongs_to", RelatedModel: "Company"},
			{Name: "tags", Type: "many_to_many", RelatedModel: "Tag"},
		},
	}
}

func TestResolveUniqueColumns(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
		want  []string
		errIs string
	}{
		{"one column", "email", []string{"email"}, ""},
		{"a composite, in the order named", "name,company", []string{"name", "company"}, ""},
		{"a relation by either spelling", "company_id", []string{"company_id"}, ""},
		{"spaces are trimmed", " name , email ", []string{"name", "email"}, ""},

		{"a column that is not there", "nope", nil, "not a field"},
		{"a many-to-many has no column here", "tags", nil, "join table"},
		{"money is two columns", "price", nil, "two\ncolumns"},
		{"a list in one column", "photos", nil, "different order"},
		{"rich text is too long to index", "notes", nil, "refuse to index"},
		{"the same column twice", "email,email", nil, "twice"},
		{"the same relation twice, spelled both ways", "company,company_id", nil, "twice"},
		{"nothing at all", ",,", nil, "no columns"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := resolveUniqueColumns(uniqueDef(), tc.input)
			if tc.errIs != "" {
				if err == nil {
					t.Fatalf("%q was accepted, and the migration would fail instead", tc.input)
				}
				// The message has to say which column and why, because the
				// alternative is somebody reading a migration error.
				for _, part := range strings.Split(tc.errIs, "\n") {
					if !strings.Contains(err.Error(), part) {
						t.Errorf("the refusal does not say %q: %v", part, err)
					}
				}
				return
			}
			if err != nil {
				t.Fatalf("%q was refused: %v", tc.input, err)
			}
			if strings.Join(got, ",") != strings.Join(tc.want, ",") {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

// The error names what the caller could have written instead. Without it the
// only way to find a column's spelling is to read the generated model.
func TestResolveUniqueColumnsListsTheColumns(t *testing.T) {
	_, err := resolveUniqueColumns(uniqueDef(), "emial")
	if err == nil {
		t.Fatal("a typo was accepted")
	}
	for _, want := range []string{"email", "company", "name"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the error does not list %q as an available column: %v", want, err)
		}
	}
}
