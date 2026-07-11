package scaffold

// Source for the apps/desktop/sync/ Go package — the offline-first
// engine that runs inside the Wails app. Held in three files:
//
//   engine.go  — Engine struct, Open, Sync(), Push(), Pull(), helpers
//   outbox.go  — Outbox + Record GORM models, squash semantics
//   local.go   — LocalCreate/Update/Delete/Get/List
//
// Together they form a self-contained package the Wails main can import
// and bind to the frontend. The package is `package sync` (separate
// from the root `package main` so it can be imported normally).

func desktopSyncEngineGo() string {
	return `// Package sync is the offline-first engine for the desktop app.
//
// Data model:
//   - records — local mirror of every row the user has touched. Reads
//     come from this table; pulls UPSERT into it.
//   - outbox  — pending local changes that haven't been pushed yet.
//     One entry per (table_name, entity_id) — multiple edits squash.
//
// Wire flow:
//   - LocalCreate/Update/Delete writes both records (so reads see the
//     change immediately) AND an outbox entry.
//   - Sync() runs Pull then Push:
//       Pull updates records from /api/sync/pull?model=X&since=...
//       Push posts the outbox to /api/sync/push and applies results.
//   - Conflicts surface as outbox rows with HasConflict=true so the UI
//     can present a per-field merge dialog. ResolveConflict(...) writes
//     the merged data back into the outbox and clears the flag.
package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	gosync "sync"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Engine owns the local SQLite database and the HTTP transport used to
// talk to /api/sync. One Engine per running app.
type Engine struct {
	DB        *gorm.DB
	APIURL    string
	GetToken  func() (string, error)
	HTTP      *http.Client
	cursors   map[string]string // last pull cursor per model
	cursorsMu gosync.RWMutex

	// Background auto-sync state (offline-hybrid UX).
	autoMu     gosync.Mutex
	stopAuto   chan struct{}
	syncModels []string
	lastSync   time.Time
	lastErr    string
	online     bool
	deviceID   string
}

// Open initializes a sync engine. dbPath is the absolute path to the
// SQLite file; apiURL is the base API URL ending in "/api". getToken
// returns the user's current bearer token (called per request).
func Open(dbPath, apiURL string, getToken func() (string, error)) (*Engine, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("creating data dir: %w", err)
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("opening local db: %w", err)
	}
	if err := db.AutoMigrate(&Record{}, &Outbox{}, &Cursor{}, &Setting{}); err != nil {
		return nil, fmt.Errorf("migrating local db: %w", err)
	}
	e := &Engine{
		DB:       db,
		APIURL:   apiURL,
		GetToken: getToken,
		HTTP:     &http.Client{Timeout: 30 * time.Second},
		cursors:  make(map[string]string),
	}
	// A stable per-install device id, generated once and persisted. Shown on
	// the Sync page and useful for correlating a device's changes server-side.
	e.deviceID = e.getSetting("device_id")
	if e.deviceID == "" {
		e.deviceID = uuid.New().String()
		_ = e.setSetting("device_id", e.deviceID)
	}
	return e, nil
}

// SyncResult summarizes one Sync run for the UI.
//
// StartedAt/FinishedAt are RFC3339 strings, not time.Time. This struct crosses
// the Wails boundary, and Wails' TypeScript binding generator cannot resolve
// time.Time ("Not found: time.Time") — it then drops the models AND every App
// method that mentions them, so the frontend loses Sync/LocalCreate/... The
// JSON shape is identical either way (time.Time already marshals to RFC3339).
type SyncResult struct {
	Pushed     int      ` + "`" + `json:"pushed"` + "`" + `
	Pulled     int      ` + "`" + `json:"pulled"` + "`" + `
	Conflicts  int      ` + "`" + `json:"conflicts"` + "`" + `
	Errors     []string ` + "`" + `json:"errors,omitempty"` + "`" + `
	StartedAt  string   ` + "`" + `json:"started_at"` + "`" + `
	FinishedAt string   ` + "`" + `json:"finished_at"` + "`" + `
}

// Sync runs Pull → Push. Pull first so the user pushes against the
// freshest server state and conflict surface area is minimized. Models
// is the list of table names to pull; pass empty to skip pull.
func (e *Engine) Sync(models []string) (*SyncResult, error) {
	res := &SyncResult{StartedAt: time.Now().Format(time.RFC3339)}

	for _, m := range models {
		n, err := e.Pull(m)
		if err != nil {
			res.Errors = append(res.Errors, fmt.Sprintf("pull %s: %v", m, err))
			continue
		}
		res.Pulled += n
	}

	pushed, conflicts, err := e.Push()
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("push: %v", err))
	}
	res.Pushed = pushed
	res.Conflicts = conflicts
	res.FinishedAt = time.Now().Format(time.RFC3339)
	return res, nil
}

// Setting is a tiny key/value table for engine flags that must survive an app
// restart — currently just the manual "Work offline" toggle.
type Setting struct {
	Key   string ` + "`" + `gorm:"primarykey"` + "`" + `
	Value string
}

func (Setting) TableName() string { return "sync_settings" }

func (e *Engine) getSetting(key string) string {
	var s Setting
	if err := e.DB.First(&s, "key = ?", key).Error; err != nil {
		return ""
	}
	return s.Value
}

func (e *Engine) setSetting(key, value string) error {
	return e.DB.Save(&Setting{Key: key, Value: value}).Error
}

// SetForceOffline persists the manual "Work offline" toggle. While on, the
// background loop stops syncing and every write just queues locally; flip it
// back off and the next tick (or SyncNow) reconciles.
func (e *Engine) SetForceOffline(v bool) error {
	val := "0"
	if v {
		val = "1"
	}
	return e.setSetting("force_offline", val)
}

// IsForceOffline reports the persisted manual-offline toggle.
func (e *Engine) IsForceOffline() bool { return e.getSetting("force_offline") == "1" }

// Reachable pings /api/health to tell "the server is down / no network" apart
// from "the user chose offline". Short timeout so the UI stays responsive.
func (e *Engine) Reachable() bool {
	req, err := http.NewRequest(http.MethodGet, e.APIURL+"/health", nil)
	if err != nil {
		return false
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode < 500
}

// StartAutoSync launches the background mirror loop: every interval, if the
// server is reachable AND the user hasn't chosen offline, it runs Sync (Pull
// then Push). This is what continuously mirrors server data locally while
// online, and what auto-reconciles offline edits the moment you're back on.
// Safe to call once at startup; a second call is a no-op.
func (e *Engine) StartAutoSync(models []string, interval time.Duration) {
	e.autoMu.Lock()
	if e.stopAuto != nil {
		e.autoMu.Unlock()
		return
	}
	stop := make(chan struct{})
	e.stopAuto = stop
	e.syncModels = models
	e.autoMu.Unlock()

	go func() {
		e.tick() // sync shortly after startup so the mirror is fresh
		t := time.NewTicker(interval)
		defer t.Stop()
		for {
			select {
			case <-stop:
				return
			case <-t.C:
				e.tick()
			}
		}
	}()
}

// StopAutoSync stops the background loop (call on shutdown).
func (e *Engine) StopAutoSync() {
	e.autoMu.Lock()
	if e.stopAuto != nil {
		close(e.stopAuto)
		e.stopAuto = nil
	}
	e.autoMu.Unlock()
}

// tick is one iteration of the background loop.
func (e *Engine) tick() {
	reachable := e.Reachable()
	e.autoMu.Lock()
	e.online = reachable
	models := e.syncModels
	e.autoMu.Unlock()

	if e.IsForceOffline() || !reachable {
		return
	}
	_, err := e.Sync(models)
	e.autoMu.Lock()
	if err != nil {
		e.lastErr = err.Error()
	} else {
		e.lastErr = ""
		e.lastSync = time.Now()
	}
	e.autoMu.Unlock()
}

// SyncNow forces an immediate Pull+Push (used when the user flips back online).
func (e *Engine) SyncNow() (*SyncResult, error) {
	if e.IsForceOffline() {
		return nil, fmt.Errorf("offline mode is on")
	}
	e.autoMu.Lock()
	models := e.syncModels
	e.autoMu.Unlock()
	res, err := e.Sync(models)
	e.autoMu.Lock()
	if err == nil {
		e.lastSync = time.Now()
		e.lastErr = ""
	}
	e.autoMu.Unlock()
	return res, err
}

// SyncStatus is the snapshot the dashboard/title-bar shows.
type SyncStatus struct {
	Reachable    bool   ` + "`" + `json:"reachable"` + "`" + `
	ForceOffline bool   ` + "`" + `json:"force_offline"` + "`" + `
	Pending      int64  ` + "`" + `json:"pending"` + "`" + `
	LastSync     string ` + "`" + `json:"last_sync,omitempty"` + "`" + `
	LastError    string ` + "`" + `json:"last_error,omitempty"` + "`" + `
	DeviceID     string ` + "`" + `json:"device_id,omitempty"` + "`" + `
	Tables       []string ` + "`" + `json:"tables,omitempty"` + "`" + `
}

// Status returns the current sync snapshot. Reachability comes from the last
// background tick (no extra network call on every poll).
func (e *Engine) Status() SyncStatus {
	pending, _ := e.PendingCount()
	e.autoMu.Lock()
	ls := ""
	if !e.lastSync.IsZero() {
		ls = e.lastSync.Format(time.RFC3339)
	}
	st := SyncStatus{
		Reachable:    e.online,
		ForceOffline: e.IsForceOffline(),
		Pending:      pending,
		LastSync:     ls,
		LastError:    e.lastErr,
		DeviceID:     e.deviceID,
		Tables:       e.syncModels,
	}
	e.autoMu.Unlock()
	return st
}

// Pull fetches every row in modelName updated after our last cursor
// and UPSERTs them into local records. Returns the number of rows seen.
func (e *Engine) Pull(modelName string) (int, error) {
	since := e.getCursor(modelName)
	url := fmt.Sprintf("%s/sync/pull?model=%s", e.APIURL, modelName)
	if since != "" {
		url += "&since=" + since
	}

	body, err := e.doGET(url)
	if err != nil {
		return 0, err
	}

	var resp struct {
		Data   []map[string]interface{} ` + "`" + `json:"data"` + "`" + `
		Cursor string                   ` + "`" + `json:"cursor"` + "`" + `
		Count  int                      ` + "`" + `json:"count"` + "`" + `
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return 0, fmt.Errorf("decoding pull response: %w", err)
	}

	for _, row := range resp.Data {
		id, _ := row["id"].(string)
		if id == "" {
			continue
		}
		version, _ := toInt(row["version"])
		deleted, _ := row["_deleted"].(bool)
		raw, err := json.Marshal(row)
		if err != nil {
			continue
		}
		updatedAt := time.Now().Unix()
		if u, ok := row["updated_at"].(string); ok {
			if t, err := time.Parse(time.RFC3339Nano, u); err == nil {
				updatedAt = t.Unix()
			}
		}
		rec := Record{
			Model:     modelName,
			ID:        id,
			Data:      raw,
			Version:   version,
			UpdatedAt: updatedAt,
			// A server-side delete arrives as a tombstone — mark the mirror
			// row deleted so reads (LocalList/LocalGet) hide it. We keep the
			// row rather than hard-deleting so a later re-create still upserts.
			Deleted: deleted,
		}
		// UPSERT: pulled state always wins over the cached version.
		if err := e.DB.Save(&rec).Error; err != nil {
			return 0, err
		}
	}
	if resp.Cursor != "" {
		e.setCursor(modelName, resp.Cursor)
	}
	return resp.Count, nil
}

// PushBatch is the JSON shape /api/sync/push expects.
type PushBatch struct {
	Changes []PushChange ` + "`" + `json:"changes"` + "`" + `
}

// PushChange mirrors the server's PushChange.
type PushChange struct {
	Op      string                 ` + "`" + `json:"op"` + "`" + `
	Model   string                 ` + "`" + `json:"model"` + "`" + `
	ID      string                 ` + "`" + `json:"id"` + "`" + `
	Version int                    ` + "`" + `json:"version"` + "`" + `
	Data    map[string]interface{} ` + "`" + `json:"data"` + "`" + `
}

// PushResult mirrors the server's PushResult.
type PushResult struct {
	OK            bool                   ` + "`" + `json:"ok"` + "`" + `
	Code          string                 ` + "`" + `json:"code,omitempty"` + "`" + `
	Message       string                 ` + "`" + `json:"message,omitempty"` + "`" + `
	ServerVersion int                    ` + "`" + `json:"server_version,omitempty"` + "`" + `
	ServerData    map[string]interface{} ` + "`" + `json:"server_data,omitempty"` + "`" + `
	NewVersion    int                    ` + "`" + `json:"new_version,omitempty"` + "`" + `
}

// Push drains the outbox into a single /api/sync/push call and applies
// each per-entry result. Successful entries are removed from the outbox
// and the local record is marked synced. Conflicts stay in the outbox
// with HasConflict=true and the server state attached so the UI can
// drive a merge dialog.
//
// Returns (pushed, conflicts, err).
func (e *Engine) Push() (int, int, error) {
	// Skip rows that already have an unresolved conflict — the user has
	// to resolve those via ResolveConflict before they're tried again.
	var entries []Outbox
	if err := e.DB.Where("has_conflict = 0").Order("created_at asc").Find(&entries).Error; err != nil {
		return 0, 0, err
	}
	if len(entries) == 0 {
		return 0, 0, nil
	}

	batch := PushBatch{Changes: make([]PushChange, 0, len(entries))}
	for _, en := range entries {
		var data map[string]interface{}
		if len(en.Data) > 0 {
			_ = json.Unmarshal(en.Data, &data)
		}
		batch.Changes = append(batch.Changes, PushChange{
			Op:      en.Op,
			Model:   en.Model,
			ID:      en.EntityID,
			Version: en.Version,
			Data:    data,
		})
	}

	body, err := e.doPOST(e.APIURL+"/sync/push", batch)
	if err != nil {
		return 0, 0, err
	}

	var resp struct {
		Results []PushResult ` + "`" + `json:"results"` + "`" + `
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return 0, 0, fmt.Errorf("decoding push response: %w", err)
	}
	if len(resp.Results) != len(entries) {
		return 0, 0, fmt.Errorf("push: server returned %d results for %d changes", len(resp.Results), len(entries))
	}

	pushed := 0
	conflicts := 0
	for i, r := range resp.Results {
		en := entries[i]
		switch {
		case r.OK:
			// Remove the outbox entry; mark the local record at the new version.
			if err := e.DB.Delete(&en).Error; err != nil {
				return pushed, conflicts, err
			}
			if en.Op == "delete" {
				e.DB.Where("model = ? AND id = ?", en.Model, en.EntityID).Delete(&Record{})
			} else {
				e.DB.Model(&Record{}).
					Where("model = ? AND id = ?", en.Model, en.EntityID).
					Updates(map[string]interface{}{"version": r.NewVersion})
			}
			pushed++

		case r.Code == "VERSION_CONFLICT":
			// Stash the server state on the outbox row for the UI.
			serverDataJSON, _ := json.Marshal(r.ServerData)
			e.DB.Model(&en).Updates(map[string]interface{}{
				"has_conflict":   true,
				"server_data":    serverDataJSON,
				"server_version": r.ServerVersion,
				"conflict_msg":   r.Message,
			})
			conflicts++

		default:
			// Other errors leave the entry alone; user can retry later.
			e.DB.Model(&en).Update("conflict_msg", fmt.Sprintf("%s: %s", r.Code, r.Message))
		}
	}
	return pushed, conflicts, nil
}

// ResolveConflict accepts the user's merge for one conflicted entity.
// mergedData becomes the new outbox payload; serverVersion is the
// version the user is overwriting (so the next push uses If-Match
// semantics correctly). Clears HasConflict so the entry is replayed
// on the next Push.
func (e *Engine) ResolveConflict(tableName, entityID string, mergedData map[string]interface{}, serverVersion int) error {
	dataJSON, err := json.Marshal(mergedData)
	if err != nil {
		return err
	}
	res := e.DB.Model(&Outbox{}).
		Where("model = ? AND entity_id = ?", tableName, entityID).
		Updates(map[string]interface{}{
			"data":           dataJSON,
			"version":        serverVersion,
			"has_conflict":   false,
			"server_data":    nil,
			"server_version": 0,
			"conflict_msg":   "",
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("no outbox entry for %s:%s", tableName, entityID)
	}
	// Also update the local record so reads see the merged state immediately.
	return e.DB.Model(&Record{}).
		Where("model = ? AND id = ?", tableName, entityID).
		Updates(map[string]interface{}{"data": dataJSON, "version": serverVersion}).Error
}

// PendingCount returns how many outbox entries are waiting (used to
// drive the title-bar badge).
func (e *Engine) PendingCount() (int64, error) {
	var n int64
	err := e.DB.Model(&Outbox{}).Count(&n).Error
	return n, err
}

// GetPendingChanges returns every outbox entry, oldest first, for the
// "review what's about to push" panel.
func (e *Engine) GetPendingChanges() ([]Outbox, error) {
	var entries []Outbox
	err := e.DB.Order("created_at asc").Find(&entries).Error
	return entries, err
}

func (e *Engine) doGET(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if err := e.attachAuth(req); err != nil {
		return nil, err
	}
	return e.do(req)
}

func (e *Engine) doPOST(url string, payload interface{}) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if err := e.attachAuth(req); err != nil {
		return nil, err
	}
	return e.do(req)
}

func (e *Engine) attachAuth(req *http.Request) error {
	if e.GetToken == nil {
		return nil
	}
	tok, err := e.GetToken()
	if err != nil {
		return err
	}
	if tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}
	return nil
}

func (e *Engine) do(req *http.Request) ([]byte, error) {
	resp, err := e.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("api %s: %d: %s", req.URL.Path, resp.StatusCode, string(body))
	}
	return body, nil
}

func (e *Engine) getCursor(model string) string {
	e.cursorsMu.RLock()
	c, ok := e.cursors[model]
	e.cursorsMu.RUnlock()
	if ok {
		return c
	}
	var row Cursor
	if err := e.DB.First(&row, "model = ?", model).Error; err == nil {
		e.cursorsMu.Lock()
		e.cursors[model] = row.Value
		e.cursorsMu.Unlock()
		return row.Value
	}
	return ""
}

func (e *Engine) setCursor(model, value string) {
	e.cursorsMu.Lock()
	e.cursors[model] = value
	e.cursorsMu.Unlock()
	e.DB.Save(&Cursor{Model: model, Value: value})
}

func toInt(v interface{}) (int, bool) {
	switch x := v.(type) {
	case float64:
		return int(x), true
	case int:
		return x, true
	case int64:
		return int(x), true
	}
	return 0, false
}
`
}

