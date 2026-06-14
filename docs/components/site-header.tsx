'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { Github, Youtube, Heart, Sun, Moon } from 'lucide-react'
import { useTheme } from 'next-themes'
import { Button } from '@/components/ui/button'
import { SearchDialog } from '@/components/search-dialog'
import { MobileNav } from '@/components/docs-sidebar'
import { cn } from '@/lib/utils'
import { GRIT_VERSION } from '@/config/site'

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
    <header className="sticky top-0 z-50 w-full border-b border-border/40 bg-sidebar-background/95 backdrop-blur-xl">
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
            <span className="hidden sm:inline-flex items-center rounded-full bg-primary/10 border border-primary/20 px-2 py-0.5 text-[10px] font-mono font-medium text-primary/80">
              v{GRIT_VERSION}
            </span>
          </Link>
        </div>

        {/* Navigation */}
        <nav className="hidden md:flex items-center gap-1">
          {[
            { label: 'Docs', href: '/docs' },
            { label: 'Tech Kits', href: '/docs/tech-kits' },
            { label: 'AI Integration', href: '/docs/ai-integration', highlight: true },
            { label: 'Courses', href: '/courses' },
            { label: 'Changelog', href: '/docs/changelog' },
            { label: 'Showcase', href: '/showcase' },
          ].map((item) => {
            const active = isActive(item.href)
            return (
              <Link
                key={item.label}
                href={item.href}
                className={cn(
                  'relative px-3 py-1.5 text-[13px] transition-colors',
                  active
                    ? 'font-semibold text-primary'
                    : item.highlight
                      ? 'font-medium text-primary hover:text-primary/80'
                      : 'text-muted-foreground hover:text-foreground',
                )}
              >
                {item.label}
                {active && (
                  <span className="absolute left-3 right-3 -bottom-0.5 h-[2px] rounded-full bg-primary" />
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
              className="h-8 w-8 text-muted-foreground hover:text-foreground hover:bg-accent/50"
              onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
            >
              {theme === 'dark' ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
              <span className="sr-only">Toggle theme</span>
            </Button>
          )}

          <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground hover:text-foreground hover:bg-accent/50" asChild>
            <Link href="https://github.com/MUKE-coder/grit" target="_blank" rel="noreferrer">
              <Github className="h-4 w-4" />
              <span className="sr-only">GitHub</span>
            </Link>
          </Button>

          <Button variant="ghost" size="icon" className="hidden sm:inline-flex h-8 w-8 text-muted-foreground hover:text-foreground hover:bg-accent/50" asChild>
            <Link href="https://www.youtube.com/@GritFramework" target="_blank" rel="noreferrer">
              <Youtube className="h-4 w-4" />
              <span className="sr-only">YouTube</span>
            </Link>
          </Button>

          <Button variant="ghost" size="icon" className="hidden sm:inline-flex h-8 w-8 text-pink-500/70 hover:text-pink-500 hover:bg-pink-500/10" asChild>
            <Link href="/donate">
              <Heart className="h-4 w-4" />
              <span className="sr-only">Sponsor</span>
            </Link>
          </Button>

          <Button size="sm" className="hidden md:inline-flex h-8 px-4 text-xs font-medium bg-primary hover:bg-primary/80 text-primary-foreground border-0 rounded-full" asChild>
            <Link href="/docs/getting-started/quick-start">
              Get Started
            </Link>
          </Button>
        </div>
      </div>
    </header>
  )
}
