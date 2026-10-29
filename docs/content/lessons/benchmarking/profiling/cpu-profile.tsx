import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        A benchmark says a function takes 14790 nanoseconds. A CPU profile
        says which of the fifty functions it calls spent them. You need both,
        and you need them in that order: profiling without a benchmark gives
        you a profile of whatever the process happened to be doing.
      </p>

      <h2>From a benchmark</h2>
      <CodeBlock
        language="bash"
        code={`cd apps/api
go test -run XXX -bench BenchmarkListEnvelope -benchtime 5s \\
  -cpuprofile cpu.out ./internal/handlers/

go tool pprof -http=: cpu.out`}
      />

      <p>
        <code>-benchtime 5s</code> matters more than it looks. A profile is a
        sample taken roughly a hundred times a second, so a benchmark that
        runs for 200 milliseconds gives you about twenty samples and a graph
        made of noise. Five seconds is five hundred, which is enough for the
        shape to be real.
      </p>

      <p>
        <code>-run XXX</code> is there to match no test. Without it you
        profile the test suite as well as the benchmark, and the test suite
        is usually the larger half.
      </p>

      <h2>Reading it: flat against cumulative</h2>
      <CodeBlock
        language="text"
        code={`      flat  flat%   sum%        cum   cum%
     1.2s  24.0%  24.0%      1.2s  24.0%  encoding/json.(*encodeState).string
     0.8s  16.0%  40.0%      3.4s  68.0%  handlers.listEnvelope
     0.1s   2.0%  42.0%      4.6s  92.0%  handlers.(*PostHandler).List`}
      />

      <ul>
        <li>
          <strong>flat</strong> is time inside that function itself.
        </li>
        <li>
          <strong>cum</strong> is time inside it and everything it called.
        </li>
      </ul>

      <p>
        The handler has 92% cumulative and 2% flat, which is exactly what you
        expect: it is a handler, it does nothing but call things. Sort by
        flat to find the work, and read cum to find which path led there. A
        function with high cum and low flat is a signpost, not a destination.
      </p>

      <TipBox tone="info">
        <p>
          In the web view, the <strong>Flame Graph</strong> is the fastest way
          in, and <strong>Source</strong> is the most convincing: it shows
          per-line timings in your own code, which usually ends an argument
          about whether a loop is the problem.
        </p>
      </TipBox>

      <h2>From the running API</h2>
      <p>
        A benchmark profiles one function. Sometimes the question is what the
        whole API is doing under load, and for that Pulse serves standard
        pprof:
      </p>

      <CodeBlock
        language="bash"
        code={`# while a load test is running
go tool pprof -http=: http://localhost:8080/pulse/api/profile/cpu`}
      />

      <p>
        This is off by default, and it takes two switches rather than one:{' '}
        <code>PULSE_PROFILING=true</code> in <code>.env</code>, and{' '}
        <code>PULSE_PROFILE_ENABLED=true</code> in the environment. A profile
        leaks the shape of your code, so one forgotten setting is not enough
        to expose it.
      </p>

      <h2>What a CPU profile cannot see</h2>
      <p>
        It samples goroutines that are <em>running</em>. A request that spends
        900 milliseconds waiting for a database, a lock or an HTTP call
        appears nowhere, because nothing was on CPU. This is the single most
        common way a profile misleads: a handler that is slow for waiting
        reasons looks cheap.
      </p>

      <CodeBlock
        language="text"
        code={`slow and CPU is high    -> CPU profile
slow and CPU is low     -> not a CPU problem. Look at pool wait,
                           the block profile, the mutex profile,
                           or whatever it calls over a network`}
      />

      <Exercise
        prompt={
          <>
            <p>
              <strong>One.</strong> Profile{' '}
              <code>BenchmarkListEnvelope</code> for five seconds and open
              it. Note the top three functions by flat time.
            </p>
            <p>
              <strong>Two.</strong> Run the same benchmark with{' '}
              <code>-benchtime 100ms</code> and profile that. Compare the top
              three.
            </p>
            <p>
              <strong>Three.</strong> Add{' '}
              <code>time.Sleep(10 * time.Millisecond)</code> inside the
              function, re-run the five second profile, and find the sleep.
            </p>
          </>
        }
        hint={
          <p>
            For step three, look at the wall clock the benchmark reports
            first, and then try to find those milliseconds in the profile.
          </p>
        }
        solution={
          <>
            <p>
              Step two usually disagrees with step one, sometimes completely.
              That is the lesson: a short profile is not a smaller version of
              a long one, it is a different and mostly fictional picture.
            </p>
            <p>
              In step three the benchmark&apos;s ns/op jumps by ten million
              and the sleep is <strong>not in the profile at all</strong>. A
              sleeping goroutine is not running, so it is never sampled. Every
              kind of waiting looks like this, which is why the question
              &quot;is CPU actually high?&quot; comes before opening a
              profiler.
            </p>
          </>
        }
      />

      <KnowledgeCheck
        question="A handler takes 400ms. Its CPU profile shows almost no time anywhere in your code. What does that tell you?"
        choices={[
          {
            label: 'The time is spent waiting, and a CPU profile cannot see it',
            correct: true,
            feedback:
              'Database round trips, lock contention, network calls. Look at queries per request in Pulse, then the block and mutex profiles.',
          },
          {
            label: 'The profile was captured wrongly',
            feedback:
              'Possible, and worth ruling out with a longer benchtime. But an empty CPU profile on a slow handler is usually a true answer to the question you asked.',
          },
          {
            label: 'The compiler optimised the work away',
            feedback:
              'That happens in microbenchmarks with unused results, not in a handler whose response somebody reads.',
          },
          {
            label: 'The handler is fine and the measurement is wrong',
            feedback:
              'The 400ms is real; the user waits it. What is wrong is the choice of instrument.',
          },
        ]}
      />
    </>
  )
}