func desktopSyncOutboxGo() string {
	return `package sync

import (
	gosync "sync"
)

// Record is one row in the local mirror — the cache reads come from.
// Model holds the logical table name ("buildings"); the physical
// SQLite table is "sync_records" (set via TableName below).
type Record struct {
	Model     string ` + "`" + `gorm:"primaryKey;size:50"` + "`" + `
	ID        string ` + "`" + `gorm:"primaryKey;size:36"` + "`" + `
	Data      []byte ` + "`" + `gorm:"type:blob"` + "`" + ` // server JSON
	Version   int
	UpdatedAt int64 // unix seconds
	Deleted   bool
}

func (Record) TableName() string { return "sync_records" }

// Outbox is one pending local change. Squashed: at most one row per
// (table_name, entity_id), enforced by a UNIQUE index.
//
// Op is "create" / "update" / "delete". Version is the server version
// the local change is based on — sent as the optimistic-lock check on push.
//
// HasConflict + ServerData/ServerVersion populated when /api/sync/push
// returned VERSION_CONFLICT. The UI builds a per-field merge dialog
// from these and calls ResolveConflict() to clear them.
type Outbox struct {
	ID            int64  ` + "`" + `gorm:"primaryKey;autoIncrement"` + "`" + `
	Model         string ` + "`" + `gorm:"size:50;uniqueIndex:idx_outbox_entity"` + "`" + `
	EntityID      string ` + "`" + `gorm:"size:36;uniqueIndex:idx_outbox_entity"` + "`" + `
	Op            string ` + "`" + `gorm:"size:10"` + "`" + `
	Data          []byte ` + "`" + `gorm:"type:blob"` + "`" + `
	Version       int
	CreatedAt     int64
	HasConflict   bool
	ServerData    []byte ` + "`" + `gorm:"type:blob"` + "`" + `
	ServerVersion int
	ConflictMsg   string ` + "`" + `gorm:"size:500"` + "`" + `
}

func (Outbox) TableName() string { return "sync_outbox" }

// Cursor stores the last pull timestamp per model.
type Cursor struct {
	Model string ` + "`" + `gorm:"primaryKey;size:50"` + "`" + `
	Value string ` + "`" + `gorm:"size:50"` + "`" + ` // RFC3339Nano
}

func (Cursor) TableName() string { return "sync_cursors" }

var enqueueMu gosync.Mutex

// enqueue applies squash semantics:
//   - new "create" + no entry            → INSERT create
//   - new "update" + existing create     → UPDATE the create's data
//   - new "update" + existing update     → UPDATE data
//   - new "delete" + existing create     → DELETE outbox row (never made it to server)
//   - new "delete" + existing update     → flip to delete, drop data
//   - new "delete" + no entry            → INSERT delete
func enqueue(e *Engine, table, id, op string, data []byte, version int) error {
	enqueueMu.Lock()
	defer enqueueMu.Unlock()

	var existing Outbox
	err := e.DB.Where("model = ? AND entity_id = ?", table, id).First(&existing).Error
	switch {
	case err == nil:
		switch {
		case op == "delete" && existing.Op == "create":
			// Squashed away — both ends cancel.
			return e.DB.Delete(&existing).Error
		case op == "delete":
			return e.DB.Model(&existing).Updates(map[string]interface{}{
				"op":           "delete",
				"data":         nil,
				"has_conflict": false,
			}).Error
		default:
			// Update or upgrade: keep original op, refresh data.
			return e.DB.Model(&existing).Updates(map[string]interface{}{
				"data":         data,
				"has_conflict": false,
			}).Error
		}
	default:
		// New entry.
		entry := Outbox{
			Model:     table,
			EntityID:  id,
			Op:        op,
			Data:      data,
			Version:   version,
			CreatedAt: nowUnix(),
		}
		return e.DB.Create(&entry).Error
	}
}
`
}

