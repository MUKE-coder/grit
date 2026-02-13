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
} from 'lucide-react'
import { useState } from 'react'

interface NavItem {
  title: string
  href?: string
  icon?: React.ReactNode
  items?: NavItem[]
}

const navItems: NavItem[] = [
  {
    title: 'Getting Started',
    icon: <Rocket className="h-3.5 w-3.5" />,
    items: [
      { title: 'Introduction', href: '/docs' },
      { title: 'Philosophy & Inspiration', href: '/docs/getting-started/philosophy' },
      { title: 'Quick Start', href: '/docs/getting-started/quick-start' },
      { title: 'Installation', href: '/docs/getting-started/installation' },
      { title: 'Project Structure', href: '/docs/getting-started/project-structure' },
      { title: 'Configuration', href: '/docs/getting-started/configuration' },
      { title: 'Troubleshooting', href: '/docs/getting-started/troubleshooting' },
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
      { title: 'API Response Format', href: '/docs/backend/response-format' },
      { title: 'Migrations', href: '/docs/backend/migrations' },
      { title: 'Seeders', href: '/docs/backend/seeders' },
      { title: 'RBAC & Roles', href: '/docs/backend/rbac' },
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
      { title: 'Build a Blog', href: '/docs/tutorials/blog' },
      { title: 'Build a SaaS', href: '/docs/tutorials/saas' },
      { title: 'Build an E-Commerce', href: '/docs/tutorials/ecommerce' },
    ],
  },
  {
    title: 'For AI Assistants',
    icon: <Lightbulb className="h-3.5 w-3.5" />,
    items: [
      { title: 'LLM Skill Guide', href: '/docs/ai-skill' },
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
    <aside className="fixed top-14 z-30 hidden h-[calc(100vh-3.5rem)] w-64 shrink-0 overflow-y-auto border-r border-border/30 bg-background/80 backdrop-blur-xl py-6 lg:block">
      <nav className="space-y-1 px-4">
        {navItems.map((item) => (
          <NavSection key={item.title} item={item} />
        ))}
      </nav>
    </aside>
  )
}
