package scaffold

// Realtime presence (A3): who is subscribed to a presence-* channel, on every
// replica, and a fix for a hub crash found while building it.
//
// presence.go is new framework code. It needs a field and two calls in hub.go
// and a few changes to channels.go, and on upgrade channels.go is written only
// when missing, so those changes are hunks shared by the templates and the
// repair (realtime_presence_repair.go), like the A1 and A2 hunks.

// apiRealtimePresenceGo emits internal/realtime/presence.go.
func apiRealtimePresenceGo() string {
	return `package realtime

import (
	"context"
	"encoding/json"
	"hash/fnv"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Presence channels know who is subscribed.
//
// A presence-* channel is authorized like a private one (see channels.go), and
// the hub also keeps its member list: one entry per user, however many tabs or
// devices that user subscribed from. A connection that subscribes receives, right
// after subscribed,
//
//	{"type":"presence.members","channel":"presence-rooms.1","payload":{"members":[{"user_id":"7","info":{"name":"Ada"},"joined_at":"..."}]}}
//
// and then, as users arrive and go,
//
//	{"type":"presence.joined","channel":"presence-rooms.1","payload":{"user_id":"8","info":{...},"joined_at":"..."}}
//	{"type":"presence.left","channel":"presence-rooms.1","payload":{"user_id":"8"}}
//
// A user's first connection announces them and their last connection to leave
// announces them gone. The connection that joins is not told about itself: the
// snapshot it receives already lists it.
//
// An authorizer decides what the other members see with SetInfo:
//
//	realtime.Channel("presence-rooms.{id}", func(c realtime.ChannelContext) bool {
//	    member, ok := roomMember(db, c.Param("id"), c.UserID)
//	    if ok {
//	        c.SetInfo(map[string]string{"name": member.Name})
//	    }
//	    return ok
//	})
//
// With the Redis backplane (WithRedis) members live in Redis, one hash per
// channel, so every replica lists the same members. Each replica rewrites its
// own entries every PresenceHeartbeat, and an entry not rewritten for
// PresenceTTL is removed by the next heartbeat of a replica holding that
// channel, which announces presence.left. So the members of a replica that dies
// without closing its sockets leave within about PresenceTTL plus one
// heartbeat. Without Redis the members are this process's own.

// Presence timing, read when a hub is built. PresenceTTL bounds how long a dead
// replica's members stay listed; PresenceHeartbeat is how often a live replica
// renews its own, and has to be well inside the TTL.
var (
	PresenceHeartbeat = 15 * time.Second
	PresenceTTL       = 45 * time.Second
)

// MaxPresenceInfoBytes bounds what an authorizer attaches with SetInfo, once
// encoded. Larger info is left out and the member joins without it.
var MaxPresenceInfoBytes = 1024

// presenceTimeout bounds each Redis call, so a Redis that stops answering
// stalls a subscribe briefly rather than for good.
const presenceTimeout = 2 * time.Second

// Member is one user in a presence channel.
type Member struct {
	UserID   string          ` + "`" + `json:"user_id"` + "`" + `
	Info     json.RawMessage ` + "`" + `json:"info,omitempty"` + "`" + `
	JoinedAt time.Time       ` + "`" + `json:"joined_at"` + "`" + `
}

type presenceLeft struct {
	UserID string ` + "`" + `json:"user_id"` + "`" + `
}

type presenceSnapshot struct {
	Members []Member ` + "`" + `json:"members"` + "`" + `
}

// SetInfo sets what the other members of a presence channel see about the user
// being authorized: a name, an avatar URL, a role. Keep it small: it is sent
// with every presence.joined and presence.members, and info larger than
// MaxPresenceInfoBytes is left out. On a private channel it does nothing.
func (c ChannelContext) SetInfo(info interface{}) {
	if c.info != nil {
		*c.info = info
	}
}

// presenceSet is a hub's presence state. The zero value keeps members in this
// process only; startPresence moves them into Redis.
type presenceSet struct {
	store     *redisPresence
	heartbeat time.Duration

	// stripes serialize the joins and leaves of one user in one channel, so
	// their Redis writes land in the order the connections made them.
	stripes [32]sync.Mutex

	mu     sync.Mutex
	joined map[*Client]map[string]struct{}      // connection -> its presence channels
	local  map[string]map[string]*localMember // channel -> user -> this process's hold
}

type localMember struct {
	member Member
	conns  int
}

func (p *presenceSet) stripe(channel, userID string) *sync.Mutex {
	sum := fnv.New32a()
	_, _ = sum.Write([]byte(channel + "\x00" + userID))
	return &p.stripes[sum.Sum32()%uint32(len(p.stripes))]
}

// add records c in channel. added is false when it was already there; first is
// true when c is the user's only connection in the channel on this process.
func (p *presenceSet) add(c *Client, channel string, info interface{}) (member Member, first, added bool) {
	if _, ok := p.joined[c][channel]; ok {
		return Member{}, false, false
	}
	if p.joined == nil {
		p.joined = make(map[*Client]map[string]struct{})
		p.local = make(map[string]map[string]*localMember)
	}
	if p.joined[c] == nil {
		p.joined[c] = make(map[string]struct{})
	}
	p.joined[c][channel] = struct{}{}
	if p.local[channel] == nil {
		p.local[channel] = make(map[string]*localMember)
	}
	held := p.local[channel][c.UserID]
	if held == nil {
		held = &localMember{member: Member{
			UserID:   c.UserID,
			Info:     encodePresenceInfo(channel, info),
			JoinedAt: time.Now().UTC().Truncate(time.Millisecond),
		}}
		p.local[channel][c.UserID] = held
	}
	held.conns++
	return held.member, held.conns == 1, true
}

// remove forgets c in channel and reports whether it was the user's last
// connection there on this process.
func (p *presenceSet) remove(c *Client, channel string) bool {
	if _, ok := p.joined[c][channel]; !ok {
		return false
	}
	delete(p.joined[c], channel)
	if len(p.joined[c]) == 0 {
		delete(p.joined, c)
	}
	held := p.local[channel][c.UserID]
	if held == nil {
		return false
	}
	if held.conns--; held.conns > 0 {
		return false
	}
	delete(p.local[channel], c.UserID)
	if len(p.local[channel]) == 0 {
		delete(p.local, channel)
	}
	return true
}

func encodePresenceInfo(channel string, info interface{}) json.RawMessage {
	if info == nil {
		return nil
	}
	raw, err := json.Marshal(info)
	if err != nil {
		log.Printf("[realtime] presence info for %s is not JSON, leaving it out: %v", channel, err)
		return nil
	}
	if string(raw) == "null" {
		return nil
	}
	if MaxPresenceInfoBytes > 0 && len(raw) > MaxPresenceInfoBytes {
		log.Printf("[realtime] presence info for %s is %d bytes, over the %d allowed, leaving it out", channel, len(raw), MaxPresenceInfoBytes)
		return nil
	}
	return raw
}

// startPresence keeps members in Redis when the backplane is Redis, and starts
// the heartbeat that renews this process's entries there.
func (h *Hub) startPresence(ctx context.Context) {
	rb, ok := h.backplane.(*redisBackplane)
	if !ok {
		return
	}
	ttl, heartbeat := PresenceTTL, PresenceHeartbeat
	if ttl <= 0 {
		ttl = 45 * time.Second
	}
	if heartbeat <= 0 || heartbeat*2 > ttl {
		// A heartbeat this close to the TTL lets live members expire between beats.
		heartbeat = ttl / 3
	}
	h.presence.store = &redisPresence{client: rb.client, prefix: rb.channel + ":presence:", ttl: ttl}
	h.presence.heartbeat = heartbeat
	go h.presenceLoop(ctx)
}

// joinPresence adds c to a presence channel it has just subscribed to. The
// user's first connection announces them to the channel's other subscribers.
func (h *Hub) joinPresence(c *Client, channel string, info interface{}) {
	p := &h.presence
	lock := p.stripe(channel, c.UserID)
	lock.Lock()
	defer lock.Unlock()

	// The registry check and the add happen under the hub's read lock, which
	// Unregister needs exclusively before it reads the channels to leave. So a
	// connection closing now is either refused here or finds this channel.
	h.mu.RLock()
	_, open := h.clients[c.UserID][c]
	var member Member
	first := false
	if open {
		p.mu.Lock()
		member, first, _ = p.add(c, channel, info)
		p.mu.Unlock()
	}
	h.mu.RUnlock()
	if !first {
		return
	}

	announce := true
	if p.store != nil {
		var err error
		if announce, err = p.store.join(channel, h.nodeID, member); err != nil {
			// Still a member here; the next heartbeat writes the entry again.
			log.Printf("[realtime] presence: joining %s in Redis: %v", channel, err)
			announce = true
		}
	}
	if announce {
		h.publishPresence(channel, "presence.joined", member, c)
	}
}

// leavePresence removes c from a presence channel. The user's last connection
// announces them gone. A channel c is not a member of is ignored.
func (h *Hub) leavePresence(c *Client, channel string) {
	if !strings.HasPrefix(channel, "presence-") {
		return
	}
	p := &h.presence
	lock := p.stripe(channel, c.UserID)
	lock.Lock()
	defer lock.Unlock()

	p.mu.Lock()
	last := p.remove(c, channel)
	p.mu.Unlock()
	if !last {
		return
	}

	announce := true
	if p.store != nil {
		var err error
		if announce, err = p.store.leave(channel, h.nodeID, c.UserID); err != nil {
			// Told anyway. The entry is no longer renewed, so it expires too.
			log.Printf("[realtime] presence: leaving %s in Redis: %v", channel, err)
			announce = true
		}
	}
	if announce {
		h.publishPresence(channel, "presence.left", presenceLeft{UserID: c.UserID}, nil)
	}
}

// leaveAllPresence removes a closing connection from every presence channel.
func (h *Hub) leaveAllPresence(c *Client) {
	p := &h.presence
	p.mu.Lock()
	channels := make([]string, 0, len(p.joined[c]))
	for channel := range p.joined[c] {
		channels = append(channels, channel)
	}
	p.mu.Unlock()
	for _, channel := range channels {
		h.leavePresence(c, channel)
	}
}

// publishPresence tells a presence channel's subscribers, here and on every
// other replica, that someone joined or left. except is the connection that
// joined, whose snapshot already lists it.
func (h *Hub) publishPresence(channel, kind string, payload interface{}, except *Client) {
	bytes, err := json.Marshal(channelEnvelope{Type: kind, Channel: channel, Payload: payload})
	if err != nil {
		log.Printf("[realtime] marshal: %v", err)
		return
	}
	h.mu.RLock()
	for c := range h.channels.members[channel] {
		if c == except {
			continue
		}
` + presencePublishSendNew + `	}
	h.mu.RUnlock()
	h.publish(fanout{Channel: channel, Event: bytes})
}

// replyPresenceMembers sends a connection that has just joined a presence
// channel the channel's members, itself included.
func (h *Hub) replyPresenceMembers(c *Client, channel string) {
	if !strings.HasPrefix(channel, "presence-") {
		return
	}
	h.presence.mu.Lock()
	_, member := h.presence.joined[c][channel]
	h.presence.mu.Unlock()
	if member {
		h.reply(c, "presence.members", channel, presenceSnapshot{Members: h.PresenceMembers(channel)})
	}
}

// PresenceMembers lists who is in a presence channel, across every replica when
// the hub uses Redis. Each user appears once, earliest to join first.
func (h *Hub) PresenceMembers(channel string) []Member {
	if !strings.HasPrefix(channel, "presence-") {
		return []Member{}
	}
	var members []Member
	var err error
	if h.presence.store != nil {
		if members, err = h.presence.store.members(channel); err != nil {
			log.Printf("[realtime] presence: reading %s from Redis, listing this process's members only: %v", channel, err)
		}
	}
	if h.presence.store == nil || err != nil {
		members = h.localPresence(channel)
	}
	sort.SliceStable(members, func(i, j int) bool {
		if !members[i].JoinedAt.Equal(members[j].JoinedAt) {
			return members[i].JoinedAt.Before(members[j].JoinedAt)
		}
		return members[i].UserID < members[j].UserID
	})
	seen := make(map[string]struct{}, len(members))
	unique := make([]Member, 0, len(members))
	for _, m := range members {
		if _, dup := seen[m.UserID]; !dup {
			seen[m.UserID] = struct{}{}
			unique = append(unique, m)
		}
	}
	return unique
}

func (h *Hub) localPresence(channel string) []Member {
	h.presence.mu.Lock()
	defer h.presence.mu.Unlock()
	members := make([]Member, 0, len(h.presence.local[channel]))
	for _, held := range h.presence.local[channel] {
		members = append(members, held.member)
	}
	return members
}

func (h *Hub) presenceLoop(ctx context.Context) {
	ticker := time.NewTicker(h.presence.heartbeat)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Each presence write carries its own presenceTimeout, because
			// join and leave also run from a connection closing, where there
			// is no request to inherit from. The loop's ctx only stops it.
			h.presenceBeat() //nolint:contextcheck // bounded per call, see above
		}
	}
}

// presenceBeat renews this process's entries, writes back any Redis lost, and
// removes the entries other replicas stopped renewing, announcing who left.
func (h *Hub) presenceBeat() {
	p := &h.presence
	h.mu.RLock()
	p.mu.Lock()
	var closed []*Client
	for c := range p.joined {
		if _, open := h.clients[c.UserID][c]; !open {
			closed = append(closed, c)
		}
	}
	held := make(map[string][]string, len(p.local))
	for channel, users := range p.local {
		for userID := range users {
			held[channel] = append(held[channel], userID)
		}
	}
	p.mu.Unlock()
	h.mu.RUnlock()

	// A connection DisconnectUser closed leaves when its read pump unregisters
	// it. This catches one that never does.
	for _, c := range closed {
		h.leaveAllPresence(c)
	}

	for channel, users := range held {
		missing, err := p.store.refresh(channel, h.nodeID, users)
		if err != nil {
			log.Printf("[realtime] presence: renewing %s in Redis: %v", channel, err)
			continue
		}
		for _, userID := range missing {
			h.rejoinPresence(channel, userID)
		}
		gone, err := p.store.sweep(channel)
		if err != nil {
			log.Printf("[realtime] presence: sweeping %s in Redis: %v", channel, err)
			continue
		}
		for _, userID := range gone {
			h.publishPresence(channel, "presence.left", presenceLeft{UserID: userID}, nil)
		}
	}
}

// rejoinPresence writes back an entry this process still holds but Redis lost,
// to a restart or to a pause longer than the TTL.
func (h *Hub) rejoinPresence(channel, userID string) {
	p := &h.presence
	lock := p.stripe(channel, userID)
	lock.Lock()
	defer lock.Unlock()

	p.mu.Lock()
	held := p.local[channel][userID]
	var member Member
	if held != nil {
		member = held.member
	}
	p.mu.Unlock()
	if held == nil {
		return // left since the heartbeat looked
	}
	announce, err := p.store.join(channel, h.nodeID, member)
	if err != nil {
		log.Printf("[realtime] presence: rejoining %s in Redis: %v", channel, err)
		return
	}
	if announce {
		h.publishPresence(channel, "presence.joined", member, nil)
	}
}

// redisPresence keeps members in Redis: a hash per channel, a field per user
// per replica ("<node id>:<user id>"), each value "<last renewed, unix ms>
// <member JSON>". The scripts take the time from Redis, so replicas whose clocks
// disagree still agree on which entries are stale, and each runs atomically, so
// two replicas cannot both announce the same user.
type redisPresence struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
}

func (r *redisPresence) key(channel string) string { return r.prefix + channel }

func (r *redisPresence) join(channel, node string, m Member) (bool, error) {
	raw, err := json.Marshal(m)
	if err != nil {
		return false, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), presenceTimeout)
	defer cancel()
	first, err := presenceJoinScript.Run(ctx, r.client, []string{r.key(channel)},
		node+":"+m.UserID, m.UserID, string(raw), r.ttl.Milliseconds()).Int()
	return first == 1, err
}

func (r *redisPresence) leave(channel, node, userID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), presenceTimeout)
	defer cancel()
	last, err := presenceLeaveScript.Run(ctx, r.client, []string{r.key(channel)},
		node+":"+userID, userID, r.ttl.Milliseconds()).Int()
	return last == 1, err
}

// refresh renews this node's entries for users and returns those Redis no
// longer has.
func (r *redisPresence) refresh(channel, node string, users []string) ([]string, error) {
	args := make([]interface{}, 0, len(users)+1)
	args = append(args, r.ttl.Milliseconds())
	for _, userID := range users {
		args = append(args, node+":"+userID)
	}
	ctx, cancel := context.WithTimeout(context.Background(), presenceTimeout)
	defer cancel()
	fields, err := presenceRefreshScript.Run(ctx, r.client, []string{r.key(channel)}, args...).StringSlice()
	if err != nil {
		return nil, err
	}
	missing := make([]string, 0, len(fields))
	for _, field := range fields {
		missing = append(missing, strings.TrimPrefix(field, node+":"))
	}
	return missing, nil
}

// sweep removes entries nobody renewed within the TTL and returns the users
// with no entry left.
func (r *redisPresence) sweep(channel string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), presenceTimeout)
	defer cancel()
	return presenceSweepScript.Run(ctx, r.client, []string{r.key(channel)}, r.ttl.Milliseconds()).StringSlice()
}

func (r *redisPresence) members(channel string) ([]Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), presenceTimeout)
	defer cancel()
	values, err := presenceMembersScript.Run(ctx, r.client, []string{r.key(channel)}, r.ttl.Milliseconds()).StringSlice()
	if err != nil {
		return nil, err
	}
	members := make([]Member, 0, len(values))
	for _, value := range values {
		var m Member
		if err := json.Unmarshal([]byte(value), &m); err != nil {
			log.Printf("[realtime] presence: skipping an unreadable member of %s: %v", channel, err)
			continue
		}
		members = append(members, m)
	}
	return members, nil
}

const presenceLua = ` + "`" + `
local function now_ms()
  local t = redis.call('TIME')
  return tonumber(t[1]) * 1000 + math.floor(tonumber(t[2]) / 1000)
end
local function seen_of(value)
  return tonumber(string.match(value, '^(%d+) ')) or 0
end
local function user_of(field)
  local at = string.find(field, ':', 1, true)
  if not at then return field end
  return string.sub(field, at + 1)
end
local function present(key, user, skip, now, ttl)
  local all = redis.call('HGETALL', key)
  for i = 1, #all, 2 do
    if all[i] ~= skip and user_of(all[i]) == user and now - seen_of(all[i + 1]) < ttl then
      return true
    end
  end
  return false
end
` + "`" + `

// KEYS[1] channel hash; ARGV field, user, member JSON, ttl ms. Returns 1 when
// no fresh entry holds the user, this replica's own included, so a heartbeat that rewrites a join still in flight does not announce it twice.
var presenceJoinScript = redis.NewScript(presenceLua + ` + "`" + `
local now, ttl = now_ms(), tonumber(ARGV[4])
local first = not present(KEYS[1], ARGV[2], '', now, ttl)
redis.call('HSET', KEYS[1], ARGV[1], string.format('%d', now) .. ' ' .. ARGV[3])
redis.call('PEXPIRE', KEYS[1], ttl)
if first then return 1 end
return 0
` + "`" + `)

// ARGV field, user, ttl ms. Returns 1 when that was the user's last entry.
var presenceLeaveScript = redis.NewScript(presenceLua + ` + "`" + `
if redis.call('HDEL', KEYS[1], ARGV[1]) == 0 then return 0 end
if present(KEYS[1], ARGV[2], '', now_ms(), tonumber(ARGV[3])) then return 0 end
return 1
` + "`" + `)

// ARGV ttl ms, then fields. Renews the fields that exist, returns the rest.
var presenceRefreshScript = redis.NewScript(presenceLua + ` + "`" + `
local now, ttl = now_ms(), tonumber(ARGV[1])
local missing, kept = {}, false
for i = 2, #ARGV do
  local value = redis.call('HGET', KEYS[1], ARGV[i])
  local space = value and string.find(value, ' ', 1, true)
  if space then
    redis.call('HSET', KEYS[1], ARGV[i], string.format('%d', now) .. string.sub(value, space))
    kept = true
  else
    table.insert(missing, ARGV[i])
  end
end
if kept then redis.call('PEXPIRE', KEYS[1], ttl) end
return missing
` + "`" + `)

// ARGV ttl ms. Removes stale fields, returns the users left with none.
var presenceSweepScript = redis.NewScript(presenceLua + ` + "`" + `
local now, ttl = now_ms(), tonumber(ARGV[1])
local all = redis.call('HGETALL', KEYS[1])
local dropped = {}
for i = 1, #all, 2 do
  if now - seen_of(all[i + 1]) >= ttl then
    redis.call('HDEL', KEYS[1], all[i])
    dropped[user_of(all[i])] = true
  end
end
local gone = {}
for user in pairs(dropped) do
  if not present(KEYS[1], user, '', now, ttl) then table.insert(gone, user) end
end
return gone
` + "`" + `)

// ARGV ttl ms. Returns the member JSON of every fresh field.
var presenceMembersScript = redis.NewScript(presenceLua + ` + "`" + `
local now, ttl = now_ms(), tonumber(ARGV[1])
local all = redis.call('HGETALL', KEYS[1])
local out = {}
for i = 1, #all, 2 do
  local value = all[i + 1]
  local space = string.find(value, ' ', 1, true)
  if space and now - seen_of(value) < ttl then
    table.insert(out, string.sub(value, space + 1))
  end
end
return out
` + "`" + `)
`
}

