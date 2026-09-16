package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// deliverLocal and broadcastLocal sent after releasing the read lock, so a send
// racing Unregister panicked on a closed channel. Presence channels had no
// member list at all.
func TestRealtimePresenceShipsInTheTemplates(t *testing.T) {
	read := apiSource(t)

	hub := read("internal", "realtime", "hub.go")
	// hubLocalSendsCounted, not hubLocalSendsNew: A4 counts the two sends inside
	// it, and reverting those two hunks gives A3's text back.
	for _, want := range []string{hubLocalSendsCounted, hubBackstopNew, hubPresenceFieldNew, hubPresenceStartNew, hubPresenceUnregisterNew} {
		if strings.Count(hub, want) != 1 {
			t.Errorf("realtime/hub.go does not hold exactly one %q", want)
		}
	}
	if strings.Contains(hub, "targets = append(targets, c)") {
		t.Error("realtime/hub.go still collects targets to send to after releasing the lock")
	}

	channels := read("internal", "realtime", "channels.go")
	for _, h := range realtimeChannelsPresenceHunks {
		if strings.Count(channels, h[1]) != 1 {
			t.Errorf("realtime/channels.go does not hold exactly one %q", h[1])
		}
	}

	presence := read("internal", "realtime", "presence.go")
	for _, want := range []string{
		"func (c ChannelContext) SetInfo(info interface{})",
		"func (h *Hub) PresenceMembers(channel string) []Member",
		`h.reply(c, "presence.members", channel, presenceSnapshot{Members: h.PresenceMembers(channel)})`,
		`h.publishPresence(channel, "presence.joined", member, c)`,
		`h.publishPresence(channel, "presence.left", presenceLeft{UserID: c.UserID}, nil)`,
		"redis.call('PEXPIRE', KEYS[1], ttl)",
		"prefix: rb.channel + \":presence:\"",
	} {
		if !strings.Contains(presence, want) {
			t.Errorf("realtime/presence.go is missing %q", want)
		}
	}
	read("internal", "realtime", "presence_test.go")
	read("internal", "realtime", "hub_send_test.go")
}

func TestRealtimePresenceRepairsProduceTheTemplate(t *testing.T) {
	read := apiSource(t)
	// The files as v3.277.0 wrote them: A4 layers its own hunks on top, and
	// reverting those is what leaves A3's text to check.
	hub := revertRealtimeHunks(read("internal", "realtime", "hub.go"), realtimeHubWhisperHunks)
	channels := revertRealtimeHunks(read("internal", "realtime", "channels.go"), realtimeChannelsWhisperHunks)
	cases := []struct {
		name     string
		template string
		hunks    [][2]string
		repair   func(string) (string, []string, []string)
		edit     [2]string
	}{
		{"hub sends", hub, realtimeHubSendsHunks, repairRealtimeHubSendsSource,
			[2]string{"\ttargets := make([]*Client, 0)\n", "\ttargets := []*Client{}\n"}},
		{"hub presence", hub, realtimeHubPresenceHunks, repairRealtimeHubPresenceSource,
			[2]string{hubPresenceUnregisterOld, "// edited by hand\n"}},
		{"channels presence", channels, realtimeChannelsPresenceHunks, repairRealtimeChannelsPresenceSource,
			[2]string{channelsSubscribedReplyOld, "// edited by hand\n"}},
	}
	for _, tc := range cases {
		old := oldRealtimeFile(t, tc.template, tc.hunks)
		if got := revertRealtimeHunks(tc.template, tc.hunks); got != old {
			t.Errorf("%s: revertRealtimeHunks does not rebuild the old file", tc.name)
		}
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
		edited := strings.Replace(old, tc.edit[0], tc.edit[1], 1)
		if edited == old {
			t.Fatalf("%s: the edit did not apply", tc.name)
		}
		if got, _, warnings := tc.repair(edited); len(warnings) == 0 || got != edited {
			t.Errorf("%s: an edited file was rewritten or raised no warning", tc.name)
		}
	}
}

