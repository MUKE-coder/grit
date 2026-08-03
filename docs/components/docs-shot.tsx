/**
 * A captioned screenshot for docs pages.
 *
 * Server component on purpose — these are static images with no interaction,
 * and docs pages should not pay for hydration to show a picture.
 *
 * Every image passed to this lives under /public/images and is a capture of a
 * real generated project. The caption is not decoration: it says which screen
 * this is, so a reader landing mid-page knows what they are looking at without
 * scrolling up for context.
 */

interface DocsShotProps {
  src: string
  alt: string
  /** Shown in the chrome bar — a URL, or an app name for native windows. */
  label: string
  /** Sentence under the image. Say what to look at, not what it obviously is. */
  caption?: string
  /** Phone captures are portrait and would be absurd at full width. */
  narrow?: boolean
}

export function DocsShot({ src, alt, label, caption, narrow = false }: DocsShotProps) {
  return (
    <figure className={`my-8 ${narrow ? 'max-w-xs mx-auto' : ''}`}>
      <div className="rounded-xl overflow-hidden border border-border bg-card/40 shadow-[0_18px_48px_-16px_rgba(2,6,23,0.4)]">
        <div className="flex items-center gap-2 px-3.5 py-2.5 bg-card/70 border-b border-border/60">
          <div className="flex items-center gap-1.5">
            <span className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
            <span className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
            <span className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
          </div>
          <span className="mx-auto text-[11px] font-mono text-muted-foreground truncate max-w-[70%]">
            {label}
          </span>
        </div>
        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img src={src} alt={alt} className="w-full h-auto block" loading="lazy" />
      </div>
      {caption && (
        <figcaption className="text-[12.5px] text-muted-foreground mt-3 leading-relaxed">
          {caption}
        </figcaption>
      )}
    </figure>
  )
}