// apiRealtimePresenceTestGo emits internal/realtime/presence_test.go. The
// cross-replica test needs a Redis and runs when REALTIME_TEST_REDIS_URL names
// one; its helpers are its own, so the file compiles beside any other test file.
func apiRealtimePresenceTestGo() string {
	return `package realtime

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
	Type    string          ` + "`" + `json:"type"` + "`" + `
	Channel string          ` + "`" + `json:"channel"` + "`" + `
	Payload json.RawMessage ` + "`" + `json:"payload"` + "`" + `
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
			if strings.Contains(string(raw), ` + "`" + `"type":"` + "`" + `+kind+` + "`" + `"` + "`" + `) {
				t.Fatalf("%s: unexpected %s", c.UserID, raw)
			}
		case <-deadline:
			return
		}
	}
}

func presenceSubscribe(t *testing.T, h *Hub, c *Client, channel string) presenceMessage {
	t.Helper()
	h.HandleClientMessage(c, []byte(` + "`" + `{"type":"subscribe","channel":"` + "`" + `+channel+` + "`" + `"}` + "`" + `))
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
	if !strings.Contains(string(snap.Payload), ` + "`" + `"name":"User alice"` + "`" + `) {
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
	hub.HandleClientMessage(stranger, []byte(` + "`" + `{"type":"subscribe","channel":"` + "`" + `+room+` + "`" + `"}` + "`" + `))
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
`
}

