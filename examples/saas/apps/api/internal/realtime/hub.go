// Package realtime is a tiny WebSocket fan-out hub. One Hub per process;
// each authenticated user can have multiple connections (e.g. desktop +
// mobile + web). The hub owns the registry and exposes safe SendToUser /
// SendToUsers / Broadcast helpers that handlers call from anywhere.
//
// Wire format on the websocket is a JSON envelope:
//
//	{ "type": "<topic>", "payload": { ... } }
//
// Topics are caller-defined strings. Suggested namespacing:
//
//	chat.message.new       — payload is a chat message
//	notification.new       — payload is a notification
//	system.connected       — server greeting on first connect
//	resource.<name>.<verb> — e.g. building.created, lease.expired
package realtime

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Event is the envelope every WS message uses on the wire.
type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Client is one open WebSocket connection bound to a user.
type Client struct {
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte

	// ExpiresAt is when the JWT that authorised this connection stops being
	// valid. The handshake is the only time the token is checked, so without
	// this the socket would outlive its own credential and keep streaming to
	// a session that has since been signed out. writePump closes the
	// connection at this instant and the client reconnects with a fresh
	// token, which is the same bound REST already operates under.
	ExpiresAt time.Time

	// clientEvents limits the whispers this connection may send. See whispers.go.
	clientEvents clientEventLimiter
}

// disconnectGrace is how long DisconnectUser waits for writePump to send its
// close frame and tear the socket down before forcing it. It matches the write
// deadline writePump already works to.
const disconnectGrace = 10 * time.Second

// Hub manages connected clients. Safe for concurrent use.
type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*Client]struct{} // userID -> set of connections

	// channels is who is subscribed to which channel. Guarded by mu. See
	// channels.go.
	channels channelIndex

	// presence is who this process holds in presence channels. See presence.go.
	presence presenceSet

	// stats counts what the hub delivered, dropped and failed to publish. See
	// stats.go.
	stats hubStats

	// nodeID identifies this process so it can ignore its own messages coming
	// back off the backplane.
	nodeID string

	// backplane is nil for a single-process deployment, which is the default
	// and is correct for most projects. See WithBackplane.
	backplane Backplane
	pub       chan []byte
	cancel    context.CancelFunc
}

// Option configures a Hub at construction.
type Option func(*Hub)

// WithBackplane fans every send out to the other API processes.
//
// Without one the Hub is an in-process registry, so a second replica silently
// halves realtime: a user connected to replica A never sees an event published
// on replica B, the push succeeds into a registry that does not contain them,
// and nothing errors. Pass this whenever more than one process serves
// websockets, including during a rolling deploy where two versions overlap.
//
// Delivery between nodes is best effort. See Backplane.
func WithBackplane(b Backplane) Option {
	return func(h *Hub) { h.backplane = b }
}

// WithRedis is WithBackplane over the Redis this project already runs.
//
// An empty url is a no-op, so the common wiring is one unconditional line:
//
//	realtime.NewHub(realtime.WithRedis(cfg.RedisURL, ""))
//
// A project with no Redis configured then keeps the single-process behaviour
// instead of failing to start.
func WithRedis(redisURL, channel string) Option {
	return func(h *Hub) {
		if redisURL == "" {
			return
		}
		bp, err := RedisBackplane(redisURL, channel)
		if err != nil {
			// Not fatal. Realtime degrades to this process's own clients,
			// which is exactly what every project had before backplanes
			// existed, and the API still serves everything else.
			log.Printf("[realtime] no backplane, staying single-process: %v", err)
			return
		}
		h.backplane = bp
	}
}

// NewHub returns a Hub, single-process unless an option says otherwise.
func NewHub(opts ...Option) *Hub {
	h := &Hub{
		clients: make(map[string]map[*Client]struct{}),
		nodeID:  newNodeID(),
	}
	for _, opt := range opts {
		opt(h)
	}
	if h.backplane != nil {
		h.pub = make(chan []byte, publishBuffer)
		ctx, cancel := context.WithCancel(context.Background())
		h.cancel = cancel
		go h.backplane.Subscribe(ctx, h.receive)
		go h.publishLoop(ctx)
		h.startPresence(ctx)
	}
	return h
}

// Close stops the backplane subscription. Local clients are unaffected.
func (h *Hub) Close() error {
	if h.cancel != nil {
		h.cancel()
	}
	if h.backplane != nil {
		return h.backplane.Close()
	}
	return nil
}

// Register adds a client to the hub. A user can have multiple registered
// clients (different devices); each gets its own slot.
func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	set, ok := h.clients[c.UserID]
	if !ok {
		set = make(map[*Client]struct{})
		h.clients[c.UserID] = set
	}
	set[c] = struct{}{}
	log.Printf("[realtime] client registered user=%s total=%d", c.UserID, len(set))
}

