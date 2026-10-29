import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Finding the bottleneck is the easy half. The hard half is the
        discipline after it: one change, re-measured the same way, kept or
        reverted on the number rather than on how clever it felt.
      </p>

      <h2>The loop</h2>
      <CodeBlock
        language="text"
        code={`1. Measure            run the test, record the number, say which stage
2. Form a hypothesis  "p95 is pool wait, because DB CPU is high and ours is not"
3. Change ONE thing   the smallest edit that tests the hypothesis
4. Measure again      same test, same data, same machine, same stage
5. Keep or revert     on the number, not on the reasoning`}
      />

      <p>
        Step three is where most performance work goes wrong. Two changes at
        once gives you one number and no attribution, and the usual outcome
        is that both get kept: the one that helped, and the one that made
        things slightly worse and is now load bearing for no reason.
      </p>

      <h2>Same test, same data, same machine</h2>
      <p>
        A comparison is only a comparison if everything else held still.
        The things that move without anybody deciding to move them:
      </p>

      <ul>
        <li>
          <strong>The data.</strong> A load test writes rows. The second run
          starts with a bigger table than the first, so a &quot;15%
          regression&quot; can be nothing but the table growing. Reset the
          database between runs, or use read only journeys for comparisons.
        </li>
        <li>
          <strong>The machine.</strong> A laptop that was idle for run one
          and compiling for run two produces a number about the compiler.
          Close the other work, and watch the load generator&apos;s own CPU.
        </li>
        <li>
          <strong>The stage.</strong> Comparing a p95 over the whole run
          against a p95 over the hold is comparing two different questions.
          Pick one and write it down.
        </li>
        <li>
          <strong>Warmth.</strong> First run after a restart has cold caches
          and an empty pool. Discard it, or always discard it, but be
          consistent.
        </li>
      </ul>

      <TipBox tone="warning">
        <p>
          Run the before measurement twice before you change anything. The
          gap between those two runs is your noise floor, and no improvement
          smaller than that gap is an improvement. People skip this and then
          celebrate 4% on a benchmark that varies by 9%.
        </p>
      </TipBox>

      <h2>k6 says what moved. Pulse says why</h2>
      <p>
        k6 is outside the process. It can tell you p95 went from 40ms to
        310ms, and it has no way to tell you which function those 310ms were
        in. Pulse is inside the process, and the two are meant to be read
        together.
      </p>

      <CodeBlock
        language="bash"
        filename="perf/README.md"
        code={`# The run registers itself in Pulse, so the slow p95 arrives with the
# profile of what was burning during it.
set -a; source .env; set +a    # PULSE_USERNAME / PULSE_PASSWORD
k6 run perf/load.js`}
      />

      <p>
        That is <code>grit generate perf</code>, which writes a k6 script
        that calls Pulse at the start and end of the run. Three things come
        out of it that k6 cannot measure:
      </p>

      <ul>
        <li>
          Traffic per route before, during and after, from the API&apos;s own
          rollups rather than from k6&apos;s view of its own requests.
        </li>
        <li>
          Queries per request, and N+1 findings, in each window. An N+1 that
          is invisible at one request a second is exactly what falls over at
          500 VUs.
        </li>
        <li>
          A CPU profile captured over the run, so you get the flame graph of
          the slow window rather than of an idle process afterwards.
        </li>
      </ul>

      <h2>Write the result down, including the failures</h2>
      <CodeBlock
        language="text"
        filename="notes.md"
        code={`2026-10-29  average-load, 100 VUs, hold stage, RESOURCE=products
  before (x2):  p95 311ms, 298ms      noise floor ~13ms
  change:       Preload("Category") on the list query
  after:        p95 104ms             queries/request 21 -> 2
  kept

2026-10-29  same test
  change:       DB_MAX_OPEN_CONNS 25 -> 50
  after:        p95 101ms             within noise
  reverted, no effect at this concurrency`}
      />

      <p>
        The second entry is the one worth keeping. A change that did nothing
        is a fact about your system, and without the note somebody tries it
        again in three months, including you.
      </p>

      <Exercise
        prompt={
          <>
            <p>
              <strong>One.</strong> Run the average load test twice without
              changing anything. The difference between the two p95s is your
              noise floor. Write it down.
            </p>
            <p>
              <strong>Two.</strong> Make one change you believe will help,
              from last lesson&apos;s shortlist. Re-run. Record whether the
              change is larger than the noise floor.
            </p>
            <p>
              <strong>Three.</strong> Revert it and run again. Confirm the
              number goes back.
            </p>
          </>
        }
        hint={
          <p>
            Step three is the one everybody skips and the only one that
            proves causation. If the number does not return, something else
            changed between your runs and step two proved nothing.
          </p>
        }
        solution={
          <>
            <p>
              Noise floors on a laptop are commonly 5 to 15% and on a quiet
              CI runner 2 to 5%. Knowing yours converts &quot;it feels
              faster&quot; into a decision rule you can apply in ten seconds.
            </p>
            <p>
              If step three does not restore the original number, the honest
              conclusion is that the experiment was invalid, not that the
              revert was also an improvement. Reset the data, close the other
              work and start again.
            </p>
          </>
        }
      />

      <KnowledgeCheck
        question="You change three things at once and p95 improves by 40%. What do you know?"
        choices={[
          {
            label: 'That the combination helped, and nothing about which part',
            correct: true,
            feedback:
              'Exactly. One of them might be doing all the work and another making it worse. Unpick them one at a time, or keep a change you cannot justify.',
          },
          {
            label: 'That all three helped',
            feedback:
              'A 40% net improvement is consistent with one change giving 50% and another costing 10%.',
          },
          {
            label: 'That the biggest change is responsible',
            feedback:
              'Size of diff and size of effect are unrelated. A one character index hint can beat a thousand line rewrite.',
          },
          {
            label: 'Enough to ship it, since the number improved',
            feedback:
              'You can ship it. You cannot maintain it: nobody knows which parts are load bearing, so the next person removes the wrong one.',
          },
        ]}
      />
    </>
  )
}
