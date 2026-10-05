package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"commerce/apps/api/internal/authz"
	"commerce/apps/api/internal/respond"
	"commerce/apps/api/internal/services"
	"commerce/apps/api/internal/sync"
)

// SyncHandler implements /api/sync/push and /api/sync/pull. The push
// endpoint applies a batch of client changes with per-change version
// checking; the pull endpoint streams server-side updates since a
// caller-supplied cursor.
type SyncHandler struct {
	DB       *gorm.DB
	Registry *sync.Registry
}

// NewSyncHandler wires the handler to the database + model registry.
func NewSyncHandler(db *gorm.DB, reg *sync.Registry) *SyncHandler {
	return &SyncHandler{DB: db, Registry: reg}
}

// PushChange is one entry in a /api/sync/push batch. Op is one of
// "create" / "update" / "delete". Version is the version the client
// believes the server has — mismatches surface as VERSION_CONFLICT.
type PushChange struct {
	Op      string                 `json:"op"`
	Model   string                 `json:"model"`
	ID      string                 `json:"id"`
	Version int                    `json:"version"`
	Data    map[string]interface{} `json:"data"`
}

// PushResult is the per-change result returned in the same order as
// the input batch. On VERSION_CONFLICT, ServerVersion + ServerData
// carry the current server state so the client can build a merge UI.
type PushResult struct {
	OK            bool        `json:"ok"`
	Code          string      `json:"code,omitempty"`
	Message       string      `json:"message,omitempty"`
	ServerVersion int         `json:"server_version,omitempty"`
	ServerData    interface{} `json:"server_data,omitempty"`
	NewVersion    int         `json:"new_version,omitempty"`
}

// Changes made while offline.
type SyncPushRequest struct {
	Changes []PushChange `json:"changes"`
}

// MaxPushChanges is the most changes one push may carry. The sync clients Grit
// ships send their outbox in pushes of this size. An unbounded push of 20,000
// changes was one request running 60,000 queries, which the client gave up on
// and sent again.
const MaxPushChanges = 500

// Push handles POST /api/sync/push. Each change is applied
// independently: one conflict does not abort the rest of the batch.
func (h *SyncHandler) Push(c *gin.Context) {
	var req SyncPushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeInvalidBody, err.Error())
		return
	}
	if len(req.Changes) > MaxPushChanges {
		respond.Fail(c, respond.CodeTooManyChanges, fmt.Sprintf("a push may carry at most %d changes; send the rest in another push", MaxPushChanges))
		return
	}

	rows := h.loadCurrent(c, req.Changes)
	results := make([]PushResult, len(req.Changes))
	for i, ch := range req.Changes {
		results[i] = h.applyChange(c, ch, rows)
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}

// rows is the service every query below goes through. Built per request because
// it is one field and a struct literal, and because building it here means the
// wiring in routes.go did not have to change.
//
// It takes destinations rather than models: this handler syncs whatever the
// project registered, so the type is the registry's and the reflection is this
// file's. What the service owns is the keyset pagination a pull depends on and
// what a save may touch. See internal/services/sync.go.
func (h *SyncHandler) rows() *services.SyncService {
	return &services.SyncService{DB: h.DB}
}

// currentRows holds the rows a push updates or deletes, keyed by model and
// then id, read with one query per model rather than one per change.
type currentRows map[string]map[string]interface{}

// loadCurrent reads the rows the push's updates and deletes name. A row it did
// not load is read by its change, so a failed preload costs queries, not
// correctness.
func (h *SyncHandler) loadCurrent(c *gin.Context, changes []PushChange) currentRows {
	ids := map[string][]string{}
	for _, ch := range changes {
		if ch.Op == "update" || ch.Op == "delete" {
			ids[ch.Model] = append(ids[ch.Model], ch.ID)
		}
	}
	rows := currentRows{}
	for model, list := range ids {
		proto, err := h.Registry.New(model)
		if err != nil {
			continue
		}
		found := reflect.New(reflect.SliceOf(reflect.TypeOf(proto)))
		if err := h.rows().ByIDs(c.Request.Context(), found.Interface(), list); err != nil {
			log.Printf("sync push: preloading %s | id=%s: %v", model, c.GetString("request_id"), err)
			continue
		}
		byID := make(map[string]interface{}, found.Elem().Len())
		for i := 0; i < found.Elem().Len(); i++ {
			row := found.Elem().Index(i).Interface()
			if id := reflect.ValueOf(row).Elem().FieldByName("ID"); id.IsValid() {
				byID[fmt.Sprint(id.Interface())] = row
			}
		}
		rows[model] = byID
	}
	return rows
}

