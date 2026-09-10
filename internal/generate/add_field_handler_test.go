package generate

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// `grit generate field` reaches the API, not just the model and the admin.
//
// It used to stop before the handler. Generated handlers bind into
// Create<X>Request and Update<X>Request and copy fields across one by one, so
// after adding notes to an Account the admin showed and sent the field, POST
// returned 201 with notes empty, PUT ignored it, PATCH refused it with 422, and
// the column sat empty in Postgres.

func accountWithHandler(t *testing.T) (*Generator, string) {
	t.Helper()
	const module = "ledger/apps/api"
	root := setupMinimalProject(t, module)
	g := newTestGenerator(root, module, mustFields(t, "Account", "code:string,name:string"))
	if err := g.writeGoHandler(g.Names()); err != nil {
		t.Fatalf("handler: %v", err)
	}
	return g, filepath.Join(root, "apps", "api", "internal", "handlers", "account.go")
}

func oneField(t *testing.T, spec string) Field {
	t.Helper()
	def := mustFields(t, "Account", spec)
	return def.Fields[len(def.Fields)-1]
}

func TestAddFieldReachesTheHandler(t *testing.T) {
	g, path := accountWithHandler(t)
	if err := g.injectHandlerField(g.Names(), oneField(t, "notes:text")); err != nil {
		t.Fatalf("inject: %v", err)
	}
	src := readTestFile(t, path)

	for what, re := range map[string]string{
		"create request": `(?s)type CreateAccountRequest struct \{[^}]*Notes\s+string\s+` + "`" + `json:"notes"` + "`",
		"update request": `(?s)type UpdateAccountRequest struct \{[^}]*Notes\s+string\s+` + "`" + `json:"notes"` + "`",
		"create literal": `Notes:\s+req\.Notes,`,
		"update map":     `updates\["notes"\] = req\.Notes`,
	} {
		if !regexp.MustCompile(re).MatchString(src) {
			t.Errorf("the field never reached the %s", what)
		}
	}
	// PATCH's allow-list and bulk edit's copy of it.
	if n := strings.Count(src, `"notes": true,`); n < 2 {
		t.Errorf(`"notes" is in %d allow-list(s), want the PATCH and bulk ones`, n)
	}
	assertParses(t, src)
}

// grit g field is re-runnable; the second run must not add the field twice.
func TestAddFieldHandlerInjectionIsIdempotent(t *testing.T) {
	g, path := accountWithHandler(t)
	f := oneField(t, "notes:text")
	for i := 0; i < 2; i++ {
		if err := g.injectHandlerField(g.Names(), f); err != nil {
			t.Fatalf("run %d: %v", i+1, err)
		}
	}
	if n := strings.Count(readTestFile(t, path), "`json:\"notes\"`"); n != 2 {
		t.Errorf("notes appears in %d request fields, want exactly 2 (create and update)", n)
	}
}

// A date field brings the import its type lives in, or the handler does not
// build.
func TestAddFieldBringsJSONTimeImport(t *testing.T) {
	g, path := accountWithHandler(t)
	if err := g.injectHandlerField(g.Names(), oneField(t, "opened_on:date")); err != nil {
		t.Fatalf("inject: %v", err)
	}
	src := readTestFile(t, path)
	if !strings.Contains(src, `"ledger/apps/api/internal/jsontime"`) {
		t.Error("the handler uses jsontime.Date and does not import it")
	}
	assertParses(t, src)
}

// A handler that has been reshaped until the create request is gone is refused
// by name, not skipped: skipping is the silent drop this exists to end.
func TestAddFieldRefusesWithoutACreateRequest(t *testing.T) {
	g, path := accountWithHandler(t)
	src := readTestFile(t, path)
	writeTestFile(t, path, strings.Replace(src, "type CreateAccountRequest struct {", "type CreateThing struct {", 1))

	err := g.injectHandlerField(g.Names(), oneField(t, "notes:text"))
	if err == nil || !strings.Contains(err.Error(), "CreateAccountRequest") {
		t.Fatalf("want a refusal naming CreateAccountRequest, got %v", err)
	}
}

// scalarHandlerParts and importAssign are copies of the generator's rules.
// These fail the day either copy stops matching what a fresh generate emits.
var addFieldTypes = "s:string,tx:text,rt:richtext,n:int,u:uint,fl:float,b:bool,tg:toggle,sel:select:a=A|b=B,d:date,dt:datetime"

func TestAddFieldMatchesTheGenerator(t *testing.T) {
	const module = "ledger/apps/api"
	root := setupMinimalProject(t, module)
	def := mustFields(t, "Account", addFieldTypes)
	g := newTestGenerator(root, module, def)
	if err := g.writeGoHandler(g.Names()); err != nil {
		t.Fatalf("handler: %v", err)
	}
	generated := squash(readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "account.go")))

	for _, f := range def.Fields {
		p := scalarHandlerParts(f)
		for what, part := range map[string]string{
			"create field": p.createField, "create assign": p.createAssign,
			"patch key": p.patchKey, "update field": p.updateField, "update set": p.updateSet,
		} {
			if !strings.Contains(generated, squash(part)) {
				t.Errorf("%s: the %s no longer matches the generator:\n  %s", f.Name, what, part)
			}
		}
	}
}

func TestAddFieldImportMatchesTheGenerator(t *testing.T) {
	const module = "ledger/apps/api"
	root := setupMinimalProject(t, module)
	def := mustFields(t, "Account", addFieldTypes)
	g := newTestGenerator(root, module, def)
	if err := g.writeGoImportHandler(g.Names()); err != nil {
		t.Fatalf("import handler: %v", err)
	}
	generated := squash(readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "account_import.go")))

	for _, f := range def.Fields {
		code, _, ok := importAssign(f)
		if !ok {
			continue
		}
		if !strings.Contains(generated, squash(code)) {
			t.Errorf("%s: the CSV import no longer matches the generator:\n%s", f.Name, code)
		}
	}
}

// squash collapses whitespace, so gofmt's column alignment does not count as a
// difference.
func squash(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func assertParses(t *testing.T, src string) {
	t.Helper()
	if _, err := parser.ParseFile(token.NewFileSet(), "handler.go", src, parser.AllErrors); err != nil {
		t.Fatalf("the handler no longer parses: %v", err)
	}
}

// A date field brings its import to the model too. Without it the project did
// not build: grit g field added "OpenedOn *jsontime.Date" to a model that had
// never imported jsontime.
func TestAddFieldModelBringsJSONTimeImport(t *testing.T) {
	const module = "ledger/apps/api"
	root := setupMinimalProject(t, module)
	g := newTestGenerator(root, module, mustFields(t, "Account", "code:string"))
	names := g.Names()
	if err := g.writeGoModel(names); err != nil {
		t.Fatalf("model: %v", err)
	}
	path := filepath.Join(root, "apps", "api", "internal", "models", "account.go")
	if strings.Contains(readTestFile(t, path), "internal/jsontime") {
		t.Fatal("precondition: the model already imports jsontime, so this proves nothing")
	}

	f := oneField(t, "opened_on:date")
	if err := g.injectModelField(names, f); err != nil {
		t.Fatalf("field: %v", err)
	}
	if err := g.ensureModelImport(names, f); err != nil {
		t.Fatalf("import: %v", err)
	}
	src := readTestFile(t, path)
	if !strings.Contains(src, "\"ledger/apps/api/internal/jsontime\"") {
		t.Error("the model uses jsontime.Date and does not import it")
	}
	assertParses(t, src)
}
