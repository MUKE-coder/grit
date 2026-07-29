import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { CommunityCTA } from '@/components/community-cta'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/frontend/ui-components')

const CATEGORIES = [
  { name: 'marketing', count: 20, examples: 'heroes, pricing, testimonials, FAQ, CTA banners' },
  { name: 'saas', count: 30, examples: 'billing, usage meters, API keys, audit logs, team rows' },
  { name: 'ecommerce', count: 20, examples: 'product cards, cart, checkout steps, variants' },
  { name: 'layout', count: 20, examples: 'navbars, modals, tabs, empty states, toasts' },
  { name: 'auth', count: 10, examples: 'login, signup, OTP, 2FA setup, OAuth buttons' },
]

export default function UIComponentsPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Frontend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">UI Components</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit UI is a registry of 100 ready-made React components. Install one
                with a command and it lands in your repo as an ordinary{' '}
                <code>.tsx</code> file you own.
              </p>
            </div>

            <div className="prose-grit">
              <div className="rounded-lg border border-primary/20 bg-primary/5 p-4 mb-8">
                <p className="text-sm text-muted-foreground leading-relaxed">
                  Browse and preview every component at{' '}
                  <Link
                    href="https://ui.gritframework.dev"
                    target="_blank"
                    rel="noreferrer"
                    className="text-primary hover:underline"
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
                code={`grit ui list                      # everything in the registry
grit ui list --category saas      # narrow it down
grit ui add billing-card-01       # install one
grit ui add login-card-01 otp-input-01   # or several`}
              />
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
                code={`npx shadcn@latest add https://ui.gritframework.dev/r/hero-split-01.json`}
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
                        <td className="px-4 py-2.5 font-mono text-xs capitalize">{c.name}</td>
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
                code={`https://ui.gritframework.dev/r/registry.json   # index of all 100
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
              <Button variant="outline" asChild className="border-border/60 bg-transparent hover:bg-accent/50">
                <Link href="/docs/frontend">
                  <ArrowLeft className="mr-1.5 h-4 w-4" />
                  Frontend
                </Link>
              </Button>
              <Button asChild className="glow-purple-sm ml-auto">
                <Link href="/docs/admin">
                  Admin Panel
                  <ArrowRight className="ml-1.5 h-4 w-4" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
