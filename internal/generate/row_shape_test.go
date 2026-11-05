package generate

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// grit#145: the row TYPE and the row SCHEMA disagreed about the same row.
//
//	packages/shared/types/product.ts    price: number;
//	packages/shared/schemas/product.ts  price: z.number().optional()
//
// Both export `Product`, so a file importing it from one module and a file
// importing it from the other got different answers about whether price can be
// undefined. Neither errored, because each file is internally consistent,
// which is exactly why it went unnoticed: there is no symptom until somebody
// relies on the difference.
//
// They disagreed because each derived presence by its own rule. buildTSType
// uses one fact, the json tag's omitempty, with nullability from the Go
// pointer. The entity schema reused the CREATE shape, which carries
// .optional() for anything a create form may leave out: a create-time rule
// applied to a response.
//
// This test reads BOTH and compares them. A test for either one alone passes
// while the pair is broken, which is the shape of bug this repository keeps
// finding.
func TestTheRowTypeAndTheRowSchemaAgree(t *testing.T) {
	// A field of each interesting kind: required, optional, a pointer, a
	// framework column, and a nullable timestamp.
	const model = `package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          string         ` + "`gorm:\"primarykey\" json:\"id\"`" + `
	Name        string         ` + "`json:\"name\" binding:\"required\"`" + `
	Price       float64        ` + "`json:\"price\"`" + `
	Active      bool           ` + "`json:\"active\"`" + `
	Note        *string        ` + "`json:\"note,omitempty\"`" + `
	Version     int            ` + "`json:\"version\"`" + `
	CreatedAt   time.Time      ` + "`json:\"created_at\"`" + `
	UpdatedAt   time.Time      ` + "`json:\"updated_at\"`" + `
	ArchivedAt  *time.Time     ` + "`json:\"archived_at,omitempty\"`" + `
	DeletedAt   gorm.DeletedAt ` + "`gorm:\"index\" json:\"-\"`" + `
}
`
	dir := t.TempDir()
	path := filepath.Join(dir, "product.go")
	if err := os.WriteFile(path, []byte(model), 0o644); err != nil {
		t.Fatal(err)
	}
	structs, err := parseGoStructs(path)
	if err != nil {
		t.Fatalf("parsing the model: %v", err)
	}
	if len(structs) == 0 {
		t.Fatal("no structs parsed, so this test is reading nothing")
	}

	ts := buildTSType(structs[0])
	zod := buildZodSchema(structs[0])

	typeShape := tsPresence(ts)
	schemaShape := zodPresence(zod)

	if len(typeShape) == 0 {
		t.Fatal("no fields read from the interface; this test is not reading what it thinks")
	}

	for field, want := range typeShape {
		got, ok := schemaShape[field]
		if !ok {
			t.Errorf("%s is in the interface and not in the row schema", field)
			continue
		}
		if got != want {
			t.Errorf("%s: the interface says %s, the row schema says %s.\n"+
				"Both are exported as the row's type, so which answer you get depends on "+
				"which module you import. Presence comes from the json tag's omitempty "+
				"and nullability from the Go pointer, in both.", field, want, got)
		}
	}
	for field := range schemaShape {
		if _, ok := typeShape[field]; !ok {
			t.Errorf("%s is in the row schema and not in the interface", field)
		}
	}
}

var rowTSFieldRe = regexp.MustCompile(`(?m)^\s{2}(\w+)(\??):\s*(.+);$`)

// tsPresence reads each interface field as "required", "optional",
// "nullable" or "optional nullable".
func tsPresence(src string) map[string]string {
	out := map[string]string{}
	for _, m := range rowTSFieldRe.FindAllStringSubmatch(src, -1) {
		name, question, tsType := m[1], m[2], m[3]
		shape := "required"
		if strings.Contains(tsType, "| null") {
			shape = "nullable"
		}
		if question == "?" {
			shape = "optional " + shape
		}
		out[name] = shape
	}
	return out
}

var rowZodFieldRe = regexp.MustCompile(`(?m)^\s{2}(\w+):\s*(.+),$`)

// zodPresence reads the ROW schema's fields the same way. Only the first
// z.object, which is the entity schema; the create and update pair below it
// describe what goes in rather than what comes back.
func zodPresence(src string) map[string]string {
	start := strings.Index(src, "Schema = z.object({")
	if start < 0 {
		return nil
	}
	body := src[start:]
	if end := strings.Index(body, "});"); end > 0 {
		body = body[:end]
	}

	out := map[string]string{}
	for _, m := range rowZodFieldRe.FindAllStringSubmatch(body, -1) {
		name, zod := m[1], m[2]
		shape := "required"
		if strings.Contains(zod, ".nullable()") {
			shape = "nullable"
		}
		if strings.Contains(zod, ".optional()") {
			shape = "optional " + shape
		}
		out[name] = shape
	}
	return out
}
