import { ChevronDown, Github } from 'lucide-react'
import type { ReactNode } from 'react'

export interface FAQItem {
  q: string
  a: ReactNode
}

interface PageHelpProps {
  /** Page-specific FAQ items. Omit to render just the GitHub CTA. */
  faqs?: FAQItem[]
  title?: string
}

/**
 * Standard end-of-page block: a short FAQ (native <details>, no JS) plus a
 * call-to-action to raise an issue on GitHub. Drop at the bottom of a docs page.
 */
export function PageHelp({ faqs, title = 'Frequently asked questions' }: PageHelpProps) {
  return (
    <div className="mt-14 border-t border-border pt-10">
      {faqs && faqs.length > 0 && (
        <div className="mb-8">
          <h2 className="text-2xl font-semibold tracking-tight mb-4">{title}</h2>
          <div className="space-y-2">
            {faqs.map((f, i) => (
              <details
                key={i}
                className="group rounded-lg border border-border/50 bg-card/30 px-4 py-3 [&_a]:text-primary [&_a:hover]:underline"
              >
                <summary className="flex cursor-pointer list-none items-center justify-between gap-3 text-[15px] font-medium text-foreground">
                  {f.q}
                  <ChevronDown className="h-4 w-4 shrink-0 text-muted-foreground transition-transform group-open:rotate-180" />
                </summary>
                <div className="mt-2.5 text-sm leading-relaxed text-muted-foreground">{f.a}</div>
              </details>
            ))}
          </div>
        </div>
      )}

      <div className="flex flex-wrap items-center justify-between gap-4 rounded-xl border border-border/50 bg-card/30 p-5">
        <div>
          <p className="text-sm font-semibold text-foreground">
            Something wrong, unclear, or missing on this page?
          </p>
          <p className="text-sm text-muted-foreground">
            Open an issue &mdash; your feedback shapes Grit, and we fix docs fast.
          </p>
        </div>
        <a
          href="https://github.com/MUKE-coder/grit/issues/new"
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex shrink-0 items-center gap-2 rounded-lg border border-primary/25 bg-primary/10 px-4 py-2 text-sm font-medium text-primary transition-colors hover:bg-primary/15"
        >
          <Github className="h-4 w-4" />
          Raise an issue on GitHub
        </a>
      </div>
    </div>
  )
}
