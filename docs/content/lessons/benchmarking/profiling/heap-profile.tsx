import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Most Go performance work is allocation work. The collector is fast,
        but it is not free, and the cheapest way to spend less time in it is
        to give it less to do. A heap profile says where the allocations come
        from.
      </p>

      <h2>Two different questions</h2>
      <CodeBlock
        language="bash"
        code={`go test -run XXX -bench BenchmarkListEnvelope -benchtime 5s \\
  -memprofile mem.out ./internal/handlers/

go tool pprof -sample_index=alloc_space -http=: mem.out   # what allocated
go tool pprof -sample_index=inuse_space -http=: mem.out   # what is still held`}
      />

      <ul>
        <li>
          <strong>alloc_space</strong> is everything allocated over the run,
          including what was freed immediately. This is the one for garbage
          pressure: a function allocating a megabyte per request and freeing
          it still costs you collector time.
        </li>
        <li>
          <strong>inuse_space</strong> is what is still reachable at the
          sample. This is the one for leaks: memory that should have gone and
          did not.
        </li>
      </ul>

      <p>
        Picking the wrong one is the usual mistake. Hunting a leak in
        alloc_space shows you your busiest function, which is almost never
        the leak.
      </p>

      <h2>Start with the numbers you already have</h2>
      <p>
        <code>-benchmem</code> is free and often enough. From a scaffolded
        Grit API:
      </p>

      <CodeBlock
        language="text"
        code={`BenchmarkHealthCheck-14        11465 ns/op    6878 B/op    25 allocs/op
BenchmarkVerifyAccessToken-14   8330 ns/op    3338 B/op    41 allocs/op
BenchmarkBindListParams-14      8355 ns/op    7974 B/op    40 allocs/op
BenchmarkListEnvelope-14       14790 ns/op   11390 B/op    36 allocs/op`}
      />

      <p>
        Note that verifying a token allocates more times than building a list
        envelope while allocating a third of the bytes. Count and size are
        different costs: many small allocations are mostly collector
        bookkeeping, one large one is mostly the copy. The profile tells you
        which you have.
      </p>

      <h2>Allocation counts as assertions</h2>
      <p>
        Grit scaffolds <code>internal/handlers/perf_test.go</code>, and the
        reasoning in it is worth more than the numbers:
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/perf_test.go"
        code={`// A benchmark tells you how long something took on the machine that ran
// it, which on a shared CI runner is a number that moves on its own. An
// allocation count does not move with the machine or the load, which is
// what makes these assertions rather than measurements, and why they can
// fail a build when a benchmark cannot.`}
      />

      <p>
        That is why the budgets exist and the timings do not get asserted.
        CI runs the benchmarks with <code>-benchtime 1x</code>, purely to
        prove they still compile and run, and the allocation budgets are the
        part that can actually fail.
      </p>

      <TipBox tone="warning">
        <p>
          Allocation counts move with the Go toolchain, by more than you
          would guess. Token verification allocates 56 times under go1.26.9
          and 41 under go1.27.2: the same code, a newer runtime. So a budget
          is set from the <em>oldest</em> toolchain the project supports,
          because that is the one that allocates most. Set from whatever is
          on your laptop, a budget passes locally and fails in CI.
        </p>
      </TipBox>

      <h2>Where the allocations usually are</h2>
      <ul>
        <li>
          <strong>Growing a slice in a loop.</strong> Each growth copies.{' '}
          <code>make([]T, 0, n)</code> when you know n turns many
          allocations into one.
        </li>
        <li>
          <strong>String concatenation in a loop.</strong> Every{' '}
          <code>+=</code> allocates a new string.{' '}
          <code>strings.Builder</code> does not.
        </li>
        <li>
          <strong>
            <code>interface&#123;&#125;</code> and <code>any</code>.
          </strong>{' '}
          Putting a value in an interface usually heap allocates it. A
          response built as <code>map[string]any</code> allocates once per
          value.
        </li>
        <li>
          <strong>Taking the address of a local.</strong> If a pointer
          escapes the function, the value moves to the heap.{' '}
          <code>go build -gcflags=-m</code> prints every such decision.
        </li>
      </ul>

      <h2>From the running API</h2>
      <CodeBlock
        language="bash"
        code={`go tool pprof -http=: http://localhost:8080/pulse/api/profile/heap

# and the one that answers "what changed", which is usually the question
go tool pprof -http=: -diff_base before.out after.out`}
      />

      <p>
        <code>-diff_base</code> is the most useful flag in pprof and the
        least used. Two captures, one graph, showing only what moved between
        them.
      </p>

      <Exercise
        prompt={
          <>
            <p>
              <strong>One.</strong> Run the handler benchmarks with{' '}
              <code>-benchmem</code> and write down allocs/op for each.
            </p>
            <p>
              <strong>Two.</strong> Take a memory profile of the one with the
              most allocations and find the top three sources by alloc_space.
            </p>
            <p>
              <strong>Three.</strong> Change one of them: preallocate a
              slice, or swap a concatenation for a{' '}
              <code>strings.Builder</code>. Re-run and compare allocs/op.
            </p>
          </>
        }
        hint={
          <p>
            In step three, watch ns/op as well. Fewer allocations usually
            means less time, and when it does not, that is a result worth
            keeping too.
          </p>
        }
        solution={
          <>
            <p>
              A preallocated slice typically removes several allocations and a
              proportional slice of the time. If allocs/op falls and ns/op
              does not, the allocations were not the bottleneck in that
              function, which is a real finding and the reason to measure
              rather than assume.
            </p>
            <p>
              Whatever the result, if the function has an entry in{' '}
              <code>perf_test.go</code>, lower its budget in the same commit.
              A budget left at the old number is a regression nobody will
              catch.
            </p>
          </>
        }
      />

      <KnowledgeCheck
        question="Memory climbs steadily over a four hour soak. Which profile, and which sample index?"
        choices={[
          {
            label: 'Heap profile, inuse_space, and compare two captures with -diff_base',
            correct: true,
            feedback:
              'inuse_space is what is still reachable, which is what a leak means. Two captures an hour apart, diffed, shows what grew rather than what is big.',
          },
          {
            label: 'Heap profile, alloc_space',
            feedback:
              'That shows your busiest allocator, which is normal behaviour and almost never the leak. Freed memory counts in alloc_space.',
          },
          {
            label: 'CPU profile, since the collector will be busy',
            feedback:
              'It may show more GC time as a symptom, but it will not tell you what is being retained.',
          },
          {
            label: 'Increase GOGC so the collector runs less often',
            feedback:
              'That makes a leak reach the memory limit sooner, with the same cause unexamined.',
          },
        ]}
      />
    </>
  )
}
