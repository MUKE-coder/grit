package realtime

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
	Type    string          `json:"type"`
	Channel string          `json:"channel"`
	Payload json.RawMessage `json:"payload"`
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
		"private-nobody.1":                   ErrChannelForbidden,
		"presence-nobody.1":                  ErrChannelForbidden,
		"invoices.42":                        ErrInvalidChannel,
		"private-":                           ErrInvalidChannel,
		"private-a b":                        ErrInvalidChannel,
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
	if !ok || msg.Type != "invoices.paid" || msg.Channel != "private-invoices.42" || string(msg.Payload) != `{"id":"42"}` {
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
		{`{"type":"subscribe","channel":"private-invoices.42"}`, "subscribed", ""},
		{`{"type":"subscribe","channel":"private-invoices.7"}`, "subscription_error", "FORBIDDEN"},
		{`{"type":"subscribe","channel":"no-prefix"}`, "subscription_error", "INVALID_CHANNEL"},
		{`{"type":"unsubscribe","channel":"private-invoices.42"}`, "unsubscribed", ""},
	} {
		hub.HandleClientMessage(c, []byte(step.send))
		msg, ok := nextMessage(t, c)
		if !ok || msg.Type != step.wantType {
			t.Fatalf("%s: got %q (ok=%v), want %s", step.send, msg.raw, ok, step.wantType)
		}
		if step.wantCode != "" && !strings.Contains(string(msg.Payload), `"code":"`+step.wantCode+`"`) {
			t.Errorf("%s: payload %s has no code %s", step.send, msg.Payload, step.wantCode)
		}
	}
	hub.HandleClientMessage(c, []byte("not json"))
	hub.HandleClientMessage(c, []byte(`{"type":"something-else"}`))
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
	if !ok || msg.raw != `{"type":"note.created","payload":null}` {
		t.Fatalf("got %q (ok=%v)", msg.raw, ok)
	}
}
