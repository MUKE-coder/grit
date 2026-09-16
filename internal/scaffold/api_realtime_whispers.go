package scaffold

import "strings"

// Realtime client events and stats (A4).
//
// Client events ("whispers") let one browser tell the other subscribers of a
// private or presence channel something that is not worth a REST round trip:
// typing, a cursor, "is viewing". Stats say how many sockets this replica holds
// and what the hub delivered, dropped and failed to publish, in /api/health.
//
// whispers.go and stats.go are new framework code. They need a field on Client
// and on Hub, a counted send in hub.go, channels.go and presence.go, and a
// counted publish error in backplane.go. Those files are written on upgrade
// only when missing, so the changes are hunks shared by the templates and the
// repair (realtime_whispers_repair.go), like the A1 to A3 hunks.

// withRealtimeHunks applies hunks to text a template splices in, so a constant
// A3 shares with its own repair can carry A4's changes without being rewritten.
// A hunk whose old text is not there exactly once is left out, and the template
// test that looks for the new text says so.
func withRealtimeHunks(src string, hunks [][2]string) string {
	for _, h := range hunks {
		if strings.Count(src, h[0]) == 1 {
			src = strings.Replace(src, h[0], h[1], 1)
		}
	}
	return src
}

// apiRealtimeWhispersGo emits internal/realtime/whispers.go.
func apiRealtimeWhispersGo() string {
	return `package realtime

import (
	"encoding/json"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Client events, or whispers, carry something one browser wants the others in a
// channel to know right now and nobody needs to keep: "Ada is typing", a cursor
// position, "is viewing this invoice". They go from socket to socket through the
// hub without a REST call, and the server stores nothing.
//
// A connection subscribed to a private or presence channel sends
//
//	{"type":"client-event","channel":"presence-rooms.1","event":"typing","payload":{"typing":true}}
//
// and every other connection subscribed to that channel, on every replica, gets
//
//	{"type":"client-event","channel":"presence-rooms.1","payload":{"event":"typing","user_id":"7","data":{"typing":true}}}
//
// user_id is set by the hub from the connection's token, so a receiver can trust
// who sent a whisper. It cannot trust what it says: data is whatever the sender's
// browser put there, so treat it as user input.
//
// The rules:
//
//   - Only private- and presence- channels. Everyone subscribed to one passed
//     its authorizer; a public- channel admits any connection that asks.
//   - Only a connection subscribed to the channel may whisper on it.
//   - Never back to the connection that sent it. The same user's other tabs do
//     receive it, with their own user_id, so they can ignore it if they like.
//   - At most ClientEventsPerSecond per connection. The rest are dropped.
//   - A payload of at most MaxClientEventBytes, and an event name of 1 to 64
//     letters, digits or _ . : -
//
// A refused whisper is answered, on the sender's socket only, with
//
//	{"type":"client_event_error","channel":"...","payload":{"code":"...","message":"..."}}
//
// where the code is one of INVALID_CHANNEL, PUBLIC_CHANNEL, NOT_SUBSCRIBED,
// INVALID_EVENT, PAYLOAD_TOO_LARGE or RATE_LIMITED. These travel on the socket
// and never as an HTTP response, so they are not in the HTTP error catalogue,
// the same as the subscription_error codes in channels.go. RATE_LIMITED is sent once per second at
// most, however many whispers that second drops.

// ClientEventsPerSecond bounds the whispers one connection may send in a second.
// Zero or less turns the limit off.
var ClientEventsPerSecond = 10

// MaxClientEventBytes bounds a whisper's payload as sent. The socket's read
// limit in handlers/realtime.go leaves room for this plus the channel and event
// names. Zero or less turns the check off, not the read limit.
var MaxClientEventBytes = 1024

var clientEventNameRe = regexp.MustCompile("^[A-Za-z0-9_.:-]{1,64}$")

// clientEventMessage is a whisper as a client sends it.
type clientEventMessage struct {
	Channel string          ` + "`" + `json:"channel"` + "`" + `
	Event   string          ` + "`" + `json:"event"` + "`" + `
	Payload json.RawMessage ` + "`" + `json:"payload"` + "`" + `
}

// clientEvent is a whisper as the other subscribers receive it.
type clientEvent struct {
	Event  string          ` + "`" + `json:"event"` + "`" + `
	UserID string          ` + "`" + `json:"user_id"` + "`" + `
	Data   json.RawMessage ` + "`" + `json:"data,omitempty"` + "`" + `
}

// clientEventLimiter counts one connection's whispers in one-second windows.
// The zero value is ready to use.
type clientEventLimiter struct {
	mu     sync.Mutex
	window time.Time
	count  int
	told   bool
}

// allow reports whether one more whisper fits in the current window, and for a
// refusal whether it is the window's first, which is the one worth answering.
func (l *clientEventLimiter) allow(now time.Time, perSecond int) (ok, firstRefusal bool) {
	if perSecond <= 0 {
		return true, false
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if now.Before(l.window) || now.Sub(l.window) >= time.Second {
		l.window, l.count, l.told = now, 0, false
	}
	if l.count < perSecond {
		l.count++
		return true, false
	}
	firstRefusal = !l.told
	l.told = true
	return false, firstRefusal
}

// relayClientEvent applies one client-event message from c.
func (h *Hub) relayClientEvent(c *Client, raw []byte) {
	var msg clientEventMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		return
	}
	// The limit comes first, so a flood of malformed whispers costs no more
	// than a flood of good ones.
	if ok, first := c.clientEvents.allow(time.Now(), ClientEventsPerSecond); !ok {
		h.stats.clientEventsLimited.Add(1)
		if first {
			h.refuseClientEvent(c, msg.Channel, "RATE_LIMITED", "too many client events: the rest of this second's are dropped")
		}
		return
	}
	switch {
	case !validChannel(msg.Channel):
		h.refuseClientEvent(c, msg.Channel, "INVALID_CHANNEL", ErrInvalidChannel.Error())
		return
	case strings.HasPrefix(msg.Channel, "public-"):
		h.refuseClientEvent(c, msg.Channel, "PUBLIC_CHANNEL", "client events go only to private- and presence- channels, whose subscribers were authorized")
		return
	case !clientEventNameRe.MatchString(msg.Event):
		h.refuseClientEvent(c, msg.Channel, "INVALID_EVENT", "an event name is 1 to 64 letters, digits or _ . : -")
		return
	case MaxClientEventBytes > 0 && len(msg.Payload) > MaxClientEventBytes:
		h.refuseClientEvent(c, msg.Channel, "PAYLOAD_TOO_LARGE", "a client event's payload is limited in size")
		return
	}
	bytes, err := json.Marshal(channelEnvelope{
		Type:    "client-event",
		Channel: msg.Channel,
		Payload: clientEvent{Event: msg.Event, UserID: c.UserID, Data: msg.Payload},
	})
	if err != nil {
		return
	}

	// The subscription check and the local sends share one read lock, so a
	// connection that unsubscribes or closes meanwhile is either not sent to or
	// still open when it is.
	h.mu.RLock()
	_, subscribed := h.channels.joined[c][msg.Channel]
	if subscribed {
		for other := range h.channels.members[msg.Channel] {
			if other != c {
				h.offer(other, bytes, "client event", msg.Channel)
			}
		}
	}
	h.mu.RUnlock()
	if !subscribed {
		h.refuseClientEvent(c, msg.Channel, "NOT_SUBSCRIBED", "subscribe to a channel before sending client events on it")
		return
	}
	h.stats.clientEvents.Add(1)
	// The sender is on this replica, so every receiver on the others is someone else.
	h.publish(fanout{Channel: msg.Channel, Event: bytes})
}

func (h *Hub) refuseClientEvent(c *Client, channel, code, message string) {
	h.reply(c, "client_event_error", channel, subscriptionError{Code: code, Message: message})
}
`
}

