import Link from 'next/link'
import { ArrowLeft, ArrowRight, ExternalLink, Check, Github } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { StarterPromptCard } from '@/components/starter-prompt-card'

export interface TechKitFeature {
  /** Lucide icon component. */
  icon: React.ComponentType<{ className?: string }>
  title: string
  body: string
  /** Optional flag indicating this is "new in this kit". */
  badge?: string
}

export interface TechKitProps {
  /** Heading shown in the hero. */
  name: string
  /** One-line subtitle under the heading. */
  tagline: string
  /** Two-three sentence pitch under the tagline. */
  pitch: React.ReactNode
  /** Exact CLI command — `grit new my-app --foo` */
  command: string
  /** Mockup / visual element rendered to the right of the hero text. */
  mockup: React.ReactNode
  /** Feature card grid. */
  features: TechKitFeature[]
  /** Path to the matching `/docs/concepts/architecture-modes/<slug>` deep-dive. */
  architectureDeepLink: string
  /**
   * Starter prompt that the user copies into claude.ai to generate the four
   * planning files (project-description, project-phases, design-style-guide,
   * prompt.md) for this kit. When provided, renders a `<StarterPromptCard>`
   * between the features grid and the "Ready to build?" strip.
   */
  starterPrompt?: string
  /** Prev / next nav at the bottom. */
  prev?: { label: string; href: string }
  next?: { label: string; href: string }
}

export function TechKitLayout({
  name,
  tagline,
  pitch,
  command,
  mockup,
  features,
  architectureDeepLink,
  starterPrompt,
  prev,
  next,
}: TechKitProps) {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-5xl mx-auto">

            {/* Hero */}
            <div className="grid grid-cols-1 lg:grid-cols-[1.1fr_1fr] gap-10 items-center mb-16">
              <div>
                <span className="tag-mono text-primary/80 mb-3 block">Tech Kit</span>
                <h1 className="text-4xl md:text-5xl font-bold tracking-tight mb-3 leading-tight">
                  {name}
                </h1>
                <p className="text-lg text-foreground/75 mb-4 leading-snug">{tagline}</p>
                <div className="prose-grit mb-6 text-base">{pitch}</div>

                {/* CLI command */}
                <CodeBlock terminal filename="Quick install" code={command} className="!mb-5" />

                <div className="flex flex-wrap gap-2.5">
                  <Button asChild className="rounded-full bg-primary hover:bg-primary/90 text-primary-foreground">
                    <Link href="/docs/getting-started/quick-start">
                      Get started <ArrowRight className="ml-1.5 h-3.5 w-3.5" />
                    </Link>
                  </Button>
                  <Button asChild variant="outline" className="border-border/60 rounded-full">
                    <Link href={architectureDeepLink}>
                      Architecture deep-dive <ExternalLink className="ml-1.5 h-3.5 w-3.5" />
                    </Link>
                  </Button>
                  <Button asChild variant="outline" className="border-border/60 rounded-full">
                    <Link href="https://github.com/MUKE-coder/grit" target="_blank" rel="noopener noreferrer">
                      <Github className="mr-1.5 h-3.5 w-3.5" /> Source
                    </Link>
                  </Button>
                </div>
              </div>

              {/* Visual */}
              <div className="relative">{mockup}</div>
            </div>

            {/* What's included */}
            <div className="mb-16">
              <h2 className="text-2xl font-semibold tracking-tight mb-2">What&apos;s included</h2>
              <p className="text-muted-foreground mb-6">
                Every batteries-included primitive the framework ships, wired into this kit by default.
              </p>
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
                {features.map((f) => {
                  const Icon = f.icon
                  return (
                    <div
                      key={f.title}
                      className="rounded-xl border border-border/50 bg-card/50 p-5 hover:border-primary/40 transition-colors"
                    >
                      <div className="flex items-center gap-2.5 mb-2">
                        <div className="h-8 w-8 rounded-md bg-primary/10 text-primary flex items-center justify-center">
                          <Icon className="h-4 w-4" />
                        </div>
                        <h3 className="text-sm font-semibold text-foreground">{f.title}</h3>
                        {f.badge && (
                          <span className="text-[9px] font-mono text-primary bg-primary/10 px-1.5 py-0.5 rounded-full ml-auto">
                            {f.badge}
                          </span>
                        )}
                      </div>
                      <p className="text-sm text-muted-foreground leading-relaxed">{f.body}</p>
                    </div>
                  )
                })}
              </div>
            </div>

            {/* Starter prompt for claude.ai (renders only when provided) */}
            {starterPrompt && (
              <StarterPromptCard
                kitLabel={name}
                command={command}
                prompt={starterPrompt}
              />
            )}

            {/* Ready-to-build CTA strip */}
            <div className="rounded-2xl border border-primary/20 bg-primary/[0.04] p-6 mb-12 text-center">
              <h3 className="text-xl font-semibold mb-2">Ready to build?</h3>
              <p className="text-muted-foreground mb-5">
                One command, your app online in 30 seconds. Read the quick-start guide for the
                rest of the wiring.
              </p>
              <div className="flex flex-wrap justify-center gap-2.5">
                <Button asChild className="rounded-full bg-primary hover:bg-primary/90 text-primary-foreground">
                  <Link href="/docs/getting-started/quick-start">
                    Read docs <ArrowRight className="ml-1.5 h-3.5 w-3.5" />
                  </Link>
                </Button>
                <Button asChild variant="outline" className="border-border/60 rounded-full">
                  <Link href="/docs/tech-kits">Explore all kits</Link>
                </Button>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href={prev?.href ?? '/docs/tech-kits'} className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  {prev?.label ?? 'All Tech Kits'}
                </Link>
              </Button>
              {next ? (
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href={next.href} className="gap-1.5">
                    {next.label}
                    <ArrowRight className="h-3.5 w-3.5" />
                  </Link>
                </Button>
              ) : null}
            </div>

          </div>
        </div>
      </main>
    </div>
  )
}

/** Compact "what's included" bullet for the hero pitch — reusable. */
export function IncludedRow({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex items-start gap-2 text-sm text-foreground/80 mb-1.5">
      <Check className="h-4 w-4 text-primary mt-0.5 shrink-0" strokeWidth={2.5} />
      <span>{children}</span>
    </div>
  )
}
