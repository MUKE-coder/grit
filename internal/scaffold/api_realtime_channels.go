package scaffold

// Realtime channels (A2): subscribe to a named stream, authorized per channel.
//
// Until now the hub could only address a user id or everyone, so a page showing
// one invoice had no way to hear about that invoice without every event for the
// user being pushed at it, and an app had no way to say "whoever may read this
// record" at all. channels.go adds subscribe and unsubscribe on the socket, a
// pattern registry of authorizers, and Hub.Publish across replicas. The hunk
// constants below are shared by the templates and the upgrade repair, so both
// write the same text.

// apiRealtimeChannelsGo emits internal/realtime/channels.go.
func apiRealtimeChannelsGo() string {
	return `package realtime

import (
	"encoding/json"
	"errors"
	"log"
	"regexp"
	"strings"
	"sync"
)

// Channels carry events to whoever subscribed to a name, instead of to a user id.
//
// A client sends
//
//	{"type":"subscribe","channel":"private-invoices.42"}
//	{"type":"unsubscribe","channel":"private-invoices.42"}
//
// and the hub answers with subscribed, subscription_error or unsubscribed, each
// naming the channel. Server code calls hub.Publish, and every subscriber on
// every replica receives
//
//	{"type":"invoices.paid","channel":"private-invoices.42","payload":{...}}
//
// The prefix decides who may subscribe:
//
//	public-*    any connection, no authorizer needed
//	private-*   only a user the channel's authorizer accepts
` + channelsPresenceDocNew + `//
// Register authorizers once at startup, before the router serves requests:
//
//	realtime.Channel("invoices.{id}", func(c realtime.ChannelContext) bool {
//	    return invoiceReadableBy(db, c.Param("id"), c.UserID)
//	})
//
// A pattern without a prefix covers the private- and presence- channels of that
// name; one with a prefix ("presence-rooms.{id}") covers only that kind. Each
// {name} matches one dot-separated segment. A private or presence channel that
// no pattern matches is refused, so a channel nobody wrote an authorizer for is
// closed rather than open.

// ChannelContext is what an authorizer decides on.
type ChannelContext struct {
	UserID  string
	Channel string            // the full name, "private-invoices.42"
` + channelsContextNew + `
// Param returns one parsed pattern parameter, "" when the pattern has none by that name.
func (c ChannelContext) Param(name string) string { return c.Params[name] }

// Why a subscribe was refused. The code in a subscription_error payload names which.
var (
	ErrInvalidChannel   = errors.New("channel names start with public-, private- or presence-, followed by up to 155 letters, digits or _ - = @ , . ;")
	ErrChannelForbidden = errors.New("not allowed to subscribe to this channel")
	ErrTooManyChannels  = errors.New("this connection holds the most channels it may")
)

// MaxChannelsPerConnection bounds the subscriptions one socket may hold, so a
// client cannot grow the hub's index without limit. Zero or less turns it off.
var MaxChannelsPerConnection = 100

// maxChannelName keeps the longest subscribe message near 200 bytes, well inside
// the socket's 1 KB read limit.
const maxChannelName = 164

var channelNameRe = regexp.MustCompile("^(public|private|presence)-[A-Za-z0-9_=@,.;-]+$")

func validChannel(name string) bool {
	return len(name) <= maxChannelName && channelNameRe.MatchString(name)
}

type channelPattern struct {
	source    string
	kind      string // "private-", "presence-", or "" for both
	segments  []string
	authorize func(ChannelContext) bool
}

var (
	channelsMu      sync.RWMutex
	channelPatterns []channelPattern
)

// Channel registers the authorizer for the private and presence channels that
// match pattern. Registering a pattern again replaces its authorizer. Patterns
// are tried in the order they were registered and the first match decides.
//
// It panics on a nil authorizer, a public- pattern or a malformed pattern: those
// are mistakes in startup code, found the first time the API starts.
func Channel(pattern string, authorize func(c ChannelContext) bool) {
	if authorize == nil {
		panic("realtime.Channel: nil authorizer for " + pattern)
	}
	if strings.HasPrefix(pattern, "public-") {
		panic("realtime.Channel: public channels need no authorizer: " + pattern)
	}
	kind, rest := "", pattern
	for _, prefix := range []string{"private-", "presence-"} {
		if strings.HasPrefix(pattern, prefix) {
			kind, rest = prefix, strings.TrimPrefix(pattern, prefix)
		}
	}
	segments := strings.Split(rest, ".")
	for _, segment := range segments {
		if segment == "" || segment == "{}" {
			panic("realtime.Channel: empty segment in pattern " + pattern)
		}
	}
	entry := channelPattern{source: pattern, kind: kind, segments: segments, authorize: authorize}

	channelsMu.Lock()
	defer channelsMu.Unlock()
	for i := range channelPatterns {
		if channelPatterns[i].source == pattern {
			channelPatterns[i] = entry
			return
		}
	}
	channelPatterns = append(channelPatterns, entry)
}

// matchChannel finds the authorizer for a private or presence channel.
func matchChannel(channel string) (func(ChannelContext) bool, map[string]string) {
	kind := "private-"
	if strings.HasPrefix(channel, "presence-") {
		kind = "presence-"
	}
	segments := strings.Split(strings.TrimPrefix(channel, kind), ".")

	channelsMu.RLock()
	defer channelsMu.RUnlock()
	for _, p := range channelPatterns {
		if (p.kind != "" && p.kind != kind) || len(p.segments) != len(segments) {
			continue
		}
		params := map[string]string{}
		matched := true
		for i, want := range p.segments {
			got := segments[i]
			if got == "" {
				matched = false
				break
			}
			if strings.HasPrefix(want, "{") && strings.HasSuffix(want, "}") {
				params[want[1:len(want)-1]] = got
				continue
			}
			if want != got {
				matched = false
				break
			}
		}
		if matched {
			return p.authorize, params
		}
	}
	return nil, nil
}

` + channelsAuthorizeNew + `
// channelIndex records subscriptions both ways, so a publish finds a channel's
// subscribers and a closing connection finds its channels without scanning
// everything. The hub's mu guards it; the zero value is ready to use.
type channelIndex struct {
	members map[string]map[*Client]struct{}
	joined  map[*Client]map[string]struct{}
}

func (x *channelIndex) add(c *Client, channel string) error {
	if _, ok := x.joined[c][channel]; ok {
		return nil
	}
	if MaxChannelsPerConnection > 0 && len(x.joined[c]) >= MaxChannelsPerConnection {
		return ErrTooManyChannels
	}
	if x.members == nil {
		x.members = make(map[string]map[*Client]struct{})
		x.joined = make(map[*Client]map[string]struct{})
	}
	if x.members[channel] == nil {
		x.members[channel] = make(map[*Client]struct{})
	}
	if x.joined[c] == nil {
		x.joined[c] = make(map[string]struct{})
	}
	x.members[channel][c] = struct{}{}
	x.joined[c][channel] = struct{}{}
	return nil
}

func (x *channelIndex) remove(c *Client, channel string) {
	if set, ok := x.members[channel]; ok {
		delete(set, c)
		if len(set) == 0 {
			delete(x.members, channel)
		}
	}
	if set, ok := x.joined[c]; ok {
		delete(set, channel)
		if len(set) == 0 {
			delete(x.joined, c)
		}
	}
}

// drop forgets every subscription a closing connection held.
func (x *channelIndex) drop(c *Client) {
	for channel := range x.joined[c] {
		x.remove(c, channel)
	}
}

` + channelsSubscribeNew + `
// Publish delivers evt to every subscriber of channel, on this process and,
// through the backplane, on every other.
//
// Nothing is checked at publish time: each subscriber passed the channel's
// authorizer when it joined. So name a channel for what it carries,
// "private-invoices.42" rather than "private-updates", and publish to it only
// what everyone the authorizer admits may see.
func (h *Hub) Publish(channel string, evt Event) {
	if !validChannel(channel) {
		log.Printf("[realtime] not publishing to the invalid channel %q", channel)
		return
	}
	bytes, err := json.Marshal(channelEnvelope{Type: evt.Type, Channel: channel, Payload: evt.Payload})
	if err != nil {
		log.Printf("[realtime] marshal: %v", err)
		return
	}
	// This node first, as SendToUsers does, so local subscribers do not wait on
	// Redis and are still served when it is down.
	h.deliverChannelLocal(channel, bytes)
	h.publish(fanout{Channel: channel, Event: bytes})
}

// deliverChannelLocal pushes an encoded event to this process's subscribers.
//
// The sends happen under the read lock. Unregister and DisconnectUser close a
// client's Send under the write lock, so no send here can land on a closed
// channel, and a non-blocking send keeps the lock short.
func (h *Hub) deliverChannelLocal(channel string, bytes []byte) {
	if len(bytes) == 0 {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.channels.members[channel] {
` + channelsDeliverSendNew + `	}
}

// channelEnvelope is an Event on the wire with the channel it arrived on.
type channelEnvelope struct {
	Type    string      ` + "`" + `json:"type"` + "`" + `
	Channel string      ` + "`" + `json:"channel"` + "`" + `
	Payload interface{} ` + "`" + `json:"payload"` + "`" + `
}

type clientMessage struct {
	Type    string ` + "`" + `json:"type"` + "`" + `
	Channel string ` + "`" + `json:"channel"` + "`" + `
}

type subscriptionError struct {
	Code    string ` + "`" + `json:"code"` + "`" + `
	Message string ` + "`" + `json:"message"` + "`" + `
}

func subscriptionErrorCode(err error) string {
	switch {
	case errors.Is(err, ErrInvalidChannel):
		return "INVALID_CHANNEL"
	case errors.Is(err, ErrTooManyChannels):
		return "TOO_MANY_CHANNELS"
	default:
		return "FORBIDDEN"
	}
}

` + channelsHandleDocNew + `func (h *Hub) HandleClientMessage(c *Client, raw []byte) {
	var msg clientMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		return
	}
	switch msg.Type {
	case "subscribe":
		if err := h.Subscribe(c, msg.Channel); err != nil {
			h.reply(c, "subscription_error", msg.Channel, subscriptionError{Code: subscriptionErrorCode(err), Message: err.Error()})
			return
		}
` + channelsSubscribedReplyNew + channelsClientEventCaseNew + `		h.reply(c, "unsubscribed", msg.Channel, nil)
	}
}

// reply answers one connection, and only while it is still registered: a
// closed connection's Send is closed too.
func (h *Hub) reply(c *Client, kind, channel string, payload interface{}) {
	bytes, err := json.Marshal(channelEnvelope{Type: kind, Channel: channel, Payload: payload})
	if err != nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	if _, open := h.clients[c.UserID][c]; !open {
		return
	}
` + channelsReplySendNew + `}
`
}

