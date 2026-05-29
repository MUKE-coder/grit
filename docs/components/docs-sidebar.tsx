'use client'

import React, { useRef, useEffect } from "react"
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { cn } from '@/lib/utils'
import {
  ChevronRight,
  Rocket,
  Box,
  Server,
  Database,
  Shield,
  Palette,
  BookOpen,
  Layers,
  Settings,
  Lightbulb,
  Wand2,
  FileText,
  GraduationCap,
  Download,
  Monitor,
  Menu,
} from 'lucide-react'
import { useState } from 'react'
// Banners removed for cleaner Tailwind-style sidebar
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet'
import { Button } from '@/components/ui/button'

interface NavItem {
  title: string
  href?: string
  icon?: React.ReactNode
  items?: NavItem[]
}

const navItems: NavItem[] = [
  {
    title: 'Prerequisites',
    icon: <GraduationCap className="h-3.5 w-3.5" />,
    items: [
      { title: 'Go for Grit Developers', href: '/docs/prerequisites/golang' },
      { title: 'Go Playground', href: '/playground' },
      { title: 'Next.js & React', href: '/docs/prerequisites/nextjs' },
      { title: 'Docker', href: '/docs/prerequisites/docker' },
    ],
  },
  {
    title: 'Demo',
    icon: <Rocket className="h-3.5 w-3.5" />,
    items: [
      { title: 'Demo Application', href: '/docs/demo' },
    ],
  },
  {
    title: 'Tech Kits',
    icon: <Box className="h-3.5 w-3.5" />,
    items: [
      { title: 'All Tech Kits', href: '/docs/tech-kits' },
      { title: 'AI Integration Helper', href: '/docs/ai-integration' },
      { title: 'Single (Next.js)', href: '/docs/tech-kits/single' },
      { title: 'Single + Vite', href: '/docs/tech-kits/single-vite' },
      { title: 'Double — Web + API', href: '/docs/tech-kits/double' },
      { title: 'Triple — Web + Admin + API', href: '/docs/tech-kits/triple' },
      { title: 'API only', href: '/docs/tech-kits/api' },
      { title: 'Mobile (Expo)', href: '/docs/tech-kits/mobile' },
      { title: 'Desktop (Wails)', href: '/docs/tech-kits/desktop' },
    ],
  },
  {
    title: 'Getting Started',
    icon: <Rocket className="h-3.5 w-3.5" />,
    items: [
      { title: 'Introduction', href: '/docs' },
      { title: 'Stack Selector — Pick a Combo', href: '/docs/stack-selector' },
      { title: 'Philosophy & Inspiration', href: '/docs/getting-started/philosophy' },
      { title: 'Quick Start', href: '/docs/getting-started/quick-start' },
      { title: 'Your First App', href: '/docs/tutorials/contact-app' },
      { title: 'Installation', href: '/docs/getting-started/installation' },
      { title: 'Project Structure', href: '/docs/getting-started/project-structure' },
      { title: 'Configuration', href: '/docs/getting-started/configuration' },
      { title: 'Troubleshooting', href: '/docs/getting-started/troubleshooting' },
      { title: 'Create without Docker', href: '/docs/getting-started/create-without-docker' },
      { title: 'CLI Cheatsheet', href: '/docs/getting-started/cli-cheatsheet' },
    ],
  },
  {
    title: 'Core Concepts',
    icon: <Box className="h-3.5 w-3.5" />,
    items: [
      { title: 'Architecture Overview', href: '/docs/concepts/architecture' },
      { title: 'Architecture Modes', href: '/docs/concepts/architecture-modes' },
      { title: 'Triple (Web+Admin+API)', href: '/docs/concepts/architecture-modes/triple' },
      { title: 'Double (Web+API)', href: '/docs/concepts/architecture-modes/double' },
      { title: 'Single (One Binary)', href: '/docs/concepts/architecture-modes/single' },
      { title: 'API Only', href: '/docs/concepts/architecture-modes/api-only' },
      { title: 'Mobile (API+Expo)', href: '/docs/concepts/architecture-modes/mobile' },
      { title: 'Multi-Client (+ Desktop)', href: '/docs/concepts/architecture-modes/multi-client' },
      { title: 'CLI Commands', href: '/docs/concepts/cli' },
      { title: 'New CLI Commands', href: '/docs/concepts/cli-commands' },
      { title: 'Code Generation', href: '/docs/concepts/code-generation' },
      { title: 'Type System', href: '/docs/concepts/type-system' },
      { title: 'Naming Conventions', href: '/docs/concepts/naming-conventions' },
      { title: 'Style Variants', href: '/docs/concepts/styles' },
      { title: 'Performance', href: '/docs/concepts/performance' },
    ],
  },
  {
    title: 'Backend (Go API)',
    icon: <Server className="h-3.5 w-3.5" />,
    items: [
      { title: 'Models & Database', href: '/docs/backend/models' },
      { title: 'Handlers', href: '/docs/backend/handlers' },
      { title: 'Services', href: '/docs/backend/services' },
      { title: 'Middleware', href: '/docs/backend/middleware' },
      { title: 'Authentication', href: '/docs/backend/authentication' },
      { title: 'Social Login (OAuth2)', href: '/docs/backend/oauth' },
      { title: 'API Response Format', href: '/docs/backend/response-format' },
      { title: 'Migrations', href: '/docs/backend/migrations' },
      { title: 'Seeders', href: '/docs/backend/seeders' },
      { title: 'RBAC & Roles', href: '/docs/backend/rbac' },
      { title: 'API Documentation', href: '/docs/backend/api-docs' },
      { title: 'Pulse (Observability)', href: '/docs/backend/pulse' },
    ],
  },
  {
    title: 'Security & Testing',
    icon: <Shield className="h-3.5 w-3.5" />,
    items: [
      { title: 'Security Guide', href: '/docs/security' },
      { title: 'Performance & Pentest Testing', href: '/docs/testing' },
    ],
  },
  {
    title: 'Admin Panel',
    icon: <Shield className="h-3.5 w-3.5" />,
    items: [
      { title: 'Admin Overview', href: '/docs/admin/overview' },
      { title: 'Resource Definitions', href: '/docs/admin/resources' },
      { title: 'DataTable', href: '/docs/admin/datatable' },
      { title: 'Form Builder', href: '/docs/admin/forms' },
      { title: 'Multi-Step Forms', href: '/docs/admin/multi-step-forms' },
      { title: 'Standalone Usage', href: '/docs/admin/standalone-usage' },
      { title: 'Relationships', href: '/docs/admin/relationships' },
      { title: 'Dashboard & Widgets', href: '/docs/admin/widgets' },
    ],
  },
  {
    title: 'Frontend',
    icon: <Layers className="h-3.5 w-3.5" />,
    items: [
      { title: 'Web App (Next.js)', href: '/docs/frontend/web-app' },
      { title: 'TanStack Router (Vite)', href: '/docs/frontend/tanstack-router' },
      { title: 'React Query Hooks', href: '/docs/frontend/hooks' },
      { title: 'Shared Package', href: '/docs/frontend/shared-package' },
    ],
  },
  {
    title: 'Batteries',
    icon: <Database className="h-3.5 w-3.5" />,
    items: [
      { title: 'File Storage', href: '/docs/batteries/storage' },
      { title: 'Email System', href: '/docs/batteries/email' },
      { title: 'Background Jobs', href: '/docs/batteries/jobs' },
      { title: 'Cron Scheduler', href: '/docs/batteries/cron' },
      { title: 'Redis Caching', href: '/docs/batteries/caching' },
      { title: 'AI Integration', href: '/docs/batteries/ai' },
      { title: 'Security (Sentinel)', href: '/docs/batteries/security' },
    ],
  },
  {
    title: 'Plugins',
    icon: <Layers className="h-3.5 w-3.5" />,
    items: [
      { title: 'Overview', href: '/docs/plugins' },
    ],
  },
  {
    title: 'Infrastructure',
    icon: <Settings className="h-3.5 w-3.5" />,
    items: [
      { title: 'Docker Setup', href: '/docs/infrastructure/docker' },
      { title: 'Docker Cheat Sheet', href: '/docs/infrastructure/docker-cheatsheet' },
      { title: 'Database & Migrations', href: '/docs/infrastructure/database' },
      { title: 'Deployment', href: '/docs/infrastructure/deployment' },
      { title: 'Deploy Command', href: '/docs/infrastructure/deploy-command' },
      { title: 'Deploy with Dokploy', href: '/docs/infrastructure/dokploy' },
    ],
  },
  {
    title: 'Design System',
    icon: <Palette className="h-3.5 w-3.5" />,
    items: [
      { title: 'Theme & Colors', href: '/docs/design/theme' },
    ],
  },
  {
    title: 'Tutorials',
    icon: <BookOpen className="h-3.5 w-3.5" />,
    items: [
      { title: 'Your First App', href: '/docs/tutorials/contact-app' },
      { title: 'Learn Grit Step by Step', href: '/docs/tutorials/learn' },
      { title: 'Build a Blog', href: '/docs/tutorials/blog' },
      { title: 'Build a SaaS', href: '/docs/tutorials/saas' },
      { title: 'Build an E-Commerce', href: '/docs/tutorials/ecommerce' },
      { title: 'Build a Product Catalog', href: '/docs/tutorials/product-catalog' },
    ],
  },
  {
    title: 'Desktop (Wails)',
    icon: <Monitor className="h-3.5 w-3.5" />,
    items: [
      { title: 'Overview', href: '/docs/desktop' },
      { title: 'Getting Started', href: '/docs/desktop/getting-started' },
      { title: 'Your First Desktop App', href: '/docs/desktop/first-app' },
      { title: 'Build a POS App', href: '/docs/desktop/pos-app' },
      { title: 'Resource Generation', href: '/docs/desktop/resource-generation' },
      { title: 'Offline-First Apps', href: '/docs/desktop/offline' },
      { title: 'Building & Distribution', href: '/docs/desktop/building' },
      { title: '20 Project Ideas', href: '/docs/desktop/project-ideas' },
      { title: 'LLM Reference', href: '/docs/desktop/llm-reference' },
    ],
  },
  {
    title: 'AI Workflows',
    icon: <Wand2 className="h-3.5 w-3.5" />,
    items: [
      { title: 'Using Grit with Claude', href: '/docs/ai-workflows/claude' },
      { title: 'Using Grit with Antigravity', href: '/docs/ai-workflows/antigravity' },
    ],
  },
  {
    title: 'For AI Assistants',
    icon: <Lightbulb className="h-3.5 w-3.5" />,
    items: [
      { title: 'LLM Skill Guide', href: '/docs/ai-skill' },
      { title: 'Complete LLM Reference', href: '/docs/ai-skill/llm-guide' },
    ],
  },
  {
    title: 'Changelog',
    icon: <FileText className="h-3.5 w-3.5" />,
    items: [
      { title: 'Release History', href: '/docs/changelog' },
    ],
  },
]

