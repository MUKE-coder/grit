'use client'

import { useState } from 'react'
import { DocsShell } from '@/components/docs/shell'
import { Button, ButtonArrow } from '@/components/grit/button'
import { Tabs, TabPanel } from '@/components/grit/tabs'
import { Terminal } from '@/components/grit/terminal'
import { Reveal } from '@/components/grit/reveal'
import { Highlight } from '@/components/grit/section'

// Every colour token, with the value beside it. A swatch without its value is
// a picture of a colour; the value is the thing somebody needs.
const COLOURS = [
  ['bg', '#08090B', '#FFFFFF', 'Page ground. The only full-bleed background.'],
  ['surface-1', '#0F1114', '#F7F7F8', 'Cards, sidebar, terminal body.'],
  ['surface-2', '#15181C', '#F0F0F2', 'Hover, inputs, tab rails, inline code.'],
  ['surface-3', '#1C2025', '#E8E8EB', 'Pressed, selected, active tab pill.'],
  ['border', '#23272D', '#E4E4E7', 'Every hairline. Borders, not shadows.'],
  ['border-strong', '#343A42', '#C8CAD0', 'Hover edges, inputs.'],
  ['fg', '#EDEEF0', '#0B0C0E', 'Headings, primary text, primary button fill.'],
  ['fg-muted', '#A1A7B0', '#555B64', 'Body copy. 7.8:1 dark, 6.3:1 light.'],
  ['fg-subtle', '#7D848E', '#6B717A', 'Captions and eyebrows. The floor.'],
  ['accent', '#A99BFF', '#5B47E0', 'Links, focus, active, one highlight per view.'],
  ['success', '#3DD68C', '#0B7A45', 'Terminal ticks, positive deltas.'],
  ['warning', '#F5A524', '#A35F00', 'Warning callouts. Sparingly.'],
  ['danger', '#FF6B6B', '#C93636', 'Destructive actions and errors.'],
] as const

const TYPE = [
  ['display-xl', '80 / 84', '-0.045em', 'Landing hero only'],
  ['display', '56 / 60', '-0.04em', 'Landing section titles'],
  ['h1', '40 / 48', '-0.035em', 'Docs page title'],
  ['h2', '28 / 36', '-0.025em', 'Docs sections'],
  ['h3', '20 / 28', '-0.015em', 'Subsections, card titles'],
  ['lead', '20 / 32', '-0.01em', 'Hero and docs intro'],
  ['body', '16 / 28', '—', 'Prose, max 72ch'],
  ['body-sm', '14 / 22', '—', 'UI, cards, sidebar'],
  ['caption', '13 / 20', '—', 'Footnotes, figure captions'],
  ['eyebrow', '12 / 16 mono', '+0.08em', 'Section labels, uppercase'],
  ['code', '13.5 / 24 mono', '—', 'Code blocks and terminals'],
] as const

const RADII = [
  ['radius-xs', 4, 'Inline code, kbd, tiny badges'],
  ['radius-sm', 6, 'Sidebar items, tab pills, inputs'],
  ['radius-md', 10, 'Buttons, code blocks, callouts'],
  ['radius-lg', 16, 'Cards and terminals'],
  ['radius-xl', 24, 'Large showcase panels'],
] as const

