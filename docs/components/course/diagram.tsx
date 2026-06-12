import { Workflow } from 'lucide-react'

interface Props {
  /** Optional caption shown below the diagram. */
  caption?: React.ReactNode
  /** Optional label shown above the diagram (e.g., "Flow", "Architecture"). */
  label?: string
  /** Diagram content. Pass ASCII art / a <pre> block / an SVG / a React tree. */
  children: React.ReactNode
}

/**
 * Diagram — a styled wrapper for any visual aid: ASCII flow, SVG
 * illustration, request-trace, mental model. Centred, dark, monospace
 * by default. Pass <svg> children when ASCII isn't enough.
 *
 * Use this for any concept where words alone are clumsy: lifecycle
 * diagrams, request flow, file-system layouts, sequence interactions.
 */
export function Diagram({ caption, label = 'Diagram', children }: Props) {
  return (
    <figure className="not-prose my-8 rounded-2xl border border-border/60 bg-card/40 overflow-hidden">
      <div className="flex items-center gap-2 px-5 py-2.5 border-b border-border/40 bg-muted/30">
        <Workflow className="h-3.5 w-3.5 text-primary" />
        <p className="text-[10px] uppercase tracking-wider text-muted-foreground font-mono">
          {label}
        </p>
      </div>
      <div className="px-5 py-6 overflow-x-auto">
        <div className="font-mono text-[12.5px] leading-[1.65] text-foreground/90 whitespace-pre">
          {children}
        </div>
      </div>
      {caption && (
        <figcaption className="px-5 py-3 border-t border-border/40 bg-muted/20 text-xs text-muted-foreground italic">
          {caption}
        </figcaption>
      )}
    </figure>
  )
}
