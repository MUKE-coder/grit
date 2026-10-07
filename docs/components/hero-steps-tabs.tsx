'use client'

import { useState } from 'react'
import { Check, Copy, Globe, Layers, Monitor, Smartphone, Server } from 'lucide-react'

/*
 * The hero: how a full-stack app gets built, in the commands that build it.
 *
 * What was here was a code block of a generated handler, which answers "what
 * does the output look like" for somebody already persuaded. The question a
 * first visitor has is "what do I type", and four tabs answer it for the four
 * shapes Grit scaffolds, all building the same thing: contacts, in groups, each
 * with a photo.
 *
 * Every command on this page is real and in this order. The whole point of the
 * hero is that somebody pastes it: a command that does not exist is worse than
 * no command, because the first person to try it finds out and then doubts
 * everything else on the site.
 */

interface Step {
  command: string
  result: string
}

interface Variant {
  key: string
  label: string
  icon: typeof Globe
  blurb: string
  steps: Step[]
  /** What is running when the last step finishes. */
  opens: string
}

// The two resources, identical in every variant, because the point is that the
// same four commands build the same app whatever you are shipping it as.
const GROUP = `grit generate resource Group \\
  --fields "name:string,description:text"`

const CONTACT = `grit generate resource Contact --fields \\
  "name:string,email:email,phone:tel,\\
   photo:file:image,group:belongs_to:Group"`

const VARIANTS: Variant[] = [
  {
    key: 'web',
    label: 'Web app',
    icon: Globe,
    blurb: 'Next.js front end, admin panel and Go API.',
    steps: [
      { command: 'grit new contacts --triple --next', result: 'API, web app and admin panel, with auth already working' },
      { command: 'cd contacts', result: '' },
      { command: 'pnpm install', result: 'The workspace, once' },
      { command: GROUP, result: 'Model, migration, API, types and an admin screen' },
      { command: CONTACT, result: 'A photo upload and a group picker, wired to the Group above' },
      { command: 'grit start', result: 'Everything, in parallel' },
    ],
    opens: 'localhost:3000 for the app, :3001 for the admin panel',
  },
  {
    key: 'desktop',
    label: 'Desktop app',
    icon: Monitor,
    blurb: 'The same API, in a Wails window, offline-capable.',
    steps: [
      { command: 'grit new contacts --triple --desktop', result: 'Adds a Wails app that shares the monorepo API' },
      { command: 'cd contacts', result: '' },
      { command: 'pnpm install', result: 'The workspace, once' },
      { command: GROUP, result: 'Model, migration, API, types and an admin screen' },
      { command: CONTACT, result: 'A photo upload and a group picker, wired to the Group above' },
      { command: 'grit start desktop', result: 'wails dev, with hot reload' },
    ],
    opens: 'a native window; grit package builds the installer',
  },
  {
    key: 'mobile',
    label: 'Mobile app',
    icon: Smartphone,
    blurb: 'Expo and React Native, against the same API.',
    steps: [
      { command: 'grit new contacts --triple --expo', result: 'Adds apps/expo, sharing the Zod schemas and types' },
      { command: 'cd contacts', result: '' },
      { command: 'pnpm install', result: 'The workspace, once' },
      { command: GROUP, result: 'Model, migration, API, types and an admin screen' },
      { command: CONTACT, result: 'A photo upload and a group picker, wired to the Group above' },
      { command: 'grit start expo', result: 'Metro, for a simulator or your phone' },
    ],
    opens: 'Expo Go, or a development build',
  },
  {
    key: 'api',
    label: 'API only',
    icon: Server,
    blurb: 'Go, Gin and GORM. No front end at all.',
    steps: [
      { command: 'grit new contacts --api', result: 'The Go API alone, with auth, docs and migrations' },
      { command: 'cd contacts', result: '' },
      { command: GROUP, result: 'Model, migration, service, handler and routes' },
      { command: CONTACT, result: 'An upload endpoint and the foreign key to Group' },
      { command: 'grit start', result: 'The API' },
    ],
    opens: 'localhost:8080, with the reference at /docs',
  },
  {
    key: 'full',
    label: 'Everything',
    icon: Layers,
    blurb: 'API, web, admin, desktop, mobile and a docs site, in one repo.',
    steps: [
      { command: 'grit new contacts --full', result: 'Every app Grit can scaffold, sharing one API and one set of types' },
      { command: 'cd contacts', result: '' },
      { command: 'pnpm install', result: 'The workspace, once' },
      { command: GROUP, result: 'Model, migration, API, types and an admin screen' },
      { command: CONTACT, result: 'A photo upload and a group picker, wired to the Group above' },
      { command: 'grit start', result: 'API, web and admin together; add expo or desktop for those' },
    ],
    opens: 'all of them at once, against the same database',
  },
]

