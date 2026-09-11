package scaffold

import (
	"strings"
	"testing"
)

// The activity middleware skipped every GET, so reads of health records went
// unrecorded. It now records the reads a handler marks, and only those.
func TestMiddlewareRecordsMarkedReads(t *testing.T) {
	src := apiActivityMiddlewareGo()
	for _, want := range []string{
		"recordRead(c)",
		"audit.ReadMarkOf(c)",
		// Search terms can be personal data; only their digest is kept.
		"digestBody([]byte(c.Request.URL.RawQuery))",
		// Writes say which record they touched.
		`ResourceIDs: c.Param("id")`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the activity middleware is missing %s", want)
		}
	}
}

func TestAuditCanMarkAndCapReads(t *testing.T) {
	src := apiAuditGo()
	for _, want := range []string{
		"func Read(c *gin.Context, resource string, ids ...string)",
		"func ReadCount(c *gin.Context, resource string, n int)",
		"func ReadMarkOf(c *gin.Context) (ReadMark, bool)",
		"const MaxReadIDs = 1000",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("audit.go is missing %s", want)
		}
	}
}

// "Who read this record" has to be a query an administrator can run.
func TestActivityListFindsARecord(t *testing.T) {
	src := apiActivityHandlerGo()
	if !strings.Contains(src, `With("resource", c.Query("resource"))`) ||
		!strings.Contains(src, `q.Where("resource_ids LIKE ?", "%"+record+"%")`) {
		t.Error("the activity list cannot be filtered to one resource or one record")
	}
	if !strings.Contains(apiActivityReadTestGo(), "func TestMarkedReadsAreRecorded(") {
		t.Error("the shipped read-audit test is missing")
	}
}
