import Link from 'next/link'
import { notFound } from 'next/navigation'
import { ArrowLeft, Package } from 'lucide-react'
import { getComponent, getComponents, readSource, baseUrl } from '@/lib/registry'
import { CopyButton } from '@/components/copy-button'

export async function generateMetadata({
  params,
}: {
  params: Promise<{ name: string }>
}) {
  const { name } = await params
  const c = getComponent(name)
  if (!c) return { title: 'Not found' }
  return { title: c.title, description: c.description }
}

export default async function ComponentPage({
  params,
}: {
  params: Promise<{ name: string }>
}) {
  const { name } = await params
  const component = getComponent(name)
  if (!component) notFound()

  const base = baseUrl()
  const source = readSource(name)
  const install = `npx shadcn@latest add ${base}/r/${name}.json`
  const gritInstall = `grit ui add ${name}`

  return (
    <div className="min-h-screen bg-background">
      <header className="sticky top-0 z-40 border-b border-border bg-background/80 backdrop-blur">
        <div className="mx-auto flex h-14 max-w-7xl items-center justify-between px-6">
          <Link href="/" className="flex items-center gap-2 text-sm text-text-secondary transition-colors hover:text-foreground">
            <ArrowLeft size={14} />
            All components
          </Link>
          <Link href="/" className="flex items-center gap-2">
            <span className="inline-flex h-6 w-6 items-center justify-center rounded-md bg-accent text-accent-foreground">
              <Package size={14} />
            </span>
            <span className="text-sm font-semibold text-foreground">Grit UI</span>
          </Link>
        </div>
      </header>

      <main className="mx-auto max-w-6xl px-6 py-12">
        <div className="mb-8 flex flex-wrap items-start justify-between gap-4">
          <div>
            <div className="mb-2 flex items-center gap-2.5">
              <span className="rounded-full border border-border bg-bg-elevated px-2.5 py-0.5 font-mono text-[10px] uppercase tracking-wider text-text-muted">
                {component.category}
              </span>
              <code className="font-mono text-xs text-text-muted">{component.name}</code>
            </div>
            <h1 className="text-3xl font-bold tracking-tight text-foreground">
              {component.title}
            </h1>
            <p className="mt-2 max-w-2xl text-sm leading-relaxed text-text-secondary">
              {component.description}
            </p>
          </div>
        </div>

        {/* Install commands */}
        <div className="mb-10 grid gap-3 md:grid-cols-2">
          <InstallBox label="Any React project" command={install} />
          <InstallBox label="Inside a Grit project" command={gritInstall} />
        </div>

        {/* Live preview */}
        <section className="mb-10">
          <h2 className="mb-3 text-sm font-semibold text-foreground">Preview</h2>
          <div className="overflow-hidden rounded-2xl border border-border bg-bg-primary">
            <iframe
              src={`/preview/${name}`}
              title={`${component.title} preview`}
              className="h-[680px] w-full border-0"
            />
          </div>
          <p className="mt-2 font-mono text-[11px] text-text-muted">
            Rendered from the same source the registry serves.
          </p>
        </section>

        {/* Dependencies */}
        {component.dependencies?.length > 0 && (
          <section className="mb-10">
            <h2 className="mb-3 text-sm font-semibold text-foreground">Dependencies</h2>
            <div className="flex flex-wrap gap-2">
              {component.dependencies.map((d) => (
                <code
                  key={d}
                  className="rounded-lg border border-border bg-bg-secondary px-2.5 py-1 font-mono text-xs text-text-secondary"
                >
                  {d}
                </code>
              ))}
            </div>
            <p className="mt-2 text-xs text-text-muted">
              Installed automatically by <code className="font-mono">shadcn add</code>.
            </p>
          </section>
        )}

        {/* Source */}
        <section>
          <div className="mb-3 flex items-center justify-between">
            <h2 className="text-sm font-semibold text-foreground">Source</h2>
            <CopyButton value={source} label="Copy source" />
          </div>
          <div className="overflow-hidden rounded-2xl border border-border bg-bg-secondary">
            <div className="flex items-center gap-2 border-b border-border bg-bg-elevated px-4 py-2">
              <code className="font-mono text-[11px] text-text-muted">
                components/grit-ui/{name}.tsx
              </code>
            </div>
            <pre className="overflow-x-auto p-5 font-mono text-xs leading-relaxed text-text-secondary">
              <code>{source}</code>
            </pre>
          </div>
        </section>
      </main>
    </div>
  )
}

function InstallBox({ label, command }: { label: string; command: string }) {
  return (
    <div className="rounded-xl border border-border bg-bg-secondary p-4">
      <p className="mb-2 font-mono text-[10px] uppercase tracking-wider text-text-muted">
        {label}
      </p>
      <div className="flex items-center justify-between gap-3">
        <code className="truncate font-mono text-xs text-text-secondary">{command}</code>
        <CopyButton value={command} />
      </div>
    </div>
  )
}

export function generateStaticParams() {
  return getComponents().map((c) => ({ name: c.name }))
}

export const dynamicParams = false