// apiRealtimeChannelsTestGo emits internal/realtime/channels_test.go. Its helpers
// are its own, so the file compiles beside any backplane_test.go or none.
func apiRealtimeChannelsTestGo() string {
	return `package realtime

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"
)

type channelBus struct {
	mu   sync.Mutex
	subs []func([]byte)
}

func (b *channelBus) Publish(ctx context.Context, msg []byte) error {
	b.mu.Lock()
	subs := append([]func([]byte){}, b.subs...)
	b.mu.Unlock()
	for _, deliver := range subs {
		deliver(msg)
	}
	return nil
}

func (b *channelBus) Subscribe(ctx context.Context, deliver func([]byte)) {
	b.mu.Lock()
	b.subs = append(b.subs, deliver)
	b.mu.Unlock()
	<-ctx.Done()
}

func (b *channelBus) Close() error { return nil }

type received struct {
	Type    string          ` + "`" + `json:"type"` + "`" + `
	Channel string          ` + "`" + `json:"channel"` + "`" + `
	Payload json.RawMessage ` + "`" + `json:"payload"` + "`" + `
	raw     string
}

func connected(h *Hub, userID string) *Client {
	c := &Client{UserID: userID, Send: make(chan []byte, 16)}
	h.Register(c)
	return c
}

func nextMessage(t *testing.T, c *Client) (received, bool) {
	t.Helper()
	select {
	case raw, ok := <-c.Send:
		if !ok {
			return received{}, false
		}
		var msg received
		if err := json.Unmarshal(raw, &msg); err != nil {
			t.Fatalf("decode %s: %v", raw, err)
		}
		msg.raw = string(raw)
		return msg, true
	case <-time.After(300 * time.Millisecond):
		return received{}, false
	}
}

// withChannel registers an authorizer for one test and restores the registry after.
func withChannel(t *testing.T, pattern string, authorize func(ChannelContext) bool) {
	t.Helper()
	channelsMu.RLock()
	saved := append([]channelPattern(nil), channelPatterns...)
	channelsMu.RUnlock()
	Channel(pattern, authorize)
	t.Cleanup(func() {
		channelsMu.Lock()
		channelPatterns = saved
		channelsMu.Unlock()
	})
}

func ownerOf42(c ChannelContext) bool { return c.UserID == "owner" && c.Param("id") == "42" }

func TestSubscribeRefusesAUserTheAuthorizerRejects(t *testing.T) {
	withChannel(t, "invoices.{id}", ownerOf42)
	hub := NewHub()
	intruder := connected(hub, "intruder")

	if err := hub.Subscribe(intruder, "private-invoices.42"); !errors.Is(err, ErrChannelForbidden) {
		t.Fatalf("an unauthorized user subscribed (err=%v)", err)
	}
	hub.Publish("private-invoices.42", Event{Type: "invoices.paid"})
	if msg, ok := nextMessage(t, intruder); ok {
		t.Fatalf("a refused subscriber received %s", msg.raw)
	}
}

func TestSubscribeRefusesChannelsWithNoAuthorizerOrABadName(t *testing.T) {
	hub := NewHub()
	c := connected(hub, "u1")
	cases := map[string]error{
		"private-nobody.1":                     ErrChannelForbidden,
		"presence-nobody.1":                    ErrChannelForbidden,
		"invoices.42":                          ErrInvalidChannel,
		"private-":                             ErrInvalidChannel,
		"private-a b":                          ErrInvalidChannel,
		"public-" + strings.Repeat("x", 200): ErrInvalidChannel,
	}
	for channel, want := range cases {
		if err := hub.Subscribe(c, channel); !errors.Is(err, want) {
			t.Errorf("%q: err=%v, want %v", channel, err, want)
		}
	}
}

func TestAuthorizedSubscriberReceivesAPublish(t *testing.T) {
	var seen ChannelContext
	withChannel(t, "invoices.{id}", func(c ChannelContext) bool {
		seen = c
		return ownerOf42(c)
	})
	hub := NewHub()
	owner := connected(hub, "owner")

	if err := hub.Subscribe(owner, "private-invoices.42"); err != nil {
		t.Fatalf("the owner was refused: %v", err)
	}
	if seen.Channel != "private-invoices.42" || seen.Param("id") != "42" || seen.UserID != "owner" {
		t.Errorf("the authorizer saw %+v", seen)
	}
	hub.Publish("private-invoices.42", Event{Type: "invoices.paid", Payload: map[string]string{"id": "42"}})
	msg, ok := nextMessage(t, owner)
	if !ok || msg.Type != "invoices.paid" || msg.Channel != "private-invoices.42" || string(msg.Payload) != ` + "`" + `{"id":"42"}` + "`" + ` {
		t.Fatalf("got %q (ok=%v)", msg.raw, ok)
	}
	hub.Publish("private-invoices.43", Event{Type: "invoices.paid"})
	if msg, ok := nextMessage(t, owner); ok {
		t.Fatalf("a publish to another channel arrived: %s", msg.raw)
	}
}

func TestPresenceChannelsAuthorizeLikePrivateOnes(t *testing.T) {
	withChannel(t, "rooms.{id}", func(c ChannelContext) bool { return c.UserID == "member" })
	withChannel(t, "presence-lobby", func(c ChannelContext) bool { return true })
	hub := NewHub()
	member, stranger := connected(hub, "member"), connected(hub, "stranger")

	if err := hub.Subscribe(member, "presence-rooms.1"); err != nil {
		t.Errorf("an authorized member was refused a presence channel: %v", err)
	}
	if err := hub.Subscribe(stranger, "presence-rooms.1"); !errors.Is(err, ErrChannelForbidden) {
		t.Errorf("a presence channel admitted a user its authorizer rejects (err=%v)", err)
	}
	if err := hub.Subscribe(stranger, "private-lobby"); !errors.Is(err, ErrChannelForbidden) {
		t.Errorf("a presence- pattern authorized a private- channel (err=%v)", err)
	}
}

func TestPublicChannelsNeedNoAuthorizer(t *testing.T) {
	hub := NewHub()
	c := connected(hub, "anyone")
	if err := hub.Subscribe(c, "public-announcements"); err != nil {
		t.Fatalf("a public channel was refused: %v", err)
	}
	hub.Publish("public-announcements", Event{Type: "release.published"})
	if msg, ok := nextMessage(t, c); !ok || msg.Type != "release.published" {
		t.Fatalf("got %q (ok=%v)", msg.raw, ok)
	}
}

func TestUnsubscribeStopsDelivery(t *testing.T) {
	hub := NewHub()
	c := connected(hub, "u1")
	if err := hub.Subscribe(c, "public-feed"); err != nil {
		t.Fatal(err)
	}
	hub.Unsubscribe(c, "public-feed")
	hub.Publish("public-feed", Event{Type: "feed.item"})
	if msg, ok := nextMessage(t, c); ok {
		t.Fatalf("delivery continued after unsubscribe: %s", msg.raw)
	}
}

// A publish on one replica reaches a subscriber on another, exactly once, and
// the publishing node does not deliver its own message twice.
func TestPublishReachesSubscribersOnEveryNode(t *testing.T) {
	bus := &channelBus{}
	a, b := NewHub(WithBackplane(bus)), NewHub(WithBackplane(bus))
	t.Cleanup(func() { _ = a.Close(); _ = b.Close() })
	time.Sleep(20 * time.Millisecond)

	onA := connected(a, "u1")
	if err := a.Subscribe(onA, "public-orders"); err != nil {
		t.Fatal(err)
	}
	b.Publish("public-orders", Event{Type: "orders.created"})
	if msg, ok := nextMessage(t, onA); !ok || msg.Channel != "public-orders" {
		t.Fatalf("a subscriber on node A got %q for a publish on node B (ok=%v)", msg.raw, ok)
	}
	a.Publish("public-orders", Event{Type: "orders.updated"})
	if _, ok := nextMessage(t, onA); !ok {
		t.Fatal("a local publish did not arrive")
	}
	if msg, ok := nextMessage(t, onA); ok {
		t.Fatalf("a publish arrived twice: %s", msg.raw)
	}
}

func TestClosingAConnectionDropsItsSubscriptions(t *testing.T) {
	hub := NewHub()
	first, second := connected(hub, "u1"), connected(hub, "u2")
	for _, c := range []*Client{first, second} {
		if err := hub.Subscribe(c, "public-feed"); err != nil {
			t.Fatal(err)
		}
	}
	// Unregister for both: DisconnectUser's backstop closes a real socket after
	// ten seconds, and these clients have none.
	hub.Unregister(first)
	hub.Unregister(second)

	hub.mu.RLock()
	defer hub.mu.RUnlock()
	if len(hub.channels.members) != 0 || len(hub.channels.joined) != 0 {
		t.Fatalf("closed connections are still indexed: %d channels, %d connections",
			len(hub.channels.members), len(hub.channels.joined))
	}
}

func TestSubscriptionsPerConnectionAreCapped(t *testing.T) {
	previous := MaxChannelsPerConnection
	MaxChannelsPerConnection = 2
	t.Cleanup(func() { MaxChannelsPerConnection = previous })
	hub := NewHub()
	c := connected(hub, "u1")
	for _, channel := range []string{"public-a", "public-b", "public-a"} {
		if err := hub.Subscribe(c, channel); err != nil {
			t.Fatalf("%s refused under the cap: %v", channel, err)
		}
	}
	if err := hub.Subscribe(c, "public-c"); !errors.Is(err, ErrTooManyChannels) {
		t.Fatalf("a subscription past the cap was accepted (err=%v)", err)
	}
}

func TestAPanickingAuthorizerRefuses(t *testing.T) {
	withChannel(t, "boom", func(ChannelContext) bool { panic("database gone") })
	hub := NewHub()
	if err := hub.Subscribe(connected(hub, "u1"), "private-boom"); !errors.Is(err, ErrChannelForbidden) {
		t.Fatalf("err=%v, want ErrChannelForbidden", err)
	}
}

func TestClientMessagesAreAnswered(t *testing.T) {
	withChannel(t, "invoices.{id}", ownerOf42)
	hub := NewHub()
	c := connected(hub, "owner")

	for _, step := range []struct{ send, wantType, wantCode string }{
		{` + "`" + `{"type":"subscribe","channel":"private-invoices.42"}` + "`" + `, "subscribed", ""},
		{` + "`" + `{"type":"subscribe","channel":"private-invoices.7"}` + "`" + `, "subscription_error", "FORBIDDEN"},
		{` + "`" + `{"type":"subscribe","channel":"no-prefix"}` + "`" + `, "subscription_error", "INVALID_CHANNEL"},
		{` + "`" + `{"type":"unsubscribe","channel":"private-invoices.42"}` + "`" + `, "unsubscribed", ""},
	} {
		hub.HandleClientMessage(c, []byte(step.send))
		msg, ok := nextMessage(t, c)
		if !ok || msg.Type != step.wantType {
			t.Fatalf("%s: got %q (ok=%v), want %s", step.send, msg.raw, ok, step.wantType)
		}
		if step.wantCode != "" && !strings.Contains(string(msg.Payload), ` + "`" + `"code":"` + "`" + `+step.wantCode+` + "`" + `"` + "`" + `) {
			t.Errorf("%s: payload %s has no code %s", step.send, msg.Payload, step.wantCode)
		}
	}
	hub.HandleClientMessage(c, []byte("not json"))
	hub.HandleClientMessage(c, []byte(` + "`" + `{"type":"something-else"}` + "`" + `))
	if msg, ok := nextMessage(t, c); ok {
		t.Fatalf("an unknown message was answered: %s", msg.raw)
	}
}

// Messages sent to a user carry no channel and still arrive, subscribed or not.
func TestSendToUsersIsUnchangedByChannels(t *testing.T) {
	hub := NewHub()
	c := connected(hub, "u1")
	if err := hub.Subscribe(c, "public-feed"); err != nil {
		t.Fatal(err)
	}
	hub.SendToUsers([]string{"u1"}, Event{Type: "note.created"})
	msg, ok := nextMessage(t, c)
	if !ok || msg.raw != ` + "`" + `{"type":"note.created","payload":null}` + "`" + ` {
		t.Fatalf("got %q (ok=%v)", msg.raw, ok)
	}
}
`
}

