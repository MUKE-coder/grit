import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Run the same benchmark twice and you get two different numbers. Run
        it before and after a change and you get two different numbers, which
        may or may not be about the change. <code>benchstat</code> is the
        tool that tells the difference, and it is the one thing standing
        between performance work and a vibe.
      </p>

      <h2>Install it, then use counts</h2>
      <CodeBlock
        language="bash"
        code={`go install golang.org/x/perf/cmd/benchstat@latest

cd apps/api
go test -run XXX -bench BenchmarkListEnvelope -benchmem \\
  -count=10 ./internal/handlers/ > before.txt

# make the change, then

go test -run XXX -bench BenchmarkListEnvelope -benchmem \\
  -count=10 ./internal/handlers/ > after.txt

benchstat before.txt after.txt`}
      />

      <p>
        <code>-count=10</code> is the point of the exercise. One run gives
        benchstat nothing to work with: with a single sample there is no
        spread, so there is no way to say whether a difference is larger than
        the noise. Ten is the usual minimum, and more for a benchmark that
        is naturally variable.
      </p>

      <h2>Reading the output</h2>
      <CodeBlock
        language="text"
        code={`                  │  before.txt  │             after.txt              │
                  │    sec/op    │   sec/op     vs base               │
ListEnvelope-14     14.79µ ± 3%    9.12µ ± 2%  -38.34% (p=0.000 n=10)

                  │  before.txt  │             after.txt              │
                  │     B/op     │    B/op      vs base               │
ListEnvelope-14    11.39Ki ± 0%   6.21Ki ± 0%  -45.48% (p=0.000 n=10)

                  │  before.txt  │            after.txt             │
                  │  allocs/op   │ allocs/op   vs base              │
ListEnvelope-14      36.00 ± 0%   22.00 ± 0%  -38.89% (p=0.000 n=10)`}
      />

      <ul>
        <li>
          <strong>± 3%</strong> is the spread across the ten runs. A large
          spread means the measurement is unstable, and the first thing to
          fix is the measurement.
        </li>
        <li>
          <strong>p=0.000</strong> is the probability of seeing this
          difference if the change did nothing. Below 0.05 it is reported;
          above it, benchstat prints <code>~</code>.
        </li>
        <li>
          <strong>n=10</strong> is how many samples survived. If benchstat
          discarded some as outliers it says so, and that is worth reading.
        </li>
      </ul>

      <TipBox tone="danger">
        <p>
          A <code>~</code> in the vs base column means &quot;no difference
          was detected&quot;. It does not mean the change is free and it does
          not mean it helped by a small amount. It means this experiment
          cannot tell, and the honest responses are more samples or a
          revert, not a sentence in the pull request saying it is slightly
          faster.
        </p>
      </TipBox>

      <h2>Making the machine quiet enough</h2>
      <p>
        A ± of 15% makes any change under 15% unmeasurable, so it is worth
        spending a minute on:
      </p>

      <CodeBlock
        language="bash"
        code={`# Close the other work: browsers, dev servers, Docker, the other build.
# Then pin the parallelism so it does not vary with what else is running.
go test -run XXX -bench . -benchmem -count=10 -cpu=4 ./internal/handlers/

# A laptop on battery clocks down mid-run. Plug it in.`}
      />

      <p>
        Two benchmarks in the same package also interfere: one leaves the
        caches warm and the allocator in a different state for the next. Use{' '}
        <code>-bench</code> with an exact name when the number matters.
      </p>

      <h2>Why CI does not do this</h2>
      <p>
        Grit&apos;s scaffolded <code>ci.yml</code> runs the benchmarks with{' '}
        <code>-benchtime 1x</code> and asserts nothing about their timing,
        with a comment saying why: timing a shared CI runner tells you about
        the runner. What CI asserts instead is the allocation budgets, which
        do not move with the machine.
      </p>

      <p>
        Timing comparisons belong on a machine you control, run deliberately,
        with the result written down.
      </p>

      <Exercise
        prompt={
          <>
            <p>
              <strong>One.</strong> Run a benchmark with{' '}
              <code>-count=10</code> twice, changing nothing, into{' '}
              <code>a.txt</code> and <code>b.txt</code>. Run{' '}
              <code>benchstat a.txt b.txt</code>.
            </p>
            <p>
              <strong>Two.</strong> Note the ± figure and whether benchstat
              reports a difference between two identical versions.
            </p>
            <p>
              <strong>Three.</strong> Repeat with a browser and a dev server
              running. Compare the ± figures.
            </p>
          </>
        }
        hint={
          <p>
            Step two should print <code>~</code>. If it reports a
            significant difference between two runs of identical code, your
            measurement environment is the thing you are measuring.
          </p>
        }
        solution={
          <>
            <p>
              Step one is a calibration: it tells you what &quot;no
              change&quot; looks like on your machine, so you can recognise
              it later.
            </p>
            <p>
              Step three usually doubles or triples the ±, and that is the
              argument for closing things. A noisy machine does not add noise
              evenly, it adds occasional very slow iterations, which is also
              why benchstat works on the distribution rather than the mean.
            </p>
          </>
        }
      />

      <KnowledgeCheck
        question="benchstat reports -4.2% with p=0.21 and n=10. What do you conclude?"
        choices={[
          {
            label: 'Nothing yet: p=0.21 means the difference is indistinguishable from noise',
            correct: true,
            feedback:
              'Above 0.05, benchstat is saying it cannot tell. More samples, or a quieter machine, or accept that the change is not measurable here.',
          },
          {
            label: 'A 4.2% improvement, small but real',
            feedback:
              'The 4.2% is the difference between two noisy means. The p value exists precisely to stop you reading it as a result.',
          },
          {
            label: 'The change made things worse',
            feedback:
              'The sign points the other way, and in any case it is not distinguishable from noise.',
          },
          {
            label: 'Re-run until p drops below 0.05',
            feedback:
              'Running until the answer is the one you want is how you manufacture a result. Increase -count as a decision made before you look.',
          },
        ]}
      />
    </>
  )
}
