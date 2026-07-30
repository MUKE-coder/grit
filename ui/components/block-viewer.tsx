'use client'

import { useState } from 'react'
import { AnimatePresence, motion } from 'motion/react'
import {
  Check,
  Code2,
  Copy,
  Eye,
  Monitor,
  MoonStar,
  Smartphone,
  Sun,
  Tablet,
} from 'lucide-react'

type Tab = 'preview' | 'code'
type Breakpoint = 'mobile' | 'tablet' | 'desktop'

const WIDTHS: Record<Breakpoint, string> = {
  mobile: '390px',
  tablet: '834px',
  desktop: '100%',
}

const BREAKPOINTS: { key: Breakpoint; icon: typeof Monitor; label: string }[] = [
  { key: 'mobile', icon: Smartphone, label: 'iPhone width' },
  { key: 'tablet', icon: Tablet, label: 'iPad width' },
  { key: 'desktop', icon: Monitor, label: 'Full width' },
]

/**
 * iOS-style segmented control.
 *
 * The selected pill is a single shared element animated between positions with
 * a layout transition, rather than a background toggled per option. That is the
 * detail that makes it feel native: the indicator travels, it does not blink.
 */
function Segmented<T extends string>({
  options,
  value,
  onChange,
  layoutId,
}: {
  options: { key: T; label?: string; icon?: typeof Monitor; title?: string }[]
  value: T
  onChange: (v: T) => void
  layoutId: string
}) {
  return (
    <div className="relative flex rounded-[10px] bg-gray-500/[0.08] p-0.5 dark:bg-white/[0.06]">
      {options.map((o) => {
        const selected = o.key === value
        return (
          <button
            key={o.key}
            type="button"
            onClick={() => onChange(o.key)}
            aria-pressed={selected}
            aria-label={o.title ?? o.label}
            title={o.title ?? o.label}
            className={`relative z-10 inline-flex items-center justify-center gap-1.5 rounded-lg px-3 py-1.5 text-[13px] font-medium transition-colors duration-200 ${
              o.icon && !o.label ? 'w-8' : ''
            } ${
              selected
                ? 'text-gray-900 dark:text-white'
                : 'text-gray-500 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200'
            }`}
          >
            {selected && (
              <motion.span
                layoutId={layoutId}
                transition={{ type: 'spring', stiffness: 520, damping: 38, mass: 0.7 }}
                className="absolute inset-0 -z-10 rounded-lg bg-white shadow-[0_1px_2px_rgb(15_23_42_/_0.10),0_1px_1px_rgb(15_23_42_/_0.04)] dark:bg-white/[0.14] dark:shadow-none"
              />
            )}
            {o.icon && <o.icon aria-hidden="true" className="size-[15px]" />}
            {o.label}
          </button>
        )
      })}
    </div>
  )
}

