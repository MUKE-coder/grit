package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeOutboxFiles writes internal/outbox: the transactional outbox.
func writeOutboxFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "outbox", "outbox.go"):      outboxGo(),
		filepath.Join(apiRoot, "internal", "outbox", "relay.go"):       outboxRelayGo(),
		filepath.Join(apiRoot, "internal", "outbox", "outbox_test.go"): outboxTestGo(),
	}
	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

func outboxGo() string {
	return `// Package outbox makes "save the row and tell the world" a single atomic act.
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

	"{{MODULE}}/internal/models"
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
`
}

func apiOutboxModelGo() string {
	return `package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"{{MODULE}}/internal/ids"
)

// Message is one thing to be delivered, and the record that it was.
//
// Kept after delivery rather than deleted. The row is the evidence that an
// event was published, which is what you want when a downstream system says it
// never received one, and it is what a replay reads from.
type OutboxMessage struct {
	ID    string ` + "`" + `gorm:"type:varchar(36);primaryKey" json:"id"` + "`" + `
	Topic string ` + "`" + `gorm:"size:255;not null;index:idx_outbox_topic" json:"topic"` + "`" + `

	// Key is the caller's idempotency key for this message, unique across the
	// table. Enqueueing the same key twice is a no-op rather than an error, so
	// a retried request cannot publish the same event twice.
	//
	// Empty means "no key", and an empty string cannot be unique across many
	// rows, so it is stored as NULL. That is why this is a pointer.
	Key *string ` + "`" + `gorm:"size:255;uniqueIndex" json:"key,omitempty"` + "`" + `

	Payload datatypes.JSON ` + "`" + `gorm:"type:json" json:"payload"` + "`" + `

	Status   string ` + "`" + `gorm:"size:20;not null;default:'pending';index:idx_outbox_claim,priority:1" json:"status"` + "`" + `
	Attempts int    ` + "`" + `gorm:"not null;default:0" json:"attempts"` + "`" + `
	LastError string ` + "`" + `gorm:"type:text" json:"last_error,omitempty"` + "`" + `

	// AvailableAt is when the relay may next try. It moves forward on each
	// failure, which is the backoff.
	AvailableAt time.Time ` + "`" + `gorm:"not null;index:idx_outbox_claim,priority:2" json:"available_at"` + "`" + `

	ClaimedBy   string     ` + "`" + `gorm:"size:64" json:"claimed_by,omitempty"` + "`" + `
	ClaimedAt   *time.Time ` + "`" + `json:"claimed_at,omitempty"` + "`" + `
	DeliveredAt *time.Time ` + "`" + `json:"delivered_at,omitempty"` + "`" + `

	CreatedAt time.Time ` + "`" + `gorm:"index" json:"created_at"` + "`" + `
	UpdatedAt time.Time ` + "`" + `json:"updated_at"` + "`" + `
}

func (OutboxMessage) TableName() string { return "outbox_messages" }

// BeforeCreate fills the id and the first availability.
func (m *OutboxMessage) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	if m.AvailableAt.IsZero() {
		m.AvailableAt = time.Now()
	}
	if m.Status == "" {
		m.Status = "pending"
	}
	return nil
}`
}
