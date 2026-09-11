package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Every Postgres deployment's activity log failed verification on its first
// entry: the hash covered a nanosecond timestamp and Postgres keeps
// microseconds. Proven on a live project by finding the 400ns offset that
// reproduced the stored hash from the stored row.
func TestChainIsStampedAtStoredPrecision(t *testing.T) {
	src := apiAuditGo()
	for _, want := range []string{
		"const Precision = time.Millisecond",
		"time.Now().UTC().Truncate(Precision)",
		// Replicas each run a writer; the lock is what stops them forking.
		"pg_advisory_xact_lock",
		"func Reseal(",
		// New canonical fields must not change the bytes of old entries.
		`json:"resource,omitempty"`,
		`json:"resource_ids,omitempty"`,
		`json:"record_count,omitempty"`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("audit.go is missing %s", want)
		}
	}
}

func TestMiddlewareUsesTheChainWriter(t *testing.T) {
	src := apiActivityMiddlewareGo()
	if !strings.Contains(src, "audit.Start(db)") || !strings.Contains(src, "audit.Enqueue(entry)") {
		t.Error("the activity middleware does not hand entries to the chain writer")
	}
	for _, gone := range []string{"startAuditWorker", "auditChan", "CreatedAt:"} {
		if strings.Contains(src, gone) {
			t.Errorf("the middleware still has %s: the writer stamps and chains entries", gone)
		}
	}
}

// Security events were inserted with no hash: the first broke verification,
// and the unique index on hash refused every one after it.
func TestSecurityEventsJoinTheChain(t *testing.T) {
	if !strings.Contains(securityLogGo(), "audit.AppendChained(db.WithContext(ctx), &entry)") {
		t.Error("security events bypass the chain")
	}
}

func TestResealIsMountedAndInjected(t *testing.T) {
	if !strings.Contains(apiActivityHandlerGo(), "func (h *ActivityHandler) Reseal(") {
		t.Error("the activity handler has no Reseal")
	}

	api := t.TempDir()
	routes := filepath.Join(api, "internal", "routes", "routes.go")
	if err := os.MkdirAll(filepath.Dir(routes), 0o755); err != nil {
		t.Fatal(err)
	}
	body := "package routes\n\nfunc mount() {\n" +
		"\t\tadmin.GET(\"/admin/activity/integrity\", activityHandler.VerifyIntegrity)\r\n" +
		"}\n"
	body = strings.ReplaceAll(strings.ReplaceAll(body, "\r\n", "\n"), "\n", "\r\n")
	if err := os.WriteFile(routes, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := ensureAuditResealRoute(api); err != nil {
			t.Fatalf("inject: %v", err)
		}
	}
	got, err := os.ReadFile(routes)
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(string(got), "activityHandler.Reseal"); n != 1 {
		t.Errorf("reseal route mounted %d times, want once", n)
	}
	if !strings.Contains(string(got), "\r\n") || strings.Contains(strings.ReplaceAll(string(got), "\r\n", ""), "\n") {
		t.Error("the file's CRLF line endings were not kept")
	}
}
