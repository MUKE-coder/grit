package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

type presenceMessage struct {
	Type    string          `json:"type"`
	Channel string          `json:"channel"`
	Payload json.RawMessage `json:"payload"`
}

func presenceClient(h *Hub, userID string) *Client {
	c := &Client{UserID: userID, Send: make(chan []byte, 64)}
	h.Register(c)
	return c
}

// awaitPresence reads c's messages until one of type kind arrives.
func awaitPresence(t *testing.T, c *Client, kind string, wait time.Duration) presenceMessage {
	t.Helper()
	deadline := time.After(wait)
	for {
		select {
		case raw, ok := <-c.Send:
			if !ok {
				t.Fatalf("%s: the connection closed while waiting for %s", c.UserID, kind)
			}
			var msg presenceMessage
			if err := json.Unmarshal(raw, &msg); err != nil {
				t.Fatalf("decode %s: %v", raw, err)
			}
			if msg.Type == kind {
				return msg
			}
		case <-deadline:
			t.Fatalf("%s: no %s within %s", c.UserID, kind, wait)
		}
	}
}

// noPresence fails if c receives a message of type kind within wait.
func noPresence(t *testing.T, c *Client, kind string, wait time.Duration) {
	t.Helper()
	deadline := time.After(wait)
	for {
		select {
		case raw := <-c.Send:
			if strings.Contains(string(raw), `"type":"`+kind+`"`) {
				t.Fatalf("%s: unexpected %s", c.UserID, raw)
			}
		case <-deadline:
			return
		}
	}
}

func presenceSubscribe(t *testing.T, h *Hub, c *Client, channel string) presenceMessage {
	t.Helper()
	h.HandleClientMessage(c, []byte(`{"type":"subscribe","channel":"`+channel+`"}`))
	awaitPresence(t, c, "subscribed", 2*time.Second)
	return awaitPresence(t, c, "presence.members", 2*time.Second)
}

func memberIDs(members []Member) string {
	ids := make([]string, 0, len(members))
	for _, m := range members {
		ids = append(ids, m.UserID)
	}
	sort.Strings(ids)
	return strings.Join(ids, ",")
}

func snapshotIDs(t *testing.T, msg presenceMessage) string {
	t.Helper()
	var snap presenceSnapshot
	if err := json.Unmarshal(msg.Payload, &snap); err != nil {
		t.Fatalf("decode %s: %v", msg.Payload, err)
	}
	return memberIDs(snap.Members)
}

func payloadUser(t *testing.T, msg presenceMessage) string {
	t.Helper()
	var m presenceLeft
	if err := json.Unmarshal(msg.Payload, &m); err != nil {
		t.Fatalf("decode %s: %v", msg.Payload, err)
	}
	return m.UserID
}

// awaitBackplanes waits until a publish on each hub reaches the other, so no
// presence event is lost to a Redis subscription still connecting.
func awaitBackplanes(t *testing.T, hubs ...*Hub) {
	t.Helper()
	for i, from := range hubs {
		to := hubs[(i+1)%len(hubs)]
		probe := presenceClient(to, fmt.Sprintf("probe-%d", i))
		channel := fmt.Sprintf("public-backplane-probe.%d", i)
		if err := to.Subscribe(probe, channel); err != nil {
			t.Fatal(err)
		}
		deadline := time.Now().Add(15 * time.Second)
		for arrived := false; !arrived; {
			if time.Now().After(deadline) {
				t.Fatalf("hub %d's publishes never reached hub %d", i, (i+1)%len(hubs))
			}
			from.Publish(channel, Event{Type: "probe"})
			select {
			case <-probe.Send:
				arrived = true
			case <-time.After(200 * time.Millisecond):
			}
		}
		to.Unregister(probe)
	}
}

