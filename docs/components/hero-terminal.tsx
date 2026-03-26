'use client'

import { useState, useEffect } from 'react'

const lines = [
  { text: '$ grit new my-saas --triple --vite', delay: 0, type: 'command' as const },
  { text: '', delay: 800, type: 'blank' as const },
  { text: '  ██████╗ ██████╗ ██╗████████╗', delay: 1000, type: 'ascii' as const },
  { text: ' ██╔════╝ ██╔══██╗██║╚══██╔══╝', delay: 1050, type: 'ascii' as const },
  { text: ' ██║  ███╗██████╔╝██║   ██║', delay: 1100, type: 'ascii' as const },
  { text: ' ██║   ██║██╔══██╗██║   ██║', delay: 1150, type: 'ascii' as const },
  { text: ' ╚██████╔╝██║  ██║██║   ██║', delay: 1200, type: 'ascii' as const },
  { text: '  ╚═════╝ ╚═╝  ╚═╝╚═╝   ╚═╝', delay: 1250, type: 'ascii' as const },
  { text: '', delay: 1400, type: 'blank' as const },
  { text: '  Creating new Grit project: my-saas', delay: 1600, type: 'info' as const },
  { text: '  Architecture: triple | Frontend: tanstack', delay: 1800, type: 'dim' as const },
  { text: '', delay: 2000, type: 'blank' as const },
  { text: '  → Creating directory structure...', delay: 2200, type: 'step' as const },
  { text: '  → Scaffolding Go API...', delay: 2600, type: 'step' as const },
  { text: '  → Adding batteries (cache, storage, mail, jobs, AI, TOTP)...', delay: 3000, type: 'step' as const },
  { text: '  → Scaffolding TanStack Router web app (Vite)...', delay: 3400, type: 'step' as const },
  { text: '  → Scaffolding TanStack Router admin panel (Vite)...', delay: 3800, type: 'step' as const },
  { text: '', delay: 4200, type: 'blank' as const },
  { text: '  ✓ Project created successfully!', delay: 4400, type: 'success' as const },
]

export function HeroTerminal() {
  const [visibleLines, setVisibleLines] = useState(0)

  useEffect(() => {
    const timers: NodeJS.Timeout[] = []
    lines.forEach((line, i) => {
      timers.push(setTimeout(() => setVisibleLines(i + 1), line.delay))
    })
    return () => timers.forEach(clearTimeout)
  }, [])

  return (
    <div className="rounded-xl border border-slate-200 dark:border-border/40 bg-slate-900 overflow-hidden shadow-2xl shadow-black/20">
      {/* Window chrome */}
      <div className="flex items-center gap-2 px-4 py-3 border-b border-slate-700/50">
        <div className="flex items-center gap-1.5">
          <div className="h-3 w-3 rounded-full bg-[#ff5f57]" />
          <div className="h-3 w-3 rounded-full bg-[#febc2e]" />
          <div className="h-3 w-3 rounded-full bg-[#28c840]" />
        </div>
        <span className="text-xs font-mono text-slate-500 ml-2">Terminal</span>
      </div>

      {/* Terminal content */}
      <div className="p-5 font-mono text-[13px] leading-6 min-h-[340px]">
        {lines.slice(0, visibleLines).map((line, i) => {
          if (line.type === 'blank') return <div key={i} className="h-5" />

          const colorClass = {
            command: 'text-slate-200',
            ascii: 'text-primary/70',
            info: 'text-primary font-semibold',
            dim: 'text-slate-500',
            step: 'text-slate-400',
            success: 'text-emerald-400 font-semibold',
          }[line.type]

          return (
            <div key={i} className={`${colorClass} transition-opacity duration-300`}>
              {line.text}
              {i === visibleLines - 1 && line.type === 'command' && (
                <span className="inline-block w-2 h-4 bg-primary/80 ml-0.5 animate-pulse align-middle" />
              )}
            </div>
          )
        })}

        {visibleLines === 0 && (
          <div className="text-slate-200">
            <span className="text-sky-400/60">$</span>
            <span className="inline-block w-2 h-4 bg-primary/80 ml-1 animate-pulse align-middle" />
          </div>
        )}
      </div>
    </div>
  )
}
