package generate

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"
)

// `grit generate field` updates the test factory.
//
// It was the one generated file the command did not touch. A new NOT NULL
// column then broke every test that creates a fixture, with a constraint
// error pointing at the factory rather than at the command that caused it;
// a nullable one got the zero value everywhere, which is worse, because the
// tests pass while describing a record nobody would save.
func TestAddingAFieldReachesTheFactory(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)
	g := newTestGenerator(root, module, mustFields(t, "Widget", "name:string,qty:int"))
	names := g.Names()
	if err := g.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	added, err := ParseInlineFields("Widget", "colour:string")
	if err != nil {
		t.Fatal(err)
	}
	if err := g.injectFactoryField(names, added.Fields[0]); err != nil {
		t.Fatalf("injectFactoryField: %v", err)
	}

	path := filepath.Join(root, "apps", "api", "internal", "factory", "widget.go")
	src := readTestFile(t, path)

	if _, err := parser.ParseFile(token.NewFileSet(), "widget.go", src, parser.AllErrors); err != nil {
		t.Fatalf("the factory no longer parses: %v\n%s", err, src)
	}
	if !strings.Contains(src, "Colour:") {
		t.Errorf("the new column is not in the fixture:\n%s", src)
	}
	// Inside the struct literal, not after it.
	literal, _, found := strings.Cut(src, "// grit:factory:auto-end")
	if !found {
		t.Fatal("the anchor is gone, so the next field has nowhere to go")
	}
	if !strings.Contains(literal, "Colour:") {
		t.Error("the field landed outside the record literal")
	}
	// Twice is once.
	before := src
	if err := g.injectFactoryField(names, added.Fields[0]); err != nil {
		t.Fatal(err)
	}
	if readTestFile(t, path) != before {
		t.Error("injecting the same field twice changed the file again")
	}
}

// A column added to a resource whose fields never interpolated the counter is
// the first one that needs it declared.
//
// The generator only writes `n := next()` when a field uses it, decided from
// the fields the resource had at generation time. Without this the injected
// value referenced an undefined n and the package stopped compiling, which
// is a generated test helper breaking the build of every test in the project.
func TestTheFactoryCounterIsDeclaredWhenAnAddedFieldNeedsIt(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)
	// A bool needs no counter, so the generated factory declares none.
	g := newTestGenerator(root, module, mustFields(t, "Flag", "active:bool"))
	names := g.Names()
	if err := g.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	path := filepath.Join(root, "apps", "api", "internal", "factory", "flag.go")
	if strings.Contains(readTestFile(t, path), "n := next()") {
		t.Skip("this resource's factory already declares the counter, so the case " +
			"this test is about cannot arise here")
	}

	added, err := ParseInlineFields("Flag", "label:string")
	if err != nil {
		t.Fatal(err)
	}
	if err := g.injectFactoryField(names, added.Fields[0]); err != nil {
		t.Fatalf("injectFactoryField: %v", err)
	}

	src := readTestFile(t, path)
	if _, err := parser.ParseFile(token.NewFileSet(), "flag.go", src, parser.AllErrors); err != nil {
		t.Fatalf("the factory no longer parses: %v\n%s", err, src)
	}
	if !strings.Contains(src, "n := next()") {
		t.Errorf("the injected value uses n and nothing declares it, so the factory "+
			"package will not compile:\n%s", src)
	}
	if !strings.Contains(src, `"fmt"`) {
		t.Errorf("the injected value calls fmt.Sprintf without the import:\n%s", src)
	}
}

// A relationship has no fixture value to inject: the generator leaves the
// foreign key empty and CreateX fills it with a fresh parent. Writing one
// here would be an id no row has.
func TestARelationshipIsNotGivenAFixtureValue(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)
	g := newTestGenerator(root, module, mustFields(t, "Widget", "name:string"))
	names := g.Names()
	if err := g.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	path := filepath.Join(root, "apps", "api", "internal", "factory", "widget.go")
	before := readTestFile(t, path)

	added, err := ParseInlineFields("Widget", "owner:belongs_to[User]")
	if err != nil {
		t.Fatal(err)
	}
	if err := g.injectFactoryField(names, added.Fields[0]); err != nil {
		t.Fatalf("injectFactoryField: %v", err)
	}
	if readTestFile(t, path) != before {
		t.Error("a belongs_to was given a fixture value, which is a parent id that " +
			"no row has")
	}
}
