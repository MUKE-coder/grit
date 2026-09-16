package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// repairRealtimeWhispers brings an existing project's realtime hub up to A4:
// client events between the subscribers of a channel, and counters the hub
// reports on /api/health.
//
// whispers.go and stats.go are new framework code and arrive whole. They only
// compile beside the channel index A2 added and the presence field A3 added,
// and every send in hub.go, channels.go and presence.go has to go through the
// counted offer for the numbers to mean anything, so all four files are checked
// first and none is touched unless all four still read as Grit wrote them. A
// project that has not taken A2 and A3 gets neither client events nor stats.
//
// The handler and routes.go follow on their own: without them the hub still
// relays client events between connections, it just reads 1 KB messages and
// reports nothing on /api/health.
func repairRealtimeWhispers(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	realtimeDir := filepath.Join(apiRoot, "internal", "realtime")
	hub := filepath.Join(realtimeDir, "hub.go")
	channels := filepath.Join(realtimeDir, "channels.go")
	presence := filepath.Join(realtimeDir, "presence.go")
	backplane := filepath.Join(realtimeDir, "backplane.go")
	core := []struct {
		path   string
		repair func(string) (string, []string, []string)
	}{
		{hub, repairRealtimeHubWhispersSource},
		{channels, repairRealtimeChannelsWhispersSource},
		{presence, repairRealtimePresenceWhispersSource},
		{backplane, repairRealtimeBackplaneWhispersSource},
	}
	for _, f := range core {
		if !fileExists(f.path) {
			return nil
		}
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	for _, f := range core {
		src, err := readNormalized(f.path)
		if err != nil {
			return err
		}
		if _, _, warnings := f.repair(src); len(warnings) > 0 {
			fmt.Printf("  ⚠ %s: %s\n", shownPath(root, f.path), warnings[0])
			fmt.Println("    Client events and hub stats were not added. Compare it with grit upgrade --diff and run grit upgrade again.")
			return nil
		}
	}
	for _, f := range core {
		if err := repairSourceFile(root, m, f.path, f.repair); err != nil {
			return err
		}
	}

	module := opts.Module()
	addFile := func(path, content string) error {
		if fileExists(path) {
			return nil
		}
		if err := writeFile(path, strings.ReplaceAll(content, "{{MODULE}}", module)); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
		fmt.Printf("  ✓ %s: added\n", shownPath(root, path))
		return nil
	}
	for _, f := range []struct{ path, content string }{
		{filepath.Join(realtimeDir, "stats.go"), apiRealtimeStatsGo()},
		{filepath.Join(realtimeDir, "whispers.go"), apiRealtimeWhispersGo()},
		{filepath.Join(realtimeDir, "whispers_test.go"), apiRealtimeWhispersTestGo()},
	} {
		if err := addFile(f.path, f.content); err != nil {
			return err
		}
	}

	// The socket's read limit has to fit a client event, and the greeting is
	// queued before the hub can close the connection underneath it.
	handler := filepath.Join(apiRoot, "internal", "handlers", "realtime.go")
	if fileExists(handler) {
		if err := repairSourceFile(root, m, handler, repairRealtimeHandlerWhispersSource); err != nil {
			return err
		}
		authTest := filepath.Join(apiRoot, "internal", "handlers", "auth_test.go")
		greetingTest := filepath.Join(apiRoot, "internal", "handlers", "realtime_greeting_test.go")
		if fileContains(handler, realtimeGreetingBlockNew) && fileContains(authTest, "func newTestAuthSvc(") &&
			fileContains(authTest, "func testCfg(") && !fileExists(greetingTest) {
			if err := addFile(greetingTest, apiRealtimeGreetingTestGo()); err != nil {
				return err
			}
		}
	}

	// The numbers, beside the event bus's, on the endpoint the admin's System
	// Health page already reads.
	if routes := filepath.Join(apiRoot, "internal", "routes", "routes.go"); fileExists(routes) {
		if err := repairSourceFile(root, m, routes, repairRealtimeHealthStatsSource); err != nil {
			return err
		}
	}
	return nil
}

func repairRealtimeHubWhispersSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "func (h *Hub) deliverLocal(") {
		return src, nil, nil
	}
	return applyRealtimeHunks(src, realtimeHubWhisperHunks,
		"the hub counts what it sends and drops, and a connection carries its client event limit",
		"is not the hub Grit wrote: give Client a clientEvents clientEventLimiter field and Hub a stats hubStats field, and send through h.offer in deliverLocal and broadcastLocal, or client events and /api/health's realtime numbers cannot be added")
}

func repairRealtimeChannelsWhispersSource(src string) (string, []string, []string) {
	return applyRealtimeHunks(src, realtimeChannelsWhisperHunks,
		"a connection can send client events to the other subscribers of a channel",
		"is not the channels.go Grit wrote: compare it with the one grit new writes to send through h.offer and to answer a client-event message in HandleClientMessage")
}

func repairRealtimePresenceWhispersSource(src string) (string, []string, []string) {
	return applyRealtimeHunks(src, realtimePresenceWhisperHunks,
		"presence messages are counted with the rest",
		"is not the presence.go Grit wrote: send through h.offer in publishPresence, or presence messages are left out of the hub's counters")
}

func repairRealtimeBackplaneWhispersSource(src string) (string, []string, []string) {
	return applyRealtimeHunks(src, realtimeBackplaneWhisperHunks,
		"an event that never reached the other replicas is counted",
		"is not the backplane Grit wrote: count h.stats.backplaneErrors where a publish fails or the queue is full, or /api/health reports no publish errors")
}

func repairRealtimeHandlerWhispersSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "RealtimeHandler") {
		return src, nil, nil
	}
	return applyRealtimeHunks(src, realtimeHandlerWhisperHunks,
		"the socket reads a client event, and the greeting is queued before the hub can close the connection",
		"is not the handler Grit wrote: raise wsMaxMessageSize to 2048 and queue the greeting before h.Hub.Admit, or a client event does not fit and a session revoked mid-handshake panics Connect")
}

// repairRealtimeHealthStatsSource puts the hub's counters on /api/health.
func repairRealtimeHealthStatsSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "realtimeHub :=") || strings.Contains(src, "realtimeHub.Stats()") {
		return src, nil, nil
	}
	if strings.Count(src, healthRealtimeOld) != 1 {
		return src, nil, []string{"does not report the event bus on /api/health the way Grit wrote it: add \"realtime\": realtimeHub.Stats() to the health payload, or the page shows no socket counts"}
	}
	return strings.Replace(src, healthRealtimeOld, healthRealtimeNew, 1),
		[]string{"/api/health reports this replica's sockets, and what the hub delivered, dropped and failed to publish"}, nil
}
