package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// repairRealtimePresence brings an existing project's realtime hub up to A3.
//
// First the send crash, which every hub with a backplane has: deliverLocal and
// broadcastLocal sent after releasing the read lock, so an Unregister in between
// closed Send first and the send panicked, taking the process down. That repair
// stands alone.
//
// Then presence. presence.go is new framework code and arrives whole, but it
// needs a field and two calls in hub.go and changes to channels.go, which upgrade
// writes only when it is missing. So both are checked first and neither is
// touched unless both still read as Grit wrote them. A project without
// channels.go has not taken the A2 repair and gets no presence.
func repairRealtimePresence(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	realtimeDir := filepath.Join(apiRoot, "internal", "realtime")
	hub := filepath.Join(realtimeDir, "hub.go")
	channels := filepath.Join(realtimeDir, "channels.go")
	if !fileExists(hub) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
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

	if err := repairSourceFile(root, m, hub, repairRealtimeHubSendsSource); err != nil {
		return err
	}
	hubSrc, err := readNormalized(hub)
	if err != nil {
		return err
	}
	if strings.Contains(hubSrc, hubLocalSendsNew) {
		if err := addFile(filepath.Join(realtimeDir, "hub_send_test.go"), apiRealtimeHubSendTestGo()); err != nil {
			return err
		}
	}

	if !fileExists(channels) {
		return nil
	}
	core := []struct {
		path   string
		repair func(string) (string, []string, []string)
	}{
		{hub, repairRealtimeHubPresenceSource},
		{channels, repairRealtimeChannelsPresenceSource},
	}
	for _, f := range core {
		src, err := readNormalized(f.path)
		if err != nil {
			return err
		}
		if _, _, warnings := f.repair(src); len(warnings) > 0 {
			fmt.Printf("  ⚠ %s: %s\n", shownPath(root, f.path), warnings[0])
			fmt.Println("    Presence channels were not added. Compare it with grit upgrade --diff and run grit upgrade again.")
			return nil
		}
	}
	for _, f := range core {
		if err := repairSourceFile(root, m, f.path, f.repair); err != nil {
			return err
		}
	}
	if err := addFile(filepath.Join(realtimeDir, "presence.go"), apiRealtimePresenceGo()); err != nil {
		return err
	}
	return addFile(filepath.Join(realtimeDir, "presence_test.go"), apiRealtimePresenceTestGo())
}

func readNormalized(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading %s: %w", path, err)
	}
	return strings.ReplaceAll(string(raw), "\r\n", "\n"), nil
}

// revertRealtimeHunks turns template text back into what Grit wrote before the
// hunks: each new text, present once, becomes its old text.
func revertRealtimeHunks(src string, hunks [][2]string) string {
	for _, h := range hunks {
		src = strings.Replace(src, h[1], h[0], 1)
	}
	return src
}

var realtimeHubSendsHunks = [][2]string{
	{hubLocalSendsOld, hubLocalSendsNew},
	{hubBackstopOld, hubBackstopNew},
}

// repairRealtimeHubSendsSource moves the local sends under the read lock, and
// stops DisconnectUser's backstop closing a nil Conn.
func repairRealtimeHubSendsSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "func (h *Hub) deliverLocal(") {
		return src, nil, nil
	}
	return applyRealtimeHunks(src, realtimeHubSendsHunks,
		"a message sent while a connection closes no longer panics the API",
		"does not send the way Grit wrote it: in deliverLocal and broadcastLocal, send while holding h.mu.RLock (with a non-blocking select), or a send racing a closing connection panics with \"send on closed channel\" and stops the API")
}

var realtimeHubPresenceHunks = [][2]string{
	{hubPresenceFieldOld, hubPresenceFieldNew},
	{hubPresenceStartOld, hubPresenceStartNew},
	{hubPresenceUnregisterOld, hubPresenceUnregisterNew},
}

var realtimeChannelsPresenceHunks = [][2]string{
	{channelsPresenceDocOld, channelsPresenceDocNew},
	{channelsContextOld, channelsContextNew},
	{channelsAuthorizeOld, channelsAuthorizeNew},
	{channelsSubscribeOld, channelsSubscribeNew},
	{channelsSubscribedReplyOld, channelsSubscribedReplyNew},
}

func repairRealtimeHubPresenceSource(src string) (string, []string, []string) {
	return applyRealtimeHunks(src, realtimeHubPresenceHunks,
		"the hub keeps presence channel members, and a closing connection leaves them",
		"is not the hub Grit wrote: give Hub a presence presenceSet field, call h.startPresence(ctx) after go h.publishLoop(ctx) in NewHub, and defer h.leaveAllPresence(c) before the lock in Unregister")
}

func repairRealtimeChannelsPresenceSource(src string) (string, []string, []string) {
	return applyRealtimeHunks(src, realtimeChannelsPresenceHunks,
		"presence-* channels list their members and announce joins and leaves",
		"is not the channels.go Grit wrote: compare it with the one grit new writes to add presence (SetInfo, joinPresence in Subscribe, leavePresence in Unsubscribe, replyPresenceMembers after subscribed)")
}
