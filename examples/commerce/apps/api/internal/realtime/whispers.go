package realtime

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
	Channel string          `json:"channel"`
	Event   string          `json:"event"`
	Payload json.RawMessage `json:"payload"`
}

// clientEvent is a whisper as the other subscribers receive it.
type clientEvent struct {
	Event  string          `json:"event"`
	UserID string          `json:"user_id"`
	Data   json.RawMessage `json:"data,omitempty"`
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
