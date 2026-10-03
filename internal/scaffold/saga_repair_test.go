package scaffold

import (
	"strings"
	"testing"
)

// preSagaRoutes is routes.go as v3.353.0 wrote it: no saga registration, no
// runner, and neither import.
func preSagaRoutes(t *testing.T) string {
	t.Helper()
	fresh := apiRoutesGo()
	old := strings.Replace(fresh, sagaRuntimeBlock, "", 1)
	old = strings.Replace(old, "\t\"{{MODULE}}/internal/saga\"\n", "", 1)
	old = strings.Replace(old, "\t\"{{MODULE}}/internal/sagas\"\n", "", 1)
	if old == fresh || strings.Contains(old, "sagas.Register()") {
		t.Fatal("could not rebuild routes.go from before the saga runner")
	}
	return old
}

// preSagaModels is the model list before the two tables joined it.
func preSagaModels(t *testing.T) string {
	t.Helper()
	fresh := apiUserModelGo()
	old := strings.Replace(fresh, sagaModelRegistration, "", 1)
	if old == fresh || strings.Contains(old, "&SagaRun{}") {
		t.Fatal("could not rebuild the model list from before the saga tables")
	}
	return old
}

func TestRepairSagaModels(t *testing.T) {
	out, fixed, warn := repairSagaModelsSource(preSagaModels(t))
	if len(warn) > 0 || len(fixed) != 1 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	// Without this the engine arrives, the migration reports nothing to do, and
	// the first saga anybody writes inserts into a table that does not exist.
	if !strings.Contains(out, "&SagaRun{}") || !strings.Contains(out, "&SagaStep{}") {
		t.Error("the tables were not added to the model list")
	}
	if out != apiUserModelGo() {
		t.Error("the repaired model list differs from a fresh one")
	}
	if again, fixed, _ := repairSagaModelsSource(out); again != out || len(fixed) > 0 {
		t.Error("a second upgrade added them again")
	}
}

func TestRepairSagaModelsWarnsWhenTheMarkerIsGone(t *testing.T) {
	src := strings.Replace(preSagaModels(t), "\t\t// grit:models\n", "", 1)
	out, fixed, warn := repairSagaModelsSource(src)
	if out != src || len(fixed) > 0 || len(warn) != 1 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
}

func TestRepairSagaRuntime(t *testing.T) {
	out, fixed, warn := repairSagaRuntimeSource("demo")(preSagaRoutes(t))
	if len(warn) > 0 || len(fixed) != 1 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	for _, want := range []string{
		"sagas.Register()",
		"saga.Runner{DB: db}",
		`"demo/internal/saga"`,
		`"demo/internal/sagas"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the repaired routes.go is missing %q", want)
		}
	}
	// Beside the event bus relay, which is where it belongs: both are background
	// delivery that every replica runs.
	relay := strings.Index(out, "events.StartRelay(db)")
	register := strings.Index(out, "sagas.Register()")
	if relay < 0 || register < relay {
		t.Error("the runner is not started after the relay")
	}
	if again, fixed, _ := repairSagaRuntimeSource("demo")(out); again != out || len(fixed) > 0 {
		t.Error("a second upgrade wired it again")
	}
}

func TestRepairSagaRuntimeWarnsOnAnEditedRoutesFile(t *testing.T) {
	src := strings.Replace(preSagaRoutes(t), "\tevents.StartRelay(db)\n", "\tevents.StartRelay(db) // mine\n", 1)
	out, fixed, warn := repairSagaRuntimeSource("demo")(src)
	if out != src || len(fixed) > 0 || len(warn) != 1 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
}

func TestTheRunnerIsNotStartedWhenNothingIsRegistered(t *testing.T) {
	// A project with no sagas would otherwise poll an empty table every five
	// seconds for the life of the process.
	if !strings.Contains(sagaRuntimeBlock, "if len(saga.Registered()) > 0 {") {
		t.Error("the runner starts unconditionally")
	}
}

func TestFreshTemplatesAlreadyHaveTheSagaWiring(t *testing.T) {
	routes := apiRoutesGo()
	for _, want := range []string{"sagas.Register()", `"{{MODULE}}/internal/saga"`, `"{{MODULE}}/internal/sagas"`} {
		if !strings.Contains(routes, want) {
			t.Errorf("a fresh routes.go is missing %q", want)
		}
	}
	if out, fixed, warn := repairSagaRuntimeSource("demo")(routes); out != routes || len(fixed) > 0 || len(warn) > 0 {
		t.Errorf("a fresh routes.go still needs the repair (fixed %v, warn %v)", fixed, warn)
	}
	models := apiUserModelGo()
	if out, fixed, warn := repairSagaModelsSource(models); out != models || len(fixed) > 0 || len(warn) > 0 {
		t.Errorf("a fresh model list still needs the repair (fixed %v, warn %v)", fixed, warn)
	}
}

func TestTheSagaEngineIsFrameworkOwned(t *testing.T) {
	// It has to travel on upgrade: a fix to the engine that reached new projects
	// only would be no fix at all for the one app with money moving through it.
	root := t.TempDir()
	opts := Options{ProjectName: "demo", Architecture: ArchTriple, Frontend: FrontendNext}
	if err := writeFrameworkOwnedFiles(root, opts); err != nil {
		t.Fatalf("writeFrameworkOwnedFiles: %v", err)
	}
	for _, rel := range []string{"internal/saga/saga.go", "internal/saga/saga_test.go", "internal/models/saga_run.go"} {
		if !fileExists(opts.APIRoot(root) + "/" + rel) {
			t.Errorf("%s was not written", rel)
		}
	}
}

func TestTheSagasPackageIsCreatedOnceAndNeverRewritten(t *testing.T) {
	// It holds the developer's sagas and the generator injects into it, so an
	// upgrade that replaced it would delete every saga in the project.
	root := t.TempDir()
	opts := Options{ProjectName: "demo", Architecture: ArchTriple, Frontend: FrontendNext}
	created, err := writeSagaRegistry(root, opts)
	if err != nil || !created {
		t.Fatalf("created %v, err %v", created, err)
	}
	path := opts.APIRoot(root) + "/internal/sagas/sagas.go"
	if err := writeFile(path, "package sagas\n\n// my sagas\nfunc Register() {\n\t// grit:sagas\n}\n"); err != nil {
		t.Fatal(err)
	}
	again, err := writeSagaRegistry(root, opts)
	if err != nil || again {
		t.Fatalf("a second call rewrote it (created %v, err %v)", again, err)
	}
	if !fileContains(path, "// my sagas") {
		t.Error("the developer's registry was overwritten")
	}
}
