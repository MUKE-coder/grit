package scaffold

// Domain events: one bus, so a handler says what happened once instead of
// telling four systems separately.
//
// Before this, a generated handler called services.LogCreate and nothing else.
// Webhooks, realtime and background jobs all existed and none of them were
// ever fired by generated code: dispatching a webhook meant hand-writing the
// call in every handler, for every operation, on every resource, and
// forgetting one produced no error at all, just a subscriber that never heard
// anything.
//
// The audit log now subscribes to the bus rather than being called directly,
// so behaviour is unchanged and three more systems can listen.

// APIEventsGo and APIEventsSubscribersGo are exported so grit generate
// resource can add the package to a project that predates it. The generated
// handler emits, so the bus has to be there.
func APIEventsGo() string { return apiEventsGo() }

// APIEventsSubscribersGo emits the default subscriber wiring.
func APIEventsSubscribersGo() string { return apiEventsSubscribersGo() }

// apiEventsGo emits internal/events/events.go.
func apiEventsGo() string {
	return `// Package events is the application's domain event bus.
//
// A handler emits once:
//
//	events.Emit(c, events.Event{
//	    Name: "invoices.paid", Resource: "invoices", ID: inv.ID, After: inv,
//	})
//
// and every interested system hears about it. The audit log, outbound
// webhooks and realtime clients are all subscribers rather than things the
// handler has to remember to call.
package events

import (
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Delivery says whether a subscriber runs inside the request or after it.
type Delivery int

const (
	// Sync runs the subscriber inline, before the handler responds. Use it
	// only for work that must be part of the request: the audit record, which
	// needs the request context and should be written before the caller is
	// told the write succeeded. A slow Sync subscriber is a slow endpoint.
	Sync Delivery = iota
	// Async runs the subscriber on the bus's worker pool after the handler has
	// returned. Everything with a network call belongs here: a webhook
	// endpoint that takes four seconds must not make the API take four
	// seconds.
	Async
)

// Event is something that happened, described once.
type Event struct {
	// Name is the topic: "invoices.created", "invoices.paid". Subscribers
	// match on it exactly, by resource prefix ("invoices.*"), or with "*".
	Name string ` + "`" + `json:"name"` + "`" + `
	// Resource is the plural snake_case resource name, so a subscriber can
	// route without parsing Name.
	Resource string ` + "`" + `json:"resource"` + "`" + `
	// Entity is the Go type name ("Invoice"), for human-facing output.
	Entity string ` + "`" + `json:"entity"` + "`" + `
	// ID is the affected row.
	ID string ` + "`" + `json:"id"` + "`" + `
	// Label is a human-readable identifier: an invoice number, a person's
	// name. Falls back to the ID.
	Label string ` + "`" + `json:"label"` + "`" + `
	// Actor is the user who caused it, empty for system events.
	Actor string ` + "`" + `json:"actor,omitempty"` + "`" + `
	// Before and After carry the row either side of the change. Before is nil
	// on a create, After is nil on a delete.
	Before interface{} ` + "`" + `json:"before,omitempty"` + "`" + `
	After  interface{} ` + "`" + `json:"after,omitempty"` + "`" + `
	// Detail is a short summary of what changed, for the activity feed.
	Detail string ` + "`" + `json:"detail,omitempty"` + "`" + `
	// Meta carries anything else a subscriber might want. Request IP and user
	// agent are put here by Emit.
	Meta map[string]interface{} ` + "`" + `json:"meta,omitempty"` + "`" + `
	// At is when it happened.
	At time.Time ` + "`" + `json:"at"` + "`" + `

	// C is the request context, and it is set only for Sync subscribers.
	//
	// Async subscribers receive it as nil, deliberately. A gin context is
	// cancelled and recycled once the handler returns, so an async subscriber
	// reading it would see a context belonging to somebody else's request.
	// Making it nil turns a subtle data race into an obvious nil pointer the
	// first time anyone tries.
	C *gin.Context ` + "`" + `json:"-"` + "`" + `
}

// Handler is a subscriber.
type Handler func(e Event) error

type subscription struct {
	pattern  string
	delivery Delivery
	name     string
	handler  Handler
}

// Bus fans events out to subscribers. One per process, created at boot.
type Bus struct {
	mu   sync.RWMutex
	subs []subscription

	queue  chan Event
	wg     sync.WaitGroup
	closed bool

	// dropped counts events discarded because the async queue was full.
	// Exposed through Stats so "did my webhook fire" has an answer other than
	// a shrug.
	dropped int64
}

// QueueSize is how many async events may be waiting at once.
//
// Bounded rather than unbounded on purpose: an unbounded queue turns a slow
// subscriber into memory exhaustion, which fails later and worse than
// dropping events and saying so.
const QueueSize = 1024

// NewBus returns a bus with workers running.
func NewBus(workers int) *Bus {
	if workers < 1 {
		workers = 4
	}
	b := &Bus{queue: make(chan Event, QueueSize)}
	for i := 0; i < workers; i++ {
		b.wg.Add(1)
		go b.worker()
	}
	return b
}

// On registers a subscriber.
//
// Patterns are exact ("invoices.paid"), resource-wide ("invoices.*") or
// everything ("*").
func (b *Bus) On(pattern string, delivery Delivery, name string, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subs = append(b.subs, subscription{
		pattern: pattern, delivery: delivery, name: name, handler: h,
	})
}

// Emit publishes an event.
//
// It never returns an error and never panics the request. A subscriber that
// fails is logged and the others still run: a webhook endpoint being down is
// not a reason to fail the write that has already happened.
func (b *Bus) Emit(c *gin.Context, e Event) {
	if e.At.IsZero() {
		e.At = time.Now().UTC()
	}
	if e.Label == "" {
		e.Label = e.ID
	}
	if c != nil {
		if e.Actor == "" {
			if v, ok := c.Get("user_id"); ok {
				if s, ok := v.(string); ok {
					e.Actor = s
				}
			}
		}
		if e.Meta == nil {
			e.Meta = map[string]interface{}{}
		}
		e.Meta["ip"] = c.ClientIP()
		e.Meta["user_agent"] = c.Request.UserAgent()
	}

	b.mu.RLock()
	subs := make([]subscription, len(b.subs))
	copy(subs, b.subs)
	b.mu.RUnlock()

	var wantAsync bool
	for _, s := range subs {
		if !matches(s.pattern, e.Name, e.Resource) {
			continue
		}
		if s.delivery == Sync {
			sync := e
			sync.C = c
			b.run(s, sync)
		} else {
			wantAsync = true
		}
	}

	if !wantAsync {
		return
	}
	// The async copy carries no request context, for the reason on Event.C.
	async := e
	async.C = nil
	b.enqueue(async)
}

func (b *Bus) enqueue(e Event) {
	b.mu.RLock()
	closed := b.closed
	b.mu.RUnlock()
	if closed {
		return
	}
	select {
	case b.queue <- e:
	default:
		// Full. Dropping is the least bad option: blocking would make the
		// request wait on a subscriber, and growing would trade a visible
		// problem for an invisible one.
		b.mu.Lock()
		b.dropped++
		b.mu.Unlock()
		log.Printf("[events] queue full, dropped %s (%s)", e.Name, e.ID)
	}
}

func (b *Bus) worker() {
	defer b.wg.Done()
	for e := range b.queue {
		b.mu.RLock()
		subs := make([]subscription, len(b.subs))
		copy(subs, b.subs)
		b.mu.RUnlock()

		for _, s := range subs {
			if s.delivery != Async || !matches(s.pattern, e.Name, e.Resource) {
				continue
			}
			b.run(s, e)
		}
	}
}

// run calls one subscriber, absorbing panics.
//
// A subscriber is third-party-ish code as far as the bus is concerned, and a
// panic in a webhook formatter should not take down the process or the
// request that emitted the event.
func (b *Bus) run(s subscription, e Event) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[events] subscriber %q panicked on %s: %v", s.name, e.Name, r)
		}
	}()
	if err := s.handler(e); err != nil {
		log.Printf("[events] subscriber %q failed on %s: %v", s.name, e.Name, err)
	}
}

// Stats reports the bus's health, for /api/health and grit doctor.
type Stats struct {
	Subscribers int   ` + "`" + `json:"subscribers"` + "`" + `
	Queued      int   ` + "`" + `json:"queued"` + "`" + `
	Capacity    int   ` + "`" + `json:"capacity"` + "`" + `
	Dropped     int64 ` + "`" + `json:"dropped"` + "`" + `
}

// Stats returns a snapshot. Dropped rising is the signal that a subscriber is
// too slow or the queue too small, and it is the only way to find that out
// from outside.
func (b *Bus) Stats() Stats {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return Stats{
		Subscribers: len(b.subs),
		Queued:      len(b.queue),
		Capacity:    cap(b.queue),
		Dropped:     b.dropped,
	}
}

// Close drains the queue and stops the workers. Called on shutdown so events
// already accepted are delivered rather than lost at the door.
func (b *Bus) Close() {
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		return
	}
	b.closed = true
	b.mu.Unlock()

	close(b.queue)
	b.wg.Wait()
}

// matches reports whether a pattern selects an event.
func matches(pattern, name, resource string) bool {
	switch {
	case pattern == "*":
		return true
	case pattern == name:
		return true
	case len(pattern) > 2 && pattern[len(pattern)-2:] == ".*":
		prefix := pattern[:len(pattern)-2]
		return prefix == resource ||
			(len(name) > len(prefix) && name[:len(prefix)+1] == prefix+".")
	}
	return false
}

// ── Package-level bus ────────────────────────────────────────────────────
//
// Handlers call events.Emit(c, ...) without holding a reference, the same way
// they call log.Printf. There is one bus per process, created at boot in
// routes.Setup, and threading it into every handler struct would be a lot of
// plumbing for a singleton.

var defaultBus *Bus

// Init creates the process bus. Call once, at boot, before Register.
func Init(workers int) *Bus {
	defaultBus = NewBus(workers)
	return defaultBus
}

// Default returns the process bus, or nil before Init.
func Default() *Bus { return defaultBus }

// Emit publishes on the process bus. A no-op before Init, so a handler can
// emit safely in a test that never set one up.
func Emit(c *gin.Context, e Event) {
	if defaultBus == nil {
		return
	}
	defaultBus.Emit(c, e)
}

// On subscribes on the process bus.
func On(pattern string, delivery Delivery, name string, h Handler) {
	if defaultBus == nil {
		return
	}
	defaultBus.On(pattern, delivery, name, h)
}

// Emitted is the common case: a CRUD or transition event on a resource.
//
// Generated handlers call this rather than building an Event literal, so the
// name and the resource cannot drift apart and every resource emits the same
// shape.
func Emitted(c *gin.Context, resource, entity, action, id, label, detail string, before, after interface{}) {
	Emit(c, Event{
		Name:     resource + "." + action,
		Resource: resource,
		Entity:   entity,
		ID:       id,
		Label:    label,
		Detail:   detail,
		Before:   before,
		After:    after,
	})
}
`
}
