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

// ── Outbox squash ────────────────────────────────────────────────────────
//
// The outbox holds at most one entry per (model, entity_id), so a user who
// edits the same row nine times offline sends one change rather than nine.
// Six rules decide what the surviving entry looks like, and each one is a bug
// with a different shape if it is wrong: a create that reaches the server as an
// update 404s, a delete of a row the server never heard of is a round trip to
// arrive at nothing, and an op that does not survive a later edit sends the
// wrong verb entirely.
//
// None of this is visible from the outside. The outbox drains, the push
// succeeds, and the damage is a row that exists on one device and not another.

func outboxFor(t *testing.T, e *Engine, table, id string) (Outbox, bool) {
	t.Helper()
	var row Outbox
	err := e.DB.Where("model = ? AND entity_id = ?", table, id).First(&row).Error
	if err != nil {
		return Outbox{}, false
	}
	return row, true
}

func outboxCount(t *testing.T, e *Engine) int64 {
	t.Helper()
	var n int64
	if err := e.DB.Model(&Outbox{}).Count(&n).Error; err != nil {
		t.Fatalf("counting outbox: %v", err)
	}
	return n
}

// Rule: new "update" + existing "create" keeps the op as create.
//
// The server has never heard of this row. Sending it an update would 404 and
// the entry would stick in the outbox forever.
func TestSquashUpdateOnTopOfCreateStaysACreate(t *testing.T) {
	e, _ := engineFor(t, pushHandler(t))

	if err := e.LocalCreate("notes", "n1", map[string]interface{}{"title": "first"}); err != nil {
		t.Fatalf("LocalCreate: %v", err)
	}
	if err := e.LocalUpdate("notes", "n1", map[string]interface{}{"title": "second"}); err != nil {
		t.Fatalf("LocalUpdate: %v", err)
	}

	if n := outboxCount(t, e); n != 1 {
		t.Fatalf("outbox holds %d entries, want 1: edits to one row must squash", n)
	}
	row, ok := outboxFor(t, e, "notes", "n1")
	if !ok {
		t.Fatal("no outbox entry for the row")
	}
	if row.Op != "create" {
		t.Errorf("op = %q, want create: the server has never seen this row, so an update would 404", row.Op)
	}
	var data map[string]interface{}
	if err := json.Unmarshal(row.Data, &data); err != nil {
		t.Fatalf("decoding outbox data: %v", err)
	}
	if data["title"] != "second" {
		t.Errorf("title = %v, want second: the squashed entry carries the latest state", data["title"])
	}
}

// Rule: new "update" + existing "update" keeps one entry with the latest data.
func TestSquashRepeatedUpdatesCollapseToOne(t *testing.T) {
	e, _ := engineFor(t, pushHandler(t))

	e.LocalCreate("notes", "n2", map[string]interface{}{"title": "v1"})
	// Pretend the create was already pushed, so the next edits are updates.
	e.DB.Where("entity_id = ?", "n2").Delete(&Outbox{})

	for _, title := range []string{"v2", "v3", "v4"} {
		if err := e.LocalUpdate("notes", "n2", map[string]interface{}{"title": title}); err != nil {
			t.Fatalf("LocalUpdate %s: %v", title, err)
		}
	}

	if n := outboxCount(t, e); n != 1 {
		t.Fatalf("outbox holds %d entries, want 1: three edits to one row are one change to push", n)
	}
	row, _ := outboxFor(t, e, "notes", "n2")
	if row.Op != "update" {
		t.Errorf("op = %q, want update", row.Op)
	}
	var data map[string]interface{}
	json.Unmarshal(row.Data, &data)
	if data["title"] != "v4" {
		t.Errorf("title = %v, want v4: the last edit wins", data["title"])
	}
}

// Rule: new "delete" + existing "create" removes the entry entirely.
//
// The user made something offline and changed their mind. The server never
// knew. Sending a create followed by a delete is two round trips to arrive at
// having done nothing, and a window where other devices see a row that was
// never meant to exist.
func TestSquashDeleteOfAnUnsentCreateCancelsBoth(t *testing.T) {
	e, _ := engineFor(t, pushHandler(t))

	if err := e.LocalCreate("notes", "n3", map[string]interface{}{"title": "regretted"}); err != nil {
		t.Fatalf("LocalCreate: %v", err)
	}
	if n := outboxCount(t, e); n != 1 {
		t.Fatalf("outbox holds %d entries after the create, want 1", n)
	}

	if err := e.LocalDelete("notes", "n3"); err != nil {
		t.Fatalf("LocalDelete: %v", err)
	}

	if n := outboxCount(t, e); n != 0 {
		t.Errorf("outbox holds %d entries, want 0: a create the server never saw, then deleted, is nothing to send", n)
	}
}