// oldRealtimePresenceProject scaffolds an API and turns hub.go and channels.go
// back into what v3.274.0 wrote, removing the files A3 adds.
func oldRealtimePresenceProject(t *testing.T) (string, string, map[string]string) {
	t.Helper()
	root := t.TempDir()
	opts := Options{ProjectName: "chat", Architecture: ArchTriple, Frontend: FrontendNext}
	if err := createDirectories(root, opts); err != nil {
		t.Fatalf("createDirectories: %v", err)
	}
	if err := writeAPIFiles(root, opts); err != nil {
		t.Fatalf("writeAPIFiles: %v", err)
	}
	dir := filepath.Join(root, "apps", "api", "internal", "realtime")
	files := map[string][2][][2]string{
		// A4's hunks first: this repair puts a file back to v3.277.0, which is
		// what it is checked against, and A4's own repair carries on from there.
		filepath.Join(dir, "hub.go"):      {realtimeHubWhisperHunks, append(append([][2]string{}, realtimeHubSendsHunks...), realtimeHubPresenceHunks...)},
		filepath.Join(dir, "channels.go"): {realtimeChannelsWhisperHunks, realtimeChannelsPresenceHunks},
	}
	want := map[string]string{}
	for path, hunks := range files {
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		beforeWhispers := revertRealtimeHunks(string(b), hunks[0])
		want[path] = gofmtSource(t, beforeWhispers)
		if err := os.WriteFile(path, []byte(oldRealtimeFile(t, beforeWhispers, hunks[1])), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, path := range realtimePresenceAddedFiles(dir) {
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
	}
	return root, dir, want
}

func realtimePresenceAddedFiles(dir string) []string {
	return []string{
		filepath.Join(dir, "presence.go"),
		filepath.Join(dir, "presence_test.go"),
		filepath.Join(dir, "hub_send_test.go"),
	}
}

var presenceTestOpts = Options{ProjectName: "chat", Architecture: ArchTriple, Frontend: FrontendNext}

// A v3.274.0 project, upgraded twice: the files come out as a fresh scaffold
// writes them, and the second run changes nothing.
func TestRepairRealtimePresenceOnAnOldProject(t *testing.T) {
	root, dir, want := oldRealtimePresenceProject(t)
	for run := 0; run < 2; run++ {
		if err := repairRealtimePresence(root, presenceTestOpts); err != nil {
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
		for _, path := range realtimePresenceAddedFiles(dir) {
			if !fileExists(path) {
				t.Errorf("run %d: %s was not added", run+1, filepath.Base(path))
			}
		}
	}
}

// A v3.273.0 project takes the channels repair and then this one, and ends with
// the channels.go a fresh scaffold writes.
func TestRepairRealtimePresenceAfterTheChannelsRepair(t *testing.T) {
	template := revertRealtimeHunks(apiSource(t)("internal", "realtime", "hub.go"), realtimeHubWhisperHunks)
	root, api, _ := oldRealtimeChannelsProject(t)
	dir := filepath.Join(api, "realtime")
	hub := filepath.Join(dir, "hub.go")
	// Presence first: its field hunk sits inside the channels field hunk.
	hunks := append(append([][2]string{}, realtimeHubSendsHunks...), realtimeHubPresenceHunks...)
	old := oldRealtimeFile(t, oldRealtimeFile(t, template, hunks), realtimeHubChannelsHunks)
	if err := os.WriteFile(hub, []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, path := range realtimePresenceAddedFiles(dir) {
		_ = os.Remove(path)
	}

	if err := repairRealtimeChannels(root, presenceTestOpts); err != nil {
		t.Fatal(err)
	}
	channels := filepath.Join(dir, "channels.go")
	if written, err := os.ReadFile(channels); err != nil || strings.Contains(string(written), "joinPresence") {
		t.Fatalf("the channels repair wrote presence code before the hub could take it (err=%v)", err)
	}
	if err := repairRealtimePresence(root, presenceTestOpts); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(channels)
	if err != nil {
		t.Fatal(err)
	}
	if gofmtSource(t, string(got)) != gofmtSource(t, revertRealtimeHunks(apiRealtimeChannelsGo(), realtimeChannelsWhisperHunks)) {
		t.Error("channels.go after both repairs differs from the template")
	}
	hubNow, err := os.ReadFile(hub)
	if err != nil {
		t.Fatal(err)
	}
	if gofmtSource(t, string(hubNow)) != gofmtSource(t, template) {
		t.Error("hub.go after both repairs differs from the template")
	}
}

// presence.go needs the hub's field and calls, so an edited hub means no
// presence: channels.go keeps compiling as it was. The send fix still applies.
func TestRepairRealtimePresenceLeavesAnEditedHubAlone(t *testing.T) {
	root, dir, _ := oldRealtimePresenceProject(t)
	hub := filepath.Join(dir, "hub.go")
	channels := filepath.Join(dir, "channels.go")
	b, err := os.ReadFile(hub)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(hub, []byte(strings.Replace(string(b), hubPresenceStartOld, "\t\tgo h.publishLoop(ctx) // edited\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(channels)
	if err != nil {
		t.Fatal(err)
	}
	if err := repairRealtimePresence(root, presenceTestOpts); err != nil {
		t.Fatal(err)
	}
	if fileExists(filepath.Join(dir, "presence.go")) {
		t.Error("presence.go was written beside a hub it does not compile with")
	}
	after, err := os.ReadFile(channels)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Error("channels.go was changed although the hub could not be")
	}
	if hubNow, err := os.ReadFile(hub); err != nil || !strings.Contains(string(hubNow), hubLocalSendsNew) {
		t.Errorf("the send fix was not applied to an otherwise edited hub (err=%v)", err)
	}
	if !fileExists(filepath.Join(dir, "hub_send_test.go")) {
		t.Error("hub_send_test.go was not added with the send fix")
	}
}

// Every client the realtime files serve tracks presence and has usePresence.
func TestRealtimeClientTracksPresence(t *testing.T) {
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
			"export function presenceMembers<Info = unknown>(channel: string): PresenceMember<Info>[]",
			"trackPresence(evt.channel, evt);",
			`evt.type === "presence.members"`,
			"presence.clear();",
		} {
			if !strings.Contains(string(lib), want) {
				t.Errorf("%s: lib/realtime.ts is missing %q", app, want)
			}
		}
		hook, err := os.ReadFile(filepath.Join(root, "apps", app, "hooks", "use-realtime.ts"))
		if err != nil {
			t.Fatalf("%s: %v", app, err)
		}
		if !strings.Contains(string(hook), "export function usePresence<Info = unknown>(channel: string | null | undefined): PresenceMember<Info>[]") {
			t.Errorf("%s: hooks/use-realtime.ts has no usePresence", app)
		}
	}
}

// A hub written by v3.281.0 or later counts its sends, so the fixed text the
// send repair looks for is not there verbatim. Upgrading such a project used to
// warn that its hub "does not send the way Grit wrote it", on every upgrade,
// about a hub that was exactly as Grit wrote it.
func TestHubSendsRepairAcceptsTheCountedHub(t *testing.T) {
	src := apiRealtimeHubGo()
	out, fixed, warnings := repairRealtimeHubSendsSource(src)
	if len(warnings) != 0 {
		t.Fatalf("the current hub template drew a warning: %v", warnings)
	}
	if len(fixed) != 0 || out != src {
		t.Fatalf("the current hub template was changed: %v", fixed)
	}
}
