package generate

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

// writeFactory emits internal/factory/<snake>.go: a builder that makes a valid
// instance of this resource for a test.
//
// # Why a generated factory rather than a struct literal per test
//
// A test that needs a widget writes one, and a test that needs a widget with a
// different price writes another, and both name every required field. Add a
// required field to the model and every one of them stops compiling, or worse,
// keeps compiling and starts failing validation at the database. The factory is
// the one place that knows what a valid row looks like, so adding a field is one
// edit here rather than one per test.
//
// # What it deliberately leaves empty
//
// A belongs_to is left unset. The factory cannot know which parent you mean, and
// inventing one by inserting a parent row would make every test that touches
// this resource also create rows in another table, which is how a test suite
// becomes impossible to reason about. The doc comment on the generated function
// says so and shows the override.
func (g *Generator) writeFactory(names Names) error {
	path := filepath.Join(g.APIRoot(), "internal", "factory", names.Snake+".go")
	return writeFileWithDirs(path, g.factorySource(names))
}

// factorySource renders the factory file.
func (g *Generator) factorySource(names Names) string {
	assignments, needsFmt, needsTime, needsMoney, needsCrypto, needsJSONTime, relations := g.factoryFields(names)

	imports := []string{}
	if needsFmt {
		imports = append(imports, `"fmt"`)
	}
	imports = append(imports, `"testing"`)
	if needsTime {
		imports = append(imports, `"time"`)
	}

	var b strings.Builder
	b.WriteString("package factory\n\n")
	b.WriteString("import (\n")
	for _, imp := range imports {
		b.WriteString("\t" + imp + "\n")
	}
	b.WriteString("\n\t\"gorm.io/gorm\"\n\n")
	if needsCrypto {
		fmt.Fprintf(&b, "\t%q\n", g.Module+"/internal/crypto")
	}
	if needsJSONTime {
		fmt.Fprintf(&b, "\t%q\n", g.Module+"/internal/jsontime")
	}
	if needsMoney {
		fmt.Fprintf(&b, "\t%q\n", g.Module+"/internal/money")
	}
	fmt.Fprintf(&b, "\t%q\n)\n\n", g.Module+"/internal/models")

	// The doc comment carries the relation note, because that is the thing
	// somebody hits ten minutes in and the answer is one line of override.
	fmt.Fprintf(&b, "// New%s builds a valid %s without saving it.\n//\n", names.Pascal, names.Pascal)
	b.WriteString("// Every call produces different values for the unique fields, so two rows made\n")
	b.WriteString("// in the same test do not collide.\n")
	if len(relations) > 0 {
		b.WriteString("//\n")
		fmt.Fprintf(&b, "// %s is left empty: the factory cannot know which parent you mean, and\n",
			humanList(relations))
		b.WriteString("// inserting one would make every test that touches this resource write rows in\n")
		b.WriteString("// another table too. Pass it in:\n//\n")
		fmt.Fprintf(&b, "//\t%s := factory.Create%s(t, db, func(m *models.%s) {\n",
			names.Camel, names.Pascal, names.Pascal)
		fmt.Fprintf(&b, "//\t\tm.%s = parent.ID\n", relations[0])
		b.WriteString("//\t})\n")
	}
	fmt.Fprintf(&b, "func New%s(overrides ...func(*models.%s)) *models.%s {\n",
		names.Pascal, names.Pascal, names.Pascal)
	// The counter is declared only when a field interpolates it. A resource
	// whose every field is a relation or a file has no use for it, and `_ = n`
	// in generated code reads as something the generator was not sure about.
	if strings.Contains(assignments, "n)") || strings.Contains(assignments, "n,") {
		b.WriteString("\tn := next()\n\n")
	}
	fmt.Fprintf(&b, "\trecord := &models.%s{\n", names.Pascal)
	b.WriteString(assignments)
	b.WriteString("\t}\n\n")
	b.WriteString("\tfor _, apply := range overrides {\n\t\tapply(record)\n\t}\n")
	b.WriteString("\treturn record\n}\n\n")

	fmt.Fprintf(&b, "// Create%s builds one and writes it, failing the test if it will not save.\n",
		names.Pascal)
	b.WriteString("//\n")
	b.WriteString("// Failing here rather than returning an error is deliberate: a fixture that\n")
	b.WriteString("// cannot be created is not a test failure to be asserted on, it is a broken\n")
	b.WriteString("// test, and the stack should point at the line that asked for it.\n")
	fmt.Fprintf(&b, "func Create%s(tb testing.TB, db *gorm.DB, overrides ...func(*models.%s)) *models.%s {\n",
		names.Pascal, names.Pascal, names.Pascal)
	b.WriteString("\ttb.Helper()\n")
	fmt.Fprintf(&b, "\trecord := New%s(overrides...)\n", names.Pascal)
	b.WriteString("\tif err := db.Create(record).Error; err != nil {\n")
	fmt.Fprintf(&b, "\t\ttb.Fatalf(\"creating a %s fixture: %%v\", err)\n", names.Lower)
	b.WriteString("\t}\n\treturn record\n}\n\n")

	fmt.Fprintf(&b, "// Create%s writes n of them, for a test about paging or sorting.\n", names.PluralPascal)
	fmt.Fprintf(&b, "func Create%s(tb testing.TB, db *gorm.DB, n int, overrides ...func(*models.%s)) []*models.%s {\n",
		names.PluralPascal, names.Pascal, names.Pascal)
	b.WriteString("\ttb.Helper()\n")
	fmt.Fprintf(&b, "\tout := make([]*models.%s, 0, n)\n", names.Pascal)
	b.WriteString("\tfor i := 0; i < n; i++ {\n")
	fmt.Fprintf(&b, "\t\tout = append(out, Create%s(tb, db, overrides...))\n", names.Pascal)
	b.WriteString("\t}\n\treturn out\n}\n")

	return b.String()
}

