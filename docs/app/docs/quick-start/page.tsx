import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CopyButton } from '@/components/copy-button'

export default function QuickStartPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Getting Started</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Quick Start
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Get up and running with OmniStack in under 5 minutes. This guide will walk you
                through installing your first component.
              </p>
            </div>

            {/* Prerequisites */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Prerequisites
              </h2>
              <ul className="space-y-2.5">
                {[
                  'Node.js 18+ installed',
                  'A Next.js 14+ project (App Router recommended)',
                  'pnpm, npm, or yarn package manager',
                  'shadcn/ui initialized in your project',
                ].map((item) => (
                  <li key={item} className="flex items-start gap-2.5 text-[14px] text-muted-foreground">
                    <span className="text-primary mt-1">&#10003;</span>
                    {item}
                  </li>
                ))}
              </ul>
            </div>

            {/* Step 1 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                  1
                </div>
                <h2 className="text-xl font-semibold tracking-tight">
                  Create a Next.js project
                </h2>
              </div>
              <p className="text-muted-foreground leading-relaxed mb-4 pl-10">
                If you don&apos;t have a project yet, create one with:
              </p>
              <div className="ml-10 rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                </div>
                <div className="p-5 font-mono text-sm">
                  <div className="flex items-center justify-between gap-4">
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">npx create-next-app@latest my-app --typescript --tailwind --app</span>
                    </div>
                    <CopyButton text="npx create-next-app@latest my-app --typescript --tailwind --app" />
                  </div>
                </div>
              </div>
            </div>

            {/* Step 2 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                  2
                </div>
                <h2 className="text-xl font-semibold tracking-tight">
                  Initialize shadcn/ui
                </h2>
              </div>
              <p className="text-muted-foreground leading-relaxed mb-4 pl-10">
                OmniStack components are built on top of shadcn/ui. Initialize it in your project:
              </p>
              <div className="ml-10 rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                </div>
                <div className="p-5 font-mono text-sm">
                  <div className="flex items-center justify-between gap-4">
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">pnpm dlx shadcn@latest init</span>
                    </div>
                    <CopyButton text="pnpm dlx shadcn@latest init" />
                  </div>
                </div>
              </div>
            </div>

            {/* Step 3 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                  3
                </div>
                <h2 className="text-xl font-semibold tracking-tight">
                  Install your first component
                </h2>
              </div>
              <p className="text-muted-foreground leading-relaxed mb-4 pl-10">
                Pick a component from the library and install it. For example, add the
                authentication system:
              </p>
              <div className="ml-10 rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-green-sm">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                </div>
                <div className="p-5 font-mono text-sm">
                  <div className="flex items-center justify-between gap-4">
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-primary font-medium">pnpm dlx shadcn@latest add</span>{' '}
                      <span className="text-foreground/80">https://better-auth-ui.desishub.com/r/auth-components.json</span>
                    </div>
                    <CopyButton text="pnpm dlx shadcn@latest add https://better-auth-ui.desishub.com/r/auth-components.json" />
                  </div>
                </div>
              </div>
              <p className="text-sm text-muted-foreground/60 mt-3 pl-10">
                The CLI will install the component files, dependencies, and any required shadcn/ui
                primitives automatically.
              </p>
            </div>

            {/* Step 4 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                  4
                </div>
                <h2 className="text-xl font-semibold tracking-tight">
                  Configure and run
                </h2>
              </div>
              <p className="text-muted-foreground leading-relaxed mb-4 pl-10">
                Each component has its own configuration requirements (environment variables,
                database setup, etc.). Check the component&apos;s detail page for specific instructions.
              </p>
              <div className="ml-10 rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                </div>
                <div className="p-5 font-mono text-sm space-y-2">
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">npx prisma migrate dev --name init</span>
                  </div>
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">pnpm dev</span>
                  </div>
                </div>
              </div>
            </div>

            {/* What's next */}
            <div className="mb-10">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                What&apos;s Next?
              </h2>
              <div className="space-y-3">
                {[
                  { label: 'Installation Guide', href: '/docs/installation', desc: 'Detailed installation instructions and troubleshooting.' },
                  { label: 'Browse Components', href: '/components', desc: 'Explore the full component library.' },
                  { label: 'JB Better Auth UI', href: '/components/jb-better-auth-ui', desc: 'Start with a complete authentication system.' },
                ].map((item) => (
                  <Link
                    key={item.href}
                    href={item.href}
                    className="flex items-center justify-between p-4 rounded-lg border border-border/30 bg-card/30 hover:bg-card/60 hover:border-primary/20 transition-all group"
                  >
                    <div>
                      <p className="text-sm font-medium group-hover:text-primary transition-colors">{item.label}</p>
                      <p className="text-xs text-muted-foreground/60 mt-0.5">{item.desc}</p>
                    </div>
                    <ArrowRight className="h-4 w-4 text-muted-foreground/40 group-hover:text-primary transition-colors" />
                  </Link>
                ))}
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Introduction
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/installation" className="gap-1.5">
                  Installation
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
