package scaffold

import (
	"strings"
	"testing"
)

// Each repair, fed the shape a project in the field actually has, then fed
// its own output again: a repair that is not idempotent corrupts a file on
// the second upgrade, and an upgrade runs more than once.

func TestRepairCORSOptimisticLocking(t *testing.T) {
	before := `package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

var allowedRequestHeaders = []string{
	"Origin",
	"Content-Type",
}

// allowedRequestHeadersValue is the joined header value, built once.
var allowedRequestHeadersValue = strings.Join(allowedRequestHeaders, ", ")

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Headers", allowedRequestHeadersValue)
	}
}
`
	out, changes, warnings := repairCORSOptimisticLockingSource(before)
	if len(warnings) > 0 {
		t.Fatalf("warnings: %v", warnings)
	}
	if len(changes) != 2 {
		t.Errorf("expected two changes reported, got %v", changes)
	}
	for _, want := range []string{
		`"If-Match",`,
		"var exposedResponseHeaders = []string{",
		`c.Header("Access-Control-Expose-Headers", exposedResponseHeadersValue)`,
		"func ExposeResponseHeader(",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing after repair: %s", want)
		}
	}

	again, changes, _ := repairCORSOptimisticLockingSource(out)
	if again != out {
		t.Error("the CORS repair is not idempotent: a second upgrade changes the file again")
	}
	if len(changes) != 0 {
		t.Errorf("a second run reported changes: %v", changes)
	}
}

// A project still on the string literal belongs to repairCORSHeaderList, and
// this one must not touch it.
func TestRepairCORSLeavesTheLiteralShapeAlone(t *testing.T) {
	before := `c.Header("Access-Control-Allow-Headers", "Origin, Content-Type")`
	out, changes, warnings := repairCORSOptimisticLockingSource(before)
	if out != before || len(changes) != 0 || len(warnings) != 0 {
		t.Error("the literal shape was touched; repairCORSHeaderList reshapes it first")
	}
}

func TestRepairAliasStatus(t *testing.T) {
	before := `package routes

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func mountLegacyAPIAlias(r *gin.Engine) {
	r.NoRoute(func(c *gin.Context) {
		if true {
			c.Header("Deprecation", "true")
			r.HandleContext(c)
			return
		}
	})
}
`
	out, changes, warnings := repairAliasStatusSource(before)
	if len(warnings) > 0 {
		t.Fatalf("warnings: %v", warnings)
	}
	if len(changes) == 0 {
		t.Error("the repair reported no change")
	}
	status := strings.Index(out, "c.Status(http.StatusOK)")
	handle := strings.Index(out, "r.HandleContext(c)")
	if status < 0 {
		t.Fatal("the pending 404 is still not cleared")
	}
	if status > handle {
		t.Error("the status is reset after HandleContext, which is too late")
	}
	if !strings.Contains(out, `"net/http"`) {
		t.Error(`the repair did not add the "net/http" import, so the file will not compile`)
	}

	if again, _, _ := repairAliasStatusSource(out); again != out {
		t.Error("the alias repair is not idempotent")
	}
}

func TestRepairCallerMessageID(t *testing.T) {
	before := `package mail

import "strings"

func buildMIME(m *Message, now time.Time) ([]byte, error) {
	id, err := messageID(from.Address)
	if err != nil {
		return nil, err
	}
	writeHeader(&out, "Message-ID", id)
	return nil, nil
}

func messageID(from string) (string, error) {
	return "", nil
}
`
	out, changes, warnings := repairCallerMessageIDSource(before)
	if len(warnings) > 0 {
		t.Fatalf("warnings: %v", warnings)
	}
	if len(changes) == 0 {
		t.Error("the repair reported no change")
	}
	for _, want := range []string{
		"resolveMessageID(m.Headers, from.Address)",
		"func resolveMessageID(",
		"func validMessageID(",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing after repair: %s", want)
		}
	}
	if strings.Contains(out, "id, err := messageID(from.Address)") {
		t.Error("the old call is still there, so the generated id still wins")
	}

	if again, _, _ := repairCallerMessageIDSource(out); again != out {
		t.Error("the Message-ID repair is not idempotent")
	}
}

func TestRepairSequenceSave(t *testing.T) {
	before := `		c.NextValue++
		next = c.NextValue
		return tx.Save(&c).Error
	})
`
	out, changes, warnings := repairSequenceSaveSource(before)
	if len(warnings) > 0 {
		t.Fatalf("warnings: %v", warnings)
	}
	if len(changes) == 0 {
		t.Error("the repair reported no change")
	}
	if strings.Contains(out, "tx.Save(&c)") {
		t.Error("Save is still there, so the second number still fails")
	}
	if !strings.Contains(out, `Where("name = ? AND bucket = ?", cfg.Name, bucket).`) {
		t.Error("the counter is not updated by its key columns")
	}

	if again, _, _ := repairSequenceSaveSource(out); again != out {
		t.Error("the sequence repair is not idempotent")
	}
}