// withPresenceAuthorizer registers an authorizer for one test and restores the registry after.
func withPresenceAuthorizer(t *testing.T, pattern string, authorize func(ChannelContext) bool) {
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

func TestPresenceListsJoinsAndLeaves(t *testing.T) {
	withPresenceAuthorizer(t, "presence-rooms.{id}", func(c ChannelContext) bool {
		c.SetInfo(map[string]string{"name": "User " + c.UserID})
		return c.UserID != "stranger"
	})
	hub := NewHub()
	const room = "presence-rooms.1"

	alice := presenceClient(hub, "alice")
	snap := presenceSubscribe(t, hub, alice, room)
	if got := snapshotIDs(t, snap); got != "alice" {
		t.Fatalf("alice's snapshot lists %q", got)
	}
	if !strings.Contains(string(snap.Payload), `"name":"User alice"`) {
		t.Errorf("the snapshot has no info from SetInfo: %s", snap.Payload)
	}

	bob := presenceClient(hub, "bob")
	if got := snapshotIDs(t, presenceSubscribe(t, hub, bob, room)); got != "alice,bob" {
		t.Fatalf("bob's snapshot lists %q", got)
	}
	if got := payloadUser(t, awaitPresence(t, alice, "presence.joined", time.Second)); got != "bob" {
		t.Fatalf("alice was told %q joined", got)
	}
	noPresence(t, bob, "presence.joined", 100*time.Millisecond)

	// A second tab is not a second member, and announces nothing.
	bobTab := presenceClient(hub, "bob")
	if got := snapshotIDs(t, presenceSubscribe(t, hub, bobTab, room)); got != "alice,bob" {
		t.Fatalf("bob's second tab lists %q", got)
	}
	noPresence(t, alice, "presence.joined", 100*time.Millisecond)

	// bob leaves with his last connection, not his first.
	hub.Unsubscribe(bobTab, room)
	noPresence(t, alice, "presence.left", 100*time.Millisecond)
	hub.Unregister(bob)
	if got := payloadUser(t, awaitPresence(t, alice, "presence.left", time.Second)); got != "bob" {
		t.Fatalf("alice was told %q left", got)
	}
	if got := memberIDs(hub.PresenceMembers(room)); got != "alice" {
		t.Fatalf("members after bob left: %q", got)
	}

	stranger := presenceClient(hub, "stranger")
	hub.HandleClientMessage(stranger, []byte(`{"type":"subscribe","channel":"`+room+`"}`))
	awaitPresence(t, stranger, "subscription_error", time.Second)
	if got := memberIDs(hub.PresenceMembers(room)); got != "alice" {
		t.Fatalf("a refused user is listed: %q", got)
	}
}

func TestPresenceInfoOverTheLimitIsLeftOut(t *testing.T) {
	withPresenceAuthorizer(t, "presence-big", func(c ChannelContext) bool {
		c.SetInfo(map[string]string{"bio": strings.Repeat("x", MaxPresenceInfoBytes)})
		return true
	})
	hub := NewHub()
	c := presenceClient(hub, "u1")
	presenceSubscribe(t, hub, c, "presence-big")
	members := hub.PresenceMembers("presence-big")
	if len(members) != 1 || members[0].Info != nil {
		t.Fatalf("members %+v, want u1 without info", members)
	}
}

// Two replicas sharing Redis list the same members, a closed socket leaves on
// both, and the members of a replica that dies leave once their entries expire.
func TestPresenceAcrossReplicas(t *testing.T) {
	url := os.Getenv("REALTIME_TEST_REDIS_URL")
	if url == "" {
		t.Skip("set REALTIME_TEST_REDIS_URL to a Redis this test may write to")
	}
	opts, err := redis.ParseURL(url)
	if err != nil {
		t.Fatal(err)
	}
	rdb := redis.NewClient(opts)
	namespace := fmt.Sprintf("grit:realtime:test:%d", time.Now().UnixNano())
	const room = "presence-rooms.1"
	t.Cleanup(func() {
		_ = rdb.Del(context.Background(), namespace+":presence:"+room).Err()
		_ = rdb.Close()
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		t.Fatalf("REALTIME_TEST_REDIS_URL: %v", err)
	}

	// Short enough to watch an expiry, long enough that a Redis slowed by a busy
	// machine does not expire a live member between heartbeats.
	previousTTL, previousBeat := PresenceTTL, PresenceHeartbeat
	PresenceTTL, PresenceHeartbeat = 4*time.Second, 400*time.Millisecond
	t.Cleanup(func() { PresenceTTL, PresenceHeartbeat = previousTTL, previousBeat })
	withPresenceAuthorizer(t, "presence-rooms.{id}", func(ChannelContext) bool { return true })

	a := NewHub(WithRedis(url, namespace))
	b := NewHub(WithRedis(url, namespace))
	t.Cleanup(func() { _ = a.Close() })
	awaitBackplanes(t, a, b)

	alice := presenceClient(a, "alice")
	presenceSubscribe(t, a, alice, room)
	bob := presenceClient(b, "bob")
	snap := snapshotIDs(t, presenceSubscribe(t, b, bob, room))
	if snap != "alice,bob" {
		t.Fatalf("bob's snapshot on replica B lists %q", snap)
	}
	if got := payloadUser(t, awaitPresence(t, alice, "presence.joined", 2*time.Second)); got != "bob" {
		t.Fatalf("alice on replica A was told %q joined", got)
	}
	onA, onB := memberIDs(a.PresenceMembers(room)), memberIDs(b.PresenceMembers(room))
	if onA != "alice,bob" || onB != onA {
		t.Fatalf("replica A lists %q, replica B lists %q", onA, onB)
	}
	t.Logf("two replicas: snapshot %q, replica A lists %q, replica B lists %q", snap, onA, onB)

	b.Unregister(bob)
	if got := payloadUser(t, awaitPresence(t, alice, "presence.left", 2*time.Second)); got != "bob" {
		t.Fatalf("alice was told %q left", got)
	}
	if got := memberIDs(a.PresenceMembers(room)); got != "alice" {
		t.Fatalf("after bob's socket closed replica A lists %q", got)
	}
	t.Logf("socket close: alice was told bob left, replica A lists %q", "alice")

	carol := presenceClient(b, "carol")
	presenceSubscribe(t, b, carol, room)
	awaitPresence(t, alice, "presence.joined", 2*time.Second)
	// Live replicas renew their entries: nobody expires while both run.
	noPresence(t, alice, "presence.left", PresenceTTL+2*PresenceHeartbeat)
	if got := memberIDs(a.PresenceMembers(room)); got != "alice,carol" {
		t.Fatalf("after %s with both replicas up replica A lists %q", PresenceTTL+2*PresenceHeartbeat, got)
	}

	// Replica B dies: no Unregister, no leave, its heartbeat just stops.
	died := time.Now()
	_ = b.Close()
	msg := awaitPresence(t, alice, "presence.left", PresenceTTL+4*PresenceHeartbeat+time.Second)
	if got := payloadUser(t, msg); got != "carol" {
		t.Fatalf("alice was told %q left", got)
	}
	if got := memberIDs(a.PresenceMembers(room)); got != "alice" {
		t.Fatalf("after replica B died replica A lists %q", got)
	}
	t.Logf("replica death: carol left %s after replica B stopped (TTL %s, heartbeat %s)",
		time.Since(died).Round(10*time.Millisecond), PresenceTTL, PresenceHeartbeat)
}
