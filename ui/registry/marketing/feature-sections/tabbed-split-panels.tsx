'use client'

import { useEffect, useState } from 'react'
import { Bot, Brain, ChevronRight, Hourglass, Plus, ShieldCheck } from 'lucide-react'

/**
 * Tabs at the top switch the entire split panel — headline, copy and mock.
 *
 * Unlike the accordion variant, here one control owns everything, so the tab is
 * the section's single piece of state. That is the right trade when each tab is a
 * genuinely different story rather than a detail of a shared one.
 *
 * The panel is always mounted and only its content swaps, so the section height
 * does not jump between tabs. A section that resizes under the pointer makes the
 * next tab move away from the cursor mid-click.
 */

const COMPLIANCE = [
  { icon: ShieldCheck, label: 'SOC 2' },
  { icon: ShieldCheck, label: 'ISO 27001' },
  { icon: ShieldCheck, label: 'GDPR' },
  { icon: Hourglass, label: '99.9%', muted: 'uptime' },
]

/**
 * A drawing of a compose window. aria-hidden because it is scenery: the
 * addresses and field labels are invented, and read aloud they arrive as if
 * this page were an inbox. Hiding it also settles the contrast question — the
 * faint field labels are decoration, not text someone has to read.
 */
function ComposeMock() {
  return (
    <div
      aria-hidden="true"
      className="w-full max-w-[23rem] rounded-2xl bg-white/95 p-4 shadow-[0_24px_60px_-20px_rgb(15_23_42_/_0.45)] backdrop-blur-sm dark:bg-gray-900/95">
      <div className="flex items-center justify-between gap-3 border-b border-gray-200 pb-3 dark:border-white/10">
        <span className="flex items-center gap-1.5 text-[13px] text-gray-500 dark:text-gray-400">
          To:
          <span className="inline-flex items-center gap-1.5 rounded-full bg-gray-100 py-0.5 pr-2.5 pl-1 dark:bg-white/10">
            <span
              aria-hidden
              className="flex size-5 items-center justify-center rounded-full bg-gradient-to-br from-violet-400 to-sky-400 text-[9px] font-bold text-white"
            >
              SH
            </span>
            <span className="text-[12.5px] font-medium text-gray-900 dark:text-white">
              Shadcn
            </span>
          </span>
        </span>
        <span
          aria-hidden
          className="flex size-7 shrink-0 items-center justify-center rounded-full bg-gray-100 text-gray-600 dark:bg-white/10 dark:text-gray-300"
        >
          <Plus className="size-3.5" />
        </span>
      </div>

      {['Cc:', 'Subject:', 'From:'].map((label) => (
        <p
          key={label}
          className="border-b border-gray-100 py-2.5 text-[13px] text-gray-400 dark:border-white/[0.06] dark:text-gray-500"
        >
          {label}
        </p>
      ))}

      <p className="mt-4 text-[13.5px]/[1.6] text-gray-700 dark:text-gray-300">
        Web applications with{' '}
        <span className="text-sky-600 dark:text-sky-400">React and TypeScript</span> using
        best practices.
      </p>
      <p className="mt-3 text-[13.5px] text-gray-500 dark:text-gray-400">
        Sent from my iPhone
      </p>
    </div>
  )
}

function AgentMock() {
  const steps = [
    { label: 'Read the design brief', state: 'done' },
    { label: 'Draft three layout options', state: 'done' },
    { label: 'Render previews', state: 'running' },
    { label: 'Open a pull request', state: 'todo' },
  ] as const

  return (
    <div className="w-full max-w-[23rem] rounded-2xl bg-white/95 p-5 shadow-[0_24px_60px_-20px_rgb(15_23_42_/_0.45)] backdrop-blur-sm dark:bg-gray-900/95">
      <div className="flex items-center gap-2.5">
        <span
          aria-hidden
          className="flex size-8 items-center justify-center rounded-xl bg-gray-900 text-white dark:bg-white dark:text-gray-900"
        >
          <Bot className="size-4" />
        </span>
        <div>
          <p className="text-[14px] font-semibold text-gray-900 dark:text-white">
            Smart Agent
          </p>
          <p className="text-[12px] text-gray-500 dark:text-gray-400">4 steps · 2 done</p>
        </div>
      </div>

      <ol className="mt-4 space-y-3">
        {steps.map((step) => (
          <li key={step.label} className="flex items-center gap-3">
            <span
              className={`flex size-4 shrink-0 items-center justify-center rounded-full text-[9px] font-bold ${
                step.state === 'done'
                  ? 'bg-emerald-500 text-white'
                  : step.state === 'running'
                    ? 'bg-sky-500 text-white'
                    : 'border border-gray-300 dark:border-white/20'
              }`}
            >
              {step.state === 'done' ? '✓' : step.state === 'running' ? '·' : ''}
            </span>
            <span
              className={`text-[13.5px] ${
                step.state === 'todo'
                  ? 'text-gray-400 dark:text-gray-500'
                  : 'font-medium text-gray-900 dark:text-white'
              }`}
            >
              {step.label}
            </span>
          </li>
        ))}
      </ol>

      <div className="mt-4 rounded-xl bg-gray-50 p-3 dark:bg-white/5">
        <p className="font-mono text-[11.5px] text-gray-500 dark:text-gray-400">
          render_preview(variant=&quot;split&quot;)
        </p>
      </div>
    </div>
  )
}