// apiRealtimeChannelsHandlerTestGo emits internal/handlers/realtime_channels_test.go:
// subscribe, publish and unsubscribe over a real socket.
func apiRealtimeChannelsHandlerTestGo() string {
	return `package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"

	"{{MODULE}}/internal/realtime"
)

func TestRealtime_ChannelsOverTheSocket(t *testing.T) {
	gin.SetMode(gin.TestMode)
	auth := newTestAuthSvc(testCfg())
	pair, err := auth.GenerateTokenPair("user-7", "seven@example.com", "USER")
	require.NoError(t, err)
	realtime.Channel("socket-test-notes.{id}", func(c realtime.ChannelContext) bool {
		return c.UserID == "user-7" && c.Param("id") == "1"
	})

	hub := realtime.NewHub()
	r := gin.New()
	r.GET("/api/ws", NewRealtimeHandler(hub, auth).Connect)
	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	header := http.Header{}
	header.Set("Authorization", "Bearer "+pair.AccessToken)
	conn, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(srv.URL, "http")+"/api/ws", header)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	type message struct {
		Type    string          ` + "`" + `json:"type"` + "`" + `
		Channel string          ` + "`" + `json:"channel"` + "`" + `
		Payload json.RawMessage ` + "`" + `json:"payload"` + "`" + `
	}
	read := func() (message, error) {
		_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		var msg message
		err := conn.ReadJSON(&msg)
		return msg, err
	}
	send := func(kind, channel string) {
		require.NoError(t, conn.WriteJSON(map[string]string{"type": kind, "channel": channel}))
	}

	msg, err := read()
	require.NoError(t, err)
	require.Equal(t, "system.connected", msg.Type)

	send("subscribe", "private-socket-test-notes.2")
	msg, err = read()
	require.NoError(t, err)
	require.Equal(t, "subscription_error", msg.Type, "a channel the authorizer rejects was subscribed")

	send("subscribe", "private-socket-test-notes.1")
	msg, err = read()
	require.NoError(t, err)
	require.Equal(t, "subscribed", msg.Type)

	hub.Publish("private-socket-test-notes.1", realtime.Event{Type: "notes.updated", Payload: map[string]string{"id": "1"}})
	msg, err = read()
	require.NoError(t, err)
	require.Equal(t, "notes.updated", msg.Type)
	require.Equal(t, "private-socket-test-notes.1", msg.Channel)

	send("unsubscribe", "private-socket-test-notes.1")
	msg, err = read()
	require.NoError(t, err)
	require.Equal(t, "unsubscribed", msg.Type)

	hub.Publish("private-socket-test-notes.1", realtime.Event{Type: "notes.updated"})
	_ = conn.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
	_, _, err = conn.ReadMessage()
	require.Error(t, err, "a publish arrived after unsubscribe")
}
`
}

