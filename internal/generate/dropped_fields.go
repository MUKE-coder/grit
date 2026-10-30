package generate

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Review 2.23: re-describing an unedited resource with fewer fields dropped
// columns silently.
//
// The guard for edited files is the good half of this: a model somebody has
// changed is refused rather than overwritten. The unedited case is the one
// nobody thought about, and it is the more likely of the two. You generate
// Contact with name, email and company; a week later you run the same command
// from your shell history with name alone, because that is the line you still
// have. The model is rewritten with one field, the handler and the service and
// the admin form follow it, and the only thing that says so is that the output
// is identical to the first run.
//
// AutoMigrate is additive, so the columns are still in the database with the
// data still in them. Nothing reads them any more, and nothing said that.
//
// So: say it. Not a refusal, because regenerating with fewer fields is a
// legitimate thing to do and the files are in git. A line naming each column
// that stopped being read, and what that means.

// frameworkFields are the columns the generator writes for every resource,
// whatever the definition says. They are not the developer's fields and a
// definition never mentions them, so they would otherwise be reported as
// dropped on every regenerate.
var frameworkFields = map[string]bool{
	"ID":         true,
	"Version":    true,
	"CreatedAt":  true,
	"UpdatedAt":  true,
	"DeletedAt":  true,
	"ArchivedAt": true,
	// The owner column, written from --owned-by rather than from a field.
	"UserID": true,
	"User":   true,
	// Multitenancy, added by the plugin rather than by the definition.
	"OrganizationID": true,
	"Organization":   true,
}

// reportDroppedFields prints the columns a regenerate is about to stop reading.
//
// Called before the model is written, because afterwards the old field list is
// gone. Silent for a resource being generated for the first time, and silent
// when nothing is dropped, which is the overwhelming majority of runs.
func (g *Generator) reportDroppedFields(names Names) {
	path := filepath.Join(g.APIRoot(), "internal", "models", names.Snake+".go")
	if !fileExists(path) {
		return
	}
	structs, err := parseGoStructs(path)
	if err != nil {
		// A model that does not parse is a model the developer is in the middle
		// of, and the write below will fail or be refused on its own terms.
		// Nothing useful to say about its fields.
		return
	}

	var existing []GoField
	for _, s := range structs {
		if s.Name == names.Pascal {
			existing = s.Fields
			break
		}
	}
	if len(existing) == 0 {
		return
	}

	// What the new definition will write. A belongs_to field named "group"
	// becomes GroupID and Group, so both spellings count as kept.
	keeping := make(map[string]bool, len(g.Definition.Fields)*2)
	for _, f := range g.Definition.Fields {
		pascal := toPascalCase(f.Name)
		keeping[pascal] = true
		keeping[pascal+"ID"] = true
		keeping[strings.TrimSuffix(pascal, "ID")] = true
	}

	var dropped []string
	for _, f := range existing {
		if frameworkFields[f.Name] || keeping[f.Name] {
			continue
		}
		// A relation struct and its key are one field to a person, so report
		// the key and not both.
		if strings.HasSuffix(f.Name, "ID") && keeping[strings.TrimSuffix(f.Name, "ID")] {
			continue
		}
		dropped = append(dropped, f.Name)
	}
	if len(dropped) == 0 {
		return
	}

	fmt.Printf("\n  ! This regenerate stops reading %d column%s: %s\n",
		len(dropped), plural(len(dropped)), strings.Join(dropped, ", "))
	fmt.Printf("    They are still in the database with their data, because migrations\n")
	fmt.Printf("    are additive. Nothing in the API, the admin or the types will read\n")
	fmt.Printf("    them after this. If that is not what you meant, the previous model\n")
	fmt.Printf("    is in git, and the field list to pass is the one in it.\n\n")
}