// apiRealtimeHubSendTestGo emits internal/realtime/hub_send_test.go.
func apiRealtimeHubSendTestGo() string {
	return `package realtime

import (
	"io"
	"log"
	"os"
	"sync"
	"testing"
)

// A send racing a closing connection must not panic.
//
// deliverLocal and broadcastLocal used to copy their targets under the read
// lock and send after releasing it. An Unregister in that gap closed Send
// first, the send panicked with "send on closed channel", and a panic on a
// goroutine takes the whole API process down, not just one socket.
func TestSendsRacingUnregisterDoNotPanic(t *testing.T) {
	log.SetOutput(io.Discard)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	hub := NewHub()
	const rounds, clientsPerRound = 50, 200
	sends := 0
	for round := 0; round < rounds; round++ {
		clients := make([]*Client, clientsPerRound)
		for i := range clients {
			clients[i] = &Client{UserID: "u1", Send: make(chan []byte, 1)}
			hub.Register(clients[i])
		}
		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			defer wg.Done()
			for _, c := range clients {
				hub.Unregister(c)
			}
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < 20; i++ {
				hub.Broadcast(Event{Type: "ping"})
			}
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < 20; i++ {
				hub.SendToUsers([]string{"u1"}, Event{Type: "ping"})
			}
		}()
		wg.Wait()
		sends += 40
	}
	t.Logf("%d sends raced %d closing connections without a panic", sends, rounds*clientsPerRound)
}
`
}