const TABS = [
  {
    key: 'models',
    label: 'AI Models',
    icon: Brain,
    title: 'Accelerate development with AI-powered assistance',
    body: 'Switch between GPT, Claude, and Gemini from a unified interface with seamless provider switching.',
    canvas:
      'bg-[linear-gradient(180deg,#a7c6e0_0%,#bcd0e2_26%,#d8cdb4_56%,#b9924f_80%,#8e6b34_100%)] dark:bg-[linear-gradient(180deg,#16273a_0%,#182839_26%,#2a2617_56%,#33240f_80%,#231708_100%)]',
    mock: <ComposeMock />,
  },
  {
    key: 'agent',
    label: 'Smart Agent',
    icon: Bot,
    title: 'Delegate the work you would rather not do by hand',
    body: 'Give the agent a goal and watch every tool call as it happens, with the option to stop it at any step.',
    canvas:
      'bg-[linear-gradient(180deg,#b6c7d8_0%,#c4cfd6_26%,#c2c9bd_56%,#7f9a7c_80%,#4f6b4e_100%)] dark:bg-[linear-gradient(180deg,#1a2530_0%,#1b2730_26%,#1d2823_56%,#1a2a1b_80%,#0f1c10_100%)]',
    mock: <AgentMock />,
  },
] as const

/**
 * Crossfades the panel copy when `token` changes.
 *
 * The content is swapped in place rather than remounted, so a plain `transition`
 * class has nothing to animate from — the new text is already at full opacity by
 * the time it paints. Dropping to 0 and flipping back on the next frame gives the
 * browser one paint at the start value, which is what makes the transition run.
 */
function useCrossfade(token: string) {
  const [shown, setShown] = useState(true)
  useEffect(() => {
    setShown(false)
    const id = requestAnimationFrame(() => setShown(true))
    return () => cancelAnimationFrame(id)
  }, [token])
  return shown
}

export default function TabbedSplitPanels({
  ctaLabel = 'Learn more',
}: {
  ctaLabel?: string
}) {
  const [active, setActive] = useState<string>(TABS[0].key)
  const current = TABS.find((t) => t.key === active) ?? TABS[0]
  const shown = useCrossfade(active)

  return (
    <section className="bg-white dark:bg-gray-950">
      {/* Dashed rules top and bottom of the tab strip, so the tabs sit in a band
          rather than floating above the content they control. */}
      <div className="border-b border-dashed border-gray-200 dark:border-white/10">
        <div
          role="tablist"
          aria-label="Feature areas"
          className="mx-auto flex max-w-7xl gap-8 px-6"
        >
          {TABS.map((tab) => {
            const isActive = active === tab.key
            return (
              <button
                key={tab.key}
                type="button"
                role="tab"
                aria-selected={isActive}
                onClick={() => setActive(tab.key)}
                className={`-mb-px inline-flex items-center gap-2 border-b-2 py-5 text-[15px] font-semibold transition-colors duration-200 ${
                  isActive
                    ? 'border-gray-900 text-gray-900 dark:border-white dark:text-white'
                    : 'border-transparent text-gray-500 hover:text-gray-900 dark:text-gray-400 dark:hover:text-gray-100'
                }`}
              >
                <tab.icon aria-hidden="true" className="size-[18px]" />
                {tab.label}
              </button>
            )
          })}
        </div>
      </div>

      <div
        role="tabpanel"
        className="mx-auto grid max-w-7xl items-start gap-14 px-6 py-16 lg:grid-cols-2 lg:gap-20 lg:py-20"
      >
        <div>
          <div
            className={`transition-all duration-300 ease-out ${
              shown ? 'translate-y-0 opacity-100' : 'translate-y-1 opacity-0'
            }`}
          >
            <h2 className="text-[2rem] leading-[1.12] font-semibold tracking-[-0.025em] text-gray-900 sm:text-[2.6rem] dark:text-white">
              {current.title}
            </h2>
            <p className="mt-5 max-w-md text-[16.5px]/[1.6] text-gray-500 dark:text-gray-400">
              {current.body}
            </p>
          </div>

          <a
            href="#"
            className="group mt-8 inline-flex h-11 items-center gap-1.5 rounded-xl border border-gray-200 bg-white px-5 text-[14.5px] font-semibold text-gray-900 shadow-[0_1px_2px_rgb(15_23_42_/_0.06)] transition-all duration-200 hover:bg-gray-50 active:scale-[0.98] dark:border-white/15 dark:bg-white/5 dark:text-white dark:hover:bg-white/10"
          >
            {ctaLabel}
            <ChevronRight
              aria-hidden="true"
              className="size-4 transition-transform duration-200 group-hover:translate-x-0.5"
            />
          </a>

          <ul className="mt-16 space-y-3">
            {COMPLIANCE.map((item) => (
              <li key={item.label} className="flex items-center gap-3 text-[15px]">
                <item.icon aria-hidden="true" className="size-[18px] text-gray-400" />
                <span className="font-semibold text-gray-900 dark:text-white">
                  {item.label}
                </span>
                {item.muted && (
                  <span className="text-gray-500 dark:text-gray-400">{item.muted}</span>
                )}
              </li>
            ))}
          </ul>
        </div>

        <div
          className={`relative aspect-[4/4.2] overflow-hidden rounded-3xl transition-colors duration-500 ${current.canvas}`}
        >
          <div
            aria-hidden
            className="pointer-events-none absolute inset-0"
            style={{
              background:
                'radial-gradient(54% 30% at 22% 97%, rgb(40 32 14 / 0.4), transparent 70%), radial-gradient(46% 26% at 82% 99%, rgb(30 26 12 / 0.35), transparent 70%)',
            }}
          />
          <div className="absolute inset-0 flex items-center justify-center p-6">
            {current.mock}
          </div>
        </div>
      </div>
    </section>
  )
}
