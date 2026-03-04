'use client'

import React from "react"
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
} from 'lucide-react'
import { useState } from 'react'
import { SidebarSponsorBanner, RightSideBanners } from '@/components/sidebar-banners'

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
    title: 'Getting Started',
    icon: <Rocket className="h-3.5 w-3.5" />,
    items: [
      { title: 'Introduction', href: '/docs' },
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
      { title: 'CLI Commands', href: '/docs/concepts/cli' },
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
    title: 'Frontend (Next.js)',
    icon: <Layers className="h-3.5 w-3.5" />,
    items: [
      { title: 'Web App', href: '/docs/frontend/web-app' },
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
    title: 'Infrastructure',
    icon: <Settings className="h-3.5 w-3.5" />,
    items: [
      { title: 'Docker Setup', href: '/docs/infrastructure/docker' },
      { title: 'Docker Cheat Sheet', href: '/docs/infrastructure/docker-cheatsheet' },
      { title: 'Database & Migrations', href: '/docs/infrastructure/database' },
      { title: 'Deployment', href: '/docs/infrastructure/deployment' },
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
      { title: 'Building & Distribution', href: '/docs/desktop/building' },
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
  const [isOpen, setIsOpen] = useState(true)
  const pathname = usePathname()

  if (!item.items) {
    return (
      <Link
        href={item.href || '#'}
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

export function DocsSidebar() {
  return (
    <>
      <aside className="fixed top-14 z-30 hidden h-[calc(100vh-3.5rem)] w-64 shrink-0 overflow-y-auto border-r border-border/30 bg-background/80 backdrop-blur-xl py-6 lg:block">
        <nav className="space-y-1 px-4">
          {navItems.map((item) => (
            <NavSection key={item.title} item={item} />
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
        <SidebarSponsorBanner />
      </aside>
      <RightSideBanners />
    </>
  )
}