// apiRealtimeStatsGo emits internal/realtime/stats.go.
func apiRealtimeStatsGo() string {
	return `package realtime

import (
	"log"
	"sync/atomic"
)

// Stats is what one hub reports about itself, for /api/health and the admin's
// System Health page.
//
// Connections, Users and Channels are this process's own: with several replicas
// each reports its share. The counters run from when the process started, so a
// number that keeps rising is the signal, not its size.
type Stats struct {
	// Connections is the open sockets this process holds.
	Connections int ` + "`" + `json:"connections"` + "`" + `
	// Users is how many distinct users those sockets belong to.
	Users int ` + "`" + `json:"users"` + "`" + `
	// Channels is how many channels have at least one subscriber here.
	Channels int ` + "`" + `json:"channels"` + "`" + `

	// MessagesSent counts messages queued for a connection: events, channel
	// publishes, presence, whispers and replies.
	MessagesSent uint64 ` + "`" + `json:"messages_sent"` + "`" + `
	// MessagesDropped counts messages not queued because the connection's send
	// buffer was full. A slow client resyncs on its next REST call; a count that
	// climbs steadily means clients cannot keep up.
	MessagesDropped uint64 ` + "`" + `json:"messages_dropped"` + "`" + `

	// ClientEvents counts whispers relayed, and ClientEventsRateLimited the ones
	// dropped for going over ClientEventsPerSecond.
	ClientEvents            uint64 ` + "`" + `json:"client_events"` + "`" + `
	ClientEventsRateLimited uint64 ` + "`" + `json:"client_events_rate_limited"` + "`" + `

	// Backplane says whether this hub fans out to other replicas.
	// BackplanePublishErrors counts messages that did not reach them: a publish
	// that failed, or one dropped because the publish queue was full.
	Backplane              bool   ` + "`" + `json:"backplane"` + "`" + `
	BackplanePublishErrors uint64 ` + "`" + `json:"backplane_publish_errors"` + "`" + `
}

// hubStats holds a hub's counters. The zero value is ready to use.
type hubStats struct {
	sent                atomic.Uint64
	dropped             atomic.Uint64
	clientEvents        atomic.Uint64
	clientEventsLimited atomic.Uint64
	backplaneErrors     atomic.Uint64
}

// Stats reports this hub's connections and counters.
func (h *Hub) Stats() Stats {
	h.mu.RLock()
	connections := 0
	for _, set := range h.clients {
		connections += len(set)
	}
	s := Stats{
		Connections: connections,
		Users:       len(h.clients),
		Channels:    len(h.channels.members),
	}
	h.mu.RUnlock()

	s.MessagesSent = h.stats.sent.Load()
	s.MessagesDropped = h.stats.dropped.Load()
	s.ClientEvents = h.stats.clientEvents.Load()
	s.ClientEventsRateLimited = h.stats.clientEventsLimited.Load()
	s.Backplane = h.backplane != nil
	s.BackplanePublishErrors = h.stats.backplaneErrors.Load()
	return s
}

// offer queues one encoded message for c without blocking, and counts it sent
// or dropped. Call it holding h.mu, read or write: Unregister and DisconnectUser
// close Send under the write lock, so a send outside it can hit a closed
// channel and panic. what and channel only name the message in the drop log.
func (h *Hub) offer(c *Client, bytes []byte, what, channel string) bool {
	select {
	case c.Send <- bytes:
		h.stats.sent.Add(1)
		return true
	default:
		h.stats.dropped.Add(1)
		if channel != "" {
			log.Printf("[realtime] dropping a %s on %s for slow client user=%s", what, channel, c.UserID)
		} else {
			log.Printf("[realtime] dropping a %s for slow client user=%s", what, c.UserID)
		}
		return false
	}
}
`
}