// current returns the row a change applies to and takes it out of rows, so a
// second change to the same row in one push reads what the first one wrote.
func (h *SyncHandler) current(c *gin.Context, rows currentRows, ch PushChange, proto interface{}) (interface{}, error) {
	if row, ok := rows[ch.Model][ch.ID]; ok {
		delete(rows[ch.Model], ch.ID)
		return row, nil
	}
	err := h.rows().ByID(c.Request.Context(), proto, ch.ID)
	return proto, err
}

// syncIdentifier picks a human-friendly label for the semantic activity feed
// from a change payload (name / title / slug / email), falling back to the id.
func syncIdentifier(data map[string]interface{}, id string) string {
	for _, k := range []string{"name", "title", "slug", "email"} {
		if v, ok := data[k]; ok {
			if s, ok := v.(string); ok && s != "" {
				return s
			}
		}
	}
	return id
}

// syncFault logs why the server could not apply a change and returns what the
// client is told instead. A database error names tables, columns and
// constraints, and the push result goes back to the device that sent it.
func syncFault(c *gin.Context, ch PushChange, err error, message string) string {
	log.Printf("sync push %s %s %s | id=%s: %v", ch.Op, ch.Model, ch.ID, c.GetString("request_id"), err)
	return message
}

func (h *SyncHandler) applyChange(c *gin.Context, ch PushChange, rows currentRows) PushResult {
	proto, err := h.Registry.New(ch.Model)
	if err != nil {
		return PushResult{OK: false, Code: "UNKNOWN_MODEL", Message: err.Error()}
	}
	policy := h.Registry.PolicyFor(ch.Model)
	// Enforced before the payload reaches a decoder: a field declared
	// local_only is a promise it never leaves the device, and a promise kept
	// only by well-behaved clients is not one.
	ch.Data = policy.StripLocalOnly(ch.Data)
	// Fields the server owns: the row's id is the change's id, the version is the
	// server's, and the timestamps are never the client's to set.
	ch.Data = stripServerFields(ch.Data)
	// The Go struct name (e.g. "Category") is the nicest entity label for the
	// activity feed — offline edits should read the same as online ones.
	entityType := reflect.TypeOf(proto).Elem().Name()

	switch ch.Op {
	case "create":
		// Decode the client payload into a fresh model struct and insert.
		// We trust the client-supplied ID (UUID) so the local outbox can
		// keep referring to the same row after the server insert.
		obj := proto
		if err := decodeInto(obj, ch.Data); err != nil {
			return PushResult{OK: false, Code: "DECODE_ERROR", Message: err.Error()}
		}
		setField(obj, "ID", ch.ID)
		if !syncMayWrite(c, ch.Model, "create", obj) {
			return PushResult{OK: false, Code: "FORBIDDEN", Message: "you may not create " + ch.Model}
		}
		if err := h.rows().Create(c.Request.Context(), obj); err != nil {
			return PushResult{OK: false, Code: "CREATE_FAILED", Message: syncFault(c, ch, err, "the server could not create the row")}
		}
		// Mirror the online handler: emit a semantic activity row so offline
		// creates surface in /system/activity, not just the raw audit log.
		services.LogCreate(h.DB, c, entityType, syncIdentifier(ch.Data, ch.ID), ch.ID, "")
		return PushResult{OK: true, NewVersion: 1}

	case "update":
		// Versioned update: load current row, compare versions, update if match.
		current, err := h.current(c, rows, ch, proto)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return PushResult{OK: false, Code: "NOT_FOUND", Message: "row was deleted on the server"}
			}
			return PushResult{OK: false, Code: "INTERNAL_ERROR", Message: syncFault(c, ch, err, "the server could not read the row")}
		}
		// Not found rather than forbidden: a row the caller may not touch is one
		// whose existence it does not get to learn.
		if !syncMayWrite(c, ch.Model, "edit", current) {
			return PushResult{OK: false, Code: "NOT_FOUND", Message: "row was deleted on the server"}
		}
		serverVersion := getIntField(current, "Version")
		if serverVersion != ch.Version {
			// The versions disagree. What happens next is the resource's
			// declared policy, decided here rather than in the client, because
			// a rule an old build can ignore is not a rule.
			switch policy.Conflict {
			case sync.ConflictServerWins:
				// The client's change is dropped and it is told to take the
				// server row. Reported distinctly from a conflict so the
				// client can apply it without asking anyone.
				return PushResult{
					OK:            false,
					Code:          "SERVER_WINS",
					Message:       fmt.Sprintf("server v%d kept over client v%d", serverVersion, ch.Version),
					ServerVersion: serverVersion,
					ServerData:    current,
				}
			case sync.ConflictClientWins:
				// Fall through and overwrite. The version check was protecting
				// nothing for this resource, and the author said so.
			default:
				return PushResult{
					OK:            false,
					Code:          "VERSION_CONFLICT",
					Message:       fmt.Sprintf("client had v%d, server has v%d", ch.Version, serverVersion),
					ServerVersion: serverVersion,
					ServerData:    current,
				}
			}
		}
		// Apply the update.
		//
		// Decode the client payload into a fresh, typed model struct rather than
		// calling .Updates(ch.Data) directly. A raw map[string]interface{} hands
		// nested values (a FileRef image, a FileRefs slice, a belongs-to relation
		// object) straight to the DB driver, which cannot encode a Go map into a
		// json column ("cannot find encode plan for OID 0") — the update fails and
		// the offline outbox entry gets stuck forever. Decoding first routes those
		// fields through their driver.Valuer implementations, exactly like create.
		obj := current
		if err := decodeInto(obj, ch.Data); err != nil {
			return PushResult{OK: false, Code: "DECODE_ERROR", Message: err.Error()}
		}
		setField(obj, "ID", ch.ID)
		// The payload is merged onto the loaded row, so check again: an owner may
		// not hand a row to somebody else by rewriting its owner field.
		if !syncMayWrite(c, ch.Model, "edit", obj) {
			return PushResult{OK: false, Code: "FORBIDDEN", Message: "you may not change who owns this row"}
		}
		// Seed Version with the server's value so the BeforeUpdate hook bumps it to
		// serverVersion+1 regardless of what the client sent in the payload.
		setIntField(obj, "Version", serverVersion)
		// Save writes every column (so cleared/zeroed fields persist) and runs the
		// BeforeUpdate hook. Omit associations so the nested relation object is not
		// upserted, and CreatedAt so the client can't rewind the original timestamp.
		if err := h.rows().Save(c.Request.Context(), obj); err != nil {
			return PushResult{OK: false, Code: "UPDATE_FAILED", Message: syncFault(c, ch, err, "the server could not update the row")}
		}
		newVersion := getIntField(obj, "Version")
		services.LogUpdate(h.DB, c, entityType, syncIdentifier(ch.Data, ch.ID), ch.ID, services.DiffSummary(ch.Data))
		return PushResult{OK: true, NewVersion: newVersion}

	case "delete":
		current, err := h.current(c, rows, ch, proto)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Already gone — treat as success so the outbox can clear.
				return PushResult{OK: true}
			}
			return PushResult{OK: false, Code: "INTERNAL_ERROR", Message: syncFault(c, ch, err, "the server could not read the row")}
		}
		if !syncMayWrite(c, ch.Model, "delete", current) {
			return PushResult{OK: false, Code: "NOT_FOUND", Message: "row was deleted on the server"}
		}
		serverVersion := getIntField(current, "Version")
		if ch.Version != 0 && serverVersion != ch.Version {
			return PushResult{
				OK:            false,
				Code:          "VERSION_CONFLICT",
				Message:       "row was modified after the client's last sync",
				ServerVersion: serverVersion,
				ServerData:    current,
			}
		}
		if err := h.rows().Delete(c.Request.Context(), current, ch.ID); err != nil {
			return PushResult{OK: false, Code: "DELETE_FAILED", Message: syncFault(c, ch, err, "the server could not delete the row")}
		}
		services.LogDelete(h.DB, c, entityType, ch.ID, ch.ID)
		return PushResult{OK: true}

	default:
		return PushResult{OK: false, Code: "INVALID_OP", Message: "op must be create, update, or delete"}
	}
}

