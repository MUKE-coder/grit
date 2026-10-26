package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The plugin told clients to send X-Organization-ID and did not tell the
// browser it was allowed to. The preflight succeeded, the real request was
// never sent, and there was nothing to see: no status, no log line, an axios
// network error with nothing behind it, and the same call working under curl.
// A user in one organization was unaffected, because the API picks that one
// without being told, so it waited for a second organization (grit#107).
func TestTenantMiddlewareRegistersItsOwnHeader(t *testing.T) {
	src := mtMiddleware(Context{Module: "demo/apps/api"})

	if !strings.Contains(src, `func init() { AllowRequestHeader("X-Organization-ID") }`) {
		t.Error("tenant.go does not register the header it reads, so the browser strips it")
	}
	// From this file on purpose. The alternative was injecting a line into
	// cors.go, which removal then has to find again; this way the header
	// arrives and leaves with the file that reads it.
	if !strings.Contains(src, "X-Organization-ID") {
		t.Fatal("the middleware no longer reads the header at all")
	}
	for _, inj := range multitenantInjections(Context{Architecture: "triple", Module: "demo/apps/api"}) {
		if strings.Contains(inj.File, "cors.go") {
			t.Errorf("the plugin injects into %s; registering from tenant.go is what makes "+
				"removal take the header out without editing that file", inj.File)
		}
	}
}

// And a project whose cors.go predates AllowRequestHeader is told what to run,
// before anything is written. Installing anyway would write a tenant.go that
// does not compile, which is worse than not installing.
func TestMultitenantPreflightRefusesAnOlderCORS(t *testing.T) {
	root := t.TempDir()
	mw := filepath.Join(root, "apps", "api", "internal", "middleware")
	if err := os.MkdirAll(mw, 0o755); err != nil {
		t.Fatal(err)
	}
	cors := filepath.Join(mw, "cors.go")
	ctx := Context{Root: root, APIRoot: filepath.Join(root, "apps", "api")}

	old := `package middleware

func CORS() {
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type")
}
`
	if err := os.WriteFile(cors, []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}
	err := multitenantPreflight(ctx)
	if err == nil {
		t.Fatal("an older cors.go was accepted")
	}
	if !strings.Contains(err.Error(), "grit upgrade") {
		t.Errorf("the error does not say what to run: %v", err)
	}

	// The current shape passes.
	if err := os.WriteFile(cors, []byte(old+"\nfunc AllowRequestHeader(name string) {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := multitenantPreflight(ctx); err != nil {
		t.Errorf("a current cors.go was refused: %v", err)
	}

	// And no cors.go at all is not this plugin's to diagnose: the injections
	// fail on their own targets with a better message than one invented here.
	if err := os.Remove(cors); err != nil {
		t.Fatal(err)
	}
	if err := multitenantPreflight(ctx); err != nil {
		t.Errorf("a missing cors.go should not be this plugin's error: %v", err)
	}
}