// apiRealtimeWhispersTestGo emits internal/realtime/whispers_test.go. Its
// helpers are its own, so the file compiles beside any other test file.
func apiRealtimeWhispersTestGo() string {
	return `package realtime

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
	Type    string          ` + "`" + `json:"type"` + "`" + `
	Channel string          ` + "`" + `json:"channel"` + "`" + `
	Payload json.RawMessage ` + "`" + `json:"payload"` + "`" + `
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
	h.HandleClientMessage(c, []byte(fmt.Sprintf(` + "`" + `{"type":"client-event","channel":%q,"event":%q,"payload":%s}` + "`" + `, channel, event, payload)))
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

		whisper(hub, alice, room, "typing", ` + "`" + `{"typing":true}` + "`" + `)

		got := collectWhispers(t, bob, 200*time.Millisecond)
		if len(got) != 1 || got[0].Type != "client-event" || got[0].Channel != room {
			t.Fatalf("%s: bob got %+v", room, got)
		}
		var evt clientEvent
		if err := json.Unmarshal(got[0].Payload, &evt); err != nil {
			t.Fatal(err)
		}
		if evt.Event != "typing" || evt.UserID != "alice" || string(evt.Data) != ` + "`" + `{"typing":true}` + "`" + ` {
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
		{"private-rooms.1", "typing", ` + "`" + `"` + "`" + ` + strings.Repeat("x", MaxClientEventBytes) + ` + "`" + `"` + "`" + `, "PAYLOAD_TOO_LARGE"},
	} {
		whisper(hub, sender, tc.channel, tc.event, tc.payload)
		msgs := collectWhispers(t, sender, 100*time.Millisecond)
		if len(msgs) != 1 || msgs[0].Type != "client_event_error" || !strings.Contains(string(msgs[0].Payload), ` + "`" + `"code":"` + "`" + `+tc.code+` + "`" + `"` + "`" + `) {
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
		whisper(hub, sender, "private-rooms.1", "cursor", fmt.Sprintf(` + "`" + `{"x":%d}` + "`" + `, i))
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

	whisper(a, alice, "private-rooms.1", "typing", ` + "`" + `{"typing":true}` + "`" + `)
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
`
}

