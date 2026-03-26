'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { Github, Sun, Moon } from 'lucide-react'
import { useTheme } from 'next-themes'
import { Button } from '@/components/ui/button'
import { SearchDialog } from '@/components/search-dialog'
import { MobileNav } from '@/components/docs-sidebar'

export function SiteHeader() {
  const { theme, setTheme } = useTheme()
  const [mounted, setMounted] = useState(false)

  useEffect(() => setMounted(true), [])

  return (
    <header className="sticky top-0 z-50 w-full border-b border-border/40 bg-sidebar-background/95 backdrop-blur-xl">
      <div className="container flex h-16 max-w-screen-2xl items-center px-4 sm:px-8">
        {/* Mobile menu */}
        <MobileNav />

        {/* Logo */}
        <div className="mr-6 lg:mr-8 flex items-center">
          <Link href="/" className="flex items-center gap-2.5 group">
            <div className="flex h-7 w-7 items-center justify-center rounded-md bg-sky-500/10 transition-colors group-hover:bg-sky-500/20">
              <span className="text-sky-400 font-mono font-bold text-xs">G</span>
            </div>
            <span className="font-semibold text-[15px] tracking-tight text-white/90">
              Grit
            </span>
            <span className="hidden sm:inline-flex items-center rounded-full bg-sky-400/10 border border-sky-400/20 px-2 py-0.5 text-[10px] font-mono font-medium text-sky-400/80">
              v3.5.0
            </span>
          </Link>
        </div>

        {/* Navigation */}
        <nav className="hidden md:flex items-center gap-1">
          {[
            { label: 'Docs', href: '/docs' },
            { label: 'Components', href: '/docs/admin/resources' },
            { label: 'Showcase', href: '/showcase' },
            { label: 'Blog', href: 'https://gritcms.com', external: true },
          ].map((item) => (
            <Link
              key={item.label}
              href={item.href}
              {...('external' in item ? { target: '_blank', rel: 'noreferrer' } : {})}
              className="px-3 py-1.5 text-[13px] text-slate-400 hover:text-white transition-colors"
            >
              {item.label}
            </Link>
          ))}
        </nav>

        {/* Right side */}
        <div className="flex flex-1 items-center justify-end gap-1.5">
          <SearchDialog />

          {mounted && (
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8 text-slate-400 hover:text-white hover:bg-white/5"
              onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
            >
              {theme === 'dark' ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
              <span className="sr-only">Toggle theme</span>
            </Button>
          )}

          <Button variant="ghost" size="icon" className="h-8 w-8 text-slate-400 hover:text-white hover:bg-white/5" asChild>
            <Link href="https://github.com/MUKE-coder/grit" target="_blank" rel="noreferrer">
              <Github className="h-4 w-4" />
              <span className="sr-only">GitHub</span>
            </Link>
          </Button>

          <Button size="sm" className="hidden md:inline-flex h-8 px-4 text-xs font-medium bg-sky-500 hover:bg-sky-400 text-white border-0 rounded-full" asChild>
            <Link href="/docs/getting-started/quick-start">
              Get Started
            </Link>
          </Button>
        </div>
      </div>
    </header>
  )
}
