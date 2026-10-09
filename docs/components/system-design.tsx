import type { ReactNode } from 'react'
import { CodeBlock } from '@/components/code-block'
import { SystemFlow } from '@/components/system-flow'
import type { ApiGroup, Estimate, Labelled, Pair, SystemDesign } from '@/config/systems-types'

/* ─────────────────────────────────────────────────────────────
   The template every system design page is rendered through.

   Sections appear in a fixed order and a section with no data is
   skipped rather than left as an empty heading, so "not here"
   reads as "does not apply" and never as "unfinished".
   ───────────────────────────────────────────────────────────── */

const MONO = 'text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded'

function Section({
  n,
  id,
  title,
  children,
}: {
  n: number
  id: string
  title: string
  children: ReactNode
}) {
  return (
    <section className="mt-12 scroll-mt-24" id={id}>
      <h2 className="text-2xl font-bold tracking-tight">
        <span className="mr-2 text-muted-foreground tabular-nums">{n}.</span>
        {title}
      </h2>
      <div className="mt-4">{children}</div>
    </section>
  )
}

function Bullets({ items }: { items: string[] }) {
  return (
    <ul className="my-3 space-y-2">
      {items.map((t, i) => (
        <li key={i} className="flex gap-2.5 text-[15px] leading-relaxed text-muted-foreground">
          <span aria-hidden className="mt-[9px] h-1 w-1 shrink-0 rounded-full bg-primary/70" />
          <span>{t}</span>
        </li>
      ))}
    </ul>
  )
}

function LabelledList({ items }: { items: Labelled[] }) {
  return (
    <ul className="my-3 space-y-2.5">
      {items.map((it, i) => (
        <li key={i} className="flex gap-2.5 text-[15px] leading-relaxed text-muted-foreground">
          <span aria-hidden className="mt-[9px] h-1 w-1 shrink-0 rounded-full bg-primary/70" />
          <span>
            <strong className="font-semibold text-foreground">{it.label}.</strong> {it.text}
          </span>
        </li>
      ))}
    </ul>
  )
}

