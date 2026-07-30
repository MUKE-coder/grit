'use client'

import { useState } from 'react'
import { Check, Clipboard, Monitor, Moon, Smartphone, Sun, Tablet } from 'lucide-react'

type Tab = 'preview' | 'code'
type Breakpoint = 'mobile' | 'tablet' | 'desktop'

/** Frame widths, chosen to sit just inside common device breakpoints. */
const WIDTHS: Record<Breakpoint, string> = {
  mobile: '390px',
  tablet: '768px',
  desktop: '100%',
}

const BREAKPOINTS: { key: Breakpoint; icon: typeof Monitor; label: string }[] = [
  { key: 'mobile', icon: Smartphone, label: 'Mobile' },
  { key: 'tablet', icon: Tablet, label: 'Tablet' },
  { key: 'desktop', icon: Monitor, label: 'Desktop' },
]

export function BlockViewer({
  name,
  title,
  source,
  installCommand,
  height = 640,
}: {
  name: string
  title: string
  source: string
  installCommand: string
  height?: number
}) {
  const [tab, setTab] = useState<Tab>('preview')
  const [breakpoint, setBreakpoint] = useState<Breakpoint>('desktop')
  const [dark, setDark] = useState(false)

  return (
    <section className="scroll-mt-20" id={name}>
      {/* Toolbar */}
      <div className="mb-3 flex flex-wrap items-center gap-3">
        <h3 className="text-sm font-semibold text-gray-900 dark:text-white">{title}</h3>

        <div className="ml-auto flex flex-wrap items-center gap-2">
          {/* Preview / Code */}
          <div className="flex rounded-lg bg-gray-100 p-0.5 dark:bg-white/10">
            {(['preview', 'code'] as Tab[]).map((t) => (
              <button
                key={t}
                type="button"
                onClick={() => setTab(t)}
                aria-pressed={tab === t}
                className={`rounded-md px-3 py-1 text-xs font-medium capitalize transition-colors ${
                  tab === t
                    ? 'bg-white text-gray-900 shadow-sm dark:bg-gray-800 dark:text-white'
                    : 'text-gray-500 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white'
                }`}
              >
                {t}
              </button>
            ))}
          </div>

          {/* Breakpoints — only meaningful while previewing */}
          {tab === 'preview' && (
            <>
              <div className="hidden items-center gap-0.5 rounded-lg bg-gray-100 p-0.5 sm:flex dark:bg-white/10">
                {BREAKPOINTS.map(({ key, icon: Icon, label }) => (
                  <button
                    key={key}
                    type="button"
                    onClick={() => setBreakpoint(key)}
                    aria-pressed={breakpoint === key}
                    aria-label={label}
                    title={label}
                    className={`inline-flex size-7 items-center justify-center rounded-md transition-colors ${
                      breakpoint === key
                        ? 'bg-white text-gray-900 shadow-sm dark:bg-gray-800 dark:text-white'
                        : 'text-gray-500 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white'
                    }`}
                  >
                    <Icon className="size-3.5" />
                  </button>
                ))}
              </div>

              <button
                type="button"
                onClick={() => setDark((d) => !d)}
                aria-pressed={dark}
                aria-label={dark ? 'Preview in light mode' : 'Preview in dark mode'}
                title={dark ? 'Preview in light mode' : 'Preview in dark mode'}
                className="inline-flex size-7 items-center justify-center rounded-md bg-gray-100 text-gray-500 transition-colors hover:text-gray-900 dark:bg-white/10 dark:text-gray-400 dark:hover:text-white"
              >
                {dark ? <Moon className="size-3.5" /> : <Sun className="size-3.5" />}
              </button>
            </>
          )}

          <CopyButton
            value={tab === 'code' ? source : installCommand}
            label={tab === 'code' ? 'Copy code' : 'Copy command'}
          />
        </div>
      </div>

      {/* Body */}
      <div className="overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900">
        {tab === 'preview' ? (
          <div className="flex justify-center bg-gray-50 p-0 dark:bg-gray-950/50">
            <iframe
              // Remounting on theme change forces the frame to reload with the
              // new search param; without a key React keeps the old document.
              key={`${name}-${dark ? 'dark' : 'light'}`}
              src={`/preview/${name}?theme=${dark ? 'dark' : 'light'}`}
              title={`${title} preview`}
              loading="lazy"
              className="border-0 bg-white transition-[width] duration-200 dark:bg-gray-900"
              style={{ width: WIDTHS[breakpoint], height }}
            />
          </div>
        ) : (
          <pre className="max-h-[640px] overflow-auto bg-gray-950 p-5 font-mono text-[13px]/6 text-gray-200">
            <code>{source}</code>
          </pre>
        )}
      </div>

      {/* Install line */}
      <div className="mt-2 flex items-center gap-2 overflow-x-auto">
        <code className="font-mono text-[11px] text-gray-500 dark:text-gray-500">
          {installCommand}
        </code>
      </div>
    </section>
  )
}

function CopyButton({ value, label }: { value: string; label: string }) {
  const [copied, setCopied] = useState(false)
  const [failed, setFailed] = useState(false)

  async function copy() {
    try {
      await navigator.clipboard.writeText(value)
      setCopied(true)
      setTimeout(() => setCopied(false), 1800)
    } catch {
      // clipboard is unavailable on insecure origins — say so rather than
      // showing a success state that did not happen.
      setFailed(true)
      setTimeout(() => setFailed(false), 2400)
    }
  }

  return (
    <button
      type="button"
      onClick={copy}
      className="inline-flex items-center gap-1.5 rounded-lg bg-gray-100 px-2.5 py-1.5 text-xs font-medium text-gray-600 transition-colors hover:text-gray-900 dark:bg-white/10 dark:text-gray-400 dark:hover:text-white"
    >
      {copied ? (
        <>
          <Check className="size-3.5 text-green-600 dark:text-green-400" />
          Copied
        </>
      ) : failed ? (
        <>
          <Clipboard className="size-3.5" />
          Press &#8984;C
        </>
      ) : (
        <>
          <Clipboard className="size-3.5" />
          {label}
        </>
      )}
    </button>
  )
}
