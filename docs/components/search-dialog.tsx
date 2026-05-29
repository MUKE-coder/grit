'use client'

import { useCallback, useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import {
  CommandDialog,
  CommandInput,
  CommandList,
  CommandEmpty,
  CommandGroup,
  CommandItem,
} from '@/components/ui/command'
import {
  Rocket,
  Box,
  Server,
  Shield,
  Layers,
  Database,
  Settings,
  Palette,
  BookOpen,
  Lightbulb,
  Wand2,
  FileText,
  Monitor,
  Plug,
} from 'lucide-react'

interface SearchItem {
  title: string
  href: string
  section: string
  keywords: string
}

const searchIndex: SearchItem[] = [
  // Getting Started
  { title: 'Introduction', href: '/docs', section: 'Getting Started', keywords: 'overview what is grit start begin home introduction welcome' },
  { title: 'Stack Selector — Pick a Combo', href: '/docs/stack-selector', section: 'Getting Started', keywords: 'choose stack combo architecture which what should I pick build single double triple api mobile desktop offline web portal saas internal tool dashboard kitchen sink multi-platform decision tree recommendation capability matrix vite next nextjs tanstack' },
  { title: 'Philosophy & Inspiration', href: '/docs/getting-started/philosophy', section: 'Getting Started', keywords: 'why design decisions laravel rails inspiration manifesto principles' },
  { title: 'Quick Start', href: '/docs/getting-started/quick-start', section: 'Getting Started', keywords: 'setup first project tutorial begin hello new install run dev' },
  { title: 'Installation', href: '/docs/getting-started/installation', section: 'Getting Started', keywords: 'install go node pnpm setup requirements prerequisites cli go install brew' },
  { title: 'Project Structure', href: '/docs/getting-started/project-structure', section: 'Getting Started', keywords: 'folders files directory layout monorepo apps tree organization' },
  { title: 'Configuration', href: '/docs/getting-started/configuration', section: 'Getting Started', keywords: 'env config environment variables settings dotenv' },
  { title: 'Troubleshooting', href: '/docs/getting-started/troubleshooting', section: 'Getting Started', keywords: 'errors fix debug issues problems help broken stuck' },
  { title: 'Create without Docker', href: '/docs/getting-started/create-without-docker', section: 'Getting Started', keywords: 'no docker without docker sqlite minimal setup local' },
  { title: 'CLI Cheatsheet', href: '/docs/getting-started/cli-cheatsheet', section: 'Getting Started', keywords: 'cli commands cheatsheet reference grit new generate sync migrate seed' },

  // Core Concepts
  { title: 'Architecture Overview', href: '/docs/concepts/architecture', section: 'Core Concepts', keywords: 'design monorepo api frontend shared structure' },
  { title: 'CLI Commands', href: '/docs/concepts/cli', section: 'Core Concepts', keywords: 'grit new generate sync migrate seed add role start client server command' },
  { title: 'Code Generation', href: '/docs/concepts/code-generation', section: 'Core Concepts', keywords: 'generate resource crud scaffold model handler' },
  { title: 'Type System', href: '/docs/concepts/type-system', section: 'Core Concepts', keywords: 'typescript go types shared zod schema sync' },
  { title: 'Naming Conventions', href: '/docs/concepts/naming-conventions', section: 'Core Concepts', keywords: 'case snake camel pascal kebab naming style' },
  { title: 'Style Variants', href: '/docs/concepts/styles', section: 'Core Concepts', keywords: 'style variant theme modern minimal glass default auth dashboard layout design' },

  // Backend
  { title: 'Models & Database', href: '/docs/backend/models', section: 'Backend (Go API)', keywords: 'gorm model struct database table fields columns version belongsto manytomany uuid pk' },
  { title: 'Handlers', href: '/docs/backend/handlers', section: 'Backend (Go API)', keywords: 'api endpoint handler controller request response gin paginate list' },
  { title: 'Services', href: '/docs/backend/services', section: 'Backend (Go API)', keywords: 'business logic service layer repository getByID update delete' },
  { title: 'Middleware', href: '/docs/backend/middleware', section: 'Backend (Go API)', keywords: 'auth cors logger rate limit middleware idempotency activity audit cache gzip security headers' },
  { title: 'Authentication', href: '/docs/backend/authentication', section: 'Backend (Go API)', keywords: 'jwt login register auth token password 2fa totp oauth google github trusted devices backup codes refresh' },
  { title: 'API Response Format', href: '/docs/backend/response-format', section: 'Backend (Go API)', keywords: 'json response error pagination format api envelope respond meta data message error code apiErrorMessage' },
  { title: 'Migrations', href: '/docs/backend/migrations', section: 'Backend (Go API)', keywords: 'migrate database schema table create alter fresh automigrate column diff verbose' },
  { title: 'Seeders', href: '/docs/backend/seeders', section: 'Backend (Go API)', keywords: 'seed data demo users populate database initial fixtures' },
  { title: 'RBAC & Roles', href: '/docs/backend/rbac', section: 'Backend (Go API)', keywords: 'roles rbac admin editor user permissions access control require role middleware' },

  // Admin Panel
  { title: 'Admin Overview', href: '/docs/admin/overview', section: 'Admin Panel', keywords: 'admin panel dashboard overview filament' },
  { title: 'Resource Definitions', href: '/docs/admin/resources', section: 'Admin Panel', keywords: 'resource define crud table form config columns fields' },
  { title: 'DataTable', href: '/docs/admin/datatable', section: 'Admin Panel', keywords: 'table list sort filter search pagination columns' },
  { title: 'Form Builder', href: '/docs/admin/forms', section: 'Admin Panel', keywords: 'form create edit fields input select toggle checkbox' },
  { title: 'Relationships', href: '/docs/admin/relationships', section: 'Admin Panel', keywords: 'relationship belongs_to many_to_many foreign key association preload select' },
  { title: 'Dashboard & Widgets', href: '/docs/admin/widgets', section: 'Admin Panel', keywords: 'dashboard stats chart widget cards analytics' },

  // Frontend
  { title: 'Web App', href: '/docs/frontend/web-app', section: 'Frontend (Next.js)', keywords: 'nextjs react web app pages routes components vite tanstack router' },
  { title: 'React Query Hooks', href: '/docs/frontend/hooks', section: 'Frontend (Next.js)', keywords: 'react query hooks tanstack fetch data mutation cache useQuery useMutation' },
  { title: 'Shared Package', href: '/docs/frontend/shared-package', section: 'Frontend (Next.js)', keywords: 'shared types schemas constants zod validation packages monorepo apiErrorMessage' },

  // Batteries
  { title: 'File Storage', href: '/docs/batteries/storage', section: 'Batteries', keywords: 's3 minio upload file image storage cloudflare r2 presigned url thumbnail' },
  { title: 'Email System', href: '/docs/batteries/email', section: 'Batteries', keywords: 'email resend mail template send smtp transactional html' },
  { title: 'Background Jobs', href: '/docs/batteries/jobs', section: 'Batteries', keywords: 'jobs queue background worker asynq redis async dlq retry' },
  { title: 'Cron Scheduler', href: '/docs/batteries/cron', section: 'Batteries', keywords: 'cron schedule periodic task timer recurring' },
  { title: 'Redis Caching', href: '/docs/batteries/caching', section: 'Batteries', keywords: 'redis cache middleware ttl performance' },
  { title: 'AI Integration', href: '/docs/batteries/ai', section: 'Batteries', keywords: 'ai claude openai gemini llm chat stream completion sonnet anthropic gateway' },
  { title: 'Idempotency Middleware', href: '/docs/backend/middleware', section: 'Batteries', keywords: 'idempotency idempotency-key safe retry replay 24h cache stripe-style middleware unsafe methods POST PUT PATCH DELETE' },
  { title: 'Activity Log + Hash Chain', href: '/docs/backend/middleware', section: 'Batteries', keywords: 'activity log audit trail tamper evident hash chain sha256 SOC2 mutation tracking compliance integrity verification' },
  { title: 'Feature Flags + A/B Testing', href: '/docs/backend/middleware', section: 'Batteries', keywords: 'feature flags ab testing rollout percentage allowlist blocklist sticky bucketing variants launchdarkly posthog kill switch experiments' },
  { title: 'Webhook Receiver', href: '/docs/backend/middleware', section: 'Batteries', keywords: 'webhook receiver inbound stripe github twilio whatsapp signature verification HMAC verify replay deduplication' },
  { title: 'Realtime WebSocket Hub', href: '/docs/backend/middleware', section: 'Batteries', keywords: 'realtime websocket hub broadcast SendToUser fan-out chat notifications useRealtimeEvent' },
  { title: 'PDF Generation', href: '/docs/batteries/storage', section: 'Batteries', keywords: 'pdf generation invoice receipt lease fpdf go-pdf doc primitives header table totals' },
  { title: 'CSV / Excel Export', href: '/docs/batteries/storage', section: 'Batteries', keywords: 'csv excel xlsx export resource download streaming excelize' },
  { title: 'Cursor-based Pagination', href: '/docs/backend/response-format', section: 'Batteries', keywords: 'cursor pagination offset paginate next_cursor has_more sticky pages stable' },

  // Security & Testing
  { title: 'Security Guide', href: '/docs/security', section: 'Security & Testing', keywords: 'security owasp top 10 2025 idor sqli xss csrf ssrf broken access control authentication authorization injection cryptography secrets headers csp hsts dependabot supply chain audit log hardening checklist sentinel waf rate limiting' },
  { title: 'Performance & Pentest Testing', href: '/docs/testing', section: 'Security & Testing', keywords: 'testing k6 load test smoke average stress spike soak breakpoint pentest penetration test methodology burp ffuf nmap sqlmap nuclei cvss audit report vulnerability scan govulncheck pnpm audit codeql ptes owasp wstg' },

  // Infrastructure
  { title: 'Docker Setup', href: '/docs/infrastructure/docker', section: 'Infrastructure', keywords: 'docker compose container postgresql redis minio' },
  { title: 'Docker Cheat Sheet', href: '/docs/infrastructure/docker-cheatsheet', section: 'Infrastructure', keywords: 'docker commands cheat sheet reference' },
  { title: 'Database & Migrations', href: '/docs/infrastructure/database', section: 'Infrastructure', keywords: 'postgresql database connection pool config' },
  { title: 'Deployment', href: '/docs/infrastructure/deployment', section: 'Infrastructure', keywords: 'deploy production hosting railway fly docker' },

  // Design System
  { title: 'Theme & Colors', href: '/docs/design/theme', section: 'Design System', keywords: 'theme dark light colors palette tailwind design' },

  // Tutorials
  { title: 'Learn Grit Step by Step', href: '/docs/tutorials/learn', section: 'Tutorials', keywords: 'tutorial beginner learn first getting started task manager step by step curriculum' },
  { title: 'Build a Blog', href: '/docs/tutorials/blog', section: 'Tutorials', keywords: 'tutorial blog post article example walkthrough' },
  { title: 'Build a SaaS', href: '/docs/tutorials/saas', section: 'Tutorials', keywords: 'tutorial saas subscription billing example' },
  { title: 'Build an E-Commerce', href: '/docs/tutorials/ecommerce', section: 'Tutorials', keywords: 'tutorial ecommerce store products orders cart' },

  // AI
  { title: 'LLM Skill Guide', href: '/docs/ai-skill', section: 'For AI Assistants', keywords: 'ai llm claude skill guide assistant prompt' },

  // AI Workflows
  { title: 'Using Grit with Claude', href: '/docs/ai-workflows/claude', section: 'AI Workflows', keywords: 'claude code ai spec workflow plan build prompt project description phases' },
  { title: 'Using Grit with Antigravity', href: '/docs/ai-workflows/antigravity', section: 'AI Workflows', keywords: 'antigravity cursor ide ai spec workflow plan build composer inline' },

  // Desktop (Wails)
  { title: 'Desktop Overview', href: '/docs/desktop', section: 'Desktop (Wails)', keywords: 'desktop wails native overview windows macos linux electron' },
  { title: 'Desktop Getting Started', href: '/docs/desktop/getting-started', section: 'Desktop (Wails)', keywords: 'desktop install wails prerequisites setup new-desktop' },
  { title: 'Your First Desktop App', href: '/docs/desktop/first-app', section: 'Desktop (Wails)', keywords: 'desktop first app tutorial wails dev hello hands on' },
  { title: 'Build a POS App', href: '/docs/desktop/pos-app', section: 'Desktop (Wails)', keywords: 'desktop pos point of sale tutorial wails inventory receipts sales transactions' },
  { title: 'Desktop Resource Generation', href: '/docs/desktop/resource-generation', section: 'Desktop (Wails)', keywords: 'desktop generate resource wails crud model service tanstack injection' },
  { title: 'Offline-First Desktop Apps', href: '/docs/desktop/offline', section: 'Desktop (Wails)', keywords: 'offline desktop sync local sqlite outbox conflict resolution wails git push pull versioned mirror' },
  { title: 'Building & Distribution', href: '/docs/desktop/building', section: 'Desktop (Wails)', keywords: 'desktop build distribution wails package executable nsis installer cross compile' },
  { title: '20 Desktop Project Ideas', href: '/docs/desktop/project-ideas', section: 'Desktop (Wails)', keywords: 'desktop project ideas examples inspiration 20 wails grit' },
  { title: 'Desktop LLM Reference', href: '/docs/desktop/llm-reference', section: 'Desktop (Wails)', keywords: 'desktop llm ai reference assistant wails complete cheat sheet' },

  // Plugins
  { title: 'Plugins Overview', href: '/docs/plugins', section: 'Plugins', keywords: 'plugins extensions websockets stripe oauth grit-websockets grit-stripe grit-oauth packages' },
]

const sectionIcons: Record<string, React.ReactNode> = {
  'Getting Started': <Rocket className="h-3.5 w-3.5" />,
  'Core Concepts': <Box className="h-3.5 w-3.5" />,
  'Backend (Go API)': <Server className="h-3.5 w-3.5" />,
  'Admin Panel': <Shield className="h-3.5 w-3.5" />,
  'Frontend (Next.js)': <Layers className="h-3.5 w-3.5" />,
  'Batteries': <Database className="h-3.5 w-3.5" />,
  'Infrastructure': <Settings className="h-3.5 w-3.5" />,
  'Design System': <Palette className="h-3.5 w-3.5" />,
  'Tutorials': <BookOpen className="h-3.5 w-3.5" />,
  'For AI Assistants': <Lightbulb className="h-3.5 w-3.5" />,
  'AI Workflows': <Wand2 className="h-3.5 w-3.5" />,
  'Desktop (Wails)': <Monitor className="h-3.5 w-3.5" />,
  'Plugins': <Plug className="h-3.5 w-3.5" />,
}

export function SearchDialog() {
  const [open, setOpen] = useState(false)
  const router = useRouter()

  useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === 'k' && (e.metaKey || e.ctrlKey)) {
        e.preventDefault()
        setOpen((open) => !open)
      }
    }

    document.addEventListener('keydown', down)
    return () => document.removeEventListener('keydown', down)
  }, [])

  const onSelect = useCallback((href: string) => {
    setOpen(false)
    router.push(href)
  }, [router])

  const sections = Array.from(new Set(searchIndex.map((item) => item.section)))

  return (
    <>
      <button
        onClick={() => setOpen(true)}
        className="inline-flex items-center gap-2 rounded-md border border-border/50 bg-accent/30 px-3 py-1.5 text-sm text-muted-foreground hover:bg-accent/50 hover:text-foreground transition-colors"
      >
        <FileText className="h-3.5 w-3.5" />
        <span className="hidden sm:inline">Search docs...</span>
        <kbd className="pointer-events-none hidden h-5 select-none items-center gap-1 rounded border border-border/50 bg-background/80 px-1.5 font-mono text-[10px] font-medium text-muted-foreground/70 sm:inline-flex">
          <span className="text-xs">⌘</span>K
        </kbd>
      </button>
      <CommandDialog open={open} onOpenChange={setOpen}>
        <CommandInput placeholder="Search documentation..." />
        <CommandList>
          <CommandEmpty>No results found.</CommandEmpty>
          {sections.map((section) => (
            <CommandGroup key={section} heading={section}>
              {searchIndex
                .filter((item) => item.section === section)
                .map((item) => (
                  <CommandItem
                    key={item.href}
                    value={`${item.title} ${item.keywords}`}
                    onSelect={() => onSelect(item.href)}
                    className="cursor-pointer"
                  >
                    <span className="text-muted-foreground/60">
                      {sectionIcons[item.section] || <FileText className="h-3.5 w-3.5" />}
                    </span>
                    <span>{item.title}</span>
                  </CommandItem>
                ))}
            </CommandGroup>
          ))}
        </CommandList>
      </CommandDialog>
    </>
  )
}
