package concurrency

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// The same guarantee with no request at all, which is how a service sees it:
// a Precondition rather than a gin context.
func TestAPreconditionNeedsNoRequest(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.AutoMigrate(&lot{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	item := lot{Bid: 100, Version: 1}
	db.Create(&item)

	read := &Precondition{Version: 1}
	res := db.Model(&lot{ID: item.ID}).Scopes(read.Scope).Updates(map[string]interface{}{"bid": 110})
	if res.Error != nil || read.Missed(res) {
		t.Fatalf("the first writer of version 1 was refused: %v", res.Error)
	}
	res = db.Model(&lot{ID: item.ID}).Scopes(read.Scope).Updates(map[string]interface{}{"bid": 105})
	if !read.Missed(res) {
		t.Fatal("a stale precondition overwrote the first writer")
	}

	var none *Precondition
	res = db.Model(&lot{ID: item.ID}).Scopes(none.Scope).Updates(map[string]interface{}{"bid": 120})
	if res.Error != nil || none.Missed(res) {
		t.Error("an update with no precondition was refused")
	}

	if FromRequest(ctxWith("")) != nil {
		t.Error("no If-Match read as a precondition")
	}
	if p := FromRequest(ctxWith(Tag(3))); p == nil || p.Version != 3 {
		t.Errorf("W/\"3\" read as %+v", p)
	}
	if p := FromRequest(ctxWith("junk")); p == nil || !p.Unreadable {
		t.Error("an unreadable If-Match was not marked unreadable")
	}
}

type lot struct {
	ID      uint
	Bid     int
	Version int `gorm:"not null;default:1"`
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

	second := ctxWith(`"1"`)
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

// current_version is a number on the wire, because a client compares it with
// the version it holds. v3.285.0 sent it as a string and broke that.
func TestAConflictNamesTheCurrentVersionAsANumber(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	WriteConflict(c, 2)

	if w.Code != 409 {
		t.Fatalf("status %d, want 409", w.Code)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("body is not JSON: %v", err)
	}
	errObj, _ := body["error"].(map[string]interface{})
	details, _ := errObj["details"].(map[string]interface{})
	if errObj["code"] != "VERSION_CONFLICT" {
		t.Errorf("code %v, want VERSION_CONFLICT", errObj["code"])
	}
	if v, ok := details["current_version"].(float64); !ok || v != 2 {
		t.Errorf("current_version is %#v, want the number 2", details["current_version"])
	}
	if got := w.Header().Get("ETag"); got != Tag(2) {
		t.Errorf("ETag %q, want %q", got, Tag(2))
	}
}
