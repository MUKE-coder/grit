// Package concurrency is optimistic locking for update routes.
//
// Two people load the same record at version 7 and both save. Without this the
// second save silently overwrites the first: last write wins, and the first
// person's change is gone with nothing to say so. A client that sends the
// version it read, as If-Match, gets its write applied only if the record is
// still at that version, and a 409 naming the current one if it is not.
//
// Nothing changes for a client that does not send the header. Generated models
// already carry a version column that every update increments; this is the
// WHERE clause that makes it mean something.
package concurrency

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const askedKey = "concurrency.if_match"

// IfMatch scopes an update to the version the client last read, when it sent
// one. Use it on the update, and check Conflicted afterwards:
//
//	res := db.Model(&item).Scopes(concurrency.IfMatch(c)).Updates(updates)
//	if concurrency.Conflicted(c, res) { ... }
//
// Accepts a bare version (7), a quoted one ("7"), or the weak ETag the API
// sends back (W/"7"). A value that is none of those matches nothing, so the
// client gets a conflict rather than a blind write.
func IfMatch(c *gin.Context) func(*gorm.DB) *gorm.DB {
	raw := strings.TrimSpace(c.GetHeader("If-Match"))
	if raw == "" || raw == "*" {
		return func(db *gorm.DB) *gorm.DB { return db }
	}
	c.Set(askedKey, true)
	version, err := strconv.Atoi(strings.Trim(strings.TrimPrefix(raw, "W/"), `"`))
	if err != nil {
		return func(db *gorm.DB) *gorm.DB { return db.Where("1 = 0") }
	}
	return func(db *gorm.DB) *gorm.DB { return db.Where("version = ?", version) }
}

// Conflicted reports whether an update scoped by IfMatch changed no row: the
// record moved on after the client read it.
func Conflicted(c *gin.Context, res *gorm.DB) bool {
	_, asked := c.Get(askedKey)
	return asked && res.Error == nil && res.RowsAffected == 0
}

// WriteConflict answers 409 with the version the record is at now, so the
// client can reload and decide.
//
// Built by hand rather than through respond.Fail: Fail's details are strings,
// and current_version is a number on the wire. v3.285.0 sent it as "2", which
// broke every client comparing it with the version it holds.
func WriteConflict(c *gin.Context, current int) {
	c.Header("ETag", Tag(current))
	c.JSON(http.StatusConflict, gin.H{
		"error": gin.H{
			"code":    "VERSION_CONFLICT",
			"message": "This record changed after you loaded it. Reload it and try again.",
			"details": gin.H{"current_version": current},
		},
	})
}

// Tag is the ETag for a version, the value a client sends back as If-Match.
func Tag(version int) string {
	return `W/"` + strconv.Itoa(version) + `"`
}

// ── Without a request ────────────────────────────────────────────────────
//
// IfMatch and Conflicted read the gin context, so only a handler can use
// them. A service has to run the same update from a job, a command or a test,
// so it takes a Precondition instead: the handler reads it with FromRequest,
// the service scopes its update with it, and a miss comes back as an
// ErrConflict for the handler to answer with WriteConflict.
//
//	pre := concurrency.FromRequest(c)              // in the handler
//	res := db.Model(&item).Scopes(pre.Scope).Updates(updates)
//	if pre.Missed(res) {                           // in the service
//	    return &concurrency.ErrConflict{Current: current}
//	}

// Precondition is the version a client said it read, sent as If-Match. Nil
// means it sent none, so the write applies whatever the version is.
type Precondition struct {
	Version int
	// Unreadable is an If-Match that was not a version. It matches no row, so
	// the client gets a conflict rather than a blind write.
	Unreadable bool
}

// FromRequest reads If-Match. Nil when the header is absent or "*".
func FromRequest(c *gin.Context) *Precondition {
	raw := strings.TrimSpace(c.GetHeader("If-Match"))
	if raw == "" || raw == "*" {
		return nil
	}
	version, err := strconv.Atoi(strings.Trim(strings.TrimPrefix(raw, "W/"), "\""))
	if err != nil {
		return &Precondition{Unreadable: true}
	}
	return &Precondition{Version: version}
}

// Scope limits an update to a row still at the expected version. A nil
// Precondition scopes nothing, so it can be passed to Scopes unconditionally.
func (p *Precondition) Scope(db *gorm.DB) *gorm.DB {
	switch {
	case p == nil:
		return db
	case p.Unreadable:
		return db.Where("1 = 0")
	}
	return db.Where("version = ?", p.Version)
}

// Missed reports whether an update scoped by p changed no row: the record
// moved on after the client read it.
func (p *Precondition) Missed(res *gorm.DB) bool {
	return p != nil && res.Error == nil && res.RowsAffected == 0
}

// ErrConflict is a write whose precondition failed. Current is the version
// the record is at now, for the answer.
type ErrConflict struct {
	Current int
}

func (e *ErrConflict) Error() string {
	return "the record is at version " + strconv.Itoa(e.Current) + ", not the one the client read"
}
