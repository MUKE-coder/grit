package docscheck

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

/*
 * "8 field types" was on the Next.js prerequisites page, naming them: text,
 * textarea, number, select, date, toggle, checkbox, radio. The admin's
 * FormBuilder takes 31. The architecture page said "10+", which is true in
 * the way that "2+" is true, and reads as a hedge from somebody who did not
 * want to count.
 *
 * A reader deciding whether the admin can render their form counts on that
 * number. Eight is a reason to write the form by hand.
 *
 * The count lives in one place that is not prose: the FieldType union in
 * admin_resource_files.go, which is what a scaffolded project gets as
 * lib/resource.ts. So the claim can be compared to it rather than re-counted
 * by hand the next time somebody wonders.
 */

var (
	fieldTypeUnion = regexp.MustCompile(`export type FieldType = ([^;]+);`)
	quotedName     = regexp.MustCompile(`"[a-z-]+"`)
	// "31 field types", "10+ field types", "8 field types".
	documentedCount = regexp.MustCompile(`\b(\d{1,3})\+? field types\b`)
)

// formBuilderFieldTypes counts the union the scaffold writes.
func formBuilderFieldTypes(t *testing.T) int {
	t.Helper()

	path := filepath.Join("..", "scaffold", "admin_resource_files.go")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}

	m := fieldTypeUnion.FindSubmatch(body)
	if m == nil {
		t.Fatalf("no `export type FieldType = ...` in %s. If the union moved, "+
			"move this test with it rather than deleting it: the docs state "+
			"this number and nothing else compares them.", path)
	}

	names := map[string]bool{}
	for _, q := range quotedName.FindAllString(string(m[1]), -1) {
		names[q] = true
	}
	if len(names) == 0 {
		t.Fatalf("read the FieldType union and found no type names in it")
	}
	return len(names)
}

func TestDocsAgreeOnHowManyFieldTypesTheFormBuilderTakes(t *testing.T) {
	want := formBuilderFieldTypes(t)

	docsRoot := filepath.Join("..", "..", "docs", "app", "docs")
	var problems []string

	err := filepath.Walk(docsRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || info.Name() != "page.tsx" {
			return nil
		}
		// The changelog records what was true at the time.
		if strings.Contains(filepath.ToSlash(path), "/changelog/") {
			return nil
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		for _, m := range documentedCount.FindAllSubmatch(body, -1) {
			got, convErr := strconv.Atoi(string(m[1]))
			if convErr != nil {
				continue
			}
			if got != want {
				rel, _ := filepath.Rel(filepath.Join("..", ".."), path)
				problems = append(problems, filepath.ToSlash(rel)+
					": says "+string(m[0])+", the FormBuilder takes "+
					strconv.Itoa(want))
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking the docs: %v", err)
	}

	for _, p := range problems {
		t.Error(p)
	}
}
