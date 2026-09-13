import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/outbox')

export default function OutboxPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Backend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">The transactional outbox</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                A webhook that fires for an order the database rolled back is worse than one that
                never fires. The outbox writes the message in the same transaction as the business
                data, so the two commit together or neither does, and a relay delivers it
                afterwards.
              </p>
            </div>

            <div className="prose-grit">
              <h2>Enqueue inside the transaction</h2>
              <p>
                <code>outbox.Enqueue</code> takes the transaction, not the database. That is the
                whole point: the row lands with the write it is about.
              </p>
              <CodeBlock
                language="go"
                filename="internal/services/order.go"
                code={`func (s *OrderService) Place(in OrderInput) error {
    return s.DB.Transaction(func(tx *gorm.DB) error {
        order := models.Order{Total: in.Total, UserID: in.UserID}
        if err := tx.Create(&order).Error; err != nil {
            return err
        }

        // Same transaction. If the order does not commit, neither does this.
        return outbox.Enqueue(tx, "orders.created", order, outbox.Key("order:"+order.ID))
    })
}`}
              />
              <p>
                <code>outbox.Key</code> makes the message idempotent: the key is unique in the
                table, so a retry of the same operation enqueues nothing the second time rather
                than sending twice. <code>outbox.After</code> delays a message.
              </p>

              <h2>A relay delivers it</h2>
              <p>
                Nothing is sent until a relay claims the row and calls your{' '}
                <code>Deliver</code>. Grit starts one relay of its own, for the durable event bus,
                and it takes only the topics that begin with <code>event:</code>. Every other topic
                is yours, which means <strong>a topic of your own needs a relay of your own</strong>.
              </p>
              <CodeBlock
                language="go"
                filename="cmd/relay/main.go"
                code={`relay := &outbox.Relay{
    DB:          db,
    TopicPrefix: "orders.",          // leave empty to take every topic
    Interval:    time.Second,        // how often to look when the last poll was empty
    Batch:       50,
    MaxAttempts: 12,                 // then the message is parked at "failed"
    BaseBackoff: time.Second,        // doubles each attempt, up to MaxBackoff
    MaxBackoff:  5 * time.Minute,
    Deliver: func(ctx context.Context, m outbox.Message) error {
        return postToWebhook(ctx, m.Topic, m.Payload)
    },
}
go relay.Start(ctx)`}
              />
              <p>
                Return an error from <code>Deliver</code> and the message is retried with an
                exponential backoff. After <code>MaxAttempts</code> it is parked at{' '}
                <code>failed</code> rather than deleted: a message nobody can deliver is evidence
                of a bug, and deleting the evidence is how the bug survives.
              </p>

              <Callout type="warning" title="A topic with no relay is a queue that only grows">
                The write succeeds, the transaction commits, and the row sits at{' '}
                <code>pending</code> with zero attempts forever. Nothing fails and nothing logs.{' '}
                <code>grit doctor</code> checks for this: it reads the topics your code enqueues
                and the relays your code starts, and reports any topic no relay covers.
              </Callout>

              <h2>Running more than one</h2>
              <p>
                A relay claims a batch under its own name and honours the claim for{' '}
                <code>ClaimTimeout</code>, so several replicas can run the same relay without
                delivering each other&apos;s messages. Set the timeout comfortably above your
                slowest delivery: too low and a slow send is retried by a second relay while the
                first is still going.
              </p>
              <p>
                Ordering is per message, not global. Messages are claimed oldest first, and a
                failing message is retried later without blocking the ones behind it, which is what
                you want for webhooks and not what you want if two messages must arrive in order.
                When order matters, put the sequence in the payload and let the receiver sort.
              </p>

              <h2>What the table looks like</h2>
              <p>
                <code>outbox_messages</code> carries the topic, the key, the JSON payload, the
                status, the attempt count, the last error, and the claim. It is ordinary SQL, so
                the backlog is a query rather than a dashboard:
              </p>
              <CodeBlock
                language="sql"
                code={`-- Anything stuck?
SELECT topic, status, count(*), max(attempts)
FROM outbox_messages
GROUP BY topic, status;

-- Parked messages, with the reason.
SELECT topic, key, attempts, last_error
FROM outbox_messages
WHERE status = 'failed'
ORDER BY updated_at DESC;`}
              />

              <Callout type="note" title="Verified under load">
                2,200 messages enqueued by 32 concurrent workers, delivered by a relay whose
                Deliver refused a third of the time: 2,200 delivered, 1,091 refusals retried, every
                key delivered exactly once, drained in ten seconds. The retry path is exercised in{' '}
                <code>internal/outbox</code>&apos;s own tests, which ship into your project.
              </Callout>

              <div className="mt-16 flex items-center justify-between border-t border-border/30 pt-8">
                <Link href="/docs/backend/append-only">
                  <Button variant="outline" size="sm" className="gap-2">
                    <ArrowLeft className="h-4 w-4" />
                    Append-only Records
                  </Button>
                </Link>
                <Link href="/docs/backend/workflows">
                  <Button variant="outline" size="sm" className="gap-2">
                    Workflows
                    <ArrowRight className="h-4 w-4" />
                  </Button>
                </Link>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
