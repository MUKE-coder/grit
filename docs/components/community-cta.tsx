import Link from 'next/link'
import { ArrowRight, MessagesSquare } from 'lucide-react'
import { siteConfig } from '@/config/site'

/**
 * The Grit Discord — the main place to get help. Channels for installing,
 * questions and bugs, so an answer stays findable instead of scrolling away.
 *
 * Re-exported from siteConfig rather than written twice: the invite also goes
 * into the page metadata and the structured data, and two copies of a URL that
 * must agree is how a dead invite link survives a rotation.
 */
export const DISCORD_INVITE_URL = siteConfig.discord

/** The Grit WhatsApp community — questions, tutorials, and help from other builders. */
export const WHATSAPP_COMMUNITY_URL = 'https://chat.whatsapp.com/HXsOWlp5W6o9lV4YhsY376'

/**
 * DiscordIcon — lucide ships no Discord glyph, so the official mark is
 * inlined. currentColor keeps it themeable alongside the other header icons.
 */
export function DiscordIcon({ className = 'h-4 w-4' }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" fill="currentColor" aria-hidden="true" className={className}>
      <path d="M20.317 4.3698a19.7913 19.7913 0 00-4.8851-1.5152.0741.0741 0 00-.0785.0371c-.211.3753-.4447.8648-.6083 1.2495-1.8447-.2762-3.68-.2762-5.4868 0-.1636-.3933-.4058-.8742-.6177-1.2495a.077.077 0 00-.0785-.037 19.7363 19.7363 0 00-4.8852 1.515.0699.0699 0 00-.0321.0277C.5334 9.0458-.319 13.5799.0992 18.0578a.0824.0824 0 00.0312.0561c2.0528 1.5076 4.0413 2.4228 5.9929 3.0294a.0777.0777 0 00.0842-.0276c.4616-.6304.8731-1.2952 1.226-1.9942a.076.076 0 00-.0416-.1057c-.6528-.2476-1.2743-.5495-1.8722-.8923a.077.077 0 01-.0076-.1277c.1258-.0943.2517-.1923.3718-.2914a.0743.0743 0 01.0776-.0105c3.9278 1.7933 8.18 1.7933 12.0614 0a.0739.0739 0 01.0785.0095c.1202.099.246.1981.3728.2924a.077.077 0 01-.0066.1276 12.2986 12.2986 0 01-1.873.8914.0766.0766 0 00-.0407.1067c.3604.698.7719 1.3628 1.225 1.9932a.076.076 0 00.0842.0286c1.961-.6067 3.9495-1.5219 6.0023-3.0294a.077.077 0 00.0313-.0552c.5004-5.177-.8382-9.6739-3.5485-13.6604a.061.061 0 00-.0312-.0286zM8.02 15.3312c-1.1825 0-2.1569-1.0857-2.1569-2.419 0-1.3332.9555-2.4189 2.157-2.4189 1.2108 0 2.1757 1.0952 2.1568 2.419 0 1.3332-.9555 2.4189-2.1569 2.4189zm7.9748 0c-1.1825 0-2.1569-1.0857-2.1569-2.419 0-1.3332.9554-2.4189 2.1569-2.4189 1.2108 0 2.1757 1.0952 2.1568 2.419 0 1.3332-.946 2.4189-2.1568 2.4189Z" />
    </svg>
  )
}

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
 * CommunityCTA — invitation to the community.
 *
 * Discord leads and WhatsApp follows, because the two rooms are not equivalent
 * for the thing being asked: a Discord answer lands in a named channel and
 * stays searchable, where the same answer in a WhatsApp group scrolls away
 * within a day and gets asked again next week. WhatsApp stays because people
 * are already in it.
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
        href={DISCORD_INVITE_URL}
        target="_blank"
        rel="noreferrer"
        className={`group flex items-center gap-3 rounded-xl border border-accent/25 bg-accent/[0.06] px-4 py-3 no-underline transition-colors hover:border-accent/40 hover:bg-accent/[0.10] ${className}`}
      >
        <span className="inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-accent/15 text-accent">
          <DiscordIcon className="h-4 w-4" />
        </span>
        <span className="flex-1 text-sm leading-snug text-muted-foreground">
          <span className="font-medium text-foreground">Stuck, or want a second opinion?</span>{' '}
          Ask in the Grit Discord.
        </span>
        <ArrowRight className="h-4 w-4 shrink-0 text-accent transition-transform group-hover:translate-x-0.5" />
      </Link>
    )
  }

  return (
    <div
      className={`relative overflow-hidden rounded-2xl border border-accent/25 bg-gradient-to-br from-accent/[0.10] via-card/40 to-card/40 p-7 md:p-8 ${className}`}
    >
      <div className="flex flex-col gap-5 md:flex-row md:items-center md:justify-between">
        <div className="max-w-xl">
          <div className="mb-3 flex items-center gap-2.5">
            <span className="inline-flex h-9 w-9 items-center justify-center rounded-xl bg-accent/15 text-accent">
              <DiscordIcon className="h-[18px] w-[18px]" />
            </span>
            <span className="tag-mono text-accent">Community</span>
          </div>
          <h3 className="mb-2 text-xl font-semibold tracking-tight text-foreground md:text-2xl">
            Build alongside other Grit developers
          </h3>
          <p className="text-sm leading-relaxed text-muted-foreground md:text-base">
            Join the Discord for questions, tutorials, and guidance from people
            shipping Grit apps, and from the person who builds it. There are
            channels for installing, for questions and for bugs, so you can ask in
            the right place and read what has already been answered. No question is
            too small.
          </p>
        </div>

        <div className="flex shrink-0 flex-col gap-2.5 self-start md:self-auto">
          <Link
            href={DISCORD_INVITE_URL}
            target="_blank"
            rel="noreferrer"
            className="group inline-flex h-11 items-center justify-center gap-2 rounded-full bg-primary px-6 text-sm font-semibold text-primary-foreground no-underline transition-colors hover:bg-accent/85"
          >
            <DiscordIcon className="h-4 w-4" />
            Join the Discord
            <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-0.5" />
          </Link>
          <Link
            href={WHATSAPP_COMMUNITY_URL}
            target="_blank"
            rel="noreferrer"
            className="inline-flex h-9 items-center justify-center gap-2 rounded-full border border-emerald-600/30 px-5 text-xs font-medium text-emerald-700 no-underline transition-colors hover:border-emerald-600/50 hover:bg-emerald-600/10 dark:border-emerald-500/30 dark:text-emerald-400 dark:hover:border-emerald-500/50 dark:hover:bg-emerald-500/10"
          >
            <WhatsAppIcon className="h-3.5 w-3.5" />
            Or the WhatsApp group
          </Link>
        </div>
      </div>
    </div>
  )
}

/** Compact sidebar entry, for navigation lists rather than page bodies. */
export function CommunitySidebarLink({ className = '' }: { className?: string }) {
  return (
    <Link
      href={DISCORD_INVITE_URL}
      target="_blank"
      rel="noreferrer"
      className={`group flex items-center gap-2.5 rounded-lg border border-accent/20 bg-accent/[0.06] px-3 py-2.5 text-xs transition-colors hover:border-accent/40 hover:bg-accent/[0.10] ${className}`}
    >
      <MessagesSquare className="h-3.5 w-3.5 shrink-0 text-accent" />
      <span className="flex-1 font-medium text-foreground/90">Ask the community</span>
      <ArrowRight className="h-3 w-3 shrink-0 text-accent/70 transition-transform group-hover:translate-x-0.5" />
    </Link>
  )
}