func desktopSyncLocalGo() string {
	return `package sync

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LocalCreate persists data locally and queues a "create" entry in the
// outbox. id is required (UUID); pass uuid.New().String() if you don't
// have one yet. Reads via LocalGet/LocalList see the new row immediately.
func (e *Engine) LocalCreate(tableName, id string, data map[string]interface{}) error {
	if id == "" {
		id = uuid.New().String()
	}
	data["id"] = id
	if data["version"] == nil {
		data["version"] = 0
	}
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := e.DB.Save(&Record{
		Model:     tableName,
		ID:        id,
		Data:      raw,
		Version:   0,
		UpdatedAt: nowUnix(),
	}).Error; err != nil {
		return err
	}
	return enqueue(e, tableName, id, "create", raw, 0)
}

// LocalUpdate merges data into the existing local row and queues an
// "update" in the outbox. Errors if the row is missing.
func (e *Engine) LocalUpdate(tableName, id string, data map[string]interface{}) error {
	var rec Record
	if err := e.DB.First(&rec, "model = ? AND id = ?", tableName, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("local update: %s/%s not found", tableName, id)
		}
		return err
	}
	current := map[string]interface{}{}
	_ = json.Unmarshal(rec.Data, &current)
	for k, v := range data {
		current[k] = v
	}
	current["id"] = id
	raw, err := json.Marshal(current)
	if err != nil {
		return err
	}
	if err := e.DB.Model(&rec).Updates(map[string]interface{}{
		"data":       raw,
		"updated_at": nowUnix(),
	}).Error; err != nil {
		return err
	}
	return enqueue(e, tableName, id, "update", raw, rec.Version)
}

// LocalDelete soft-deletes from the local mirror and queues a "delete"
// in the outbox. The mirror row is removed too so reads stop seeing it.
func (e *Engine) LocalDelete(tableName, id string) error {
	var rec Record
	err := e.DB.First(&rec, "model = ? AND id = ?", tableName, id).Error
	knownVersion := 0
	if err == nil {
		knownVersion = rec.Version
		e.DB.Delete(&rec)
	} else if err != gorm.ErrRecordNotFound {
		return err
	}
	return enqueue(e, tableName, id, "delete", nil, knownVersion)
}

// LocalGet returns the cached record decoded from JSON, or nil if missing or
// tombstoned (deleted locally or via a pulled server delete).
func (e *Engine) LocalGet(tableName, id string) (map[string]interface{}, error) {
	var rec Record
	if err := e.DB.First(&rec, "model = ? AND id = ? AND deleted = ?", tableName, id, false).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	out := map[string]interface{}{}
	if err := json.Unmarshal(rec.Data, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// LocalList returns every cached record for tableName, decoded. The
// caller can filter / sort in memory; for an MVP we don't push
// SQL-shaped queries through to SQLite.
func (e *Engine) LocalList(tableName string) ([]map[string]interface{}, error) {
	var rows []Record
	if err := e.DB.Where("model = ? AND deleted = ?", tableName, false).Order("updated_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(rows))
	for _, r := range rows {
		m := map[string]interface{}{}
		if err := json.Unmarshal(r.Data, &m); err == nil {
			out = append(out, m)
		}
	}
	return out, nil
}

func nowUnix() int64 { return time.Now().Unix() }
`
}