export default function DesignPage() {
  const [tab, setTab] = useState('segmented')

  return (
    <DocsShell>
          <p className="g-eyebrow mb-4">{'{ DESIGN SYSTEM }'}</p>
          <h1 className="text-[40px] font-semibold leading-[48px] tracking-[-0.035em] text-fg">
            Every component, <Highlight>beside its value</Highlight>
          </h1>
          <p className="mt-4 text-lead text-fg-muted">
            The reference the design review runs against. If something here does not match the
            system, the system is what is right.
          </p>

          {/* ------------------------------------------------------ colour */}
          <Reveal>
            <h2 className="mt-16 text-[28px] font-semibold leading-9 tracking-[-0.025em] text-fg">
              Colour
            </h2>
            <p className="mt-2 text-body-sm text-fg-muted">
              Neutral graphite plus one violet. The accent is rationed: at most three accent
              fills in a viewport.
            </p>
            <div className="mt-6 overflow-hidden rounded-lg border border-border">
              <table className="w-full text-body-sm">
                <thead>
                  <tr className="border-b border-border bg-surface-1 text-left">
                    <th className="px-4 py-2.5 font-medium text-fg">Token</th>
                    <th className="px-4 py-2.5 font-medium text-fg">Dark</th>
                    <th className="px-4 py-2.5 font-medium text-fg">Light</th>
                    <th className="px-4 py-2.5 font-medium text-fg">Use</th>
                  </tr>
                </thead>
                <tbody>
                  {COLOURS.map(([name, dark, light, use]) => (
                    <tr key={name} className="border-b border-border/60 last:border-0">
                      <td className="px-4 py-2.5">
                        <span className="flex items-center gap-2.5">
                          <span
                            aria-hidden="true"
                            className="h-4 w-4 shrink-0 rounded-[4px] border border-border"
                            style={{ background: `var(--${name})` }}
                          />
                          <code className="font-mono text-[13px] text-fg">{name}</code>
                        </span>
                      </td>
                      <td className="px-4 py-2.5 font-mono text-[12px] text-fg-subtle">{dark}</td>
                      <td className="px-4 py-2.5 font-mono text-[12px] text-fg-subtle">{light}</td>
                      <td className="px-4 py-2.5 text-fg-muted">{use}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </Reveal>

          {/* -------------------------------------------------------- type */}
          <Reveal>
            <h2 className="mt-16 text-[28px] font-semibold leading-9 tracking-[-0.025em] text-fg">
              Type
            </h2>
            <p className="mt-2 text-body-sm text-fg-muted">
              Geist and Geist Mono. Eleven styles. Hierarchy is size and weight, never colour.
            </p>
            <div className="mt-6 overflow-hidden rounded-lg border border-border">
              <table className="w-full text-body-sm">
                <thead>
                  <tr className="border-b border-border bg-surface-1 text-left">
                    <th className="px-4 py-2.5 font-medium text-fg">Style</th>
                    <th className="px-4 py-2.5 font-medium text-fg">Size / line</th>
                    <th className="px-4 py-2.5 font-medium text-fg">Tracking</th>
                    <th className="px-4 py-2.5 font-medium text-fg">Use</th>
                  </tr>
                </thead>
                <tbody>
                  {TYPE.map(([name, size, tracking, use]) => (
                    <tr key={name} className="border-b border-border/60 last:border-0">
                      <td className="px-4 py-2.5">
                        <code className="font-mono text-[13px] text-accent">{name}</code>
                      </td>
                      <td className="px-4 py-2.5 font-mono text-[12px] text-fg-subtle">{size}</td>
                      <td className="px-4 py-2.5 font-mono text-[12px] text-fg-subtle">
                        {tracking}
                      </td>
                      <td className="px-4 py-2.5 text-fg-muted">{use}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </Reveal>

          {/* ----------------------------------------------------- buttons */}
          <Reveal>
            <h2 className="mt-16 text-[28px] font-semibold leading-9 tracking-[-0.025em] text-fg">
              Button
            </h2>
            <p className="mt-2 text-body-sm text-fg-muted">
              Primary is inverted, not violet. That is what lets violet mean brand everywhere
              else. One primary per view.
            </p>
            <div className="mt-6 flex flex-wrap items-center gap-3 rounded-lg border border-border bg-surface-1 p-6">
              <Button variant="primary">
                Get started
                <ButtonArrow />
              </Button>
              <Button variant="secondary">Why Grit</Button>
              <Button variant="ghost">Read the docs</Button>
              <Button variant="accent">Sponsor</Button>
            </div>
            <div className="mt-3 flex flex-wrap items-center gap-3 rounded-lg border border-border bg-surface-1 p-6">
              <Button variant="primary" size="lg">
                Large
              </Button>
              <Button variant="secondary">Medium</Button>
              <Button variant="secondary" size="sm">
                Small
              </Button>
            </div>
          </Reveal>

          {/* -------------------------------------------------------- tabs */}
          <Reveal>
            <h2 className="mt-16 text-[28px] font-semibold leading-9 tracking-[-0.025em] text-fg">
              Tabs
            </h2>
            <p className="mt-2 text-body-sm text-fg-muted">
              The indicator is measured from the selected tab, not animated by index: the labels
              are not equal width. Arrow keys move selection, and focus moves with it.
            </p>
            <div className="mt-6 rounded-lg border border-border bg-surface-1 p-6">
              <Tabs
                label="Component examples"
                value={tab}
                onChange={setTab}
                items={[
                  { id: 'segmented', label: 'Segmented' },
                  { id: 'underline', label: 'Underline' },
                  { id: 'keyboard', label: 'Keyboard' },
                ]}
              />
              <div className="mt-5">
                <TabPanel id="segmented" active={tab === 'segmented'}>
                  <p className="text-body-sm text-fg-muted">
                    A pill on surface-3 with an inset hairline, sliding over 240ms.
                  </p>
                </TabPanel>
                <TabPanel id="underline" active={tab === 'underline'}>
                  <Tabs
                    label="Underline variant"
                    variant="line"
                    value="a"
                    onChange={() => {}}
                    items={[
                      { id: 'a', label: 'macOS / Linux' },
                      { id: 'b', label: 'Windows' },
                      { id: 'c', label: 'Go install' },
                    ]}
                  />
                </TabPanel>
                <TabPanel id="keyboard" active={tab === 'keyboard'}>
                  <p className="text-body-sm text-fg-muted">
                    Tab reaches the tablist once, then{' '}
                    <kbd className="g-kbd">←</kbd> <kbd className="g-kbd">→</kbd> move between
                    tabs and <kbd className="g-kbd">Home</kbd> / <kbd className="g-kbd">End</kbd>{' '}
                    jump to the ends.
                  </p>
                </TabPanel>
              </div>
            </div>
          </Reveal>

          {/* ---------------------------------------------------- terminal */}
          <Reveal>
            <h2 className="mt-16 text-[28px] font-semibold leading-9 tracking-[-0.025em] text-fg">
              Terminal
            </h2>
            <p className="mt-2 text-body-sm text-fg-muted">
              Dark in both themes, so a screenshot of one is stable. Copy confirms in place for
              1.6s rather than raising a toast.
            </p>
            <div className="mt-6">
              <Terminal
                title="~/contacts"
                lines={[
                  { run: 'grit new contacts --triple --next' },
                  { ok: 'API, web app and admin panel, with auth already working' },
                  { run: 'cd contacts && pnpm install' },
                  { ok: 'The workspace, once' },
                  { run: 'grit start' },
                  { out: 'Opens localhost:3000 for the app, :3001 for the admin' },
                ]}
              />
            </div>
          </Reveal>

          {/* ---------------------------------------------------- callouts */}
          <Reveal>
            <h2 className="mt-16 text-[28px] font-semibold leading-9 tracking-[-0.025em] text-fg">
              Callout
            </h2>
            <div className="mt-6 space-y-3">
              <div className="g-callout g-callout--tip">
                <svg viewBox="0 0 20 20" fill="none" stroke="currentColor" strokeWidth="1.5" aria-hidden="true">
                  <circle cx="10" cy="10" r="7.5" />
                  <path d="M10 13.5v-4M10 7h.01" strokeLinecap="round" />
                </svg>
                <p>
                  <strong>Tip.</strong> The accent-soft fill and accent-line border are the only
                  place a callout carries the brand colour.
                </p>
              </div>
              <div className="g-callout g-callout--warn">
                <svg viewBox="0 0 20 20" fill="none" stroke="currentColor" strokeWidth="1.5" aria-hidden="true">
                  <path d="M10 3 2.5 16.5h15L10 3Z" strokeLinejoin="round" />
                  <path d="M10 8v3.5M10 14h.01" strokeLinecap="round" />
                </svg>
                <p>
                  <strong>Careful.</strong> Status colours appear in terminals, callouts and
                  table cells. Never in decoration.
                </p>
              </div>
            </div>
          </Reveal>

          {/* ------------------------------------------------------ shape */}
          <Reveal>
            <h2 className="mt-16 text-[28px] font-semibold leading-9 tracking-[-0.025em] text-fg">
              Shape
            </h2>
            <p className="mt-2 text-body-sm text-fg-muted">
              Radius scales with the element. Nested radius is the outer radius minus the
              padding.
            </p>
            <div className="mt-6 flex flex-wrap gap-4">
              {RADII.map(([name, px, use]) => (
                <div key={name} className="flex-1 basis-40">
                  <div
                    className="flex h-20 items-center justify-center border border-border bg-surface-1 font-mono text-[12px] text-fg-subtle"
                    style={{ borderRadius: px }}
                  >
                    {px}px
                  </div>
                  <p className="mt-2 font-mono text-[12px] text-accent">{name}</p>
                  <p className="mt-0.5 text-caption text-fg-subtle">{use}</p>
                </div>
              ))}
            </div>
          </Reveal>

          {/* -------------------------------------------------- atmosphere */}
          <Reveal>
            <h2 className="mt-16 text-[28px] font-semibold leading-9 tracking-[-0.025em] text-fg">
              Atmosphere
            </h2>
            <p className="mt-2 text-body-sm text-fg-muted">
              Structure that is drawn, not implied. Corner marks are the signature, and what
              makes them a signature is that every section has them.
            </p>
            <div className="mt-6 grid gap-4 sm:grid-cols-2">
              <div className="relative rounded-lg border border-border bg-surface-1 p-6">
                <span className="g-cnr g-cnr--tl" aria-hidden="true" />
                <span className="g-cnr g-cnr--tr" aria-hidden="true" />
                <p className="font-mono text-[12px] text-accent">g-cnr</p>
                <p className="mt-1 text-caption text-fg-subtle">
                  7×7 accent square where a rail meets a rule.
                </p>
              </div>
              <div className="g-dots relative rounded-lg border border-border p-6">
                <p className="font-mono text-[12px] text-accent">g-dots</p>
                <p className="mt-1 text-caption text-fg-subtle">
                  16px dot field, behind diagrams only.
                </p>
              </div>
              <div className="g-gridbg relative rounded-lg border border-border p-6">
                <p className="font-mono text-[12px] text-accent">g-gridbg</p>
                <p className="mt-1 text-caption text-fg-subtle">
                  56px line grid, radially masked. Hero and CTA only.
                </p>
              </div>
              <div className="relative overflow-hidden rounded-lg border border-border p-6">
                <div className="g-wash absolute -left-20 -top-20 h-60 w-60" aria-hidden="true" />
                <p className="relative font-mono text-[12px] text-accent">g-wash</p>
                <p className="relative mt-1 text-caption text-fg-subtle">
                  One radial violet glow per section, at most.
                </p>
              </div>
            </div>
          </Reveal>

          <p className="mt-16 border-t border-border pt-6 text-caption text-fg-subtle">
            Run the audit against this page before changing it:{' '}
            <code className="g-code-inline">
              eval(await (await fetch(&apos;/design-audit.js&apos;)).text())
            </code>
          </p>
    </DocsShell>
  )
}
