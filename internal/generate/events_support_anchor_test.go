package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// routes.go stopped containing "realtimeHub := realtime.NewHub()" when the hub
// gained options, so a project without the event bus got a warning instead of
// the boot lines. Both current shapes of the line must be found.
func TestEnsureEventsBootFindsTheCurrentHubLine(t *testing.T) {
	for name, hub := range map[string]string{
		"with redis": "\trealtimeHub := realtime.NewHub(realtime.WithRedis(cfg.RedisURL, \"\"))\n",
		"gated":      "\tvar realtimeOptions []realtime.Option\n\trealtimeHub := realtime.NewHub(realtimeOptions...)\n",
	} {
		root := t.TempDir()
		dir := filepath.Join(root, "apps", "api", "internal", "routes")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		src := "package routes\n\nimport (\n\t\"app/internal/handlers\"\n)\n\nfunc Setup() {\n" + hub + "\t_ = handlers.X\n}\n"
		path := filepath.Join(dir, "routes.go")
		if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
		g := &Generator{Root: root, Module: "app", Architecture: "triple"}
		if err := g.ensureEventsBoot(); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		out := string(b)
		boot := strings.Index(out, "events.Init(4)")
		if boot < 0 {
			t.Fatalf("%s: the event bus was not started:\n%s", name, out)
		}
		if boot < strings.Index(out, "realtimeHub := realtime.NewHub(") {
			t.Errorf("%s: events.Init was inserted before the hub exists", name)
		}
		if !strings.Contains(out, "realtime.NewHub(") || strings.Contains(out, "NewHub(\n") {
			t.Errorf("%s: the hub line was split:\n%s", name, out)
		}
	}
}
