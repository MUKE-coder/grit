package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WriteConcurrencyPackage writes internal/concurrency under apiRoot.
// overwrite=false leaves an existing copy alone, for the generator on an older
// project; grit upgrade passes true, and the write is manifest-guarded.
func WriteConcurrencyPackage(apiRoot, module string, overwrite bool) error {
	files := map[string]string{
		filepath.Join(apiRoot, "internal", "concurrency", "concurrency.go"):      apiConcurrencyGo(),
		filepath.Join(apiRoot, "internal", "concurrency", "concurrency_test.go"): apiConcurrencyTestGo(),
	}
	for path, content := range files {
		if !overwrite {
			if _, err := os.Stat(path); err == nil {
				continue
			}
		}
		if err := writeFile(path, strings.ReplaceAll(content, "{{MODULE}}", module)); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

func apiConcurrencyGo() string {
	return `// Package concurrency is optimistic locking for update routes.
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
	version, err := strconv.Atoi(strings.Trim(strings.TrimPrefix(raw, "W/"), ` + "`" + `"` + "`" + `))
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
	return ` + "`" + `W/"` + "`" + ` + strconv.Itoa(version) + ` + "`" + `"` + "`" + `
}
`
}

func apiConcurrencyTestGo() string {
	return `package concurrency

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type lot struct {
	ID      uint
	Bid     int
	Version int ` + "`" + `gorm:"not null;default:1"` + "`" + `
}

func (l *lot) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}

func ctxWith(ifMatch string) *gin.Context {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("PATCH", "/", nil)
	if ifMatch != "" {
		c.Request.Header.Set("If-Match", ifMatch)
	}
	return c
}

// Two bidders read the lot at version 1 and both bid. Only the first lands;
// the second is told the lot moved on instead of overwriting the first bid.
func TestTheSecondWriterOfAVersionConflicts(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.AutoMigrate(&lot{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	item := lot{Bid: 100, Version: 1}
	db.Create(&item)

	first := ctxWith(Tag(1))
	res := db.Model(&lot{ID: item.ID}).Scopes(IfMatch(first)).Updates(map[string]interface{}{"bid": 110})
	if res.Error != nil || Conflicted(first, res) {
		t.Fatalf("the first writer of version 1 was refused: %v", res.Error)
	}

	second := ctxWith(` + "`" + `"1"` + "`" + `)
	res = db.Model(&lot{ID: item.ID}).Scopes(IfMatch(second)).Updates(map[string]interface{}{"bid": 105})
	if !Conflicted(second, res) {
		t.Fatal("the second writer of version 1 overwrote the first")
	}

	var got lot
	db.First(&got, item.ID)
	if got.Bid != 110 || got.Version != 2 {
		t.Errorf("got bid %d at version %d, want 110 at version 2", got.Bid, got.Version)
	}

	// No header: last write wins, as it always has.
	plain := ctxWith("")
	res = db.Model(&lot{ID: item.ID}).Scopes(IfMatch(plain)).Updates(map[string]interface{}{"bid": 120})
	if res.Error != nil || Conflicted(plain, res) {
		t.Error("an update without If-Match was refused")
	}

	// Garbage never writes blindly.
	junk := ctxWith("not-a-version")
	res = db.Model(&lot{ID: item.ID}).Scopes(IfMatch(junk)).Updates(map[string]interface{}{"bid": 1})
	if !Conflicted(junk, res) {
		t.Error("an unreadable If-Match wrote anyway")
	}
}
`
}
