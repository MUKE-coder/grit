package scaffold

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDiscoverSyncModels(t *testing.T) {
	root := t.TempDir()
	routes := filepath.Join(root, "apps", "api", "internal", "routes", "routes.go")
	if err := os.MkdirAll(filepath.Dir(routes), 0o755); err != nil {
		t.Fatal(err)
	}
	body := `package routes

func Setup() {
	syncRegistry.Register("users", &models.User{})
	syncRegistry.Register("blog_posts", &models.BlogPost{})
	syncRegistry.Register( "products", &models.Product{})
	// grit:sync
}
`
	if err := os.WriteFile(routes, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	got := discoverSyncModels(root)
	want := []string{"blog_posts", "products", "users"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v (sorted)", got, want)
		}
	}
}

// A single-app project keeps routes.go at the root rather than under apps/api.
func TestDiscoverSyncModelsFindsSingleAppLayout(t *testing.T) {
	root := t.TempDir()
	routes := filepath.Join(root, "internal", "routes", "routes.go")
	if err := os.MkdirAll(filepath.Dir(routes), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(routes, []byte(`syncRegistry.Register("orders", &models.Order{})`), 0o644); err != nil {
		t.Fatal(err)
	}

	got := discoverSyncModels(root)
	if len(got) != 1 || got[0] != "orders" {
		t.Fatalf("got %v, want [orders]", got)
	}
}

func TestDiscoverSyncModelsOnAProjectWithNoResources(t *testing.T) {
	if got := discoverSyncModels(t.TempDir()); len(got) != 0 {
		t.Fatalf("got %v, want none", got)
	}
}

func TestAddWorkspaceDependency(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "package.json")
	original := `{
  "name": "web",
  "dependencies": {
    "next": "15.0.0",
    "react": "19.2.7"
  },
  "devDependencies": {
    "typescript": "5.6.0"
  }
}
`
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	added, err := addWorkspaceDependency(path, "@myapp/sync")
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	if !added {
		t.Fatal("expected the dependency to be added")
	}

	data, _ := os.ReadFile(path)
	var parsed struct {
		Name         string            `json:"name"`
		Dependencies map[string]string `json:"dependencies"`
		DevDeps      map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("the edited package.json no longer parses: %v\n%s", err, data)
	}
	if parsed.Dependencies["@myapp/sync"] != "workspace:*" {
		t.Fatalf("dependency missing: %v", parsed.Dependencies)
	}
	// The point of editing text rather than re-encoding: everything else is
	// exactly where it was, including key order and the dev dependencies.
	if parsed.Dependencies["next"] != "15.0.0" || parsed.DevDeps["typescript"] != "5.6.0" {
		t.Fatalf("the rest of the file was disturbed: %s", data)
	}

	// Running it twice must not add it twice.
	added, err = addWorkspaceDependency(path, "@myapp/sync")
	if err != nil {
		t.Fatalf("second add: %v", err)
	}
	if added {
		t.Fatal("the dependency was added a second time")
	}
	after, _ := os.ReadFile(path)
	if strings.Count(string(after), "@myapp/sync") != 1 {
		t.Fatalf("duplicate entry:\n%s", after)
	}
}

func TestAddWorkspaceDependencyRefusesAMalformedPackageJSON(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "package.json")
	broken := `{"name": "web", "dependencies": {`
	if err := os.WriteFile(path, []byte(broken), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := addWorkspaceDependency(path, "@myapp/sync"); err == nil {
		t.Fatal("expected an error rather than an edit to an unparseable file")
	}
	data, _ := os.ReadFile(path)
	if string(data) != broken {
		t.Fatalf("the broken file was modified: %s", data)
	}
}

// The templates carry TypeScript inside Go raw strings, where a stray
// backtick would silently terminate the literal and take the rest of the file
// with it. Nothing else in the build catches that.
func TestSyncTemplatesContainNoBackticks(t *testing.T) {
	templates := map[string]string{
		"types.ts":            syncTypesTS(),
		"engine.ts":           syncEngineTS(),
		"react.tsx":           syncReactTS(),
		"adapters/memory.ts":  syncMemoryAdapterTS(),
		"adapters/indexed.ts": syncIndexedAdapterTS(),
		"adapters/sqlite.ts":  syncSQLiteAdapterTS(),
	}
	for name, body := range templates {
		if strings.Contains(body, "`") {
			t.Errorf("%s contains a backtick", name)
		}
		if strings.ContainsRune(body, 0) {
			t.Errorf("%s contains a NUL byte", name)
		}
		if body == "" {
			t.Errorf("%s is empty", name)
		}
	}
}

// The engine and the adapters have to agree on the StorageAdapter surface, and
// a method added to one and not the others fails at typecheck rather than
// here. What this checks is cheaper and still worth having: every method the
// interface declares is implemented by all three adapters.
func TestEveryAdapterImplementsTheWholeInterface(t *testing.T) {
	methods := []string{
		"open(", "close(",
		"putRecord(", "getRecord(", "listRecords(", "deleteRecord(", "clearModel(",
		"putOutbox(", "getOutbox(", "listOutbox(", "deleteOutbox(",
		"getMeta(", "setMeta(",
	}
	adapters := map[string]string{
		"memory":  syncMemoryAdapterTS(),
		"indexed": syncIndexedAdapterTS(),
		"sqlite":  syncSQLiteAdapterTS(),
	}
	for name, body := range adapters {
		for _, method := range methods {
			if !strings.Contains(body, method) {
				t.Errorf("%s adapter is missing %s", name, method)
			}
		}
	}
}

// The client and the server have to agree on the wire, and they are written in
// different languages in different files. These are the field names that
// carry the conflict state; a rename on one side and not the other turns every
// conflict into a silent overwrite.
func TestClientSpeaksTheServersWireFormat(t *testing.T) {
	types := syncTypesTS()
	engine := syncEngineTS()
	server := apiSyncHandlerGo()

	for _, field := range []string{
		"server_version", "server_data", "new_version", "_deleted",
	} {
		if !strings.Contains(types+engine, field) {
			t.Errorf("the TypeScript client never mentions %q", field)
		}
		if !strings.Contains(server, field) {
			t.Errorf("the Go server never mentions %q", field)
		}
	}

	if !strings.Contains(engine, "VERSION_CONFLICT") {
		t.Error("the client does not handle VERSION_CONFLICT")
	}
	if !strings.Contains(server, "VERSION_CONFLICT") {
		t.Error("the server does not emit VERSION_CONFLICT")
	}
}

func TestSyncSetupFilesNameTheDiscoveredModels(t *testing.T) {
	models := []string{"products", "orders"}

	web := webSyncSetupTS("shopkit", models)
	expo := expoSyncSetupTS("shopkit", models)

	for _, body := range []string{web, expo} {
		for _, model := range models {
			if !strings.Contains(body, "\""+model+"\"") {
				t.Errorf("setup file does not mirror %q:\n%s", model, body)
			}
		}
		if !strings.Contains(body, "@shopkit/sync") {
			t.Errorf("setup file does not import the workspace package:\n%s", body)
		}
	}

	// The web app renders on a server that has no IndexedDB, so it has to pick
	// its adapter at runtime rather than assume the browser.
	if !strings.Contains(web, "MemoryAdapter") {
		t.Error("the web setup has no fallback for server rendering")
	}
	// Expo opens SQLite asynchronously, so its setup is a factory.
	if !strings.Contains(expo, "async function createSyncEngine") {
		t.Error("the expo setup should be an async factory")
	}
}