function NavSection({ item }: { item: NavItem }) {
  const pathname = usePathname()

  if (!item.items) {
    return (
      <Link
        href={item.href || '#'}
        className={cn(
          'block px-3 py-1 text-sm transition-colors',
          pathname === item.href
            ? 'text-primary font-medium'
            : 'text-muted-foreground hover:text-foreground'
        )}
      >
        {item.title}
      </Link>
    )
  }

  return (
    <div className="mb-5">
      <h5 className="px-3 mb-2 text-[11px] font-semibold tracking-wider text-foreground/50 uppercase">
        {item.title}
      </h5>
      <div className="space-y-0.5">
        {item.items.map((subItem) => (
          <Link
            key={subItem.href}
            href={subItem.href || '#'}
            data-active={pathname === subItem.href ? 'true' : undefined}
            className={cn(
              'block px-3 py-1 text-sm transition-colors border-l-2',
              pathname === subItem.href
                ? 'text-primary font-medium border-primary'
                : 'text-muted-foreground hover:text-foreground border-transparent hover:border-muted-foreground/30'
            )}
          >
            {subItem.title}
          </Link>
        ))}
      </div>
    </div>
  )
}

export function DocsSidebar() {
  const sidebarRef = useRef<HTMLElement>(null)
  const pathname = usePathname()

  useEffect(() => {
    const el = sidebarRef.current?.querySelector('[data-active="true"]')
    if (el) {
      el.scrollIntoView({ block: 'nearest', behavior: 'instant' })
    }
  }, [pathname])

  return (
    <aside ref={sidebarRef} className="fixed top-16 z-30 hidden h-[calc(100vh-4rem)] w-64 shrink-0 overflow-y-auto border-r border-border/40 bg-sidebar-background py-8 lg:block">
      <nav className="px-2">
        {navItems.map((item) => (
          <NavSection key={item.title} item={item} />
        ))}
      </nav>
    </aside>
  )
}

