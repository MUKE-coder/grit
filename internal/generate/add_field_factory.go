package generate

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// injectFactoryField adds a column to the resource's test factory.
//
// `grit generate field` updated the model, the handler, the service, the
// schemas, the type and the admin screen, and left this one alone. A new
// NOT NULL column then broke every test that creates a fixture, with a
// constraint error pointing at the factory rather than at the command that
// caused it; a nullable one got the zero value in every test, which is worse,
// because the tests pass while describing a record nobody would save.
//
// The value comes from factoryValue, the same function the generator uses, so
// a field type added tomorrow produces the same fixture whichever command
// wrote it. Reimplementing it here is how the two would drift.
func (g *Generator) injectFactoryField(names Names, f Field) error {
	path := filepath.Join(g.APIRoot(), "internal", "factory", names.Snake+".go")
	if !fileExists(path) {
		// An older project, or a resource generated before factories existed.
		// Not a failure: the column is already everywhere that matters.
		return nil
	}

	// A relationship is not a fixture value. The generator leaves the foreign
	// key empty and CreateX fills it with a fresh parent, which is a decision
	// made in the template rather than per field: injecting a value here would
	// be a parent id that does not exist.
	if f.IsRelationship() || f.IsSlug() {
		return nil
	}

	l, err := readGoLines(path)
	if err != nil {
		return fmt.Errorf("reading the %s factory: %w", names.Pascal, err)
	}

	field := toPascalCase(f.Name)
	if l.contains("\t\t" + field + ":") {
		return nil // already there
	}

	anchor := l.find(0, "// grit:factory:auto-end")
	if anchor < 0 {
		// A factory from before the anchor, or one somebody has rewritten.
		// Reported rather than forced: a struct literal is not a thing to
		// guess the shape of.
		fmt.Printf("  • internal/factory/%s.go has no marker, so add this by hand: %s\n",
			names.Snake, field)
		return nil
	}

	value, needsFmt, needsTime := factoryValue(f, names)
	l.insert(anchor, fmt.Sprintf("\t\t%s: %s,", field, value))

	// n is only declared when a field interpolates it, and that was decided
	// from the fields the resource had at generation time. A column added
	// later can be the first one that needs it.
	if strings.Contains(value, "n)") || strings.Contains(value, "n,") || value == "n" {
		if !l.contains("n := next()") {
			if open := l.findPrefix(0, "record := &models."); open >= 0 {
				l.insert(open, "\tn := next()")
				l.insert(open+1, "")
			}
		}
	}

	if err := ensureFactoryImports(l, g.Module, needsFmt, needsTime, f); err != nil {
		return err
	}

	if err := l.write(path); err != nil {
		return fmt.Errorf("writing the %s factory: %w", names.Pascal, err)
	}
	manifest.Refresh(path)
	return nil
}

// ensureFactoryImports adds what the new value needs.
//
// The factory imports only what its fields use, which is right for a
// generated file and means an added field can be the first to need one.
func ensureFactoryImports(l *goLines, module string, needsFmt, needsTime bool, f Field) error {
	wanted := map[string]string{}
	if needsFmt {
		wanted["fmt"] = "\t\"fmt\""
	}
	if needsTime {
		wanted["time"] = "\t\"time\""
	}
	if f.Encrypted {
		wanted["crypto"] = fmt.Sprintf("\t%q", module+"/internal/crypto")
	}
	if f.NeedsJSONTimeImport() {
		wanted["jsontime"] = fmt.Sprintf("\t%q", module+"/internal/jsontime")
	}
	if FieldType(f.Type) == FieldMoney {
		wanted["money"] = fmt.Sprintf("\t%q", module+"/internal/money")
	}
	if len(wanted) == 0 {
		return nil
	}

	open := l.find(0, "import (")
	if open < 0 {
		return fmt.Errorf("the factory has no import block")
	}
	for name, line := range wanted {
		if l.contains(strings.TrimSpace(line)) {
			continue
		}
		_ = name
		// At the top of the block. gofmt, which l.write runs, sorts and
		// groups it afterwards.
		l.insert(open+1, line)
	}
	return nil
}
