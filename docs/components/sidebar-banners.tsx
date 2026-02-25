'use client'

import Link from 'next/link'
import { ExternalLink, Heart, Users } from 'lucide-react'

export function SidebarBanners() {
  return (
    <div className="mt-6 space-y-3 px-4 pb-6">
      {/* Divider */}
      <div className="h-px bg-border/30" />

      {/* GritCMS Banner */}
      <Link
        href="https://gritcms.com"
        target="_blank"
        rel="noreferrer"
        className="group block rounded-lg border border-primary/20 bg-primary/[0.06] p-3 transition-all hover:border-primary/35 hover:bg-primary/[0.10]"
      >
        <div className="flex items-center gap-2 mb-1.5">
          <div className="flex h-5 w-5 items-center justify-center rounded bg-primary/20">
            <span className="text-[10px] font-mono font-bold text-primary">G</span>
          </div>
          <span className="text-[12px] font-semibold text-foreground/90 group-hover:text-primary transition-colors">
            GritCMS
          </span>
          <ExternalLink className="h-2.5 w-2.5 text-muted-foreground/40 ml-auto opacity-0 group-hover:opacity-100 transition-opacity" />
        </div>
        <p className="text-[11px] leading-relaxed text-muted-foreground/60">
          The self-hostable creator platform. Website builder, email marketing,
          courses, community &mdash; all in one.
        </p>
        <span className="mt-2 inline-flex items-center rounded-full bg-primary/15 px-2 py-0.5 text-[10px] font-mono font-medium text-primary/80">
          Built with Grit
        </span>
      </Link>

      {/* Hire Us Banner */}
      <Link
        href="/hire"
        className="group block rounded-lg border border-border/30 bg-card/30 p-3 transition-all hover:border-primary/20 hover:bg-card/60"
      >
        <div className="flex items-center gap-2 mb-1.5">
          <Users className="h-3.5 w-3.5 text-primary/70" />
          <span className="text-[12px] font-semibold text-foreground/90 group-hover:text-primary transition-colors">
            Hire Grit Developers
          </span>
        </div>
        <p className="text-[11px] leading-relaxed text-muted-foreground/60">
          Need a production app built with Grit? Our team ships fast, full-stack
          Go + React projects.
        </p>
      </Link>

      {/* Sponsor Banner */}
      <Link
        href="https://github.com/sponsors/MUKE-coder"
        target="_blank"
        rel="noreferrer"
        className="group flex items-center gap-2.5 rounded-lg border border-border/30 bg-card/30 px-3 py-2.5 transition-all hover:border-pink-500/20 hover:bg-pink-500/[0.04]"
      >
        <Heart className="h-3.5 w-3.5 text-pink-500/70 group-hover:text-pink-500 transition-colors" />
        <div>
          <span className="text-[12px] font-semibold text-foreground/90 group-hover:text-pink-500 transition-colors block">
            Sponsor Grit
          </span>
          <span className="text-[10px] text-muted-foreground/50">
            Support open-source development
          </span>
        </div>
      </Link>
    </div>
  )
}