// SyncPolicyResponse is what GET /api/sync/policy answers with.
type SyncPolicyResponse struct {
	Models map[string]sync.Policy `json:"models"`
}

// Policy handles GET /api/sync/policy.
//
// Clients configure themselves from this rather than from a generated copy.
// A copy is a second thing to keep in sync, and an offline client running
// last month's build against this month's conflict rules is exactly the
// silent failure this whole feature exists to prevent.
func (h *SyncHandler) Policy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": SyncPolicyResponse{Models: h.Registry.Policies()}})
}

// Pull handles GET /api/sync/pull?since=<cursor>&model=<table>. Returns the
// rows of the table changed after the cursor, oldest first, deletes included
// as tombstones. The client sends the response's cursor as the next ?since.
func (h *SyncHandler) Pull(c *gin.Context) {
	model := c.Query("model")
	if model == "" {
		respond.Fail(c, respond.CodeMissingModel, "?model is required")
		return
	}
	sinceStr := c.DefaultQuery("since", "")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "500"))
	if limit < 1 || limit > 5000 {
		limit = 500
	}

	proto, err := h.Registry.New(model)
	if err != nil {
		respond.Fail(c, respond.CodeUnknownModel, err.Error())
		return
	}
	policy := h.Registry.PolicyFor(model)
	if policy.Mode == sync.ModeOnlineOnly {
		// Answering with an empty page would look like "nothing has changed",
		// and the client would mirror an empty table forever.
		respond.Fail(c, respond.CodeNotSyncable, model+" is registered online_only and is not mirrored")
		return
	}

	// Who reads what. A holder of <model>.view reads the table; anybody else
	// reads only the rows they own, and only from a table whose rows have owners.
	// Before this, any signed-in account pulled every row of every synced table.
	seeAll := syncGranted(c, model, "view")
	if _, owned := proto.(ownedRow); !seeAll && !owned {
		respond.Fail(c, respond.CodeForbidden, "you do not have permission to read "+model)
		return
	}

	// Build a slice of the right type via reflection.
	sliceType := reflect.SliceOf(reflect.TypeOf(proto).Elem())
	results := reflect.New(sliceType)

	// Where this page starts. The cursor is the request's, so it is parsed here;
	// how a page is read from it is the service's, and the reasons it is keyset
	// rather than ordered by a computed column are written there.
	var since time.Time
	var afterID string
	if sinceStr != "" {
		parsedSince, parsedID, err := parseSyncCursor(sinceStr)
		if err != nil {
			respond.Fail(c, respond.CodeInvalidSince, err.Error())
			return
		}
		since, afterID = parsedSince, parsedID
	}
	if err := h.rows().Page(c.Request.Context(), proto, results.Interface(), since, afterID, limit); err != nil {
		respond.ServerError(c, "INTERNAL_ERROR", err, "Internal server error")
		return
	}

	rs := results.Elem()
	rows := make([]map[string]interface{}, 0, rs.Len())
	cursor := sinceStr
	for i := 0; i < rs.Len(); i++ {
		item := rs.Index(i).Addr().Interface()
		// The cursor passes every row read, including ones this caller may not
		// see, so the next page starts after them.
		if next, ok := syncCursor(item); ok {
			cursor = next
		}
		if !seeAll && !syncOwns(c, item) {
			continue
		}
		b, merr := json.Marshal(item)
		if merr != nil {
			continue
		}
		var m map[string]interface{}
		if uerr := json.Unmarshal(b, &m); uerr != nil {
			continue
		}
		m["_deleted"] = isSyncDeleted(item)
		rows = append(rows, policy.Projects(m))
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   rows,
		"cursor": cursor,
		"count":  len(rows),
	})
}

