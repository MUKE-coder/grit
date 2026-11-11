package scaffold

import (
	"path/filepath"
	"strings"
	"testing"
)

// grit#210, and the trap it very nearly walked into.
//
// internal/sync/policy.go was written to new projects only, while
// internal/handlers/sync.go travels on upgrade. Adding Policy.Writable() and
// calling it from the handler would have delivered a handler calling a method
// the project's own package does not declare, so every upgraded project with
// offline sync would have stopped compiling. The two are one cluster now.
func TestTheSyncPackageTravelsWithItsHandler(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "probe", Architecture: ArchTriple}
	apiRoot := opts.APIRoot(root)

	if err := writeFrameworkOwnedFiles(root, opts); err != nil {
		t.Fatalf("writing the owned set: %v", err)
	}

	handler := filepath.Join(apiRoot, "internal", "handlers", "sync.go")
	if !fileExists(handler) {
		t.Skip("the sync handler is no longer in the set that travels on upgrade")
	}
	// Whatever the handler calls into has to arrive with it.
	for _, name := range []string{"policy.go", "registry.go"} {
		path := filepath.Join(apiRoot, "internal", "sync", name)
		if !fileExists(path) {
			t.Errorf("internal/handlers/sync.go travels on upgrade and internal/sync/%s does "+
				"not, so an upgraded project gets a handler calling into a package the "+
				"framework has moved on from (grit#210)", name)
		}
	}
	// The specific pair: the method and its caller.
	if fileContains(handler, "policy.Writable()") &&
		!fileContains(filepath.Join(apiRoot, "internal", "sync", "policy.go"), "func (p Policy) Writable()") {
		t.Error("the handler calls Writable() and the delivered policy.go does not declare it")
	}
}

// The mode exists, in both languages, with the same list of what cannot be
// written. Two copies of that list is how the server and the client come to
// disagree about which models are writable.
func TestReadOnlyIsDeclaredOnBothSides(t *testing.T) {
	goSrc := apiSyncPolicyGo()
	if !strings.Contains(goSrc, `ModeReadOnly = "read_only"`) {
		t.Error("the Go policy package has no read_only mode")
	}
	if !strings.Contains(goSrc, "func (p Policy) Writable() bool") {
		t.Error("the Go policy has no Writable(), so the push handler has nothing to ask")
	}

	types := syncTypesTS()
	if !strings.Contains(types, `"read_only"`) {
		t.Error("the TypeScript SyncPolicy union has no read_only, so the client cannot type " +
			"the mode the server sends and tsc refuses the comparison")
	}

	engine := syncEngineTS()
	if !strings.Contains(engine, "writable(model: string): boolean") {
		t.Error("the client engine cannot answer whether a model is writable")
	}
	// The three write methods and the queue all ask, so a write that can never
	// succeed never reaches the outbox.
	if n := strings.Count(engine, "assertWritable("); n < 5 {
		t.Errorf("assertWritable appears %d times; create, update, remove, enqueue and its own "+
			"definition is five, and a missing one is a write queued to fail", n)
	}
	// Checked before the mirror is touched, not only before the queue:
	// writing locally and then refusing to queue leaves a row that never
	// reaches the server and never will.
	create := engine[strings.Index(engine, "async create("):]
	mirror := strings.Index(create, "putRecord")
	guard := strings.Index(create, "assertWritable")
	if guard < 0 || mirror < 0 || guard > mirror {
		t.Error("create writes to the local mirror before checking whether the model is writable")
	}
}

// The push is refused by the server too, so an old client build is still told.
func TestThePushHandlerRefusesAReadOnlyModel(t *testing.T) {
	src := apiSyncHandlerGo()
	if !strings.Contains(src, `Code: "READ_ONLY"`) {
		t.Error("a push to a read_only model is not refused with a code of its own")
	}
	// Before the ownership and permission checks: the answer does not depend
	// on who is asking, and a 403 would suggest a different account might
	// manage it.
	body := src[strings.Index(src, "func (h *SyncHandler) applyChange("):]
	refusal := strings.Index(body, `Code: "READ_ONLY"`)
	forbidden := strings.Index(body, `"you may not create "`)
	if refusal < 0 || forbidden < 0 || refusal > forbidden {
		t.Error("the read_only refusal comes after the permission check, so a caller who may " +
			"not create is told that rather than that nobody may")
	}
}

// A read_only model is still mirrored, which is the entire point: she can open
// her own report on a train. Only online_only is refused a pull.
func TestAReadOnlyModelIsStillPullable(t *testing.T) {
	src := apiSyncHandlerGo()
	pull := src[strings.Index(src, "func (h *SyncHandler) Pull("):]
	if strings.Contains(pull, "ModeReadOnly") {
		t.Error("the pull handler refuses read_only, which leaves the user unable to read " +
			"records about her own work offline: the thing online_only already did wrong")
	}
	if !strings.Contains(pull, "ModeOnlineOnly") {
		t.Error("the pull handler no longer refuses online_only")
	}
}
