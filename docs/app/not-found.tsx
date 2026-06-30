import Link from 'next/link'
import type { Metadata } from 'next'
import { ArrowRight, BookOpen, Rocket, GraduationCap, Box, Sparkles, FileText, LayoutGrid, Home } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'

export const metadata: Metadata = {
  title: 'Page not found',
  description: 'That page does not exist — here are some useful places to go instead.',
  robots: { index: false, follow: true },
}

const links = [
  { icon: Rocket, title: 'Quick Start', desc: 'Scaffold a full-stack app in 5 minutes.', href: '/docs/getting-started/quick-start' },
  { icon: BookOpen, title: 'Documentation', desc: 'Guides for backend, admin, frontend, and batteries.', href: '/docs' },
  { icon: GraduationCap, title: 'Courses', desc: 'Self-paced paths and focused tutorials.', href: '/courses' },
  { icon: Box, title: 'Tech Kits', desc: 'Single, double, triple, API, mobile, desktop.', href: '/docs/tech-kits' },
  { icon: Sparkles, title: 'The Pitch', desc: 'Why Grit — one command to a production app.', href: '/pitch' },
  { icon: LayoutGrid, title: 'Showcase', desc: 'Real apps built with the framework.', href: '/showcase' },
  { icon: FileText, title: 'Changelog', desc: 'What shipped in each release.', href: '/docs/changelog' },
  { icon: Home, title: 'Home', desc: 'Back to the start.', href: '/' },
]

export default function NotFound() {
  return (
    <div className="relative min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-24 md:py-32">
        {/* Header */}
        <div className="text-center mb-14">
          <div className="font-display text-[6rem] md:text-[9rem] font-bold leading-none tracking-tight text-foreground/10 select-none">
            404
          </div>
          <span className="tag-mono text-primary -mt-6 block mb-5">Page not found</span>
          <h1 className="font-display text-3xl md:text-4xl font-bold tracking-tight text-foreground mb-4">
            This page took a different architecture mode
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed max-w-xl mx-auto">
            The page you&apos;re looking for doesn&apos;t exist, moved, or never shipped. Try one
            of these instead — or press <kbd className="px-1.5 py-0.5 rounded border border-border bg-card text-xs font-mono">⌘ K</kbd> to
            search the docs.
          </p>
          <div className="mt-8 flex items-center justify-center gap-3">
            <Link
              href="/"
              className="group inline-flex items-center gap-2 h-11 px-6 rounded-full bg-primary text-primary-foreground font-semibold text-sm hover:bg-primary/90 transition-colors glow-primary-sm"
            >
              Back home
              <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-0.5" />
            </Link>
            <Link
              href="/docs"
              className="inline-flex items-center h-11 px-6 rounded-full border border-border bg-card/60 text-foreground font-medium text-sm hover:bg-accent/40 transition-colors"
            >
              Read the docs
            </Link>
          </div>
        </div>

        {/* Useful links — bordered grid */}
        <div className="grid sm:grid-cols-2 rounded-xl border border-foreground/15 overflow-hidden">
          {links.map(({ icon: Icon, title, desc, href }) => (
            <Link
              key={href}
              href={href}
              className="group border-b border-r border-foreground/15 p-5 hover:bg-card/50 transition-colors"
            >
              <div className="flex items-center gap-2.5 mb-1.5">
                <Icon className="h-4 w-4 text-primary" />
                <span className="font-semibold text-foreground text-sm group-hover:text-primary transition-colors">{title}</span>
                <ArrowRight className="h-3.5 w-3.5 text-muted-foreground/50 ml-auto opacity-0 group-hover:opacity-100 transition-opacity" />
              </div>
              <p className="text-sm text-muted-foreground leading-relaxed">{desc}</p>
            </Link>
          ))}
        </div>
      </main>
    </div>
  )
}
