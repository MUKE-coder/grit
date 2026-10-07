package scaffold

import (
	"strings"
	"testing"
)

// A refused origin says so.
//
// An origin outside CORS_ORIGINS gets a preflight with no
// Access-Control-Allow-Origin, which is right and silent: the browser never
// sends the real request, so nothing is logged and the page shows a button
// that does nothing. Running the admin on a port other than 3001 is enough.
func TestRefusedOriginIsReported(t *testing.T) {
	src := apiCorsMiddlewareGo()
	if !strings.Contains(src, "reportRefusedOrigin(origin)") {
		t.Error("a refused origin is still silent")
	}
	if !strings.Contains(src, "CORS_ORIGINS in .env") {
		t.Error("the message does not say what to change")
	}
	// Once per origin: a page that retries must not fill the log.
	if !strings.Contains(src, "refusedOrigins.LoadOrStore") {
		t.Error("the message repeats on every request")
	}
	// The allowlist itself is unchanged: this adds a log line, not an origin.
	if !strings.Contains(src, "if candidate == origin && origin != \"\"") {
		t.Error("the allowlist check changed")
	}
	if strings.Contains(src, `strings.HasPrefix(origin, "http://localhost")`) {
		t.Error("localhost became blanket-allowed, which is a different decision")
	}
}
