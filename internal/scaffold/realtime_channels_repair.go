package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// repairRealtimeChannels brings an existing project's realtime hub up to A2:
// channel subscribe and unsubscribe on the socket, per-pattern authorizers,
// Hub.Publish across replicas, and services.RealtimeChannels for resource events.
//
// channels.go is new framework code and arrives whole. It needs a field in Hub
// and a Channel in the backplane fanout, so hub.go and backplane.go are checked
// together first and neither is touched unless both still read as Grit wrote
// them. The handler and event_subscribers.go are then repaired on their own:
// without them the hub compiles and simply has no subscribers.
func repairRealtimeChannels(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	realtimeDir := filepath.Join(apiRoot, "internal", "realtime")
	hub := filepath.Join(realtimeDir, "hub.go")
	backplane := filepath.Join(realtimeDir, "backplane.go")
	if !fileExists(hub) || !fileExists(backplane) {
		return nil
	}

	core := []struct {
		path   string
		repair func(string) (string, []string, []string)
	}{
		{hub, repairRealtimeHubChannelsSource},
		{backplane, repairRealtimeBackplaneChannelsSource},
	}
	for _, f := range core {
		raw, err := os.ReadFile(f.path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", f.path, err)
		}
		if _, _, warnings := f.repair(strings.ReplaceAll(string(raw), "\r\n", "\n")); len(warnings) > 0 {
			fmt.Printf("  ⚠ %s: %s\n", shownPath(root, f.path), warnings[0])
			fmt.Println("    Realtime channels were not added. Compare it with grit upgrade --diff and run grit upgrade again.")
			return nil
		}
	}

	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	for _, f := range core {
		if err := repairSourceFile(root, m, f.path, f.repair); err != nil {
			return err
		}
	}

	module := opts.Module()
	for _, f := range []struct{ path, content string }{
		{filepath.Join(realtimeDir, "channels.go"), apiRealtimeChannelsGo()},
		{filepath.Join(realtimeDir, "channels_test.go"), apiRealtimeChannelsTestGo()},
	} {
		if fileExists(f.path) {
			continue
		}
		if err := writeFile(f.path, strings.ReplaceAll(f.content, "{{MODULE}}", module)); err != nil {
			return fmt.Errorf("writing %s: %w", f.path, err)
		}
		fmt.Printf("  ✓ %s: added\n", shownPath(root, f.path))
	}

	handler := filepath.Join(apiRoot, "internal", "handlers", "realtime.go")
	if fileExists(handler) {
		if err := repairSourceFile(root, m, handler, repairRealtimeReadPumpSource); err != nil {
			return err
		}
		authTest := filepath.Join(apiRoot, "internal", "handlers", "auth_test.go")
		handlerTest := filepath.Join(apiRoot, "internal", "handlers", "realtime_channels_test.go")
		if fileContains(handler, realtimeReadLoopNew) && fileContains(authTest, "func newTestAuthSvc(") &&
			fileContains(authTest, "func testCfg(") && !fileExists(handlerTest) {
			if err := writeFile(handlerTest, strings.ReplaceAll(apiRealtimeChannelsHandlerTestGo(), "{{MODULE}}", module)); err != nil {
				return fmt.Errorf("writing %s: %w", handlerTest, err)
			}
			fmt.Printf("  ✓ %s: added\n", shownPath(root, handlerTest))
		}
	}

	if subs := filepath.Join(apiRoot, "internal", "services", "event_subscribers.go"); fileExists(subs) {
		if err := repairSourceFile(root, m, subs, repairRealtimeChannelsSubscribersSource); err != nil {
			return err
		}
	}
	return nil
}

// applyRealtimeHunks applies each hunk not already applied. All or nothing: one
// hunk whose old text is missing or repeated leaves the file alone with warning.
func applyRealtimeHunks(src string, hunks [][2]string, fixed, warning string) (string, []string, []string) {
	out := src
	for _, h := range hunks {
		if strings.Contains(out, h[1]) {
			continue
		}
		if strings.Count(out, h[0]) != 1 {
			return src, nil, []string{warning}
		}
		out = strings.Replace(out, h[0], h[1], 1)
	}
	if out == src {
		return src, nil, nil
	}
	return out, []string{fixed}, nil
}

var realtimeHubChannelsHunks = [][2]string{
	{hubChannelsFieldOld, hubChannelsFieldNew},
	{hubUnregisterOld, hubUnregisterNew},
	{hubDisconnectOld, hubDisconnectNew},
}

var realtimeBackplaneChannelsHunks = [][2]string{
	{backplaneFanoutOld, backplaneFanoutNew},
	{backplaneReceiveOld, backplaneReceiveNew},
}

func repairRealtimeHubChannelsSource(src string) (string, []string, []string) {
	return applyRealtimeHunks(src, realtimeHubChannelsHunks,
		"a closing connection leaves the channels it subscribed to",
		"is not the hub Grit wrote: give Hub a channels channelIndex field and call h.channels.drop(c) where Unregister and disconnectLocal close a client's Send")
}

func repairRealtimeBackplaneChannelsSource(src string) (string, []string, []string) {
	return applyRealtimeHunks(src, realtimeBackplaneChannelsHunks,
		"a channel publish reaches subscribers on every replica",
		"is not the backplane Grit wrote: add a Channel field to fanout and deliver it with h.deliverChannelLocal in receive")
}

// repairRealtimeReadPumpSource hands what a client sends to the hub. The comment
// above readPump is updated when it is still Grit's; the loop is what matters.
func repairRealtimeReadPumpSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "func readPump(") || strings.Contains(src, realtimeReadLoopNew) {
		return src, nil, nil
	}
	if !strings.Contains(src, "func readPump(hub *realtime.Hub, c *realtime.Client) {") || strings.Count(src, realtimeReadLoopOld) != 1 {
		return src, nil, []string{"does not read the socket the way Grit wrote it: pass each message readPump reads to hub.HandleClientMessage(c, msg), or clients cannot subscribe to channels"}
	}
	out := strings.Replace(src, realtimeReadLoopOld, realtimeReadLoopNew, 1)
	if strings.Count(out, realtimeReadCommentOld) == 1 {
		out = strings.Replace(out, realtimeReadCommentOld, realtimeReadCommentNew, 1)
	}
	return out, []string{"clients subscribe to channels on the socket"}, nil
}

var realtimeSubscribersChannelsHunks = [][2]string{
	{subscribersAudienceOld, subscribersAudienceNew},
	{subscribersPushOld, subscribersPushNew},
}

func repairRealtimeChannelsSubscribersSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "func registerRealtime(") || strings.Contains(src, "RealtimeChannels") {
		return src, nil, nil
	}
	return applyRealtimeHunks(src, realtimeSubscribersChannelsHunks,
		"resource events can go to channels through services.RealtimeChannels",
		"does not push resource events the way Grit wrote it: declare var RealtimeChannels func(e events.Event) []string and hub.Publish each channel it returns in registerRealtime, if you want resource events on channels")
}
