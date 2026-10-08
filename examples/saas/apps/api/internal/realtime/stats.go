package realtime

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
	Connections int `json:"connections"`
	// Users is how many distinct users those sockets belong to.
	Users int `json:"users"`
	// Channels is how many channels have at least one subscriber here.
	Channels int `json:"channels"`

	// MessagesSent counts messages queued for a connection: events, channel
	// publishes, presence, whispers and replies.
	MessagesSent uint64 `json:"messages_sent"`
	// MessagesDropped counts messages not queued because the connection's send
	// buffer was full. A slow client resyncs on its next REST call; a count that
	// climbs steadily means clients cannot keep up.
	MessagesDropped uint64 `json:"messages_dropped"`

	// ClientEvents counts whispers relayed, and ClientEventsRateLimited the ones
	// dropped for going over ClientEventsPerSecond.
	ClientEvents            uint64 `json:"client_events"`
	ClientEventsRateLimited uint64 `json:"client_events_rate_limited"`

	// Backplane says whether this hub fans out to other replicas.
	// BackplanePublishErrors counts messages that did not reach them: a publish
	// that failed, or one dropped because the publish queue was full.
	Backplane              bool   `json:"backplane"`
	BackplanePublishErrors uint64 `json:"backplane_publish_errors"`
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
