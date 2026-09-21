package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Every formatted type, the way a person declares it.
const formattedFieldSpec = "name:string,email:email:unique,website:url,host:domain,phone:tel:UG:unique," +
	"country:country:KE,brand:color,discount:percent,score:rating:10,opens_at:time,meta:json"

func TestFormattedFieldsParse(t *testing.T) {
	def, err := ParseInlineFields("Company", formattedFieldSpec)
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	byName := map[string]Field{}
	for _, f := range def.Fields {
		byName[f.Name] = f
	}
	if f := byName["phone"]; f.Type != "tel" || f.DefaultCountry != "UG" || !f.Unique {
		t.Errorf("phone:tel:UG:unique parsed as %+v", f)
	}
	if f := byName["country"]; f.DefaultCountry != "KE" {
		t.Errorf("country:country:KE parsed as %+v", f)
	}
	if f := byName["score"]; f.Max != 10 || f.RatingMax() != 10 {
		t.Errorf("score:rating:10 parsed as %+v", f)
	}
	if f := byName["email"]; !f.Unique || f.DefaultCountry != "" {
		t.Errorf("email:email:unique parsed as %+v", f)
	}

	// The option is optional, and a modifier may take its place.
	if f, err := parseFieldInput("phone:tel:unique"); err != nil || !f.Unique || f.DefaultCountry != "" {
		t.Errorf("phone:tel:unique = %+v, %v", f, err)
	}
	if f, err := parseFieldInput("score:rating"); err != nil || f.RatingMax() != DefaultRatingMax {
		t.Errorf("score:rating = %+v, %v", f, err)
	}
	for _, bad := range []string{"phone:tel:XX", "score:rating:11", "score:rating:1", "meta:json:unique", "phone:tel:UG:bogus"} {
		if _, err := parseFieldInput(bad); err == nil {
			t.Errorf("%q was accepted", bad)
		}
	}
}

func TestFormattedFieldMappings(t *testing.T) {
	cases := []struct {
		field                        Field
		goType, gorm, format, ts, ui string
	}{
		{Field{Name: "email", Type: "email"}, "string", "size:254", "email", "string", "email"},
		{Field{Name: "site", Type: "url"}, "string", "size:2048", "url", "string", "url"},
		{Field{Name: "host", Type: "domain"}, "string", "size:253", "domain", "string", "domain"},
		{Field{Name: "phone", Type: "tel", DefaultCountry: "UG"}, "string", "size:32", "tel:UG", "string", "tel"},
		{Field{Name: "country", Type: "country"}, "string", "size:2", "country", "string", "country"},
		{Field{Name: "brand", Type: "color"}, "string", "size:7", "color", "string", "color"},
		{Field{Name: "discount", Type: "percent"}, "float64", "type:decimal(5,2)", "percent", "number", "percent"},
		{Field{Name: "score", Type: "rating"}, "int", "", "rating:5", "number", "rating"},
		{Field{Name: "opens_at", Type: "time"}, "string", "size:5", "time", "string", "time"},
		{Field{Name: "meta", Type: "json"}, "datatypes.JSON", "", "json", "unknown", "json"},
	}
	for _, c := range cases {
		f := c.field
		if got := f.GoType(); got != c.goType {
			t.Errorf("%s GoType = %q, want %q", f.Type, got, c.goType)
		}
		if got := f.GORMTag(); got != c.gorm {
			t.Errorf("%s GORMTag = %q, want %q", f.Type, got, c.gorm)
		}
		if got := f.FormatTag(); got != c.format {
			t.Errorf("%s FormatTag = %q, want %q", f.Type, got, c.format)
		}
		if got := f.TSType(); got != c.ts {
			t.Errorf("%s TSType = %q, want %q", f.Type, got, c.ts)
		}
		if got := f.FormFieldType(); got != c.ui {
			t.Errorf("%s FormFieldType = %q, want %q", f.Type, got, c.ui)
		}
		if f.ColumnFormat() == "text" {
			t.Errorf("%s renders as plain text in the table", f.Type)
		}
		if f.DocsTag() == "" {
			t.Errorf("%s has no OpenAPI docs tag", f.Type)
		}
		if module, _ := f.SharedSchema(); module == "" {
			t.Errorf("%s has no shared Zod schema", f.Type)
		}
		if strings.Contains(f.DocsTag(), `"`) {
			t.Errorf("%s docs tag would break the struct tag: %s", f.Type, f.DocsTag())
		}
	}
	if !(Field{Type: "json"}).NeedsDatatypesImport() {
		t.Error("json does not pull in gorm.io/datatypes")
	}
	for _, typ := range []string{"email", "url", "domain", "tel", "country", "color", "percent", "rating", "time", "json"} {
		if !isValidType(typ) {
			t.Errorf("%s is not a valid field type", typ)
		}
	}
}