// ═══════════════════════════════════════════════════════════════════
// Frontend — sync client + UI templates
// ═══════════════════════════════════════════════════════════════════

func desktopClientSyncClientTS() string {
	return `// Typed wrappers around the Wails sync bindings exposed by App. Use
// these instead of touching window.go directly so the call sites are
// type-checked.

export interface OutboxEntry {
  ID: number;
  Model: string;
  EntityID: string;
  Op: "create" | "update" | "delete";
  Data: string | null;            // JSON-encoded byte slice from Go
  Version: number;
  CreatedAt: number;
  HasConflict: boolean;
  ServerData: string | null;       // JSON-encoded server state on conflict
  ServerVersion: number;
  ConflictMsg: string;
}

export interface SyncResult {
  pushed: number;
  pulled: number;
  conflicts: number;
  errors?: string[];
  started_at: string;
  finished_at: string;
}

const isWails = typeof window !== "undefined" && !!window.go?.main?.App;

function notWailsError(): never {
  throw new Error("Sync requires the Wails desktop runtime (run via 'wails dev').");
}

export async function localCreate(
  table: string,
  id: string,
  data: Record<string, unknown>,
): Promise<void> {
  if (!isWails) notWailsError();
  return window.go!.main.App.LocalCreate(table, id, data);
}

export async function localUpdate(
  table: string,
  id: string,
  data: Record<string, unknown>,
): Promise<void> {
  if (!isWails) notWailsError();
  return window.go!.main.App.LocalUpdate(table, id, data);
}

export async function localDelete(table: string, id: string): Promise<void> {
  if (!isWails) notWailsError();
  return window.go!.main.App.LocalDelete(table, id);
}

export async function localGet(
  table: string,
  id: string,
): Promise<Record<string, unknown> | null> {
  if (!isWails) notWailsError();
  return window.go!.main.App.LocalGet(table, id);
}

export async function localList(
  table: string,
): Promise<Record<string, unknown>[]> {
  if (!isWails) notWailsError();
  return window.go!.main.App.LocalList(table);
}

export async function sync(tables: string[]): Promise<SyncResult> {
  if (!isWails) notWailsError();
  return window.go!.main.App.Sync(tables);
}

export async function pendingCount(): Promise<number> {
  if (!isWails) return 0;
  return window.go!.main.App.PendingCount();
}

export async function getPendingChanges(): Promise<OutboxEntry[]> {
  if (!isWails) return [];
  return window.go!.main.App.GetPendingChanges();
}

export async function resolveConflict(
  table: string,
  entityID: string,
  mergedData: Record<string, unknown>,
  serverVersion: number,
): Promise<void> {
  if (!isWails) notWailsError();
  return window.go!.main.App.ResolveConflict(table, entityID, mergedData, serverVersion);
}

// Helper: parse the JSON-encoded byte slices Go sends through Wails.
// Wails serializes []byte as a base64 string OR a stringified JSON
// depending on version; this handles both.
export function parseGoBytes(s: string | null): Record<string, unknown> | null {
  if (!s) return null;
  try {
    return JSON.parse(s);
  } catch {
    try {
      return JSON.parse(atob(s));
    } catch {
      return null;
    }
  }
}
`
}

