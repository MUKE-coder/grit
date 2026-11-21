import Link from 'next/link'
import type { Metadata } from 'next'
import { GuidePage, Section, H3, Outcome, Note, Warn, Checklist, Faq, Related } from '@/components/guide'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Sagas in Go: compensating transactions without a workflow engine',
  description:
    'Do and undo steps, persisted state so a crash resumes, idempotency on both directions, and an honest account of when to reach for Temporal instead.',
  alternates: { canonical: 'https://gritframework.dev/guides/go-saga-pattern' },
  openGraph: {
    title: 'Sagas in Go: compensating transactions without a workflow engine',
    description: 'Multi-step work that can be undone, in about two hundred lines.',
    url: 'https://gritframework.dev/guides/go-saga-pattern',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="go-saga-pattern"
      lead={
        <>
          <p className="mb-4">
            A database transaction gives you all-or-nothing across tables in one database. The
            moment a business operation also charges a card, sends an email and calls a shipping
            provider, there is no transaction that covers it. Step three failing leaves steps one
            and two done.
          </p>
          <p>
            A saga is the answer, and it is simpler than its reputation: each step has an undo,
            the progress is written down, and a failure runs the undos backwards. You can build
            one in an afternoon, and most teams should, before considering a workflow engine.
          </p>
        </>
      }
    >
      <Outcome>
        <p>A saga runner with persisted state, so a crash resumes rather than restarting.</p>
        <p>Compensations that are safe to run twice, which they will be.</p>
        <p>A clear threshold for when this stops being enough.</p>
      </Outcome>

      <Section id="problem" title="The operation that cannot be a transaction">
        <CodeBlock
          code={`func (s *OrderService) Place(ctx context.Context, req PlaceOrder) error {
	order, err := s.orders.Create(ctx, req)          // our database
	if err != nil {
		return err
	}
	if err := s.inventory.Reserve(ctx, order.Items); err != nil {
		return err                                   // order exists, nothing reserved
	}
	charge, err := s.payments.Charge(ctx, order.Total)
	if err != nil {
		return err                                   // reserved, never released
	}
	if err := s.shipping.Book(ctx, order); err != nil {
		return err                                   // CHARGED, and no shipment
	}
	return nil
}`}
        />
        <p>
          Each early return leaves the world in a different half-finished state, and the last one
          leaves a customer charged for something that will never arrive. Wrapping it in{' '}
          <code>db.Transaction</code> does nothing: a rollback cannot un-charge a card, because
          the card was charged by a company that has never heard of your transaction.
        </p>
      </Section>

      <Section id="shape" title="The shape">
        <CodeBlock
          language="text"
          code={`Forward:   reserve stock  →  charge card  →  book shipment  →  send email
Backward:  release stock ←  refund card  ←  cancel shipment ←

If step 3 fails, run the compensations for steps 2 and 1, in reverse
order, and leave the saga in a terminal failed state.`}
        />
        <p>Three properties follow from that, and all three are load-bearing.</p>
        <Checklist
          items={[
            <>
              <strong>A saga is not atomic.</strong> There are moments where stock is reserved
              and the card is not yet charged. Anything observing the system must tolerate those
              intermediate states, which usually means the order has a status rather than simply
              existing or not.
            </>,
            <>
              <strong>Compensation is not rollback.</strong> A refund is a second transaction,
              visible on the customer&apos;s statement. Some things cannot be compensated at all:
              an email that has been sent is sent, so steps that cannot be undone go last.
            </>,
            <>
              <strong>Every step and every compensation must be idempotent.</strong> The runner
              retries, and a crash between doing the work and recording that it was done is
              normal rather than exceptional.
            </>,
          ]}
        />
      </Section>

      <Section id="runner" title="A runner, in about a hundred lines">
        <CodeBlock
          filename="internal/saga/saga.go"
          code={`package saga

type Step struct {
	Name string
	// Do performs the step. The returned map is merged into the saga's
	// state, so a later step can read what an earlier one produced,
	// such as a payment id needed by the refund.
	Do func(ctx context.Context, state State) (State, error)
	// Undo reverses it. It receives the state as it was AFTER Do, so
	// the refund can find the charge id. It must tolerate being called
	// when Do only partly succeeded.
	Undo func(ctx context.Context, state State) error
}

type Saga struct {
	Name  string
	Steps []Step
}`}
        />
        <CodeBlock
          filename="internal/saga/run.go"
          code={`func (r *Runner) Run(ctx context.Context, s Saga, id string, initial State) error {
	// Load or create. This is what makes a crash resumable: the saga's
	// progress is a row, not a position in a goroutine's stack.
	inst, err := r.load(ctx, id, s.Name, initial)
	if err != nil {
		return fmt.Errorf("load saga %s: %w", id, err)
	}
	if inst.Status == StatusDone || inst.Status == StatusFailed {
		return nil   // already finished; a retry is a no-op
	}
	if inst.Status == StatusCompensating {
		return r.compensate(ctx, s, inst)
	}

	for i := inst.Position; i < len(s.Steps); i++ {
		step := s.Steps[i]

		out, err := step.Do(ctx, inst.State)
		if err != nil {
			// Record the cause, then unwind everything before this step.
			inst.Status, inst.LastError = StatusCompensating, err.Error()
			if err := r.save(ctx, inst); err != nil {
				return fmt.Errorf("persist failure: %w", err)
			}
			return r.compensate(ctx, s, inst)
		}

		// Advance and persist BEFORE the next step. A crash here means
		// the step runs again, which is why Do must be idempotent, and
		// that is a much better failure than losing the record of it.
		inst.State = inst.State.Merge(out)
		inst.Position = i + 1
		if err := r.save(ctx, inst); err != nil {
			return fmt.Errorf("persist progress at %s: %w", step.Name, err)
		}
	}

	inst.Status = StatusDone
	return r.save(ctx, inst)
}

func (r *Runner) compensate(ctx context.Context, s Saga, inst *Instance) error {
	// Backwards from the last step that completed.
	for i := inst.Position - 1; i >= 0; i-- {
		step := s.Steps[i]
		if step.Undo == nil {
			continue      // nothing to reverse, e.g. a read
		}
		if err := step.Undo(ctx, inst.State); err != nil {
			// A failed compensation is the serious case: the system is
			// now inconsistent and no amount of retrying inside this
			// process will fix it. Park it and tell a human.
			inst.Status, inst.LastError = StatusStuck,
				fmt.Sprintf("undo %s: %v", step.Name, err)
			_ = r.save(ctx, inst)
			return fmt.Errorf("compensation failed at %s: %w", step.Name, err)
		}
		inst.Position = i
		if err := r.save(ctx, inst); err != nil {
			return fmt.Errorf("persist compensation: %w", err)
		}
	}
	inst.Status = StatusFailed
	return r.save(ctx, inst)
}`}
        />
        <Warn title="StatusStuck is the state that matters">
          <p>
            A saga that failed and compensated cleanly is a normal outcome. A saga whose
            compensation failed is an inconsistency that will not resolve itself: money taken and
            not refunded, stock reserved and not released. That state needs an alert and a
            screen, not a retry counter. If you build only one piece of operational tooling for
            sagas, build this one.
          </p>
        </Warn>
      </Section>

      <Section id="state" title="The table">
        <CodeBlock
          language="sql"
          code={`CREATE TABLE saga_instances (
  id          text PRIMARY KEY,         -- deterministic: "order-8821-place"
  name        text        NOT NULL,
  status      text        NOT NULL,     -- running|compensating|done|failed|stuck
  position    int         NOT NULL DEFAULT 0,
  state       jsonb       NOT NULL,     -- what the steps have produced
  last_error  text,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX ON saga_instances (status) WHERE status IN ('running', 'compensating', 'stuck');`}
        />
        <p>
          The id should be derived from the thing being done rather than generated, so that
          starting the same saga twice is the same saga. <code>order-8821-place</code> can only
          exist once, which means a duplicated request resumes rather than running a parallel
          copy.
        </p>
      </Section>

      <Section id="steps" title="Writing a step">
        <CodeBlock
          code={`var PlaceOrder = saga.Saga{
	Name: "order.place",
	Steps: []saga.Step{
		{
			Name: "reserve-stock",
			Do: func(ctx context.Context, st saga.State) (saga.State, error) {
				// Idempotent by key: calling it twice reserves once.
				res, err := inventory.Reserve(ctx, inventory.Request{
					IdempotencyKey: st.String("order_id") + ":reserve",
					Items:          st.Items("items"),
				})
				if err != nil {
					return nil, fmt.Errorf("reserve stock: %w", err)
				}
				return saga.State{"reservation_id": res.ID}, nil
			},
			Undo: func(ctx context.Context, st saga.State) error {
				id := st.String("reservation_id")
				if id == "" {
					// Do never got far enough to produce one. Nothing to
					// release, and that is a success rather than an error.
					return nil
				}
				err := inventory.Release(ctx, id)
				if errors.Is(err, inventory.ErrAlreadyReleased) {
					return nil      // compensations run twice; accept it
				}
				return err
			},
		},
		{
			Name: "charge-card",
			Do: func(ctx context.Context, st saga.State) (saga.State, error) {
				ch, err := payments.Charge(ctx, payments.Request{
					// The provider's own idempotency key, derived from the
					// order rather than generated per attempt. A fresh key
					// per attempt is how a retry becomes a second charge.
					IdempotencyKey: st.String("order_id") + ":charge",
					Amount:         st.Int64("total"),
				})
				if err != nil {
					return nil, fmt.Errorf("charge: %w", err)
				}
				return saga.State{"charge_id": ch.ID}, nil
			},
			Undo: func(ctx context.Context, st saga.State) error {
				if id := st.String("charge_id"); id != "" {
					return payments.Refund(ctx, id, st.String("order_id")+":refund")
				}
				return nil
			},
		},
		{
			Name: "send-confirmation",
			Do: func(ctx context.Context, st saga.State) (saga.State, error) {
				return nil, mail.SendOrderConfirmation(ctx, st.Uint("order_id"))
			},
			// No Undo: an email cannot be unsent. Steps that cannot be
			// compensated go LAST, so nothing after them can fail and
			// leave this one stranded.
			Undo: nil,
		},
	},
}`}
        />
        <Note title="The ordering rule">
          <p>
            Put the reversible steps first and the irreversible ones last. If the email is step
            two and the shipment booking is step three, a shipping failure leaves a customer
            holding a confirmation for an order you are about to cancel. Sorting by reversibility
            removes a whole class of awkward conversation.
          </p>
        </Note>
      </Section>

      <Section id="retries" title="Retries, and the two kinds of failure">
        <p>
          Not every failure should unwind the saga. A payment provider returning 503 is probably
          worth retrying in thirty seconds; a card declined is not, and retrying it looks like
          fraud.
        </p>
        <CodeBlock
          code={`// Distinguish them explicitly rather than retrying everything.
type Retryable struct{ error }

func (r Retryable) Unwrap() error { return r.error }

// In a step:
if isTimeout(err) || isServerError(err) {
	return nil, Retryable{fmt.Errorf("charge: %w", err)}
}
return nil, fmt.Errorf("charge declined: %w", err)   // terminal, compensate

// In the runner:
var retry Retryable
if errors.As(err, &retry) && inst.Attempts < maxAttempts {
	inst.Attempts++
	inst.NextAttemptAt = time.Now().Add(backoff(inst.Attempts))
	return r.save(ctx, inst)    // a scheduler picks it up again
}`}
        />
        <p>
          Drive the saga from your job queue rather than from the request. A saga that lives in a
          goroutine started by an HTTP handler dies with the deploy, which defeats the point of
          having written the state down.
        </p>
      </Section>

      <Section id="when" title="When to use something else">
        <p>
          A hand-rolled saga is right for a handful of workflows with a handful of steps each. It
          stops being right at recognisable points.
        </p>
        <H3>Reach for Temporal or Cadence when</H3>
        <Checklist
          items={[
            <>
              <strong>A workflow waits for days.</strong> &quot;Charge on day fourteen unless
              cancelled&quot; as a saga is a scheduler, a timer table and a lot of edge cases. It
              is a durable timer in a workflow engine, and that is a genuine difference in kind.
            </>,
            <>
              <strong>There are many long-running instances and you need visibility into each.</strong>{' '}
              The query, inspect and replay tooling is most of what a workflow engine sells.
            </>,
            <>
              <strong>Humans are in the loop.</strong> Approvals, escalations and timeouts
              between steps are what these systems are designed for.
            </>,
            <>
              <strong>Versioning matters.</strong> Changing a saga while instances are mid-flight
              is genuinely hard. Workflow engines have an answer for it. Your runner has
              &quot;hope none are running&quot;.
            </>,
          ]}
        />
        <H3>Stay with the hand-rolled version when</H3>
        <Checklist
          items={[
            <>The workflow completes in seconds and has fewer than about six steps.</>,
            <>You have a handful of workflows rather than dozens.</>,
            <>
              The operational cost of another distributed system, with its own cluster,
              deployment and failure modes, exceeds the cost of the code you would avoid writing.
            </>,
          ]}
        />
        <p>
          The version above is about two hundred lines and uses infrastructure you already
          operate. That is not a reason to avoid a workflow engine forever; it is a reason not to
          adopt one for a three-step order placement.
        </p>
      </Section>

      <Section id="testing" title="Testing it">
        <p>
          The forward path is the easy half and the one that gets tested. The tests that matter
          are the failures.
        </p>
        <CodeBlock
          code={`func TestPaymentFailureReleasesStock(t *testing.T) {
	inv := &fakeInventory{}
	pay := &fakePayments{failOn: 1}      // fail the first charge

	err := runner.Run(ctx, PlaceOrder, "order-1-place", initial)

	require.Error(t, err)
	require.True(t, inv.Released, "stock must be released when payment fails")
	require.Equal(t, saga.StatusFailed, loadStatus(t, "order-1-place"))
}

func TestCrashBetweenStepsResumes(t *testing.T) {
	// Run with a runner that stops after step 2, then run again with a
	// complete one and assert step 2 did NOT happen twice.
	_ = partialRunner(2).Run(ctx, PlaceOrder, "order-2-place", initial)
	require.NoError(t, runner.Run(ctx, PlaceOrder, "order-2-place", initial))
	require.Equal(t, 1, pay.Charges, "the charge must not be repeated")
}

func TestCompensationIsIdempotent(t *testing.T) {
	// Run every Undo twice and assert the end state is the same. This is
	// the test that catches a refund issued twice.
}`}
        />
      </Section>

      <Section id="with-grit" title="Generated">
        <CodeBlock
          terminal
          code={`grit generate workflow order.place --steps "reserve-stock,charge-card,book-shipment,send-confirmation"`}
        />
        <p>
          That writes the saga definition with a Do and an Undo stub per step, the instance table
          and its migration, the runner wired to the job queue with retry classification, and an
          admin screen listing running, failed and stuck instances with their state, their last
          error and a resume button.
        </p>
        <p>
          The stuck list is the screen worth having, for the reason above: a compensation that
          failed is money or stock in a state nothing will resolve on its own.
        </p>
        <p>
          <Link href="/docs/backend/sagas" className="text-accent hover:underline">
            The sagas reference
          </Link>{' '}
          covers the generated runner and the admin screen.
        </p>
      </Section>

      <Faq
        slug="go-saga-pattern"
        items={[
          {
            q: 'Is a saga just a transaction across services?',
            a: (
              <>
                No, and the difference is the point. A transaction is invisible until it commits;
                a saga&apos;s intermediate states are visible to everything. Stock really is
                reserved before the card is charged, and anything reading that data sees it. The
                design work is deciding what those intermediate states mean, usually by giving
                the entity a status.
              </>
            ),
            text:
              'No. A transaction is invisible until it commits, while a saga exposes its intermediate states: stock really is reserved before the card is charged, and anything reading that data sees it. The design work is deciding what those intermediate states mean, usually with a status on the entity.',
          },
          {
            q: 'What if a compensation fails?',
            a: (
              <>
                The system is inconsistent and no amount of retrying inside the process will
                change that. Retry a few times with backoff, then park the instance in a terminal
                stuck state, alert, and resolve it by hand. Designing as though compensations
                always succeed is how a refund quietly never happens.
              </>
            ),
            text:
              'The system is inconsistent and retrying will not fix it. Retry a few times with backoff, then park the instance in a terminal stuck state, alert on it and resolve by hand. Assuming compensations always succeed is how a refund quietly never happens.',
          },
          {
            q: 'Orchestration or choreography?',
            a: (
              <>
                Orchestration, which is what this page describes: one place that knows the whole
                sequence. Choreography, where each service reacts to events and emits its own,
                has less coupling and no single place that describes what happens, which means
                nobody can answer &quot;what does placing an order do&quot; without reading six
                services. Choose choreography when teams must deploy independently, and accept
                the cost in comprehensibility.
              </>
            ),
            text:
              'Orchestration, meaning one place that knows the sequence. Choreography, where services react to events, has less coupling and no single description of the flow, so nobody can say what placing an order does without reading several services. Choose it when teams must deploy independently, and accept the cost in comprehensibility.',
          },
          {
            q: 'Do I need a saga if everything is in one database?',
            a: (
              <>
                No. If every step is a write to your own Postgres, use a transaction: it is
                stronger, simpler and atomic. A saga is specifically for the case where at least
                one step is outside your transactional boundary, which in practice means a third
                party or another service.
              </>
            ),
            text:
              'No. If every step writes to your own database, use a transaction: it is stronger, simpler and atomic. A saga is for when at least one step is outside your transactional boundary, which in practice means a third party or another service.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/go-transactional-outbox',
            label: 'The transactional outbox',
            blurb: 'Making sure the event about a step actually leaves the building.',
          },
          {
            href: '/guides/go-background-jobs-asynq',
            label: 'Background jobs in Go',
            blurb: 'What drives a saga between steps, and why not a goroutine.',
          },
          {
            href: '/docs/backend/sagas',
            label: 'Sagas reference',
            blurb: 'The generated runner, the instance table and the admin screen.',
          },
          {
            href: '/guides/go-webhooks-hmac-idempotency',
            label: 'Idempotency on the receiving end',
            blurb: 'The same discipline, applied to events arriving from outside.',
          },
        ]}
      />
    </GuidePage>
  )
}
