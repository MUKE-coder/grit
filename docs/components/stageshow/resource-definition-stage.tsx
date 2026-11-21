'use client'

import { useTheme } from 'next-themes'

import { Bar, Collapse, Stage, Text, type Act } from './kit'
import './resource-definition-stage.css'

// One file describes the whole screen.
//
// The section's claim, as a picture: the resource definition above the admin
// screen it produces, and a column added to one appearing in the other. A
// screenshot cannot show that, because the subject is the relationship rather
// than either surface.
//
// Copied from the skill's stacked-cards gold: two components at one width, a
// small gap, the lower one off the bottom. The first attempt was
// panel-annotation, the panel-over-product gold, and it does not survive a
// table: a satellite straddling the window's right edge lands on the columns,
// and the rule is that a satellite covers chrome. Stacking also puts the two
// in the order the sentence reads, file then screen.
//
// No leader. The gap is 12px, which is no path for one, and the gold's own
// note is that where there is no clear path the shared change carries the
// story: the line and the column arrive together and share the accent.
//
// Honest by construction: adding an entry to a resource definition's columns
// array is what puts a column in the generated table. This illustrates a real
// mechanism rather than inventing a product.
//
// Real strings: tasks.ts, Task, the four column names, Tasks, New task,
// Priority, and the two rows the eye lands on. Everything else is a bar.

type S = { added: boolean }

const initial: S = { added: false }
// The payoff, because most people see only this frame: the column is in the
// file and in the table, and nothing is caught mid-change.
const poster: S = { added: true }

const ROWS = [
  { title: 'Kitchen extractor', due: '4 Mar', priority: 'High' },
  { title: 'Pearson loft quote', due: '6 Mar', priority: 'Medium' },
  { title: null, due: null, priority: 'Low' },
  { title: null, due: null, priority: 'Medium' },
  { title: null, due: null, priority: 'Low' },
]

async function script(a: Act<S>) {
  await a.wait(1100)

  // Where somebody would add it: at the end of the columns array, under the
  // three already there.
  a.mark('add')
  await a.move({ at: 'columns', x: 0.42, y: 0.9 })
  a.put({ added: true })

  // Long enough to read the column that arrived, which is the whole point of
  // the frame.
  await a.wait(3200)
}

export function ResourceDefinitionStage() {
  // The stage takes the site's theme and accent rather than choosing its
  // own. A white demo on the dark landing page is a lamp in the middle of
  // the section, and a dark one in light mode is the same mistake the
  // other way; the kit's look is a set of tokens for exactly this.
  const { resolvedTheme } = useTheme()
  const dark = resolvedTheme !== 'light'

  return (
    <Stage
      width={620}
      height={500}
      compact={{ below: 580, width: 358, height: 580 }}
      initial={initial}
      poster={poster}
      script={script}
      play="once"
      cursor={{ at: { x: 566, y: 462 } }}
      theme={dark ? 'dark' : 'light'}
      accent={dark ? '#a99bff' : '#5b47e0'}
      className="rd"
    >
      {(s) => (
        <>
          {/* The file */}
          <section className="rd-file">
            <header className="rd-file__head">
              <span className="rd-file__dot" />
              <Text className="rd-file__name">tasks.ts</Text>
            </header>

            <div className="rd-code">
            {/* A code line is a plain span, not the kit's Text, wherever it
                carries <b> or <i> for syntax colour. Text measures
                children.length to size its filler bar in minimal fidelity,
                and the length of an array of three nodes is 3, so those lines
                came out as three-character bars. At the fidelity this stage
                runs at, Text renders exactly this span, so nothing moved. The
                lines whose whole content is one string still use Text. */}
              <span className="ss-text rd-ln"><i>export default</i> defineResource({'{'}</span>
              <span className="ss-text rd-ln rd-ln--in">name: <b>&apos;Task&apos;</b>,</span>
              <Text className="rd-ln rd-ln--in">columns: [</Text>
              <div data-target="columns">
                <span className="ss-text rd-ln rd-ln--in2">{'{'} name: <b>&apos;title&apos;</b>, label: <b>&apos;Title&apos;</b> {'}'},</span>
                <span className="ss-text rd-ln rd-ln--in2">{'{'} name: <b>&apos;status&apos;</b>, label: <b>&apos;Status&apos;</b><span className="rd-wide-only">, filter: <i>true</i></span> {'}'},</span>
                <span className="ss-text rd-ln rd-ln--in2">{'{'} name: <b>&apos;due_at&apos;</b>, label: <b>&apos;Due&apos;</b> {'}'},</span>
                <Collapse open={s.added} reveal="clip">
                  <span className="ss-text rd-ln rd-ln--in2">{'{'} name: <b>&apos;priority&apos;</b>, label: <b>&apos;Priority&apos;</b> {'}'},</span>
                </Collapse>
              </div>
              <Text className="rd-ln rd-ln--in">],</Text>
              <span className="ss-text rd-ln">{'}'})</span>
            </div>
          </section>

          {/* The screen it describes */}
          <section className="rd-screen">
            <header className="rd-screen__head">
              <Text className="rd-title">Tasks</Text>
              <span className="rd-new">New task</span>
            </header>

            <div className="rd-table">
              <div className="rd-tr rd-tr--head">
                <Text className="rd-th">Title</Text>
                <Text className="rd-th">Status</Text>
                <Text className="rd-th rd-due">Due</Text>
                <Collapse open={s.added} reveal="clip">
                  <Text className="rd-th rd-th--new">Priority</Text>
                </Collapse>
              </div>

              {ROWS.map((row, i) => (
                <div className="rd-tr" key={i}>
                  {row.title
                    ? <Text className="rd-cell">{row.title}</Text>
                    : <Bar w={i === 3 ? 150 : 124} h={8} r={4} tone={2} />}
                  <span className="rd-pill" data-tone={i % 3} />
                  {row.due
                    ? <Text className="rd-cell rd-cell--meta rd-due">{row.due}</Text>
                    : <Bar w={38} h={8} r={4} tone={3} className="rd-due" />}
                  <Collapse open={s.added} reveal="clip">
                    <Text className="rd-cell rd-cell--pri">{row.priority}</Text>
                  </Collapse>
                </div>
              ))}
            </div>
          </section>
        </>
      )}
    </Stage>
  )
}