func desktopClientUseSync() string {
	return `import { useEffect, useState } from "react";
import {
  pendingCount,
  getPendingChanges,
  sync,
  resolveConflict,
  type OutboxEntry,
  type SyncResult,
} from "@/lib/sync-client";

// useSync wraps the Wails sync bindings with React state. We avoid React
// Query here because the source of truth is the Wails Go process — not
// the API — and React Query's network-defaults don't fit the
// IPC-not-HTTP nature of the bindings.

// usePendingCount polls the engine every POLL_MS for an updated count.
// Cheap because PendingCount is a single SELECT COUNT(*) on a small
// table. Wires to the title-bar Sync button badge.
const POLL_MS = 2_000;

export function usePendingCount() {
  const [count, setCount] = useState(0);
  useEffect(() => {
    let cancelled = false;
    const tick = async () => {
      try {
        const n = await pendingCount();
        if (!cancelled) setCount(n);
      } catch {
        // ignore — Wails may not be initialized yet on first render
      }
    };
    tick();
    const id = window.setInterval(tick, POLL_MS);
    return () => {
      cancelled = true;
      window.clearInterval(id);
    };
  }, []);
  return count;
}

// usePendingChanges fetches the full outbox once and exposes a refresh
// fn. Used by the pending-changes panel.
export function usePendingChanges() {
  const [entries, setEntries] = useState<OutboxEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const refresh = async () => {
    setLoading(true);
    try {
      const rows = await getPendingChanges();
      setEntries(rows || []);
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    refresh();
  }, []);

  return { entries, loading, error, refresh };
}

// useSyncMutation runs a Sync against the supplied tables and tracks
// running / result state. The caller passes the list of tables the app
// cares about (e.g. ["buildings", "tenants"]).
export function useSyncMutation(tables: string[]) {
  const [running, setRunning] = useState(false);
  const [result, setResult] = useState<SyncResult | null>(null);
  const [error, setError] = useState<string | null>(null);

  const run = async () => {
    setRunning(true);
    setError(null);
    try {
      const r = await sync(tables);
      setResult(r);
      return r;
    } catch (e) {
      const msg = e instanceof Error ? e.message : String(e);
      setError(msg);
      throw e;
    } finally {
      setRunning(false);
    }
  };

  return { run, running, result, error };
}

// useResolveConflict applies a merge for a single conflicted entry.
// The hook itself is a thin wrapper; consumers should refresh the
// pending list after a successful resolve.
export function useResolveConflict() {
  const [resolving, setResolving] = useState(false);
  const resolve = async (
    table: string,
    entityID: string,
    mergedData: Record<string, unknown>,
    serverVersion: number,
  ) => {
    setResolving(true);
    try {
      await resolveConflict(table, entityID, mergedData, serverVersion);
    } finally {
      setResolving(false);
    }
  };
  return { resolve, resolving };
}
`
}

