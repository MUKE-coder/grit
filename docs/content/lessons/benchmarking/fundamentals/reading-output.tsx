import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        You wrote a benchmark and got numbers. Now make them mean
        something. This lesson is fluency with{' '}
        <code>ns/op</code>, <code>B/op</code>, and{' '}
        <code>allocs/op</code> — what each tells you, when each
        matters, and what red flags to look for.
      </p>

      <h2>ns/op — the headline</h2>
      <p>
        Nanoseconds per operation. The number you optimise for, most of
        the time.
      </p>
      <ul>
        <li>
          <strong>Under 100ns</strong> — fast. Probably a pure function
          with no syscalls.
        </li>
        <li>
          <strong>100ns – 1µs (1000ns)</strong> — typical for string
          ops, simple JSON marshal, small loops.
        </li>
        <li>
          <strong>1µs – 100µs</strong> — non-trivial work. Allocation,
          regex, complex parsing.
        </li>
        <li>
          <strong>100µs – 1ms</strong> — usually some I/O or a big
          loop. If you didn&apos;t expect this, investigate.
        </li>
        <li>
          <strong>1ms+</strong> — likely the wrong target for
          microbenchmarking. Profile the broader system instead.
        </li>
      </ul>

      <h2>B/op — bytes per op</h2>
      <p>
        How many bytes of heap memory the operation allocates per
        iteration.
      </p>
      <ul>
        <li>
          <strong>0 B/op</strong> — zero allocations. Often achievable
          for pure ops with stack-only data. The dream for hot paths.
        </li>
        <li>
          <strong>16-64 B/op</strong> — typical for ops that build a
          small string, return a struct on the heap, or take a slice.
        </li>
        <li>
          <strong>1KB+ /op</strong> — large; either you copied data or
          allocated a buffer. Usually a flag for optimisation.
        </li>
      </ul>

      <h2>allocs/op — number of heap allocations</h2>
      <p>
        Sometimes more important than bytes. Each allocation:
      </p>
      <ul>
        <li>Costs CPU to satisfy.</li>
        <li>Increases the work the GC has to do later.</li>
        <li>Causes a tiny pause when the GC runs.</li>
      </ul>
      <p>
        For libraries used a million times per second (logging,
        request routing), driving <code>allocs/op</code> to 0 or 1 is a
        common target.
      </p>

      <h2>Reading a realistic output</h2>
      <CodeBlock
        language="text"
        code={`BenchmarkPluralize-10        5102937	   226 ns/op   64 B/op   2 allocs/op
BenchmarkGoType-10           7891234	   151 ns/op   48 B/op   1 allocs/op
BenchmarkZodType-10           892145	  1342 ns/op  256 B/op   8 allocs/op   ← red flag
BenchmarkGORMTag-10          2912034	   411 ns/op  128 B/op   3 allocs/op
BenchmarkInjectBefore-10      503456	  2385 ns/op  512 B/op  12 allocs/op  ← bigger flag
BenchmarkParseInline-10      1056432	  1147 ns/op  192 B/op   5 allocs/op`}
      />
      <p>
        Two red flags worth noting:
      </p>
      <ul>
        <li>
          <code>BenchmarkZodType</code> at 1.3µs with 8 allocs. Likely
          building intermediate strings.{' '}
          <code>strings.Builder</code> would drop both numbers.
        </li>
        <li>
          <code>BenchmarkInjectBefore</code> at 2.4µs with 12 allocs.
          12 allocations PER OP is a lot — almost certainly a string
          search-and-replace with many intermediate values. A scanner
          or single buffer would help.
        </li>
      </ul>

      <h2>Spotting hidden allocations</h2>
      <CodeBlock
        language="go"
        code={`// Allocates a new string every call (string is immutable in Go)
func concat1(s string) string {
  return s + " · v1"
}
// Bench: 30 ns/op, 16 B/op, 1 allocs/op

// Allocates a new slice on every append (cap growth)
func concat2(words []string) string {
  out := ""
  for _, w := range words {
    out += w + " "
  }
  return out
}
// Bench (10 words): 290 ns/op, 192 B/op, 9 allocs/op  ← yikes

// strings.Builder reuses one growing buffer
func concat3(words []string) string {
  var b strings.Builder
  for _, w := range words {
    b.WriteString(w)
    b.WriteByte(' ')
  }
  return b.String()
}
// Bench (10 words): 110 ns/op, 32 B/op, 2 allocs/op  ← much better`}
      />
      <p>
        Same logical work, ~3x faster + 6x fewer allocations. The
        benchmark numbers told you exactly where the inefficiency was.
      </p>

      <h2>What &quot;noise&quot; looks like</h2>
      <CodeBlock
        language="text"
        code={`# Run 5 times, results swinging 50%:
BenchmarkX-10  5102937   226 ns/op
BenchmarkX-10  4992014   289 ns/op
BenchmarkX-10  5023412   232 ns/op
BenchmarkX-10  4123912   312 ns/op
BenchmarkX-10  5202345   228 ns/op`}
      />
      <p>
        Noise. Your laptop is doing something else under load (browser,
        Spotlight, antivirus). Don&apos;t trust these numbers. Fix by:
      </p>
      <ul>
        <li>Closing background apps.</li>
        <li>Running on a quiet machine (CI, dedicated box).</li>
        <li>Using <code>-count=10</code> + benchstat (chapter 3).</li>
        <li>Avoiding battery / thermal throttling.</li>
      </ul>

      <TipBox tone="info">
        <strong>The compiler can elide your work.</strong> If your
        benchmark&apos;s return value is unused, the Go compiler may
        skip the call entirely — making your benchmark say 0ns. The
        fix: assign the result to a package-level variable so the
        compiler can&apos;t prove the work is dead.
        <CodeBlock
          language="go"
          code={`var sink string  // package-level

func BenchmarkPluralize(b *testing.B) {
  var s string
  for i := 0; i < b.N; i++ { s = Pluralize("cat") }
  sink = s  // prevent dead-code elimination
}`}
        />
      </TipBox>

      <h2>What ns/op alone can&apos;t tell you</h2>
      <p>
        Benchmarks measure ONE thing in isolation. They don&apos;t
        tell you:
      </p>
      <ul>
        <li>
          <strong>How often this function is called.</strong> A 5µs
          function called 10× a day is fine. A 50ns function called a
          million times per request matters.
        </li>
        <li>
          <strong>What the rest of the system is doing while it
          runs.</strong> Lock contention, GC pauses, network jitter all
          dwarf microbenchmark differences.
        </li>
        <li>
          <strong>Real-world data shapes.</strong> Your benchmark hits
          one input; production hits a distribution.
        </li>
      </ul>
      <p>
        Profiling (next chapter) helps with the &quot;how often&quot;
        question. Load testing (the K6 course) helps with the
        &quot;real world&quot; question. Use the right tool for the
        right question.
      </p>

      <KnowledgeCheck
        question="Benchmark A: 50 ns/op, 0 B/op, 0 allocs/op. Benchmark B: 30 ns/op, 32 B/op, 1 allocs/op. Which is faster?"
        choices={[
          {
            label: 'B — it has lower ns/op',
            feedback:
              "True at the per-op level, but the answer is &quot;it depends on call frequency&quot;.",
          },
          {
            label: 'A — zero allocations is always better',
            feedback:
              "Not always. If you call this once, A is slower overall.",
          },
          {
            label: 'It depends — for a hot loop run millions of times, A wins because GC pressure adds up; for an infrequent call B wins',
            correct: true,
            feedback:
              "Right — context matters. ns/op is the headline; allocations are the long-tail cost. Choose based on call frequency. Profile in real conditions to know.",
          },
          {
            label: 'They are the same speed',
            feedback: 'Different ns/op means different speed at the call site.',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Interpret real numbers:</p>
            <ol>
              <li>
                Take the benchmarks you wrote last lesson. Look at
                each <code>B/op</code> number.
              </li>
              <li>
                Find the benchmark with the highest{' '}
                <code>allocs/op</code>. Read the source. Make ONE
                hypothesis about why it allocates that much.
              </li>
              <li>
                Try one quick change (use <code>strings.Builder</code>,
                pre-allocate a slice with <code>make([]T, 0, cap)</code>,
                or avoid an intermediate string). Re-run.
              </li>
              <li>
                Did the numbers move? Write down: how much, in which
                direction, and your guess at why.
              </li>
              <li>Paste before / after in notes.md.</li>
            </ol>
          </>
        }
        hint={
          <>
            For <code>strings.Builder</code>: in any code where you do{' '}
            <code>s += foo</code> inside a loop, replace with a Builder
            and call <code>.String()</code> at the end. Almost
            guaranteed to drop both ns/op and allocs/op.
          </>
        }
        solution={
          <>
            <p>
              You should have moved at least one number. The
              important part is the DISCIPLINE: hypothesis →
              measurement → conclusion. Most performance work is
              vibes; you&apos;re doing it the legit way.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 2 — <strong>Profiling with pprof</strong>. Benchmarks
        tell you ONE function&apos;s cost; pprof tells you where a
        whole program spends its time. The CPU + heap flame graph is
        the highest-bandwidth tool in Go.
      </p>
    </>
  )
}
