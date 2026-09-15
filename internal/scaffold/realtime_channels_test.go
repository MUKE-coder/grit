package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The hub could address a user id or everyone and nothing else, so a page could
// not follow one record, and an app could not say "whoever may read this".
func TestRealtimeChannelsShipInTheTemplates(t *testing.T) {
	read := apiSource(t)

	channels := read("internal", "realtime", "channels.go")
	for _, want := range []string{
		"func Channel(pattern string, authorize func(c ChannelContext) bool)",
		"func (h *Hub) Publish(channel string, evt Event)",
		"func (h *Hub) HandleClientMessage(c *Client, raw []byte)",
		`h.reply(c, "subscribed", msg.Channel, nil)`,
		`"subscription_error"`,
		`strings.HasPrefix(channel, "public-")`,
		"h.publish(fanout{Channel: channel, Event: bytes})",
	} {
		if !strings.Contains(channels, want) {
			t.Errorf("realtime/channels.go is missing %q", want)
		}
	}
	read("internal", "realtime", "channels_test.go")
	read("internal", "handlers", "realtime_channels_test.go")

	for file, wants := range map[string][]string{
		"hub.go":       {hubChannelsFieldNew, hubUnregisterNew, hubDisconnectNew},
		"backplane.go": {backplaneFanoutNew, backplaneReceiveNew},
	} {
		src := read("internal", "realtime", file)
		for _, want := range wants {
			if !strings.Contains(src, want) {
				t.Errorf("realtime/%s is missing %q", file, want)
			}
		}
	}

	handler := read("internal", "handlers", "realtime.go")
	if !strings.Contains(handler, "hub.HandleClientMessage(c, msg)") {
		t.Error("readPump does not hand client messages to the hub")
	}
	// The longest subscribe message is about 200 bytes, so the limit stays.
	if !strings.Contains(handler, "wsMaxMessageSize = 1024") {
		t.Error("the socket read limit changed; channel messages do not need more than 1 KB")
	}

	subs := read("internal", "services", "event_subscribers.go")
	for _, want := range []string{
		"hub.SendToUsers(audience, evt)",
		"var RealtimeChannels func(e events.Event) []string",
		"hub.Publish(channel, evt)",
	} {
		if !strings.Contains(subs, want) {
			t.Errorf("event_subscribers.go is missing %q", want)
		}
	}
}

var realtimeReadPumpHunks = [][2]string{
	{realtimeReadCommentOld, realtimeReadCommentNew},
	{realtimeReadLoopOld, realtimeReadLoopNew},
}

func TestRealtimeChannelRepairsProduceTheTemplate(t *testing.T) {
	read := apiSource(t)
	cases := []struct {
		name     string
		template string
		hunks    [][2]string
		repair   func(string) (string, []string, []string)
	}{
		{"hub", read("internal", "realtime", "hub.go"), realtimeHubChannelsHunks, repairRealtimeHubChannelsSource},
		{"backplane", read("internal", "realtime", "backplane.go"), realtimeBackplaneChannelsHunks, repairRealtimeBackplaneChannelsSource},
		{"handler", read("internal", "handlers", "realtime.go"), realtimeReadPumpHunks, repairRealtimeReadPumpSource},
		{"subscribers", read("internal", "services", "event_subscribers.go"), realtimeSubscribersChannelsHunks, repairRealtimeChannelsSubscribersSource},
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
		edited := strings.Replace(old, tc.hunks[1][0], "// edited by hand\n", 1)
		if got, _, warnings := tc.repair(edited); len(warnings) == 0 || got != edited {
			t.Errorf("%s: an edited file was rewritten or raised no warning", tc.name)
		}
	}
}

// A project from v3.273.0, upgraded: every file comes out as a fresh scaffold
// writes it, and a second upgrade changes nothing.
func TestRepairRealtimeChannelsOnAnOldProject(t *testing.T) {
	root, api, want := oldRealtimeChannelsProject(t)
	added := realtimeChannelsAddedFiles(api)

	for run := 0; run < 2; run++ {
		if err := repairRealtimeChannels(root, Options{ProjectName: "chat", Architecture: ArchTriple, Frontend: FrontendNext}); err != nil {
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

// channels.go needs the hub field and the fanout Channel, so an edited hub.go
// means nothing is written: not channels.go, not the backplane.
func TestRepairRealtimeChannelsLeavesAnEditedHubAlone(t *testing.T) {
	root, api, _ := oldRealtimeChannelsProject(t)
	hub := filepath.Join(api, "realtime", "hub.go")
	backplane := filepath.Join(api, "realtime", "backplane.go")
	b, err := os.ReadFile(hub)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(hub, []byte(strings.Replace(string(b), hubUnregisterOld, "\t\t\tdelete(set, c)\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(backplane)
	if err != nil {
		t.Fatal(err)
	}
	if err := repairRealtimeChannels(root, Options{ProjectName: "chat", Architecture: ArchTriple, Frontend: FrontendNext}); err != nil {
		t.Fatal(err)
	}
	if fileExists(filepath.Join(api, "realtime", "channels.go")) {
		t.Error("channels.go was written beside a hub it does not compile with")
	}
	after, err := os.ReadFile(backplane)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Error("the backplane was changed although the hub could not be")
	}
}

func realtimeChannelsAddedFiles(api string) []string {
	return []string{
		filepath.Join(api, "realtime", "channels.go"),
		filepath.Join(api, "realtime", "channels_test.go"),
		filepath.Join(api, "handlers", "realtime_channels_test.go"),
	}
}

// oldRealtimeChannelsProject scaffolds an API and turns its realtime files back
// into what v3.273.0 wrote. It returns the root, the internal directory and the
// expected (gofmt'd template) content of each reverted file.
func oldRealtimeChannelsProject(t *testing.T) (string, string, map[string]string) {
	t.Helper()
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
		filepath.Join(api, "realtime", "hub.go"):               realtimeHubChannelsHunks,
		filepath.Join(api, "realtime", "backplane.go"):         realtimeBackplaneChannelsHunks,
		filepath.Join(api, "handlers", "realtime.go"):          realtimeReadPumpHunks,
		filepath.Join(api, "services", "event_subscribers.go"): realtimeSubscribersChannelsHunks,
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
	for _, path := range realtimeChannelsAddedFiles(api) {
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
	}
	return root, api, want
}

// Every client the realtime files serve can subscribe to a channel, and asks
// again for its channels after a reconnect.
func TestRealtimeClientSubscribesToChannels(t *testing.T) {
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
			"export function subscribe(channel: string, handlers: ChannelHandlers): () => void;",
			"channels.forEach((_, channel) => send({ type: \"subscribe\", channel }));",
			"send({ type: \"unsubscribe\", channel });",
			"dispatchChannel(evt.channel, evt);",
		} {
			if !strings.Contains(string(lib), want) {
				t.Errorf("%s: lib/realtime.ts is missing %q", app, want)
			}
		}
		hook, err := os.ReadFile(filepath.Join(root, "apps", app, "hooks", "use-realtime.ts"))
		if err != nil {
			t.Fatalf("%s: %v", app, err)
		}
		if !strings.Contains(string(hook), "export function useChannel(") {
			t.Errorf("%s: hooks/use-realtime.ts has no useChannel", app)
		}
	}
}