// isSyncDeleted reports whether a model row is soft-deleted (a tombstone).
func isSyncDeleted(obj interface{}) bool {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	f := v.FieldByName("DeletedAt")
	if !f.IsValid() {
		return false
	}
	if d, ok := f.Interface().(gorm.DeletedAt); ok {
		return d.Valid
	}
	return false
}

// syncCursor is where the next pull starts after row: its updated_at and id,
// as "<RFC3339Nano>~<id>". "~" needs no escaping in a query string.
func syncCursor(row interface{}) (string, bool) {
	v := reflect.ValueOf(row)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	updated, id := v.FieldByName("UpdatedAt"), v.FieldByName("ID")
	if !updated.IsValid() || !id.IsValid() {
		return "", false
	}
	t, ok := updated.Interface().(time.Time)
	if !ok || t.IsZero() {
		return "", false
	}
	return t.Format(time.RFC3339Nano) + "~" + fmt.Sprint(id.Interface()), true
}

// parseSyncCursor reads a cursor made by syncCursor, or a bare RFC3339 time from
// a client that stored one before cursors carried the id.
func parseSyncCursor(s string) (time.Time, string, error) {
	at, id := s, ""
	if i := strings.LastIndex(s, "~"); i >= 0 {
		at, id = s[:i], s[i+1:]
	}
	t, err := time.Parse(time.RFC3339Nano, at)
	return t, id, err
}

