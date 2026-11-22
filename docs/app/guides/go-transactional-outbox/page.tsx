import Link from 'next/link'
import type { Metadata } from 'next'
import { GuidePage, Section, H3, Outcome, Note, Warn, Checklist, Faq, Related } from '@/components/guide'
import { CodeBlock } from '@/components/code-block'
import { Diagram, DiagramBox, DiagramArrow } from '@/components/diagram'

export const metadata: Metadata = {
  title: 'The transactional outbox pattern in Go with Postgres',
  description:
    'Why dual writes lose events, a minimal outbox in Postgres, SKIP LOCKED for the relay, idempotency keys, backoff and parked failures.',
  alternates: { canonical: 'https://gritframework.dev/guides/go-transactional-outbox' },
  openGraph: {
    title: 'The transactional outbox pattern in Go with Postgres',
    description: 'One transaction, two outcomes, and the table that makes them agree.',
    url: 'https://gritframework.dev/guides/go-transactional-outbox',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="go-transactional-outbox"
      lead={
        <>
          <p className="mb-4">
            You write a row and publish an event. Two systems, no shared transaction, so there is
            a moment where one has happened and the other has not. The process can die in that
            moment, and at any volume, it eventually will.
          </p>
          <p>
            The outbox makes the event part of the same transaction as the write, so they commit
            together or not at all, and a separate process delivers it afterwards. It is one
            table and about a hundred lines, and it converts a silent data-loss bug into a
            visible queue.
          </p>
        </>
      }
    >
      <Outcome>
        <p>An outbox table, a relay with <code>SKIP LOCKED</code>, and delivery you can inspect.</p>
        <p>
          At-least-once delivery, with the consumer-side idempotency that makes at-least-once
          usable.
        </p>
        <p>A clear sense of when this is overkill, which is often.</p>
      </Outcome>

      <Section id="problem" title="The dual write">
        <CodeBlock
          code={`func (s *OrderService) Place(ctx context.Context, o *models.Order) error {
	if err := s.db.WithContext(ctx).Create(o).Error; err != nil {
		return fmt.Errorf("create order: %w", err)
	}

	// ←── the process can die here ──→

	return s.queue.Enqueue(ctx, jobs.NewOrderPlaced(o.ID))
}`}
        />
        <p>Three outcomes, and only one of them is the one you wanted:</p>
        <Checklist
          items={[
            <>
              Both succeed. Fine, and this is almost always what happens, which is why the bug
              survives testing.
            </>,
            <>
              The insert fails. Fine: nothing happened, the caller gets an error.
            </>,
            <>
              <strong>The insert succeeds and the enqueue does not.</strong> The order exists. No
              confirmation email, no warehouse notification, no analytics event, and nothing
              anywhere records that this happened.
            </>,
          ]}
        />
        <p>
          Reversing the order does not help: now you can send a confirmation for an order that was
          never created, which is worse. Retrying the enqueue in a loop does not help either,
          because the process can die during the retries.
        </p>
        <Warn title="Enqueuing inside the transaction is a different bug">
          <p>
            Moving the enqueue inside <code>db.Transaction</code> looks like a fix and is not:
            Redis knows nothing about your transaction, so a worker can pick the job up before
            the commit, fail to find the row, and consume a retry. Worse, if the transaction then
            rolls back, you have published an event for something that never existed.
          </p>
        </Warn>
      </Section>

      <Section id="pattern" title="The pattern">
        <Diagram>
          <DiagramBox title="One transaction" sub="insert the order AND insert the outbox row" tone="primary" />
          <DiagramArrow label="commit, or neither" />
          <DiagramBox title="outbox_events" sub="pending rows, in the same database" />
          <DiagramArrow label="a relay polls" />
          <DiagramBox title="Publish" sub="queue, webhook, broker" />
          <DiagramArrow label="on success" />
          <DiagramBox title="Mark published" sub="or increment attempts and back off" tone="primary" />
        </Diagram>
        <p>
          The insight is small: the event is a row, so it is covered by the same transaction as
          everything else. Delivery becomes a separate, retryable problem that cannot lose
          anything, because the record of what is owed survives a crash.
        </p>
      </Section>

      <Section id="table" title="The table">
        <CodeBlock
          language="sql"
          filename="migrations/0043_outbox.up.sql"
          code={`CREATE TABLE outbox_events (
  id             bigserial PRIMARY KEY,
  aggregate_type text        NOT NULL,       -- 'order'
  aggregate_id   bigint      NOT NULL,       -- 8821
  event_type     text        NOT NULL,       -- 'order.placed'
  payload        jsonb       NOT NULL,
  created_at     timestamptz NOT NULL DEFAULT now(),

  published_at   timestamptz,
  attempts       int         NOT NULL DEFAULT 0,
  next_attempt_at timestamptz NOT NULL DEFAULT now(),
  last_error     text
);

-- The relay's only query: unpublished, due, oldest first. A partial
-- index keeps it small, because the published rows are the vast
-- majority and the relay never looks at them.
CREATE INDEX idx_outbox_pending
  ON outbox_events (next_attempt_at, id)
  WHERE published_at IS NULL;`}
        />
        <p>
          Keeping the published rows is worth it: they are an audit trail of everything the system
          emitted, and the question &quot;did we send that&quot; is asked often. Prune them on a
          schedule rather than deleting on publish.
        </p>
      </Section>

      <Section id="write" title="Writing the event">
        <CodeBlock
          filename="internal/outbox/outbox.go"
          code={`// The tx argument is the whole point. A Write that opened its own
// transaction would be a dual write with extra steps.
func Write(ctx context.Context, tx *gorm.DB, e Event) error {
	payload, err := json.Marshal(e.Payload)
	if err != nil {
		return fmt.Errorf("marshal %s: %w", e.Type, err)
	}
	return tx.WithContext(ctx).Create(&models.OutboxEvent{
		AggregateType: e.AggregateType,
		AggregateID:   e.AggregateID,
		EventType:     e.Type,
		Payload:       payload,
	}).Error
}`}
        />
        <CodeBlock
          code={`func (s *OrderService) Place(ctx context.Context, o *models.Order) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(o).Error; err != nil {
			return fmt.Errorf("create order: %w", err)
		}
		// Same transaction. Both, or neither. There is no window.
		return outbox.Write(ctx, tx, outbox.Event{
			AggregateType: "order",
			AggregateID:   o.ID,
			Type:          "order.placed",
			Payload:       OrderPlaced{OrderID: o.ID, Total: o.Total},
		})
	})
}`}
        />
        <Note title="Put the data in the payload, not just the id">
          <p>
            This is the one place the usual advice reverses. A background job should carry an id
            and read the current row. An event describes something that happened at a moment, and
            a consumer reading the row later sees a state that may have changed three times
            since. Put what the event means in the payload, and version it.
          </p>
        </Note>
      </Section>

      <Section id="relay" title="The relay">
        <CodeBlock
          filename="internal/outbox/relay.go"
          code={`func (r *Relay) tick(ctx context.Context) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var batch []models.OutboxEvent

		// FOR UPDATE SKIP LOCKED is what lets several relay instances run
		// at once: each takes rows nobody else is holding, instead of
		// queueing behind them. Without it, two relays serialise, and
		// without the lock at all, two relays publish the same event.
		err := tx.Raw(\`
			SELECT * FROM outbox_events
			WHERE published_at IS NULL AND next_attempt_at <= now()
			ORDER BY id
			LIMIT 100
			FOR UPDATE SKIP LOCKED\`).Scan(&batch).Error
		if err != nil {
			return fmt.Errorf("claim batch: %w", err)
		}

		for i := range batch {
			e := &batch[i]
			if err := r.publish(ctx, e); err != nil {
				// Not a failure of the tick. Record it, back off, carry on
				// with the rest of the batch.
				r.backoff(tx, e, err)
				continue
			}
			now := time.Now()
			if err := tx.Model(e).Update("published_at", now).Error; err != nil {
				return fmt.Errorf("mark published %d: %w", e.ID, err)
			}
		}
		return nil
	})
}

func (r *Relay) backoff(tx *gorm.DB, e *models.OutboxEvent, cause error) {
	e.Attempts++
	// Exponential with a ceiling, and jitter so a provider coming back
	// up is not hit by every parked event in the same second.
	delay := time.Duration(math.Pow(2, float64(e.Attempts))) * time.Second
	if delay > 15*time.Minute {
		delay = 15 * time.Minute
	}
	delay += time.Duration(rand.Int63n(int64(delay / 4)))

	tx.Model(e).Updates(map[string]any{
		"attempts":        e.Attempts,
		"next_attempt_at": time.Now().Add(delay),
		"last_error":      cause.Error(),
	})
}`}
        />
        <p>
          Poll every second or two. For lower latency, have the writer send a notification after
          commit so the relay wakes immediately; keep the poll as the fallback, because the
          notification is the part that can be lost and the poll is the part that cannot.
        </p>
        <CodeBlock
          language="sql"
          code={`-- Optional: the relay LISTENs, the writer NOTIFYs after commit.
-- Latency drops from seconds to milliseconds, and correctness still
-- rests on the poll underneath.
NOTIFY outbox_new;`}
        />
      </Section>

      <Section id="consumer" title="At-least-once means the consumer must cope">
        <p>
          An outbox guarantees that an event is delivered at least once. It cannot guarantee
          exactly once: the publish can succeed and the process die before the row is marked, and
          the next tick will send it again. That is not a flaw to engineer around, it is the
          contract.
        </p>
        <CodeBlock
          code={`// Every event carries its id, and the consumer records what it has
// processed. The unique index does the arbitration, so two concurrent
// deliveries cannot both proceed.
type ProcessedEvent struct {
	EventID     uint64    \`gorm:"primaryKey"\`
	ProcessedAt time.Time
}

func (c *Consumer) Handle(ctx context.Context, e Event) error {
	return c.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&ProcessedEvent{EventID: e.ID, ProcessedAt: time.Now()}).Error
		if isUniqueViolation(err) {
			return nil   // already done, by us or by a duplicate delivery
		}
		if err != nil {
			return fmt.Errorf("claim event %d: %w", e.ID, err)
		}
		// The work, in the same transaction as the claim, so a crash
		// rolls back both and the event is retried cleanly.
		return c.apply(ctx, tx, e)
	})
}`}
        />
        <p>
          Claim and work in one transaction. Claim first and work after, in a separate
          transaction, and a crash between them loses the event permanently, which is the bug the
          outbox existed to prevent, reintroduced at the other end.
        </p>
      </Section>

      <Section id="operating" title="Operating it">
        <Checklist
          items={[
            <>
              <strong>Alert on the oldest pending event,</strong> not on the count. A thousand
              events enqueued a second ago is healthy; one event pending for two hours is a stuck
              relay or a permanently failing publish.
            </>,
            <>
              <strong>Park after a threshold.</strong> After roughly ten attempts, stop retrying
              and flag it. An event retrying forever consumes capacity and hides the queue depth
              that would have told you something was wrong.
            </>,
            <>
              <strong>Prune published rows on a schedule.</strong> The partial index keeps the
              relay fast regardless, but the table itself will grow without bound.
            </>,
            <>
              <strong>Order is per aggregate, not global.</strong> If a consumer needs{' '}
              <code>order.placed</code> before <code>order.shipped</code>, either publish in id
              order within an aggregate and process serially per aggregate, or make the consumer
              tolerate arrival out of order. Global ordering across a parallel relay does not
              exist.
            </>,
            <>
              <strong>Watch the payload size.</strong> An outbox row with a megabyte of JSON in it
              turns your main database into a message broker with poor ergonomics.
            </>,
          ]}
        />
      </Section>

      <Section id="when" title="When not to do this">
        <p>
          The outbox is for when losing an event is a real problem. Much of what people put
          through one is not.
        </p>
        <Checklist
          items={[
            <>
              <strong>A welcome email</strong> can be lost occasionally. Log the enqueue failure
              and move on; an outbox here is machinery for an incident that costs an apology.
            </>,
            <>
              <strong>Cache invalidation</strong> usually has a TTL underneath it, which is a
              correctness floor the outbox is not adding to.
            </>,
            <>
              <strong>Anything with a reconciliation job</strong> already has a recovery path. A
              nightly sweep that finds orders with no shipment is simpler than an outbox and
              covers more failure modes.
            </>,
          ]}
        />
        <p>
          Use it for the events that move money, notify another system of record, or cannot be
          reconstructed later: payment captured, order placed, subscription cancelled, document
          signed.
        </p>
      </Section>

      <Section id="with-grit" title="Generated">
        <CodeBlock
          terminal
          code={`# The outbox table, its migration and the relay are scaffolded into
# every project. Writing an event is a call, not a command:
#
#   outbox.Write(ctx, tx, outbox.Event{Type: "order.placed", ...})
grit migrate`}
        />
        <p>
          That writes the table, the migration, the relay with <code>SKIP LOCKED</code> and
          jittered backoff, the parked-event handling, and an admin screen showing pending,
          retrying and parked events with their last error and a retry button.
        </p>
        <p>
          The screen is the part worth having. An outbox with no visibility is a queue you cannot
          reason about, and the first time something is stuck you will want the last error rather
          than a row count.
        </p>
        <p>
          <Link href="/docs/backend/outbox" className="text-accent hover:underline">
            The outbox reference
          </Link>{' '}
          covers the generated relay and the admin screen.
        </p>
      </Section>

      <Faq
        slug="go-transactional-outbox"
        items={[
          {
            q: 'Why not use a queue that supports transactions?',
            a: (
              <>
                Very few do across a different database, and those that do, such as a
                Postgres-backed queue like River, are in effect an outbox with the relay built in.
                That is a good answer: if your queue lives in the same Postgres as your data, you
                can enqueue inside the transaction and the whole problem disappears. The pattern
                here is what you need when the queue is Redis, a broker, or someone else&apos;s
                HTTP endpoint.
              </>
            ),
            text:
              'Few queues support a transaction across a different database. Those that do, such as a Postgres-backed queue, are effectively an outbox with the relay built in, which is a good answer: if the queue lives in the same Postgres as the data, you can enqueue inside the transaction. The pattern is for when the queue is Redis, a broker or an external endpoint.',
          },
          {
            q: 'Does this give me exactly-once delivery?',
            a: (
              <>
                No, and nothing does across a network. It gives at-least-once with no loss, which
                is the strongest guarantee available. Exactly-once <em>processing</em> is
                achievable, and it is achieved at the consumer, with a processed-events table and
                a unique index, as above.
              </>
            ),
            text:
              'No, and nothing does across a network. It provides at-least-once delivery with no loss, which is the strongest available guarantee. Exactly-once processing is achievable at the consumer, using a processed-events table with a unique index.',
          },
          {
            q: 'How do I keep events in order?',
            a: (
              <>
                Per aggregate, by publishing in id order and processing one aggregate at a time.
                Global ordering is not achievable with a parallel relay, and wanting it is usually
                a sign the consumer should be tolerant of order instead. Designing consumers to
                handle arrival out of order is cheaper than enforcing a total order.
              </>
            ),
            text:
              'Per aggregate, by publishing in id order and processing one aggregate at a time. Global ordering is not achievable with a parallel relay, and needing it usually means the consumer should tolerate out-of-order arrival instead.',
          },
          {
            q: 'What about change data capture instead?',
            a: (
              <>
                Reading the write-ahead log with Debezium captures every change with no
                application code, which is powerful and couples your consumers to your schema:
                they see table rows rather than domain events, so a column rename becomes a
                breaking change for them. An outbox publishes an event you designed and can
                version. Worth the extra table for anything consumers outside your team rely on.
              </>
            ),
            text:
              'Reading the write-ahead log captures every change with no application code, but couples consumers to your schema: they see rows rather than domain events, so a column rename breaks them. An outbox publishes an event you designed and can version, which is worth the extra table when consumers outside your team depend on it.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/go-background-jobs-asynq',
            label: 'Background jobs in Go',
            blurb: 'The queue this delivers into, and why enqueuing in a transaction is wrong.',
          },
          {
            href: '/guides/go-saga-pattern',
            label: 'Sagas and compensating transactions',
            blurb: 'When the work spans several services and some of it has to be undone.',
          },
          {
            href: '/docs/backend/outbox',
            label: 'Outbox reference',
            blurb: 'The generated table, relay and admin screen.',
          },
          {
            href: '/guides/go-webhooks-hmac-idempotency',
            label: 'Webhooks: verification and idempotency',
            blurb: 'The same at-least-once problem, seen from the receiving end.',
          },
        ]}
      />
    </GuidePage>
  )
}
