import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { CommunityCTA } from '@/components/community-cta'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/frontend/ui-components')

// Counted from `grit ui list`, which is the registry itself. These numbers
// move every time components are added, so they are worth re-checking rather
// than trusting: the previous set had been wrong for long enough that every
// example on this page named a component that no longer existed.
const CATEGORIES = [
  {
    name: 'marketing',
    count: 93,
    examples: 'hero sections, pricing, FAQs, feature sections, CTAs, headers, stats, whole landing pages',
  },
  {
    name: 'ecommerce',
    count: 34,
    examples: 'product grids, product detail, checkout, store banners, category tiles, store navigation',
  },
  {
    name: 'application-ui',
    count: 14,
    examples: 'sign-in and register cards, sidebar layouts, tables, forms, inputs, buttons, upload',
  },
]

export default function UIComponentsPage() {
  return (
    <DocsShell>
            <div className="mb-10">
              <span className="tag-mono text-accent/80 mb-3 block">Frontend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">UI Components</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit UI is a registry of 141 ready-made React components. Install one
                with a command and it lands in your repo as an ordinary{' '}
                <code>.tsx</code> file you own.
              </p>
            </div>

            <div className="prose-grit">
              <div className="rounded-lg border border-accent/20 bg-accent/5 p-4 mb-8">
                <p className="text-sm text-muted-foreground leading-relaxed">
                  Browse and preview every component at{' '}
                  <Link
                    href="https://ui.gritframework.dev"
                    target="_blank"
                    rel="noreferrer"
                    className="text-accent hover:underline"
                  >
                    ui.gritframework.dev
                  </Link>
                  .
                </p>
              </div>

              {/* ============================================================ */}
              <h2>Installing a component</h2>
              <p>
                Inside a Grit project, the CLI writes the component into the right app
                for your architecture:
              </p>
              <CodeBlock
                language="bash"
                code={`grit ui list                           # everything in the registry
grit ui list --category ecommerce      # narrow it down
grit ui add marketing-pricing-sections-three-tiers-with-toggle
grit ui add application-ui-tables-sortable-table-with-selection \\
  application-ui-authentication-sign-in-card-with-oauth`}
              />
              <p>
                A name is its category, then its group, then the component:{' '}
                <code>ecommerce-product-grids-grid-with-ratings</code>. Running{' '}
                <code>grit ui list</code> is the only reliable way to spell one.
              </p>
              <p>
                Files land in <code>components/grit-ui/</code> inside{' '}
                <code>apps/web</code>, falling back to <code>apps/admin</code> or a flat{' '}
                <code>frontend/</code> depending on how the project was scaffolded. An
                existing file is never overwritten without <code>--force</code> — once
                you have edited a component, it is your code.
              </p>

              {/* ============================================================ */}
              <h2>Outside a Grit project</h2>
              <p>
                Every component is a{' '}
                <Link href="https://ui.shadcn.com" target="_blank" rel="noreferrer">
                  shadcn
                </Link>{' '}
                registry item, so the registry works with any React project that has
                Tailwind — no Grit required:
              </p>
              <CodeBlock
                language="bash"
                code={`npx shadcn@latest add https://ui.gritframework.dev/r/marketing-hero-sections-split-with-code.json`}
              />
              <p>
                That single command writes the component, merges Grit&apos;s design
                tokens into your <code>globals.css</code>, and adds the colour scale to
                your <code>tailwind.config.ts</code>. Both are needed: the tokens supply
                the palette, and the scale is what makes classes like{' '}
                <code>text-text-muted</code> resolve.
              </p>

              {/* ============================================================ */}
              <h2>What is in the registry</h2>
              <div className="not-prose my-6 overflow-hidden rounded-xl border border-border/40">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/40 bg-accent/20">
                      <th className="px-4 py-2.5 text-left font-medium">Category</th>
                      <th className="px-4 py-2.5 text-left font-medium w-16">Count</th>
                      <th className="px-4 py-2.5 text-left font-medium">Examples</th>
                    </tr>
                  </thead>
                  <tbody>
                    {CATEGORIES.map((c) => (
                      <tr key={c.name} className="border-b border-border/20 last:border-0">
                        <td className="px-4 py-2.5 font-mono text-xs">{c.name}</td>
                        <td className="px-4 py-2.5 font-mono text-xs">{c.count}</td>
                        <td className="px-4 py-2.5 text-muted-foreground">{c.examples}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              {/* ============================================================ */}
              <h2>The registry endpoints</h2>
              <p>
                The registry is plain JSON over HTTP, so anything that speaks the shadcn
                format can consume it — your own tooling included.
              </p>
              <CodeBlock
                language="bash"
                code={`https://ui.gritframework.dev/r/registry.json   # the index
https://ui.gritframework.dev/r/<name>.json     # one item, source inlined`}
              />
              <p>
                Each item carries its source in <code>files[0].content</code> rather than
                a path to fetch separately, which is what lets a single request produce a
                working component.
              </p>

              {/* ============================================================ */}
              <h2>Components are yours</h2>
              <p>
                Nothing here is a dependency. There is no <code>grit-ui</code> package in
                your <code>package.json</code>, no version to keep up with, and no
                upgrade that can change a component out from under you. The trade is the
                usual one for generated code:{' '}
                <strong>improvements to a component do not reach copies you already
                installed</strong>. Re-run <code>grit ui add --force</code> if you want
                the newer version and have not edited yours.
              </p>
              <p>
                They use Grit&apos;s design tokens, so in a scaffolded project they match
                the admin panel and auth pages out of the box.
              </p>
            </div>

            <CommunityCTA className="mt-10" />

            <div className="flex flex-wrap gap-3 mt-12 pt-6 border-t border-border/30">
              <Button variant="outline" asChild className="border-border/60 bg-transparent hover:bg-surface-2">
                <Link href="/docs/frontend/web-app">
                  <ArrowLeft className="mr-1.5 h-4 w-4" />
                  Frontend
                </Link>
              </Button>
              <Button asChild className="glow-purple-sm ml-auto">
                <Link href="/docs/admin/overview">
                  Admin Panel
                  <ArrowRight className="ml-1.5 h-4 w-4" />
                </Link>
              </Button>
            </div>
    </DocsShell>
  )
}
