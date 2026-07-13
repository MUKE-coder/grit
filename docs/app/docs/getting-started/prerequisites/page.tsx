import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { PageHelp } from '@/components/page-help'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/getting-started/prerequisites')

const TILES: { title: string; blurb: string; href: string; tag: string; art: string }[] = [
  {
    title: 'Go for Grit Developers',
    blurb: 'The 20% of Go you need for Grit — structs, methods, interfaces, error handling.',
    href: '/docs/prerequisites/golang',
    tag: 'Backend',
    art: 'func main() {\n  r := gin.Default()\n  r.Run(":8080")\n}',
  },
  {
    title: 'Go Playground',
    blurb: 'Run Go in the browser — experiment with the snippets from the primer, no install.',
    href: '/playground',
    tag: 'Interactive',
    art: 'package main\n\nimport "fmt"\n\nfmt.Println("Hello, Grit")',
  },
  {
    title: 'Next.js & React',
    blurb: 'App Router, components, and hooks — the frontend concepts Grit builds on.',
    href: '/docs/prerequisites/nextjs',
    tag: 'Frontend',
    art: 'export default function Page() {\n  return <h1>Hello</h1>\n}',
  },
  {
    title: 'Docker',
    blurb: 'Just enough Docker to run Postgres, Redis and MinIO locally with one command.',
    href: '/docs/prerequisites/docker',
    tag: 'Infra',
    art: '$ docker compose up -d\n  ✓ postgres  redis  minio',
  },
]

export default function PrerequisitesPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-4xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Getting Started</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Prerequisites</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit assumes you know Go and can read React/TypeScript &mdash; it teaches{' '}
                <em>Grit&apos;s</em> conventions on top, not the languages themselves. If any of
                these are new, start with the matching primer. Otherwise, jump straight to{' '}
                <Link href="/docs/getting-started/create-a-project" className="text-primary hover:underline">Create a project</Link>.
              </p>
            </div>

            <div className="grid gap-5 sm:grid-cols-2">
              {TILES.map((t) => (
                <Link
                  key={t.href}
                  href={t.href}
                  className="group rounded-xl border border-border/50 bg-card/40 overflow-hidden hover:border-primary/30 transition-colors"
                >
                  {/* Mockup art */}
                  <div className="relative border-b border-border/40 bg-[#0d1117] px-4 py-3.5 overflow-hidden">
                    <div className="flex items-center gap-1.5 mb-2">
                      <span className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <span className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <span className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                      <span className="ml-2 text-[10px] font-mono uppercase tracking-wider text-primary/50">{t.tag}</span>
                    </div>
                    <pre className="text-[11px] leading-5 font-mono text-slate-300/80 whitespace-pre-wrap">{t.art}</pre>
                  </div>
                  {/* Body */}
                  <div className="p-5">
                    <h2 className="text-base font-semibold mb-1.5 flex items-center gap-1.5 group-hover:text-primary transition-colors">
                      {t.title}
                      <ArrowRight className="h-3.5 w-3.5 opacity-0 -translate-x-1 group-hover:opacity-100 group-hover:translate-x-0 transition-all" />
                    </h2>
                    <p className="text-sm text-muted-foreground/70 leading-relaxed">{t.blurb}</p>
                  </div>
                </Link>
              ))}
            </div>

            <div className="max-w-3xl">
              <PageHelp
                faqs={[
                  {
                    q: 'Do I really need to know Go?',
                    a: 'Enough to read and write handlers, services and GORM models. The primer covers exactly that slice — you do not need advanced Go.',
                  },
                  {
                    q: 'Can I skip Docker?',
                    a: 'Yes — use managed Neon (Postgres), Upstash (Redis) and Cloudflare R2 instead, or the Desktop kit which runs on local SQLite.',
                  },
                ]}
              />
            </div>

            <div className="flex items-center justify-between border-t border-border pt-8 mt-12 max-w-3xl">
              <Button variant="ghost" asChild>
                <Link href="/docs/getting-started/coming-from" className="gap-2">
                  <ArrowLeft className="h-4 w-4" />
                  Coming from Laravel / Django
                </Link>
              </Button>
              <Button variant="ghost" asChild>
                <Link href="/docs/getting-started/create-a-project" className="gap-2">
                  Create a project
                  <ArrowRight className="h-4 w-4" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