export function BlockViewer({
  name,
  title,
  source,
  highlighted,
  installCommand,
  height = 660,
}: {
  name: string
  title: string
  source: string
  /** Shiki output, produced at build time. */
  highlighted: string
  installCommand: string
  height?: number
}) {
  const [tab, setTab] = useState<Tab>('preview')
  const [breakpoint, setBreakpoint] = useState<Breakpoint>('desktop')
  const [dark, setDark] = useState(false)

  return (
    <section className="scroll-mt-24" id={name}>
      {/* Toolbar */}
      <div className="mb-3.5 flex flex-wrap items-center gap-3">
        <h3 className="display text-[15px] text-gray-900 dark:text-white">{title}</h3>

        <div className="ml-auto flex flex-wrap items-center gap-2">
          <Segmented
            layoutId={`${name}-tab`}
            value={tab}
            onChange={setTab}
            options={[
              { key: 'preview', label: 'Preview', icon: Eye },
              { key: 'code', label: 'Code', icon: Code2 },
            ]}
          />

          <AnimatePresence initial={false}>
            {tab === 'preview' && (
              <motion.div
                initial={{ opacity: 0, scale: 0.94 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0.94 }}
                transition={{ duration: 0.18, ease: [0.22, 1, 0.36, 1] }}
                className="flex items-center gap-2"
              >
                <div className="hidden sm:block">
                  <Segmented
                    layoutId={`${name}-bp`}
                    value={breakpoint}
                    onChange={setBreakpoint}
                    options={BREAKPOINTS.map((b) => ({
                      key: b.key,
                      icon: b.icon,
                      title: b.label,
                    }))}
                  />
                </div>
                <Segmented
                  layoutId={`${name}-theme`}
                  value={dark ? 'dark' : 'light'}
                  onChange={(v) => setDark(v === 'dark')}
                  options={[
                    { key: 'light', icon: Sun, title: 'Preview in light' },
                    { key: 'dark', icon: MoonStar, title: 'Preview in dark' },
                  ]}
                />
              </motion.div>
            )}
          </AnimatePresence>

          <CopyButton
            value={tab === 'code' ? source : installCommand}
            label={tab === 'code' ? 'Copy code' : 'Copy'}
          />
        </div>
      </div>

      {/* Body */}
      <div className="hairline overflow-hidden rounded-2xl border bg-white lift dark:bg-gray-900">
        {tab === 'preview' ? (
          <div className="flex justify-center bg-gray-50/70 py-0 dark:bg-black/25">
            <motion.iframe
              key={`${name}-${dark ? 'dark' : 'light'}`}
              src={`/preview/${name}?theme=${dark ? 'dark' : 'light'}`}
              title={`${title} preview`}
              loading="lazy"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ duration: 0.25 }}
              // Width is transitioned in CSS rather than by motion: the frame
              // reflows its document as it resizes, and a spring on width makes
              // the content inside jitter. An eased width reads as a window
              // being dragged.
              style={{ width: WIDTHS[breakpoint], height }}
              className="border-0 bg-white transition-[width] duration-[450ms] ease-[cubic-bezier(0.22,1,0.36,1)] dark:bg-gray-900"
            />
          </div>
        ) : (
          <div
            className="shiki-wrap thin-scroll max-h-[660px] overflow-auto bg-gray-50/60 dark:bg-black/30"
            // Shiki output is generated at build time from our own files.
            dangerouslySetInnerHTML={{ __html: highlighted }}
          />
        )}
      </div>

      {/* Install line */}
      <div className="mt-2.5 flex items-center gap-2">
        <code className="truncate font-mono text-[11px] text-gray-400 dark:text-gray-600">
          {installCommand}
        </code>
      </div>
    </section>
  )
}

function CopyButton({ value, label }: { value: string; label: string }) {
  const [state, setState] = useState<'idle' | 'copied' | 'failed'>('idle')

  async function copy() {
    try {
      await navigator.clipboard.writeText(value)
      setState('copied')
    } catch {
      setState('failed')
    }
    setTimeout(() => setState('idle'), 1900)
  }

  return (
    <button
      type="button"
      onClick={copy}
      className="group inline-flex items-center gap-1.5 rounded-[10px] bg-gray-500/[0.08] px-3 py-1.5 text-[13px] font-medium text-gray-600 transition-all duration-200 hover:bg-gray-500/[0.14] active:scale-[0.97] dark:bg-white/[0.06] dark:text-gray-300 dark:hover:bg-white/[0.12]"
    >
      <AnimatePresence mode="wait" initial={false}>
        <motion.span
          key={state}
          initial={{ opacity: 0, y: -4 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: 4 }}
          transition={{ duration: 0.14 }}
          className="inline-flex items-center gap-1.5"
        >
          {state === 'copied' ? (
            <>
              <Check className="size-[15px] text-emerald-500" />
              Copied
            </>
          ) : state === 'failed' ? (
            <>
              <Copy className="size-[15px]" />
              Press &#8984;C
            </>
          ) : (
            <>
              <Copy className="size-[15px]" />
              {label}
            </>
          )}
        </motion.span>
      </AnimatePresence>
    </button>
  )
}
