import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'

export default function Lesson() {
  return (
    <>
      <p>
        k6 prints about twenty numbers at the end of a run. Four of them
        answer questions you have, and the rest are context. This lesson is
        the four, and the one number everybody quotes that should be ignored.
      </p>

      <h2>The summary, annotated</h2>
      <CodeBlock
        language="text"
        code={`checks_succeeded...: 100.00%  60 out of 60     <- did it do the right thing
http_req_failed....: 0.00%    0 out of 62       <- did it answer at all
http_req_duration..: avg=1.66ms  p(95)=3.42ms   <- how long it took
                     p(99)=69.28ms
iteration_duration.: avg=2.05s                  <- includes your sleep()
http_reqs..........: 62       1.98/s            <- throughput`}
      />

      <p>
        <code>iteration_duration</code> catches people out: it includes the
        think time in the journey, so it is about 2 seconds whatever your API
        does. It measures the script, not the system. The number that
        measures the system is <code>http_req_duration</code>.
      </p>

      <h2>Ignore the average</h2>
      <p>
        The average is the one number almost every report leads with, and it
        is the least informative thing in the summary. An average of 50ms is
        consistent with every request taking 50ms, and equally consistent
        with 95% taking 10ms and 5% taking a second.
      </p>

      <Diagram
        label="Two systems, the same average"
        caption="Latency distributions are not symmetric. Averaging one is throwing the shape away."
      >
{`System A                        System B
avg 50ms, p95 55ms              avg 50ms, p95 800ms

 ####                            ########
 ####                            ########
 ####                            ########
 ####                            ########_______
 0   50   100  ...  1000         0   50   100  ...  1000

Everybody gets 50ms.            One user in twenty waits
                                most of a second.`}
      </Diagram>

      <h2>What each percentile is for</h2>
      <ul>
        <li>
          <strong>p50, the median.</strong> The typical experience. Useful
          for &quot;did my change make the common case faster&quot;, and
          useless for finding problems, because problems are not typical.
        </li>
        <li>
          <strong>p95.</strong> The usual SLO number. One request in twenty
          is worse than this. On a page that makes ten API calls, p95 per
          request means roughly <em>half</em> of page loads contain at least
          one slow call, which is why p95 matters more than it sounds.
        </li>
        <li>
          <strong>p99.</strong> Where the pathologies live: GC pauses, cache
          misses, lock contention, a connection being opened. It is also the
          noisiest, and it needs volume. At 60 requests, p99 is one request
          and tells you nothing. At 60,000 it is the most interesting line in
          the summary.
        </li>
      </ul>

      <TipBox tone="warning">
        <p>
          Never average percentiles, and never compare a p95 computed over a
          whole run against one computed over a steady state. Both are
          arithmetic that looks reasonable and produces a number that
          describes nothing. Report the percentile from the hold stage, and
          say which stage it came from.
        </p>
      </TipBox>

      <h2>Failed is not the same as wrong</h2>
      <p>
        <code>http_req_failed</code> counts responses that were not 2xx or
        3xx. <code>checks</code> counts your own assertions. They answer
        different questions, and a gap between them is informative:
      </p>

      <CodeBlock
        language="text"
        code={`http_req_failed 0%, checks 100%    the system is healthy
http_req_failed 8%, checks 92%     it is returning errors
http_req_failed 0%, checks 60%     it is answering 200 with the wrong body`}
      />

      <p>
        The third line is the one worth having checks for. An endpoint that
        returns 200 with an empty list because a query timed out and the
        error was swallowed is invisible to <code>http_req_failed</code> and
        completely visible to a check that asserts the shape of the body.
      </p>

      <TipBox tone="danger">
        <p>
          Which is why the check has to be specific. The suite Grit shipped
          until v3.432.0 checked{' '}
          <code>r.body &amp;&amp; r.body.length &gt; 0</code>, which an error
          envelope passes. It reported a successful check on a 500 for
          months. The replacement asserts{' '}
          <code>Array.isArray(r.json(&apos;data&apos;))</code>.
        </p>
      </TipBox>

      <h2>Throughput, and why it is a result and not a target</h2>
      <p>
        <code>http_reqs</code> per second is a function of your virtual user
        count and your think time, not of how fast the system is. A hundred
        VUs with a 1 to 3 second sleep produce about 50 requests a second
        whatever your API does, until the API becomes slow enough to be the
        limiting factor. Throughput only becomes a measurement of the system
        at the point where latency starts rising, and that point is what the
        breakpoint test finds.
      </p>

      <Exercise
        prompt={
          <>
            <p>
              <strong>One.</strong> Run the average load test and record avg,
              p50, p95 and p99 for <code>http_req_duration</code>.
            </p>
            <p>
              <strong>Two.</strong> Work out the ratio p99 divided by p50.
              Write it down.
            </p>
            <p>
              <strong>Three.</strong> Add a deliberate stall to one handler,
              <code>time.Sleep(500 * time.Millisecond)</code> behind a 1 in
              20 condition. Re-run and recompute the ratio.
            </p>
          </>
        }
        hint={
          <p>
            Watch what the average does in step three compared with what p95
            does. One of them notices.
          </p>
        }
        solution={
          <>
            <p>
              A p99 to p50 ratio under about 3 is a system doing the same
              amount of work for every request. Above 10 there is a second
              population in there: a cache that misses sometimes, a code path
              taken occasionally, a lock.
            </p>
            <p>
              In step three the average moves by roughly 25ms, which looks
              like noise and would survive any review. p95 moves to over
              500ms, because one request in twenty is exactly what p95
              measures. That is the entire argument for percentiles in one
              experiment.
            </p>
          </>
        }
      />

      <KnowledgeCheck
        question="A run reports avg=45ms, p50=12ms, p95=60ms, p99=2.4s. What does the shape tell you?"
        choices={[
          {
            label: 'The system is healthy: the median is 12ms',
            feedback:
              'The median says the common case is fast. It says nothing about the 1% waiting 2.4 seconds, who are real users.',
          },
          {
            label: 'Most requests are fast, and a small group takes a different and much slower path',
            correct: true,
            feedback:
              'A p99 forty times p95 is two populations, not one slow system. Cache misses, a cold connection, a lock, or a retry. The next step is finding which.',
          },
          {
            label: 'The average is the honest summary at 45ms',
            feedback:
              'The average sits between two populations and describes neither. Almost nobody experienced 45ms.',
          },
          {
            label: 'p99 is always noise and should be ignored',
            feedback:
              'It is noisy at low volume. At 2.4 seconds against a 12ms median it is a signal, and the way to confirm it is more requests, not fewer.',
          },
        ]}
      />
    </>
  )
}
