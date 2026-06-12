import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Go has the cleanest built-in benchmarking I&apos;ve ever used.
        Write a function with a name starting <code>Benchmark</code>,
        loop over <code>b.N</code>, run <code>go test -bench=.</code>.
        Done. This lesson is that whole loop, end-to-end.
      </p>

      <h2>The shape of a benchmark</h2>
      <CodeBlock
        language="go"
        filename="internal/generate/pluralize_test.go"
        code={`package generate

import "testing"

func BenchmarkPluralize(b *testing.B) {
  for i := 0; i < b.N; i++ {
    Pluralize("category")
  }
}`}
      />
      <p>
        Three things to absorb:
      </p>
      <ul>
        <li>
          <strong>Name starts with <code>Benchmark</code></strong> — Go
          test tool finds it.
        </li>
        <li>
          <strong>Takes <code>*testing.B</code></strong>.
        </li>
        <li>
          <strong>Loops <code>b.N</code> times</strong>. Go decides{' '}
          <code>b.N</code> at runtime — it ramps until the benchmark
          runs long enough to be statistically meaningful (usually a
          second).
        </li>
      </ul>

      <h2>Run it</h2>
      <CodeBlock
        language="bash"
        code={`go test -bench=BenchmarkPluralize -benchmem -count=3 ./internal/generate/`}
      />
      <p>Three flags worth knowing:</p>
      <ul>
        <li>
          <code>-bench=BenchmarkXxx</code> — which benchmarks to run.
          <code>=.</code> matches all.
        </li>
        <li>
          <code>-benchmem</code> — also report bytes/op and allocs/op.
          Use it every time.
        </li>
        <li>
          <code>-count=3</code> — run each benchmark 3 times. You want
          this — single runs are noisy.
        </li>
      </ul>

      <h2>The output</h2>
      <CodeBlock
        language="text"
        code={`goos: darwin
goarch: arm64
pkg: github.com/MUKE-coder/grit/internal/generate
BenchmarkPluralize-10  	5012461	   228.4 ns/op	    64 B/op	   2 allocs/op
BenchmarkPluralize-10  	5102937	   226.1 ns/op	    64 B/op	   2 allocs/op
BenchmarkPluralize-10  	5098102	   227.8 ns/op	    64 B/op	   2 allocs/op
PASS
ok  	github.com/MUKE-coder/grit/internal/generate	4.523s`}
      />
      <p>Decode each column:</p>
      <ul>
        <li>
          <strong><code>BenchmarkPluralize-10</code></strong> — the{' '}
          <code>-10</code> is GOMAXPROCS (cores). Same number across runs
          means same machine.
        </li>
        <li>
          <strong><code>5012461</code></strong> — the b.N Go settled
          on. ~5 million iterations.
        </li>
        <li>
          <strong><code>228.4 ns/op</code></strong> — average{' '}
          nanoseconds per operation. THIS is the headline metric.
        </li>
        <li>
          <strong><code>64 B/op</code></strong> — bytes allocated per
          op. Lower = less GC pressure.
        </li>
        <li>
          <strong><code>2 allocs/op</code></strong> — number of heap
          allocations per op. Lower = even less GC pressure.
        </li>
      </ul>

      <h2>Why &quot;per op&quot; not &quot;total time&quot;</h2>
      <p>
        b.N changes between runs (or even between machines), so
        absolute total time is meaningless. What matters is the per-op
        average, which IS comparable across runs and machines (with
        caveats — see &quot;noise&quot; below).
      </p>

      <h2>Multiple benchmarks in one file</h2>
      <CodeBlock
        language="go"
        code={`func BenchmarkPluralize_short(b *testing.B) {
  for i := 0; i < b.N; i++ { Pluralize("dog") }
}

func BenchmarkPluralize_long(b *testing.B) {
  for i := 0; i < b.N; i++ { Pluralize("administrator") }
}

func BenchmarkPluralize_irregular(b *testing.B) {
  for i := 0; i < b.N; i++ { Pluralize("child") }
}`}
      />
      <p>
        Compare different input shapes side by side. Often the
        revelation: &quot;long words are 3x slower&quot;, &quot;irregular
        words trigger a fallback path&quot;.
      </p>

      <h2>What about setup that shouldn&apos;t count?</h2>
      <CodeBlock
        language="go"
        code={`func BenchmarkParseInline(b *testing.B) {
  // Setup — NOT measured
  input := buildLargeStringFromFile()

  b.ResetTimer()   // restart the clock from here

  for i := 0; i < b.N; i++ {
    ParseInlineFields(input)
  }
}`}
      />
      <p>
        <code>b.ResetTimer()</code> says &quot;everything before
        this is setup; start measuring now&quot;. Critical for
        benchmarks where the setup is expensive.
      </p>

      <h2>Sub-benchmarks (table-driven)</h2>
      <CodeBlock
        language="go"
        code={`func BenchmarkPluralize_table(b *testing.B) {
  for _, name := range []string{"dog", "cat", "administrator", "child", "deer"} {
    b.Run(name, func(b *testing.B) {
      for i := 0; i < b.N; i++ {
        Pluralize(name)
      }
    })
  }
}`}
      />
      <p>
        Output shows each as its own line:
      </p>
      <CodeBlock
        language="text"
        code={`BenchmarkPluralize_table/dog-10           5102937   226.1 ns/op
BenchmarkPluralize_table/cat-10           5012461   228.4 ns/op
BenchmarkPluralize_table/administrator-10 4012461   299.8 ns/op
...`}
      />
      <p>
        One file, comprehensive coverage of the function&apos;s input
        space. The Go convention for serious benchmarking.
      </p>

      <TipBox tone="warning">
        <strong>Benchmarks are noisy on a laptop.</strong> Background
        Slack, Spotlight indexing, thermal throttling — all add jitter.
        For serious numbers: kill background apps, close browser
        tabs, plug in to power, run <code>-count=10</code> and use{' '}
        <code>benchstat</code> (next chapter) to get a statistically
        sound comparison.
      </TipBox>

      <KnowledgeCheck
        question="You write `for i := 0; i < 100; i++` instead of `for i := 0; i < b.N; i++` in a benchmark. What happens?"
        choices={[
          {
            label: 'It runs faster because you control the iterations',
            feedback: 'It doesn’t run faster — and you lose accuracy badly.',
          },
          {
            label: 'Go can\'t adapt the iteration count, so the result is unreliable: too few iterations get rounded noise; the per-op number is meaningless',
            correct: true,
            feedback:
              "Right — Go ramps b.N to get a statistically meaningful sample. Hard-coding 100 means you might run for 22µs total, which is mostly noise. Use b.N. Always.",
          },
          {
            label: 'Go panics',
            feedback: "It doesn't — it just gives you a number you can't trust.",
          },
          {
            label: 'Nothing — Go ignores hand-coded loops',
            feedback: 'It runs your loop. But b.N is the contract you must honor for the result to be meaningful.',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Write your first three benchmarks:</p>
            <ol>
              <li>
                Pick a function in your service layer. Anything pure
                and side-effect-free is easiest (a parser, a
                formatter).
              </li>
              <li>
                Add{' '}
                <code>func BenchmarkX(b *testing.B) &#123; for i := 0; i &lt; b.N; i++ &#123; X() &#125; &#125;</code>{' '}
                to the test file.
              </li>
              <li>
                Run with{' '}
                <code>go test -bench=. -benchmem -count=3 ./your/pkg/</code>.
              </li>
              <li>
                Make a 3-input table-driven version (sub-benchmarks).
              </li>
              <li>
                Paste the output in <code>notes.md</code>. Circle the
                most interesting line — the one that&apos;s slower
                than you expected.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            For your first benchmark, pick something simple but
            non-trivial — e.g., your <code>Pluralize</code>, your{' '}
            <code>SlugifyEnabled</code>, your JSON-tag parser.
            Numbers under 1µs are reasonable for these.
          </>
        }
        solution={
          <>
            <p>
              You should have 3 working benchmarks and a feel for the
              output. Next lesson: how to read ns/op, B/op, allocs/op
              for what they really tell you about your code.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Next lesson — <strong>Reading the output</strong>. The numbers
        tell a story. After this lesson you&apos;ll be fluent enough to
        spot &quot;there&apos;s a hidden allocation&quot; or &quot;this
        function is O(n²) and we didn&apos;t notice&quot;.
      </p>
    </>
  )
}