func desktopClientSyncButton() string {
	return `import { useState } from "react";
import { RefreshCw, AlertCircle } from "lucide-react";
import { cn } from "@/lib/utils";
import { usePendingCount } from "@/hooks/use-sync";
import { PendingChangesPanel } from "@/components/pending-changes";

// SyncButton sits in the title-bar next to the connection indicator.
// Shows a pending-count badge; click opens the PendingChangesPanel
// (drawer over the right edge) where the user can review and push.
//
// tables is the list of model names this app cares about — passed down
// to the panel which forwards it to Sync(). For most apps this is a
// const at the layout level: ["buildings", "tenants", "leases", ...].
export function SyncButton({ tables }: { tables: string[] }) {
  const [open, setOpen] = useState(false);
  const count = usePendingCount();
  const hasPending = count > 0;

  return (
    <>
      <button
        type="button"
        onClick={() => setOpen(true)}
        className={cn(
          "no-drag flex h-titlebar items-center gap-1.5 px-3 transition-colors",
          hasPending
            ? "text-warning hover:text-warning"
            : "text-foreground-muted hover:text-foreground",
        )}
        title={hasPending ? count + " pending change" + (count === 1 ? "" : "s") + " - click to sync" : "All synced"}
      >
        {hasPending ? <AlertCircle className="h-3.5 w-3.5" /> : <RefreshCw className="h-3.5 w-3.5" />}
        {hasPending && (
          <span className="rounded-full bg-warning/15 px-1.5 py-0.5 text-[10px] font-semibold leading-none">
            {count}
          </span>
        )}
      </button>
      {open && <PendingChangesPanel tables={tables} onClose={() => setOpen(false)} />}
    </>
  );
}
`
}