// apiRealtimeGreetingTestGo emits internal/handlers/realtime_greeting_test.go.
func apiRealtimeGreetingTestGo() string {
	return `package handlers

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"

	"{{MODULE}}/internal/realtime"
)

// disconnectOnRegister signs the user out the moment the hub logs that it
// admitted their socket, which is the instant Connect used to send its greeting.
type disconnectOnRegister struct {
	hub *realtime.Hub
	mu  sync.Mutex
	buf bytes.Buffer
}

func (w *disconnectOnRegister) Write(p []byte) (int, error) {
	if i := bytes.Index(p, []byte("client registered user=")); i >= 0 {
		user := strings.Fields(string(p[i+len("client registered user="):]))[0]
		go w.hub.DisconnectUser(user)
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.Write(p)
}

// A session revoked while a socket is being admitted must not panic Connect.
//
// Connect sent its greeting to client.Send after Admit. From Admit on, the hub
// holds the client, and DisconnectUser closes Send; a send on a closed channel
// panics. The greeting is now queued before Admit, while nothing else can reach
// the client.
func TestRealtime_GreetingDoesNotRaceADisconnect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	realtime.SetConnectionLimits(0, 0)
	t.Cleanup(func() { realtime.SetConnectionLimits(10, 10000) })
	auth := newTestAuthSvc(testCfg())
	pair, err := auth.GenerateTokenPair("user-g", "greeting@example.com", "USER")
	require.NoError(t, err)

	hub := realtime.NewHub()
	log.SetOutput(&disconnectOnRegister{hub: hub})
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	var panics atomic.Int64
	r := gin.New()
	r.Use(func(c *gin.Context) {
		defer func() {
			if recover() != nil {
				panics.Add(1)
			}
		}()
		c.Next()
	})
	r.GET("/api/ws", NewRealtimeHandler(hub, auth).Connect)
	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	url := "ws" + strings.TrimPrefix(srv.URL, "http") + "/api/ws"
	header := http.Header{}
	header.Set("Authorization", "Bearer "+pair.AccessToken)
	const sockets = 300
	for i := 0; i < sockets; i++ {
		conn, _, err := websocket.DefaultDialer.Dial(url, header)
		require.NoError(t, err)
		// Read until the server closes it: the greeting, then the sign-out. A
		// Connect that panicked leaves the socket open, hence the short deadline.
		_ = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
		_ = conn.Close()
	}
	t.Logf("%d sockets signed out as they were admitted: %d panics", sockets, panics.Load())
	require.Zero(t, panics.Load(), "Connect panicked sending to a connection DisconnectUser had closed")
}
`
}