export function HeroStepsTabs() {
  const [active, setActive] = useState(VARIANTS[0].key)
  const [copied, setCopied] = useState<number | null>(null)
  const variant = VARIANTS.find((v) => v.key === active) ?? VARIANTS[0]

  async function copy(text: string, index: number) {
    try {
      await navigator.clipboard.writeText(text.replace(/\\\n\s*/g, ' '))
      setCopied(index)
      setTimeout(() => setCopied(null), 1600)
    } catch {
      // A browser that refuses the clipboard still shows the command to read.
    }
  }

  return (
    <div className="relative rounded-2xl overflow-hidden bg-white dark:bg-[#0d1117] border border-border shadow-[0_24px_64px_-16px_rgba(2,6,23,0.5)]">
      <div
        role="tablist"
        aria-label="What you are building"
        className="flex items-center gap-0 overflow-x-auto bg-[#f6f8fa] dark:bg-[#161b22] border-b border-[#d0d7de] dark:border-white/[0.08] [-ms-overflow-style:none] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
      >
        {VARIANTS.map((v) => {
          const selected = v.key === active
          return (
            <button
              key={v.key}
              type="button"
              role="tab"
              aria-selected={selected}
              onClick={() => setActive(v.key)}
              className={`inline-flex shrink-0 items-center gap-1.5 border-b-2 px-3.5 py-2.5 text-[12px] font-medium whitespace-nowrap transition-colors ${
                selected
                  ? 'border-primary bg-white text-[#24292f] dark:bg-[#0d1117] dark:text-slate-100'
                  : 'border-transparent text-[#57606a] hover:text-[#24292f] dark:text-slate-500 dark:hover:text-slate-300'
              }`}
            >
              <v.icon className="h-3.5 w-3.5" />
              {v.label}
            </button>
          )
        })}
      </div>

      <div className="flex items-baseline justify-between gap-3 px-4 py-2 bg-white dark:bg-[#0d1117] border-b border-[#d0d7de]/60 dark:border-white/[0.06]">
        <span className="text-[11px] text-[#57606a] dark:text-slate-500">{variant.blurb}</span>
        <span className="text-[11px] font-mono text-[#57606a] dark:text-slate-500">contacts</span>
      </div>

      {/* A fixed height, so switching tabs does not resize the hero and shove
          the page under the reader's cursor. */}
      <div className="h-[430px] overflow-y-auto px-4 py-3.5">
        <ol className="space-y-2.5">
          {variant.steps.map((step, i) => (
            <li key={step.command} className="flex gap-3">
              <span className="mt-1.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-primary/10 text-[10px] font-semibold text-primary">
                {i + 1}
              </span>
              <div className="min-w-0 flex-1">
                <div className="group relative rounded-lg border border-[#d0d7de] dark:border-white/[0.08] bg-[#f6f8fa] dark:bg-[#161b22] px-3 py-2">
                  <pre className="overflow-x-auto text-[12px] leading-relaxed font-mono text-[#24292f] dark:text-slate-200">
                    <code>{step.command}</code>
                  </pre>
                  <button
                    type="button"
                    onClick={() => void copy(step.command, i)}
                    aria-label={`Copy: ${step.command.split('\n')[0]}`}
                    className="absolute right-1.5 top-1.5 rounded-md p-1 text-[#57606a] opacity-0 transition-opacity hover:bg-black/5 focus-visible:opacity-100 group-hover:opacity-100 dark:text-slate-500 dark:hover:bg-white/10"
                  >
                    {copied === i ? <Check className="h-3.5 w-3.5 text-emerald-500" /> : <Copy className="h-3.5 w-3.5" />}
                  </button>
                </div>
                {step.result && (
                  <p className="mt-1 flex items-start gap-1.5 text-[11.5px] leading-snug text-[#57606a] dark:text-slate-500">
                    <Check className="mt-0.5 h-3 w-3 shrink-0 text-emerald-500" aria-hidden="true" />
                    {step.result}
                  </p>
                )}
              </div>
            </li>
          ))}
        </ol>
      </div>

      <div className="flex items-center gap-2 border-t border-[#d0d7de] dark:border-white/[0.08] bg-[#f6f8fa] dark:bg-[#161b22] px-4 py-2.5">
        <Check className="h-3.5 w-3.5 shrink-0 text-emerald-500" aria-hidden="true" />
        <span className="text-[11.5px] text-[#57606a] dark:text-slate-400">Opens {variant.opens}</span>
      </div>
    </div>
  )
}
