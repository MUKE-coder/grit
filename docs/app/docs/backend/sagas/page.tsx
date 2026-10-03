import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/sagas')

export default function SagasPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Backend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Sagas</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                A transaction is the right tool when every write is in one database. It is no help
                when the steps are in four places: charge a card, reserve stock, book a courier,
                send the receipt. A saga is the answer, and{' '}
                <code>grit generate workflow</code> writes one.
              </p>
            </div>

            <div className="prose-grit">
              <h2>The problem it solves</h2>
              <p>
                The card does not roll back when the courier refuses. The process holding all of
                this in its head is exactly the thing that crashes, usually between the charge and
                the reservation, and what is left is a customer charged for nothing and no record of
                why.
              </p>
              <p>
                A saga gives each step a <code>Do</code> and an <code>Undo</code>. The steps run in
                order, and when one fails for good the completed ones are undone newest first. Where
                the run got to is a row in <code>saga_runs</code>, not a stack frame, so a process
                that dies resumes rather than losing the thread.
              </p>

              <h2>Generating one</h2>
              <CodeBlock
                language="bash"
                code={`grit generate workflow Checkout --steps "charge,reserve_stock,book_courier,send_receipt"`}
              />
              <p>
                That writes <code>internal/sagas/checkout.go</code> with the steps stubbed out and
                registers it. Each <code>Do</code> returns an error until you write it, so a
                half-built saga fails loudly rather than reporting success and doing nothing: the
                first you would otherwise hear of it is a customer saying the parcel never came.
              </p>
              <CodeBlock
                language="go"
                filename="internal/sagas/checkout.go"
                code={`func Checkout() saga.Definition {
    return saga.Definition{
        Name: "checkout",
        Steps: []saga.Step{
            {
                Name: "charge",
                Do: func(ctx context.Context, r *saga.Run) error {
                    var in CheckoutInput
                    if err := r.Input(&in); err != nil {
                        return err
                    }
                    charge, err := payments.Charge(ctx, in.OrderID, payments.Idempotency(r.IdempotencyKey()))
                    if err != nil {
                        if payments.Declined(err) {
                            return saga.Fatal(err) // retrying will not fix a declined card
                        }
                        return err
                    }
                    r.Set("charge_id", charge.ID)
                    return nil
                },
                Undo: func(ctx context.Context, r *saga.Run) error {
                    id := r.GetString("charge_id")
                    if id == "" {
                        return nil // the charge never landed, so there is nothing to refund
                    }
                    return payments.Refund(ctx, id)
                },
            },
            // ...
        },
    }
}`}
              />

              <h2>Starting a run</h2>
              <CodeBlock
                language="go"
                code={`run, err := saga.Start(ctx, db, "checkout", CheckoutInput{OrderID: order.ID},
    saga.Key("order:"+order.ID))`}
              />
              <p>
                It records the run and returns at once. A runner advances it, and every replica runs
                one: a run is claimed before it is touched, so two of them cannot execute the same
                step and charge a card twice. The key makes the start idempotent, which is what you
                want behind a retried request or a redelivered webhook.
              </p>
              <p>
                Start it in the same transaction as the write that justifies it and the two commit
                together, the way <a href="/docs/backend/outbox">outbox.Enqueue</a> does.
              </p>

              <h2>The three rules</h2>
              <ol>
                <li>
                  <strong>Every Do must be idempotent.</strong> A crash between the side effect and
                  the record of it is not preventable, so a resumed run will sometimes repeat a step.
                  Pass <code>r.IdempotencyKey()</code> to whatever you are calling: it is stable for
                  a given run and step, so a provider that deduplicates on it will not charge twice.
                </li>
                <li>
                  <strong>Every Undo must be idempotent too</strong>, and must tolerate a{' '}
                  <code>Do</code> that never finished. Compensation runs <em>because</em> something
                  went wrong, and that includes not knowing whether the charge landed. Check before
                  you reverse.
                </li>
                <li>
                  <strong>Steps you cannot undo go last.</strong> An email cannot be unsent, so{' '}
                  <code>Undo: nil</code> is a legitimate thing to write. The consequence is the
                  ordering: a step after it that fails will compensate everything before it and leave
                  that email sent. Sending the receipt before the parcel is booked is a bug you only
                  find in production.
                </li>
              </ol>

              <Callout type="warning" title="State is a column, not a closure">
                The charge id the refund needs has to come from <code>r.Set</code> and{' '}
                <code>r.Get</code>, because by the time the refund runs the process that made the
                charge is usually gone. What a step writes is saved even when that step then fails,
                which is deliberate: a charge id obtained just before a timeout is exactly what the
                refund needs.
              </Callout>

              <h2>What the statuses mean</h2>
              <ul>
                <li>
                  <code>running</code>: working through the steps.
                </li>
                <li>
                  <code>compensating</code>: a step failed for good and the completed ones are being
                  undone.
                </li>
                <li>
                  <code>done</code>: every step completed.
                </li>
                <li>
                  <code>compensated</code>: a step failed and everything before it was undone. The
                  world is back where it started, which is the good outcome of a bad run.
                </li>
                <li>
                  <code>stuck</code>: a compensation itself failed past its attempts. This is the one
                  that needs a person. Something happened that could not be taken back, and the run
                  says so instead of carrying on with a status that reads as resolved.
                </li>
              </ul>
              <p>
                A compensation is never skipped, because skipping one is how money goes missing. A
                failed step is retried with exponential backoff and jitter;{' '}
                <code>saga.Fatal(err)</code> skips the retries for a failure that retrying will not
                fix.
              </p>

              <h2>Not the same as a workflow field</h2>
              <p>
                <a href="/docs/backend/workflows">Workflows</a> turn a status column into a state
                machine: the states are one record&apos;s own, and the transitions are usually a
                person clicking a button. A saga is the other thing: work that crosses systems, with
                nobody watching, that has to either finish or be undone. A project often wants both,
                and they do not overlap.
              </p>

              <div className="mt-16 flex items-center justify-between border-t border-border/30 pt-8">
                <Link href="/docs/backend/outbox">
                  <Button variant="outline" size="sm" className="gap-2">
                    <ArrowLeft className="h-4 w-4" />
                    Transactional Outbox
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
