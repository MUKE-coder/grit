'use client'

import { useState } from 'react'
import { ArrowUp, ChevronDown, Globe, Play, Plus, Sparkles } from 'lucide-react'

/** The Grit UI mark, inlined so the block stays self-contained. */
function GritMark({ className = 'size-7' }: { className?: string }) {
  return (
    <svg viewBox="0 0 32 32" fill="none" aria-hidden="true" className={className}>
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M8 0H24A8 8 0 0 1 32 8V24A8 8 0 0 1 24 32H8A8 8 0 0 1 0 24V8A8 8 0 0 1 8 0ZM16 9.4A6.6 6.6 0 1 0 21.4 19.9V17.4H17.2A1.7 1.7 0 0 1 17.2 14H23.1A1.7 1.7 0 0 1 24.8 15.7V20.6A1.7 1.7 0 0 1 24.4 21.7A10 10 0 1 1 22.6 8.2A1.7 1.7 0 0 1 20.4 10.8A6.6 6.6 0 0 0 16 9.4Z"
        fill="currentColor"
      />
    </svg>
  )
}

/** Invented names — real marks in a demo logo row imply endorsements. */
const LOGOS = ['NORTHWIND', 'ACME', 'GLOBEX', 'BEACON']

export default function WithAiChat({
  title = 'Powerful analytics for smarter decisions',
  subtitle = "With Grit's built-in assistant, get your projects to the finish line faster and with no context switching.",
  primaryLabel = 'Start using',
  secondaryLabel = 'Watch video',
  assistantName = 'Grit AI',
  assistantTagline = 'Your personal assistant, when it really matters',
}: {
  title?: string
  subtitle?: string
  primaryLabel?: string
  secondaryLabel?: string
  assistantName?: string
  assistantTagline?: string
}) {
  const [prompt, setPrompt] = useState('')

  return (
    <div className="relative isolate min-h-screen overflow-hidden bg-gray-50 dark:bg-gray-950">
      {/* Warm bloom behind the card, the signature of this layout */}
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-x-0 top-1/3 -z-10 h-[36rem] opacity-70 dark:opacity-40"
        style={{
          background:
            'radial-gradient(32rem 20rem at 22% 40%, rgba(251,191,36,0.16), transparent 65%), radial-gradient(32rem 20rem at 78% 45%, rgba(129,140,248,0.20), transparent 65%)',
        }}
      />

      {/* Floating pill nav */}
      <div className="px-6 pt-6">
        <nav
          aria-label="Global"
          className="mx-auto flex max-w-6xl items-center justify-between"
        >
          <a href="#" className="flex items-center gap-2 text-gray-900 dark:text-white">
            <GritMark className="size-7" />
            <span className="text-lg font-semibold tracking-tight">grit</span>
          </a>

          <div className="hidden items-center gap-1 rounded-full border border-gray-200 bg-white/80 px-2 py-1.5 shadow-sm backdrop-blur lg:flex dark:border-white/10 dark:bg-white/5">
            {['Product', 'Solutions'].map((item) => (
              <button
                key={item}
                type="button"
                className="inline-flex items-center gap-1 rounded-full px-4 py-1.5 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-white/10"
              >
                {item}
                <ChevronDown aria-hidden="true" className="size-3.5 text-gray-400" />
              </button>
            ))}
            {['Pricing', 'Customers'].map((item) => (
              <a
                key={item}
                href="#"
                className="rounded-full px-4 py-1.5 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-white/10"
              >
                {item}
              </a>
            ))}
          </div>

          <div className="flex items-center gap-3">
            <a
              href="#"
              className="hidden text-sm font-medium text-gray-700 transition-colors hover:text-gray-900 sm:block dark:text-gray-300 dark:hover:text-white"
            >
              Sign in
            </a>
            <a
              href="#"
              className="rounded-full bg-indigo-600 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-indigo-500"
            >
              Start for free
            </a>
          </div>
        </nav>
      </div>

      {/* Hero */}
      <div className="mx-auto max-w-3xl px-6 pt-20 text-center sm:pt-28">
        <h1 className="text-5xl font-bold tracking-tight text-balance text-gray-900 sm:text-6xl dark:text-white">
          {title}
        </h1>
        <p className="mx-auto mt-6 max-w-xl text-lg text-gray-600 dark:text-gray-400">
          {subtitle}
        </p>
        <div className="mt-9 flex flex-wrap items-center justify-center gap-3">
          <a
            href="#"
            className="rounded-xl bg-indigo-600 px-6 py-3 text-sm font-semibold text-white transition-colors hover:bg-indigo-500"
          >
            {primaryLabel}
          </a>
          <a
            href="#"
            className="inline-flex items-center gap-2 rounded-xl border border-gray-300 bg-white px-6 py-3 text-sm font-semibold text-gray-900 transition-colors hover:bg-gray-50 dark:border-white/15 dark:bg-white/5 dark:text-white dark:hover:bg-white/10"
          >
            <Play aria-hidden="true" className="size-3.5 fill-current" />
            {secondaryLabel}
          </a>
        </div>
      </div>

      {/* Assistant card */}
      <div className="mx-auto max-w-md px-6 pt-20 pb-16">
        <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-2xl dark:border-white/10 dark:bg-gray-900">
          <div className="flex gap-1.5 px-4 pt-4">
            <span className="size-2.5 rounded-full bg-gray-300 dark:bg-gray-700" />
            <span className="size-2.5 rounded-full bg-gray-300 dark:bg-gray-700" />
            <span className="size-2.5 rounded-full bg-gray-300 dark:bg-gray-700" />
          </div>

          <div className="px-5 pt-5 pb-3 text-center">
            <div className="mb-1.5 inline-flex items-center gap-2">
              <span className="inline-flex size-5 items-center justify-center rounded-full bg-indigo-600 text-white">
                <Sparkles aria-hidden="true" className="size-3" />
              </span>
              <span className="text-base font-semibold text-gray-900 dark:text-white">
                {assistantName}
              </span>
            </div>
            <p className="text-sm text-gray-500 dark:text-gray-400">{assistantTagline}</p>
          </div>

          <div className="mx-5 border-t border-gray-100 dark:border-white/5" />

          {/* Conversation */}
          <div className="space-y-4 p-5">
            <div className="flex flex-col items-end">
              <div className="max-w-[85%] rounded-2xl rounded-br-md bg-gray-100 px-4 py-3 text-sm text-gray-800 dark:bg-white/10 dark:text-gray-200">
                Which of my resources are missing an owner check?
              </div>
              <span className="mt-1 text-[11px] text-gray-600 dark:text-gray-400">Now</span>
            </div>

            <div className="flex gap-2.5">
              <span className="mt-0.5 inline-flex size-6 shrink-0 items-center justify-center rounded-full bg-indigo-600/10 text-indigo-600 dark:bg-indigo-500/15 dark:text-indigo-400">
                <Sparkles aria-hidden="true" className="size-3" />
              </span>
              <p className="text-sm/6 text-gray-700 dark:text-gray-300">
                Three handlers query by ID without scoping to the current user:{' '}
                <code className="rounded bg-gray-100 px-1 py-0.5 font-mono text-[12px] dark:bg-white/10">
                  invoice
                </code>
                ,{' '}
                <code className="rounded bg-gray-100 px-1 py-0.5 font-mono text-[12px] dark:bg-white/10">
                  upload
                </code>{' '}
                and{' '}
                <code className="rounded bg-gray-100 px-1 py-0.5 font-mono text-[12px] dark:bg-white/10">
                  note
                </code>
                . Want me to add the ownership guard to each?
              </p>
            </div>
          </div>

          {/* Composer */}
          <div className="m-3 rounded-xl bg-gray-50 p-3 dark:bg-white/5">
            <label htmlFor="grit-ai-prompt" className="sr-only">
              Ask anything
            </label>
            <input
              id="grit-ai-prompt"
              value={prompt}
              onChange={(e) => setPrompt(e.target.value)}
              placeholder="Ask anything"
              className="w-full bg-transparent text-sm text-gray-900 placeholder:text-gray-400 focus:outline-none dark:text-white dark:placeholder:text-gray-500"
            />
            <div className="mt-3 flex items-center gap-2">
              <button
                type="button"
                aria-label="Add attachment"
                className="inline-flex size-8 items-center justify-center rounded-full border border-gray-200 text-gray-500 transition-colors hover:text-gray-900 dark:border-white/10 dark:text-gray-400 dark:hover:text-white"
              >
                <Plus aria-hidden="true" className="size-4" />
              </button>
              <button
                type="button"
                aria-label="Search the web"
                className="inline-flex size-8 items-center justify-center rounded-full border border-gray-200 text-gray-500 transition-colors hover:text-gray-900 dark:border-white/10 dark:text-gray-400 dark:hover:text-white"
              >
                <Globe aria-hidden="true" className="size-4" />
              </button>
              <button
                type="button"
                aria-label="Send"
                className="ml-auto inline-flex size-8 items-center justify-center rounded-full bg-gray-900 text-white transition-colors hover:bg-gray-700 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-200"
              >
                <ArrowUp aria-hidden="true" className="size-4" />
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Logos */}
      <div className="px-6 pb-20 text-center">
        <p className="text-sm font-medium text-gray-600 dark:text-gray-400">
          Trusted by teams at
        </p>
        <div className="mt-6 flex flex-wrap items-center justify-center gap-x-14 gap-y-6">
          {LOGOS.map((logo) => (
            <span
              key={logo}
              className="font-mono text-sm tracking-[0.15em] text-gray-600 dark:text-gray-400"
            >
              {logo}
            </span>
          ))}
        </div>
      </div>
    </div>
  )
}