// ownedRow is a model whose rows belong to a user.
type ownedRow interface {
	GetOwnerID() string
}

// syncGranted reports whether the caller holds <model>.<action> or is ADMIN:
// the same test the resource's own routes make with RequireRole.
func syncGranted(c *gin.Context, model, action string) bool {
	if role, _ := c.Get("user_role"); role == "ADMIN" {
		return true
	}
	if grants, ok := c.Get("user_grants"); ok {
		if list, ok := grants.([]string); ok {
			return authz.Granted(list, model+"."+action)
		}
	}
	return false
}

// syncOwns reports whether row belongs to the caller.
func syncOwns(c *gin.Context, row interface{}) bool {
	o, ok := row.(ownedRow)
	if !ok {
		return false
	}
	uid := c.GetString("user_id")
	return uid != "" && o.GetOwnerID() == uid
}

// syncMayWrite: a holder of the permission may, and so may the owner of a row in
// a table whose rows have owners. Nobody else.
func syncMayWrite(c *gin.Context, model, action string, row interface{}) bool {
	return syncGranted(c, model, action) || syncOwns(c, row)
}

// serverFields are never taken from a pushed payload.
var serverFields = []string{"id", "version", "created_at", "updated_at", "deleted_at"}

func stripServerFields(data map[string]interface{}) map[string]interface{} {
	for _, k := range serverFields {
		delete(data, k)
	}
	return data
}

// decodeInto round-trips a map through JSON into the target struct so
// gorm field tags + types are respected. Cheap; the maps are small.
func decodeInto(target interface{}, data map[string]interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, target)
}

// setField sets a string field on a struct via reflection. Used for ID.
func setField(obj interface{}, name, value string) {
	v := reflect.ValueOf(obj).Elem()
	f := v.FieldByName(name)
	if f.IsValid() && f.CanSet() && f.Kind() == reflect.String {
		f.SetString(value)
	}
}

func getIntField(obj interface{}, name string) int {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	f := v.FieldByName(name)
	if !f.IsValid() {
		return 0
	}
	return int(f.Int())
}

// setIntField sets an int field on a struct via reflection. Used to seed
// Version before an update save so the BeforeUpdate hook bumps from the
// server's value rather than whatever the client happened to send.
func setIntField(obj interface{}, name string, value int) {
	v := reflect.ValueOf(obj).Elem()
	f := v.FieldByName(name)
	if f.IsValid() && f.CanSet() && f.CanInt() {
		f.SetInt(int64(value))
	}
}
