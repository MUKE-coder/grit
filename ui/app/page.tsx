import Link from 'next/link'
import { ArrowRight, Package } from 'lucide-react'
import { getComponents, categories, baseUrl } from '@/lib/registry'
import { Gallery } from '@/components/gallery'
import { CopyButton } from '@/components/copy-button'

export default function HomePage() {
  const components = getComponents()
  const base = baseUrl()
  const cats = categories()
  const exampleInstall = `npx shadcn@latest add ${base}/r/hero-split-01.json`

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="sticky top-0 z-40 border-b border-border bg-background/80 backdrop-blur">
        <div className="mx-auto flex h-14 max-w-7xl items-center justify-between px-6">
          <Link href="/" className="flex items-center gap-2">
            <span className="inline-flex h-6 w-6 items-center justify-center rounded-md bg-accent text-accent-foreground">
              <Package size={14} />
            </span>
            <span className="text-sm font-semibold text-foreground">Grit UI</span>
          </Link>
          <nav className="flex items-center gap-5 text-xs">
            <a
              href={`${base}/r/registry.json`}
              className="font-mono text-text-muted transition-colors hover:text-foreground"
            >
              registry.json
            </a>
            <Link
              href="https://gritframework.dev"
              className="text-text-secondary transition-colors hover:text-foreground"
            >
              Grit Framework
            </Link>
            <Link
              href="https://github.com/MUKE-coder/grit"
              target="_blank"
              rel="noreferrer"
              className="text-text-secondary transition-colors hover:text-foreground"
            >
              GitHub
            </Link>
          </nav>
        </div>
      </header>

      {/* Hero */}
      <section className="relative overflow-hidden border-b border-border">
        <div aria-hidden className="pointer-events-none absolute inset-0 grit-grid mask-fade-b" />
        <div
          aria-hidden
          className="pointer-events-none absolute inset-0"
          style={{
            background:
              'radial-gradient(ellipse 55% 45% at 50% -10%, rgba(108,92,231,0.16), transparent 60%)',
          }}
        />
        <div className="relative mx-auto max-w-7xl px-6 py-20 md:py-28">
          <div className="max-w-3xl">
            <span className="font-mono text-xs uppercase tracking-wider text-accent">
              Grit UI
            </span>
            <h1 className="mt-4 text-4xl font-bold leading-[1.08] tracking-tight text-foreground md:text-6xl">
              {components.length} React components,
              <br />
              installed with one command
            </h1>
            <p className="mt-6 max-w-2xl text-lg leading-relaxed text-text-secondary">
              Marketing sections, SaaS dashboards, ecommerce, auth flows and app
              layout. Every one is a shadcn registry item, so it drops into any
              React project — Grit or not — as ordinary source you own.
            </p>

            <div className="mt-8 flex max-w-2xl items-center justify-between gap-3 rounded-xl border border-border bg-bg-secondary px-4 py-3">
              <code className="truncate font-mono text-xs text-text-secondary md:text-sm">
                {exampleInstall}
              </code>
              <CopyButton value={exampleInstall} />
            </div>

            <div className="mt-6 flex flex-wrap items-center gap-x-6 gap-y-2 font-mono text-xs text-text-muted">
              <span>MIT licensed</span>
              <span>·</span>
              <span>No signup</span>
              <span>·</span>
              <span>Tailwind + lucide-react</span>
              <span>·</span>
              <Link href="/install" className="inline-flex items-center gap-1 text-accent hover:underline">
                Setup guide
                <ArrowRight size={11} />
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* Gallery */}
      <main className="mx-auto max-w-7xl px-6 py-14">
        <Gallery
          items={components.map((c) => ({
            name: c.name,
            title: c.title,
            description: c.description,
            category: c.category,
          }))}
          categories={cats}
          baseUrl={base}
        />
      </main>

      <footer className="border-t border-border">
        <div className="mx-auto flex max-w-7xl flex-col items-center justify-between gap-3 px-6 py-8 text-xs text-text-muted sm:flex-row">
          <p>Part of the Grit Framework. MIT licensed.</p>
          <div className="flex items-center gap-5">
            <Link href="/install" className="transition-colors hover:text-foreground">
              Install
            </Link>
            <a href={`${base}/r/registry.json`} className="font-mono transition-colors hover:text-foreground">
              registry.json
            </a>
            <Link href="https://gritframework.dev/docs" className="transition-colors hover:text-foreground">
              Docs
            </Link>
          </div>
        </div>
      </footer>
    </div>
  )
}