// hub.go hunks. Each Old is exactly what Grit wrote in v3.274.0.

const hubLocalSendsOld = `// deliverLocal pushes an encoded event to this process's own connections.
// It never touches the backplane, so it is also what a received message runs.
func (h *Hub) deliverLocal(userIDs []string, bytes []byte) {
	if len(bytes) == 0 {
		return
	}
	h.mu.RLock()
	targets := make([]*Client, 0, len(userIDs))
	for _, uid := range userIDs {
		for c := range h.clients[uid] {
			targets = append(targets, c)
		}
	}
	h.mu.RUnlock()
	for _, c := range targets {
		select {
		case c.Send <- bytes:
		default:
			log.Printf("[realtime] dropping message for slow client user=%s", c.UserID)
		}
	}
}

// broadcastLocal is deliverLocal for every connection on this node.
func (h *Hub) broadcastLocal(bytes []byte) {
	if len(bytes) == 0 {
		return
	}
	h.mu.RLock()
	targets := make([]*Client, 0)
	for _, set := range h.clients {
		for c := range set {
			targets = append(targets, c)
		}
	}
	h.mu.RUnlock()
	for _, c := range targets {
		select {
		case c.Send <- bytes:
		default:
			log.Printf("[realtime] dropping broadcast for slow client user=%s", c.UserID)
		}
	}
}
`

