'use client'

import { useState } from 'react'
import { Check, Copy, Globe, Layers, Monitor, Smartphone, Server } from 'lucide-react'
import { Tabs } from '@/components/grit/tabs'

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
    blurb: 'Next.js + admin + Go API',
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
    label: 'Desktop',
    icon: Monitor,
    blurb: 'Wails window, offline-capable',
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
    label: 'Mobile',
    icon: Smartphone,
    blurb: 'Expo and React Native',
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
    blurb: 'Go, Gin and GORM. No front end',
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
    blurb: 'Every app, one repo',
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

/** The binary at the head of a command, coloured apart from its arguments. */
function Command({ text }: { text: string }) {
  const space = text.indexOf(' ')
  const bin = space === -1 ? text : text.slice(0, space)
  const rest = space === -1 ? '' : text.slice(space)
  return (
    <code className="block whitespace-pre-wrap break-words font-mono text-[13px] leading-[22px] text-code-fg">
      <span className="text-code-fn">{bin}</span>
      {rest}
    </code>
  )
}

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
    <div className="min-w-0">
      <Tabs
        label="What are you building?"
        value={active}
        onChange={setActive}
        className="max-w-full overflow-x-auto"
        items={VARIANTS.map((v) => ({
          id: v.key,
          label: (
            <>
              <v.icon className="h-3.5 w-3.5" aria-hidden="true" />
              {v.label}
            </>
          ),
        }))}
      />

      <figure
        className="mt-3 overflow-hidden rounded-2xl border border-border bg-surface-1"
        style={{
          boxShadow:
            '0 40px 100px -40px rgba(169,155,255,0.18), 0 30px 60px -30px rgba(0,0,0,0.8)',
        }}
      >
        <div className="flex h-11 items-center gap-3 border-b border-surface-3 px-4">
          <span aria-hidden="true" className="g-terminal__dots">
            <i />
            <i />
            <i />
          </span>
          <span className="truncate font-mono text-[12px] font-medium text-fg-subtle">
            ~/contacts · {variant.blurb}
          </span>
        </div>

        {/* A fixed height, so switching tabs does not resize the hero and shove
            the page under the reader's cursor. */}
        <ol className="h-[430px] overflow-y-auto px-2.5 pb-1.5 pt-3.5">
          {variant.steps.map((step, i) => (
            <li
              key={step.command}
              className="group grid grid-cols-[26px_minmax(0,1fr)_32px] gap-2 rounded-md px-1.5 py-[7px] transition-colors hover:bg-bg/60"
            >
              <span className="text-right font-mono text-[11px] leading-[22px] text-fg-subtle">
                {i + 1}
              </span>
              <span className="min-w-0">
                <Command text={step.command} />
                {step.result && (
                  <span className="mt-0.5 block text-[12.5px] leading-[18px] text-fg-subtle">
                    <span className="text-success" aria-hidden="true">
                      ✓{' '}
                    </span>
                    {step.result}
                  </span>
                )}
              </span>
              <button
                type="button"
                onClick={() => void copy(step.command, i)}
                aria-label={`Copy: ${step.command.split('\n')[0]}`}
                className="g-copy h-[30px] w-[30px] opacity-0 transition-opacity focus-visible:opacity-100 group-hover:opacity-100"
              >
                {copied === i ? (
                  <Check className="h-3.5 w-3.5 text-success" />
                ) : (
                  <Copy className="h-3.5 w-3.5" />
                )}
              </button>
            </li>
          ))}
        </ol>

        <div className="flex items-center gap-2.5 border-t border-surface-3 px-[22px] pb-4 pt-3 font-mono text-[12.5px] text-fg-muted">
          <span
            aria-hidden="true"
            className="g-ping h-[7px] w-[7px] flex-none rounded-full bg-success"
          />
          Opens {variant.opens}
          <span className="g-caret" aria-hidden="true" />
        </div>
      </figure>
    </div>
  )
}
