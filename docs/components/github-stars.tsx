'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { Github, Star } from 'lucide-react'
import { siteConfig } from '@/config/site'

// Compact "GitHub · ★ count" pill for the site header. Reads the cached
// /api/github-stars route on mount; if that fails it gracefully degrades
// to an icon-only link (no layout shift, no error surfaced to the user).
function formatStars(n: number): string {
  if (n >= 1000) return (n / 1000).toFixed(n >= 10000 ? 0 : 1).replace(/\.0$/, '') + 'k'
  return String(n)
}

export function GitHubStars({ className = '' }: { className?: string }) {
  const [stars, setStars] = useState<number | null>(null)

  useEffect(() => {
    let alive = true
    fetch('/api/github-stars')
      .then((r) => r.json())
      .then((d) => {
        if (alive && typeof d?.stars === 'number') setStars(d.stars)
      })
      .catch(() => {})
    return () => {
      alive = false
    }
  }, [])

  return (
    <Link
      href={siteConfig.github}
      target="_blank"
      rel="noreferrer"
      aria-label={stars === null ? 'Star Grit on GitHub' : `Grit has ${stars} stars on GitHub`}
      className={`group inline-flex items-center gap-1.5 h-8 rounded-full border border-border/60 bg-background/40 pl-2.5 pr-2.5 text-muted-foreground transition-colors hover:text-foreground hover:border-primary/40 ${className}`}
    >
      <Github className="h-4 w-4" />
      <span className="flex items-center gap-1 text-xs font-medium tabular-nums">
        <Star className="h-3 w-3 fill-amber-400 text-amber-400 transition-transform group-hover:scale-110" />
        {stars === null ? <span className="opacity-60">Star</span> : formatStars(stars)}
      </span>
    </Link>
  )
}
