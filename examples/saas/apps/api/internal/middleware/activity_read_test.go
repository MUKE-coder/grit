package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"saas/apps/api/internal/audit"
	"saas/apps/api/internal/models"
)

// A read of a resource generated with --audit-reads is recorded with the rows
// it served. An unmarked read is not, and neither is one that failed. Reads
// went unrecorded entirely before v3.216.0.
func TestMarkedReadsAreRecorded(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&models.ActivityLog{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_id", "reader-1"); c.Next() })
	r.Use(ActivityLogger(db))
	r.GET("/notes/:id", func(c *gin.Context) {
		audit.Read(c, "notes", c.Param("id"))
		c.JSON(http.StatusOK, gin.H{})
	})
	r.GET("/plain", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) })
	r.GET("/gone/:id", func(c *gin.Context) {
		audit.Read(c, "notes", c.Param("id"))
		c.JSON(http.StatusNotFound, gin.H{})
	})
	for _, path := range []string{"/notes/n1?q=Jane+Roe", "/plain", "/gone/n2"} {
		r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, path, nil))
	}

	// The writer is asynchronous: wait for the read to land, then a moment
	// more, so an entry that should not exist has had the chance to.
	var reads []models.ActivityLog
	for deadline := time.Now().Add(3 * time.Second); time.Now().Before(deadline); time.Sleep(20 * time.Millisecond) {
		if db.Where("method = ?", "GET").Find(&reads); len(reads) > 0 {
			break
		}
	}
	time.Sleep(200 * time.Millisecond)
	db.Where("method = ?", "GET").Find(&reads)

	if len(reads) != 1 {
		t.Fatalf("recorded %d reads, want only the marked, successful one", len(reads))
	}
	got := reads[0]
	if got.Resource != "notes" || got.ResourceIDs != "n1" || got.RecordCount != 1 || got.UserID != "reader-1" {
		t.Errorf("the read was recorded wrongly: %+v", got)
	}
	sum := sha256.Sum256([]byte("q=Jane+Roe"))
	if got.PayloadDigest != hex.EncodeToString(sum[:]) {
		t.Error("the query was not recorded as its digest")
	}
}