// Rule: new "delete" + existing "update" becomes a delete and drops the data.
func TestSquashDeleteAfterUpdateBecomesADelete(t *testing.T) {
	e, _ := engineFor(t, pushHandler(t))

	e.LocalCreate("notes", "n4", map[string]interface{}{"title": "v1"})
	e.DB.Where("entity_id = ?", "n4").Delete(&Outbox{}) // the create was pushed

	if err := e.LocalUpdate("notes", "n4", map[string]interface{}{"title": "v2"}); err != nil {
		t.Fatalf("LocalUpdate: %v", err)
	}
	if err := e.LocalDelete("notes", "n4"); err != nil {
		t.Fatalf("LocalDelete: %v", err)
	}

	if n := outboxCount(t, e); n != 1 {
		t.Fatalf("outbox holds %d entries, want 1", n)
	}
	row, _ := outboxFor(t, e, "notes", "n4")
	if row.Op != "delete" {
		t.Errorf("op = %q, want delete: the edit is moot once the row is gone", row.Op)
	}
	if len(row.Data) != 0 {
		t.Errorf("data = %q, want empty: a delete carries no payload", string(row.Data))
	}
}

// Rule: a delete with no existing entry is inserted as a delete.
func TestSquashDeleteOfASyncedRowQueuesADelete(t *testing.T) {
	e, _ := engineFor(t, pushHandler(t))

	e.LocalCreate("notes", "n5", map[string]interface{}{"title": "v1"})
	e.DB.Where("entity_id = ?", "n5").Delete(&Outbox{}) // the create was pushed

	if err := e.LocalDelete("notes", "n5"); err != nil {
		t.Fatalf("LocalDelete: %v", err)
	}

	row, ok := outboxFor(t, e, "notes", "n5")
	if !ok {
		t.Fatal("deleting a synced row must queue a delete, or the server never hears about it")
	}
	if row.Op != "delete" {
		t.Errorf("op = %q, want delete", row.Op)
	}
}

// Squashing is per row. Two rows edited offline are two changes, not one.
func TestSquashIsPerRowNotGlobal(t *testing.T) {
	e, _ := engineFor(t, pushHandler(t))

	e.LocalCreate("notes", "a", map[string]interface{}{"title": "A"})
	e.LocalCreate("notes", "b", map[string]interface{}{"title": "B"})
	e.LocalUpdate("notes", "a", map[string]interface{}{"title": "A2"})

	if n := outboxCount(t, e); n != 2 {
		t.Errorf("outbox holds %d entries, want 2: squashing collapses edits to one row, not the whole queue", n)
	}
}

// The same id in two different models is two different rows.
func TestSquashKeysOnModelAndID(t *testing.T) {
	e, _ := engineFor(t, pushHandler(t))

	e.LocalCreate("notes", "shared-id", map[string]interface{}{"title": "a note"})
	e.LocalCreate("tasks", "shared-id", map[string]interface{}{"title": "a task"})

	if n := outboxCount(t, e); n != 2 {
		t.Errorf("outbox holds %d entries, want 2: the key is (model, id), so the same id in two tables is two rows", n)
	}
}

// A local write is readable immediately and survives the engine being reopened,
// which is the whole reason the outbox is on disk rather than in memory.
func TestOutboxSurvivesAReopen(t *testing.T) {
	dir := t.TempDir()
	srv := httptest.NewServer(pushHandler(t))
	t.Cleanup(srv.Close)

	open := func() *Engine {
		e, err := Open(filepath.Join(dir, "local.db"), srv.URL+"/api",
			func() (string, error) { return "test-token", nil })
		if err != nil {
			t.Fatalf("Open: %v", err)
		}
		return e
	}

	first := open()
	if err := first.LocalCreate("notes", "persist", map[string]interface{}{"title": "written offline"}); err != nil {
		t.Fatalf("LocalCreate: %v", err)
	}
	if sqlDB, err := first.DB.DB(); err == nil {
		_ = sqlDB.Close()
	}

	second := open()
	t.Cleanup(func() {
		if sqlDB, err := second.DB.DB(); err == nil {
			_ = sqlDB.Close()
		}
	})

	if n := outboxCount(t, second); n != 1 {
		t.Errorf("outbox holds %d entries after reopening, want 1: an in-memory outbox loses the work it exists to protect", n)
	}
	row, err := second.LocalGet("notes", "persist")
	if err != nil {
		t.Fatalf("LocalGet after reopen: %v", err)
	}
	if row == nil {
		t.Fatal("the local mirror lost the row across a reopen")
	}
	if row["title"] != "written offline" {
		t.Errorf("title = %v, want the offline edit", row["title"])
	}
}
`
}
