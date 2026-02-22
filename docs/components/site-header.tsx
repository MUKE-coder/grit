'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { Github, Youtube, Linkedin, Globe, Sun, Moon } from 'lucide-react'
import { useTheme } from 'next-themes'
import { Button } from '@/components/ui/button'
import { SearchDialog } from '@/components/search-dialog'

export function SiteHeader() {
  const { theme, setTheme } = useTheme()
  const [mounted, setMounted] = useState(false)

  useEffect(() => setMounted(true), [])

  return (
    <header className="sticky top-0 z-50 w-full border-b border-border/40 bg-background/60 backdrop-blur-2xl">
      <div className="container flex h-14 max-w-screen-2xl items-center px-6">
        {/* Logo */}
        <div className="mr-8 flex items-center">
          <Link href="/" className="flex items-center gap-2 group">
            <div className="flex h-7 w-7 items-center justify-center rounded-md bg-primary/15 border border-primary/20 transition-colors group-hover:bg-primary/25">
              <span className="text-primary font-mono font-bold text-xs">G</span>
            </div>
            <span className="font-semibold text-[15px] tracking-tight text-foreground/90">
              Grit
            </span>
            <span className="text-xs text-muted-foreground/50 font-mono hidden sm:inline">docs</span>
            <span className="hidden sm:inline-flex items-center rounded-full bg-primary/10 border border-primary/20 px-2 py-0.5 text-[10px] font-mono font-medium text-primary/70">
              v0.17.0
            </span>
          </Link>
        </div>

        {/* Navigation */}
        <nav className="hidden md:flex items-center gap-1">
          {[
            { label: 'Docs', href: '/docs' },
            { label: 'Tutorials', href: '/docs/tutorials/blog' },
            { label: 'Playground', href: '/playground' },
            { label: 'API Reference', href: '/docs/backend/response-format' },
          ].map((item) => (
            <Link
              key={item.label}
              href={item.href}
              className="px-3 py-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors rounded-md hover:bg-accent/50"
            >
              {item.label}
            </Link>
          ))}
        </nav>

        {/* Right side */}
        <div className="flex flex-1 items-center justify-end gap-2">
          <SearchDialog />

          {mounted && (
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8 text-muted-foreground hover:text-foreground"
              onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
            >
              {theme === 'dark' ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
              <span className="sr-only">Toggle theme</span>
            </Button>
          )}

          <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground hover:text-foreground" asChild>
            <Link href="https://www.youtube.com/@JBWEBDEVELOPER" target="_blank" rel="noreferrer">
              <Youtube className="h-4 w-4" />
              <span className="sr-only">YouTube</span>
            </Link>
          </Button>
          <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground hover:text-foreground" asChild>
            <Link href="https://www.linkedin.com/in/muke-johnbaptist-95bb82198/" target="_blank" rel="noreferrer">
              <Linkedin className="h-4 w-4" />
              <span className="sr-only">LinkedIn</span>
            </Link>
          </Button>
          <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground hover:text-foreground" asChild>
            <Link href="https://jb.desishub.com" target="_blank" rel="noreferrer">
              <Globe className="h-4 w-4" />
              <span className="sr-only">Website</span>
            </Link>
          </Button>
          <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground hover:text-foreground" asChild>
            <Link href="https://github.com/MUKE-coder/grit" target="_blank" rel="noreferrer">
              <Github className="h-4 w-4" />
              <span className="sr-only">GitHub</span>
            </Link>
          </Button>

          <Button size="sm" className="hidden md:inline-flex h-8 px-3 text-xs font-medium bg-primary/90 hover:bg-primary" asChild>
            <Link href="/docs/getting-started/quick-start">
              Get Started
            </Link>
          </Button>
        </div>
      </div>
    </header>
  )
}
