'use client'

import { useState } from 'react'
import { AnimatePresence, motion } from 'motion/react'
import { Check, Copy, Terminal } from 'lucide-react'
import { trackBlockCopy } from '@/lib/track'

/**
 * The install command for a block, as a first-class surface.
 *
 * It sits directly above the preview rather than as a caption below it. The
 * command IS the product here — the whole point of the page is to get a block
 * into someone's project — and a small grey line under a large preview is
 * something people genuinely do not see.
 */
export function InstallCommand({
  command,
  block,
  category,
  subcategory,
}: {
  command: string
  block: string
  category: string
  subcategory: string
}) {
  const [state, setState] = useState<'idle' | 'copied' | 'failed'>('idle')

  async function copy() {
    try {
      await navigator.clipboard.writeText(command)
      setState('copied')
      trackBlockCopy({ block, category, subcategory, kind: 'command' })
    } catch {
      setState('failed')
    }
    setTimeout(() => setState('idle'), 2000)
  }

  return (
    <button
      type="button"
      onClick={copy}
      // The whole bar is the button. A separate small copy target next to a wide
      // strip of text invites a miss, and there is nothing else to do here.
      className="group hairline flex w-full items-center gap-3 rounded-xl border bg-gray-500/[0.045] px-4 py-3 text-left transition-all duration-200 hover:bg-gray-500/[0.075] active:scale-[0.995] dark:bg-white/[0.035] dark:hover:bg-white/[0.06]"
    >
      <Terminal
        aria-hidden="true"
        className="size-4 shrink-0 text-indigo-500 dark:text-indigo-400"
      />

      <code
        className="no-scroll min-w-0 flex-1 overflow-x-auto font-mono text-[12.5px] whitespace-nowrap text-gray-700 dark:text-gray-300"
        style={{
          maskImage: 'linear-gradient(to right, black calc(100% - 1.5rem), transparent)',
          WebkitMaskImage:
            'linear-gradient(to right, black calc(100% - 1.5rem), transparent)',
        }}
      >
        {command}
      </code>

      <span className="inline-flex h-7 shrink-0 items-center gap-1.5 rounded-lg bg-white px-2.5 text-[12px] font-semibold text-gray-700 shadow-[0_1px_2px_rgb(15_23_42_/_0.08)] transition-colors group-hover:text-gray-900 dark:bg-white/[0.12] dark:text-gray-200 dark:shadow-none dark:group-hover:text-white">
        <AnimatePresence mode="wait" initial={false}>
          <motion.span
            key={state}
            initial={{ opacity: 0, y: -3 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: 3 }}
            transition={{ duration: 0.14 }}
            className="inline-flex items-center gap-1.5"
          >
            {state === 'copied' ? (
              <>
                <Check className="size-3.5 text-emerald-500" />
                Copied
              </>
            ) : state === 'failed' ? (
              <>
                <Copy className="size-3.5" />
                &#8984;C
              </>
            ) : (
              <>
                <Copy className="size-3.5" />
                Copy
              </>
            )}
          </motion.span>
        </AnimatePresence>
      </span>
    </button>
  )
}