// hub.go hunks. Each Old is exactly what Grit wrote in v3.277.0.

const hubClientEventsFieldOld = "\tExpiresAt time.Time\n}\n"

const hubClientEventsFieldNew = `	ExpiresAt time.Time

	// clientEvents limits the whispers this connection may send. See whispers.go.
	clientEvents clientEventLimiter
}
`

const hubStatsFieldOld = "\tpresence presenceSet\n"

// hubStatsFieldAdd is what A4 adds after A3's presence field, so hub.go's
// template can splice it on the end of hubPresenceFieldAdd and leave A3's
// constant a substring of the result.
const hubStatsFieldAdd = `
	// stats counts what the hub delivered, dropped and failed to publish. See
	// stats.go.
	stats hubStats
`

const hubStatsFieldNew = hubStatsFieldOld + hubStatsFieldAdd

const hubDeliverSendOld = `			select {
			case c.Send <- bytes:
			default:
				log.Printf("[realtime] dropping message for slow client user=%s", c.UserID)
			}
`

const hubDeliverSendNew = "\t\t\th.offer(c, bytes, \"message\", \"\")\n"

const hubBroadcastSendOld = `			select {
			case c.Send <- bytes:
			default:
				log.Printf("[realtime] dropping broadcast for slow client user=%s", c.UserID)
			}
`

const hubBroadcastSendNew = "\t\t\th.offer(c, bytes, \"broadcast\", \"\")\n"

// realtimeHubLocalSendHunks are the two sends inside the deliverLocal and
// broadcastLocal that A3 wrote (hubLocalSendsNew), kept apart from the field
// hunks so hub.go's template can splice the counted text without breaking A3.
var realtimeHubLocalSendHunks = [][2]string{
	{hubDeliverSendOld, hubDeliverSendNew},
	{hubBroadcastSendOld, hubBroadcastSendNew},
}

// hubLocalSendsCounted is A3's deliverLocal and broadcastLocal with their sends
// counted. hub.go's template splices this in place of hubLocalSendsNew; revert
// these two hunks and A3's constant is back, so A3's repair still applies to a
// project that has not taken A4 yet.
var hubLocalSendsCounted = withRealtimeHunks(hubLocalSendsNew, realtimeHubLocalSendHunks)

var realtimeHubWhisperHunks = append([][2]string{
	{hubClientEventsFieldOld, hubClientEventsFieldNew},
	{hubStatsFieldOld, hubStatsFieldNew},
}, realtimeHubLocalSendHunks...)

// backplane.go hunks.

const backplaneBufferFullOld = "\t\tlog.Printf(\"[realtime] backplane buffer full, dropping an event\")\n"

const backplaneBufferFullNew = "\t\th.stats.backplaneErrors.Add(1)\n" + backplaneBufferFullOld

const backplanePublishErrOld = "\t\t\t\tlog.Printf(\"[realtime] backplane publish: %v\", err)\n"

const backplanePublishErrNew = "\t\t\t\th.stats.backplaneErrors.Add(1)\n" + backplanePublishErrOld

var realtimeBackplaneWhisperHunks = [][2]string{
	{backplaneBufferFullOld, backplaneBufferFullNew},
	{backplanePublishErrOld, backplanePublishErrNew},
}

// channels.go hunks.

