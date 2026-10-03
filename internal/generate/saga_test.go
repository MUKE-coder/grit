package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// sagaProject builds the part of a project `grit generate workflow` needs: the
// engine it checks for, and the registry it injects into.
func sagaProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	write := func(rel, content string) {
		path := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write(filepath.Join("internal", "saga", "saga.go"), "package saga\n\nfunc Register(d Definition) {}\n")
	write(filepath.Join("internal", "sagas", "sagas.go"),
		"package sagas\n\nfunc Register() {\n\t// grit:sagas\n}\n")
	return root
}

func readSaga(t *testing.T, root, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		t.Fatalf("reading %s: %v", rel, err)
	}
	return string(data)
}

func TestGenerateWorkflowWritesASagaAndRegistersIt(t *testing.T) {
	root := sagaProject(t)
	err := generateWorkflowAt(root, "demo/apps/api", WorkflowOptions{
		Name: "Checkout", Steps: "charge,reserve_stock,book_courier",
	})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	src := readSaga(t, root, filepath.Join("internal", "sagas", "checkout.go"))
	mustParse(t, "checkout.go", src)
	for _, want := range []string{
		`Name: "checkout"`,
		`Name: "charge"`,
		`Name: "reserve_stock"`,
		`Name: "book_courier"`,
		`"demo/apps/api/internal/saga"`,
		"func Checkout() saga.Definition {",
		"type CheckoutInput struct {",
		"Undo: func(ctx context.Context, r *saga.Run) error {",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the generated saga is missing %q", want)
		}
	}

	// The steps are in the order they were given. A saga whose steps ran in map
	// order would compensate in the wrong order too, which is the one thing it
	// exists to get right.
	charge := strings.Index(src, `Name: "charge"`)
	reserve := strings.Index(src, `Name: "reserve_stock"`)
	courier := strings.Index(src, `Name: "book_courier"`)
	if !(charge < reserve && reserve < courier) {
		t.Error("the steps are not in the order they were given")
	}

	registry := readSaga(t, root, filepath.Join("internal", "sagas", "sagas.go"))
	mustParse(t, "sagas.go", registry)
	if !strings.Contains(registry, "saga.Register(Checkout())") {
		t.Error("the saga was written and never registered, so it would never run")
	}
	// The first saga is the one that adds the import, because an empty registry
	// with an unused import does not compile.
	if !strings.Contains(registry, `"demo/apps/api/internal/saga"`) {
		t.Error("the registry has no import for the saga package")
	}
	if !strings.Contains(registry, "// grit:sagas") {
		t.Error("the marker was consumed, so a second workflow could not be added")
	}
}

func TestGeneratedStepsFailLoudlyUntilTheyAreWritten(t *testing.T) {
	// A stub returning nil is a saga that reports success and does nothing, and
	// the first you hear of it is a customer saying the parcel never came.
	root := sagaProject(t)
	if err := generateWorkflowAt(root, "demo", WorkflowOptions{Name: "Checkout", Steps: "charge"}); err != nil {
		t.Fatal(err)
	}
	src := readSaga(t, root, filepath.Join("internal", "sagas", "checkout.go"))
	if !strings.Contains(src, `return fmt.Errorf("saga checkout: step charge is not implemented")`) {
		t.Error("the Do stub does not fail")
	}
	// The Undo stub does return nil, and that is right: compensating something
	// that was never done is a no-op, and an Undo that failed would park the
	// run on a step nobody has written yet.
	if !strings.Contains(src, "\t\t\t\t\treturn nil\n") {
		t.Error("the Undo stub should be a no-op")
	}
}

func TestASecondWorkflowJoinsTheFirst(t *testing.T) {
	root := sagaProject(t)
	for _, name := range []string{"Checkout", "Refund"} {
		if err := generateWorkflowAt(root, "demo", WorkflowOptions{Name: name, Steps: "one,two"}); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
	}
	registry := readSaga(t, root, filepath.Join("internal", "sagas", "sagas.go"))
	mustParse(t, "sagas.go", registry)
	if !strings.Contains(registry, "saga.Register(Checkout())") || !strings.Contains(registry, "saga.Register(Refund())") {
		t.Fatalf("both sagas should be registered:\n%s", registry)
	}
	if strings.Count(registry, `"demo/internal/saga"`) != 1 {
		t.Error("the import was added twice")
	}
}

func TestGenerateWorkflowRefusesWhatCannotWork(t *testing.T) {
	for name, opts := range map[string]WorkflowOptions{
		"no name":                          {Steps: "a,b"},
		"a name that is not an identifier": {Name: "my-saga!", Steps: "a"},
		"no steps":                         {Name: "Checkout"},
		"blank steps":                      {Name: "Checkout", Steps: " , , "},
		"a step that is not an identifier": {Name: "Checkout", Steps: "charge,2nd step!"},
		// Two steps with one name share a row in saga_steps, and a compensation
		// would undo whichever the query happened to return.
		"two steps with one name": {Name: "Checkout", Steps: "charge,charge"},
	} {
		if err := generateWorkflowAt(sagaProject(t), "demo", opts); err == nil {
			t.Errorf("%s was accepted", name)
		}
	}
}

func TestGenerateWorkflowRefusesToOverwriteYourSaga(t *testing.T) {
	root := sagaProject(t)
	if err := generateWorkflowAt(root, "demo", WorkflowOptions{Name: "Checkout", Steps: "a"}); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "internal", "sagas", "checkout.go")
	mine := readSaga(t, root, filepath.Join("internal", "sagas", "checkout.go")) + "\n// my work\n"
	if err := os.WriteFile(path, []byte(mine), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := generateWorkflowAt(root, "demo", WorkflowOptions{Name: "Checkout", Steps: "a,b"}); err == nil {
		t.Fatal("regenerating overwrote the saga")
	}
	if !strings.Contains(readSaga(t, root, filepath.Join("internal", "sagas", "checkout.go")), "// my work") {
		t.Error("the edit was lost")
	}
}

func TestGenerateWorkflowRefusesAProjectWithoutTheEngine(t *testing.T) {
	// A project from before sagas. Writing the file anyway would leave it
	// importing a package that is not there, which is a build failure with
	// nothing in it to point at the cause.
	root := sagaProject(t)
	if err := os.Remove(filepath.Join(root, "internal", "saga", "saga.go")); err != nil {
		t.Fatal(err)
	}
	err := generateWorkflowAt(root, "demo", WorkflowOptions{Name: "Checkout", Steps: "a"})
	if err == nil || !strings.Contains(err.Error(), "grit upgrade") {
		t.Fatalf("err = %v, want one that says to upgrade", err)
	}
}

func TestGenerateWorkflowRefusesAProjectWithoutTheRegistry(t *testing.T) {
	root := sagaProject(t)
	if err := os.Remove(filepath.Join(root, "internal", "sagas", "sagas.go")); err != nil {
		t.Fatal(err)
	}
	err := generateWorkflowAt(root, "demo", WorkflowOptions{Name: "Checkout", Steps: "a"})
	if err == nil || !strings.Contains(err.Error(), "never run") {
		t.Fatalf("err = %v, want one that says the saga would never run", err)
	}
}
