import Link from 'next/link'
import { ArrowRight, MessagesSquare } from 'lucide-react'

/** The Grit WhatsApp community — questions, tutorials, and help from other builders. */
export const WHATSAPP_COMMUNITY_URL = 'https://chat.whatsapp.com/HXsOWlp5W6o9lV4YhsY376'

/**
 * WhatsAppIcon — lucide ships no WhatsApp glyph, so the official mark is
 * inlined. currentColor keeps it themeable alongside the other header icons.
 */
export function WhatsAppIcon({ className = 'h-4 w-4' }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" fill="currentColor" aria-hidden="true" className={className}>
      <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51a12.8 12.8 0 0 0-.57-.01c-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.872.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 0 1-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 0 1-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884a9.82 9.82 0 0 1 6.988 2.896 9.83 9.83 0 0 1 2.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.8 11.8 0 0 0 12.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.9 11.9 0 0 0 5.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.82 11.82 0 0 0-3.48-8.413Z" />
    </svg>
  )
}

/**
 * CommunityCTA — invitation to the WhatsApp community.
 *
 * Two variants because the same ask reads differently depending on where the
 * reader is. "card" is a standalone block for the end of a page, where someone
 * has finished reading and is deciding what to do next. "inline" is a quieter
 * strip for mid-document, where the reader is mid-task and a full call to
 * action would interrupt rather than help.
 */
export function CommunityCTA({
  variant = 'card',
  className = '',
}: {
  variant?: 'card' | 'inline'
  className?: string
}) {
  if (variant === 'inline') {
    return (
      <Link
        href={WHATSAPP_COMMUNITY_URL}
        target="_blank"
        rel="noreferrer"
        className={`group flex items-center gap-3 rounded-xl border border-emerald-500/25 bg-emerald-500/[0.06] px-4 py-3 no-underline transition-colors hover:border-emerald-500/40 hover:bg-emerald-500/[0.10] ${className}`}
      >
        <span className="inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-emerald-500/15 text-emerald-500">
          <WhatsAppIcon className="h-4 w-4" />
        </span>
        <span className="flex-1 text-sm leading-snug text-muted-foreground">
          <span className="font-medium text-foreground">Stuck, or want a second opinion?</span>{' '}
          Ask in the Grit WhatsApp community.
        </span>
        <ArrowRight className="h-4 w-4 shrink-0 text-emerald-500 transition-transform group-hover:translate-x-0.5" />
      </Link>
    )
  }

  return (
    <div
      className={`relative overflow-hidden rounded-2xl border border-emerald-500/25 bg-gradient-to-br from-emerald-500/[0.10] via-card/40 to-card/40 p-7 md:p-8 ${className}`}
    >
      <div className="flex flex-col gap-5 md:flex-row md:items-center md:justify-between">
        <div className="max-w-xl">
          <div className="mb-3 flex items-center gap-2.5">
            <span className="inline-flex h-9 w-9 items-center justify-center rounded-xl bg-emerald-500/15 text-emerald-500">
              <WhatsAppIcon className="h-[18px] w-[18px]" />
            </span>
            <span className="tag-mono text-emerald-500">Community</span>
          </div>
          <h3 className="mb-2 text-xl font-semibold tracking-tight text-foreground md:text-2xl">
            Build alongside other Grit developers
          </h3>
          <p className="text-sm leading-relaxed text-muted-foreground md:text-base">
            Join the WhatsApp community for questions, tutorials, and guidance — from
            people shipping Grit apps, and from the person who builds it. No question
            is too small.
          </p>
        </div>

        <Link
          href={WHATSAPP_COMMUNITY_URL}
          target="_blank"
          rel="noreferrer"
          className="group inline-flex h-11 shrink-0 items-center gap-2 self-start rounded-full bg-emerald-500 px-6 text-sm font-semibold text-white no-underline transition-colors hover:bg-emerald-600 md:self-auto"
        >
          <WhatsAppIcon className="h-4 w-4" />
          Join the community
          <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-0.5" />
        </Link>
      </div>
    </div>
  )
}

/** Compact sidebar entry, for navigation lists rather than page bodies. */
export function CommunitySidebarLink({ className = '' }: { className?: string }) {
  return (
    <Link
      href={WHATSAPP_COMMUNITY_URL}
      target="_blank"
      rel="noreferrer"
      className={`group flex items-center gap-2.5 rounded-lg border border-emerald-500/20 bg-emerald-500/[0.06] px-3 py-2.5 text-xs transition-colors hover:border-emerald-500/40 hover:bg-emerald-500/[0.10] ${className}`}
    >
      <MessagesSquare className="h-3.5 w-3.5 shrink-0 text-emerald-500" />
      <span className="flex-1 font-medium text-foreground/90">Ask the community</span>
      <ArrowRight className="h-3 w-3 shrink-0 text-emerald-500/70 transition-transform group-hover:translate-x-0.5" />
    </Link>
  )
}