func TestFormattedZod(t *testing.T) {
	cases := map[string]string{
		"email:required": "EmailSchema",
		"email":          `EmailSchema.or(z.literal("")).optional()`,
		"tel:required":   "PhoneSchema",
		"rating":         "ratingSchema(5).nullable().optional()",
		"percent":        "PercentSchema.nullable().optional()",
		"json":           "JsonValueSchema.optional()",
	}
	for spec, want := range cases {
		f, err := parseFieldInput("x:" + spec)
		if err != nil {
			t.Fatalf("%s: %v", spec, err)
		}
		if got := f.ZodType(); got != want {
			t.Errorf("%s ZodType = %q, want %q", spec, got, want)
		}
	}
}

func TestFormattedSeederValues(t *testing.T) {
	def, err := ParseInlineFields("Company", formattedFieldSpec)
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := &Generator{Module: "acme/apps/api", Definition: def}

	faker := g.seederContent(MakeNames("Company"), SeederOptions{Faker: true, Count: 1000})
	assertContains(t, "faker seeder", faker,
		"acme/apps/api/internal/fieldtypes",
		"acme/apps/api/internal/phone",
		"fieldtypes.SampleEmail(gofakeit.FirstName(), gofakeit.LastName(), gofakeit.DomainName(), i)",
		"phone.Sample(i)",
		"fieldtypes.SampleRating(i, 10)",
		"fieldtypes.SampleJSON(i)",
		"const n = 1000",
	)
	if strings.Contains(faker, "gofakeit.Phone()") {
		t.Error("the faker seeder fills a tel column with gofakeit.Phone(), which the API refuses")
	}

	static := g.seederContent(MakeNames("Company"), SeederOptions{})
	assertContains(t, "static seeder", static,
		`"+256772123456"`, `"ada@example.com"`, `"KE"`, `datatypes.JSON(`, `"gorm.io/datatypes"`)
	if strings.Contains(static, "/internal/phone") {
		t.Error("the static seeder imports internal/phone for a literal")
	}
}

func TestFormatTagRoundTrip(t *testing.T) {
	for _, tag := range []string{"email", "tel:UG", "rating:10", "json"} {
		f, ok := fieldFromFormatTag("x", tag)
		if !ok || f.FormatTag() != tag {
			t.Errorf("format tag %q came back as %q", tag, f.FormatTag())
		}
	}
	if _, ok := fieldFromFormatTag("x", "html"); ok {
		t.Error("an unknown format was taken for a field type")
	}
}

func TestSyncMapsFormatTags(t *testing.T) {
	got := fieldZod(GoField{GoType: "string", Format: "tel:UG"})
	if !strings.Contains(got, `\+[1-9]`) {
		t.Errorf("tel syncs to %s", got)
	}
	if got := fieldZod(GoField{GoType: "int", Format: "rating:10"}); !strings.Contains(got, "max(10)") {
		t.Errorf("rating:10 syncs to %s", got)
	}
	if got := fieldZod(GoField{GoType: "string"}); got != "z.string()" {
		t.Errorf("a plain string syncs to %s", got)
	}
}

// End to end: the model carries the format tags and the phone import, the
// schema imports the shared rules, the admin definition carries the options,
// and the project gets internal/fieldtypes and internal/phone.
func TestFormattedFieldsGenerate(t *testing.T) {
	const module = "acme/apps/api"
	root := setupMinimalProject(t, module)

	var fetched []string
	prev := goGetModule
	goGetModule = func(dir, spec string) error { fetched = append(fetched, spec); return nil }
	t.Cleanup(func() { goGetModule = prev })

	def, err := ParseInlineFields("Company", "name:string,phone:tel:UG,email:email,score:rating:10,brand:color")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	if err := newTestGenerator(root, module, def).Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	api := filepath.Join(root, "apps", "api")
	model := readTestFile(t, filepath.Join(api, "internal", "models", "company.go"))
	assertContains(t, "model", model,
		`format:"tel:UG"`, `format:"email"`, `format:"rating:10"`, `format:"color"`,
		`docs:"format:e164,example:+256772123456"`,
		`_ "acme/apps/api/internal/phone"`,
	)
	for _, pkg := range []string{"fieldtypes/fieldtypes.go", "fieldtypes/samples.go", "phone/phone.go", "phone/phone_test.go"} {
		if _, err := os.Stat(filepath.Join(api, "internal", filepath.FromSlash(pkg))); err != nil {
			t.Errorf("internal/%s was not written", pkg)
		}
	}
	if len(fetched) != 1 || !strings.HasPrefix(fetched[0], "github.com/nyaruka/phonenumbers@") {
		t.Errorf("go get = %v, want the phonenumbers module once", fetched)
	}

	schema := readTestFile(t, filepath.Join(root, "packages", "shared", "schemas", "company.ts"))
	assertContains(t, "schema", schema,
		`import { ColorSchema, EmailSchema, ratingSchema } from "./field-formats";`,
		`import { PhoneSchema } from "./phone";`,
		"ratingSchema(10)",
	)
	for _, f := range []string{"field-formats.ts", "phone.ts"} {
		if _, err := os.Stat(filepath.Join(root, "packages", "shared", "schemas", f)); err != nil {
			t.Errorf("packages/shared/schemas/%s was not written", f)
		}
	}
}
