package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// toolsFor generates one resource's tools file as a string.
func toolsFor(t *testing.T, fields []Field) string {
	t.Helper()
	g := &Generator{
		Definition: &ResourceDefinition{Name: "Invoice", Fields: fields},
		Module:     "example.com/app",
	}
	return g.mcpToolsSource()
}

// The file registers itself, which is what makes removal a delete.
//
// grit#224 proposed a "// grit:mcp-tools" marker to register at. The registry
// already made one unnecessary, and the difference matters: every injection
// anchor in this generator is a thing a removal has to unpick, and the ones
// that go wrong leave a project that does not compile. An init() needs no
// anchor and no reverse edit.
func TestTheToolsFileRegistersItself(t *testing.T) {
	src := toolsFor(t, []Field{{Name: "number", Type: "string", Required: true}})

	if !strings.Contains(src, "func init() {") || !strings.Contains(src, "Register(") {
		t.Error("the tools file does not register itself, so it needs an injection " +
			"somewhere and a removal has something to unpick")
	}
	// Top-level, after init closes: a func declaration inside Register's
	// argument list does not parse, which is how the first version failed.
	initAt := strings.Index(src, "func init() {")
	bodyAt := strings.Index(src, "func invoiceBody(")
	if bodyAt == -1 {
		t.Fatal("no shared body builder")
	}
	if bodyAt < initAt {
		t.Error("the body builder is above init")
	}
	if strings.Contains(src[initAt:bodyAt], "func invoiceBody") {
		t.Error("the body builder is inside Register's argument list, which does not parse")
	}
}

// Every tool declares the routes it calls, because the key says so.
//
// A connection's API key is restricted to the union of its tools' routes, so a
// route a tool omits is one the app refuses at runtime. That makes Routes a
// security boundary rather than documentation, and the case that gets it wrong
// is a relation: resolving a client by name is a GET on /api/v1/clients, which
// is not obviously part of "create an invoice".
func TestARelationsLookupRouteIsDeclared(t *testing.T) {
	src := toolsFor(t, []Field{
		{Name: "number", Type: "string", Required: true},
		{Name: "client", Type: "belongs_to", RelatedModel: "Client", Required: true},
	})

	// The create tool resolves a name, so it needs the related read.
	for _, want := range []string{
		`Routes:     []string{"POST /api/v1/invoices", "GET /api/v1/clients"}`,
		`"PATCH /api/v1/invoices/:id", "GET /api/v1/clients"`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("missing %s, so the connection's key refuses the name lookup "+
				"the tool makes", want)
		}
	}
	// And a resource with no relation does not ask for a route it never calls:
	// a wider key for no reason.
	plain := toolsFor(t, []Field{{Name: "number", Type: "string"}})
	if strings.Contains(plain, "GET /api/v1/clients") {
		t.Error("a resource with no relation declares a relation route")
	}
}

// A belongs_to writes to the _id column, whatever the argument is called.
//
// The model names the foreign key client_id and the create schema requires it
// under that name, so a body keyed on the bare field is refused with 422
// "client_id is required" for a request that named the client perfectly well.
// Found live, by making one.
func TestARelationWritesToTheIDColumn(t *testing.T) {
	src := toolsFor(t, []Field{
		{Name: "client", Type: "belongs_to", RelatedModel: "Client", Required: true},
	})

	if !strings.Contains(src, `body["client_id"] = related`) {
		t.Error(`the body is keyed on "client" rather than "client_id", so every ` +
			"create is refused by the schema")
	}
	// The argument stays bare: "client" is what a person says, and accepting a
	// name there is the point of the lookup.
	if !strings.Contains(src, `OptionalString(args, "client")`) {
		t.Error("the argument is named client_id, which asks a model for an id " +
			"when the whole point is that it can pass a name")
	}
}

// A select's options reach the schema, so a model picks from them.
func TestASelectsOptionsBecomeAnEnum(t *testing.T) {
	src := toolsFor(t, []Field{
		{Name: "status", Type: "select", Default: "draft",
			Options: []FieldOption{{Value: "draft"}, {Value: "sent"}}},
	})
	if !strings.Contains(src, `"enum": []string{"draft", "sent"}`) {
		t.Errorf("a select's options are not in the schema:\n%s", src)
	}
}

// What an agent must never be offered.
//
// Each of these is a field that must not reach a model or a request body, and
// the reasons differ, which is why they are listed rather than inferred from
// one rule.
func TestTheFieldsAnAgentIsNeverOffered(t *testing.T) {
	src := toolsFor(t, []Field{
		{Name: "number", Type: "string"},
		{Name: "secret", Type: "text", Encrypted: true},
		{Name: "recovery", Type: "string", Sensitive: true},
		{Name: "closed_at", Type: "datetime", ServerSet: true},
		{Name: "contract", Type: "file"},
		{Name: "gallery", Type: "files"},
		{Name: "slug", Type: "slug", SlugSource: "number"},
	})

	for _, never := range []string{"secret", "recovery", "closed_at", "contract", "gallery", "slug"} {
		if strings.Contains(src, `"`+never+`"`) {
			t.Errorf("%q is offered to an agent", never)
		}
	}
	if !strings.Contains(src, `"number"`) {
		t.Error("an ordinary field was excluded too, so the filter is too wide")
	}
}

// The condition in a generated file is never a constant.
//
// "else if creating && false" is dead code and "&& true" is noise. Generated
// code is read by people, and a condition that is a literal tells the reader
// the generator could not be bothered.
func TestNoGeneratedConditionIsAConstant(t *testing.T) {
	src := toolsFor(t, []Field{
		{Name: "number", Type: "string", Required: true},
		{Name: "total", Type: "money"},
		{Name: "client", Type: "belongs_to", RelatedModel: "Client"},
	})

	for _, dead := range []string{"creating && true", "creating && false"} {
		if strings.Contains(src, dead) {
			t.Errorf("the generated file contains %q", dead)
		}
	}
}

// Nothing is written when the plugin is not installed.
//
// Probed on the package existing rather than on a flag, because a project that
// installed the plugin and removed it has the flag either way, and the
// question that decides it is whether the import would compile.
func TestNoToolsWithoutThePlugin(t *testing.T) {
	root := t.TempDir()
	g := &Generator{
		Definition:   &ResourceDefinition{Name: "Invoice"},
		Root:         root,
		Architecture: "api",
	}
	if path := g.mcpToolsPath(); path != "" {
		t.Errorf("would write %s into a project with no mcp package", path)
	}

	// And it does once the package is there.
	pkg := filepath.Join(g.APIRoot(), "internal", "mcp")
	if err := os.MkdirAll(pkg, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkg, "tool.go"), []byte("package mcp\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if path := g.mcpToolsPath(); path == "" {
		t.Error("the plugin is installed and no tools would be written")
	}
}

// Removing the resource removes its tools.
//
// Nothing imports the file and the tools reach the app over HTTP, so leaving
// it behind still compiles. It would go on offering an agent five tools whose
// routes answer 404, and tools/list would advertise them: a removal that
// leaves a working build and a lying catalogue is worse than one that breaks,
// because nothing says so.
func TestRemovingAResourceRemovesItsTools(t *testing.T) {
	body, err := os.ReadFile("remove.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), `"mcp", "tools_"+names.PluralSnake+".go"`) {
		t.Error("grit remove resource leaves the tools file behind")
	}
}
