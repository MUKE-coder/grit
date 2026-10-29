import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'

export default function Lesson() {
  return (
    <>
      <p>
        Smoke and average load ask whether the system works. Stress and spike
        ask what it does when it stops working, which is a different question
        and usually a more useful one, because every system stops working
        eventually and only some of them do it gracefully.
      </p>

      <h2>Stress: four times normal, held</h2>
      <CodeBlock
        language="javascript"
        filename="tests/k6/stress.js"
        code={`export const options = {
  stages: [
    { duration: '2m', target: 100 },  // normal
    { duration: '5m', target: 400 },  // 4x normal, the stress
    { duration: '2m', target: 0 },    // ramp down
  ],
  thresholds: aggressiveThresholds,
}`}
      />

      <p>
        Four times is a convention, not a law. Pick the multiple that
        represents your worst realistic day: the Monday after a launch, the
        hour after an email goes out, Black Friday. The two minutes at normal
        load first are there so you have a before to compare against, in the
        same run and on the same machine.
      </p>

      <h2>Why the thresholds are looser</h2>
      <CodeBlock
        language="javascript"
        filename="tests/k6/lib/common.js"
        code={`export const aggressiveThresholds = {
  http_req_duration: ['p(95)<2000'],
  http_req_failed:   ['rate<0.10'],
}`}
      />

      <p>
        A stress test is supposed to hurt. With the normal thresholds it
        would breach in the first thirty seconds and the run would report
        failure, which tells you only what you already intended. The looser
        set lets the run finish so you can watch <em>how</em> it degrades.
      </p>

      <h2>The thing you are actually looking for</h2>
      <Diagram
        label="Two systems at 4x load"
        caption="Both are over capacity. Only one of them comes back."
      >
{`graceful                        collapsing

p95                             p95
 |            ____               |              /
 |       ____/                   |             /
 |______/                        |____________/
 +----------------- time         +----------------- time
 errors stay near zero,          errors climb, then the
 latency rises and holds         process restarts and the
                                 graph starts over at zero`}
      </Diagram>

      <p>
        Rising latency with a flat error rate is a queue doing its job:
        requests wait, and they are served. Rising latency <em>and</em>{' '}
        rising errors means something is giving up: a connection pool
        timeout, a context deadline, an out of memory kill. The second is
        the one that wakes you up, and it is invisible in any single number.
      </p>

      <TipBox tone="info">
        <p>
          Watch the API&apos;s own memory and the database&apos;s CPU during
          the run, not just k6&apos;s output. A system that collapses at 4x
          almost never does so because of the API process: it runs out of
          database connections, or the database runs out of CPU.
        </p>
      </TipBox>

      <h2>Spike: the surge, and the recovery</h2>
      <CodeBlock
        language="javascript"
        filename="tests/k6/spike.js"
        code={`export const options = {
  stages: [
    { duration: '1m',  target: 50 },    // baseline
    { duration: '10s', target: 1000 },  // SLAM
    { duration: '1m',  target: 1000 },  // hold the spike
    { duration: '10s', target: 50 },    // drop back
    { duration: '1m',  target: 50 },    // recovery window, watch p95
    { duration: '10s', target: 0 },
  ],
  thresholds: aggressiveThresholds,
}`}
      />

      <p>
        Most people read this script as the 10 second ramp to a thousand
        users. The important stage is the last minute, back at fifty.
      </p>

      <p>
        A system that survives a spike and returns to its baseline latency
        has absorbed it. A system whose p95 is still elevated a minute after
        the surge has not: something is still draining. A retry storm, a
        queue that filled and is being worked off, a cache that was evicted
        and is being rebuilt one request at a time. That is the finding, and
        it only exists in the recovery window.
      </p>

      <TipBox tone="warning">
        <p>
          A thousand virtual users from one laptop is usually a test of your
          laptop. If CPU on the k6 machine saturates, every number it reports
          includes its own queueing. Watch the load generator as carefully as
          the system under test, and move k6 to a bigger box before
          concluding anything about the system.
        </p>
      </TipBox>

      <Exercise
        prompt={
          <>
            <p>
              <strong>One.</strong> Run the stress test with{' '}
              <code>-e RESOURCE=</code> set to your heaviest list endpoint.
              Record p95 at the two minute mark and at the seven minute mark.
            </p>
            <p>
              <strong>Two.</strong> Record the error rate at the same two
              points. Classify the system: graceful or collapsing.
            </p>
            <p>
              <strong>Three.</strong> Run the spike test. Record p95 during
              the hold, and again in the last thirty seconds of the recovery
              window.
            </p>
          </>
        }
        hint={
          <p>
            For step three, compare the recovery p95 against the first
            minute&apos;s baseline, not against the spike. The question is
            whether it came back, not whether it was bad while the spike was
            on.
          </p>
        }
        solution={
          <>
            <p>
              Step two is the classification that matters. Latency up, errors
              flat is a system under pressure and working. Latency up, errors
              up is a system failing, and the error response tells you where:
              a 503 from your own middleware, a 500 from a pool timeout and a
              connection reset are three different problems.
            </p>
            <p>
              In step three, a recovery p95 within noise of the baseline means
              the spike cost you nothing after it ended. Anything elevated by
              more than about 20% a minute later is a system still working
              something off, and the next chapter is how to find out what.
            </p>
          </>
        }
      />

      <KnowledgeCheck
        question="During a spike test, p95 goes from 40ms to 1.8s at peak, then settles at 900ms for the full recovery minute at baseline load. What have you found?"
        choices={[
          {
            label: 'Nothing: 1.8s at ten times the load is reasonable',
            feedback:
              'The peak is not the finding. The problem is the 900ms afterwards, at a load the system handled at 40ms a minute earlier.',
          },
          {
            label: 'The system did not recover, and something is still draining after the surge',
            correct: true,
            feedback:
              'Same load, twenty times the latency. A queue being worked off, a cache rebuilding one request at a time, or clients retrying. Whatever it is, it outlived the spike.',
          },
          {
            label: 'The thresholds are too loose; tighten them and re-run',
            feedback:
              'Tightening them makes the run fail earlier, which hides the recovery window that contains the finding.',
          },
          {
            label: 'The load generator was saturated',
            feedback:
              'Worth ruling out, and it would show as elevated latency during the peak rather than after it, while k6 CPU was pinned.',
          },
        ]}
      />
    </>
  )
}
