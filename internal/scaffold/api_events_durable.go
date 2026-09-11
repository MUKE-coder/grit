package scaffold

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Durable event delivery: the event bus's third mode, carried by the outbox.
//
// Sync and Async subscribers live in the process. An Async event is gone if the
// process stops, dropped when the queue is full, and not retried when the
// subscriber fails, which is right for a realtime push and wrong for "reserve
// the stock this purchase order was approved for". The outbox already had the
// delivery half (claims, backoff, parking), and nothing in a scaffolded project
// ever started its relay.

// APIEventsDurableGo, APIOutboxGo, APIOutboxRelayGo and APIOutboxModelGo are
// exported so grit generate can bring a project that predates them up to date:
// a generated transition service calls events.EmitTx.
func APIEventsDurableGo() string { return apiEventsDurableGo() }

// APIOutboxGo emits internal/outbox/outbox.go.
func APIOutboxGo() string { return outboxGo() }

// APIOutboxRelayGo emits internal/outbox/relay.go.
func APIOutboxRelayGo() string { return outboxRelayGo() }

// APIOutboxModelGo emits internal/models/outbox_message.go.
func APIOutboxModelGo() string { return apiOutboxModelGo() }

// apiEventsDurableGo emits internal/events/durable.go.
func apiEventsDurableGo() string {
	return `package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/outbox"
)

// Durable subscribers are delivered through the transactional outbox.
//
// Sync and Async subscribers live in this process's memory. An Async event
// still in the queue is gone if the process stops, dropped if the queue is
// full, and not retried if the subscriber fails. That is fine for a realtime
// push and wrong for work another module depends on, such as "record the spend
// against the department's budget", where a lost event leaves two modules
// disagreeing and nothing to say so.
//
// A Durable subscriber gets each event as a row in outbox_messages. When the
// emitter has a transaction (EmitTx, which every workflow transition uses) the
// row is written in it, so the event reaches the subscriber if and only if the
// change commits. The relay StartRelay runs then delivers it: retried with
// backoff until the subscriber returns nil, parked as "failed" after the
// relay's MaxAttempts, and picked up by another replica if this one dies
// holding it.
//
//	events.On("purchase_requests.fulfill", events.Durable, "record-spend", recordSpend)
//
// The name is the delivery address, so keep it unique and stable: a queued
// event is delivered to the subscriber registered under that name when it is
// taken off the queue, which may be after a deploy.
//
// Delivery is at least once. A subscriber that crashes after doing its work
// but before returning is run again, so it must be idempotent: check whether
// the work is already done, or key it on the event (e.Name + e.ID).
//
// The subscriber receives the event decoded from JSON: Before and After are
// maps rather than your model, and C is nil. DecodeAfter turns After back into
// a struct.

const durableTopicPrefix = "event:"

var (
	durableMu   sync.RWMutex
	durableDB   *gorm.DB
	relayWarned sync.Once
)

// StartRelay starts delivering Durable subscribers. Call once at boot, after
// Init. Every replica runs one; the outbox's row claims keep two of them from
// delivering the same message at once.
func StartRelay(db *gorm.DB) {
	durableMu.Lock()
	durableDB = db
	durableMu.Unlock()
	relay := &outbox.Relay{DB: db, Deliver: deliverDurable, TopicPrefix: durableTopicPrefix}
	go relay.Start(context.Background())
}

// OnDurable registers a Durable subscriber that writes to the database, which
// most do. It runs in a transaction the relay opens, so its writes commit only
// if it returns nil and a failed attempt leaves nothing half done for the retry
// to trip over. Registering from init() is fine: the database is supplied at
// delivery, by StartRelay.
//
//	func init() {
//	    events.OnDurable("purchase_requests.fulfill", "record-spend", recordSpend)
//	}
//
//	func recordSpend(tx *gorm.DB, e events.Event) error { ... }
func OnDurable(pattern, name string, h func(tx *gorm.DB, e Event) error) {
	On(pattern, Durable, name, func(e Event) error {
		durableMu.RLock()
		db := durableDB
		durableMu.RUnlock()
		if db == nil {
			return fmt.Errorf("the event relay has no database: call events.StartRelay(db) at boot")
		}
		return db.Transaction(func(tx *gorm.DB) error { return h(tx, e) })
	})
}

// EmitTx queues the event for every matching Durable subscriber in tx, the
// transaction that made the change. It marks e so the Emit after the commit,
// which runs the Sync and Async subscribers, does not queue it again:
//
//	var ev events.Event
//	err := db.Transaction(func(tx *gorm.DB) error {
//	    // ...the change...
//	    ev = events.Event{Name: "orders.paid", Resource: "orders", ID: order.ID, After: order}
//	    return events.EmitTx(tx, c, &ev)
//	})
//	if err == nil {
//	    events.Emit(c, ev)
//	}
//
// An error means the event could not be queued, and returning it rolls the
// change back with it, which is the point.
func EmitTx(tx *gorm.DB, c *gin.Context, e *Event) error {
	e.fill(c)
	bus := defaultBus
	if bus == nil {
		return nil
	}
	names := bus.durableFor(*e)
	for _, name := range names {
		if err := outbox.Enqueue(tx, durableTopicPrefix+name, e); err != nil {
			return fmt.Errorf("events: queueing %s for %s: %w", e.Name, name, err)
		}
	}
	e.durable = true
	if len(names) > 0 {
		warnWithoutRelay()
	}
	return nil
}

// enqueueDurable is the path for an event emitted without a transaction, such
// as a generated create, update or delete, which has already committed. The
// rows are written in their own transaction straight afterwards, so the event
// is lost only if the process dies in between. EmitTx closes that gap too.
func (b *Bus) enqueueDurable(e Event) {
	durableMu.RLock()
	db := durableDB
	durableMu.RUnlock()
	if db == nil {
		warnWithoutRelay()
		b.mu.Lock()
		b.dropped++
		b.mu.Unlock()
		return
	}
	names := b.durableFor(e)
	err := db.Transaction(func(tx *gorm.DB) error {
		for _, name := range names {
			if err := outbox.Enqueue(tx, durableTopicPrefix+name, e); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[events] could not queue %s for its Durable subscribers: %v", e.Name, err)
	}
}

// durableFor names the Durable subscribers an event goes to.
func (b *Bus) durableFor(e Event) []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	var names []string
	for _, s := range b.subs {
		if s.delivery == Durable && matches(s.pattern, e.Name, e.Resource) {
			names = append(names, s.name)
		}
	}
	return names
}

// durableHandler finds a Durable subscriber by name.
func (b *Bus) durableHandler(name string) Handler {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, s := range b.subs {
		if s.delivery == Durable && s.name == name {
			return s.handler
		}
	}
	return nil
}

// deliverDurable is the relay's delivery function. An error, a panic included,
// leaves the message queued for another attempt.
func deliverDurable(ctx context.Context, m outbox.Message) (err error) {
	name := strings.TrimPrefix(m.Topic, durableTopicPrefix)
	var e Event
	if err := json.Unmarshal(m.Payload, &e); err != nil {
		return fmt.Errorf("decoding outbox message %s: %w", m.ID, err)
	}
	bus := defaultBus
	if bus == nil {
		return fmt.Errorf("the event bus is not running")
	}
	h := bus.durableHandler(name)
	if h == nil {
		// Kept and retried rather than dropped: a subscriber that was renamed or
		// removed in a deploy leaves its queue as evidence, parked as failed.
		return fmt.Errorf("no Durable subscriber named %q is registered", name)
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("subscriber %q panicked: %v", name, r)
		}
	}()
	return h(e)
}

// warnWithoutRelay says once that Durable events are queueing with nothing to
// deliver them.
func warnWithoutRelay() {
	durableMu.RLock()
	running := durableDB != nil
	durableMu.RUnlock()
	if running {
		return
	}
	relayWarned.Do(func() {
		log.Println("[events] Durable subscribers are registered but the relay is not running. " +
			"Add events.StartRelay(db) after events.Init in routes.go; events are kept in outbox_messages until then.")
	})
}

// DecodeAfter fills v with the event's After, for a Durable subscriber that
// wants its model back rather than a map:
//
//	var pr models.PurchaseRequest
//	if err := e.DecodeAfter(&pr); err != nil {
//	    return err
//	}
func (e Event) DecodeAfter(v interface{}) error {
	raw, err := json.Marshal(e.After)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, v)
}

// fill completes an event from the request that caused it.
func (e *Event) fill(c *gin.Context) {
	if e.At.IsZero() {
		e.At = time.Now().UTC()
	}
	if e.Label == "" {
		e.Label = e.ID
	}
	if c == nil {
		return
	}
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
`
}

