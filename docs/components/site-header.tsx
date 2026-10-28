'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { Sun, Moon } from 'lucide-react'
import { useTheme } from 'next-themes'
import { Button } from '@/components/ui/button'
import { SearchDialog } from '@/components/search-dialog'
import { MobileNav } from '@/components/docs-sidebar'
import { GitHubStars } from '@/components/github-stars'
import { DiscordIcon, DISCORD_INVITE_URL } from '@/components/community-cta'
import { cn } from '@/lib/utils'
import { GRIT_VERSION, RELEASE_CHANNEL_LABEL, releaseChannel } from '@/config/site'

export function SiteHeader() {
  const { theme, setTheme } = useTheme()
  const [mounted, setMounted] = useState(false)
  const pathname = usePathname() ?? ''

  // A navbar item is active when the current path is the item or any nested
  // route under it (e.g., /courses, /courses/concepts/welcome/what-is-grit
  // both highlight "Courses").
  const isActive = (href: string) => pathname === href || pathname.startsWith(href + '/')

  useEffect(() => setMounted(true), [])

  return (
    <header className="sticky top-0 z-50 h-16 w-full border-b border-[#1a1d22] bg-bg/80 backdrop-blur-[14px]">
      <div className="container flex h-16 max-w-screen-2xl items-center px-4 sm:px-8">
        {/* Mobile menu */}
        <MobileNav />

        {/* Logo */}
        <div className="mr-6 lg:mr-8 flex items-center">
          <Link href="/" className="flex items-center gap-2.5 group">
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src="/grit_logo.png"
              alt="Grit"
              className="h-7 w-7 rounded-md transition-transform group-hover:scale-105"
            />
            <span className="font-semibold text-[15px] tracking-tight text-foreground">
              Grit
            </span>
            {/* The channel is shown only when it is not stable, because a
                badge saying "stable" on every page is noise, and its absence
                has to mean stable for that to work. */}
            <span
              className="hidden sm:inline-flex items-center rounded-full bg-accent/10 border border-accent/20 px-2 py-0.5 text-[10px] font-mono font-medium text-accent/80"
              title={`Grit ${GRIT_VERSION}, a ${RELEASE_CHANNEL_LABEL[releaseChannel()]} release`}
            >
              v{GRIT_VERSION}
            </span>
            {releaseChannel() !== 'stable' && (
              <Link
                href="/docs/stability#release-channels"
                className="hidden sm:inline-flex items-center rounded-full bg-amber-500/15 border border-amber-500/30 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wide text-amber-600 dark:text-amber-400"
              >
                {RELEASE_CHANNEL_LABEL[releaseChannel()]}
              </Link>
            )}
          </Link>
        </div>

        {/* Navigation */}
        <nav className="hidden md:flex items-center gap-1">
          {/* Five destinations.
              It was nine, then seven, and it still read as a toolbar rather
              than a navigation bar: by the time a reader has scanned seven
              labels and six coloured icons they have stopped looking for the
              one they wanted. Sponsors, Courses, YouTube and WhatsApp are all
              in the footer and in the docs sidebar, which is where somebody
              goes looking for them. Discord stays as a single icon because it
              is the one place to ask a question, and Build with AI folds into
              Docs, where the page itself lives. */}
          {[
            { label: 'Docs', href: '/docs' },
            { label: 'Builder', href: '/builder' },
            { label: 'Grit UI', href: 'https://ui.gritframework.dev', external: true },
            { label: 'Showcase', href: '/showcase' },
            { label: 'Blog', href: '/blog' },
          ].map((item) => {
            // An off-domain entry can never be the active route, and running it
            // through isActive() would prefix-match "https://…" against the
            // pathname for no reason.
            const active = !item.external && isActive(item.href)
            return (
              <Link
                key={item.label}
                href={item.href}
                {...(item.external ? { target: '_blank', rel: 'noreferrer' } : {})}
                className={cn(
                  'relative px-3 py-1.5 text-[14px] transition-colors',
                  active ? 'text-fg' : 'text-fg-muted hover:text-fg',
                )}
              >
                {item.label}
                {active && (
                  <span className="absolute left-3 right-3 -bottom-0.5 h-[2px] rounded-full bg-accent" />
                )}
              </Link>
            )
          })}
        </nav>

        {/* Right side */}
        <div className="flex flex-1 items-center justify-end gap-1">
          <SearchDialog />

          {mounted && (
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8 text-fg-muted hover:text-fg hover:bg-surface-2"
              onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
            >
              {theme === 'dark' ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
              <span className="sr-only">Toggle theme</span>
            </Button>
          )}

          <GitHubStars className="mr-0.5" />

          {/* The one community link that stays in the bar. The rest (YouTube,
              WhatsApp, Sponsors) are in the footer, where a reader looks for
              them, rather than competing with the search box. */}
          <Button
            variant="ghost"
            size="icon"
            className="hidden sm:inline-flex h-8 w-8 text-fg-muted hover:text-fg hover:bg-surface-2"
            asChild
          >
            <Link href={DISCORD_INVITE_URL} target="_blank" rel="noreferrer">
              <DiscordIcon className="h-4 w-4" />
              <span className="sr-only">Join the Discord server</span>
            </Link>
          </Button>

          <Link
            href="/docs/getting-started/quick-start"
            className="g-btn g-btn--primary g-btn--sm ml-1 hidden md:inline-flex"
          >
            Get started
          </Link>
        </div>
      </div>
    </header>
  )
}
