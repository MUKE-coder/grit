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
              <span className="tag-mono text-accent/80 mb-3 block">Backend</span>
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

              <h2>Finding a stuck run</h2>
              <p>
                The admin has a screen at <code>/system/sagas</code>: every run, filtered by
                status, with its steps in order and the error that stopped each one. A stuck run
                gets a Retry button, which is the only state where a person pressing something is
                the right answer. It resumes compensating from the step it stopped on, because the
                thing that could not be undone is still not undone and whoever fixed the refund API
                wants that same undo attempted again.
              </p>
              <p>
                Without the screen the only way to find a stuck run is SQL against{' '}
                <code>saga_runs</code>, which means nobody finds one until a customer complains.
              </p>

              <h2 id="worked-example">A whole saga, end to end</h2>
              <p>
                Checkout: take the money, hold the stock, book the courier, email the receipt.
                Four systems, none of which rolls back because another one failed. The order
                matters and is explained under it.
              </p>

              <CodeBlock
                language="go"
                filename="internal/sagas/checkout.go"
                code={`package sagas

import (
    "context"

    "myapp/apps/api/internal/saga"
)

type CheckoutInput struct {
    OrderID string \`json:"order_id"\`
    Email   string \`json:"email"\`
}

func Checkout() saga.Definition {
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
                    // IdempotencyKey is stable for this run and this step, so a
                    // replay after a crash reaches the same charge rather than a
                    // second one.
                    charge, err := payments.Charge(ctx, in.OrderID,
                        payments.Idempotency(r.IdempotencyKey()))
                    if err != nil {
                        if payments.Declined(err) {
                            // Retrying a declined card spends five attempts to
                            // arrive at the same answer.
                            return saga.Fatal(err)
                        }
                        return err
                    }
                    r.Set("charge_id", charge.ID)
                    return nil
                },
                Undo: func(ctx context.Context, r *saga.Run) error {
                    id := r.GetString("charge_id")
                    if id == "" {
                        // The Do never got far enough to charge anything.
                        return nil
                    }
                    return payments.Refund(ctx, id)
                },
            },
            {
                Name: "reserve_stock",
                Do: func(ctx context.Context, r *saga.Run) error {
                    var in CheckoutInput
                    if err := r.Input(&in); err != nil {
                        return err
                    }
                    ref, err := inventory.Reserve(ctx, in.OrderID)
                    if err != nil {
                        return err
                    }
                    r.Set("reservation", ref)
                    return nil
                },
                Undo: func(ctx context.Context, r *saga.Run) error {
                    ref := r.GetString("reservation")
                    if ref == "" {
                        return nil
                    }
                    return inventory.Release(ctx, ref)
                },
            },
            {
                Name:        "book_courier",
                MaxAttempts: 8, // their API is flaky in a way ours is not
                Do: func(ctx context.Context, r *saga.Run) error {
                    var in CheckoutInput
                    if err := r.Input(&in); err != nil {
                        return err
                    }
                    booking, err := courier.Book(ctx, in.OrderID)
                    if err != nil {
                        return err
                    }
                    r.Set("booking", booking.Reference)
                    return nil
                },
                Undo: func(ctx context.Context, r *saga.Run) error {
                    ref := r.GetString("booking")
                    if ref == "" {
                        return nil
                    }
                    return courier.Cancel(ctx, ref)
                },
            },
            {
                Name: "send_receipt",
                // No Undo: an email cannot be unsent. That is why it is last.
                Do: func(ctx context.Context, r *saga.Run) error {
                    var in CheckoutInput
                    if err := r.Input(&in); err != nil {
                        return err
                    }
                    return mail.Receipt(ctx, in.Email, r.GetString("charge_id"))
                },
            },
        },
    }
}`}
              />

              <p>
                <strong>Why the receipt is last.</strong> A step with no{' '}
                <code>Undo</code> is allowed, and the ordering is the design. Put the receipt
                second and a courier failure compensates the charge and the stock, and leaves a
                customer holding an email for an order that no longer exists. The rule is: the
                things you cannot take back go after the things you can.
              </p>

              <p>
                <strong>Why every Undo checks first.</strong> Compensation exists because
                something went wrong, and &quot;went wrong&quot; includes &quot;I do not know
                whether the charge landed&quot;. Each <code>Undo</code> above reads what the{' '}
                <code>Do</code> recorded and returns cleanly when it finds nothing, because it
                will be called for steps that never finished.
              </p>

              <p>
                <strong>Why state goes through <code>r.Set</code>.</strong> By the time a
                compensation runs, the process that charged the card is usually gone. The run&apos;s
                state is a column, so <code>r.Get</code> is the only thing an <code>Undo</code>{' '}
                has. A closure variable would be a different process&apos;s memory.
              </p>

              <h2 id="mechanics">What the runner actually does</h2>
              <p>
                A background runner ticks every five seconds, claims the runs that are due, and
                advances each by one step. The claim is a row update, so with four replicas
                exactly one of them is advancing a given run: a step is never racing another copy
                of itself. The runner does not start at all when no saga is registered, rather
                than polling an empty table for the life of the process.
              </p>
              <p>
                A failing step is retried with exponential backoff, five attempts by default,
                starting at five seconds and doubling to a fifteen-minute ceiling, with jitter so
                a thousand runs waiting on one outage do not retry in the same millisecond. Those
                are defaults on the definition; a step can override{' '}
                <code>MaxAttempts</code> for itself, and <code>1</code> means do not retry.
              </p>
              <p>
                Past its attempts the run turns to compensating and the completed steps are undone
                newest first. If a compensation itself fails past its attempts the run goes{' '}
                <strong>stuck</strong>, which is the status that needs a person. A compensation is
                never skipped on the way past, because skipping one is how money goes missing.
              </p>

              <h2 id="temporal">Is this Temporal?</h2>
              <p>
                No, and the difference is worth being precise about, because the word
                &quot;saga&quot; names a <em>pattern</em> and Temporal is a <em>product</em>.
              </p>
              <p>
                Both implement the same pattern, and the same flavour of it:{' '}
                <strong>orchestrated</strong>, where one coordinator calls each step in turn and
                decides when to compensate. The alternative is choreography, where services react
                to each other&apos;s events and no single component knows the whole sequence.
                Grit&apos;s engine and Temporal are both the first kind, so the shape of the code
                is recognisably similar: steps, compensations registered as you go, compensations
                run in reverse.
              </p>
              <p>
                What differs is where the durability lives, and what you have to run.
              </p>

              <div className="overflow-x-auto my-6">
                <table className="w-full text-sm table-fixed">
                  <colgroup>
                    <col className="w-[22%]" />
                    <col className="w-[39%]" />
                    <col className="w-[39%]" />
                  </colgroup>
                  <thead>
                    <tr className="border-b border-border/40 text-left">
                      <th className="px-3 py-2 font-medium" />
                      <th className="px-3 py-2 font-medium text-accent">internal/saga</th>
                      <th className="px-3 py-2 font-medium text-accent">Temporal</th>
                    </tr>
                  </thead>
                  <tbody>
                    {[
                      ['To run it', 'Nothing. Two tables in the database you already have.', 'A Temporal cluster, plus worker processes polling a task queue.'],
                      ['Where progress lives', 'A saga_runs row and a saga_steps row per step.', 'An event history in the cluster, replayed to rebuild workflow state.'],
                      ['How code resumes', 'The next step is read from the row and called.', 'The whole workflow function is replayed from its history.'],
                      ['Constraints on your code', 'A step must be idempotent. Nothing else.', 'Workflow code must be deterministic: no time.Now, no rand, no direct I/O.'],
                      ['Waiting and timers', 'Retry backoff. No durable sleeps or signals.', 'Durable timers, signals, queries, child workflows, versioning.'],
                      ['Scale it is built for', 'One application\u2019s processes, sharing one database.', 'Many services and languages, millions of concurrent executions.'],
                    ].map(([label, grit, temporal]) => (
                      <tr key={label} className="border-b border-border/20 align-top">
                        <td className="px-3 py-2 font-medium text-foreground/80">{label}</td>
                        <td className="px-3 py-2 text-muted-foreground leading-relaxed">{grit}</td>
                        <td className="px-3 py-2 text-muted-foreground leading-relaxed">{temporal}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              <p>
                The determinism row is the one that bites people. Temporal rebuilds a
                workflow&apos;s state by replaying its code against a recorded history, so that
                code may not do anything that could answer differently the second time: no{' '}
                <code>time.Now()</code>, no <code>rand</code>, no database calls outside an
                activity. Grit replays nothing. It reads which step is next out of a row and calls
                it, so a step is ordinary Go that may do whatever it likes, and the single rule
                is that running it twice must be safe.
              </p>
              <p>
                That is the trade in both directions. If you need durable timers, signals into a
                running process, child workflows or many services in several languages sharing one
                engine, that is what Temporal is for and this is not a substitute. If what you
                have is one Go application whose checkout touches four APIs, a cluster is a large
                piece of infrastructure to operate for a problem two tables solve.
              </p>

              <h2 id="commands">Every command</h2>
              <CodeBlock
                language="bash"
                code={`# Write a saga with its steps stubbed out, and register it
grit generate workflow Checkout --steps "charge,reserve_stock,book_courier,send_receipt"

# Create the two tables (saga_runs, saga_steps)
grit migrate

# What the project has, and what is running
grit routes | grep sagas`}
              />
              <p>
                There is no command to start a run: a saga is started from your own code, because
                the thing that knows an order was placed is the handler that placed it.
              </p>
              <CodeBlock
                language="go"
                code={`run, err := saga.Start(ctx, db, "checkout", sagas.CheckoutInput{
    OrderID: order.ID,
    Email:   order.Email,
})`}
              />

              <h2 id="watching">Watching a run</h2>
              <p>
                The admin has a screen at <strong>System &rarr; Sagas</strong>, listing runs by
                status with each run&apos;s steps and the error that stopped it. It is behind{' '}
                <code>ADMIN</code> or <code>jobs.view</code>. Finding a stuck run should not mean
                writing SQL against <code>saga_runs</code>, because nobody does that until a
                customer complains.
              </p>
              <CodeBlock
                language="bash"
                code={`GET  /api/v1/admin/sagas              # runs, filterable by status
GET  /api/v1/admin/sagas/definitions  # what this project has registered
GET  /api/v1/admin/sagas/:id          # one run, with its steps
POST /api/v1/admin/sagas/:id/retry    # push a stuck run on again`}
              />
              <p>
                <code>retry</code> is for the stuck ones: you have fixed whatever made the
                compensation fail, and you want the run to carry on undoing rather than stay
                half-undone forever.
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
