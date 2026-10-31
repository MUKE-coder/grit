package scaffold

import (
	"sort"
	"strings"
	"testing"
)

// Two things that have to agree, with nothing comparing them until now.
//
// sharedModelTypeExports writes the export lines into types/index.ts in list
// order, as part of the template. repairSharedTypes adds a missing one at the
// // grit:types marker, and addMissingExports sorts by name to do it. The two
// produce the same file only while the list is alphabetical, and the list was
// alphabetical by accident: every entry had been appended in a release whose
// model name happened to come next.
//
// Adding Blog and Upload at the head of it broke a test three files away, with
// a diff of a whole index file and no hint of the cause. This says it instead.
func TestSharedModelTypesIsSortedByName(t *testing.T) {
	names := make([]string, 0, len(SharedModelTypes))
	for _, typ := range SharedModelTypes {
		names = append(names, typ.Name)
	}
	sorted := append([]string(nil), names...)
	sort.Strings(sorted)
	for i := range names {
		if names[i] != sorted[i] {
			t.Fatalf("SharedModelTypes is not sorted by Name, so the template and the repair\n"+
				"would write types/index.ts differently.\n  is:   %v\n  want: %v",
				names, sorted)
		}
	}
}

// The other half of the same agreement: the template's export block and what
// the repair would add for a project that has none of them.
func TestTheTemplateAndTheRepairExportTheSameThings(t *testing.T) {
	fromTemplate := sharedModelTypeExports()

	exports := map[string]string{}
	for _, typ := range SharedModelTypes {
		exports[typ.Name] = "export type { " + typ.Name + " } from \"./" + typ.Kebab + "\";\n"
	}
	// An index with the marker and nothing else: every export is missing, so
	// the repair adds all of them, in the order it would use.
	out, added := addMissingExports("// grit:types\n", "// grit:types", exports)
	if len(added) != len(SharedModelTypes) {
		t.Fatalf("the repair added %d exports for %d models", len(added), len(SharedModelTypes))
	}
	fromRepair := strings.TrimSuffix(out, "// grit:types\n")

	if fromTemplate != fromRepair {
		t.Errorf("a new project and a repaired one would not have the same types/index.ts.\ntemplate:\n%s\nrepair:\n%s",
			fromTemplate, fromRepair)
	}
}

// No name is exported twice. types/index.ts used to name Upload and Blog
// itself and also range over this list, which is two export lines for one
// identifier and fails the type-check of @repo/shared.
func TestNoSharedTypeIsExportedTwice(t *testing.T) {
	index := sharedTypesIndex()
	for _, typ := range SharedModelTypes {
		if n := strings.Count(index, "export type { "+typ.Name+" }"); n > 1 {
			t.Errorf("types/index.ts exports %s %d times", typ.Name, n)
		}
	}
}
