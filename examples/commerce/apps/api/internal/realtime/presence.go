package realtime

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
	UserID   string          `json:"user_id"`
	Info     json.RawMessage `json:"info,omitempty"`
	JoinedAt time.Time       `json:"joined_at"`
}

type presenceLeft struct {
	UserID string `json:"user_id"`
}

type presenceSnapshot struct {
	Members []Member `json:"members"`
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
	joined map[*Client]map[string]struct{}    // connection -> its presence channels
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
		h.offer(c, bytes, kind, channel)
	}
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

const presenceLua = `
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
`

// KEYS[1] channel hash; ARGV field, user, member JSON, ttl ms. Returns 1 when
// no fresh entry holds the user, this replica's own included, so a heartbeat that rewrites a join still in flight does not announce it twice.
var presenceJoinScript = redis.NewScript(presenceLua + `
local now, ttl = now_ms(), tonumber(ARGV[4])
local first = not present(KEYS[1], ARGV[2], '', now, ttl)
redis.call('HSET', KEYS[1], ARGV[1], string.format('%d', now) .. ' ' .. ARGV[3])
redis.call('PEXPIRE', KEYS[1], ttl)
if first then return 1 end
return 0
`)

// ARGV field, user, ttl ms. Returns 1 when that was the user's last entry.
var presenceLeaveScript = redis.NewScript(presenceLua + `
if redis.call('HDEL', KEYS[1], ARGV[1]) == 0 then return 0 end
if present(KEYS[1], ARGV[2], '', now_ms(), tonumber(ARGV[3])) then return 0 end
return 1
`)

// ARGV ttl ms, then fields. Renews the fields that exist, returns the rest.
var presenceRefreshScript = redis.NewScript(presenceLua + `
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
`)

// ARGV ttl ms. Removes stale fields, returns the users left with none.
var presenceSweepScript = redis.NewScript(presenceLua + `
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
`)

// ARGV ttl ms. Returns the member JSON of every fresh field.
var presenceMembersScript = redis.NewScript(presenceLua + `
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
`)
