package generate

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// The public surface has to make a deliberate decision about every field type
// there is, and emit the Go type the model actually declares for the ones it
// publishes.
//
// Nothing checked either, and both went wrong at once. A Product priced in
// money.Money published no price, because money was missing from the
// publishable switch and the switch's default is "no"; a resource with a date
// on it produced a public handler that did not compile, because the view struct
// said time.Time where the model says *jsontime.Date. Six field types were
// affected, and the existing tests passed throughout: they asserted things
// about the types somebody had remembered to handle.
//
// So this table lists every field type, and the test below reads field.go to
// find any that the table does not mention. Adding a field type and forgetting
// the public surface now fails here rather than in a generated project.
var publicFieldTypeExpectations = map[FieldType]struct {
	published bool
	goType    string
	why       string
}{
	FieldString:   {true, "string", ""},
	FieldText:     {true, "string", ""},
	FieldRichtext: {true, "string", ""},
	FieldSlug:     {true, "string", ""},
	FieldSelect:   {true, "string", ""},
	FieldRadio:    {true, "string", "a select wearing a different admin widget"},
	FieldURL:      {true, "string", ""},
	FieldDomain:   {true, "string", ""},
	FieldCountry:  {true, "string", ""},
	FieldColor:    {true, "string", ""},
	FieldTime:     {true, "string", "a HH:MM string column, not an instant"},

	FieldInt:     {true, "int", ""},
	FieldRating:  {true, "int", ""},
	FieldUint:    {true, "uint", ""},
	FieldFloat:   {true, "float64", ""},
	FieldPercent: {true, "float64", ""},
	FieldBool:    {true, "bool", ""},
	FieldToggle:  {true, "bool", "a bool wearing a different admin widget"},

	FieldMoney:    {true, "money.Money", "a price is the most public thing a product has"},
	FieldDate:     {true, "*jsontime.Date", "the model's type, not time.Time"},
	FieldDatetime: {true, "*jsontime.DateTime", "the model's type, not time.Time"},

	FieldFile:        {true, "*files.FileRef", ""},
	FieldFiles:       {true, "files.FileRefs", ""},
	FieldStringArray: {true, "datatypes.JSONSlice[string]", ""},
	FieldCheck:       {true, "datatypes.JSONSlice[string]", "a string list wearing a different admin widget"},

	FieldEmail: {false, "", "personal data, published deliberately or not at all"},
	FieldTel:   {false, "", "personal data, published deliberately or not at all"},
	FieldJSON:  {false, "", "holds whatever was put in it"},

	FieldBelongsTo:  {false, "", "a relation publishes a whole record nobody vetted"},
	FieldManyToMany: {false, "", "a relation publishes a whole record nobody vetted"},
	FieldOneToOne:   {false, "", "a relation publishes a whole record nobody vetted"},
}

// Every field type declared in field.go appears in the table above.
func TestPublicSurfaceDecidesEveryFieldType(t *testing.T) {
	src, err := os.ReadFile("field.go")
	if err != nil {
		t.Fatalf("reading field.go: %v", err)
	}
	decl := regexp.MustCompile(`(?m)^\s*Field[A-Za-z]+\s+(?:FieldType\s+)?=\s*"([a-z_]+)"`)
	matches := decl.FindAllStringSubmatch(string(src), -1)
	if len(matches) < 25 {
		t.Fatalf("found only %d field types in field.go; the regex has stopped matching", len(matches))
	}
	for _, m := range matches {
		ft := FieldType(m[1])
		if _, ok := publicFieldTypeExpectations[ft]; !ok {
			t.Errorf("field type %q is not in publicFieldTypeExpectations: decide whether the "+
				"public surface publishes it, and if it does, which Go type the model declares for it", ft)
		}
	}
}

// Published or held back, as the table says, and with the declared Go type.
func TestPublicFieldsMatchExpectations(t *testing.T) {
	for ft, want := range publicFieldTypeExpectations {
		f := Field{Name: "thing", Type: string(ft)}
		if ft == FieldSelect || ft == FieldRadio || ft == FieldCheck {
			f.Options = []FieldOption{{Value: "a", Label: "A"}}
		}

		included, _ := PublicFields([]Field{f})
		got := len(included) == 1
		if got != want.published {
			verb := "held back"
			if want.published {
				verb = "published"
			}
			t.Errorf("%s: want %s, got the opposite (%s)", ft, verb, want.why)
			continue
		}
		if !want.published {
			continue
		}
		if gt := publicGoType(f); gt != want.goType {
			t.Errorf("%s: public view struct says %s, model declares %s", ft, gt, want.goType)
		}
	}
}

// A published type that needs an import has to bring it, and one that does not
// must not: an unused import is a build failure.
func TestPublicImportsCoverThePublishedTypes(t *testing.T) {
	const module = "shop/apps/api"
	needs := map[FieldType]string{
		FieldDate:        "internal/jsontime",
		FieldDatetime:    "internal/jsontime",
		FieldMoney:       "internal/money",
		FieldFile:        "internal/files",
		FieldFiles:       "internal/files",
		FieldStringArray: "gorm.io/datatypes",
		FieldCheck:       "gorm.io/datatypes",
	}
	for ft, want := range needs {
		f := Field{Name: "thing", Type: string(ft)}
		if ft == FieldCheck {
			f.Options = []FieldOption{{Value: "a", Label: "A"}}
		}
		std, third, local := publicImports(module, []Field{f})
		all := std + third + local
		if !strings.Contains(all, want) {
			t.Errorf("%s publishes as %s but publicImports emitted %q, which does not import %s",
				ft, publicGoType(f), all, want)
		}
	}

	// A resource of plain strings needs none of them.
	std, third, local := publicImports(module, []Field{{Name: "title", Type: string(FieldString)}})
	if got := std + third + local; strings.TrimSpace(got) != "" {
		t.Errorf("a string-only resource should need no extra imports, got %q", got)
	}
}

// The name rule that holds back stock counts must not hold back the derived
// boolean it tells you to publish instead.
func TestAvailableBoolIsPublishedButAvailableCountIsNot(t *testing.T) {
	cases := []struct {
		name, fieldType string
		want            bool
	}{
		{"available", string(FieldBool), true},
		{"available", string(FieldToggle), true},
		{"available", string(FieldInt), false},
		{"stock", string(FieldInt), false},
		{"stock", string(FieldBool), true},
		{"quantity", string(FieldInt), false},
		{"in_stock", string(FieldBool), true},
		{"cost_price", string(FieldMoney), false},
		{"price", string(FieldMoney), true},
	}
	for _, c := range cases {
		included, _ := PublicFields([]Field{{Name: c.name, Type: c.fieldType}})
		if got := len(included) == 1; got != c.want {
			t.Errorf("%s:%s published=%v, want %v", c.name, c.fieldType, got, c.want)
		}
	}
}
