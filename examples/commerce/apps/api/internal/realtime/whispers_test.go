package realtime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

type whisperMessage struct {
	Type    string          `json:"type"`
	Channel string          `json:"channel"`
	Payload json.RawMessage `json:"payload"`
}

func whisperClient(t *testing.T, h *Hub, userID string, channels ...string) *Client {
	t.Helper()
	c := &Client{UserID: userID, Send: make(chan []byte, 64)}
	h.Register(c)
	for _, channel := range channels {
		if err := h.Subscribe(c, channel); err != nil {
			t.Fatalf("%s subscribing to %s: %v", userID, channel, err)
		}
	}
	drainWhispers(c) // presence snapshots and joins
	return c
}

func drainWhispers(c *Client) {
	for {
		select {
		case <-c.Send:
		case <-time.After(50 * time.Millisecond):
			return
		}
	}
}

// collectWhispers reads c's messages for wait and returns them.
func collectWhispers(t *testing.T, c *Client, wait time.Duration) []whisperMessage {
	t.Helper()
	var out []whisperMessage
	deadline := time.After(wait)
	for {
		select {
		case raw := <-c.Send:
			var msg whisperMessage
			if err := json.Unmarshal(raw, &msg); err != nil {
				t.Fatalf("decode %s: %v", raw, err)
			}
			out = append(out, msg)
		case <-deadline:
			return out
		}
	}
}

func countType(msgs []whisperMessage, kind string) int {
	n := 0
	for _, m := range msgs {
		if m.Type == kind {
			n++
		}
	}
	return n
}

func whisper(h *Hub, c *Client, channel, event, payload string) {
	h.HandleClientMessage(c, []byte(fmt.Sprintf(`{"type":"client-event","channel":%q,"event":%q,"payload":%s}`, channel, event, payload)))
}

func withWhisperRooms(t *testing.T) {
	t.Helper()
	channelsMu.RLock()
	saved := append([]channelPattern(nil), channelPatterns...)
	channelsMu.RUnlock()
	Channel("rooms.{id}", func(c ChannelContext) bool { return true })
	t.Cleanup(func() {
		channelsMu.Lock()
		channelPatterns = saved
		channelsMu.Unlock()
	})
}

func TestWhisperReachesTheOtherSubscribersOnly(t *testing.T) {
	withWhisperRooms(t)
	hub := NewHub()
	for _, room := range []string{"private-rooms.1", "presence-rooms.1"} {
		alice := whisperClient(t, hub, "alice", room)
		aliceTab := whisperClient(t, hub, "alice", room)
		bob := whisperClient(t, hub, "bob", room)
		carol := whisperClient(t, hub, "carol", strings.Replace(room, ".1", ".2", 1))
		// Again, now everyone has joined: a presence channel tells the ones
		// already in it about each new member.
		for _, c := range []*Client{alice, aliceTab, bob, carol} {
			drainWhispers(c)
		}

		whisper(hub, alice, room, "typing", `{"typing":true}`)

		got := collectWhispers(t, bob, 200*time.Millisecond)
		if len(got) != 1 || got[0].Type != "client-event" || got[0].Channel != room {
			t.Fatalf("%s: bob got %+v", room, got)
		}
		var evt clientEvent
		if err := json.Unmarshal(got[0].Payload, &evt); err != nil {
			t.Fatal(err)
		}
		if evt.Event != "typing" || evt.UserID != "alice" || string(evt.Data) != `{"typing":true}` {
			t.Fatalf("%s: bob got %+v", room, evt)
		}
		if n := len(collectWhispers(t, aliceTab, 100*time.Millisecond)); n != 1 {
			t.Errorf("%s: alice's other tab got %d messages, want 1", room, n)
		}
		if msgs := collectWhispers(t, alice, 100*time.Millisecond); len(msgs) != 0 {
			t.Errorf("%s: the sender got its own whisper back: %+v", room, msgs)
		}
		if msgs := collectWhispers(t, carol, 100*time.Millisecond); len(msgs) != 0 {
			t.Errorf("%s: a subscriber of another channel got %+v", room, msgs)
		}
		for _, c := range []*Client{alice, aliceTab, bob, carol} {
			hub.Unregister(c)
		}
	}
	if s := hub.Stats(); s.ClientEvents != 2 {
		t.Errorf("stats count %d client events, want 2", s.ClientEvents)
	}
}

func TestWhispersAreRefusedWhereTheyDoNotBelong(t *testing.T) {
	withWhisperRooms(t)
	hub := NewHub()
	sender := whisperClient(t, hub, "alice", "public-lobby", "private-rooms.1")
	listener := whisperClient(t, hub, "bob", "public-lobby", "private-rooms.1", "private-rooms.9")

	for _, tc := range []struct{ channel, event, payload, code string }{
		{"public-lobby", "typing", "{}", "PUBLIC_CHANNEL"},
		{"private-rooms.9", "typing", "{}", "NOT_SUBSCRIBED"},
		{"rooms.1", "typing", "{}", "INVALID_CHANNEL"},
		{"private-rooms.1", "bad event", "{}", "INVALID_EVENT"},
		{"private-rooms.1", "typing", `"` + strings.Repeat("x", MaxClientEventBytes) + `"`, "PAYLOAD_TOO_LARGE"},
	} {
		whisper(hub, sender, tc.channel, tc.event, tc.payload)
		msgs := collectWhispers(t, sender, 100*time.Millisecond)
		if len(msgs) != 1 || msgs[0].Type != "client_event_error" || !strings.Contains(string(msgs[0].Payload), `"code":"`+tc.code+`"`) {
			t.Errorf("%s %q: sender got %+v, want client_event_error %s", tc.channel, tc.event, msgs, tc.code)
		}
	}
	if msgs := collectWhispers(t, listener, 100*time.Millisecond); len(msgs) != 0 {
		t.Errorf("a refused whisper was delivered: %+v", msgs)
	}
}