// EnsureEventRelay starts the Durable relay in a routes.go that predates it,
// on the line after the subscribers are registered. Only once durable.go is
// there, so the edit never leaves routes.go calling something undefined.
// Returns whether it changed the file.
func EnsureEventRelay(apiRoot string) (bool, error) {
	path := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	raw, err := os.ReadFile(path)
	if err != nil {
		return false, nil
	}
	if !fileContains(filepath.Join(apiRoot, "internal", "events", "durable.go"), "func StartRelay(") {
		return false, nil
	}
	crlf := strings.Contains(string(raw), "\r\n")
	content := strings.ReplaceAll(string(raw), "\r\n", "\n")
	if strings.Contains(content, "events.StartRelay(") {
		return false, nil
	}
	loc := registerSubscribersLine.FindStringSubmatchIndex(content)
	if loc == nil {
		return false, fmt.Errorf("could not find services.RegisterEventSubscribers in routes.go, so Durable " +
			"subscribers will queue and never be delivered. Add events.StartRelay(db) after it")
	}
	indent := content[loc[2]:loc[3]]
	content = content[:loc[1]] + indent + "events.StartRelay(db)\n" + content[loc[1]:]
	if _, err := format.Source([]byte(content)); err != nil {
		return false, fmt.Errorf("routes.go left alone: starting the event relay did not produce valid Go (%v)", err)
	}
	if crlf {
		content = strings.ReplaceAll(content, "\n", "\r\n")
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return false, fmt.Errorf("writing %s: %w", path, err)
	}
	manifest.Refresh(path)
	return true, nil
}

var registerSubscribersLine = regexp.MustCompile(`(?m)^([ \t]*)services\.RegisterEventSubscribers\([^\n]*\)\n`)