const hubLocalSendsNew = `// deliverLocal pushes an encoded event to this process's own connections.
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
			select {
			case c.Send <- bytes:
			default:
				log.Printf("[realtime] dropping message for slow client user=%s", c.UserID)
			}
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
			select {
			case c.Send <- bytes:
			default:
				log.Printf("[realtime] dropping broadcast for slow client user=%s", c.UserID)
			}
		}
	}
}
`

// The backstop closed every connection it collected, and a Client built without
// a socket (tests, server-side code) has a nil Conn, whose Close panics ten
// seconds later on a goroutine nothing recovers.
const hubBackstopOld = "\t\tfor _, conn := range conns {\n\t\t\t_ = conn.Close()\n\t\t}\n"

const hubBackstopNew = "\t\tfor _, conn := range conns {\n\t\t\tif conn != nil { // a Client made without a socket has no Conn\n\t\t\t\t_ = conn.Close()\n\t\t\t}\n\t\t}\n"

const hubPresenceFieldOld = "\tchannels channelIndex\n"

const hubPresenceFieldAdd = `
	// presence is who this process holds in presence channels. See presence.go.
	presence presenceSet
`

const hubPresenceFieldNew = hubPresenceFieldOld + hubPresenceFieldAdd

const hubPresenceStartOld = "\t\tgo h.publishLoop(ctx)\n"

