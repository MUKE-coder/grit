import Link from 'next/link'
import { Github, Globe, Star, Heart, Plus } from 'lucide-react'
import { cn } from '@/lib/utils'
import {
  type Sponsor,
  type SponsorTierId,
  activeSponsors,
  formatContribution,
  formatSince,
  getTier,
  specialSponsors,
} from '@/config/sponsors'

/**
 * Presentation layer for sponsorship. All data comes from config/sponsors.ts —
 * these components never hardcode a name, price or perk.
 */

const TIER_STYLES: Record<SponsorTierId, string> = {
  supporter: 'text-muted-foreground border-border/60',
  backer: 'text-sky-400 border-sky-500/30 bg-sky-500/5',
  special: 'text-amber-400 border-amber-500/30 bg-amber-500/5',
  partner: 'text-fuchsia-400 border-fuchsia-500/30 bg-fuchsia-500/5',
}

/** The SPECIAL · SINCE JUL 2026 chip along the top of a sponsor card. */
export function TierBadge({
  tier,
  since,
  className,
}: {
  tier: SponsorTierId
  since?: string
  className?: string
}) {
  const isTop = tier === 'special' || tier === 'partner'
  return (
    <div className={cn('flex items-center justify-between gap-2', className)}>
      {isTop ? (
        <Star className="h-3.5 w-3.5 shrink-0 fill-current text-amber-400" />
      ) : (
        <Plus className="h-3 w-3 shrink-0 text-muted-foreground/40" />
      )}
      <div className="flex items-center gap-2 font-mono text-[10px] uppercase tracking-wider">
        <span className={cn('rounded border px-1.5 py-0.5', TIER_STYLES[tier])}>
          {getTier(tier).name}
        </span>
        {since && (
          <span className="text-muted-foreground/50">since {formatSince(since)}</span>
        )}
      </div>
    </div>
  )
}

/** Logo when we have one, monogram when we don't — never a broken image. */
function SponsorAvatar({
  sponsor,
  size = 'md',
}: {
  sponsor: Sponsor
  size?: 'sm' | 'md'
}) {
  const dim = size === 'sm' ? 'h-10 w-10' : 'h-16 w-16'
  if (sponsor.avatar) {
    return (
      // eslint-disable-next-line @next/next/no-img-element -- sponsor logos are
      // arbitrary remote hosts; next/image would need every one in remotePatterns.
      <img
        src={sponsor.avatar}
        alt={sponsor.name}
        className={cn(dim, 'shrink-0 rounded-lg border border-border/40 object-contain')}
      />
    )
  }
  return (
    <div
      className={cn(
        dim,
        'flex shrink-0 items-center justify-center rounded-lg border border-border/40 bg-card font-display text-lg font-bold text-muted-foreground'
      )}
    >
      {sponsor.name.charAt(0).toUpperCase()}
    </div>
  )
}