// Unregister removes a client and closes its Send channel. Safe to call
// once per client (e.g. from the read pump's defer).
func (h *Hub) Unregister(c *Client) {
	// Deferred before the unlock, so it runs after it: leaving a presence
	// channel can wait on Redis, which must not hold up the hub.
	defer h.leaveAllPresence(c)
	h.mu.Lock()
	defer h.mu.Unlock()
	if set, ok := h.clients[c.UserID]; ok {
		if _, exists := set[c]; exists {
			delete(set, c)
			h.channels.drop(c)
			close(c.Send)
		}
		if len(set) == 0 {
			delete(h.clients, c.UserID)
		}
	}
}

// DisconnectUser closes every connection a user holds, immediately.
//
// Call this when a session is revoked, a password changes, or an account is
// deactivated. Waiting for the token to expire leaves a revoked device
// receiving live data for the rest of the access-token lifetime, which is not
// what "sign out of all devices" tells the user happened.
func (h *Hub) DisconnectUser(userID string) {
	h.disconnectLocal(userID)
	// And on every other node. A revoked connection on a replica that did not
	// handle the revocation request would otherwise stay open and keep
	// receiving, which is the failure this function exists to prevent.
	h.publish(fanout{Kick: userID})
}

// disconnectLocal closes this process's connections for a user. It is what
// both DisconnectUser and an incoming kick run.
func (h *Hub) disconnectLocal(userID string) {
	h.mu.Lock()
	set, ok := h.clients[userID]
	if !ok {
		h.mu.Unlock()
		return
	}
	conns := make([]*websocket.Conn, 0, len(set))
	for c := range set {
		conns = append(conns, c.Conn)
		h.channels.drop(c)
		// Closing Send makes writePump emit a proper close frame and tear the
		// connection down itself, so the client sees a clean 1000 rather than
		// an abnormal 1006 and can distinguish "you were signed out" from
		// "the network dropped".
		close(c.Send)
	}
	delete(h.clients, userID)
	h.mu.Unlock()

	// Backstop, outside the lock. If writePump is wedged on a dead socket it
	// will not get to its own deferred Close, and a revoked connection that
	// stays open is the thing this function exists to prevent. wsWriteWait is
	// the bound writePump is already operating under.
	go func() {
		time.Sleep(disconnectGrace)
		for _, conn := range conns {
			if conn != nil { // a Client made without a socket has no Conn
				_ = conn.Close()
			}
		}
	}()
}

// SendToUser delivers an event to every connection bound to userID.
// If a connection's send buffer is full the message is dropped for that
// connection only — we never block the entire hub on a slow client.
// The slow client will resync on its next REST poll/refetch.
func (h *Hub) SendToUser(userID string, evt Event) {
	h.SendToUsers([]string{userID}, evt)
}

// deliverLocal pushes an encoded event to this process's own connections.
// It never touches the backplane, so it is also what a received message runs.
//
// The sends happen under the read lock. Unregister and DisconnectUser close a
// client's Send under the write lock, so a send made after releasing the lock
// could land on a channel closed in between and panic, and a panic here stops
// the whole process. A non-blocking send keeps the lock short.
func (h *Hub) deliverLocal(userIDs []string, bytes []byte) {
	if len(bytes) == 0 {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, uid := range userIDs {
		for c := range h.clients[uid] {
			h.offer(c, bytes, "message", "")
		}
	}
}

// broadcastLocal is deliverLocal for every connection on this node, under the
// read lock for the same reason.
func (h *Hub) broadcastLocal(bytes []byte) {
	if len(bytes) == 0 {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, set := range h.clients {
		for c := range set {
			h.offer(c, bytes, "broadcast", "")
		}
	}
}

// SendToUsers fans out to a slice of user IDs.
func (h *Hub) SendToUsers(userIDs []string, evt Event) {
	if len(userIDs) == 0 {
		return
	}
	bytes, err := json.Marshal(evt)
	if err != nil {
		log.Printf("[realtime] marshal: %v", err)
		return
	}
	// This node first, and unconditionally: local clients must not wait on
	// Redis, and must still be served when Redis is down.
	h.deliverLocal(userIDs, bytes)
	// Then once for every other node, whatever the size of the audience. A
	// per-user publish would send the same payload N times over the wire.
	h.publish(fanout{Users: userIDs, Event: bytes})
}

// Broadcast delivers an event to every connected client, regardless of
// user. Use sparingly — for system-wide announcements, maintenance
// notices, etc.
func (h *Hub) Broadcast(evt Event) {
	bytes, err := json.Marshal(evt)
	if err != nil {
		log.Printf("[realtime] marshal: %v", err)
		return
	}
	h.broadcastLocal(bytes)
	h.publish(fanout{All: true, Event: bytes})
}
