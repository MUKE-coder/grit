package scaffold

// Sync policy: the declaration that turns offline behaviour from convention
// into something a developer states and a tool can check.
//
// Before this, every registered model was mirrored in full, conflicts were
// always manual, there was no age limit and no way to keep a field off the
// wire. Those are defensible defaults, but nobody chose them and nobody could
// find out what they had. That is the difference between a feature and magic.
//
// The policy lives on the server registry rather than in client config on
// purpose. A conflict rule enforced only in the client is not a rule: an old
// build, a second client, or a crafted request ignores it. The client reads
// the policy to render the right UI; the server reads the same policy to
// decide what actually happens.

// APISyncPolicyGo is exported so grit generate resource can add this file to
// a project generated before policies existed. Everything policy-related lives
// here rather than in registry.go precisely so that is one new file and not an
// edit to an existing one.
func APISyncPolicyGo() string { return apiSyncPolicyGo() }

// APISyncHandlerGo is exported for the same reason, though this one replaces
// an existing file and so is only written when the manifest proves nobody has
// edited it.
func APISyncHandlerGo() string { return apiSyncHandlerGo() }

// apiSyncPolicyGo emits internal/sync/policy.go.
func apiSyncPolicyGo() string {
	return `package sync

import (
	"fmt"
	"sync"
	"time"
)

// Conflict strategies. What happens when a client pushes an update whose
// version does not match the row on disk.
const (
	// ConflictManual parks the change and hands both versions back so a human
	// decides. The default, because silently discarding somebody's work should
	// be something you opt into.
	ConflictManual = "manual"
	// ConflictServerWins discards the client's change and returns the server
	// row. For records the server is authoritative about: stock levels,
	// prices, anything a back office owns.
	ConflictServerWins = "server_wins"
	// ConflictClientWins applies the client's change over the server's. For
	// records whose author is the only person who edits them, where the
	// version check is protecting nothing.
	ConflictClientWins = "client_wins"
)

// Sync modes.
const (
	// ModeOfflineFirst mirrors the model locally and queues writes.
	ModeOfflineFirst = "offline_first"
	// ModeOnlineOnly leaves the model out of the mirror entirely. Registered
	// so it can be named, and so a client asking for it gets a clear answer
	// rather than an empty list.
	ModeOnlineOnly = "online_only"
)

// Policy is how one model behaves offline.
//
// The zero value is the behaviour every project had before policies existed:
// offline-first, manual conflicts, every field, no age limit. Adding a policy
// to a project changes nothing until somebody sets a field.
type Policy struct {
	Mode     string ` + "`" + `json:"mode"` + "`" + `
	Conflict string ` + "`" + `json:"conflict"` + "`" + `

	// Fields, when set, is the allowlist sent to clients. Identity and
	// bookkeeping columns (id, version, created_at, updated_at, deleted_at)
	// are always included: without them the client cannot version a row, order
	// it, or know it was deleted.
	Fields []string ` + "`" + `json:"fields,omitempty"` + "`" + `

	// LocalOnly names fields that exist only on the device. They are stripped
	// from anything a client pushes, so a draft note or a scratch flag can
	// live in the mirror without ever reaching the server.
	LocalOnly []string ` + "`" + `json:"local_only,omitempty"` + "`" + `

	// MaxOfflineAge is how long a client may go without a successful sync
	// before its data should be treated as stale. Zero means no limit.
	//
	// This is advisory by design. The server cannot enforce it, because a
	// client that has not synced is by definition not talking to the server.
	// What it buys is a client that can say so instead of showing
	// three-day-old stock levels as though they were current.
	MaxOfflineAge time.Duration ` + "`" + `json:"-"` + "`" + `

	// MaxOfflineAgeSeconds is the wire form of MaxOfflineAge.
	MaxOfflineAgeSeconds int ` + "`" + `json:"max_offline_age_seconds,omitempty"` + "`" + `
}

// DefaultPolicy is what a model registered without one gets.
func DefaultPolicy() Policy {
	return Policy{Mode: ModeOfflineFirst, Conflict: ConflictManual}
}

// alwaysIncluded are the columns a client cannot work without, whatever the
// field allowlist says.
var alwaysIncluded = map[string]bool{
	"id":         true,
	"version":    true,
	"created_at": true,
	"updated_at": true,
	"deleted_at": true,
	"_deleted":   true,
}

// Normalise fills in defaults and returns an error for anything unusable.
//
// Called at registration rather than at request time, so a typo in a conflict
// strategy fails at boot with the model name attached rather than on the
// first push from somebody's phone.
func (p *Policy) Normalise(table string) error {
	if p.Mode == "" {
		p.Mode = ModeOfflineFirst
	}
	if p.Conflict == "" {
		p.Conflict = ConflictManual
	}

	switch p.Mode {
	case ModeOfflineFirst, ModeOnlineOnly:
	default:
		return fmt.Errorf("sync: %s has unknown mode %q (want %s or %s)",
			table, p.Mode, ModeOfflineFirst, ModeOnlineOnly)
	}

	switch p.Conflict {
	case ConflictManual, ConflictServerWins, ConflictClientWins:
	default:
		return fmt.Errorf("sync: %s has unknown conflict strategy %q (want %s, %s or %s)",
			table, p.Conflict, ConflictManual, ConflictServerWins, ConflictClientWins)
	}

	for _, f := range p.LocalOnly {
		if alwaysIncluded[f] {
			return fmt.Errorf("sync: %s cannot mark %q local_only: the client needs it to version and order rows", table, f)
		}
	}

	// A field named in both lists is a contradiction the author has to
	// resolve. Guessing which one they meant is how a local_only field ends
	// up on the wire.
	inFields := map[string]bool{}
	for _, f := range p.Fields {
		inFields[f] = true
	}
	for _, f := range p.LocalOnly {
		if inFields[f] {
			return fmt.Errorf("sync: %s lists %q in both fields and local_only", table, f)
		}
	}

	if p.MaxOfflineAge > 0 {
		p.MaxOfflineAgeSeconds = int(p.MaxOfflineAge / time.Second)
	}
	return nil
}

// The policy store lives beside the registry rather than inside it.
//
// That is what makes this feature one new file instead of an edit to
// registry.go, which in turn is what lets grit generate resource add policy
// support to a project generated before policies existed. There is exactly one
// Registry per process, created once at boot from routes.Setup, so a
// package-level store and a field on the struct are the same thing here. A
// test building two registries would see them share policies; nothing in a
// scaffolded project does that.
var policyStore = struct {
	mu sync.RWMutex
	m  map[string]Policy
}{m: make(map[string]Policy)}

// RegisterWithPolicy registers a model and states how it behaves offline.
//
// A policy that will not normalise panics rather than being quietly corrected.
// This runs at boot, so the failure is loud, immediate, and names the model,
// instead of surfacing as a wrong conflict decision on somebody's phone weeks
// later.
func (r *Registry) RegisterWithPolicy(table string, proto interface{}, policy Policy) {
	if err := policy.Normalise(table); err != nil {
		panic(err)
	}
	r.Register(table, proto)
	policyStore.mu.Lock()
	defer policyStore.mu.Unlock()
	policyStore.m[table] = policy
}

// PolicyFor returns a table's policy, or the default for one registered
// without one. Never fails: an unregistered table is caught by New.
func (r *Registry) PolicyFor(table string) Policy {
	policyStore.mu.RLock()
	defer policyStore.mu.RUnlock()
	if p, ok := policyStore.m[table]; ok {
		return p
	}
	return DefaultPolicy()
}

// Policies returns every declared policy, for GET /api/sync/policy. Clients
// configure themselves from this rather than from a copy that can drift.
func (r *Registry) Policies() map[string]Policy {
	policyStore.mu.RLock()
	defer policyStore.mu.RUnlock()
	out := make(map[string]Policy, len(policyStore.m))
	for k, v := range policyStore.m {
		out[k] = v
	}
	// Tables registered without a policy still have one, and a client that
	// only heard about the explicit ones would apply defaults it was never
	// told about. Say all of them.
	for _, table := range r.Tables() {
		if _, ok := out[table]; !ok {
			out[table] = DefaultPolicy()
		}
	}
	return out
}

// Projects returns the row filtered to what this policy sends to clients.
//
// A nil or empty allowlist means everything, which is what every project had
// before policies existed.
func (p Policy) Projects(row map[string]interface{}) map[string]interface{} {
	if len(p.Fields) == 0 {
		return row
	}
	allowed := make(map[string]bool, len(p.Fields))
	for _, f := range p.Fields {
		allowed[f] = true
	}
	out := make(map[string]interface{}, len(p.Fields)+len(alwaysIncluded))
	for k, v := range row {
		if allowed[k] || alwaysIncluded[k] {
			out[k] = v
		}
	}
	return out
}

// StripLocalOnly removes device-only fields from an incoming push payload.
//
// Enforced here rather than trusted to the client: a field declared local_only
// is a promise that it never reaches the server, and a promise kept only by
// well-behaved clients is not one.
func (p Policy) StripLocalOnly(data map[string]interface{}) map[string]interface{} {
	if len(p.LocalOnly) == 0 || data == nil {
		return data
	}
	out := make(map[string]interface{}, len(data))
	local := make(map[string]bool, len(p.LocalOnly))
	for _, f := range p.LocalOnly {
		local[f] = true
	}
	for k, v := range data {
		if !local[k] {
			out[k] = v
		}
	}
	return out
}
`
}
