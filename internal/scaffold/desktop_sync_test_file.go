package scaffold

// desktopSyncEngineTestGo emits apps/desktop/sync/engine_test.go.
//
// The offline engine is the least testable-by-hand and most consequential part
// of the desktop app: it decides what happens to a change made on a plane, and
// what happens when the server moved underneath it. It shipped with no tests,
// so every one of those decisions was only ever verified by using the app and
// hoping.
//
// These run against an httptest server standing in for /api/sync, so the whole
// push path is exercised without a backend: the happy case, a version
// conflict, the fact that a conflicted row is not retried behind the user's
// back, and that resolving one replays it at the version the user actually
// saw.
func desktopSyncEngineTestGo(opts Options) string {
	return `package sync

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

// engineFor opens an engine against a stub server and a scratch database.
func engineFor(t *testing.T, handler http.HandlerFunc) (*Engine, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	e, err := Open(filepath.Join(t.TempDir(), "local.db"), srv.URL+"/api",
		func() (string, error) { return "test-token", nil })
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	// Release the SQLite file before t.TempDir() removes the directory.
	// Cleanups run last-registered-first, so this closes ahead of the delete.
	// On Windows an open handle makes the removal fail and the test with it.
	t.Cleanup(func() {
		if sqlDB, err := e.DB.DB(); err == nil {
			_ = sqlDB.Close()
		}
	})
	return e, srv
}

// pushHandler answers /sync/push with one canned result per change.
func pushHandler(t *testing.T, results ...PushResult) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/sync/push" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(` + "`" + `{"changes":[],"cursor":""}` + "`" + `))
			return
		}
		var batch PushBatch
		if err := json.NewDecoder(r.Body).Decode(&batch); err != nil {
			t.Errorf("decoding push batch: %v", err)
		}
		if len(batch.Changes) != len(results) {
			t.Errorf("server got %d changes, the test canned %d results",
				len(batch.Changes), len(results))
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"results": results})
	}
}

func outboxRow(t *testing.T, e *Engine, model, id string) Outbox {
	t.Helper()
	var row Outbox
	if err := e.DB.Where("model = ? AND entity_id = ?", model, id).First(&row).Error; err != nil {
		t.Fatalf("no outbox row for %s/%s: %v", model, id, err)
	}
	return row
}

// A change made offline is readable immediately and queued for later.
//
// Both halves matter. Writing only the outbox would mean the row vanished from
// the user's own screen until they reconnected, which is the opposite of what
// offline-first is for.
func TestLocalWriteIsVisibleBeforeItIsPushed(t *testing.T) {
	e, _ := engineFor(t, pushHandler(t))

	const id = "task-1"
	if err := e.LocalCreate("tasks", id, map[string]interface{}{"title": "Written on a plane"}); err != nil {
		t.Fatalf("LocalCreate: %v", err)
	}

	got, err := e.LocalGet("tasks", id)
	if err != nil {
		t.Fatalf("LocalGet: %v", err)
	}
	if got["title"] != "Written on a plane" {
		t.Errorf("the local mirror does not show the change: %v", got)
	}

	if n, _ := e.PendingCount(); n != 1 {
		t.Errorf("pending count is %d, want 1: the change was not queued", n)
	}
}

// A successful push clears the queue and records the server's new version.
//
// Leaving the row behind would push it again on the next sync, turning one
// edit into two.
func TestSuccessfulPushClearsTheQueue(t *testing.T) {
	e, _ := engineFor(t, pushHandler(t, PushResult{OK: true, NewVersion: 7}))
	const id = "task-1"
	if err := e.LocalCreate("tasks", id, map[string]interface{}{"title": "x"}); err != nil {
		t.Fatalf("LocalCreate: %v", err)
	}

	pushed, conflicts, err := e.Push()
	if err != nil {
		t.Fatalf("Push: %v", err)
	}
	if pushed != 1 || conflicts != 0 {
		t.Fatalf("pushed=%d conflicts=%d, want 1 and 0", pushed, conflicts)
	}
	if n, _ := e.PendingCount(); n != 0 {
		t.Errorf("%d change(s) still queued after a successful push; they would "+
			"be pushed again and applied twice", n)
	}

	var rec Record
	if err := e.DB.Where("model = ? AND id = ?", "tasks", id).First(&rec).Error; err != nil {
		t.Fatalf("the local mirror lost the row after a successful push: %v", err)
	}
	if rec.Version != 7 {
		t.Errorf("local version is %d, want the server's 7: the next edit would "+
			"be sent with a stale version and conflict for no reason", rec.Version)
	}
}

// A version conflict parks the change with the server's state attached.
//
// The user has to be able to see what they are choosing between, so the
// server's copy is stored rather than only the fact that it differed.
func TestVersionConflictParksTheChangeWithServerState(t *testing.T) {
	e, _ := engineFor(t, pushHandler(t, PushResult{
		OK:            false,
		Code:          "VERSION_CONFLICT",
		Message:       "someone else edited this",
		ServerVersion: 4,
		ServerData:    map[string]interface{}{"title": "Their title"},
	}))
	const id = "task-1"
	if err := e.LocalCreate("tasks", id, map[string]interface{}{"title": "My title"}); err != nil {
		t.Fatalf("LocalCreate: %v", err)
	}

	pushed, conflicts, err := e.Push()
	if err != nil {
		t.Fatalf("Push: %v", err)
	}
	if pushed != 0 || conflicts != 1 {
		t.Fatalf("pushed=%d conflicts=%d, want 0 and 1", pushed, conflicts)
	}

	row := outboxRow(t, e, "tasks", id)
	if !row.HasConflict {
		t.Error("the row is not flagged, so the UI has no way to raise it")
	}
	if row.ServerVersion != 4 {
		t.Errorf("server version is %d, want 4", row.ServerVersion)
	}
	var server map[string]interface{}
	_ = json.Unmarshal(row.ServerData, &server)
	if server["title"] != "Their title" {
		t.Errorf("the server's copy was not kept, so the merge dialog has "+
			"nothing to show: %v", server)
	}
	// And the local edit must survive: losing it would mean the conflict
	// resolved itself in the server's favour without asking.
	var mine map[string]interface{}
	_ = json.Unmarshal(row.Data, &mine)
	if mine["title"] != "My title" {
		t.Errorf("the local edit was lost: %v", mine)
	}
}

// A conflicted change is not pushed again until it is resolved.
//
// Retrying would either fail forever or, worse, succeed once the versions
// happened to line up and silently overwrite the other edit.
func TestConflictedChangeIsNotRetriedSilently(t *testing.T) {
	// Counts what the server is actually asked to push, across both attempts.
	var sent int
	e, _ := engineFor(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/sync/push" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(` + "`" + `{"changes":[],"cursor":""}` + "`" + `))
			return
		}
		var batch PushBatch
		_ = json.NewDecoder(r.Body).Decode(&batch)
		sent += len(batch.Changes)

		results := make([]PushResult, len(batch.Changes))
		for i := range results {
			results[i] = PushResult{
				OK: false, Code: "VERSION_CONFLICT", ServerVersion: 2,
				ServerData: map[string]interface{}{"title": "Theirs"},
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"results": results})
	})

	if err := e.LocalCreate("tasks", "task-1", map[string]interface{}{"title": "Mine"}); err != nil {
		t.Fatalf("LocalCreate: %v", err)
	}
	if _, _, err := e.Push(); err != nil {
		t.Fatalf("first push: %v", err)
	}
	if sent != 1 {
		t.Fatalf("the first push sent %d change(s), want 1", sent)
	}

	pushed, conflicts, err := e.Push()
	if err != nil {
		t.Fatalf("second push: %v", err)
	}
	if sent != 1 {
		t.Errorf("the conflicted row was sent again (%d total): retrying either "+
			"fails forever or succeeds once the versions happen to line up, "+
			"silently overwriting the other edit", sent)
	}
	if pushed != 0 || conflicts != 0 {
		t.Errorf("second push reported pushed=%d conflicts=%d, want 0 and 0",
			pushed, conflicts)
	}
}

// Resolving a conflict replays the merge at the version the user saw.
//
// Sending the old local version would conflict again immediately; sending no
// version would overwrite whatever arrived in between. The version the user
// was shown is the only correct one.
func TestResolvingAConflictReplaysAtTheVersionTheUserSaw(t *testing.T) {
	e, _ := engineFor(t, pushHandler(t, PushResult{
		OK: false, Code: "VERSION_CONFLICT", ServerVersion: 9,
		ServerData: map[string]interface{}{"title": "Theirs"},
	}))
	const id = "task-1"
	if err := e.LocalCreate("tasks", id, map[string]interface{}{"title": "Mine"}); err != nil {
		t.Fatalf("LocalCreate: %v", err)
	}
	if _, _, err := e.Push(); err != nil {
		t.Fatalf("push: %v", err)
	}

	merged := map[string]interface{}{"title": "Mine, keeping theirs in mind"}
	if err := e.ResolveConflict("tasks", id, merged, 9); err != nil {
		t.Fatalf("ResolveConflict: %v", err)
	}

	row := outboxRow(t, e, "tasks", id)
	if row.HasConflict {
		t.Error("the row is still flagged after being resolved, so it will " +
			"never be pushed")
	}
	if row.Version != 9 {
		t.Errorf("the replay is queued at version %d, want the server's 9: "+
			"it would conflict again on the next push", row.Version)
	}
	var data map[string]interface{}
	_ = json.Unmarshal(row.Data, &data)
	if data["title"] != merged["title"] {
		t.Errorf("the merged value was not stored: %v", data)
	}
}
`
}