export function SponsorCard({ sponsor }: { sponsor: Sponsor }) {
  const host = sponsor.website?.replace(/^https?:\/\//, '').replace(/\/$/, '')
  return (
    <div className="card-grit flex flex-col rounded-xl border border-border/40 bg-card/50 p-4">
      <TierBadge tier={sponsor.tier} since={sponsor.since} className="mb-4" />

      <div className="flex flex-1 items-start gap-4">
        <SponsorAvatar sponsor={sponsor} />
        <div className="min-w-0 flex-1">
          <h3 className="truncate font-display text-base font-bold text-foreground">
            {sponsor.name}
          </h3>
          <p className="mt-0.5 font-mono text-sm text-primary">
            {formatContribution(sponsor)}
          </p>
          <div className="mt-2 space-y-1">
            {sponsor.github && (
              <Link
                href={`https://github.com/${sponsor.github}`}
                target="_blank"
                rel="noreferrer"
                className="flex items-center gap-1.5 text-xs text-muted-foreground transition-colors hover:text-foreground"
              >
                <Github className="h-3 w-3 shrink-0" />
                <span className="truncate">{sponsor.github}</span>
              </Link>
            )}
            {host && (
              <Link
                href={sponsor.website!.startsWith('http') ? sponsor.website! : `https://${host}`}
                target="_blank"
                rel="noreferrer"
                className="flex items-center gap-1.5 text-xs text-muted-foreground transition-colors hover:text-foreground"
              >
                <Globe className="h-3 w-3 shrink-0" />
                <span className="truncate">{host}</span>
              </Link>
            )}
          </div>
        </div>
      </div>

      {sponsor.lifetime !== undefined && (
        <div className="mt-4 flex items-center justify-between border-t border-border/30 pt-3">
          <span className="font-mono text-[10px] uppercase tracking-wider text-muted-foreground/50">
            Lifetime support
          </span>
          <span className="font-mono text-sm font-semibold text-foreground">
            ${sponsor.lifetime.toLocaleString(undefined, { minimumFractionDigits: 2 })}
          </span>
        </div>
      )}
    </div>
  )
}

/**
 * Compact logo strip — the home page and docs footer. Renders nothing when the
 * roster is empty so we never ship an empty "Sponsored by" rail.
 */
export function SponsorsStrip({
  sponsors,
  className,
}: {
  sponsors: Sponsor[]
  className?: string
}) {
  if (sponsors.length === 0) return null
  return (
    <div className={cn('flex flex-wrap items-center justify-center gap-x-8 gap-y-4', className)}>
      {sponsors.map((sponsor) => {
        const href = sponsor.website
          ? sponsor.website.startsWith('http')
            ? sponsor.website
            : `https://${sponsor.website}`
          : sponsor.github
            ? `https://github.com/${sponsor.github}`
            : undefined
        const body = (
          <div className="flex items-center gap-2.5 opacity-70 transition-opacity hover:opacity-100">
            <SponsorAvatar sponsor={sponsor} size="sm" />
            <span className="font-display text-sm font-semibold text-foreground">
              {sponsor.name}
            </span>
          </div>
        )
        return href ? (
          <Link key={sponsor.name} href={href} target="_blank" rel="noreferrer">
            {body}
          </Link>
        ) : (
          <div key={sponsor.name}>{body}</div>
        )
      })}
    </div>
  )
}

/**
 * The home page block. Leads with the top tiers as cards, falls back to the
 * founding pitch while the roster is empty, so the section is never dead space.
 */
export function HomeSponsors() {
  const special = specialSponsors()
  const everyone = activeSponsors()

  if (everyone.length === 0) return <FoundingSponsorCTA />

  const rest = everyone.filter((s) => !special.includes(s))
  return (
    <div className="space-y-10">
      {special.length > 0 && (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {special.map((sponsor) => (
            <SponsorCard key={sponsor.name} sponsor={sponsor} />
          ))}
        </div>
      )}
      {rest.length > 0 && (
        <div>
          <p className="mb-6 text-center font-mono text-[10px] uppercase tracking-wider text-muted-foreground/50">
            Also sponsored by
          </p>
          <SponsorsStrip sponsors={rest} />
        </div>
      )}
    </div>
  )
}

/**
 * Shown wherever the roster is empty. With zero sponsors a grid of zeros reads
 * as neglect — this asks for the founding slot instead.
 */
export function FoundingSponsorCTA({ className }: { className?: string }) {
  return (
    <div
      className={cn(
        'relative overflow-hidden rounded-xl border border-dashed border-primary/30 bg-primary/[0.03] px-6 py-10 text-center',
        className
      )}
    >
      <div className="pointer-events-none absolute inset-0 -z-10 bg-grit-dots mask-fade-center opacity-40" />
      <Heart className="mx-auto mb-3 h-5 w-5 text-primary" />
      <h3 className="font-display text-lg font-bold text-foreground">
        This spot is open
      </h3>
      <p className="mx-auto mt-2 max-w-md text-sm text-muted-foreground">
        Grit has no sponsors yet. Be the first — your logo goes at the top of this
        page, on the home page, in the README, and inside the CLI.
      </p>
      <Link
        href="/sponsor"
        className="mt-5 inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground transition-opacity hover:opacity-90"
      >
        <Star className="h-4 w-4" />
        Become the founding sponsor
      </Link>
    </div>
  )
}
