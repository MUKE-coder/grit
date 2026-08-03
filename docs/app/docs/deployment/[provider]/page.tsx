import Link from 'next/link'
import { notFound } from 'next/navigation'
import { AlertCircle, ArrowLeft, ArrowRight, Check, ExternalLink, Minus } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import {
  DEPLOYMENT_PROVIDERS,
  PROVIDER_KIND_LABEL,
  getProvider,
} from '@/config/deployment-providers'

/**
 * One page per hosting provider, rendered from config/deployment-providers.ts.
 *
 * A shared template rather than seven files, because the useful thing about these
 * pages is that they answer the SAME questions in the same order — cost, effort,
 * who it suits, the steps, what breaks. Hand-written pages diverge, and by the
 * fourth one nobody can tell whether a missing section means "not applicable" or
 * "nobody got round to it".
 */

export async function generateMetadata({
  params,
}: {
  params: Promise<{ provider: string }>
}) {
  const { provider: slug } = await params
  const provider = getProvider(slug)
  if (!provider) return { title: 'Not found' }
  return {
    title: `Deploy Grit to ${provider.name}`,
    description: `${provider.tagline} Cost, trade-offs, step-by-step setup and the things that catch people out.`,
  }
}

export default async function ProviderPage({
  params,
}: {
  params: Promise<{ provider: string }>
}) {
  const { provider: slug } = await params
  const provider = getProvider(slug)
  if (!provider) notFound()

  const index = DEPLOYMENT_PROVIDERS.findIndex((p) => p.slug === slug)
  const next = DEPLOYMENT_PROVIDERS[index + 1]
  const external = provider.docsUrl.startsWith('http')

  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <Link
              href="/docs/deployment"
              className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors mb-8"
            >
              <ArrowLeft className="h-3.5 w-3.5" />
              All deployment targets
            </Link>

            <div className="mb-8">
              <span className="tag-mono text-primary/80 mb-3 block">
                {PROVIDER_KIND_LABEL[provider.kind]}
              </span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Deploy to {provider.name}
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">{provider.tagline}</p>
            </div>

            {/* Facts up front, so someone can rule it out in ten seconds rather
                than reading a walkthrough for a platform that cannot do what
                they need. */}
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-8">
              {[
                { label: 'From', value: provider.costFrom },
                { label: 'Ops effort', value: provider.effort },
                { label: 'Managed Postgres', value: provider.managedPostgres },
                { label: 'Persistent disk', value: provider.persistentDisk },
              ].map((f) => (
                <div key={f.label} className="rounded-xl border border-border/50 bg-card/40 p-4">
                  <div className="text-[11px] uppercase tracking-wider text-muted-foreground mb-1.5">
                    {f.label}
                  </div>
                  {typeof f.value === 'boolean' ? (
                    f.value ? (
                      <Check className="h-4 w-4 text-emerald-500" />
                    ) : (
                      <Minus className="h-4 w-4 text-muted-foreground/50" />
                    )
                  ) : (
                    <div className="text-sm font-semibold">{f.value}</div>
                  )}
                </div>
              ))}
            </div>

            <div className="grid gap-3 sm:grid-cols-2 mb-12">
              <div className="rounded-xl border border-emerald-500/25 bg-emerald-500/[0.04] p-5">
                <div className="text-sm font-semibold text-emerald-600 dark:text-emerald-400 mb-1.5">
                  Best for
                </div>
                <p className="text-sm text-muted-foreground leading-relaxed">{provider.bestFor}</p>
              </div>
              {/* Saying who a platform is wrong for is the most useful sentence
                  on the page, and the one every vendor's own docs omit. */}
              <div className="rounded-xl border border-amber-500/25 bg-amber-500/[0.04] p-5">
                <div className="text-sm font-semibold text-amber-600 dark:text-amber-400 mb-1.5">
                  Not for
                </div>
                <p className="text-sm text-muted-foreground leading-relaxed">{provider.notFor}</p>
              </div>
            </div>

            <h2 className="text-2xl font-bold tracking-tight mb-6">Setup</h2>
            <ol className="space-y-8 mb-14">
              {provider.steps.map((step, i) => (
                <li key={step.title} className="relative pl-11">
                  <span className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-full bg-primary/10 text-[13px] font-bold text-primary">
                    {i + 1}
                  </span>
                  <h3 className="font-semibold mb-2 leading-tight">{step.title}</h3>
                  <p className="text-sm text-muted-foreground leading-relaxed mb-3">{step.body}</p>
                  {step.code && <CodeBlock language={step.code.language} code={step.code.code} />}
                </li>
              ))}
            </ol>

            <h2 className="text-2xl font-bold tracking-tight mb-4">What catches people out</h2>
            <div className="space-y-3 mb-14">
              {provider.gotchas.map((g) => (
                <div
                  key={g}
                  className="flex gap-3 rounded-xl border border-border/50 bg-card/40 p-4"
                >
                  <AlertCircle className="h-4 w-4 shrink-0 text-amber-500 mt-0.5" />
                  <p className="text-sm text-muted-foreground leading-relaxed">{g}</p>
                </div>
              ))}
            </div>

            <div className="rounded-xl border border-border/50 bg-card/40 p-5 mb-12">
              <p className="text-sm text-muted-foreground leading-relaxed">
                Platform dashboards and CLI flags change faster than these docs. For anything
                that looks different from what is written here, {provider.name}&apos;s own
                documentation is the authority:{' '}
                <Link
                  href={provider.docsUrl}
                  target={external ? '_blank' : undefined}
                  rel={external ? 'noreferrer' : undefined}
                  className="inline-flex items-center gap-1 text-primary hover:underline"
                >
                  {provider.docsUrl.replace(/^https?:\/\//, '')}
                  {external && <ExternalLink className="h-3 w-3" />}
                </Link>
              </p>
            </div>

            <div className="flex flex-wrap items-center justify-between gap-4 border-t border-border/40 pt-8">
              <Link
                href="/docs/deployment/checklist"
                className="inline-flex items-center gap-1.5 text-sm font-medium text-primary hover:underline"
              >
                Go-live checklist
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
              {next && (
                <Link
                  href={`/docs/deployment/${next.slug}`}
                  className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
                >
                  Next: {next.name}
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              )}
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}

export function generateStaticParams() {
  return DEPLOYMENT_PROVIDERS.map((p) => ({ provider: p.slug }))
}

export const dynamicParams = false
