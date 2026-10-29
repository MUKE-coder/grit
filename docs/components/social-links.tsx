import Link from 'next/link'
import { Github, Linkedin, Youtube } from 'lucide-react'

import { siteConfig } from '@/config/site'
import { DiscordIcon, WhatsAppIcon } from '@/components/community-cta'

/**
 * Where to find Grit, in one row.
 *
 * Every URL comes from config/site.ts. None is written here, because the
 * Discord invite already learned that lesson: an invite pasted into a
 * component is an invite nobody can find when it expires.
 *
 * Discord and WhatsApp have no lucide glyph, so community-cta.tsx holds the
 * official marks and this imports them rather than carrying a second copy.
 *
 * `label` is what a screen reader announces and what the tooltip shows. Five
 * icon links in a row with no accessible name is a row of five "link"
 * announcements, which is the most common way a footer like this fails.
 */
const LINKS = [
  { label: 'GitHub', href: siteConfig.github, Icon: Github },
  { label: 'YouTube', href: siteConfig.youtube, Icon: Youtube },
  { label: 'LinkedIn', href: siteConfig.linkedin, Icon: Linkedin },
  { label: 'Discord', href: siteConfig.discord, Icon: DiscordIcon },
  { label: 'WhatsApp', href: siteConfig.whatsapp, Icon: WhatsAppIcon },
] as const

export function SocialLinks({
  className = '',
  size = 'default',
}: {
  className?: string
  /** 'compact' for a footer rail, 'default' for a section of its own. */
  size?: 'default' | 'compact'
}) {
  const box = size === 'compact' ? 'h-8 w-8' : 'h-11 w-11'
  const glyph = size === 'compact' ? 'h-[15px] w-[15px]' : 'h-[18px] w-[18px]'

  return (
    <ul className={`flex flex-wrap items-center justify-center gap-1.5 ${className}`}>
      {LINKS.map(({ label, href, Icon }) => (
        <li key={label}>
          <Link
            href={href}
            target="_blank"
            rel="noreferrer"
            aria-label={label}
            title={label}
            className={`group grid ${box} place-items-center rounded-xl border border-border bg-surface-1 text-fg-muted transition-colors hover:border-accent-line hover:bg-surface-2 hover:text-accent focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent focus-visible:ring-offset-2 focus-visible:ring-offset-bg`}
          >
            <Icon className={glyph} />
          </Link>
        </li>
      ))}
    </ul>
  )
}

/**
 * The same links with the sentence that makes somebody click one.
 *
 * A bare row of icons at the end of a landing page is decoration. Saying
 * what is on the other side of them is the difference between a logo and an
 * invitation.
 */
export function SocialStrip({ className = '' }: { className?: string }) {
  return (
    <div className={className}>
      <p className="text-[14px] leading-[1.7] text-fg-muted">
        Questions, releases and what people are building: the Discord and the
        WhatsApp group are where the answers happen.
      </p>
      <SocialLinks className="mt-4" />
    </div>
  )
}
