import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Diagram, DiagramBox, DiagramRow, DiagramArrow } from '@/components/diagram'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/workflows')

const C = 'text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded'

export default function WorkflowsPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Backend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Workflows</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                A status column that anything can set to anything is not a status. It is a text
                field with suggestions. A workflow turns one select field into a process the
                server enforces, declared in the same place the field is.
              </p>
            </div>

            <div className="prose-grit">
              <h2 id="problem">What a status field is without one</h2>
              <p>
                Generate an order with a <code className={C}>status</code> select and the admin
                shows a dropdown with every value on every record. Any client that can update
                an order can send this:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock
                language="bash"
                code={`PATCH /api/v1/orders/8f2a
{ "status": "delivered" }`}
              />
            </div>

            <div className="prose-grit">
              <p>
                On an unpaid order. Nothing stops it. A delivered order can go back to pending,
                an unpaid one can be marked shipped, and the only thing between your fulfilment
                process and nonsense is that people mostly click the right option. That is the
                default behaviour of a status dropdown in every CRUD admin, Grit included,
                until you declare the process.
              </p>

              <h2 id="declare">Declaring it</h2>
              <p>
                A workflow is a block on the select field itself. That needs more structure
                than a command-line flag, so the resource goes in a YAML file:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock
                filename="order.yaml"
                code={`- name: status
  type: select
  options:
    - { value: pending,   label: Pending payment }
    - { value: paid,      label: Paid }
    - { value: packed,    label: Packed }
    - { value: shipped,   label: Shipped }
    - { value: delivered, label: Delivered }
    - { value: cancelled, label: Cancelled }
  workflow:
    initial: pending
    terminal: [delivered, cancelled]
    transitions:
      - action: mark_paid
        from: [pending]
        to: paid
      - action: pack
        from: [paid]
        to: packed
        permission: orders.fulfil
      - action: ship
        from: [packed]
        to: shipped
        permission: orders.fulfil
      - action: deliver
        from: [shipped]
        to: delivered
      - action: cancel
        from: [pending, paid, packed]
        to: cancelled
        confirm: true`}
              />
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock
                language="bash"
                code={`grit g resource Order --from order.yaml --force
grit migrate`}
              />
            </div>

            <h2 className="text-2xl font-semibold tracking-tight mb-4 mt-12">The shape of it</h2>

            <Diagram>
              <DiagramRow>
                <DiagramBox title="pending" sub="initial" tone="amber" />
                <DiagramBox title="paid" sub="mark_paid" tone="cyan" />
                <DiagramBox title="packed" sub="pack · orders.fulfil" tone="cyan" />
              </DiagramRow>
              <DiagramArrow label="ship · orders.fulfil" />
              <DiagramRow>
                <DiagramBox title="shipped" sub="deliver" tone="blue" />
                <DiagramBox title="delivered" sub="terminal" tone="green" />
                <DiagramBox
                  title="cancelled"
                  sub="terminal · from pending, paid, packed"
                  tone="rose"
                />
              </DiagramRow>
            </Diagram>

            <div className="prose-grit">
              <h2 id="decisions">Three decisions worth understanding</h2>

              <h3>States come from the field&apos;s own options</h3>
              <p>
                They are not listed twice. If they were, the dropdown and the state machine
                would drift, and you would get a transition to a state the UI never offers, or
                an option nothing can reach.
              </p>

              <h3>It is a directed graph, not a ladder</h3>
              <p>
                This is the one people get wrong. <code className={C}>paid</code> back to{' '}
                <code className={C}>pending</code> is illegal only because you did not declare
                it. Add a transition and it becomes legal:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock
                language="yaml"
                code={`- action: reopen
  from: [paid]
  to: pending`}
              />
            </div>

            <div className="prose-grit">
              <p>Only three things are actually enforced:</p>
              <ul>
                <li>
                  <strong>Nothing leaves a terminal state.</strong> Whatever you list under{' '}
                  <code className={C}>terminal</code> is an end.
                </li>
                <li>
                  <strong>
                    An empty <code className={C}>from</code> means from anywhere.
                  </strong>{' '}
                  Useful for something like <code className={C}>archive</code> that applies
                  whatever the current state is.
                </li>
                <li>
                  <strong>Generation refuses a dead end.</strong> A non-terminal state with no
                  way out is rejected when you generate, rather than discovered in production
                  when an order lands in it and nobody can move it.
                </li>
              </ul>

              <h3>Permission sits on the transition, not the resource</h3>
              <p>
                <code className={C}>pack</code> and <code className={C}>ship</code> require{' '}
                <code className={C}>orders.fulfil</code>. <code className={C}>deliver</code> does
                not. Warehouse staff advance fulfilment without being able to refund, which is
                not expressible if permission is a property of the whole resource.
              </p>

              <h2 id="generated">What it generates</h2>
              <p>Four things, none of which you write.</p>

              <h3>1. A guarded service method</h3>
              <p>
                The check runs in the service, so a warehouse account calling{' '}
                <code className={C}>ship</code> without the permission gets a 403 from the
                server. Not a hidden button: hiding a button is a UI preference, and the request
                still works if somebody sends it by hand.
              </p>

              <h3>2. Per-action endpoints, and no general status write</h3>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock
                language="bash"
                code={`POST /api/v1/orders/:id/transitions/mark_paid
POST /api/v1/orders/:id/transitions/pack
POST /api/v1/orders/:id/transitions/cancel

GET  /api/v1/orders/workflow      the definition, for a client that draws it`}
              />
            </div>

            <div className="prose-grit">
              <p>
                This is the core of the whole feature. There is no endpoint that sets{' '}
                <code className={C}>status</code> to an arbitrary value, so an illegal jump is
                not rejected. It is <strong>unrepresentable</strong>.
              </p>
              <p>
                An action that exists but is not legal from here returns a 422 saying what state
                the record is in and which actions <em>are</em> available from it, so a client
                can show the next step without a second request.
              </p>
              <p>
                <code className={C}>GET /workflow</code> returns the definition itself, which is
                what lets a client draw only the legal transitions as buttons instead of a
                dropdown of everything.
              </p>

              <h3>3. A domain event per transition</h3>
              <p>
                Each transition emits{' '}
                <code className={C}>&lt;resource&gt;.&lt;action&gt;</code>, so{' '}
                <code className={C}>orders.mark_paid</code>, on the shared bus. Audit, webhooks
                and realtime are already subscribers, and so is anything you add:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock
                language="go"
                code={`events.On("orders.mark_paid", events.Async, "send-confirmation", sendOrderConfirmation)`}
              />
            </div>

            <div className="prose-grit">
              <p>
                The transition is its own event rather than a generic update, because a
                subscriber that cares about orders being paid should not have to diff two
                versions of a record to work out that is what happened. It is also how a
                confirmation email gets sent without a line of email code in the checkout
                handler.
              </p>

              <h3>
                4. <code className={C}>confirm: true</code> in the admin
              </h3>
              <p>
                Marks the action as needing a confirmation step, because{' '}
                <code className={C}>cancel</code> is not undoable and a misclick on a table row
                is easy.
              </p>

              <h2 id="trap">The trap those events create</h2>
              <p>
                Read this before adding a back-edge. Say your{' '}
                <code className={C}>orders.mark_paid</code> subscriber sends the customer
                confirmation email, and you add a <code className={C}>reopen</code> transition
                from <code className={C}>paid</code> to <code className={C}>pending</code>{' '}
                because a payment sometimes needs redoing.
              </p>
              <p>
                Now <code className={C}>paid → pending → paid</code> sends that email twice. The
                state machine is behaving exactly as declared. The mistake is reversing into a
                state whose entry has a side effect.
              </p>
              <p>
                Use a compensating state instead. A <code className={C}>refunded</code> or{' '}
                <code className={C}>payment_failed</code> state, with its own transitions and
                its own event, says what actually happened rather than pretending the order went
                back in time.
              </p>

              <h2 id="not-yet">What is not built yet</h2>
              <p>
                The admin does not render workflow state as a badge with the legal transitions
                as buttons. Today you get the generated dropdown plus the transition endpoints,
                and a custom cell is how to draw it properly. The definition endpoint exists so
                that component can be written without hardcoding the graph.
              </p>
            </div>

            <div className="mt-16 flex items-center justify-between border-t border-border/50 pt-8">
              <Link href="/docs/backend/rbac">
                <Button variant="ghost" className="gap-2">
                  <ArrowLeft className="h-4 w-4" />
                  RBAC &amp; Roles
                </Button>
              </Link>
              <Link href="/docs/backend/variants">
                <Button variant="ghost" className="gap-2">
                  Product Variants
                  <ArrowRight className="h-4 w-4" />
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