// hub.go hunks. The Old text is exactly what Grit wrote before channels.

const hubChannelsFieldOld = "\tclients map[string]map[*Client]struct{} // userID -> set of connections\n"

const hubChannelsFieldNew = hubChannelsFieldOld + `
	// channels is who is subscribed to which channel. Guarded by mu. See
	// channels.go.
	channels channelIndex
`

const hubUnregisterOld = "\t\t\tdelete(set, c)\n\t\t\tclose(c.Send)\n"

const hubUnregisterNew = "\t\t\tdelete(set, c)\n\t\t\th.channels.drop(c)\n\t\t\tclose(c.Send)\n"

const hubDisconnectOld = "\t\tconns = append(conns, c.Conn)\n"

const hubDisconnectNew = hubDisconnectOld + "\t\th.channels.drop(c)\n"

// backplane.go hunks.

const backplaneFanoutOld = "\tEvent json.RawMessage `json:\"e,omitempty\"`\n"

const backplaneFanoutNew = backplaneFanoutOld + `
	// Channel sends Event to the subscribers of one channel. See Hub.Publish.
	Channel string ` + "`json:\"c,omitempty\"`" + `
`

const backplaneReceiveOld = "\tcase f.All:\n\t\th.broadcastLocal(f.Event)\n"

const backplaneReceiveNew = "\tcase f.Channel != \"\":\n\t\th.deliverChannelLocal(f.Channel, f.Event)\n" + backplaneReceiveOld