/* ── Mobile Navigation Sheet ─────────────────────────────────── */

function MobileNavSection({ item, onNavigate }: { item: NavItem; onNavigate: () => void }) {
  const [isOpen, setIsOpen] = useState(true)
  const pathname = usePathname()

  if (!item.items) {
    return (
      <Link
        href={item.href || '#'}
        onClick={onNavigate}
        className={cn(
          'flex items-center gap-2.5 rounded-md px-2.5 py-1.5 text-[13px] transition-all',
          pathname === item.href
            ? 'bg-primary/10 text-primary font-medium'
            : 'text-muted-foreground hover:bg-accent/50 hover:text-foreground'
        )}
      >
        {item.icon}
        {item.title}
      </Link>
    )
  }

  return (
    <div className="mb-1">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex w-full items-center justify-between rounded-md px-2.5 py-1.5 text-[13px] font-medium text-foreground/80 hover:bg-accent/30 transition-colors"
      >
        <div className="flex items-center gap-2.5">
          <span className="text-muted-foreground/60">{item.icon}</span>
          {item.title}
        </div>
        <ChevronRight
          className={cn(
            'h-3 w-3 text-muted-foreground/40 transition-transform duration-200',
            isOpen && 'rotate-90'
          )}
        />
      </button>
      {isOpen && (
        <div className="mt-0.5 space-y-0.5 ml-3 pl-3 border-l border-border/30">
          {item.items.map((subItem) => (
            <Link
              key={subItem.href}
              href={subItem.href || '#'}
              onClick={onNavigate}
              className={cn(
                'block rounded-md px-2.5 py-1.5 text-[13px] transition-all',
                pathname === subItem.href
                  ? 'text-primary font-medium bg-primary/5'
                  : 'text-muted-foreground/70 hover:text-foreground hover:bg-accent/30'
              )}
            >
              {subItem.title}
            </Link>
          ))}
        </div>
      )}
    </div>
  )
}

