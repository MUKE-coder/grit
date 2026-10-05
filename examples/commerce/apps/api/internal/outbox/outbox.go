// Package outbox makes "save the row and tell the world" a single atomic act.
//
// The problem it solves has two halves, and doing the obvious thing gets one
// of them wrong whichever order you pick.
//
// Publish first, then commit: the publish succeeds, the commit fails, and a
// webhook has now announced an order that does not exist. Downstream systems
// act on it. There is nothing to reconcile against, because the row was never
// written.
//
// Commit first, then publish: the commit succeeds, the process is killed
// before the publish, and the order exists with nobody told. No error is
// logged anywhere, because from the process's point of view nothing failed.
//
// The outbox removes the choice. The message is written to a table in the same
// transaction as the business data, so it commits or rolls back with it, and a
// relay reads committed messages and delivers them afterwards. Either both
// happened or neither did.
//
// The cost is delivery semantics: at-least-once, not exactly-once. A message
// can be delivered twice if the process dies between the send and the status
// update. Consumers have to be idempotent, which is what a dedup key on the
// receiving end is for.
package outbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
)

// Status values. Strings rather than an enum so a human reading the table can
// see what happened without a lookup.
const (
	StatusPending   = "pending"
	StatusClaimed   = "claimed"
	StatusDelivered = "delivered"
	StatusFailed    = "failed"
)

// Message is one thing to be delivered, and the record that it was.
//
// An alias rather than its own type: the table is declared in
// internal/models like every other table, so AutoMigrate creates it and
// the backup writer includes it, while calling code still reads
// outbox.Message.
type Message = models.OutboxMessage

// ErrNoTransaction is returned when Enqueue is handed something that is not a
// transaction.
//
// Enqueueing outside a transaction is not a smaller version of the right
// thing; it is the "commit first, then publish" bug with extra steps, and it
// fails silently in production and never in a test. So it is refused.
var ErrNoTransaction = errors.New("outbox: Enqueue needs the transaction that writes the data")

// Enqueue writes a message using the caller's transaction.
//
//	err := db.Transaction(func(tx *gorm.DB) error {
//	    if err := tx.Create(&order).Error; err != nil {
//	        return err
//	    }
//	    return outbox.Enqueue(tx, "orders.created", order, outbox.Key(order.ID))
//	})
//
// tx must be the same handle the business write used. Passing the root *gorm.DB
// instead gives the message its own implicit transaction, which commits
// independently, which is the bug this package exists to prevent.
func Enqueue(tx *gorm.DB, topic string, payload any, opts ...Option) error {
	if tx == nil {
		return ErrNoTransaction
	}
	// GORM sets this on the handle inside a Transaction callback. Checking it
	// is the only way to tell a transaction from the root handle, and the
	// distinction is the entire correctness argument for this package.
	if !inTransaction(tx) {
		return ErrNoTransaction
	}
	if topic == "" {
		return errors.New("outbox: a message needs a topic")
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("outbox: encoding %s: %w", topic, err)
	}

	m := Message{Topic: topic, Payload: datatypes.JSON(body)}
	for _, opt := range opts {
		opt(&m)
	}

	// A key already present means this exact message is queued or sent
	// already, which is success from the caller's point of view: the event
	// gets published once either way. Checked first because a failed INSERT
	// marks the transaction aborted on Postgres, taking the business write
	// with it; the constraint below is the backstop for the race.
	if m.Key != nil {
		var existing int64
		if err := tx.Model(&Message{}).Where("key = ?", *m.Key).Count(&existing).Error; err == nil && existing > 0 {
			return nil
		}
	}
	if err := tx.Create(&m).Error; err != nil {
		if isDuplicateKey(err) {
			return nil
		}
		return fmt.Errorf("outbox: queueing %s: %w", topic, err)
	}
	return nil
}

// Option configures a message.
type Option func(*Message)

// Key sets the idempotency key. Enqueueing the same key twice writes one row.
func Key(k string) Option {
	return func(m *Message) {
		if k != "" {
			m.Key = &k
		}
	}
}

// After delays the first delivery attempt.
func After(d time.Duration) Option {
	return func(m *Message) { m.AvailableAt = time.Now().Add(d) }
}

// inTransaction reports whether this handle is inside one.
//
// GORM does not expose that directly. It does put the transaction's
// *sql.Tx in the statement's ConnPool, and the root handle holds an *sql.DB,
// so the two are distinguishable by type. It is not pretty and it is checked
// by TestEnqueueRefusesTheRootHandle, which is the part that matters.
func inTransaction(db *gorm.DB) bool {
	if db.Statement == nil || db.Statement.ConnPool == nil {
		return false
	}
	_, isTx := db.Statement.ConnPool.(gorm.TxCommitter)
	return isTx
}

func isDuplicateKey(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	// Not every driver maps to ErrDuplicatedKey, and the ones that do not say
	// so in prose. Matching on the message is unpleasant; getting a duplicate
	// wrong is worse, because it turns an idempotent retry into a 500.
	s := strings.ToLower(err.Error())
	for _, frag := range []string{"unique constraint failed", "duplicate key value", "duplicate entry"} {
		if strings.Contains(s, frag) {
			return true
		}
	}
	return false
}