const hubPresenceStartNew = hubPresenceStartOld + "\t\th.startPresence(ctx)\n"

const hubPresenceUnregisterOld = "func (h *Hub) Unregister(c *Client) {\n\th.mu.Lock()\n\tdefer h.mu.Unlock()\n"

const hubPresenceUnregisterNew = `func (h *Hub) Unregister(c *Client) {
	// Deferred before the unlock, so it runs after it: leaving a presence
	// channel can wait on Redis, which must not hold up the hub.
	defer h.leaveAllPresence(c)
	h.mu.Lock()
	defer h.mu.Unlock()
`

// channels.go hunks.

const channelsPresenceDocOld = "//\tpresence-*  authorized exactly like private-*\n"

const channelsPresenceDocNew = "//\tpresence-*  authorized like private-*, and lists who is subscribed (presence.go)\n"

const channelsContextOld = `	Params  map[string]string // "id" -> "42" for the pattern "invoices.{id}"
}
`

const channelsContextNew = `	Params  map[string]string // "id" -> "42" for the pattern "invoices.{id}"

	info *interface{} // where SetInfo writes, see presence.go
}
`

const channelsAuthorizeOld = `// authorizeChannel decides whether userID may subscribe to channel.
func authorizeChannel(userID, channel string) (err error) {
	if !validChannel(channel) {
		return ErrInvalidChannel
	}
	if strings.HasPrefix(channel, "public-") {
		return nil
	}
	authorize, params := matchChannel(channel)
	if authorize == nil {
		return ErrChannelForbidden
	}
	// The authorizer runs on the connection's read goroutine, where a panic
	// would take the whole process down rather than one request.
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[realtime] the authorizer for %s panicked: %v", channel, r)
			err = ErrChannelForbidden
		}
	}()
	if !authorize(ChannelContext{UserID: userID, Channel: channel, Params: params}) {
		return ErrChannelForbidden
	}
	return nil
}
`

const channelsAuthorizeNew = `// authorizeChannel decides whether userID may subscribe to channel, and returns
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
`

const channelsSubscribeOld = `// Subscribe authorizes c for channel and adds it. A channel the connection
// already holds is a no-op. The authorizer runs outside the hub's lock, so a
// slow database check does not stall every other connection.
func (h *Hub) Subscribe(c *Client, channel string) error {
	if err := authorizeChannel(c.UserID, channel); err != nil {
		return err
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, open := h.clients[c.UserID][c]; !open {
		return nil // closed while the authorizer ran: nothing left to deliver to
	}
	return h.channels.add(c, channel)
}

// Unsubscribe removes c from channel. Unknown channels are ignored.
func (h *Hub) Unsubscribe(c *Client, channel string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.channels.remove(c, channel)
}
`

const channelsSubscribeNew = `// Subscribe authorizes c for channel and adds it. A channel the connection
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
`

const channelsSubscribedReplyOld = "\t\th.reply(c, \"subscribed\", msg.Channel, nil)\n"

const channelsSubscribedReplyNew = channelsSubscribedReplyOld + "\t\th.replyPresenceMembers(c, msg.Channel)\n"