export function MobileNav() {
  const [open, setOpen] = useState(false)

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button
          variant="ghost"
          size="icon"
          className="h-8 w-8 text-muted-foreground hover:text-foreground lg:hidden"
        >
          <Menu className="h-4 w-4" />
          <span className="sr-only">Toggle navigation menu</span>
        </Button>
      </SheetTrigger>
      <SheetContent side="left" className="w-72 p-0 bg-background border-border/30">
        <SheetHeader className="px-4 py-3 border-b border-border/30">
          <SheetTitle className="flex items-center gap-2 text-left">
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img src="/grit_logo.png" alt="Grit" className="h-6 w-6 rounded-md" />
            <span className="text-sm font-semibold tracking-tight">Grit Docs</span>
          </SheetTitle>
        </SheetHeader>
        <div className="overflow-y-auto h-[calc(100vh-3.5rem)] py-4">
          {/* Top-level nav links (Docs, Showcase, etc.) */}
          <div className="px-4 mb-4 space-y-1">
            {[
              { label: 'Docs', href: '/docs' },
              { label: 'Showcase', href: '/showcase' },
              { label: 'GritCMS', href: 'https://gritcms.com' },
              { label: 'Hire Us', href: '/hire' },
            ].map((item) => (
              <Link
                key={item.label}
                href={item.href}
                onClick={() => setOpen(false)}
                className="block rounded-md px-2.5 py-1.5 text-sm text-muted-foreground hover:text-foreground hover:bg-accent/50 transition-colors"
              >
                {item.label}
              </Link>
            ))}
          </div>
          <div className="mx-4 mb-3 border-t border-border/30" />
          {/* Full sidebar nav */}
          <nav className="space-y-1 px-4">
            {navItems.map((item) => (
              <MobileNavSection key={item.title} item={item} onNavigate={() => setOpen(false)} />
            ))}
          </nav>
          {/* Handbook download */}
          <div className="px-4 mt-4 mb-2">
            <a
              href="https://14j7oh8kso.ufs.sh/f/HLxTbDBCDLwfeHHJl34ZKSqNhOvVj6p9rg3Icmo05TAEwQ4a"
              target="_blank"
              rel="noopener noreferrer"
            >
              <div className="flex items-center gap-2 rounded-lg border border-primary/25 bg-primary/8 px-3 py-2.5 text-xs font-medium text-primary/80 hover:bg-primary/15 hover:text-primary transition-colors cursor-pointer">
                <Download className="h-3.5 w-3.5 shrink-0" />
                <span>Download Handbook PDF</span>
              </div>
            </a>
          </div>
        </div>
      </SheetContent>
    </Sheet>
  )
}
