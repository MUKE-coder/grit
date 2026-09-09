package scaffold

// apiRealtimeBackplaneTestGo emits internal/realtime/backplane_test.go.
//
// The interesting behaviour is between two processes, which no amount of
// reading one process's code shows you. These wire two Hubs through an
// in-memory backplane so the cross-node paths are exercised without Redis, and
// they ship with the project so a change to the hub has to keep them true.
func apiRealtimeBackplaneTestGo() string {
	return `package realtime

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"
)

// memBackplane joins hubs in one process the way Redis joins them across
// several: everything published reaches every subscriber, including the
// publisher, which is what makes the origin check worth testing.
type memBackplane struct {
	mu   sync.Mutex
	subs []func([]byte)
}

func (m *memBackplane) Publish(ctx context.Context, msg []byte) error {
	m.mu.Lock()
	subs := append([]func([]byte){}, m.subs...)
	m.mu.Unlock()
	for _, s := range subs {
		s(msg)
	}
	return nil
}

func (m *memBackplane) Subscribe(ctx context.Context, deliver func([]byte)) {
	m.mu.Lock()
	m.subs = append(m.subs, deliver)
	m.mu.Unlock()
	<-ctx.Done()
}

func (m *memBackplane) Close() error { return nil }

// joinedHubs returns two hubs on one backplane, as two API replicas would be.
func joinedHubs(t *testing.T) (*Hub, *Hub) {
	t.Helper()
	bp := &memBackplane{}
	a := NewHub(WithBackplane(bp))
	b := NewHub(WithBackplane(bp))
	t.Cleanup(func() { _ = a.Close(); _ = b.Close() })
	// Subscribe runs in a goroutine; give both a moment to register.
	time.Sleep(20 * time.Millisecond)
	return a, b
}

// fakeClient registers a connection without a socket, so delivery can be
// observed by reading its Send channel.
func fakeClient(h *Hub, userID string) *Client {
	c := &Client{UserID: userID, Send: make(chan []byte, 8)}
	h.Register(c)
	return c
}

// waitFor reads one message, or reports that none came.
func waitFor(t *testing.T, c *Client) (Event, bool) {
	t.Helper()
	select {
	case raw, ok := <-c.Send:
		if !ok {
			return Event{}, false
		}
		var evt Event
		if err := json.Unmarshal(raw, &evt); err != nil {
			t.Fatalf("decode: %v", err)
		}
		return evt, true
	case <-time.After(500 * time.Millisecond):
		return Event{}, false
	}
}

// An event published on one replica must reach a client held by another.
//
// This is the whole reason the backplane exists. Without one the push succeeds
// into a registry that does not contain the recipient: no error, no log, and a
// user who simply never sees the message. It is invisible on one instance and
// appears the first time anything runs two.
func TestBackplaneDeliversAcrossNodes(t *testing.T) {
	a, b := joinedHubs(t)
	client := fakeClient(a, "user-1")

	b.SendToUser("user-1", Event{Type: "note.created", Payload: map[string]string{"id": "n1"}})

	evt, ok := waitFor(t, client)
	if !ok {
		t.Fatal("a client on node A received nothing for an event published on node B")
	}
	if evt.Type != "note.created" {
		t.Errorf("got %q, want note.created", evt.Type)
	}
}

// The publishing node must not deliver its own message twice.
//
// It delivers locally before publishing, and the backplane hands the message
// back to every subscriber including the publisher. Without the origin check
// every user sees their own events doubled.
func TestBackplaneDoesNotEchoToTheSender(t *testing.T) {
	a, _ := joinedHubs(t)
	client := fakeClient(a, "user-1")

	a.SendToUser("user-1", Event{Type: "note.created"})

	if _, ok := waitFor(t, client); !ok {
		t.Fatal("the local client got nothing at all")
	}
	if evt, ok := waitFor(t, client); ok {
		t.Errorf("the same event arrived twice (second was %q): the node "+
			"delivered locally and then again off its own publish", evt.Type)
	}
}

// Broadcast has to cross nodes too, or a system-wide notice reaches whichever
// replica happened to serve the request.
func TestBackplaneBroadcastReachesEveryNode(t *testing.T) {
	a, b := joinedHubs(t)
	onA := fakeClient(a, "user-1")
	onB := fakeClient(b, "user-2")

	a.Broadcast(Event{Type: "system.maintenance"})

	for name, c := range map[string]*Client{"node A": onA, "node B": onB} {
		if evt, ok := waitFor(t, c); !ok || evt.Type != "system.maintenance" {
			t.Errorf("%s did not receive the broadcast", name)
		}
	}
}

// Revocation has to close sockets on every node.
//
// Sessions are revoked by whichever replica serves the request. If the kick
// stays local, "sign out of all devices" signs out of the devices that happen
// to share a replica with that request and leaves the rest streaming, while
// the UI reports success. That is worse than not offering the button.
func TestBackplaneDisconnectReachesEveryNode(t *testing.T) {
	a, b := joinedHubs(t)
	onA := fakeClient(a, "user-1")

	b.DisconnectUser("user-1")

	// Unregister closes Send, so a closed channel is the signal.
	select {
	case _, open := <-onA.Send:
		if open {
			t.Error("the connection on node A is still delivering after a " +
				"revocation handled by node B")
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("a revocation on node B never closed the connection on node A: " +
			"a signed-out device keeps receiving events")
	}
}

// Without a backplane the hub behaves exactly as it always did.
func TestHubWithoutBackplaneIsUnchanged(t *testing.T) {
	h := NewHub()
	t.Cleanup(func() { _ = h.Close() })
	client := fakeClient(h, "user-1")

	h.SendToUser("user-1", Event{Type: "note.created"})

	if _, ok := waitFor(t, client); !ok {
		t.Error("single-process delivery stopped working")
	}
}

// A publish must never block the caller.
//
// SendToUser is called straight from handlers as well as from the async event
// bus, so a backplane that has stopped answering must cost the request nothing
// beyond a dropped event.
func TestPublishDoesNotBlockOnAStuckBackplane(t *testing.T) {
	h := NewHub(WithBackplane(stuckBackplane{}))
	t.Cleanup(func() { _ = h.Close() })

	done := make(chan struct{})
	go func() {
		for i := 0; i < publishBuffer*2; i++ {
			h.SendToUser("user-1", Event{Type: "note.created"})
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("SendToUser blocked on a backplane that never returns; a " +
			"handler calling it would hang with it")
	}
}

// stuckBackplane accepts a publish and never finishes it.
type stuckBackplane struct{}

func (stuckBackplane) Publish(ctx context.Context, msg []byte) error {
	<-ctx.Done()
	return ctx.Err()
}
func (stuckBackplane) Subscribe(ctx context.Context, deliver func([]byte)) { <-ctx.Done() }
func (stuckBackplane) Close() error                                       { return nil }
`
}