function TwoCol({ headers, rows }: { headers: [string, string]; rows: Pair[] }) {
  return (
    <div className="my-4 overflow-x-auto rounded-lg border border-border/50">
      <table className="w-full text-left text-sm">
        <thead className="bg-card/60">
          <tr>
            {headers.map((h) => (
              <th key={h} className="px-4 py-2.5 font-semibold">
                {h}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map(([a, b], i) => (
            <tr key={i} className="border-t border-border/40">
              <td className="px-4 py-2.5 align-top font-medium">{a}</td>
              <td className="px-4 py-2.5 align-top text-muted-foreground">{b}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function Working({ e }: { e: Estimate }) {
  return (
    <div className="my-4">
      <h3 className="text-base font-semibold">{e.label}</h3>
      <div className="my-2 space-y-1 rounded-lg border-l-2 border-primary/50 bg-card/40 px-4 py-3 font-mono text-[13px] leading-relaxed">
        {e.working.map((line, i) => (
          <div key={i}>{line}</div>
        ))}
      </div>
      {e.note && <p className="text-[15px] leading-relaxed text-muted-foreground">{e.note}</p>}
    </div>
  )
}

function ApiTable({ group }: { group: ApiGroup }) {
  const tone: Record<string, string> = {
    GET: 'text-emerald-500',
    POST: 'text-blue-500',
    PUT: 'text-amber-500',
    PATCH: 'text-amber-500',
    DELETE: 'text-rose-500',
  }
  return (
    <div className="my-5">
      <h3 className="text-base font-semibold">{group.title}</h3>
      <div className="mt-2 overflow-x-auto rounded-lg border border-border/50">
        <table className="w-full text-left text-sm">
          <thead className="bg-card/60">
            <tr>
              <th className="px-4 py-2.5 font-semibold">Method</th>
              <th className="px-4 py-2.5 font-semibold">Endpoint</th>
              <th className="px-4 py-2.5 font-semibold">What it does</th>
            </tr>
          </thead>
          <tbody>
            {group.rows.map((r, i) => (
              <tr key={i} className="border-t border-border/40">
                <td className={`px-4 py-2.5 align-top font-mono text-xs font-bold ${tone[r.method]}`}>
                  {r.method}
                </td>
                <td className="px-4 py-2.5 align-top font-mono text-xs">{r.path}</td>
                <td className="px-4 py-2.5 align-top text-muted-foreground">{r.what}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

const SECTION_NAMES: Record<string, string> = {
  problem: 'Problem statement',
  requirements: 'System requirements',
  capacity: 'Capacity estimation',
  'high-level-design': 'High level design',
  stack: 'Technology stack',
  'data-model': 'Data model',
  api: 'API design',
  'low-level-design': 'Low level design',
  scaling: 'Scalability and performance',
  bottlenecks: 'Bottlenecks and improvements',
}

/**
 * The interview questions, before anything else and outside the numbering.
 *
 * Collapsed by default with native <details>, so the page still opens on its
 * problem statement rather than on a wall of question text, and so this works
 * with no JavaScript at all.
 */
function InterviewQuestions({ s }: { s: SystemDesign }) {
  if (!s.interview) return null
  return (
    <section className="mt-10 scroll-mt-24" id="interview-questions">
      <h2 className="text-2xl font-bold tracking-tight">Questions this page answers</h2>
      {s.interview.intro && (
        <p className="mt-3 text-[15px] leading-relaxed text-muted-foreground">
          {s.interview.intro}
        </p>
      )}
      <div className="mt-4 divide-y divide-border/50 overflow-hidden rounded-xl border border-border/50">
        {s.interview.questions.map((item, i) => (
          <details key={i} className="group bg-card/20 open:bg-card/40">
            <summary className="flex cursor-pointer list-none items-start gap-3 px-4 py-3 text-[15px] font-medium hover:bg-card/60">
              <span
                aria-hidden
                className="mt-0.5 shrink-0 text-xs font-semibold tabular-nums text-primary"
              >
                {String(i + 1).padStart(2, '0')}
              </span>
              <span className="flex-1">{item.q}</span>
              <span
                aria-hidden
                className="mt-1 shrink-0 text-muted-foreground transition-transform group-open:rotate-90"
              >
                ›
              </span>
            </summary>
            <div className="px-4 pb-4 pl-11">
              <p className="text-[15px] leading-relaxed text-muted-foreground">{item.a}</p>
              {item.see && SECTION_NAMES[item.see] && (
                <a
                  href={`#${item.see}`}
                  className="mt-2 inline-block text-sm font-medium text-primary hover:underline"
                >
                  Worked through in {SECTION_NAMES[item.see]} &rarr;
                </a>
              )}
            </div>
          </details>
        ))}
      </div>
    </section>
  )
}

export function SystemDesignBody({ s }: { s: SystemDesign }) {
  let n = 0
  return (
    <>
      <InterviewQuestions s={s} />

      <Section n={++n} id="problem" title="Problem statement">
        {s.problem.text.map((p, i) => (
          <p key={i} className="mb-3 text-[15px] leading-relaxed text-muted-foreground">
            {p}
          </p>
        ))}
        <p className="mt-4 text-[15px] font-medium">The system has to be able to:</p>
        <Bullets items={s.problem.capabilities} />
      </Section>

      <Section n={++n} id="requirements" title="System requirements">
        <h3 className="text-base font-semibold">Functional requirements</h3>
        <Bullets items={s.functional} />
        <h3 className="mt-5 text-base font-semibold">Non-functional requirements</h3>
        <LabelledList items={s.nonFunctional} />
      </Section>

      {s.capacity && (
        <Section n={++n} id="capacity" title="Capacity estimation">
          <p className="text-[15px] leading-relaxed text-muted-foreground">
            Numbers for a mid-sized deployment. They are here to size the thing, not to predict
            your traffic: change an assumption and the sums below move with it.
          </p>
          <h3 className="mt-5 text-base font-semibold">Assumptions</h3>
          <TwoCol headers={['Parameter', 'Value']} rows={s.capacity.assumptions} />
          {s.capacity.estimates.map((e, i) => (
            <Working key={i} e={e} />
          ))}
        </Section>
      )}

      <Section n={++n} id="high-level-design" title="High level design">
        {s.highLevel.intro && (
          <p className="mb-3 text-[15px] leading-relaxed text-muted-foreground">
            {s.highLevel.intro}
          </p>
        )}
        <h3 className="text-base font-semibold">Core components</h3>
        <LabelledList items={s.highLevel.components} />

        <h3 className="mt-6 text-base font-semibold">Request flow</h3>
        <SystemFlow
          title={s.highLevel.flow.title}
          nodes={s.highLevel.flow.nodes}
          edges={s.highLevel.flow.edges}
        />
        <ol className="my-3 space-y-2">
          {s.highLevel.flow.steps.map((t, i) => (
            <li key={i} className="flex gap-3 text-[15px] leading-relaxed text-muted-foreground">
              <span className="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-primary text-[11px] font-bold text-primary-foreground tabular-nums">
                {i + 1}
              </span>
              <span>{t}</span>
            </li>
          ))}
        </ol>

        {s.highLevel.dataFlow && (
          <>
            <h3 className="mt-6 text-base font-semibold">Data flow</h3>
            <Bullets items={s.highLevel.dataFlow} />
          </>
        )}
      </Section>

      <Section n={++n} id="stack" title="Technology stack">
        <TwoCol headers={['Component', 'What it is']} rows={s.stack} />
      </Section>

      {s.dataModel && (
        <Section n={++n} id="data-model" title="Data model">
          {s.dataModel.intro && (
            <p className="mb-3 text-[15px] leading-relaxed text-muted-foreground">
              {s.dataModel.intro}
            </p>
          )}
          {s.dataModel.entities.map((e) => (
            <div key={e.name} className="my-5">
              <h3 className="font-mono text-base font-semibold">{e.name}</h3>
              {e.note && (
                <p className="mt-1 text-[15px] leading-relaxed text-muted-foreground">{e.note}</p>
              )}
              <TwoCol headers={['Column', 'Holds']} rows={e.fields} />
            </div>
          ))}
          {s.dataModel.storage && (
            <>
              <h3 className="mt-5 text-base font-semibold">Where it lives</h3>
              <Bullets items={s.dataModel.storage} />
            </>
          )}
        </Section>
      )}

      {s.api && (
        <Section n={++n} id="api" title="API design">
          {s.api.groups.map((g) => (
            <ApiTable key={g.title} group={g} />
          ))}
          {s.api.samples?.map((sample) => (
            <div key={sample.title} className="my-5">
              <h3 className="text-base font-semibold">{sample.title}</h3>
              <div className="mt-2">
                <CodeBlock language={sample.language} code={sample.code} />
              </div>
            </div>
          ))}
        </Section>
      )}

      {s.lowLevel && (
        <Section n={++n} id="low-level-design" title="Low level design">
          <h3 className="text-base font-semibold">Core types</h3>
          <div className="my-3 space-y-4">
            {s.lowLevel.classes.map((c) => (
              <div key={c.name} className="rounded-lg border border-border/50 bg-card/30 p-4">
                <div className="flex flex-wrap items-baseline gap-x-3 gap-y-1">
                  <span className="font-mono text-sm font-semibold">{c.name}</span>
                  {c.file && <span className={MONO}>{c.file}</span>}
                </div>
                <p className="mt-2 text-[15px] leading-relaxed text-muted-foreground">{c.what}</p>
                {c.methods && (
                  <div className="mt-2.5 flex flex-wrap gap-1.5">
                    {c.methods.map((m) => (
                      <code key={m} className={MONO}>
                        {m}
                      </code>
                    ))}
                  </div>
                )}
              </div>
            ))}
          </div>
          {s.lowLevel.principles && (
            <>
              <h3 className="mt-6 text-base font-semibold">Design principles applied</h3>
              <LabelledList items={s.lowLevel.principles} />
            </>
          )}
          {s.lowLevel.patterns && (
            <>
              <h3 className="mt-6 text-base font-semibold">Patterns</h3>
              <TwoCol headers={['Pattern', 'Where it is used']} rows={s.lowLevel.patterns} />
            </>
          )}
        </Section>
      )}

      <Section n={++n} id="scaling" title="Scalability and performance">
        <Bullets items={s.scaling} />
      </Section>

      <Section n={++n} id="bottlenecks" title="Bottlenecks and improvements">
        <h3 className="text-base font-semibold">What breaks first</h3>
        <LabelledList items={s.bottlenecks.problems} />
        <h3 className="mt-6 text-base font-semibold">What to do about it</h3>
        <LabelledList items={s.bottlenecks.improvements} />
      </Section>
    </>
  )
}
