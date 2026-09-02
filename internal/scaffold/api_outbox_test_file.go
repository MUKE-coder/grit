package scaffold

func outboxTestGo() string {
	return `package outbox

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// order stands in for whatever business row the caller is writing.
type order struct {
	ID     string ` + "`" + `gorm:"primarykey;size:36"` + "`" + `
	Number string
}

func newDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	// One connection: each connection to ":memory:" gets its own empty
	// database, and a relay on a second one would see no messages at all.
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db handle: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&Message{}, &order{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func count(t *testing.T, db *gorm.DB, model any, where ...any) int64 {
	t.Helper()
	var n int64
	q := db.Model(model)
	if len(where) > 0 {
		q = q.Where(where[0], where[1:]...)
	}
	if err := q.Count(&n).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	return n
}

// The half that "publish, then commit" gets wrong: a rollback must take the
// message with it, or a webhook announces an order that does not exist.
func TestRollbackTakesTheMessageWithIt(t *testing.T) {
	db := newDB(t)
	boom := errors.New("something failed after the write")

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order{ID: "o1", Number: "INV-1"}).Error; err != nil {
			return err
		}
		if err := Enqueue(tx, "orders.created", map[string]string{"id": "o1"}); err != nil {
			return err
		}
		return boom
	})
	if !errors.Is(err, boom) {
		t.Fatalf("transaction error = %v, want the injected failure", err)
	}

	if n := count(t, db, &order{}); n != 0 {
		t.Errorf("%d orders survived a rollback", n)
	}
	if n := count(t, db, &Message{}); n != 0 {
		t.Errorf("%d messages survived a rollback: a webhook would announce an order that does not exist", n)
	}
}

// The other half: a commit must keep both.
func TestCommitKeepsBoth(t *testing.T) {
	db := newDB(t)

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order{ID: "o1", Number: "INV-1"}).Error; err != nil {
			return err
		}
		return Enqueue(tx, "orders.created", map[string]string{"id": "o1"})
	})
	if err != nil {
		t.Fatalf("transaction: %v", err)
	}

	if n := count(t, db, &order{}); n != 1 {
		t.Errorf("orders = %d, want 1", n)
	}
	if n := count(t, db, &Message{}, "status = ?", StatusPending); n != 1 {
		t.Errorf("pending messages = %d, want 1", n)
	}
}

// Enqueue outside a transaction is the bug it exists to prevent, so it is
// refused rather than quietly doing the wrong thing.
func TestEnqueueRefusesTheRootHandle(t *testing.T) {
	db := newDB(t)

	err := Enqueue(db, "orders.created", map[string]string{"id": "o1"})
	if !errors.Is(err, ErrNoTransaction) {
		t.Fatalf("error = %v, want ErrNoTransaction", err)
	}
	if n := count(t, db, &Message{}); n != 0 {
		t.Errorf("%d messages written outside a transaction", n)
	}
}

// A retried request must not publish the same event twice.
func TestSameKeyEnqueuesOnce(t *testing.T) {
	db := newDB(t)

	for i := 0; i < 3; i++ {
		err := db.Transaction(func(tx *gorm.DB) error {
			return Enqueue(tx, "orders.created", map[string]int{"n": i}, Key("order-42-created"))
		})
		if err != nil {
			t.Fatalf("attempt %d: %v", i, err)
		}
	}

	if n := count(t, db, &Message{}); n != 1 {
		t.Errorf("messages = %d, want 1: the same key was queued more than once", n)
	}
}

// Messages without a key are all distinct. An empty key stored as "" would
// collide on the unique index and silently drop every message after the first.
func TestMessagesWithoutAKeyAreIndependent(t *testing.T) {
	db := newDB(t)

	for i := 0; i < 5; i++ {
		if err := db.Transaction(func(tx *gorm.DB) error {
			return Enqueue(tx, "orders.created", map[string]int{"n": i})
		}); err != nil {
			t.Fatalf("enqueue %d: %v", i, err)
		}
	}

	if n := count(t, db, &Message{}); n != 5 {
		t.Errorf("messages = %d, want 5: keyless messages collided", n)
	}
}

func TestRelayDeliversAndMarks(t *testing.T) {
	db := newDB(t)
	for i := 0; i < 3; i++ {
		if err := db.Transaction(func(tx *gorm.DB) error {
			return Enqueue(tx, "orders.created", map[string]int{"n": i})
		}); err != nil {
			t.Fatalf("enqueue: %v", err)
		}
	}

	var got []string
	r := &Relay{DB: db, Deliver: func(ctx context.Context, m Message) error {
		got = append(got, m.Topic)
		return nil
	}}

	n, err := r.Tick(context.Background())
	if err != nil {
		t.Fatalf("Tick: %v", err)
	}
	if n != 3 || len(got) != 3 {
		t.Fatalf("delivered %d (Tick reported %d), want 3", len(got), n)
	}
	if n := count(t, db, &Message{}, "status = ?", StatusDelivered); n != 3 {
		t.Errorf("delivered rows = %d, want 3", n)
	}

	// A second tick has nothing to do. Without that, a relay loop would
	// redeliver the same messages forever.
	if n, _ := r.Tick(context.Background()); n != 0 {
		t.Errorf("second tick claimed %d messages, want 0", n)
	}
}

func TestFailedDeliveryRetriesWithBackoff(t *testing.T) {
	db := newDB(t)
	if err := db.Transaction(func(tx *gorm.DB) error {
		return Enqueue(tx, "orders.created", map[string]int{"n": 1})
	}); err != nil {
		t.Fatalf("enqueue: %v", err)
	}

	r := &Relay{
		DB:          db,
		BaseBackoff: time.Minute,
		Deliver: func(ctx context.Context, m Message) error {
			return errors.New("receiver said no")
		},
	}
	if _, err := r.Tick(context.Background()); err != nil {
		t.Fatalf("Tick: %v", err)
	}

	var m Message
	if err := db.First(&m).Error; err != nil {
		t.Fatalf("reload: %v", err)
	}
	if m.Status != StatusPending {
		t.Errorf("status = %q, want pending so it is tried again", m.Status)
	}
	if m.Attempts != 1 {
		t.Errorf("attempts = %d, want 1", m.Attempts)
	}
	if m.LastError == "" {
		t.Error("last_error is empty; the reason for the retry was not recorded")
	}
	if !m.AvailableAt.After(time.Now().Add(30 * time.Second)) {
		t.Errorf("available_at = %v, want it pushed out by the backoff", m.AvailableAt)
	}

	// And it is not picked up again before the backoff elapses. A relay that
	// ignored available_at would hammer a failing endpoint at poll speed.
	if n, _ := r.Tick(context.Background()); n != 0 {
		t.Errorf("claimed %d messages during the backoff window, want 0", n)
	}
}

func TestGivesUpAfterMaxAttempts(t *testing.T) {
	db := newDB(t)
	if err := db.Transaction(func(tx *gorm.DB) error {
		return Enqueue(tx, "orders.created", map[string]int{"n": 1})
	}); err != nil {
		t.Fatalf("enqueue: %v", err)
	}

	r := &Relay{
		DB:          db,
		MaxAttempts: 3,
		BaseBackoff: time.Nanosecond,
		Deliver: func(ctx context.Context, m Message) error {
			return errors.New("still no")
		},
	}
	for i := 0; i < 3; i++ {
		if _, err := r.Tick(context.Background()); err != nil {
			t.Fatalf("tick %d: %v", i, err)
		}
		// Wait out the backoff, which is the point of the field: a relay that
		// ignored available_at would hammer a failing endpoint at poll speed.
		//
		// A millisecond rather than the nanosecond configured above because
		// the clock is not that precise. Windows advances time.Now() in steps
		// of about 100ns, so two calls inside one step return the same value
		// and a 1ns backoff has provably not elapsed. That is the clock being
		// coarse, not the relay being wrong, and it cost an afternoon to
		// establish.
		time.Sleep(time.Millisecond)
	}

	var m Message
	if err := db.First(&m).Error; err != nil {
		t.Fatalf("reload: %v", err)
	}
	if m.Status != StatusFailed {
		t.Errorf("status = %q after 3 attempts with MaxAttempts=3, want failed", m.Status)
	}
	// Parked, not deleted: the row is the evidence.
	if n := count(t, db, &Message{}); n != 1 {
		t.Errorf("messages = %d; a message that could not be delivered was thrown away", n)
	}

	// Retry puts it back for an operator who has fixed the receiver.
	if err := Retry(db, m.ID); err != nil {
		t.Fatalf("Retry: %v", err)
	}
	if n := count(t, db, &Message{}, "status = ?", StatusPending); n != 1 {
		t.Error("Retry did not return the message to pending")
	}
}

// A relay that dies holding claims must not strand them.
func TestStaleClaimsAreReclaimed(t *testing.T) {
	db := newDB(t)
	if err := db.Transaction(func(tx *gorm.DB) error {
		return Enqueue(tx, "orders.created", map[string]int{"n": 1})
	}); err != nil {
		t.Fatalf("enqueue: %v", err)
	}

	// Simulate the process that took it and never came back.
	long := time.Now().Add(-time.Hour)
	if err := db.Model(&Message{}).Where("1 = 1").Updates(map[string]any{
		"status":     StatusClaimed,
		"claimed_by": "a-process-that-died",
		"claimed_at": long,
	}).Error; err != nil {
		t.Fatalf("stale claim: %v", err)
	}

	delivered := 0
	r := &Relay{DB: db, ClaimTimeout: time.Minute, Deliver: func(ctx context.Context, m Message) error {
		delivered++
		return nil
	}}
	if _, err := r.Tick(context.Background()); err != nil {
		t.Fatalf("Tick: %v", err)
	}
	if delivered != 1 {
		t.Errorf("delivered %d, want 1: a claim from a dead process stranded the message", delivered)
	}
}

// A claim that is still fresh belongs to someone else and must be left alone.
func TestFreshClaimsAreLeftAlone(t *testing.T) {
	db := newDB(t)
	if err := db.Transaction(func(tx *gorm.DB) error {
		return Enqueue(tx, "orders.created", map[string]int{"n": 1})
	}); err != nil {
		t.Fatalf("enqueue: %v", err)
	}
	if err := db.Model(&Message{}).Where("1 = 1").Updates(map[string]any{
		"status":     StatusClaimed,
		"claimed_by": "another-relay",
		"claimed_at": time.Now(),
	}).Error; err != nil {
		t.Fatalf("claim: %v", err)
	}

	r := &Relay{DB: db, ClaimTimeout: time.Hour, Deliver: func(ctx context.Context, m Message) error {
		t.Error("delivered a message another relay is holding")
		return nil
	}}
	if n, _ := r.Tick(context.Background()); n != 0 {
		t.Errorf("claimed %d messages held by another relay, want 0", n)
	}
}

func TestPruneKeepsWhatMatters(t *testing.T) {
	db := newDB(t)
	for i := 0; i < 3; i++ {
		if err := db.Transaction(func(tx *gorm.DB) error {
			return Enqueue(tx, "orders.created", map[string]int{"n": i})
		}); err != nil {
			t.Fatalf("enqueue: %v", err)
		}
	}

	var all []Message
	if err := db.Find(&all).Error; err != nil {
		t.Fatalf("load: %v", err)
	}
	old := time.Now().Add(-48 * time.Hour)
	// One delivered long ago, one delivered just now, one still failed.
	db.Model(&Message{}).Where("id = ?", all[0].ID).
		Updates(map[string]any{"status": StatusDelivered, "delivered_at": old})
	db.Model(&Message{}).Where("id = ?", all[1].ID).
		Updates(map[string]any{"status": StatusDelivered, "delivered_at": time.Now()})
	db.Model(&Message{}).Where("id = ?", all[2].ID).
		Updates(map[string]any{"status": StatusFailed, "delivered_at": old})

	n, err := Prune(db, 24*time.Hour)
	if err != nil {
		t.Fatalf("Prune: %v", err)
	}
	if n != 1 {
		t.Errorf("pruned %d, want 1", n)
	}
	if count(t, db, &Message{}, "status = ?", StatusFailed) != 1 {
		t.Error("Prune deleted a failed message; that is a bug nobody has looked at yet")
	}
	if count(t, db, &Message{}, "status = ?", StatusDelivered) != 1 {
		t.Error("Prune deleted a message delivered inside the retention window")
	}
}

func TestCounters(t *testing.T) {
	db := newDB(t)
	for i := 0; i < 4; i++ {
		if err := db.Transaction(func(tx *gorm.DB) error {
			return Enqueue(tx, "orders.created", map[string]int{"n": i})
		}); err != nil {
			t.Fatalf("enqueue: %v", err)
		}
	}

	pending, err := Pending(db)
	if err != nil || pending != 4 {
		t.Errorf("Pending() = %d, %v; want 4, nil", pending, err)
	}
	failed, err := Failed(db)
	if err != nil || failed != 0 {
		t.Errorf("Failed() = %d, %v; want 0, nil", failed, err)
	}
}
`
}
