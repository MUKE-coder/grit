package generate

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// Two functions write a model field line: templates.go for a new resource and
// add_field.go for `grit generate field`. This reads both.
//
// It is the shape of bug this project keeps finding. Each writer compiles,
// each writer's own test passes, and the pair disagrees: a column added later
// with :sensitive would have carried no tag and no MarshalJSON, so it would be
// filterable, exported and returned in every response, while the modifier said
// otherwise and nothing said anything.
func TestBothModelWritersHideASensitiveColumn(t *testing.T) {
	const module = "myapp/apps/api"

	// Writer one: a new resource.
	root := setupMinimalProject(t, module)
	def, err := ParseInlineFields("Wallet", "label:string,phrase:text:sensitive")
	if err != nil {
		t.Fatal(err)
	}
	if err := newTestGenerator(root, module, def).Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}
	model := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "wallet.go"))
	assertHides(t, "grit generate resource", model, "Wallet", "phrase")

	// Writer two: `grit generate field`, into a model that already exists and
	// already hides one column. A method naming only the new field would put
	// the first one back in every response.
	root2 := setupMinimalProject(t, module)
	def2, err := ParseInlineFields("Vault", "label:string,phrase:text:sensitive")
	if err != nil {
		t.Fatal(err)
	}
	g2 := newTestGenerator(root2, module, def2)
	if err := g2.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}
	added, err := ParseInlineFields("Vault", "answer:string:sensitive")
	if err != nil {
		t.Fatal(err)
	}
	names2 := MakeNames("Vault")
	if err := g2.injectModelField(names2, added.Fields[0]); err != nil {
		t.Fatalf("injectModelField: %v", err)
	}
	model2 := readTestFile(t, filepath.Join(root2, "apps", "api", "internal", "models", "vault.go"))
	assertHides(t, "grit generate field", model2, "Vault", "phrase", "answer")

	// And exactly one method, not one per added field.
	if n := strings.Count(model2, "func (m Vault) MarshalJSON()"); n != 1 {
		t.Errorf("got %d MarshalJSON methods, want 1; the model will not compile:\n%s", n, model2)
	}

	// The three frontend files the same command writes. The create and update
	// schemas keep the column, because a form has to be able to set it; the
	// row type and the admin table do not, because the API sends nothing.
	if err := g2.writeResourceDefinition(names2); err != nil {
		t.Fatalf("writeResourceDefinition: %v", err)
	}
	if err := g2.injectFrontendField(names2, added.Fields[0]); err != nil {
		t.Fatalf("injectFrontendField: %v", err)
	}
	ts := readTestFile(t, filepath.Join(root2, "packages", "shared", "types", "vault.ts"))
	if strings.Contains(ts, "answer") {
		t.Errorf("the row type says the API sends it:\n%s", ts)
	}
	zod := readTestFile(t, filepath.Join(root2, "packages", "shared", "schemas", "vault.ts"))
	entity, writes, found := strings.Cut(zod, "export const CreateVaultSchema")
	if !found {
		t.Fatalf("no create schema in:\n%s", zod)
	}
	if !strings.Contains(writes, "answer") {
		t.Errorf("the create and update schemas dropped it, so a form cannot set it:\n%s", writes)
	}
	if strings.Contains(entity, "answer") {
		t.Errorf("the entity schema expects a key no response carries:\n%s", entity)
	}
	_, adminPath := scaffold.ResourceDefPath(filepath.Join(g2.AdminCodeRoot(), "resources"), names2.PluralKebab)
	admin := readTestFile(t, adminPath)
	cols, form, found := strings.Cut(admin, "form:")
	if !found {
		t.Fatalf("no form block in the admin definition:\n%s", admin)
	}
	if strings.Contains(cols, "answer") {
		t.Errorf("the admin lists a column of blanks:\n%s", cols)
	}
	if !strings.Contains(form, "answer") {
		t.Errorf("the admin form has no input for it, so it can never be set:\n%s", form)
	}
}

func assertHides(t *testing.T, writer, model, pascal string, columns ...string) {
	t.Helper()
	if _, err := parser.ParseFile(token.NewFileSet(), "model.go", model, parser.AllErrors); err != nil {
		t.Fatalf("%s wrote a model that does not parse: %v\n%s", writer, err, model)
	}
	if !strings.Contains(model, `"encoding/json"`) {
		t.Errorf("%s: the model marshals without importing encoding/json:\n%s", writer, model)
	}
	if !strings.Contains(model, "func (m "+pascal+") MarshalJSON()") {
		t.Fatalf("%s: no MarshalJSON, so every response carries the column:\n%s", writer, model)
	}
	for _, col := range columns {
		if !strings.Contains(model, `sensitive:"true"`) {
			t.Errorf("%s: no sensitive tag, so the rules engine cannot know:\n%s", writer, model)
			break
		}
		shadow := `*struct{} ` + "`" + `json:"` + col + `,omitempty"` + "`"
		if !strings.Contains(model, shadow) {
			t.Errorf("%s: MarshalJSON does not drop %q:\n%s", writer, col, model)
		}
	}
}
