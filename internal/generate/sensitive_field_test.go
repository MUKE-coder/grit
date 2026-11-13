package generate

import (
	"encoding/json"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// The :sensitive modifier makes a column write-only.
//
// It exists because the rules engine shipped with a list of column names it
// would never filter on, and a list of names is a guess. `?filter=
// recovery_phrase ~ "aband"` reads a column one character at a time, and
// whether that works depended on whether somebody had thought of the word
// "recovery_phrase" when writing the denylist.

func TestSensitiveModifierParses(t *testing.T) {
	def, err := ParseInlineFields("Wallet", "label:string,recovery_phrase:text:sensitive")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if def.Fields[0].Sensitive {
		t.Error("an ordinary field came back sensitive")
	}
	if !def.Fields[1].Sensitive {
		t.Fatal("recovery_phrase:text:sensitive did not mark the field sensitive")
	}
	for spec, why := range map[string]string{
		"balance:int:sensitive":         "a number is not something this can hide meaningfully",
		"avatar:file:sensitive":         "a file is served by URL, so hiding the column hides nothing",
		"owner_id:belongs_to:sensitive": "an association the admin still renders a picker for",
		"phrase:text:sensitve":          "a typo must not be ignored",
	} {
		if _, err := ParseInlineFields("Wallet", spec); err == nil {
			t.Errorf("%s was accepted (%s)", spec, why)
		}
	}
}

// A sensitive column must never be searchable: `?search=abandon` returning
// the row is the same oracle a filter would be, with a friendlier spelling.
func TestSensitiveFieldIsNotSearchable(t *testing.T) {
	plain := Field{Name: "notes", Type: "text"}
	if !plain.IsSearchable() {
		t.Fatal("a plain text column stopped being searchable")
	}
	secret := Field{Name: "recovery_phrase", Type: "text", Sensitive: true}
	if secret.IsSearchable() {
		t.Error("a sensitive column is searchable")
	}
}

// The promise is kept by the model, not by whoever writes the next handler.
//
// This is the test that matters: it marshals a struct of exactly the shape
// the generator emits and checks the key is gone. A test that only grepped
// the generated source for "MarshalJSON" would pass on a method that returned
// the field anyway.
func TestTheGeneratedMarshalJSONDropsTheKey(t *testing.T) {
	type wallet struct {
		ID             string `json:"id"`
		Label          string `json:"label"`
		RecoveryPhrase string `json:"recovery_phrase"`
	}
	// The emitted shape, by hand.
	marshal := func(m wallet) string {
		type alias wallet
		out, err := json.Marshal(&struct {
			*alias
			RecoveryPhrase *struct{} `json:"recovery_phrase,omitempty"`
		}{alias: (*alias)(&m)})
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		return string(out)
	}

	got := marshal(wallet{ID: "w1", Label: "Savings", RecoveryPhrase: "abandon ability able"})
	if strings.Contains(got, "abandon") {
		t.Fatalf("the value is in the response: %s", got)
	}
	if strings.Contains(got, "recovery_phrase") {
		t.Fatalf("the key is in the response, so a client sees the column exists and is empty: %s", got)
	}
	if !strings.Contains(got, `"label":"Savings"`) {
		t.Fatalf("the rest of the record did not survive: %s", got)
	}
}

func TestSensitiveFieldIsHiddenEverywhereItCouldLeak(t *testing.T) {
	const module = "vault/apps/api"
	root := setupMinimalProject(t, module)
	g := newTestGenerator(root, module,
		mustFields(t, "Wallet", "label:string,recovery_phrase:text:sensitive"))
	names := g.Names()

	if err := g.writeGoModel(names); err != nil {
		t.Fatalf("model: %v", err)
	}
	if err := g.writeGoHandler(names); err != nil {
		t.Fatalf("handler: %v", err)
	}
	if err := g.writeTSTypes(names); err != nil {
		t.Fatalf("types: %v", err)
	}
	if err := g.writeZodSchema(names); err != nil {
		t.Fatalf("schemas: %v", err)
	}
	if err := g.writeResourceDefinition(names); err != nil {
		t.Fatalf("admin: %v", err)
	}

	api := filepath.Join(root, "apps", "api", "internal")
	model := readTestFile(t, filepath.Join(api, "models", "wallet.go"))
	handler := readTestFile(t, filepath.Join(api, "handlers", "wallet.go"))

	if _, err := parser.ParseFile(token.NewFileSet(), "wallet.go", model, parser.AllErrors); err != nil {
		t.Fatalf("the model does not parse: %v", err)
	}
	if !strings.Contains(model, `sensitive:"true"`) {
		t.Error("the model does not declare the column sensitive, so the rules engine cannot know")
	}
	if !strings.Contains(model, "func (m Wallet) MarshalJSON()") {
		t.Error("the model has no MarshalJSON, so every response carries the column")
	}
	if !strings.Contains(model, `"encoding/json"`) {
		t.Error("the model marshals without importing encoding/json")
	}
	if !strings.Contains(model, `RecoveryPhrase *struct{} `+"`"+`json:"recovery_phrase,omitempty"`+"`") {
		t.Error("MarshalJSON does not shadow the column it is supposed to drop")
	}
	// A model with nothing to hide pays nothing.
	plain := newTestGenerator(root, module, mustFields(t, "Note", "body:text"))
	if err := plain.writeGoModel(plain.Names()); err != nil {
		t.Fatalf("plain model: %v", err)
	}
	note := readTestFile(t, filepath.Join(api, "models", "note.go"))
	if strings.Contains(note, "MarshalJSON") || strings.Contains(note, `"encoding/json"`) {
		t.Error("a model with no sensitive column carries a MarshalJSON anyway")
	}

	// The column is still writable: a value nobody can set is a comment.
	if !strings.Contains(handler, "RecoveryPhrase") {
		t.Error("the create request dropped the column, so it can never be set")
	}

	// Not in the export. export reads Go fields by reflection, so a column
	// the API refuses to send would otherwise go out in a spreadsheet.
	if cols := g.crud(names).exportCols; strings.Contains(cols, "RecoveryPhrase") {
		t.Errorf("the export carries the column:\n%s", cols)
	}

	// Not a type the frontend thinks it receives.
	tsPath := filepath.Join(root, "packages", "shared", "types", "wallet.ts")
	ts := readTestFile(t, tsPath)
	if strings.Contains(ts, "recovery_phrase") {
		t.Errorf("the row type says the API sends it:\n%s", ts)
	}

	// In the create schema, out of the entity schema: one describes what goes
	// in, the other what comes back.
	zod := readTestFile(t, filepath.Join(root, "packages", "shared", "schemas", "wallet.ts"))
	// The entity schema is emitted first, then the create and update pair.
	entity, writes, found := strings.Cut(zod, "export const CreateWalletSchema")
	if !found {
		t.Fatalf("no create schema in:\n%s", zod)
	}
	if !strings.Contains(writes, "recovery_phrase") {
		t.Errorf("the create and update schemas dropped the column, so a form cannot set it:\n%s", writes)
	}
	if strings.Contains(entity, "recovery_phrase") {
		t.Errorf("the entity schema expects a key no response carries:\n%s", entity)
	}

	// Not an admin column, and not a filter.
	adminDir, adminPath := scaffold.ResourceDefPath(filepath.Join(g.AdminCodeRoot(), "resources"), names.PluralKebab)
	_ = adminDir
	admin := readTestFile(t, adminPath)
	cols, _, found := strings.Cut(admin, "form:")
	if !found {
		cols = admin
	}
	if strings.Contains(cols, "recovery_phrase") {
		t.Errorf("the admin lists a column of blanks:\n%s", cols)
	}
}
