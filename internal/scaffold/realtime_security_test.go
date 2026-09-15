package scaffold

import (
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The socket authenticates with the grit_access cookie, which a browser sends
// with a handshake from any page. CheckOrigin returned true for every origin, so
// a page on another origin opened a socket as the signed-in user and read their
// events. Reproduced from a headless browser against a stock scaffold.
func TestRealtimeSocketChecksOriginAndCapsConnections(t *testing.T) {
	read := apiSource(t)

	h := read("internal", "handlers", "realtime.go")
	if strings.Contains(h, "func(r *http.Request) bool { return true }") {
		t.Error("the upgrader still accepts every origin")
	}
	for _, want := range []string{
		"CheckOrigin:     realtime.CheckOrigin",
		"if !realtime.CheckOrigin(c.Request) {",
		"h.Hub.Admit(client)",
		"websocket.CloseTryAgainLater",
	} {
		if !strings.Contains(h, want) {
			t.Errorf("handlers/realtime.go is missing %q", want)
		}
	}

	guard := read("internal", "realtime", "guard.go")
	for _, want := range []string{
		"func CheckOrigin(r *http.Request) bool",
		"func (h *Hub) Admit(c *Client) error",
		`"REALTIME_MAX_CONNECTIONS_PER_USER", 10`,
		`"REALTIME_MAX_CONNECTIONS", 10000`,
		`allowed != "*"`,
	} {
		if !strings.Contains(guard, want) {
			t.Errorf("realtime/guard.go is missing %q", want)
		}
	}
	read("internal", "realtime", "guard_test.go")
	read("internal", "handlers", "realtime_test.go")

	routes := read("internal", "routes", "routes.go")
	for _, want := range []string{
		"r.Use(middleware.CORSDynamic(corsOrigins))",
		"realtime.AllowedOrigins = corsOrigins",
		"if cfg.Modules.Realtime {\n\t\trealtimeOptions = append(realtimeOptions, realtime.WithRedis(cfg.RedisURL, \"\"))",
		"if cfg.Modules.Realtime {\n\t\tr.GET(\"/api/ws\", realtimeHandler.Connect)",
	} {
		if !strings.Contains(routes, want) {
			t.Errorf("routes.go is missing %q", want)
		}
	}
}

func gofmtSource(t *testing.T, src string) string {
	t.Helper()
	out, err := format.Source([]byte(src))
	if err != nil {
		t.Fatalf("gofmt: %v", err)
	}
	return string(out)
}

// oldRealtimeFile rebuilds the file Grit wrote before the fix from the template.
func oldRealtimeFile(t *testing.T, template string, hunks [][2]string) string {
	t.Helper()
	old := template
	for _, h := range hunks {
		if n := strings.Count(old, h[1]); n != 1 {
			t.Fatalf("the template holds the new text %d times, want 1:\n%s", n, h[1])
		}
		old = strings.Replace(old, h[1], h[0], 1)
	}
	return old
}

var realtimeHandlerHunks = [][2]string{
	{realtimeHandlerImportOld, realtimeHandlerImportNew},
	{realtimeUpgraderOld, realtimeUpgraderNew},
	{realtimeTokenOld, realtimeTokenNew},
	{realtimeRegisterOld, realtimeRegisterNew},
}

var realtimeRoutesHunks = [][2]string{
	{routesCORSOld, routesCORSNew},
	{routesRealtimeHubOld, routesRealtimeHubNew},
	{routesRealtimeHandlerOld, routesRealtimeHandlerNew},
	{routesRealtimeRouteOld, routesRealtimeRouteNew},
}

func TestRealtimeRepairsProduceTheTemplate(t *testing.T) {
	read := apiSource(t)
	cases := []struct {
		name     string
		template string
		hunks    [][2]string
		repair   func(string) (string, []string, []string)
	}{
		{"handler", read("internal", "handlers", "realtime.go"), realtimeHandlerHunks, repairRealtimeHandlerSource},
		{"routes", read("internal", "routes", "routes.go"), realtimeRoutesHunks, func(s string) (string, []string, []string) {
			return repairRealtimeRoutesSource(s, true)
		}},
	}
	for _, tc := range cases {
		old := oldRealtimeFile(t, tc.template, tc.hunks)
		out, fixed, warnings := tc.repair(old)
		if len(warnings) != 0 || len(fixed) == 0 {
			t.Fatalf("%s: fixed=%v warnings=%v", tc.name, fixed, warnings)
		}
		if gofmtSource(t, out) != gofmtSource(t, tc.template) {
			t.Errorf("%s: the repaired file differs from the template", tc.name)
		}
		again, fixed, warnings := tc.repair(out)
		if again != out || len(fixed) != 0 || len(warnings) != 0 {
			t.Errorf("%s: not idempotent (fixed=%v warnings=%v)", tc.name, fixed, warnings)
		}
		// An edited file is warned about, not rewritten.
		edited := strings.Replace(old, tc.hunks[1][0], "// edited by hand\n", 1)
		if _, _, warnings := tc.repair(edited); len(warnings) == 0 {
			t.Errorf("%s: an edited file raised no warning", tc.name)
		}
	}
}

// A project scaffolded before the fix, upgraded: every file comes out as a fresh
// scaffold would write it.
func TestRepairRealtimeSecurityOnAnOldProject(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "chat", Architecture: ArchTriple, Frontend: FrontendNext}
	if err := createDirectories(root, opts); err != nil {
		t.Fatalf("createDirectories: %v", err)
	}
	if err := writeAPIFiles(root, opts); err != nil {
		t.Fatalf("writeAPIFiles: %v", err)
	}
	api := filepath.Join(root, "apps", "api", "internal")
	files := map[string][][2]string{
		filepath.Join(api, "handlers", "realtime.go"): realtimeHandlerHunks,
		filepath.Join(api, "routes", "routes.go"):     realtimeRoutesHunks,
	}
	want := map[string]string{}
	for path, hunks := range files {
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		want[path] = gofmtSource(t, string(b))
		if err := os.WriteFile(path, []byte(oldRealtimeFile(t, string(b), hunks)), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	added := []string{
		filepath.Join(api, "realtime", "guard.go"),
		filepath.Join(api, "realtime", "guard_test.go"),
		filepath.Join(api, "handlers", "realtime_test.go"),
	}
	for _, path := range added {
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
	}

	for run := 0; run < 2; run++ {
		if err := repairRealtimeSecurity(root, opts); err != nil {
			t.Fatalf("run %d: %v", run+1, err)
		}
		for path, content := range want {
			b, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(b) != content {
				t.Errorf("run %d: %s differs from the template", run+1, filepath.Base(path))
			}
		}
		for _, path := range added {
			if !fileExists(path) {
				t.Errorf("run %d: %s was not added", run+1, filepath.Base(path))
			}
		}
	}
}