func desktopClientPendingChanges() string {
	return `import { useState } from "react";
import { X, RefreshCw, AlertTriangle, Plus, Pencil, Trash2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { parseGoBytes, type OutboxEntry } from "@/lib/sync-client";
import { usePendingChanges, useSyncMutation } from "@/hooks/use-sync";
import { ConflictDialog } from "@/components/conflict-dialog";

// PendingChangesPanel is a right-edge drawer listing every outbox
// entry. The user sees what's about to push, then clicks "Sync now".
// On success, conflict entries surface a per-entry "Resolve" button
// that opens the ConflictDialog.
//
// The drawer fixed-positions itself over the rest of the app — clean
// and predictable on desktop. No animation, no tricks.
export function PendingChangesPanel({
  tables,
  onClose,
}: {
  tables: string[];
  onClose: () => void;
}) {
  const { entries, refresh } = usePendingChanges();
  const { run, running } = useSyncMutation(tables);
  const [conflictEntry, setConflictEntry] = useState<OutboxEntry | null>(null);

  const handleSync = async () => {
    try {
      await run();
      await refresh();
    } catch {
      // useSyncMutation surfaces the error; refresh anyway so the user
      // sees what's still pending.
      await refresh();
    }
  };

  const conflicts = entries.filter((e) => e.HasConflict);
  const clean = entries.filter((e) => !e.HasConflict);

  return (
    <>
      <div className="fixed inset-0 z-40 bg-black/40" onClick={onClose} />
      <aside className="fixed right-0 top-0 z-50 h-full w-[420px] bg-surface border-l border-border flex flex-col">
        <header className="flex items-center justify-between px-4 py-3 border-b border-border">
          <div>
            <h2 className="text-[14px] font-semibold text-foreground">Pending changes</h2>
            <p className="text-[12px] text-foreground-muted mt-0.5">
              {entries.length === 0 ? "Nothing to sync" : entries.length + " in outbox"}
              {conflicts.length > 0 && ", " + conflicts.length + " in conflict"}
            </p>
          </div>
          <button onClick={onClose} className="text-foreground-muted hover:text-foreground">
            <X className="h-4 w-4" />
          </button>
        </header>

        <div className="flex-1 overflow-y-auto">
          {conflicts.length > 0 && (
            <section className="border-b border-border-subtle">
              <h3 className="px-4 pt-3 pb-1 text-[11px] font-semibold uppercase tracking-wider text-warning">
                Needs review ({conflicts.length})
              </h3>
              {conflicts.map((e) => (
                <PendingRow key={e.ID} entry={e} onResolve={() => setConflictEntry(e)} />
              ))}
            </section>
          )}
          {clean.length > 0 && (
            <section>
              <h3 className="px-4 pt-3 pb-1 text-[11px] font-semibold uppercase tracking-wider text-foreground-muted">
                Ready to push ({clean.length})
              </h3>
              {clean.map((e) => (
                <PendingRow key={e.ID} entry={e} />
              ))}
            </section>
          )}
          {entries.length === 0 && (
            <div className="p-8 text-center text-[13px] text-foreground-muted">
              No pending changes. Everything is synced.
            </div>
          )}
        </div>

        <footer className="border-t border-border px-4 py-3 flex items-center justify-between">
          <button
            onClick={refresh}
            className="text-[12px] text-foreground-muted hover:text-foreground"
          >
            Refresh
          </button>
          <button
            onClick={handleSync}
            disabled={running || entries.length === 0}
            className={cn(
              "inline-flex items-center gap-2 h-9 px-3.5 rounded-lg text-[13px] font-medium transition-colors",
              "bg-accent text-white hover:bg-accent-hover",
              "disabled:opacity-50 disabled:cursor-not-allowed",
            )}
          >
            <RefreshCw className={cn("h-3.5 w-3.5", running && "animate-spin")} />
            {running ? "Syncing..." : "Sync now"}
          </button>
        </footer>
      </aside>

      {conflictEntry && (
        <ConflictDialog
          entry={conflictEntry}
          onClose={() => setConflictEntry(null)}
          onResolved={() => {
            setConflictEntry(null);
            refresh();
          }}
        />
      )}
    </>
  );
}

function PendingRow({
  entry,
  onResolve,
}: {
  entry: OutboxEntry;
  onResolve?: () => void;
}) {
  const data = parseGoBytes(entry.Data);
  const Icon = entry.Op === "create" ? Plus : entry.Op === "delete" ? Trash2 : Pencil;
  const opColor =
    entry.Op === "create"
      ? "text-success"
      : entry.Op === "delete"
        ? "text-danger"
        : "text-info";

  return (
    <div className="px-4 py-2.5 border-b border-border-subtle hover:bg-surface-hover">
      <div className="flex items-start gap-3">
        <div className={cn("mt-0.5", opColor)}>
          <Icon className="h-3.5 w-3.5" />
        </div>
        <div className="flex-1 min-w-0">
          <div className="text-[13px] font-medium text-foreground">
            {entry.Op} {entry.Model}
          </div>
          <div className="text-[11.5px] text-foreground-muted truncate font-mono">
            {entry.EntityID.slice(0, 8)}{data && data.name ? " · " + String(data.name) : ""}
          </div>
          {entry.HasConflict && (
            <div className="mt-1 flex items-center gap-1 text-[11px] text-warning">
              <AlertTriangle className="h-3 w-3" />
              {entry.ConflictMsg || "Server has a newer version"}
            </div>
          )}
        </div>
        {entry.HasConflict && onResolve && (
          <button
            onClick={onResolve}
            className="shrink-0 h-7 px-2.5 rounded-md bg-warning/10 hover:bg-warning/20 text-warning text-[11.5px] font-medium"
          >
            Resolve
          </button>
        )}
      </div>
    </div>
  );
}
`
}

