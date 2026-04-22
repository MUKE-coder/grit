import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/styles')

export default function StyleVariantsPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Core Concepts</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Style Variants
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit ships with 4 admin panel style variants so every scaffolded project
                doesn&apos;t look identical. Choose a style when you create your project and it
                applies to both the auth screens and the dashboard.
              </p>
            </div>

            <div className="prose-grit">
              {/* Usage */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Choosing a Style
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Pass the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--style</code> flag
                  when creating a new project:
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm">
                    <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit new my-app --style modern</span></div>
                  </div>
                </div>

                <p className="text-sm text-muted-foreground/60 mt-3">
                  If you omit the flag, the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">default</code> style
                  is used.
                </p>
              </div>

              {/* Available Styles */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Available Styles
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-6">
                  Each style applies a different visual treatment to the auth pages (login, sign-up,
                  forgot password) and the admin dashboard. All styles share the same functionality,
                  API calls, form validation, and component library &mdash; only the layout and
                  Tailwind classes differ.
                </p>

                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Style</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Auth Pages</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Dashboard</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">default</td>
                        <td className="px-4 py-2.5">Split-screen with branding panel</td>
                        <td className="px-4 py-2.5">Stats grid + resource cards + quick links</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">modern</td>
                        <td className="px-4 py-2.5">Centered card on gradient background</td>
                        <td className="px-4 py-2.5">Bento grid with varied card sizes</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">minimal</td>
                        <td className="px-4 py-2.5">Clean full-width form (Vercel/Linear style)</td>
                        <td className="px-4 py-2.5">Compact data-dense with table layout</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">glass</td>
                        <td className="px-4 py-2.5">Frosted-glass card with backdrop blur</td>
                        <td className="px-4 py-2.5">Hero banner with gradient mesh</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">centered</td>
                        <td className="px-4 py-2.5">Single centered card on radial gradient (Linear school)</td>
                        <td className="px-4 py-2.5">Uses default dashboard</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* default */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  default
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The original Grit layout. Auth pages use a split-screen design with a
                  branded gradient panel on the left and the form on the right. The dashboard
                  shows a welcome banner, stats widgets, resource cards, and quick links in a
                  clean grid.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm">
                    <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit new my-app</span></div>
                    <div className="text-muted-foreground/50 mt-1"># or explicitly:</div>
                    <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit new my-app --style default</span></div>
                  </div>
                </div>

                <div className="mt-6 space-y-3">
                  <h3 className="text-xl font-semibold tracking-tight">Auth layout</h3>
                  <ul className="list-disc pl-6 space-y-2 text-muted-foreground">
                    <li>Left panel: accent gradient background with Grit logo, tagline, and branding text</li>
                    <li>Right panel: centered form with React Hook Form + Zod validation</li>
                    <li>Responsive: left panel hidden on mobile, form fills the screen</li>
                  </ul>
                  <h3 className="text-xl font-semibold tracking-tight mt-6">Dashboard layout</h3>
                  <ul className="list-disc pl-6 space-y-2 text-muted-foreground">
                    <li>Welcome banner with gradient background and user greeting</li>
                    <li>Stats widgets from resource definitions (or fallback stats)</li>
                    <li>Two-column grid: resource cards on the left, quick links on the right</li>
                  </ul>
                </div>
              </div>

              {/* modern */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  modern
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  A centered card layout inspired by contemporary SaaS products. Auth pages
                  place a single card on a subtle gradient background with decorative blurred
                  orbs. The dashboard uses a bento grid with varying card sizes.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm">
                    <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit new my-app --style modern</span></div>
                  </div>
                </div>

                <div className="mt-6 space-y-3">
                  <h3 className="text-xl font-semibold tracking-tight">Auth layout</h3>
                  <ul className="list-disc pl-6 space-y-2 text-muted-foreground">
                    <li>Full-screen gradient background (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">from-accent/5 via-background to-background</code>)</li>
                    <li>Decorative blurred accent orbs for visual depth</li>
                    <li>Single centered card with <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">rounded-2xl</code>, shadow, and border</li>
                    <li>Grit logo centered at the top of the card</li>
                  </ul>
                  <h3 className="text-xl font-semibold tracking-tight mt-6">Dashboard layout</h3>
                  <ul className="list-disc pl-6 space-y-2 text-muted-foreground">
                    <li>Full-width welcome banner with accent gradient</li>
                    <li>Bento-style grid: 2-column stats + 1-column system status card</li>
                    <li>Resource cards in a large panel, quick links in a compact sidebar</li>
                  </ul>
                </div>
              </div>

              {/* minimal */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  minimal
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Inspired by Vercel and Linear &mdash; clean, sparse, and information-first.
                  No decorative backgrounds, no cards on auth pages. Just the form on a flat
                  background with generous whitespace. The dashboard trades visual flair for
                  a compact, data-dense table view.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm">
                    <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit new my-app --style minimal</span></div>
                  </div>
                </div>

                <div className="mt-6 space-y-3">
                  <h3 className="text-xl font-semibold tracking-tight">Auth layout</h3>
                  <ul className="list-disc pl-6 space-y-2 text-muted-foreground">
                    <li>Flat <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">bg-background</code> &mdash; no cards, no gradients, no panels</li>
                    <li>Narrow form area (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">max-w-sm</code>) centered on the page</li>
                    <li>Subtle text logo at the top, tight vertical spacing</li>
                    <li>Tighter label spacing and smaller heading sizes</li>
                  </ul>
                  <h3 className="text-xl font-semibold tracking-tight mt-6">Dashboard layout</h3>
                  <ul className="list-disc pl-6 space-y-2 text-muted-foreground">
                    <li>No welcome banner &mdash; just a heading and subtitle</li>
                    <li>Stats in a compact 4-column row</li>
                    <li>Resources displayed as a table (icon, name, endpoint, &ldquo;Manage&rdquo; link)</li>
                    <li>Quick links as horizontal pill buttons instead of stacked cards</li>
                  </ul>
                </div>
              </div>

              {/* glass */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  glass
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The most visually striking variant. Auth pages use a layered gradient
                  background with a frosted-glass card (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">backdrop-blur-2xl</code>).
                  The dashboard features a large hero banner with decorative SVG mesh patterns
                  and floating blurred orbs.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm">
                    <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit new my-app --style glass</span></div>
                  </div>
                </div>

                <div className="mt-6 space-y-3">
                  <h3 className="text-xl font-semibold tracking-tight">Auth layout</h3>
                  <ul className="list-disc pl-6 space-y-2 text-muted-foreground">
                    <li>Multi-layer gradient background with radial gradients for depth</li>
                    <li>Frosted card: <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">bg-bg-secondary/60 backdrop-blur-2xl border-white/10</code></li>
                    <li>Inputs use <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">bg-white/5</code> and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">border-white/10</code> for consistency</li>
                    <li>Softer label colors (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">text-foreground/80</code>)</li>
                  </ul>
                  <h3 className="text-xl font-semibold tracking-tight mt-6">Dashboard layout</h3>
                  <ul className="list-disc pl-6 space-y-2 text-muted-foreground">
                    <li>Large hero section with layered gradients, SVG mesh overlay, and floating orbs</li>
                    <li>Inline glass stats cards embedded in the hero</li>
                    <li>Resource cards use glass styling (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">backdrop-blur-xl border-white/5</code>)</li>
                    <li>Quick links in a 4-column glass panel at the bottom</li>
                  </ul>
                </div>
              </div>

              {/* centered */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  centered
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Linear / Vercel school auth pages. A single centered card (~420px wide) floats on
                  a subtle radial gradient background. No split-screen branding panel. Minimal chrome,
                  generous whitespace. The dashboard uses the default layout — only the auth pages
                  change. Added in v3.8.0.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm">
                    <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit new my-app --style centered</span></div>
                  </div>
                </div>

                <div className="mt-6 space-y-3">
                  <h3 className="text-xl font-semibold tracking-tight">Auth layout</h3>
                  <ul className="list-disc pl-6 space-y-2 text-muted-foreground">
                    <li>Full-screen container: <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">flex items-center justify-center</code></li>
                    <li>Card: <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">max-w-[420px] rounded-2xl p-10 shadow-sm</code></li>
                    <li>Grit logo mark (purple squircle with white &ldquo;G&rdquo;) centered at top</li>
                    <li>44px input height, 44px button height, 20px field gap</li>
                    <li>Subtle radial glow at top of background (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">rgba(108, 92, 231, 0.10)</code>)</li>
                    <li>OAuth buttons in a 2-column grid with &ldquo;or continue with&rdquo; divider</li>
                  </ul>
                  <h3 className="text-xl font-semibold tracking-tight mt-6">Dashboard layout</h3>
                  <p className="text-muted-foreground">
                    Uses the same dashboard as <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">default</code>.
                    Only the auth pages change. Good fit when you want Linear-feeling onboarding
                    but a standard dashboard experience.
                  </p>
                </div>
              </div>

              {/* Persistence */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  How the Style is Persisted
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The selected style is saved in your project&apos;s <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit.config.ts</code> file:
                </p>

                <CodeBlock language="typescript" filename="grit.config.ts" code={`// Grit Framework Configuration
export default {
  name: "my-app",
  style: "modern",
  api: {
    port: 8080,
    prefix: "/api",
  },
  // ...
};`} />

                <p className="text-muted-foreground leading-relaxed mt-4">
                  When you run <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit upgrade</code>,
                  the CLI reads this value and regenerates the auth pages and dashboard using the
                  correct style variant. Your choice is preserved across upgrades.
                </p>
              </div>

              {/* What changes */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  What Changes Between Styles
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Style variants <strong>only</strong> affect the visual layout. Everything else stays the same:
                </p>

                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-4">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Changes</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Stays the same</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5">JSX layout structure</td>
                        <td className="px-4 py-2.5">Form validation (React Hook Form + Zod)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5">Tailwind CSS classes</td>
                        <td className="px-4 py-2.5">API endpoints and data fetching</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5">Background / gradient effects</td>
                        <td className="px-4 py-2.5">Authentication hooks (useLogin, useRegister)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5">Card / panel arrangements</td>
                        <td className="px-4 py-2.5">Shared schemas and types</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5">Dashboard widget layout</td>
                        <td className="px-4 py-2.5">All other admin components (tables, forms, sidebar)</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <p className="text-muted-foreground leading-relaxed">
                  The sidebar, navbar, data table, form builder, resource pages, and system pages
                  are identical across all styles. Only the 4 files below change:
                </p>

                <ul className="list-disc pl-6 mt-3 space-y-2 text-muted-foreground">
                  <li><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">app/(auth)/login/page.tsx</code></li>
                  <li><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">app/(auth)/sign-up/page.tsx</code></li>
                  <li><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">app/(auth)/forgot-password/page.tsx</code></li>
                  <li><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">app/(dashboard)/dashboard/page.tsx</code></li>
                </ul>
              </div>

              {/* Customizing */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Customizing After Scaffolding
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The generated pages are plain React components with Tailwind CSS. After
                  scaffolding, you own the code and can modify anything:
                </p>
                <ul className="list-disc pl-6 space-y-2 text-muted-foreground">
                  <li>Change colors, gradients, or spacing in the Tailwind classes</li>
                  <li>Add your own logo or branding images</li>
                  <li>Mix elements from different styles (e.g., use the glass auth with the minimal dashboard)</li>
                  <li>Add new fields to auth forms by updating the Zod schemas in <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">packages/shared</code></li>
                </ul>
                <p className="text-muted-foreground leading-relaxed mt-4">
                  Keep in mind that running <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit upgrade</code> will
                  overwrite these 4 files with the latest framework version of your chosen style.
                  If you&apos;ve heavily customized them, back them up before upgrading.
                </p>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/concepts/naming-conventions" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Naming Conventions
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/backend/models" className="gap-1.5">
                  Models & Database
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
