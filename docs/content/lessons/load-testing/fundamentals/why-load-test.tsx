import { TipBox } from '@/components/course/tip-box'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'

export default function Lesson() {
  return (
    <>
      <p>
        Load testing means deliberately hammering your API to find out
        what happens when traffic spikes. Most teams discover their
        bottleneck on launch day from a screaming Slack channel. This
        course is the alternative: find it on Tuesday afternoon, with
        coffee, and fix it.
      </p>

      <h2>What load testing answers</h2>
      <ul>
        <li>
          <strong>How many concurrent users can my API handle?</strong>{' '}
          50? 5000? You should be able to answer this with a number.
        </li>
        <li>
          <strong>Where does it break?</strong> DB CPU, connection
          pool, app CPU, memory, downstream service?
        </li>
        <li>
          <strong>What does the user experience look like under
          load?</strong> p95 = 200ms? 2s? 20s? The difference matters.
        </li>
        <li>
          <strong>Does my change make things WORSE?</strong> Run before
          + after; the numbers tell you.
        </li>
      </ul>

      <h2>The five test types</h2>
      <Diagram label="Test types over time" caption="Each test asks a different question. Pick based on what you need to learn.">
{`     VUs
      │
      │                              ──── Stress (find the break point) ────
   500├────────────────────────────/
      │                          /
      │                        /
   200├──────────────────────/────── Load (expected steady state) ──
      │                    /─────────────────────────────────────
      │                   /
    10│──── Smoke ────────/                                       Soak (hours, find leaks)
      │                                                          ──────────
      │   ───────► time                                          ──────────`}
      </Diagram>
      <ul>
        <li>
          <strong>Smoke</strong> — 10 VUs for 1 min. &quot;Does this
          basically work?&quot; Run on every PR.
        </li>
        <li>
          <strong>Load</strong> — expected production traffic for 5-10
          min. &quot;Does it hold up at planned scale?&quot;
        </li>
        <li>
          <strong>Stress</strong> — ramp until it breaks. &quot;Where&apos;s
          the cliff?&quot;
        </li>
        <li>
          <strong>Spike</strong> — sudden 10x surge. &quot;Does a
          marketing email crash us?&quot;
        </li>
        <li>
          <strong>Soak</strong> — moderate load for hours. &quot;Do we
          leak memory / connections over time?&quot;
        </li>
      </ul>

      <h2>When NOT to load test</h2>
      <p>
        Load testing has a cost: time to write, infrastructure to run,
        signal to interpret. Skip it if:
      </p>
      <ul>
        <li>
          <strong>You have &lt;100 daily active users.</strong> Your DB
          isn&apos;t the bottleneck — feature velocity is. Build first.
        </li>
        <li>
          <strong>You haven&apos;t shipped anything.</strong> &quot;What if
          a million users come?&quot; is a great problem to have later.
          Don&apos;t pre-optimise.
        </li>
        <li>
          <strong>The endpoint is rarely called.</strong> A
          quarterly-reporting endpoint doesn&apos;t need to handle 1000
          rps. Test the hot paths, ignore the cold ones.
        </li>
      </ul>

      <TipBox tone="info">
        <strong>The right time</strong> is one of: (a) you&apos;re
        about to launch / get featured / get press, (b) you just
        shipped a major change to a hot path, (c) you noticed slowness
        in production and want to repro locally, (d) you&apos;re
        capacity-planning for next quarter.
      </TipBox>

      <h2>What we&apos;ll build together</h2>
      <ul>
        <li>A k6 suite in <code>tests/k6/</code> with one script per scenario type.</li>
        <li>Thresholds that fail the build on regression.</li>
        <li>A workflow that runs k6 on every PR via GitHub Actions.</li>
        <li>One real measured optimisation against your own API.</li>
      </ul>

      <h2>The mindset shift</h2>
      <p>
        Most performance work is vibes-based: &quot;feels faster&quot;,
        &quot;should be quicker&quot;. Load testing is the discipline
        of replacing vibes with numbers. You make a change, you run
        the test, the number went down — or it didn&apos;t. The graph
        is brutal honest. That&apos;s the whole point.
      </p>

      <KnowledgeCheck
        question="A 2-person team is 3 weeks from launching. They've never load-tested. Should they?"
        choices={[
          {
            label: "Yes — they need to know they'll survive launch",
            correct: true,
            feedback:
              "Right — launch is the canonical right time. Even one smoke + one load test catches the cliff before customers do. 2 hours of work, often saves a 6-hour outage.",
          },
          {
            label: "No — they don't have users yet",
            feedback:
              "But they're about to. Launch is exactly the wrong moment to discover your API tips over at 100 RPS.",
          },
          {
            label: 'No — load testing is for big teams',
            feedback: 'k6 takes 5 minutes to install and 20 lines to write a smoke test. Team size irrelevant.',
          },
          {
            label: 'Only if they have Kubernetes',
            feedback: 'k6 runs anywhere. Local, Docker, CI. No infrastructure assumptions.',
          },
        ]}
      />

      <h2>What&apos;s next</h2>
      <p>
        Next lesson — <strong>Install k6</strong>. One binary, no Node,
        no deps. We&apos;ll have you running a test by the end of the
        next 3 minutes.
      </p>
    </>
  )
}
