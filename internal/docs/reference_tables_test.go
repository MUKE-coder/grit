package docs

import (
	"strings"
	"testing"
)

// The reference pages render their tables from a const array and a .map(), so
// every <td> holds an expression and the text lives in the array. The indexer
// read <p>, <li> and <blockquote> and stripped expressions, so those tables
// were indexed as nothing: the field-types page called itself the single
// source of truth for field mapping and contained no field types (grit#99).
//
// This is the page every agent is told to read before generating a resource,
// and the instruction was useless. These assertions are the content that has
// to survive, not the mechanism that currently produces it.
func TestFieldTypesPageCarriesItsMappings(t *testing.T) {
	page, err := Read("/docs/concepts/field-types")
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{
		// The type mapping, which is the reason to open the page.
		"money | price:money",
		"string | name:string",
		"belongs_to",
		// The modifier tags.
		"unique",
		// The flags, which are a different table further down the page and
		// were the first thing lost when prose and tables shared one budget.
		"--owned-by",
		"--append-only",
		"--faker",
	} {
		if !strings.Contains(page.Text, want) {
			t.Errorf("the field-types page does not carry %q", want)
		}
	}
}

// A reference page that indexes to prose alone is the failure above, wearing a
// different page's name. Any page whose job is a table should carry rows.
func TestReferencePagesCarryRows(t *testing.T) {
	for _, url := range []string{
		"/docs/concepts/field-types",
	} {
		page, err := Read(url)
		if err != nil {
			t.Errorf("%s: %v", url, err)
			continue
		}
		if !strings.Contains(page.Text, " | ") {
			t.Errorf("%s indexes no table rows, so looking a value up in it "+
				"returns prose about values", url)
		}
	}
}