// factoryFields renders the struct literal body, and reports what it needed.
func (g *Generator) factoryFields(names Names) (assignments string, needsFmt, needsTime, needsMoney, needsCrypto, needsJSONTime bool, relations []string) {
	var b strings.Builder

	for _, f := range g.Definition.Fields {
		// A slug is derived in a hook, an auto field comes from a sequence, and
		// a many-to-many is a join table written after the row exists. Setting
		// any of them here would be overwritten or rejected.
		if f.IsSlug() || f.Auto || f.IsManyToMany() {
			continue
		}

		goName := toPascalCase(f.Name)

		if f.IsBelongsTo() {
			relations = append(relations, toPascalCase(f.FKColumnName()))
			continue
		}

		value, usesFmt, usesTime := factoryValue(f, names)
		if value == "" {
			continue
		}
		needsFmt = needsFmt || usesFmt
		needsTime = needsTime || usesTime
		needsMoney = needsMoney || FieldType(f.Type) == FieldMoney
		switch FieldType(f.Type) {
		case FieldDate, FieldDatetime:
			needsJSONTime = true
		}

		// An :encrypted column is stored as crypto.EncryptedString:
		// AES-256-GCM at rest, a plain string in code. The field is still
		// declared with that type, so the literal has to say so.
		if f.Encrypted {
			needsCrypto = true
			value = "crypto.EncryptedString(" + value + ")"
		}
		fmt.Fprintf(&b, "\t\t%s: %s,\n", goName, value)
	}

	return b.String(), needsFmt, needsTime, needsMoney, needsCrypto, needsJSONTime, relations
}

