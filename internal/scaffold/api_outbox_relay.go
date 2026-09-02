package scaffold

func outboxRelayGo() string {
	return `package outbox

import (
	"context"
	"log"
	"math"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Deliver sends one message. Returning an error schedules a retry.
type Deliver func(ctx context.Context, m Message) error

// Relay drains the outbox.
//
// One goroutine, polling. Not a queue subscription, because the whole point is
// that the outbox lives in the same database as the business data: if the
// database is reachable the relay works, and there is no second system to be
// up, configured, or drifting out of sync.
type Relay struct {
	DB      *gorm.DB
	Deliver Deliver

	// Name identifies this instance in claimed_by. Defaults to the hostname,
	// which is what you want in a container: it tells you which replica has a
	// message when one is stuck.
	Name string

	// Interval is how often to look when the last poll found nothing. A poll
	// that found work goes straight round again, so a busy outbox drains at
	// the speed of delivery rather than the speed of this.
	Interval time.Duration

	// Batch is how many messages to claim at once.
	Batch int

	// MaxAttempts is how many times to try before parking a message at
	// "failed". Parked rather than deleted: a message nobody can deliver is
	// evidence of a bug, and deleting the evidence is how the bug survives.
	MaxAttempts int

	// BaseBackoff is the first retry delay; it doubles each attempt up to
	// MaxBackoff.
	BaseBackoff time.Duration
	MaxBackoff  time.Duration

	// ClaimTimeout is how long a claim is honoured before another relay may
	// take the message. It exists because a process can die holding claims,
	// and without it those messages are stuck forever.
	//
	// Set it comfortably above your slowest delivery. Too low means two relays
	// deliver the same message concurrently, which at-least-once permits but
	// nobody enjoys.
	ClaimTimeout time.Duration
}

func (r *Relay) defaults() {
	if r.Interval <= 0 {
		r.Interval = 2 * time.Second
	}
	if r.Batch <= 0 {
		r.Batch = 50
	}
	if r.MaxAttempts <= 0 {
		r.MaxAttempts = 12
	}
	if r.BaseBackoff <= 0 {
		r.BaseBackoff = 5 * time.Second
	}
	if r.MaxBackoff <= 0 {
		r.MaxBackoff = 30 * time.Minute
	}
	if r.ClaimTimeout <= 0 {
		r.ClaimTimeout = 5 * time.Minute
	}
	if r.Name == "" {
		r.Name = hostname()
	}
}

// Start runs until the context is cancelled.
func (r *Relay) Start(ctx context.Context) {
	r.defaults()
	if r.Deliver == nil {
		log.Println("outbox: no Deliver function; relay not started")
		return
	}

	log.Printf("outbox: relay %s started (batch %d, every %s)", r.Name, r.Batch, r.Interval)
	for {
		n, err := r.Tick(ctx)
		if err != nil {
			log.Printf("outbox: %v", err)
		}
		// Work found means there is probably more. Only idle polls wait.
		if n > 0 && ctx.Err() == nil {
			continue
		}
		select {
		case <-ctx.Done():
			log.Printf("outbox: relay %s stopped", r.Name)
			return
		case <-time.After(r.Interval):
		}
	}
}

// Tick claims one batch and delivers it. Returns how many it handled.
//
// Exported so a test can drive the relay a step at a time instead of racing a
// goroutine, and so an operator can drain the outbox from a one-off command.
func (r *Relay) Tick(ctx context.Context) (int, error) {
	r.defaults()

	batch, err := r.claim(ctx)
	if err != nil {
		return 0, err
	}
	for _, m := range batch {
		r.attempt(ctx, m)
	}
	return len(batch), nil
}

// claim marks a batch as ours, atomically.
//
// The select and the update are in one transaction with a row lock, because
// two relays polling the same table will otherwise read the same rows and both
// deliver them. Postgres and MySQL take the lock; SQLite serialises writes
// anyway, so the transaction alone is enough there.
//
// SKIP LOCKED would be better on Postgres and does not exist on the other two,
// and a locking clause that silently changes meaning per dialect is worse than
// a plain lock: the contended case here is two relays, not two hundred.
func (r *Relay) claim(ctx context.Context) ([]Message, error) {
	var claimed []Message
	now := time.Now()

	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var batch []Message
		q := tx.Model(&Message{}).
			Where("status = ? AND available_at <= ?", StatusPending, now).
			// A claim older than the timeout belonged to a process that is no
			// longer running. Without this the messages it held are stranded.
			Or("status = ? AND claimed_at < ?", StatusClaimed, now.Add(-r.ClaimTimeout)).
			Order("created_at").
			Limit(r.Batch)

		if r.DB.Dialector.Name() != "sqlite" {
			q = q.Clauses(clause.Locking{Strength: "UPDATE"})
		}
		if err := q.Find(&batch).Error; err != nil {
			return err
		}
		if len(batch) == 0 {
			return nil
		}

		ids := make([]string, len(batch))
		for i, m := range batch {
			ids[i] = m.ID
		}
		if err := tx.Model(&Message{}).Where("id IN ?", ids).Updates(map[string]any{
			"status":     StatusClaimed,
			"claimed_by": r.Name,
			"claimed_at": now,
		}).Error; err != nil {
			return err
		}
		claimed = batch
		return nil
	})
	if err != nil {
		return nil, err
	}
	return claimed, nil
}

// attempt delivers one message and records what happened.
func (r *Relay) attempt(ctx context.Context, m Message) {
	err := r.Deliver(ctx, m)
	now := time.Now()

	if err == nil {
		r.DB.Model(&Message{}).Where("id = ?", m.ID).Updates(map[string]any{
			"status":       StatusDelivered,
			"delivered_at": now,
			"attempts":     m.Attempts + 1,
			"last_error":   "",
			"claimed_by":   "",
		})
		return
	}

	attempts := m.Attempts + 1
	update := map[string]any{
		"attempts":   attempts,
		"last_error": err.Error(),
		"claimed_by": "",
		"claimed_at": nil,
	}
	if attempts >= r.MaxAttempts {
		update["status"] = StatusFailed
		log.Printf("outbox: %s (%s) failed after %d attempts: %v", m.ID, m.Topic, attempts, err)
	} else {
		update["status"] = StatusPending
		update["available_at"] = now.Add(r.backoff(attempts))
	}
	r.DB.Model(&Message{}).Where("id = ?", m.ID).Updates(update)
}

// backoff doubles per attempt, capped.
func (r *Relay) backoff(attempt int) time.Duration {
	d := float64(r.BaseBackoff) * math.Pow(2, float64(attempt-1))
	if d > float64(r.MaxBackoff) || math.IsInf(d, 0) {
		return r.MaxBackoff
	}
	return time.Duration(d)
}

// Retry puts a failed message back in the queue, immediately.
//
// For the operator who has fixed whatever was rejecting it. Attempts are not
// reset, so the history of how hard this fought stays readable.
func Retry(db *gorm.DB, id string) error {
	return db.Model(&Message{}).Where("id = ? AND status = ?", id, StatusFailed).
		Updates(map[string]any{
			"status":       StatusPending,
			"available_at": time.Now(),
			"claimed_by":   "",
			"claimed_at":   nil,
		}).Error
}

// Prune deletes delivered messages older than the given age.
//
// Only delivered ones. A failed message is a bug that has not been looked at
// yet, and a pending one has not been sent; deleting either would lose an
// event nobody ever saw.
func Prune(db *gorm.DB, olderThan time.Duration) (int64, error) {
	res := db.Where("status = ? AND delivered_at < ?", StatusDelivered, time.Now().Add(-olderThan)).
		Delete(&Message{})
	return res.RowsAffected, res.Error
}

// hostname names this relay in claimed_by, so a stuck message says which
// replica is holding it.
func hostname() string {
	h, err := os.Hostname()
	if err != nil || h == "" {
		return "relay"
	}
	return h
}

// Pending counts what is waiting, for the health card and the admin page.
func Pending(db *gorm.DB) (int64, error) {
	var n int64
	err := db.Model(&Message{}).Where("status IN ?", []string{StatusPending, StatusClaimed}).Count(&n).Error
	return n, err
}

// Failed counts what has given up.
func Failed(db *gorm.DB) (int64, error) {
	var n int64
	err := db.Model(&Message{}).Where("status = ?", StatusFailed).Count(&n).Error
	return n, err
}
`
}