// handlers/realtime.go hunks.

const realtimeReadCommentOld = `// readPump pumps messages from the client → hub. We don't currently
// accept commands from clients (mutations go through the REST API), so
// this loop just services ping/pong and cleans up on disconnect.
`

const realtimeReadCommentNew = `// readPump pumps messages from the client to the hub. Mutations go through the
// REST API; what a client sends here is channel subscribe and unsubscribe, which
// the hub answers on the same socket. It also services ping/pong and cleans up
// on disconnect.
`

const realtimeReadLoopOld = "\t\tif _, _, err := c.Conn.ReadMessage(); err != nil {\n\t\t\treturn\n\t\t}\n"

const realtimeReadLoopNew = "\t\t_, msg, err := c.Conn.ReadMessage()\n\t\tif err != nil {\n\t\t\treturn\n\t\t}\n\t\thub.HandleClientMessage(c, msg)\n"

// services/event_subscribers.go hunks.

const subscribersAudienceOld = `var RealtimeAudience = func(e events.Event) []string {
	if e.Actor == "" {
		return nil
	}
	return []string{e.Actor}
}
`

const subscribersAudienceNew = subscribersAudienceOld + `
// RealtimeChannels sends a resource event to channels as well as to users.
//
// Nil by default, so resource events reach only RealtimeAudience. Set it when a
// page subscribes to a record rather than waiting for events addressed to its
// user:
//
//	services.RealtimeChannels = func(e events.Event) []string {
//	    if e.Resource == "invoices" {
//	        return []string{"private-invoices." + e.ID}
//	    }
//	    return nil
//	}
//
// Every subscriber of a returned channel receives the event, label included, so
// return a private- or presence- channel whose authorizer (realtime.Channel)
// admits only users allowed to read the record. A public- channel reaches any
// connection that asks for it.
var RealtimeChannels func(e events.Event) []string
`

const subscribersPushOld = `		audience := RealtimeAudience(e)
		if len(audience) == 0 {
			return nil
		}
		hub.SendToUsers(audience, realtime.Event{
			Type: e.Name,
			Payload: map[string]interface{}{
				"resource": e.Resource,
				"id":       e.ID,
				"label":    e.Label,
				"actor":    e.Actor,
				"at":       e.At,
			},
		})
		return nil
`

const subscribersPushNew = `		evt := realtime.Event{
			Type: e.Name,
			Payload: map[string]interface{}{
				"resource": e.Resource,
				"id":       e.ID,
				"label":    e.Label,
				"actor":    e.Actor,
				"at":       e.At,
			},
		}
		if audience := RealtimeAudience(e); len(audience) > 0 {
			hub.SendToUsers(audience, evt)
		}
		if RealtimeChannels != nil {
			for _, channel := range RealtimeChannels(e) {
				hub.Publish(channel, evt)
			}
		}
		return nil
`