// factoryValue is a plausible value for one field.
//
// Plausible, not random. A factory that produced random values would make a
// failing test's output different on every run, and "expected 3, got 7" is only
// useful when running it again gives the same 7. The counter is what keeps
// unique columns distinct.
func factoryValue(f Field, names Names) (value string, needsFmt, needsTime bool) {
	switch FieldType(f.Type) {
	case FieldString, FieldText, FieldRichtext:
		return fmt.Sprintf("fmt.Sprintf(%q, n)", titleFor(f, names)+" %d"), true, false
	case FieldEmail:
		return `fmt.Sprintf("fixture-%d@example.com", n)`, true, false
	case FieldURL:
		return `fmt.Sprintf("https://example.com/%d", n)`, true, false
	case FieldDomain:
		return `fmt.Sprintf("fixture-%d.example.com", n)`, true, false
	case FieldTel:
		return `fmt.Sprintf("+1415555%04d", n%10000)`, true, false
	case FieldColor:
		return `"#6c5ce7"`, false, false
	case FieldCountry:
		return `"US"`, false, false
	case FieldInt:
		return "n", false, false
	case FieldUint:
		return "uint(n)", false, false
	case FieldMoney:
		// money.Money is two columns, an amount in minor units and a currency,
		// so this is 19.99 and not 1999 dollars.
		return `money.New(1999+int64(n), "USD")`, false, false
	case FieldFloat:
		return "float64(n) + 0.5", false, false
	case FieldPercent:
		return "50", false, false
	case FieldRating:
		return "4", false, false
	case FieldBool, FieldToggle:
		return "true", false, false
	case FieldDatetime:
		// A fixed point rather than time.Now(): a test asserting on a date
		// should not change its answer at midnight, and one asserting on
		// ordering needs the values to be predictable.
		//
		// A *jsontime.DateTime and not a time.Time, because the admin sends
		// "2001-08-06T14:30", which is not RFC3339 and time.Time refuses.
		return "&jsontime.DateTime{Time: time.Date(2026, 1, 15, 9, 30, 0, 0, time.UTC)}", false, true
	case FieldDate:
		return "&jsontime.Date{Time: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)}", false, true
	case FieldTime:
		return `"09:30"`, false, false
	case FieldSelect, FieldRadio:
		if len(f.Options) > 0 {
			return strconv.Quote(f.Options[0].Value), false, false
		}
		return "", false, false
	case FieldJSON, FieldCheck, FieldStringArray, FieldFile, FieldFiles, FieldOneToOne:
		// Left to an override. A JSON blob, a file reference or a set of
		// checkboxes has no value that is right more often than it is wrong,
		// and a wrong default is worse than an empty one because it looks
		// deliberate.
		return "", false, false
	default:
		return "", false, false
	}
}

// titleFor is the human-ish prefix a string field's fixture value gets.
func titleFor(f Field, names Names) string {
	lower := strings.ToLower(f.Name)
	switch {
	case strings.Contains(lower, "name"), strings.Contains(lower, "title"):
		return names.Pascal
	case strings.Contains(lower, "description"), strings.Contains(lower, "body"),
		strings.Contains(lower, "content"), strings.Contains(lower, "notes"):
		return "A description for fixture"
	default:
		return toPascalCase(f.Name)
	}
}

// humanList renders names as "A", "A and B", or "A, B and C".
func humanList(values []string) string {
	switch len(values) {
	case 0:
		return ""
	case 1:
		return values[0]
	case 2:
		return values[0] + " and " + values[1]
	default:
		return strings.Join(values[:len(values)-1], ", ") + " and " + values[len(values)-1]
	}
}

// factoryCounterSource is internal/factory/factory.go: the shared counter.
//
// Written once per project rather than once per resource, because every factory
// draws from the same sequence and two resources sharing it is what keeps a test
// that creates one of each from seeing the same number twice.
func factoryCounterSource() string {
	return `package factory

// Package factory builds valid records for tests.
//
// One builder per resource, generated beside the model, so the knowledge of
// what a valid row looks like lives in one place rather than in every test that
// needs one. Adding a required field is then one edit here instead of one per
// test file.

import "sync/atomic"

// counter makes each fixture distinct.
//
// Shared across every resource in the package on purpose: two resources with
// their own counters would both start at 1, and a test creating one of each
// would get two rows claiming to be the first, which is exactly the collision
// a unique column then rejects.
var counter atomic.Int64

// next is the number the builders interpolate into unique fields.
func next() int {
	return int(counter.Add(1))
}

// Reset puts the counter back, for a test that asserts on an exact generated
// value. Rarely what you want: asserting on "Widget 1" rather than on the value
// the factory returned is a test that breaks when somebody adds a fixture above
// it.
func Reset() {
	counter.Store(0)
}
`
}

// ensureFactoryCounter writes internal/factory/factory.go when it is not
// there yet.
//
// Created rather than overwritten, like every other file this generator puts
// in the user tree that it does not own line by line: somebody may have added
// a helper to the package, and taking that away to deliver a counter that has
// not changed since it was written is the trade nobody wants.
func (g *Generator) ensureFactoryCounter() error {
	path := filepath.Join(g.APIRoot(), "internal", "factory", "factory.go")
	if fileExists(path) {
		return nil
	}
	if err := writeFileWithDirs(path, factoryCounterSource()); err != nil {
		return fmt.Errorf("writing the factory counter: %w", err)
	}
	return nil
}
