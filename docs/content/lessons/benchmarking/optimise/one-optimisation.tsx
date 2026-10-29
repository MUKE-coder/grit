import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Everything in this course has been instruments. This lesson is the
        whole loop once, end to end, on one real change, with the parts
        people skip left in.
      </p>

      <h2>Pick the thing that is worth it</h2>
      <p>
        Two questions before any measurement, and they rule out most
        candidates:
      </p>

      <CodeBlock
        language="text"
        code={`How often does it run?     a function called once at boot is not
                           worth a day, whatever the profile says
What fraction is it?       Amdahl: making 5% of the time twice as fast
                           saves 2.5%. Find the big slice first`}
      />

      <p>
        In a Grit API the per request hot path is the place to look, because
        everything there is multiplied by your whole traffic. That is the
        argument <code>internal/handlers/perf_test.go</code> makes for the
        paths it budgets: token verification, parameter binding, the list
        envelope.
      </p>

      <h2>The loop, with the numbers</h2>
      <CodeBlock
        language="bash"
        code={`# 1. Baseline, twice, so you know your noise floor.
go test -run XXX -bench BenchmarkListEnvelope -benchmem \\
  -count=10 ./internal/handlers/ > before.txt

# 2. Find out where it goes.
go test -run XXX -bench BenchmarkListEnvelope -benchtime 5s \\
  -memprofile mem.out ./internal/handlers/
go tool pprof -sample_index=alloc_space -http=: mem.out

# 3. ONE change.

# 4. Same command, same machine, same count.
go test -run XXX -bench BenchmarkListEnvelope -benchmem \\
  -count=10 ./internal/handlers/ > after.txt
benchstat before.txt after.txt`}
      />

      <p>
        Step three is one change. Two at once gives one number and no
        attribution, and the usual outcome is that both get kept: the one
        that helped and the one that quietly cost something.
      </p>

      <h2>Worked: the list envelope</h2>
      <p>
        Baseline <code>14790 ns/op, 11390 B/op, 36 allocs/op</code>. The
        alloc_space profile points at the slice the handler builds, grown
        one row at a time:
      </p>

      <CodeBlock
        language="go"
        code={`// Before: grows and copies as it goes.
var out []PostResponse
for _, p := range posts {
    out = append(out, toResponse(p))
}

// After: one allocation, because the length is known.
out := make([]PostResponse, 0, len(posts))
for _, p := range posts {
    out = append(out, toResponse(p))
}`}
      />

      <p>
        A one line change, and the kind that is worth the ceremony only
        because the loop runs on every list request in the application.
      </p>

      <h2>Then write it down, and move the budget</h2>
      <CodeBlock
        language="text"
        filename="notes.md"
        code={`2026-10-29  BenchmarkListEnvelope, -count=10, laptop on mains
  before:  14.79µ ± 3%   11.39Ki   36 allocs
  change:  preallocate the response slice with len(posts)
  after:    9.12µ ± 2%    6.21Ki   22 allocs   p=0.000
  kept, and the perf_test budget lowered 36 -> 22 in the same commit`}
      />

      <p>
        Lowering the budget in the same commit is the part that makes the
        work stick. Left at 36, the next change can put fourteen allocations
        back and nothing will say so.
      </p>

      <TipBox tone="warning">
        <p>
          If a change raises a budget, that can be correct: a feature costs
          what it costs. Raise it in the same commit, with the reason in the
          comment, so the next person sees the new number was a decision and
          not a drift.
        </p>
      </TipBox>

      <h2>When to stop</h2>
      <ul>
        <li>
          <strong>The change is within the noise floor.</strong> Revert it.
          An unmeasurable gain is a maintenance cost with no benefit.
        </li>
        <li>
          <strong>The code got significantly harder to read.</strong> A 3%
          win for an unreadable function is a bad trade, and most of the
          time somebody undoes it later without knowing why it was there.
        </li>
        <li>
          <strong>You are optimising something that is not the bottleneck.</strong>{' '}
          A function at 2% of the time cannot give you more than 2%.
        </li>
        <li>
          <strong>The real fix is elsewhere.</strong> Shaving allocations
          off a handler that issues 21 queries per request is polishing the
          wrong surface. Fix the N+1 first.
        </li>
      </ul>

      <Exercise
        prompt={
          <>
            <p>
              <strong>One.</strong> Run the whole loop above on one function
              in your own project, start to finish, including the two
              baseline runs.
            </p>
            <p>
              <strong>Two.</strong> Write the notes entry: before, change,
              after, p value, and kept or reverted.
            </p>
            <p>
              <strong>Three.</strong> If you kept it and the function has a
              budget in <code>perf_test.go</code>, lower it and confirm the
              test still passes.
            </p>
          </>
        }
        hint={
          <p>
            If the change is not significant, write the entry anyway and
            revert. A recorded null result is worth more than an unrecorded
            win, because it stops the next person repeating it.
          </p>
        }
        solution={
          <>
            <p>
              The deliverable is the notes entry, not the speedup. A change
              with a number, a p value and a date can be defended, reproduced
              and reversed. A change without one is a thing in the codebase
              that nobody can evaluate.
            </p>
            <p>
              With the budget lowered, the improvement is defended by CI: the
              allocation test fails if it comes back. That is the difference
              between making something faster once and keeping it fast.
            </p>
          </>
        }
      />

      <KnowledgeCheck
        question="Your change cuts allocations by 40%, and benchstat reports ~ on sec/op. Do you keep it?"
        choices={[
          {
            label: 'Yes, if it is no harder to read: fewer allocations is less collector work under real concurrency',
            correct: true,
            feedback:
              'A single-threaded benchmark hides GC pressure that shows up with many goroutines. Keep it, lower the budget, and verify under load rather than claiming the time saving.',
          },
          {
            label: 'No: if the time did not move, there is no benefit',
            feedback:
              'The benchmark runs one goroutine. Allocation cost is paid collectively, and the place it shows is a load test, not here.',
          },
          {
            label: 'Yes, and report it as a 40% speedup',
            feedback:
              'Keep it, but a 40% allocation reduction is not a 40% speedup, and benchstat has just said the time did not move measurably.',
          },
          {
            label: 'Re-run with -count=50 until sec/op becomes significant',
            feedback:
              'More samples is reasonable once, decided in advance. Running until the number agrees with you is not a measurement.',
          },
        ]}
      />
    </>
  )
}
