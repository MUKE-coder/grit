package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A channel could carry only what the server published, so "Ada is typing" cost
// a REST round trip per keystroke or was not built. And a hub that held every
// socket in the deployment reported no number at all: not how many it held, not
// how many messages it had thrown away for slow clients.
func TestRealtimeWhispersShipInTheTemplates(t *testing.T) {
	read := apiSource(t)

	whispers := read("internal", "realtime", "whispers.go")
	for _, want := range []string{
		"func (h *Hub) relayClientEvent(c *Client, raw []byte)",
		"var ClientEventsPerSecond = 10",
		"var MaxClientEventBytes = 1024",
		`h.refuseClientEvent(c, msg.Channel, "PUBLIC_CHANNEL"`,
		`h.refuseClientEvent(c, msg.Channel, "NOT_SUBSCRIBED"`,
		`h.refuseClientEvent(c, msg.Channel, "RATE_LIMITED"`,
		"h.publish(fanout{Channel: msg.Channel, Event: bytes})",
	} {
		if !strings.Contains(whispers, want) {
			t.Errorf("realtime/whispers.go is missing %q", want)
		}
	}

	stats := read("internal", "realtime", "stats.go")
	for _, want := range []string{
		"func (h *Hub) Stats() Stats",
		"func (h *Hub) offer(c *Client, bytes []byte, what, channel string) bool",
		`MessagesDropped uint64 ` + "`json:\"messages_dropped\"`",
		`BackplanePublishErrors uint64 ` + "`json:\"backplane_publish_errors\"`",
	} {
		if !strings.Contains(stats, want) {
			t.Errorf("realtime/stats.go is missing %q", want)
		}
	}
	read("internal", "realtime", "whispers_test.go")
	read("internal", "handlers", "realtime_greeting_test.go")

	for _, tc := range realtimeWhisperRepairs {
		src := read(tc.file...)
		for _, h := range tc.hunks {
			if strings.Count(src, h[1]) != 1 {
				t.Errorf("%s does not hold exactly one %q", strings.Join(tc.file, "/"), h[1])
			}
		}
	}

	// The greeting is queued before the hub can reach the connection at all.
	handler := read("internal", "handlers", "realtime.go")
	if strings.Index(handler, "client.Send <- greeting") > strings.Index(handler, "h.Hub.Admit(client)") {
		t.Error("the greeting is still queued after Admit, where DisconnectUser can close Send under it")
	}

	routes := read("internal", "routes", "routes.go")
	if !strings.Contains(routes, `"realtime": realtimeHub.Stats(),`) {
		t.Error("/api/health does not report the hub's stats")
	}
}

var realtimeWhisperRepairs = []struct {
	name   string
	file   []string
	hunks  [][2]string
	repair func(string) (string, []string, []string)
}{
	{"hub", []string{"internal", "realtime", "hub.go"}, realtimeHubWhisperHunks, repairRealtimeHubWhispersSource},
	{"channels", []string{"internal", "realtime", "channels.go"}, realtimeChannelsWhisperHunks, repairRealtimeChannelsWhispersSource},
	{"presence", []string{"internal", "realtime", "presence.go"}, realtimePresenceWhisperHunks, repairRealtimePresenceWhispersSource},
	{"backplane", []string{"internal", "realtime", "backplane.go"}, realtimeBackplaneWhisperHunks, repairRealtimeBackplaneWhispersSource},
	{"handler", []string{"internal", "handlers", "realtime.go"}, realtimeHandlerWhisperHunks, repairRealtimeHandlerWhispersSource},
	{"health", []string{"internal", "routes", "routes.go"}, [][2]string{{healthRealtimeOld, healthRealtimeNew}}, repairRealtimeHealthStatsSource},
}

func TestRealtimeWhisperRepairsProduceTheTemplate(t *testing.T) {
	read := apiSource(t)
	for _, tc := range realtimeWhisperRepairs {
		template := read(tc.file...)
		old := oldRealtimeFile(t, template, tc.hunks)
		out, fixed, warnings := tc.repair(old)
		if len(warnings) != 0 || len(fixed) == 0 {
			t.Fatalf("%s: fixed=%v warnings=%v", tc.name, fixed, warnings)
		}
		if gofmtSource(t, out) != gofmtSource(t, template) {
			t.Errorf("%s: the repaired file differs from the template", tc.name)
		}
		again, fixed, warnings := tc.repair(out)
		if again != out || len(fixed) != 0 || len(warnings) != 0 {
			t.Errorf("%s: not idempotent (fixed=%v warnings=%v)", tc.name, fixed, warnings)
		}
		edited := strings.Replace(old, tc.hunks[0][0], "// edited by hand\n", 1)
		if edited == old {
			t.Fatalf("%s: the edit did not apply", tc.name)
		}
		if got, _, warnings := tc.repair(edited); len(warnings) == 0 || got != edited {
			t.Errorf("%s: an edited file was rewritten or raised no warning", tc.name)
		}
	}
}

var whisperTestOpts = Options{ProjectName: "chat", Architecture: ArchTriple, Frontend: FrontendNext}

