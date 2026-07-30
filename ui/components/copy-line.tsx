'use client'

import { useState } from 'react'
import { AnimatePresence, motion } from 'motion/react'
import { Check, Copy, Terminal } from 'lucide-react'

/** A copyable command line, styled as a glass surface. */
export function CopyLine({ value }: { value: string }) {
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
    <div className="glass flex items-center gap-3 rounded-xl px-4 py-3">
      <Terminal aria-hidden="true" className="size-4 shrink-0 text-gray-400" />
      {/* Masked fade rather than a scrollbar. A long command is copied, not
          read, so a visible track is noise on an otherwise clean surface. */}
      <code
        className="no-scroll min-w-0 flex-1 overflow-x-auto font-mono text-[12.5px] whitespace-nowrap text-gray-700 dark:text-gray-300"
        style={{
          maskImage: 'linear-gradient(to right, black calc(100% - 2rem), transparent)',
          WebkitMaskImage: 'linear-gradient(to right, black calc(100% - 2rem), transparent)',
        }}
      >
        {value}
      </code>
      <button
        type="button"
        onClick={copy}
        aria-label="Copy command"
        className="inline-flex size-7 shrink-0 items-center justify-center rounded-lg text-gray-400 transition-all duration-200 hover:bg-gray-500/[0.10] hover:text-gray-900 active:scale-[0.92] dark:hover:bg-white/[0.10] dark:hover:text-white"
      >
        <AnimatePresence mode="wait" initial={false}>
          <motion.span
            key={state}
            initial={{ opacity: 0, scale: 0.7 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.7 }}
            transition={{ duration: 0.15 }}
          >
            {state === 'copied' ? (
              <Check className="size-4 text-emerald-500" />
            ) : (
              <Copy className="size-4" />
            )}
          </motion.span>
        </AnimatePresence>
      </button>
    </div>
  )
}
