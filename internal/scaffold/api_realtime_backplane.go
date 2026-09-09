package scaffold

// apiRealtimeBackplaneGo emits internal/realtime/backplane.go.
//
// Without this the Hub is an in-process registry and realtime stops working the
// moment there are two API replicas: a user connected to replica A never
// receives an event published on replica B. The push succeeds, into a registry
// that does not contain them. Nothing errors and nothing is logged, so it is
// invisible in development, invisible on a single instance, and shows up as
// "messages sometimes do not arrive" during the first rolling deploy where two
// versions overlap.
func apiRealtimeBackplaneGo() string {
	return `package realtime

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Backplane carries hub traffic between processes.
//
// Delivery is best effort, and deliberately so. A realtime event is a hint
// that something changed, not the record of it: the database holds the truth
// and every client resyncs on its next REST call. Guaranteeing delivery here
// would mean per-subscriber queues, acknowledgements and retention, which is a
// message broker, and one that fails in ways a chat notification does not
// justify. A message published while a node is reconnecting is lost, and that
// is the right trade.
//
// Implement this to sit on something other than Redis. Publish must not block:
// it is called on the request path.
type Backplane interface {
	// Publish sends one already-encoded envelope to every other node.
	Publish(ctx context.Context, msg []byte) error

	// Subscribe delivers envelopes from other nodes until ctx is cancelled.
	// It is expected to keep reconnecting on its own.
	Subscribe(ctx context.Context, deliver func([]byte))

	Close() error
}

// fanout is what crosses the wire between nodes.
//
// Short field names because this is machine-to-machine and every event pays
// for them. Payload is the already-marshalled Event, so a message is encoded
// once no matter how many nodes receive it.
type fanout struct {
	// Origin is the node that published. A node skips its own messages: it
	// delivered to its local clients before publishing, and doing it again on
	// the way back would show every user their own events twice.
	Origin string ` + "`" + `json:"o"` + "`" + `

	Users []string        ` + "`" + `json:"u,omitempty"` + "`" + `
	All   bool            ` + "`" + `json:"a,omitempty"` + "`" + `
	Event json.RawMessage ` + "`" + `json:"e,omitempty"` + "`" + `

	// Kick closes every connection a user holds. Revocation has to cross the
	// backplane or "sign out of all devices" only signs out of the devices
	// that happen to be connected to the replica handling the request, which
	// is worse than not offering it.
	Kick string ` + "`" + `json:"k,omitempty"` + "`" + `
}

// newNodeID identifies this process for the lifetime of the process.
func newNodeID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		// A collision here would make a node ignore another node's messages.
		// Falling back to a timestamp keeps that vanishingly unlikely without
		// making hub construction fail over an entropy hiccup.
		return "node-" + time.Now().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(b)
}

// publishBuffer is how many envelopes may be waiting for the backplane.
//
// Deep enough to ride out a slow round trip, shallow enough that a Redis that
// has stopped answering costs a bounded amount of memory rather than growing
// until the process dies.
const publishBuffer = 256

// publish queues a fanout for the other nodes.
//
// Non-blocking. SendToUser is called directly from handlers as well as from the
// async event bus, and a handler must not wait on a network round trip for a
// message its own clients already have.
func (h *Hub) publish(f fanout) {
	if h.backplane == nil {
		return
	}
	f.Origin = h.nodeID
	msg, err := json.Marshal(f)
	if err != nil {
		log.Printf("[realtime] backplane marshal: %v", err)
		return
	}
	select {
	case h.pub <- msg:
	default:
		// Same policy as a client whose send buffer is full: drop it. The
		// alternative is blocking a request on a backplane that is not
		// keeping up, and the receiving clients resync on their next REST call.
		log.Printf("[realtime] backplane buffer full, dropping an event")
	}
}

// publishLoop drains the queue on one goroutine, which keeps publishes in the
// order they were made. A goroutine per publish would not.
func (h *Hub) publishLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-h.pub:
			// Bounded so a wedged Redis stalls this goroutine briefly rather
			// than forever; the queue absorbs the difference.
			pctx, cancel := context.WithTimeout(ctx, 2*time.Second)
			if err := h.backplane.Publish(pctx, msg); err != nil {
				log.Printf("[realtime] backplane publish: %v", err)
			}
			cancel()
		}
	}
}

// receive applies a fanout that arrived from another node.
func (h *Hub) receive(msg []byte) {
	var f fanout
	if err := json.Unmarshal(msg, &f); err != nil {
		log.Printf("[realtime] backplane decode: %v", err)
		return
	}
	if f.Origin == h.nodeID {
		return // our own message, already delivered locally
	}
	switch {
	case f.Kick != "":
		h.disconnectLocal(f.Kick)
	case f.All:
		h.broadcastLocal(f.Event)
	default:
		h.deliverLocal(f.Users, f.Event)
	}
}

// ── Redis ────────────────────────────────────────────────────────────────

// redisBackplane is pub/sub over the Redis this project already runs for cache
// and background jobs, so the common case adds no new infrastructure.
type redisBackplane struct {
	client  *redis.Client
	channel string
}

// RedisBackplane connects to Redis for cross-process fan-out.
//
// channel namespaces the traffic; pass "" for the default. Two apps sharing one
// Redis need different channels or they will deliver each other's events.
func RedisBackplane(redisURL, channel string) (Backplane, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	if channel == "" {
		channel = "grit:realtime"
	}
	return &redisBackplane{client: redis.NewClient(opts), channel: channel}, nil
}

func (r *redisBackplane) Publish(ctx context.Context, msg []byte) error {
	return r.client.Publish(ctx, r.channel, msg).Err()
}

func (r *redisBackplane) Subscribe(ctx context.Context, deliver func([]byte)) {
	// go-redis resubscribes internally when the connection drops, but it does
	// not survive Redis being unreachable at startup or being restarted long
	// enough for the subscription to be torn down, so the outer loop retries.
	backoff := time.Second
	for {
		if ctx.Err() != nil {
			return
		}
		sub := r.client.Subscribe(ctx, r.channel)

		// Confirm the subscription BEFORE taking the channel. Channel() starts
		// a goroutine that consumes the connection, so a Receive() after it
		// races that goroutine for the same reply and go-redis throws the
		// connection away: "discarding bad PubSub connection: invalid reply".
		// One or the other reads the socket, never both.
		if _, err := sub.Receive(ctx); err != nil {
			_ = sub.Close()
			if ctx.Err() != nil {
				return
			}
			log.Printf("[realtime] backplane subscribe failed, retrying in %s: %v", backoff, err)
			select {
			case <-ctx.Done():
				return
			case <-time.After(backoff):
			}
			if backoff *= 2; backoff > 30*time.Second {
				backoff = 30 * time.Second
			}
			continue
		}

		log.Printf("[realtime] backplane connected on %q", r.channel)
		backoff = time.Second

		ch := sub.Channel()
		for msg := range ch {
			deliver([]byte(msg.Payload))
		}
		// Channel closed: the subscription is gone. Loop and rebuild it.
		_ = sub.Close()
	}
}

func (r *redisBackplane) Close() error { return r.client.Close() }
`
}