func realtimeWhisperAddedFiles(api string) []string {
	return []string{
		filepath.Join(api, "realtime", "whispers.go"),
		filepath.Join(api, "realtime", "whispers_test.go"),
		filepath.Join(api, "realtime", "stats.go"),
		filepath.Join(api, "handlers", "realtime_greeting_test.go"),
	}
}

// oldRealtimeWhispersProject scaffolds an API and turns it back into what
// v3.277.0 wrote, removing the files A4 adds.
func oldRealtimeWhispersProject(t *testing.T) (string, string, map[string]string) {
	t.Helper()
	root := t.TempDir()
	if err := createDirectories(root, whisperTestOpts); err != nil {
		t.Fatalf("createDirectories: %v", err)
	}
	if err := writeAPIFiles(root, whisperTestOpts); err != nil {
		t.Fatalf("writeAPIFiles: %v", err)
	}
	api := filepath.Join(root, "apps", "api", "internal")
	want := map[string]string{}
	for _, tc := range realtimeWhisperRepairs {
		path := filepath.Join(append([]string{root, "apps", "api"}, tc.file...)...)
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		want[path] = gofmtSource(t, string(b))
		if err := os.WriteFile(path, []byte(oldRealtimeFile(t, string(b), tc.hunks)), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, path := range realtimeWhisperAddedFiles(api) {
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
	}
	return root, api, want
}

// A v3.277.0 project, upgraded twice: the files come out as a fresh scaffold
// writes them, and the second run changes nothing.
func TestRepairRealtimeWhispersOnAnOldProject(t *testing.T) {
	root, api, want := oldRealtimeWhispersProject(t)
	for run := 0; run < 2; run++ {
		if err := repairRealtimeWhispers(root, whisperTestOpts); err != nil {
			t.Fatalf("run %d: %v", run+1, err)
		}
		for path, content := range want {
			b, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if gofmtSource(t, string(b)) != content {
				t.Errorf("run %d: %s differs from the template", run+1, filepath.Base(path))
			}
		}
		for _, path := range realtimeWhisperAddedFiles(api) {
			if !fileExists(path) {
				t.Errorf("run %d: %s was not added", run+1, filepath.Base(path))
			}
		}
	}
}

// whispers.go and stats.go only compile beside a hub that has the counters, so
// an edited hub means nothing is written and no other file is touched.
func TestRepairRealtimeWhispersLeavesAnEditedHubAlone(t *testing.T) {
	root, api, _ := oldRealtimeWhispersProject(t)
	hub := filepath.Join(api, "realtime", "hub.go")
	channels := filepath.Join(api, "realtime", "channels.go")
	b, err := os.ReadFile(hub)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(hub, []byte(strings.Replace(string(b), hubStatsFieldOld, "\tpresence presenceSet // edited\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(channels)
	if err != nil {
		t.Fatal(err)
	}
	if err := repairRealtimeWhispers(root, whisperTestOpts); err != nil {
		t.Fatal(err)
	}
	for _, path := range realtimeWhisperAddedFiles(api) {
		if fileExists(path) {
			t.Errorf("%s was written beside a hub it does not compile with", filepath.Base(path))
		}
	}
	after, err := os.ReadFile(channels)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Error("channels.go was changed although the hub could not be")
	}
}

// A project that never took the channels repair has no channels.go, so there is
// nothing for a client event to travel on and nothing to count.
func TestRepairRealtimeWhispersSkipsAProjectWithoutChannels(t *testing.T) {
	root, api, _ := oldRealtimeWhispersProject(t)
	if err := os.Remove(filepath.Join(api, "realtime", "channels.go")); err != nil {
		t.Fatal(err)
	}
	if err := repairRealtimeWhispers(root, whisperTestOpts); err != nil {
		t.Fatal(err)
	}
	for _, path := range realtimeWhisperAddedFiles(api) {
		if fileExists(path) {
			t.Errorf("%s was written into a project with no channels", filepath.Base(path))
		}
	}
}

// Every client the realtime files serve can send a client event and read one.
func TestRealtimeClientSendsWhispers(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "chat", Architecture: ArchTriple, Frontend: FrontendNext, IncludeExpo: true}
	if err := createDirectories(root, opts); err != nil {
		t.Fatalf("createDirectories: %v", err)
	}
	if err := writeRealtimeClientFiles(root, opts); err != nil {
		t.Fatal(err)
	}
	for _, app := range []string{"web", "admin", "expo"} {
		lib, err := os.ReadFile(filepath.Join(root, "apps", app, "lib", "realtime.ts"))
		if err != nil {
			t.Fatalf("%s: %v", app, err)
		}
		for _, want := range []string{
			"export function whisper(channel: string, event: string, data?: unknown): boolean",
			`type: "client-event"`,
			`"client-event:" + name`,
		} {
			if !strings.Contains(string(lib), want) {
				t.Errorf("%s: lib/realtime.ts is missing %q", app, want)
			}
		}
		hook, err := os.ReadFile(filepath.Join(root, "apps", app, "hooks", "use-realtime.ts"))
		if err != nil {
			t.Fatalf("%s: %v", app, err)
		}
		if !strings.Contains(string(hook), "export function useWhisper(channel: string | null | undefined)") {
			t.Errorf("%s: hooks/use-realtime.ts has no useWhisper", app)
		}
	}
}