func TestWhispersOverTheRateLimitAreDropped(t *testing.T) {
	withWhisperRooms(t)
	previous := ClientEventsPerSecond
	ClientEventsPerSecond = 10
	t.Cleanup(func() { ClientEventsPerSecond = previous })
	hub := NewHub()
	sender := whisperClient(t, hub, "alice", "private-rooms.1")
	listener := whisperClient(t, hub, "bob", "private-rooms.1")

	const sent = 25
	for i := 0; i < sent; i++ {
		whisper(hub, sender, "private-rooms.1", "cursor", fmt.Sprintf(`{"x":%d}`, i))
	}
	delivered := countType(collectWhispers(t, listener, 200*time.Millisecond), "client-event")
	refusals := countType(collectWhispers(t, sender, 100*time.Millisecond), "client_event_error")
	s := hub.Stats()
	t.Logf("%d whispers in a burst: %d delivered, %d dropped, %d RATE_LIMITED replies", sent, delivered, s.ClientEventsRateLimited, refusals)
	if delivered != 10 || s.ClientEvents != 10 || s.ClientEventsRateLimited != sent-10 || refusals != 1 {
		t.Fatalf("delivered=%d relayed=%d limited=%d refusals=%d, want 10, 10, %d, 1",
			delivered, s.ClientEvents, s.ClientEventsRateLimited, sent-10, refusals)
	}
}

type whisperBus struct {
	mu   sync.Mutex
	subs []func([]byte)
	fail bool
}

func (b *whisperBus) Publish(ctx context.Context, msg []byte) error {
	b.mu.Lock()
	subs, fail := append([]func([]byte){}, b.subs...), b.fail
	b.mu.Unlock()
	if fail {
		return errors.New("backplane down")
	}
	for _, deliver := range subs {
		deliver(msg)
	}
	return nil
}

func (b *whisperBus) Subscribe(ctx context.Context, deliver func([]byte)) {
	b.mu.Lock()
	b.subs = append(b.subs, deliver)
	b.mu.Unlock()
	<-ctx.Done()
}

func (b *whisperBus) Close() error { return nil }

func TestWhisperCrossesReplicas(t *testing.T) {
	withWhisperRooms(t)
	bus := &whisperBus{}
	a, b := NewHub(WithBackplane(bus)), NewHub(WithBackplane(bus))
	t.Cleanup(func() { _ = a.Close(); _ = b.Close() })
	time.Sleep(20 * time.Millisecond)

	alice := whisperClient(t, a, "alice", "private-rooms.1")
	bobOnA := whisperClient(t, a, "bob", "private-rooms.1")
	carolOnB := whisperClient(t, b, "carol", "private-rooms.1")

	whisper(a, alice, "private-rooms.1", "typing", `{"typing":true}`)
	onA := countType(collectWhispers(t, bobOnA, 200*time.Millisecond), "client-event")
	onB := countType(collectWhispers(t, carolOnB, 200*time.Millisecond), "client-event")
	back := len(collectWhispers(t, alice, 100*time.Millisecond))
	t.Logf("one whisper: %d on the same replica, %d on the other, %d back to the sender", onA, onB, back)
	if onA != 1 || onB != 1 || back != 0 {
		t.Fatalf("same replica %d, other replica %d, sender %d; want 1, 1, 0", onA, onB, back)
	}
}

func TestStatsCountConnectionsAndDeliveries(t *testing.T) {
	log.SetOutput(io.Discard)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })
	bus := &whisperBus{fail: true}
	hub := NewHub(WithBackplane(bus))
	t.Cleanup(func() { _ = hub.Close() })

	u1a := &Client{UserID: "u1", Send: make(chan []byte, 8)}
	u1b := &Client{UserID: "u1", Send: make(chan []byte, 8)}
	slow := &Client{UserID: "u2", Send: make(chan []byte, 1)}
	for _, c := range []*Client{u1a, u1b, slow} {
		hub.Register(c)
	}
	for _, channel := range []string{"public-a", "public-b"} {
		if err := hub.Subscribe(u1a, channel); err != nil {
			t.Fatal(err)
		}
	}

	hub.SendToUsers([]string{"u1", "u2"}, Event{Type: "one"}) // 3 queued
	hub.SendToUsers([]string{"u2"}, Event{Type: "two"})       // slow's buffer is full: dropped

	deadline := time.Now().Add(3 * time.Second)
	for hub.Stats().BackplanePublishErrors < 2 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	s := hub.Stats()
	t.Logf("stats: %+v", s)
	if s.Connections != 3 || s.Users != 2 || s.Channels != 2 {
		t.Errorf("connections=%d users=%d channels=%d, want 3, 2, 2", s.Connections, s.Users, s.Channels)
	}
	if s.MessagesSent != 3 || s.MessagesDropped != 1 {
		t.Errorf("sent=%d dropped=%d, want 3 and 1", s.MessagesSent, s.MessagesDropped)
	}
	if !s.Backplane || s.BackplanePublishErrors != 2 {
		t.Errorf("backplane=%v publish errors=%d, want true and 2", s.Backplane, s.BackplanePublishErrors)
	}

	hub.Unregister(u1b)
	if s := hub.Stats(); s.Connections != 2 || s.Users != 2 {
		t.Errorf("after one socket closed: connections=%d users=%d, want 2 and 2", s.Connections, s.Users)
	}
}
