import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'

export default function Lesson() {
  return (
    <>
      <p>
        Every test so far finishes in under ten minutes, and there is a
        whole category of fault that ten minutes cannot see. A goroutine that
        leaks one per request, a slice that is appended to and never reset, a
        connection that is opened and not returned: at ten minutes they are
        invisible, and at four hours they are the outage.
      </p>

      <h2>Soak: a hundred users for four hours</h2>
      <CodeBlock
        language="javascript"
        filename="tests/k6/soak.js"
        code={`export const options = {
  stages: [
    { duration: '5m', target: 100 },  // ramp up
    { duration: '4h', target: 100 },  // long hold, the soak
    { duration: '5m', target: 0 },    // ramp down
  ],
  thresholds: defaultThresholds,
}`}
      />

      <p>
        The load is deliberately ordinary. A soak test is not about
        pressure, it is about duration: the same hundred users doing the
        same thing for long enough that anything accumulating has time to
        accumulate. Run it overnight, or on a CI runner, before a release
        that matters.
      </p>

      <h2>The signature</h2>
      <Diagram
        label="What a leak looks like"
        caption="Not a spike, not a cliff. A line that should be flat and is not."
      >
{`p95 over four hours

  |                                        ___/
  |                               ____/----
  |                      ____/----
  |             ____/----
  |____/--------
  +--------------------------------------------- time
   0h          1h          2h          3h      4h

A flat line that slowly tilts upward is the classic
memory leak. Nothing is wrong at any single moment,
which is why no short test finds it.`}
      </Diagram>

      <p>
        The tilt is the whole finding. Any single ten minute window inside
        that graph looks healthy, and that is exactly why the leak reached
        production: it passed every test anyone was willing to wait for.
      </p>

      <TipBox tone="info">
        <p>
          Open <code>/pulse/ui/</code> in another tab for the duration. k6
          tells you latency drifted. Pulse tells you heap in use climbed with
          it, or that goroutine count never came back down, which is the
          difference between knowing something is wrong and knowing what.
        </p>
      </TipBox>

      <h2>What leaks, in practice</h2>
      <ul>
        <li>
          <strong>Goroutines.</strong> The common shape is a goroutine
          waiting on a channel nobody will write to again, one per request.
          It costs a few kilobytes each, which is nothing until it is four
          hours of traffic.
        </li>
        <li>
          <strong>Database connections.</strong> Rows that are never closed,
          or a transaction that neither commits nor rolls back on an error
          path. The pool fills, and the symptom is latency, because every
          request is now waiting for a connection.
        </li>
        <li>
          <strong>Caches with no bound.</strong> A map keyed by something
          unbounded, such as a user id or a URL, with no eviction. It is not
          a leak in the strict sense; it behaves exactly like one.
        </li>
        <li>
          <strong>Nothing at all.</strong> Latency that drifts upward with a
          flat heap is usually the database: a table growing through the run,
          an index no longer fitting in cache, autovacuum falling behind.
        </li>
      </ul>

      <h2>The sixth test: breakpoint</h2>
      <p>
        The suite has one more script, and it answers a question the other
        five do not: not &quot;does it hold at 400&quot; but &quot;at what
        number does it stop holding&quot;.
      </p>

      <CodeBlock
        language="javascript"
        filename="tests/k6/breakpoint.js"
        code={`export const options = {
  // Slow climb to 5000 VUs over 1h. Adjust the ceiling for your gear.
  stages: [
    { duration: '1h', target: 5000 },
  ],
  thresholds: {
    // Abort the test the moment thresholds breach. The VU count when
    // that happens is the breaking point.
    http_req_duration: [{ threshold: 'p(95)<1000', abortOnFail: true }],
    http_req_failed:   [{ threshold: 'rate<0.05',  abortOnFail: true }],
  },
}`}
      />

      <p>
        <code>abortOnFail</code> is what makes this different from a stress
        test. The run stops at the first breach, and the virtual user count
        at that moment is your capacity number: the one to put in a launch
        plan, divide by your expected concurrency, and use to decide how many
        replicas to pay for.
      </p>

      <Exercise
        prompt={
          <>
            <p>
              <strong>One.</strong> Start the soak test with{' '}
              <code>-e RESOURCE=</code> set, and note p95 at the fifteen
              minute mark.
            </p>
            <p>
              <strong>Two.</strong> Leave it running and check p95 every
              thirty minutes. Four hours is the script&apos;s default; one
              hour is enough to see a tilt if there is one.
            </p>
            <p>
              <strong>Three.</strong> Separately, run the breakpoint test and
              record the VU count at the abort. Divide by your expected peak
              concurrency and write down the headroom.
            </p>
          </>
        }
        hint={
          <p>
            For step two, compare each reading against the fifteen minute
            one rather than against the previous reading. A drift of a few
            percent per half hour is invisible step by step and obvious end
            to end.
          </p>
        }
        solution={
          <>
            <p>
              A flat line over an hour is a real result: write it down, since
              the value of a soak test is mostly in the runs that find
              nothing. A line that has climbed 20% by the end is a leak, and
              Pulse&apos;s heap and goroutine graphs over the same window
              narrow it to one of the four causes above in a minute.
            </p>
            <p>
              The breakpoint number is the one to be careful with. It is
              capacity for <em>this</em> journey on <em>this</em> hardware
              with <em>this</em> data volume. It is not a property of your
              application, and it moves when any of those three do.
            </p>
          </>
        }
      />

      <KnowledgeCheck
        question="A four hour soak shows p95 climbing steadily from 45ms to 70ms, with the error rate flat at zero. The heap graph is flat. Where do you look?"
        choices={[
          {
            label: 'The database. A flat heap rules out the application leaking, and the data grew through the run',
            correct: true,
            feedback:
              'Four hours of writes is a bigger table, a colder index and more work for autovacuum. The application is doing the same thing to more data.',
          },
          {
            label: 'Nowhere: 70ms is well inside the threshold, so the test passed',
            feedback:
              'It passed, and it still found something. A 55% drift in four hours is a trend, and the question is where it stops.',
          },
          {
            label: 'Goroutine leak',
            feedback:
              'Goroutines show up in the heap graph, since each one holds a stack. A flat heap makes this unlikely, and the goroutine count would confirm it either way.',
          },
          {
            label: 'The load generator degrading over four hours',
            feedback:
              'Worth ruling out, and it is the rarer explanation. k6 holding a hundred VUs is close to idle on any modern machine.',
          },
        ]}
      />
    </>
  )
}