const channelsDeliverSendOld = `		select {
		case c.Send <- bytes:
		default:
			log.Printf("[realtime] dropping a %s message for slow client user=%s", channel, c.UserID)
		}
`

const channelsDeliverSendNew = "\t\th.offer(c, bytes, \"message\", channel)\n"

const channelsReplySendOld = `	select {
	case c.Send <- bytes:
	default:
		log.Printf("[realtime] dropping a %s reply for slow client user=%s", kind, c.UserID)
	}
`

const channelsReplySendNew = "\th.offer(c, bytes, kind+\" reply\", channel)\n"

const channelsHandleDocOld = `// HandleClientMessage applies one message a client sent on the socket:
// subscribe or unsubscribe. Anything else is ignored, so a client newer than
// the server does not break the connection.
`

const channelsHandleDocNew = `// HandleClientMessage applies one message a client sent on the socket:
// subscribe, unsubscribe, or a client event (whispers.go). Anything else is
// ignored, so a client newer than the server does not break the connection.
`

const channelsClientEventCaseOld = "\tcase \"unsubscribe\":\n\t\th.Unsubscribe(c, msg.Channel)\n"

const channelsClientEventCaseNew = "\tcase \"client-event\":\n\t\th.relayClientEvent(c, raw)\n" + channelsClientEventCaseOld

var realtimeChannelsWhisperHunks = [][2]string{
	{channelsDeliverSendOld, channelsDeliverSendNew},
	{channelsReplySendOld, channelsReplySendNew},
	{channelsHandleDocOld, channelsHandleDocNew},
	{channelsClientEventCaseOld, channelsClientEventCaseNew},
}

// presence.go hunks.

const presencePublishSendOld = `		select {
		case c.Send <- bytes:
		default:
			log.Printf("[realtime] dropping a %s message for slow client user=%s", kind, c.UserID)
		}
`

const presencePublishSendNew = "\t\th.offer(c, bytes, kind, channel)\n"

var realtimePresenceWhisperHunks = [][2]string{
	{presencePublishSendOld, presencePublishSendNew},
}

// handlers/realtime.go hunks.

const realtimeReadLimitOld = "\twsMaxMessageSize = 1024 // we don't expect clients to send anything large\n"

const realtimeReadLimitNew = "\twsMaxMessageSize = 2048 // subscribe messages, and client events with a payload of up to realtime.MaxClientEventBytes\n"

const realtimeGreetingBlockOld = `	// Greeting so the client knows the link is live.
	greeting, _ := json.Marshal(realtime.Event{
		Type:    "system.connected",
		Payload: gin.H{"user_id": claims.UserID},
	})
	select {
	case client.Send <- greeting:
	default:
	}
`

const realtimeGreetingOld = realtimeRegisterNew + "\n" + realtimeGreetingBlockOld

const realtimeGreetingBlockNew = `	// Greeting so the client knows the link is live. Queued before Admit, while
	// nothing else can reach this client: once the hub holds it, DisconnectUser
	// may close Send at any moment, and a send on a closed channel panics.
	// writePump starts after Admit, so the greeting still arrives only once the
	// socket is admitted.
	greeting, _ := json.Marshal(realtime.Event{
		Type:    "system.connected",
		Payload: gin.H{"user_id": claims.UserID},
	})
	select {
	case client.Send <- greeting:
	default:
	}
`

const realtimeGreetingNew = realtimeGreetingBlockNew + "\n" + realtimeRegisterNew

var realtimeHandlerWhisperHunks = [][2]string{
	{realtimeReadLimitOld, realtimeReadLimitNew},
	{realtimeGreetingOld, realtimeGreetingNew},
}

// routes.go hunk: the hub's stats beside the event bus in /api/health.

const healthRealtimeOld = "\t\t\t\"events\": eventBusStatus(),\n"

const healthRealtimeNew = healthRealtimeOld + `
			// This replica's sockets, and what its hub delivered, dropped and
			// failed to publish to the others.
			"realtime": realtimeHub.Stats(),
`
