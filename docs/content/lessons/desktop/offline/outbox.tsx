import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'

export default function Lesson() {
  return (
    <>
      <p>
        The outbox is a tiny table that says &quot;these rows haven&apos;t been
        sent to the server yet&quot;. When you write a sale locally, you also
        write an outbox row. When the sync engine runs (next lesson), it
        drains the outbox.
      </p>

      <h2>The pattern</h2>
      <Diagram label="Outbox write flow" caption="Transactionally: the business write + the outbox entry happen together. Either both succeed or both fail.">
{`   ┌─────────────────────────────────────────────┐
   │                                             │
   │   React calls RecordSale(items)             │
   │                                             │
   │              │                              │
   │              ▼                              │
   │   ┌──────────────────────┐                  │
   │   │   BEGIN TRANSACTION  │                  │
   │   ├──────────────────────┤                  │
   │   │ INSERT INTO sales    │  ◄── business    │
   │   ├──────────────────────┤                  │
   │   │ INSERT INTO outbox   │  ◄── sync intent │
   │   │  (entity='sale',     │                  │
   │   │   entity_id=<uuid>,  │                  │
   │   │   op='create')       │                  │
   │   ├──────────────────────┤                  │
   │   │   COMMIT             │                  │
   │   └──────────────────────┘                  │
   │              │                              │
   │              ▼                              │
   │      Return sale to React                   │
   │      (sync happens later, in background)    │
   │                                             │
   └─────────────────────────────────────────────┘`}
      </Diagram>

      <h2>The outbox model</h2>
      <CodeBlock
        language="go"
        filename="internal/models/outbox.go"
        code={`type OutboxEntry struct {
    ID         int64     \`gorm:"primaryKey;autoIncrement"\`   // local-only
    Entity     string    \`gorm:"not null;index"\`            // 'sale', 'product', …
    EntityID   string    \`gorm:"not null;index"\`            // the business-row UUID
    Operation  string    \`gorm:"not null"\`                  // 'create' | 'update' | 'delete'
    Payload    datatypes.JSON                                // the row, serialised
    Attempts   int       \`gorm:"default:0"\`
    LastError  string
    CreatedAt  time.Time
}`}
      />
      <p>
        Notice: <code>ID</code> is auto-increment here (local-only) — gives
        a clean FIFO order. <code>EntityID</code> is the UUID of the
        actual business row (a Sale, a Product), used to look up the
        latest state when syncing.
      </p>

      <h2>Writing through the outbox</h2>
      <CodeBlock
        language="go"
        filename="internal/service/sales.go"
        code={`func (s *Service) RecordSale(sale *Sale) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(sale).Error; err != nil {
            return err
        }
        payload, _ := json.Marshal(sale)
        return tx.Create(&OutboxEntry{
            Entity:    "sale",
            EntityID:  sale.ID,
            Operation: "create",
            Payload:   payload,
        }).Error
    })
}`}
      />
      <p>
        Both inserts in one transaction. If either fails, both roll back.
        You can never have a sale without a matching outbox entry, or vice
        versa.
      </p>

      <h2>The unsynced indicator</h2>
      <p>
        For UX, show the user what hasn&apos;t synced yet:
      </p>
      <CodeBlock
        language="go"
        code={`func (s *Service) PendingCount() (int64, error) {
    var n int64
    err := s.db.Model(&OutboxEntry{}).Where("attempts < 10").Count(&n).Error
    return n, err
}`}
      />
      <CodeBlock
        language="tsx"
        code={`function SyncStatus() {
  const { data } = useQuery({
    queryKey: ['pending'],
    queryFn: () => PendingCount(),
    refetchInterval: 3000,
  })
  return data === 0
    ? <span className="text-green-500">✓ synced</span>
    : <span className="text-yellow-500">{data} pending sync</span>
}`}
      />
      <p>
        Cashier glances at the corner of the screen, knows whether they&apos;re
        caught up.
      </p>

      <h2>Compacting + squashing</h2>
      <p>
        If a product gets updated 50 times offline, do you really need 50
        outbox rows? No — only the final state matters. A compaction step
        squashes the trail:
      </p>
      <CodeBlock
        language="text"
        code={`Before squash:                After squash:
  outbox: product update 1     outbox: product update (latest)
  outbox: product update 2
  outbox: product update 3
  outbox: product update 4`}
      />
      <p>
        For create-then-delete, both rows can be dropped — the server never
        needed to know. Grit ships a squash function; call it before sync.
      </p>

      <TipBox tone="warning">
        <strong>Don&apos;t delete outbox rows on success without a marker.</strong>{' '}
        If the sync confirms but the local commit fails, the outbox row
        sticks around and re-sends next time. That&apos;s the desired
        behaviour for at-least-once delivery. The server-side handler must
        be idempotent (use the EntityID for dedup).
      </TipBox>

      <KnowledgeCheck
        question="A sale write succeeded but the outbox insert failed (disk full). What's the right behaviour?"
        choices={[
          {
            label: 'Roll back the sale — both writes succeed or both fail',
            correct: true,
            feedback:
              "Right — that's exactly why both inserts are in one transaction. A sale without an outbox entry will never sync; you'd have ghost rows the server never sees. Roll back the whole thing and tell the user.",
          },
          {
            label: 'Keep the sale, log the outbox failure, sync later',
            feedback:
              "Wrong — without the outbox row, the sync engine doesn\'t know the sale exists. It stays local forever.",
          },
          {
            label: 'Insert a placeholder outbox row',
            feedback:
              "If the disk is full, the placeholder won\'t insert either. And lying about state is worse than failing.",
          },
          {
            label: 'Retry the outbox insert separately',
            feedback:
              "Possible but messy. The transaction is the right tool — atomic, simple, correct.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Add the outbox to your Note workflow:</p>
            <ol>
              <li>Add the <code>OutboxEntry</code> model + AutoMigrate it.</li>
              <li>
                Wrap your <code>AddNote</code> service in a transaction
                that inserts both the note AND the outbox entry.
              </li>
              <li>
                Wire a <code>PendingCount</code> method bound to React.
                Show the count in the corner.
              </li>
              <li>
                Add 5 notes. Confirm the count goes from 0 → 1 → 2 → … → 5.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            <code>db.Transaction(func(tx *gorm.DB) error {`{`} ... {`}`})</code>{' '}
            wraps both inserts. Return the error and GORM rolls back.
          </>
        }
        solution={
          <>
            <p>
              After 5 notes you should see &quot;5 pending sync&quot; in the corner
              (because nothing&apos;s syncing yet). Once we wire the sync
              engine next lesson, the count drops back to 0 when online.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Outbox is filling. Last lesson of the chapter — the{' '}
        <strong>sync engine</strong> that drains it.
      </p>
    </>
  )
}