func desktopClientConflictDialog() string {
	return `import { useMemo, useState } from "react";
import { X } from "lucide-react";
import { cn } from "@/lib/utils";
import { parseGoBytes, type OutboxEntry } from "@/lib/sync-client";
import { useResolveConflict } from "@/hooks/use-sync";

// ConflictDialog is the field-level merge UI. When push hits a
// VERSION_CONFLICT, the outbox entry stores both the local payload
// (Data) and the server's current state (ServerData). We diff them
// and let the user pick a side per field.
//
// "Local" = the value the user typed offline.
// "Server" = the value some other client (or this user on another
// device) wrote since the last sync.
//
// On Resolve we build the merged record by walking each field's choice
// and call ResolveConflict() with the result + the new ServerVersion.
// The next Sync replays the entry with that version as the optimistic-
// lock check, so it'll succeed cleanly.
export function ConflictDialog({
  entry,
  onClose,
  onResolved,
}: {
  entry: OutboxEntry;
  onClose: () => void;
  onResolved: () => void;
}) {
  const local = useMemo(() => parseGoBytes(entry.Data) ?? {}, [entry.Data]);
  const server = useMemo(() => parseGoBytes(entry.ServerData) ?? {}, [entry.ServerData]);
  const { resolve, resolving } = useResolveConflict();

  // Hidden / system fields we don't want to expose in the merge UI.
  const HIDE_FIELDS = new Set([
    "id",
    "version",
    "created_at",
    "updated_at",
    "deleted_at",
  ]);

  const allFields = useMemo(() => {
    const set = new Set<string>();
    Object.keys(local).forEach((k) => set.add(k));
    Object.keys(server).forEach((k) => set.add(k));
    HIDE_FIELDS.forEach((k) => set.delete(k));
    return Array.from(set).sort();
  }, [local, server]);

  // Per-field choice: "local" or "server". Default to "local" because
  // the user just made those edits — they probably want to keep them.
  const [choices, setChoices] = useState<Record<string, "local" | "server">>(
    () => {
      const initial: Record<string, "local" | "server"> = {};
      for (const f of allFields) initial[f] = "local";
      return initial;
    },
  );

  const handleResolve = async () => {
    const merged: Record<string, unknown> = { ...server }; // start from server state
    for (const field of allFields) {
      merged[field] = choices[field] === "local" ? local[field] : server[field];
    }
    // Preserve the ID so the wire payload is consistent.
    merged.id = entry.EntityID;
    await resolve(entry.Model, entry.EntityID, merged, entry.ServerVersion);
    onResolved();
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-6">
      <div className="w-full max-w-3xl max-h-[85vh] bg-surface border border-border rounded-xl flex flex-col">
        <header className="flex items-start justify-between px-5 py-4 border-b border-border">
          <div>
            <h2 className="text-[15px] font-semibold text-foreground">
              Resolve conflict: {entry.Op} {entry.Model}
            </h2>
            <p className="text-[12.5px] text-foreground-muted mt-1">
              The server has a newer version (v{entry.ServerVersion}) than what you edited locally.
              Pick which side wins for each field.
            </p>
          </div>
          <button onClick={onClose} className="text-foreground-muted hover:text-foreground">
            <X className="h-4 w-4" />
          </button>
        </header>

        <div className="flex-1 overflow-y-auto px-5 py-4 space-y-2">
          <div className="grid grid-cols-[1fr,160px,160px] gap-3 text-[11px] uppercase tracking-wider text-foreground-muted pb-2 border-b border-border-subtle">
            <div>Field</div>
            <div>Local</div>
            <div>Server (v{entry.ServerVersion})</div>
          </div>

          {allFields.length === 0 && (
            <div className="py-6 text-center text-[13px] text-foreground-muted">
              No editable fields differ — nothing to merge.
            </div>
          )}

          {allFields.map((field) => {
            const localVal = local[field];
            const serverVal = server[field];
            const sameValue =
              JSON.stringify(localVal) === JSON.stringify(serverVal);
            const choice = choices[field];

            return (
              <div
                key={field}
                className={cn(
                  "grid grid-cols-[1fr,160px,160px] gap-3 py-2 items-center",
                  sameValue ? "opacity-50" : "",
                )}
              >
                <div className="text-[13px] font-medium text-foreground-secondary">
                  {field}
                  {sameValue && (
                    <span className="ml-2 text-[10px] text-foreground-muted">unchanged</span>
                  )}
                </div>
                <FieldCell
                  value={localVal}
                  selected={choice === "local"}
                  disabled={sameValue}
                  onClick={() => !sameValue && setChoices({ ...choices, [field]: "local" })}
                />
                <FieldCell
                  value={serverVal}
                  selected={choice === "server"}
                  disabled={sameValue}
                  onClick={() => !sameValue && setChoices({ ...choices, [field]: "server" })}
                />
              </div>
            );
          })}
        </div>

        <footer className="flex items-center justify-end gap-2 px-5 py-3 border-t border-border">
          <button
            onClick={onClose}
            className="h-9 px-3.5 rounded-lg border border-border bg-surface text-[13px] font-medium text-foreground-secondary hover:bg-surface-hover"
          >
            Cancel
          </button>
          <button
            onClick={handleResolve}
            disabled={resolving}
            className="h-9 px-3.5 rounded-lg bg-accent text-white text-[13px] font-medium hover:bg-accent-hover disabled:opacity-60"
          >
            {resolving ? "Resolving..." : "Apply merge"}
          </button>
        </footer>
      </div>
    </div>
  );
}

function FieldCell({
  value,
  selected,
  disabled,
  onClick,
}: {
  value: unknown;
  selected: boolean;
  disabled: boolean;
  onClick: () => void;
}) {
  const display =
    value === null || value === undefined
      ? "(empty)"
      : typeof value === "object"
        ? JSON.stringify(value)
        : String(value);

  return (
    <button
      type="button"
      onClick={onClick}
      disabled={disabled}
      className={cn(
        "px-2.5 py-1.5 rounded-md text-[12.5px] text-left truncate border transition-colors",
        selected
          ? "border-accent bg-accent/10 text-foreground"
          : "border-border bg-surface-2 text-foreground-secondary hover:bg-surface-hover",
        disabled && "cursor-not-allowed",
      )}
      title={display}
    >
      {display}
    </button>
  );
}
`
}
