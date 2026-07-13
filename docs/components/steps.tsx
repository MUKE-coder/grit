import type { ReactNode } from 'react'

/**
 * Numbered steps with a vertical connector line and auto-incrementing badges
 * (CSS counters — no props to keep in sync):
 *
 *   <Steps>
 *     <Step title="Install the CLI">…</Step>
 *     <Step title="Create a project">…</Step>
 *   </Steps>
 */

export function Steps({ children }: { children: ReactNode }) {
  return (
    <div className="my-6 ml-3 border-l-2 border-border/50 pl-9 [counter-reset:step]">
      {children}
    </div>
  )
}

export function Step({ title, children }: { title?: ReactNode; children: ReactNode }) {
  return (
    <div className="relative pb-8 last:pb-1 [counter-increment:step]">
      <div
        aria-hidden
        className="absolute top-0 -left-[49px] flex h-7 w-7 items-center justify-center rounded-full border border-border bg-card text-[13px] font-semibold text-foreground shadow-sm before:[content:counter(step)]"
      />
      {title && (
        <div className="mb-2 text-lg font-semibold tracking-tight text-foreground">{title}</div>
      )}
      <div className="text-[15px] leading-relaxed text-muted-foreground [&_a]:text-primary [&_a:hover]:underline">
        {children}
      </div>
    </div>
  )
}
