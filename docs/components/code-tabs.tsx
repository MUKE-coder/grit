'use client'

import { useState } from 'react'
import { CodeBlock } from './code-block'
import { CopyButton } from './copy-button'

export interface CodeTabFile {
  filename: string
  language?: string
  code: string
}

interface CodeTabsProps {
  files: CodeTabFile[]
  className?: string
}

/**
 * Multi-file tabbed code block — Grit's signature "one command, whole stack"
 * moment. Show the model, service, handler, hooks and page a single
 * `grit generate resource` emits, all in one switchable block.
 */
export function CodeTabs({ files, className = 'mb-6' }: CodeTabsProps) {
  const [active, setActive] = useState(0)
  if (files.length === 0) return null
  const current = files[active] ?? files[0]

  return (
    <div className={className}>
      {/* Tab strip = the header. The filenames are the tabs. */}
      <div className="flex items-center justify-between rounded-t-lg border-2 border-b border-slate-200/80 bg-slate-50/70 dark:border-white/10 dark:bg-white/[0.02] pl-1.5 pr-1">
        <div className="flex overflow-x-auto">
          {files.map((f, i) => (
            <button
              key={f.filename + i}
              type="button"
              onClick={() => setActive(i)}
              className={`whitespace-nowrap px-3 py-2 text-[12px] font-mono transition-colors border-b-2 -mb-px ${
                i === active
                  ? 'border-primary text-foreground'
                  : 'border-transparent text-slate-500 dark:text-slate-400 hover:text-foreground'
              }`}
            >
              {f.filename}
            </button>
          ))}
        </div>
        <CopyButton
          text={current.code.trim()}
          className="h-7 w-7 shrink-0 text-slate-500 hover:text-slate-900 dark:hover:text-slate-100"
        />
      </div>

      {/* Headerless code body — squared top edge so it connects to the tabs. */}
      <CodeBlock
        code={current.code}
        language={current.language ?? 'go'}
        className="mb-0 rounded-t-none border-t-0"
      />
    </div>
  )
}
