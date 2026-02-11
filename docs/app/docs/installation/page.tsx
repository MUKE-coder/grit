import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CopyButton } from '@/components/copy-button'
import { components } from '@/lib/components-data'

export default function InstallationPage() {
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
                Installation
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                All OmniStack components are installed using the shadcn CLI. This page covers
                the installation process, requirements, and available components.
              </p>
            </div>

            {/* Requirements */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Requirements
              </h2>
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Dependency</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Version</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Node.js</td>
                      <td className="px-4 py-2.5 font-mono text-xs">18+</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Next.js</td>
                      <td className="px-4 py-2.5 font-mono text-xs">14+</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">React</td>
                      <td className="px-4 py-2.5 font-mono text-xs">18+</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">TypeScript</td>
                      <td className="px-4 py-2.5 font-mono text-xs">5+</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5">shadcn/ui</td>
                      <td className="px-4 py-2.5 font-mono text-xs">initialized</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            {/* Install shadcn */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Initialize shadcn/ui
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                If you haven&apos;t already, initialize shadcn/ui in your project. This sets up the
                component configuration and base styles:
              </p>
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
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

            {/* Database setup */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Database Setup
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Most components require a PostgreSQL database with Prisma ORM. Set up your
                database connection in your <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">.env</code> file:
              </p>
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">.env</span>
                </div>
                <div className="p-5 font-mono text-sm">
                  <div className="flex items-center justify-between gap-4">
                    <code className="text-foreground/70">DATABASE_URL=&quot;postgresql://user:password@localhost:5432/mydb&quot;</code>
                    <CopyButton text='DATABASE_URL="postgresql://user:password@localhost:5432/mydb"' />
                  </div>
                </div>
              </div>
              <p className="text-sm text-muted-foreground/60 mt-3">
                After installing a component, run <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">npx prisma migrate dev</code> to
                create the required database tables.
              </p>
            </div>

            {/* All components */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-6">
                Available Components
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-6">
                Install any component using the shadcn CLI. Click on a component to view its
                full documentation and usage guide.
              </p>
              <div className="space-y-4">
                {components.map((component) => (
                  <div key={component.id} className="rounded-xl border border-border/40 bg-card/50 overflow-hidden">
                    <div className="p-4">
                      <div className="flex items-start justify-between gap-3 mb-2">
                        <Link href={`/components/${component.id}`} className="text-sm font-semibold hover:text-primary transition-colors">
                          {component.name}
                        </Link>
                        <span className="tag-mono text-[10px] text-primary/60 px-2 py-0.5 rounded-md bg-primary/5 border border-primary/10 shrink-0">
                          {component.category}
                        </span>
                      </div>
                      <p className="text-xs text-muted-foreground/60 mb-3 line-clamp-1">
                        {component.description}
                      </p>
                      <div className="rounded-lg border border-border/30 bg-accent/20 overflow-hidden">
                        <div className="px-4 py-2.5 font-mono text-xs">
                          <div className="flex items-center justify-between gap-4">
                            <code className="text-foreground/70 truncate">{component.installCommand}</code>
                            <CopyButton text={component.installCommand} />
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Troubleshooting */}
            <div className="mb-10">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Troubleshooting
              </h2>
              <div className="space-y-4">
                {[
                  {
                    q: 'shadcn CLI not found',
                    a: 'Make sure you\'re using pnpm dlx, npx, or bunx to run the shadcn CLI. Example: pnpm dlx shadcn@latest add <url>',
                  },
                  {
                    q: 'Prisma migration errors',
                    a: 'Ensure your DATABASE_URL is correct and the PostgreSQL server is running. Run npx prisma generate after migrations.',
                  },
                  {
                    q: 'Component conflicts',
                    a: 'If a component installs files that conflict with existing ones, the CLI will prompt you. You can choose to overwrite or skip.',
                  },
                ].map((item) => (
                  <div key={item.q} className="p-4 rounded-lg border border-border/30 bg-card/30">
                    <h3 className="text-sm font-semibold mb-1.5">{item.q}</h3>
                    <p className="text-xs text-muted-foreground/70 leading-relaxed">{item.a}</p>
                  </div>
                ))}
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/quick-start" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Quick Start
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/components" className="gap-1.5">
                  Components
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
