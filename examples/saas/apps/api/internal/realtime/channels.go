package realtime

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
//	presence-*  authorized like private-*, and lists who is subscribed (presence.go)
//
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
	Params  map[string]string // "id" -> "42" for the pattern "invoices.{id}"

	info *interface{} // where SetInfo writes, see presence.go
}

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

// authorizeChannel decides whether userID may subscribe to channel, and returns
// what the authorizer passed to SetInfo.
func authorizeChannel(userID, channel string) (info interface{}, err error) {
	if !validChannel(channel) {
		return nil, ErrInvalidChannel
	}
	if strings.HasPrefix(channel, "public-") {
		return nil, nil
	}
	authorize, params := matchChannel(channel)
	if authorize == nil {
		return nil, ErrChannelForbidden
	}
	// The authorizer runs on the connection's read goroutine, where a panic
	// would take the whole process down rather than one request.
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[realtime] the authorizer for %s panicked: %v", channel, r)
			info, err = nil, ErrChannelForbidden
		}
	}()
	if !authorize(ChannelContext{UserID: userID, Channel: channel, Params: params, info: &info}) {
		return nil, ErrChannelForbidden
	}
	return info, nil
}

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

// Subscribe authorizes c for channel and adds it. A channel the connection
// already holds is a no-op. The authorizer runs outside the hub's lock, so a
// slow database check does not stall every other connection, and so does
// joining a presence channel, which can wait on Redis.
func (h *Hub) Subscribe(c *Client, channel string) error {
	info, err := authorizeChannel(c.UserID, channel)
	if err != nil {
		return err
	}
	h.mu.Lock()
	if _, open := h.clients[c.UserID][c]; !open {
		h.mu.Unlock()
		return nil // closed while the authorizer ran: nothing left to deliver to
	}
	err = h.channels.add(c, channel)
	h.mu.Unlock()
	if err != nil {
		return err
	}
	if strings.HasPrefix(channel, "presence-") {
		h.joinPresence(c, channel, info)
	}
	return nil
}

// Unsubscribe removes c from channel. Unknown channels are ignored.
func (h *Hub) Unsubscribe(c *Client, channel string) {
	h.mu.Lock()
	h.channels.remove(c, channel)
	h.mu.Unlock()
	h.leavePresence(c, channel)
}

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
		h.offer(c, bytes, "message", channel)
	}
}

// channelEnvelope is an Event on the wire with the channel it arrived on.
type channelEnvelope struct {
	Type    string      `json:"type"`
	Channel string      `json:"channel"`
	Payload interface{} `json:"payload"`
}

type clientMessage struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
}

type subscriptionError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
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

// HandleClientMessage applies one message a client sent on the socket:
// subscribe, unsubscribe, or a client event (whispers.go). Anything else is
// ignored, so a client newer than the server does not break the connection.
func (h *Hub) HandleClientMessage(c *Client, raw []byte) {
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
		h.reply(c, "subscribed", msg.Channel, nil)
		h.replyPresenceMembers(c, msg.Channel)
	case "client-event":
		h.relayClientEvent(c, raw)
	case "unsubscribe":
		h.Unsubscribe(c, msg.Channel)
		h.reply(c, "unsubscribed", msg.Channel, nil)
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
	h.offer(c, bytes, kind+" reply", channel)
}
